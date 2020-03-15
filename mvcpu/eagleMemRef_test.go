// +build physical !virtual

// mvemg project eagleMemRef_test.go

// Copyright (C) 2017,2019,2020 Steve Merrony

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

	"github.com/SMerrony/dgemug/dg"
	"github.com/SMerrony/dgemug/memory"
)

func TestWBTZ(t *testing.T) {
	cpu := new(CPUT)
	var iPtr decodedInstrT
	var twoAcc1Word twoAcc1WordT
	iPtr.ix = instrWBTZ
	memory.MemInit(10000, false)

	// case where acs == acd
	twoAcc1Word.acs = 0
	twoAcc1Word.acd = 0
	iPtr.variant = twoAcc1Word
	var wordOffset dg.DwordT = 73 << 4
	var bitNum dg.DwordT = 3
	cpu.ac[0] = wordOffset | bitNum
	memory.WriteWord(73, 0xffff)
	if !eagleMemRef(cpu, &iPtr) {
		t.Error("Failed to execute WBTZ 1")
	}
	w := memory.ReadWord(73)
	if w != 0xefff {
		t.Errorf("Expected %x, got %x", 0xefff, w)
	}

	// case where acs != acd
	twoAcc1Word.acs = 1
	twoAcc1Word.acd = 0
	iPtr.variant = twoAcc1Word
	wordOffset = 33 << 4
	bitNum = 3
	cpu.ac[0] = wordOffset | bitNum
	cpu.ac[1] = 40
	memory.WriteWord(73, 0xffff)
	if !eagleMemRef(cpu, &iPtr) {
		t.Error("Failed to execute WBTZ 2")
	}
	w = memory.ReadWord(73)
	if w != 0xefff {
		t.Errorf("Expected %x, got %x", 0xefff, w)
	}

	// case where acs != acd and acs is indirect
	twoAcc1Word.acs = 1
	twoAcc1Word.acd = 0
	iPtr.variant = twoAcc1Word
	wordOffset = 33 << 4
	bitNum = 3
	cpu.ac[0] = wordOffset | bitNum
	// put an indirect address in ac1 pointing to 60
	cpu.ac[1] = 0x80000000 | 60
	// put 40 in location 60
	memory.WriteDWord(60, 40) // DWord!!!
	memory.WriteWord(73, 0xffff)
	if !eagleMemRef(cpu, &iPtr) {
		t.Error("Failed to execute WBTZ 3")
	}
	w = memory.ReadWord(73)
	if w != 0xefff {
		t.Errorf("Expected %x, got %x", 0xefff, w)
	}
}

func TestWBLM(t *testing.T) {
	cpu := new(CPUT)
	var iPtr decodedInstrT
	iPtr.ix = instrWBLM
	memory.MemInit(1000, false)

	// FORWARDS

	// put own location into each mem location
	var wdaddr dg.PhysAddrT
	for wdaddr = 0; wdaddr < 1000; wdaddr++ {
		memory.WriteWord(wdaddr, dg.WordT(wdaddr))
	}
	cpu.ac[0] = 77
	cpu.ac[1] = 5
	cpu.ac[2] = 50 // src
	cpu.ac[3] = 40 // dest
	if !eagleMemRef(cpu, &iPtr) {
		t.Error("Failed to execute forwards WBLM")
	}
	for wdaddr = 40; wdaddr < 45; wdaddr++ {
		if memory.ReadWord(wdaddr) != dg.WordT(wdaddr)+10 {
			t.Errorf("Expected %d, got %d", wdaddr+10, memory.ReadWord(wdaddr))
		}
	}
	if cpu.ac[0] != 77 {
		t.Error("Expected AC0 == 77")
	}
	if cpu.ac[1] != 0 {
		t.Error("Expected AC1 = 0")
	}
	if cpu.ac[2] != 55 {
		t.Error("Expected AC2 = 55")
	}
	if cpu.ac[3] != 45 {
		t.Error("Expected AC3 = 45")
	}

	// BACKWARDS

	// put own location into each mem location
	for wdaddr = 0; wdaddr < 1000; wdaddr++ {
		memory.WriteWord(wdaddr, dg.WordT(wdaddr))
	}
	cpu.ac[0] = 77
	cpu.ac[1] = 0xffff_fffb // -5
	cpu.ac[2] = 50          // src
	cpu.ac[3] = 40          // dest
	if !eagleMemRef(cpu, &iPtr) {
		t.Error("Failed to execute reverse WBLM")
	}
	for wdaddr = 40; wdaddr > 35; wdaddr-- {
		if memory.ReadWord(wdaddr) != dg.WordT(wdaddr)+10 {
			t.Errorf("Expected %d, got %d", wdaddr+10, memory.ReadWord(wdaddr))
		}
	}
	if cpu.ac[0] != 77 {
		t.Error("Expected AC0 == 77")
	}
	if cpu.ac[1] != 0 {
		t.Error("Expected AC1 = 0")
	}
	if cpu.ac[2] != 45 {
		t.Errorf("Expected AC2 = 45, got %d", cpu.ac[2])
	}
	if cpu.ac[3] != 35 {
		t.Error("Expected AC3 = 35")
	}
}

