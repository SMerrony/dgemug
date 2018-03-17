// tto - console output

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

package devices

import (
	"log"
	"net"

	dg "github.com/SMerrony/dgemug"
)

const (
	asciiNL = 0x0a
	asciiFF = 0x0c
)

var (
	tto    net.Conn
	devNum int
)

func TtoInit(dev int, c net.Conn) {
	devNum = dev
	tto = c
	BusSetResetFunc(devNum, ttoReset)
	BusSetDataOutFunc(devNum, ttoDataOut)
}

func TtoPutChar(c byte) {
	tto.Write([]byte{c})
}

func TtoPutString(s string) {
	tto.Write([]byte(s))
}

func TtoPutStringNL(s string) {
	tto.Write([]byte(s))
	tto.Write([]byte{asciiNL})
}

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
			BusSetBusy(devNum, true)
			BusSetDone(devNum, false)
		}
		TtoPutChar(ascii)
		BusSetBusy(devNum, false)
		BusSetDone(devNum, true)
		// send IRQ if not masked out
		if !BusIsDevMasked(devNum) {
			// InterruptingDev[devNum] = true
			// IRQ = true
			BusSendInterrupt(devNum)
		}
	case 'N':
		switch flag {
		case 'S':
			BusSetBusy(devNum, true)
			BusSetDone(devNum, false)
		case 'C':
			BusSetBusy(devNum, false)
			BusSetDone(devNum, false)
		}
	default:
		log.Fatalf("ERROR: unexpected source buffer <%c> for DOx ac,TTO instruction\n", abc)
	}
}
