// disk.go

// Copyright (C) 2019  Steve Merrony

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

// Here we are emulating the disk4231a device, specifically model 4231a
// controller/drive combination which provides c.190MB of formatted capacity.

package devices

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/SMerrony/dgemug/dg"
	"github.com/SMerrony/dgemug/logging"
	"github.com/SMerrony/dgemug/memory"
)

// Physical characteristics of the emulated disk
const (
	disk4231aSurfPerDisk  = 19
	disk4231aSectPerTrack = 23
	disk4231aWordsPerSect = 256
	disk4231aBytesPerSect = disk4231aWordsPerSect * 2
	disk4231aPhysCyls     = 411
	disk4231aPhysByteSize = disk4231aSurfPerDisk * disk4231aSectPerTrack * disk4231aBytesPerSect * disk4231aPhysCyls
)

const (
	disk4231aCmdRead = iota
	disk4231aCmdWrite
	disk4231aCmdSeek
	disk4231aCmdRecal
)

const (
	// drive statuses
	disk4231aStatusError = 1 << iota
	disk4231aStatusDataLate
	disk4231aStatusCheckError
	disk4231aStatusUnsafe
	disk4231aStatusEndError
	disk4231aStatusSeekError
	disk4231aStatusDiscReady
	disk4231aStatusAddressError
	disk4231aStatusHeadError
	disk4231aStatusSectorError
	disk4231aStatusDualProcessor
	disk4231aStatusDrv3Done
	disk4231aStatusDrv2Done
	disk4231aStatusDrv1Done
	disk4231aStatusDrv0Done
	disk4231aStatusDPDone
)

// disk4231aStatsPeriodMs is the number of milliseconds between sending status updates
const disk4231aStatsPeriodMs = 333

// Disk4231aT holds the current state of the Type-6231A Moving-Head Disk
type Disk4231aT struct {
	// MV/Em internals...
	bus                 *BusT
	ImageAttached       bool
	disk4231aMu         sync.RWMutex
	devNum              int
	logID               int
	debugLogging        bool
	imageFileName       string
	imageFile           *os.File
	reads, writes       uint64
	readBuff, writeBuff []byte
	cmdDecode           [4]string
	// DG data...
	cmdDrvAddr byte     // 6-bit?
	command    int8     // 4-bit
	drive      uint8    // 2-bit
	mapEnabled bool     // is the BMC addressing physical (0) or Mapped (1)
	memAddr    dg.WordT // 15-bit self-incrementing on DG
	statusReg  dg.WordT
	ema        uint8    // 5-bit
	cylinder   dg.WordT // 10-bit
	surface    uint8    // 5-bit - increments post-op
	sector     uint8    // 5-bit - increments mid-op
	sectCnt    int8     // 4-bit - 2's complement of secs remaining, incrememts mid-op - signed
}

// Disk4231aStatT holds the data reported to the status collector
type Disk4231aStatT struct {
	ImageAttached bool
	Cylinder      dg.WordT
	Head, Sector  uint8
	Reads, Writes uint64
}

// Disk4231aInit must be called to initialise the emulated disk4231a controller
func (disk *Disk4231aT) Disk4231aInit(dev int, bus *BusT, statsChann chan Disk4231aStatT, logID int, logging bool) {
	disk.disk4231aMu.Lock()
	defer disk.disk4231aMu.Unlock()
	disk.bus = bus
	disk.devNum = dev
	disk.logID = logID
	disk.debugLogging = logging

	disk.cmdDecode = [...]string{"READ", "WRITE", "SEEK", "RECAL"}

	go disk.disk4231aStatsSender(statsChann)

	bus.SetResetFunc(disk.devNum, disk.disk4231aReset)
	bus.SetDataInFunc(disk.devNum, disk.disk4231aDataIn)
	bus.SetDataOutFunc(disk.devNum, disk.disk4231aDataOut)
	disk.ImageAttached = false
	disk.statusReg = disk4231aStatusDiscReady
	disk.mapEnabled = false
	disk.readBuff = make([]byte, disk4231aBytesPerSect)
	disk.writeBuff = make([]byte, disk4231aBytesPerSect)
}

