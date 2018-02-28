// disk6239.go

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

// Here we are emulating the disk6239 device, specifically model 6239/6240
// controller/drive combination with 14-inch platters which provide 592MB of formatted capacity.
//
// All communication with the drive is via CPU PIO instructions and memory
// accessed via the BMC interface running at 2.2MB/sec in mapped or physical mode.
// There is also a small set of flags and pulses shared between the controller and the CPU.

// ASYNCHRONOUS interrupts occur on completion of a CB (list), or when an error
// occurs during CB processing.

// SYNCHRONOUS interrupts occur after a PIO command executes.

// N.B. Assembler mnemonic: DSKP, AOS/VS mnemonic: DPJ

package devices

import (
	"bufio"
	"log"
	"os"
	"sync"
	"time"

	"github.com/SMerrony/dgemug/logging"

	"github.com/SMerrony/dgemug/util"

	"github.com/SMerrony/dgemug/memory"

	"github.com/SMerrony/dgemug"
)

const (
	// Physical disk characteristics
	disk6239SurfacesPerDisk   = 8
	disk6239HeadsPerSurface   = 2
	disk6239SectorsPerTrack   = 75
	disk6239WordsPerSector    = 256
	disk6239BytesPerSector    = disk6239WordsPerSector * 2
	disk6239PhysicalCylinders = 981
	disk6239UserCylinders     = 978
	disk6239LogicalBlocks     = 1157952 // ??? 1147943 17<<16 | 43840
	disk6239LogicalBlocksH    = disk6239LogicalBlocks >> 16
	disk6239LogicalBlocksL    = disk6239LogicalBlocks & 0x0ffff
	disk6239UcodeRev          = 99

	disk6239MaxQueuedCBs = 30 // See p.2-13

	disk6239IntInfBlkSize   = 8
	disk6239CtrlrInfBlkSize = 2
	disk6239UnitInfBlkSize  = 7
	disk6239CbMaxSize       = 21
	disk6239CbMinSize       = 10 //12 // Was 10

	disk6239AsynchStatRetryInterval = time.Millisecond

	statXecStateResetting = 0x00
	statXecStateResetDone = 0x01
	statXecStateBegun     = 0x08
	statXecStateMapped    = 0x0c
	statXecStateDiagMode  = 0x04

	statCcsAsync        = 0
	statCcsPioInvCmd    = 1
	statCcsPioCmdFailed = 2
	statCcsPioCmdOk     = 3

	statAsyncNoErrors = 5

	// disk6239 PIO Command Set
	disk6239PioProgLoad        = 000
	disk6239PioBegin           = 002
	disk6239PioSysgen          = 025
	disk6239DiagMode           = 024
	disk6239SetMapping         = 026
	disk6239GetMapping         = 027
	disk6239SetInterface       = 030
	disk6239GetInterface       = 031
	disk6239SetController      = 032
	disk6239GetController      = 033
	disk6239SetUnit            = 034
	disk6239GetUnit            = 035
	disk6239GetExtendedStatus0 = 040
	disk6239GetExtendedStatus1 = 041
	disk6239GetExtendedStatus2 = 042
	disk6239GetExtendedStatus3 = 043
	disk6239StartList          = 0100
	disk6239StartListHp        = 0103
	disk6239Restart            = 0116
	disk6239CancelList         = 0123
	disk6239UnitStatus         = 0131
	disk6239Trespass           = 0132
	disk6239GetListStatus      = 0133
	disk6239PioReset           = 0777

	// disk6239 CB Command Set/OpCodes
	disk6239CbOpNoOp             = 0
	disk6239CbOpWrite            = 0100
	disk6239CbOpWriteVerify      = 0101
	disk6239CbOpWrite1Word       = 0104
	disk6239CbOpWriteVerify1Word = 0105
	disk6239CbOpWriteModBitmap   = 0142
	disk6239CbOpRead             = 0200
	disk6239CbOpReadVerify       = 0201
	disk6239CbOpReadVerify1Word  = 0205
	disk6239CbOpReadRawData      = 0210
	disk6239CbOpReadHeaders      = 0220
	disk6239CbOpReadModBitmap    = 0242
	disk6239CbOpRecalibrateDisk  = 0400

	// disk6239 CB FIELDS
	disk6239CbLINK_ADDR_HIGH        = 0
	disk6239CbLINK_ADDR_LOW         = 1
	disk6239CbINA_FLAGS_OPCODE      = 2
	disk6239CbPAGENO_LIST_ADDR_HIGH = 3
	disk6239CbPAGENO_LIST_ADDR_LOW  = 4
	disk6239CbTXFER_ADDR_HIGH       = 5
	disk6239CbTXFER_ADDR_LOW        = 6
	disk6239CbDEV_ADDR_HIGH         = 7
	disk6239CbDEV_ADDR_LOW          = 8
	disk6239CbUNIT_NO               = 9
	disk6239CbTXFER_COUNT           = 10
	disk6239CbCB_STATUS             = 11
	disk6239CbRES1                  = 12
	disk6239CbRES2                  = 13
	disk6239CbERR_STATUS            = 14
	disk6239CbUNIT_STATUS           = 15
	disk6239CbRETRIES_DONE          = 16
	disk6239CbSOFT_RTN_TXFER_COUNT  = 17
	disk6239CbPHYS_CYL              = 18
	disk6239CbPHYS_HEAD_SECT        = 19
	disk6239CbDISK_ERR_CODE         = 20

	// Mapping bits
	disk6239MapSlotLoadInts = 1 << 15
	disk6239MapIntBmcPhys   = 1 << 14
	disk6239MapUpstreamLoad = 1 << 13
	disk6239MapUpstreamHpt  = 1 << 12

	// calculated consts
	// disk6239PhysicalByteSize is the total  # bytes on a disk6239-type disk
	disk6239PhysicalByteSize = disk6239SurfacesPerDisk * disk6239HeadsPerSurface * disk6239SectorsPerTrack * disk6239BytesPerSector * disk6239PhysicalCylinders
	// disk6239PhysicalBlockSize is the total # blocks on a disk6239-type disk
	disk6239PhysicalBlockSize = disk6239SurfacesPerDisk * disk6239HeadsPerSurface * disk6239SectorsPerTrack * disk6239PhysicalCylinders
)

