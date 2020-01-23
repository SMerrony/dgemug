// bitHandling.go

// Copyright (C) 2018,2020 Steve Merrony

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

// N.B. These funcs may be used very frequently, ensure test coverage even for trivial ones
//      and keep them efficient (review from time to time).

package memory

import "github.com/SMerrony/dgemug/dg"

// GetWbits - in the DG world, the first (leftmost) bit is numbered zero...
// extract nbits from value starting at leftBit
func GetWbits(value dg.WordT, leftBit uint, nbits uint) dg.WordT {
	var mask dg.WordT
	if leftBit >= 16 {
		return 0
	}
	value >>= 16 - (leftBit + nbits)
	mask = (1 << nbits) - 1
	return value & mask
}

// FlipWbit - flips a single bit in a Word using DG numbering
func FlipWbit(word *dg.WordT, bitNum uint) {
	*word = *word ^ (1 << (15 - bitNum))
}

// SetWbit sets a single bit in a DG word
func SetWbit(word *dg.WordT, bitNum uint) {
	*word = *word | (1 << (15 - bitNum))
}

// SetDwbit - sets a bit in a Double-Word
func SetDwbit(dw *dg.DwordT, bitNum uint) {
	*dw = *dw | (1 << (31 - bitNum))
}

// SetQwbit - sets a bit in a Quad-Word
func SetQwbit(qw *dg.QwordT, bitNum uint) {
	*qw = *qw | (1 << (63 - bitNum))
}

// ClearWbit clears a single bit in a DG word
func ClearWbit(word *dg.WordT, bitNum uint) {
	*word = *word &^ (1 << (15 - bitNum))
}

// GetDwbits - in the DG world, the first (leftmost) bit is numbered zero...
// extract nbits from value starting at leftBit
func GetDwbits(value dg.DwordT, leftBit uint, nbits uint) dg.DwordT {
	var mask dg.DwordT
	if leftBit >= 32 {
		return 0
	}
	value >>= 32 - (leftBit + nbits)
	mask = (1 << nbits) - 1
	return value & mask
}

// GetQwbits - in the DG world, the first (leftmost) bit is numbered zero...
// extract nbits from value starting at leftBit
func GetQwbits(value dg.QwordT, leftBit uint, nbits uint) dg.QwordT {
	var mask dg.QwordT
	if leftBit >= 64 {
		return 0
	}
	value >>= 64 - (leftBit + nbits)
	mask = (1 << nbits) - 1
	return value & mask
}

// TestWbit - does word w have bit b set?
func TestWbit(w dg.WordT, b int) bool {
	return (w & (1 << (15 - uint8(b)))) != 0
}

// TestDwbit - does dword dw have bit b set?
func TestDwbit(dw dg.DwordT, b int) bool {
	return ((dw & (1 << (31 - uint8(b)))) != 0)
}

// TestQwbit - does qword qw have bit b set?
func TestQwbit(qw dg.QwordT, b int) bool {
	return ((qw & (1 << (63 - uint8(b)))) != 0)
}
