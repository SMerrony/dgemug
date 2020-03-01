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

func eaglePC(cpuPtr *MvCPUT, iPtr *decodedInstrT) bool {

	switch iPtr.ix {

	case instrDSZTS, instrISZTS:
		// tmpAddr := dg.PhysAddrT(memory.ReadDWord(cpuPtr.wsp))
		tmpAddr := dg.PhysAddrT(cpuPtr.wsp)
		var dwd dg.DwordT
		if iPtr.ix == instrDSZTS {
			dwd = memory.ReadDWord(tmpAddr) - 1
		} else {
			dwd = memory.ReadDWord(tmpAddr) + 1
		}
		memory.WriteDWord(tmpAddr, dwd)
		cpuPtr.CPUSetOVR(false)
		if dwd == 0 {
			cpuPtr.pc += 2
		} else {
			cpuPtr.pc++
		}
		if debugLogging {
			logging.DebugPrint(logging.DebugLog, "..... wrote %d to %d\n", dwd, tmpAddr)
		}

	case instrLCALL: // FIXME - LCALL only handling trivial case, no checking
		noAccModeInd4Word := iPtr.variant.(noAccModeInd4WordT)
		cpuPtr.ac[3] = dg.DwordT(cpuPtr.pc) + 4
		var dwd dg.DwordT
		if noAccModeInd4Word.argCount >= 0 {
			dwd = memory.DwordFromTwoWords(cpuPtr.psr, dg.WordT(noAccModeInd4Word.argCount))
		} else {
			dwd = dg.DwordT(noAccModeInd4Word.argCount) & 0x00007fff
		}
		wsPush(cpuPtr, 0, dwd)
		cpuPtr.CPUSetOVR(false)
		cpuPtr.pc = resolve32bitEffAddr(cpuPtr, noAccModeInd4Word.ind, noAccModeInd4Word.mode, noAccModeInd4Word.disp31, iPtr.dispOffset)

	case instrLDSP:
		oneAccModeInd3Word := iPtr.variant.(oneAccModeInd3WordT)
		value := int32(cpuPtr.ac[oneAccModeInd3Word.acd])
		tableAddr := resolve32bitEffAddr(cpuPtr, oneAccModeInd3Word.ind, oneAccModeInd3Word.mode, oneAccModeInd3Word.disp31, iPtr.dispOffset)
		h := int32(memory.ReadDWord(tableAddr - 2))
		l := int32(memory.ReadDWord(tableAddr - 4))
		if value < l || value > h {
			cpuPtr.pc += 3
		} else {
			tableIndex := tableAddr + (2 * dg.PhysAddrT(value)) - (2 * dg.PhysAddrT(l))
			tableVal := memory.ReadDWord(tableIndex)
			if tableVal == 0xFFFFFFFF {
				cpuPtr.pc += 3
			} else {
				cpuPtr.pc = dg.PhysAddrT(tableVal) + tableIndex
			}
		}

	case instrLJMP:
		noAccModeInd3Word := iPtr.variant.(noAccModeInd3WordT)
		cpuPtr.pc = resolve32bitEffAddr(cpuPtr, noAccModeInd3Word.ind, noAccModeInd3Word.mode, noAccModeInd3Word.disp31, iPtr.dispOffset)

	case instrLJSR:
		noAccModeInd3Word := iPtr.variant.(noAccModeInd3WordT)
		cpuPtr.ac[3] = dg.DwordT(cpuPtr.pc) + 3
		cpuPtr.pc = resolve32bitEffAddr(cpuPtr, noAccModeInd3Word.ind, noAccModeInd3Word.mode, noAccModeInd3Word.disp31, iPtr.dispOffset)

	case instrLNISZ:
		noAccModeInd3Word := iPtr.variant.(noAccModeInd3WordT)
		// unsigned narrow increment and skip if zero
		tmpAddr := resolve32bitEffAddr(cpuPtr, noAccModeInd3Word.ind, noAccModeInd3Word.mode, noAccModeInd3Word.disp31, iPtr.dispOffset)
		wd := memory.ReadWord(tmpAddr) + 1
		memory.WriteWord(tmpAddr, wd)
		if wd == 0 {
			cpuPtr.pc += 4
		} else {
			cpuPtr.pc += 3
		}

	case instrLPSHJ:
		noAccModeInd3Word := iPtr.variant.(noAccModeInd3WordT)
		wsPush(cpuPtr, 0, dg.DwordT(cpuPtr.pc)+3)
		cpuPtr.pc = resolve32bitEffAddr(cpuPtr, noAccModeInd3Word.ind, noAccModeInd3Word.mode, noAccModeInd3Word.disp31, iPtr.dispOffset)

	case instrLWDSZ:
		noAccModeInd3Word := iPtr.variant.(noAccModeInd3WordT)
		// unsigned wide decrement and skip if zero
		tmpAddr := resolve32bitEffAddr(cpuPtr, noAccModeInd3Word.ind, noAccModeInd3Word.mode, noAccModeInd3Word.disp31, iPtr.dispOffset)
		tmp32b := memory.ReadDWord(tmpAddr) - 1
		memory.WriteDWord(tmpAddr, tmp32b)
		if tmp32b == 0 {
			cpuPtr.pc += 4
		} else {
			cpuPtr.pc += 3
		}

	case instrNSALA:
		oneAccImm2Word := iPtr.variant.(oneAccImm2WordT)
		wd := ^memory.DwordGetLowerWord(cpuPtr.ac[oneAccImm2Word.acd])
		if dg.WordT(oneAccImm2Word.immS16)&wd == 0 {
			cpuPtr.pc += 3
		} else {
			cpuPtr.pc += 2
		}

	case instrNSANA:
		oneAccImm2Word := iPtr.variant.(oneAccImm2WordT)
		wd := memory.DwordGetLowerWord(cpuPtr.ac[oneAccImm2Word.acd])
		if dg.WordT(oneAccImm2Word.immS16)&wd == 0 {
			cpuPtr.pc += 3
		} else {
			cpuPtr.pc += 2
		}

	case instrSNOVR:
		if cpuPtr.CPUGetOVR() {
			cpuPtr.pc++
		} else {
			cpuPtr.pc += 2
		}

	case instrWBR:
		//		if iPtr.disp > 0 {
		//			cpuPtr.pc += dg_phys_addr(iPtr.disp)
		//		} else {
		//			cpuPtr.pc -= dg_phys_addr(iPtr.disp)
		//		}
		split8bitDisp := iPtr.variant.(split8bitDispT)
		cpuPtr.pc += dg.PhysAddrT(int32(split8bitDisp.disp8))

		// case WPOPB: // FIXME - not yet decoded!
		// 	wpopb(cpuPtr)

	case instrWCLM:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		var h, l int32
		v := int32(cpuPtr.ac[twoAcc1Word.acs])
		if twoAcc1Word.acs != twoAcc1Word.acd {
			l = int32(memory.ReadDWord(dg.PhysAddrT(cpuPtr.ac[twoAcc1Word.acd])))
			h = int32(memory.ReadDWord(dg.PhysAddrT(cpuPtr.ac[twoAcc1Word.acd+2])))
			if v >= l && v <= h {
				cpuPtr.pc += 2
			} else {
				cpuPtr.pc++
			}
		} else {
			l = int32(memory.ReadDWord(cpuPtr.pc + 1))
			h = int32(memory.ReadDWord(cpuPtr.pc + 3))
			if v >= l && v <= h {
				cpuPtr.pc += 6
			} else {
				cpuPtr.pc += 5
			}
		}

	case instrWPOPJ:
		dwd := wsPop(cpuPtr, 0)
		cpuPtr.pc = (cpuPtr.pc & 0xf000_0000) | (dg.PhysAddrT(dwd) & 0x0fff_ffff)
		cpuPtr.CPUSetOVR(false)

	case instrWRTN: // FIXME incomplete: handle PSR and rings
		// set WSP equal to WFP
		cpuPtr.wsp = cpuPtr.wfp
		wpopb(cpuPtr)

	case instrWSEQ: // Signedness doen't matter for equality testing
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		var dwd dg.DwordT
		if twoAcc1Word.acd == twoAcc1Word.acs {
			dwd = 0
		} else {
			dwd = cpuPtr.ac[twoAcc1Word.acd]
		}
		if cpuPtr.ac[twoAcc1Word.acs] == dwd {
			cpuPtr.pc += 2
		} else {
			cpuPtr.pc++
		}

	case instrWSEQI:
		oneAccImm2Word := iPtr.variant.(oneAccImm2WordT)
		if cpuPtr.ac[oneAccImm2Word.acd] == dg.DwordT(int32(oneAccImm2Word.immS16)) {
			cpuPtr.pc += 3
		} else {
			cpuPtr.pc += 2
		}

	case instrWSGE: // wide signed
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		var s32s, s32d int32
		if twoAcc1Word.acd == twoAcc1Word.acs {
			s32d = 0
		} else {
			s32d = int32(cpuPtr.ac[twoAcc1Word.acd]) // this does the right thing in Go
		}
		s32s = int32(cpuPtr.ac[twoAcc1Word.acs])
		if s32s >= s32d {
			cpuPtr.pc += 2
		} else {
			cpuPtr.pc++
		}

	case instrWSGT:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		var s32s, s32d int32
		if twoAcc1Word.acd == twoAcc1Word.acs {
			s32d = 0
		} else {
			s32d = int32(cpuPtr.ac[twoAcc1Word.acd]) // this does the right thing in Go
		}
		s32s = int32(cpuPtr.ac[twoAcc1Word.acs])
		if s32s > s32d {
			cpuPtr.pc += 2
		} else {
			cpuPtr.pc++
		}

	case instrWSKBO:
		wskb := iPtr.variant.(wskbT)
		if memory.TestDwbit(cpuPtr.ac[0], wskb.bitNum) {
			cpuPtr.pc += 2
		} else {
			cpuPtr.pc++
		}

	case instrWSGTI:
		oneAccImm2Word := iPtr.variant.(oneAccImm2WordT)
		if int32(cpuPtr.ac[oneAccImm2Word.acd]) > int32(oneAccImm2Word.immS16) {
			cpuPtr.pc += 3
		} else {
			cpuPtr.pc += 2
		}

	case instrWSKBZ:
		wskb := iPtr.variant.(wskbT)
		if !memory.TestDwbit(cpuPtr.ac[0], wskb.bitNum) {
			cpuPtr.pc += 2
		} else {
			cpuPtr.pc++
		}

	case instrWSLE:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		var s32s, s32d int32
		if twoAcc1Word.acd == twoAcc1Word.acs {
			s32d = 0
		} else {
			s32d = int32(cpuPtr.ac[twoAcc1Word.acd]) // this does the right thing in Go
		}
		s32s = int32(cpuPtr.ac[twoAcc1Word.acs])
		if s32s <= s32d {
			cpuPtr.pc += 2
		} else {
			cpuPtr.pc++
		}

	case instrWSLEI:
		oneAccImm2Word := iPtr.variant.(oneAccImm2WordT)
		if int32(cpuPtr.ac[oneAccImm2Word.acd]) <= int32(oneAccImm2Word.immS16) {
			cpuPtr.pc += 3
		} else {
			cpuPtr.pc += 2
		}

	case instrWSLT:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		var s32s, s32d int32
		if twoAcc1Word.acd == twoAcc1Word.acs {
			s32d = 0
		} else {
			s32d = int32(cpuPtr.ac[twoAcc1Word.acd]) // this does the right thing in Go
		}
		s32s = int32(cpuPtr.ac[twoAcc1Word.acs])
		if s32s < s32d {
			cpuPtr.pc += 2
		} else {
			cpuPtr.pc++
		}

	case instrWSNB:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		tmpAddr, bit := resolveEagleBitAddr(cpuPtr, &twoAcc1Word)
		wd := memory.ReadWord(tmpAddr)
		if memory.TestWbit(wd, int(bit)) {
			cpuPtr.pc += 2
		} else {
			cpuPtr.pc++
		}
		if debugLogging {
			logging.DebugPrint(logging.DebugLog, ".... Wd Addr: %d., word: %0X, bit #: %d\n", tmpAddr, wd, bit)
		}

	case instrWSNE:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		var dwd dg.DwordT
		if twoAcc1Word.acd == twoAcc1Word.acs {
			dwd = 0
		} else {
			dwd = cpuPtr.ac[twoAcc1Word.acd]
		}
		if cpuPtr.ac[twoAcc1Word.acs] != dwd {
			cpuPtr.pc += 2
		} else {
			cpuPtr.pc++
		}

	case instrWSNEI:
		oneAccImm2Word := iPtr.variant.(oneAccImm2WordT)
		tmp32b := dg.DwordT(int32(oneAccImm2Word.immS16))
		if cpuPtr.ac[oneAccImm2Word.acd] != tmp32b {
			cpuPtr.pc += 3
		} else {
			cpuPtr.pc += 2
		}

	case instrWSZB:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		tmpAddr, bit := resolveEagleBitAddr(cpuPtr, &twoAcc1Word)
		wd := memory.ReadWord(tmpAddr)
		if !memory.TestWbit(wd, int(bit)) {
			cpuPtr.pc += 2
		} else {
			cpuPtr.pc++
		}
		if debugLogging {
			logging.DebugPrint(logging.DebugLog, ".... Wd Addr: %d., word: %0X, bit #: %d\n", tmpAddr, wd, bit)
		}

	case instrWUGTI:
		oneAccImm3Word := iPtr.variant.(oneAccImm3WordT)
		if uint32(cpuPtr.ac[oneAccImm3Word.acd]) > oneAccImm3Word.immU32 {
			cpuPtr.pc += 4
		} else {
			cpuPtr.pc += 3
		}

	case instrWUSGT:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		if twoAcc1Word.acs == twoAcc1Word.acd {
			if cpuPtr.ac[twoAcc1Word.acs] > 0 {
				cpuPtr.pc += 2
			} else {
				cpuPtr.pc++
			}
		} else {
			if cpuPtr.ac[twoAcc1Word.acs] > cpuPtr.ac[twoAcc1Word.acd] {
				cpuPtr.pc += 2
			} else {
				cpuPtr.pc++
			}
		}

	case instrXCALL:
		noAccModeInd3WordXcall := iPtr.variant.(noAccModeInd3WordXcallT)
		// FIXME - only handling the trivial case so far
		cpuPtr.ac[3] = dg.DwordT(cpuPtr.pc) + 3
		var dwd dg.DwordT
		if noAccModeInd3WordXcall.argCount >= 0 {
			dwd = dg.DwordT(cpuPtr.psr) << 16
			dwd |= dg.DwordT(noAccModeInd3WordXcall.argCount)
		} else {
			dwd = dg.DwordT(noAccModeInd3WordXcall.argCount) & 0x00007fff
		}
		wsPush(cpuPtr, 0, dwd)
		cpuPtr.pc = resolve15bitDisplacement(cpuPtr, noAccModeInd3WordXcall.ind, noAccModeInd3WordXcall.mode,
			dg.WordT(noAccModeInd3WordXcall.disp15), iPtr.dispOffset)

	case instrXJMP:
		noAccModeInd2Word := iPtr.variant.(noAccModeInd2WordT)
		cpuPtr.pc = resolve15bitDisplacement(cpuPtr, noAccModeInd2Word.ind, noAccModeInd2Word.mode, dg.WordT(noAccModeInd2Word.disp15), iPtr.dispOffset)

	case instrXJSR:
		noAccModeInd2Word := iPtr.variant.(noAccModeInd2WordT)
		cpuPtr.ac[3] = dg.DwordT(cpuPtr.pc + 2) // TODO Check this, PoP is self-contradictory on p.11-642
		cpuPtr.pc = resolve15bitDisplacement(cpuPtr, noAccModeInd2Word.ind, noAccModeInd2Word.mode, dg.WordT(noAccModeInd2Word.disp15), iPtr.dispOffset)

	case instrXNDO: // Narrow Do Until Greater Than
		threeWordDo := iPtr.variant.(threeWordDoT)
		loopVarAddr := resolve15bitDisplacement(cpuPtr, threeWordDo.ind, threeWordDo.mode, dg.WordT(threeWordDo.disp15), iPtr.dispOffset)
		loopVar := int32(memory.SexWordToDword(memory.DwordGetLowerWord(memory.ReadDWord(loopVarAddr))))
		loopVar++
		memory.WriteDWord(loopVarAddr, dg.DwordT(loopVar))
		acVar := int32(cpuPtr.ac[threeWordDo.acd])
		cpuPtr.ac[threeWordDo.acd] = dg.DwordT(loopVar)
		if loopVar > acVar {
			// loop ends
			cpuPtr.pc = cpuPtr.pc + 1 + dg.PhysAddrT(threeWordDo.offsetU16)
		} else {
			cpuPtr.pc += dg.PhysAddrT(iPtr.instrLength)
		}

	case instrXNDSZ: // unsigned narrow increment and skip if zero
		noAccModeInd2Word := iPtr.variant.(noAccModeInd2WordT)
		tmpAddr := resolve15bitDisplacement(cpuPtr, noAccModeInd2Word.ind, noAccModeInd2Word.mode, dg.WordT(noAccModeInd2Word.disp15), iPtr.dispOffset)
		wd := memory.ReadWord(tmpAddr)
		wd-- // N.B. have checked that 0xffff + 1 == 0 in Go
		memory.WriteWord(tmpAddr, wd)
		if wd == 0 {
			cpuPtr.pc += 3
		} else {
			cpuPtr.pc += 2
		}

	case instrXNISZ: // unsigned narrow increment and skip if zero
		noAccModeInd2Word := iPtr.variant.(noAccModeInd2WordT)
		tmpAddr := resolve15bitDisplacement(cpuPtr, noAccModeInd2Word.ind, noAccModeInd2Word.mode, dg.WordT(noAccModeInd2Word.disp15), iPtr.dispOffset)
		wd := memory.ReadWord(tmpAddr)
		wd++ // N.B. have checked that 0xffff + 1 == 0 in Go
		memory.WriteWord(tmpAddr, wd)
		if wd == 0 {
			cpuPtr.pc += 3
		} else {
			cpuPtr.pc += 2
		}

	case instrXWDSZ:
		noAccModeInd2Word := iPtr.variant.(noAccModeInd2WordT)
		tmpAddr := resolve15bitDisplacement(cpuPtr, noAccModeInd2Word.ind, noAccModeInd2Word.mode, dg.WordT(noAccModeInd2Word.disp15), iPtr.dispOffset)
		dwd := memory.ReadDWord(tmpAddr)
		dwd--
		memory.WriteDWord(tmpAddr, dwd)
		if dwd == 0 {
			cpuPtr.pc += 3
		} else {
			cpuPtr.pc += 2
		}

	case instrXWISZ:
		noAccModeInd2Word := iPtr.variant.(noAccModeInd2WordT)
		tmpAddr := resolve15bitDisplacement(cpuPtr, noAccModeInd2Word.ind, noAccModeInd2Word.mode, dg.WordT(noAccModeInd2Word.disp15), iPtr.dispOffset)
		dwd := memory.ReadDWord(tmpAddr)
		dwd++
		memory.WriteDWord(tmpAddr, dwd)
		if dwd == 0 {
			cpuPtr.pc += 3
		} else {
			cpuPtr.pc += 2
		}

	default:
		log.Fatalf("ERROR: EAGLE_PC instruction <%s> not yet implemented\n", iPtr.mnemonic)
		return false
	}

	return true
}
