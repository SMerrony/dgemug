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
	virtualRam   map[int]pageT
	virtualRamMu sync.RWMutex
)

func mapPage(page int) {
	virtualRamMu.Lock()
	if _, mapped := virtualRam[page]; mapped {
		log.Fatalf("ERROR: Attempt to map already-mapped memory page %#o", page)
	}
	var emptyPage pageT
	virtualRam[page] = emptyPage
	virtualRamMu.Unlock()
}

func isAddrMapped(addr dg.PhysAddrT) (mapped bool) {
	virtualRamMu.RLock()
	_, mapped = virtualRam[int(addr>>10)]
	virtualRamMu.RUnlock()
	return mapped
}

func MapSlice(addr dg.PhysAddrT, wds []dg.WordT) {
	for offset, word := range wds {
		loc := addr + dg.PhysAddrT(offset)
		// check each time we hit a page boundary to see if it's mapped
		if ((loc & 0x3ff) == 0) && !isAddrMapped(loc) {
			mapPage(int(loc >> 10))
		}
		WriteWord(loc, word)
	}
}

func unmapPage(page int) {
	virtualRamMu.Lock()
	if _, mapped := virtualRam[page]; !mapped {
		log.Fatalf("ERROR: Attempt to unmap a non-mapped memory page %#o", page)
	}
	delete(virtualRam, page)
	virtualRamMu.Unlock()
}

// MemInit must be called when the virtual machine is started
func MemInit() {
	virtualRam = make(map[int]pageT)
	// always map user page 0
	mapPage(ring7page0)
}

func ReadWord(addr dg.PhysAddrT) (wd dg.WordT) {
	virtualRamMu.RLock()
	page, found := virtualRam[int(addr>>10)]
	if !found {
		log.Fatalf("ERROR: Attempt to read from unmapped page")
	}
	wd = page.words[int(addr&0x3ff)]
	virtualRamMu.RUnlock()

	return wd
}

func WriteWord(addr dg.PhysAddrT, datum dg.WordT) {
	virtualRamMu.Lock()
	page, found := virtualRam[int(addr>>10)]
	if !found {
		log.Fatalf("ERROR: Attempt to write to unmapped page")
	}
	page.words[int(addr&0x3ff)] = datum
	virtualRam[int(addr>>10)] = page
	virtualRamMu.Unlock()
}

func ReadDWord(addr dg.PhysAddrT) dg.DwordT {
	var hiWd, loWd dg.WordT
	hiWd = ReadWord(addr)
	loWd = ReadWord(addr + 1)
	return DwordFromTwoWords(hiWd, loWd)
}

func WriteDWord(wordAddr dg.PhysAddrT, dwd dg.DwordT) {
	WriteWord(wordAddr, DwordGetUpperWord(dwd))
	WriteWord(wordAddr+1, DwordGetLowerWord(dwd))
}
