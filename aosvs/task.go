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
	"github.com/SMerrony/dgemug/memory"
	"github.com/SMerrony/dgemug/mvcpu"
)

type taskT struct {
	pid, tid                              int
	agentChan                             chan AgentReqT
	dir                                   string
	startAddr                             dg.PhysAddrT
	wfp, wsp, wsb, wsl, stackFaultHandler dg.PhysAddrT
}

func createTask(pid int, tid int, agent chan AgentReqT, startAddr, wfp, wsp, wsb, wsl, sfh dg.PhysAddrT) *taskT {
	var task taskT
	task.pid = pid
	task.tid = tid
	task.agentChan = agent
	task.startAddr = startAddr
	task.wfp = wfp
	task.wsp = wsp
	task.wsb = wsb
	task.wsl = wsl
	task.stackFaultHandler = sfh

	log.Printf("DEBUG: Task %d Created, Initial PC=%#o\n", tid, startAddr)
	return &task
}

func (task *taskT) run() (errDetail string, instrCounts [750]int) {
	var (
		cpu         mvcpu.CPUT
		syscallTrap bool
	)

	cpu.CPUInit(077, nil, nil)
	cpu.SetupStack(task.wfp, task.wsp, task.wsb, task.wsl)
	cpu.SetATU(true)
	cpu.SetPC(task.startAddr)

	// log.Println(cpu.DisassembleRange(0x7000_0000, 0x7000_0020))
	// log.Println(cpu.DisassembleRange(0x7007_fc00, 0x7007_fc00+0400))

	for {
		syscallTrap, errDetail, instrCounts = cpu.Vrun()
		if syscallTrap {
			returnAddr := dg.PhysAddrT(cpu.GetAc(3))
			callID := memory.ReadWord(cpu.GetPC() + 2)
			//log.Printf("DEBUG: Trapped System Call: %#o, return addr: %#o\n", callID, returnAddr)
			if syscall(callID, task.agentChan, &cpu) {
				cpu.SetPC(returnAddr + 2)
			} else {
				cpu.SetPC(returnAddr + 1)
			}
		} else {
			// Vrun has stopped and we're not at a system call
			break
		}
	}
	return errDetail, instrCounts
}
