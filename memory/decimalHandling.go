// decimalHandling.go - handle DG packed (BCD) and unpacked (ASCII) decimal formats

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
	"strconv"
	"strings"

	"github.com/SMerrony/dgemug/dg"
)

// Decimal Data Types - values are significant
const (
	UnpackedDecTSC = 0
	UnpackedDecLSC = 1
	UnpackedDecTS  = 2
	UnpackedDecLS  = 3 // <sign><zeroes><int>
	UnpackedDecU   = 4 // <zeroes><int>
	PackedDec      = 5
	TwosCompDec    = 6
	FPDec          = 7
)

// DecodeDecDataType extracts the scale, type and length from a Decimal Type Indicator
func DecodeDecDataType(dti dg.DwordT) (scaleFactor int8, decType int, size int) {
	scaleFactor = int8(GetDwbits(dti, 0, 8))
	decType = int(uint8(GetDwbits(dti, 24, 3)))
	size = int(uint8(GetDwbits(dti, 27, 5)))
	if decType != 5 {
		size++
	}
	return scaleFactor, decType, size
}

// DecIntToInt returns an int value of the DG Decimal integer supplied
func DecIntToInt(decType int, dec string) (i int) {
	var err error
	switch decType {
	case UnpackedDecLS, UnpackedDecU:
		i, err = strconv.Atoi(strings.TrimSpace(dec))
		if err != nil {
			log.Panicf("ERROR: Could not parse Decimal <%s> as integer\n", dec)
		}
	default:
		log.Panicf("DecIntToInt does not yet handle data type %d.\n", decType)
	}
	return i
}

// ReadDec returns a string of the Decimal value pointed to by the given byte address
func ReadDec(ba dg.PhysAddrT, size int) (dec string) {
	bytes := ReadNBytes(ba, size)
	dec = strings.TrimSpace(string(bytes))
	return dec
}
