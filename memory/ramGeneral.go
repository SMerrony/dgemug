// Routines that apply to both physical and virtual RAM

// Copyright ©2017-2020  Steve Merrony

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
	"bytes"
	"fmt"
	"os"

	"github.com/SMerrony/dgemug/dg"
)

// GetSegment - return the segment number for the supplied address
func GetSegment(addr dg.PhysAddrT) int {
	return int((addr & 0x70000000) >> 28)
}

// ReadByte - read a byte from memory using word address and low-byte flag (true => lower (rightmost) byte)
func ReadByte(wordAddr dg.PhysAddrT, loByte bool) dg.ByteT {
	wd := ReadWord(wordAddr)
	if !loByte {
		wd >>= 8
	}
	return dg.ByteT(wd)
}

// ReadNBytes loads n bytes into a slice which is returned - no ring-forcing is performed
func ReadNBytes(ba dg.PhysAddrT, n int) []byte {
	buff := bytes.NewBufferString("")
	lobyte := (ba & 0x0001) == 1
	wdAddr := dg.PhysAddrT(ba >> 1)
	for b := 0; b < n; b++ {
		c := ReadByte(wdAddr, lobyte)
		buff.WriteByte(byte(c))
		if lobyte {
			wdAddr++
		}
		lobyte = !lobyte
	}
	return buff.Bytes()
}

// ReadByteEclipseBA - read a byte - special version for Eclipse Byte-Addressing
func ReadByteEclipseBA(pcAddr dg.PhysAddrT, byteAddr16 dg.WordT) dg.ByteT {
	var (
		hiLo bool
		addr dg.PhysAddrT
	)
	hiLo = TestWbit(byteAddr16, 15) // determine which byte to get
	addr = dg.PhysAddrT(byteAddr16) >> 1
	addr |= (pcAddr & 0x7000_0000)
	return ReadByte(addr, hiLo)
}

// WriteByte takes a normal word addr, low-byte flag and datum byte

func WriteByte(wordAddr dg.PhysAddrT, loByte bool, b dg.ByteT) {
	wd := ReadWord(wordAddr)
	if loByte {
		wd = (wd & 0xff00) | dg.WordT(b)
	} else {
		wd = dg.WordT(b)<<8 | (wd & 0x00ff)
	}
	WriteWord(wordAddr, wd)
}

// WriteByteBA writes a byte to a standard Byte Addressed location
func WriteByteBA(byteAddr dg.DwordT, b dg.ByteT) {
	loByte := (byteAddr & 0x01) == 1
	WriteByte(dg.PhysAddrT(byteAddr>>1), loByte, b)
}

// ReadDWord does what it says on the tin
func ReadDWord(addr dg.PhysAddrT) dg.DwordT {
	var hiWd, loWd dg.WordT
	hiWd = ReadWord(addr)
	loWd = ReadWord(addr + 1)
	return DwordFromTwoWords(hiWd, loWd)
}

// WriteDWord writes a doubleword into memory at the given physical address
func WriteDWord(wordAddr dg.PhysAddrT, dwd dg.DwordT) {
	WriteWord(wordAddr, DwordGetUpperWord(dwd))
	WriteWord(wordAddr+1, DwordGetLowerWord(dwd))
}

// DumpToFile writes out usefully greppable text representation of memory.
func DumpToFile(fn string) bool {
	f, err := os.Create(fn)
	if err != nil {
		return false
	}
	defer f.Close()

	for a := 0; a < int(memSizeWords); a++ {
		_, err := f.WriteString(fmt.Sprintf("%12o %04X\n", a, ReadWord(dg.PhysAddrT(a))))
		if err != nil {
			return false
		}
	}
	return true
}
