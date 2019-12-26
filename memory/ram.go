// ram.go

// Copyright (C) 2017,2019  Steve Merrony

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
	"fmt"
	"log"
	"runtime/debug"
	"sync"

	"github.com/SMerrony/dgemug/dg"
	"github.com/SMerrony/dgemug/logging"
)

var (
	ram          []dg.WordT
	ramMu        sync.RWMutex
	atuEnabled   bool
	memSizeWords dg.PhysAddrT // just for efficiency
)

// MemInit should be called at machine start
func MemInit(wordSize int, doLog bool) {
	ram = make([]dg.WordT, wordSize)
	memSizeWords = dg.PhysAddrT(wordSize)
	atuEnabled = false
	bmcdchInit(doLog)
	for addr := range ram {
		ram[addr] = 0
	}
	log.Printf("INFO: Initialised %#o words of main memory\n", wordSize)
}

// GetSegment - return the segment number for the supplied address
func GetSegment(addr dg.PhysAddrT) int {
	return int((addr & 0x70000000) >> 28)
}

// ReadByte - read a byte from memory using word address and low-byte flag (true => lower (rightmost) byte)
func ReadByte(wordAddr dg.PhysAddrT, loByte bool) dg.ByteT {
	var res dg.ByteT
	wd := ReadWord(wordAddr)
	if loByte {
		res = dg.ByteT(wd & 0xff)
	} else {
		res = dg.ByteT(wd >> 8)
	}
	return res
}

