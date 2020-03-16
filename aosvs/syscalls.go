// syscalls.go - map of AOS/VS system calls

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
	"bytes"
	"log"

	"github.com/SMerrony/dgemug/memory"

	"github.com/SMerrony/dgemug/dg"
	"github.com/SMerrony/dgemug/mvcpu"
)

type syscallDescT struct {
	name        string                                 // AOS/VS System Call Name
	alias       string                                 // AOS/VS 4-char Name Alias
	syscallType int                                    // groupings as per Table 2-1 in Sys Call Dict
	fn          func(*mvcpu.CPUT, chan AgentReqT) bool // implementation
}

const (
	scMemory = iota
	scProcess
	scFileManage
	scFileIO
	scDebugging
	scWindowing
	scMultitasking
	scIPC
	scConnection
	scMultiproc
	scClass
	scSystem
	scUserDev
	scBisync
	sc16Bit
)

var syscalls = map[dg.WordT]syscallDescT{
	0:    {"?CREATE", "?CREA", scFileManage, nil},
	1:    {"?DELETE", "?DELE", scFileManage, nil},
	0300: {"?OPEN", "?OPEN", scFileIO, scOpen},
	0301: {"?CLOSE", "?CLOS", scFileIO, scClose},
	0302: {"?READ", "?READ", scFileIO, nil},
	0303: {"?WRITE", "?WRIT", scFileIO, scWrite},
	0310: {"?RETURN", "?RETU", scFileIO, nil},
	0542: {"?IFPU", "?IFPU", scMultitasking, scIfpu},
}

func syscall(callID dg.WordT, agent chan AgentReqT, cpu *mvcpu.CPUT) (ok bool) {
	call, defined := syscalls[callID]
	if !defined {
		log.Fatalf("ERROR: System call No. %#o not yet defined at PC=%#x", callID, cpu.GetPC())
	}
	if call.fn == nil {
		log.Fatalf("ERROR: System call No. %#o not yet implemented at PC=%#x", callID, cpu.GetPC())
	}
	return call.fn(cpu, agent)
}

// readPacket just loads a chunk of memory into a slice of words
// TODO maybe this should be in ram_virtual.go as 'ReadWords' for efficiency?
func readPacket(addr dg.PhysAddrT, pktLen int) (pkt []dg.WordT) {
	pkt = make([]dg.WordT, pktLen, pktLen)
	for w := range pkt {
		pkt[w] = memory.ReadWord(addr + dg.PhysAddrT(w))
	}
	return pkt
}

// readBytes reads characters up to the first NUL from the given doubleword byte address
func readBytes(bpAddr dg.DwordT) []byte {
	buff := bytes.NewBufferString("")
	lobyte := (bpAddr & 0x0001) == 1
	wdAddr := dg.PhysAddrT(bpAddr >> 1)
	c := memory.ReadByte(wdAddr, lobyte)
	for c != 0 {
		buff.WriteByte(byte(c))
		if lobyte {
			wdAddr++
		}
		lobyte = !lobyte
		c = memory.ReadByte(wdAddr, lobyte)
	}
	return buff.Bytes()
}

// readString reads characters up to the first NUL from the given doubleword byte address
func readString(bpAddr dg.DwordT) string {
	buff := bytes.NewBufferString("")
	lobyte := (bpAddr & 0x0001) == 1
	wdAddr := dg.PhysAddrT(bpAddr >> 1)
	c := memory.ReadByte(wdAddr, lobyte)
	for c != 0 {
		buff.WriteByte(byte(c))
		if lobyte {
			wdAddr++
		}
		lobyte = !lobyte
		c = memory.ReadByte(wdAddr, lobyte)
	}
	return buff.String()
}
