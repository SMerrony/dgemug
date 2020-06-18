// +build virtual !physical

// scMemory.go - Memory-related System Call Emulation

// Copyright ©2020 Steve Merrony

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package aosvs

import (
	"log"
	"strings"

	"github.com/SMerrony/dgemug/dg"
	"github.com/SMerrony/dgemug/logging"
	"github.com/SMerrony/dgemug/memory"
)

func scGshpt(p syscallParmsT) bool {
	p.cpu.SetAc(0, memory.GetFirstSharedPage()&0x0003_ffff)
	p.cpu.SetAc(1, dg.DwordT(memory.GetNumSharedPages()))
	return true
}

func scMem(p syscallParmsT) bool {
	highestUnsharedInUse := memory.GetLastUnsharedPage() //& 0x0003_ffff // assumed in current ring
	lowShared := memory.GetFirstSharedPage()             //& 0x0003_ffff
	unusedUnshared := int32(lowShared) - int32(highestUnsharedInUse) - 4
	if unusedUnshared < 0 {
		unusedUnshared = 0
	}
	inUsePages := highestUnsharedInUse & (0x0fff_ffff >> 10)
	p.cpu.SetAc(0, dg.DwordT(unusedUnshared))
	p.cpu.SetAc(1, inUsePages)
	p.cpu.SetAc(2, highestUnsharedInUse<<10|dg.DwordT(p.ringMask))
	//p.cpu.SetAc(2, inUsePages+dg.DwordT(p.ringMask)-1)
	logging.DebugPrint(logging.ScLog, "\tMax Unshared Available: %d., Unshared in Use: %d., Highest in Use: %#o (%#x)\n",
		unusedUnshared, inUsePages, highestUnsharedInUse<<10, highestUnsharedInUse<<10)
	logging.DebugPrint(logging.ScLog, "\tN.B. Lowest shared page is %#o (%#x)\n", lowShared<<10, lowShared<<10)
	return true
}

func scMemi(p syscallParmsT) bool {
	numPages := int32(int16(p.cpu.GetAc(0)))
	logging.DebugPrint(logging.ScLog, "\tRequested page count: %d, (%#x)\n", numPages, p.cpu.GetAc(0))
	//var lastPage int
	switch {
	case numPages > 0: // add pages
		logging.DebugPrint(logging.ScLog, "\tAdding %d. page(s)\n", numPages)
		for numPages > 0 {
			//lastPage = memory.AddUnsharedPage()
			memory.AddUnsharedPage()
			numPages--
		}
		//p.cpu.SetAc(1, (dg.DwordT(lastPage<<10)|dg.DwordT(p.ringMask))-1)
		highestUnsharedInUse := memory.GetLastUnsharedPage()
		p.cpu.SetAc(1, highestUnsharedInUse<<10|dg.DwordT(p.ringMask))
		logging.DebugPrint(logging.ScLog, "\tHighest in use is now  %#o (%#x)\n", memory.GetLastUnsharedPage()<<10, memory.GetLastUnsharedPage()<<10)
	case numPages < 0: // remove pages
		log.Panicln("ERROR: Unmapping via ?MEMI not yet supported")
	}
	return true
}

// AOS/VS treats this as a memory operation, not a file one...
func scSopen(p syscallParmsT) bool {
	bpFilename := p.cpu.GetAc(0)
	filename := strings.ToUpper(readString(bpFilename, p.ringMask))
	if p.cpu.GetAc(1) != 0xffff_ffff {
		log.Panicln("ERROR: ?SOPEN of specific channel not yet implemented")
	}
	var sopenReq = agSharedOpenReqT{p.PID, filename, p.cpu.GetAc(2) == 0}
	var areq = AgentReqT{agentSharedOpen, sopenReq, nil}
	p.agentChan <- areq
	areq = <-p.agentChan
	if areq.result.(agSharedOpenRespT).ac0 != 0 {
		p.cpu.SetAc(0, areq.result.(agSharedOpenRespT).ac0)
		return false
	}
	p.cpu.SetAc(1, areq.result.(agSharedOpenRespT).channelNo)
	return true
}

func scSpage(p syscallParmsT) bool {
	fileChan := int(p.cpu.GetAc(1))
	pktAddr := dg.PhysAddrT(p.cpu.GetAc(2))
	diskBytesCount := int(memory.ReadWord(pktAddr+psti)) * 512
	diskStartBlock := int64(memory.ReadDWord(pktAddr+prnh)) * 512
	memStartAddr := dg.PhysAddrT(memory.ReadDWord(pktAddr + pcad))
	asrRec := agSharedReadReqT{fileChan, diskBytesCount, diskStartBlock}
	areq := AgentReqT{agentSharedRead, asrRec, nil}
	p.agentChan <- areq
	areq = <-p.agentChan
	if areq.result.(agSharedReadRespT).ac0 != 0 {
		p.cpu.SetAc(0, areq.result.(agSharedReadRespT).ac0)
		return false
	}
	//writeBytes(memStartAddr, p.ringMask, areq.result.(agSharedReadRespT).data)
	words := memory.WordsFromBytes(areq.result.(agSharedReadRespT).data)
	memory.MapSlice(memStartAddr, words, true)
	return true
}

func scSshpt(p syscallParmsT) bool { // TODO removing pages
	firstPageNo := p.cpu.GetAc(0)
	newSize := p.cpu.GetAc(1)
	//initialLastPage := firstPageNo + dg.DwordT(memory.GetNumSharedPages())
	if firstPageNo < memory.GetFirstSharedPage()&0x0003_ffff {
		// try to add pages at start
		var pg int
		for pg = int(firstPageNo); pg < int(memory.GetFirstSharedPage()&0x0003_ffff); pg++ {
			if memory.IsPageMapped(pg) {
				p.cpu.SetAc(0, ermem)
				return false
			}
			memory.MapPage(pg, true)
		}
	}
	// now at end
	for memory.GetNumSharedPages() < int(newSize) {
		memory.MapPage(int(memory.GetLastSharedPage())+1, true)
	}
	return true
}
