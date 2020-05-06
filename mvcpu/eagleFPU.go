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
	"fmt"
	"log"
	"math"
	"strconv"

	"github.com/SMerrony/dgemug/dg"
	"github.com/SMerrony/dgemug/logging"
	"github.com/SMerrony/dgemug/memory"
)

func eagleFPU(cpu *CPUT, iPtr *decodedInstrT) bool {
	switch iPtr.ix {

	case instrLFAMD, instrLFDMD, instrLFMMD, instrLFSMD:
		oneAccModeInd3Word := iPtr.variant.(oneAccModeInd3WordT)
		addr := resolve31bitDisplacement(cpu, oneAccModeInd3Word.ind, oneAccModeInd3Word.mode, oneAccModeInd3Word.disp31, iPtr.dispOffset)
		qwd := dg.QwordT(memory.ReadDWord(addr))<<32 | dg.QwordT(memory.ReadDWord(addr+2))
		switch iPtr.ix {
		case instrLFAMD:
			cpu.fpac[oneAccModeInd3Word.acd] += memory.DGdoubleToFloat64(qwd)
		case instrLFDMD:
			cpu.fpac[oneAccModeInd3Word.acd] /= memory.DGdoubleToFloat64(qwd)
		case instrLFMMD:
			cpu.fpac[oneAccModeInd3Word.acd] *= memory.DGdoubleToFloat64(qwd)
		case instrLFSMD:
			cpu.fpac[oneAccModeInd3Word.acd] -= memory.DGdoubleToFloat64(qwd)
		}
		cpu.SetZ(cpu.fpac[oneAccModeInd3Word.acd] == 0.0)
		cpu.SetN(cpu.fpac[oneAccModeInd3Word.acd] < 0.0)

	case instrLFLDD:
		oneAccModeInd3Word := iPtr.variant.(oneAccModeInd3WordT)
		addr := resolve31bitDisplacement(cpu, oneAccModeInd3Word.ind, oneAccModeInd3Word.mode, oneAccModeInd3Word.disp31, iPtr.dispOffset)
		qwd := dg.QwordT(memory.ReadDWord(addr))<<32 | dg.QwordT(memory.ReadDWord(addr+2))
		cpu.fpac[oneAccModeInd3Word.acd] = memory.DGdoubleToFloat64(qwd)
		cpu.SetZ(cpu.fpac[oneAccModeInd3Word.acd] == 0.0)
		cpu.SetN(cpu.fpac[oneAccModeInd3Word.acd] < 0.0)

	case instrLFLDS:
		oneAccModeInd3Word := iPtr.variant.(oneAccModeInd3WordT)
		addr := resolve31bitDisplacement(cpu, oneAccModeInd3Word.ind, oneAccModeInd3Word.mode, oneAccModeInd3Word.disp31, iPtr.dispOffset)
		dwd := memory.ReadDWord(addr)
		cpu.fpac[oneAccModeInd3Word.acd] = memory.DGsingleToFloat64(dwd)
		cpu.SetZ(cpu.fpac[oneAccModeInd3Word.acd] == 0.0)
		cpu.SetN(cpu.fpac[oneAccModeInd3Word.acd] < 0.0)

	case instrLFMMS:
		oneAccModeInd3Word := iPtr.variant.(oneAccModeInd3WordT)
		addr := resolve31bitDisplacement(cpu, oneAccModeInd3Word.ind, oneAccModeInd3Word.mode, oneAccModeInd3Word.disp31, iPtr.dispOffset)
		single := memory.DGsingleToFloat64(memory.ReadDWord(addr))
		cpu.fpac[oneAccModeInd3Word.acd] *= single
		cpu.SetZ(cpu.fpac[oneAccModeInd3Word.acd] == 0.0)
		cpu.SetN(cpu.fpac[oneAccModeInd3Word.acd] < 0.0)

	case instrLFSTD:
		oneAccModeInd3Word := iPtr.variant.(oneAccModeInd3WordT)
		addr := resolve31bitDisplacement(cpu, oneAccModeInd3Word.ind, oneAccModeInd3Word.mode, oneAccModeInd3Word.disp31, iPtr.dispOffset)
		qwd := memory.Float64toDGdouble(cpu.fpac[oneAccModeInd3Word.acd])
		memory.WriteDWord(addr, dg.DwordT(qwd>>32))
		memory.WriteDWord(addr+2, dg.DwordT(qwd))

	case instrLFSTS:
		oneAccModeInd3Word := iPtr.variant.(oneAccModeInd3WordT)
		addr := resolve31bitDisplacement(cpu, oneAccModeInd3Word.ind, oneAccModeInd3Word.mode, oneAccModeInd3Word.disp31, iPtr.dispOffset)
		memory.WriteDWord(addr, memory.Float64toDGsingle(cpu.fpac[oneAccModeInd3Word.acd]))

	case instrWFFAD:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		cpu.ac[twoAcc1Word.acs] = dg.DwordT(int32(math.Round(cpu.fpac[twoAcc1Word.acd])))

	case instrWFLAD:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		cpu.fpac[twoAcc1Word.acd] = float64(int32(cpu.ac[twoAcc1Word.acs])) // N.B INT32 conversion required!!!
		cpu.SetZ(cpu.fpac[twoAcc1Word.acd] == 0.0)
		cpu.SetN(cpu.fpac[twoAcc1Word.acd] < 0.0)

	case instrWLDI:
		_, dataType, size := memory.DecodeDecDataType(cpu.ac[1]) // "WLDI does not use the scale factor..."
		cpu.ac[2] = cpu.ac[3]
		//bs := memory.ReadNBytes(dg.PhysAddrT(cpu.ac[3]), size)
		bs := memory.ReadDec(dg.PhysAddrT(cpu.ac[3]), size)
		switch dataType {
		// case memory.UnpackedDecTSC:
		// case memory.UnpackedDecLSC:
		// case memory.UnpackedDecTS:
		// case memory.UnpackedDecLS:
		case memory.UnpackedDecU:
			//intNum, err := strconv.Atoi(strings.TrimSpace(string(bs)))
			intNum, err := strconv.Atoi(bs)
			if err != nil {
				log.Panicf("ERROR: Could not parse Decimal <%v> as integer\n", bs)
			}
			logging.DebugPrint(logging.DebugLog, "... decoded %d from %s\n", intNum, bs)
			cpu.fpac[iPtr.ac] = float64(intNum)
		default:
			log.Panicf("ERROR: Packed data types not yet implemented, type is: %d.\n", dataType)
		}
		cpu.SetZ(cpu.fpac[iPtr.ac] == 0.0)
		cpu.SetN(cpu.fpac[iPtr.ac] < 0.0)

	case instrWSTI:
		cpu.ac[2] = cpu.ac[3]
		// TODO a lot of this should be moved into a func...
		unconverted := cpu.fpac[iPtr.ac]
		logging.DebugPrint(logging.DebugLog, "... FPAC %d = %f\n", iPtr.ac, unconverted)
		scaleFactor, dataType, size := memory.DecodeDecDataType(cpu.ac[1])
		if scaleFactor != 0 {
			log.Panicf("ERROR: Non-zero (%d) scale factors not yet supported\n", scaleFactor)
		}
		switch dataType {
		case memory.UnpackedDecLS:
			// if unconverted < 0 {
			// 	size++
			// }
			converted := fmt.Sprintf("%+0*.f", size, unconverted)
			for c := 0; c < size; c++ {
				memory.WriteByteBA(cpu.ac[3], dg.ByteT(converted[c]))
				logging.DebugPrint(logging.DebugLog, "... %c -> %#x\n", converted[c], cpu.ac[3])
				cpu.ac[3]++
			}
		case memory.UnpackedDecU:
			// if unconverted < 0 {
			// 	size++
			// }
			converted := fmt.Sprintf("%0*.f", size, unconverted)
			for c := 0; c < size; c++ {
				memory.WriteByteBA(cpu.ac[3], dg.ByteT(converted[c]))
				logging.DebugPrint(logging.DebugLog, "... %c -> %#x\n", converted[c], cpu.ac[3])
				cpu.ac[3]++
			}
		default:
			log.Panicf("ERROR: Decimal data type %d not yet supported\n", dataType)
		}

	case instrXFAMD:
		oneAccModeInd2Word := iPtr.variant.(oneAccModeInd2WordT)
		addr := resolve15bitDisplacement(cpu, oneAccModeInd2Word.ind, oneAccModeInd2Word.mode, dg.WordT(oneAccModeInd2Word.disp15), iPtr.dispOffset)
		fpQuad := dg.QwordT(memory.ReadDWord(addr))<<32 | dg.QwordT(memory.ReadDWord(addr+2))
		cpu.fpac[oneAccModeInd2Word.acd] += memory.DGdoubleToFloat64(fpQuad)
		cpu.SetZ(cpu.fpac[oneAccModeInd2Word.acd] == 0.0)
		cpu.SetN(cpu.fpac[oneAccModeInd2Word.acd] < 0.0)

	case instrXFAMS:
		oneAccModeInd2Word := iPtr.variant.(oneAccModeInd2WordT)
		addr := resolve15bitDisplacement(cpu, oneAccModeInd2Word.ind, oneAccModeInd2Word.mode, dg.WordT(oneAccModeInd2Word.disp15), iPtr.dispOffset)
		fpDoub := memory.ReadDWord(addr)
		cpu.fpac[oneAccModeInd2Word.acd] += memory.DGsingleToFloat64(fpDoub)
		cpu.SetZ(cpu.fpac[oneAccModeInd2Word.acd] == 0.0)
		cpu.SetN(cpu.fpac[oneAccModeInd2Word.acd] < 0.0)

	case instrXFLDD:
		oneAccModeInd2Word := iPtr.variant.(oneAccModeInd2WordT)
		addr := resolve15bitDisplacement(cpu, oneAccModeInd2Word.ind, oneAccModeInd2Word.mode, dg.WordT(oneAccModeInd2Word.disp15), iPtr.dispOffset)
		fpQuad := dg.QwordT(memory.ReadDWord(addr))<<32 | dg.QwordT(memory.ReadDWord(addr+2))
		cpu.fpac[oneAccModeInd2Word.acd] = memory.DGdoubleToFloat64(fpQuad)
		cpu.SetZ(cpu.fpac[oneAccModeInd2Word.acd] == 0.0)
		cpu.SetN(cpu.fpac[oneAccModeInd2Word.acd] < 0.0)

	case instrXFLDS:
		oneAccModeInd2Word := iPtr.variant.(oneAccModeInd2WordT)
		addr := resolve15bitDisplacement(cpu, oneAccModeInd2Word.ind, oneAccModeInd2Word.mode, dg.WordT(oneAccModeInd2Word.disp15), iPtr.dispOffset)
		fpSingle := memory.ReadDWord(addr)
		cpu.fpac[oneAccModeInd2Word.acd] = memory.DGsingleToFloat64(fpSingle)
		cpu.SetZ(cpu.fpac[oneAccModeInd2Word.acd] == 0.0)
		cpu.SetN(cpu.fpac[oneAccModeInd2Word.acd] < 0.0)

	case instrXFMMD:
		oneAccModeInd2Word := iPtr.variant.(oneAccModeInd2WordT)
		addr := resolve15bitDisplacement(cpu, oneAccModeInd2Word.ind, oneAccModeInd2Word.mode, dg.WordT(oneAccModeInd2Word.disp15), iPtr.dispOffset)
		fpQuad := dg.QwordT(memory.ReadDWord(addr))<<32 | dg.QwordT(memory.ReadDWord(addr+2))
		cpu.fpac[oneAccModeInd2Word.acd] *= memory.DGdoubleToFloat64(fpQuad)
		cpu.SetZ(cpu.fpac[oneAccModeInd2Word.acd] == 0.0)
		cpu.SetN(cpu.fpac[oneAccModeInd2Word.acd] < 0.0)

	case instrXFMMS:
		oneAccModeInd2Word := iPtr.variant.(oneAccModeInd2WordT)
		addr := resolve15bitDisplacement(cpu, oneAccModeInd2Word.ind, oneAccModeInd2Word.mode, dg.WordT(oneAccModeInd2Word.disp15), iPtr.dispOffset)
		fpSingle := memory.ReadDWord(addr)
		cpu.fpac[oneAccModeInd2Word.acd] *= memory.DGsingleToFloat64(fpSingle)
		cpu.SetZ(cpu.fpac[oneAccModeInd2Word.acd] == 0.0)
		cpu.SetN(cpu.fpac[oneAccModeInd2Word.acd] < 0.0)

	case instrXFSTD:
		oneAccModeInd2Word := iPtr.variant.(oneAccModeInd2WordT)
		addr := resolve15bitDisplacement(cpu, oneAccModeInd2Word.ind, oneAccModeInd2Word.mode, dg.WordT(oneAccModeInd2Word.disp15), iPtr.dispOffset)
		fpQuad := memory.Float64toDGdouble(cpu.fpac[oneAccModeInd2Word.acd])
		memory.WriteDWord(addr, dg.DwordT(fpQuad>>32))
		memory.WriteDWord(addr+2, dg.DwordT(fpQuad))

	default:
		log.Panicf("ERROR: EAGLE_FPU instruction <%s> not yet implemented\n", iPtr.mnemonic)
		return false
	}

	cpu.pc += dg.PhysAddrT(iPtr.instrLength)
	return true
}
