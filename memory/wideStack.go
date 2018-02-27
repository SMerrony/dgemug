// wideStack.go

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
	"mvemg/logging"

	"github.com/SMerrony/dgemug"
)

// Some Page Zero special locations for the Wide Stack...
const (
	WfpLoc  = 020 // 16.
	WspLoc  = 022 // 18.
	WslLoc  = 024 // 20.
	WsbLoc  = 026 // 22.
	WpfhLoc = 030
	CbpLoc  = 032
)

// WsPush - PUSH a doubleword onto the Wide Stack
func WsPush(seg dg.PhysAddrT, data dg.DwordT) {
	// TODO segment handling
	// TODO overflow/underflow handling - either here or in instruction?
	wsp := ReadDWord(WspLoc) + 2
	WriteDWord(WspLoc, wsp)
	WriteDWord(dg.PhysAddrT(wsp), data)
	logging.DebugPrint(logging.DebugLog, "... memory.WsPush pushed %8d onto the Wide Stack at location: %d\n", data, wsp)
}

// WsPop - POP a doubleword off the Wide Stack
func WsPop(seg dg.PhysAddrT) dg.DwordT {
	// TODO segment handling
	// TODO overflow/underflow handling - either here or in instruction?
	wsp := ReadDWord(WspLoc)
	dword := ReadDWord(dg.PhysAddrT(wsp))
	WriteDWord(WspLoc, wsp-2)
	logging.DebugPrint(logging.DebugLog, "... memory.WsPop  popped %8d off  the Wide Stack at location: %d\n", dword, wsp)
	return dword
}

// WsPopQWord - POP a Quad-word off the Wide Stack
func WsPopQWord(seg dg.PhysAddrT) dg.QwordT {
	// TODO segment handling
	var qw dg.QwordT
	rhDWord := WsPop(seg)
	lhDWord := WsPop(seg)
	qw = dg.QwordT(lhDWord)<<32 | dg.QwordT(rhDWord)
	return qw
}

// AdvanceWSP increases the WSP by the given amount of DWords
func AdvanceWSP(dwdCnt uint) {
	wsp := ReadDWord(WspLoc) + dg.DwordT(dwdCnt*2)
	WriteDWord(WspLoc, wsp)
	logging.DebugPrint(logging.DebugLog, "... WSP advanced by %d DWords to %d\n", dwdCnt, wsp)
}
