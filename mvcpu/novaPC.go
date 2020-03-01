// novaPC.go

// Copyright (C) 2017,2020 Steve Merrony

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
)

func novaPC(cpuPtr *MvCPUT, iPtr *decodedInstrT) bool {

	var novaNoAccEffAddr novaNoAccEffAddrT

	switch iPtr.ix {

	case instrJMP:
		novaNoAccEffAddr = iPtr.variant.(novaNoAccEffAddrT)
		cpuPtr.pc = resolve8bitDisplacement(cpuPtr, novaNoAccEffAddr.ind, novaNoAccEffAddr.mode, novaNoAccEffAddr.disp15) & 0x7fff

	case instrJSR:
		novaNoAccEffAddr = iPtr.variant.(novaNoAccEffAddrT)
		tmpPC := dg.DwordT(cpuPtr.pc + 1)
		cpuPtr.pc = resolve8bitDisplacement(cpuPtr, novaNoAccEffAddr.ind, novaNoAccEffAddr.mode, novaNoAccEffAddr.disp15) & 0x7fff
		cpuPtr.ac[3] = tmpPC

	default:
		log.Fatalf("ERROR: NOVA_PC instruction <%s> not yet implemented\n", iPtr.mnemonic)
		return false
	}
	return true
}
