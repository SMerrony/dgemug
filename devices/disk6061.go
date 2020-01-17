// disk.go

// Copyright (C) 2017,2018,2019 Steve Merrony

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

// Here we are emulating the disk6061 device, specifically model 6061
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
	disk6061SurfPerDisk  = 19 //5 // 19
	disk6061SectPerTrack = 24
	disk6061WordsPerSect = 256
	disk6061BytesPerSect = disk6061WordsPerSect * 2
	disk6061PhysCyls     = 815
	disk6061PhysByteSize = disk6061SurfPerDisk * disk6061SectPerTrack * disk6061BytesPerSect * disk6061PhysCyls
)

const (
	disk6061CmdRead = iota
	disk6061CmdRecal
	disk6061CmdSeek
	disk6061CmdStop
	disk6061CmdOffsetFwd
	disk6061CmdOffsetRev
	disk6061CmdWriteDisable
	disk6061CmdRelease
	disk6061CmdTrespass
	disk6061CmdSetAltMode1
	disk6061CmdSetAltMode2
	disk6061CmdNoOp
	disk6061CmdVerify
	disk6061CmdReadBuffs
	disk6061CmdWrite
	disk6061CmdFormat
)
const (
	disk6061InsModeNormal = iota
	disk6061InsModeAlt1
	disk6061InsModeAlt2
)
const (
	// drive statuses
	disk6061Drivefault = 1 << iota
	disk6061Writefault
	disk6061Clockfault
	disk6061Posnfault
	disk6061Packunsafe
	disk6061Powerfault
	disk6061Illegalcmd
	disk6061Invalidaddr
	disk6061Unused
	disk6061Writedis
	disk6061Offset
	disk6061Busy
	disk6061Ready
	disk6061Trespassed
	disk6061Reserved
	disk6061Invalid
)
const (
	// R/W statuses
	disk6061Rwfault = 1 << iota
	disk6061late
	disk6061Rwtimeout
	disk6061Verify
	disk6061Surfsect
	disk6061Cylinder
	disk6061Badsector
	disk6061Ecc
	disk6061Illegalsector
	disk6061Parity
	disk6061Drive3Done
	disk6061Drive2Done
	disk6061Drive1Done
	disk6061Drive0Done
	disk6061Rwdone
	disk6061Controlfull
)

// disk6061StatsPeriodMs is the number of milliseconds between sending status updates
const disk6061StatsPeriodMs = 500

// Disk6061T holds the current state of a Type 6061 Moving-Head Disk
type Disk6061T struct {
	// MV/Em internals...
	bus                 *BusT
	ImageAttached       bool
	disk6061Mu          sync.RWMutex
	devNum              int
	logID               int
	imageFileName       string
	imageFile           *os.File
	reads, writes       uint64
	readBuff, writeBuff []byte
	cmdDecode           [disk6061CmdFormat + 1]string
	debugLogging        bool
	// DG data...
	cmdDrvAddr      byte     // 6-bit?
	command         int8     // 4-bit
	drive           uint8    // 2-bit
	mapEnabled      bool     // is the BMC addressing physical (0) or Mapped (1)
	memAddr         dg.WordT // self-incrementing on DG
	ema             uint8    // 5-bit
	cylinder        dg.WordT // 10-bit
	surface         uint8    // 5-bit - increments post-op
	sector          uint8    // 5-bit - increments mid-op
	sectCnt         int8     // 5-bit - incrememts mid-op - signed
	ecc             dg.DwordT
	driveStatus     dg.WordT
	rwStatus        dg.WordT
	instructionMode int
	lastDOAwasSeek  bool
}

// Disk6061StatT holds the data reported to the status collector
type Disk6061StatT struct {
	ImageAttached bool
	Cylinder      dg.WordT
	Head, Sector  uint8
	Reads, Writes uint64
}

