// eclipsePC.go

// Copyright ©2017-2020  Steve Merrony

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

func eclipsePC(cpu *CPUT, iPtr *decodedInstrT) bool {

	ring := cpu.pc & 0x7000_0000

	switch iPtr.ix {

	case instrCLM: // signed compare to limits
		var (
			l, h int16
			inc  dg.PhysAddrT
		)
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		acs := int16(memory.DwordGetLowerWord(cpu.ac[twoAcc1Word.acs]))
		if twoAcc1Word.acs == twoAcc1Word.acd {
			l = int16(memory.ReadWord(cpu.pc + 1))
			h = int16(memory.ReadWord(cpu.pc + 2))
			if acs < l || acs > h {
				inc = 3
			} else {
				inc = 4
			}
		} else {
			l = int16(memory.ReadWord(dg.PhysAddrT(memory.DwordGetLowerWord(cpu.ac[twoAcc1Word.acd])) | ring))
			h = int16(memory.ReadWord(dg.PhysAddrT(memory.DwordGetLowerWord(cpu.ac[twoAcc1Word.acd])+1) | ring))
			if acs < l || acs > h {
				inc = 1
			} else {
				inc = 2
			}
		}
		if cpu.debugLogging {
			logging.DebugPrint(logging.DebugLog, "CLM compared %d with limits %d and %d, moving PC by %d\n", acs, l, h, inc)
		}
		cpu.pc += inc
		cpu.pc = (cpu.pc & 0x7fff) | ring

	case instrDSPA:
		oneAccModeInt2Word := iPtr.variant.(oneAccModeInd2WordT)
		tableStart := resolve15bitDisplacement(cpu, oneAccModeInt2Word.ind, oneAccModeInt2Word.mode, dg.WordT(oneAccModeInt2Word.disp15), iPtr.dispOffset)
		tableStart &= 0x7fff
		tableStart |= ring
		offset := memory.DwordGetLowerWord(cpu.ac[oneAccModeInt2Word.acd])
		lowLimit := memory.ReadWord(tableStart - 2)
		hiLimit := memory.ReadWord(tableStart - 1)
		if cpu.debugLogging {
			logging.DebugPrint(logging.DebugLog, "DSPA called with table at %d, offset %d, lo %d hi %d\n",
				tableStart, offset, lowLimit, hiLimit)
		}
		if offset < lowLimit || offset > hiLimit {
			log.Fatalf("ERROR: DPSA called with out of bounds offset %d", offset)
		}
		entry := tableStart - dg.PhysAddrT(lowLimit) + dg.PhysAddrT(offset)
		addr := dg.PhysAddrT(memory.ReadWord(entry))
		if addr == 0xffffffff {
			cpu.pc += 2
		} else {
			cpu.pc = (addr & 0x7fff) | ring
		}

	case instrEISZ:
		addr := resolve15bitDisplacement(cpu, iPtr.ind, iPtr.mode, iPtr.disp15, iPtr.dispOffset)
		addr &= 0x7fff
		addr |= ring
		wd := memory.ReadWord(addr)
		wd++
		memory.WriteWord(addr, wd)
		if wd == 0 {
			cpu.pc += 3
		} else {
			cpu.pc += 2
		}

	case instrEJMP:
		addr := resolve15bitDisplacement(cpu, iPtr.ind, iPtr.mode, iPtr.disp15, iPtr.dispOffset)
		addr &= 0x7fff
		addr |= ring
		cpu.pc = addr

	case instrEJSR:
		cpu.ac[3] = dg.DwordT(cpu.pc) + 2
		addr := resolve15bitDisplacement(cpu, iPtr.ind, iPtr.mode, iPtr.disp15, iPtr.dispOffset)
		addr &= 0x7fff
		addr |= ring
		cpu.pc = addr

	case instrFNS:
		cpu.pc++
		cpu.pc = (cpu.pc & 0x7fff) | ring

	case instrSGE: //16-bit signed numbers
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		acs := int16(memory.DwordGetLowerWord(cpu.ac[twoAcc1Word.acs]))
		acd := int16(memory.DwordGetLowerWord(cpu.ac[twoAcc1Word.acd]))
		if acs >= acd {
			cpu.pc += 2
		} else {
			cpu.pc++
		}
		cpu.pc = (cpu.pc & 0x7fff) | ring

	case instrSGT: //16-bit signed numbers
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		acs := int16(memory.DwordGetLowerWord(cpu.ac[twoAcc1Word.acs]))
		acd := int16(memory.DwordGetLowerWord(cpu.ac[twoAcc1Word.acd]))
		if acs > acd {
			cpu.pc += 2
		} else {
			cpu.pc++
		}
		cpu.pc = (cpu.pc & 0x7fff) | ring

	case instrSNB:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		addr, bit := resolveEclipseBitAddr(cpu, &twoAcc1Word)
		addr |= ring
		wd := memory.ReadWord(addr)
		if memory.TestWbit(wd, int(bit)) {
			cpu.pc += 2
		} else {
			cpu.pc++
		}
		cpu.pc = (cpu.pc & 0x7fff) | ring
		if cpu.debugLogging {
			logging.DebugPrint(logging.DebugLog, "SNB: Wd Addr: %d., word: %0X, bit #: %d\n", addr, wd, bit)
		}

	case instrSZB, instrSZBO:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		addr, bit := resolveEclipseBitAddr(cpu, &twoAcc1Word)
		addr |= (cpu.pc & ringMask32)
		wd := memory.ReadWord(addr)
		if !memory.TestWbit(wd, int(bit)) {
			if iPtr.ix == instrSZBO {
				memory.SetWbit(&wd, bit)
				memory.WriteWord(addr, wd)
			}
			cpu.pc += 2
		} else {
			cpu.pc++
		}
		cpu.pc = (cpu.pc & 0x7fff) | ring
		if cpu.debugLogging {
			logging.DebugPrint(logging.DebugLog, "SZB(O): Wd Addr: %d., word: %0X, bit #: %d\n", addr, wd, bit)
		}

	default:
		log.Fatalf("ERROR: ECLIPSE_PC instruction <%s> not yet implemented\n", iPtr.mnemonic)
		return false
	}

	return true
}
