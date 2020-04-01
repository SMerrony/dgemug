// cpu.go

// Copyright ©2017-2020 Steve Merrony

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

package mvcpu

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/SMerrony/dgemug/devices"
	"github.com/SMerrony/dgemug/dg"
	"github.com/SMerrony/dgemug/logging"
	"github.com/SMerrony/dgemug/memory"
)

const (
	cpuModelNo = 0x224C // => MV/10000 according to p.2-19 of AOS/VS Internals
	ucodeRev   = 0x04

	// MemSizeWords defines the size of MV/Em's emulated RAM in 16-bit words
	MemSizeWords = 8388608 // = 040 000 000 (8) = 0x80 0000
	// MemSizeLCPID is the code returned by the LCPID to indicate the size of RAM in half megabytes
	MemSizeLCPID = ((MemSizeWords * 2) / (256 * 1024)) - 1 // 0x3F
	// MemSizeNCLID is the code returned by NCLID to indicate size of RAM in 32Kb increments
	MemSizeNCLID = ((MemSizeWords * 2) / (32 * 1024)) - 1
)

// Useful signed int limits
const (
	maxPosS16 = 1<<15 - 1
	minNegS16 = -(maxPosS16 + 1)
	maxPosS32 = 1<<31 - 1
	minNegS32 = -(maxPosS32 + 1)
)

// TODO sbrBits is currently an abstraction of the Segment Base Registers - may need to represent physically
// via a 32-bit DWord in the future
type sbrBits struct {
	v, len, lef, io bool
	physAddr        uint32 // 19 bits used
}

// CPUT holds the current state of a CPUT
type CPUT struct {
	cpuMu sync.RWMutex
	// representations of physical attributes
	pc dg.PhysAddrT // 32-bit PC
	ac [4]dg.DwordT // 4 x 32-bit Accumulators
	// mask                    dg.WordT     // interrupt mask - moved to bus
	psr                     dg.WordT     // Processor Status Register - see PoP 2-11 & A-4
	carry, atu, ion, pfflag bool         // flag bits
	sbr                     [8]sbrBits   // SBRs (see above)
	fpac                    [4]float64   // 4 x 64-bit Floating Point Acs N.B Not same internal fmt as DG
	fpsr                    dg.QwordT    // 64-bit Floating-Point Status Register
	sr                      dg.WordT     // Not sure about this... fake Switch Register
	wfp, wsp, wsl, wsb      dg.PhysAddrT // Active Wide Stack values

	devNum int
	bus    *devices.BusT

	// emulator internals
	debugLogging bool
	instrCount   uint64 // how many instructions executed during the current run, running at 2 MIPS this will loop round roughly every 100 million years!
	scpIO        bool   // true if console I/O is directed to the SCP
}

// CPUStatT defines the data we will send to the statusCollector monitor
type CPUStatT struct {
	Pc              dg.PhysAddrT
	Ac              [4]dg.DwordT
	Carry, Atu, Ion bool
	InstrCount      uint64
	GoVersion       string
	GoroutineCount  int
	HostCPUCount    int
	HeapSizeMB      int
}

const cpuStatPeriodMs = 333 // 125 // i.e. we send stats every 1/8th of a second

// CPUInit sets up an MV-Class CPU
func (cpu *CPUT) CPUInit(devNum int, bus *devices.BusT, statsChan chan CPUStatT) {
	cpu.devNum = devNum
	cpu.bus = bus
	cpu.Reset()
	decoderGenAllPossOpcodes()
	if statsChan != nil {
		go cpu.statSender(statsChan)
	}
}

// Reset sets sane initial values for a CPU
func (cpu *CPUT) Reset() {
	cpu.cpuMu.Lock()
	cpu.pc = 0
	for a := 0; a < 4; a++ {
		cpu.ac[a] = 0
		cpu.fpac[a] = 0
	}
	cpu.psr = 0
	cpu.carry = false
	cpu.atu = false
	cpu.ion = false
	cpu.pfflag = false
	cpu.SetOVR(false)
	cpu.instrCount = 0
	cpu.cpuMu.Unlock()
}