// Disk6061Init must be called to initialise the emulated disk6061 controller
func (disk *Disk6061T) Disk6061Init(dev int, bus *BusT, statsChann chan Disk6061StatT, logID int, logging bool) {
	disk.disk6061Mu.Lock()
	defer disk.disk6061Mu.Unlock()
	disk.devNum = dev
	disk.bus = bus
	disk.logID = logID
	disk.debugLogging = logging

	disk.cmdDecode = [...]string{"READ", "RECAL", "SEEK", "STOP", "OFFSET FWD", "OFFSET REV",
		"WRITE DISABLE", "RELEASE", "TRESPASS", "SET ALT MODE 1", "SET ALT MODE 2",
		"NO OP", "VERIFY", "READ BUFFERS", "WRITE", "FORMAT"}

	go disk.disk6061StatsSender(statsChann)

	bus.SetResetFunc(disk.devNum, disk.disk6061Reset)
	bus.SetDataInFunc(disk.devNum, disk.disk6061In)
	bus.SetDataOutFunc(disk.devNum, disk.disk6061Out)
	disk.ImageAttached = false
	disk.instructionMode = disk6061InsModeNormal
	disk.driveStatus = disk6061Ready
	disk.mapEnabled = false
	disk.readBuff = make([]byte, disk6061BytesPerSect)
	disk.writeBuff = make([]byte, disk6061BytesPerSect)
}

// Disk6061Attach attempts to attach an extant 6061 disk image to the running emulator
func (disk *Disk6061T) Disk6061Attach(dNum int, imgName string) bool {
	// TODO Disk Number not currently used
	logging.DebugPrint(disk.logID, "disk6061Attach called for disk #%d with image <%s>\n", dNum, imgName)
	disk.disk6061Mu.Lock()
	var err error
	disk.imageFile, err = os.OpenFile(imgName, os.O_RDWR, 0755)
	if err != nil {
		logging.DebugPrint(disk.logID, "Failed to open image for attaching\n")
		logging.DebugPrint(logging.DebugLog, "WARN: Failed to open disk6061 image <%s> for ATTach\n", imgName)
		disk.disk6061Mu.Unlock()
		return false
	}
	disk.imageFileName = imgName
	disk.ImageAttached = true
	disk.disk6061Mu.Unlock()
	disk.bus.SetAttached(disk.devNum, imgName)
	return true
}

// Disk6061SetLogging sets the disk's internal; debug logging flag as specified
// N.B. The disk runs slower with this set.
func (disk *Disk6061T) Disk6061SetLogging(log bool) {
	disk.disk6061Mu.Lock()
	disk.debugLogging = log
	disk.disk6061Mu.Unlock()
}

func (disk *Disk6061T) disk6061StatsSender(sChan chan Disk6061StatT) {
	var stats Disk6061StatT
	for {
		disk.disk6061Mu.RLock()
		if disk.ImageAttached {
			stats.ImageAttached = true
			stats.Cylinder = disk.cylinder
			stats.Head = disk.surface
			stats.Sector = disk.sector
			stats.Reads = disk.reads
			stats.Writes = disk.writes
		} else {
			stats = Disk6061StatT{}
		}
		disk.disk6061Mu.RUnlock()
		select {
		case sChan <- stats:
		default:
		}
		time.Sleep(time.Millisecond * disk6061StatsPeriodMs)
	}
}

// Disk6061CreateBlank creates an empty disk file of the correct size for the disk6061 emulator to use
func (disk *Disk6061T) Disk6061CreateBlank(imgName string) bool {
	newFile, err := os.Create(imgName)
	if err != nil {
		return false
	}
	defer newFile.Close()
	logging.DebugPrint(disk.logID, "disk6061CreateBlank attempting to write %d bytes\n", disk6061PhysByteSize)
	w := bufio.NewWriter(newFile)
	for b := 0; b < disk6061PhysByteSize; b++ {
		w.WriteByte(0)
	}
	w.Flush()
	return true
}

// Disk6061LoadDKBT - This func mimics a system ROM routine to boot from disk.
// Rather than copying a ROM routine (!) we simply mimic its basic actions...
// Load 1st block from disk into location 0
func (disk *Disk6061T) Disk6061LoadDKBT() {
	logging.DebugPrint(disk.logID, "Disk6961LoadDKBT() called\n")
	// set posn
	disk.command = disk6061CmdRecal
	disk.disk6061DoCommand()
	disk.memAddr = 0
	disk.sectCnt = -1
	disk.command = disk6061CmdRead
	disk.disk6061DoCommand()
	logging.DebugPrint(disk.logID, "Disk6961LoadDKBT() completed\n")
}