// ReadByteEclipseBA - read a byte - special version for Eclipse Byte-Addressing
func ReadByteEclipseBA(byteAddr16 dg.WordT) dg.ByteT {
	var (
		hiLo bool
		addr dg.PhysAddrT
	)
	hiLo = TestWbit(byteAddr16, 15) // determine which byte to get
	addr = dg.PhysAddrT(byteAddr16) >> 1
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

// ReadWord16 returns the DG Word at the specified physical address
func ReadWord16(wordAddr dg.WordT) dg.WordT {
	var wd dg.WordT
	ramMu.RLock()
	wd = ram[wordAddr]
	ramMu.RUnlock()
	return wd
}

// ReadWord returns the DG Word at the specified physical address
func ReadWord(wordAddr dg.PhysAddrT) dg.WordT {
	var wd dg.WordT
	if wordAddr >= memSizeWords {
		logging.DebugLogsDump("logs/")
		debug.PrintStack()
		log.Fatalf("ERROR: Attempt to read word beyond end of physical memory (%#o) using address: %#o", memSizeWords, wordAddr)
	}
	ramMu.RLock()
	wd = ram[wordAddr]
	ramMu.RUnlock()
	return wd
}

// ReadWordTrap returns the DG Word at the specified physical address
func ReadWordTrap(wordAddr dg.PhysAddrT) (dg.WordT, bool) {
	var wd dg.WordT
	if wordAddr >= memSizeWords {
		logging.DebugLogsDump("logs/")
		debug.PrintStack()
		log.Printf("ERROR: Attempt to read word beyond end of physical memory (%#o) using address: %#o", memSizeWords, wordAddr)
		return 0, false
	}
	ramMu.RLock()
	wd = ram[wordAddr]
	ramMu.RUnlock()
	return wd, true
}

// ReadEclipseWordTrap returns the DG Word at the specified 16-bit physical address
func ReadEclipseWordTrap(wordAddr dg.WordT) (dg.WordT, bool) {
	var wd dg.WordT
	if dg.PhysAddrT(wordAddr) >= memSizeWords {
		logging.DebugLogsDump("logs/")
		debug.PrintStack()
		log.Printf("ERROR: Attempt to read word beyond end of physical memory (%#o) using address: %#o", memSizeWords, wordAddr)
		return 0, false
	}
	ramMu.RLock()
	wd = ram[wordAddr]
	ramMu.RUnlock()
	return wd, true
}

// WriteWord16 - For the 16-bit emulators ALL memory-writing should ultimately go through this function
// N.B. minor exceptions may be made for NsPush() and NsPop()
func WriteWord16(wordAddr dg.WordT, datum dg.WordT) {
	ramMu.Lock()
	ram[wordAddr] = datum
	ramMu.Unlock()
}

// WriteWord - For the 32-bit emulator ALL memory-writing should ultimately go through this function
// N.B. minor exceptions may be made for NsPush() and NsPop()
func WriteWord(wordAddr dg.PhysAddrT, datum dg.WordT) {
	// if wordAddr == 6 {
	// 	runtime.Breakpoint()
	// }
	if wordAddr >= memSizeWords {
		debug.PrintStack()
		logging.DebugLogsDump("logs/")
		log.Fatalf("ERROR: Attempt to write word beyond end of physical memory (%#o) using address: %#o", memSizeWords, wordAddr)
	}
	ramMu.Lock()
	ram[wordAddr] = datum
	ramMu.Unlock()
}

// ReadDWord returns the doubleword at the given physical address
func ReadDWord(wordAddr dg.PhysAddrT) dg.DwordT {
	var hiWd, loWd dg.WordT
	if wordAddr >= memSizeWords {
		debug.PrintStack()
		logging.DebugLogsDump("logs/")
		log.Fatalf("ERROR: Attempt to read doubleword beyond end of physical memory (%#o) using address: %#o", memSizeWords, wordAddr)
	}
	ramMu.RLock()
	hiWd = ram[wordAddr]
	loWd = ram[wordAddr+1]
	ramMu.RUnlock()
	return DwordFromTwoWords(hiWd, loWd)
}

// ReadDwordTrap returns the doubleword at the given physical address
func ReadDwordTrap(wordAddr dg.PhysAddrT) (dg.DwordT, bool) {
	var hiWd, loWd dg.WordT
	if wordAddr >= memSizeWords {
		debug.PrintStack()
		logging.DebugLogsDump("logs/")
		log.Printf("ERROR: Attempt to read doubleword beyond end of physical memory (%#o) using address: %#o\n", memSizeWords, wordAddr)
		return 0, false
	}
	ramMu.RLock()
	hiWd = ram[wordAddr]
	loWd = ram[wordAddr+1]
	ramMu.RUnlock()
	return DwordFromTwoWords(hiWd, loWd), true
}

// WriteDWord writes a doubleword into memory at the given physical address
func WriteDWord(wordAddr dg.PhysAddrT, dwd dg.DwordT) {
	WriteWord(wordAddr, DwordGetUpperWord(dwd))
	WriteWord(wordAddr+1, DwordGetLowerWord(dwd))
}

// BoolToInt converts a bool to 1 or 0
func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// BoolToYN converts a bool to Y or N
func BoolToYN(b bool) byte {
	if b {
		return 'Y'
	}
	return 'N'
}

// BoolToOnOff converts a bool to "On" or "Off"
func BoolToOnOff(b bool) string {
	if b {
		return "On"
	}
	return "Off"
}

// BoolToOZ converts a boolean to a O(ne) or Z(ero) byte
func BoolToOZ(b bool) byte {
	if b {
		return 'O'
	}
	return 'Z'
}

// DwordGetLowerWord gets the DG-lower word (RHS) of a doubleword
// Called VERY often, hopefully inlined!
func DwordGetLowerWord(dwd dg.DwordT) dg.WordT {
	return dg.WordT(dwd) // & 0x0000ffff mask unneccessary
}

// DwordGetUpperWord gets the DG-higher word (LHS) of a doubleword
// Called VERY often, hopefully inlined!
func DwordGetUpperWord(dwd dg.DwordT) dg.WordT {
	return dg.WordT(dwd >> 16)
}

// DwordFromTwoWords - catenate two DG Words into a DG DoubleWord
func DwordFromTwoWords(hw dg.WordT, lw dg.WordT) dg.DwordT {
	return dg.DwordT(hw)<<16 | dg.DwordT(lw)
}

// QwordFromTwoDwords - catenate two DG DoubleWords into a DG QuadWord
func QwordFromTwoDwords(hdw dg.DwordT, ldw dg.DwordT) dg.QwordT {
	return dg.QwordT(hdw)<<32 | dg.QwordT(ldw)
}

// WordToBinStr - get a pretty-printable string of a word
func WordToBinStr(w dg.WordT) string {
	return fmt.Sprintf("%08b %08b", w>>8, w&0x0ff)
}

// SexWordToDword - sign-extend a DG word to a DG DoubleWord
func SexWordToDword(wd dg.WordT) dg.DwordT {
	var dwd dg.DwordT
	if TestWbit(wd, 0) {
		dwd = dg.DwordT(wd) | 0xffff0000
	} else {
		dwd = dg.DwordT(wd) & 0x0000ffff
	}
	return dwd
}

// SwapBytes - swap over the two bytes in a dg_word
func SwapBytes(wd dg.WordT) (res dg.WordT) {
	res = (wd >> 8) | ((wd & 0x00ff) << 8)
	return res
}
