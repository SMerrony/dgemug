// magtape6026.go

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

const mtStatsPeriodMs = 333 // Will update status this often

const maxTapes = 8

// MagTape6026T contains the current state of a type 6026 Magnetic Tape Drive
type MagTape6026T struct {
	mtMu                   sync.RWMutex
	bus                    *BusT
	devNum                 int
	imageAttached          [maxTapes]bool
	fileName               [maxTapes]string
	simhFile               [maxTapes]*os.File
	commandSet             [mtCmdCount]dg.WordT
	logID                  int
	debugLogging           bool
	statusReg1, statusReg2 dg.WordT
	memAddrReg             dg.PhysAddrT
	negWordCntReg          int16
	currentCmd             int
	currentUnit            int
}

// MtStatT holds data for the status collector
type MtStatT struct {
	// Internals
	ImageAttached [maxTapes]bool
	FileName      [maxTapes]string
	// DG device state
	MemAddrReg             dg.PhysAddrT
	CurrentCmd             int
	StatusReg1, StatusReg2 dg.WordT
}

// MtInit sets the initial state of the (unmounted) tape drive(s)
func (tape *MagTape6026T) MtInit(dev int, bus *BusT, statsChan chan MtStatT, logID int, debugLogging bool) bool {
	tape.mtMu.Lock()
	tape.devNum = dev
	tape.bus = bus
	tape.logID = logID
	tape.debugLogging = debugLogging
	tape.commandSet[mtCmdRead] = mtCmdReadBits
	tape.commandSet[mtCmdRewind] = mtCmdRewindBits
	tape.commandSet[mtCmdCtrlMode] = mtCmdCtrlModeBits
	tape.commandSet[mtCmdSpaceFwd] = mtCmdSpaceFwdBits
	tape.commandSet[mtCmdSpaceRev] = mtCmdSpaceRevBits
	tape.commandSet[mtCmdWrite] = mtCmdWiteBits
	tape.commandSet[mtCmdWriteEOF] = mtCmdWriteEOFBits
	tape.commandSet[mtCmdErase] = mtCmdEraseBits
	tape.commandSet[mtCmdReadNonStop] = mtCmdReadNonStopBits
	tape.commandSet[mtCmdUnload] = mtCmdUnloadBits
	tape.commandSet[mtCmdDriveMode] = mtCmdDriveModeBits

	tape.bus.SetResetFunc(tape.devNum, tape.MtReset)
	tape.bus.SetDataInFunc(tape.devNum, tape.mtDataIn)
	tape.bus.SetDataOutFunc(tape.devNum, tape.mtDataOut)

	go tape.mtStatSender(statsChan)

	tape.statusReg1 = mtSr1HiDensity | mtSr19Track | mtSr1UnitReady
	tape.statusReg2 = mtSr2PEMode
	tape.mtMu.Unlock()

	logging.DebugPrint(tape.logID, "mt Initialised via call to mtInit()\n")
	return true
}

