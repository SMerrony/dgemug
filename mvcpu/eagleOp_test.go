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
	var cpu MvCPUT
	cpuPtr := &cpu
	var iPtr decodedInstrT
	var twoAcc1Word twoAcc1WordT
	iPtr.ix = instrNADD
	twoAcc1Word.acs = 0
	twoAcc1Word.acd = 1
	// test neg + neg
	cpuPtr.ac[0] = 0xffff // -1
	cpuPtr.ac[1] = 0xffff // -1
	cpuPtr.carry = true
	iPtr.variant = twoAcc1Word
	if !eagleOp(cpuPtr, &iPtr) {
		t.Error("Failed to execute NADD")
	}
	if cpuPtr.ac[1] != 0xfffffffe { // sign-extended
		t.Errorf("Expected %x, got %x", 0xfffffffe, cpuPtr.ac[1])
	}
	if cpuPtr.carry {
		t.Error("Unexpected CARRY")
	}

	// test neg + pos
	cpuPtr.ac[0] = 0x0001 //
	cpuPtr.ac[1] = 0xffff // -1

	if !eagleOp(cpuPtr, &iPtr) {
		t.Error("Failed to execute NADD")
	}
	if cpuPtr.ac[1] != 0 {
		t.Errorf("Expected %x, got %x", 0, cpuPtr.ac[1])
	}
	if cpuPtr.carry {
		t.Error("Unexpected CARRY")
	}

	// test CARRY
	cpuPtr.ac[0] = maxPosS16
	cpuPtr.ac[1] = 10
	if !eagleOp(cpuPtr, &iPtr) {
		t.Error("Failed to execute NADD")
	}
	if !cpuPtr.carry {
		t.Error("Should have set CARRY")
	}
}

func TestNSUB(t *testing.T) {
	var cpu MvCPUT
	cpuPtr := &cpu
	var iPtr decodedInstrT
	var twoAcc1Word twoAcc1WordT
	iPtr.ix = instrNSUB
	twoAcc1Word.acs = 0
	twoAcc1Word.acd = 1
	iPtr.variant = twoAcc1Word
	// test neg - neg
	cpuPtr.ac[0] = 0xffff // -1
	cpuPtr.ac[1] = 0xffff // -1

	if !eagleOp(cpuPtr, &iPtr) {
		t.Error("Failed to execute NSUB")
	}
	if cpuPtr.ac[1] != 0 {
		t.Errorf("Expected %x, got %x", 0, cpuPtr.ac[1])
	}

	// test neg - pos
	cpuPtr.ac[0] = 0x0001 // 1
	cpuPtr.ac[1] = 0xffff // -1

	if !eagleOp(cpuPtr, &iPtr) {
		t.Error("Failed to execute NADD")
	}
	if cpuPtr.ac[1] != 0xfffffffe {
		t.Errorf("Expected %x, got %x", 0xfffffffe, cpuPtr.ac[1])
	}
}

func TestWADC(t *testing.T) {
	var cpu MvCPUT
	cpuPtr := &cpu
	var iPtr decodedInstrT
	var twoAcc1Word twoAcc1WordT
	iPtr.ix = instrWADC
	twoAcc1Word.acs = 1
	twoAcc1Word.acd = 1
	iPtr.variant = twoAcc1Word
	// test neg - neg
	cpuPtr.ac[0] = 0
	cpuPtr.ac[1] = 1
	if !eagleOp(cpuPtr, &iPtr) {
		t.Error("Failed to execute WADC")
	}
	if int32(cpuPtr.ac[1]) != -1 {
		t.Errorf("Expected %x, got %x", -1, cpuPtr.ac[1])
	}
}

func TestWADI(t *testing.T) {
	var cpu MvCPUT
	cpuPtr := &cpu
	var iPtr decodedInstrT
	var immOneAcc immOneAccT
	iPtr.ix = instrWADI
	immOneAcc.acd = 1
	immOneAcc.immU16 = 4
	iPtr.variant = immOneAcc

	cpuPtr.ac[0] = 0
	cpuPtr.ac[1] = 76
	if !eagleOp(cpuPtr, &iPtr) {
		t.Error("Failed to execute WADI")
	}
	if int32(cpuPtr.ac[1]) != 80 {
		t.Errorf("Expected %d, got %d", 80, cpuPtr.ac[1])
	}
	if cpuPtr.carry {
		t.Error("Unexpected CARRY")
	}

	cpuPtr.ac[0] = 0
	cpuPtr.ac[1] = maxPosS32 - 2
	immOneAcc.immU16 = 4
	if !eagleOp(cpuPtr, &iPtr) {
		t.Error("Failed to execute WADI")
	}
	// if int32(cpuPtr.ac[1]) != 80 {
	// 	t.Errorf("Expected %d, got %d", 80, cpuPtr.ac[1])
	// }
	if !cpuPtr.carry {
		t.Error("Expected CARRY")
	}
}

