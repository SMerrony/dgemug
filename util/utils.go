// utils.go

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

// N.B. These funcs may be used very frequently, ensure test coverage even for trivial ones
//      and keep them efficient (review from time to time).

package util

import (
	"fmt"
	"math/bits"

	"github.com/SMerrony/dgemug"
)

// BoolToInt converts a bool to 1 or 0
func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// BoolToYN converts a bool to Y or N
func BoolToYN(b bool) byte {
	if b {
		return 'Y'
	}
	return 'N'
}

// BoolToOnOff converts a bool to "On" or "Off"
func BoolToOnOff(b bool) string {
	if b {
		return "On"
	}
	return "Off"
}

// BoolToOZ converts a boolean to a O(ne) or Z(ero) byte
func BoolToOZ(b bool) byte {
	if b {
		return 'O'
	}
	return 'Z'
}

// DwordGetLowerWord gets the DG-lower word (RHS) of a doubleword
// Called VERY often, hopefully inlined!
func DwordGetLowerWord(dwd dg.DwordT) dg.WordT {
	return dg.WordT(dwd) // & 0x0000ffff mask unneccessary
}

// DwordGetUpperWord gets the DG-higher word (LHS) of a doubleword
// Called VERY often, hopefully inlined!
func DwordGetUpperWord(dwd dg.DwordT) dg.WordT {
	return dg.WordT(dwd >> 16)
}

// DwordFromTwoWords - catenate two DG Words into a DG DoubleWord
func DwordFromTwoWords(hw dg.WordT, lw dg.WordT) dg.DwordT {
	return dg.DwordT(hw)<<16 | dg.DwordT(lw)
}

// QwordFromTwoDwords - catenate two DG DoubleWords into a DG QuadWord
func QwordFromTwoDwords(hdw dg.DwordT, ldw dg.DwordT) dg.QwordT {
	return dg.QwordT(hdw)<<32 | dg.QwordT(ldw)
}

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

// WordToBinStr - get a pretty-printable string of a word
func WordToBinStr(w dg.WordT) string {
	return fmt.Sprintf("%08b %08b", w>>8, w&0x0ff)
}

// SexWordToDword - sign-extend a DG word to a DG DoubleWord
func SexWordToDword(wd dg.WordT) dg.DwordT {
	var dwd dg.DwordT
	if TestWbit(wd, 0) {
		dwd = dg.DwordT(wd) | 0xffff0000
	} else {
		dwd = dg.DwordT(wd) & 0x0000ffff
	}
	return dwd
}

// SwapBytes - swap over the two bytes in a dg_word
func SwapBytes(wd dg.WordT) dg.WordT {
	// var res dg.WordT
	// res = (wd >> 8) | ((wd & 0x00ff) << 8)
	// return res
	return dg.WordT(bits.ReverseBytes16(uint16(wd)))
}

// GetSegment - return the segment number for the supplied address
func GetSegment(addr dg.PhysAddrT) int {
	return int(GetDwbits(dg.DwordT(addr), 1, 3))
}
