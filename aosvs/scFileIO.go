// scFileIO.go - File I/O System Call Emulation

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

func scClose(p syscallParmsT) bool {
	pktAddr := dg.PhysAddrT(p.cpu.GetAc(2))
	channel := memory.ReadWord(pktAddr + ich)
	var creq = agCloseReqT{channel}
	var areq = AgentReqT{agentFileClose, creq, nil}
	p.agentChan <- areq
	areq = <-p.agentChan
	return true
}

func scGchr(p syscallParmsT) bool {
	var areq AgentReqT
	var gchrReq agGchrReqT
	if memory.TestDwbit(p.cpu.GetAc(1), 0) {
		// AC0 should contain a channel #
		if memory.TestDwbit(p.cpu.GetAc(1), 1) {
			// get default chars
			gchrReq = agGchrReqT{p.PID, true, true, dg.WordT(p.cpu.GetAc(0)), ""}
		} else {
			// get current chars
			gchrReq = agGchrReqT{p.PID, false, true, dg.WordT(p.cpu.GetAc(0)), ""}
		}
	} else {
		// AC0 should contain BP to device name
		bpPathname := p.cpu.GetAc(0)
		path := strings.ToUpper(readString(bpPathname, p.ringMask))
		if memory.TestDwbit(p.cpu.GetAc(1), 1) {
			// get default chars
			gchrReq = agGchrReqT{p.PID, true, false, 0, path}
		} else {
			// get current chars
			gchrReq = agGchrReqT{p.PID, false, false, 0, path}
		}
	}
	areq.action = agentGetChars
	areq.reqParms = gchrReq
	p.agentChan <- areq
	areq = <-p.agentChan
	wrAddr := dg.PhysAddrT(p.cpu.GetAc(2)) | p.ringMask
	memory.WriteWord(dg.PhysAddrT(wrAddr), areq.result.(agGchrRespT).words[0])
	memory.WriteWord(dg.PhysAddrT(wrAddr+1), areq.result.(agGchrRespT).words[1])
	memory.WriteWord(dg.PhysAddrT(wrAddr+2), areq.result.(agGchrRespT).words[2])
	return true
}

func scGechr(p syscallParmsT) bool {

	return true
}

func scOpen(p syscallParmsT) bool {
	pktAddr := dg.PhysAddrT(p.cpu.GetAc(2)) | (p.ringMask)
	options := memory.ReadWord(pktAddr + isti)
	fileType := memory.ReadWord(pktAddr + isto)
	// blockSize := memory.ReadWord(pktAddr+imrs)
	// recLen := memory.ReadWord(pktAddr+ircl)
	bpPathname := memory.ReadDWord(pktAddr + ifnp)
	path := strings.ToUpper(readString(bpPathname, p.ringMask))
	logging.DebugPrint(logging.ScLog, "?OPEN Pathname: %s, Type: %#x, Options: %#x\n", path, fileType, options)
	var areq AgentReqT
	var openReq = agOpenReqT{p.PID, path, options}
	areq.action = agentFileOpen
	areq.reqParms = openReq
	p.agentChan <- areq
	areq = <-p.agentChan
	if areq.result.(agOpenRespT).ac0 != 0 {
		p.cpu.SetAc(0, areq.result.(agOpenRespT).ac0)
		return false
	}
	memory.WriteWord(pktAddr+ich, areq.result.(agOpenRespT).channelNo)
	logging.DebugPrint(logging.ScLog, "----- Returned channel # %d\n", areq.result.(agOpenRespT).channelNo)
	return true
}

func scOpen16(p syscallParmsT) bool {
	pktAddr := dg.PhysAddrT(p.cpu.GetAc(2)) | (p.ringMask)
	options := memory.ReadWord(pktAddr + isti16)
	fileType := memory.ReadWord(pktAddr + isto16)
	// blockSize := memory.ReadWord(pktAddr+imrs16)
	// recLen := memory.ReadWord(pktAddr+ircl)
	bpPathname := dg.DwordT(memory.ReadWord(pktAddr + ifnp16))
	path := strings.ToUpper(readString(bpPathname, p.ringMask))
	logging.DebugPrint(logging.ScLog, "?OPEN Pathname: %s, Type: %#x, Options: %#x\n", path, fileType, options)
	var areq AgentReqT
	var openReq = agOpenReqT{p.PID, path, options}
	areq.action = agentFileOpen
	areq.reqParms = openReq
	p.agentChan <- areq
	areq = <-p.agentChan
	if areq.result.(agOpenRespT).ac0 != 0 {
		p.cpu.SetAc(0, areq.result.(agOpenRespT).ac0)
		return false
	}
	memory.WriteWord(pktAddr+ich16, areq.result.(agOpenRespT).channelNo)
	return true
}