func TestWANDI(t *testing.T) {
	var cpu MvCPUT
	cpuPtr := &cpu
	var iPtr decodedInstrT
	var oneAccImmDwd3Word oneAccImmDwd3WordT
	iPtr.ix = instrWANDI
	oneAccImmDwd3Word.immDword = 0x7fffffff
	oneAccImmDwd3Word.acd = 0
	iPtr.variant = oneAccImmDwd3Word
	cpuPtr.ac[0] = 0x3171
	if !eagleOp(cpuPtr, &iPtr) {
		t.Error("Failed to execute WANDI")
	}
	if cpuPtr.ac[0] != 0x3171 {
		t.Errorf("Expected %x, got %x", 0x3171, cpuPtr.ac[0])
	}
	oneAccImmDwd3Word.immDword = 0x7fffffff
	oneAccImmDwd3Word.acd = 0
	iPtr.variant = oneAccImmDwd3Word
	cpuPtr.ac[0] = 0x20202020
	if !eagleOp(cpuPtr, &iPtr) {
		t.Error("Failed to execute WANDI")
	}
	if cpuPtr.ac[0] != 0x20202020 {
		t.Errorf("Expected %x, got %x", 0x20202020, cpuPtr.ac[0])
	}
}

func TestWINC(t *testing.T) {
	var cpu MvCPUT
	cpuPtr := &cpu
	var iPtr decodedInstrT
	var twoAcc1Word twoAcc1WordT
	iPtr.ix = instrWINC
	twoAcc1Word.acs = 1
	twoAcc1Word.acd = 1
	iPtr.variant = twoAcc1Word
	// test neg - neg
	cpuPtr.ac[0] = 0
	cpuPtr.ac[1] = 1
	if !eagleOp(cpuPtr, &iPtr) {
		t.Error("Failed to execute WADC")
	}
	if int32(cpuPtr.ac[1]) != 2 {
		t.Errorf("Expected %x, got %x", 2, cpuPtr.ac[1])
	}
	if cpuPtr.carry {
		t.Error("Unexpected CARRY")
	}

	cpuPtr.ac[1] = 0xffffffff
	if !eagleOp(cpuPtr, &iPtr) {
		t.Error("Failed to execute WADC")
	}
	if int32(cpuPtr.ac[1]) != 0 {
		t.Errorf("Expected %x, got %x", 0, cpuPtr.ac[1])
	}
	if !cpuPtr.carry {
		t.Error("Expected CARRY")
	}
}

func TestWLSH(t *testing.T) {
	var cpu MvCPUT
	cpuPtr := &cpu
	var iPtr decodedInstrT
	var twoAcc1Word twoAcc1WordT
	iPtr.ix = instrWLSH
	twoAcc1Word.acs = 1
	twoAcc1Word.acd = 2
	iPtr.variant = twoAcc1Word

	cpuPtr.ac[2] = 8
	cpuPtr.ac[1] = 0
	if !eagleOp(cpuPtr, &iPtr) {
		t.Error("Failed to execute WLSH")
	}
	if cpuPtr.ac[2] != 8 {
		t.Errorf("Expected 8 got %d", cpuPtr.ac[2])
	}

	cpuPtr.ac[2] = 8
	cpuPtr.ac[1] = 1
	if !eagleOp(cpuPtr, &iPtr) {
		t.Error("Failed to execute WLSH")
	}
	if cpuPtr.ac[2] != 16 {
		t.Errorf("Expected 16 got %d", cpuPtr.ac[2])
	}

	cpuPtr.ac[2] = 8
	cpuPtr.ac[1] = 0xff // -1
	if !eagleOp(cpuPtr, &iPtr) {
		t.Error("Failed to execute WLSH")
	}
	if cpuPtr.ac[2] != 4 {
		t.Errorf("Expected 4 got %d", cpuPtr.ac[2])
	}
}

