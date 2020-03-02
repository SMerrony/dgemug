// mvemg project eagleOp_test.go

// Copyright (C) 2017,2019  Steve Merrony

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

func TestNADD(t *testing.T) {
	cpu := new(CPUT)
	var iPtr decodedInstrT
	var twoAcc1Word twoAcc1WordT
	iPtr.ix = instrNADD
	twoAcc1Word.acs = 0
	twoAcc1Word.acd = 1
	// test neg + neg
	cpu.ac[0] = 0xffff // -1
	cpu.ac[1] = 0xffff // -1
	cpu.carry = true
	iPtr.variant = twoAcc1Word
	if !eagleOp(cpu, &iPtr) {
		t.Error("Failed to execute NADD")
	}
	if cpu.ac[1] != 0xfffffffe { // sign-extended
		t.Errorf("Expected %x, got %x", 0xfffffffe, cpu.ac[1])
	}
	if cpu.carry {
		t.Error("Unexpected CARRY")
	}

	// test neg + pos
	cpu.ac[0] = 0x0001 //
	cpu.ac[1] = 0xffff // -1

	if !eagleOp(cpu, &iPtr) {
		t.Error("Failed to execute NADD")
	}
	if cpu.ac[1] != 0 {
		t.Errorf("Expected %x, got %x", 0, cpu.ac[1])
	}
	if cpu.carry {
		t.Error("Unexpected CARRY")
	}

	// test CARRY
	cpu.ac[0] = maxPosS16
	cpu.ac[1] = 10
	if !eagleOp(cpu, &iPtr) {
		t.Error("Failed to execute NADD")
	}
	if !cpu.carry {
		t.Error("Should have set CARRY")
	}
}

func TestNSUB(t *testing.T) {
	cpu := new(CPUT)
	var iPtr decodedInstrT
	var twoAcc1Word twoAcc1WordT
	iPtr.ix = instrNSUB
	twoAcc1Word.acs = 0
	twoAcc1Word.acd = 1
	iPtr.variant = twoAcc1Word
	// test neg - neg
	cpu.ac[0] = 0xffff // -1
	cpu.ac[1] = 0xffff // -1

	if !eagleOp(cpu, &iPtr) {
		t.Error("Failed to execute NSUB")
	}
	if cpu.ac[1] != 0 {
		t.Errorf("Expected %x, got %x", 0, cpu.ac[1])
	}

	// test neg - pos
	cpu.ac[0] = 0x0001 // 1
	cpu.ac[1] = 0xffff // -1

	if !eagleOp(cpu, &iPtr) {
		t.Error("Failed to execute NADD")
	}
	if cpu.ac[1] != 0xfffffffe {
		t.Errorf("Expected %x, got %x", 0xfffffffe, cpu.ac[1])
	}
}

func TestWADC(t *testing.T) {
	cpu := new(CPUT)
	var iPtr decodedInstrT
	var twoAcc1Word twoAcc1WordT
	iPtr.ix = instrWADC
	twoAcc1Word.acs = 1
	twoAcc1Word.acd = 1
	iPtr.variant = twoAcc1Word
	// test neg - neg
	cpu.ac[0] = 0
	cpu.ac[1] = 1
	if !eagleOp(cpu, &iPtr) {
		t.Error("Failed to execute WADC")
	}
	if int32(cpu.ac[1]) != -1 {
		t.Errorf("Expected %x, got %x", -1, cpu.ac[1])
	}
}

func TestWADI(t *testing.T) {
	cpu := new(CPUT)
	var iPtr decodedInstrT
	var immOneAcc immOneAccT
	iPtr.ix = instrWADI
	immOneAcc.acd = 1
	immOneAcc.immU16 = 4
	iPtr.variant = immOneAcc

	cpu.ac[0] = 0
	cpu.ac[1] = 76
	if !eagleOp(cpu, &iPtr) {
		t.Error("Failed to execute WADI")
	}
	if int32(cpu.ac[1]) != 80 {
		t.Errorf("Expected %d, got %d", 80, cpu.ac[1])
	}
	if cpu.carry {
		t.Error("Unexpected CARRY")
	}

	cpu.ac[0] = 0
	cpu.ac[1] = maxPosS32 - 2
	immOneAcc.immU16 = 4
	if !eagleOp(cpu, &iPtr) {
		t.Error("Failed to execute WADI")
	}
	// if int32(cpu.ac[1]) != 80 {
	// 	t.Errorf("Expected %d, got %d", 80, cpu.ac[1])
	// }
	if !cpu.carry {
		t.Error("Expected CARRY")
	}
}