// disk6061In implements the DIA/B/C I/O instructions for this device
func (disk *Disk6061T) disk6061In(abc byte, flag byte) (data dg.WordT) {
	disk.disk6061Mu.RLock()
	switch abc {
	case 'A':
		switch disk.instructionMode {
		case disk6061InsModeNormal:
			data = disk.rwStatus
			if disk.debugLogging {
				logging.DebugPrint(disk.logID, "DIA [Read Data Txfr Status] (Normal mode returning %s for DRV=%d\n",
					memory.WordToBinStr(disk.rwStatus), disk.drive)
			}
		case disk6061InsModeAlt1:
			data = disk.memAddr // ???
			if disk.debugLogging {
				logging.DebugPrint(disk.logID, "DIA [Read Memory Addr] (Alt Mode 1) returning %#0o for DRV=%d\n",
					data, disk.drive)
			}
		case disk6061InsModeAlt2:
			data = 0
			if disk.debugLogging {
				logging.DebugPrint(disk.logID, "DIA [Read 1st ECC Word] (Alt Mode 2) returning %#0o for DRV=%d\n",
					data, disk.drive)
			}
		}
	case 'B':
		switch disk.instructionMode {
		case disk6061InsModeNormal:
			data = disk.driveStatus & 0xfeff
			if disk.debugLogging {
				logging.DebugPrint(disk.logID, "DIB [Read Drive Status] (Normal mode) returning %s for DRV=%d\n", memory.WordToBinStr(data), disk.drive)
			}
		case disk6061InsModeAlt1:
			data = dg.WordT(0x8000) | dg.WordT(disk.ema)&0x01f
			//			if disk.mapEnabled {
			//				data = dg_dword(disk.ema&0x1f) | 0x8000
			//			} else {
			//				data = dg_dword(disk.ema & 0x1f)
			//			}
			if disk.debugLogging {
				logging.DebugPrint(disk.logID, "DIB [Read EMA] (Alt Mode 1) returning: %#0o\n", data)
			}
		case disk6061InsModeAlt2:
			data = 0
			if disk.debugLogging {
				logging.DebugPrint(disk.logID, "DIB [Read 2nd ECC Word] (Alt Mode 2) returning %#0o for DRV=%d\n",
					data, disk.drive)
			}
		}
	case 'C':
		var ssc dg.WordT
		if disk.mapEnabled {
			ssc = 1 << 15
		}
		ssc |= (dg.WordT(disk.surface) & 0x1f) << 10
		ssc |= (dg.WordT(disk.sector) & 0x1f) << 5
		ssc |= (dg.WordT(disk.sectCnt) & 0x1f)
		data = ssc
		if disk.debugLogging {
			logging.DebugPrint(disk.logID, "disk6061 DIC returning: %s\n", memory.WordToBinStr(ssc))
		}
	}
	disk.disk6061Mu.RUnlock()

	disk.disk6061HandleFlag(flag)

	return data
}

