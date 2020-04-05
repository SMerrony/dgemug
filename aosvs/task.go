// task.go - abstraction of an AOS/VS task

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
	"github.com/SMerrony/dgemug/logging"
	"github.com/SMerrony/dgemug/memory"
	"github.com/SMerrony/dgemug/mvcpu"
)

type taskT struct {
	pid, tid                 int
	agentChan                chan AgentReqT
	dir                      string
	startAddr                dg.PhysAddrT
	wfp, wsp, wsb, wsl, wsfh dg.PhysAddrT
	killAddr                 dg.PhysAddrT
	debugLogging             bool
}

func createTask(pid int, tid int, agent chan AgentReqT, startAddr, wfp, wsp, wsb, wsl, sfh dg.PhysAddrT, debugLogging bool) *taskT {
	var task taskT
	task.pid = pid
	task.tid = tid
	task.agentChan = agent
	task.startAddr = startAddr
	task.wfp = wfp
	task.wsp = wsp
	task.wsb = wsb
	task.wsl = wsl
	task.wsfh = sfh
	task.debugLogging = debugLogging

	log.Printf("DEBUG: Task %d Created, Initial PC=%#o\n", tid, startAddr)
	log.Printf("-----  Start Addr: %#o, WFP: %#o, WSP: %#o, WSB: %#o, WSL: %#o, WSFH: %#o\n", startAddr, wfp, wsp, wsb, wsl, sfh)
	return &task
}

func (task *taskT) run() (errorCode dg.DwordT, termMessage string, flags dg.ByteT) {
	var (
		cpu         mvcpu.CPUT
		syscallTrap int
	)

	cpu.CPUInit(077, nil, nil)
	cpu.SetPC(task.startAddr) // must be done before stack set up
	cpu.SetupStack(task.wfp, task.wsp, task.wsb, task.wsl, task.wsfh)
	adjustedWsfh := (cpu.GetPC() & 0x7000_0000) | dg.PhysAddrT(memory.ReadWord((cpu.GetPC()&0x7000_0000)|014)) // just for debugging
	log.Printf("----- Wide Stack Fault Handler reset to: %#x (%#o)\n", adjustedWsfh, adjustedWsfh)
	cpu.SetATU(true)
	cpu.SetDebugLogging(task.debugLogging)

	// log.Println(cpu.DisassembleRange(0x7000_0000, 0x7000_0020))
	// log.Println(cpu.DisassembleRange(0x7007_fc00, 0x7007_fc00+0400))

	for {
		syscallTrap, _, _ = cpu.Vrun()
		if syscallTrap != mvcpu.SyscallNot {
			returnAddr := cpu.GetPC() + 2 // dg.PhysAddrT(cpu.GetAc(3))
			callID := memory.ReadWord(returnAddr)
			// special handling for the ?RETURN system call
			if callID == scReturn {
				if task.debugLogging {
					logging.DebugPrint(logging.DebugLog, "?RETURN")
				}
				log.Println("INFO: ?RETURN")
				errorCode = cpu.GetAc(0)
				flags = dg.ByteT(memory.GetDwbits(cpu.GetAc(2), 16, 8))
				msgLen := int(uint8(memory.GetDwbits(cpu.GetAc(2), 24, 8)))
				if msgLen > 0 {
					termMessage = string(memory.ReadBytes(cpu.GetAc(1), cpu.GetPC(), msgLen))
				}
				break
			}
			var scOk bool
			switch syscallTrap {
			case mvcpu.Syscall32Trap:
				scOk = syscall(callID, task.agentChan, &cpu)
				cpu.SetAc(3, dg.DwordT(cpu.GetWFP()))
			case mvcpu.Syscall16Trap:
				scOk = syscall16(callID, task.agentChan, &cpu)
			}
			if scOk {
				cpu.SetPC(returnAddr + 2)
			} else {
				cpu.SetPC(returnAddr + 1)
			}
			//cpu.SetAc(3, dg.DwordT(cpu.GetWFP()))
		} else {
			// Vrun has stopped and we're not at a system call
			break
		}
	}
	return errorCode, termMessage, flags
}