type disk6239DataT struct {
	// MV/Em internals...
	disk6239DataMu sync.RWMutex
	devNum         int
	imageAttached  bool
	imageFileName  string
	imageFile      *os.File
	reads, writes  uint64
	logID          int
	// DG data...
	commandRegA, commandRegB, commandRegC dg.WordT
	statusRegA, statusRegB, statusRegC    dg.WordT
	isMapped                              bool
	mappingRegA, mappingRegB              dg.WordT
	intInfBlock                           [disk6239IntInfBlkSize]dg.WordT
	ctrlInfBlock                          [disk6239CtrlrInfBlkSize]dg.WordT
	unitInfBlock                          [disk6239UnitInfBlkSize]dg.WordT
	// cylinder, head, sector                dg_word
	sectorNo dg.DwordT
}

const disk6239StatsPeriodMs = 500 // Will send status update this often

type Disk6239StatT struct {
	ImageAttached                      bool
	StatusRegA, StatusRegB, StatusRegC dg.WordT
	//	cylinder, head, sector             dg_word
	SectorNo      dg.DwordT
	Reads, Writes uint64
}

var (
	disk6239Data disk6239DataT
	cbChan       chan dg.PhysAddrT
)

// disk6239Init is called once by the main routine to initialise this disk6239 emulator
func Disk6239Init(dev int, statsChann chan Disk6239StatT, logId int, logging bool) {

	disk6239Data.devNum = dev

	go disk6239StatSender(statsChann)

	BusSetResetFunc(disk6239Data.devNum, disk6239Reset)
	BusSetDataInFunc(disk6239Data.devNum, disk6239DataIn)
	BusSetDataOutFunc(disk6239Data.devNum, disk6239DataOut)

	disk6239Data.disk6239DataMu.Lock()
	disk6239Data.logID = logId

	disk6239Data.imageAttached = false
	disk6239Data.disk6239DataMu.Unlock()
	cbChan = make(chan dg.PhysAddrT, disk6239MaxQueuedCBs)
	go disk6239CBprocessor(&disk6239Data)

	disk6239Reset()
}

// attempt to attach an extant MV/Em disk image to the running emulator
func Disk6239Attach(dNum int, imgName string) bool {
	// TODO Disk Number not currently used
	logging.DebugPrint(disk6239Data.logID, "disk6239Attach called for disk #%d with image <%s>\n", dNum, imgName)

	disk6239Data.disk6239DataMu.Lock()

	disk6239Data.imageFile, err = os.OpenFile(imgName, os.O_RDWR, 0755)
	if err != nil {
		logging.DebugPrint(disk6239Data.logID, "Failed to open image for attaching\n")
		logging.DebugPrint(logging.DebugLog, "WARN: Failed to open disk6239 image <%s> for ATTach\n", imgName)
		return false
	}
	disk6239Data.imageFileName = imgName
	disk6239Data.imageAttached = true

	disk6239Data.disk6239DataMu.Unlock()

	BusSetAttached(disk6239Data.devNum, imgName)
	return true
}

