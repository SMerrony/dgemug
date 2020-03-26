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

	"github.com/SMerrony/dgemug/dg"
	"github.com/SMerrony/dgemug/memory"
	"github.com/SMerrony/dgemug/mvcpu"
)

func scClose(cpu *mvcpu.CPUT, agentChan chan AgentReqT) bool {
	pktAddr := dg.PhysAddrT(cpu.GetAc(2))
	channel := memory.ReadWord(pktAddr + ich)
	var creq = agCloseReqT{channel}
	var areq = AgentReqT{agentFileClose, creq, nil}
	agentChan <- areq
	areq = <-agentChan
	return true
}

func scOpen(cpu *mvcpu.CPUT, agentChan chan AgentReqT) bool {
	pktAddr := dg.PhysAddrT(cpu.GetAc(2)) | (cpu.GetPC() & 0x7000_0000)
	// pkt := readPacket(pktAddr, iosz)
	// _ = pkt
	options := memory.ReadWord(pktAddr + isti)
	fileType := memory.ReadWord(pktAddr + isto)
	bpPathname := memory.ReadDWord(pktAddr + ifnp)
	path := readString(bpPathname)
	log.Printf("DEBUG: ?OPEN Pathname: %s, Type: %#x, Options: %#x\n", path, fileType, options)
	var areq AgentReqT
	var openReq = agOpenReqT{path, options}
	areq.action = agentFileOpen
	areq.reqParms = openReq
	agentChan <- areq
	areq = <-agentChan
	memory.WriteWord(pktAddr+ich, areq.result.(agOpenRespT).channelNo)
	return true
}

func scWrite(cpu *mvcpu.CPUT, agentChan chan AgentReqT) bool {
	pktAddr := dg.PhysAddrT(cpu.GetAc(2))
	channel := memory.ReadWord(pktAddr + ich)
	bytes := readBytes(memory.ReadDWord(pktAddr + ibad))
	// log.Println("DEBUG: ?WRITE")
	var writeReq = agWriteReqT{channel, bytes}
	var areq = AgentReqT{agentFileWrite, writeReq, nil}
	agentChan <- areq
	areq = <-agentChan
	memory.WriteWord(pktAddr+irlr, areq.result.(agWriteRespT).bytesTxfrd)
	return true
}
