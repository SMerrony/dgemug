// scMemory.go - Memory-related System Call Emulation

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

func scMem(cpu *mvcpu.CPUT, agentChan chan AgentReqT) bool {
	highUnshared := memory.GetLastUnsharedPage()
	lowShared := memory.GetFirstSharedPage() & (0x0fff_ffff >> 10)
	cpu.SetAc(0, lowShared-highUnshared)
	cpu.SetAc(1, highUnshared)
	cpu.SetAc(2, ((dg.DwordT(cpu.GetPC())&0x7000_0000)|highUnshared<<10)-1)
	return true
}

func scMemi(cpu *mvcpu.CPUT, agentChan chan AgentReqT) bool {
	numPages := int32(cpu.GetAc(0))
	var lastPage int
	switch {
	case numPages > 0: // add pages
		for numPages > 0 {
			lastPage = memory.AddUnsharedPage()
			numPages--
		}
		cpu.SetAc(1, dg.DwordT(lastPage<<10)|dg.DwordT(cpu.GetPC()&0x7000_0000)-1)
	case numPages < 0: // remove pages
		log.Fatalln("ERROR: Unmapping via ?MEMI not yet supported")
	}
	return true
}
