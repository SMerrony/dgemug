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

	"github.com/SMerrony/dgemug/dg"
	"github.com/SMerrony/dgemug/memory"
)

func eclipseFPU(cpu *CPUT, iPtr *decodedInstrT) bool {
	switch iPtr.ix {

	case instrFCLE:
		cpu.fpsr = 0 // TODO check - PoP contradicts itself

	case instrFLAS:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		cpu.fpac[twoAcc1Word.acd] = float64(int16(twoAcc1Word.acs)) // TODO not quite right...

	case instrFLDS:
		oneAccModeInd2Word := iPtr.variant.(oneAccModeInd2WordT)
		addr := resolve15bitDisplacement(cpu, oneAccModeInd2Word.ind, oneAccModeInd2Word.mode, dg.WordT(oneAccModeInd2Word.disp15), iPtr.dispOffset)
		addr &= 0x7fff
		addr |= (cpu.pc & ringMask32)
		cpu.fpac[oneAccModeInd2Word.acd] = float64(memory.ReadDWord(addr))

	case instrFNEG:
		cpu.fpac[iPtr.ac] = -cpu.fpac[iPtr.ac]

	case instrFSTS:
		oneAccModeInd2Word := iPtr.variant.(oneAccModeInd2WordT)
		addr := resolve15bitDisplacement(cpu, oneAccModeInd2Word.ind, oneAccModeInd2Word.mode, dg.WordT(oneAccModeInd2Word.disp15), iPtr.dispOffset)
		addr &= 0x7fff
		addr |= (cpu.pc & ringMask32)
		memory.WriteDWord(addr, dg.DwordT(cpu.fpac[oneAccModeInd2Word.acd]))

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
