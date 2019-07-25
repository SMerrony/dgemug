// bus.go

// Copyright (C) 2018,2019  Steve Merrony

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
	// devMu           sync.RWMutex
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

// BusT holds the system bus and all its associated devices
type BusT struct {
	busMu           sync.RWMutex
	devices         [devMax]device
	irqMask         dg.WordT
	irq             bool
	irqsByPriority  [16]bool
	devsByPriority  [16][]int
	interruptingDev [devMax]bool
}

// SendInterrupt triggers an IRQ for the given device
func (bus *BusT) SendInterrupt(devNum int) {
	bus.busMu.Lock()
	bus.interruptingDev[devNum] = true
	bus.irqsByPriority[bus.devices[devNum].priorityMaskBit] = true
	bus.irq = true
	bus.busMu.Unlock()
}

// ClearInterrupt clears the IRQ for the given device
func (bus *BusT) ClearInterrupt(devNum int) {
	bus.busMu.Lock()
	bus.interruptingDev[devNum] = false
	bus.irqsByPriority[bus.devices[devNum].priorityMaskBit] = false
	bus.busMu.Unlock()
}

// GetHighestPriorityInt returns the device number of the highest priority device
// that has an outstanding interrupt
func (bus *BusT) GetHighestPriorityInt() (devNum int) {
	bus.busMu.RLock()
	defer bus.busMu.RUnlock()
	for p, i := range bus.irqsByPriority {
		if i {
			for _, d := range bus.devsByPriority[p] {
				if bus.interruptingDev[d] {
					return d
				}
			}
		}
	}
	return 0 // ?
}

// BusInit must be called before attaching any devices
func (bus *BusT) BusInit() {
	bus.busMu.Lock()
	defer bus.busMu.Unlock()
	for dev := range bus.devices {
		bus.devices[dev].mnemonic = ""
		bus.devices[dev].priorityMaskBit = 0
		bus.devices[dev].dataInFunc = nil
		bus.devices[dev].dataOutFunc = nil
		bus.devices[dev].simAttached = false
		bus.devices[dev].ioDevice = false
		bus.devices[dev].bootable = false
		bus.devices[dev].busy = false
		bus.devices[dev].done = false
	}
}

// AddDevice adds a new device to the system bus
func (bus *BusT) AddDevice(devMap DeviceMapT, devNum int, att bool) {
	if devNum >= devMax {
		log.Fatalf("ERROR: Attempt to add device with too-high device number: %#o", devNum)
	}
	bus.busMu.Lock()
	bus.devices[devNum].mnemonic = devMap[devNum].DgMnemonic
	bus.devices[devNum].priorityMaskBit = devMap[devNum].PMB
	bus.devices[devNum].simAttached = att
	bus.devices[devNum].ioDevice = devMap[devNum].IsIO
	bus.devices[devNum].bootable = devMap[devNum].IsBootable
	// N.B. The relative priority of devs with the same PMB is established
	//      here by the order they are added to the bus
	if devMap[devNum].PMB <= 32 {
		bus.devsByPriority[devMap[devNum].PMB] = append(bus.devsByPriority[devMap[devNum].PMB], devNum)
	}
	logging.DebugPrint(logging.DebugLog, "INFO: Device %#o added to bus\n", devNum)
	bus.busMu.Unlock()
}

// SetDataInFunc sets the Data In callback fror the given device
func (bus *BusT) SetDataInFunc(devNum int, fn DataInFunc) {
	bus.busMu.Lock()
	bus.devices[devNum].dataInFunc = fn
	logging.DebugPrint(logging.DebugLog, "INFO: Bus Data In function set for dev %#o (%d.)\n", devNum, devNum)
	bus.busMu.Unlock()
}

// DataIn forwards a Data In command to the given device
func (bus *BusT) DataIn(devNum int, abc byte, flag byte) (datum dg.WordT) {
	// bus.busMu.RLock()
	if bus.devices[devNum].dataInFunc == nil {
		log.Fatalf("ERROR: busDataIn called for device %#o with no function set", devNum)
	}
	// logging.DebugPrint(logging.DebugLog, "Data In Command: DI-%c to device %#o\n", abc, devNum)
	wd := bus.devices[devNum].dataInFunc(abc, flag)
	// bus.busMu.RUnlock()
	return wd
}

// SetDataOutFunc sets the Data Out callback fror the given device
func (bus *BusT) SetDataOutFunc(devNum int, fn DataOutFunc) {
	bus.busMu.Lock()
	bus.devices[devNum].dataOutFunc = fn
	bus.busMu.Unlock()
	logging.DebugPrint(logging.DebugLog, "INFO: Bus Data Out function set for dev %#o (%d.)\n", devNum, devNum)
}

// DataOut forwards a Data Out command to the given device
func (bus *BusT) DataOut(devNum int, datum dg.WordT, abc byte, flag byte) {
	// bus.busMu.RLock()
	if bus.devices[devNum].dataOutFunc == nil {
		logging.DebugLogsDump("logs/")
		log.Fatalf("ERROR: busDataOut called for device %#o with no function set", devNum)
	}
	bus.devices[devNum].dataOutFunc(datum, abc, flag)
	// bus.busMu.RUnlock()
}

// SetResetFunc sets the device reset callback for the given device
func (bus *BusT) SetResetFunc(devNum int, resetFn ResetFunc) {
	bus.busMu.Lock()
	bus.devices[devNum].resetFunc = resetFn
	logging.DebugPrint(logging.DebugLog, "INFO: Bus reset function set for dev %#o\n", devNum)
	bus.busMu.Unlock()
}

