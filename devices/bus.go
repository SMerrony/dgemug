// bus.go

// Copyright (C) 2018  Steve Merrony

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
	"fmt"
	"log"
	"sync"

	"github.com/SMerrony/dgemug"

	"github.com/SMerrony/dgemug/logging"
	"github.com/SMerrony/dgemug/util"
)

const devMax = 64

type (
	// ResetFunc stores an I/O reset func pointer
	ResetFunc func()

	// DataOutFunc stores a DOx func pointer
	DataOutFunc func(datum dg.WordT, abc byte, flag byte)

	// DataInFunc stores a DIx func pointer
	DataInFunc func(abc byte, flag byte) (datum dg.WordT)
)

type device struct {
	devMu           sync.RWMutex
	mnemonic        string
	priorityMaskBit uint
	resetFunc       ResetFunc
	dataOutFunc     DataOutFunc
	dataInFunc      DataInFunc
	simAttached     bool
	simImageName    string
	ioDevice        bool
	bootable        bool
	busy            bool
	done            bool
}

type devices [devMax]device

var (
	d               devices  // not exported
	irqMask         dg.WordT // not exported, use setter and getter funcs below
	IRQ             bool
	InterruptingDev [devMax]bool
)

// BusInit must be called before attaching any devices
func BusInit() {
	for dev := range d {
		d[dev].devMu.Lock()
		d[dev].mnemonic = ""
		d[dev].priorityMaskBit = 0
		d[dev].dataInFunc = nil
		d[dev].dataOutFunc = nil
		d[dev].simAttached = false
		d[dev].ioDevice = false
		d[dev].bootable = false
		d[dev].busy = false
		d[dev].done = false
		d[dev].devMu.Unlock()
	}
}

func BusAddDevice(devNum int, mnem string, pmb uint, att bool, io bool, boot bool) {
	if devNum >= devMax {
		log.Fatalf("ERROR: Attempt to add device with too-high device number: %#o", devNum)
	}
	d[devNum].devMu.Lock()
	d[devNum].mnemonic = mnem
	d[devNum].priorityMaskBit = pmb
	d[devNum].simAttached = att
	d[devNum].ioDevice = io
	d[devNum].bootable = boot
	logging.DebugPrint(logging.DebugLog, "INFO: Device %#o added to bus\n", devNum)
	d[devNum].devMu.Unlock()
}

func BusSetDataInFunc(devNum int, fn DataInFunc) {
	d[devNum].devMu.Lock()
	d[devNum].dataInFunc = fn
	logging.DebugPrint(logging.DebugLog, "INFO: Bus Data In function set for dev %#o (%d.)\n", devNum, devNum)
	d[devNum].devMu.Unlock()
}

func BusDataIn(devNum int, abc byte, flag byte) (datum dg.WordT) {
	if d[devNum].dataInFunc == nil {
		log.Fatalf("ERROR: busDataIn called for device %#o with no function set", devNum)
	}
	//cpuPtr.ac[iPtr.acd] = dg.DwordT(d[iPtr.devNum].dataInFunc(abc, iPtr.f))
	return d[devNum].dataInFunc(abc, flag)
}

func BusSetDataOutFunc(devNum int, fn DataOutFunc) {
	d[devNum].devMu.Lock()
	d[devNum].dataOutFunc = fn
	d[devNum].devMu.Unlock()
	logging.DebugPrint(logging.DebugLog, "INFO: Bus Data Out function set for dev %#o (%d.)\n", devNum, devNum)
}

func BusDataOut(devNum int, datum dg.WordT, abc byte, flag byte) {
	if d[devNum].dataOutFunc == nil {
		logging.DebugLogsDump()
		log.Fatalf("ERROR: busDataOut called for device %#o with no function set", devNum)
	}
	d[devNum].dataOutFunc(datum, abc, flag)
}

func BusSetResetFunc(devNum int, resetFn ResetFunc) {
	d[devNum].devMu.Lock()
	d[devNum].resetFunc = resetFn
	logging.DebugPrint(logging.DebugLog, "INFO: Bus reset function set for dev %#o\n", devNum)
	d[devNum].devMu.Unlock()
}

