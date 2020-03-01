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

func eagleStack(cpuPtr *CPUT, iPtr *decodedInstrT) bool {

	switch iPtr.ix {

	// N.B. DSZTS and ISZTS are in eaglePC

	case instrLDAFP:
		oneAcc1Word := iPtr.variant.(oneAcc1WordT)
		cpuPtr.ac[oneAcc1Word.acd] = dg.DwordT(cpuPtr.wfp)
		cpuPtr.CPUSetOVR(false)

	case instrLDASB:
		oneAcc1Word := iPtr.variant.(oneAcc1WordT)
		cpuPtr.ac[oneAcc1Word.acd] = dg.DwordT(cpuPtr.wsb)
		cpuPtr.CPUSetOVR(false)

	case instrLDASL:
		oneAcc1Word := iPtr.variant.(oneAcc1WordT)
		cpuPtr.ac[oneAcc1Word.acd] = dg.DwordT(cpuPtr.wsl)
		cpuPtr.CPUSetOVR(false)

	case instrLDASP:
		oneAcc1Word := iPtr.variant.(oneAcc1WordT)
		cpuPtr.ac[oneAcc1Word.acd] = dg.DwordT(cpuPtr.wsp)
		cpuPtr.CPUSetOVR(false)

	case instrLDATS:
		oneAcc1Word := iPtr.variant.(oneAcc1WordT)
		cpuPtr.ac[oneAcc1Word.acd] = memory.ReadDWord(cpuPtr.wsp)
		cpuPtr.CPUSetOVR(false)

	case instrLPEF:
		noAccModeInd3Word := iPtr.variant.(noAccModeInd3WordT)
		wsPush(cpuPtr, 0, dg.DwordT(resolve32bitEffAddr(cpuPtr, noAccModeInd3Word.ind, noAccModeInd3Word.mode, noAccModeInd3Word.disp31, iPtr.dispOffset)))
		cpuPtr.CPUSetOVR(false)

	case instrLPEFB:
		noAccMode3Word := iPtr.variant.(noAccMode3WordT)
		eff := dg.DwordT(noAccMode3Word.immU32)
		switch noAccMode3Word.mode {
		case absoluteMode: // do nothing
		case pcMode:
			eff += dg.DwordT(cpuPtr.pc)
		case ac2Mode:
			eff += cpuPtr.ac[2]
		case ac3Mode:
			eff += cpuPtr.ac[3]
		}
		wsPush(cpuPtr, 0, eff)
		cpuPtr.CPUSetOVR(false)

	case instrSTAFP:
		oneAcc1Word := iPtr.variant.(oneAcc1WordT)
		// FIXME handle segments
		cpuPtr.wfp = dg.PhysAddrT(cpuPtr.ac[oneAcc1Word.acd])
		// according the PoP does not write through to page zero...
		//memory.WriteDWord(memory.WfpLoc, cpuPtr.ac[oneAcc1Word.acd])
		cpuPtr.CPUSetOVR(false)

	case instrSTASB:
		oneAcc1Word := iPtr.variant.(oneAcc1WordT)
		// FIXME handle segments
		cpuPtr.wsb = dg.PhysAddrT(cpuPtr.ac[oneAcc1Word.acd])
		memory.WriteDWord(WsbLoc, cpuPtr.ac[oneAcc1Word.acd]) // write-through to p.0
		cpuPtr.CPUSetOVR(false)

	case instrSTASL:
		oneAcc1Word := iPtr.variant.(oneAcc1WordT)
		// FIXME handle segments
		cpuPtr.wsl = dg.PhysAddrT(cpuPtr.ac[oneAcc1Word.acd])
		memory.WriteDWord(WslLoc, cpuPtr.ac[oneAcc1Word.acd]) // write-through to p.0
		cpuPtr.CPUSetOVR(false)

	case instrSTASP:
		oneAcc1Word := iPtr.variant.(oneAcc1WordT)
		// FIXME handle segments
		cpuPtr.wsp = dg.PhysAddrT(cpuPtr.ac[oneAcc1Word.acd])
		// according the PoP does not write through to page zero...
		// memory.WriteDWord(memory.WspLoc, cpuPtr.ac[oneAcc1Word.acd])
		cpuPtr.CPUSetOVR(false)
		if debugLogging {
			logging.DebugPrint(logging.DebugLog, "... STASP set WSP to %#o\n", cpuPtr.ac[oneAcc1Word.acd])
		}

	case instrSTATS:
		oneAcc1Word := iPtr.variant.(oneAcc1WordT)
		// FIXME handle segments
		memory.WriteDWord(dg.PhysAddrT(memory.ReadDWord(cpuPtr.wsp)), cpuPtr.ac[oneAcc1Word.acd])
		cpuPtr.CPUSetOVR(false)

	case instrWFPOP:
		cpuPtr.fpac[3] = wsPopQWord(cpuPtr, 0)
		cpuPtr.fpac[2] = wsPopQWord(cpuPtr, 0)
		cpuPtr.fpac[1] = wsPopQWord(cpuPtr, 0)
		cpuPtr.fpac[0] = wsPopQWord(cpuPtr, 0)
		tmpQwd := wsPopQWord(cpuPtr, 0)
		cpuPtr.fpsr = 0
		any := false
		// set the ANY bit?
		if memory.GetQwbits(tmpQwd, 1, 4) != 0 {
			memory.SetQwbit(&cpuPtr.fpsr, 0)
			any = true
		}
		// copy bits 1-11
		for b := 1; b <= 11; b++ {
			if memory.TestQwbit(tmpQwd, b) {
				memory.SetQwbit(&cpuPtr.fpsr, uint(b))
			}
		}
		// bits 28-31
		if any {
			for b := 28; b <= 31; b++ {
				if memory.TestQwbit(tmpQwd, b) {
					memory.SetQwbit(&cpuPtr.fpsr, uint(b))
				}
			}
			for b := 33; b <= 63; b++ {
				if memory.TestQwbit(tmpQwd, b) {
					memory.SetQwbit(&cpuPtr.fpsr, uint(b))
				}
			}
		}

	case instrWMSP:
		oneAcc1Word := iPtr.variant.(oneAcc1WordT)
		tmpDwd := cpuPtr.ac[oneAcc1Word.acd] << 1
		tmpDwd += dg.DwordT(cpuPtr.wsp) // memory.WspLoc)
		// FIXME - handle overflow
		cpuPtr.wsp = dg.PhysAddrT(tmpDwd)
		cpuPtr.CPUSetOVR(false)

	case instrWPOP:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		firstAc := twoAcc1Word.acs
		lastAc := twoAcc1Word.acd
		if lastAc > firstAc {
			firstAc += 4
		}
		var acsUp = [8]int{0, 1, 2, 3, 0, 1, 2, 3}
		for thisAc := firstAc; thisAc >= lastAc; thisAc-- {
			if debugLogging {
				logging.DebugPrint(logging.DebugLog, "... wide popping AC%d\n", acsUp[thisAc])
			}
			cpuPtr.ac[acsUp[thisAc]] = wsPop(cpuPtr, 0)
		}
		cpuPtr.CPUSetOVR(false)

	case instrWPSH:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		firstAc := twoAcc1Word.acs
		lastAc := twoAcc1Word.acd
		if lastAc < firstAc {
			lastAc += 4
		}
		var acsUp = [8]int{0, 1, 2, 3, 0, 1, 2, 3}
		for thisAc := firstAc; thisAc <= lastAc; thisAc++ {
			if debugLogging {
				logging.DebugPrint(logging.DebugLog, "... wide pushing AC%d\n", acsUp[thisAc])
			}
			wsPush(cpuPtr, 0, cpuPtr.ac[acsUp[thisAc]])
		}
		cpuPtr.CPUSetOVR(false)

	// N.B. WRTN is in eaglePC

	case instrWSAVR:
		unique2Word := iPtr.variant.(unique2WordT)
		wsav(cpuPtr, &unique2Word)
		cpuPtr.CPUSetOVK(false)

	case instrWSAVS:
		unique2Word := iPtr.variant.(unique2WordT)
		wsav(cpuPtr, &unique2Word)
		cpuPtr.CPUSetOVK(true)

	case instrWSSVR:
		unique2Word := iPtr.variant.(unique2WordT)
		wssav(cpuPtr, &unique2Word)
		cpuPtr.CPUSetOVK(false)
		cpuPtr.CPUSetOVR(false)

	case instrXPEF:
		noAccModeInd2Word := iPtr.variant.(noAccModeInd2WordT)
		wsPush(cpuPtr, 0, dg.DwordT(resolve15bitDisplacement(cpuPtr, noAccModeInd2Word.ind, noAccModeInd2Word.mode, noAccModeInd2Word.disp15, iPtr.dispOffset)))

	case instrXPEFB:
		noAccMode2Word := iPtr.variant.(noAccMode2WordT)
		// FIXME check for overflow
		eff := dg.DwordT(noAccMode2Word.disp16)
		switch noAccMode2Word.mode {
		case absoluteMode: // do nothing
		case pcMode:
			eff += dg.DwordT(cpuPtr.pc)
		case ac2Mode:
			eff += cpuPtr.ac[2]
		case ac3Mode:
			eff += cpuPtr.ac[3]
		}
		wsPush(cpuPtr, 0, eff)

	case instrXPSHJ:
		// FIXME check for overflow
		immMode2Word := iPtr.variant.(immMode2WordT)
		wsPush(cpuPtr, 0, dg.DwordT(cpuPtr.pc+2))
		cpuPtr.pc = resolve32bitEffAddr(cpuPtr, immMode2Word.ind, immMode2Word.mode, int32(immMode2Word.disp15), iPtr.dispOffset)
		return true

	default:
		log.Fatalf("ERROR: EAGLE_STACK instruction <%s> not yet implemented\n", iPtr.mnemonic)
		return false
	}

	cpuPtr.pc += dg.PhysAddrT(iPtr.instrLength)
	return true
}

