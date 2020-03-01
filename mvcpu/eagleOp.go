// eagleOp.go

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
	"github.com/SMerrony/dgemug/memory"
)

func eagleOp(cpuPtr *CPUT, iPtr *decodedInstrT) bool {

	switch iPtr.ix {

	case instrCRYTC:
		cpuPtr.carry = !cpuPtr.carry

	case instrCRYTO:
		cpuPtr.carry = true

	case instrCRYTZ:
		cpuPtr.carry = false

	case instrCVWN:
		oneAcc1Word := iPtr.variant.(oneAcc1WordT)
		dwd := cpuPtr.ac[oneAcc1Word.acd]
		if dwd>>16 != 0 && dwd>>16 != 0xffff {
			cpuPtr.CPUSetOVR(true)
		}
		if memory.TestDwbit(dwd, 16) {
			cpuPtr.ac[oneAcc1Word.acd] |= 0xffff0000
		} else {
			cpuPtr.ac[oneAcc1Word.acd] &= 0x0000ffff
		}

	case instrLLDB:
		oneAccMode3Word := iPtr.variant.(oneAccMode3WordT)
		addr := resolve32bitEffAddr(cpuPtr, ' ', oneAccMode3Word.mode, oneAccMode3Word.disp31>>1, iPtr.dispOffset)
		lobyte := memory.TestDwbit(dg.DwordT(oneAccMode3Word.disp31), 31)
		cpuPtr.ac[oneAccMode3Word.acd] = dg.DwordT(memory.ReadByte(addr, lobyte))

	case instrLLEF:
		oneAccModeInd3Word := iPtr.variant.(oneAccModeInd3WordT)
		cpuPtr.ac[oneAccModeInd3Word.acd] = dg.DwordT(
			resolve32bitEffAddr(cpuPtr, oneAccModeInd3Word.ind, oneAccModeInd3Word.mode, oneAccModeInd3Word.disp31, iPtr.dispOffset))

	case instrLLEFB:
		oneAccMode3Word := iPtr.variant.(oneAccMode3WordT)
		addr := resolve32bitEffAddr(cpuPtr, ' ', oneAccMode3Word.mode, oneAccMode3Word.disp31>>1, iPtr.dispOffset)
		addr <<= 1
		if memory.TestDwbit(dg.DwordT(oneAccMode3Word.disp31), 31) {
			addr |= 1
		}
		cpuPtr.ac[oneAccMode3Word.acd] = dg.DwordT(addr)

	case instrLPSR:
		cpuPtr.ac[0] = dg.DwordT(cpuPtr.psr)

	case instrNADD: // signed add
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		acDs16 := int16(cpuPtr.ac[twoAcc1Word.acd])
		acSs16 := int16(cpuPtr.ac[twoAcc1Word.acs])
		s32 := int32(acDs16) + int32(acSs16)
		cpuPtr.carry = (s32 > maxPosS16) || (s32 < minNegS16) // TODO handle overflow flag
		cpuPtr.ac[twoAcc1Word.acd] = dg.DwordT(s32)

	case instrNADDI:
		oneAccImm2Word := iPtr.variant.(oneAccImm2WordT)
		acDs16 := int16(cpuPtr.ac[oneAccImm2Word.acd])
		s16 := oneAccImm2Word.immS16
		s32 := int32(acDs16) + int32(s16)
		cpuPtr.carry = (s32 > maxPosS16) || (s32 < minNegS16) // TODO handle overflow flag
		cpuPtr.ac[oneAccImm2Word.acd] = dg.DwordT(s32)

	case instrNLDAI:
		oneAccImm2Word := iPtr.variant.(oneAccImm2WordT)
		cpuPtr.ac[oneAccImm2Word.acd] = dg.DwordT(int32(oneAccImm2Word.immS16))

	case instrNSUB: // signed subtract
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		acDs16 := int16(cpuPtr.ac[twoAcc1Word.acd])
		acSs16 := int16(cpuPtr.ac[twoAcc1Word.acs])
		s32 := int32(acDs16) - int32(acSs16)
		cpuPtr.carry = (s32 > maxPosS16) || (s32 < minNegS16) // TODO handle overflow flag
		cpuPtr.ac[twoAcc1Word.acd] = dg.DwordT(s32)

	case instrSEX: // Sign EXtend
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		cpuPtr.ac[twoAcc1Word.acd] = memory.SexWordToDword(memory.DwordGetLowerWord(cpuPtr.ac[twoAcc1Word.acs]))

	case instrSSPT: /* NO-OP - see p.8-5 of MV/10000 Sys Func Chars */
		log.Println("INFO: SSPT is a No-Op on this machine, continuing")

	case instrWADC:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		acDs32 := int32(cpuPtr.ac[twoAcc1Word.acd])
		acSs32 := int32(^cpuPtr.ac[twoAcc1Word.acs])
		s64 := int64(acSs32) + int64(acDs32)
		cpuPtr.carry = (s64 > maxPosS32) || (s64 < minNegS32) // TODO handle overflow flag
		cpuPtr.ac[twoAcc1Word.acd] = dg.DwordT(s64)

	case instrWADD:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		acDs32 := int32(cpuPtr.ac[twoAcc1Word.acd])
		acSs32 := int32(cpuPtr.ac[twoAcc1Word.acs])
		s64 := int64(acSs32) + int64(acDs32)
		cpuPtr.carry = (s64 > maxPosS32) || (s64 < minNegS32) // TODO handle overflow flag
		cpuPtr.ac[twoAcc1Word.acd] = dg.DwordT(s64)

	case instrWADDI:
		oneAccImm3Word := iPtr.variant.(oneAccImm3WordT)
		acDs32 := int32(cpuPtr.ac[oneAccImm3Word.acd])
		s32i := int32(oneAccImm3Word.immU32)
		s64 := int64(s32i) + int64(acDs32)
		cpuPtr.carry = (s64 > maxPosS32) || (s64 < minNegS32) // TODO handle overflow flag
		cpuPtr.ac[oneAccImm3Word.acd] = dg.DwordT(s64)

	case instrWADI:
		immOneAcc := iPtr.variant.(immOneAccT)
		acDs32 := int32(cpuPtr.ac[immOneAcc.acd])
		s32 := int32(immOneAcc.immU16)
		s64 := int64(s32) + int64(acDs32)
		cpuPtr.carry = (s64 > maxPosS32) || (s64 < minNegS32) // TODO handle overflow flag
		cpuPtr.ac[immOneAcc.acd] = dg.DwordT(s64)

	case instrWAND:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		cpuPtr.ac[twoAcc1Word.acd] &= cpuPtr.ac[twoAcc1Word.acs]

	case instrWANDI:
		oneAccImmDwd3Word := iPtr.variant.(oneAccImmDwd3WordT)
		cpuPtr.ac[oneAccImmDwd3Word.acd] &= oneAccImmDwd3Word.immDword

	case instrWCOM:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		cpuPtr.ac[twoAcc1Word.acd] = ^cpuPtr.ac[twoAcc1Word.acs]

	case instrWDIVS:
		s64 := int64(memory.QwordFromTwoDwords(cpuPtr.ac[0], cpuPtr.ac[1]))
		if cpuPtr.ac[2] == 0 {
			cpuPtr.CPUSetOVR(true)
		} else {
			s32 := int32(cpuPtr.ac[2])
			if s64/int64(s32) < -2147483648 || s64/int64(s32) > 2147483647 {
				cpuPtr.CPUSetOVR(true)
			} else {
				cpuPtr.ac[0] = dg.DwordT(s64 % int64(s32))
				cpuPtr.ac[1] = dg.DwordT(s64 / int64(s32))
			}
		}

	case instrWINC:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		cpuPtr.carry = cpuPtr.ac[twoAcc1Word.acs] == 0xffffffff // TODO handle overflow flag
		cpuPtr.ac[twoAcc1Word.acd] = cpuPtr.ac[twoAcc1Word.acs] + 1

	case instrWIOR:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		cpuPtr.ac[twoAcc1Word.acd] |= cpuPtr.ac[twoAcc1Word.acs]

	case instrWIORI:
		oneAccImmDwd3Word := iPtr.variant.(oneAccImmDwd3WordT)
		cpuPtr.ac[oneAccImmDwd3Word.acd] |= oneAccImmDwd3Word.immDword

	case instrWLDAI:
		oneAccImmDwd3Word := iPtr.variant.(oneAccImmDwd3WordT)
		cpuPtr.ac[oneAccImmDwd3Word.acd] = oneAccImmDwd3Word.immDword

	case instrWLSH:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		shiftAmt8 := int8(cpuPtr.ac[twoAcc1Word.acs] & 0x0ff)
		switch { // do nothing if shift of zero was specified
		case shiftAmt8 < 0: // shift right
			shiftAmt8 *= -1
			cpuPtr.ac[twoAcc1Word.acd] >>= uint(shiftAmt8)
		case shiftAmt8 > 0: // shift left
			cpuPtr.ac[twoAcc1Word.acd] <<= uint(shiftAmt8)
		}

	case instrWLSHI:
		oneAccImm2Word := iPtr.variant.(oneAccImm2WordT)
		shiftAmt8 := int8(oneAccImm2Word.immS16 & 0x0ff)
		switch { // do nothing if shift of zero was specified
		case shiftAmt8 < 0: // shift right
			shiftAmt8 *= -1
			cpuPtr.ac[oneAccImm2Word.acd] >>= uint(shiftAmt8)
		case shiftAmt8 > 0: // shift left
			cpuPtr.ac[oneAccImm2Word.acd] <<= uint(shiftAmt8)
		}

	case instrWLSI:
		immOneAcc := iPtr.variant.(immOneAccT)
		cpuPtr.ac[immOneAcc.acd] = cpuPtr.ac[immOneAcc.acd] << immOneAcc.immU16

	case instrWMOV:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		cpuPtr.ac[twoAcc1Word.acd] = cpuPtr.ac[twoAcc1Word.acs]

	case instrWMOVR:
		oneAcc1Word := iPtr.variant.(oneAcc1WordT)
		cpuPtr.ac[oneAcc1Word.acd] = cpuPtr.ac[oneAcc1Word.acd] >> 1

	case instrWMUL:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		res := int32(cpuPtr.ac[twoAcc1Word.acd]) * int32(cpuPtr.ac[twoAcc1Word.acs])
		// FIXME - handle overflow
		cpuPtr.ac[twoAcc1Word.acd] = dg.DwordT(res)

	case instrWMULS:
		s64 := int64(cpuPtr.ac[1])*int64(cpuPtr.ac[2]) + int64(cpuPtr.ac[0])
		cpuPtr.ac[0] = memory.QwordGetUpperDword(dg.QwordT(s64))
		cpuPtr.ac[1] = memory.QwordGetLowerDword(dg.QwordT(s64))

	case instrWNADI: //signed 16-bit
		oneAccImm2Word := iPtr.variant.(oneAccImm2WordT)
		acDs32 := int32(cpuPtr.ac[oneAccImm2Word.acd])
		s16 := oneAccImm2Word.immS16
		s64 := int64(acDs32) + int64(s16)
		cpuPtr.carry = (s64 > maxPosS32) || (s64 < minNegS32) // TODO handle overflow flag
		cpuPtr.ac[oneAccImm2Word.acd] = dg.DwordT(s64)

	case instrWNEG:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		cpuPtr.carry = int32(cpuPtr.ac[twoAcc1Word.acs]) == minNegS32 // TODO handle overflow flag
		cpuPtr.ac[twoAcc1Word.acd] = (^cpuPtr.ac[twoAcc1Word.acs]) + 1

	case instrWSBI:
		immOneAcc := iPtr.variant.(immOneAccT)
		acDs32 := int32(cpuPtr.ac[immOneAcc.acd])
		s32 := int32(immOneAcc.immU16)
		s64 := int64(acDs32) - int64(s32)
		cpuPtr.carry = (s64 > maxPosS32) || (s64 < minNegS32) // TODO handle overflow flag
		cpuPtr.ac[immOneAcc.acd] = dg.DwordT(s64)

	case instrWSUB:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		acDs32 := int32(cpuPtr.ac[twoAcc1Word.acd])
		acSs32 := int32(cpuPtr.ac[twoAcc1Word.acs])
		s64 := int64(acDs32) - int64(acSs32)
		cpuPtr.carry = (s64 > maxPosS32) || (s64 < minNegS32) // TODO handle overflow flag
		cpuPtr.ac[twoAcc1Word.acd] = dg.DwordT(s64)

	case instrWXCH:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		dwd := cpuPtr.ac[twoAcc1Word.acs]
		cpuPtr.ac[twoAcc1Word.acs] = cpuPtr.ac[twoAcc1Word.acd]
		cpuPtr.ac[twoAcc1Word.acd] = dwd

	case instrZEX:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		cpuPtr.ac[twoAcc1Word.acd] = 0 | dg.DwordT(memory.DwordGetLowerWord(cpuPtr.ac[twoAcc1Word.acs]))

	default:
		log.Fatalf("ERROR: EAGLE_OP instruction <%s> not yet implemented\n", iPtr.mnemonic)
		return false
	}

	cpuPtr.pc += dg.PhysAddrT(iPtr.instrLength)
	return true
}