// PrepToRun is called prior to a normal run
func (cpu *CPUT) PrepToRun() {
	cpu.cpuMu.Lock()
	cpu.instrCount = 0
	cpu.scpIO = false
	cpu.cpuMu.Unlock()
}

// Boot sets up the CPU for booting
func (cpu *CPUT) Boot(devNum int, pc dg.PhysAddrT) {
	cpu.cpuMu.Lock()
	cpu.sr = 0x8000 | dg.WordT(devNum)
	cpu.ac[0] = dg.DwordT(devNum)
	cpu.pc = pc
	cpu.cpuMu.Unlock()
}

// PrintableStatus returns a verbose status of the CPU
func (cpu *CPUT) PrintableStatus() string {
	cpu.cpuMu.RLock()
	res := fmt.Sprintf("%c         AC0          AC1         AC2          AC3           PC CRY LEF ATU ION%c", dg.ASCIINL, dg.ASCIINL)
	res += fmt.Sprintf("%#12o %#12o %#12o %#12o %#12o", cpu.ac[0], cpu.ac[1], cpu.ac[2], cpu.ac[3], cpu.pc)
	res += fmt.Sprintf("  %d   %d   %d   %d",
		memory.BoolToInt(cpu.carry),
		memory.BoolToInt(cpu.sbr[memory.GetSegment(cpu.pc)].lef),
		memory.BoolToInt(cpu.atu),
		memory.BoolToInt(cpu.ion))
	cpu.cpuMu.RUnlock()
	return res
}

// DisassembleRange returns a DASHER-formatted string of the disassembled specified region
func (cpu *CPUT) DisassembleRange(lowAddr, highAddr dg.PhysAddrT) (disassembly string) {
	if lowAddr > highAddr {
		return ("%c *** Invalid address range for disassembly ***")
	}

	var skipDecode int

	for addr := lowAddr; addr <= highAddr; addr++ {
		word := memory.ReadWord(addr)
		byte1 := dg.ByteT(word >> 8)
		byte2 := dg.ByteT(word & 0x00ff)
		display := fmt.Sprintf("%c%#x: %02X %02X %06o %s \"", dg.ASCIINL, addr, byte1, byte2, word, memory.WordToBinStr(word))
		if byte1 >= ' ' && byte1 <= '~' {
			display += string(byte1)
		} else {
			display += " "
		}
		if byte2 >= ' ' && byte2 <= '~' {
			display += string(byte2)
		} else {
			display += " "
		}
		display += "\" "
		if skipDecode == 0 {
			instrTmp, ok := InstructionDecode(word, addr, true, false, true, true, nil)
			if ok {
				display += instrTmp.GetDisassembly()
				if instrTmp.GetLength() > 1 {
					skipDecode = instrTmp.GetLength() - 1
				}
			} else {
				display += " *** Could not decode ***"
			}
		} else {
			skipDecode--
		}
		disassembly = disassembly + display
	}
	return disassembly
}

// CompactPrintableStatus returns a concise CPU status
func (cpu *CPUT) CompactPrintableStatus() string {
	cpu.cpuMu.RLock()
	res := fmt.Sprintf("AC0=%-12o AC1=%-12o AC2=%-12o AC3=%-12o C:%d I:%d PC=%-12o",
		cpu.ac[0], cpu.ac[1], cpu.ac[2], cpu.ac[3],
		memory.BoolToInt(cpu.carry), memory.BoolToInt(cpu.ion), cpu.pc)
	cpu.cpuMu.RUnlock()
	return res
}

// GetAc is a getter for the ACs
func (cpu *CPUT) GetAc(ac int) (contents dg.DwordT) {
	cpu.cpuMu.RLock()
	contents = cpu.ac[ac]
	cpu.cpuMu.RUnlock()
	return contents
}

// SetAc is a setter for the ACs
func (cpu *CPUT) SetAc(ac int, val dg.DwordT) {
	cpu.cpuMu.Lock()
	cpu.ac[ac] = val
	cpu.cpuMu.Unlock()
}

