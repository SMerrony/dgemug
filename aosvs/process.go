// +build virtual !physical

// process.go - abstraction of an AOS/VS process

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
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"

	"github.com/SMerrony/dgemug/dg"
	"github.com/SMerrony/dgemug/memory"
)

const (
	maxTasksPerProc = 32
	firstTask       = 1
	pcInPr          = 0574
	page8offset     = 8192
	sfhInPr         = page8offset + 12
	wfpInPr         = page8offset + 16
	wspInPr         = page8offset + 18
	wslInPr         = page8offset + 20
	wsbInPr         = page8offset + 22
)

type ustT struct {
	extVarWdCnt         dg.WordT
	extVarP0Start       dg.WordT
	symsStart           dg.DwordT
	symsEnd             dg.DwordT
	debugAddr           dg.PhysAddrT
	revision            dg.DwordT
	taskCount           dg.WordT
	impureBlocks        dg.WordT
	sharedStartBlock    dg.DwordT // this seems to be a page #, not a block #
	intAddr             dg.PhysAddrT
	sharedBlockCount    dg.WordT // this seems to be a page count, not blocks
	prType              dg.WordT
	sharedStartPageInPR dg.DwordT
}

// ProcessT represents an AOS/VS Process which will contain one or more Tasks
type ProcessT struct {
	PID             dg.WordT
	name            string
	programFileName string
	console         net.Conn
	ust             ustT
}

var debugLogging bool

// CreateProcess creates, but does not start, an emulated AOS/VS Process
func CreateProcess(args []string, vRoot string, prName string, ring int, con net.Conn, agentChan chan AgentReqT, debugLog bool) (err error) {
	debugLogging = debugLog
	progWds, err := readProgram(prName)
	if err != nil {
		return err
	}
	var proc ProcessT
	proc.programFileName = prName
	proc.console = con

	// get info from PR preamble
	proc.loadUST(progWds)
	proc.printUST()

	// Announce ourself to pseudo-Agent and get PID
	var areq AgentReqT
	areq.action = agentAllocatePID
	areq.reqParms = agAllocatePIDReqT{invocationArgs: args, virtualRoot: vRoot, sixteenBit: proc.ust.prType&0x8000 != 0}
	agentChan <- areq
	areq = <-agentChan
	if !areq.result.(agAllocatePIDRespT).ok {
		log.Panic("ERROR: Could not get PID from pseudo-Agent")
	}
	proc.PID = areq.result.(agAllocatePIDRespT).PID
	log.Printf("INFO: Obtained PID %d for process\n", proc.PID)
	log.Printf("INFO: Preparing ring %d process with up to %d tasks\n", ring, proc.ust.taskCount)
	log.Printf("----  PR: %s  Args: %v\n", prName, args)
	if proc.ust.sharedBlockCount != 0 {
		log.Println("WARNING: Shared pages not yet fully supported")
	}

	segBase := dg.PhysAddrT(ring) << 28

	// map (load) program into RAM
	// unshared portion
	memory.MapSlice(segBase, progWds[8192:proc.ust.sharedStartPageInPR<<10-8], false)
	// shared portion
	memory.MapSlice(segBase+dg.PhysAddrT(proc.ust.sharedStartBlock)<<10, progWds[proc.ust.sharedStartPageInPR<<10:], true)

	// set up initial task
	var taskReq agTaskReqT
	taskReq.PID = proc.PID
	taskReq.TID = firstTask
	taskReq.agentChan = agentChan
	taskReq.conn = con
	taskReq.initAC2 = 0
	taskReq.priority = 0
	taskReq.startAddr = dg.PhysAddrT(memory.DwordFromTwoWords(progWds[pcInPr], progWds[pcInPr+1]))
	taskReq.wfp = dg.PhysAddrT(memory.DwordFromTwoWords(progWds[wfpInPr], progWds[wfpInPr+1]))
	taskReq.wsb = dg.PhysAddrT(memory.DwordFromTwoWords(progWds[wsbInPr], progWds[wsbInPr+1]))
	taskReq.wsfh = dg.PhysAddrT(progWds[sfhInPr])
	taskReq.wsl = dg.PhysAddrT(memory.DwordFromTwoWords(progWds[wslInPr], progWds[wslInPr+1]))
	taskReq.wsp = dg.PhysAddrT(memory.DwordFromTwoWords(progWds[wspInPr], progWds[wspInPr+1]))

	areq = AgentReqT{agentTask, taskReq, nil}
	agentChan <- areq
	areq = <-agentChan

	return nil
}

