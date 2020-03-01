// mvemg project eaglePC_test.go

// Copyright ©2017-2020  Steve Merrony

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
	"testing"

	"github.com/SMerrony/dgemug/memory"
)

func TestISZTS(t *testing.T) {
	cpuPtr := new(MvCPUT)
	var iPtr decodedInstrT
	iPtr.ix = instrISZTS
	memory.MemInit(1000, false)

	cpuPtr.wsp = 100              // set the WSP
	wsPush(cpuPtr, 0, 0xfffffffe) // push -2
	cpuPtr.pc = 7000

	if !eaglePC(cpuPtr, &iPtr) {
		t.Error("Failed to execute ISZTS")
	}
	if cpuPtr.pc != 7001 {
		t.Errorf("Expected PC to be 7001, got %d", cpuPtr.pc)
	}
	v := memory.ReadDWord(cpuPtr.wsp)
	if v != 0xffffffff {
		t.Errorf("Expected 0xffffffff at WSP, got: %#x", v)
	}

	cpuPtr.pc = 7000
	if !eaglePC(cpuPtr, &iPtr) {
		t.Error("Failed to execute ISZTS")
	}
	if cpuPtr.pc != 7002 {
		t.Errorf("Expected PC to be 7002, got %d", cpuPtr.pc)
	}
	v = memory.ReadDWord(cpuPtr.wsp)
	if v != 0 {
		t.Errorf("Expected 0 at WSP, got: %#x", v)
	}
}

func TestWSKBO(t *testing.T) {
	cpuPtr := new(MvCPUT)
	var iPtr decodedInstrT
	var wskb wskbT
	iPtr.ix = instrWSKBO
	cpuPtr.ac[0] = 0x55555555 // 0101010101010101...
	wskb.bitNum = 1           // 2nd from left
	iPtr.variant = wskb
	cpuPtr.pc = 1000
	if !eaglePC(cpuPtr, &iPtr) {
		t.Error("Failed to execute WSKBO")
	}
	if cpuPtr.pc != 1002 {
		t.Errorf("Expected PC: 1002., got %d.", cpuPtr.pc)
	}
	wskb.bitNum = 20
	iPtr.variant = wskb
	cpuPtr.pc = 1000
	if !eaglePC(cpuPtr, &iPtr) {
		t.Error("Failed to execute WSKBO")
	}
	if cpuPtr.pc != 1001 {
		t.Errorf("Expected PC: 1001., got %d.", cpuPtr.pc)
	}
}

func TestWSLE(t *testing.T) {
	cpuPtr := new(MvCPUT)
	var iPtr decodedInstrT
	var twoAcc1Word twoAcc1WordT
	twoAcc1Word.acs = 0
	twoAcc1Word.acd = 1
	iPtr.ix = instrWSLE
	iPtr.variant = twoAcc1Word
	cpuPtr.ac[0] = 0
	cpuPtr.ac[1] = 1
	cpuPtr.pc = 1000
	if !eaglePC(cpuPtr, &iPtr) {
		t.Error("Failed to execute WSLE")
	}
	if cpuPtr.pc != 1002 {
		t.Errorf("Expected PC: 1002., got %d.", cpuPtr.pc)
	}

	cpuPtr.ac[0] = 1
	cpuPtr.ac[1] = 0
	cpuPtr.pc = 1000
	if !eaglePC(cpuPtr, &iPtr) {
		t.Error("Failed to execute WSLE")
	}
	if cpuPtr.pc != 1001 {
		t.Errorf("Expected PC: 1001., got %d.", cpuPtr.pc)
	}

	cpuPtr.ac[0] = 1
	cpuPtr.ac[1] = 1
	cpuPtr.pc = 1000
	if !eaglePC(cpuPtr, &iPtr) {
		t.Error("Failed to execute WSLE")
	}
	if cpuPtr.pc != 1002 {
		t.Errorf("Expected PC: 1002., got %d.", cpuPtr.pc)
	}

	cpuPtr.ac[0] = 0xfffffffe // -1
	cpuPtr.ac[1] = 1
	cpuPtr.pc = 1000
	if !eaglePC(cpuPtr, &iPtr) {
		t.Error("Failed to execute WSLE")
	}
	if cpuPtr.pc != 1002 {
		t.Errorf("Expected PC: 1002., got %d.", cpuPtr.pc)
	}

	cpuPtr.ac[0] = 0xfffffffe // -1
	cpuPtr.ac[1] = 0xfffffffd
	cpuPtr.pc = 1000
	if !eaglePC(cpuPtr, &iPtr) {
		t.Error("Failed to execute WSLE")
	}
	if cpuPtr.pc != 1001 {
		t.Errorf("Expected PC: 1001., got %d.", cpuPtr.pc)
	}
}

