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
	wsfhLoc = 014 // 12.
	wfpLoc  = 020 // WFP 16.
	wspLoc  = 022 // WSP 18.
	wslLoc  = 024 // WSL 20.
	wsbLoc  = 026 // WSB 22.
)

// Wide Stack Fault codes, values are significant
const (
	wsfOverflow       = 0
	wsfPending        = 1
	wsfTooManyArgs    = 2
	wsfUnderflow      = 3
	wsfReturnOverflow = 4
)

func eagleStack(cpu *CPUT, iPtr *decodedInstrT) bool {

	switch iPtr.ix {

	// N.B. DSZTS and ISZTS are in eaglePC

	case instrLDAFP:
		cpu.ac[iPtr.ac] = dg.DwordT(cpu.wfp)
		cpu.SetOVR(false)

	case instrLDASB:
		cpu.ac[iPtr.ac] = dg.DwordT(cpu.wsb)
		cpu.SetOVR(false)

	case instrLDASL:
		cpu.ac[iPtr.ac] = dg.DwordT(cpu.wsl)
		cpu.SetOVR(false)

	case instrLDASP:
		cpu.ac[iPtr.ac] = dg.DwordT(cpu.wsp)
		cpu.SetOVR(false)

	case instrLDATS:
		cpu.ac[iPtr.ac] = memory.ReadDWord(cpu.wsp)
		cpu.SetOVR(false)

	case instrLPEF:
		noAccModeInd3Word := iPtr.variant.(noAccModeInd3WordT)
		wsPush(cpu, dg.DwordT(resolve31bitDisplacement(cpu, noAccModeInd3Word.ind, noAccModeInd3Word.mode, noAccModeInd3Word.disp31, iPtr.dispOffset)))
		cpu.SetOVR(false)

	case instrLPEFB:
		noAccMode3Word := iPtr.variant.(noAccMode3WordT)
		eff := dg.DwordT(noAccMode3Word.immU32)
		switch noAccMode3Word.mode {
		case absoluteMode: // do nothing
		case pcMode:
			eff |= (dg.DwordT(cpu.pc) << 1)
		case ac2Mode:
			eff |= (cpu.ac[2] << 1)
		case ac3Mode:
			eff |= (cpu.ac[3] << 1)
		}
		wsPush(cpu, eff)
		cpu.SetOVR(false)

	case instrSTAFP:
		// FIXME handle segments
		cpu.wfp = dg.PhysAddrT(cpu.ac[iPtr.ac])
		// according the PoP does not write through to page zero...
		//memory.WriteDWord(memory.wfpLoc, cpu.ac[iPtr.ac])
		cpu.SetOVR(false)

	case instrSTASB:
		cpu.wsb = dg.PhysAddrT(cpu.ac[iPtr.ac])
		memory.WriteDWord((cpu.pc&0x7000_0000)|wsbLoc, cpu.ac[iPtr.ac]) // write-through to p.0
		cpu.SetOVR(false)

	case instrSTASL:
		cpu.wsl = dg.PhysAddrT(cpu.ac[iPtr.ac])
		memory.WriteDWord((cpu.pc&0x7000_0000)|wslLoc, cpu.ac[iPtr.ac]) // write-through to p.0
		cpu.SetOVR(false)

	case instrSTASP:
		// FIXME handle segments
		cpu.wsp = dg.PhysAddrT(cpu.ac[iPtr.ac])
		// according the PoP does not write through to page zero...
		// memory.WriteDWord(memory.wspLoc, cpu.ac[iPtr.ac])
		cpu.SetOVR(false)
		if cpu.debugLogging {
			logging.DebugPrint(logging.DebugLog, "... STASP set WSP to %#o\n", cpu.ac[iPtr.ac])
		}

	case instrSTATS:
		memory.WriteDWord(cpu.wsp, cpu.ac[iPtr.ac])
		cpu.SetOVR(false)

	case instrWFPOP:
		cpu.fpac[3] = memory.DGdoubleToFloat64(wsPopQWord(cpu, 0))
		cpu.fpac[2] = memory.DGdoubleToFloat64(wsPopQWord(cpu, 0))
		cpu.fpac[1] = memory.DGdoubleToFloat64(wsPopQWord(cpu, 0))
		cpu.fpac[0] = memory.DGdoubleToFloat64(wsPopQWord(cpu, 0))
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

	case instrWFPSH:
		wsPushQWord(cpu, cpu.fpsr) // TODO Is this right?
		wsPushQWord(cpu, memory.Float64toDGdouble(cpu.fpac[0]))
		wsPushQWord(cpu, memory.Float64toDGdouble(cpu.fpac[1]))
		wsPushQWord(cpu, memory.Float64toDGdouble(cpu.fpac[2]))
		wsPushQWord(cpu, memory.Float64toDGdouble(cpu.fpac[3]))

	case instrWMSP:
		sMove := int(int32(cpu.ac[iPtr.ac]) * 2)
		ok, faultCode, secondaryFault := wspCheckBounds(cpu, sMove, true)
		if !ok {
			log.Printf("DEBUG: Stack fault trapped by WMSP, codes %d and %d", faultCode, secondaryFault)
			wspHandleFault(cpu, iPtr.instrLength, faultCode, secondaryFault)
			return true // we have set PC
		}
		cpu.wsp = dg.PhysAddrT(sMove) + cpu.wsp
		cpu.SetOVR(false)

	case instrWPOP:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		firstAc := twoAcc1Word.acs
		lastAc := twoAcc1Word.acd
		thisAc := firstAc
		for {
			cpu.ac[thisAc] = WsPop(cpu, 0)
			if thisAc == lastAc {
				break
			}
			thisAc--
			if thisAc == -1 {
				thisAc = 3
			}
		}
		cpu.SetOVR(false)

	// N.B. WPOPB & WPOPJ are in eaglePC.go

	case instrWPSH:
		twoAcc1Word := iPtr.variant.(twoAcc1WordT)
		firstAc := twoAcc1Word.acs
		lastAc := twoAcc1Word.acd
		thisAc := firstAc
		for {
			wsPush(cpu, cpu.ac[thisAc])
			if thisAc == lastAc {
				break
			}
			thisAc++
			if thisAc == 4 {
				thisAc = 0
			}
		}
		cpu.SetOVR(false)

	// N.B. WRTN is in eaglePC.go

	case instrWSAVR, instrWSAVS:
		unique2Word := iPtr.variant.(unique2WordT)
		ok, faultCode, secondaryFault := wspCheckBounds(cpu, int(unique2Word.immU16)*2+12, true)
		if !ok {
			log.Printf("DEBUG: Stack fault trapped by WSAVR/S, codes %d and %d", faultCode, secondaryFault)
			wspHandleFault(cpu, iPtr.instrLength, faultCode, secondaryFault)
			return true // we have set PC
		}
		wsav(cpu, &unique2Word)
		switch iPtr.ix {
		case instrWSAVR:
			cpu.SetOVK(false)
		case instrWSAVS:
			cpu.SetOVK(true)
		}

	case instrWSSVR, instrWSSVS:
		unique2Word := iPtr.variant.(unique2WordT)
		ok, faultCode, secondaryFault := wspCheckBounds(cpu, int(unique2Word.immU16)*2+12, true)
		if !ok {
			log.Printf("DEBUG: Stack fault trapped by WSSVR/S, codes %d and %d", faultCode, secondaryFault)
			wspHandleFault(cpu, iPtr.instrLength, faultCode, secondaryFault)
			return true // we have set PC
		}
		wssav(cpu, &unique2Word)
		switch iPtr.ix {
		case instrWSSVR:
			cpu.SetOVK(false)
			cpu.SetOVR(false)
		case instrWSSVS:
			cpu.SetOVK(true)
			cpu.SetOVR(false)
		}

	case instrXPEF:
		wsPush(cpu, dg.DwordT(resolve15bitDisplacement(cpu, iPtr.ind, iPtr.mode, iPtr.disp15, iPtr.dispOffset)))

	case instrXPEFB:
		noAccMode2Word := iPtr.variant.(noAccMode2WordT)
		// FIXME check for overflow
		eff := resolve16bitByteAddr(cpu, noAccMode2Word.mode, noAccMode2Word.disp16, noAccMode2Word.lowByte)
		wsPush(cpu, dg.DwordT(eff))

	case instrXPSHJ:
		// FIXME check for overflow
		immMode2Word := iPtr.variant.(immMode2WordT)
		wsPush(cpu, dg.DwordT(cpu.pc+2))
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
	// ok, faultCode, secondaryFault := wspCheckBounds(cpu, int(u2wd.immU16)*2+12, true)
	// if !ok {
	// 	log.Printf("DEBUG: Stack fault trapped by wsav(), codes %d and %d", faultCode, secondaryFault)
	// 	wspHandleFault(cpu, faultCode, secondaryFault)
	// }
	dwd := cpu.ac[3] & 0x7fffffff
	if cpu.carry {
		dwd |= 0x80000000
	}
	wsPush(cpu, cpu.ac[0])          // 1
	wsPush(cpu, cpu.ac[1])          // 2
	wsPush(cpu, cpu.ac[2])          // 3
	wsPush(cpu, dg.DwordT(cpu.wfp)) // 4
	wsPush(cpu, dwd)                // 5
	cpu.wfp = cpu.wsp
	cpu.ac[3] = dg.DwordT(cpu.wsp)
	dwdCnt := uint(u2wd.immU16)
	if dwdCnt > 0 {
		advanceWSP(cpu, dwdCnt)
	}
}

func wsPushSpecialReturnBlock(cpu *CPUT) {
	dwd := cpu.ac[3] & 0x7fffffff
	if cpu.carry {
		dwd |= 0x80000000
	}
	wsPush(cpu, memory.DwordFromTwoWords(cpu.psr, 0)) // 1
	wsPush(cpu, cpu.ac[0])                            // 2
	wsPush(cpu, cpu.ac[1])                            // 3
	wsPush(cpu, cpu.ac[2])                            // 4
	wsPush(cpu, dg.DwordT(cpu.wfp))                   // 5
	wsPush(cpu, dwd)                                  // 6
}

// wssav is common to WSSVR and WSSVS
func wssav(cpu *CPUT, u2wd *unique2WordT) {
	ok, faultCode, secondaryFault := wspCheckBounds(cpu, int(u2wd.immU16)*2+12, true)
	if !ok {
		log.Fatalf("DEBUG: Stack fault trapped by wssav(), codes %d and %d", faultCode, secondaryFault)
	}
	wsPushSpecialReturnBlock(cpu)
	cpu.wfp = cpu.wsp
	cpu.ac[3] = dg.DwordT(cpu.wsp)
	dwdCnt := uint(u2wd.immU16)
	if dwdCnt > 0 {
		advanceWSP(cpu, dwdCnt)
	}
}

// wsPush - PUSH a doubleword onto the Wide Stack
func wsPush(cpu *CPUT, data dg.DwordT) {
	// TODO overflow/underflow handling - either here or in instruction?
	cpu.wsp += 2
	memory.WriteDWord(cpu.wsp, data)
	if cpu.debugLogging {
		logging.DebugPrint(logging.DebugLog, "... wsPush pushed %#o onto the Wide Stack at location: %#o\n", data, cpu.wsp)
	}
}

func wsPushQWord(cpu *CPUT, qw dg.QwordT) {
	cpu.wsp += 2
	memory.WriteDWord(cpu.wsp, dg.DwordT(qw>>32))
	cpu.wsp += 2
	memory.WriteDWord(cpu.wsp, dg.DwordT(qw))
}

// WsPop - POP a doubleword off the Wide Stack
func WsPop(cpu *CPUT, seg dg.PhysAddrT) (dword dg.DwordT) {
	dword = memory.ReadDWord(cpu.wsp)
	cpu.wsp -= 2
	if cpu.debugLogging {
		logging.DebugPrint(logging.DebugLog, "... wsPop  popped %#o off  the Wide Stack at location: %#o\n", dword, cpu.wsp+2)
	}
	return dword
}

// used by WPOPB, WRTN
func wpopb(cpu *CPUT) {
	wspSav := cpu.wsp
	// pop off 6 double words
	dwd := WsPop(cpu, 0) // 1
	cpu.carry = memory.TestDwbit(dwd, 0)
	cpu.pc = dg.PhysAddrT(dwd & 0x7fffffff)
	cpu.ac[3] = WsPop(cpu, 0) // 2
	// replace WFP with popped value of AC3
	cpu.wfp = dg.PhysAddrT(cpu.ac[3])
	cpu.ac[2] = WsPop(cpu, 0) // 3
	cpu.ac[1] = WsPop(cpu, 0) // 4
	cpu.ac[0] = WsPop(cpu, 0) // 5
	dwd = WsPop(cpu, 0)       // 6
	cpu.psr = memory.DwordGetUpperWord(dwd)
	// TODO Set WFP is crossing rings
	wsFramSz2 := (int(dwd&0x00007fff) * 2) + 12
	cpu.wsp = wspSav - dg.PhysAddrT(wsFramSz2)
}

// wsPopQWord - POP a Quad-word off the Wide Stack
func wsPopQWord(cpu *CPUT, seg dg.PhysAddrT) dg.QwordT {
	var qw dg.QwordT
	rhDWord := WsPop(cpu, seg)
	lhDWord := WsPop(cpu, seg)
	qw = dg.QwordT(lhDWord)<<32 | dg.QwordT(rhDWord)
	return qw
}

// advanceWSP increases the WSP by the given amount of DWords
func advanceWSP(cpu *CPUT, dwdCnt uint) {
	cpu.wsp += dg.PhysAddrT(dwdCnt * 2)
	if cpu.debugLogging {
		logging.DebugPrint(logging.DebugLog, "... WSP advanced by %#o DWords to %#o\n", dwdCnt, cpu.wsp)
	}
}

// wspCheckBounds does a pre-flight check to see if the intended change of WSP would cause a stack fault
// isSave must be set by WMSP, WSSVR, WSSVS, WSAVR & WSAVS
func wspCheckBounds(cpu *CPUT, wspChangeWds int, isSave bool) (ok bool, primaryFault, secondaryFault int) {
	ok = true
	if wspChangeWds > 0 {
		if cpu.wsp+dg.PhysAddrT(wspChangeWds) > cpu.wsl {
			if isSave {
				return false, wsfPending, wsfOverflow
			}
			return false, wsfOverflow, 0
		}
	} else {
		posMove := -wspChangeWds
		if cpu.wsp-dg.PhysAddrT(posMove) < cpu.wsb {
			if isSave {
				return false, wsfPending, wsfUnderflow
			}
			return false, wsfUnderflow, 0
		}
	}
	return ok, 0, 0
}

func wspHandleFault(cpu *CPUT, instrLen int, primaryFault, secondaryFault int) {
	// from pp.5-23 of PoP
	// Step 1
	if primaryFault == wsfUnderflow {
		cpu.wsp = cpu.wsl
	}
	// Step 2
	dwd := dg.DwordT(cpu.pc)
	if primaryFault != wsfPending {
		dwd += dg.DwordT(instrLen)
	}
	if cpu.carry {
		dwd |= 0x80000000
	}
	wsPush(cpu, memory.DwordFromTwoWords(cpu.psr, 0)) // 1
	wsPush(cpu, cpu.ac[0])                            // 2
	wsPush(cpu, cpu.ac[1])                            // 3
	wsPush(cpu, cpu.ac[2])                            // 4
	wsPush(cpu, dg.DwordT(cpu.wfp))                   // 5
	wsPush(cpu, dwd)                                  // 6
	// Step 3
	memory.ClearWbit(&cpu.psr, 0) // OVK
	memory.ClearWbit(&cpu.psr, 1) // OVR
	// TODO IRES???
	// Step 4
	cpu.wsp = cpu.wsp &^ (1 << 31)
	// Step 5
	cpu.wsl |= 0x8000_0000
	// Step 6
	wsSaveToMemory(cpu)
	// Step 7
	cpu.ac[0] = dg.DwordT(cpu.pc)
	// Step 8
	cpu.ac[1] = dg.DwordT(primaryFault)
	// Step 9
	wsfhAddr := dg.PhysAddrT(memory.ReadWord((cpu.pc & 0x7000_0000) | wsfhLoc))
	wsfhAddr |= (cpu.pc & 0x7000_0000)
	log.Printf("DEBUG: Calling Wide Stack Fault Handler at %#x (#%o)", wsfhAddr, wsfhAddr)
	cpu.pc = wsfhAddr
}

func wsSaveToMemory(cpu *CPUT) {
	seg := (cpu.pc & 0x7000_0000)
	memory.WriteDWord(seg+wfpLoc, dg.DwordT(cpu.wfp))
	memory.WriteDWord(seg+wspLoc, dg.DwordT(cpu.wsp))
	memory.WriteDWord(seg+wslLoc, dg.DwordT(cpu.wsl))
	memory.WriteDWord(seg+wsbLoc, dg.DwordT(cpu.wsb))
}
