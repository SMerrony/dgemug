// ram.go

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
	"log"
	"runtime/debug"
	"sync"

	"github.com/SMerrony/dgemug"
	"github.com/SMerrony/dgemug/logging"
	"github.com/SMerrony/dgemug/util"
)

const (
	// MemSizeWords defines the size of MV/Em's emulated RAM in 16-bit words
	MemSizeWords = 8388608
	// MemSizeLCPID is the code returned by the LCPID to indicate the size of RAM
	MemSizeLCPID = 0x3F
)

// The memoryT structure holds our representation of system RAM.
// It is not exported and should not be directly accessed other than within this package.
type memoryT struct {
	ram        [MemSizeWords]dg.WordT
	ramMu      sync.RWMutex
	atuEnabled bool
}

var (
	memory memoryT // memory is NOT exported
)

// MemInit should be called at machine start
func MemInit(doLog bool) {
	memory.atuEnabled = false
	bmcdchInit(doLog)
	for addr := range memory.ram {
		memory.ram[addr] = 0
	}
	log.Printf("INFO: Initialised %d words of main memory\n", MemSizeWords)
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
	hiLo = util.TestWbit(byteAddr16, 15) // determine which byte to get
	addr = dg.PhysAddrT(byteAddr16) >> 1
	return ReadByte(addr, hiLo)
}

// WriteByte takes a normal word addr, low-byte flag and datum byte
func WriteByte(wordAddr dg.PhysAddrT, loByte bool, b dg.ByteT) {
	memory.ramMu.RLock()
	wd := memory.ram[wordAddr]
	memory.ramMu.RUnlock()
	if loByte {
		wd = (wd & 0xff00) | dg.WordT(b)
	} else {
		wd = dg.WordT(b)<<8 | (wd & 0x00ff)
	}
	WriteWord(wordAddr, wd)
}

// ReadWord returns the DG Word at the specified physical address
func ReadWord(wordAddr dg.PhysAddrT) dg.WordT {
	var wd dg.WordT
	if wordAddr >= MemSizeWords {
		logging.DebugLogsDump()
		debug.PrintStack()
		log.Fatalf("ERROR: Attempt to read word beyond end of physical memory using address: %d", wordAddr)
	}
	memory.ramMu.RLock()
	wd = memory.ram[wordAddr]
	memory.ramMu.RUnlock()
	return wd
}

// ReadWordTrap returns the DG Word at the specified physical address
func ReadWordTrap(wordAddr dg.PhysAddrT) (dg.WordT, bool) {
	var wd dg.WordT
	if wordAddr >= MemSizeWords {
		logging.DebugLogsDump()
		debug.PrintStack()
		log.Printf("ERROR: Attempt to read word beyond end of physical memory using address: %d", wordAddr)
		return 0, false
	}
	memory.ramMu.RLock()
	wd = memory.ram[wordAddr]
	memory.ramMu.RUnlock()
	return wd, true
}

// WriteWord - ALL memory-writing should ultimately go through this function
// N.B. minor exceptions may be made for memory.NsPush() and memory.NsPop()
func WriteWord(wordAddr dg.PhysAddrT, datum dg.WordT) {
	// if wordAddr == 2891 {
	// 	debugCatcher()
	// }
	if wordAddr >= MemSizeWords {
		debug.PrintStack()
		log.Fatalf("ERROR: Attempt to write word beyond end of physical memory using address: %d", wordAddr)
	}
	memory.ramMu.Lock()
	memory.ram[wordAddr] = datum
	memory.ramMu.Unlock()
}

// ReadDWord returns the doubleword at the given physical address
func ReadDWord(wordAddr dg.PhysAddrT) dg.DwordT {
	var hiWd, loWd dg.WordT
	if wordAddr >= MemSizeWords {
		debug.PrintStack()
		log.Fatalf("ERROR: Attempt to read doubleword beyond end of physical memory using address: %d", wordAddr)
	}
	memory.ramMu.RLock()
	hiWd = memory.ram[wordAddr]
	loWd = memory.ram[wordAddr+1]
	memory.ramMu.RUnlock()
	return util.DwordFromTwoWords(hiWd, loWd)
}

// ReadDwordTrap returns the doubleword at the given physical address
func ReadDwordTrap(wordAddr dg.PhysAddrT) (dg.DwordT, bool) {
	var hiWd, loWd dg.WordT
	if wordAddr >= MemSizeWords {
		logging.DebugLogsDump()
		log.Printf("ERROR: Attempt to read doubleword beyond end of physical memory using address: %d\n", wordAddr)
		return 0, false
	}
	memory.ramMu.RLock()
	hiWd = memory.ram[wordAddr]
	loWd = memory.ram[wordAddr+1]
	memory.ramMu.RUnlock()
	return util.DwordFromTwoWords(hiWd, loWd), true
}

// WriteDWord writes a doubleword into memory at the given physical address
func WriteDWord(wordAddr dg.PhysAddrT, dwd dg.DwordT) {
	if wordAddr >= MemSizeWords {
		log.Fatalf("ERROR: Attempt to write doubleword beyond end of physical memory using address: %d", wordAddr)
	}
	WriteWord(wordAddr, util.DwordGetUpperWord(dwd))
	WriteWord(wordAddr+1, util.DwordGetLowerWord(dwd))
}
