// +build virtual

// REPRESENTATION OF VIRTUAL MEMORY USED IN THE OS-LEVEL EMULATOR(S)

// Copyright ©2020  Steve Merrony

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
	"sync"

	"github.com/SMerrony/dgemug/dg"
)

const (
	memPageSizeWords = 1024
	ring7page0       = 0x7000_0000 >> 10
)

type pageT struct {
	words [memPageSizeWords]dg.WordT
}

var (
	virtualRam       map[int]pageT
	virtualRamMu     sync.RWMutex
	lastUnsharedPage int = -1
	firstSharedPage  int = 0x7fff_ffff // dummy high value
	// lastSharedPage   int
)

// MapPage maps (allocates) a 1kW page of virtual memory for the process
func MapPage(page int, shared bool) {
	virtualRamMu.Lock()
	if _, mapped := virtualRam[page]; mapped {
		log.Fatalf("ERROR: Attempt to map already-mapped memory page %#o", page)
	}
	var emptyPage pageT
	virtualRam[page] = emptyPage
	if !shared {
		lastUnsharedPage = page
	}
	if shared && (page < firstSharedPage) {
		firstSharedPage = page
	}
	virtualRamMu.Unlock()
	log.Printf("DEBUG: Mapped page %#x for %#x", page, page<<10)
	if !shared {
		log.Printf("DEBUG: ...Last unshared page is now: %#x\n", lastUnsharedPage)
	}
}

// GetFirstSharedPage is a getter for the lowest shared page currently mapped
func GetFirstSharedPage() dg.DwordT {
	virtualRamMu.RLock()
	p := firstSharedPage
	virtualRamMu.RUnlock()
	return dg.DwordT(p)
}

// AddUnsharedPage appends an unshared page to virtual memory
func AddUnsharedPage() int {
	virtualRamMu.RLock()
	nextPage := lastUnsharedPage + 1
	virtualRamMu.RUnlock()
	MapPage(nextPage, false)
	return nextPage
}

// GetLastUnsharedPage is a getter for the highest unshared page currently mapped
func GetLastUnsharedPage() dg.DwordT {
	virtualRamMu.RLock()
	p := lastUnsharedPage
	virtualRamMu.RUnlock()
	return dg.DwordT(p)
}

func isAddrMapped(addr dg.PhysAddrT) (mapped bool) {
	virtualRamMu.RLock()
	_, mapped = virtualRam[int(addr>>10)]
	virtualRamMu.RUnlock()
	return mapped
}

// MapSlice maps (copies) the provided slice to virtual memory starting at the given address
func MapSlice(addr dg.PhysAddrT, wds []dg.WordT, shared bool) {
	for offset, word := range wds {
		loc := addr + dg.PhysAddrT(offset)
		// check each time we hit a page boundary to see if it's mapped
		if ((loc & 0x3ff) == 0) && !isAddrMapped(loc) {
			MapPage(int(loc>>10), shared)
		}
		WriteWord(loc, word)
	}
}

// UnmapPage unmaps (deallocates) a 1kW page of virtual memory from the process
func UnmapPage(page int, shared bool) {
	virtualRamMu.Lock()
	if _, mapped := virtualRam[page]; !mapped {
		log.Fatalf("ERROR: Attempt to unmap a non-mapped memory page %#o", page)
	}
	delete(virtualRam, page)
	if !shared {
		lastUnsharedPage--
	}
	virtualRamMu.Unlock()
	log.Printf("DEBUG: Unpapped page %#x", page)
	if !shared {
		log.Printf("DEBUG: ...Last unshared page is now: %#x\n", lastUnsharedPage)
	}
}

// MemInit must be called when the virtual machine is started
func MemInit() {
	virtualRam = make(map[int]pageT)
	// always map user page 0
	MapPage(ring7page0, false)
}

// ReadBytes - read specified # of bytes from 32-bit BA into slice
func ReadBytes(ba32 dg.DwordT, num int) (res []byte) {
	var c dg.DwordT
	for c = 0; c < dg.DwordT(num); c++ {
		if (ba32+c)&0x01 == 1 {
			res = append(res, byte(ReadByte(dg.PhysAddrT((ba32+c)>>1), true)))
		} else {
			res = append(res, byte(ReadByte(dg.PhysAddrT((ba32+c)>>1), false)))
		}
	}
	return res
}

// WriteStringBA copies a string to the specified address
func WriteStringBA(s string, byteAddr dg.DwordT) {
	for c := 0; c < len(s); c++ {
		WriteByteBA(byteAddr+dg.DwordT(c), dg.ByteT(s[c]))
	}
}

// ReadWord reads a single 16-bit word from the specified address
func ReadWord(addr dg.PhysAddrT) (wd dg.WordT) {
	virtualRamMu.RLock()
	page, found := virtualRam[int(addr>>10)]
	if !found {
		log.Fatalf("ERROR: Attempt to read from unmapped page %#x at address: %#x", addr>>10, addr)
	}
	wd = page.words[int(addr&0x3ff)]
	virtualRamMu.RUnlock()

	return wd
}

func ReadWordTrap(addr dg.PhysAddrT) (dg.WordT, bool) {
	if !isAddrMapped(addr) {
		log.Printf("ERROR: Attempt to read unmapped word at %#x\n", addr)
		return 0, false
	}
	return ReadWord(addr), true
}

func WriteWord(addr dg.PhysAddrT, datum dg.WordT) {
	virtualRamMu.Lock()
	page, found := virtualRam[int(addr>>10)]
	if !found {
		log.Fatalf("ERROR: Attempt to write to unmapped page %#x for addr %#x", addr>>10, addr)
	}
	page.words[int(addr&0x3ff)] = datum
	virtualRam[int(addr>>10)] = page
	virtualRamMu.Unlock()
}

func ReadDwordTrap(addr dg.PhysAddrT) (dg.DwordT, bool) {
	if !isAddrMapped(addr) {
		log.Printf("ERROR: Attempt to read unmapped doubleword at %#x\n", addr)
		return 0, false
	}
	return ReadDWord(addr), true
}