// wsav is common to WSAVR and WSAVS
func wsav(cpuPtr *CPUT, u2wd *unique2WordT) {
	wfpSav := dg.DwordT(cpuPtr.wfp)
	wsPush(cpuPtr, 0, cpuPtr.ac[0]) // 1
	wsPush(cpuPtr, 0, cpuPtr.ac[1]) // 2
	wsPush(cpuPtr, 0, cpuPtr.ac[2]) // 3
	wsPush(cpuPtr, 0, wfpSav)       // 4
	dwd := cpuPtr.ac[3] & 0x7fffffff
	if cpuPtr.carry {
		dwd |= 0x80000000
	}
	wsPush(cpuPtr, 0, dwd) // 5
	cpuPtr.wfp = cpuPtr.wsp
	cpuPtr.ac[3] = dg.DwordT(cpuPtr.wsp)
	dwdCnt := uint(u2wd.immU16)
	if dwdCnt > 0 {
		advanceWSP(cpuPtr, dwdCnt)
	}
}

// wssav is common to WSSVR and WSSVS
func wssav(cpuPtr *CPUT, u2wd *unique2WordT) {
	wfpSav := dg.DwordT(cpuPtr.wfp)
	wsPush(cpuPtr, 0, memory.DwordFromTwoWords(cpuPtr.psr, 0)) // 1
	wsPush(cpuPtr, 0, cpuPtr.ac[0])                            // 2
	wsPush(cpuPtr, 0, cpuPtr.ac[1])                            // 3
	wsPush(cpuPtr, 0, cpuPtr.ac[2])                            // 4
	wsPush(cpuPtr, 0, wfpSav)                                  // 5
	dwd := cpuPtr.ac[3] & 0x7fffffff
	if cpuPtr.carry {
		dwd |= 0x80000000
	}
	wsPush(cpuPtr, 0, dwd) // 6
	cpuPtr.wfp = cpuPtr.wsp
	cpuPtr.ac[3] = dg.DwordT(cpuPtr.wsp)
	dwdCnt := uint(u2wd.immU16)
	if dwdCnt > 0 {
		advanceWSP(cpuPtr, dwdCnt)
	}
}