// disk6239StatSender provides a near real-time view of the disk6239 status and should be run as a Goroutine
func disk6239StatSender(sChan chan Disk6239StatT) {
	var stats Disk6239StatT
	logging.DebugPrint(logging.DebugLog, "disk6239StatSender() started\n")
	for {
		disk6239Data.disk6239DataMu.RLock()
		if disk6239Data.imageAttached {
			stats.ImageAttached = true
			//stats.cylinder = disk6239Data.cylinder
			//stats.head = disk6239Data.head
			//stats.sector = disk6239Data.sector
			stats.StatusRegA = disk6239Data.statusRegA
			stats.StatusRegB = disk6239Data.statusRegB
			stats.StatusRegC = disk6239Data.statusRegC
			stats.SectorNo = disk6239Data.sectorNo
			stats.Reads = disk6239Data.reads
			stats.Writes = disk6239Data.writes
		} else {
			stats = Disk6239StatT{}
		}
		disk6239Data.disk6239DataMu.RUnlock()
		// Non-blocking send of stats
		select {
		case sChan <- stats:
		default:
		}
		time.Sleep(time.Millisecond * disk6239StatsPeriodMs)
	}
}

// Create an empty disk file of the correct size for the disk6239 emulator to use
func Disk6239CreateBlank(imgName string) bool {
	newFile, err := os.Create(imgName)
	if err != nil {
		return false
	}
	defer newFile.Close()
	logging.DebugPrint(disk6239Data.logID, "disk6239CreateBlank attempting to write %d bytes\n", disk6239PhysicalByteSize)
	w := bufio.NewWriter(newFile)
	for b := 0; b < disk6239PhysicalByteSize; b++ {
		w.WriteByte(0)
	}
	w.Flush()
	return true
}

// Handle the DIA/B/C PIO commands
func disk6239DataIn(abc byte, flag byte) (datum dg.WordT) {
	disk6239Data.disk6239DataMu.Lock()
	switch abc {
	case 'A':
		datum = disk6239Data.statusRegA
		if debugLogging {
			logging.DebugPrint(disk6239Data.logID, "DIA [Read Status A] returning %s for DRV=%d\n", util.WordToBinStr(disk6239Data.statusRegA), 0)
		}
	case 'B':
		datum = disk6239Data.statusRegB
		if debugLogging {
			logging.DebugPrint(disk6239Data.logID, "DIB [Read Status B] returning %s for DRV=%d\n", util.WordToBinStr(disk6239Data.statusRegB), 0)
		}
	case 'C':
		datum = disk6239Data.statusRegC
		if debugLogging {
			logging.DebugPrint(disk6239Data.logID, "DIC [Read Status C] returning %s for DRV=%d\n", util.WordToBinStr(disk6239Data.statusRegC), 0)
		}
	}
	disk6239Data.disk6239DataMu.Unlock()
	disk6239HandleFlag(flag)
	return datum
}

// Handle the DOA/B/C PIO commands
func disk6239DataOut(datum dg.WordT, abc byte, flag byte) {
	disk6239Data.disk6239DataMu.Lock()
	switch abc {
	case 'A':
		disk6239Data.commandRegA = datum
	case 'B':
		disk6239Data.commandRegB = datum
	case 'C':
		disk6239Data.commandRegC = datum
	}
	disk6239Data.disk6239DataMu.Unlock()
	disk6239HandleFlag(flag)
}