func readProgram(prName string) (wordImg []dg.WordT, err error) {
	userProgramBytes, err := ioutil.ReadFile(prName)
	if err != nil {
		return nil, fmt.Errorf("could not read PR file %s", prName)
	}
	if len(userProgramBytes) == 0 {
		return nil, errors.New("empty PR file supplied")
	}

	progWds := make([]dg.WordT, len(userProgramBytes)/2)
	var word int
	for word = 0; word < len(userProgramBytes)/2; word++ {
		progWds[word] = dg.WordT(userProgramBytes[word*2])<<8 | dg.WordT(userProgramBytes[word*2+1])
	}
	log.Printf("DEBUG: Read user program %s containing %#x words\n", prName, word)
	return progWds, nil
}

// // Run launches a new AOS/VS process and its initial Task
// func (proc *ProcessT) Run() (errorCode dg.DwordT, termMessage string, flags dg.ByteT) {

// 	errorCode, termMessage, flags = proc.tasks[firstTask].run()
// 	return errorCode, termMessage, flags

// }

func (proc *ProcessT) loadUST(progWds []dg.WordT) {
	proc.ust.extVarWdCnt = progWds[ust+ustez]
	proc.ust.extVarP0Start = progWds[ust+ustes]
	proc.ust.symsStart = memory.DwordFromTwoWords(progWds[ust+ustss], progWds[ust+ustss+1])
	proc.ust.symsEnd = memory.DwordFromTwoWords(progWds[ust+ustse], progWds[ust+ustse+1])
	proc.ust.debugAddr = dg.PhysAddrT(memory.DwordFromTwoWords(progWds[ust+ustda], progWds[ust+ustda+1]))
	proc.ust.revision = memory.DwordFromTwoWords(progWds[ust+ustrv], progWds[ust+ustrv+1])
	proc.ust.taskCount = progWds[ust+usttc] // 266.
	//proc.ust.impureBlocks = progWds[ust+ustbl]
	proc.ust.impureBlocks = progWds[ust+ustbl+1]
	proc.ust.sharedStartBlock = memory.DwordFromTwoWords(progWds[ust+ustst], progWds[ust+ustst+1]) // 0x10f 271. 0417
	proc.ust.intAddr = dg.PhysAddrT(memory.DwordFromTwoWords(progWds[ust+ustit], progWds[ust+ustit+1]))
	proc.ust.sharedBlockCount = progWds[ust+ustsz]                                                    // 0x113 275. 0423
	proc.ust.sharedStartPageInPR = memory.DwordFromTwoWords(progWds[ust+ustsh], progWds[ust+ustsh+1]) // 0x11a 282. 0432

	proc.ust.prType = progWds[ust+ustpr]
	// proc.ust.sharedStartBlock = dg.DwordT(progWds[0x10f])    // 271. 0417
	// proc.ust.sharedBlockCount = progWds[0x113]               // 275. 0423
	// //proc.ust.sharedStartPageInPR = dg.DwordT(progWds[0x11a]) // 282. 0432
	// //proc.ust.taskCount = progWds[0x10a]                      // 266. 0412
}

func (proc *ProcessT) printUST() {
	log.Printf("UST: Exended Variables - word count: %d., page 0 start: %#x\n", proc.ust.extVarWdCnt, proc.ust.extVarP0Start)
	log.Printf("     Symbols - start: %#x, end: %#x, debug addr: %#x\n", proc.ust.symsStart, proc.ust.symsEnd, proc.ust.debugAddr)
	log.Printf("     Program - revision: %#x, task count: %d., type: %d.", proc.ust.revision, proc.ust.taskCount, proc.ust.prType)
	log.Printf("     Shared - start block: %#x, # blocks: %#x, start page in .PR: %#x\n", proc.ust.sharedStartBlock, proc.ust.sharedBlockCount, proc.ust.sharedStartPageInPR)
	log.Printf("     Impure Blocks: %d., Interrupt Addr: %#x, ", proc.ust.impureBlocks, proc.ust.intAddr)
}
