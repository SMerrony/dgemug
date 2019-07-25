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

import (
	"fmt"
	"testing"
)

func TestIRQmasking(t *testing.T) {
	var bus BusT
	bus.BusInit()
	var testDevMap = DeviceMapT{
		1: {"TEST", 2, true, false},
		2: {"TEST2", 15, true, false},
	}
	bus.AddDevice(testDevMap, 1, true)
	bus.AddDevice(testDevMap, 2, true)

	fmt.Printf("Device Map:\n%s\n", bus.GetPrintableDevList())

	if bus.IsDevMasked(1) {
		t.Error("Device 1 should not be masked")
	}
	bus.SetIrqMask(1)
	if bus.IsDevMasked(1) {
		t.Error("Device 1 should not be masked")
	}
	if !bus.IsDevMasked(2) {
		t.Error("Device 2 should be masked")
	}
	bus.SetIrqMask(0x7000)
	if !bus.IsDevMasked(1) {
		t.Error("Device 1 should be masked")
	}
}
