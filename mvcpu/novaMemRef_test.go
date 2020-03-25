// +build physical !virtual

// novaMemRef_test.go

// Copyright (C) 2019,2020 Steve Merrony

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

func TestDSZ(t *testing.T) {
	cpu := new(CPUT)
	var iPtr decodedInstrT
	iPtr.ix = instrDSZ
	memory.MemInit(10000, false)
	memory.WriteWord(100, 2)
	cpu.pc = 10
	iPtr.disp15 = 100
	iPtr.ind = ' '
	iPtr.mode = absoluteMode

	if !novaMemRef(cpu, &iPtr) {
		t.Error("Failed to execute DSZ")
	}
	if cpu.pc != 11 {
		t.Errorf("Expected PC 11, got %d.", cpu.pc)
	}
	w := memory.ReadWord(100)
	if w != 1 {
		t.Errorf("Expected loc 100 to contain 1, got: %x", w)
	}

	if !novaMemRef(cpu, &iPtr) {
		t.Error("Failed to execute DSZ")
	}
	if cpu.pc != 13 {
		t.Errorf("Expected PC 13, got %d.", cpu.pc)
	}
	w = memory.ReadWord(100)
	if w != 0 {
		t.Errorf("Expected loc 100 to contain 0, got: %x", w)
	}
}

func TestISZ(t *testing.T) {
	cpu := new(CPUT)
	var iPtr decodedInstrT
	iPtr.ix = instrISZ
	memory.MemInit(10000, false)
	memory.WriteWord(100, 0xfffe)
	cpu.pc = 10
	iPtr.disp15 = 100
	iPtr.ind = ' '
	iPtr.mode = absoluteMode

	if !novaMemRef(cpu, &iPtr) {
		t.Error("Failed to execute ISZ")
	}
	if cpu.pc != 11 {
		t.Errorf("Expected PC 11, got %d.", cpu.pc)
	}
	w := memory.ReadWord(100)
	if w != 0xffff {
		t.Errorf("Expected loc 100 to contain 0xffff, got: %x", w)
	}

	if !novaMemRef(cpu, &iPtr) {
		t.Error("Failed to execute ISZ")
	}
	if cpu.pc != 13 {
		t.Errorf("Expected PC 13, got %d.", cpu.pc)
	}
	w = memory.ReadWord(100)
	if w != 0 {
		t.Errorf("Expected loc 500 to contain 0, got: %x", w)
	}
}

func TestSTA(t *testing.T) {
	cpu := new(CPUT)
	var iPtr decodedInstrT
	var novaOneAccEffAddr novaOneAccEffAddrT
	iPtr.ix = instrSTA
	cpu.ac[1] = 0x12345678
	novaOneAccEffAddr.acd = 1
	novaOneAccEffAddr.disp15 = 100
	novaOneAccEffAddr.ind = ' '
	novaOneAccEffAddr.mode = absoluteMode
	memory.MemInit(10000, false)
	memory.WriteWord(100, 0xfffe)
	iPtr.variant = novaOneAccEffAddr

	if !novaMemRef(cpu, &iPtr) {
		t.Error("Failed to execute STA")
	}
	w := memory.ReadWord(100)
	if w != 0x5678 {
		t.Errorf("Expected loc 100 to contain 0x5678, got: %x", w)
	}
}