func TestWLSHI(t *testing.T) {
	var cpu MvCPUT
	cpuPtr := &cpu
	var iPtr decodedInstrT
	var oneAccImm2Word oneAccImm2WordT
	iPtr.ix = instrWLSHI
	oneAccImm2Word.acd = 0
	oneAccImm2Word.immS16 = 8 // should shift 1 byte left
	iPtr.variant = oneAccImm2Word
	cpuPtr.ac[0] = 0x00001234
	if !eagleOp(cpuPtr, &iPtr) {
		t.Error("Failed to execute WLSHI")
	}
	if cpuPtr.ac[0] != 0x00123400 {
		t.Errorf("Expected %x, got %x", 0x00123400, cpuPtr.ac[0])
	}

	oneAccImm2Word.immS16 = -8 // should shift 1 byte right
	iPtr.variant = oneAccImm2Word
	cpuPtr.ac[0] = 0x00001234
	if !eagleOp(cpuPtr, &iPtr) {
		t.Error("Failed to execute WLSHI")
	}
	if cpuPtr.ac[0] != 0x00000012 {
		t.Errorf("Expected %x, got %x", 0x00000012, cpuPtr.ac[0])
	}
}

func TestWNADI(t *testing.T) {
	var cpu MvCPUT
	cpuPtr := &cpu
	var iPtr decodedInstrT
	var oneAccImm2Word oneAccImm2WordT
	iPtr.ix = instrWNADI
	oneAccImm2Word.acd = 0
	oneAccImm2Word.immS16 = -32
	iPtr.variant = oneAccImm2Word
	cpuPtr.ac[0] = 'x'

	if !eagleOp(cpuPtr, &iPtr) {
		t.Error("Failed to execute WNADI")
	}
	if cpuPtr.ac[0] != 'X' {
		t.Errorf("Expected %d, got %d", 'X', cpuPtr.ac[0])
	}
	if cpuPtr.carry {
		t.Error("Unexpected CARRY")
	}
}

func TestWNEG(t *testing.T) {
	var cpu MvCPUT
	cpuPtr := &cpu
	var iPtr decodedInstrT
	var twoAcc1Word twoAcc1WordT
	iPtr.ix = instrWNEG
	twoAcc1Word.acs = 0
	twoAcc1Word.acd = 1
	iPtr.variant = twoAcc1Word
	cpuPtr.ac[0] = 37
	// test cpnversion to negative
	if !eagleOp(cpuPtr, &iPtr) {
		t.Error("Failed to execute WNEG")
	}
	if cpuPtr.ac[1] != 0xffffffdb {
		t.Errorf("Expected 0xffffffdb, got %x", cpuPtr.ac[1])
	}
	if cpuPtr.carry {
		t.Error("Unexpected CARRY")
	}
	// convert back to test conversion from negative
	cpuPtr.ac[0] = cpuPtr.ac[1]
	if !eagleOp(cpuPtr, &iPtr) {
		t.Error("Failed to execute WNEG")
	}
	if cpuPtr.ac[1] != 37 {
		t.Errorf("Expected 37, got %d", cpuPtr.ac[1])
	}
	if cpuPtr.carry {
		t.Error("Unexpected CARRY")
	}
}

func TestWSBI(t *testing.T) {
	var cpu MvCPUT
	cpuPtr := &cpu
	var iPtr decodedInstrT
	var immOneAcc immOneAccT
	iPtr.ix = instrWSBI
	immOneAcc.acd = 1
	immOneAcc.immU16 = 4
	iPtr.variant = immOneAcc

	cpuPtr.ac[0] = 0
	cpuPtr.ac[1] = 76
	if !eagleOp(cpuPtr, &iPtr) {
		t.Error("Failed to execute WSBI")
	}
	if int32(cpuPtr.ac[1]) != 72 {
		t.Errorf("Expected %d, got %d", 80, cpuPtr.ac[1])
	}
	if cpuPtr.carry {
		t.Error("Unexpected CARRY")
	}

	// cpuPtr.ac[0] = 0
	// cpuPtr.ac[1] = uint32(minNegS32) + 2
	// immOneAcc.immU16 = 4
	// if !eagleOp(cpuPtr, &iPtr) {
	// 	t.Error("Failed to execute WSBI")
	// }
	// // if int32(cpuPtr.ac[1]) != 80 {
	// // 	t.Errorf("Expected %d, got %d", 80, cpuPtr.ac[1])
	// // }
	// if !cpuPtr.carry {
	// 	t.Error("Expected CARRY")
	// }
}

func TestZEX(t *testing.T) {
	var cpu MvCPUT
	cpuPtr := &cpu
	var iPtr decodedInstrT
	var twoAcc1Word twoAcc1WordT
	iPtr.ix = instrZEX
	cpuPtr.ac[0] = 0x12345678
	cpuPtr.ac[1] = 0
	twoAcc1Word.acs = 0
	twoAcc1Word.acd = 1
	iPtr.variant = twoAcc1Word
	if !eagleOp(cpuPtr, &iPtr) {
		t.Error("Failed to execute ZEX")
	}
	if cpuPtr.ac[1] != 0x00005678 {
		t.Errorf("Expected 0x5678, got %x", cpuPtr.ac[1])
	}
}
