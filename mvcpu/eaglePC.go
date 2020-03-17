// eaglePC.go

// Copyright (C) 2017,2019,2020 Steve Merrony

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

func eaglePC(cpu *CPUT, iPtr *decodedInstrT) bool {

	switch iPtr.ix {

	case instrDSZTS, instrISZTS:
		// tmpAddr := dg.PhysAddrT(memory.ReadDWord(cpu.wsp))
		tmpAddr := dg.PhysAddrT(cpu.wsp)
		var dwd dg.DwordT
		if iPtr.ix == instrDSZTS {
			dwd = memory.ReadDWord(tmpAddr) - 1
		} else {
			dwd = memory.ReadDWord(tmpAddr) + 1
		}
		memory.WriteDWord(tmpAddr, dwd)
		cpu.SetOVR(false)
		if dwd == 0 {
			cpu.pc += 2
		} else {
			cpu.pc++
		}
		if cpu.debugLogging {
			logging.DebugPrint(logging.DebugLog, "..... wrote %d to %d\n", dwd, tmpAddr)
		}

	case instrLCALL: // FIXME - LCALL only handling trivial case, no checking
		noAccModeInd4Word := iPtr.variant.(noAccModeInd4WordT)
		cpu.ac[3] = dg.DwordT(cpu.pc) + 4
		var dwd dg.DwordT
		if noAccModeInd4Word.argCount >= 0 {
			dwd = memory.DwordFromTwoWords(cpu.psr, dg.WordT(noAccModeInd4Word.argCount))
		} else {
			dwd = dg.DwordT(noAccModeInd4Word.argCount) & 0x00007fff
		}
		wsPush(cpu, 0, dwd)
		cpu.SetOVR(false)
		cpu.pc = resolve32bitEffAddr(cpu, noAccModeInd4Word.ind, noAccModeInd4Word.mode, noAccModeInd4Word.disp31, iPtr.dispOffset)

	case instrLDSP:
		oneAccModeInd3Word := iPtr.variant.(oneAccModeInd3WordT)
		value := int32(cpu.ac[oneAccModeInd3Word.acd])
		tableAddr := resolve32bitEffAddr(cpu, oneAccModeInd3Word.ind, oneAccModeInd3Word.mode, oneAccModeInd3Word.disp31, iPtr.dispOffset)
		h := int32(memory.ReadDWord(tableAddr - 2))
		l := int32(memory.ReadDWord(tableAddr - 4))
		if value < l || value > h {
			cpu.pc += 3
		} else {
			tableIndex := tableAddr + (2 * dg.PhysAddrT(value)) - (2 * dg.PhysAddrT(l))
			tableVal := memory.ReadDWord(tableIndex)
			if tableVal == 0xFFFFFFFF {
				cpu.pc += 3
			} else {
				cpu.pc = dg.PhysAddrT(tableVal) + tableIndex
			}
		}

	case instrLJMP:
		noAccModeInd3Word := iPtr.variant.(noAccModeInd3WordT)
		cpu.pc = cpu.pc&ringMask32 | resolve32bitEffAddr(cpu, noAccModeInd3Word.ind, noAccModeInd3Word.mode, noAccModeInd3Word.disp31, iPtr.dispOffset)

	case instrLJSR:
		noAccModeInd3Word := iPtr.variant.(noAccModeInd3WordT)
		cpu.ac[3] = dg.DwordT(cpu.pc) + 3
		cpu.pc = cpu.pc&ringMask32 | resolve32bitEffAddr(cpu, noAccModeInd3Word.ind, noAccModeInd3Word.mode, noAccModeInd3Word.disp31, iPtr.dispOffset)

	case instrLNISZ:
		noAccModeInd3Word := iPtr.variant.(noAccModeInd3WordT)
		// unsigned narrow increment and skip if zero
		tmpAddr := resolve32bitEffAddr(cpu, noAccModeInd3Word.ind, noAccModeInd3Word.mode, noAccModeInd3Word.disp31, iPtr.dispOffset)
		wd := memory.ReadWord(tmpAddr) + 1
		memory.WriteWord(tmpAddr, wd)
		if wd == 0 {
			cpu.pc += 4
		} else {
			cpu.pc += 3
		}

	case instrLPSHJ:
		noAccModeInd3Word := iPtr.variant.(noAccModeInd3WordT)
		wsPush(cpu, 0, dg.DwordT(cpu.pc)+3)
		cpu.pc = resolve32bitEffAddr(cpu, noAccModeInd3Word.ind, noAccModeInd3Word.mode, noAccModeInd3Word.disp31, iPtr.dispOffset)

	case instrLWDSZ:
		noAccModeInd3Word := iPtr.variant.(noAccModeInd3WordT)
		// unsigned wide decrement and skip if zero
		tmpAddr := resolve32bitEffAddr(cpu, noAccModeInd3Word.ind, noAccModeInd3Word.mode, noAccModeInd3Word.disp31, iPtr.dispOffset)
		tmp32b := memory.ReadDWord(tmpAddr) - 1
		memory.WriteDWord(tmpAddr, tmp32b)
		if tmp32b == 0 {
			cpu.pc += 4
		} else {
			cpu.pc += 3
		}

	case instrLWISZ:
		noAccModeInd3Word := iPtr.variant.(noAccModeInd3WordT)
		// unsigned wide increment and skip if zero
		tmpAddr := resolve32bitEffAddr(cpu, noAccModeInd3Word.ind, noAccModeInd3Word.mode, noAccModeInd3Word.disp31, iPtr.dispOffset)
		tmp32b := memory.ReadDWord(tmpAddr) + 1
		memory.WriteDWord(tmpAddr, tmp32b)
		if tmp32b == 0 {
			cpu.pc += 4
		} else {
			cpu.pc += 3
		}

	case instrNSALA:
		oneAccImm2Word := iPtr.variant.(oneAccImm2WordT)
		wd := ^memory.DwordGetLowerWord(cpu.ac[oneAccImm2Word.acd])
		if dg.WordT(oneAccImm2Word.immS16)&wd == 0 {
			cpu.pc += 3
		} else {
			cpu.pc += 2
		}

	case instrNSANA:
		oneAccImm2Word := iPtr.variant.(oneAccImm2WordT)
		wd := memory.DwordGetLowerWord(cpu.ac[oneAccImm2Word.acd])
		if dg.WordT(oneAccImm2Word.immS16)&wd == 0 {
			cpu.pc += 3
		} else {
			cpu.pc += 2
		}

	case instrSNOVR:
		if cpu.GetOVR() {
			cpu.pc++
		} else {
			cpu.pc += 2
		}

	case instrWBR:
		//		if iPtr.disp > 0 {
		//			cpu.pc += dg_phys_addr(iPtr.disp)
		//		} else {
		//			cpu.pc -= dg_phys_addr(iPtr.disp)
		//		}
		split8bitDisp := iPtr.variant.(split8bitDispT)
		cpu.pc += dg.PhysAddrT(int32(split8bitDisp.disp8))

		// case WPOPB: // FIXME - not yet decoded!
		// 	wpopb(cpu)

	case instrWCLM:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		var h, l int32
		v := int32(cpu.ac[twoAcc1Word.acs])
		if twoAcc1Word.acs != twoAcc1Word.acd {
			l = int32(memory.ReadDWord(dg.PhysAddrT(cpu.ac[twoAcc1Word.acd])))
			h = int32(memory.ReadDWord(dg.PhysAddrT(cpu.ac[twoAcc1Word.acd+2])))
			if v >= l && v <= h {
				cpu.pc += 2
			} else {
				cpu.pc++
			}
		} else {
			l = int32(memory.ReadDWord(cpu.pc + 1))
			h = int32(memory.ReadDWord(cpu.pc + 3))
			if v >= l && v <= h {
				cpu.pc += 6
			} else {
				cpu.pc += 5
			}
		}

	case instrWPOPJ:
		dwd := wsPop(cpu, 0)
		cpu.pc = cpu.pc&ringMask32 | (dg.PhysAddrT(dwd) & 0x0fff_ffff)
		cpu.SetOVR(false)

	case instrWRTN: // FIXME incomplete: handle PSR and rings
		// set WSP equal to WFP
		cpu.wsp = cpu.wfp
		wpopb(cpu)

	case instrWSEQ: // Signedness doen't matter for equality testing
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		var dwd dg.DwordT
		if twoAcc1Word.acd == twoAcc1Word.acs {
			dwd = 0
		} else {
			dwd = cpu.ac[twoAcc1Word.acd]
		}
		if cpu.ac[twoAcc1Word.acs] == dwd {
			cpu.pc += 2
		} else {
			cpu.pc++
		}

	case instrWSEQI:
		oneAccImm2Word := iPtr.variant.(oneAccImm2WordT)
		if cpu.ac[oneAccImm2Word.acd] == dg.DwordT(int32(oneAccImm2Word.immS16)) {
			cpu.pc += 3
		} else {
			cpu.pc += 2
		}

	case instrWSGE: // wide signed
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		var s32s, s32d int32
		if twoAcc1Word.acd == twoAcc1Word.acs {
			s32d = 0
		} else {
			s32d = int32(cpu.ac[twoAcc1Word.acd]) // this does the right thing in Go
		}
		s32s = int32(cpu.ac[twoAcc1Word.acs])
		if s32s >= s32d {
			cpu.pc += 2
		} else {
			cpu.pc++
		}

	case instrWSGT:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		var s32s, s32d int32
		if twoAcc1Word.acd == twoAcc1Word.acs {
			s32d = 0
		} else {
			s32d = int32(cpu.ac[twoAcc1Word.acd]) // this does the right thing in Go
		}
		s32s = int32(cpu.ac[twoAcc1Word.acs])
		if s32s > s32d {
			cpu.pc += 2
		} else {
			cpu.pc++
		}

	case instrWSKBO:
		wskb := iPtr.variant.(wskbT)
		if memory.TestDwbit(cpu.ac[0], wskb.bitNum) {
			cpu.pc += 2
		} else {
			cpu.pc++
		}

	case instrWSGTI:
		oneAccImm2Word := iPtr.variant.(oneAccImm2WordT)
		if int32(cpu.ac[oneAccImm2Word.acd]) > int32(oneAccImm2Word.immS16) {
			cpu.pc += 3
		} else {
			cpu.pc += 2
		}

	case instrWSKBZ:
		wskb := iPtr.variant.(wskbT)
		if !memory.TestDwbit(cpu.ac[0], wskb.bitNum) {
			cpu.pc += 2
		} else {
			cpu.pc++
		}

	case instrWSLE:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		var s32s, s32d int32
		if twoAcc1Word.acd == twoAcc1Word.acs {
			s32d = 0
		} else {
			s32d = int32(cpu.ac[twoAcc1Word.acd]) // this does the right thing in Go
		}
		s32s = int32(cpu.ac[twoAcc1Word.acs])
		if s32s <= s32d {
			cpu.pc += 2
		} else {
			cpu.pc++
		}

	case instrWSLEI:
		oneAccImm2Word := iPtr.variant.(oneAccImm2WordT)
		if int32(cpu.ac[oneAccImm2Word.acd]) <= int32(oneAccImm2Word.immS16) {
			cpu.pc += 3
		} else {
			cpu.pc += 2
		}

	case instrWSLT:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		var s32s, s32d int32
		if twoAcc1Word.acd == twoAcc1Word.acs {
			s32d = 0
		} else {
			s32d = int32(cpu.ac[twoAcc1Word.acd]) // this does the right thing in Go
		}
		s32s = int32(cpu.ac[twoAcc1Word.acs])
		if s32s < s32d {
			cpu.pc += 2
		} else {
			cpu.pc++
		}

	case instrWSNB:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		tmpAddr, bit := resolveEagleBitAddr(cpu, &twoAcc1Word)
		wd := memory.ReadWord(tmpAddr)
		if memory.TestWbit(wd, int(bit)) {
			cpu.pc += 2
		} else {
			cpu.pc++
		}
		if cpu.debugLogging {
			logging.DebugPrint(logging.DebugLog, ".... Wd Addr: %d., word: %0X, bit #: %d\n", tmpAddr, wd, bit)
		}

	case instrWSNE:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		var dwd dg.DwordT
		if twoAcc1Word.acd == twoAcc1Word.acs {
			dwd = 0
		} else {
			dwd = cpu.ac[twoAcc1Word.acd]
		}
		if cpu.ac[twoAcc1Word.acs] != dwd {
			cpu.pc += 2
		} else {
			cpu.pc++
		}

	case instrWSNEI:
		oneAccImm2Word := iPtr.variant.(oneAccImm2WordT)
		tmp32b := dg.DwordT(int32(oneAccImm2Word.immS16))
		if cpu.ac[oneAccImm2Word.acd] != tmp32b {
			cpu.pc += 3
		} else {
			cpu.pc += 2
		}

	case instrWSZB:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		tmpAddr, bit := resolveEagleBitAddr(cpu, &twoAcc1Word)
		wd := memory.ReadWord(tmpAddr)
		if !memory.TestWbit(wd, int(bit)) {
			cpu.pc += 2
		} else {
			cpu.pc++
		}
		if cpu.debugLogging {
			logging.DebugPrint(logging.DebugLog, ".... Wd Addr: %d., word: %0X, bit #: %d\n", tmpAddr, wd, bit)
		}

	case instrWUGTI:
		oneAccImm3Word := iPtr.variant.(oneAccImm3WordT)
		if uint32(cpu.ac[oneAccImm3Word.acd]) > oneAccImm3Word.immU32 {
			cpu.pc += 4
		} else {
			cpu.pc += 3
		}

	case instrWUSGT:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		if twoAcc1Word.acs == twoAcc1Word.acd {
			if cpu.ac[twoAcc1Word.acs] > 0 {
				cpu.pc += 2
			} else {
				cpu.pc++
			}
		} else {
			if cpu.ac[twoAcc1Word.acs] > cpu.ac[twoAcc1Word.acd] {
				cpu.pc += 2
			} else {
				cpu.pc++
			}
		}

	case instrXCALL:
		noAccModeInd3WordXcall := iPtr.variant.(noAccModeInd3WordXcallT)
		// FIXME - only handling the trivial case so far
		cpu.ac[3] = dg.DwordT(cpu.pc) + 3
		var dwd dg.DwordT
		if noAccModeInd3WordXcall.argCount >= 0 {
			dwd = dg.DwordT(cpu.psr) << 16
			dwd |= dg.DwordT(noAccModeInd3WordXcall.argCount)
		} else {
			dwd = dg.DwordT(noAccModeInd3WordXcall.argCount) & 0x00007fff
		}
		wsPush(cpu, 0, dwd)
		cpu.pc = resolve15bitDisplacement(cpu, noAccModeInd3WordXcall.ind, noAccModeInd3WordXcall.mode,
			dg.WordT(noAccModeInd3WordXcall.disp15), iPtr.dispOffset)

	case instrXJMP:
		noAccModeInd2Word := iPtr.variant.(noAccModeInd2WordT)
		cpu.pc = cpu.pc&ringMask32 | resolve15bitDisplacement(cpu, noAccModeInd2Word.ind, noAccModeInd2Word.mode, dg.WordT(noAccModeInd2Word.disp15), iPtr.dispOffset)

	case instrXJSR:
		noAccModeInd2Word := iPtr.variant.(noAccModeInd2WordT)
		cpu.ac[3] = dg.DwordT(cpu.pc + 2) // TODO Check this, PoP is self-contradictory on p.11-642
		cpu.pc = cpu.pc&ringMask32 | resolve15bitDisplacement(cpu, noAccModeInd2Word.ind, noAccModeInd2Word.mode, dg.WordT(noAccModeInd2Word.disp15), iPtr.dispOffset)

	case instrXNDO: // Narrow Do Until Greater Than
		threeWordDo := iPtr.variant.(threeWordDoT)
		loopVarAddr := resolve15bitDisplacement(cpu, threeWordDo.ind, threeWordDo.mode, dg.WordT(threeWordDo.disp15), iPtr.dispOffset)
		loopVar := int32(memory.SexWordToDword(memory.DwordGetLowerWord(memory.ReadDWord(loopVarAddr))))
		loopVar++
		memory.WriteDWord(loopVarAddr, dg.DwordT(loopVar))
		acVar := int32(cpu.ac[threeWordDo.acd])
		cpu.ac[threeWordDo.acd] = dg.DwordT(loopVar)
		if loopVar > acVar {
			// loop ends
			cpu.pc = cpu.pc + 1 + dg.PhysAddrT(threeWordDo.offsetU16)
		} else {
			cpu.pc += dg.PhysAddrT(iPtr.instrLength)
		}

	case instrXNDSZ: // unsigned narrow increment and skip if zero
		noAccModeInd2Word := iPtr.variant.(noAccModeInd2WordT)
		tmpAddr := resolve15bitDisplacement(cpu, noAccModeInd2Word.ind, noAccModeInd2Word.mode, dg.WordT(noAccModeInd2Word.disp15), iPtr.dispOffset)
		wd := memory.ReadWord(tmpAddr)
		wd-- // N.B. have checked that 0xffff + 1 == 0 in Go
		memory.WriteWord(tmpAddr, wd)
		if wd == 0 {
			cpu.pc += 3
		} else {
			cpu.pc += 2
		}

	case instrXNISZ: // unsigned narrow increment and skip if zero
		noAccModeInd2Word := iPtr.variant.(noAccModeInd2WordT)
		tmpAddr := resolve15bitDisplacement(cpu, noAccModeInd2Word.ind, noAccModeInd2Word.mode, dg.WordT(noAccModeInd2Word.disp15), iPtr.dispOffset)
		wd := memory.ReadWord(tmpAddr)
		wd++ // N.B. have checked that 0xffff + 1 == 0 in Go
		memory.WriteWord(tmpAddr, wd)
		if wd == 0 {
			cpu.pc += 3
		} else {
			cpu.pc += 2
		}

	case instrXWDSZ:
		noAccModeInd2Word := iPtr.variant.(noAccModeInd2WordT)
		tmpAddr := resolve15bitDisplacement(cpu, noAccModeInd2Word.ind, noAccModeInd2Word.mode, dg.WordT(noAccModeInd2Word.disp15), iPtr.dispOffset)
		dwd := memory.ReadDWord(tmpAddr)
		dwd--
		memory.WriteDWord(tmpAddr, dwd)
		if dwd == 0 {
			cpu.pc += 3
		} else {
			cpu.pc += 2
		}

	case instrXWISZ:
		noAccModeInd2Word := iPtr.variant.(noAccModeInd2WordT)
		tmpAddr := resolve15bitDisplacement(cpu, noAccModeInd2Word.ind, noAccModeInd2Word.mode, dg.WordT(noAccModeInd2Word.disp15), iPtr.dispOffset)
		dwd := memory.ReadDWord(tmpAddr)
		dwd++
		memory.WriteDWord(tmpAddr, dwd)
		if dwd == 0 {
			cpu.pc += 3
		} else {
			cpu.pc += 2
		}

	default:
		log.Fatalf("ERROR: EAGLE_PC instruction <%s> not yet implemented\n", iPtr.mnemonic)
		return false
	}

	return true
}
