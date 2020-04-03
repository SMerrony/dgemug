// eagleOp.go

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
	"github.com/SMerrony/dgemug/memory"
)

func eagleOp(cpu *CPUT, iPtr *decodedInstrT) bool {

	switch iPtr.ix {

	case instrCRYTC:
		cpu.carry = !cpu.carry

	case instrCRYTO:
		cpu.carry = true

	case instrCRYTZ:
		cpu.carry = false

	case instrCVWN:
		dwd := cpu.ac[iPtr.ac]
		if dwd>>16 != 0 && dwd>>16 != 0xffff {
			cpu.SetOVR(true)
		}
		if memory.TestDwbit(dwd, 16) {
			cpu.ac[iPtr.ac] |= 0xffff0000
		} else {
			cpu.ac[iPtr.ac] &= 0x0000ffff
		}

	case instrLPSR:
		cpu.ac[0] = dg.DwordT(cpu.psr)

	case instrNADD: // signed add
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		acDs16 := int16(cpu.ac[twoAcc1Word.acd])
		acSs16 := int16(cpu.ac[twoAcc1Word.acs])
		s32 := int32(acDs16) + int32(acSs16)
		cpu.carry = (s32 > maxPosS16) || (s32 < minNegS16) // TODO handle overflow flag
		cpu.ac[twoAcc1Word.acd] = dg.DwordT(s32)

	case instrNADDI:
		oneAccImm2Word := iPtr.variant.(oneAccImm2WordT)
		acDs16 := int16(cpu.ac[oneAccImm2Word.acd])
		s16 := oneAccImm2Word.immS16
		s32 := int32(acDs16) + int32(s16)
		cpu.carry = (s32 > maxPosS16) || (s32 < minNegS16) // TODO handle overflow flag
		cpu.ac[oneAccImm2Word.acd] = dg.DwordT(s32)

	case instrNADI, instrNSBI:
		immOneAcc := iPtr.variant.(immOneAccT)
		s32 := int32(int16(cpu.ac[immOneAcc.acd] & 0xffff))
		switch iPtr.ix {
		case instrNADI:
			cpu.ac[immOneAcc.acd] = dg.DwordT(s32 + int32(immOneAcc.immU16))
		case instrNSBI:
			cpu.ac[immOneAcc.acd] = dg.DwordT(s32 - int32(immOneAcc.immU16))
		}
		if s32 < minNegS16 || s32 > maxPosS16 {
			cpu.carry = true
			cpu.SetOVR(true)
		}

	case instrNLDAI:
		oneAccImm2Word := iPtr.variant.(oneAccImm2WordT)
		cpu.ac[oneAccImm2Word.acd] = dg.DwordT(int32(oneAccImm2Word.immS16))

	case instrNSUB: // signed subtract
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		acDs16 := int16(cpu.ac[twoAcc1Word.acd])
		acSs16 := int16(cpu.ac[twoAcc1Word.acs])
		s32 := int32(acDs16) - int32(acSs16)
		cpu.carry = (s32 > maxPosS16) || (s32 < minNegS16) // TODO handle overflow flag
		cpu.ac[twoAcc1Word.acd] = dg.DwordT(s32)

	case instrSEX: // Sign EXtend
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		cpu.ac[twoAcc1Word.acd] = memory.SexWordToDword(memory.DwordGetLowerWord(cpu.ac[twoAcc1Word.acs]))

	case instrSSPT: /* NO-OP - see p.8-5 of MV/10000 Sys Func Chars */
		log.Println("INFO: SSPT is a No-Op on this machine, continuing")

	case instrWADC:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		acDs32 := int32(cpu.ac[twoAcc1Word.acd])
		acSs32 := int32(^cpu.ac[twoAcc1Word.acs])
		s64 := int64(acSs32) + int64(acDs32)
		cpu.carry = (s64 > maxPosS32) || (s64 < minNegS32) // TODO handle overflow flag
		cpu.ac[twoAcc1Word.acd] = dg.DwordT(s64)

	case instrWADD:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		acDs32 := int32(cpu.ac[twoAcc1Word.acd])
		acSs32 := int32(cpu.ac[twoAcc1Word.acs])
		s64 := int64(acSs32) + int64(acDs32)
		cpu.carry = (s64 > maxPosS32) || (s64 < minNegS32) // TODO handle overflow flag
		cpu.ac[twoAcc1Word.acd] = dg.DwordT(s64)

	case instrWADDI:
		oneAccImm3Word := iPtr.variant.(oneAccImm3WordT)
		acDs32 := int32(cpu.ac[oneAccImm3Word.acd])
		s32i := int32(oneAccImm3Word.immU32)
		s64 := int64(s32i) + int64(acDs32)
		cpu.carry = (s64 > maxPosS32) || (s64 < minNegS32) // TODO handle overflow flag
		cpu.ac[oneAccImm3Word.acd] = dg.DwordT(s64)

	case instrWADI:
		immOneAcc := iPtr.variant.(immOneAccT)
		acDs32 := int32(cpu.ac[immOneAcc.acd])
		s32 := int32(immOneAcc.immU16)
		s64 := int64(s32) + int64(acDs32)
		cpu.carry = (s64 > maxPosS32) || (s64 < minNegS32) // TODO handle overflow flag
		cpu.ac[immOneAcc.acd] = dg.DwordT(s64)

	case instrWANC:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		cpu.ac[twoAcc1Word.acd] &= ^cpu.ac[twoAcc1Word.acs]

	case instrWAND:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		cpu.ac[twoAcc1Word.acd] &= cpu.ac[twoAcc1Word.acs]

	case instrWANDI:
		oneAccImmDwd3Word := iPtr.variant.(oneAccImmDwd3WordT)
		cpu.ac[oneAccImmDwd3Word.acd] &= oneAccImmDwd3Word.immDword

	case instrWASHI:
		oneAccImm2Word := iPtr.variant.(oneAccImm2WordT)
		s32 := int32(cpu.ac[oneAccImm2Word.acd])
		switch {
		case oneAccImm2Word.immS16 == 0: // nothing happens
		case oneAccImm2Word.immS16 < 0: // shift right
			s32 >>= (oneAccImm2Word.immS16 * -1)
			cpu.ac[oneAccImm2Word.acd] = dg.DwordT(s32)
		case oneAccImm2Word.immS16 > 0: // shift left
			s32 <<= oneAccImm2Word.immS16
			cpu.ac[oneAccImm2Word.acd] = dg.DwordT(s32)
		}

	case instrWCOM:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		cpu.ac[twoAcc1Word.acd] = ^cpu.ac[twoAcc1Word.acs]

	case instrWDIV:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		if cpu.ac[twoAcc1Word.acs] != 0 {
			s64 := int64(int32(cpu.ac[twoAcc1Word.acd])) / int64(int32(cpu.ac[twoAcc1Word.acs]))
			if s64 <= maxPosS32 && s64 >= minNegS32 {
				cpu.ac[twoAcc1Word.acd] = dg.DwordT(s64)
				cpu.SetOVR(false)
			} else {
				cpu.SetOVR(true)
			}
		} else {
			cpu.SetOVR(true)
		}

	case instrWDIVS:
		s64 := int64(memory.QwordFromTwoDwords(cpu.ac[0], cpu.ac[1]))
		if cpu.ac[2] == 0 {
			cpu.SetOVR(true)
		} else {
			s32 := int32(cpu.ac[2])
			if s64/int64(s32) < -2147483648 || s64/int64(s32) > 2147483647 {
				cpu.SetOVR(true)
			} else {
				cpu.ac[0] = dg.DwordT(s64 % int64(s32))
				cpu.ac[1] = dg.DwordT(s64 / int64(s32))
			}
		}

	case instrWHLV:
		s32 := int32(cpu.ac[iPtr.ac]) / 2
		cpu.ac[iPtr.ac] = dg.DwordT(s32)

	case instrWINC:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		cpu.carry = cpu.ac[twoAcc1Word.acs] == 0xffffffff // TODO handle overflow flag
		cpu.ac[twoAcc1Word.acd] = cpu.ac[twoAcc1Word.acs] + 1

	case instrWIOR:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		cpu.ac[twoAcc1Word.acd] |= cpu.ac[twoAcc1Word.acs]

	case instrWIORI:
		oneAccImmDwd3Word := iPtr.variant.(oneAccImmDwd3WordT)
		cpu.ac[oneAccImmDwd3Word.acd] |= oneAccImmDwd3Word.immDword

	case instrWLDAI:
		oneAccImmDwd3Word := iPtr.variant.(oneAccImmDwd3WordT)
		cpu.ac[oneAccImmDwd3Word.acd] = oneAccImmDwd3Word.immDword

	case instrWLSH:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		shiftAmt8 := int8(cpu.ac[twoAcc1Word.acs] & 0x0ff)
		switch { // do nothing if shift of zero was specified
		case shiftAmt8 < 0: // shift right
			shiftAmt8 *= -1
			cpu.ac[twoAcc1Word.acd] >>= uint(shiftAmt8)
		case shiftAmt8 > 0: // shift left
			cpu.ac[twoAcc1Word.acd] <<= uint(shiftAmt8)
		}

	case instrWLSHI:
		oneAccImm2Word := iPtr.variant.(oneAccImm2WordT)
		shiftAmt8 := int8(oneAccImm2Word.immS16 & 0x0ff)
		switch { // do nothing if shift of zero was specified
		case shiftAmt8 < 0: // shift right
			shiftAmt8 *= -1
			cpu.ac[oneAccImm2Word.acd] >>= uint(shiftAmt8)
		case shiftAmt8 > 0: // shift left
			cpu.ac[oneAccImm2Word.acd] <<= uint(shiftAmt8)
		}

	case instrWLSI:
		immOneAcc := iPtr.variant.(immOneAccT)
		cpu.ac[immOneAcc.acd] = cpu.ac[immOneAcc.acd] << immOneAcc.immU16

	case instrWMOV:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		cpu.ac[twoAcc1Word.acd] = cpu.ac[twoAcc1Word.acs]

	case instrWMOVR:
		cpu.ac[iPtr.ac] >>= 1

	case instrWMUL:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		res := int32(cpu.ac[twoAcc1Word.acd]) * int32(cpu.ac[twoAcc1Word.acs])
		// FIXME - handle overflow
		cpu.ac[twoAcc1Word.acd] = dg.DwordT(res)

	case instrWMULS:
		s64 := int64(cpu.ac[1])*int64(cpu.ac[2]) + int64(cpu.ac[0])
		cpu.ac[0] = memory.QwordGetUpperDword(dg.QwordT(s64))
		cpu.ac[1] = memory.QwordGetLowerDword(dg.QwordT(s64))

	case instrWNADI: //signed 16-bit
		oneAccImm2Word := iPtr.variant.(oneAccImm2WordT)
		acDs32 := int32(cpu.ac[oneAccImm2Word.acd])
		s16 := oneAccImm2Word.immS16
		s64 := int64(acDs32) + int64(s16)
		cpu.carry = (s64 > maxPosS32) || (s64 < minNegS32) // TODO handle overflow flag
		cpu.ac[oneAccImm2Word.acd] = dg.DwordT(s64)

	case instrWNEG:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		cpu.carry = int32(cpu.ac[twoAcc1Word.acs]) == minNegS32 // TODO handle overflow flag
		cpu.ac[twoAcc1Word.acd] = (^cpu.ac[twoAcc1Word.acs]) + 1

	case instrWSBI:
		immOneAcc := iPtr.variant.(immOneAccT)
		acDs32 := int32(cpu.ac[immOneAcc.acd])
		s32 := int32(immOneAcc.immU16)
		s64 := int64(acDs32) - int64(s32)
		cpu.carry = (s64 > maxPosS32) || (s64 < minNegS32) // TODO handle overflow flag
		cpu.ac[immOneAcc.acd] = dg.DwordT(s64)

	case instrWSUB:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		acDs32 := int32(cpu.ac[twoAcc1Word.acd])
		acSs32 := int32(cpu.ac[twoAcc1Word.acs])
		s64 := int64(acDs32) - int64(acSs32)
		cpu.carry = (s64 > maxPosS32) || (s64 < minNegS32) // TODO handle overflow flag
		cpu.ac[twoAcc1Word.acd] = dg.DwordT(s64)

	case instrWXCH:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		dwd := cpu.ac[twoAcc1Word.acs]
		cpu.ac[twoAcc1Word.acs] = cpu.ac[twoAcc1Word.acd]
		cpu.ac[twoAcc1Word.acd] = dwd

	case instrWXORI:
		oneAccImm3Word := iPtr.variant.(oneAccImm3WordT)
		cpu.ac[oneAccImm3Word.acd] ^= dg.DwordT(oneAccImm3Word.immU32)

	case instrZEX:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		cpu.ac[twoAcc1Word.acd] = 0 | dg.DwordT(memory.DwordGetLowerWord(cpu.ac[twoAcc1Word.acs]))

	default:
		log.Fatalf("ERROR: EAGLE_OP instruction <%s> not yet implemented\n", iPtr.mnemonic)
		return false
	}

	cpu.pc += dg.PhysAddrT(iPtr.instrLength)
	return true
}
