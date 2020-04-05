// eagleMemRef.go

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

func eagleMemRef(cpu *CPUT, iPtr *decodedInstrT) bool {

	switch iPtr.ix {

	case instrLLDB:
		oneAccMode3Word := iPtr.variant.(oneAccMode3WordT)
		addr := resolve32bitEffAddr(cpu, ' ', oneAccMode3Word.mode, int32(oneAccMode3Word.u32>>1), iPtr.dispOffset)
		lobyte := memory.TestDwbit(dg.DwordT(oneAccMode3Word.u32), 31)
		cpu.ac[oneAccMode3Word.acd] = dg.DwordT(memory.ReadByte(addr, lobyte))

	case instrLLEF:
		oneAccModeInd3Word := iPtr.variant.(oneAccModeInd3WordT)
		cpu.ac[oneAccModeInd3Word.acd] = dg.DwordT(
			resolve31bitDisplacement(cpu, oneAccModeInd3Word.ind, oneAccModeInd3Word.mode, oneAccModeInd3Word.disp31, iPtr.dispOffset))

	case instrLLEFB:
		oneAccMode3Word := iPtr.variant.(oneAccMode3WordT)
		addr := resolve32bitEffAddr(cpu, ' ', oneAccMode3Word.mode, int32(oneAccMode3Word.u32>>1), iPtr.dispOffset)
		addr <<= 1
		if memory.TestDwbit(dg.DwordT(oneAccMode3Word.u32), 31) {
			addr |= 1
		}
		cpu.ac[oneAccMode3Word.acd] = dg.DwordT(addr)

	case instrLNADI:
		noAccModeImmInd3Word := iPtr.variant.(noAccModeImmInd3WordT)
		addr := resolve31bitDisplacement(cpu, noAccModeImmInd3Word.ind, noAccModeImmInd3Word.mode, noAccModeImmInd3Word.disp31, iPtr.dispOffset)
		wd := memory.ReadWord(addr)
		wd += dg.WordT(noAccModeImmInd3Word.immU16)
		memory.WriteWord(addr, wd)

	case instrLNLDA:
		oneAccModeInd3Word := iPtr.variant.(oneAccModeInd3WordT)
		addr := resolve31bitDisplacement(cpu, oneAccModeInd3Word.ind, oneAccModeInd3Word.mode, oneAccModeInd3Word.disp31, iPtr.dispOffset)
		cpu.ac[oneAccModeInd3Word.acd] = memory.SexWordToDword(memory.ReadWord(addr))

	case instrLNSTA:
		oneAccModeInd3Word := iPtr.variant.(oneAccModeInd3WordT)
		addr := resolve31bitDisplacement(cpu, oneAccModeInd3Word.ind, oneAccModeInd3Word.mode, oneAccModeInd3Word.disp31, iPtr.dispOffset)
		wd := memory.DwordGetLowerWord(cpu.ac[oneAccModeInd3Word.acd])
		memory.WriteWord(addr, wd)

	case instrLSTB:
		oneAccMode3Word := iPtr.variant.(oneAccMode3WordT)
		addr := resolve32bitEffAddr(cpu, ' ', oneAccMode3Word.mode, int32(oneAccMode3Word.u32>>1), iPtr.dispOffset)
		lobyte := memory.TestDwbit(dg.DwordT(oneAccMode3Word.u32), 31)
		memory.WriteByte(addr, lobyte, dg.ByteT(cpu.ac[oneAccMode3Word.acd]))

	case instrLWADD, instrLWSUB:
		oneAccModeInd3Word := iPtr.variant.(oneAccModeInd3WordT)
		addr := resolve31bitDisplacement(cpu, oneAccModeInd3Word.ind, oneAccModeInd3Word.mode, oneAccModeInd3Word.disp31, iPtr.dispOffset)
		var s32 int32
		switch iPtr.ix {
		case instrLWADD:
			s32 = int32(memory.ReadDWord(addr)) + int32(cpu.ac[oneAccModeInd3Word.acd])
		case instrLWSUB:
			s32 = int32(cpu.ac[oneAccModeInd3Word.acd]) - int32(memory.ReadDWord(addr))
		}
		cpu.ac[oneAccModeInd3Word.acd] = dg.DwordT(s32)

	case instrLWLDA:
		oneAccModeInd3Word := iPtr.variant.(oneAccModeInd3WordT)
		addr := resolve31bitDisplacement(cpu, oneAccModeInd3Word.ind, oneAccModeInd3Word.mode, oneAccModeInd3Word.disp31, iPtr.dispOffset)
		cpu.ac[oneAccModeInd3Word.acd] = memory.ReadDWord(addr)

	case instrLWSTA:
		oneAccModeInd3Word := iPtr.variant.(oneAccModeInd3WordT)
		addr := resolve31bitDisplacement(cpu, oneAccModeInd3Word.ind, oneAccModeInd3Word.mode, oneAccModeInd3Word.disp31, iPtr.dispOffset)
		memory.WriteDWord(addr, cpu.ac[oneAccModeInd3Word.acd])

	case instrWBLM:
		wblm(cpu)

	case instrWBTO:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		var addr dg.PhysAddrT
		if twoAcc1Word.acs == twoAcc1Word.acd {
			addr = cpu.pc & ringMask32
		} else {
			addr = resolve32bitIndirectableAddr(cpu, cpu.ac[twoAcc1Word.acs])
		}
		offset := dg.PhysAddrT(cpu.ac[twoAcc1Word.acd]) >> 4
		bitNum := uint(cpu.ac[twoAcc1Word.acd] & 0x0f)
		wd := memory.ReadWord(addr + offset)
		memory.SetWbit(&wd, bitNum)
		memory.WriteWord(addr+offset, wd)

	case instrWBTZ:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		var addr dg.PhysAddrT
		if twoAcc1Word.acs == twoAcc1Word.acd {
			addr = cpu.pc & ringMask32
		} else {
			addr = resolve32bitIndirectableAddr(cpu, cpu.ac[twoAcc1Word.acs])
		}
		offset := dg.PhysAddrT(cpu.ac[twoAcc1Word.acd]) >> 4
		bitNum := uint(cpu.ac[twoAcc1Word.acd] & 0x0f)
		wd := memory.ReadWord(addr + offset)
		memory.ClearWbit(&wd, bitNum)
		memory.WriteWord(addr+offset, wd)

	case instrWCMV:
		wcmv(cpu)

	case instrWCMP:
		wcmp(cpu)

	case instrWCST:
		wcst(cpu)

	case instrWCTR:
		wctr(cpu)

	case instrWLDB:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		wordAddr := dg.PhysAddrT(cpu.ac[twoAcc1Word.acs]) >> 1
		lowByte := memory.TestDwbit(cpu.ac[twoAcc1Word.acs], 31)
		cpu.ac[twoAcc1Word.acd] = dg.DwordT(memory.ReadByte(wordAddr, lowByte))

	case instrWSTB:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		memWriteByteBA(dg.ByteT(cpu.ac[twoAcc1Word.acd]&0x0ff), cpu.ac[twoAcc1Word.acs])

	case instrXLDB:
		oneAccMode2Word := iPtr.variant.(oneAccMode2WordT)
		disp := int32(oneAccMode2Word.disp16 >> 1)
		switch oneAccMode2Word.mode {
		case absoluteMode:
			disp &= 0x1fff_ffff
			disp |= int32(cpu.pc & 0x7000_0000)
			// case ac2Mode:
			// 	cpu.ac[2] >>= 1
			// case ac3Mode:
			// 	cpu.ac[3] >>= 1

		}
		addr := resolve32bitEffAddr(cpu, ' ', oneAccMode2Word.mode, disp, iPtr.dispOffset)
		cpu.ac[oneAccMode2Word.acd] = dg.DwordT(memory.ReadByte(addr, oneAccMode2Word.bitLow)) & 0x00ff

	case instrXLEF:
		oneAccModeInd2Word := iPtr.variant.(oneAccModeInd2WordT)
		cpu.ac[oneAccModeInd2Word.acd] = dg.DwordT(resolve15bitDisplacement(cpu, oneAccModeInd2Word.ind, oneAccModeInd2Word.mode, dg.WordT(oneAccModeInd2Word.disp15), iPtr.dispOffset))

	case instrXLEFB:
		oneAccMode2Word := iPtr.variant.(oneAccMode2WordT)
		disp := int32(oneAccMode2Word.disp16)
		if oneAccMode2Word.mode == absoluteMode {
			disp &= 0x1fff_ffff
			disp |= int32(cpu.pc & 0x7000_0000)
		}
		addr := resolve32bitEffAddr(cpu, 0, oneAccMode2Word.mode, disp, iPtr.dispOffset)
		addr <<= 1
		if oneAccMode2Word.bitLow {
			addr++
		}
		cpu.ac[oneAccMode2Word.acd] = dg.DwordT(addr)

	case instrXNADI, instrXNSBI:
		immMode2Word := iPtr.variant.(immMode2WordT)
		addr := resolve15bitDisplacement(cpu, immMode2Word.ind, immMode2Word.mode, dg.WordT(immMode2Word.disp15), iPtr.dispOffset)
		var s32 int32
		if iPtr.ix == instrXNADI {
			s32 = int32(int16(memory.ReadWord(addr))) + int32(immMode2Word.immU16)
		} else {
			s32 = int32(int16(memory.ReadWord(addr))) - int32(immMode2Word.immU16)
		}
		if (s32 > maxPosS16) || (s32 < minNegS16) {
			cpu.carry = true
			cpu.SetOVR(true)
		}
		memory.WriteWord(addr, dg.WordT(s32))

	case instrXNADD, instrXNSUB:
		oneAccModeInd2Word := iPtr.variant.(oneAccModeInd2WordT)
		addr := resolve15bitDisplacement(cpu, oneAccModeInd2Word.ind, oneAccModeInd2Word.mode, dg.WordT(oneAccModeInd2Word.disp15), iPtr.dispOffset)
		i16mem := int16(memory.ReadWord(addr))
		i16ac := int16(memory.DwordGetLowerWord(cpu.ac[oneAccModeInd2Word.acd]))
		var t32 int32
		if iPtr.ix == instrXNADD {
			i16ac += i16mem
			t32 = int32(i16ac) + int32(i16mem)
		} else {
			i16ac -= i16mem
			t32 = int32(i16ac) - int32(i16mem)
		}
		if t32 > maxPosS16 || t32 < minNegS16 {
			cpu.carry = true
			cpu.SetOVR(true)
		}
		cpu.ac[oneAccModeInd2Word.acd] = memory.SexWordToDword(dg.WordT(i16mem))

	case instrXNLDA:
		oneAccModeInd2Word := iPtr.variant.(oneAccModeInd2WordT)
		addr := resolve15bitDisplacement(cpu, oneAccModeInd2Word.ind, oneAccModeInd2Word.mode, dg.WordT(oneAccModeInd2Word.disp15), iPtr.dispOffset)
		wd, ok := memory.ReadWordTrap(addr)
		if !ok {
			return false
		}
		cpu.ac[oneAccModeInd2Word.acd] = memory.SexWordToDword(wd)

	case instrXNSTA:
		oneAccModeInd2Word := iPtr.variant.(oneAccModeInd2WordT)
		addr := resolve15bitDisplacement(cpu, oneAccModeInd2Word.ind, oneAccModeInd2Word.mode, dg.WordT(oneAccModeInd2Word.disp15), iPtr.dispOffset)
		memory.WriteWord(addr, memory.DwordGetLowerWord(cpu.ac[oneAccModeInd2Word.acd]))

	case instrXSTB:
		oneAccMode2Word := iPtr.variant.(oneAccMode2WordT)
		byt := dg.ByteT(cpu.ac[oneAccMode2Word.acd])
		disp := int32(oneAccMode2Word.disp16)
		if oneAccMode2Word.mode == absoluteMode {
			disp &= 0x1fff_ffff
			disp |= int32(cpu.pc & 0x7000_0000)
		}
		memory.WriteByte(resolve32bitEffAddr(cpu, ' ', oneAccMode2Word.mode, disp, iPtr.dispOffset), oneAccMode2Word.bitLow, byt)

	case instrXWADD, instrXWSUB, instrXWMUL:
		oneAccModeInd2Word := iPtr.variant.(oneAccModeInd2WordT)
		addr := resolve15bitDisplacement(cpu, oneAccModeInd2Word.ind, oneAccModeInd2Word.mode, dg.WordT(oneAccModeInd2Word.disp15), iPtr.dispOffset)
		s64mem := int64(int32(memory.ReadDWord(addr)))
		s64ac := int64(int32(cpu.ac[oneAccModeInd2Word.acd]))
		var t64 int64
		switch iPtr.ix {
		case instrXWADD:
			t64 = s64ac + s64mem
		case instrXWMUL:
			t64 = s64ac * s64mem
		case instrXWSUB:
			t64 = s64ac - s64mem
		}
		if t64 > maxPosS32 || t64 < minNegS32 {
			cpu.carry = true
			cpu.SetOVR(true)
		}
		cpu.ac[oneAccModeInd2Word.acd] = dg.DwordT(t64)

	case instrXWADI, instrXWSBI:
		immMode2Word := iPtr.variant.(immMode2WordT)
		addr := resolve15bitDisplacement(cpu, immMode2Word.ind, immMode2Word.mode, dg.WordT(immMode2Word.disp15), iPtr.dispOffset)
		var s64 int64
		switch iPtr.ix {
		case instrXWADI:
			s64 = int64(int32(memory.ReadDWord(addr))) + int64(immMode2Word.immU16)
		case instrXWSBI:
			s64 = int64(int32(memory.ReadDWord(addr))) - int64(immMode2Word.immU16)
		}
		if (s64 > maxPosS32) || (s64 < minNegS32) {
			cpu.carry = true
			cpu.SetOVR(true)
		}
		memory.WriteDWord(addr, dg.DwordT(s64))

	case instrXWLDA:
		oneAccModeInd2Word := iPtr.variant.(oneAccModeInd2WordT)
		addr := resolve15bitDisplacement(cpu, oneAccModeInd2Word.ind, oneAccModeInd2Word.mode, dg.WordT(oneAccModeInd2Word.disp15), iPtr.dispOffset)
		cpu.ac[oneAccModeInd2Word.acd] = memory.ReadDWord(addr)

	case instrXWSTA:
		oneAccModeInd2Word := iPtr.variant.(oneAccModeInd2WordT)
		addr := resolve15bitDisplacement(cpu, oneAccModeInd2Word.ind, oneAccModeInd2Word.mode, dg.WordT(oneAccModeInd2Word.disp15), iPtr.dispOffset)
		memory.WriteDWord(addr, cpu.ac[oneAccModeInd2Word.acd])

	default:
		log.Fatalf("ERROR: EAGLE_MEMREF instruction <%s> not yet implemented\n", iPtr.mnemonic)
		return false
	}

	cpu.pc += dg.PhysAddrT(iPtr.instrLength)
	return true
}

