// syscalls.go - map of AOS/VS system calls

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
	"github.com/SMerrony/dgemug/mvcpu"
)

type syscallDescT struct {
	name  string                 // AOS/VS System Call Name
	alias string                 // AOS/VS 4-char Name Alias
	fn    func(*mvcpu.CPUT) bool // implementation
}

var syscalls = map[dg.WordT]syscallDescT{
	0:    {"?CREATE", "?CREA", nil},
	1:    {"?DELETE", "?DELE", nil},
	0300: {"?OPEN", "?OPEN", nil},
	0301: {"?CLOSE", "?CLOS", nil},
	0302: {"?READ", "?READ", nil},
	0303: {"?WRITE", "?WRIT", nil},
	0310: {"?RETURN", "?RETU", nil},
}

func syscall(callID dg.WordT, cpu *mvcpu.CPUT) (ok bool) {
	call, defined := syscalls[callID]
	if !defined {
		log.Fatalf("ERROR: System call No. %#od not yet defined at PC=%#od", callID, cpu.GetPC())
	}
	if call.fn == nil {
		log.Fatalf("ERROR: System call No. %#od not yet implemented at PC=%#od", callID, cpu.GetPC())
	}
	return call.fn(cpu)
}
