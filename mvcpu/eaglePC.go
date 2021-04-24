// eaglePC.go

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
	"log"

	"github.com/SMerrony/dgemug/dg"
	"github.com/SMerrony/dgemug/logging"
	"github.com/SMerrony/dgemug/memory"
)

func eaglePC(cpu *CPUT, iPtr *decodedInstrT) bool {

	switch iPtr.ix {

	case instrDERR: // TODO stack overflow checking
		derr := iPtr.variant.(derrT)
		wsPush(cpu, dg.DwordT(cpu.pc))
		wsPush(cpu, dg.DwordT(derr.errCode))
		cpu.pc = cpu.pc&ringMask32 | dg.PhysAddrT(memory.ReadWord(cpu.pc&ringMask32|047))
		if cpu.debugLogging {
			logging.DebugPrint(logging.DebugLog, "..... DERR handler at: %#x\n", cpu.pc)
		}

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
		pc_plus_4 := dg.DwordT(cpu.pc) + 4
		var dwd dg.DwordT
		if iPtr.argCount >= 0 {
			dwd = dg.DwordT(iPtr.argCount)
		} else {
			//dwd = dg.DwordT(iPtr.argCount) & 0x00007fff
			dwd = memory.ReadDWord(cpu.wsp) & 0x0000_7fff
		}
		dwd |= dg.DwordT(cpu.psr) << 16
		ok, faultCode, secondaryFault := wspCheckBounds(cpu, 2, false)
		if !ok {
			//log.Panicf("DEBUG: Stack fault trapped in LCALL, codes %d and %d", faultCode, secondaryFault)
			wspHandleFault(cpu, iPtr.instrLength, faultCode, secondaryFault)
		}
		wsPush(cpu, dwd)
		cpu.SetOVR(false)
		cpu.pc = resolve31bitDisplacement(cpu, iPtr.ind, iPtr.mode, iPtr.disp31, iPtr.dispOffset)
		cpu.ac[3] = pc_plus_4

	case instrLDSP:
		oneAccModeInd3Word := iPtr.variant.(oneAccModeInd3WordT)
		value := int32(cpu.ac[oneAccModeInd3Word.acd])
		tableAddr := resolve31bitDisplacement(cpu, oneAccModeInd3Word.ind, oneAccModeInd3Word.mode, oneAccModeInd3Word.disp31, iPtr.dispOffset)
		h := int32(memory.ReadDWord(tableAddr - 2))
		l := int32(memory.ReadDWord(tableAddr - 4))
		if value < l || value > h {
			cpu.pc += 3
		} else {
			tableIndex := tableAddr + (2 * dg.PhysAddrT(value)) - (2 * dg.PhysAddrT(l))
			tableVal := memory.ReadDWord(tableIndex)
			if memory.TestDwbit(tableVal, 4) { // sign-extend from 28-bits
				tableVal |= 0xF000_0000
			}
			if tableVal == 0xFFFF_FFFF {
				cpu.pc += 3
			} else {
				cpu.pc = dg.PhysAddrT(int32(tableIndex) + int32(tableVal))
			}
		}

	case instrLJMP:
		noAccModeInd3Word := iPtr.variant.(noAccModeInd3WordT)
		cpu.pc = cpu.pc&ringMask32 | resolve31bitDisplacement(cpu, noAccModeInd3Word.ind, noAccModeInd3Word.mode, noAccModeInd3Word.disp31, iPtr.dispOffset)

	case instrLJSR:
		noAccModeInd3Word := iPtr.variant.(noAccModeInd3WordT)
		cpu.ac[3] = dg.DwordT(cpu.pc) + 3
		cpu.pc = cpu.pc&ringMask32 | resolve31bitDisplacement(cpu, noAccModeInd3Word.ind, noAccModeInd3Word.mode, noAccModeInd3Word.disp31, iPtr.dispOffset)

	case instrLNDO:
		lndo4Word := iPtr.variant.(lndo4WordT)
		count := int32(cpu.ac[lndo4Word.acd])
		memVarAddr := resolve31bitDisplacement(cpu, lndo4Word.ind, lndo4Word.mode, lndo4Word.disp31, iPtr.dispOffset)
		memVar := int32(int16(memory.ReadWord(memVarAddr))) + 1
		memory.WriteWord(memVarAddr, dg.WordT(memVar))
		cpu.ac[lndo4Word.acd] = dg.DwordT(memVar)
		if memVar > count {
			// loop ends
			cpu.pc += dg.PhysAddrT(lndo4Word.offsetU16) + 1
		} else {
			// loop continues
			cpu.pc += 4
		}

	case instrLNDSZ, instrLNISZ:
		noAccModeInd3Word := iPtr.variant.(noAccModeInd3WordT)
		tmpAddr := resolve31bitDisplacement(cpu, noAccModeInd3Word.ind, noAccModeInd3Word.mode, noAccModeInd3Word.disp31, iPtr.dispOffset)
		wd := memory.ReadWord(tmpAddr)
		if iPtr.ix == instrLNDSZ {
			wd--
		} else {
			wd++
		}
		memory.WriteWord(tmpAddr, wd)
		if wd == 0 {
			cpu.pc += 4
		} else {
			cpu.pc += 3
		}

	case instrLPSHJ:
		noAccModeInd3Word := iPtr.variant.(noAccModeInd3WordT)
		wsPush(cpu, dg.DwordT(cpu.pc)+3)
		cpu.pc = resolve31bitDisplacement(cpu, noAccModeInd3Word.ind, noAccModeInd3Word.mode, noAccModeInd3Word.disp31, iPtr.dispOffset)

	case instrLWDO: // Wide Do Until Greater Than
		lndo4Word := iPtr.variant.(lndo4WordT)
		count := int32(cpu.ac[lndo4Word.acd])
		memVarAddr := resolve31bitDisplacement(cpu, lndo4Word.ind, lndo4Word.mode, lndo4Word.disp31, iPtr.dispOffset)
		memVar := int32(memory.ReadDWord(memVarAddr)) + 1
		memory.WriteDWord(memVarAddr, dg.DwordT(memVar))
		cpu.ac[lndo4Word.acd] = dg.DwordT(memVar)
		if memVar > count {
			// loop ends
			cpu.pc += dg.PhysAddrT(lndo4Word.offsetU16) + 1
		} else {
			// loop continues
			cpu.pc += 4
		}

	case instrLWDSZ, instrLWISZ:
		noAccModeInd3Word := iPtr.variant.(noAccModeInd3WordT)
		tmpAddr := resolve31bitDisplacement(cpu, noAccModeInd3Word.ind, noAccModeInd3Word.mode, noAccModeInd3Word.disp31, iPtr.dispOffset)
		tmp32b := memory.ReadDWord(tmpAddr)
		if iPtr.ix == instrLWDSZ {
			tmp32b--
		} else {
			tmp32b++
		}
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
		split8bitDisp := iPtr.variant.(split8bitDispT)
		cpu.pc += dg.PhysAddrT(int32(split8bitDisp.disp8))

	case instrWCLM:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		var h, l int32
		v := int32(cpu.ac[twoAcc1Word.acs])
		if twoAcc1Word.acs != twoAcc1Word.acd {
			l = int32(memory.ReadDWord(dg.PhysAddrT(cpu.ac[twoAcc1Word.acd])))
			h = int32(memory.ReadDWord(dg.PhysAddrT(cpu.ac[twoAcc1Word.acd] + 2)))
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

	case instrWMESS:
		dwd := memory.ReadDWord(dg.PhysAddrT(cpu.ac[2]))
		ord := dwd ^ cpu.ac[0]
		if ord&cpu.ac[3] == 0 {
			memory.WriteDWord(dg.PhysAddrT(cpu.ac[2]), cpu.ac[1])
			cpu.ac[1] = dwd
			cpu.pc += 2
		} else {
			cpu.ac[1] = dwd
			cpu.pc++
		}

	case instrWSANA:
		oneAccImm3Word := iPtr.variant.(oneAccImm3WordT)
		if uint32(cpu.ac[oneAccImm3Word.acd])&oneAccImm3Word.immU32 != 0 {
			cpu.pc += 4
		} else {
			cpu.pc += 3
		}

	case instrWSEQ, instrWSNE: // Signedness doen't matter for equality testing
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		var dwd dg.DwordT
		if twoAcc1Word.acd == twoAcc1Word.acs {
			dwd = 0
		} else {
			dwd = cpu.ac[twoAcc1Word.acd]
		}
		var skip bool
		switch iPtr.ix {
		case instrWSEQ:
			skip = cpu.ac[twoAcc1Word.acs] == dwd
		case instrWSNE:
			skip = cpu.ac[twoAcc1Word.acs] != dwd
		}
		if skip {
			cpu.pc += 2
		} else {
			cpu.pc++
		}

	case instrWSEQI, instrWSGTI, instrWSLEI, instrWSNEI:
		oneAccImm2Word := iPtr.variant.(oneAccImm2WordT)
		var skip bool
		switch iPtr.ix {
		case instrWSEQI:
			skip = int32(cpu.ac[oneAccImm2Word.acd]) == int32(oneAccImm2Word.immS16)
		case instrWSGTI:
			skip = int32(cpu.ac[oneAccImm2Word.acd]) >= int32(oneAccImm2Word.immS16)
		case instrWSLEI:
			skip = int32(cpu.ac[oneAccImm2Word.acd]) <= int32(oneAccImm2Word.immS16)
		case instrWSNEI:
			skip = int32(cpu.ac[oneAccImm2Word.acd]) != int32(oneAccImm2Word.immS16)
		}
		if skip {
			cpu.pc += 3
		} else {
			cpu.pc += 2
		}

	case instrWSGE, instrWSGT, instrWSLE, instrWSLT: // wide signed
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		var s32s, s32d int32
		if twoAcc1Word.acd == twoAcc1Word.acs {
			s32d = 0
		} else {
			s32d = int32(cpu.ac[twoAcc1Word.acd]) // this does the right thing in Go
		}
		s32s = int32(cpu.ac[twoAcc1Word.acs])
		var skip bool
		switch iPtr.ix {
		case instrWSGE:
			skip = s32s >= s32d
		case instrWSGT:
			skip = s32s > s32d
		case instrWSLE:
			skip = s32s <= s32d
		case instrWSLT:
			skip = s32s < s32d
		}
		if skip {
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

	case instrWSKBZ:
		wskb := iPtr.variant.(wskbT)
		if !memory.TestDwbit(cpu.ac[0], wskb.bitNum) {
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

	case instrWSZBO:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		tmpAddr, bit := resolveEagleBitAddr(cpu, &twoAcc1Word)
		wd := memory.ReadWord(tmpAddr)
		if !memory.TestWbit(wd, int(bit)) {
			memory.SetWbit(&wd, uint(bit))
			memory.WriteWord(tmpAddr, wd)
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

	case instrWULEI:
		oneAccImm3Word := iPtr.variant.(oneAccImm3WordT)
		if uint32(cpu.ac[oneAccImm3Word.acd]) <= oneAccImm3Word.immU32 {
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
		wsPush(cpu, dwd)
		cpu.pc = resolve15bitDisplacement(cpu, noAccModeInd3WordXcall.ind, noAccModeInd3WordXcall.mode,
			dg.WordT(noAccModeInd3WordXcall.disp15), iPtr.dispOffset)

	case instrXJMP:
		cpu.pc = cpu.pc&ringMask32 | resolve15bitDisplacement(cpu, iPtr.ind, iPtr.mode, dg.WordT(iPtr.disp15), iPtr.dispOffset)

	case instrXJSR:
		cpu.ac[3] = dg.DwordT(cpu.pc + 2) // TODO Check this, PoP is self-contradictory on p.11-642
		cpu.pc = cpu.pc&ringMask32 | resolve15bitDisplacement(cpu, iPtr.ind, iPtr.mode, dg.WordT(iPtr.disp15), iPtr.dispOffset)

	case instrXNDO: // Narrow Do Until Greater Than
		threeWordDo := iPtr.variant.(threeWordDoT)
		loopVarAddr := resolve15bitDisplacement(cpu, threeWordDo.ind, threeWordDo.mode, dg.WordT(threeWordDo.disp15), iPtr.dispOffset)
		//loopVar := int32(int16(memory.ReadWord(loopVarAddr + 1)))
		loopVar := int32(int16(memory.ReadWord(loopVarAddr)))
		loopVar++
		//memory.WriteDWord(loopVarAddr, dg.DwordT(loopVar))
		memory.WriteWord(loopVarAddr, dg.WordT(loopVar))
		acVar := int32(cpu.ac[threeWordDo.acd])
		//log.Printf("\t loopVar: %#x, acVar: %#x\n", loopVar, acVar)
		cpu.ac[threeWordDo.acd] = dg.DwordT(loopVar)
		if loopVar > acVar {
			// loop ends
			cpu.pc = cpu.pc + 1 + dg.PhysAddrT(threeWordDo.offsetU16)
			//log.Println("\tExiting loop")
		} else {
			cpu.pc += dg.PhysAddrT(iPtr.instrLength)
			//log.Println("\tLooping...")
		}

	case instrXNDSZ, instrXNISZ: // unsigned narrow inc/decrement and skip if zero
		tmpAddr := resolve15bitDisplacement(cpu, iPtr.ind, iPtr.mode, dg.WordT(iPtr.disp15), iPtr.dispOffset)
		wd := memory.ReadWord(tmpAddr)
		if iPtr.ix == instrXNDSZ {
			wd-- // N.B. have checked that 0xffff + 1 == 0 in Go
		} else {
			wd++
		}
		memory.WriteWord(tmpAddr, wd)
		if wd == 0 {
			cpu.pc += 3
		} else {
			cpu.pc += 2
		}

	case instrXWDO:
		threeWordDo := iPtr.variant.(threeWordDoT)
		loopVarAddr := resolve15bitDisplacement(cpu, threeWordDo.ind, threeWordDo.mode, dg.WordT(threeWordDo.disp15), iPtr.dispOffset)
		loopVar := int32(memory.ReadDWord(loopVarAddr))
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

	case instrXWDSZ:
		tmpAddr := resolve15bitDisplacement(cpu, iPtr.ind, iPtr.mode, dg.WordT(iPtr.disp15), iPtr.dispOffset)
		dwd := memory.ReadDWord(tmpAddr)
		dwd--
		memory.WriteDWord(tmpAddr, dwd)
		if dwd == 0 {
			cpu.pc += 3
		} else {
			cpu.pc += 2
		}

	case instrXWISZ:
		tmpAddr := resolve15bitDisplacement(cpu, iPtr.ind, iPtr.mode, dg.WordT(iPtr.disp15), iPtr.dispOffset)
		dwd := memory.ReadDWord(tmpAddr)
		dwd++
		memory.WriteDWord(tmpAddr, dwd)
		if dwd == 0 {
			cpu.pc += 3
		} else {
			cpu.pc += 2
		}

	default:
		log.Panicf("ERROR: EAGLE_PC instruction <%s> not yet implemented\n", iPtr.mnemonic)
		return false
	}

	return true
}
