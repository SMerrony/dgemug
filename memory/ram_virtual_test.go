// +build virtual !physical

// TESTS FOR REPRESENTATION OF VIRTUAL MEMORY USED IN THE OS-LEVEL EMULATOR(S)

// Copyright ©2020  Steve Merrony

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
// THE SOFTWARE.package memory

package memory

import (
	"testing"

	"github.com/SMerrony/dgemug/dg"
)

func TestReadWriteWord(t *testing.T) {
	MemInit()
	wd := ReadWord(0x7000_0001)
	if wd != 0 {
		t.Errorf("Expected zero")
	}

	WriteWord(0x7000_0001, 5)
	wd = ReadWord(0x7000_0001)
	if wd != 5 {
		t.Errorf("Expected 5, got %#x", wd)
	}
}

func TestReadWriteDWord(t *testing.T) {
	MemInit()
	var (
		addr dg.PhysAddrT
		dwd  dg.DwordT = 0x1234_5678
	)
	addr = 0x7000_0010 // normal case
	WriteDWord(addr, dwd)
	res := ReadDWord(addr)
	if res != dwd {
		t.Errorf("Expected %#x, got %#x", dwd, res)
	}

	// test across page boundary
	MapPage(ring7page0+1, true)
	addr = 0x7000_03ff // last word of page 0
	WriteDWord(addr, dwd)
	res = ReadDWord(addr)
	if res != dwd {
		t.Errorf("Expected %#x, got %#x", dwd, res)
	}
}
