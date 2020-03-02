// eagleIO.go

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

func eagleIO(cpu *CPUT, iPtr *decodedInstrT) bool {

	switch iPtr.ix {

	case instrCIO:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		word := memory.DwordGetLowerWord(cpu.ac[twoAcc1Word.acs])
		mapRegAddr := int(word & 0x0fff)
		ioChan := memory.GetWbits(word, 1, 3)
		if debugLogging {
			logging.DebugPrint(logging.DebugLog, "... Channel: %d.\n", ioChan)
		}
		// N.B. Channel 7 => all channels
		if ioChan != 0 && ioChan != 7 {
			log.Fatalf("ERROR: Attempt to use CIO on unsupported IO Channel %d.", ioChan)
		}
		if memory.TestWbit(word, 0) { // write command
			memory.BmcdchWriteReg(mapRegAddr, memory.DwordGetLowerWord(cpu.ac[twoAcc1Word.acd]))
			if debugLogging {
				logging.DebugPrint(logging.MapLog, "CIO write to register %#o with %#o\n", mapRegAddr, memory.DwordGetLowerWord(cpu.ac[twoAcc1Word.acd]))
				logging.DebugPrint(logging.MapLog, "... Written %#o to register %#o\n", memory.DwordGetLowerWord(cpu.ac[twoAcc1Word.acd]), mapRegAddr)
			}
		} else { // read command
			cpu.ac[twoAcc1Word.acd] = dg.DwordT(memory.BmcdchReadReg(mapRegAddr))
			if debugLogging {
				logging.DebugPrint(logging.DebugLog, "... Read %#o from register %#o\n", memory.BmcdchReadReg(mapRegAddr), mapRegAddr)
			}
		}

	case instrCIOI:
		// TODO handle I/O channel
		twoAccImm2Word := iPtr.variant.(twoAccImm2WordT)
		var cmd dg.WordT
		if twoAccImm2Word.acs == twoAccImm2Word.acd {
			cmd = twoAccImm2Word.immWord
		} else {
			cmd = twoAccImm2Word.immWord | memory.DwordGetLowerWord(cpu.ac[twoAccImm2Word.acs])
		}
		mapRegAddr := int(cmd & 0x0fff)
		if memory.TestWbit(cmd, 0) { // write command
			memory.BmcdchWriteReg(mapRegAddr, memory.DwordGetLowerWord(cpu.ac[twoAccImm2Word.acd]))
			if debugLogging {
				logging.DebugPrint(logging.MapLog, "CIOI write to register %#o with %#o\n", mapRegAddr, memory.DwordGetLowerWord(cpu.ac[twoAccImm2Word.acd]))
				logging.DebugPrint(logging.MapLog, "... Written %#o to register %#o\n", memory.DwordGetLowerWord(cpu.ac[twoAccImm2Word.acd]), mapRegAddr)
			}
		} else { // read command
			cpu.ac[twoAccImm2Word.acd] = dg.DwordT(memory.BmcdchReadReg(mapRegAddr))
		}

	case instrECLID: // seems to be the same as LCPID
		dwd := dg.DwordT((cpuModelNo & 0xffff)) << 16
		dwd |= (ucodeRev & 0x0f) << 8
		dwd |= MemSizeLCPID & 0x0f
		cpu.ac[0] = dwd

	case instrINTDS:
		return intds(cpu)

	case instrINTEN:
		return inten(cpu)

	case instrLCPID: // seems to be the same as ECLID
		dwd := dg.DwordT((cpuModelNo & 0xffff)) << 16
		dwd |= (ucodeRev & 0x0f) << 8
		dwd |= MemSizeLCPID & 0x0f
		cpu.ac[0] = dwd

		// MSKO is handled via DOB n,CPU

	case instrNCLID:
		cpu.ac[0] = cpuModelNo & 0xffff
		cpu.ac[1] = ucodeRev & 0xffff
		cpu.ac[2] = MemSizeNCLID & 0x00ff

	case instrPRTSEL:
		if debugLogging {
			logging.DebugPrint(logging.DebugLog, "INFO: PRTSEL AC0: %d, PC: %d\n", cpu.ac[0], cpu.pc)
		}
		// only handle the query mode, setting is a no-op on this 'single-channel' machine
		if memory.DwordGetLowerWord(cpu.ac[0]) == 0xffff {
			// return default I/O channel if -1 passed in
			cpu.ac[0] = 0
		}

	case instrREADS:
		oneAcc1Word := iPtr.variant.(oneAcc1WordT)
		return reads(cpu, oneAcc1Word.acd)

	case instrWLMP:
		if cpu.ac[1] == 0 {
			mapRegAddr := int(cpu.ac[0] & 0x7ff)
			wAddr := dg.PhysAddrT(cpu.ac[2])
			if debugLogging {
				logging.DebugPrint(logging.DebugLog, "WLMP called with AC1 = 0 - MapRegAddr was %#o, 1st DWord was %#o\n",
					mapRegAddr, memory.ReadDWord(wAddr))
				logging.DebugPrint(logging.MapLog, "WLMP called with AC1 = 0 - MapRegAddr was %#o, 1st DWord was %#o\n",
					mapRegAddr, memory.ReadDWord(wAddr))
			}
			// memory.BmcdchWriteSlot(mapRegAddr, memory.ReadDWord(wAddr))
			// cpu.ac[0]++
			// cpu.ac[2] += 2
		} else {
			for {
				dwd, ok := memory.ReadDwordTrap(dg.PhysAddrT(cpu.ac[2]))
				if !ok {
					log.Fatalf("ERROR: Memory access failed at PC: %#o\n", cpu.pc)
				}
				memory.BmcdchWriteSlot(int(cpu.ac[0]&0x07ff), dwd)
				if debugLogging {
					logging.DebugPrint(logging.DebugLog, "WLMP written slot: %#o, data: %#o\n", cpu.ac[0]&0x7ff, dwd)
					logging.DebugPrint(logging.MapLog, "WLMP written slot: %#o, data: %#o\n", cpu.ac[0]&0x7ff, dwd)
				}
				cpu.ac[2] += 2
				cpu.ac[0]++
				cpu.ac[1]--
				if cpu.ac[1] == 0 {
					break
				}
			}
		}

	default:
		log.Fatalf("ERROR: EAGLE_IO instruction <%s> not yet implemented\n", iPtr.mnemonic)
		return false
	}

	cpu.pc += dg.PhysAddrT(iPtr.instrLength)
	return true
}
