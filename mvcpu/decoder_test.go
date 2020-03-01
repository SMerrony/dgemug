// decoder_test.go

// Copyright (C) 2017,2019  Steve Merrony

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
)

func Test2bitImm(t *testing.T) {
	ttable := []struct {
		i dg.WordT
		o uint16
	}{
		{0, 1},
		{1, 2},
		{3, 4},
	}
	for _, tt := range ttable {
		res := decode2bitImm(tt.i)
		if res != tt.o {
			t.Errorf("Expected %d, got %d", res, tt.o)
		}
	}
}

func TestDecode8bitDisp(t *testing.T) {
	var db dg.ByteT
	var md int
	var res int16

	db = 7
	md = absoluteMode
	res = decode8bitDisp(db, md)
	if res != 7 {
		t.Error("Expected 7, got ", res)
	}

	db = 7
	md = pcMode
	res = decode8bitDisp(db, md)
	if res != 7 {
		t.Error("Expected 7, got ", res)
	}

	db = 0xff
	md = pcMode
	res = decode8bitDisp(db, md)
	if res != -1 {
		t.Error("Expected -1, got ", res)
	}
}

func TestDecode15bitDisp(t *testing.T) {
	var disp15 dg.WordT
	var m int
	var res int16

	disp15 = 300
	m = absoluteMode
	res = decode15bitDisp(disp15, m)
	if res != 300 {
		t.Error("Absolute Mode: Expected 300, got ", res)
	}

	disp15 = 300
	m = ac2Mode
	res = decode15bitDisp(disp15, m)
	if res != 300 {
		t.Error("AC2 Mode: Expected 300, got ", res)
	}

	disp15 = 0x7ed4
	m = ac2Mode
	res = decode15bitDisp(disp15, m)
	if res != -300 {
		t.Error("AC2 Mode: Expected -300, got ", res)
	}

	disp15 = 0x7ed4
	m = pcMode
	res = decode15bitDisp(disp15, m)
	if res != -300 {
		t.Error("PC Mode: Expected -300, got ", res)
	}
}

func TestDecode31bitDisp(t *testing.T) {
	var lWord, hWord dg.WordT
	var res, ex int32

	res = decode31bitDisp(hWord, lWord, absoluteMode)
	if res != 0 {
		t.Error("Absolute mode: expected 0, got ", res)
	}

	hWord, lWord = 0, 10
	res = decode31bitDisp(hWord, lWord, absoluteMode)
	if res != 10 {
		t.Error("Absolute mode: expected 10, got ", res)
	}

	hWord, lWord = 10, 11
	res = decode31bitDisp(hWord, lWord, absoluteMode)
	ex = 10<<16 + 11
	if res != ex {
		t.Errorf("Absolute mode: expected %d, got %d", ex, res)
	}

	hWord, lWord = 0, 0
	res = decode31bitDisp(hWord, lWord, pcMode)
	if res != 0 {
		t.Error("PC mode: expected 0, got ", res)
	}

	hWord, lWord = 0x7fff, 0xffff
	res = decode31bitDisp(hWord, lWord, pcMode)
	ex = -1
	if res != ex {
		t.Errorf("PC mode: expected %d, got %d", ex, res)
	}

}
