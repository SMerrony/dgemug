// wideStack_test.go

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

	"github.com/SMerrony/dgemug"
)

func TestAdvanceWSP(t *testing.T) {
	wsp1 := ReadDWord(WspLoc)
	AdvanceWSP(1)
	wsp2 := ReadDWord(WspLoc)
	if wsp2-wsp1 != 2 {
		t.Errorf("Expected %d, got %d", 2, wsp2-wsp1)
	}
}

func TestWsPushAndPop(t *testing.T) {
	WsPush(0, 1)
	wsp := dg.PhysAddrT(ReadDWord(WspLoc))
	dw := ReadDWord(wsp)
	if dw != 1 {
		t.Errorf("Expected WspLoc+1 to contain 1, contains %x", dw)
	}
	dw = WsPop(0)
	if dw != 1 {
		t.Errorf("Expected POP to produce 1, got %x", dw)
	}
}

func TestPopQWord(t *testing.T) {
	WsPush(0, 0x11112222)
	WsPush(0, 0x33334444)
	r := WsPopQWord(0)
	if r != 0x1111222233334444 {
		t.Errorf("Expected 0x1111222233334444, got %x", r)
	}
}
