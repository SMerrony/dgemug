// decoder.go

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

// BEWARE: DO NOT FALL INTO THE TRAP OF TRYING TO RESOLVE ADDRESSES HERE - DECODE ONLY!

package mvcpu

import (
	"fmt"

	"github.com/SMerrony/dgemug/devices"

	"github.com/SMerrony/dgemug/dg"
	"github.com/SMerrony/dgemug/logging"
	"github.com/SMerrony/dgemug/memory"
)

// decodedInstrT defines the MV/Em internal decode of an opcode and any
// parameters.
type decodedInstrT struct {
	ix          int
	mnemonic    string
	instrFmt    int
	instrType   int
	instrLength int
	dispOffset  int
	disassembly string
	variant     interface{}
}

// here are the types for the variant portion of the decoded instruction...
type derrT struct {
	errCode uint32
}
type immMode2WordT struct {
	immU16 uint16
	mode   int
	ind    byte
	disp15 int16
}
type immOneAccT struct {
	immU16 uint16
	acd    int
}
type ioFlagsDevT struct {
	f     byte
	ioDev int
}
type ioTestDevT struct {
	t     int
	ioDev int
}
type lndo4WordT struct {
	acd       int
	mode      int
	ind       byte
	disp31    int32
	offsetU16 uint16
}
type noAccMode2WordT struct {
	mode   int
	disp16 uint16
}
type noAccMode3WordT struct {
	mode   int
	immU32 uint32
}
type noAccModeInd2WordT struct {
	mode   int
	ind    byte
	disp15 dg.WordT
}
type noAccModeInd3WordT struct {
	mode   int
	ind    byte
	disp31 int32
}
type noAccModeInd3WordXcallT struct {
	mode     int
	ind      byte
	disp15   int16
	argCount int
}
type noAccModeImmInd3WordT struct {
	immU16 uint16
	mode   int
	ind    byte
	disp31 int32
}
type noAccModeInd4WordT struct {
	mode     int
	ind      byte
	disp31   int32
	argCount int
}
type novaDataIoT struct {
	acd   int
	f     byte
	ioDev int
}
type novaNoAccEffAddrT struct {
	mode   int
	ind    byte
	disp15 int16
}
type novaOneAccEffAddrT struct {
	acd    int
	mode   int
	ind    byte
	disp15 int16
}
type novaTwoAccMultOpT struct {
	acd, acs  int
	c, sh, nl byte
	skip      int
}
type oneAcc1WordT struct {
	acd int
}
type oneAccImm2WordT struct {
	acd    int
	immS16 int16
}
type oneAccImmWd2WordT struct {
	acd     int
	immWord dg.WordT
}
type oneAccImm3WordT struct {
	acd    int
	immU32 uint32
}
type oneAccImmDwd3WordT struct {
	acd      int
	immDword dg.DwordT
}
type oneAccMode2WordT struct {
	acd    int
	mode   int
	disp16 int16
	bitLow bool
}
type oneAccMode3WordT struct {
	acd    int
	mode   int
	disp31 int32
}
type oneAccModeInd2WordT struct {
	acd    int
	mode   int
	ind    byte
	disp15 int16
}
type oneAccModeInd3WordT struct {
	acd    int
	mode   int
	ind    byte
	disp31 int32
}
type split8bitDispT struct {
	disp8 int8
}
type threeWordDoT struct {
	acd       int
	mode      int
	ind       byte
	disp15    int16
	offsetU16 uint16
}
type twoAcc1WordT struct {
	acd, acs int
}
type twoAccImm2WordT struct {
	acd, acs int
	immWord  dg.WordT
}
type unique2WordT struct {
	immU16 uint16
}
type wskbT struct {
	bitNum int
}

const numPosOpcodes = 65536

var opCodeLookup [numPosOpcodes]int

// decoderGenAllPossOpcodes builds an array keyed by every possible DG Word
// containing the corresponding Op Code.  LEF is not included or handled here.
func decoderGenAllPossOpcodes() {
	for opcode := 0; opcode < numPosOpcodes; opcode++ {
		mnem, found := instructionMatch(dg.WordT(opcode), false, false, false)
		if found {
			opCodeLookup[opcode] = mnem
		} else {
			opCodeLookup[opcode] = -1
		}
	}
}

// InstructionFind looks up an opcode in the opcode lookup table and returns
// the corresponding mnemonic.  This needs to be as quick as possible
func instructionLookup(opcode dg.WordT, lefMode bool, ioOn bool, atuOn bool) int {
	if lefMode {
		// special case, if LEF mode is enabled then ALL I/O instructions
		// must be interpreted as LEF
		if memory.GetWbits(opcode, 0, 3) == 3 { // an I/O instruction
			logging.DebugPrint(logging.DebugLog, "DEBUG: instructionLookup() returning LEF\n")
			return instrLEF
		}
	}
	return opCodeLookup[opcode]
}

