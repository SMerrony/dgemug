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

	"github.com/SMerrony/dgemug/dg"
	"github.com/SMerrony/dgemug/memory"
)

func eagleFPU(cpu *CPUT, iPtr *decodedInstrT) bool {
	switch iPtr.ix {

	case instrWSTI:
		cpu.ac[2] = cpu.ac[3]
		// TODO a lot of this should be moved into a func...
		unconverted := cpu.fpac[iPtr.ac]
		scaleFactor := int(int8(memory.GetDwbits(cpu.ac[1], 0, 8)))
		if scaleFactor != 0 {
			log.Fatalf("ERROR: Non-zero (%d) scale factors not yet supported\n", scaleFactor)
		}
		dataType := uint8(memory.GetDwbits(cpu.ac[1], 24, 3))
		size := int(uint8(memory.GetDwbits(cpu.ac[1], 27, 5)))
		switch dataType {
		case 3: // <sign><zeroes><int>
			if unconverted < 0 {
				size++
			}
			converted := fmt.Sprintf("%+*.f", size, unconverted)
			for c := 0; c < size; c++ {
				memory.WriteByteBA(cpu.ac[3], dg.ByteT(converted[c]))
				cpu.ac[3]++
			}
		default:
			log.Fatalf("ERROR: Decimal data type %d not yet supported\n", dataType)
		}

	case instrWFLAD:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		cpu.fpac[twoAcc1Word.acd] = float64(int32(cpu.ac[twoAcc1Word.acs])) // N.B INT32 conversion required!!!
		cpu.SetZ(cpu.fpac[twoAcc1Word.acd] == 0.0)
		cpu.SetN(cpu.fpac[twoAcc1Word.acd] < 0.0)

	default:
		log.Fatalf("ERROR: EAGLE_FPU instruction <%s> not yet implemented\n", iPtr.mnemonic)
		return false
	}

	cpu.pc += dg.PhysAddrT(iPtr.instrLength)
	return true
}