func TestWANDI(t *testing.T) {
	cpu := new(CPUT)
	var iPtr decodedInstrT
	var oneAccImmDwd3Word oneAccImmDwd3WordT
	iPtr.ix = instrWANDI
	oneAccImmDwd3Word.immDword = 0x7fffffff
	oneAccImmDwd3Word.acd = 0
	iPtr.variant = oneAccImmDwd3Word
	cpu.ac[0] = 0x3171
	if !eagleOp(cpu, &iPtr) {
		t.Error("Failed to execute WANDI")
	}
	if cpu.ac[0] != 0x3171 {
		t.Errorf("Expected %x, got %x", 0x3171, cpu.ac[0])
	}
	oneAccImmDwd3Word.immDword = 0x7fffffff
	oneAccImmDwd3Word.acd = 0
	iPtr.variant = oneAccImmDwd3Word
	cpu.ac[0] = 0x20202020
	if !eagleOp(cpu, &iPtr) {
		t.Error("Failed to execute WANDI")
	}
	if cpu.ac[0] != 0x20202020 {
		t.Errorf("Expected %x, got %x", 0x20202020, cpu.ac[0])
	}
}

func TestWINC(t *testing.T) {
	cpu := new(CPUT)
	var iPtr decodedInstrT
	var twoAcc1Word twoAcc1WordT
	iPtr.ix = instrWINC
	twoAcc1Word.acs = 1
	twoAcc1Word.acd = 1
	iPtr.variant = twoAcc1Word
	// test neg - neg
	cpu.ac[0] = 0
	cpu.ac[1] = 1
	if !eagleOp(cpu, &iPtr) {
		t.Error("Failed to execute WADC")
	}
	if int32(cpu.ac[1]) != 2 {
		t.Errorf("Expected %x, got %x", 2, cpu.ac[1])
	}
	if cpu.carry {
		t.Error("Unexpected CARRY")
	}

	cpu.ac[1] = 0xffffffff
	if !eagleOp(cpu, &iPtr) {
		t.Error("Failed to execute WADC")
	}
	if int32(cpu.ac[1]) != 0 {
		t.Errorf("Expected %x, got %x", 0, cpu.ac[1])
	}
	if !cpu.carry {
		t.Error("Expected CARRY")
	}
}

func TestWLSH(t *testing.T) {
	cpu := new(CPUT)
	var iPtr decodedInstrT
	var twoAcc1Word twoAcc1WordT
	iPtr.ix = instrWLSH
	twoAcc1Word.acs = 1
	twoAcc1Word.acd = 2
	iPtr.variant = twoAcc1Word

	cpu.ac[2] = 8
	cpu.ac[1] = 0
	if !eagleOp(cpu, &iPtr) {
		t.Error("Failed to execute WLSH")
	}
	if cpu.ac[2] != 8 {
		t.Errorf("Expected 8 got %d", cpu.ac[2])
	}

	cpu.ac[2] = 8
	cpu.ac[1] = 1
	if !eagleOp(cpu, &iPtr) {
		t.Error("Failed to execute WLSH")
	}
	if cpu.ac[2] != 16 {
		t.Errorf("Expected 16 got %d", cpu.ac[2])
	}

	cpu.ac[2] = 8
	cpu.ac[1] = 0xff // -1
	if !eagleOp(cpu, &iPtr) {
		t.Error("Failed to execute WLSH")
	}
	if cpu.ac[2] != 4 {
		t.Errorf("Expected 4 got %d", cpu.ac[2])
	}
}

