// floatHandling_test.go

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

package memory

import (
	"log"
	"testing"

	"github.com/SMerrony/dgemug/dg"
)

// func TestF64toSingle(t *testing.T) {
// 	dw := Float64toDGsingle(1.5)
// 	if dw != 0x4000_0001 {
// 		t.Errorf("Expected 0x180_0000, got %#x", dw)
// 	}
// }

func TestDGdoubleToFloat64(t *testing.T) {
	var dgd dg.QwordT
	dgd = 0
	f := DGdoubleToFloat64(dgd)
	if f != 0.0 {
		t.Errorf("Expected 0.0, got %f", f)
	}
}

func TestToAndFro(t *testing.T) {
	f1 := 0.0
	q := Float64toDGdouble(f1)
	f2 := DGdoubleToFloat64(q)
	if f1 != f2 {
		t.Errorf("Expected %f, got %f", f1, f2)
	}
	f1 = -1.5
	q = Float64toDGdouble(f1)
	log.Printf("q: %v", q)
	f2 = DGdoubleToFloat64(q)
	if f1 != f2 {
		t.Errorf("Expected %f, got %f", f1, f2)
	}
}