// Disk4231aAttach attempts to attach an extant 4231a disk image to the running emulator
func (disk *Disk4231aT) Disk4231aAttach(dNum int, imgName string) bool {
	// TODO Disk Number not currently used
	logging.DebugPrint(disk.logID, "disk4231aAttach called for disk #%d with image <%s>\n", dNum, imgName)
	disk.disk4231aMu.Lock()
	var err error
	disk.imageFile, err = os.OpenFile(imgName, os.O_RDWR, 0755)
	if err != nil {
		logging.DebugPrint(disk.logID, "Failed to open image for attaching\n")
		logging.DebugPrint(logging.DebugLog, "WARN: Failed to open disk4231a image <%s> for ATTach\n", imgName)
		disk.disk4231aMu.Unlock()
		return false
	}
	disk.imageFileName = imgName
	disk.ImageAttached = true
	disk.disk4231aMu.Unlock()
	disk.bus.SetAttached(disk.devNum, imgName)
	return true
}

func (disk *Disk4231aT) disk4231aStatsSender(sChan chan Disk4231aStatT) {
	var stats Disk4231aStatT
	for {
		disk.disk4231aMu.RLock()
		if disk.ImageAttached {
			stats.ImageAttached = true
			stats.Cylinder = disk.cylinder
			stats.Head = disk.surface
			stats.Sector = disk.sector
			stats.Reads = disk.reads
			stats.Writes = disk.writes
		} else {
			stats = Disk4231aStatT{}
		}
		disk.disk4231aMu.RUnlock()
		select {
		case sChan <- stats:
		default:
		}
		time.Sleep(time.Millisecond * disk4231aStatsPeriodMs)
	}
}

// Disk4231aCreateBlank creates an empty disk file of the correct size for the disk4231a emulator to use
func (disk *Disk4231aT) Disk4231aCreateBlank(imgName string) bool {
	newFile, err := os.Create(imgName)
	if err != nil {
		return false
	}
	defer newFile.Close()
	logging.DebugPrint(disk.logID, "disk4231aCreateBlank attempting to write %d bytes\n", disk4231aPhysByteSize)
	w := bufio.NewWriter(newFile)
	for b := 0; b < disk4231aPhysByteSize; b++ {
		w.WriteByte(0)
	}
	w.Flush()
	return true
}

// Disk4231aLoadDKBT - This func mimics a system ROM routine to boot from disk.
// Rather than copying a ROM routine (!) we simply mimic its basic actions...
// Load 1st two blocks from disk into location 0
func (disk *Disk4231aT) Disk4231aLoadDKBT() {
	logging.DebugPrint(disk.logID, "Disk6961LoadDKBT() called\n")
	// set posn
	disk.command = disk4231aCmdRecal
	disk.disk4231aDoCommand()
	disk.memAddr = 0
	disk.sectCnt = -2
	disk.command = disk4231aCmdRead
	disk.disk4231aDoCommand()
	logging.DebugPrint(disk.logID, "Disk6961LoadDKBT() completed\n")
}

// disk4231aDataIn implements the DIA/B/C I/O instructions for this device
func (disk *Disk4231aT) disk4231aDataIn(abc byte, flag byte) (data dg.WordT) {
	disk.disk4231aMu.RLock()
	switch abc {
	case 'A': // Read Status
		data = disk.statusReg
		if disk.debugLogging {
			logging.DebugPrint(disk.logID, "DIA [Read Status] returning %s for DRV=%d\n",
				memory.WordToBinStr(data), disk.drive)
		}
	case 'B': // Read Memory Address Counter
		data = disk.memAddr & 0x7fff
		if disk.debugLogging {
			logging.DebugPrint(disk.logID, "DIB [Read Memory Address] returning %s for DRV=%d\n", memory.WordToBinStr(data), disk.drive)
		}
	case 'C': // Read Disc Address
		data = dg.WordT(disk.drive) << 14
		data |= (dg.WordT(disk.surface) & 0x1f) << 9
		data |= (dg.WordT(disk.sector) & 0x1f) << 4
		data |= dg.WordT(disk.sectCnt) & 0x0f
		if disk.debugLogging {
			logging.DebugPrint(disk.logID, "DIC [Read Disc Address] returning %s\n", memory.WordToBinStr(data))
		}
	}
	disk.disk4231aMu.RUnlock()

	disk.disk4231aHandleFlag(flag)

	return data
}

