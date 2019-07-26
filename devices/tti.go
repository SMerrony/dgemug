//  TTI - console output

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

package devices

import (
	"log"
	"sync"

	dg "github.com/SMerrony/dgemug/dg"
)

// TtiT describes the current state of the TTI device
type TtiT struct {
	ttiMu      sync.RWMutex
	oneCharBuf byte
	bus        *BusT
	devNum     int
}

// Init performs iniial setup of the TTI device
func (tti *TtiT) Init(dev int, bus *BusT) {
	tti.ttiMu.Lock()
	tti.devNum = dev
	tti.bus = bus
	bus.SetResetFunc(dev, tti.reset)
	bus.SetDataInFunc(dev, tti.dataIn)
	bus.SetDataOutFunc(dev, tti.dataOut)
	tti.ttiMu.Unlock()
}

// InsertChar places one byte in the TTI buffer for handling by the CPU
func (tti *TtiT) InsertChar(c byte) {
	tti.ttiMu.Lock()
	tti.oneCharBuf = c
	tti.ttiMu.Unlock()
	tti.bus.SetDone(tti.devNum, true)
	// send IRQ if not masked out
	if !tti.bus.IsDevMasked(tti.devNum) {
		tti.bus.SendInterrupt(tti.devNum)
	}
}

// Reset is a stub
func (tti *TtiT) reset() {
	log.Println("INFO: TTI Reset")
}

// This is called from Bus to implement DIA from the TTI device
func (tti *TtiT) dataIn(abc byte, flag byte) (datum dg.WordT) {
	tti.ttiMu.Lock()
	datum = dg.WordT(tti.oneCharBuf) // grab the char from the buffer
	switch abc {
	case 'A':
		switch flag {
		case 'S':
			tti.bus.SetBusy(tti.devNum, true)
			tti.bus.SetDone(tti.devNum, false)
		case 'C':
			tti.bus.SetBusy(tti.devNum, false)
			tti.bus.SetDone(tti.devNum, false)
		}
	default:
		log.Fatalf("ERROR: unexpected source buffer <%c> for DOx ac,TTO instruction\n", abc)
	}
	tti.ttiMu.Unlock()
	return datum
}

// this is only here to support NIO commands to TTI
func (tti *TtiT) dataOut(datum dg.WordT, abc byte, flag byte) {
	tti.ttiMu.RLock()
	switch abc {
	case 'N':
		switch flag {
		case 'S':
			tti.bus.SetBusy(tti.devNum, true)
			tti.bus.SetDone(tti.devNum, false)
		case 'C':
			tti.bus.SetBusy(tti.devNum, false)
			tti.bus.SetDone(tti.devNum, false)
		}
	default:
		log.Fatalf("ERROR: unexpected call to ttiDataOut with abc(n) flag set to %c\n", abc)
	}
	tti.ttiMu.RUnlock()
}
