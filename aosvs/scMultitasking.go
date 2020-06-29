// +build virtual !physical

// scM.go - File I/O System Call Emulation

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
	"github.com/SMerrony/dgemug/logging"
	"github.com/SMerrony/dgemug/memory"
)

func scIfpu(p syscallParmsT) bool {
	// TODO should reserve FPU save area
	return true
}

// func scKilad(cpu *mvcpu.CPUT, agentChan chan AgentReqT) bool {

// 	return true
// }

func scTask(p syscallParmsT) bool {
	tpa := dg.PhysAddrT(p.cpu.GetAc(2))
	if memory.ReadDWord(tpa+dlnk) == 0 {
		log.Panicln("?TASK extended packets not yet implemented")
	}
	var tskData agTaskReqT
	tskData.priority = memory.ReadWord(tpa + dpri)
	tskData.TID = memory.ReadWord(tpa + did)
	tskData.startAddr = dg.PhysAddrT(memory.ReadDWord(tpa + dpc))
	tskData.initAC2 = memory.ReadDWord(tpa + dac2)
	tskData.wsb = dg.PhysAddrT(memory.ReadDWord(tpa + dstb))
	tskData.wsfh = (p.cpu.GetPC() & 0x7000_0000) | dg.PhysAddrT(memory.ReadWord(tpa+dsflt))
	tskData.wsl = tskData.wsb + dg.PhysAddrT(memory.ReadDWord(tpa+dssz))

	return true
}

func scUidstat(p syscallParmsT) bool {
	reqTID := p.cpu.GetAc(1)
	if reqTID != 0xffff_ffff {
		log.Panicln("?UIDSTAT request for another TID not yet implemented")
	}
	retPacketAddr := dg.PhysAddrT(p.cpu.GetAc(2))
	memory.WriteWord(retPacketAddr, dg.WordT(p.TID))
	memory.WriteWord(retPacketAddr+1, 0)
	memory.WriteWord(retPacketAddr+2, dg.WordT(p.TID))
	memory.WriteWord(retPacketAddr+3, 0)
	logging.DebugPrint(logging.ScLog, "-------- Returning UTID: %#o, STID: %#o\n", dg.WordT(p.TID), p.TID)
	return true
}

func scWdelay(p syscallParmsT) bool {
	delayMs := int(p.cpu.GetAc(0))
	time.Sleep(time.Millisecond * time.Duration(delayMs))
	return true
}
