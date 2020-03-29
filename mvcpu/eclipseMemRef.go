// eclipseMemRef.go

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
	"fmt"
	"log"

	"github.com/SMerrony/dgemug/dg"
	"github.com/SMerrony/dgemug/logging"
	"github.com/SMerrony/dgemug/memory"
)

func eclipseMemRef(cpu *CPUT, iPtr *decodedInstrT) bool {

	switch iPtr.ix {

	case instrBLM:
		/* AC0 - unused, AC1 - no. wds to move, AC2 - src, AC3 - dest */
		numWds := memory.DwordGetLowerWord(cpu.ac[1])
		if numWds == 0 || numWds > 32768 {
			if cpu.debugLogging {
				logging.DebugPrint(logging.DebugLog, "BLM called with AC1 out-of-bounds, not moving anything\n")
			}
			break
		}
		src := (cpu.pc & ringMask32) | dg.PhysAddrT(memory.DwordGetLowerWord(cpu.ac[2]))
		dest := (cpu.pc & ringMask32) | dg.PhysAddrT(memory.DwordGetLowerWord(cpu.ac[3]))
		if cpu.debugLogging {
			logging.DebugPrint(logging.DebugLog, fmt.Sprintf("BLM moving %d words from %d to %d\n", numWds, src, dest))
		}
		for numWds != 0 {
			memory.WriteWord(dest, memory.ReadWord(src))
			numWds--
			src++
			dest++
		}
		cpu.ac[1] = 0
		cpu.ac[2] = dg.DwordT(src) // TODO confirm this is right, doc ambiguous
		cpu.ac[3] = dg.DwordT(dest)

	case instrBTO:
		// TODO Handle segment and indirection...
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		addr, bitNum := resolveEclipseBitAddr(cpu, &twoAcc1Word)
		addr |= (cpu.pc & ringMask32)
		wd := memory.ReadWord(addr)
		if cpu.debugLogging {
			logging.DebugPrint(logging.DebugLog, "... BTO Addr: %d, Bit: %d, Before: %s\n",
				addr, bitNum, memory.WordToBinStr(wd))
		}
		memory.SetWbit(&wd, bitNum)
		memory.WriteWord(addr, wd)
		if cpu.debugLogging {
			logging.DebugPrint(logging.DebugLog, "... BTO                     Result: %s\n", memory.WordToBinStr(wd))
		}

	case instrBTZ:
		// TODO Handle segment and indirection...
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		addr, bitNum := resolveEclipseBitAddr(cpu, &twoAcc1Word)
		addr |= (cpu.pc & ringMask32)
		wd := memory.ReadWord(addr)
		if cpu.debugLogging {
			logging.DebugPrint(logging.DebugLog, "... BTZ Addr: %d, Bit: %d, Before: %s\n", addr, bitNum, memory.WordToBinStr(wd))
		}
		memory.ClearWbit(&wd, bitNum)
		memory.WriteWord(addr, wd)
		if cpu.debugLogging {
			logging.DebugPrint(logging.DebugLog, "... BTZ                     Result: %s\n",
				memory.WordToBinStr(wd))
		}
	case instrCMP:
		cmp(cpu)

	case instrCMV:
		cmv(cpu)

	case instrELDA:
		oneAccModeInt2Word := iPtr.variant.(oneAccModeInd2WordT)
		addr := resolve15bitDisplacement(cpu, oneAccModeInt2Word.ind, oneAccModeInt2Word.mode, dg.WordT(oneAccModeInt2Word.disp15), iPtr.dispOffset)
		addr &= 0x7fff
		addr |= (cpu.pc & ringMask32)
		cpu.ac[oneAccModeInt2Word.acd] = dg.DwordT(memory.ReadWord(addr))

	case instrELEF:
		oneAccModeInt2Word := iPtr.variant.(oneAccModeInd2WordT)
		addr := resolve15bitDisplacement(cpu, oneAccModeInt2Word.ind, oneAccModeInt2Word.mode, dg.WordT(oneAccModeInt2Word.disp15), iPtr.dispOffset)
		addr &= 0x7fff
		addr |= (cpu.pc & ringMask32)
		cpu.ac[oneAccModeInt2Word.acd] = dg.DwordT(addr)

	case instrESTA:
		oneAccModeInt2Word := iPtr.variant.(oneAccModeInd2WordT)
		// addr := resolve16bitEffAddr(cpu, oneAccModeInt2Word.ind, oneAccModeInt2Word.mode, oneAccModeInt2Word.disp15, iPtr.dispOffset)
		addr := resolve15bitDisplacement(cpu, oneAccModeInt2Word.ind, oneAccModeInt2Word.mode, dg.WordT(oneAccModeInt2Word.disp15), iPtr.dispOffset)
		addr &= 0x7fff
		addr |= (cpu.pc & ringMask32)
		memory.WriteWord(addr, memory.DwordGetLowerWord(cpu.ac[oneAccModeInt2Word.acd]))

	case instrLDB:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		cpu.ac[twoAcc1Word.acd] = dg.DwordT(memory.ReadByteEclipseBA(cpu.pc, memory.DwordGetLowerWord(cpu.ac[twoAcc1Word.acs])))

	case instrLEF:
		novaOneAccEffAddr := iPtr.variant.(novaOneAccEffAddrT)
		addr := resolve8bitDisplacement(cpu, novaOneAccEffAddr.ind, novaOneAccEffAddr.mode, novaOneAccEffAddr.disp15)
		addr &= 0x7fff
		addr |= (cpu.pc & ringMask32)
		cpu.ac[novaOneAccEffAddr.acd] = dg.DwordT(addr)

	case instrSTB:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		hiLo := memory.TestDwbit(cpu.ac[twoAcc1Word.acs], 31)
		addr := dg.PhysAddrT(memory.DwordGetLowerWord(cpu.ac[twoAcc1Word.acs])) >> 1
		addr &= 0x7fff
		addr |= (cpu.pc & ringMask32)
		byt := dg.ByteT(cpu.ac[twoAcc1Word.acd])
		memory.WriteByte(addr, hiLo, byt)

	default:
		log.Printf("ERROR: ECLIPSE_MEMREF instruction <%s> not yet implemented\n", iPtr.mnemonic)
		return false
	}

	cpu.pc += dg.PhysAddrT(iPtr.instrLength)
	return true
}

