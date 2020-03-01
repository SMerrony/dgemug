// eagleMemRef.go

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

func eagleMemRef(cpuPtr *CPUT, iPtr *decodedInstrT) bool {

	switch iPtr.ix {

	case instrLNLDA:
		oneAccModeInd3Word := iPtr.variant.(oneAccModeInd3WordT)
		addr := resolve32bitEffAddr(cpuPtr, oneAccModeInd3Word.ind, oneAccModeInd3Word.mode, oneAccModeInd3Word.disp31, iPtr.dispOffset)
		cpuPtr.ac[oneAccModeInd3Word.acd] = memory.SexWordToDword(memory.ReadWord(addr))

	case instrLNSTA:
		oneAccModeInd3Word := iPtr.variant.(oneAccModeInd3WordT)
		addr := resolve32bitEffAddr(cpuPtr, oneAccModeInd3Word.ind, oneAccModeInd3Word.mode, oneAccModeInd3Word.disp31, iPtr.dispOffset)
		wd := memory.DwordGetLowerWord(cpuPtr.ac[oneAccModeInd3Word.acd])
		memory.WriteWord(addr, wd)

	case instrLWLDA:
		oneAccModeInd3Word := iPtr.variant.(oneAccModeInd3WordT)
		addr := resolve32bitEffAddr(cpuPtr, oneAccModeInd3Word.ind, oneAccModeInd3Word.mode, oneAccModeInd3Word.disp31, iPtr.dispOffset)
		cpuPtr.ac[oneAccModeInd3Word.acd] = memory.ReadDWord(addr)

	case instrLWSTA:
		oneAccModeInd3Word := iPtr.variant.(oneAccModeInd3WordT)
		addr := resolve32bitEffAddr(cpuPtr, oneAccModeInd3Word.ind, oneAccModeInd3Word.mode, oneAccModeInd3Word.disp31, iPtr.dispOffset)
		memory.WriteDWord(addr, cpuPtr.ac[oneAccModeInd3Word.acd])

	case instrWBLM:
		wblm(cpuPtr)

	case instrWBTO:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		var addr dg.PhysAddrT
		if twoAcc1Word.acs == twoAcc1Word.acd {
			addr = 0
		} else {
			addr = resolve32bitIndirectableAddr(cpuPtr, cpuPtr.ac[twoAcc1Word.acs])
		}
		offset := dg.PhysAddrT(cpuPtr.ac[twoAcc1Word.acd]) >> 4
		bitNum := uint(cpuPtr.ac[twoAcc1Word.acd] & 0x0f)
		wd := memory.ReadWord(addr + offset)
		memory.SetWbit(&wd, bitNum)
		memory.WriteWord(addr+offset, wd)

	case instrWBTZ:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		var addr dg.PhysAddrT
		if twoAcc1Word.acs == twoAcc1Word.acd {
			addr = 0
		} else {
			addr = resolve32bitIndirectableAddr(cpuPtr, cpuPtr.ac[twoAcc1Word.acs])
		}
		offset := dg.PhysAddrT(cpuPtr.ac[twoAcc1Word.acd]) >> 4
		bitNum := uint(cpuPtr.ac[twoAcc1Word.acd] & 0x0f)
		wd := memory.ReadWord(addr + offset)
		memory.ClearWbit(&wd, bitNum)
		memory.WriteWord(addr+offset, wd)

	case instrWCMV:
		wcmv(cpuPtr)

	case instrWCMP:
		wcmp(cpuPtr)

	case instrWCST:
		wcst(cpuPtr)

	case instrWCTR:
		wctr(cpuPtr)

	case instrWLDB:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		wordAddr := dg.PhysAddrT(cpuPtr.ac[twoAcc1Word.acs]) >> 1
		lowByte := memory.TestDwbit(cpuPtr.ac[twoAcc1Word.acs], 31)
		cpuPtr.ac[twoAcc1Word.acd] = dg.DwordT(memory.ReadByte(wordAddr, lowByte))

	case instrWSTB:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		memWriteByteBA(dg.ByteT(cpuPtr.ac[twoAcc1Word.acd]&0x0ff), cpuPtr.ac[twoAcc1Word.acs])

	case instrXLDB:
		oneAccMode2Word := iPtr.variant.(oneAccMode2WordT)
		disp := int32(oneAccMode2Word.disp16 >> 1)
		switch oneAccMode2Word.mode {
		case absoluteMode:
			disp &= 0x1fff_ffff
			disp |= int32(cpuPtr.pc & 0x7000_0000)
			// case ac2Mode:
			// 	cpuPtr.ac[2] >>= 1
			// case ac3Mode:
			// 	cpuPtr.ac[3] >>= 1

		}
		addr := resolve32bitEffAddr(cpuPtr, ' ', oneAccMode2Word.mode, disp, iPtr.dispOffset)
		cpuPtr.ac[oneAccMode2Word.acd] = dg.DwordT(memory.ReadByte(addr, oneAccMode2Word.bitLow)) & 0x00ff

	case instrXLEF:
		oneAccModeInd2Word := iPtr.variant.(oneAccModeInd2WordT)
		cpuPtr.ac[oneAccModeInd2Word.acd] = dg.DwordT(resolve15bitDisplacement(cpuPtr, oneAccModeInd2Word.ind, oneAccModeInd2Word.mode, dg.WordT(oneAccModeInd2Word.disp15), iPtr.dispOffset))

	case instrXLEFB:
		oneAccMode2Word := iPtr.variant.(oneAccMode2WordT)
		disp := int32(oneAccMode2Word.disp16)
		if oneAccMode2Word.mode == absoluteMode {
			disp &= 0x1fff_ffff
			disp |= int32(cpuPtr.pc & 0x7000_0000)
		}
		addr := resolve32bitEffAddr(cpuPtr, 0, oneAccMode2Word.mode, disp, iPtr.dispOffset)
		addr <<= 1
		if !oneAccMode2Word.bitLow {
			addr++
		}
		cpuPtr.ac[oneAccMode2Word.acd] = dg.DwordT(addr)

	case instrXNADD, instrXNSUB:
		oneAccModeInd2Word := iPtr.variant.(oneAccModeInd2WordT)
		addr := resolve15bitDisplacement(cpuPtr, oneAccModeInd2Word.ind, oneAccModeInd2Word.mode, dg.WordT(oneAccModeInd2Word.disp15), iPtr.dispOffset)
		i16mem := int16(memory.ReadWord(addr))
		i16ac := int16(memory.DwordGetLowerWord(cpuPtr.ac[oneAccModeInd2Word.acd]))
		var t32 int32
		if iPtr.ix == instrXNADD {
			i16ac += i16mem
			t32 = int32(i16ac) + int32(i16mem)
		} else {
			i16ac -= i16mem
			t32 = int32(i16ac) - int32(i16mem)
		}
		if t32 > maxPosS16 || t32 < minNegS16 {
			cpuPtr.carry = true
			cpuPtr.CPUSetOVR(true)
		}
		cpuPtr.ac[oneAccModeInd2Word.acd] = memory.SexWordToDword(dg.WordT(i16mem))

	case instrXNLDA:
		oneAccModeInd2Word := iPtr.variant.(oneAccModeInd2WordT)
		addr := resolve15bitDisplacement(cpuPtr, oneAccModeInd2Word.ind, oneAccModeInd2Word.mode, dg.WordT(oneAccModeInd2Word.disp15), iPtr.dispOffset)
		wd, ok := memory.ReadWordTrap(addr)
		if !ok {
			return false
		}
		cpuPtr.ac[oneAccModeInd2Word.acd] = memory.SexWordToDword(wd)

	case instrXNSTA:
		oneAccModeInd2Word := iPtr.variant.(oneAccModeInd2WordT)
		addr := resolve15bitDisplacement(cpuPtr, oneAccModeInd2Word.ind, oneAccModeInd2Word.mode, dg.WordT(oneAccModeInd2Word.disp15), iPtr.dispOffset)
		memory.WriteWord(addr, memory.DwordGetLowerWord(cpuPtr.ac[oneAccModeInd2Word.acd]))

	case instrXSTB:
		oneAccMode2Word := iPtr.variant.(oneAccMode2WordT)
		byt := dg.ByteT(cpuPtr.ac[oneAccMode2Word.acd])
		disp := int32(oneAccMode2Word.disp16)
		if oneAccMode2Word.mode == absoluteMode {
			disp &= 0x1fff_ffff
			disp |= int32(cpuPtr.pc & 0x7000_0000)
		}
		memory.WriteByte(resolve32bitEffAddr(cpuPtr, ' ', oneAccMode2Word.mode, disp, iPtr.dispOffset), oneAccMode2Word.bitLow, byt)

	case instrXWADI:
		immMode2Word := iPtr.variant.(immMode2WordT)
		addr := resolve15bitDisplacement(cpuPtr, immMode2Word.ind, immMode2Word.mode, dg.WordT(immMode2Word.disp15), iPtr.dispOffset)
		s64 := int64(memory.ReadDWord(addr)) + int64(immMode2Word.immU16)
		if (s64 > maxPosS32) || (s64 < minNegS32) {
			cpuPtr.carry = true
			cpuPtr.CPUSetOVR(true)
		}
		memory.WriteDWord(addr, dg.DwordT(s64))

	case instrXWLDA:
		oneAccModeInd2Word := iPtr.variant.(oneAccModeInd2WordT)
		addr := resolve15bitDisplacement(cpuPtr, oneAccModeInd2Word.ind, oneAccModeInd2Word.mode, dg.WordT(oneAccModeInd2Word.disp15), iPtr.dispOffset)
		cpuPtr.ac[oneAccModeInd2Word.acd] = memory.ReadDWord(addr)

	case instrXWSTA:
		oneAccModeInd2Word := iPtr.variant.(oneAccModeInd2WordT)
		addr := resolve15bitDisplacement(cpuPtr, oneAccModeInd2Word.ind, oneAccModeInd2Word.mode, dg.WordT(oneAccModeInd2Word.disp15), iPtr.dispOffset)
		memory.WriteDWord(addr, cpuPtr.ac[oneAccModeInd2Word.acd])

	default:
		log.Fatalf("ERROR: EAGLE_MEMREF instruction <%s> not yet implemented\n", iPtr.mnemonic)
		return false
	}

	cpuPtr.pc += dg.PhysAddrT(iPtr.instrLength)
	return true
}

