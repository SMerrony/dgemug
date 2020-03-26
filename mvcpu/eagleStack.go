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

////////////////////////////////////////////////////////////////
// N.B. Be sure to use Double-Word memory references here... //
//////////////////////////////////////////////////////////////

package mvcpu

import (
	"log"

	"github.com/SMerrony/dgemug/dg"
	"github.com/SMerrony/dgemug/logging"
	"github.com/SMerrony/dgemug/memory"
)

// Some Page Zero special locations for the Wide Stack...
const (
	WsfhLoc = 014 // 12.
	WfpLoc  = 020 // WFP 16.
	WspLoc  = 022 // WSP 18.
	WslLoc  = 024 // WSL 20.
	WsbLoc  = 026 // WSB 22.
)

func eagleStack(cpu *CPUT, iPtr *decodedInstrT) bool {

	switch iPtr.ix {

	// N.B. DSZTS and ISZTS are in eaglePC

	case instrLDAFP:
		oneAcc1Word := iPtr.variant.(oneAcc1WordT)
		cpu.ac[oneAcc1Word.acd] = dg.DwordT(cpu.wfp)
		cpu.SetOVR(false)

	case instrLDASB:
		oneAcc1Word := iPtr.variant.(oneAcc1WordT)
		cpu.ac[oneAcc1Word.acd] = dg.DwordT(cpu.wsb)
		cpu.SetOVR(false)

	case instrLDASL:
		oneAcc1Word := iPtr.variant.(oneAcc1WordT)
		cpu.ac[oneAcc1Word.acd] = dg.DwordT(cpu.wsl)
		cpu.SetOVR(false)

	case instrLDASP:
		oneAcc1Word := iPtr.variant.(oneAcc1WordT)
		cpu.ac[oneAcc1Word.acd] = dg.DwordT(cpu.wsp)
		cpu.SetOVR(false)

	case instrLDATS:
		oneAcc1Word := iPtr.variant.(oneAcc1WordT)
		cpu.ac[oneAcc1Word.acd] = memory.ReadDWord(cpu.wsp)
		cpu.SetOVR(false)

	case instrLPEF:
		noAccModeInd3Word := iPtr.variant.(noAccModeInd3WordT)
		wsPush(cpu, 0, dg.DwordT(resolve31bitDisplacement(cpu, noAccModeInd3Word.ind, noAccModeInd3Word.mode, noAccModeInd3Word.disp31, iPtr.dispOffset)))
		cpu.SetOVR(false)

	case instrLPEFB:
		noAccMode3Word := iPtr.variant.(noAccMode3WordT)
		eff := dg.DwordT(noAccMode3Word.immU32)
		switch noAccMode3Word.mode {
		case absoluteMode: // do nothing
		case pcMode:
			eff += dg.DwordT(cpu.pc)
		case ac2Mode:
			eff += cpu.ac[2]
		case ac3Mode:
			eff += cpu.ac[3]
		}
		wsPush(cpu, 0, eff)
		cpu.SetOVR(false)

	case instrSTAFP:
		oneAcc1Word := iPtr.variant.(oneAcc1WordT)
		// FIXME handle segments
		cpu.wfp = dg.PhysAddrT(cpu.ac[oneAcc1Word.acd])
		// according the PoP does not write through to page zero...
		//memory.WriteDWord(memory.WfpLoc, cpu.ac[oneAcc1Word.acd])
		cpu.SetOVR(false)

	case instrSTASB:
		oneAcc1Word := iPtr.variant.(oneAcc1WordT)
		// FIXME handle segments
		cpu.wsb = dg.PhysAddrT(cpu.ac[oneAcc1Word.acd])
		memory.WriteDWord(WsbLoc, cpu.ac[oneAcc1Word.acd]) // write-through to p.0
		cpu.SetOVR(false)

	case instrSTASL:
		oneAcc1Word := iPtr.variant.(oneAcc1WordT)
		// FIXME handle segments
		cpu.wsl = dg.PhysAddrT(cpu.ac[oneAcc1Word.acd])
		memory.WriteDWord(WslLoc, cpu.ac[oneAcc1Word.acd]) // write-through to p.0
		cpu.SetOVR(false)

	case instrSTASP:
		oneAcc1Word := iPtr.variant.(oneAcc1WordT)
		// FIXME handle segments
		cpu.wsp = dg.PhysAddrT(cpu.ac[oneAcc1Word.acd])
		// according the PoP does not write through to page zero...
		// memory.WriteDWord(memory.WspLoc, cpu.ac[oneAcc1Word.acd])
		cpu.SetOVR(false)
		if cpu.debugLogging {
			logging.DebugPrint(logging.DebugLog, "... STASP set WSP to %#o\n", cpu.ac[oneAcc1Word.acd])
		}

	case instrSTATS:
		oneAcc1Word := iPtr.variant.(oneAcc1WordT)
		// FIXME handle segments
		memory.WriteDWord(dg.PhysAddrT(memory.ReadDWord(cpu.wsp)), cpu.ac[oneAcc1Word.acd])
		cpu.SetOVR(false)

	case instrWFPOP:
		cpu.fpac[3] = float64(int(wsPopQWord(cpu, 0)))
		cpu.fpac[2] = float64(int(wsPopQWord(cpu, 0)))
		cpu.fpac[1] = float64(int(wsPopQWord(cpu, 0)))
		cpu.fpac[0] = float64(int(wsPopQWord(cpu, 0)))
		tmpQwd := wsPopQWord(cpu, 0)
		cpu.fpsr = 0
		any := false
		// set the ANY bit?
		if memory.GetQwbits(tmpQwd, 1, 4) != 0 {
			memory.SetQwbit(&cpu.fpsr, 0)
			any = true
		}
		// copy bits 1-11
		for b := 1; b <= 11; b++ {
			if memory.TestQwbit(tmpQwd, b) {
				memory.SetQwbit(&cpu.fpsr, uint(b))
			}
		}
		// bits 28-31
		if any {
			for b := 28; b <= 31; b++ {
				if memory.TestQwbit(tmpQwd, b) {
					memory.SetQwbit(&cpu.fpsr, uint(b))
				}
			}
			for b := 33; b <= 63; b++ {
				if memory.TestQwbit(tmpQwd, b) {
					memory.SetQwbit(&cpu.fpsr, uint(b))
				}
			}
		}

	case instrWMSP:
		oneAcc1Word := iPtr.variant.(oneAcc1WordT)
		tmpDwd := cpu.ac[oneAcc1Word.acd] << 1
		tmpDwd += dg.DwordT(cpu.wsp) // memory.WspLoc)
		// FIXME - handle overflow
		cpu.wsp = dg.PhysAddrT(tmpDwd)
		cpu.SetOVR(false)

	case instrWPOP:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		firstAc := twoAcc1Word.acs
		lastAc := twoAcc1Word.acd
		thisAc := firstAc
		for {
			cpu.ac[thisAc] = wsPop(cpu, 0)
			if thisAc == lastAc {
				break
			}
			thisAc--
			if thisAc == -1 {
				thisAc = 3
			}
		}

		// if lastAc > firstAc {
		// 	firstAc += 4
		// }
		// var acsUp = [8]int{0, 1, 2, 3, 0, 1, 2, 3}
		// for thisAc := firstAc; thisAc >= lastAc; thisAc-- {
		// 	if cpu.debugLogging {
		// 		logging.DebugPrint(logging.DebugLog, "... wide popping AC%d\n", acsUp[thisAc])
		// 	}
		// 	cpu.ac[acsUp[thisAc]] = wsPop(cpu, 0)
		// }
		cpu.SetOVR(false)

	case instrWPSH:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		firstAc := twoAcc1Word.acs
		lastAc := twoAcc1Word.acd
		thisAc := firstAc
		for {
			wsPush(cpu, 0, cpu.ac[thisAc])
			if thisAc == lastAc {
				break
			}
			thisAc++
			if thisAc == 4 {
				thisAc = 0
			}
		}
		// if lastAc < firstAc {
		// 	lastAc += 4
		// }
		// var acsUp = [8]int{0, 1, 2, 3, 0, 1, 2, 3}
		// for thisAc := firstAc; thisAc <= lastAc; thisAc++ {
		// 	if cpu.debugLogging {
		// 		logging.DebugPrint(logging.DebugLog, "... wide pushing AC%d\n", acsUp[thisAc])
		// 	}
		// 	wsPush(cpu, 0, cpu.ac[acsUp[thisAc]])
		// }
		cpu.SetOVR(false)

	// N.B. WRTN is in eaglePC

	case instrWSAVR:
		unique2Word := iPtr.variant.(unique2WordT)
		wsav(cpu, &unique2Word)
		cpu.SetOVK(false)

	case instrWSAVS:
		unique2Word := iPtr.variant.(unique2WordT)
		wsav(cpu, &unique2Word)
		cpu.SetOVK(true)

	case instrWSSVR:
		unique2Word := iPtr.variant.(unique2WordT)
		wssav(cpu, &unique2Word)
		cpu.SetOVK(false)
		cpu.SetOVR(false)

	case instrWSSVS:
		unique2Word := iPtr.variant.(unique2WordT)
		wssav(cpu, &unique2Word)
		cpu.SetOVK(true)
		cpu.SetOVR(false)

	case instrXPEF:
		wsPush(cpu, 0, dg.DwordT(resolve15bitDisplacement(cpu, iPtr.ind, iPtr.mode, iPtr.disp15, iPtr.dispOffset)))

	case instrXPEFB:
		noAccMode2Word := iPtr.variant.(noAccMode2WordT)
		// FIXME check for overflow
		eff := dg.DwordT(noAccMode2Word.disp16)
		switch noAccMode2Word.mode {
		case absoluteMode: // do nothing
		case pcMode:
			eff += dg.DwordT(cpu.pc)
		case ac2Mode:
			eff += cpu.ac[2]
		case ac3Mode:
			eff += cpu.ac[3]
		}
		wsPush(cpu, 0, eff)

	case instrXPSHJ:
		// FIXME check for overflow
		immMode2Word := iPtr.variant.(immMode2WordT)
		wsPush(cpu, 0, dg.DwordT(cpu.pc+2))
		//cpu.pc = resolve32bitEffAddr(cpu, immMode2Word.ind, immMode2Word.mode, int32(immMode2Word.disp15), iPtr.dispOffset)
		cpu.pc = (cpu.pc & ringMask32) | resolve15bitDisplacement(cpu, immMode2Word.ind, immMode2Word.mode, dg.WordT(immMode2Word.disp15), iPtr.dispOffset)
		return true

	default:
		log.Fatalf("ERROR: EAGLE_STACK instruction <%s> not yet implemented\n", iPtr.mnemonic)
		return false
	}

	cpu.pc += dg.PhysAddrT(iPtr.instrLength)
	return true
}