// wsPush - PUSH a doubleword onto the Wide Stack
func wsPush(cpuPtr *CPUT, seg dg.PhysAddrT, data dg.DwordT) {
	// TODO segment handling
	// TODO overflow/underflow handling - either here or in instruction?
	cpuPtr.wsp += 2
	memory.WriteDWord(cpuPtr.wsp, data)
	logging.DebugPrint(logging.DebugLog, "... wsPush pushed %#o onto the Wide Stack at location: %#o\n", data, cpuPtr.wsp)
}

// WsPop - POP a doubleword off the Wide Stack
func wsPop(cpuPtr *CPUT, seg dg.PhysAddrT) (dword dg.DwordT) {
	// TODO segment handling
	dword = memory.ReadDWord(cpuPtr.wsp)
	cpuPtr.wsp -= 2
	logging.DebugPrint(logging.DebugLog, "... wsPop  popped %#o off  the Wide Stack at location: %#o\n", dword, cpuPtr.wsp+2)
	return dword
}

// used by WPOPB, WRTN
func wpopb(cpuPtr *CPUT) {
	wspSav := cpuPtr.wsp
	// pop off 6 double words
	dwd := wsPop(cpuPtr, 0) // 1
	cpuPtr.carry = memory.TestDwbit(dwd, 0)
	cpuPtr.pc = dg.PhysAddrT(dwd & 0x7fffffff)
	cpuPtr.ac[3] = wsPop(cpuPtr, 0) // 2
	// replace WFP with popped value of AC3
	cpuPtr.wfp = dg.PhysAddrT(cpuPtr.ac[3])
	cpuPtr.ac[2] = wsPop(cpuPtr, 0) // 3
	cpuPtr.ac[1] = wsPop(cpuPtr, 0) // 4
	cpuPtr.ac[0] = wsPop(cpuPtr, 0) // 5
	dwd = wsPop(cpuPtr, 0)          // 6
	cpuPtr.psr = memory.DwordGetUpperWord(dwd)
	// TODO Set WFP is crossing rings
	wsFramSz2 := (int(dwd&0x00007fff) * 2) + 12
	cpuPtr.wsp = wspSav - dg.PhysAddrT(wsFramSz2)
}

// wsPopQWord - POP a Quad-word off the Wide Stack
func wsPopQWord(cpuPtr *CPUT, seg dg.PhysAddrT) dg.QwordT {
	// TODO segment handling
	var qw dg.QwordT
	rhDWord := wsPop(cpuPtr, seg)
	lhDWord := wsPop(cpuPtr, seg)
	qw = dg.QwordT(lhDWord)<<32 | dg.QwordT(rhDWord)
	return qw
}

// advanceWSP increases the WSP by the given amount of DWords
func advanceWSP(cpuPtr *CPUT, dwdCnt uint) {
	cpuPtr.wsp += dg.PhysAddrT(dwdCnt * 2)
	logging.DebugPrint(logging.DebugLog, "... WSP advanced by %#o DWords to %#o\n", dwdCnt, cpuPtr.wsp)
}
