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

	dg "github.com/SMerrony/dgemug/dg"
	"github.com/SMerrony/dgemug/logging"
	"github.com/SMerrony/dgemug/memory"
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

// DeviceDesc holds basic config info for a device, a VM will have a map of these to
// describe its known devices
type DeviceDesc struct {
	DgMnemonic string
	PMB        uint // Priority Mask Bit number
	IsIO       bool
	IsBootable bool
}

// DeviceMapT describes the Device Map used by each VM
type DeviceMapT map[int]DeviceDesc

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
	d       devices  // not exported
	irqMask dg.WordT // not exported, use setter and getter funcs below
	// IRQ indicates if we are currently being interrupted
	IRQ             bool
	irqsByPriority  [16]bool
	devsByPriority  [16][]int
	interruptingDev [devMax]bool
)

// BusSendInterrupt triggers an IRQ for the given device
func BusSendInterrupt(devNum int) {
	interruptingDev[devNum] = true
	irqsByPriority[d[devNum].priorityMaskBit] = true
	IRQ = true
}

// BusClearInterrupt clears the IRQ for the given device
func BusClearInterrupt(devNum int) {
	interruptingDev[devNum] = false
	irqsByPriority[d[devNum].priorityMaskBit] = false
}

// BusGetHighestPriorityInt returns the device number of the highest priority device
// that has an outstanding interrupt
func BusGetHighestPriorityInt() (devNum int) {
	for p, i := range irqsByPriority {
		if i {
			for _, d := range devsByPriority[p] {
				if interruptingDev[d] {
					return d
				}
			}
		}
	}

	return 0 // ?
}

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

// BusAddDevice adds a new device to the system bus
func BusAddDevice(devMap DeviceMapT, devNum int, att bool) {
	if devNum >= devMax {
		log.Fatalf("ERROR: Attempt to add device with too-high device number: %#o", devNum)
	}
	d[devNum].devMu.Lock()
	d[devNum].mnemonic = devMap[devNum].DgMnemonic
	d[devNum].priorityMaskBit = devMap[devNum].PMB
	d[devNum].simAttached = att
	d[devNum].ioDevice = devMap[devNum].IsIO
	d[devNum].bootable = devMap[devNum].IsBootable
	// N.B. The relative priority of devs with the same PMB is established
	//      here by the order they are added to the bus
	if devMap[devNum].PMB <= 32 {
		devsByPriority[devMap[devNum].PMB] = append(devsByPriority[devMap[devNum].PMB], devNum)
	}
	logging.DebugPrint(logging.DebugLog, "INFO: Device %#o added to bus\n", devNum)
	d[devNum].devMu.Unlock()
}

// BusSetDataInFunc sets the Data In callback fror the given device
func BusSetDataInFunc(devNum int, fn DataInFunc) {
	d[devNum].devMu.Lock()
	d[devNum].dataInFunc = fn
	logging.DebugPrint(logging.DebugLog, "INFO: Bus Data In function set for dev %#o (%d.)\n", devNum, devNum)
	d[devNum].devMu.Unlock()
}

// BusDataIn forwards a Data In command to the given device
func BusDataIn(devNum int, abc byte, flag byte) (datum dg.WordT) {
	if d[devNum].dataInFunc == nil {
		log.Fatalf("ERROR: busDataIn called for device %#o with no function set", devNum)
	}
	//cpuPtr.ac[iPtr.acd] = dg.DwordT(d[iPtr.devNum].dataInFunc(abc, iPtr.f))
	return d[devNum].dataInFunc(abc, flag)
}

// BusSetDataOutFunc sets the Data Out callback fror the given device
func BusSetDataOutFunc(devNum int, fn DataOutFunc) {
	d[devNum].devMu.Lock()
	d[devNum].dataOutFunc = fn
	d[devNum].devMu.Unlock()
	logging.DebugPrint(logging.DebugLog, "INFO: Bus Data Out function set for dev %#o (%d.)\n", devNum, devNum)
}

// BusDataOut forwards a Data Out command to the given device
func BusDataOut(devNum int, datum dg.WordT, abc byte, flag byte) {
	if d[devNum].dataOutFunc == nil {
		logging.DebugLogsDump("logs/")
		log.Fatalf("ERROR: busDataOut called for device %#o with no function set", devNum)
	}
	d[devNum].dataOutFunc(datum, abc, flag)
}

