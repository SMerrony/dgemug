// resolve.go

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

const (
	physMask16 = 0x7fff
	physMask32 = 0x7fff_ffff
	ringMask32 = 0x7000_0000
)

func resolve31bitDisplacement(cpu *CPUT, ind byte, mode int, disp int32, dispOffset int) (eff dg.PhysAddrT) {
	ring := cpu.pc & 0x7000_0000
	switch mode {
	case absoluteMode:
		// zero-extend to 28 bits, force to current ring...
		eff = dg.PhysAddrT(disp) // | (cpu.pc & 0x7000_0000)
	case pcMode:
		eff = dg.PhysAddrT(int32(cpu.pc) + disp + int32(dispOffset))
	case ac2Mode:
		eff = dg.PhysAddrT(int32(cpu.ac[2]) + disp)
	case ac3Mode:
		eff = dg.PhysAddrT(int32(cpu.ac[3]) + disp)
	}
	// handle indirection
	if ind == '@' { // down the rabbit hole...
		eff |= ring
		indAddr, ok := memory.ReadDwordTrap(eff)
		if !ok {
			log.Panicln("Terminating")
		}
		for memory.TestDwbit(indAddr, 0) {
			indAddr, ok = memory.ReadDwordTrap(dg.PhysAddrT(indAddr & physMask32))
			if !ok {
				log.Panicln("Terminating")
			}
		}
		eff = dg.PhysAddrT(indAddr) | ring
	}
	// check ATU
	if cpu.atu == false {
		// constrain result to 1st 32MB
		eff &= 0x1ff_ffff
	}

	if cpu.debugLogging {
		logging.DebugPrint(logging.DebugLog, "... resolve31bitDsiplacement got: %#o %s, returning %#o\n", disp, modeToString(mode), eff)
	}
	return eff & physMask32
}

func resolve15bitDisplacement(cpu *CPUT, ind byte, mode int, disp dg.WordT, dispOffset int) (eff dg.PhysAddrT) {
	var dispS32 int32
	ring := cpu.pc & 0x7000_0000
	if mode != absoluteMode {
		// relative mode
		// sign-extend to 32-bits
		dispS32 = int32(int16(disp<<1) >> 1)
	}
	switch mode {
	case absoluteMode:
		// zero-extend to 28 bits, force to current ring...
		eff = dg.PhysAddrT(disp) | ring
	case pcMode:
		eff = dg.PhysAddrT(int32(cpu.pc) + dispS32 + int32(dispOffset))
	case ac2Mode:
		eff = dg.PhysAddrT(int32(cpu.ac[2]) + dispS32)
	case ac3Mode:
		eff = dg.PhysAddrT(int32(cpu.ac[3]) + dispS32)
	}
	// handle indirection
	if ind == '@' { // down the rabbit hole...
		eff |= ring
		indAddr, ok := memory.ReadDwordTrap(eff)
		if cpu.debugLogging {
			logging.DebugPrint(logging.DebugLog, "... resolve15bitDisplacement got: @%#o %s, reading %#o, got %#o\n", disp, modeToString(mode), eff, indAddr)
		}
		if !ok {
			log.Panicln("Terminating")
		}
		for memory.TestDwbit(indAddr, 0) {
			indAddr, ok = memory.ReadDwordTrap(dg.PhysAddrT(indAddr & physMask32))
			if cpu.debugLogging {
				logging.DebugPrint(logging.DebugLog, "... resolve15bitDisplacement ... reading %#o\n", indAddr)
			}
			if !ok {
				log.Panicln("Terminating")
			}
		}
		eff = dg.PhysAddrT(indAddr) | ring
	}
	// check ATU
	if cpu.atu == false {
		// constrain result to 1st 32MB
		eff &= 0x1ff_ffff
	}
	// if cpu.debugLogging {
	// 	logging.DebugPrint(logging.DebugLog, "... resolve15bitDisplacement got: %#o %s, returning %#o\n", disp, modeToString(mode), eff)
	// }
	return eff
}

func resolve8bitDisplacement(cpu *CPUT, ind byte, mode int, disp int16) (eff dg.PhysAddrT) {
	if mode == absoluteMode {
		// zero-extend to 28 bits, force to current ring...
		eff = dg.PhysAddrT(disp) | (cpu.pc & 0x7000_0000)
	} else {
		// relative mode
		// sign-extend to 31-bits
		eff = dg.PhysAddrT(disp)
		if disp < 0 {
			eff |= 0xffff_f800
		}
	}
	switch mode {
	case pcMode:
		eff += cpu.pc
	case ac2Mode:
		eff += dg.PhysAddrT(cpu.ac[2])
	case ac3Mode:
		eff += dg.PhysAddrT(cpu.ac[3])
	}

	// handle indirection
	if ind == '@' { // down the rabbit hole...
		eff |= (cpu.pc & 0x7000_0000)
		indAddr, ok := memory.ReadWordTrap(eff)
		if !ok {
			log.Panicln("Terminating")
		}
		for memory.TestWbit(indAddr, 0) {
			indAddr, ok = memory.ReadWordTrap(dg.PhysAddrT(indAddr & physMask16))
			if !ok {
				log.Panicln("Terminating")
			}
		}
		eff = dg.PhysAddrT(indAddr)
	}
	// check ATU
	if cpu.atu == false {
		// constrain result to 1st 32MB
		eff &= 0x1ff_ffff
	}
	if cpu.debugLogging {
		logging.DebugPrint(logging.DebugLog, "... resolve8bitDisplacement got: %#o %s, returning %#o\n", disp, modeToString(mode), eff)
	}
	return eff
}

