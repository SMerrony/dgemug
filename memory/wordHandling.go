// wordHandling.go

// Copyright ©2020 Steve Merrony

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

import (
	"fmt"

	"github.com/SMerrony/dgemug/dg"
)

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

// QwordGetLowerDword gets the DG-lower (RHS) of a quadword
func QwordGetLowerDword(qwd dg.QwordT) dg.DwordT {
	return dg.DwordT(qwd)
}

// QwordGetUpperDword gets the DG-upper (LHS) of a quadword
func QwordGetUpperDword(qwd dg.QwordT) dg.DwordT {
	return dg.DwordT(qwd >> 32)
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
func SwapBytes(wd dg.WordT) (res dg.WordT) {
	res = (wd >> 8) | ((wd & 0x00ff) << 8)
	return res
}

// WordsFromBytes converts a slice of (Go) bytes into a slice of DG Words
func WordsFromBytes(ba []byte) (wa []dg.WordT) {
	wordsLen := len(ba) / 2
	wa = make([]dg.WordT, wordsLen)
	for wordAddr := 0; wordAddr < wordsLen; wordAddr++ {
		wa[wordAddr] = dg.WordT(ba[wordAddr*2])<<8 | dg.WordT(ba[wordAddr*2+1])
	}
	if len(ba)%2 == 1 {
		wa = append(wa, dg.WordT(ba[len(ba)-1])<<8)
	}
	return wa
}