// wsav is common to WSAVR and WSAVS
func wsav(cpu *CPUT, u2wd *unique2WordT) {
	wsPush(cpu, 0, cpu.ac[0])          // 1
	wsPush(cpu, 0, cpu.ac[1])          // 2
	wsPush(cpu, 0, cpu.ac[2])          // 3
	wsPush(cpu, 0, dg.DwordT(cpu.wfp)) // 4
	dwd := cpu.ac[3] & 0x7fffffff
	if cpu.carry {
		dwd |= 0x80000000
	}
	wsPush(cpu, 0, dwd) // 5
	cpu.wfp = cpu.wsp
	cpu.ac[3] = dg.DwordT(cpu.wsp)
	dwdCnt := uint(u2wd.immU16)
	if dwdCnt > 0 {
		advanceWSP(cpu, dwdCnt)
	}
}

// wssav is common to WSSVR and WSSVS
func wssav(cpu *CPUT, u2wd *unique2WordT) {
	wsPush(cpu, 0, memory.DwordFromTwoWords(cpu.psr, 0)) // 1
	wsPush(cpu, 0, cpu.ac[0])                            // 2
	wsPush(cpu, 0, cpu.ac[1])                            // 3
	wsPush(cpu, 0, cpu.ac[2])                            // 4
	wsPush(cpu, 0, dg.DwordT(cpu.wfp))                   // 5
	dwd := cpu.ac[3] & 0x7fffffff
	if cpu.carry {
		dwd |= 0x80000000
	}
	wsPush(cpu, 0, dwd) // 6
	cpu.wfp = cpu.wsp
	cpu.ac[3] = dg.DwordT(cpu.wsp)
	dwdCnt := uint(u2wd.immU16)
	if dwdCnt > 0 {
		advanceWSP(cpu, dwdCnt)
	}
}