func disk6239DoPioCommand() {

	var addr, w dg.PhysAddrT

	disk6239Data.disk6239DataMu.Lock()

	pioCmd := disk6239ExtractPioCommand(disk6239Data.commandRegC)
	switch pioCmd {
	case disk6239PioProgLoad:
		log.Panicln("disk6239ProgLoad command not yet implemented")

	case disk6239PioBegin:
		if debugLogging {
			logging.DebugPrint(disk6239Data.logID, "... BEGIN command, unit # %d\n", disk6239Data.commandRegA)
		}
		// pretend we have succesfully booted ourself
		disk6239Data.statusRegB = 0
		disk6239SetPioStatusRegC(statXecStateBegun, statCcsPioCmdOk, disk6239PioBegin, util.TestWbit(disk6239Data.commandRegC, 15))

	case disk6239GetMapping:
		if debugLogging {
			logging.DebugPrint(disk6239Data.logID, "... GET MAPPING command\n")
		}
		disk6239Data.statusRegA = disk6239Data.mappingRegA
		disk6239Data.statusRegB = disk6239Data.mappingRegB
		if debugLogging {
			logging.DebugPrint(disk6239Data.logID, "... ... Status Reg A set to %s\n", util.WordToBinStr(disk6239Data.statusRegA))
			logging.DebugPrint(disk6239Data.logID, "... ... Status Reg B set to %s\n", util.WordToBinStr(disk6239Data.statusRegB))
		}
		disk6239SetPioStatusRegC(0, statCcsPioCmdOk, disk6239GetMapping, util.TestWbit(disk6239Data.commandRegC, 15))

	case disk6239SetMapping:
		if debugLogging {
			logging.DebugPrint(disk6239Data.logID, "... SET MAPPING command\n")
		}
		disk6239Data.mappingRegA = disk6239Data.commandRegA
		disk6239Data.mappingRegB = disk6239Data.commandRegB
		disk6239Data.isMapped = true
		if debugLogging {
			logging.DebugPrint(disk6239Data.logID, "... ... Mapping Reg A set to %s\n", util.WordToBinStr(disk6239Data.commandRegA))
			logging.DebugPrint(disk6239Data.logID, "... ... Mapping Reg B set to %s\n", util.WordToBinStr(disk6239Data.commandRegB))
		}
		disk6239SetPioStatusRegC(statXecStateMapped, statCcsPioCmdOk, disk6239SetMapping, util.TestWbit(disk6239Data.commandRegC, 15))

	case disk6239GetInterface:
		addr = dg.PhysAddrT(util.DwordFromTwoWords(disk6239Data.commandRegA, disk6239Data.commandRegB))
		if debugLogging {
			logging.DebugPrint(disk6239Data.logID, "... GET INTERFACE INFO command\n")
			logging.DebugPrint(disk6239Data.logID, "... ... Destination Start Address: %d\n", addr)
		}
		for w = 0; w < disk6239IntInfBlkSize; w++ {
			memory.WriteWordBmcChan(&addr, disk6239Data.intInfBlock[w])
			if debugLogging {
				logging.DebugPrint(disk6239Data.logID, "... ... Word %d: %s\n", w, util.WordToBinStr(disk6239Data.intInfBlock[w]))
			}
		}
		disk6239SetPioStatusRegC(0, statCcsPioCmdOk, disk6239GetInterface, util.TestWbit(disk6239Data.commandRegC, 15))

	case disk6239SetInterface:
		addr = dg.PhysAddrT(util.DwordFromTwoWords(disk6239Data.commandRegA, disk6239Data.commandRegB))
		if debugLogging {
			logging.DebugPrint(disk6239Data.logID, "... SET INTERFACE INFO command\n")
			logging.DebugPrint(disk6239Data.logID, "... ... Origin Start Address: %d\n", addr)
		}
		// only a few fields can be changed...
		addr += 5
		disk6239Data.intInfBlock[w] = memory.ReadWordBmcChan(&addr) // word 5
		disk6239Data.intInfBlock[w] &= 0xff00
		disk6239Data.intInfBlock[w] = memory.ReadWordBmcChan(&addr) // word 6
		disk6239Data.intInfBlock[w] = memory.ReadWordBmcChan(&addr) // word 7
		if debugLogging {
			logging.DebugPrint(disk6239Data.logID, "... ... Word 5: %s\n", util.WordToBinStr(disk6239Data.intInfBlock[5]))
			logging.DebugPrint(disk6239Data.logID, "... ... Word 6: %s\n", util.WordToBinStr(disk6239Data.intInfBlock[6]))
			logging.DebugPrint(disk6239Data.logID, "... ... Word 7: %s\n", util.WordToBinStr(disk6239Data.intInfBlock[7]))
		}
		disk6239SetPioStatusRegC(0, statCcsPioCmdOk, disk6239SetInterface, util.TestWbit(disk6239Data.commandRegC, 15))

	case disk6239GetUnit:
		addr = dg.PhysAddrT(util.DwordFromTwoWords(disk6239Data.commandRegA, disk6239Data.commandRegB))
		if debugLogging {
			logging.DebugPrint(disk6239Data.logID, "... GET UNIT INFO command\n")
			logging.DebugPrint(disk6239Data.logID, "... ... Destination Start Address: %d\n", addr)
		}
		for w = 0; w < disk6239UnitInfBlkSize; w++ {
			memory.WriteWordBmcChan(&addr, disk6239Data.unitInfBlock[w])
			if debugLogging {
				logging.DebugPrint(disk6239Data.logID, "... ... Word %d: %s\n", w, util.WordToBinStr(disk6239Data.unitInfBlock[w]))
			}
		}
		disk6239SetPioStatusRegC(0, statCcsPioCmdOk, disk6239GetUnit, util.TestWbit(disk6239Data.commandRegC, 15))

	case disk6239SetUnit:
		addr = dg.PhysAddrT(util.DwordFromTwoWords(disk6239Data.commandRegA, disk6239Data.commandRegB))
		if debugLogging {
			logging.DebugPrint(disk6239Data.logID, "... SET UNIT INFO command\n")
			logging.DebugPrint(disk6239Data.logID, "... ... Origin Start Address: %d\n", addr)
		}
		// only the first word is writable according to p.2-16
		// TODO check no active CBs first
		disk6239Data.unitInfBlock[0] = memory.ReadWord(addr)
		if debugLogging {
			logging.DebugPrint(disk6239Data.logID, "... ... Overwrote word 0 of UIB with: %s\n", util.WordToBinStr(disk6239Data.unitInfBlock[0]))
		}
		disk6239SetPioStatusRegC(0, statCcsPioCmdOk, disk6239SetUnit, util.TestWbit(disk6239Data.commandRegC, 15))

	case disk6239PioReset:
		// disk6239Reset() has to do its own locking...
		disk6239Data.disk6239DataMu.Unlock()
		disk6239Reset()
		disk6239Data.disk6239DataMu.Lock()

	case disk6239SetController:
		addr = dg.PhysAddrT(util.DwordFromTwoWords(disk6239Data.commandRegA, disk6239Data.commandRegB))
		if debugLogging {
			logging.DebugPrint(disk6239Data.logID, "... SET CONTROLLER INFO command\n")
			logging.DebugPrint(disk6239Data.logID, "... ... Origin Start Address: %d\n", addr)
		}
		disk6239Data.ctrlInfBlock[0] = memory.ReadWord(addr)
		disk6239Data.ctrlInfBlock[1] = memory.ReadWord(addr + 1)
		if debugLogging {
			logging.DebugPrint(disk6239Data.logID, "... ... Word 0: %s\n", util.WordToBinStr(disk6239Data.ctrlInfBlock[0]))
			logging.DebugPrint(disk6239Data.logID, "... ... Word 1: %s\n", util.WordToBinStr(disk6239Data.ctrlInfBlock[1]))
		}
		disk6239SetPioStatusRegC(0, statCcsPioCmdOk, disk6239SetController, util.TestWbit(disk6239Data.commandRegC, 15))

	case disk6239StartList:
		addr = dg.PhysAddrT(util.DwordFromTwoWords(disk6239Data.commandRegA, disk6239Data.commandRegB))
		if debugLogging {
			logging.DebugPrint(disk6239Data.logID, "... START LIST command\n")
			logging.DebugPrint(disk6239Data.logID, "... ..... First CB Address: %d\n", addr)
			logging.DebugPrint(disk6239Data.logID, "... ..... CB Channel Q length: %d\n", len(cbChan))
		}
		// TODO should check addr validity before starting processing
		//disk6239ProcessCB(addr)
		cbChan <- addr
		disk6239Data.statusRegA = util.DwordGetUpperWord(dg.DwordT(addr)) // return address of 1st CB processed
		disk6239Data.statusRegB = util.DwordGetLowerWord(dg.DwordT(addr))
		disk6239SetPioStatusRegC(0, statCcsPioCmdOk, disk6239StartList, util.TestWbit(disk6239Data.commandRegC, 15))

	default:
		log.Panicf("disk6239 command %d not yet implemented\n", pioCmd)
	}
	disk6239Data.disk6239DataMu.Unlock()
}

