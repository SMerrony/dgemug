// magtape6026.go

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
	"log"
	"os"
	"sync"
	"time"

	"github.com/SMerrony/dgemug/dg"
	"github.com/SMerrony/dgemug/logging"
	"github.com/SMerrony/dgemug/memory"
	"github.com/SMerrony/simhtape/pkg/simhtape"
)

const (
	mtMaxRecordSizeW = 16384
	mtMaxRecordSizeB = mtMaxRecordSizeW * 2
	mtEOF            = 0
	mtCmdCount       = 11
	mtCmdMask        = 0x00b8

	mtCmdReadBits        = 0x0000
	mtCmdRewindBits      = 0x0008
	mtCmdCtrlModeBits    = 0x0010
	mtCmdSpaceFwdBits    = 0x0018
	mtCmdSpaceRevBits    = 0x0020
	mtCmdWiteBits        = 0x0028
	mtCmdWriteEOFBits    = 0x0030
	mtCmdEraseBits       = 0x0038
	mtCmdReadNonStopBits = 0x0080
	mtCmdUnloadBits      = 0x0088
	mtCmdDriveModeBits   = 0x0090

	mtCmdRead        = 0
	mtCmdRewind      = 1
	mtCmdCtrlMode    = 2
	mtCmdSpaceFwd    = 3
	mtCmdSpaceRev    = 4
	mtCmdWrite       = 5
	mtCmdWriteEOF    = 6
	mtCmdErase       = 7
	mtCmdReadNonStop = 8
	mtCmdUnload      = 9
	mtCmdDriveMode   = 10

	mtSr1Error         = 1 << 15
	mtSr1DataLate      = 1 << 14
	mtSr1Rewinding     = 1 << 13
	mtSr1Illegal       = 1 << 12
	mtSr1HiDensity     = 1 << 11
	mtSr1DataError     = 1 << 10
	mtSr1EOT           = 1 << 9
	mtSr1EOF           = 1 << 8
	mtSr1BOT           = 1 << 7
	mtSr19Track        = 1 << 6
	mtSr1BadTape       = 1 << 5
	mtSr1Reserved      = 1 << 4
	mtSr1StatusChanged = 1 << 3
	mtSr1WriteLock     = 1 << 2
	mtSr1OddChar       = 1 << 1
	mtSr1UnitReady     = 1

	mtSr1Readable = "ELRIHDEFB9TrSWOR"

	mtSr2Error  = 1 << 15
	mtSr2PEMode = 1
)

const mtStatsPeriodMs = 500 // Will update status this often

const maxTapes = 8

type mtT struct {
	mtDataMu               sync.RWMutex
	devNum                 int
	imageAttached          [maxTapes]bool
	fileName               [maxTapes]string
	simhFile               [maxTapes]*os.File
	statusReg1, statusReg2 dg.WordT
	memAddrReg             dg.PhysAddrT
	negWordCntReg          int16
	currentCmd             int
	currentUnit            int
	// debug                  bool
}

// MtStatT holds data for the status collector
type MtStatT struct {
	ImageAttached          [maxTapes]bool
	FileName               [maxTapes]string
	MemAddrReg             dg.PhysAddrT
	CurrentCmd             int
	StatusReg1, StatusReg2 dg.WordT
}

var (
	mt         mtT
	commandSet [mtCmdCount]dg.WordT
	logID      int
)

