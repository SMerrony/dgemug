// bmcdch_test.go

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
	"testing"

	"github.com/SMerrony/dgemug/dg"
)

func TestBmcdchReset(t *testing.T) {
	var wd dg.WordT
	bmcdchInit(false)
	wd = regs[iochanDefReg]
	if wd != ioccdr1 {
		t.Error("Got ", wd)
	}
}

func TestWriteReadMapSlot(t *testing.T) {
	var dwd1, dwd2 dg.DwordT
	dwd1 = 0x11223344
	BmcdchWriteSlot(17, dwd1)
	dwd2 = BmcdchReadSlot(17)
	if dwd2 != 0x11223344 {
		t.Error("Expected 0x11223344, got ", dwd2)
	}

}

func TestBmcMapAddr(t *testing.T) {
	var addr1, addr2, page dg.PhysAddrT
	BmcdchWriteSlot(0, 0)
	addr1 = 1
	addr2, page = getBmcMapAddr(addr1)
	if addr2 != 1 {
		t.Error("Expected 1, got ", addr2, page)
	}
	BmcdchWriteSlot(0, 3)
	addr1 = 1
	addr2, page = getBmcMapAddr(addr1)
	// 3 << 10 is 3072
	if addr2 != 3073 {
		t.Error("Expected 3073, got ", addr2, page)
	}
}
func TestDchMapAddr(t *testing.T) {
	var addr1, addr2, page dg.PhysAddrT
	BmcdchWriteSlot(0, 0)
	addr1 = 1
	addr2, page = getBmcMapAddr(addr1)
	if addr2 != 1 {
		t.Error("Expected 1, got ", addr2, page)
	}
	BmcdchWriteSlot(0, 3)
	addr1 = 1
	addr2, page = getDchMapAddr(addr1)
	// 3 << 10 is 3072
	if addr2 != 3073 {
		t.Error("Expected 3073, got ", addr2, page)
	}
}

func TestGetDchMode(t *testing.T) {
	bmcdchInit(false)
	icdr := BmcdchReadReg(iochanDefReg)
	if icdr != 1 {
		t.Error("Expected initial IOCDR == 1, got", icdr)
	}
	r := getDchMode()
	if r {
		t.Error("Unecpected return from getDchMode()")
	}
	BmcdchWriteReg(iochanDefReg, 2)
	r = getDchMode()
	if !r {
		t.Error("Unecpected return from getDchMode()")
	}
}