// instructionMatch looks for a match for the opcode in the instruction set and returns
// the corresponding mnemonic.  It is used only by the decoderGenAllPossOpcodes() above when
// MV/Em is initialising.
// N.B. LEF is ignored here.
func instructionMatch(opcode dg.WordT, lefMode bool, ioOn bool, atuOn bool) (int, bool) {
	var tail dg.WordT
	//for mnem, insChar := range instructionSet {
	for i := 0; i < len(instructionSet); i++ {
		mnem := i
		insChar := instructionSet[i]
		if (opcode & insChar.mask) == insChar.bits {
			// there are some exceptions to the normal decoding...
			switch mnem {
			case instrLEF:
				if lefMode {
					return -1, false
				}
			case instrADC, instrADD, instrAND, instrCOM, instrINC, instrMOV, instrNEG, instrSUB:
				// these instructions are not allowed to end in 1000(2) or 1001(2)
				// as those patterns are used for Eagle instructions
				tail = opcode & 0x000f
				if tail != 0b1000 && tail != 0b1001 {
					return mnem, true
				}
			default:
				return mnem, true

			}
		}
	}
	return 0, false
}

// InstructionDecode decodes an opcode
func InstructionDecode(opcode dg.WordT, pc dg.PhysAddrT, lefMode bool, ioOn bool, atuOn bool, disassemble bool, devMap devices.DeviceMapT) (*decodedInstrT, bool) {
	var decodedInstr decodedInstrT
	var secondWord, thirdWord, fourthWord dg.WordT

	decodedInstr.disassembly = "; Unknown instruction"

	ix := instructionLookup(opcode, lefMode, ioOn, atuOn)
	if ix == -1 {
		logging.DebugPrint(logging.DebugLog, "INFO: instructionDecode failed to find anything with instructionLookup for location %d., containing 0x%X\n", pc, opcode)
		return &decodedInstr, false
	}
	decodedInstr.ix = ix
	decodedInstr.mnemonic = instructionSet[ix].mnemonic
	decodedInstr.disassembly = instructionSet[ix].mnemonic
	decodedInstr.instrFmt = instructionSet[ix].instrFmt
	decodedInstr.instrType = instructionSet[ix].instrType
	decodedInstr.instrLength = instructionSet[ix].instrLen
	decodedInstr.dispOffset = instructionSet[ix].dispOffset

	switch decodedInstr.instrFmt {

	case DERR_FMT: // DERR has a unique format
		var derr derrT
		derr.errCode = uint32(memory.GetWbits(opcode, 1, 3)<<2) + uint32(memory.GetWbits(opcode, 10, 1))
		decodedInstr.variant = derr
		if disassemble {
			decodedInstr.disassembly += fmt.Sprintf(" %#o", derr.errCode)
		}
	case IMM_MODE_2_WORD_FMT: // eg. XNADI, XNSBI, XNSUB, XWADI, XWSBI
		var immMode2Word immMode2WordT
		immMode2Word.immU16 = decode2bitImm(memory.GetWbits(opcode, 1, 2))
		immMode2Word.mode = int(memory.GetWbits(opcode, 3, 2))
		secondWord = memory.ReadWord(pc + 1)
		immMode2Word.ind = decodeIndirect(memory.TestWbit(secondWord, 0))
		immMode2Word.disp15 = decode15bitDisp(secondWord, immMode2Word.mode)
		decodedInstr.variant = immMode2Word
		if disassemble {
			decodedInstr.disassembly += fmt.Sprintf(" %#o,%#o%s [2-Word OpCode]",
				immMode2Word.immU16, immMode2Word.disp15, modeToString(immMode2Word.mode))
		}
	case IMM_ONEACC_FMT: // eg. ADI, HXL, NADI, SBI, WADI, WLSI, WSBI
		// N.B. Immediate value is encoded by assembler to be one less than required
		//      This is handled by decode2bitImm()
		var immOneAcc immOneAccT
		immOneAcc.immU16 = decode2bitImm(memory.GetWbits(opcode, 1, 2))
		immOneAcc.acd = int(memory.GetWbits(opcode, 3, 2))
		decodedInstr.variant = immOneAcc
		if disassemble {
			decodedInstr.disassembly += fmt.Sprintf(" %#o,%d", immOneAcc.immU16, immOneAcc.acd)
		}
	case IO_FLAGS_DEV_FMT:
		var ioFlagsDev ioFlagsDevT
		ioFlagsDev.f = decodeIOFlags(memory.GetWbits(opcode, 8, 2))
		ioFlagsDev.ioDev = int(memory.GetWbits(opcode, 10, 6))
		decodedInstr.variant = ioFlagsDev
		if disassemble {
			decodedInstr.disassembly += fmt.Sprintf("%c %s",
				ioFlagsDev.f, deviceToString(devMap, ioFlagsDev.ioDev))
		}
	case IO_TEST_DEV_FMT:
		var ioTestDev ioTestDevT
		ioTestDev.t = int(memory.GetWbits(opcode, 8, 2))
		ioTestDev.ioDev = int(memory.GetWbits(opcode, 10, 6))
		decodedInstr.variant = ioTestDev
		if disassemble {
			decodedInstr.disassembly += fmt.Sprintf("%s %s", testToString(ioTestDev.t), deviceToString(devMap, ioTestDev.ioDev))
		}
	case LNDO_4_WORD_FMT:
		var lndo4Word lndo4WordT
		lndo4Word.acd = int(memory.GetWbits(opcode, 1, 2))
		lndo4Word.mode = int(memory.GetWbits(opcode, 3, 2))
		secondWord = memory.ReadWord(pc + 1)
		thirdWord = memory.ReadWord(pc + 2)
		fourthWord = memory.ReadWord(pc + 3)
		lndo4Word.ind = decodeIndirect(memory.TestWbit(secondWord, 0))
		lndo4Word.disp31 = decode31bitDisp(secondWord, thirdWord, lndo4Word.mode)
		lndo4Word.offsetU16 = uint16(fourthWord)
		decodedInstr.variant = lndo4Word
		if disassemble {
			decodedInstr.disassembly += fmt.Sprintf(" %d,%#o,%c%#o%s [4-Word OpCode]",
				lndo4Word.acd, lndo4Word.offsetU16, lndo4Word.ind, lndo4Word.disp31,
				modeToString(lndo4Word.mode))
		}
	case NOACC_MODE_2_WORD_FMT: // eg. XPEFB
		var noAccMode2Word noAccMode2WordT
		noAccMode2Word.mode = int(memory.GetWbits(opcode, 3, 2))
		noAccMode2Word.disp16 = uint16(memory.ReadWord(pc + 1))
		decodedInstr.variant = noAccMode2Word
		if disassemble {
			decodedInstr.disassembly += fmt.Sprintf(" %#o,%s [2-Word OpCode]",
				noAccMode2Word.disp16, modeToString(noAccMode2Word.mode))
		}
	case NOACC_MODE_3_WORD_FMT: // eg. LPEFB,
		var noAccMode3Word noAccMode3WordT
		noAccMode3Word.mode = int(memory.GetWbits(opcode, 3, 2))
		noAccMode3Word.immU32 = uint32(memory.ReadDWord(pc + 1))
		decodedInstr.variant = noAccMode3Word
		if disassemble {
			decodedInstr.disassembly += fmt.Sprintf(" %#o,%s [3-Word OpCode]",
				noAccMode3Word.immU32, modeToString(noAccMode3Word.mode))
		}
	case NOACC_MODE_IND_2_WORD_E_FMT, NOACC_MODE_IND_2_WORD_X_FMT:
		var noAccModeInd2Word noAccModeInd2WordT
		// if debugLogging {
		// 	logging.DebugPrint(logging.DebugLog, "X_FMT: Mnemonic is <%s>\n", decodedInstr.mnemonic)
		// }
		switch ix {
		case instrXJMP, instrXJSR, instrXNDSZ, instrXNISZ, instrXPEF, instrXPSHJ, instrXWDSZ:
			noAccModeInd2Word.mode = int(memory.GetWbits(opcode, 3, 2))
		case instrEDSZ, instrEISZ, instrEJMP, instrEJSR, instrPSHJ:
			noAccModeInd2Word.mode = int(memory.GetWbits(opcode, 6, 2))
		}
		secondWord = memory.ReadWord(pc + 1)
		noAccModeInd2Word.ind = decodeIndirect(memory.TestWbit(secondWord, 0))
		noAccModeInd2Word.disp15 = dg.WordT(decode15bitDisp(secondWord, noAccModeInd2Word.mode))
		decodedInstr.variant = noAccModeInd2Word
		if disassemble {
			decodedInstr.disassembly += fmt.Sprintf(" %c0%o%s [2-Word OpCode]",
				noAccModeInd2Word.ind, noAccModeInd2Word.disp15, modeToString(noAccModeInd2Word.mode))
		}
	case NOACC_MODE_IND_3_WORD_FMT: // eg. LJMP/LJSR, LNISZ, LNDSZ, LWDS
		var noAccModeInd3Word noAccModeInd3WordT
		noAccModeInd3Word.mode = int(memory.GetWbits(opcode, 3, 2))
		secondWord = memory.ReadWord(pc + 1)
		thirdWord = memory.ReadWord(pc + 2)
		noAccModeInd3Word.ind = decodeIndirect(memory.TestWbit(secondWord, 0))
		noAccModeInd3Word.disp31 = decode31bitDisp(secondWord, thirdWord, noAccModeInd3Word.mode)
		decodedInstr.variant = noAccModeInd3Word
		if disassemble {
			decodedInstr.disassembly += fmt.Sprintf(" %c0%o%s [3-Word OpCode]",
				noAccModeInd3Word.ind, noAccModeInd3Word.disp31, modeToString(noAccModeInd3Word.mode))
		}
	case NOACC_MODE_IND_3_WORD_XCALL_FMT: // XCALL
		var noAccModeInd3WordXcall noAccModeInd3WordXcallT
		noAccModeInd3WordXcall.mode = int(memory.GetWbits(opcode, 3, 2))
		secondWord = memory.ReadWord(pc + 1)
		thirdWord = memory.ReadWord(pc + 2)
		noAccModeInd3WordXcall.ind = decodeIndirect(memory.TestWbit(secondWord, 0))
		noAccModeInd3WordXcall.disp15 = decode15bitDisp(secondWord, noAccModeInd3WordXcall.mode)
		noAccModeInd3WordXcall.argCount = int(thirdWord)
		decodedInstr.variant = noAccModeInd3WordXcall
		if disassemble {
			decodedInstr.disassembly += fmt.Sprintf(" %c%#o%s, 0%o [3-Word OpCode]",
				noAccModeInd3WordXcall.ind, noAccModeInd3WordXcall.disp15,
				modeToString(noAccModeInd3WordXcall.mode), noAccModeInd3WordXcall.argCount)
		}
	case NOACC_MODE_IMM_IND_3_WORD_FMT: // eg. LNADI, LNSBI
		var noAccModeImmInd3Word noAccModeImmInd3WordT
		noAccModeImmInd3Word.immU16 = decode2bitImm(memory.GetWbits(opcode, 1, 2))
		noAccModeImmInd3Word.mode = int(memory.GetWbits(opcode, 3, 2))
		secondWord = memory.ReadWord(pc + 1)
		thirdWord = memory.ReadWord(pc + 2)
		noAccModeImmInd3Word.ind = decodeIndirect(memory.TestWbit(secondWord, 0))
		noAccModeImmInd3Word.disp31 = decode31bitDisp(secondWord, thirdWord, noAccModeImmInd3Word.mode)
		decodedInstr.variant = noAccModeImmInd3Word
		if disassemble {
			decodedInstr.disassembly += fmt.Sprintf(" %#o,%c%#o%s [3-Word OpCode]",
				noAccModeImmInd3Word.immU16, noAccModeImmInd3Word.ind, noAccModeImmInd3Word.disp31,
				modeToString(noAccModeImmInd3Word.mode))
		}
	case NOACC_MODE_IND_4_WORD_FMT: // eg. LCALL
		var noAccModeInd4Word noAccModeInd4WordT
		noAccModeInd4Word.mode = int(memory.GetWbits(opcode, 3, 2))
		secondWord = memory.ReadWord(pc + 1)
		thirdWord = memory.ReadWord(pc + 2)
		fourthWord = memory.ReadWord(pc + 3)
		noAccModeInd4Word.ind = decodeIndirect(memory.TestWbit(secondWord, 0))
		noAccModeInd4Word.disp31 = decode31bitDisp(secondWord, thirdWord, noAccModeInd4Word.mode)
		noAccModeInd4Word.argCount = int(fourthWord)
		decodedInstr.variant = noAccModeInd4Word
		if disassemble {
			decodedInstr.disassembly += fmt.Sprintf(" %c%#o%s,%#o [4-Word OpCode]",
				noAccModeInd4Word.ind, noAccModeInd4Word.disp31, modeToString(noAccModeInd4Word.mode),
				noAccModeInd4Word.argCount)
		}
	case NOVA_DATA_IO_FMT: // eg. DOA/B/C, DIA/B/C
		var novaDataIo novaDataIoT
		novaDataIo.acd = int(memory.GetWbits(opcode, 3, 2))
		novaDataIo.f = decodeIOFlags(memory.GetWbits(opcode, 8, 2))
		novaDataIo.ioDev = int(memory.GetWbits(opcode, 10, 6))
		decodedInstr.variant = novaDataIo
		if disassemble {
			decodedInstr.disassembly += fmt.Sprintf("%c %d,%s",
				novaDataIo.f, novaDataIo.acd, deviceToString(devMap, novaDataIo.ioDev))
		}
	case NOVA_NOACC_EFF_ADDR_FMT: // eg. DSZ, ISZ, JMP, JSR
		var novaNoAccEffAddr novaNoAccEffAddrT
		novaNoAccEffAddr.ind = decodeIndirect(memory.TestWbit(opcode, 5))
		novaNoAccEffAddr.mode = int(memory.GetWbits(opcode, 6, 2))
		novaNoAccEffAddr.disp15 = decode8bitDisp(dg.ByteT(opcode&0x00ff), novaNoAccEffAddr.mode) // NB
		decodedInstr.variant = novaNoAccEffAddr
		if disassemble {
			decodedInstr.disassembly += fmt.Sprintf(" %c%#o%s",
				novaNoAccEffAddr.ind, novaNoAccEffAddr.disp15, modeToString(novaNoAccEffAddr.mode))
		}
	case NOVA_ONEACC_EFF_ADDR_FMT:
		var novaOneAccEffAddr novaOneAccEffAddrT
		novaOneAccEffAddr.acd = int(memory.GetWbits(opcode, 3, 2))
		novaOneAccEffAddr.ind = decodeIndirect(memory.TestWbit(opcode, 5))
		novaOneAccEffAddr.mode = int(memory.GetWbits(opcode, 6, 2))
		novaOneAccEffAddr.disp15 = decode8bitDisp(dg.ByteT(opcode&0x00ff), novaOneAccEffAddr.mode) // NB
		decodedInstr.variant = novaOneAccEffAddr
		if disassemble {
			decodedInstr.disassembly += fmt.Sprintf(" %d,%c%#o%s",
				novaOneAccEffAddr.acd, novaOneAccEffAddr.ind, novaOneAccEffAddr.disp15,
				modeToString(novaOneAccEffAddr.mode))
		}
	case NOVA_TWOACC_MULT_OP_FMT: // eg. ADC, ADD, AND, COM
		var novaTwoAccMultOp novaTwoAccMultOpT
		novaTwoAccMultOp.acs = int(memory.GetWbits(opcode, 1, 2))
		novaTwoAccMultOp.acd = int(memory.GetWbits(opcode, 3, 2))
		novaTwoAccMultOp.sh = decodeShift(memory.GetWbits(opcode, 8, 2))
		novaTwoAccMultOp.c = decodeCarry(memory.GetWbits(opcode, 10, 2))
		novaTwoAccMultOp.nl = decodeNoLoad(memory.TestWbit(opcode, 12))
		novaTwoAccMultOp.skip = int(memory.GetWbits(opcode, 13, 3))
		decodedInstr.variant = novaTwoAccMultOp
		if disassemble {
			decodedInstr.disassembly += fmt.Sprintf("%c%c%c %d,%d %s",
				novaTwoAccMultOp.c, novaTwoAccMultOp.sh, novaTwoAccMultOp.nl, novaTwoAccMultOp.acs,
				novaTwoAccMultOp.acd, skipToString(novaTwoAccMultOp.skip))
		}
	case ONEACC_1_WORD_FMT: // eg. CVWN, HLV, LDAFP
		var oneAcc1Word oneAcc1WordT
		oneAcc1Word.acd = int(memory.GetWbits(opcode, 3, 2))
		decodedInstr.variant = oneAcc1Word
		if disassemble {
			decodedInstr.disassembly += fmt.Sprintf(" %d", oneAcc1Word.acd)
		}
	case ONEACC_IMM_2_WORD_FMT: // eg. ADDI, NADDI, NLDAI, , WSEQI, WLSHI, WNADI
		var oneAccImm2Word oneAccImm2WordT
		oneAccImm2Word.acd = int(memory.GetWbits(opcode, 3, 2))
		secondWord = memory.ReadWord(pc + 1)
		oneAccImm2Word.immS16 = int16(secondWord)
		decodedInstr.variant = oneAccImm2Word
		if disassemble {
			decodedInstr.disassembly += fmt.Sprintf(" %#o,%d [2-Word OpCode]", oneAccImm2Word.immS16, oneAccImm2Word.acd)
		}
	case ONEACC_IMMWD_2_WORD_FMT: // eg. ANDI, IORI
		var oneAccImmWd2Word oneAccImmWd2WordT
		oneAccImmWd2Word.acd = int(memory.GetWbits(opcode, 3, 2))
		secondWord = memory.ReadWord(pc + 1)
		oneAccImmWd2Word.immWord = secondWord
		decodedInstr.variant = oneAccImmWd2Word
		if disassemble {
			decodedInstr.disassembly += fmt.Sprintf(" %#o,%d [2-Word OpCode]", oneAccImmWd2Word.immWord, oneAccImmWd2Word.acd)
		}
	case ONEACC_IMM_3_WORD_FMT: // eg. WADDI, WUGTI
		var oneAccImm3Word oneAccImm3WordT
		oneAccImm3Word.acd = int(memory.GetWbits(opcode, 3, 2))
		oneAccImm3Word.immU32 = uint32(memory.ReadDWord(pc + 1))
		decodedInstr.variant = oneAccImm3Word
		if disassemble {
			decodedInstr.disassembly += fmt.Sprintf(" %#o,%d [3-Word OpCode]", oneAccImm3Word.immU32, oneAccImm3Word.acd)
		}
	case ONEACC_IMMDWD_3_WORD_FMT: // eg. WANDI, WIORI, WLDAI
		var oneAccImmDwd3Word oneAccImmDwd3WordT
		oneAccImmDwd3Word.acd = int(memory.GetWbits(opcode, 3, 2))
		oneAccImmDwd3Word.immDword = memory.ReadDWord(pc + 1)
		decodedInstr.variant = oneAccImmDwd3Word
		if disassemble {
			decodedInstr.disassembly += fmt.Sprintf(" %#o,%d [3-Word OpCode]", oneAccImmDwd3Word.immDword, oneAccImmDwd3Word.acd)
		}
	case ONEACC_MODE_2_WORD_X_B_FMT: // eg. XLDB, XLEFB, XSTB
		var oneAccMode2Word oneAccMode2WordT
		oneAccMode2Word.mode = int(memory.GetWbits(opcode, 1, 2))
		oneAccMode2Word.acd = int(memory.GetWbits(opcode, 3, 2))
		secondWord = memory.ReadWord(pc + 1)
		oneAccMode2Word.disp16, oneAccMode2Word.bitLow = decode16bitByteDisp(secondWord)
		decodedInstr.variant = oneAccMode2Word
		if disassemble {
			decodedInstr.disassembly += fmt.Sprintf(" %d,%#o+%c%s [2-Word OpCode]",
				oneAccMode2Word.acd, oneAccMode2Word.disp16*2, loHiToByte(oneAccMode2Word.bitLow), modeToString(oneAccMode2Word.mode))
		}
	case ONEACC_MODE_2_WORD_E_FMT: // eg. ELDB, ESTB
		var oneAccMode2Word oneAccMode2WordT
		oneAccMode2Word.mode = int(memory.GetWbits(opcode, 6, 2))
		oneAccMode2Word.acd = int(memory.GetWbits(opcode, 3, 2))
		secondWord = memory.ReadWord(pc + 1)
		oneAccMode2Word.disp16, oneAccMode2Word.bitLow = decode16bitByteDisp(secondWord)
		decodedInstr.variant = oneAccMode2Word
		if disassemble {
			decodedInstr.disassembly += fmt.Sprintf(" %d,%#o+%c%s [2-Word OpCode]",
				oneAccMode2Word.acd, oneAccMode2Word.disp16*2, loHiToByte(oneAccMode2Word.bitLow), modeToString(oneAccMode2Word.mode))
		}
	case ONEACC_MODE_3_WORD_FMT: // eg. LLDB, LLEFB
		var oneAccMode3Word oneAccMode3WordT
		oneAccMode3Word.mode = int(memory.GetWbits(opcode, 1, 2))
		oneAccMode3Word.acd = int(memory.GetWbits(opcode, 3, 2))
		secondWord = memory.ReadWord(pc + 1)
		thirdWord = memory.ReadWord(pc + 2)
		oneAccMode3Word.disp31 = decode31bitDisp(secondWord, thirdWord, oneAccMode3Word.mode)
		decodedInstr.variant = oneAccMode3Word
		if disassemble {
			decodedInstr.disassembly += fmt.Sprintf(" %d,%#o%s [3-Word OpCode]",
				oneAccMode3Word.acd, oneAccMode3Word.disp31, modeToString(oneAccMode3Word.mode))
		}
	case ONEACC_MODE_IND_2_WORD_E_FMT: // eg. DSPA, ELDA, ELDB, ELEF, ESTA
		var oneAccModeInd2Word oneAccModeInd2WordT
		oneAccModeInd2Word.mode = int(memory.GetWbits(opcode, 6, 2))
		oneAccModeInd2Word.acd = int(memory.GetWbits(opcode, 3, 2))
		secondWord = memory.ReadWord(pc + 1)
		oneAccModeInd2Word.ind = decodeIndirect(memory.TestWbit(secondWord, 0))
		oneAccModeInd2Word.disp15 = decode15bitDisp(secondWord, oneAccModeInd2Word.mode)
		decodedInstr.variant = oneAccModeInd2Word
		if disassemble {
			decodedInstr.disassembly += fmt.Sprintf(" %d,%c%#o%s [2-Word OpCode]",
				oneAccModeInd2Word.acd, oneAccModeInd2Word.ind, oneAccModeInd2Word.disp15,
				modeToString(oneAccModeInd2Word.mode))
		}
	case ONEACC_MODE_IND_2_WORD_X_FMT: // eg. XNADD/SUB, XNLDA/XWSTA, XLEF
		var oneAccModeInd2Word oneAccModeInd2WordT
		oneAccModeInd2Word.mode = int(memory.GetWbits(opcode, 1, 2))
		oneAccModeInd2Word.acd = int(memory.GetWbits(opcode, 3, 2))
		secondWord = memory.ReadWord(pc + 1)
		oneAccModeInd2Word.ind = decodeIndirect(memory.TestWbit(secondWord, 0))
		oneAccModeInd2Word.disp15 = decode15bitDisp(secondWord, oneAccModeInd2Word.mode)
		decodedInstr.variant = oneAccModeInd2Word
		if disassemble {
			decodedInstr.disassembly += fmt.Sprintf(" %d,%c%#o%s [2-Word OpCode]",
				oneAccModeInd2Word.acd, oneAccModeInd2Word.ind, oneAccModeInd2Word.disp15, modeToString(oneAccModeInd2Word.mode))
		}
	case ONEACC_MODE_IND_3_WORD_FMT: // eg. LDSP, LLEF, LNADD/SUB LNDIV, LNLDA/STA, LNMUL, LWLDA/LWSTA,LNLDA
		var oneAccModeInd3Word oneAccModeInd3WordT
		oneAccModeInd3Word.mode = int(memory.GetWbits(opcode, 1, 2))
		oneAccModeInd3Word.acd = int(memory.GetWbits(opcode, 3, 2))
		secondWord = memory.ReadWord(pc + 1)
		oneAccModeInd3Word.ind = decodeIndirect(memory.TestWbit(secondWord, 0))
		thirdWord = memory.ReadWord(pc + 2)
		oneAccModeInd3Word.disp31 = decode31bitDisp(secondWord, thirdWord, oneAccModeInd3Word.mode)
		decodedInstr.variant = oneAccModeInd3Word
		if disassemble {
			decodedInstr.disassembly += fmt.Sprintf(" %d,%c%#o%s [3-Word OpCode]",
				oneAccModeInd3Word.acd, oneAccModeInd3Word.ind, oneAccModeInd3Word.disp31, modeToString(oneAccModeInd3Word.mode))
		}
	case TWOACC_1_WORD_FMT: // eg. ANC, BTO, WSUB and MANY others
		var twoAcc1Word twoAcc1WordT
		twoAcc1Word.acs = int(memory.GetWbits(opcode, 1, 2))
		twoAcc1Word.acd = int(memory.GetWbits(opcode, 3, 2))
		decodedInstr.variant = twoAcc1Word
		if disassemble {
			decodedInstr.disassembly += fmt.Sprintf(" %d,%d", twoAcc1Word.acs, twoAcc1Word.acd)
		}
	case SPLIT_8BIT_DISP_FMT: // eg. WBR, always a signed disp
		var split8bitDisp split8bitDispT
		tmp8bit := dg.ByteT(memory.GetWbits(opcode, 1, 4) & 0xff)
		tmp8bit = tmp8bit << 4
		tmp8bit |= dg.ByteT(memory.GetWbits(opcode, 6, 4) & 0xff)
		split8bitDisp.disp8 = int8(decode8bitDisp(tmp8bit, pcMode))
		decodedInstr.variant = split8bitDisp
		if disassemble {
			decodedInstr.disassembly += fmt.Sprintf(" %#o", int32(split8bitDisp.disp8))
		}
	case THREE_WORD_DO_FMT: // eg. XNDO
		var threeWordDo threeWordDoT
		threeWordDo.acd = int(memory.GetWbits(opcode, 1, 2))
		threeWordDo.mode = int(memory.GetWbits(opcode, 3, 2))
		secondWord = memory.ReadWord(pc + 1)
		threeWordDo.ind = decodeIndirect(memory.TestWbit(secondWord, 0))
		threeWordDo.disp15 = decode15bitDisp(secondWord, threeWordDo.mode)
		thirdWord = memory.ReadWord(pc + 2)
		threeWordDo.offsetU16 = uint16(thirdWord)
		decodedInstr.variant = threeWordDo
		if disassemble {
			decodedInstr.disassembly += fmt.Sprintf(" %d,%#o %c%#o%s [3-Word OpCode]",
				threeWordDo.acd, threeWordDo.offsetU16, threeWordDo.ind, threeWordDo.disp15,
				modeToString(threeWordDo.mode))
		}
	case TWOACC_IMM_2_WORD_FMT: // eg. CIOI
		var twoAccImm2Word twoAccImm2WordT
		twoAccImm2Word.acs = int(memory.GetWbits(opcode, 1, 2))
		twoAccImm2Word.acd = int(memory.GetWbits(opcode, 3, 2))
		twoAccImm2Word.immWord = memory.ReadWord(pc + 1)
		decodedInstr.variant = twoAccImm2Word
		if disassemble {
			decodedInstr.disassembly += fmt.Sprintf(" %#o,%d,%d", twoAccImm2Word.immWord, twoAccImm2Word.acs,
				twoAccImm2Word.acd)
		}
	case UNIQUE_1_WORD_FMT:
		// nothing to do in this case, no associated variant

	case UNIQUE_2_WORD_FMT: // eg.SAVE, WSAVR, WSAVS
		var unique2Word unique2WordT
		unique2Word.immU16 = uint16(memory.ReadWord(pc + 1))
		decodedInstr.variant = unique2Word
		if disassemble {
			decodedInstr.disassembly += fmt.Sprintf(" %#o [2-Word OpCode]", unique2Word.immU16)
		}
	case WSKB_FMT: // eg. WSKBO/Z
		var wskb wskbT
		tmp8bit := dg.ByteT(memory.GetWbits(opcode, 1, 3) & 0xff)
		tmp8bit = tmp8bit << 2
		tmp8bit |= dg.ByteT(memory.GetWbits(opcode, 10, 2) & 0xff)
		wskb.bitNum = int(uint8(tmp8bit))
		decodedInstr.variant = wskb
		if disassemble {
			decodedInstr.disassembly += fmt.Sprintf(" %#o", wskb.bitNum)
		}
	default:
		logging.DebugPrint(logging.DebugLog, "ERROR: Invalid instruction BB format (%d) for instruction <%s>\n",
			decodedInstr.instrFmt, decodedInstr.mnemonic)
		return nil, false
	}

	return &decodedInstr, true
}

