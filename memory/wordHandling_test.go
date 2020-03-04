// wordHandling_test.go

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
// THE SOFTWARE.package memory

package memory

import (
	"testing"

	"github.com/SMerrony/dgemug/dg"
)

func TestDWordFromTwoWords(t *testing.T) {
	var hi dg.WordT = 0x1122
	var lo dg.WordT = 0x3344
	r := DwordFromTwoWords(hi, lo)
	if r != 0x11223344 {
		t.Error("Expected 287454020, got ", r)
	}
}

func TestQWordFromTwoDwords(t *testing.T) {
	var hi dg.DwordT = 0x11223344
	var lo dg.DwordT = 0x55667788
	r := QwordFromTwoDwords(hi, lo)
	if r != 0x1122334455667788 {
		t.Error("Expected 0x1122334455667788, got ", r)
	}
}

func TestDWordGetLowerWord(t *testing.T) {
	var dwd dg.DwordT = 0x11223344
	r := DwordGetLowerWord(dwd)
	if r != 0x3344 {
		t.Error("Expected 0x3344, got ", r)
	}
}
func TestDWordGetUpperWord(t *testing.T) {
	var dwd dg.DwordT = 0x11223344
	r := DwordGetUpperWord(dwd)
	if r != 0x1122 {
		t.Error("Expected 0x1122, got ", r)
	}
}

func TestSexWordToDWord(t *testing.T) {
	var wd dg.WordT
	var dwd dg.DwordT

	wd = 79
	dwd = SexWordToDword(wd)
	if dwd != 79 {
		t.Error("Expected 79, got ", dwd)
	}

	wd = 0xfff4
	dwd = SexWordToDword(wd)
	if dwd != 0xfffffff4 {
		t.Error("Expected -12, got ", dwd)
	}

}

func TestSwapBytes(t *testing.T) {
	var wd, s dg.WordT

	wd = 0x1234
	s = SwapBytes(wd)
	if s != 0x3412 {
		t.Error("Expected 13330., got ", s)
	}

	wd = SwapBytes(s)
	if wd != 0x1234 {
		t.Error("Expected 4660., got ", wd)
	}
}
