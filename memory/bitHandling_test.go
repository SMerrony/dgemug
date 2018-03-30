// bitHandling_test.go

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

package memory

import (
	"testing"

	"github.com/SMerrony/dgemug/dg"
)

func TestGetWbits(t *testing.T) {
	var w dg.WordT = 0xb38f
	r := GetWbits(w, 5, 3)
	if r != 3 {
		t.Error("Expected 3, got ", r)
	}
	w = 0xb38f
	r = GetWbits(w, 2, 2)
	if r != 3 {
		t.Error("Expected 3, got ", r)
	}
	w = 0xbeef
	r = GetWbits(w, 8, 8)
	if r != 0xef {
		t.Error("Expected 0xEF, got ", r)
	}
}

func TestFlipWbit(t *testing.T) {
	var w dg.WordT = 0x0010
	FlipWbit(&w, 0)
	if w != 0x8010 {
		t.Errorf("Expected 0x8010, got %x", w)
	}
	FlipWbit(&w, 0)
	if w != 0x0010 {
		t.Errorf("Expected 0x0010, got %x", w)
	}
}

func TestSetWbit(t *testing.T) {
	var w dg.WordT = 0x0001
	SetWbit(&w, 1)
	if w != 0x4001 {
		t.Error("Expected 16385, got ", w)
	}
	// repeat - should have no effect
	SetWbit(&w, 1)
	if w != 0x4001 {
		t.Error("Expected 16385, got ", w)
	}
}

func TestSetQWbit(t *testing.T) {
	var w dg.QwordT = 0
	SetQwbit(&w, 63)
	if w != 1 {
		t.Errorf("Expected 0x01, got %x", w)
	}
	// repeat - should have no effect
	SetQwbit(&w, 63)
	if w != 1 {
		t.Errorf("Expected 0x01, got %x", w)
	}
}
func TestClearWbit(t *testing.T) {
	var w dg.WordT = 0x4001
	ClearWbit(&w, 1)
	if w != 1 {
		t.Error("Expected 1, got ", w)
	}
	// repeat - should have no effect
	ClearWbit(&w, 1)
	if w != 1 {
		t.Error("Expected 1, got ", w)
	}
}
func TestGetDWbits(t *testing.T) {
	var w dg.DwordT = 0x11112222
	r := GetDwbits(w, 15, 2)
	if r != 2 {
		t.Error("Expected 2, got ", r)
	}

}

func TestGetQWbits(t *testing.T) {
	var q dg.QwordT = 0x1111222233334444
	r := GetQwbits(q, 12, 8)
	if r != 0x12 {
		t.Errorf("Expected 0x12, got %x", r)
	}
}

func TestTestWbit(t *testing.T) {
	var wd dg.WordT = 0x4000
	r := TestWbit(wd, 1)
	if !r {
		t.Error("Expected true")
	}
	r = TestWbit(wd, 2)
	if r {
		t.Error("Expected false")
	}
}
func TestTestDWbit(t *testing.T) {
	var wd dg.DwordT = 0x40000000
	r := TestDwbit(wd, 1)
	if !r {
		t.Error("Expected true")
	}
	r = TestDwbit(wd, 2)
	if r {
		t.Error("Expected false")
	}
}
func TestTestQWbit(t *testing.T) {
	var wd dg.QwordT = 0x4000000000000000
	r := TestQwbit(wd, 1)
	if !r {
		t.Error("Expected true")
	}
	r = TestQwbit(wd, 2)
	if r {
		t.Error("Expected false")
	}
}