// MtInit sets the initial state of the (unmounted) tape drive(s)
func MtInit(dev int, statsChan chan MtStatT, logId int) bool {
	mt.devNum = dev
	logID = logId
	commandSet[mtCmdRead] = mtCmdReadBits
	commandSet[mtCmdRewind] = mtCmdRewindBits
	commandSet[mtCmdCtrlMode] = mtCmdCtrlModeBits
	commandSet[mtCmdSpaceFwd] = mtCmdSpaceFwdBits
	commandSet[mtCmdSpaceRev] = mtCmdSpaceRevBits
	commandSet[mtCmdWrite] = mtCmdWiteBits
	commandSet[mtCmdWriteEOF] = mtCmdWriteEOFBits
	commandSet[mtCmdErase] = mtCmdEraseBits
	commandSet[mtCmdReadNonStop] = mtCmdReadNonStopBits
	commandSet[mtCmdUnload] = mtCmdUnloadBits
	commandSet[mtCmdDriveMode] = mtCmdDriveModeBits

	BusSetResetFunc(mt.devNum, MtReset)
	BusSetDataInFunc(mt.devNum, mtDataIn)
	BusSetDataOutFunc(mt.devNum, mtDataOut)

	go mtStatSender(statsChan)

	mt.mtDataMu.Lock()
	mt.statusReg1 = mtSr1HiDensity | mtSr19Track | mtSr1UnitReady
	mt.statusReg2 = mtSr2PEMode
	mt.mtDataMu.Unlock()

	logging.DebugPrint(logID, "mt Initialised via call to mtInit()\n")
	return true
}

// mtStatSender provides a near real-time view of mt status and should be run as a Goroutine
// TODO = only handles Unit 0 at the moment
func mtStatSender(sChan chan MtStatT) {
	var stats MtStatT
	logging.DebugPrint(logging.DebugLog, "dskpStatSender() started\n")
	for {
		mt.mtDataMu.RLock()
		if mt.imageAttached[0] {
			stats.ImageAttached[0] = true
			stats.FileName[0] = mt.fileName[0]
			stats.MemAddrReg = mt.memAddrReg
			stats.CurrentCmd = mt.currentCmd // Could decode this
			stats.StatusReg1 = mt.statusReg1
			stats.StatusReg2 = mt.statusReg2
		} else {
			stats = MtStatT{}
		}
		mt.mtDataMu.RUnlock()
		// Non-blocking send of stats
		select {
		case sChan <- stats:
			// stats sent
		default:
			// do not block
		}
		time.Sleep(mtStatsPeriodMs * time.Millisecond)
	}
}

// MtReset Resets the mt to known state (any tapes are rewound)
func MtReset() {
	mt.mtDataMu.Lock()
	for t := 0; t < maxTapes; t++ {
		if mt.imageAttached[t] {
			simhtape.Rewind(mt.simhFile[t])
		}
	}
	// BOT is an error state...
	mt.statusReg1 = mtSr1Error | mtSr1HiDensity | mtSr19Track | mtSr1BOT | mtSr1StatusChanged | mtSr1UnitReady
	mt.statusReg2 = mtSr2PEMode
	mt.memAddrReg = 0
	mt.negWordCntReg = 0
	mt.currentCmd = 0
	mt.currentUnit = 0
	mt.mtDataMu.Unlock()
	logging.DebugPrint(logID, "mt Reset via call to mtReset()\n")
}

// MTAttach attaches a SimH tape image file to the emulated tape drive
func MtAttach(tNum int, imgName string) bool {
	logging.DebugPrint(logID, "mtAttach called on unit #%d with image file: %s\n", tNum, imgName)
	f, err := os.Open(imgName)
	if err != nil {
		logging.DebugPrint(logID, "ERROR: Could not open simH Tape Image file: %s, due to: %s\n", imgName, err.Error())
		return false
	}
	mt.mtDataMu.Lock()
	mt.fileName[tNum] = imgName
	mt.simhFile[tNum] = f
	mt.imageAttached[tNum] = true
	mt.statusReg1 = mtSr1Error | mtSr1HiDensity | mtSr19Track | mtSr1BOT | mtSr1StatusChanged | mtSr1UnitReady
	mt.statusReg2 = mtSr2PEMode
	mt.mtDataMu.Unlock()
	BusSetAttached(mt.devNum, imgName)
	return true

}

// MTDetach disassociates a tape file image from the drive
func MtDetach(tNum int) bool {
	logging.DebugPrint(logID, "mtDetach called on unit #%d\n", tNum)
	mt.mtDataMu.Lock()
	mt.fileName[tNum] = ""
	mt.simhFile[tNum] = nil
	mt.imageAttached[tNum] = false
	mt.statusReg1 = mtSr1Error | mtSr1HiDensity | mtSr19Track | mtSr1BOT | mtSr1StatusChanged | mtSr1UnitReady
	mt.statusReg2 = mtSr2PEMode
	mt.mtDataMu.Unlock()
	BusSetDetached(mt.devNum)
	return true
}

