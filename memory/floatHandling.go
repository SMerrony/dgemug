// floatHandling.go

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
	"math"

	"github.com/SMerrony/dgemug/dg"
)

// Float64toDGsingle converts a standard Go float64 into the DG Single Precision format
// func Float64toDGsingle(f float64) (s dg.DwordT) {
// 	// // sign - first (LH) bit
// 	// if f < 0.0 {
// 	// 	SetDwbit(&s, 0)
// 	// 	f *= -1.0
// 	// 	log.Println("Negative")
// 	// }
// 	// // exponent - next 7 bits in excess-64
// 	// exp := int8(math.Floor(math.Log10(f)))
// 	// log.Printf("Exp: %d", exp)
// 	// s |= dg.DwordT((exp+64)&0x7f) << 24
// 	// // mantissa
// 	// mantissa := int32(f / math.Pow10(int(exp)))
// 	// log.Printf("Mantissa: %d", mantissa)
// 	// s |= dg.DwordT(mantissa) >> 8 & 0x00ff_ffff
// 	// return s

// }

// Float64toDGdouble converts a standard Go float64 into the DG Double Precision format
func Float64toDGdouble(f float64) (q dg.QwordT) {
	var bits, sign, exponent uint64
	if f == 0.0 {
		return 0
	}
	bits = math.Float64bits(f)
	sign = bits & 0x8000_0000_0000_0000
	exponent = ((bits >> 52) & 0x7ff) - 1019
	bits = (bits & 0x000f_ffff_ffff_ffff) | 0x0010_0000_0000_0000
	bits <<= (exponent & 0x03)
	exponent >>= 2
	q = dg.QwordT(sign | ((exponent + 64) << 56) | bits)
	return q
}

// DGsingleToFloat64 converts a DG Single Precison number to Go float64
func DGsingleToFloat64(s dg.DwordT) (f float64) {
	exp := int8((s&0x7f00_0000)>>24) - 64
	mantissa := int64(s & 0x00ff_ffff)
	f = float64(mantissa) * math.Pow10(int(exp))
	if TestDwbit(s, 0) {
		f = f * -1.0
	}
	return f
}

// DGdoubleToFloat64 converts a DG Double Precison number to Go float64
func DGdoubleToFloat64(q dg.QwordT) (f float64) {
	// exp := int8((d&0x7f00_0000_0000_0000)>>56) - 64
	// mantissa := int64(d & 0x00ff_ffff_ffff_ffff)
	// f = float64(mantissa) * math.Pow10(int(exp))
	// if TestQwbit(d, 0) {
	// 	f = f * -1.0
	// }
	// return f
	var left, mantissa, exponent uint64
	mantissa = uint64(q & 0x00ff_ffff_ffff_ffff)
	if mantissa == 0 {
		return 0.0
	}
	for (mantissa & 0x00f0_0000_0000_0000) == 0 {
		mantissa <<= 4
		left++
	}
	exponent = (uint64((q>>56)&0x7f)-65-left)*4 + 1023
	for (mantissa & 0x00e0_0000_0000_0000) != 0 {
		mantissa >>= 1
		exponent++
	}
	log.Printf("exp: %v, mant: %v\n", exponent, mantissa)
	// d = (d & 0x8000_0000_0000_0000) | (exponent << 52) | (mantissa & 0x000f_ffff_ffff_ffff)
	r := uint64(q&0x8000_0000_0000_0000) | (exponent << 52) | (mantissa & 0x000f_ffff_ffff_ffff)
	f = math.Float64frombits(r)
	return f
}