func resolve16bitByteAddr(cpu *CPUT, mode int, disp16 int16, loByte bool) (eff dg.PhysAddrT) {
	switch mode {
	case absoluteMode:
		eff = ((cpu.pc & 0x7000_0000) << 1) | dg.PhysAddrT(disp16)
	case pcMode:
		eff = dg.PhysAddrT(int(cpu.pc<<1) + int(disp16))
	case ac2Mode:
		eff = dg.PhysAddrT(int(cpu.ac[2]<<1) + int(disp16))
		eff &= 0x1fff_ffff
		eff |= (cpu.pc & 0x7000_0000) << 1
	case ac3Mode:
		eff = dg.PhysAddrT(int(cpu.ac[3]<<1) + int(disp16))
		eff &= 0x1fff_ffff
		eff |= (cpu.pc & 0x7000_0000) << 1
	}
	if loByte {
		eff++
	}
	return eff
}

func resolve32bitEffAddr(cpu *CPUT, ind byte, mode int, disp int32, dispOffset int) (eff dg.PhysAddrT) {
	switch mode {
	case absoluteMode:
		eff = dg.PhysAddrT(disp)
	case pcMode:
		eff = dg.PhysAddrT(int32(cpu.pc) + disp + int32(dispOffset))
	case ac2Mode:
		eff = dg.PhysAddrT(int32(cpu.ac[2]) + disp)
	case ac3Mode:
		eff = dg.PhysAddrT(int32(cpu.ac[3]) + disp)
	}
	// handle indirection
	if ind == '@' { //|| memory.TestDwbit(dg.DwordT(eff), 0) { // down the rabbit hole...
		indAddr, ok := memory.ReadDwordTrap(eff)
		if !ok {
			log.Panicln("Terminating")
		}
		for memory.TestDwbit(indAddr, 0) {
			indAddr, ok = memory.ReadDwordTrap(dg.PhysAddrT(indAddr & physMask32))
			if !ok {
				log.Panicln("Terminating")
			}
		}
		eff = dg.PhysAddrT(indAddr)
	}
	// check ATU
	if cpu.atu == false {
		// constrain result to 1st 32MB
		eff &= 0x1ff_ffff
	}
	// log.Printf("... resolve32bitEffAddr got: %d. %s, returning %#x\n", disp, modeToString(mode), eff)
	if cpu.debugLogging {
		logging.DebugPrint(logging.DebugLog, "... resolve32bitEffAddr got: %#o %s, returning %#o\n", disp, modeToString(mode), eff)
	}
	return eff
}

func resolve32bitIndirectableAddr(cpu *CPUT, iAddr dg.DwordT) dg.PhysAddrT {
	eff := iAddr
	// handle indirection
	for memory.TestDwbit(eff, 0) {
		eff = memory.ReadDWord(dg.PhysAddrT(eff & physMask32))
	}
	// check ATU
	if cpu.atu == false {
		// constrain result to 1st 32MB
		eff &= 0x1ff_ffff
	}
	return dg.PhysAddrT(eff)
}

// resolveEclipseBitAddr as per page 10-8 of Pop
// Used by BTO, BTZ, SNB, SZB, SZBO
func resolveEclipseBitAddr(cpu *CPUT, twoAcc1Word *twoAcc1WordT) (wordAddr dg.PhysAddrT, bitNum uint) {
	// TODO handle segments and indirection
	if twoAcc1Word.acd == twoAcc1Word.acs {
		wordAddr = 0
	} else {
		if memory.TestDwbit(cpu.ac[twoAcc1Word.acs], 0) {
			log.Panicln("ERROR: Indirect 16-bit BIT pointers not yet supported")
		}
		wordAddr = dg.PhysAddrT(cpu.ac[twoAcc1Word.acs]) & physMask16 // mask off lower 15 bits
	}
	offset := dg.PhysAddrT(cpu.ac[twoAcc1Word.acd]) >> 4
	wordAddr += offset // add unsigned offset
	bitNum = uint(cpu.ac[twoAcc1Word.acd] & 0x000f)
	return wordAddr, bitNum
}

// resolveEagleeBitAddr as per page 1-17 of Pop
// Used by eg. WSZB
func resolveEagleBitAddr(cpu *CPUT, twoAcc1Word *twoAcc1WordT) (wordAddr dg.PhysAddrT, bitNum uint) {
	// TODO handle segments and indirection
	if twoAcc1Word.acd == twoAcc1Word.acs {
		wordAddr = 0
	} else {
		if memory.TestDwbit(cpu.ac[twoAcc1Word.acs], 0) {
			log.Panicln("ERROR: Indirect 32-bit BIT pointers not yet supported")
		}
		wordAddr = dg.PhysAddrT(cpu.ac[twoAcc1Word.acs])
	}
	offset := dg.PhysAddrT(cpu.ac[twoAcc1Word.acd]) >> 4
	wordAddr += offset // add unsigned offset
	bitNum = uint(cpu.ac[twoAcc1Word.acd] & 0x000f)
	return wordAddr, bitNum
}
