// novaMemRef.go

// Copyright (C) 2017  Steve Merrony

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

package mvcpu

import (
	"log"

	"github.com/SMerrony/dgemug/dg"
	"github.com/SMerrony/dgemug/memory"
)

func novaMath(cpuPtr *MvCPUT, iPtr *decodedInstrT) bool {

	switch iPtr.ix {
	case instrDIV: // unsigned divide
		uw := memory.DwordGetLowerWord(cpuPtr.ac[0])
		lw := memory.DwordGetLowerWord(cpuPtr.ac[1])
		dwd := memory.DwordFromTwoWords(uw, lw)
		quot := memory.DwordGetLowerWord(cpuPtr.ac[2])
		if uw >= quot || quot == 0 {
			cpuPtr.carry = true
		} else {
			cpuPtr.carry = false
			cpuPtr.ac[0] = (dwd % dg.DwordT(quot)) & 0x0ffff
			cpuPtr.ac[1] = (dwd / dg.DwordT(quot)) & 0x0ffff
		}

	case instrMUL: // unsigned 16-bit multiply with add: (AC1 * AC2) + AC0 => AC0(h) and AC1(l)
		ac0 := memory.DwordGetLowerWord(cpuPtr.ac[0])
		ac1 := memory.DwordGetLowerWord(cpuPtr.ac[1])
		ac2 := memory.DwordGetLowerWord(cpuPtr.ac[2])
		dwd := (dg.DwordT(ac1) * dg.DwordT(ac2)) + dg.DwordT(ac0)
		cpuPtr.ac[0] = dg.DwordT(memory.DwordGetUpperWord(dwd))
		cpuPtr.ac[1] = dg.DwordT(memory.DwordGetLowerWord(dwd))

	default:
		log.Fatalf("ERROR: NOVA_MATH instruction <%s> not yet implemented\n", iPtr.mnemonic)
		return false
	}

	cpuPtr.pc++
	return true
}