// GetAtu returns the current ATU setting
func (cpu *CPUT) GetAtu() (atu bool) {
	cpu.cpuMu.RLock()
	atu = cpu.atu
	cpu.cpuMu.RUnlock()
	return atu
}

// SetATU is a setter for the ATU
func (cpu *CPUT) SetATU(atu bool) {
	cpu.cpuMu.Lock()
	cpu.atu = atu
	cpu.cpuMu.Unlock()
}

// GetDebugLogging is a getter for the debug logging flag
func (cpu *CPUT) GetDebugLogging() (logging bool) {
	cpu.cpuMu.RLock()
	logging = cpu.debugLogging
	cpu.cpuMu.RUnlock()
	return logging
}

// SetDebugLogging is a setter for debug logging
func (cpu *CPUT) SetDebugLogging(logging bool) {
	cpu.cpuMu.Lock()
	cpu.debugLogging = logging
	cpu.cpuMu.Unlock()
}

// GetLef returns the current LEF mode bit
func (cpu *CPUT) GetLef(segment int) (lef bool) {
	cpu.cpuMu.RLock()
	lef = cpu.sbr[segment].lef
	cpu.cpuMu.RUnlock()
	return lef
}

// SetLef sets the LEF mode bit for the current (PC) segment
func (cpu *CPUT) SetLef(lef bool) {
	cpu.cpuMu.Lock()
	cpu.sbr[(cpu.pc&0x7000_0000)>>28].lef = lef
	cpu.cpuMu.Unlock()
}

// GetIO returns the current IO bit for a segment
func (cpu *CPUT) GetIO(segment int) (io bool) {
	cpu.cpuMu.RLock()
	io = cpu.sbr[segment].io
	cpu.cpuMu.RUnlock()
	return io
}

// GetInstrCount returns the instruction-counting array
func (cpu *CPUT) GetInstrCount() (ic uint64) {
	cpu.cpuMu.RLock()
	ic = cpu.instrCount
	cpu.cpuMu.RUnlock()
	return ic
}

// GetOVR is a getter for the OVR flag embedded in the PSR
func (cpu *CPUT) GetOVR() bool {
	return memory.TestWbit(cpu.psr, 1)
}

// SetOVR is a setter for the OVR flag embedded in the PSR
func (cpu *CPUT) SetOVR(newOVR bool) {
	if newOVR {
		memory.SetWbit(&cpu.psr, 1)
	} else {
		memory.ClearWbit(&cpu.psr, 1)
	}
}

// GetOVK is a getter for the OVK mask embedded in the PSR
func (cpu *CPUT) GetOVK() bool {
	return memory.TestWbit(cpu.psr, 0)
}

// SetOVK is a setter for the OVK flag embedded in the PSR
func (cpu *CPUT) SetOVK(newOVK bool) {
	if newOVK {
		memory.SetWbit(&cpu.psr, 0)
	} else {
		memory.ClearWbit(&cpu.psr, 0)
	}
}

// GetPC is a getter for the PC
func (cpu *CPUT) GetPC() (pc dg.PhysAddrT) {
	cpu.cpuMu.RLock()
	pc = cpu.pc
	cpu.cpuMu.RUnlock()
	return pc
}

// SetPC sets the Program Counter
func (cpu *CPUT) SetPC(addr dg.PhysAddrT) {
	cpu.cpuMu.Lock()
	cpu.pc = addr
	cpu.cpuMu.Unlock()
}

// GetSCPIO is a getter for the SCP I/O flag
func (cpu *CPUT) GetSCPIO() (scp bool) {
	cpu.cpuMu.RLock()
	scp = cpu.scpIO
	cpu.cpuMu.RUnlock()
	return scp
}

// SetSCPIO is a setter for the SCP I/O flag
func (cpu *CPUT) SetSCPIO(scp bool) {
	cpu.cpuMu.Lock()
	cpu.scpIO = scp
	cpu.cpuMu.Unlock()
}

