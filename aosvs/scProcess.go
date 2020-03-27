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

func scSysprv(cpu *mvcpu.CPUT, agentChan chan AgentReqT) bool {
	pktAddr := dg.PhysAddrT(cpu.GetAc(2))
	funcCode := memory.ReadWord(pktAddr + sysprvPktFunc)
	switch funcCode {
	case sysprvGet:
		memory.WriteWord(pktAddr+sysprvPktFlags, 0)
	case sysprvEnter:
		log.Fatalln("ERROR: Enter func not yet implemented in ?SYSPRV")
	case sysprvEnterExcl:
		log.Fatalln("ERROR: Enter Exclusive func not yet implemented in ?SYSPRV")
	case sysprvLeave:
		log.Println("WARNING: Leave func not yet implemented in ?SYSPRV - Ignoring")
	default:
		log.Fatalf("ERROR: ?SYSPRV called with unknown function code %d", funcCode)
	}
	return true
}