func TestWCMV(t *testing.T) {
	cpu := new(CPUT)
	var iPtr decodedInstrT
	iPtr.ix = instrWCMV
	memory.MemInit(1000, false)
	memory.WriteByte(100, false, 'A')
	memory.WriteByte(100, true, 'B')
	memory.WriteByte(101, false, 'C')
	memory.WriteByte(101, true, 'D')
	memory.WriteByte(102, false, 'E')
	memory.WriteByte(102, true, 'F')
	memory.WriteByte(103, false, 'G')

	// simple, word-aligned fwd move
	destNoBytes := 7
	srcNoBytes := 7
	destBytePtr := 200 << 1
	srcBytePtr := 100 << 1
	cpu.ac[0] = dg.DwordT(destNoBytes)
	cpu.ac[1] = dg.DwordT(srcNoBytes)
	cpu.ac[2] = dg.DwordT(destBytePtr)
	cpu.ac[3] = dg.DwordT(srcBytePtr)
	if !eagleMemRef(cpu, &iPtr) {
		t.Error("Failed to execute WCMV")
	}
	r := memory.ReadByte(200, false)
	if r != 'A' {
		t.Errorf("Expected 'A', got '%c'", r)
	}
	r = memory.ReadByte(200, true)
	if r != 'B' {
		t.Errorf("Expected 'B', got '%c'", r)
	}
	r = memory.ReadByte(202, true)
	if r != 'F' {
		t.Errorf("Expected 'F', got '%c'", r)
	}

	// non-word-aligned fwd move
	srcBytePtr++
	cpu.ac[0] = dg.DwordT(destNoBytes)
	cpu.ac[1] = dg.DwordT(srcNoBytes)
	cpu.ac[2] = dg.DwordT(destBytePtr)
	cpu.ac[3] = dg.DwordT(srcBytePtr)
	if !eagleMemRef(cpu, &iPtr) {
		t.Error("Failed to execute WCMV")
	}
	r = memory.ReadByte(200, false)
	if r != 'B' {
		t.Errorf("Expected 'B', got '%c'", r)
	}
	r = memory.ReadByte(200, true)
	if r != 'C' {
		t.Errorf("Expected 'C', got '%c'", r)
	}

	// src backwards
	srcBytePtr = 100 << 1
	destNoBytes = -7
	cpu.ac[0] = dg.DwordT(destNoBytes)
	cpu.ac[1] = dg.DwordT(srcNoBytes)
	cpu.ac[2] = dg.DwordT(destBytePtr)
	cpu.ac[3] = dg.DwordT(srcBytePtr)
	if !eagleMemRef(cpu, &iPtr) {
		t.Error("Failed to execute WCMV")
	}
	r = memory.ReadByte(200, false)
	if r != 'A' {
		t.Errorf("Expected 'A', got '%c'", r)
	}
	r = memory.ReadByte(199, true)
	if r != 'B' {
		t.Errorf("Expected 'B', got '%c'", r)
	}
	r = memory.ReadByte(199, false)
	if r != 'C' {
		t.Errorf("Expected 'C', got '%c'", r)
	}
}