func scRead(p syscallParmsT) bool {
	pktAddr := dg.PhysAddrT(p.cpu.GetAc(2)) | (p.ringMask)
	channel := int(memory.ReadWord(pktAddr + ich))
	specs := memory.ReadWord(pktAddr + isti)
	length := int(memory.ReadWord(pktAddr + ircl))
	dest := memory.ReadDWord(pktAddr + ibad)
	readLine := (memory.ReadWord(pktAddr+isti) & ibin) == 0
	if specs&ipkl != 0 {
		log.Panic("ERROR: ?READ (32-bit) extended packet not yet implemented")
	}
	logging.DebugPrint(logging.ScLog, "?READ (32-bit) Channel: %#x, Specs: %#x, Bytes: %#x, Dest: %#x, Line Mode: %v\n", channel, specs, length, dest, readLine)
	var readReq = agReadReqT{channel, specs, length, readLine}
	var areq = AgentReqT{agentFileRead, readReq, nil}
	p.agentChan <- areq
	areq = <-p.agentChan
	resp := areq.result.(agReadRespT)
	memory.WriteWord(pktAddr+irlr, dg.WordT(len(resp.data)))
	writeBytes(dest, p.ringMask, resp.data)
	if resp.ac0 != 0 {
		p.cpu.SetAc(0, dg.DwordT(resp.ac0))
	}
	return true
}

func scRead16(p syscallParmsT) bool {
	pktAddr := dg.PhysAddrT(p.cpu.GetAc(2)) | p.ringMask
	channel := int(memory.ReadWord(pktAddr + ich16))
	specs := memory.ReadWord(pktAddr + isti16)
	length := int(memory.ReadWord(pktAddr + ircl16))
	dest := dg.DwordT(memory.ReadWord(pktAddr + ibad16))
	readLine := (memory.ReadWord(pktAddr+isti16) & ibin) == 0
	if specs&ipkl != 0 {
		log.Panic("ERROR: ?READ (16-bit) extended packet not yet implemented")
	}
	logging.DebugPrint(logging.ScLog, "?READ (16-bit) Channel: %#x, Specs: %#x, Bytes: %#x, Dest: %#x, Line Mode: %v\n", channel, specs, length, dest, readLine)
	var readReq = agReadReqT{channel, specs, length, readLine}
	var areq = AgentReqT{agentFileRead, readReq, nil}
	p.agentChan <- areq
	areq = <-p.agentChan
	resp := areq.result.(agReadRespT)
	memory.WriteWord(pktAddr+irlr16, dg.WordT(len(resp.data)))
	writeBytes(dest, p.ringMask, resp.data)
	return true
}

func scWrite(p syscallParmsT) bool {
	pktAddr := dg.PhysAddrT(p.cpu.GetAc(2))
	channel := int(memory.ReadWord(pktAddr + ich))
	bytes := readBytes(memory.ReadDWord(pktAddr+ibad), p.cpu.GetPC())
	// log.Println("DEBUG: ?WRITE")
	var writeReq = agWriteReqT{channel, bytes}
	var areq = AgentReqT{agentFileWrite, writeReq, nil}
	p.agentChan <- areq
	areq = <-p.agentChan
	memory.WriteWord(pktAddr+irlr, areq.result.(agWriteRespT).bytesTxfrd)
	return true
}

func scWrite16(p syscallParmsT) bool {
	pktAddr := dg.PhysAddrT(p.cpu.GetAc(2)) | p.ringMask
	channel := int(memory.ReadWord(pktAddr + ich16))
	bytes := readBytes(dg.DwordT(memory.ReadWord(pktAddr+ibad16)), p.ringMask)
	// log.Println("DEBUG: ?WRITE")
	var writeReq = agWriteReqT{channel, bytes}
	var areq = AgentReqT{agentFileWrite, writeReq, nil}
	p.agentChan <- areq
	areq = <-p.agentChan
	memory.WriteWord(pktAddr+irlr, areq.result.(agWriteRespT).bytesTxfrd)
	return true
}