/* decoders for (parts of) operands below here... */

func decode2bitImm(i dg.WordT) uint16 {
	// to expand range (by 1!) 1 is subtracted from operand
	return uint16(i + 1)
}

// Decode8BitDisp must return signed 16-bit as the result could be
// either 8-bit signed or 8-bit unsigned
func decode8bitDisp(d8 dg.ByteT, mode int) (disp16 int16) {
	if mode == absoluteMode {
		disp16 = int16(d8) & 0x00ff // unsigned offset
	} else {
		// signed offset...
		disp16 = int16(int8(d8)) // this should sign-extend
	}
	return disp16
}

func decode15bitDisp(d15 dg.WordT, mode int) (disp16 int16) {
	if mode == absoluteMode {
		disp16 = int16(d15 & 0x7fff) // zero extend
	} else {
		if memory.TestWbit(d15, 1) {
			disp16 = int16(d15 | 0x8000) // sign extend
		} else {
			disp16 = int16(d15 & 0x7fff) // zero extend
		}
	}
	if debugLogging {
		logging.DebugPrint(logging.DebugLog, "... decode15bitDisp got: %#o, returning: %#o\n", d15, disp16)
	}
	return disp16
}

func decode16bitByteDisp(d16 dg.WordT) (disp16 int16, loHi bool) {
	loHi = memory.TestWbit(d16, 15)
	disp16 = int16(d16 >> 1)
	if debugLogging {
		logging.DebugPrint(logging.DebugLog, "... decode16bitByteDisp got: %#o, returning %#o\n", d16, disp16)
	}
	return disp16, loHi
}