// BusSetResetFunc sets the device reset callback for the given device
func BusSetResetFunc(devNum int, resetFn ResetFunc) {
	d[devNum].devMu.Lock()
	d[devNum].resetFunc = resetFn
	logging.DebugPrint(logging.DebugLog, "INFO: Bus reset function set for dev %#o\n", devNum)
	d[devNum].devMu.Unlock()
}

// BusResetDevice forwards a Reset to the given device
func BusResetDevice(devNum int) {
	d[devNum].devMu.RLock()
	doIt := d[devNum].ioDevice && (d[devNum].resetFunc != nil)
	d[devNum].devMu.RUnlock()
	if doIt {
		d[devNum].resetFunc()
	} else {
		log.Printf("INFO: Ignoring attempt to reset non-I/O/resetable device %#o\n", devNum)
	}

}

// BusResetAllIODevices calls the Reset func for each I/O device
func BusResetAllIODevices() {
	for dev := range d {
		d[dev].devMu.RLock()
		io := d[dev].ioDevice
		d[dev].devMu.RUnlock()
		if io {
			BusResetDevice(dev)
		}
	}
}

// BusSetAttached sets the MV/Em Attached flag for a device indicating the virtual
// device is attached to a host image file
func BusSetAttached(devNum int, imgName string) {
	d[devNum].devMu.Lock()
	d[devNum].simAttached = true
	d[devNum].simImageName = imgName
	d[devNum].devMu.Unlock()
}

// BusSetDetached clears the MV/Em Attached flag for a device
func BusSetDetached(devNum int) {
	d[devNum].devMu.Lock()
	d[devNum].simAttached = false
	d[devNum].simImageName = ""
	d[devNum].devMu.Unlock()
}

// BusIsAttached returns the attached state of the given device
func BusIsAttached(devNum int) bool {
	d[devNum].devMu.RLock()
	att := d[devNum].simAttached
	d[devNum].devMu.RUnlock()
	return att
}

// BusSetBusy sets/clears the given device's BUSY flag
func BusSetBusy(devNum int, f bool) {
	d[devNum].devMu.Lock()
	d[devNum].busy = f
	d[devNum].devMu.Unlock()
}

// BusSetDone sets/clears the given device#s DONE flag
func BusSetDone(devNum int, f bool) {
	d[devNum].devMu.Lock()
	d[devNum].done = f
	d[devNum].devMu.Unlock()
}

// BusGetBusy returns a device's BUSY flag
func BusGetBusy(devNum int) bool {
	d[devNum].devMu.RLock()
	bz := d[devNum].busy
	d[devNum].devMu.RUnlock()
	return bz
}

// BusGetDone returns a device's DONE flag
func BusGetDone(devNum int) bool {
	d[devNum].devMu.RLock()
	dn := d[devNum].done
	d[devNum].devMu.RUnlock()
	return dn
}

// BusIsBootable returns true if the device can be booted from
// this is not a guarantee that it WILL boot!
func BusIsBootable(devNum int) bool {
	d[devNum].devMu.RLock()
	bt := d[devNum].bootable
	d[devNum].devMu.RUnlock()
	return bt
}

// BusIsIODevice returns true if this is an IO device
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
	return memory.TestWbit(irqMask, int(d[devNum].priorityMaskBit))
}

// // BusSetDevMasked is a setter to make the device masked out from sending IRQs
// func BusSetDevMasked(devNum int) {
// 	memory.SetWbit(&irqMask, d[devNum].priorityMaskBit)
// }

// // BusClearDevMasked is a setter to make the device able to send IRQs
// func BusClearDevMasked(devNum int) {
// 	memory.ClearWbit(&irqMask, d[devNum].priorityMaskBit)
// }

// BusGetPrintableDevList is used by the console SHOW DEV command to display
// device statuses
func BusGetPrintableDevList() string {
	lst := fmt.Sprintf(" #  Mnem   PMB  I/O Busy Done Status\012")
	var line string
	for dev := range d {
		d[dev].devMu.RLock()
		if d[dev].mnemonic != "" {
			line = fmt.Sprintf("%#3o %-6s %2d. %3d %4d %4d  ",
				dev, d[dev].mnemonic, d[dev].priorityMaskBit,
				memory.BoolToInt(d[dev].ioDevice), memory.BoolToInt(d[dev].busy), memory.BoolToInt(d[dev].done))
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
