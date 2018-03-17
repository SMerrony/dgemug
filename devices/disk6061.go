// disk6061.go

// Copyright (C) 2017,2018  Steve Merrony

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

	"github.com/SMerrony/dgemug"
	"github.com/SMerrony/dgemug/logging"
	"github.com/SMerrony/dgemug/memory"
	"github.com/SMerrony/dgemug/util"
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

type disk6061T struct {
	// MV/Em internals...
	ImageAttached       bool
	disk6061Mu          sync.RWMutex
	devNum              int
	logID               int
	imageFileName       string
	imageFile           *os.File
	reads, writes       uint64
	readBuff, writeBuff []byte
	// DG data...
	cmdDrvAddr      byte // 6-bit?
	command         int8 // 4-bit
	rwCommand       int8
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

var (
	disk6061                     disk6061T
	wd                           dg.WordT
	ssc                          dg.WordT
	bytesRead, bytesWritten, wIx int
	err                          error
	cmdDecode                    [disk6061CmdFormat + 1]string
	debugLogging                 bool
)

// disk6061Init must be called to initialise the emulated disk6061 controller
func Disk6061Init(dev int, statsChann chan Disk6061StatT, logId int, logging bool) {
	disk6061.disk6061Mu.Lock()
	defer disk6061.disk6061Mu.Unlock()
	disk6061.devNum = dev
	disk6061.logID = logId
	debugLogging = logging

	cmdDecode = [...]string{"READ", "RECAL", "SEEK", "STOP", "OFFSET FWD", "OFFSET REV",
		"WRITE DISABLE", "RELEASE", "TRESPASS", "SET ALT MODE 1", "SET ALT MODE 2",
		"NO OP", "VERIFY", "READ BUFFERS", "WRITE", "FORMAT"}

	go disk6061StatsSender(statsChann)

	BusSetResetFunc(disk6061.devNum, disk6061Reset)
	BusSetDataInFunc(disk6061.devNum, disk6061In)
	BusSetDataOutFunc(disk6061.devNum, disk6061Out)
	disk6061.ImageAttached = false
	disk6061.instructionMode = disk6061InsModeNormal
	disk6061.driveStatus = disk6061Ready
	disk6061.mapEnabled = false
	disk6061.readBuff = make([]byte, disk6061BytesPerSect)
	disk6061.writeBuff = make([]byte, disk6061BytesPerSect)
}

// attempt to attach an extant MV/Em disk image to the running emulator
func Disk6061Attach(dNum int, imgName string) bool {
	// TODO Disk Number not currently used
	logging.DebugPrint(disk6061.logID, "disk6061Attach called for disk #%d with image <%s>\n", dNum, imgName)
	disk6061.disk6061Mu.Lock()
	disk6061.imageFile, err = os.OpenFile(imgName, os.O_RDWR, 0755)
	if err != nil {
		logging.DebugPrint(disk6061.logID, "Failed to open image for attaching\n")
		logging.DebugPrint(logging.DebugLog, "WARN: Failed to open disk6061 image <%s> for ATTach\n", imgName)
		disk6061.disk6061Mu.Unlock()
		return false
	}
	disk6061.imageFileName = imgName
	disk6061.ImageAttached = true
	disk6061.disk6061Mu.Unlock()
	BusSetAttached(disk6061.devNum, imgName)
	return true
}

func disk6061StatsSender(sChan chan Disk6061StatT) {
	var stats Disk6061StatT
	for {
		disk6061.disk6061Mu.RLock()
		if disk6061.ImageAttached {
			stats.ImageAttached = true
			stats.Cylinder = disk6061.cylinder
			stats.Head = disk6061.surface
			stats.Sector = disk6061.sector
			stats.Reads = disk6061.reads
			stats.Writes = disk6061.writes
		} else {
			stats = Disk6061StatT{}
		}
		disk6061.disk6061Mu.RUnlock()
		select {
		case sChan <- stats:
		default:
		}
		time.Sleep(time.Millisecond * disk6061StatsPeriodMs)
	}
}

// Create an empty disk file of the correct size for the disk6061 emulator to use
func Disk6061CreateBlank(imgName string) bool {
	newFile, err := os.Create(imgName)
	if err != nil {
		return false
	}
	defer newFile.Close()
	logging.DebugPrint(disk6061.logID, "disk6061CreateBlank attempting to write %d bytes\n", disk6061PhysByteSize)
	w := bufio.NewWriter(newFile)
	for b := 0; b < disk6061PhysByteSize; b++ {
		w.WriteByte(0)
	}
	w.Flush()
	return true
}

// disk6061In implements the DIA/B/C I/O instructions for this device
func disk6061In(abc byte, flag byte) (data dg.WordT) {
	disk6061.disk6061Mu.RLock()
	switch abc {
	case 'A':
		switch disk6061.instructionMode {
		case disk6061InsModeNormal:
			data = disk6061.rwStatus
			if debugLogging {
				logging.DebugPrint(disk6061.logID, "DIA [Read Data Txfr Status] (Normal mode returning %s for DRV=%d\n",
					util.WordToBinStr(disk6061.rwStatus), disk6061.drive)
			}
		case disk6061InsModeAlt1:
			data = disk6061.memAddr // ???
			if debugLogging {
				logging.DebugPrint(disk6061.logID, "DIA [Read Memory Addr] (Alt Mode 1) returning %#0o for DRV=%d\n",
					data, disk6061.drive)
			}
		case disk6061InsModeAlt2:
			data = 0
			if debugLogging {
				logging.DebugPrint(disk6061.logID, "DIA [Read 1st ECC Word] (Alt Mode 2) returning %#0o for DRV=%d\n",
					data, disk6061.drive)
			}
		}
	case 'B':
		switch disk6061.instructionMode {
		case disk6061InsModeNormal:
			data = disk6061.driveStatus & 0xfeff
		case disk6061InsModeAlt1:
			data = dg.WordT(0x8000) | dg.WordT(disk6061.ema)&0x01f
			//			if disk6061.mapEnabled {
			//				data = dg_dword(disk6061.ema&0x1f) | 0x8000
			//			} else {
			//				data = dg_dword(disk6061.ema & 0x1f)
			//			}
			if debugLogging {
				logging.DebugPrint(disk6061.logID, "DIB [Read EMA] (Alt Mode 1) returning: %d\n",
					data)
			}
		case disk6061InsModeAlt2:
			data = 0
			if debugLogging {
				logging.DebugPrint(disk6061.logID, "DIB [Read 2nd ECC Word] (Alt Mode 2) returning %#0o for DRV=%d\n",
					data, disk6061.drive)
			}
		}
	case 'C':
		ssc = 0
		if disk6061.mapEnabled {
			ssc = 1 << 15
		}
		ssc |= (dg.WordT(disk6061.surface) & 0x1f) << 10
		ssc |= (dg.WordT(disk6061.sector) & 0x1f) << 5
		ssc |= (dg.WordT(disk6061.sectCnt) & 0x1f)
		data = ssc
		if debugLogging {
			logging.DebugPrint(disk6061.logID, "disk6061 DIC returning: %s\n", util.WordToBinStr(ssc))
		}
	}
	disk6061.disk6061Mu.RUnlock()

	disk6061HandleFlag(flag)

	return data
}

// disk6061Out implements the DOA/B/C instructions for this device
// NIO is also routed here with a dummy abc flag value of N
func disk6061Out(datum dg.WordT, abc byte, flag byte) {
	disk6061.disk6061Mu.Lock()
	switch abc {
	case 'A':
		disk6061.command = extractdisk6061Command(datum)
		disk6061.drive = extractdisk6061DriveNo(datum)
		disk6061.ema = extractdisk6061EMA(datum)
		if util.TestWbit(datum, 0) {
			disk6061.rwStatus &= ^dg.WordT(disk6061Rwdone)
		}
		if util.TestWbit(datum, 1) {
			disk6061.rwStatus &= ^dg.WordT(disk6061Drive0Done)
		}
		if util.TestWbit(datum, 2) {
			disk6061.rwStatus &= ^dg.WordT(disk6061Drive1Done)
		}
		if util.TestWbit(datum, 3) {
			disk6061.rwStatus &= ^dg.WordT(disk6061Drive2Done)
		}
		if util.TestWbit(datum, 4) {
			disk6061.rwStatus &= ^dg.WordT(disk6061Drive3Done)
		}
		disk6061.instructionMode = disk6061InsModeNormal
		if disk6061.command == disk6061CmdSetAltMode1 {
			disk6061.instructionMode = disk6061InsModeAlt1
		}
		if disk6061.command == disk6061CmdSetAltMode2 {
			disk6061.instructionMode = disk6061InsModeAlt2
		}
		if disk6061.command == disk6061CmdNoOp {
			disk6061.instructionMode = disk6061InsModeNormal
			disk6061.rwStatus = 0
			disk6061.driveStatus = disk6061Ready
			if debugLogging {
				logging.DebugPrint(disk6061.logID, "... NO OP command done\n")
			}
		}
		disk6061.lastDOAwasSeek = (disk6061.command == disk6061CmdSeek)
		if debugLogging {
			logging.DebugPrint(disk6061.logID, "DOA [Specify Cmd,Drv,EMA] to DRV=%d with data %s\n",
				disk6061.drive, util.WordToBinStr(datum))
			logging.DebugPrint(disk6061.logID, "... CMD: %s, DRV: %d, EMA: %d\n",
				cmdDecode[disk6061.command], disk6061.drive, disk6061.ema)
		}
	case 'B':
		if util.TestWbit(datum, 0) {
			disk6061.ema |= 0x01
		} else {
			disk6061.ema &= 0xfe
		}
		disk6061.memAddr = datum & 0x7fff
		if debugLogging {
			logging.DebugPrint(disk6061.logID, "DOB [Specify Memory Addr] with data %s\n",
				util.WordToBinStr(datum))
			logging.DebugPrint(disk6061.logID, "... MEM Addr: %d\n", disk6061.memAddr)
			logging.DebugPrint(disk6061.logID, "... EMA: %d\n", disk6061.ema)
		}
	case 'C':
		if disk6061.lastDOAwasSeek {
			disk6061.cylinder = datum & 0x03ff // mask off lower 10 bits
			if debugLogging {
				logging.DebugPrint(disk6061.logID, "DOC [Specify Cylinder] after SEEK with data %s\n",
					util.WordToBinStr(datum))
				logging.DebugPrint(disk6061.logID, "... CYL: %d\n", disk6061.cylinder)
			}
		} else {
			disk6061.mapEnabled = util.TestWbit(datum, 0)
			disk6061.surface = extractsurface(datum)
			disk6061.sector = extractSector(datum)
			disk6061.sectCnt = extractSectCnt(datum)
			if debugLogging {
				logging.DebugPrint(disk6061.logID, "DOC [Specify Surf,Sect,Cnt] (not after seek) with data %s\n",
					util.WordToBinStr(datum))
				logging.DebugPrint(disk6061.logID, "... MAP: %d, SURF: %d, SECT: %d, SECCNT: %d\n",
					util.BoolToInt(disk6061.mapEnabled), disk6061.surface, disk6061.sector, disk6061.sectCnt)
			}
		}
	case 'N': // dummy value for NIO - we just handle the flag below
		if debugLogging {
			logging.DebugPrint(disk6061.logID, "NIO%c received\n", flag)
		}
	}
	disk6061.disk6061Mu.Unlock()

	disk6061HandleFlag(flag)
}

func disk6061DoCommand() {

	disk6061.disk6061Mu.Lock()

	disk6061.instructionMode = disk6061InsModeNormal

	switch disk6061.command {

	// RECALibrate (goto pos. 0)
	case disk6061CmdRecal:
		disk6061.cylinder = 0
		disk6061.surface = 0
		disk6061PositionDiskImage()
		disk6061.driveStatus = disk6061Ready
		disk6061.rwStatus = disk6061Rwdone | disk6061Drive0Done
		if debugLogging {
			logging.DebugPrint(disk6061.logID, "... RECAL done, %s\n", disk6061PrintableAddr())
		}

	// SEEK
	case disk6061CmdSeek:
		// action the seek
		disk6061PositionDiskImage()
		disk6061.driveStatus = disk6061Ready
		disk6061.rwStatus = disk6061Rwdone | disk6061Drive0Done
		if debugLogging {
			logging.DebugPrint(disk6061.logID, "... SEEK done, %s\n", disk6061PrintableAddr())
		}

	// ===== READ from disk6061 =====
	case disk6061CmdRead:
		if debugLogging {
			logging.DebugPrint(disk6061.logID, "... READ command invoked %s\n", disk6061PrintableAddr())
			logging.DebugPrint(disk6061.logID, "... .... Start Address: %d\n", disk6061.memAddr)
		}
		disk6061.rwStatus = 0

		for disk6061.sectCnt != 0 {
			// check CYL
			if disk6061.cylinder >= disk6061PhysCyls {
				disk6061.driveStatus = disk6061Ready
				disk6061.rwStatus = disk6061Rwdone | disk6061Rwfault | disk6061Cylinder
				disk6061.disk6061Mu.Unlock()
				return
			}
			// check SECT
			if disk6061.sector >= disk6061SectPerTrack {
				disk6061.sector = 0
				disk6061.surface++
				if debugLogging {
					logging.DebugPrint(disk6061.logID, "Sector read overflow, advancing to surface %d",
						disk6061.surface)
				}
				// disk6061.driveStatus = disk6061Ready
				// disk6061.rwStatus = disk6061Rwdone | disk6061Rwfault | disk6061_ILLEGALSECTOR
				// disk6061.disk6061Mu.Unlock()
				// return
			}
			// check SURF (head)
			if disk6061.surface >= disk6061SurfPerDisk {
				disk6061.driveStatus = disk6061Ready
				disk6061.rwStatus = disk6061Rwdone | disk6061Rwfault | disk6061Illegalsector
				disk6061.disk6061Mu.Unlock()
				return
			}
			disk6061PositionDiskImage()
			bytesRead, err = disk6061.imageFile.Read(disk6061.readBuff)

			if bytesRead != disk6061BytesPerSect || err != nil {
				log.Fatalf("ERROR: unexpected return from disk6061 Image File Read: %s", err)
			}
			for wIx = 0; wIx < disk6061WordsPerSect; wIx++ {
				wd = (dg.WordT(disk6061.readBuff[wIx*2]) << 8) | dg.WordT(disk6061.readBuff[(wIx*2)+1])
				memory.WriteWordBmcChan16bit(&disk6061.memAddr, wd)
			}
			disk6061.sector++
			disk6061.sectCnt++
			disk6061.reads++

			if debugLogging {
				logging.DebugPrint(disk6061.logID, "Buffer: %X\n", disk6061.readBuff)
			}

		}
		if debugLogging {
			logging.DebugPrint(disk6061.logID, "... .... READ command finished %s\n", disk6061PrintableAddr())
			logging.DebugPrint(disk6061.logID, "\n... .... Last Address: %d\n", disk6061.memAddr)
		}
		disk6061.rwStatus = disk6061Rwdone //| disk6061Drive0Done

	case disk6061CmdRelease:
		// I think this is a NOP on a single-processor machine

	case disk6061CmdWrite:
		if debugLogging {
			logging.DebugPrint(disk6061.logID, "... WRITE command invoked %s\n", disk6061PrintableAddr())
			logging.DebugPrint(disk6061.logID, "... .....  Start Address: %d\n", disk6061.memAddr)
		}
		disk6061.rwStatus = 0

		for disk6061.sectCnt != 0 {
			// check CYL
			if disk6061.cylinder >= disk6061PhysCyls {
				disk6061.driveStatus = disk6061Ready
				disk6061.rwStatus = disk6061Rwdone | disk6061Rwfault | disk6061Cylinder
				disk6061.disk6061Mu.Unlock()
				return
			}
			// check SECT
			if disk6061.sector >= disk6061SectPerTrack {
				disk6061.sector = 0
				disk6061.surface++
				if debugLogging {
					logging.DebugPrint(disk6061.logID, "Sector write overflow, advancing to surface %d",
						disk6061.surface)
				}
				// disk6061.driveStatus = disk6061Ready
				// disk6061.rwStatus = disk6061Rwdone | disk6061Rwfault | disk6061_ILLEGALSECTOR
				// disk6061.disk6061Mu.Unlock()
				// return
			}
			// check SURF (head)/SECT
			if disk6061.surface >= disk6061SurfPerDisk {
				disk6061.driveStatus = disk6061Ready
				disk6061.rwStatus = disk6061Rwdone | disk6061Rwfault | disk6061Illegalsector
				disk6061.disk6061Mu.Unlock()
				return
			}
			disk6061PositionDiskImage()
			for wIx = 0; wIx < disk6061WordsPerSect; wIx++ {
				wd = memory.ReadWordBmcChan16bit(&disk6061.memAddr)
				disk6061.writeBuff[wIx*2] = byte((wd & 0xff00) >> 8)
				disk6061.writeBuff[(wIx*2)+1] = byte(wd & 0x00ff)
			}
			bytesWritten, err = disk6061.imageFile.Write(disk6061.writeBuff)
			if bytesWritten != disk6061BytesPerSect || err != nil {
				log.Fatalf("ERROR: unexpected return from disk6061 Image File Write: %s", err)
			}
			disk6061.sector++
			disk6061.sectCnt++
			disk6061.writes++

			if debugLogging {
				logging.DebugPrint(disk6061.logID, "Buffer: %X\n", disk6061.writeBuff)
			}
		}
		if debugLogging {
			logging.DebugPrint(disk6061.logID, "... ..... WRITE command finished %s\n", disk6061PrintableAddr())
			logging.DebugPrint(disk6061.logID, "... ..... Last Address: %d\n", disk6061.memAddr)
		}
		disk6061.driveStatus = disk6061Ready
		disk6061.rwStatus = disk6061Rwdone //| disk6061Drive0Done

	default:
		log.Fatalf("disk6061 Disk R/W Command %d not yet implemented\n", disk6061.command)
	}
	disk6061.disk6061Mu.Unlock()
}

func disk6061HandleFlag(f byte) {
	switch f {
	case 'S':
		BusSetBusy(disk6061.devNum, true)
		BusSetDone(disk6061.devNum, false)
		// TODO stop any I/O
		disk6061.disk6061Mu.Lock()
		disk6061.rwStatus = 0
		// TODO start I/O timeout
		disk6061.rwCommand = disk6061.command
		if debugLogging {
			logging.DebugPrint(disk6061.logID, "... S flag set\n")
		}
		disk6061.disk6061Mu.Unlock()
		disk6061DoCommand()
		BusSetBusy(disk6061.devNum, false)
		BusSetDone(disk6061.devNum, true)
		// send IRQ if not masked out
		if !BusIsDevMasked(disk6061.devNum) {
			// InterruptingDev[disk6061.devNum] = true
			// IRQ = true
			BusSendInterrupt(disk6061.devNum)
		}

	case 'C':
		BusSetBusy(disk6061.devNum, false)
		BusSetDone(disk6061.devNum, false)
		disk6061.disk6061Mu.Lock()
		disk6061.rwStatus = 0
		disk6061.rwCommand = 0
		disk6061.disk6061Mu.Unlock()

	case 'P':
		BusSetBusy(disk6061.devNum, false)
		disk6061.disk6061Mu.Lock()
		if debugLogging {
			logging.DebugPrint(disk6061.logID, "... P flag set\n")
		}
		disk6061.rwStatus = 0
		disk6061.disk6061Mu.Unlock()
		disk6061DoCommand()
		//disk6061.rwStatus = disk6061Drive0Done
		BusSetBusy(disk6061.devNum, false)
		BusSetDone(disk6061.devNum, true)
		// send IRQ if not masked out
		if !BusIsDevMasked(disk6061.devNum) {
			// InterruptingDev[disk6061.devNum] = true
			// IRQ = true
			BusSendInterrupt(disk6061.devNum)
		}

	default:
		// no/empty flag - nothing to do
	}
}

// set the MV/Em disk image file postion according to current C/H/S
func disk6061PositionDiskImage() {
	var offset, r int64
	//lba = ((int64(disk6061.cylinder*disk6061SurfPerDisk) + int64(disk6061.surface)) * int64(disk6061SectPerTrack)) + int64(disk6061.sector)
	offset = (((int64(disk6061.cylinder*disk6061SurfPerDisk) + int64(disk6061.surface)) * int64(disk6061SectPerTrack)) + int64(disk6061.sector)) * disk6061BytesPerSect
	r, err = disk6061.imageFile.Seek(offset, 0)
	if r != offset || err != nil {
		log.Fatal("disk6061 could not postition disk image via seek()")
	}
}

func disk6061PrintableAddr() string {
	// MUST BE LOCKED BY CALLER
	pa := fmt.Sprintf("DRV: %d, CYL: %d, SURF: %d, SECT: %d, SECCNT: %d",
		disk6061.drive, disk6061.cylinder,
		disk6061.surface, disk6061.sector, disk6061.sectCnt)
	return pa
}

// reset the disk6061 controller
func disk6061Reset() {
	disk6061.disk6061Mu.Lock()
	disk6061.instructionMode = disk6061InsModeNormal
	disk6061.rwStatus = 0
	disk6061.driveStatus = disk6061Ready
	if debugLogging {
		logging.DebugPrint(disk6061.logID, "disk6061 Reset\n")
	}
	disk6061.disk6061Mu.Unlock()
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
