// mvemg project novaMath_test.go

// Copyright (C) 2018  Steve Merrony

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

func TestDIV(t *testing.T) {
	cpu := new(CPUT)
	var iPtr decodedInstrT
	iPtr.ix = instrDIV
	cpu.ac[0] = 0 // hi dividend
	cpu.ac[1] = 6 // lo dividend
	cpu.ac[2] = 2 // divisor
	if !novaMath(cpu, &iPtr) {
		t.Error("Failed to execute DIV")
	}
	if cpu.ac[1] != 3 || cpu.ac[0] != 0 || cpu.ac[2] != 2 {
		t.Errorf("Expected 3, 0, 2, got: %d, %d, %d",
			cpu.ac[1], cpu.ac[0], cpu.ac[2])
	}

	cpu.ac[0] = 0 // hi dividend
	cpu.ac[1] = 6 // lo dividend
	cpu.ac[2] = 4 // divisor
	if !novaMath(cpu, &iPtr) {
		t.Error("Failed to execute DIV")
	}
	if cpu.ac[1] != 1 || cpu.ac[0] != 2 || cpu.ac[2] != 4 {
		t.Errorf("Expected 1, 2, 4, got: %d, %d, %d",
			cpu.ac[1], cpu.ac[0], cpu.ac[2])
	}

	cpu.ac[0] = 0      // hi dividend
	cpu.ac[1] = 0xf000 // lo dividend
	cpu.ac[2] = 2      // divisor
	if !novaMath(cpu, &iPtr) {
		t.Error("Failed to execute DIV")
	}
	if cpu.ac[1] != 0x7800 || cpu.ac[0] != 0 || cpu.ac[2] != 2 {
		t.Errorf("Expected 30720, 0, 2, got: %d, %d, %d",
			cpu.ac[1], cpu.ac[0], cpu.ac[2])
	}

	cpu.ac[0] = 1      // hi dividend
	cpu.ac[1] = 0xf000 // lo dividend
	cpu.ac[2] = 2      // divisor
	if !novaMath(cpu, &iPtr) {
		t.Error("Failed to execute DIV")
	}
	if cpu.ac[1] != 0xf800 || cpu.ac[0] != 0 || cpu.ac[2] != 2 {
		t.Errorf("Expected 63488, 0, 2, got: %d, %d, %d",
			cpu.ac[1], cpu.ac[0], cpu.ac[2])
	}

	cpu.ac[0] = 0xf000 // hi dividends- SHOULD CAUSE EXCEPTION
	cpu.ac[1] = 0xf000 // lo dividend
	cpu.ac[2] = 512    // divisor
	if !novaMath(cpu, &iPtr) {
		t.Error("Failed to execute DIV")
	}
	if cpu.ac[1] != 61440 || cpu.ac[0] != 61440 || cpu.ac[2] != 512 {
		t.Errorf("Expected 61440, 61440, 512, got: %d, %d, %d",
			cpu.ac[1], cpu.ac[0], cpu.ac[2])
	}
}