// disk4231aDataOut implements the DOA/B/C instructions for this device
// NIO is also routed here with a dummy abc flag value of N
func (disk *Disk4231aT) disk4231aDataOut(datum dg.WordT, abc byte, flag byte) {
	disk.disk4231aMu.Lock()
	switch abc {
	case 'A': // Specify Command and Cylinder
		// disk.ccsReg = datum
		disk.command = extractDisk4231aCommand(datum)
		disk.cylinder = extractDisk4231aCylinder(datum)
		if disk.debugLogging {
			logging.DebugPrint(disk.logID, "DOA [Specify Command & Cylinder] to DRV=%d with data %s, CMD: %d, CYL: %d.\n",
				disk.drive, memory.WordToBinStr(datum), disk.command, disk.cylinder)
		}
		if memory.TestWbit(datum, 0) {
			disk.statusReg &= ^dg.WordT(disk4231aStatusDPDone)
			disk.statusReg &= ^dg.WordT(disk4231aStatusAddressError)
			disk.statusReg &= ^dg.WordT(disk4231aStatusEndError)
			disk.statusReg &= ^dg.WordT(disk4231aStatusCheckError)
			disk.statusReg &= ^dg.WordT(disk4231aStatusSectorError)
		}
		if memory.TestWbit(datum, 1) {
			disk.statusReg &= ^dg.WordT(disk4231aStatusDrv0Done)
		}
		if memory.TestWbit(datum, 2) {
			disk.statusReg &= ^dg.WordT(disk4231aStatusDrv1Done)
		}
		if memory.TestWbit(datum, 3) {
			disk.statusReg &= ^dg.WordT(disk4231aStatusDrv2Done)
		}
		if memory.TestWbit(datum, 4) {
			disk.statusReg &= ^dg.WordT(disk4231aStatusDrv3Done)
		}
	case 'B': // Load Memory Address Counter
		if memory.TestWbit(datum, 0) {
			log.Fatalln("Not Yet Implemented: Disk 4231A Format Mode")
		}
		disk.memAddr = datum & 0x7fff
		if disk.debugLogging {
			logging.DebugPrint(disk.logID, "DOB [Load Memory Addr] with data %s\n",
				memory.WordToBinStr(datum))
			logging.DebugPrint(disk.logID, "... MEM Addr: %#o\n", disk.memAddr)
		}
	case 'C': // Specify Disc Address and Sector Count
		// disk.ccsReg = datum
		disk.drive = extractDisk4231aDriveNo(datum)
		disk.surface = extractDisk4231aSurface(datum)
		disk.sector = extractDisk4231aSector(datum)
		disk.sectCnt = extractDisk4231aSectCnt(datum)
		if disk.debugLogging {
			logging.DebugPrint(disk.logID, "DOC [Specify Disk Addr & Sect Cnt with data %s\n",
				memory.WordToBinStr(datum))
			logging.DebugPrint(disk.logID, "... DRV: %d, SURF: %d., SECT: %d., SECCNT: %d.\n",
				disk.drive, disk.surface, disk.sector, disk.sectCnt)
		}
	case 'N': // dummy value for NIO - we just handle the flag below
		if disk.debugLogging {
			logging.DebugPrint(disk.logID, "NIO%c received\n", flag)
		}
	}
	disk.disk4231aMu.Unlock()

	disk.disk4231aHandleFlag(flag)
}