// MtScanImage scans an attached SimH tape image to ensure it makes sense
// (This is just a pass-through to the equivalent function in simhtape)
func MtScanImage(tNum int) string {
	mt.mtDataMu.RLock()
	imageName := mt.fileName[tNum]
	att := mt.imageAttached[tNum]
	mt.mtDataMu.RUnlock()
	if !att {
		return "WARNING: No image attached"
	}
	return simhtape.ScanImage(imageName, false)
}

// MtLoadTBoot - This function fakes the ROM/SCP boot-from-tape routine.
// Rather than copying a ROM and executing that, we simply mimic its basic actions..
// Load file 0 from tape (1 x 2k block)
// Put the loaded code at physical location 0
// ...
func MtLoadTBoot() {
	const (
	// tbootSizeB = 2048
	// tbootSizeW = 1024
	)
	var (
		byte0, byte1 byte
		word         dg.WordT
		wdix, memix  dg.PhysAddrT
	)
	logging.DebugPrint(logID, "mtLoadTBoot() called\n")
	tNum := 0
	mt.mtDataMu.Lock()
	defer mt.mtDataMu.Unlock()
	simhtape.Rewind(mt.simhFile[tNum])
	logging.DebugPrint(logID, "... tape rewound\n")

readLoop:
	for {
		hdr, ok := simhtape.ReadMetaData(mt.simhFile[tNum])
		// if !ok || hdr != tbootSizeB {
		if !ok {
			logging.DebugPrint(logging.DebugLog, "WARN: mtLoadTBoot called when no bootable tape image attached\n")
			return
		}
		logging.DebugPrint(logID, "... header read, size is %d\n", hdr)
		switch hdr {
		case simhtape.SimhMtrTmk: // Tape Mark (separates files)
			break readLoop
		default:
			tbootSizeW := hdr / 2
			tapeData, ok := simhtape.ReadRecordData(mt.simhFile[tNum], int(hdr))
			if ok {
				logging.DebugPrint(logID, "... data read\n")
			} else {
				logging.DebugPrint(logID, "... error reading data\n")
				logging.DebugPrint(logging.DebugLog, "WARNING: Could not read data in mtLoadTBoot()\n")
				return
			}

			logging.DebugPrint(logID, "... loading data into memory starting at address %d\n", memix)
			for wdix = 0; wdix < dg.PhysAddrT(tbootSizeW); wdix++ {
				byte1 = tapeData[wdix*2]
				byte0 = tapeData[wdix*2+1]
				word = dg.WordT(byte1)<<8 | dg.WordT(byte0)
				memory.WriteWord(memix+wdix, word)
			}
			memix += dg.PhysAddrT(tbootSizeW)
			logging.DebugPrint(logID, "... finished loading data at address %d\n", memix+wdix)
			trailer, ok := simhtape.ReadMetaData(mt.simhFile[tNum])
			if hdr != trailer || !ok {
				logging.DebugPrint(logging.DebugLog, "WARN: mtLoadTBoot found mismatched trailer in TBOOT file\n")
			}
		}
	}
	simhtape.Rewind(mt.simhFile[tNum])
	logging.DebugPrint(logID, "... tape rewound\n")
	logging.DebugPrint(logID, "... mtLoadTBoot completed\n")
}

// MtDataIn is called from Bus to implement DIx from the mt device
func mtDataIn(abc byte, flag byte) (data dg.WordT) {

	mt.mtDataMu.RLock()
	switch abc {
	case 'A': /* Read status register 1 - see p.IV-18 of Peripherals guide */
		data = mt.statusReg1
		logging.DebugPrint(logID, "DIA - Read SR1 - returning: %#o = %s\n", data, mtReadableSR1())
	case 'B': /* Read memory addr register 1 - see p.IV-19 of Peripherals guide */
		data = dg.WordT(mt.memAddrReg)
		logging.DebugPrint(logID, "DIB - Read MA - returning: %#o\n", data)
	case 'C': /* Read status register 2 - see p.IV-18 of Peripherals guide */
		data = mt.statusReg2
		logging.DebugPrint(logID, "DIC - Read SR2 - returning: %#o\n", data)
	}
	mt.mtDataMu.RUnlock()

	mtHandleFlag(flag)

	return data
}