func readByteBA(ba dg.DwordT) dg.ByteT {
	return memory.ReadByte(resolve32bitByteAddr(ba))
}

// memWriteByte writes the supplied byte to the address derived from the given byte addr
func memWriteByteBA(b dg.ByteT, ba dg.DwordT) {
	wordAddr, lowByte := resolve32bitByteAddr(ba)
	memory.WriteByte(wordAddr, lowByte, b)
	// if debugLogging {
	// 	logging.DebugPrint(logging.DebugLog, "DEBUG: memWriteByte wrote %c to word addr: %#o\n", b, wordAddr)
	// }
}

func copyByte(srcBA, destBA dg.DwordT) {
	memWriteByteBA(readByteBA(srcBA), destBA)
}

func wblm(cpuPtr *CPUT) {
	/* AC0 - unused, AC1 - no. wds to move (if neg then descending order), AC2 - src, AC3 - dest */
	if cpuPtr.ac[1] == 0 {
		log.Println("INFO: WBLM called with AC1 == 0, not moving anything")
		return
	}
	if debugLogging {
		logging.DebugPrint(logging.DebugLog, "DEBUG: WBLM moving %#o words from %#o to %#o\n",
			int32(cpuPtr.ac[1]), cpuPtr.ac[2], cpuPtr.ac[3])
	}
	for cpuPtr.ac[1] != 0 {
		memory.WriteWord(dg.PhysAddrT(cpuPtr.ac[3]), memory.ReadWord(dg.PhysAddrT(cpuPtr.ac[2])))
		if memory.TestDwbit(cpuPtr.ac[1], 0) {
			cpuPtr.ac[1]++
			cpuPtr.ac[2]--
			cpuPtr.ac[3]--
		} else {
			cpuPtr.ac[1]--
			cpuPtr.ac[2]++
			cpuPtr.ac[3]++
		}
	}
}

