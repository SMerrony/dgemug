// tto - console output

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

	dg "github.com/SMerrony/dgemug/dg"
)

const (
	asciiNL = 0x0a
	asciiFF = 0x0c
)

var (
	tto    net.Conn // FIXME should be a type...
	bus    *BusT
	devNum int
)

// TtoInit performs iniial setup of the TTO device
func TtoInit(dev int, busp *BusT, c net.Conn) {
	devNum = dev
	bus = busp
	tto = c
	bus.SetResetFunc(devNum, ttoReset)
	bus.SetDataOutFunc(devNum, ttoDataOut)
}

// TtoPutChar outputs a single byte to TTO
func TtoPutChar(c byte) {
	tto.Write([]byte{c})
}

// TtoPutString outputs a string to TTO (no NL appended)
func TtoPutString(s string) {
	tto.Write([]byte(s))
}

// TtoPutStringNL outputs a string followed by a NL to TTO
func TtoPutStringNL(s string) {
	tto.Write([]byte(s))
	tto.Write([]byte{asciiNL})
}

// TtoPutNLString outputs a NL followed by a string to TTO
func TtoPutNLString(s string) {
	tto.Write([]byte{asciiNL})
	tto.Write([]byte(s))
}

func ttoReset() {
	TtoPutChar(asciiFF)
	log.Println("INFO: TTO Reset")
}

// This is called from Bus to implement DOA to the TTO device
func ttoDataOut(datum dg.WordT, abc byte, flag byte) {
	var ascii byte
	switch abc {
	case 'A':
		ascii = byte(datum)
		if flag == 'S' {
			bus.SetBusy(devNum, true)
			bus.SetDone(devNum, false)
		}
		TtoPutChar(ascii)
		bus.SetBusy(devNum, false)
		bus.SetDone(devNum, true)
		// send IRQ if not masked out
		if !bus.IsDevMasked(devNum) {
			// InterruptingDev[devNum] = true
			// IRQ = true
			bus.SendInterrupt(devNum)
		}
	case 'N':
		switch flag {
		case 'S':
			bus.SetBusy(devNum, true)
			bus.SetDone(devNum, false)
		case 'C':
			bus.SetBusy(devNum, false)
			bus.SetDone(devNum, false)
		}
	default:
		log.Fatalf("ERROR: unexpected source buffer <%c> for DOx ac,TTO instruction\n", abc)
	}
}