func TestWLSHI(t *testing.T) {
	cpu := new(CPUT)
	var iPtr decodedInstrT
	var oneAccImm2Word oneAccImm2WordT
	iPtr.ix = instrWLSHI
	oneAccImm2Word.acd = 0
	oneAccImm2Word.immS16 = 8 // should shift 1 byte left
	iPtr.variant = oneAccImm2Word
	cpu.ac[0] = 0x00001234
	if !eagleOp(cpu, &iPtr) {
		t.Error("Failed to execute WLSHI")
	}
	if cpu.ac[0] != 0x00123400 {
		t.Errorf("Expected %x, got %x", 0x00123400, cpu.ac[0])
	}

	oneAccImm2Word.immS16 = -8 // should shift 1 byte right
	iPtr.variant = oneAccImm2Word
	cpu.ac[0] = 0x00001234
	if !eagleOp(cpu, &iPtr) {
		t.Error("Failed to execute WLSHI")
	}
	if cpu.ac[0] != 0x00000012 {
		t.Errorf("Expected %x, got %x", 0x00000012, cpu.ac[0])
	}
}

func TestWNADI(t *testing.T) {
	cpu := new(CPUT)
	var iPtr decodedInstrT
	var oneAccImm2Word oneAccImm2WordT
	iPtr.ix = instrWNADI
	oneAccImm2Word.acd = 0
	oneAccImm2Word.immS16 = -32
	iPtr.variant = oneAccImm2Word
	cpu.ac[0] = 'x'

	if !eagleOp(cpu, &iPtr) {
		t.Error("Failed to execute WNADI")
	}
	if cpu.ac[0] != 'X' {
		t.Errorf("Expected %d, got %d", 'X', cpu.ac[0])
	}
	if cpu.carry {
		t.Error("Unexpected CARRY")
	}
}

func TestWNEG(t *testing.T) {
	cpu := new(CPUT)
	var iPtr decodedInstrT
	var twoAcc1Word twoAcc1WordT
	iPtr.ix = instrWNEG
	twoAcc1Word.acs = 0
	twoAcc1Word.acd = 1
	iPtr.variant = twoAcc1Word
	cpu.ac[0] = 37
	// test cpnversion to negative
	if !eagleOp(cpu, &iPtr) {
		t.Error("Failed to execute WNEG")
	}
	if cpu.ac[1] != 0xffffffdb {
		t.Errorf("Expected 0xffffffdb, got %x", cpu.ac[1])
	}
	if cpu.carry {
		t.Error("Unexpected CARRY")
	}
	// convert back to test conversion from negative
	cpu.ac[0] = cpu.ac[1]
	if !eagleOp(cpu, &iPtr) {
		t.Error("Failed to execute WNEG")
	}
	if cpu.ac[1] != 37 {
		t.Errorf("Expected 37, got %d", cpu.ac[1])
	}
	if cpu.carry {
		t.Error("Unexpected CARRY")
	}
}

func TestWSBI(t *testing.T) {
	cpu := new(CPUT)
	var iPtr decodedInstrT
	var immOneAcc immOneAccT
	iPtr.ix = instrWSBI
	immOneAcc.acd = 1
	immOneAcc.immU16 = 4
	iPtr.variant = immOneAcc

	cpu.ac[0] = 0
	cpu.ac[1] = 76
	if !eagleOp(cpu, &iPtr) {
		t.Error("Failed to execute WSBI")
	}
	if int32(cpu.ac[1]) != 72 {
		t.Errorf("Expected %d, got %d", 80, cpu.ac[1])
	}
	if cpu.carry {
		t.Error("Unexpected CARRY")
	}

	// cpu.ac[0] = 0
	// cpu.ac[1] = uint32(minNegS32) + 2
	// immOneAcc.immU16 = 4
	// if !eagleOp(cpu, &iPtr) {
	// 	t.Error("Failed to execute WSBI")
	// }
	// // if int32(cpu.ac[1]) != 80 {
	// // 	t.Errorf("Expected %d, got %d", 80, cpu.ac[1])
	// // }
	// if !cpu.carry {
	// 	t.Error("Expected CARRY")
	// }
}

func TestZEX(t *testing.T) {
	cpu := new(CPUT)
	var iPtr decodedInstrT
	var twoAcc1Word twoAcc1WordT
	iPtr.ix = instrZEX
	cpu.ac[0] = 0x12345678
	cpu.ac[1] = 0
	twoAcc1Word.acs = 0
	twoAcc1Word.acd = 1
	iPtr.variant = twoAcc1Word
	if !eagleOp(cpu, &iPtr) {
		t.Error("Failed to execute ZEX")
	}
	if cpu.ac[1] != 0x00005678 {
		t.Errorf("Expected 0x5678, got %x", cpu.ac[1])
	}
}