// MtDataOut is called from Bus to implement DOx from the mt device
func mtDataOut(datum dg.WordT, abc byte, flag byte) {
	mt.mtDataMu.Lock()
	switch abc {
	case 'A': // Specify Command and Drive - p.IV-17
		// which command?
		for c := 0; c < mtCmdCount; c++ {
			if (datum & mtCmdMask) == commandSet[c] {
				mt.currentCmd = c
				break
			}
		}
		// which unit?
		mt.currentUnit = mtExtractUnit(datum)
		logging.DebugPrint(logID, "DOA - Specify Command and Drive - internal cmd #: %0o, unit: %0o\n",
			mt.currentCmd, mt.currentUnit)
	case 'B':
		mt.memAddrReg = dg.PhysAddrT(datum)
		logging.DebugPrint(logID, "DOB - MA set to %0o\n", mt.memAddrReg)
	case 'C':
		mt.negWordCntReg = int16(datum)
		logging.DebugPrint(logID, "DOC - Set (neg) Word Count to 0%o (%d.)\n", datum, mt.negWordCntReg)
	case 'N': // special handling for NIOx...
		logging.DebugPrint(logID, "NIO - Flag is %c\n", flag)
	}
	mt.mtDataMu.Unlock()

	mtHandleFlag(flag)
}

func mtExtractUnit(word dg.WordT) int {
	return int(word & 0x07)
}

// mtHandleFlag actions the flag/pulse to the mt controller
func mtHandleFlag(f byte) {
	switch f {
	case 'S':
		logging.DebugPrint(logID, "... S flag set\n")
		mt.mtDataMu.RLock()
		if mt.currentCmd != mtCmdRewind {
			BusSetBusy(mt.devNum, true)
		}
		mt.mtDataMu.RUnlock()
		BusSetDone(mt.devNum, false)
		mtDoCommand()
		BusSetBusy(mt.devNum, false)
		BusSetDone(mt.devNum, true)

	case 'C':
		// if we were performing mt operations in a Goroutine, this would interrupt them...
		logging.DebugPrint(logID, "... C flag set\n")
		//mt.mtDataMu.Lock()
		//mt.statusReg1 = mtSr1HiDensity | mtSr19Track | mtSr1UnitReady // ???
		//mt.statusReg2 = mtSr2PEMode                                   // ???
		//mt.mtDataMu.Unlock()
		BusSetBusy(mt.devNum, false)
		BusSetDone(mt.devNum, false)

	case 'P':
		// 'Reserved'
		logging.DebugPrint(logID, "WARNING: Received 'P' flag - which is reserved")

	default:
		// empty flag - nothing to do
	}
}

