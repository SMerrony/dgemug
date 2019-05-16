// bus_test.go

// Copyright (C) 2018 Steve Merrony

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

package devices

// import (
// 	"testing"
// )

// func TestIRQmasking(t *testing.T) {
// 	BusInit()
// 	BusAddDevice(1, "TEST", 2, true, true, false)   // PMB 01000000000000
// 	BusAddDevice(2, "TEST2", 15, true, true, false) // PMB 00000000000001
// 	if BusIsDevMasked(1) {
// 		t.Error("Device 1 should not be masked")
// 	}
// 	BusSetIrqMask(1)
// 	if BusIsDevMasked(1) {
// 		t.Error("Device 1 should not be masked")
// 	}
// 	if !BusIsDevMasked(2) {
// 		t.Error("Device 2 should be masked")
// 	}
// 	BusSetIrqMask(0x7000)
// 	if !BusIsDevMasked(1) {
// 		t.Error("Device 1 should be masked")
// 	}
// }