// GetWFP is a getter for the Wide Frame Pointer
func (cpu *CPUT) GetWFP() (wfp dg.PhysAddrT) {
	cpu.cpuMu.RLock()
	wfp = cpu.wfp
	cpu.cpuMu.RUnlock()
	return wfp
}

// GetWSP is a getter for the Wide Stack Pointer
func (cpu *CPUT) GetWSP() (wsp dg.PhysAddrT) {
	cpu.cpuMu.RLock()
	wsp = cpu.wsp
	cpu.cpuMu.RUnlock()
	return wsp
}

// SetupStack is a group-setter for the Wide Stack
func (cpu *CPUT) SetupStack(wfp, wsp, wsb, wsl, wsfh dg.PhysAddrT) {
	cpu.cpuMu.Lock()
	cpu.wfp = wfp
	cpu.wsp = wsp
	cpu.wsb = wsb
	cpu.wsl = wsl
	memory.WriteWord((cpu.pc&0x7000_0000)|wsfhLoc, dg.WordT(wsfh))
	cpu.cpuMu.Unlock()
}

// Execute runs a single instruction
// A false return means failure, the VM should stop
func (cpu *CPUT) Execute(iPtr *decodedInstrT) (rc bool) {
	cpu.cpuMu.Lock()
	switch iPtr.instrType {
	case NOVA_MEMREF:
		rc = novaMemRef(cpu, iPtr)
	case NOVA_OP:
		rc = novaOp(cpu, iPtr)
	case NOVA_IO:
		rc = novaIO(cpu, iPtr)
	case NOVA_MATH:
		rc = novaMath(cpu, iPtr)
	case NOVA_PC:
		rc = novaPC(cpu, iPtr)
	case ECLIPSE_MEMREF:
		rc = eclipseMemRef(cpu, iPtr)
	case ECLIPSE_OP:
		rc = eclipseOp(cpu, iPtr)
	case ECLIPSE_PC:
		rc = eclipsePC(cpu, iPtr)
	case ECLIPSE_STACK:
		rc = eclipseStack(cpu, iPtr)
	case EAGLE_FPU:
		rc = eagleFPU(cpu, iPtr)
	case EAGLE_IO:
		rc = eagleIO(cpu, iPtr)
	case EAGLE_OP:
		rc = eagleOp(cpu, iPtr)
	case EAGLE_MEMREF:
		rc = eagleMemRef(cpu, iPtr)
	case EAGLE_PC:
		rc = eaglePC(cpu, iPtr)
	case EAGLE_STACK:
		rc = eagleStack(cpu, iPtr)
	default:
		log.Println("ERROR: Unimplemented instruction type in Execute()")
		rc = false
	}
	cpu.instrCount++
	cpu.cpuMu.Unlock()
	return rc
}