// ResetDevice forwards a Reset to the given device
func (bus *BusT) ResetDevice(devNum int) {
	// bus.busMu.RLock()
	if bus.devices[devNum].ioDevice && (bus.devices[devNum].resetFunc != nil) {
		bus.devices[devNum].resetFunc()
	} else {
		log.Printf("INFO: Ignoring attempt to reset non-I/O/resetable device %#o\n", devNum)
	}
	// bus.busMu.RUnlock()
}

// ResetAllIODevices calls the Reset func for each I/O device
func (bus *BusT) ResetAllIODevices() {
	bus.busMu.RLock()
	for dev := range bus.devices {
		if bus.devices[dev].ioDevice {
			bus.ResetDevice(dev)
		}
	}
	bus.busMu.RUnlock()
}

// SetAttached sets the MV/Em Attached flag for a device indicating the virtual
// device is attached to a host image file
func (bus *BusT) SetAttached(devNum int, imgName string) {
	bus.busMu.Lock()
	bus.devices[devNum].simAttached = true
	bus.devices[devNum].simImageName = imgName
	bus.busMu.Unlock()
}

// SetDetached clears the MV/Em Attached flag for a device
func (bus *BusT) SetDetached(devNum int) {
	bus.busMu.Lock()
	bus.devices[devNum].simAttached = false
	bus.devices[devNum].simImageName = ""
	bus.busMu.Unlock()
}

// IsAttached returns the attached state of the given device
func (bus *BusT) IsAttached(devNum int) bool {
	bus.busMu.RLock()
	att := bus.devices[devNum].simAttached
	bus.busMu.RUnlock()
	return att
}

// SetBusy sets/clears the given device's BUSY flag
func (bus *BusT) SetBusy(devNum int, f bool) {
	bus.busMu.Lock()
	bus.devices[devNum].busy = f
	bus.busMu.Unlock()
}

// SetDone sets/clears the given device#s DONE flag
func (bus *BusT) SetDone(devNum int, f bool) {
	bus.busMu.Lock()
	bus.devices[devNum].done = f
	bus.busMu.Unlock()
}

// GetBusy returns a device's BUSY flag
func (bus *BusT) GetBusy(devNum int) bool {
	bus.busMu.RLock()
	bz := bus.devices[devNum].busy
	bus.busMu.RUnlock()
	return bz
}

// GetDone returns a device's DONE flag
func (bus *BusT) GetDone(devNum int) bool {
	bus.busMu.RLock()
	dn := bus.devices[devNum].done
	bus.busMu.RUnlock()
	return dn
}

// IsBootable returns true if the device can be booted from
// It is is not a guarantee that it WILL boot!
func (bus *BusT) IsBootable(devNum int) bool {
	bus.busMu.RLock()
	bt := bus.devices[devNum].bootable
	bus.busMu.RUnlock()
	return bt
}

// IsIODevice returns true if this is an IO device
func (bus *BusT) IsIODevice(devNum int) bool {
	bus.busMu.RLock()
	io := bus.devices[devNum].ioDevice
	bus.busMu.RUnlock()
	return io
}

// SetIrqMask is a setter for the (whole) IRQ mask
func (bus *BusT) SetIrqMask(newMask dg.WordT) {
	bus.busMu.Lock()
	bus.irqMask = newMask
	bus.busMu.Unlock()
}

// IsDevMasked is a getter to see if the device is masked out from sending IRQs
func (bus *BusT) IsDevMasked(devNum int) (masked bool) {
	return memory.TestWbit(bus.irqMask, int(bus.devices[devNum].priorityMaskBit))
}

// GetIRQ is a getter for IRQ
func (bus *BusT) GetIRQ() bool {
	bus.busMu.RLock()
	defer bus.busMu.RUnlock()
	return bus.irq
}

// SetIRQ is a getter for IRQ
func (bus *BusT) SetIRQ(i bool) {
	bus.busMu.Lock()
	bus.irq = i
	bus.busMu.Unlock()
}

// BusSetDevMasked is a setter to make the device masked out from sending IRQs
// func (bus *BusT) SetDevMasked(devNum int) {
// 	memory.SetWbit(&irqMask, bus.devices[devNum].priorityMaskBit)
// }

// BusClearDevMasked is a setter to make the device able to send IRQs
// func (bus *BusT) ClearDevMasked(devNum int) {
// 	memory.ClearWbit(&irqMask, bus.devices[devNum].priorityMaskBit)
// }

// GetPrintableDevList is used by the console SHOW DEV command to display
// device statuses
func (bus *BusT) GetPrintableDevList() string {
	lst := fmt.Sprintf(" #  Mnem   PMB  I/O Busy Done Status\012")
	var line string
	bus.busMu.RLock()
	for dev := range bus.devices {
		if bus.devices[dev].mnemonic != "" {
			line = fmt.Sprintf("%#3o %-6s %2d. %3d %4d %4d  ",
				dev, bus.devices[dev].mnemonic, bus.devices[dev].priorityMaskBit,
				memory.BoolToInt(bus.devices[dev].ioDevice), memory.BoolToInt(bus.devices[dev].busy), memory.BoolToInt(bus.devices[dev].done))
			if bus.devices[dev].simAttached {
				line += "Attached"
				if bus.devices[dev].simImageName != "" {
					line += " to image: " + bus.devices[dev].simImageName
				}
			} else {
				line += "Not Attached"
			}
			// Commented out the below for now as it's a bit misleading...
			// if bus.devices[dev].bootable {
			// 	line += ", Bootable"
			// }
			line += "\012"
			lst += line
		}
	}
	bus.busMu.RUnlock()
	return lst
}
