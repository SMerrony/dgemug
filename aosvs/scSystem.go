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
	"log"
	"time"

	"github.com/SMerrony/dgemug/dg"
	"github.com/SMerrony/dgemug/memory"
	"github.com/SMerrony/dgemug/mvcpu"
)

func scExec(cpu *mvcpu.CPUT, agentChan chan AgentReqT) bool {
	log.Println("WARNING: ?EXEC system call not yet implemented")
	return true
}

func scGday(cpu *mvcpu.CPUT, agentChan chan AgentReqT) bool {
	now := time.Now()
	cpu.SetAc(0, dg.DwordT(now.Day()))
	cpu.SetAc(1, dg.DwordT(now.Month()))
	cpu.SetAc(2, dg.DwordT(now.Year()-1900))
	return true
}

func scGhrz(cpu *mvcpu.CPUT, agentChan chan AgentReqT) bool {
	cpu.SetAc(0, 2) // 2 => 100Hz
	return true
}

func scGtmes(cpu *mvcpu.CPUT, agentChan chan AgentReqT) bool {
	pktAddr := dg.PhysAddrT(cpu.GetAc(2))
	var gtMesReq = agGtMesReqT{memory.ReadWord(pktAddr + greq), memory.ReadWord(pktAddr + gnum), memory.ReadDWord(pktAddr + gsw)}
	var areq = AgentReqT{agentGetMessage, gtMesReq, nil}
	agentChan <- areq
	areq = <-agentChan
	cpu.SetAc(0, areq.result.(agGtMesRespT).ac0)
	cpu.SetAc(1, areq.result.(agGtMesRespT).ac1)
	gresBA := memory.ReadDWord(pktAddr+gres) | dg.DwordT((cpu.GetPC()&0x7000_0000)<<1)
	if gresBA != 0xffff_ffff && len(areq.result.(agGtMesRespT).result) > 0 {
		memory.WriteStringBA(areq.result.(agGtMesRespT).result, gresBA)
	}
	return true
}

func scGtod(cpu *mvcpu.CPUT, agentChan chan AgentReqT) bool {
	now := time.Now()
	cpu.SetAc(0, dg.DwordT(now.Second()))
	cpu.SetAc(1, dg.DwordT(now.Minute()))
	cpu.SetAc(2, dg.DwordT(now.Hour()))
	return true
}

func scInfo(cpu *mvcpu.CPUT, agentChan chan AgentReqT) bool {
	pktAddr := dg.PhysAddrT(cpu.GetAc(2))
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

func scXpstat(cpu *mvcpu.CPUT, agentChan chan AgentReqT) bool {

	return true
}