func disk6239ExtractPioCommand(word dg.WordT) uint {
	res := uint((word & 01776) >> 1) // mask penultimate 9 bits
	return res
}

func disk6239GetCBextendedStatusSize() int {
	word := disk6239Data.intInfBlock[5]
	word >>= 8
	word &= 0x0f
	return int(word)
}

// Handle flag/pulse to disk6239
func disk6239HandleFlag(f byte) {
	switch f {
	case 'S':
		BusSetBusy(disk6239Data.devNum, true)
		BusSetDone(disk6239Data.devNum, false)
		if debugLogging {
			logging.DebugPrint(disk6239Data.logID, "... S flag set\n")
		}
		disk6239DoPioCommand()

		BusSetBusy(disk6239Data.devNum, false)
		// set the DONE flag if the return bit was set
		disk6239Data.disk6239DataMu.RLock()
		if util.TestWbit(disk6239Data.commandRegC, 15) {
			BusSetDone(disk6239Data.devNum, true)
		}
		disk6239Data.disk6239DataMu.RUnlock()

	case 'C':
		if debugLogging {
			logging.DebugPrint(disk6239Data.logID, "... C flag set, clearing DONE flag\n")
		}
		BusSetDone(disk6239Data.devNum, false)
		// TODO clear pending interrupt
		//disk6239Data.statusRegC = 0
		disk6239SetPioStatusRegC(statXecStateMapped,
			statCcsPioCmdOk,
			dg.WordT(disk6239ExtractPioCommand(disk6239Data.commandRegC)),
			util.TestWbit(disk6239Data.commandRegC, 15))

	case 'P':
		if debugLogging {
			logging.DebugPrint(disk6239Data.logID, "... P flag set\n")
		}
		log.Fatalln("P flag not yet implemented in disk6239")

	default:
		// no/empty flag - nothing to do
	}
}