// wsPush - PUSH a doubleword onto the Wide Stack
func wsPush(cpu *CPUT, seg dg.PhysAddrT, data dg.DwordT) {
	// TODO overflow/underflow handling - either here or in instruction?
	cpu.wsp += 2
	memory.WriteDWord(cpu.wsp, data)
	logging.DebugPrint(logging.DebugLog, "... wsPush pushed %#o onto the Wide Stack at location: %#o\n", data, cpu.wsp)
}

// WsPop - POP a doubleword off the Wide Stack
func wsPop(cpu *CPUT, seg dg.PhysAddrT) (dword dg.DwordT) {
	dword = memory.ReadDWord(cpu.wsp)
	cpu.wsp -= 2
	logging.DebugPrint(logging.DebugLog, "... wsPop  popped %#o off  the Wide Stack at location: %#o\n", dword, cpu.wsp+2)
	return dword
}

// used by WPOPB, WRTN
func wpopb(cpu *CPUT) {
	wspSav := cpu.wsp
	// pop off 6 double words
	dwd := wsPop(cpu, 0) // 1
	cpu.carry = memory.TestDwbit(dwd, 0)
	cpu.pc = dg.PhysAddrT(dwd & 0x7fffffff)
	cpu.ac[3] = wsPop(cpu, 0) // 2
	// replace WFP with popped value of AC3
	cpu.wfp = dg.PhysAddrT(cpu.ac[3])
	cpu.ac[2] = wsPop(cpu, 0) // 3
	cpu.ac[1] = wsPop(cpu, 0) // 4
	cpu.ac[0] = wsPop(cpu, 0) // 5
	dwd = wsPop(cpu, 0)       // 6
	cpu.psr = memory.DwordGetUpperWord(dwd)
	// TODO Set WFP is crossing rings
	wsFramSz2 := (int(dwd&0x00007fff) * 2) + 12
	cpu.wsp = wspSav - dg.PhysAddrT(wsFramSz2)
}

// wsPopQWord - POP a Quad-word off the Wide Stack
func wsPopQWord(cpu *CPUT, seg dg.PhysAddrT) dg.QwordT {
	var qw dg.QwordT
	rhDWord := wsPop(cpu, seg)
	lhDWord := wsPop(cpu, seg)
	qw = dg.QwordT(lhDWord)<<32 | dg.QwordT(rhDWord)
	return qw
}

// advanceWSP increases the WSP by the given amount of DWords
func advanceWSP(cpu *CPUT, dwdCnt uint) {
	cpu.wsp += dg.PhysAddrT(dwdCnt * 2)
	logging.DebugPrint(logging.DebugLog, "... WSP advanced by %#o DWords to %#o\n", dwdCnt, cpu.wsp)
}