// mtStatSender provides a near real-time view of mt status and should be run as a Goroutine
// TODO = only handles Unit 0 at the moment
func (tape *MagTape6026T) mtStatSender(sChan chan MtStatT) {
	var stats MtStatT
	logging.DebugPrint(logging.DebugLog, "dskpStatSender() started\n")
	for {
		tape.mtMu.RLock()
		if tape.imageAttached[0] {
			stats.ImageAttached[0] = true
			stats.FileName[0] = tape.fileName[0]
			stats.MemAddrReg = tape.memAddrReg
			stats.CurrentCmd = tape.currentCmd // Could decode this
			stats.StatusReg1 = tape.statusReg1
			stats.StatusReg2 = tape.statusReg2
		} else {
			stats = MtStatT{}
		}
		tape.mtMu.RUnlock()
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
func (tape *MagTape6026T) MtReset() {
	tape.mtMu.Lock()
	for t := 0; t < maxTapes; t++ {
		if tape.imageAttached[t] {
			simhtape.Rewind(tape.simhFile[t])
		}
	}
	// BOT is an error state...
	tape.statusReg1 = mtSr1Error | mtSr1HiDensity | mtSr19Track | mtSr1BOT | mtSr1StatusChanged | mtSr1UnitReady
	tape.statusReg2 = mtSr2PEMode
	tape.memAddrReg = 0
	tape.negWordCntReg = 0
	tape.currentCmd = 0
	tape.currentUnit = 0
	tape.mtMu.Unlock()
	logging.DebugPrint(tape.logID, "mt Reset via call to mtReset()\n")
}

// MtAttach attaches a SimH tape image file to the emulated tape drive
func (tape *MagTape6026T) MtAttach(tNum int, imgName string) bool {
	logging.DebugPrint(tape.logID, "mtAttach called on unit #%d with image file: %s\n", tNum, imgName)
	f, err := os.Open(imgName)
	if err != nil {
		logging.DebugPrint(tape.logID, "ERROR: Could not open simH Tape Image file: %s, due to: %s\n", imgName, err.Error())
		return false
	}
	tape.mtMu.Lock()
	tape.fileName[tNum] = imgName
	tape.simhFile[tNum] = f
	tape.imageAttached[tNum] = true
	tape.statusReg1 = mtSr1Error | mtSr1HiDensity | mtSr19Track | mtSr1BOT | mtSr1StatusChanged | mtSr1UnitReady
	tape.statusReg2 = mtSr2PEMode
	tape.mtMu.Unlock()
	tape.bus.SetAttached(tape.devNum, imgName)
	return true

}

// MtDetach disassociates a tape file image from the drive
func (tape *MagTape6026T) MtDetach(tNum int) bool {
	logging.DebugPrint(tape.logID, "mtDetach called on unit #%d\n", tNum)
	tape.mtMu.Lock()
	tape.fileName[tNum] = ""
	tape.simhFile[tNum] = nil
	tape.imageAttached[tNum] = false
	tape.statusReg1 = mtSr1Error | mtSr1HiDensity | mtSr19Track | mtSr1BOT | mtSr1StatusChanged | mtSr1UnitReady
	tape.statusReg2 = mtSr2PEMode
	tape.mtMu.Unlock()
	tape.bus.SetDetached(tape.devNum)
	return true
}

// MtScanImage scans an attached SimH tape image to ensure it makes sense
// (This is just a pass-through to the equivalent function in simhtape)
func (tape *MagTape6026T) MtScanImage(tNum int) string {
	tape.mtMu.RLock()
	imageName := tape.fileName[tNum]
	att := tape.imageAttached[tNum]
	tape.mtMu.RUnlock()
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
func (tape *MagTape6026T) MtLoadTBoot() {
	var (
		byte0, byte1 byte
		word         dg.WordT
		wdix, memix  dg.PhysAddrT
	)
	logging.DebugPrint(tape.logID, "mtLoadTBoot() called\n")
	tNum := 0
	tape.mtMu.Lock()
	defer tape.mtMu.Unlock()
	simhtape.Rewind(tape.simhFile[tNum])
	logging.DebugPrint(tape.logID, "... tape rewound\n")

readLoop:
	for {
		hdr, ok := simhtape.ReadMetaData(tape.simhFile[tNum])
		// if !ok || hdr != tbootSizeB {
		if !ok {
			logging.DebugPrint(logging.DebugLog, "WARN: mtLoadTBoot called when no bootable tape image attached\n")
			return
		}
		logging.DebugPrint(tape.logID, "... header read, size is %d\n", hdr)
		switch hdr {
		case simhtape.SimhMtrTmk: // Tape Mark (separates files)
			break readLoop
		default:
			tbootSizeW := hdr / 2
			tapeData, ok := simhtape.ReadRecordData(tape.simhFile[tNum], int(hdr))
			if ok {
				logging.DebugPrint(tape.logID, "... data read\n")
			} else {
				logging.DebugPrint(tape.logID, "... error reading data\n")
				logging.DebugPrint(logging.DebugLog, "WARNING: Could not read data in mtLoadTBoot()\n")
				return
			}

			// logging.DebugPrint(tape.logID, "... loading data into memory starting at address %d\n", memix)
			for wdix = 0; wdix < dg.PhysAddrT(tbootSizeW); wdix++ {
				byte1 = tapeData[wdix*2]
				byte0 = tapeData[wdix*2+1]
				word = dg.WordT(byte1)<<8 | dg.WordT(byte0)
				memory.WriteWord(memix+wdix, word)
			}
			memix += dg.PhysAddrT(tbootSizeW)
			logging.DebugPrint(tape.logID, "... finished loading data at address %d\n", memix+wdix)
			trailer, ok := simhtape.ReadMetaData(tape.simhFile[tNum])
			if hdr != trailer || !ok {
				logging.DebugPrint(logging.DebugLog, "WARN: mtLoadTBoot found mismatched trailer in TBOOT file\n")
			}
		}
	}
	simhtape.Rewind(tape.simhFile[tNum])
	logging.DebugPrint(tape.logID, "... tape rewound\n")
	logging.DebugPrint(tape.logID, "... mtLoadTBoot completed\n")
}

// MtDataIn is called from Bus to implement DIx from the mt device
func (tape *MagTape6026T) mtDataIn(abc byte, flag byte) (data dg.WordT) {

	tape.mtMu.RLock()
	switch abc {
	case 'A': /* Read status register 1 - see p.IV-18 of Peripherals guide */
		data = tape.statusReg1
		logging.DebugPrint(tape.logID, "DIA - Read SR1 - returning: %#o = %s\n", data, tape.mtReadableSR1())
	case 'B': /* Read memory addr register 1 - see p.IV-19 of Peripherals guide */
		data = dg.WordT(tape.memAddrReg)
		logging.DebugPrint(tape.logID, "DIB - Read MA - returning: %#o\n", data)
	case 'C': /* Read status register 2 - see p.IV-18 of Peripherals guide */
		data = tape.statusReg2
		logging.DebugPrint(tape.logID, "DIC - Read SR2 - returning: %#o\n", data)
	}
	tape.mtMu.RUnlock()

	tape.mtHandleFlag(flag)

	return data
}

// MtDataOut is called from Bus to implement DOx from the mt device
func (tape *MagTape6026T) mtDataOut(datum dg.WordT, abc byte, flag byte) {
	tape.mtMu.Lock()
	switch abc {
	case 'A': // Specify Command and Drive - p.IV-17
		// which command?
		for c := 0; c < mtCmdCount; c++ {
			if (datum & mtCmdMask) == tape.commandSet[c] {
				tape.currentCmd = c
				break
			}
		}
		// which unit?
		tape.currentUnit = mtExtractUnit(datum)
		logging.DebugPrint(tape.logID, "DOA - Specify Command and Drive - internal cmd #: %#o, unit: %#o\n",
			tape.currentCmd, tape.currentUnit)
	case 'B':
		tape.memAddrReg = dg.PhysAddrT(datum)
		logging.DebugPrint(tape.logID, "DOB - MA set to %#o\n", tape.memAddrReg)
	case 'C':
		tape.negWordCntReg = int16(datum)
		logging.DebugPrint(tape.logID, "DOC - Set (neg) Word Count to #%o (%d.)\n", datum, tape.negWordCntReg)
	case 'N': // special handling for NIOx...
		logging.DebugPrint(tape.logID, "NIO - Flag is %c\n", flag)
	}
	tape.mtMu.Unlock()

	tape.mtHandleFlag(flag)
}

func mtExtractUnit(word dg.WordT) int {
	return int(word & 0x07)
}

// mtHandleFlag actions the flag/pulse to the mt controller
func (tape *MagTape6026T) mtHandleFlag(f byte) {
	switch f {
	case 'S':
		logging.DebugPrint(tape.logID, "... S flag set\n")
		tape.mtMu.RLock()
		if tape.currentCmd != mtCmdRewind {
			tape.bus.SetBusy(tape.devNum, true)
		}
		tape.mtMu.RUnlock()
		tape.bus.SetDone(tape.devNum, false)
		tape.mtDoCommand()
		tape.bus.SetBusy(tape.devNum, false)
		tape.bus.SetDone(tape.devNum, true)

	case 'C':
		// if we were performing mt operations in a Goroutine, this would interrupt them...
		logging.DebugPrint(tape.logID, "... C flag set\n")
		//tape.mtMu.Lock()
		//tape.statusReg1 = mtSr1HiDensity | mtSr19Track | mtSr1UnitReady // ???
		//tape.statusReg2 = mtSr2PEMode                                   // ???
		//tape.mtMu.Unlock()
		tape.bus.SetBusy(tape.devNum, false)
		tape.bus.SetDone(tape.devNum, false)

	case 'P':
		// 'Reserved'
		logging.DebugPrint(tape.logID, "WARNING: Received 'P' flag - which is reserved")

	default:
		// empty flag - nothing to do
	}
}

func (tape *MagTape6026T) mtDoCommand() {
	tape.mtMu.Lock()
	defer tape.mtMu.Unlock()

	switch tape.currentCmd {
	case mtCmdRead:
		logging.DebugPrint(tape.logID, "*READ* command ---- Unit: %d ---- Word Count: %d Location: %d\n", tape.currentUnit, tape.negWordCntReg, tape.memAddrReg)
		hdrLen, _ := simhtape.ReadMetaData(tape.simhFile[tape.currentUnit])
		logging.DebugPrint(tape.logID, " ----  Header read giving length: %d\n", hdrLen)
		if hdrLen == mtEOF {
			logging.DebugPrint(tape.logID, " ----  Header is EOF indicator\n")
			tape.statusReg1 = mtSr1HiDensity | mtSr19Track | mtSr1UnitReady | mtSr1EOF | mtSr1Error
		} else {
			// logging.DebugPrint(tape.logID, " ----  Calling simhtape.ReadRecord with length: %d\n", hdrLen)
			var w uint32
			var wd dg.WordT
			var pAddr dg.PhysAddrT
			rec, _ := simhtape.ReadRecordData(tape.simhFile[tape.currentUnit], int(hdrLen))
			for w = 0; w < hdrLen; w += 2 {
				wd = (dg.WordT(rec[w]) << 8) | dg.WordT(rec[w+1])
				pAddr = memory.WriteWordDchChan(&tape.memAddrReg, wd)
				if w == 0 || w == (hdrLen-2) {
					logging.DebugPrint(tape.logID, " ----  Written word %#04x to logical address: %#o, physical: %#o\n", wd, tape.memAddrReg-1, pAddr)
				}
				// memAddrReg is auto-incremented for every word written  *******
				// auto-incremement the (two's complement) word count
				tape.negWordCntReg++
				if tape.negWordCntReg == 0 {
					break
				}
			}
			trailer, _ := simhtape.ReadMetaData(tape.simhFile[tape.currentUnit])
			logging.DebugPrint(tape.logID, " ----  %d bytes loaded\n", w)
			logging.DebugPrint(tape.logID, " ----  Read SimH Trailer: %d\n", trailer)
			// TODO Need to verify how status should be set here...
			tape.statusReg1 = mtSr1HiDensity | mtSr19Track | mtSr1StatusChanged | mtSr1UnitReady
		}

	case mtCmdRewind:
		logging.DebugPrint(tape.logID, "*REWIND* command\n ------ Unit: #%d\n", tape.currentUnit)
		simhtape.Rewind(tape.simhFile[tape.currentUnit])
		tape.statusReg1 = mtSr1Error | mtSr1HiDensity | mtSr19Track | mtSr1UnitReady | mtSr1StatusChanged | mtSr1BOT

	case mtCmdSpaceFwd:
		logging.DebugPrint(tape.logID, "*SPACE FORWARD* command\n ----- ------- Unit: #%d\n", tape.currentUnit)
		if tape.negWordCntReg == 0 { // one whole file
			stat := simhtape.SpaceFwd(tape.simhFile[tape.currentUnit], tape.negWordCntReg)
			// according to the simH source, MA should be set to # files/recs skipped
			// can't find any reference to this in the Periph Pgmrs Guide but it lets INSTL
			// progress further...
			// It seems to need the two's complement of the number...
			tape.memAddrReg = 0xffffffff
			if stat == simhtape.SimhMtStatOk {
				//tape.statusReg1 = mtSr1HiDensity | mtSr19Track | mtSr1UnitReady | mtSr1EOF | mtSr1StatusChanged | mtSr1Error
				tape.statusReg1 = mtSr1HiDensity | mtSr19Track | mtSr1UnitReady | mtSr1EOF | mtSr1Error
			} else {
				tape.statusReg1 = mtSr1HiDensity | mtSr19Track | mtSr1UnitReady | mtSr1EOT | mtSr1StatusChanged | mtSr1Error
			}
		} else {
			stat := simhtape.SpaceFwd(tape.simhFile[tape.currentUnit], tape.negWordCntReg)
			switch stat {
			case simhtape.SimhMtStatOk:
				//tape.memAddrReg = 0
				tape.statusReg1 = mtSr1HiDensity | mtSr19Track | mtSr1UnitReady | mtSr1StatusChanged
			case simhtape.SimhMtStatTmk:
				tape.statusReg1 = mtSr1HiDensity | mtSr19Track | mtSr1UnitReady | mtSr1EOF | mtSr1StatusChanged | mtSr1Error
			case simhtape.SimhMtStatInvRec:
				tape.statusReg1 = mtSr1HiDensity | mtSr19Track | mtSr1UnitReady | mtSr1DataError | mtSr1StatusChanged | mtSr1Error
			default:
				log.Fatalf("ERROR: Unexpected return from simhTape.SpaceFwd %d", stat)
			}
			tape.memAddrReg = dg.PhysAddrT(tape.negWordCntReg)
		}

	case mtCmdSpaceRev:
		log.Fatalln("ERROR: mtDoCommand - SPACE REVERSE command Not Yet Implemented")
	case mtCmdUnload:
		log.Fatalln("ERROR: mtDoCommand - UNLOAD command Not Yet Implemented")
	default:
		log.Fatalf("ERROR: mtDoCommand - Command #%d Not Yet Implemented\n", tape.currentCmd)
	}
}

func (tape *MagTape6026T) mtReadableSR1() (res string) {
	res = mtSr1Readable
	for b := 0; b < 16; b++ {
		if (tape.statusReg1 & (1 << (15 - uint8(b)))) == 0 {
			res = res[:b] + "-" + res[b+1:]
		}
	}
	return res
}