// seek to the disk position according to sector number in disk6239Data structure
func disk6239PositionDiskImage() {
	var offset = int64(disk6239Data.sectorNo) * disk6239BytesPerSector
	_, err := disk6239Data.imageFile.Seek(offset, 0)
	if err != nil {
		log.Fatalln("disk6239 could not position disk image")
	}
	// TODO Set C/H/S???
}

// CB processing in a goroutine
func disk6239CBprocessor(dataPtr *disk6239DataT) {
	var (
		cb            [disk6239CbMaxSize]dg.WordT
		w, cbLength   int
		nextCB        dg.PhysAddrT
		sect          dg.DwordT
		physTransfers bool
		physAddr      dg.PhysAddrT
		readBuff      = make([]byte, disk6239BytesPerSector)
		writeBuff     = make([]byte, disk6239BytesPerSector)
		tmpWd         dg.WordT
	)
	for {
		cbAddr := <-cbChan
		cbLength = disk6239CbMinSize + disk6239GetCBextendedStatusSize()
		if debugLogging {
			logging.DebugPrint(disk6239Data.logID, "... Processing CB, extended status size is: %d\n", disk6239GetCBextendedStatusSize())
		}
		// copy CB contents from host memory
		addr := cbAddr
		for w = 0; w < cbLength; w++ {
			cb[w] = memory.ReadWordBmcChan(&addr)
			if debugLogging {
				logging.DebugPrint(disk6239Data.logID, "... CB[%d]: %d\n", w, cb[w])
			}
		}

		opCode := cb[disk6239CbINA_FLAGS_OPCODE] & 0x03ff
		nextCB = dg.PhysAddrT(util.DwordFromTwoWords(cb[disk6239CbLINK_ADDR_HIGH], cb[disk6239CbLINK_ADDR_LOW]))
		if debugLogging {
			logging.DebugPrint(disk6239Data.logID, "... CB OpCode: %d\n", opCode)
			logging.DebugPrint(disk6239Data.logID, "... .. Next CB Addr: %d\n", nextCB)
		}
		switch opCode {

		case disk6239CbOpRecalibrateDisk:
			dataPtr.disk6239DataMu.Lock()
			if debugLogging {
				logging.DebugPrint(disk6239Data.logID, "... .. RECALIBRATE\n")
			}
			//dataPtr.cylinder = 0
			//dataPtr.head = 0
			//dataPtr.sector = 0
			dataPtr.sectorNo = 0
			disk6239PositionDiskImage()
			dataPtr.disk6239DataMu.Unlock()
			if cbLength >= disk6239CbERR_STATUS+1 {
				cb[disk6239CbERR_STATUS] = 0
			}
			if cbLength >= disk6239CbUNIT_STATUS+1 {
				cb[disk6239CbUNIT_STATUS] = 1 << 13 // b0010000000000000; // Ready
			}
			if cbLength >= disk6239CbCB_STATUS+1 {
				cb[disk6239CbCB_STATUS] = 1 // finally, set Done bit
			}

		case disk6239CbOpRead:
			dataPtr.disk6239DataMu.Lock()
			dataPtr.sectorNo = util.DwordFromTwoWords(cb[disk6239CbDEV_ADDR_HIGH], cb[disk6239CbDEV_ADDR_LOW])
			if util.TestWbit(cb[disk6239CbPAGENO_LIST_ADDR_HIGH], 0) {
				// logical premapped host address
				physTransfers = false
				log.Fatal("disk6239 - CB READ from premapped logical addresses  Not Yet Implemented")
			} else {
				physTransfers = true
				physAddr = dg.PhysAddrT(util.DwordFromTwoWords(cb[disk6239CbTXFER_ADDR_HIGH], cb[disk6239CbTXFER_ADDR_LOW]))
			}
			if debugLogging {
				logging.DebugPrint(disk6239Data.logID, "... .. CB READ command, SECCNT: %d\n", cb[disk6239CbTXFER_COUNT])
				logging.DebugPrint(disk6239Data.logID, "... .. .. .... from sector:     %d\n", dataPtr.sectorNo)
				logging.DebugPrint(disk6239Data.logID, "... .. .. .... from phys addr:  %d\n", physAddr)
				logging.DebugPrint(disk6239Data.logID, "... .. .. .... physical txfer?: %d\n", util.BoolToInt(physTransfers))
			}
			for sect = 0; sect < dg.DwordT(cb[disk6239CbTXFER_COUNT]); sect++ {
				dataPtr.sectorNo += sect
				disk6239PositionDiskImage()
				dataPtr.imageFile.Read(readBuff)
				addr = physAddr + (dg.PhysAddrT(sect) * disk6239WordsPerSector)
				for w = 0; w < disk6239WordsPerSector; w++ {
					tmpWd = (dg.WordT(readBuff[w*2]) << 8) | dg.WordT(readBuff[(w*2)+1])
					memory.WriteWordBmcChan(&addr, tmpWd)
				}
				dataPtr.reads++
			}
			if cbLength >= disk6239CbERR_STATUS+1 {
				cb[disk6239CbERR_STATUS] = 0
			}
			if cbLength >= disk6239CbUNIT_STATUS+1 {
				cb[disk6239CbUNIT_STATUS] = 1 << 13 // b0010000000000000; // Ready
			}
			if cbLength >= disk6239CbCB_STATUS+1 {
				cb[disk6239CbCB_STATUS] = 1 // finally, set Done bit
			}

			if debugLogging {
				logging.DebugPrint(disk6239Data.logID, "... .. .... READ command finished\n")
				logging.DebugPrint(disk6239Data.logID, "Last buffer: %X\n", readBuff)
			}
			dataPtr.disk6239DataMu.Unlock()

		case disk6239CbOpWrite:
			dataPtr.disk6239DataMu.Lock()
			dataPtr.sectorNo = util.DwordFromTwoWords(cb[disk6239CbDEV_ADDR_HIGH], cb[disk6239CbDEV_ADDR_LOW])
			if util.TestWbit(cb[disk6239CbPAGENO_LIST_ADDR_HIGH], 0) {
				// logical premapped host address
				physTransfers = false
				log.Fatal("disk6239 - CB WRITE from premapped logical addresses  Not Yet Implemented")
			} else {
				physTransfers = true
				physAddr = dg.PhysAddrT(util.DwordFromTwoWords(cb[disk6239CbTXFER_ADDR_HIGH], cb[disk6239CbTXFER_ADDR_LOW]))
			}
			if debugLogging {
				logging.DebugPrint(disk6239Data.logID, "... .. CB WRITE command, SECCNT: %d\n", cb[disk6239CbTXFER_COUNT])
				logging.DebugPrint(disk6239Data.logID, "... .. .. ..... to sector:       %d\n", dataPtr.sectorNo)
				logging.DebugPrint(disk6239Data.logID, "... .. .. ..... from phys addr:  %d\n", physAddr)
				logging.DebugPrint(disk6239Data.logID, "... .. .. ..... physical txfer?: %d\n", util.BoolToInt(physTransfers))
			}
			for sect = 0; sect < dg.DwordT(cb[disk6239CbTXFER_COUNT]); sect++ {
				dataPtr.sectorNo += sect
				disk6239PositionDiskImage()
				memAddr := physAddr + (dg.PhysAddrT(sect) * disk6239WordsPerSector)
				for w = 0; w < disk6239WordsPerSector; w++ {
					tmpWd = memory.ReadWordBmcChan(&memAddr)
					writeBuff[w*2] = byte(tmpWd >> 8)
					writeBuff[(w*2)+1] = byte(tmpWd & 0x00ff)
				}
				dataPtr.imageFile.Write(writeBuff)
				if debugLogging {
					logging.DebugPrint(disk6239Data.logID, "Wrote buffer: %X\n", writeBuff)
				}
				dataPtr.writes++
			}
			if cbLength >= disk6239CbERR_STATUS+1 {
				cb[disk6239CbERR_STATUS] = 0
			}
			if cbLength >= disk6239CbUNIT_STATUS+1 {
				cb[disk6239CbUNIT_STATUS] = 1 << 13 // b0010000000000000; // Ready
			}
			if cbLength >= disk6239CbCB_STATUS+1 {
				cb[disk6239CbCB_STATUS] = 1 // finally, set Done bit
			}
			dataPtr.disk6239DataMu.Unlock()

		default:
			log.Fatalf("disk6239 CB Command %d not yet implemented\n", opCode)
		}

		// write back CB
		addr = cbAddr
		for w = 0; w < cbLength; w++ {
			memory.WriteWordBmcChan(&addr, cb[w])
		}

		if nextCB == 0 {
			// send ASYNCH status. See p.4-15
			if debugLogging {
				logging.DebugPrint(disk6239Data.logID, "...ready to set ASYNC status\n")
			}
			for BusGetBusy(disk6239Data.devNum) || BusGetDone(disk6239Data.devNum) {
				time.Sleep(disk6239AsynchStatRetryInterval)
			}
			dataPtr.disk6239DataMu.Lock()
			dataPtr.statusRegC = dg.WordT(statXecStateMapped) << 12
			dataPtr.statusRegC |= (statAsyncNoErrors & 0x03ff)
			if debugLogging {
				logging.DebugPrint(disk6239Data.logID, "disk6239 ASYNCHRONOUS status C set to: %s\n",
					util.WordToBinStr(dataPtr.statusRegC))
			}
			dataPtr.disk6239DataMu.Unlock()
			if debugLogging {
				logging.DebugPrint(disk6239Data.logID, "...set ASYNC status\n")
			}
			BusSetDone(disk6239Data.devNum, true)
		} else {
			// chain to next CB
			//disk6239ProcessCB(nextCB)
			cbChan <- nextCB
		}
	}
}

