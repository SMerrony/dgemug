// eclipseStack.go

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
	"github.com/SMerrony/dgemug/logging"
	"github.com/SMerrony/dgemug/memory"
)

func eclipseStack(cpu *CPUT, iPtr *decodedInstrT) bool {

	ring := cpu.pc & 0x7000_0000

	switch iPtr.ix {

	case instrMSP:
		// TODO handle overflow
		s16 := int16(cpu.ac[iPtr.ac])
		nsp := int16(memory.ReadWord(memory.NspLoc|ring)) + s16
		memory.WriteWord(memory.NspLoc|ring, dg.WordT(nsp))

	case instrPOP:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		first := twoAcc1Word.acs
		last := twoAcc1Word.acd
		if last > first {
			first += 4
		}
		acsUp := [8]int{0, 1, 2, 3, 0, 1, 2, 3}
		for thisAc := first; thisAc >= last; thisAc-- {
			if cpu.debugLogging {
				logging.DebugPrint(logging.DebugLog, "... narrow popping AC%d\n", acsUp[thisAc])
			}

			cpu.ac[acsUp[thisAc]] = dg.DwordT(memory.NsPop(ring, cpu.debugLogging))
		}

	case instrPOPJ:
		addr := dg.PhysAddrT(memory.NsPop(ring, cpu.debugLogging))
		cpu.pc = (addr & 0x7fff) | ring
		return true // because PC set

	case instrPSH:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		first := twoAcc1Word.acs
		last := twoAcc1Word.acd
		if last < first {
			last += 4
		}
		acsUp := [8]int{0, 1, 2, 3, 0, 1, 2, 3}
		for thisAc := first; thisAc <= last; thisAc++ {
			if cpu.debugLogging {
				logging.DebugPrint(logging.DebugLog, "... narrow pushing AC%d\n", acsUp[thisAc])
			}
			memory.NsPush(ring, memory.DwordGetLowerWord(cpu.ac[acsUp[thisAc]]), cpu.debugLogging)
		}

	case instrPSHJ:
		memory.NsPush(ring, dg.WordT(cpu.pc)+2, cpu.debugLogging)
		addr := resolve15bitDisplacement(cpu, iPtr.ind, iPtr.mode, iPtr.disp15, iPtr.dispOffset)
		addr &= 0x7fff
		addr |= ring
		cpu.pc = addr
		return true // because PC set

	case instrRTN:
		// // complement of SAVE
		// memory.WriteWord(memory.NspLoc, memory.ReadWord(memory.NfpLoc)) // ???
		// //memory.WriteWord(memory.NfpLoc, memory.ReadWord(memory.NspLoc)) // ???
		// word := memory.NsPop(0, cpu.debugLogging)
		// cpu.carry = memory.TestWbit(word, 0)
		// cpu.pc = dg.PhysAddrT(word) & 0x7fff
		// //nfpSave := memory.NsPop(0)               // 1
		// cpu.ac[3] = dg.DwordT(memory.NsPop(0, cpu.debugLogging)) // 2
		// cpu.ac[2] = dg.DwordT(memory.NsPop(0, cpu.debugLogging)) // 3
		// cpu.ac[1] = dg.DwordT(memory.NsPop(0, cpu.debugLogging)) // 4
		// cpu.ac[0] = dg.DwordT(memory.NsPop(0, cpu.debugLogging)) // 5
		// memory.WriteWord(memory.NfpLoc, memory.DWordGetLowerWord(cpu.ac[3]))
		// return true // because PC set

		nfpSav := memory.ReadWord(memory.NfpLoc | ring)
		memory.WriteWord(memory.NspLoc|(ring), nfpSav)
		pwd1 := memory.NsPop(ring, cpu.debugLogging) // 1
		cpu.carry = memory.TestWbit(pwd1, 0)
		cpu.pc = dg.PhysAddrT((pwd1 & 0x07fff)) | ring
		cpu.ac[3] = dg.DwordT(memory.NsPop(ring, cpu.debugLogging)) // 2
		cpu.ac[2] = dg.DwordT(memory.NsPop(ring, cpu.debugLogging)) // 3
		cpu.ac[1] = dg.DwordT(memory.NsPop(ring, cpu.debugLogging)) // 4
		cpu.ac[0] = dg.DwordT(memory.NsPop(ring, cpu.debugLogging)) // 5
		//memory.WriteWord(memory.NspLoc, nfpSav-5)
		memory.WriteWord(memory.NfpLoc|ring, memory.DwordGetLowerWord(cpu.ac[3]))

		return true // because PC set

	case instrSAVE:
		unique2Word := iPtr.variant.(unique2WordT)
		i := dg.WordT(unique2Word.immU16)
		nfpSav := memory.ReadWord(memory.NfpLoc | ring)
		nspSav := memory.ReadWord(memory.NspLoc | ring)

		// // version based in simH Nova SAVn
		// memory.WriteWord(memory.NspLoc, nspSav+i)
		// memory.NsPush(ring, memory.DWordGetLowerWord(cpu.ac[0]), cpu.debugLogging) // 1
		// memory.NsPush(ring, memory.DWordGetLowerWord(cpu.ac[1]), cpu.debugLogging) // 2
		// memory.NsPush(ring, memory.DWordGetLowerWord(cpu.ac[2]), cpu.debugLogging) // 3
		// memory.NsPush(ring, nfpSav, cpu.debugLogging)                               // 4
		// word := memory.DWordGetLowerWord(cpu.ac[3])
		// if cpu.carry {
		// 	word |= 0x8000
		// } else {
		// 	word &= 0x7fff
		// }
		// memory.NsPush(ring, word, cpu.debugLogging) // 5
		// cpu.ac[3] = dg.DwordT(memory.ReadWord(memory.NspLoc))
		// memory.WriteWord(memory.NfpLoc, memory.DWordGetLowerWord(cpu.ac[3]))

		// version based on 32-bit PoP
		memory.NsPush(ring, memory.DwordGetLowerWord(cpu.ac[0]), cpu.debugLogging) // 1
		memory.NsPush(ring, memory.DwordGetLowerWord(cpu.ac[1]), cpu.debugLogging) // 2
		memory.NsPush(ring, memory.DwordGetLowerWord(cpu.ac[2]), cpu.debugLogging) // 3
		memory.NsPush(ring, nfpSav, cpu.debugLogging)                              // 4
		word := memory.DwordGetLowerWord(cpu.ac[3])
		if cpu.carry {
			word |= 0x8000
		} else {
			word &= 0x7fff
		}
		memory.NsPush(ring, word, cpu.debugLogging) // 5
		memory.WriteWord(memory.NspLoc|ring, nspSav+5+i)
		memory.WriteWord(memory.NfpLoc|ring, nspSav+5)
		cpu.ac[3] = dg.DwordT(nspSav + 5)

	default:
		log.Panicf("ERROR: ECLIPSE_STACK instruction <%s> not yet implemented\n", iPtr.mnemonic)
		return false
	}

	cpu.pc += dg.PhysAddrT(iPtr.instrLength)
	return true
}
