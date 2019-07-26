//  - console output

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
	"net"
	"sync"

	dg "github.com/SMerrony/dgemug/dg"
)

const (
	asciiNL = 0x0a
	asciiFF = 0x0c
)

// TtoT describes the current state of the TTO device
type TtoT struct {
	ttoMu  sync.RWMutex
	conn   net.Conn
	bus    *BusT
	devNum int
}

// Init performs iniial setup of the TTO device
func (tto *TtoT) Init(dev int, bus *BusT, c net.Conn) {
	tto.ttoMu.Lock()
	tto.devNum = dev
	tto.bus = bus
	tto.conn = c
	bus.SetResetFunc(dev, tto.Reset)
	bus.SetDataOutFunc(dev, tto.dataOut)
	tto.ttoMu.Unlock()
}

// PutChar outputs a single byte to TTO
func (tto *TtoT) PutChar(c byte) {
	tto.ttoMu.Lock()
	tto.conn.Write([]byte{c})
	tto.ttoMu.Unlock()
}

// PutString outputs a string to TTO (no NL appended)
func (tto *TtoT) PutString(s string) {
	tto.ttoMu.Lock()
	tto.conn.Write([]byte(s))
	tto.ttoMu.Unlock()
}

// PutStringNL outputs a string followed by a NL to TTO
func (tto *TtoT) PutStringNL(s string) {
	tto.ttoMu.Lock()
	tto.conn.Write([]byte(s))
	tto.conn.Write([]byte{asciiNL})
	tto.ttoMu.Unlock()
}

// PutNLString outputs a NL followed by a string to TTO
func (tto *TtoT) PutNLString(s string) {
	tto.ttoMu.Lock()
	tto.conn.Write([]byte{asciiNL})
	tto.conn.Write([]byte(s))
	tto.ttoMu.Unlock()
}

// Reset simply clears the screen or throws a page
func (tto *TtoT) Reset() {
	tto.PutChar(asciiFF)
	log.Println("INFO: TTO Reset")
}

// This is called from Bus to implement DOA to the TTO device
func (tto *TtoT) dataOut(datum dg.WordT, abc byte, flag byte) {
	var ascii byte
	switch abc {
	case 'A':
		ascii = byte(datum)
		if flag == 'S' {
			tto.bus.SetBusy(tto.devNum, true)
			tto.bus.SetDone(tto.devNum, false)
		}
		tto.PutChar(ascii)
		tto.bus.SetBusy(tto.devNum, false)
		tto.bus.SetDone(tto.devNum, true)
		// send IRQ if not masked out
		if !tto.bus.IsDevMasked(tto.devNum) {
			// InterruptingDev[tto.devNum] = true
			// IRQ = true
			tto.bus.SendInterrupt(tto.devNum)
		}
	case 'N':
		switch flag {
		case 'S':
			tto.bus.SetBusy(tto.devNum, true)
			tto.bus.SetDone(tto.devNum, false)
		case 'C':
			tto.bus.SetBusy(tto.devNum, false)
			tto.bus.SetDone(tto.devNum, false)
		}
	default:
		log.Fatalf("ERROR: unexpected source buffer <%c> for DOx ac,TTO instruction\n", abc)
	}
}
