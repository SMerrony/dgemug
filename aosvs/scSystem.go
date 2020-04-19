// scSystem.go - 'System'-related System Call Emulation

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
	"time"

	"github.com/SMerrony/dgemug/dg"
	"github.com/SMerrony/dgemug/logging"
	"github.com/SMerrony/dgemug/memory"
)

func scErmsg(p syscallParmsT) bool {
	// fake an error message for now
	msg := "VSemuG Dummy Error Message"
	bp := p.cpu.GetAc(2)
	memory.WriteStringBA(msg, bp)
	p.cpu.SetAc(0, dg.DwordT(len(msg)))
	return true
}

func scExec(p syscallParmsT) bool {
	pktAddr := dg.PhysAddrT(p.cpu.GetAc(2))
	execFunc := memory.ReadWord(pktAddr)
	switch execFunc {
	case xfxts:
		memory.WriteWord(pktAddr+xfp1, 9) // PID 9
		bp := memory.ReadDWord(pktAddr + xfp2)
		if bp != 0 {
			memory.WriteStringBA("CON10", bp) // Claim to be @CON10
		}
	default:
		logging.DebugPrint(logging.ScLog, "WARNING: ?EXEC system call not yet implemented - fn code was: %#o\n", execFunc)
	}
	return true
}

func scGday(p syscallParmsT) bool {
	now := time.Now()
	p.cpu.SetAc(0, dg.DwordT(now.Day()))
	p.cpu.SetAc(1, dg.DwordT(now.Month()))
	p.cpu.SetAc(2, dg.DwordT(now.Year()-1900))
	return true
}

func scGhrz(p syscallParmsT) bool {
	p.cpu.SetAc(0, 2) // 2 => 100Hz
	return true
}

func scGtmes(p syscallParmsT) bool {
	pktAddr := dg.PhysAddrT(p.cpu.GetAc(2))
	var gtMesReq = agGtMesReqT{p.PID, memory.ReadWord(pktAddr + greq), memory.ReadWord(pktAddr + gnum), memory.ReadDWord(pktAddr + gsw)}
	var areq = AgentReqT{agentGetMessage, gtMesReq, nil}
	p.agentChan <- areq
	areq = <-p.agentChan
	p.cpu.SetAc(0, areq.result.(agGtMesRespT).ac0)
	p.cpu.SetAc(1, areq.result.(agGtMesRespT).ac1)
	gresBA := memory.ReadDWord(pktAddr+gres) | dg.DwordT((p.ringMask)<<1)
	if gresBA != 0xffff_ffff && len(areq.result.(agGtMesRespT).result) > 0 {
		memory.WriteStringBA(areq.result.(agGtMesRespT).result, gresBA|dg.DwordT((p.ringMask)<<1))
	}
	return true
}
func scGtmes16(p syscallParmsT) bool {
	pktAddr := dg.PhysAddrT(p.cpu.GetAc(2)) | (p.cpu.GetPC() & 0x7000_0000)
	var gtMesReq = agGtMesReqT{p.PID, memory.ReadWord(pktAddr + greq16), memory.ReadWord(pktAddr + gnum16), dg.DwordT(memory.ReadWord(pktAddr + gsw16))}
	var areq = AgentReqT{agentGetMessage, gtMesReq, nil}
	p.agentChan <- areq
	areq = <-p.agentChan
	p.cpu.SetAc(0, areq.result.(agGtMesRespT).ac0)
	p.cpu.SetAc(1, areq.result.(agGtMesRespT).ac1)
	gresBA := dg.DwordT(memory.ReadWord(pktAddr + gres16))
	if gresBA != 0xffff && len(areq.result.(agGtMesRespT).result) > 0 {
		memory.WriteStringBA(areq.result.(agGtMesRespT).result, gresBA|dg.DwordT((p.ringMask)<<1))
	}
	return true
}

func scGtod(p syscallParmsT) bool {
	now := time.Now()
	p.cpu.SetAc(0, dg.DwordT(now.Second()))
	p.cpu.SetAc(1, dg.DwordT(now.Minute()))
	p.cpu.SetAc(2, dg.DwordT(now.Hour()))
	return true
}

func scInfo(p syscallParmsT) bool {
	pktAddr := dg.PhysAddrT(p.cpu.GetAc(2))
	memory.WriteWord(pktAddr+sirn, 0x0746) // system rev - faked to 7.70
	if memory.ReadDWord(pktAddr+siln) != 0 {
		memory.WriteStringBA("MASTERLDU", memory.ReadDWord(pktAddr+siln)) // fake master LDU name
	}
	if memory.ReadDWord(pktAddr+siid) != 0 {
		memory.WriteStringBA("VSEMUG", memory.ReadDWord(pktAddr+siid)) // fake System ID
	}
	if memory.ReadDWord(pktAddr+sios) != 0 {
		memory.WriteStringBA(":VSEMUG", memory.ReadDWord(pktAddr+siid)) // fake OS pathname
	}
	memory.WriteWord(pktAddr+ssin, savs) // claim to be AOS/VS!
	return true
}

func scXpstat(p syscallParmsT) bool {

	return true
}