func (disk *Disk4231aT) disk4231aDoCommand() {
	var (
		bytesRead, bytesWritten int
		wd, wIx                 dg.WordT
		err                     error
	)

	disk.disk4231aMu.Lock()

	switch disk.command {
	// RECALibrate (goto pos. 0)
	case disk4231aCmdRecal:
		disk.cylinder = 0
		disk.surface = 0
		disk.disk4231aPositionDiskImage()
		disk.statusReg = disk4231aStatusDiscReady | disk4231aStatusDrv0Done | disk4231aStatusDPDone
		if disk.debugLogging {
			logging.DebugPrint(disk.logID, "... RECAL done, %s\n", disk.disk4231aPrintableAddr())
		}

	// SEEK
	case disk4231aCmdSeek:
		// action the seek
		disk.disk4231aPositionDiskImage()
		disk.statusReg = disk4231aStatusDiscReady | disk4231aStatusDrv0Done | disk4231aStatusDPDone
		if disk.debugLogging {
			logging.DebugPrint(disk.logID, "... SEEK done, %s\n", disk.disk4231aPrintableAddr())
		}

	// ===== READ from disk4231a =====
	case disk4231aCmdRead:
		if disk.debugLogging {
			logging.DebugPrint(disk.logID, "... READ command invoked %s\n", disk.disk4231aPrintableAddr())
			logging.DebugPrint(disk.logID, "... .... Start Address: %#o\n", disk.memAddr)
		}

		for disk.sectCnt != 0 {
			// check CYL
			if disk.cylinder >= disk4231aPhysCyls {
				disk.statusReg = disk4231aStatusDiscReady | disk4231aStatusAddressError | disk4231aStatusError
				disk.disk4231aMu.Unlock()
				return
			}
			// check SECT
			if disk.sector >= disk4231aSectPerTrack {
				disk.sector = 0
				disk.surface++
				if disk.debugLogging {
					logging.DebugPrint(disk.logID, "Sector read overflow, advancing to surface %d.", disk.surface)
				}
			}
			// check SURF (head)
			if disk.surface >= disk4231aSurfPerDisk {
				disk.statusReg = disk4231aStatusDiscReady | disk4231aStatusAddressError | disk4231aStatusHeadError | disk4231aStatusError
				disk.disk4231aMu.Unlock()
				return
			}
			disk.disk4231aPositionDiskImage()
			bytesRead, err = disk.imageFile.Read(disk.readBuff)

			if bytesRead != disk4231aBytesPerSect || err != nil {
				log.Fatalf("ERROR: unexpected return from disk4231a Image File Read: %s", err)
			}
			for wIx = 0; wIx < disk4231aWordsPerSect; wIx++ {
				wd = (dg.WordT(disk.readBuff[(wIx*2)+1]) << 8) | dg.WordT(disk.readBuff[wIx*2])
				memory.WriteWordBmcChan16bit(&disk.memAddr, wd)
			}
			disk.sector++
			disk.sectCnt++
			disk.reads++

			if disk.debugLogging {
				logging.DebugPrint(disk.logID, "Buffer: %X\n", disk.readBuff)
			}

		}
		if disk.debugLogging {
			logging.DebugPrint(disk.logID, "... .... READ command finished %s\n", disk.disk4231aPrintableAddr())
			logging.DebugPrint(disk.logID, "\n... .... Last Address: %#o\n", disk.memAddr)
		}
		disk.statusReg = disk4231aStatusDiscReady | disk4231aStatusDPDone

		// case disk4231aCmdRelease:
	// 	// I think this is a NOP on a single-processor machine

	case disk4231aCmdWrite:
		if disk.debugLogging {
			logging.DebugPrint(disk.logID, "... WRITE command invoked %s\n", disk.disk4231aPrintableAddr())
			logging.DebugPrint(disk.logID, "... .....  Start Address: %#o\n", disk.memAddr)
		}

		for disk.sectCnt != 0 {
			// check CYL
			if disk.cylinder >= disk4231aPhysCyls {
				disk.statusReg = disk4231aStatusDiscReady | disk4231aStatusAddressError | disk4231aStatusError
				if disk.debugLogging {
					logging.DebugPrint(disk.logID, "Cylinder overflow: %d.", disk.cylinder)
				}
				disk.disk4231aMu.Unlock()
				return
			}
			// check SECT
			if disk.sector >= disk4231aSectPerTrack {
				disk.sector = 0
				disk.surface++
				if disk.debugLogging {
					logging.DebugPrint(disk.logID, "Sector write overflow, advancing to surface %d.", disk.surface)
				}
			}
			// check SURF (head)
			if disk.surface >= disk4231aSurfPerDisk {
				disk.statusReg = disk4231aStatusDiscReady | disk4231aStatusAddressError | disk4231aStatusHeadError | disk4231aStatusError
				if disk.debugLogging {
					logging.DebugPrint(disk.logID, "Surface overflow: %d.", disk.surface)
				}
				disk.disk4231aMu.Unlock()
				return
			}
			disk.disk4231aPositionDiskImage()
			for wIx = 0; wIx < disk4231aWordsPerSect; wIx++ {
				wd = memory.ReadWordBmcChan16bit(&disk.memAddr)
				disk.writeBuff[(wIx*2)+1] = byte(wd >> 8)
				disk.writeBuff[wIx*2] = byte(wd)
			}
			bytesWritten, err = disk.imageFile.Write(disk.writeBuff)
			if bytesWritten != disk4231aBytesPerSect || err != nil {
				log.Fatalf("ERROR: unexpected return from disk4231a Image File Write: %s", err)
			}
			disk.sector++
			disk.sectCnt++
			disk.writes++

			if disk.debugLogging {
				logging.DebugPrint(disk.logID, "Buffer: %X\n", disk.writeBuff)
			}
		}
		if disk.debugLogging {
			logging.DebugPrint(disk.logID, "... ..... WRITE command finished %s\n", disk.disk4231aPrintableAddr())
			logging.DebugPrint(disk.logID, "... ..... Last Address: %#o\n", disk.memAddr)
		}

		disk.statusReg = disk4231aStatusDiscReady | disk4231aStatusDPDone

	default:
		log.Fatalf("disk4231a Disk R/W Command %d not yet implemented\n", disk.command)
	}
	disk.disk4231aMu.Unlock()
}