func BusResetDevice(devNum int) {
	d[devNum].devMu.Lock()
	io := d[devNum].ioDevice
	d[devNum].devMu.Unlock()
	if io {
		d[devNum].resetFunc()
	} else {
		log.Fatalf("ERROR: Attempt to reset non-I/O device %#o\n", devNum)
	}

}

func BusResetAllIODevices() {
	for dev := range d {
		d[dev].devMu.Lock()
		io := d[dev].ioDevice
		d[dev].devMu.Unlock()
		if io {
			BusResetDevice(dev)
		}
	}
}

func BusSetAttached(devNum int, imgName string) {
	d[devNum].devMu.Lock()
	d[devNum].simAttached = true
	d[devNum].simImageName = imgName
	d[devNum].devMu.Unlock()
}
func BusSetDetached(devNum int) {
	d[devNum].devMu.Lock()
	d[devNum].simAttached = false
	d[devNum].simImageName = ""
	d[devNum].devMu.Unlock()
}
func BusIsAttached(devNum int) bool {
	d[devNum].devMu.RLock()
	att := d[devNum].simAttached
	d[devNum].devMu.RUnlock()
	return att
}

func BusSetBusy(devNum int, f bool) {
	d[devNum].devMu.Lock()
	d[devNum].busy = f
	d[devNum].devMu.Unlock()
}

func BusSetDone(devNum int, f bool) {
	d[devNum].devMu.Lock()
	d[devNum].done = f
	d[devNum].devMu.Unlock()
}

func BusGetBusy(devNum int) bool {
	d[devNum].devMu.RLock()
	bz := d[devNum].busy
	d[devNum].devMu.RUnlock()
	return bz
}

func BusGetDone(devNum int) bool {
	d[devNum].devMu.RLock()
	dn := d[devNum].done
	d[devNum].devMu.RUnlock()
	return dn
}

func BusIsBootable(devNum int) bool {
	d[devNum].devMu.RLock()
	bt := d[devNum].bootable
	d[devNum].devMu.RUnlock()
	return bt
}

func BusIsIODevice(devNum int) bool {
	d[devNum].devMu.RLock()
	io := d[devNum].ioDevice
	d[devNum].devMu.RUnlock()
	return io
}

// BusSetIrqMask is a setter for the (whole) IRQ mask
func BusSetIrqMask(newMask dg.WordT) {
	irqMask = newMask
}

// BusIsDevMasked is a getter to see if the device is masked out from sending IRQs
func BusIsDevMasked(devNum int) (masked bool) {
	return util.TestWbit(irqMask, int(d[devNum].priorityMaskBit))
}

// BusSetDevMasked is a setter to make the device masked out from sending IRQs
func BusSetDevMasked(devNum int) {
	util.SetWbit(&irqMask, d[devNum].priorityMaskBit)
}

// BusClearDevMasked is a setter to make the device able to send IRQs
func BusClearDevMasked(devNum int) {
	util.ClearWbit(&irqMask, d[devNum].priorityMaskBit)
}

// BusGetHighestPriorityInt returns the device number of the highest priority device
// that has an outstanding interrupt
func BusGetHighestPriorityInt() (devNum int) {
	for devNum = range InterruptingDev {
		if InterruptingDev[devNum] {
			return devNum
		}
	}
	return 0 // ?
}

func BusGetPrintableDevList() string {
	lst := fmt.Sprintf(" #  Mnem   PMB  I/O Busy Done Status\012")
	var line string
	for dev := range d {
		d[dev].devMu.RLock()
		if d[dev].mnemonic != "" {
			line = fmt.Sprintf("%#3o %-6s %2d. %3d %4d %4d  ",
				dev, d[dev].mnemonic, d[dev].priorityMaskBit,
				util.BoolToInt(d[dev].ioDevice), util.BoolToInt(d[dev].busy), util.BoolToInt(d[dev].done))
			if d[dev].simAttached {
				line += "Attached"
				if d[dev].simImageName != "" {
					line += " to image: " + d[dev].simImageName
				}
			} else {
				line += "Not Attached"
			}
			// Commented out the below for now as it's a bit misleading...
			// if d[dev].bootable {
			// 	line += ", Bootable"
			// }
			line += "\012"
			lst += line
		}
		d[dev].devMu.RUnlock()
	}
	return lst
}
