// eclipseStack_test.go

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

import "testing"
import "github.com/SMerrony/dgemug/memory"

func TestSave(t *testing.T) {
	cpuPtr := new(MvCPUT)
	var iPtr decodedInstrT
	var unique2Word unique2WordT
	memory.MemInit(10000, false)
	// set nsb & nfp to 256
	memory.WriteWord(memory.NspLoc, 256)
	memory.WriteWord(memory.NfpLoc, 256)
	iPtr.ix = instrSAVE
	unique2Word.immU16 = 0 // "SAVE 0"
	cpuPtr.ac[0] = 0
	cpuPtr.ac[1] = 1
	cpuPtr.ac[2] = 2
	cpuPtr.ac[3] = 3
	iPtr.variant = unique2Word
	if !eclipseStack(cpuPtr, &iPtr) {
		t.Error("Failed to execute SAVE")
	}
	newSP := memory.ReadWord(memory.NspLoc)
	if newSP != 261 {
		t.Errorf("Expected NSP to be 261, got %d", newSP)
	}
	ac2 := memory.ReadWord(259)
	if ac2 != 2 {
		t.Errorf("Expected 2 from ac2 in NSP-3, got %d", ac2)
	}
	if cpuPtr.ac[3] != 256+5 {
		t.Errorf("Expected AC3 to contain 261, got %d", cpuPtr.ac[3])
	}
	newFP := memory.ReadWord(memory.NfpLoc)
	if newFP != 256+5 {
		t.Errorf("Expected NFP to be 261, got %d", newFP)
	}
	t.Logf("SAVE 0 OK\n")

	memory.WriteWord(memory.NspLoc, 256)
	memory.WriteWord(memory.NfpLoc, 256)
	iPtr.ix = instrSAVE
	unique2Word.immU16 = 5 // "SAVE 5"
	cpuPtr.ac[0] = 0
	cpuPtr.ac[1] = 1
	cpuPtr.ac[2] = 2
	cpuPtr.ac[3] = 3
	iPtr.variant = unique2Word
	if !eclipseStack(cpuPtr, &iPtr) {
		t.Error("Failed to execute SAVE")
	}
	newSP = memory.ReadWord(memory.NspLoc)
	if newSP != 266 {
		t.Errorf("Expected NSP to be 266, got %d", newSP)
	}
	ac2 = memory.ReadWord(259)
	if ac2 != 2 {
		t.Errorf("Expected 2 from ac2 in NSP-3, got %d", ac2)
	}
	if cpuPtr.ac[3] != 256+5 {
		t.Errorf("Expected AC3 to contain 261, got %d", cpuPtr.ac[3])
	}
	newFP = memory.ReadWord(memory.NfpLoc)
	if newFP != 256+5 {
		t.Errorf("Expected NFP to be 261, got %d", newFP)
	}
}
