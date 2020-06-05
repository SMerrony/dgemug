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
	"sort"

	"github.com/SMerrony/dgemug/dg"
	"github.com/SMerrony/dgemug/logging"
	"github.com/SMerrony/dgemug/memory"
	"github.com/SMerrony/dgemug/mvcpu"
)

type taskT struct {
	pid, tid                 int
	sixteenBit               bool
	agentChan                chan AgentReqT
	dir                      string
	startAddr, ringMask      dg.PhysAddrT
	wfp, wsp, wsb, wsl, wsfh dg.PhysAddrT
	killAddr                 dg.PhysAddrT
	debugLogging             bool
}

func createTask(pid int, bit16 bool, agentChan chan AgentReqT, startAddr, wfp, wsp, wsb, wsl, sfh dg.PhysAddrT, debugLogging bool) *taskT {
	var task taskT
	task.pid = pid
	task.sixteenBit = bit16
	// task.tid = tid
	task.agentChan = agentChan
	task.startAddr = startAddr
	task.ringMask = startAddr & 0x7000_0000
	if wfp != 0 {
		task.wfp = wfp
	} else {
		task.wfp = wsp
	}
	task.wsp = wsp
	task.wsb = wsb
	task.wsl = wsl
	task.wsfh = sfh
	task.debugLogging = debugLogging

	// get pseudo-Agent to allocate TID
	atreq := agAllocateTIDReqT{pid}
	areq := AgentReqT{agentAllocateTID, atreq, nil}
	agentChan <- areq
	areq = <-agentChan
	if areq.result.(agAllocateTIDRespT).standardTID == 0 {
		log.Panicln("ERROR: Could not allocat TID for new task")
	}
	task.tid = int(areq.result.(agAllocateTIDRespT).standardTID)
	log.Printf("DEBUG: Task %d Created, Initial PC=%#o\n", task.tid, startAddr)
	log.Printf("-----  Start Addr: %#o, WFP: %#o, WSP: %#o, WSB: %#o, WSL: %#o, WSFH: %#o\n", startAddr, wfp, wsp, wsb, wsl, sfh)
	return &task
}

func (task *taskT) run() (errorCode dg.DwordT, termMessage string, flags dg.ByteT) {
	var (
		cpu         mvcpu.CPUT
		syscallTrap bool
		instrCounts [750]int
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
		syscallTrap, _ = cpu.Vrun(&instrCounts)
		if syscallTrap {
			returnAddr := dg.PhysAddrT(cpu.GetAc(3))
			var callID dg.WordT
			if task.sixteenBit {
				ss := memory.NsPop(task.ringMask, false)
				callID = memory.ReadWord(task.ringMask | dg.PhysAddrT(ss))
				memory.NsPush(task.ringMask, ss, false)
			} else {
				callID = memory.ReadWord(dg.PhysAddrT(memory.ReadDWord(cpu.GetWSP() - 2)))
			}
			//log.Printf("DEBUG: System Call %#o return addr %#o\n", callID, returnAddr)
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
					termMessage = string(memory.ReadBytes(cpu.GetAc(1), task.ringMask, msgLen))
				}
				break
			}
			var scOk bool
			if task.sixteenBit {
				scOk = syscall16(callID, task.pid, task.tid, task.ringMask, task.agentChan, &cpu)
				nfp := memory.ReadWord(task.ringMask | memory.NfpLoc)
				cpu.SetAc(3, dg.DwordT(nfp)|dg.DwordT(task.ringMask))
			} else {
				scOk = syscall(callID, task.pid, task.tid, task.ringMask, task.agentChan, &cpu)
				cpu.SetAc(3, dg.DwordT(cpu.GetWFP()))
			}
			mvcpu.WsPop(&cpu)
			if scOk {
				cpu.SetPC(returnAddr + 1)
			} else {
				cpu.SetPC(returnAddr)
			}
			//cpu.SetAc(3, dg.DwordT(cpu.GetWFP()))
		} else {
			// Vrun has stopped and we're not at a system call
			break
		}
	}

	// instruction counts, first by Mnemonic, then by count
	m := make(map[int]string)
	keys := make([]int, 0)

	log.Println("Instruction Execution Count by Mnemonic")
	for i, c := range instrCounts {
		if instrCounts[i] > 0 {
			log.Printf("%s\t%d\n", mvcpu.GetMnemonic(i), c)
			if m[c] == "" {
				m[c] = mvcpu.GetMnemonic(i)
				keys = append(keys, c)
			} else {
				m[c] += ", " + mvcpu.GetMnemonic(i)
			}
		}
	}
	log.Println("instructions by Count")
	sort.Ints(keys)
	for _, c := range keys {
		log.Printf("%d\t%s\n", c, m[c])
	}

	return errorCode, termMessage, flags
}