func readByteBA(ba dg.DwordT) dg.ByteT {
	return memory.ReadByte(resolve32bitByteAddr(ba))
}

// memWriteByte writes the supplied byte to the address derived from the given byte addr
func memWriteByteBA(b dg.ByteT, ba dg.DwordT) {
	wordAddr, lowByte := resolve32bitByteAddr(ba)
	memory.WriteByte(wordAddr, lowByte, b)
	// if cpu.debugLogging {
	// 	logging.DebugPrint(logging.DebugLog, "DEBUG: memWriteByte wrote %c to word addr: %#o\n", b, wordAddr)
	// }
}

func copyByte(srcBA, destBA dg.DwordT) {
	memWriteByteBA(readByteBA(srcBA), destBA)
}

func wblm(cpu *CPUT) {
	/* AC0 - unused, AC1 - no. wds to move (if neg then descending order), AC2 - src, AC3 - dest */
	if cpu.ac[1] == 0 {
		log.Println("INFO: WBLM called with AC1 == 0, not moving anything")
		return
	}
	if cpu.debugLogging {
		logging.DebugPrint(logging.DebugLog, "DEBUG: WBLM moving %#o words from %#o to %#o\n",
			int32(cpu.ac[1]), cpu.ac[2], cpu.ac[3])
	}
	for cpu.ac[1] != 0 {
		memory.WriteWord(dg.PhysAddrT(cpu.ac[3]), memory.ReadWord(dg.PhysAddrT(cpu.ac[2])))
		if memory.TestDwbit(cpu.ac[1], 0) {
			cpu.ac[1]++
			cpu.ac[2]--
			cpu.ac[3]--
		} else {
			cpu.ac[1]--
			cpu.ac[2]++
			cpu.ac[3]++
		}
	}
}

