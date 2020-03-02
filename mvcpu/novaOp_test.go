// novaOp_test.go

// Copyright (C) 2018, 2019 Steve Merrony

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

func TestADC(t *testing.T) {
	cpu := new(CPUT)
	var iPtr decodedInstrT
	var novaTwoAccMultOp novaTwoAccMultOpT
	iPtr.ix = instrADC
	novaTwoAccMultOp.acs = 0
	novaTwoAccMultOp.acd = 1
	iPtr.variant = novaTwoAccMultOp
	cpu.ac[0] = 0xffff
	cpu.ac[1] = 3
	cpu.carry = false
	if !novaOp(cpu, &iPtr) {
		t.Error("Failed to execute ADC")
	}
	if cpu.ac[1] != 3 {
		t.Errorf("Expected 3, got %d", cpu.ac[1])
	}

	cpu.ac[0] = 0
	cpu.ac[1] = 0
	cpu.carry = false
	if !novaOp(cpu, &iPtr) {
		t.Error("Failed to execute ADC")
	}
	if cpu.ac[1] != 0xffff {
		t.Errorf("Expected %d, got %d", 0xffff, cpu.ac[1])
	}
}

func TestADD(t *testing.T) {
	cpu := new(CPUT)
	var iPtr decodedInstrT
	var novaTwoAccMultOp novaTwoAccMultOpT
	iPtr.ix = instrADD

	// simple ADD
	novaTwoAccMultOp.acs = 1
	novaTwoAccMultOp.acd = 2
	cpu.ac[1] = 1
	cpu.ac[2] = 2
	iPtr.variant = novaTwoAccMultOp
	if !novaOp(cpu, &iPtr) {
		t.Error("Failed to execute MOV")
	}
	if cpu.ac[2] != 3 {
		t.Errorf("Expected 3, got %d", cpu.ac[2])
	}

	// simple ADD that should set CARRY
	novaTwoAccMultOp.acs = 1
	novaTwoAccMultOp.acd = 2
	cpu.ac[1] = 1
	cpu.ac[2] = 2
	cpu.carry = false
	iPtr.variant = novaTwoAccMultOp
	if !novaOp(cpu, &iPtr) {
		t.Error("Failed to execute MOV")
	}
	if cpu.ac[2] != 3 {
		t.Errorf("Expected 3, got %d", cpu.ac[2])
	}

	// simple ADD to self
	novaTwoAccMultOp.acs = 1
	novaTwoAccMultOp.acd = 1
	cpu.ac[1] = 1
	iPtr.variant = novaTwoAccMultOp
	if !novaOp(cpu, &iPtr) {
		t.Error("Failed to execute MOV")
	}
	if cpu.ac[1] != 2 {
		t.Errorf("Expected 2 got %d", cpu.ac[1])
	}

	// ADDR to self
	novaTwoAccMultOp.acs = 1
	novaTwoAccMultOp.acd = 1
	cpu.ac[1] = 1
	cpu.carry = false
	novaTwoAccMultOp.sh = 'R'
	iPtr.variant = novaTwoAccMultOp
	if !novaOp(cpu, &iPtr) {
		t.Error("Failed to execute MOV")
	}
	if cpu.ac[1] != 1 {
		t.Errorf("Expected 1 got %d", cpu.ac[1])
	}
	if cpu.carry {
		t.Error("Expected CARRY to be clear")
	}

	// ADDR to self with carry set
	novaTwoAccMultOp.acs = 1
	novaTwoAccMultOp.acd = 1
	cpu.ac[1] = 1
	cpu.carry = true
	novaTwoAccMultOp.sh = 'R'
	iPtr.variant = novaTwoAccMultOp
	if !novaOp(cpu, &iPtr) {
		t.Error("Failed to execute MOV")
	}
	if cpu.ac[1] != 0x8001 {
		t.Errorf("Expected %#x got %#x", 0x8001, cpu.ac[1])
	}
	if cpu.carry {
		t.Error("Expected CARRY to be clear")
	}
}