func decode31bitDisp(d1, d2 dg.WordT, mode int) int32 {
	// FIXME Test this!
	dwd := memory.DwordFromTwoWords(d1&0x7fff, d2)
	// sign-extend if not absolute mode
	if mode != absoluteMode && memory.TestDwbit(dwd, 1) {
		memory.SetDwbit(&dwd, 0)
	}
	return int32(dwd)
}

func decodeCarry(cry dg.WordT) byte {
	switch cry {
	case 0:
		return ' '
	case 1:
		return 'Z'
	case 2:
		return 'O' // Letter 'O' for One
	case 3:
		return 'C'
	}
	return '*'
}

func decodeIndirect(i bool) byte {
	if i {
		return '@'
	}
	return ' '
}

func decodeIOFlags(fl dg.WordT) byte {
	return ioFlags[fl]
}

func decodeNoLoad(n bool) byte {
	if n {
		return '#'
	}
	return ' '
}

func decodeShift(sh dg.WordT) byte {
	switch sh {
	case 0:
		return ' '
	case 1:
		return 'L'
	case 2:
		return 'R'
	case 3:
		return 'S'
	}
	return '*'
}

// deviceToString is used for disassembly
// If the device code is known, then an assember mnemonic is returned, otherwise just the code
func deviceToString(deviceMap devices.DeviceMapT, devNum int) string {
	de, known := deviceMap[devNum]
	if known {
		return de.DgMnemonic
	}
	return fmt.Sprintf("%#o", devNum)
}

func loHiToByte(loHi bool) byte {
	if loHi {
		return 'H'
	}
	return 'L'
}

func modeToString(mode int) string {
	var modes = [...]string{"", "PC", "AC2", "AC3"}
	if mode == absoluteMode {
		return ""
	}
	return "," + modes[mode]
}

func skipToString(s int) string {
	var skips = [...]string{"NONE", "SKP", "SZC", "SNC", "SZR", "SNR", "SEZ", "SBN"}
	if s == 0 {
		return ""
	}
	return skips[s]
}

func testToString(t int) string {
	var ioTests = [...]string{"BN", "BZ", "DN", "DZ"}
	return ioTests[t]
}

func (decoded *decodedInstrT) GetDisassembly() string {
	return decoded.disassembly
}

func (decoded *decodedInstrT) GetLength() int {
	return decoded.instrLength
}