func wcmv(cpuPtr *CPUT) {
	// ACO destCount, AC1 srcCount, AC2 dest byte ptr, AC3 src byte ptr
	destCount := int32(cpuPtr.ac[0])
	if destCount == 0 {
		if debugLogging {
			logging.DebugPrint(logging.DebugLog, ".... WCMV called with AC0 == 0, not moving anything\n")
		}
		return
	}
	destAscend := (destCount > 0)
	srcCount := int32(cpuPtr.ac[1])
	srcAscend := (srcCount > 0)
	if debugLogging {
		logging.DebugPrint(logging.DebugLog, ".... WCMV moving %#o chars from %#o to %#o chars at %#o\n",
			srcCount, cpuPtr.ac[3], destCount, cpuPtr.ac[2])
	}
	// set carry if length of src is greater than length of dest
	if cpuPtr.ac[1] > cpuPtr.ac[2] {
		cpuPtr.carry = true
	}
	// 1st move srcCount bytes
	for {
		copyByte(cpuPtr.ac[3], cpuPtr.ac[2])
		if srcAscend {
			cpuPtr.ac[3]++
			srcCount--
		} else {
			cpuPtr.ac[3]--
			srcCount++
		}
		if destAscend {
			cpuPtr.ac[2]++
			destCount--
		} else {
			cpuPtr.ac[2]--
			destCount++
		}
		if srcCount == 0 || destCount == 0 {
			break
		}
	}
	// now fill any excess bytes with ASCII spaces
	if destCount != 0 {
		for {
			memWriteByteBA(asciiSPC, cpuPtr.ac[2])
			if destAscend {
				cpuPtr.ac[2]++
				destCount--
			} else {
				cpuPtr.ac[2]--
				destCount++
			}
			if destCount == 0 {
				break
			}
		}
	}
	cpuPtr.ac[0] = 0
	cpuPtr.ac[1] = 0 // we're treating this as atomic...
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

func wcmp(cpuPtr *CPUT) {
	// AC0 String2 length and dir (bwd if -ve)
	// AC1 String1 length and dir (bwd if -ve)
	// AC2 Byte Pointer to first byte of String2 to be compared
	// AC3 Byte Pointer to first byte of String1 to be compared
	str2dir := getDirection(cpuPtr.ac[0])
	str1dir := getDirection(cpuPtr.ac[1])
	var str1char, str2char dg.ByteT
	if str1dir == 0 && str2dir == 0 {
		return
	}
	for cpuPtr.ac[1] != 0 && cpuPtr.ac[0] != 0 {
		// read the two bytes to compare, substitute with a space if one string has run out
		if cpuPtr.ac[1] != 0 {
			str1char = readByteBA(cpuPtr.ac[3])
		} else {
			str1char = ' '
		}
		if cpuPtr.ac[0] != 0 {
			str2char = readByteBA(cpuPtr.ac[2])
		} else {
			str2char = ' '
		}
		// compare
		if str1char < str2char {
			cpuPtr.ac[1] = 0xFFFFFFFF
			return
		}
		if str1char > str2char {
			cpuPtr.ac[1] = 1
			return
		}
		// they were equal, so adjust remaining lengths, move pointers, and loop round
		if cpuPtr.ac[0] != 0 {
			cpuPtr.ac[0] = dg.DwordT(int32(cpuPtr.ac[0]) + str2dir)
		}
		if cpuPtr.ac[1] != 0 {
			cpuPtr.ac[1] = dg.DwordT(int32(cpuPtr.ac[1]) + str1dir)
		}
		cpuPtr.ac[2] = dg.DwordT(int32(cpuPtr.ac[2]) + str2dir)
		cpuPtr.ac[3] = dg.DwordT(int32(cpuPtr.ac[3]) + str1dir)
	}
}

func wcst(cpuPtr *CPUT) {
	strLenDir := int(int32(cpuPtr.ac[1]))
	if strLenDir == 0 {
		if debugLogging {
			logging.DebugPrint(logging.DebugLog, ".... WCST called with AC1 == 0, not scanning anything\n")
		}
		return
	}
	delimTabAddr := resolve32bitIndirectableAddr(cpuPtr, cpuPtr.ac[0])
	cpuPtr.ac[0] = dg.DwordT(delimTabAddr)
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
		thisChar := readByteBA(cpuPtr.ac[3])
		if table[int(thisChar)] {
			// match, so set AC1 and return
			cpuPtr.ac[1] = 0
			return
		}
		cpuPtr.ac[1] = dg.DwordT(int32(cpuPtr.ac[1]) + dir)
		cpuPtr.ac[3] = dg.DwordT(int32(cpuPtr.ac[3]) + dir)
		strLenDir += int(dir)
	}
}

