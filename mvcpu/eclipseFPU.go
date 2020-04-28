// Copyright ©2020  Steve Merrony

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
	"math"

	"github.com/SMerrony/dgemug/dg"
	"github.com/SMerrony/dgemug/memory"
)

func eclipseFPU(cpu *CPUT, iPtr *decodedInstrT) bool {
	switch iPtr.ix {

	case instrFAB:
		cpu.fpac[iPtr.ac] = math.Abs(cpu.fpac[iPtr.ac])
		cpu.SetN(false)
		cpu.SetZ(cpu.fpac[iPtr.ac] == 0.0)

	case instrFAD:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		cpu.fpac[twoAcc1Word.acd] += cpu.fpac[twoAcc1Word.acs]
		cpu.SetZ(cpu.fpac[twoAcc1Word.acd] == 0.0)
		cpu.SetN(cpu.fpac[twoAcc1Word.acd] < 0.0)

	case instrFAS:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		cpu.fpac[twoAcc1Word.acd] += cpu.fpac[twoAcc1Word.acs]
		cpu.SetZ(cpu.fpac[twoAcc1Word.acd] == 0.0)
		cpu.SetN(cpu.fpac[twoAcc1Word.acd] < 0.0)

	case instrFCLE:
		cpu.fpsr = 0 // TODO check - PoP contradicts itself

	case instrFCMP:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		switch {
		case cpu.fpac[twoAcc1Word.acs] == cpu.fpac[twoAcc1Word.acd]:
			memory.ClearQwbit(&cpu.fpsr, fpsrN)
			memory.SetQwbit(&cpu.fpsr, fpsrZ)
		case cpu.fpac[twoAcc1Word.acs] > cpu.fpac[twoAcc1Word.acd]:
			memory.SetQwbit(&cpu.fpsr, fpsrN)
			memory.ClearQwbit(&cpu.fpsr, fpsrZ)
		case cpu.fpac[twoAcc1Word.acs] < cpu.fpac[twoAcc1Word.acd]:
			memory.ClearQwbit(&cpu.fpsr, fpsrN)
			memory.ClearQwbit(&cpu.fpsr, fpsrZ)
		}

	case instrFDS:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		cpu.fpac[twoAcc1Word.acd] /= cpu.fpac[twoAcc1Word.acs]
		cpu.SetZ(cpu.fpac[twoAcc1Word.acd] == 0.0)
		cpu.SetN(cpu.fpac[twoAcc1Word.acd] < 0.0)

	case instrFEXP:
		qwd := memory.Float64toDGdouble(cpu.fpac[iPtr.ac])
		qwd &= 0x80FF_FFFF_FFFF_FFFF
		qwd |= dg.QwordT(cpu.ac[0]&0x0000_7f00) << 48
		cpu.fpac[iPtr.ac] = memory.DGdoubleToFloat64(qwd)
		cpu.SetZ(cpu.fpac[iPtr.ac] == 0.0)
		cpu.SetN(cpu.fpac[iPtr.ac] < 0.0)

	case instrFFAS:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT) // N.B. Not the usual AC order
		s32 := int32(cpu.fpac[twoAcc1Word.acd])
		if s32 < minNegS16 || s32 > maxPosS16 {
			memory.SetQwbit(&cpu.fpsr, fpsrMof)
		}
		cpu.ac[twoAcc1Word.acs] = dg.DwordT(s32)

	case instrFHLV:
		cpu.fpac[iPtr.ac] /= 2.0
		cpu.SetZ(cpu.fpac[iPtr.ac] == 0.0)
		cpu.SetN(cpu.fpac[iPtr.ac] < 0.0)

	case instrFINT:
		cpu.fpac[iPtr.ac] = math.Trunc(cpu.fpac[iPtr.ac])
		cpu.SetZ(cpu.fpac[iPtr.ac] == 0.0)
		cpu.SetN(cpu.fpac[iPtr.ac] < 0.0)

	case instrFLAS:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		cpu.fpac[twoAcc1Word.acd] = float64(int16(twoAcc1Word.acs))
		cpu.SetZ(cpu.fpac[twoAcc1Word.acd] == 0.0)
		cpu.SetN(cpu.fpac[twoAcc1Word.acd] < 0.0)

	case instrFLDS:
		oneAccModeInd2Word := iPtr.variant.(oneAccModeInd2WordT)
		addr := resolve15bitDisplacement(cpu, oneAccModeInd2Word.ind, oneAccModeInd2Word.mode, dg.WordT(oneAccModeInd2Word.disp15), iPtr.dispOffset)
		addr &= 0x7fff
		addr |= (cpu.pc & ringMask32)
		cpu.fpac[oneAccModeInd2Word.acd] = memory.DGsingleToFloat64(memory.ReadDWord(addr))
		cpu.SetZ(cpu.fpac[oneAccModeInd2Word.acd] == 0.0)
		cpu.SetN(cpu.fpac[oneAccModeInd2Word.acd] < 0.0)

	case instrFMD:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		cpu.fpac[twoAcc1Word.acd] *= cpu.fpac[twoAcc1Word.acs]
		cpu.SetZ(cpu.fpac[twoAcc1Word.acd] == 0.0)
		cpu.SetN(cpu.fpac[twoAcc1Word.acd] < 0.0)

	case instrFMOV:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		cpu.fpac[twoAcc1Word.acd] = cpu.fpac[twoAcc1Word.acs]
		cpu.SetZ(cpu.fpac[twoAcc1Word.acd] == 0.0)
		cpu.SetN(cpu.fpac[twoAcc1Word.acd] < 0.0)

	case instrFMS:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		cpu.fpac[twoAcc1Word.acd] *= cpu.fpac[twoAcc1Word.acs]
		cpu.SetZ(cpu.fpac[twoAcc1Word.acd] == 0.0)
		cpu.SetN(cpu.fpac[twoAcc1Word.acd] < 0.0)

	case instrFNEG:
		cpu.fpac[iPtr.ac] = -cpu.fpac[iPtr.ac]
		cpu.SetZ(cpu.fpac[iPtr.ac] == 0.0)
		cpu.SetN(cpu.fpac[iPtr.ac] < 0.0)

	case instrFRDS:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		if memory.TestQwbit(cpu.fpsr, 8) { // FIXME the first assignment is not right...
			cpu.fpac[twoAcc1Word.acd] = math.Float64frombits(math.Float64bits(cpu.fpac[twoAcc1Word.acs]) & 0xffff_ffff_0000_0000)
		} else {
			cpu.fpac[twoAcc1Word.acd] = math.Float64frombits(math.Float64bits(cpu.fpac[twoAcc1Word.acs]) & 0xffff_ffff_0000_0000)
		}
		cpu.SetZ(cpu.fpac[twoAcc1Word.acd] == 0.0)
		cpu.SetN(cpu.fpac[twoAcc1Word.acd] < 0.0)

	case instrFRH:
		cpu.ac[0] = memory.Float64toDGsingle(cpu.fpac[iPtr.ac]) >> 16

	case instrFSD:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		cpu.fpac[twoAcc1Word.acd] -= cpu.fpac[twoAcc1Word.acs]
		cpu.SetZ(cpu.fpac[twoAcc1Word.acd] == 0.0)
		cpu.SetN(cpu.fpac[twoAcc1Word.acd] < 0.0)

	case instrFSGE:
		if !memory.TestQwbit(cpu.fpsr, fpsrN) {
			cpu.pc++
		}

	case instrFSGT:
		if !memory.TestQwbit(cpu.fpsr, fpsrZ) && !memory.TestQwbit(cpu.fpsr, fpsrN) {
			cpu.pc++
		}

	case instrFSLE:
		if memory.TestQwbit(cpu.fpsr, fpsrZ) || memory.TestQwbit(cpu.fpsr, fpsrN) {
			cpu.pc++
		}

	case instrFSLT:
		if memory.TestQwbit(cpu.fpsr, fpsrN) {
			cpu.pc++
		}

	case instrFSNE:
		if !memory.TestQwbit(cpu.fpsr, fpsrZ) {
			cpu.pc++
		}

	case instrFSS:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		cpu.fpac[twoAcc1Word.acd] -= cpu.fpac[twoAcc1Word.acs]
		cpu.SetZ(cpu.fpac[twoAcc1Word.acd] == 0.0)
		cpu.SetN(cpu.fpac[twoAcc1Word.acd] < 0.0)

	case instrFSST:
		addr := resolve15bitDisplacement(cpu, iPtr.ind, iPtr.mode, dg.WordT(iPtr.disp15), iPtr.dispOffset)
		addr &= 0x7fff
		addr |= (cpu.pc & ringMask32)
		memory.WriteWord(addr, dg.WordT(cpu.fpsr>>48))
		memory.WriteWord(addr+1, dg.WordT(cpu.fpsr)) // last word of FPSR

	case instrFSTS:
		oneAccModeInd2Word := iPtr.variant.(oneAccModeInd2WordT)
		addr := resolve15bitDisplacement(cpu, oneAccModeInd2Word.ind, oneAccModeInd2Word.mode, dg.WordT(oneAccModeInd2Word.disp15), iPtr.dispOffset)
		addr &= 0x7fff
		addr |= (cpu.pc & ringMask32)
		memory.WriteDWord(addr, memory.Float64toDGsingle(cpu.fpac[oneAccModeInd2Word.acd]))

	case instrFTD:
		memory.ClearQwbit(&cpu.fpsr, fpsrTe)

	case instrFTE:
		memory.SetQwbit(&cpu.fpsr, fpsrTe)

	default:
		log.Fatalf("ERROR: ECLIPSE_FPU instruction <%s> not yet implemented\n", iPtr.mnemonic)
		return false
	}

	cpu.pc += dg.PhysAddrT(iPtr.instrLength)
	return true
}