func TestXNISZ(t *testing.T) {
	cpuPtr := new(MvCPUT)
	var iPtr decodedInstrT
	var noAccModeInd2Word noAccModeInd2WordT
	iPtr.ix = instrXNISZ
	memory.MemInit(10000, false)
	memory.WriteWord(100, 0xfffe) // write max - 1 into Word at normal addr 100
	noAccModeInd2Word.disp15 = 100
	noAccModeInd2Word.ind = ' '
	noAccModeInd2Word.mode = absoluteMode
	iPtr.variant = noAccModeInd2Word
	cpuPtr.pc = 1000
	if !eaglePC(cpuPtr, &iPtr) {
		t.Error("Failed to execute XNISZ")
	}
	// 1st time should simply increment contents
	if cpuPtr.pc != 1002 {
		t.Errorf("Expected PC: 1002., got %d.", cpuPtr.pc)
	}
	w := memory.ReadWord(100)
	if w != 0xffff {
		t.Errorf("Expected loc 100. to contain 0xffff, got %x", w)
	}
	// Again...
	cpuPtr.pc = 1000
	if !eaglePC(cpuPtr, &iPtr) {
		t.Error("Failed to execute XNISZ")
	}
	// 2nd time should roll over and skip contents
	if cpuPtr.pc != 1003 {
		t.Errorf("Expected PC: 1003., got %d.", cpuPtr.pc)
	}
	w = memory.ReadWord(100)
	if w != 0 {
		t.Errorf("Expected loc 100. to contain 0, got %x", w)
	}
}

func TestXWDSZ(t *testing.T) {
	cpuPtr := new(MvCPUT)
	var iPtr decodedInstrT
	var noAccModeInd2Word noAccModeInd2WordT
	iPtr.ix = instrXWDSZ
	memory.MemInit(10000, false)
	memory.WriteDWord(100, 2) // write 2 into DWord at normal addr 100
	noAccModeInd2Word.disp15 = 100
	noAccModeInd2Word.ind = ' '
	noAccModeInd2Word.mode = absoluteMode
	iPtr.variant = noAccModeInd2Word
	cpuPtr.pc = 1000
	if !eaglePC(cpuPtr, &iPtr) {
		t.Error("Failed to execute XWDSZ")
	}
	// 1st time should simply decrement contents
	if cpuPtr.pc != 1002 {
		t.Errorf("Expected PC: 1002., got %d.", cpuPtr.pc)
	}
	w := memory.ReadDWord(100)
	if w != 1 {
		t.Errorf("Expected loc 100. to contain 1, got %d", w)
	}
	// 2nd time should dec and skip
	cpuPtr.pc = 1000
	if !eaglePC(cpuPtr, &iPtr) {
		t.Error("Failed to execute XWDSZ")
	}
	if cpuPtr.pc != 1003 {
		t.Errorf("Expected PC: 1003., got %d.", cpuPtr.pc)
	}
	w = memory.ReadDWord(100)
	if w != 0 {
		t.Errorf("Expected loc 100. to contain 0, got %d", w)
	}
}
