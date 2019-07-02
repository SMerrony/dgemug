// narrowStack.go

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

package memory

import (
	"github.com/SMerrony/dgemug/dg"
	"github.com/SMerrony/dgemug/logging"
)

// Some Page Zero special locations for the Narrow Stack
const (
	NspLoc  = 040 // 32. Narrow Stack Pointer
	NfpLoc  = 041
	NslLoc  = 042
	NsfaLoc = 043
)

// NsPush - PUSH a word onto the Narrow Stack
func NsPush(seg dg.PhysAddrT, data dg.WordT, debugging bool) {
	// TODO segment handling
	// TODO overflow/underflow handling - either here or in instruction?
	ramMu.Lock()
	ram[NspLoc]++ // we allow this direct write to a fixed location for performance
	addr := dg.PhysAddrT(ram[NspLoc])
	ramMu.Unlock()
	WriteWord(addr, data)
	if debugging {
		logging.DebugPrint(logging.DebugLog, "... NsPush pushed %#o onto the Narrow Stack at location: %#o\n", data, addr)
	}
}

// NsPop - POP a word off the Narrow Stack
func NsPop(seg dg.PhysAddrT, debugging bool) dg.WordT {
	// TODO segment handling
	// TODO overflow/underflow handling - either here or in instruction?
	ramMu.RLock()
	addr := dg.PhysAddrT(ram[NspLoc])
	ramMu.RUnlock()
	data := ReadWord(addr)
	ramMu.Lock()
	ram[NspLoc]-- // we allow this direct write to a fixed location for performance
	ramMu.Unlock()
	if debugging {
		logging.DebugPrint(logging.DebugLog, "... NsPop  popped %#o off  the Narrow Stack at location: %#o\n", data, addr)
	}
	return data
}
