// scProcess.go - 'Process Management'-related System Call Emulation

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

func scDadid(cpu *mvcpu.CPUT, agentChan chan AgentReqT) bool {
	// TODO this should really get passed off to a central 'process manager' (pseudo-EXEC/PMGR?)
	// fake PIDs
	switch cpu.GetAc(0) {
	case 1:
		cpu.SetAc(1, 0)
	case 2:
		cpu.SetAc(1, 1)
	case 3:
		cpu.SetAc(1, 2)
	case 5:
		cpu.SetAc(1, 3)
	case 6:
		cpu.SetAc(1, 5)
	case 7:
		cpu.SetAc(1, 6)
	case 8:
		cpu.SetAc(1, 7)
	case 9:
		cpu.SetAc(1, 8)
	case 10:
		cpu.SetAc(1, 9)
	default:
		cpu.SetAc(1, 10)
	}
	return true
}

func scGunm(cpu *mvcpu.CPUT, agentChan chan AgentReqT) bool {
	// TODO this should really get passed off to a central 'process manager' (pseudo-EXEC/PMGR?)
	if dg.WordT(cpu.GetAc(0)) == 0xffff {
		cpu.SetAc(0, 1)      // Claim not to be in SU mode
		cpu.SetAc(1, 0x001f) // Claim to have nearly all privileges
		memory.WriteStringBA("XYZZY", cpu.GetAc(2))
	} else {
		log.Panic("ERROR: ?GUNM request type not yet implemented")
	}
	return true
}

func scSysprv(cpu *mvcpu.CPUT, agentChan chan AgentReqT) bool {
	pktAddr := dg.PhysAddrT(cpu.GetAc(2))
	funcCode := memory.ReadWord(pktAddr + sysprvPktFunc)
	switch funcCode {
	case sysprvGet:
		memory.WriteWord(pktAddr+sysprvPktFlags, 0)
	case sysprvEnter:
		log.Panicln("ERROR: Enter func not yet implemented in ?SYSPRV")
	case sysprvEnterExcl:
		log.Panicln("ERROR: Enter Exclusive func not yet implemented in ?SYSPRV")
	case sysprvLeave:
		log.Println("WARNING: Leave func not yet implemented in ?SYSPRV - Ignoring")
	default:
		log.Panicf("ERROR: ?SYSPRV called with unknown function code %d", funcCode)
	}
	return true
}