func wctr(cpuPtr *CPUT) {
	// AC0 Wide Byte addr of translation table - unchanged
	// AC1 # of bytes in each string, NB. -ve => translate-and-move mode, +ve => translate-and-compare mode
	// AC2 destination string ("string2") Byte addr
	// AC3 source string ("string1") byte addr
	if cpuPtr.ac[1] == 0 {
		if debugLogging {
			logging.DebugPrint(logging.DebugLog, "INFO: WCTR called with AC1 == 0, not translating anything\n")
		}
		return
	}
	transTablePtr := dg.DwordT(resolve32bitIndirectableAddr(cpuPtr, cpuPtr.ac[0]))
	// build an array representation of the table
	var transTable [256]dg.ByteT
	var c dg.DwordT
	for c = 0; c < 256; c++ {
		transTable[c] = readByteBA(transTablePtr + c)
	}

	for cpuPtr.ac[1] != 0 {
		srcByte := readByteBA(cpuPtr.ac[3])
		cpuPtr.ac[3]++
		transByte := transTable[int(srcByte)]
		if int32(cpuPtr.ac[1]) < 0 {
			// move mode
			memWriteByteBA(transByte, cpuPtr.ac[2])
			cpuPtr.ac[2]++
			cpuPtr.ac[1]++
		} else {
			// compare mode
			str2byte := readByteBA(cpuPtr.ac[2])
			cpuPtr.ac[2]++
			trans2byte := transTable[int(str2byte)]
			if srcByte < trans2byte {
				cpuPtr.ac[1] = 0xffffffff
				break
			}
			if srcByte > trans2byte {
				cpuPtr.ac[1] = 1
				break
			}
			cpuPtr.ac[1]--
		}
	}
}