func cmp(cpu *CPUT) {
	var str1len, str2len int16
	str2len = int16(memory.DwordGetLowerWord(cpu.ac[0]))
	str1len = int16(memory.DwordGetLowerWord(cpu.ac[1]))
	if str1len == 0 && str2len == 0 {
		cpu.ac[1] = 0
		return
	}
	str1bp := memory.DwordGetLowerWord(cpu.ac[3])
	str2bp := memory.DwordGetLowerWord(cpu.ac[2])
	var byte1, byte2 dg.ByteT
	res := 0
	for {
		if str1len != 0 {
			byte1 = memory.ReadByteEclipseBA(cpu.pc, str1bp)
		} else {
			byte1 = ' '
		}
		if str2len != 0 {
			byte2 = memory.ReadByteEclipseBA(cpu.pc, str2bp)
		} else {
			byte2 = ' '
		}
		if byte1 > byte2 {
			res = 1
			break
		}
		if byte1 < byte2 {
			res = -1
			break
		}
		if str1len > 0 {
			str1bp++
			str1len--
		}
		if str1len < 0 {
			str1bp--
			str1len++
		}
		if str2len > 0 {
			str2bp++
			str2len--
		}
		if str2len < 0 {
			str2bp--
			str2len++
		}
		if str1len == 0 && str2len == 0 {
			break
		}
	}
	cpu.ac[0] = dg.DwordT(str2len)
	cpu.ac[1] = dg.DwordT(res)
	cpu.ac[2] = dg.DwordT(str2bp)
	cpu.ac[3] = dg.DwordT(str1bp)
}

func cmv(cpu *CPUT) {
	// ACO destCount, AC1 srcCount, AC2 dest byte ptr, AC3 src byte ptr
	var destAscend, srcAscend bool
	destCount := int16(memory.DwordGetLowerWord(cpu.ac[0]))
	if destCount == 0 {
		log.Println("INFO: CMV called with AC0 == 0, not moving anything")
		cpu.carry = false
		return
	}
	destAscend = (destCount > 0)
	srcCount := int16(memory.DwordGetLowerWord(cpu.ac[3]))
	srcAscend = (srcCount > 0)
	if cpu.debugLogging {
		logging.DebugPrint(logging.DebugLog, "DEBUG: CMV moving %d chars from %d to %d\n",
			srcCount, cpu.ac[3], cpu.ac[2])
	}
	// set carry if length of src is greater than length of dest
	if cpu.ac[1] > cpu.ac[2] {
		cpu.carry = true
	}
	// 1st move srcCount bytes
	for {
		copyByte(cpu.ac[3], cpu.ac[2])
		if srcAscend {
			cpu.ac[3]++
			srcCount--
		} else {
			cpu.ac[3]--
			srcCount++
		}
		if destAscend {
			cpu.ac[2]++
			destCount--
		} else {
			cpu.ac[2]--
			destCount++
		}
		if srcCount == 0 || destCount == 0 {
			break
		}
	}
	// now fill any excess bytes with ASCII spaces
	if destCount != 0 {
		for {
			memWriteByteBA(dg.ASCIISPC, cpu.ac[2])
			if destAscend {
				cpu.ac[2]++
				destCount--
			} else {
				cpu.ac[2]--
				destCount++
			}
			if destCount == 0 {
				break
			}
		}
	}
	cpu.ac[0] = 0
	cpu.ac[1] = dg.DwordT(srcCount)
}