func TestWLDB(t *testing.T) {
	cpu := new(CPUT)
	var iPtr decodedInstrT
	var twoAcc1Word twoAcc1WordT
	iPtr.ix = instrWLDB
	memory.MemInit(1000, false)
	memory.WriteByte(100, false, 'A')
	memory.WriteByte(100, true, 'B')
	twoAcc1Word.acs = 1
	twoAcc1Word.acd = 2
	cpu.ac[1] = 100 << 1
	iPtr.variant = twoAcc1Word
	if !eagleMemRef(cpu, &iPtr) {
		t.Error("Failed to execute WLDB")
	}
	if cpu.ac[2] != 'A' {
		t.Errorf("Expected %d, got %d", 'A', cpu.ac[2])
	}
	twoAcc1Word.acs = 1
	twoAcc1Word.acd = 2
	cpu.ac[1] = 100<<1 + 1
	iPtr.variant = twoAcc1Word
	if !eagleMemRef(cpu, &iPtr) {
		t.Error("Failed to execute WLDB")
	}
	if cpu.ac[2] != 'B' {
		t.Errorf("Expected %d, got %d", 'B', cpu.ac[2])
	}
}

func TestXLDB(t *testing.T) {
	cpu := new(CPUT)
	var iPtr decodedInstrT
	var oneAccMode2Word oneAccMode2WordT
	iPtr.ix = instrXLDB
	memory.MemInit(1000, false)
	memory.WriteByte(100, false, 'A')
	memory.WriteByte(100, true, 'B')
	oneAccMode2Word.acd = 2
	oneAccMode2Word.mode = absoluteMode
	cpu.ac[2] = 400 << 1 // should not be used in absolute mode
	oneAccMode2Word.disp16 = 100 << 1
	oneAccMode2Word.bitLow = false // get high byte 'A'
	iPtr.variant = oneAccMode2Word
	if !eagleMemRef(cpu, &iPtr) {
		t.Error("Failed to execute XLDB [1]")
	}
	if cpu.ac[2] != 'A' {
		t.Errorf("[1] Expected %d, got %d", 'A', cpu.ac[2])
	}

	oneAccMode2Word.acd = 2
	oneAccMode2Word.mode = absoluteMode
	cpu.ac[2] = 400 << 1 // should not be used in absolute mode
	oneAccMode2Word.disp16 = 100 << 1
	oneAccMode2Word.bitLow = true // get low byte 'B'
	iPtr.variant = oneAccMode2Word
	if !eagleMemRef(cpu, &iPtr) {
		t.Error("Failed to execute XLDB [2]")
	}
	if cpu.ac[2] != 'B' {
		t.Errorf("[2] Expected %d, got %d", 'B', cpu.ac[2])
	}

	oneAccMode2Word.acd = 3
	oneAccMode2Word.mode = ac3Mode
	cpu.ac[3] = 100
	oneAccMode2Word.disp16 = 0
	oneAccMode2Word.bitLow = true // get low byte 'B'
	iPtr.variant = oneAccMode2Word
	if !eagleMemRef(cpu, &iPtr) {
		t.Error("Failed to execute XLDB [3]")
	}
	if cpu.ac[3] != 'B' {
		t.Errorf("[3] Expected %d, got %d", 'B', cpu.ac[3])
	}
}

func TestXSTB(t *testing.T) {
	cpu := new(CPUT)
	var iPtr decodedInstrT
	var oneAccMode2Word oneAccMode2WordT
	iPtr.ix = instrXSTB
	memory.MemInit(10000, false)

	// test high (left) byte write
	memory.WriteWord(7, 0) // write 0 into Word at normal addr 7
	oneAccMode2Word.disp16 = 7
	oneAccMode2Word.mode = absoluteMode
	oneAccMode2Word.bitLow = false
	oneAccMode2Word.acd = 1
	iPtr.variant = oneAccMode2Word
	cpu.ac[1] = 0x11223344
	if !eagleMemRef(cpu, &iPtr) {
		t.Error("Failed to execute XSTB")
	}
	w := memory.ReadWord(7)
	if w != 0x4400 {
		t.Errorf("Expected %d, got %d", 0x4400, w)
	}

	// test low (right) byte write
	memory.WriteWord(7, 0) // write 0 into Word at normal addr 7
	oneAccMode2Word.disp16 = 7
	oneAccMode2Word.mode = absoluteMode
	oneAccMode2Word.bitLow = true
	oneAccMode2Word.acd = 1
	iPtr.variant = oneAccMode2Word
	cpu.ac[1] = 0x11223344
	if !eagleMemRef(cpu, &iPtr) {
		t.Error("Failed to execute XSTB")
	}
	w = memory.ReadWord(7)
	if w != 0x0044 {
		t.Errorf("Expected %d, got %d", 0x0044, w)
	}
}