func (disk *Disk4231aT) disk4231aHandleFlag(f byte) {
	switch f {
	case 'S': // initiate Read, Write, Seek or Recalibrate
		disk.bus.SetBusy(disk.devNum, true)
		disk.bus.SetDone(disk.devNum, false)
		// TODO stop any I/O
		// disk.disk4231aMu.Lock()
		// TODO start I/O timeout
		if disk.debugLogging {
			logging.DebugPrint(disk.logID, "... S flag set\n")
		}
		// disk.disk4231aMu.Unlock()
		disk.disk4231aDoCommand()
		disk.bus.SetBusy(disk.devNum, false)
		disk.bus.SetDone(disk.devNum, true)
		// send IRQ if not masked out
		//if !BusIsDevMasked(disk.devNum) {
		// InterruptingDev[disk.devNum] = true
		// IRQ = true
		disk.bus.SendInterrupt(disk.devNum)
		//}

	case 'C': // stop all positioning and txfer ops
		disk.bus.SetBusy(disk.devNum, false)
		disk.bus.SetDone(disk.devNum, false)
		if disk.debugLogging {
			logging.DebugPrint(disk.logID, "... C flag set\n")
		}
		disk.disk4231aMu.Lock()
		disk.statusReg = 0
		disk.disk4231aMu.Unlock()
		//disk.bus.SendInterrupt(disk.devNum)
		// send IRQ if not masked out
		//if !BusIsDevMasked(disk.devNum) {
		// InterruptingDev[disk.devNum] = true
		// IRQ = true
	//	disk.bus.SendInterrupt(disk.devNum)
	//}

	case 'P': // Initiate a Seek or Recalibrate operation
		disk.bus.SetBusy(disk.devNum, false)
		disk.disk4231aMu.Lock()
		if disk.debugLogging {
			logging.DebugPrint(disk.logID, "... P flag set\n")
		}
		disk.statusReg = 0
		disk.disk4231aMu.Unlock()
		disk.disk4231aDoCommand()
		//disk.rwStatus = disk4231aDrive0Done
		//BusSetBusy(disk.devNum, false)
		//BusSetDone(disk.devNum, true)
		// send IRQ if not masked out
		//if !BusIsDevMasked(disk.devNum) {
		// InterruptingDev[disk.devNum] = true
		// IRQ = true
		disk.bus.SendInterrupt(disk.devNum)
		//}

	default:
		// no/empty flag - nothing to do
	}
}

// set the MV/Em disk image file postion according to current C/H/S
func (disk *Disk4231aT) disk4231aPositionDiskImage() {
	//lba = ((int64(disk.cylinder*disk4231aSurfPerDisk) + int64(disk.surface)) * int64(disk4231aSectPerTrack)) + int64(disk.sector)
	offset := (((int64(disk.cylinder*disk4231aSurfPerDisk) + int64(disk.surface)) * int64(disk4231aSectPerTrack)) + int64(disk.sector)) * disk4231aBytesPerSect
	r, err := disk.imageFile.Seek(offset, 0)
	if r != offset || err != nil {
		log.Fatal("disk4231a could not postition disk image via seek()")
	}
}

func (disk *Disk4231aT) disk4231aPrintableAddr() string {
	// MUST BE LOCKED BY CALLER
	pa := fmt.Sprintf("DRV: %d, CYL: %d, SURF: %d, SECT: %d, SECCNT: %d",
		disk.drive, disk.cylinder,
		disk.surface, disk.sector, disk.sectCnt)
	return pa
}

// reset the disk4231a controller
func (disk *Disk4231aT) disk4231aReset() {
	disk.disk4231aMu.Lock()
	disk.command = disk4231aCmdRead
	disk.cylinder = 0
	disk.surface = 0
	disk.sector = 0
	disk.sectCnt = 0
	disk.statusReg = disk4231aStatusDiscReady
	if disk.debugLogging {
		logging.DebugPrint(disk.logID, "disk4231a Reset\n")
	}
	disk.disk4231aMu.Unlock()
}

func extractDisk4231aCommand(word dg.WordT) int8 {
	return int8((word >> 9) & 0x03)
}

func extractDisk4231aCylinder(word dg.WordT) dg.WordT {
	return word & 0x1ff
}

func extractDisk4231aDriveNo(word dg.WordT) uint8 {
	return uint8((word & 0x03) >> 14)
}

func extractDisk4231aSector(word dg.WordT) uint8 {
	return uint8((word & 0x01f0) >> 4)
}

func extractDisk4231aSectCnt(word dg.WordT) int8 {
	tmpWd := word & 0x0f
	if tmpWd != 0 { // sign-extend
		tmpWd |= 0xf0
	}
	return int8(tmpWd)
}

func extractDisk4231aSurface(word dg.WordT) uint8 {
	return uint8((word & 0x3e00) >> 9)
}
