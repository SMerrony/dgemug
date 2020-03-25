// mvemg project eclipseOp.go

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

func eclipseOp(cpu *CPUT, iPtr *decodedInstrT) bool {

	switch iPtr.ix {

	case instrADDI:
		oneAccImm2Word := iPtr.variant.(oneAccImm2WordT)
		// signed 16-bit add immediate
		s16 := int16(memory.DwordGetLowerWord(cpu.ac[oneAccImm2Word.acd]))
		s16 += oneAccImm2Word.immS16
		cpu.ac[oneAccImm2Word.acd] = dg.DwordT(s16) & 0X0000FFFF

	case instrANDI:
		oneAccImmWd2Word := iPtr.variant.(oneAccImmWd2WordT)
		wd := memory.DwordGetLowerWord(cpu.ac[oneAccImmWd2Word.acd])
		cpu.ac[oneAccImmWd2Word.acd] = dg.DwordT(wd&oneAccImmWd2Word.immWord) & 0x0000ffff

	case instrADI: // 16-bit unsigned Add Immediate
		immOneAcc := iPtr.variant.(immOneAccT)
		wd := memory.DwordGetLowerWord(cpu.ac[immOneAcc.acd])
		wd += dg.WordT(immOneAcc.immU16) // unsigned arithmetic does wraparound in Go
		cpu.ac[immOneAcc.acd] = dg.DwordT(wd)

	case instrDHXL:
		immOneAcc := iPtr.variant.(immOneAccT)
		dplus1 := immOneAcc.acd + 1
		if dplus1 == 4 {
			dplus1 = 0
		}
		dwd := memory.DwordFromTwoWords(memory.DwordGetLowerWord(cpu.ac[immOneAcc.acd]), memory.DwordGetLowerWord(cpu.ac[dplus1]))
		dwd <<= (immOneAcc.immU16 * 4)
		cpu.ac[immOneAcc.acd] = dg.DwordT(memory.DwordGetUpperWord(dwd))
		cpu.ac[dplus1] = dg.DwordT(memory.DwordGetLowerWord(dwd))

	case instrDLSH:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		dplus1 := twoAcc1Word.acd + 1
		if dplus1 == 4 {
			dplus1 = 0
		}
		dwd := dlsh(cpu.ac[twoAcc1Word.acs], cpu.ac[twoAcc1Word.acd], cpu.ac[dplus1])
		cpu.ac[twoAcc1Word.acd] = dg.DwordT(memory.DwordGetUpperWord(dwd))
		cpu.ac[dplus1] = dg.DwordT(memory.DwordGetLowerWord(dwd))

	case instrHXL:
		immOneAcc := iPtr.variant.(immOneAccT)
		dwd := cpu.ac[immOneAcc.acd] << (uint32(immOneAcc.immU16) * 4)
		cpu.ac[immOneAcc.acd] = dwd & 0x0ffff

	case instrHXR:
		immOneAcc := iPtr.variant.(immOneAccT)
		dwd := cpu.ac[immOneAcc.acd] >> (uint32(immOneAcc.immU16) * 4)
		cpu.ac[immOneAcc.acd] = dwd & 0x0ffff

	case instrIOR:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		wd := memory.DwordGetLowerWord(cpu.ac[twoAcc1Word.acd]) | memory.DwordGetLowerWord(cpu.ac[twoAcc1Word.acs])
		cpu.ac[twoAcc1Word.acd] = dg.DwordT(wd)

	case instrIORI:
		oneAccImmWd2Word := iPtr.variant.(oneAccImmWd2WordT)
		wd := memory.DwordGetLowerWord(cpu.ac[oneAccImmWd2Word.acd]) | oneAccImmWd2Word.immWord
		cpu.ac[oneAccImmWd2Word.acd] = dg.DwordT(wd)

	case instrLSH:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		cpu.ac[twoAcc1Word.acd] = lsh(cpu.ac[twoAcc1Word.acs], cpu.ac[twoAcc1Word.acd])

	case instrSBI: // unsigned
		immOneAcc := iPtr.variant.(immOneAccT)
		wd := memory.DwordGetLowerWord(cpu.ac[immOneAcc.acd])
		if immOneAcc.immU16 < 1 || immOneAcc.immU16 > 4 {
			log.Fatal("Invalid immediate value in SBI")
		}
		wd -= dg.WordT(immOneAcc.immU16)
		cpu.ac[immOneAcc.acd] = dg.DwordT(wd)

	case instrXCH:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		dwd := cpu.ac[twoAcc1Word.acs]
		cpu.ac[twoAcc1Word.acs] = cpu.ac[twoAcc1Word.acd] & 0x0ffff
		cpu.ac[twoAcc1Word.acd] = dwd & 0x0ffff

	default:
		log.Fatalf("ERROR: ECLIPSE_OP instruction <%s> not yet implemented\n", iPtr.mnemonic)
		return false
	}

	cpu.pc += dg.PhysAddrT(iPtr.instrLength)
	return true
}

func dlsh(acS, acDh, acDl dg.DwordT) dg.DwordT {
	var shft = int8(acS)
	var dwd = memory.DwordFromTwoWords(memory.DwordGetLowerWord(acDh), memory.DwordGetLowerWord(acDl))
	if shft != 0 {
		if shft < -31 || shft > 31 {
			dwd = 0
		} else {
			if shft > 0 {
				dwd >>= uint8(shft)
			} else {
				shft *= -1
				dwd >>= uint8(shft)
			}
		}
	}
	return dwd
}

func lsh(acS, acD dg.DwordT) dg.DwordT {
	var shft = int8(acS)
	var wd = memory.DwordGetLowerWord(acD)
	if shft == 0 {
		wd = memory.DwordGetLowerWord(acD) // do nothing
	} else {
		if shft < -15 || shft > 15 {
			wd = 0 // 16+ bit shift clears word
		} else {
			if shft > 0 {
				wd >>= uint8(shft)
			} else {
				shft *= -1
				wd >>= uint8(shft)
			}
		}
	}
	return dg.DwordT(wd)
}