func wcmv(cpu *CPUT) {
	// ACO destCount, AC1 srcCount, AC2 dest byte ptr, AC3 src byte ptr
	destCount := int32(cpu.ac[0])
	if destCount == 0 {
		if cpu.debugLogging {
			logging.DebugPrint(logging.DebugLog, ".... WCMV called with AC0 == 0, not moving anything\n")
		}
		return
	}
	destAscend := (destCount > 0)
	srcCount := int32(cpu.ac[1])
	srcAscend := (srcCount > 0)
	if cpu.debugLogging {
		logging.DebugPrint(logging.DebugLog, ".... WCMV moving %#o chars from %#o to %#o chars at %#o\n",
			srcCount, cpu.ac[3], destCount, cpu.ac[2])
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
	cpu.ac[1] = 0 // we're treating this as atomic...
}

func getDirection(ac dg.DwordT) int32 {
	if int32(ac) > 0 {
		return 1
	}
	if int32(ac) < 0 {
		return -1
	}
	return 0
}

func wcmp(cpu *CPUT) {
	// AC0 String2 length and dir (bwd if -ve)
	// AC1 String1 length and dir (bwd if -ve)
	// AC2 Byte Pointer to first byte of String2 to be compared
	// AC3 Byte Pointer to first byte of String1 to be compared
	str2dir := getDirection(cpu.ac[0])
	str1dir := getDirection(cpu.ac[1])
	var str1char, str2char dg.ByteT
	if str1dir == 0 && str2dir == 0 {
		return
	}
	for cpu.ac[1] != 0 && cpu.ac[0] != 0 {
		// read the two bytes to compare, substitute with a space if one string has run out
		if cpu.ac[1] != 0 {
			str1char = readByteBA(cpu.ac[3])
		} else {
			str1char = ' '
		}
		if cpu.ac[0] != 0 {
			str2char = readByteBA(cpu.ac[2])
		} else {
			str2char = ' '
		}
		// compare
		if str1char < str2char {
			cpu.ac[1] = 0xFFFFFFFF
			return
		}
		if str1char > str2char {
			cpu.ac[1] = 1
			return
		}
		// they were equal, so adjust remaining lengths, move pointers, and loop round
		if cpu.ac[0] != 0 {
			cpu.ac[0] = dg.DwordT(int32(cpu.ac[0]) + str2dir)
		}
		if cpu.ac[1] != 0 {
			cpu.ac[1] = dg.DwordT(int32(cpu.ac[1]) + str1dir)
		}
		cpu.ac[2] = dg.DwordT(int32(cpu.ac[2]) + str2dir)
		cpu.ac[3] = dg.DwordT(int32(cpu.ac[3]) + str1dir)
	}
}

func wcst(cpu *CPUT) {
	strLenDir := int(int32(cpu.ac[1]))
	if strLenDir == 0 {
		if cpu.debugLogging {
			logging.DebugPrint(logging.DebugLog, ".... WCST called with AC1 == 0, not scanning anything\n")
		}
		return
	}
	delimTabAddr := resolve32bitIndirectableAddr(cpu, cpu.ac[0])
	cpu.ac[0] = dg.DwordT(delimTabAddr)
	// load the table which is 256 bits stored as 16 words
	var table [256]bool
	var tIx dg.PhysAddrT
	for tIx = 0; tIx < 16; tIx++ {
		wd := memory.ReadWord(delimTabAddr + tIx)
		for bit := 0; bit < 16; bit++ {
			if memory.TestWbit(wd, bit) {
				table[(int(tIx)*16)+bit] = true
			}
		}
	}
	// table[] now contains true for any delimiter
	var dir int32 = 1
	if strLenDir < 0 {
		dir = -1
	}

	for strLenDir != 0 {
		thisChar := readByteBA(cpu.ac[3])
		if table[int(thisChar)] {
			// match, so set AC1 and return
			cpu.ac[1] = 0
			return
		}
		cpu.ac[1] = dg.DwordT(int32(cpu.ac[1]) + dir)
		cpu.ac[3] = dg.DwordT(int32(cpu.ac[3]) + dir)
		strLenDir += int(dir)
	}
}

func wctr(cpu *CPUT) {
	// AC0 Wide Byte addr of translation table - unchanged
	// AC1 # of bytes in each string, NB. -ve => translate-and-move mode, +ve => translate-and-compare mode
	// AC2 destination string ("string2") Byte addr
	// AC3 source string ("string1") byte addr
	if cpu.ac[1] == 0 {
		if cpu.debugLogging {
			logging.DebugPrint(logging.DebugLog, "INFO: WCTR called with AC1 == 0, not translating anything\n")
		}
		return
	}
	transTablePtr := dg.DwordT(resolve32bitIndirectableAddr(cpu, cpu.ac[0]))
	// build an array representation of the table
	var transTable [256]dg.ByteT
	var c dg.DwordT
	for c = 0; c < 256; c++ {
		transTable[c] = readByteBA(transTablePtr + c)
	}

	for cpu.ac[1] != 0 {
		srcByte := readByteBA(cpu.ac[3])
		cpu.ac[3]++
		transByte := transTable[int(srcByte)]
		if int32(cpu.ac[1]) < 0 {
			// move mode
			memWriteByteBA(transByte, cpu.ac[2])
			cpu.ac[2]++
			cpu.ac[1]++
		} else {
			// compare mode
			str2byte := readByteBA(cpu.ac[2])
			cpu.ac[2]++
			trans2byte := transTable[int(str2byte)]
			if srcByte < trans2byte {
				cpu.ac[1] = 0xffffffff
				break
			}
			if srcByte > trans2byte {
				cpu.ac[1] = 1
				break
			}
			cpu.ac[1]--
		}
	}
}