func TestCOM(t *testing.T) {
	cpu := new(CPUT)
	var iPtr decodedInstrT
	var novaTwoAccMultOp novaTwoAccMultOpT
	iPtr.ix = instrCOM
	novaTwoAccMultOp.acs = 0
	novaTwoAccMultOp.acd = 1
	iPtr.variant = novaTwoAccMultOp
	cpu.ac[0] = 0xffff
	cpu.ac[1] = 3
	cpu.carry = false
	if !novaOp(cpu, &iPtr) {
		t.Error("Failed to execute COM")
	}
	if cpu.ac[1] != 0 {
		t.Errorf("Expected 0, got %d", cpu.ac[1])
	}

	cpu.ac[0] = 0
	cpu.ac[1] = 0xffff
	cpu.carry = false
	if !novaOp(cpu, &iPtr) {
		t.Error("Failed to execute COM")
	}
	if cpu.ac[1] != 0xffff {
		t.Errorf("Expected %d, got %d", 0xffff, cpu.ac[1])
	}
}

func TestMOV(t *testing.T) {
	cpu := new(CPUT)
	var iPtr decodedInstrT
	var novaTwoAccMultOp novaTwoAccMultOpT
	iPtr.ix = instrMOV

	// simple MOV
	novaTwoAccMultOp.acs = 1
	novaTwoAccMultOp.acd = 2
	cpu.ac[1] = 1
	cpu.ac[2] = 2
	iPtr.variant = novaTwoAccMultOp
	if !novaOp(cpu, &iPtr) {
		t.Error("Failed to execute MOV")
	}
	if cpu.ac[2] != 1 {
		t.Errorf("Expected 1, got %d", cpu.ac[2])
	}

	// simple MOV to self
	novaTwoAccMultOp.acs = 1
	novaTwoAccMultOp.acd = 1
	cpu.ac[1] = 1
	iPtr.variant = novaTwoAccMultOp
	if !novaOp(cpu, &iPtr) {
		t.Error("Failed to execute MOV")
	}
	if cpu.ac[1] != 1 {
		t.Errorf("Expected 1, got %d", cpu.ac[1])
	}

	// MOVR to self, no carry
	novaTwoAccMultOp.acs = 1
	novaTwoAccMultOp.acd = 1
	novaTwoAccMultOp.sh = 'R'
	cpu.carry = false
	cpu.ac[1] = 1
	iPtr.variant = novaTwoAccMultOp
	if !novaOp(cpu, &iPtr) {
		t.Error("Failed to execute MOV")
	}
	if cpu.ac[1] != 0 {
		t.Errorf("Expected 0, got %d", cpu.ac[1])
	}
	if !cpu.carry {
		t.Error("Expected CARRY to be set")
	}

	// MOVL to self, no carry
	novaTwoAccMultOp.acs = 1
	novaTwoAccMultOp.acd = 1
	novaTwoAccMultOp.sh = 'L'
	cpu.carry = false
	cpu.ac[1] = 1
	iPtr.variant = novaTwoAccMultOp
	if !novaOp(cpu, &iPtr) {
		t.Error("Failed to execute MOV")
	}
	if cpu.ac[1] != 2 {
		t.Errorf("Expected 2, got %d", cpu.ac[1])
	}
	if cpu.carry {
		t.Error("Expected CARRY to be clear")
	}

	// MOVR to self, with carry
	novaTwoAccMultOp.acs = 1
	novaTwoAccMultOp.acd = 1
	novaTwoAccMultOp.sh = 'R'
	cpu.carry = true
	cpu.ac[1] = 1
	iPtr.variant = novaTwoAccMultOp
	if !novaOp(cpu, &iPtr) {
		t.Error("Failed to execute MOV")
	}
	if cpu.ac[1] != 0x8000 {
		t.Errorf("Expected %x, got %x", 0x8000, cpu.ac[1])
	}
	if !cpu.carry {
		t.Error("Expected CARRY to be set")
	}

	// MOVL to self, with carry
	novaTwoAccMultOp.acs = 1
	novaTwoAccMultOp.acd = 1
	novaTwoAccMultOp.sh = 'L'
	cpu.carry = true
	cpu.ac[1] = 1
	iPtr.variant = novaTwoAccMultOp
	if !novaOp(cpu, &iPtr) {
		t.Error("Failed to execute MOV")
	}
	if cpu.ac[1] != 3 {
		t.Errorf("Expected %x, got %x", 3, cpu.ac[1])
	}
	if cpu.carry {
		t.Error("Expected CARRY to be clear")
	}

	// MOVL to self, with carry clear, should set
	novaTwoAccMultOp.acs = 1
	novaTwoAccMultOp.acd = 1
	novaTwoAccMultOp.sh = 'L'
	cpu.carry = false
	cpu.ac[1] = 0xffff
	iPtr.variant = novaTwoAccMultOp
	if !novaOp(cpu, &iPtr) {
		t.Error("Failed to execute MOV")
	}
	if cpu.ac[1] != 0xfffe {
		t.Errorf("Expected %x, got %x", 0xfffe, cpu.ac[1])
	}
	if !cpu.carry {
		t.Error("Expected CARRY to be set")
	}

	// MOVL to self, with carry clear, should set
	novaTwoAccMultOp.acs = 1
	novaTwoAccMultOp.acd = 1
	novaTwoAccMultOp.sh = 'L'
	cpu.carry = false
	cpu.ac[1] = 0126356
	iPtr.variant = novaTwoAccMultOp
	if !novaOp(cpu, &iPtr) {
		t.Error("Failed to execute MOV")
	}
	if cpu.ac[1] != 054734 {
		t.Errorf("Expected %x, got %x", 054734, cpu.ac[1])
	}
	if !cpu.carry {
		t.Error("Expected CARRY to be set")
	}

	// MOVL to self, with carry clear, should set
	novaTwoAccMultOp.acs = 1
	novaTwoAccMultOp.acd = 1
	novaTwoAccMultOp.sh = 'L'
	cpu.carry = false
	cpu.ac[1] = 0xacf9
	iPtr.variant = novaTwoAccMultOp
	if !novaOp(cpu, &iPtr) {
		t.Error("Failed to execute MOV")
	}
	if cpu.ac[1] != 0x59f2 {
		t.Errorf("Expected %#x, got %x#", 0x59f2, cpu.ac[1])
	}
	if !cpu.carry {
		t.Error("Expected CARRY to be set")
	}

	// specific test for possibly-failing instruction...
	// MOV# 0,0,SZR # skip if AC0 == 0
	cpu.pc = 100
	cpu.ac[0] = 1
	novaTwoAccMultOp.acs = 0
	novaTwoAccMultOp.acd = 0
	novaTwoAccMultOp.nl = '#'
	novaTwoAccMultOp.sh = ' '
	novaTwoAccMultOp.skip = szrSkip
	iPtr.variant = novaTwoAccMultOp
	if !novaOp(cpu, &iPtr) {
		t.Error("Failed to execute MOV")
	}
	if cpu.pc != 101 {
		t.Errorf("Expected PC = 101. got PC = %d", cpu.pc)
	}
	if cpu.ac[0] != 1 {
		t.Error("AC0 changed!")
	}

	cpu.ac[0] = 0
	cpu.pc = 100
	if !novaOp(cpu, &iPtr) {
		t.Error("Failed to execute MOV")
	}
	if cpu.pc != 102 {
		t.Errorf("Expected PC = 102. got PC = %d", cpu.pc)
	}
	if cpu.ac[0] != 0 {
		t.Error("AC0 changed!")
	}
}

func TestNEG(t *testing.T) {
	cpu := new(CPUT)
	var iPtr decodedInstrT
	var novaTwoAccMultOp novaTwoAccMultOpT
	iPtr.ix = instrNEG
	novaTwoAccMultOp.acs = 0
	novaTwoAccMultOp.acd = 1
	iPtr.variant = novaTwoAccMultOp
	cpu.ac[0] = 0
	cpu.ac[1] = 0
	cpu.carry = false
	if !novaOp(cpu, &iPtr) {
		t.Error("Failed to execute NEG")
	}
	if cpu.ac[1] != 0 {
		t.Errorf("Expected 0, got %d", cpu.ac[1])
	}

	cpu.ac[0] = 0xffff
	cpu.ac[1] = 0
	cpu.carry = false
	if !novaOp(cpu, &iPtr) {
		t.Error("Failed to execute NEG")
	}
	if cpu.ac[1] != 1 {
		t.Errorf("Expected %d, got %d", 1, cpu.ac[1])
	}
}