func disk6239Reset() {
	disk6239Data.disk6239DataMu.Lock()
	disk6239ResetMapping()
	disk6239ResetIntInfBlk()
	disk6239ResetCtrlrInfBlock()
	disk6239ResetUnitInfBlock()
	disk6239Data.statusRegB = 0
	disk6239SetPioStatusRegC(statXecStateResetDone, 0, disk6239PioReset, util.TestWbit(disk6239Data.commandRegC, 15))
	disk6239Data.disk6239DataMu.Unlock()
	if debugLogging {
		logging.DebugPrint(disk6239Data.logID, "disk6239 ***Reset*** via call to disk6239Reset()\n")
	}

}

// N.B. We assume disk6239Data is LOCKED before calling ANY of the following functions

// setup the controller information block to power-up defaults p.2-15
func disk6239ResetCtrlrInfBlock() {
	disk6239Data.ctrlInfBlock[0] = 0
	disk6239Data.ctrlInfBlock[1] = 0
}

// setup the interface information block to power-up defaults
func disk6239ResetIntInfBlk() {
	disk6239Data.intInfBlock[0] = 0101
	disk6239Data.intInfBlock[1] = disk6239UcodeRev
	disk6239Data.intInfBlock[2] = 3
	disk6239Data.intInfBlock[3] = 8<<11 | disk6239MaxQueuedCBs
	disk6239Data.intInfBlock[4] = 0
	disk6239Data.intInfBlock[5] = 11 << 8
	disk6239Data.intInfBlock[6] = 0
	disk6239Data.intInfBlock[7] = 0
}

