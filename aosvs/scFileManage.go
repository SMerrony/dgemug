// scFileManage.go - File Management System Call Emulation

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
	"github.com/SMerrony/dgemug/memory"
	"github.com/SMerrony/dgemug/mvcpu"
)

func scDacl(cpu *mvcpu.CPUT, PID int, agentChan chan AgentReqT) bool {
	// TODO this is all faked...
	switch dg.WordT(cpu.GetAc(0)) { // make 16-bit safe
	case 0xffff:
		log.Fatal("ERROR: Setting new DefACL not yet implemented in ?DACL")
	case 0:
		defacl := []byte("XYZZY")
		defacl = append(defacl, 0)
		defacl = append(defacl, faco+facw+faca+facw+face)
		defacl = append(defacl, 0)
		memory.WriteBytesBA(defacl, cpu.GetAc(1))
	case 1:
		log.Fatal("ERROR: Turning off DefACL not yet implemented in ?DACL")
	}
	return true
}

func scGname(cpu *mvcpu.CPUT, PID int, agentChan chan AgentReqT) bool {
	bpFilename := cpu.GetAc(0)
	filename := strings.ToUpper(readString(bpFilename, cpu.GetPC()))
	if filename[0] == '@' {
		filename = ":PER:" + filename[1:] // convert @ to :PER
	}
	filename = strings.ReplaceAll(filename, "/", ":") // convert any / to :
	bpPathname := cpu.GetAc(1)
	writeBytes(bpPathname, cpu.GetPC(), []byte(filename))
	cpu.SetAc(2, dg.DwordT(len(filename)))
	log.Printf("DEBUG: ?GNAME returning %s for %s\n", readString(bpFilename, cpu.GetPC()), filename)
	return true
}

func scRecreate(cpu *mvcpu.CPUT, PID int, agentChan chan AgentReqT) bool {
	bpFilename := cpu.GetAc(0)
	filename := strings.ToUpper(readString(bpFilename, cpu.GetPC()))
	var recReq = agRecreateReqT{PID: PID, aosFilename: filename}
	var areq = AgentReqT{agentFileRecreate, recReq, nil}
	agentChan <- areq
	areq = <-agentChan
	resp := areq.result.(agRecreateRespT)
	if !resp.ok {
		cpu.SetAc(0, resp.errCode)
		return false
	}
	return true
}
