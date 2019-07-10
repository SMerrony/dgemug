// loader.go

// Copyright (C) 2019  Steve Merrony

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
	"bufio"
	"encoding/csv"
	"io"
	"os"
	"strconv"

	"github.com/SMerrony/dgemug/dg"
)

// LoadFromASCIIFile reads a CSV-formatted ASCII file representing the contents
// of memory in octal words and loads it directly into memory.
// It can be used to directly load assembled listings such as diagnostics.
func LoadFromASCIIFile(asciiOctalFilename string) string {
	asciiOctalFile, err := os.Open(asciiOctalFilename)
	if err != nil {
		return "*** ERROR: Could not access ASCII Octal CSV load file " + asciiOctalFilename + err.Error()
	}
	defer asciiOctalFile.Close()
	csvReader := csv.NewReader(bufio.NewReader(asciiOctalFile))
	csvReader.Comment = ';'
	csvReader.FieldsPerRecord = 2

	var count int
	var contents dg.WordT

	for {
		csvRec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "*** ERROR: Could not parse CSV line - " + err.Error() + " ***"
		}

		// 1st field is the word address, 2nd is the word contents
		thisAddr64, err := strconv.ParseUint(csvRec[0], 8, 16)
		if err != nil {
			return "*** ERROR: Could not parse Address - " + csvRec[0] + " ***"
		}
		thisAddr := dg.PhysAddrT(thisAddr64)

		contents64, err := strconv.ParseUint(csvRec[1], 8, 16)
		if err != nil {
			return "*** ERROR: Could not parse Contents - " + csvRec[1] + " ***"
		}
		contents = dg.WordT(contents64)

		WriteWord(thisAddr, contents)
		count++
	}
	return "Words loaded: " + strconv.Itoa(count)
}