func mtDoCommand() {
	mt.mtDataMu.Lock()
	defer mt.mtDataMu.Unlock()

	switch mt.currentCmd {
	case mtCmdRead:
		logging.DebugPrint(logID, "*READ* command\n ---- Unit: %d\n ---- Word Count: %d Location: %d\n", mt.currentUnit, mt.negWordCntReg, mt.memAddrReg)
		hdrLen, _ := simhtape.ReadMetaData(mt.simhFile[mt.currentUnit])
		logging.DebugPrint(logID, " ----  Header read giving length: %d\n", hdrLen)
		if hdrLen == mtEOF {
			logging.DebugPrint(logID, " ----  Header is EOF indicator\n")
			mt.statusReg1 = mtSr1HiDensity | mtSr19Track | mtSr1UnitReady | mtSr1EOF | mtSr1Error
		} else {
			logging.DebugPrint(logID, " ----  Calling simhtape.ReadRecord with length: %d\n", hdrLen)
			var w uint32
			var wd dg.WordT
			var pAddr dg.PhysAddrT
			rec, _ := simhtape.ReadRecordData(mt.simhFile[mt.currentUnit], int(hdrLen))
			for w = 0; w < hdrLen; w += 2 {
				wd = (dg.WordT(rec[w]) << 8) | dg.WordT(rec[w+1])
				pAddr = memory.WriteWordDchChan(&mt.memAddrReg, wd)
				logging.DebugPrint(logID, " ----  Written word (%02X | %02X := %04X) to logical address: %d, physical: %d\n", rec[w], rec[w+1], wd, mt.memAddrReg, pAddr)
				// memAddrReg is auto-incremented for every word written  *******
				// auto-incremement the (two's complement) word count
				mt.negWordCntReg++
				if mt.negWordCntReg == 0 {
					break
				}
			}
			trailer, _ := simhtape.ReadMetaData(mt.simhFile[mt.currentUnit])
			logging.DebugPrint(logID, " ----  %d bytes loaded\n", w)
			logging.DebugPrint(logID, " ----  Read SimH Trailer: %d\n", trailer)
			// TODO Need to verify how status should be set here...
			mt.statusReg1 = mtSr1HiDensity | mtSr19Track | mtSr1StatusChanged | mtSr1UnitReady
		}

	case mtCmdRewind:
		logging.DebugPrint(logID, "*REWIND* command\n ------ Unit: #%d\n", mt.currentUnit)
		simhtape.Rewind(mt.simhFile[mt.currentUnit])
		mt.statusReg1 = mtSr1Error | mtSr1HiDensity | mtSr19Track | mtSr1UnitReady | mtSr1StatusChanged | mtSr1BOT

	case mtCmdSpaceFwd:
		logging.DebugPrint(logID, "*SPACE FORWARD* command\n ----- ------- Unit: #%d\n", mt.currentUnit)
		if mt.negWordCntReg == 0 { // one whole file
			stat := simhtape.SpaceFwd(mt.simhFile[mt.currentUnit], mt.negWordCntReg)
			mt.memAddrReg = 0xffffffff //  or 0 ???
			if stat == simhtape.SimhMtStatOk {
				mt.statusReg1 = mtSr1HiDensity | mtSr19Track | mtSr1UnitReady | mtSr1EOF | mtSr1StatusChanged | mtSr1Error
			} else {
				mt.statusReg1 = mtSr1HiDensity | mtSr19Track | mtSr1UnitReady | mtSr1EOT | mtSr1StatusChanged | mtSr1Error
			}
		} else {
			stat := simhtape.SpaceFwd(mt.simhFile[mt.currentUnit], mt.negWordCntReg)
			switch stat {
			case simhtape.SimhMtStatOk:
				mt.memAddrReg = 0
				mt.statusReg1 = mtSr1HiDensity | mtSr19Track | mtSr1UnitReady | mtSr1StatusChanged
			case simhtape.SimhMtStatTmk:
				mt.statusReg1 = mtSr1HiDensity | mtSr19Track | mtSr1UnitReady | mtSr1EOF | mtSr1StatusChanged | mtSr1Error
			case simhtape.SimhMtStatInvRec:
				mt.statusReg1 = mtSr1HiDensity | mtSr19Track | mtSr1UnitReady | mtSr1DataError | mtSr1StatusChanged | mtSr1Error
			default:
				log.Fatalf("ERROR: Unexpected return from simhTape.SpaceFwd %d", stat)
			}
		}

	case mtCmdSpaceRev:
		log.Fatalln("ERROR: mtDoCommand - SPACE REVERSE command Not Yet Implemented")
	case mtCmdUnload:
		log.Fatalln("ERROR: mtDoCommand - UNLOAD command Not Yet Implemented")
	default:
		log.Fatalf("ERROR: mtDoCommand - Command #%d Not Yet Implemented\n", mt.currentCmd)
	}
}

func mtReadableSR1() (res string) {
	res = mtSr1Readable
	for b := 0; b < 16; b++ {
		if (mt.statusReg1 & (1 << (15 - uint8(b)))) == 0 {
			res = res[:b] + "-" + res[b+1:]
		}
	}
	return res
}