// disk6061Out implements the DOA/B/C instructions for this device
// NIO is also routed here with a dummy abc flag value of N
func (disk *Disk6061T) disk6061Out(datum dg.WordT, abc byte, flag byte) {
	disk.disk6061Mu.Lock()
	switch abc {
	case 'A':
		disk.command = extractdisk6061Command(datum)
		disk.drive = extractdisk6061DriveNo(datum)
		disk.ema = extractdisk6061EMA(datum)
		if disk.debugLogging {
			logging.DebugPrint(disk.logID, "DOA [Specify Cmd,Drv,EMA] to DRV=%d with data %s\n",
				disk.drive, memory.WordToBinStr(datum))
		}
		if memory.TestWbit(datum, 0) {
			disk.rwStatus &= ^dg.WordT(disk6061Rwdone)
			disk.rwStatus &= ^dg.WordT(disk6061Rwfault)
			disk.rwStatus &= ^dg.WordT(disk6061late)
			disk.rwStatus &= ^dg.WordT(disk6061Verify)
			disk.rwStatus &= ^dg.WordT(disk6061Surfsect)
			disk.rwStatus &= ^dg.WordT(disk6061Cylinder)
			disk.rwStatus &= ^dg.WordT(disk6061Badsector)
			disk.rwStatus &= ^dg.WordT(disk6061Ecc)
			disk.rwStatus &= ^dg.WordT(disk6061Illegalsector)
			disk.rwStatus &= ^dg.WordT(disk6061Parity)
			if disk.debugLogging {
				logging.DebugPrint(disk.logID, "... Clear R/W Done et al.\n")
			}
		}
		if memory.TestWbit(datum, 1) {
			disk.rwStatus &= ^dg.WordT(disk6061Drive0Done)
		}
		if memory.TestWbit(datum, 2) {
			disk.rwStatus &= ^dg.WordT(disk6061Drive1Done)
		}
		if memory.TestWbit(datum, 3) {
			disk.rwStatus &= ^dg.WordT(disk6061Drive2Done)
		}
		if memory.TestWbit(datum, 4) {
			disk.rwStatus &= ^dg.WordT(disk6061Drive3Done)
		}
		disk.instructionMode = disk6061InsModeNormal
		if disk.command == disk6061CmdSetAltMode1 {
			disk.instructionMode = disk6061InsModeAlt1
			if disk.debugLogging {
				logging.DebugPrint(disk.logID, "... Alt Mode 1 set\n")
			}
		}
		if disk.command == disk6061CmdSetAltMode2 {
			disk.instructionMode = disk6061InsModeAlt2
			if disk.debugLogging {
				logging.DebugPrint(disk.logID, "... Alt Mode 2 set\n")
			}
		}
		if disk.command == disk6061CmdNoOp {
			disk.instructionMode = disk6061InsModeNormal
			disk.rwStatus = 0
			disk.driveStatus = disk6061Ready
			if disk.debugLogging {
				logging.DebugPrint(disk.logID, "... NO OP command done\n")
			}
		}
		disk.lastDOAwasSeek = (disk.command == disk6061CmdSeek)
		if disk.debugLogging {
			logging.DebugPrint(disk.logID, "... CMD: %s, DRV: %d, EMA: %#o\n",
				disk.cmdDecode[disk.command], disk.drive, disk.ema)
		}
	case 'B':
		if memory.TestWbit(datum, 0) {
			disk.ema |= 0x01
		} else {
			disk.ema &= 0xfe
		}
		disk.memAddr = datum & 0x7fff
		if disk.debugLogging {
			logging.DebugPrint(disk.logID, "DOB [Specify Memory Addr] with data %s\n",
				memory.WordToBinStr(datum))
			logging.DebugPrint(disk.logID, "... MEM Addr: %#o\n", disk.memAddr)
			logging.DebugPrint(disk.logID, "... EMA: %#o\n", disk.ema)
		}
	case 'C':
		if disk.lastDOAwasSeek {
			disk.cylinder = datum & 0x03ff // mask off lower 10 bits
			if disk.debugLogging {
				logging.DebugPrint(disk.logID, "DOC [Specify Cylinder] after SEEK with data %s\n",
					memory.WordToBinStr(datum))
				logging.DebugPrint(disk.logID, "... CYL: %d\n", disk.cylinder)
			}
		} else {
			disk.mapEnabled = memory.TestWbit(datum, 0)
			disk.surface = extractsurface(datum)
			disk.sector = extractSector(datum)
			disk.sectCnt = extractSectCnt(datum)
			if disk.debugLogging {
				logging.DebugPrint(disk.logID, "DOC [Specify Surf,Sect,Cnt] (not after seek) with data %s\n",
					memory.WordToBinStr(datum))
				logging.DebugPrint(disk.logID, "... MAP: %d., SURF: %d., SECT: %d., SECCNT: %d.\n",
					memory.BoolToInt(disk.mapEnabled), disk.surface, disk.sector, disk.sectCnt)
			}
		}
	case 'N': // dummy value for NIO - we just handle the flag below
		if disk.debugLogging {
			logging.DebugPrint(disk.logID, "NIO%c received\n", flag)
		}
	}
	disk.disk6061Mu.Unlock()

	disk.disk6061HandleFlag(flag)
}