// Run is the main activity loop for the virtual CPU
func (cpu *CPUT) Run(disassembly bool,
	deviceMap devices.DeviceMapT,
	breakpoints []dg.PhysAddrT,
	inputRadix int,
	tto *devices.TtoT) (errDetail string, instrCounts [maxInstrs]int) {

	var (
		thisOp dg.WordT
		prevPC dg.PhysAddrT
		iPtr   *decodedInstrT
		ok     bool
		indIrq byte
	)

	// initial read lock taken before loop starts to eliminate one lock/unlock per cycle
	cpu.cpuMu.RLock()

RunLoop: // performance-critical section starts here
	for {
		// FETCH
		thisOp = memory.ReadWord(cpu.pc)

		// DECODE
		iPtr, ok = InstructionDecode(thisOp, cpu.pc, cpu.sbr[cpu.pc>>29].lef, cpu.sbr[cpu.pc>>29].io, cpu.atu, disassembly, deviceMap)
		cpu.cpuMu.RUnlock()
		if !ok || iPtr.ix == -1 {
			errDetail = " *** Error: could not decode instruction ***"
			break
		}

		if cpu.debugLogging {
			logging.DebugPrint(logging.DebugLog, "%s  %s\n", cpu.CompactPrintableStatus(), iPtr.disassembly)
		}

		// EXECUTE
		if !cpu.Execute(iPtr) {
			errDetail = " *** Error: could not execute instruction (or CPU HALT encountered) ***"
			break
		}

		// INTERRUPT?
		cpu.cpuMu.Lock()
		if cpu.ion && cpu.bus.GetIRQ() {
			if cpu.debugLogging {
				logging.DebugPrint(logging.DebugLog, "<<< Interrupt >>>\n")
			}
			// disable further interrupts, reset the irq
			cpu.ion = false
			cpu.bus.SetIRQ(false)
			// TODO - disable User MAP
			// store PC in location zero
			memory.WriteWord(0, dg.WordT(cpu.pc))
			// fetch service routine address from location one
			if memory.TestWbit(memory.ReadWord(1), 0) {
				indIrq = '@'
			} else {
				indIrq = ' '
			}
			cpu.pc = resolve15bitDisplacement(cpu, indIrq, absoluteMode, memory.ReadWord(1), 0)
			// next time round RunLoop the interrupt service routine will be started...
		}
		cpu.cpuMu.Unlock()

		// BREAKPOINT?
		if len(breakpoints) > 0 {
			cpu.cpuMu.Lock()
			for _, bAddr := range breakpoints {
				if bAddr == cpu.pc {
					cpu.scpIO = true
					cpu.cpuMu.Unlock()
					msg := fmt.Sprintf(" *** BREAKpoint hit at physical address "+
						fmtRadixVerb(inputRadix)+
						" (previous PC "+fmtRadixVerb(inputRadix)+
						") ***",
						cpu.pc, prevPC)
					tto.PutNLString(msg)
					log.Println(msg)

					break RunLoop
				}
			}
			cpu.cpuMu.Unlock()
		}

		// Console interrupt?
		cpu.cpuMu.RLock()
		if cpu.scpIO {
			cpu.cpuMu.RUnlock()
			errDetail = " *** Console ESCape ***"
			break
		}

		// instruction counting
		instrCounts[iPtr.ix]++

		prevPC = cpu.pc

		// N.B. RLock still in effect as we loop around
	}

	return errDetail, instrCounts
}

// System call trap types
const (
	SyscallNot = iota
	Syscall32Trap
	Syscall16Trap
)

