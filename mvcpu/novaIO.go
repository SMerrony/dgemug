// novaIO.go

// Copyright (C) 2017,2019  Steve Merrony

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
	"log"

	"github.com/SMerrony/dgemug/dg"
	"github.com/SMerrony/dgemug/logging"
	"github.com/SMerrony/dgemug/memory"
)

func novaIO(cpu *CPUT, iPtr *decodedInstrT) bool {

	// The Eclipse LEF instruction is handled funkily...
	if cpu.atu && cpu.sbr[memory.GetSegment(cpu.pc)].lef {
		iPtr.ix = instrLEF
		log.Fatalf("ERROR: LEF not yet implemented, location %d\n", cpu.pc)
	}

	switch iPtr.ix {

	case instrDIA, instrDIB, instrDIC, instrDOA, instrDOB, instrDOC:
		novaDataIo := iPtr.variant.(novaDataIoT)

		// catch CPU I/O instructions
		if novaDataIo.ioDev == cpu.devNum {
			switch iPtr.ix {
			case instrDIA: // READS
				logging.DebugPrint(logging.DebugLog, "INFO: Interpreting DIA n,CPU as READS n instruction\n")
				return reads(cpu, novaDataIo.acd)
			case instrDIB: // INTA
				logging.DebugPrint(logging.DebugLog, "INFO: Interpreting DIB n,CPU as INTA n instruction\n")
				inta(cpu, novaDataIo.acd)
				switch novaDataIo.f {
				case 'S':
					cpu.ion = true
				case 'C':
					cpu.ion = false
				}
				return true
			case instrDIC: // IORST
				logging.DebugPrint(logging.DebugLog, "INFO: I/O Reset due to DIC 0,CPU instruction\n")
				return iorst(cpu)
			case instrDOB: // MKSO
				novaDataIo := iPtr.variant.(novaDataIoT)
				logging.DebugPrint(logging.DebugLog, "INFO: Handling DOB %d, CPU instruction as MSKO with flags\n", novaDataIo.acd)
				msko(cpu, novaDataIo.acd)
				switch novaDataIo.f {
				case 'S':
					cpu.ion = true
				case 'C':
					cpu.ion = false
				}
				return true
			case instrDOC: // HALT
				logging.DebugPrint(logging.DebugLog, "INFO: CPU Halting due to DOC %d,CPU (HALT) instruction\n", novaDataIo.acd)
				return halt()
			}
		}

		if cpu.bus.IsAttached(novaDataIo.ioDev) && cpu.bus.IsIODevice(novaDataIo.ioDev) {
			var abc byte
			switch iPtr.ix {
			case instrDOA, instrDIA:
				abc = 'A'
			case instrDOB, instrDIB:
				abc = 'B'
			case instrDOC, instrDIC:
				abc = 'C'
			}
			switch iPtr.ix {
			case instrDIA, instrDIB, instrDIC:
				cpu.ac[novaDataIo.acd] = dg.DwordT(cpu.bus.DataIn(novaDataIo.ioDev, abc, novaDataIo.f))
				//busDataIn(cpu, &novaDataIo, abc)
			case instrDOA, instrDOB, instrDOC:
				cpu.bus.DataOut(novaDataIo.ioDev, memory.DwordGetLowerWord(cpu.ac[novaDataIo.acd]), abc, novaDataIo.f)
				//busDataOut(cpu, &novaDataIo, abc)
			}
		} else {
			logging.DebugPrint(logging.DebugLog, "WARN: I/O attempted to unattached or non-I/O capable device 0%o\n", novaDataIo.ioDev)
			switch novaDataIo.ioDev {
			// case 0:
			// 	switch iPtr.ix {
			// 	case instrDIA:
			// 		cpu.ac[0] = 056
			// 	case instrDIB:
			// 		cpu.ac[0] = 0xffff
			// 	}
			case 2, 012, 013: // TODO - ignore for now
			default:
				return false
			}
		}

	case instrHALT:
		logging.DebugPrint(logging.DebugLog, "INFO: CPU Halting due to HALT instruction\n")
		return halt()

	case instrINTA:
		return inta(cpu, iPtr.ac)

	case instrINTDS:
		return intds(cpu)

	case instrINTEN:
		return inten(cpu)

	case instrIORST:
		// oneAcc1Word := iPtr.variant.(oneAcc1WordT) // <== this is just an assertion really
		cpu.bus.ResetAllIODevices()
		cpu.ion = false
		// TODO More to do for SMP support - HaHa!

	case instrNIO:
		ioFlagsDev := iPtr.variant.(ioFlagsDevT)

		if ioFlagsDev.ioDev == cpu.devNum {
			switch ioFlagsDev.f {
			case 'C': // INTDS
				return intds(cpu)
			case 'S': // INTEN
				return inten(cpu)
			}

		}
		if cpu.debugLogging {
			logging.DebugPrint(logging.DebugLog, "Sending NIO to device #%d.\n", ioFlagsDev.ioDev)
		}
		var novaDataIo novaDataIoT
		novaDataIo.f = ioFlagsDev.f
		novaDataIo.ioDev = ioFlagsDev.ioDev
		cpu.bus.DataOut(novaDataIo.ioDev, memory.DwordGetLowerWord(cpu.ac[novaDataIo.acd]), 'N', novaDataIo.f) // DUMMY FLAG

	case instrSKP:
		var busy, done bool
		ioTestDev := iPtr.variant.(ioTestDevT)
		switch ioTestDev.ioDev {
		case cpu.devNum:
			busy = cpu.ion
			done = cpu.pfflag
		case 012, 013: // TODO - ignore for now
			cpu.pc += 2
			return true
		default:
			busy = cpu.bus.GetBusy(ioTestDev.ioDev)
			done = cpu.bus.GetDone(ioTestDev.ioDev)
		}
		switch ioTestDev.t {
		case bnTest:
			if busy {
				cpu.pc++
			}
		case bzTest:
			if !busy {
				cpu.pc++
			}
		case dnTest:
			if done {
				cpu.pc++
			}
		case dzTest:
			if !done {
				cpu.pc++
			}
		}

	default:
		log.Fatalf("ERROR: NOVA_IO instruction <%s> not yet implemented\n", iPtr.mnemonic)
		return false
	}

	cpu.pc++
	return true
}

func halt() bool {
	// do not advance PC
	return false // stop processing
}

func intds(cpu *CPUT) bool {
	cpu.ion = false
	cpu.pc++
	return true
}

func inta(cpu *CPUT, destAc int) bool {
	// load the AC with the device code of the highest priority interrupt
	intDevNum := cpu.bus.GetHighestPriorityInt()
	cpu.ac[destAc] = dg.DwordT(intDevNum)
	// and clear it - I THINK this is the right place to do this...
	cpu.bus.ClearInterrupt(intDevNum)
	cpu.pc++
	return true
}

func inten(cpu *CPUT) bool {
	cpu.ion = true
	cpu.pc++
	return true
}

func iorst(cpu *CPUT) bool {
	cpu.bus.ResetAllIODevices()
	cpu.pc++
	return true
}

func msko(cpu *CPUT, destAc int) bool {
	//cpu.mask = memory.DwordGetLowerWord(cpu.ac[destAc])
	cpu.bus.SetIrqMask(memory.DwordGetLowerWord(cpu.ac[destAc]))
	cpu.pc++
	return true
}

func reads(cpu *CPUT, destAc int) bool {
	// load the AC with the contents of the dummy CPU register 'SR'
	cpu.ac[destAc] = dg.DwordT(cpu.sr)
	cpu.pc++
	return true
}
