// mvemg project eclipseOp_test.go

// Copyright (C) 2017,2019 Steve Merrony

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

import "testing"

import "github.com/SMerrony/dgemug/dg"

func TestADDI(t *testing.T) {
	cpu := new(CPUT)
	var iPtr decodedInstrT
	var oneAccImm2Word oneAccImm2WordT
	iPtr.ix = instrADDI
	oneAccImm2Word.immS16 = 3
	oneAccImm2Word.acd = 0
	cpu.ac[0] = 0xffff // -1
	iPtr.variant = oneAccImm2Word
	if !eclipseOp(cpu, &iPtr) {
		t.Error("Failed to execute ADDI")
	}
	if cpu.ac[0] != 2 {
		t.Errorf("Expected %x, got %x", 2, cpu.ac[0])
	}

	oneAccImm2Word.immS16 = -3
	oneAccImm2Word.acd = 0
	cpu.ac[0] = 0xffff // -1
	iPtr.variant = oneAccImm2Word
	if !eclipseOp(cpu, &iPtr) {
		t.Error("Failed to execute ADDI")
	}
	if cpu.ac[0] != 0xfffc {
		t.Errorf("Expected %x, got %x", -4, cpu.ac[0])
	}
}

func TestANDI(t *testing.T) {
	cpu := new(CPUT)
	var iPtr decodedInstrT
	var oneAccImmWd2Word oneAccImmWd2WordT
	iPtr.ix = instrANDI
	oneAccImmWd2Word.immWord = 3
	oneAccImmWd2Word.acd = 0
	cpu.ac[0] = 0
	iPtr.variant = oneAccImmWd2Word
	if !eclipseOp(cpu, &iPtr) {
		t.Error("Failed to execute ANDI")
	}
	if cpu.ac[0] != 0 {
		t.Errorf("Expected %x, got %x", 0, cpu.ac[0])
	}

	oneAccImmWd2Word.immWord = 0x5555
	oneAccImmWd2Word.acd = 0
	cpu.ac[0] = 0x00ff
	iPtr.variant = oneAccImmWd2Word
	if !eclipseOp(cpu, &iPtr) {
		t.Error("Failed to execute ANDI")
	}
	if cpu.ac[0] != 0x0055 {
		t.Errorf("Expected %x, got %x", 0x0055, cpu.ac[0])
	}
}

func TestHXL(t *testing.T) {
	cpu := new(CPUT)
	var iPtr decodedInstrT
	var immOneAcc immOneAccT
	iPtr.ix = instrHXL
	immOneAcc.acd = 0
	cpu.ac[0] = 0x0123
	immOneAcc.immU16 = 2
	iPtr.variant = immOneAcc
	expd := dg.DwordT(0x2300)
	if !eclipseOp(cpu, &iPtr) {
		t.Error("Failed to execute HXL")
	}
	if cpu.ac[0] != expd {
		t.Errorf("Expected %x, got %x", expd, cpu.ac[0])
	}

	cpu.ac[0] = 0x0123
	immOneAcc.immU16 = 4
	expd = dg.DwordT(0x0)
	iPtr.variant = immOneAcc
	if !eclipseOp(cpu, &iPtr) {
		t.Error("Failed to execute HXL")
	}
	if cpu.ac[0] != expd {
		t.Errorf("Expected %x, got %x", expd, cpu.ac[0])
	}
}
func TestHXR(t *testing.T) {
	cpu := new(CPUT)
	var iPtr decodedInstrT
	var immOneAcc immOneAccT
	iPtr.ix = instrHXR
	immOneAcc.acd = 0
	cpu.ac[0] = 0x0123
	immOneAcc.immU16 = 2
	iPtr.variant = immOneAcc
	expd := dg.DwordT(0x0001)
	if !eclipseOp(cpu, &iPtr) {
		t.Error("Failed to execute HXL")
	}
	if cpu.ac[0] != expd {
		t.Errorf("Expected %x, got %x", expd, cpu.ac[0])
	}

	cpu.ac[0] = 0x0123
	immOneAcc.immU16 = 4
	iPtr.variant = immOneAcc
	expd = dg.DwordT(0x0)
	if !eclipseOp(cpu, &iPtr) {
		t.Error("Failed to execute HXL")
	}
	if cpu.ac[0] != expd {
		t.Errorf("Expected %x, got %x", expd, cpu.ac[0])
	}
}

func TestSBI(t *testing.T) {
	cpu := new(CPUT)
	var iPtr decodedInstrT
	var immOneAcc immOneAccT
	iPtr.ix = instrSBI
	immOneAcc.immU16 = 3
	immOneAcc.acd = 0
	iPtr.variant = immOneAcc
	cpu.ac[0] = 0xffff // 65535
	if !eclipseOp(cpu, &iPtr) {
		t.Error("Failed to execute SBI")
	}
	if cpu.ac[0] != 65532 {
		t.Errorf("Expected %x, got %x", 65532, cpu.ac[0])
	}

	// test 'negative' wraparound
	immOneAcc.immU16 = 3
	immOneAcc.acd = 0
	iPtr.variant = immOneAcc
	cpu.ac[0] = 2
	if !eclipseOp(cpu, &iPtr) {
		t.Error("Failed to execute SBI")
	}
	if cpu.ac[0] != 65535 {
		t.Errorf("Expected %x, got %x", 65535, cpu.ac[0])
	}
}