// Vrun is a simplified runloop for a Virtual CPU
// It should run until a system call is encountered
func (cpu *CPUT) Vrun() (syscallTrap int, errDetail string, instrCounts [maxInstrs]int) {
	var (
		thisOp dg.WordT
		// prevPC dg.PhysAddrT
		iPtr *decodedInstrT
		ok   bool
		// indIrq byte
	)

	// initial read lock taken before loop starts to eliminate one lock/unlock per cycle
	cpu.cpuMu.RLock()

	// RunLoop: // performance-critical section starts here
	for {
		// FETCH
		thisOp = memory.ReadWord(cpu.pc)

		// DECODE
		iPtr, ok = InstructionDecode(thisOp, cpu.pc, true, false, true, cpu.debugLogging, nil)
		cpu.cpuMu.RUnlock()
		if !ok || iPtr.ix == -1 {
			errDetail = " *** Error: could not decode instruction ***"
			break
		}

		if cpu.debugLogging {
			logging.DebugPrint(logging.DebugLog, "%s  %s\n", cpu.CompactPrintableStatus(), iPtr.disassembly)
			log.Printf("%s %s\n", cpu.CompactPrintableStatus(), iPtr.disassembly)
		}

		// Trap System Calls, first 32-bit style...
		if iPtr.ix == instrXJSR && iPtr.disp15 == 0x06 && iPtr.ind == '@' && iPtr.mode == absoluteMode {
			syscallTrap = Syscall32Trap
			break
		}
		// ...now 16-bit style
		if iPtr.ix == instrJSR && iPtr.disp15 == 017 && iPtr.ind == '@' && iPtr.mode == absoluteMode {
			cpu.cpuMu.Lock()
			cpu.pc-- // Fudge!
			cpu.cpuMu.Unlock()
			syscallTrap = Syscall16Trap
			break
		}

		// EXECUTE
		if !cpu.Execute(iPtr) {
			errDetail = " *** Error: could not execute instruction (or CPU HALT encountered) ***"
			break
		}

		// INTERRUPT?
		// cpu.cpuMu.Lock()
		// if cpu.ion && cpu.bus.GetIRQ() {
		// 	if cpu.debugLogging {
		// 		logging.DebugPrint(logging.DebugLog, "<<< Interrupt >>>\n")
		// 	}
		// 	// disable further interrupts, reset the irq
		// 	cpu.ion = false
		// 	cpu.bus.SetIRQ(false)
		// 	// TODO - disable User MAP
		// 	// store PC in location zero
		// 	memory.WriteWord(0, dg.WordT(cpu.pc))
		// 	// fetch service routine address from location one
		// 	if memory.TestWbit(memory.ReadWord(1), 0) {
		// 		indIrq = '@'
		// 	} else {
		// 		indIrq = ' '
		// 	}
		// 	cpu.pc = resolve15bitDisplacement(cpu, indIrq, absoluteMode, memory.ReadWord(1), 0)
		// 	// next time round RunLoop the interrupt service routine will be started...
		// }
		// cpu.cpuMu.Unlock()

		// BREAKPOINT?
		// if len(breakpoints) > 0 {
		// 	cpu.cpuMu.Lock()
		// 	for _, bAddr := range breakpoints {
		// 		if bAddr == cpu.pc {
		// 			cpu.scpIO = true
		// 			cpu.cpuMu.Unlock()
		// 			msg := fmt.Sprintf(" *** BREAKpoint hit at physical address "+
		// 				fmtRadixVerb(inputRadix)+
		// 				" (previous PC "+fmtRadixVerb(inputRadix)+
		// 				") ***",
		// 				cpu.pc, prevPC)
		// 			tto.PutNLString(msg)
		// 			log.Println(msg)

		// 			break RunLoop
		// 		}
		// 	}
		// 	cpu.cpuMu.Unlock()
		// }

		// Console interrupt?
		cpu.cpuMu.RLock()

		if cpu.pc == 0x7000_0000 {
			break
		}
		// if cpu.scpIO {
		// 	cpu.cpuMu.RUnlock()
		// 	errDetail = " *** Console ESCape ***"
		// 	break
		// }

		syscallTrap = SyscallNot

		// instruction counting
		instrCounts[iPtr.ix]++

		// prevPC = cpu.pc

		// N.B. RLock still in effect as we loop around
	}

	return syscallTrap, errDetail, instrCounts
}

func (cpu *CPUT) statSender(sChan chan CPUStatT) {
	var stats CPUStatT
	var memStats runtime.MemStats
	stats.GoVersion = runtime.Version()
	stats.HostCPUCount = runtime.NumCPU()
	for {
		cpu.cpuMu.RLock()
		stats.Pc = cpu.pc
		stats.Ac[0] = cpu.ac[0]
		stats.Ac[1] = cpu.ac[1]
		stats.Ac[2] = cpu.ac[2]
		stats.Ac[3] = cpu.ac[3]
		stats.Ion = cpu.ion
		stats.Atu = cpu.atu
		stats.Carry = cpu.carry
		stats.InstrCount = cpu.instrCount
		cpu.cpuMu.RUnlock()
		stats.GoroutineCount = runtime.NumGoroutine()
		runtime.ReadMemStats(&memStats)
		stats.HeapSizeMB = int(memStats.HeapAlloc / 1048576)
		select {
		case sChan <- stats:
		default:
		}
		time.Sleep(time.Millisecond * cpuStatPeriodMs)
	}
}

func fmtRadixVerb(inputRadix int) string {
	switch inputRadix {
	case 2:
		return "%b"
	case 8:
		return "%#o"
	case 10:
		return "%d."
	case 16:
		return "%#x"
	default:
		log.Fatalf("ERROR: Invalid input radix %d", inputRadix)
		return ""
	}
}
