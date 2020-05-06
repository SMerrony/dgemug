// eagleDecimal.go

// Copyright ©2020 Steve Merrony

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

func eagleDecimal(cpu *CPUT, iPtr *decodedInstrT) bool {

	switch iPtr.ix {

	case instrWDecOp: // Funkiness ahead...
		switch iPtr.word2 {
		case 0x0000: // WDMOV
			log.Panicf("ERROR: EAGLE_DECIMAL instruction WDMOV not yet implemented")
		case 0x0001: // WDCMP
			// log.Panicf("ERROR: EAGLE_DECIMAL instruction WDCMP not yet implemented")
			arg1Type := cpu.ac[0]
			arg2Type := cpu.ac[1]
			arg1BA := cpu.ac[2]
			arg2BA := cpu.ac[3]
			if (arg1Type == arg2Type) && (arg1BA == arg2BA) { // short-circuit certain equality...
				logging.DebugPrint(logging.DebugLog, "... assuming equality\n")
				cpu.ac[1] = 0
			} else {
				sf1, dt1, sz1 := memory.DecodeDecDataType(arg1Type)
				str1 := memory.ReadDec(dg.PhysAddrT(arg1BA), sz1)
				sf2, dt2, sz2 := memory.DecodeDecDataType(arg2Type)
				str2 := memory.ReadDec(dg.PhysAddrT(arg2BA), sz2)
				logging.DebugPrint(logging.DebugLog, "Arg 1 - SF: %d., Type: %d., Size: %d., String: %s\n", sf1, dt1, sz1, str1)
				logging.DebugPrint(logging.DebugLog, "Arg 2 - SF: %d., Type: %d., Size: %d.\n", sf2, dt2, sz2)
				i1 := memory.DecIntToInt(dt1, str1)
				i2 := memory.DecIntToInt(dt2, str2)
				if sf1 == 0 && sf2 == 0 {
					switch {
					case i1 < i2:
						cpu.ac[1] = 0xffff_ffff
					case i1 == i2:
						cpu.ac[1] = 0
					case i1 > i2:
						cpu.ac[1] = 1
					}
				} else {
					log.Panicf("ERROR: EAGLE_DECIMAL instruction WDCMP not yet fully implemented")
				}
			}

		case 0x0002: // WDINC
			log.Panicf("ERROR: EAGLE_DECIMAL instruction WDINC not yet implemented")

		case 0x0003: // WDDEC
			log.Panicf("ERROR: EAGLE_DECIMAL instruction WDDEC not yet implemented")
		}
	default:
		log.Panicf("ERROR: EAGLE_DECIMAL instruction <%s> not yet implemented\n", iPtr.mnemonic)
		return false
	}

	cpu.pc += dg.PhysAddrT(iPtr.instrLength)
	return true
}