// set mapping options after IORST, power-up or Reset
func disk6239ResetMapping() {
	disk6239Data.mappingRegA = 0x4000 // DMA over the BMC
	disk6239Data.mappingRegB = disk6239MapIntBmcPhys | disk6239MapUpstreamLoad | disk6239MapUpstreamHpt
	disk6239Data.isMapped = false
}

// setup the unit information block to power-up defaults pp.2-16
func disk6239ResetUnitInfBlock() {
	disk6239Data.unitInfBlock[0] = 0
	disk6239Data.unitInfBlock[1] = 9<<12 | disk6239UcodeRev
	disk6239Data.unitInfBlock[2] = dg.WordT(disk6239LogicalBlocksH) // 17.
	disk6239Data.unitInfBlock[3] = dg.WordT(disk6239LogicalBlocksL) // 43840.
	disk6239Data.unitInfBlock[4] = disk6239BytesPerSector
	disk6239Data.unitInfBlock[5] = disk6239UserCylinders
	disk6239Data.unitInfBlock[6] = ((disk6239SurfacesPerDisk * disk6239HeadsPerSurface) << 8) | (0x00ff & disk6239SectorsPerTrack)
}

// this is used to set the SYNCHRONOUS standard return as per p.3-22
func disk6239SetPioStatusRegC(stat byte, ccs byte, cmdEcho dg.WordT, rr bool) {
	if stat == 0 && disk6239Data.isMapped {
		stat = statXecStateMapped
	}
	if rr || cmdEcho == disk6239PioReset {
		disk6239Data.statusRegC = dg.WordT(stat) << 12
		disk6239Data.statusRegC |= (dg.WordT(ccs) & 3) << 10
		disk6239Data.statusRegC |= (cmdEcho & 0x01ff) << 1
		if rr {
			disk6239Data.statusRegC |= 1
		}
		if debugLogging {
			logging.DebugPrint(disk6239Data.logID, "disk6239 PIO (SYNCH) status C set to: %s\n",
				util.WordToBinStr(disk6239Data.statusRegC))
		}
	}
}