func (disk *Disk6061T) disk6061DoCommand() {

	var (
		wd                           dg.WordT
		bytesRead, bytesWritten, wIx int
		err                          error
	)

	disk.disk6061Mu.Lock()

	disk.instructionMode = disk6061InsModeNormal

	switch disk.command {

	// RECALibrate (goto pos. 0)
	case disk6061CmdRecal:
		disk.cylinder = 0
		disk.surface = 0
		disk.disk6061PositionDiskImage()
		disk.driveStatus = disk6061Ready
		disk.rwStatus = disk6061Rwdone | disk6061Drive0Done
		//disk.rwStatus = disk6061Drive0Done
		if disk.debugLogging {
			logging.DebugPrint(disk.logID, "... RECAL done, %s\n", disk.disk6061PrintableAddr())
		}

	// SEEK
	case disk6061CmdSeek:
		// action the seek
		disk.disk6061PositionDiskImage()
		disk.driveStatus = disk6061Ready
		disk.rwStatus = disk6061Rwdone | disk6061Drive0Done
		if disk.debugLogging {
			logging.DebugPrint(disk.logID, "... SEEK done, %s\n", disk.disk6061PrintableAddr())
		}

	// ===== READ from disk6061 =====
	case disk6061CmdRead:
		if disk.debugLogging {
			logging.DebugPrint(disk.logID, "... READ command invoked %s\n", disk.disk6061PrintableAddr())
			logging.DebugPrint(disk.logID, "... .... Start Address: %#o\n", disk.memAddr)
		}
		disk.rwStatus = 0

		for disk.sectCnt != 0 {
			// check CYL
			if disk.cylinder >= disk6061PhysCyls {
				disk.driveStatus = disk6061Ready
				disk.rwStatus = disk6061Rwdone | disk6061Rwfault | disk6061Cylinder
				disk.disk6061Mu.Unlock()
				return
			}
			// check SECT
			if disk.sector >= disk6061SectPerTrack {
				disk.sector = 0
				disk.surface++
				if disk.debugLogging {
					logging.DebugPrint(disk.logID, "Sector read overflow, advancing to surface %d.", disk.surface)
				}
			}
			// check SURF (head)
			if disk.surface >= disk6061SurfPerDisk {
				disk.driveStatus = disk6061Ready
				disk.rwStatus = disk6061Rwdone | disk6061Rwfault | disk6061Illegalsector
				disk.disk6061Mu.Unlock()
				return
			}
			disk.disk6061PositionDiskImage()
			bytesRead, err = disk.imageFile.Read(disk.readBuff)

			if bytesRead != disk6061BytesPerSect || err != nil {
				log.Fatalf("ERROR: unexpected return from disk6061 Image File Read: %s", err)
			}
			for wIx = 0; wIx < disk6061WordsPerSect; wIx++ {
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
			logging.DebugPrint(disk.logID, "... .... READ command finished %s\n", disk.disk6061PrintableAddr())
			logging.DebugPrint(disk.logID, "\n... .... Last Address: %#o\n", disk.memAddr)
		}
		disk.rwStatus = disk6061Rwdone | disk6061Drive0Done

	case disk6061CmdRelease:
		// I think this is a NOP on a single-processor machine

	case disk6061CmdWrite:
		if disk.debugLogging {
			logging.DebugPrint(disk.logID, "... WRITE command invoked %s\n", disk.disk6061PrintableAddr())
			logging.DebugPrint(disk.logID, "... .....  Start Address: %#o\n", disk.memAddr)
		}
		disk.rwStatus = 0

		for disk.sectCnt != 0 {
			// check CYL
			if disk.cylinder >= disk6061PhysCyls {
				disk.driveStatus = disk6061Ready
				disk.rwStatus = disk6061Rwdone | disk6061Rwfault | disk6061Cylinder
				disk.disk6061Mu.Unlock()
				return
			}
			// check SECT
			if disk.sector >= disk6061SectPerTrack {
				disk.sector = 0
				disk.surface++
				if disk.debugLogging {
					logging.DebugPrint(disk.logID, "Sector write overflow, advancing to surface %d.", disk.surface)
				}
			}
			// check SURF (head)/SECT
			if disk.surface >= disk6061SurfPerDisk {
				disk.driveStatus = disk6061Ready
				disk.rwStatus = disk6061Rwdone | disk6061Rwfault | disk6061Illegalsector
				disk.disk6061Mu.Unlock()
				return
			}
			disk.disk6061PositionDiskImage()
			for wIx = 0; wIx < disk6061WordsPerSect; wIx++ {
				wd = memory.ReadWordBmcChan16bit(&disk.memAddr)
				disk.writeBuff[(wIx*2)+1] = byte(wd >> 8)
				disk.writeBuff[wIx*2] = byte(wd)
			}
			bytesWritten, err = disk.imageFile.Write(disk.writeBuff)
			if bytesWritten != disk6061BytesPerSect || err != nil {
				log.Fatalf("ERROR: unexpected return from disk6061 Image File Write: %s", err)
			}
			disk.sector++
			disk.sectCnt++
			disk.writes++

			if disk.debugLogging {
				logging.DebugPrint(disk.logID, "Buffer: %X\n", disk.writeBuff)
			}
		}
		if disk.debugLogging {
			logging.DebugPrint(disk.logID, "... ..... WRITE command finished %s\n", disk.disk6061PrintableAddr())
			logging.DebugPrint(disk.logID, "... ..... Last Address: %#o\n", disk.memAddr)
		}
		disk.driveStatus = disk6061Ready
		disk.rwStatus = disk6061Rwdone //| disk6061Drive0Done

	default:
		log.Fatalf("disk6061 Disk R/W Command %d not yet implemented\n", disk.command)
	}
	disk.disk6061Mu.Unlock()
}

func (disk *Disk6061T) disk6061HandleFlag(f byte) {
	switch f {
	case 'S':
		disk.bus.SetBusy(disk.devNum, true)
		disk.bus.SetDone(disk.devNum, false)
		// TODO stop any I/O
		disk.disk6061Mu.Lock()
		disk.rwStatus = 0
		// TODO start I/O timeout
		if disk.debugLogging {
			logging.DebugPrint(disk.logID, "... S flag set\n")
		}
		disk.disk6061Mu.Unlock()
		disk.disk6061DoCommand()
		disk.bus.SetBusy(disk.devNum, false)
		disk.bus.SetDone(disk.devNum, true)
		// send IRQ if not masked out
		//if !BusIsDevMasked(disk.devNum) {
		// InterruptingDev[disk.devNum] = true
		// IRQ = true
		disk.bus.SendInterrupt(disk.devNum)
		//}

	case 'C':
		disk.bus.SetBusy(disk.devNum, false)
		disk.bus.SetDone(disk.devNum, false)
		disk.disk6061Mu.Lock()
		disk.rwStatus = 0
		disk.disk6061Mu.Unlock()

	case 'P':
		disk.bus.SetBusy(disk.devNum, false)
		disk.disk6061Mu.Lock()
		if disk.debugLogging {
			logging.DebugPrint(disk.logID, "... P flag set\n")
		}
		disk.rwStatus = 0
		disk.disk6061Mu.Unlock()
		disk.disk6061DoCommand()
		disk.bus.SendInterrupt(disk.devNum)

	default:
		// no/empty flag - nothing to do
	}
}

// set the MV/Em disk image file postion according to current C/H/S
func (disk *Disk6061T) disk6061PositionDiskImage() {
	var offset, r int64
	var err error
	offset = (((int64(disk.cylinder*disk6061SurfPerDisk) + int64(disk.surface)) * int64(disk6061SectPerTrack)) + int64(disk.sector)) * disk6061BytesPerSect
	r, err = disk.imageFile.Seek(offset, 0)
	if r != offset || err != nil {
		log.Fatal("disk6061 could not postition disk image via seek()")
	}
}

func (disk *Disk6061T) disk6061PrintableAddr() string {
	// MUST BE LOCKED BY CALLER
	pa := fmt.Sprintf("DRV: %d, CYL: %d, SURF: %d, SECT: %d, SECCNT: %d",
		disk.drive, disk.cylinder,
		disk.surface, disk.sector, disk.sectCnt)
	return pa
}

// reset the disk6061 controller
func (disk *Disk6061T) disk6061Reset() {
	disk.disk6061Mu.Lock()
	disk.instructionMode = disk6061InsModeNormal
	disk.rwStatus = 0
	disk.command = disk6061CmdRead
	disk.cylinder = 0
	disk.surface = 0
	disk.sector = 0
	disk.sectCnt = 0
	disk.driveStatus = disk6061Ready
	if disk.debugLogging {
		logging.DebugPrint(disk.logID, "disk6061 Reset\n")
	}
	disk.disk6061Mu.Unlock()
}

func extractdisk6061Command(word dg.WordT) int8 {
	return int8((word & 0x0780) >> 7)
}

func extractdisk6061DriveNo(word dg.WordT) uint8 {
	return uint8((word & 0x60) >> 5)
}

func extractdisk6061EMA(word dg.WordT) uint8 {
	return uint8(word & 0x1f)
}

func extractSector(word dg.WordT) uint8 {
	return uint8((word & 0x03e0) >> 5)
}

func extractSectCnt(word dg.WordT) int8 {
	tmpWd := word & 0x01f
	if tmpWd != 0 { // sign-extend
		tmpWd |= 0xe0
	}
	return int8(tmpWd)
}

func extractsurface(word dg.WordT) uint8 {
	return uint8((word & 0x7c00) >> 10)
}
