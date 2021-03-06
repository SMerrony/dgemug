// +build physical !virtual

// REPRESENTATION OF PHYSICAL MEMORY USED IN THE HARDWARE EMULATOR(S)

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
	"fmt"
	"log"
	"os"
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
