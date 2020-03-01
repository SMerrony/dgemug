// resolve.go

// Copyright (C) 2017,2019  Steve Merrony

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
	physMask32 = 0x7fffffff
)

func resolve15bitDisplacement(cpuPtr *MvCPUT, ind byte, mode int, disp dg.WordT, dispOffset int) (eff dg.PhysAddrT) {
	if mode == absoluteMode {
		// zero-extend to 28 bits, force to current ring...
		eff = dg.PhysAddrT(disp) | (cpuPtr.pc & 0x7000_0000)
	} else {
		// relative mode
		// sign-extend to 31-bits
		eff = dg.PhysAddrT(disp)
		if memory.TestWbit(disp, 1) {
			eff |= 0xffff_8000
		}
	}
	switch mode {
	case pcMode:
		eff += cpuPtr.pc + dg.PhysAddrT(dispOffset)
	case ac2Mode:
		eff += dg.PhysAddrT(cpuPtr.ac[2])
	case ac3Mode:
		eff += dg.PhysAddrT(cpuPtr.ac[3])
	}
	// handle indirection
	if ind == '@' { // down the rabbit hole...
		indAddr, ok := memory.ReadDwordTrap(eff)
		if !ok {
			log.Fatalln("Terminating")
		}
		for memory.TestDwbit(indAddr, 0) {
			indAddr, ok = memory.ReadDwordTrap(dg.PhysAddrT(indAddr & physMask32))
			if !ok {
				log.Fatalln("Terminating")
			}
		}
		eff = dg.PhysAddrT(indAddr)
	}
	// check ATU
	if cpuPtr.atu == false {
		// constrain result to 1st 32MB
		eff &= 0x1ff_ffff
	}
	if debugLogging {
		logging.DebugPrint(logging.DebugLog, "... resolve15bitDsiplacement got: %#o %s, returning %#o\n", disp, modeToString(mode), eff)
	}
	return eff
}

func resolve8bitDisplacement(cpuPtr *MvCPUT, ind byte, mode int, disp int16) (eff dg.PhysAddrT) {
	if mode == absoluteMode {
		// zero-extend to 28 bits, force to current ring...
		eff = dg.PhysAddrT(disp) | (cpuPtr.pc & 0x7000_0000)
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
		eff += cpuPtr.pc
	case ac2Mode:
		eff += dg.PhysAddrT(cpuPtr.ac[2])
	case ac3Mode:
		eff += dg.PhysAddrT(cpuPtr.ac[3])
	}
	// handle indirection
	if ind == '@' { // down the rabbit hole...
		indAddr, ok := memory.ReadWordTrap(eff)
		if !ok {
			log.Fatalln("Terminating")
		}
		for memory.TestWbit(indAddr, 0) {
			indAddr, ok = memory.ReadWordTrap(dg.PhysAddrT(indAddr & physMask16))
			if !ok {
				log.Fatalln("Terminating")
			}
		}
		eff = dg.PhysAddrT(indAddr)
	}
	// check ATU
	if cpuPtr.atu == false {
		// constrain result to 1st 32MB
		eff &= 0x1ff_ffff
	}
	if debugLogging {
		logging.DebugPrint(logging.DebugLog, "... resolve8bitDsiplacement got: %#o %s, returning %#o\n", disp, modeToString(mode), eff)
	}
	return eff
}

// resolve32bitByteAddr returns the word address and low-byte flag for a given 32-bit byte address
func resolve32bitByteAddr(byteAddr dg.DwordT) (wordAddr dg.PhysAddrT, loByte bool) {
	wordAddr = dg.PhysAddrT(byteAddr) >> 1
	loByte = memory.TestDwbit(byteAddr, 31)
	return wordAddr, loByte
}

func resolve32bitEffAddr(cpuPtr *MvCPUT, ind byte, mode int, disp int32, dispOffset int) (eff dg.PhysAddrT) {

	eff = dg.PhysAddrT(disp)

	// handle addressing mode...
	switch mode {
	case absoluteMode:
		// nothing to do
	case pcMode:
		eff += cpuPtr.pc + dg.PhysAddrT(dispOffset)
	case ac2Mode:
		eff += dg.PhysAddrT(cpuPtr.ac[2])
	case ac3Mode:
		eff += dg.PhysAddrT(cpuPtr.ac[3])
	}

	// handle indirection
	if ind == '@' || memory.TestDwbit(dg.DwordT(eff), 0) { // down the rabbit hole...
		indAddr, ok := memory.ReadDwordTrap(eff)
		if !ok {
			log.Fatalln("Terminating")
		}
		for memory.TestDwbit(indAddr, 0) {
			indAddr, ok = memory.ReadDwordTrap(dg.PhysAddrT(indAddr & physMask32))
			if !ok {
				log.Fatalln("Terminating")
			}
		}
		eff = dg.PhysAddrT(indAddr)
	}

	// check ATU
	if cpuPtr.atu == false {
		// constrain result to 1st 32MB
		eff &= 0x1ff_ffff
	}

	if debugLogging {
		logging.DebugPrint(logging.DebugLog, "... resolve32bitEffAddr got: %#o %s, returning %#o\n", disp, modeToString(mode), eff)
	}
	return eff
}

func resolve32bitIndirectableAddr(cpuPtr *MvCPUT, iAddr dg.DwordT) dg.PhysAddrT {
	eff := iAddr
	// handle indirection
	for memory.TestDwbit(eff, 0) {
		eff = memory.ReadDWord(dg.PhysAddrT(eff & physMask32))
	}
	// check ATU
	if cpuPtr.atu == false {
		// constrain result to 1st 32MB
		eff &= 0x1ff_ffff
	}
	return dg.PhysAddrT(eff)
}

// resolveEclipseBitAddr as per page 10-8 of Pop
// Used by BTO, BTZ, SNB, SZB, SZBO
func resolveEclipseBitAddr(cpuPtr *MvCPUT, twoAcc1Word *twoAcc1WordT) (wordAddr dg.PhysAddrT, bitNum uint) {
	// TODO handle segments and indirection
	if twoAcc1Word.acd == twoAcc1Word.acs {
		wordAddr = 0
	} else {
		if memory.TestDwbit(cpuPtr.ac[twoAcc1Word.acs], 0) {
			log.Fatal("ERROR: Indirect 16-bit BIT pointers not yet supported")
		}
		wordAddr = dg.PhysAddrT(cpuPtr.ac[twoAcc1Word.acs]) & physMask16 // mask off lower 15 bits
	}
	offset := dg.PhysAddrT(cpuPtr.ac[twoAcc1Word.acd]) >> 4
	wordAddr += offset // add unsigned offset
	bitNum = uint(cpuPtr.ac[twoAcc1Word.acd] & 0x000f)
	return wordAddr, bitNum
}

// resolveEagleeBitAddr as per page 1-17 of Pop
// Used by eg. WSZB
func resolveEagleBitAddr(cpuPtr *MvCPUT, twoAcc1Word *twoAcc1WordT) (wordAddr dg.PhysAddrT, bitNum uint) {
	// TODO handle segments and indirection
	if twoAcc1Word.acd == twoAcc1Word.acs {
		wordAddr = 0
	} else {
		if memory.TestDwbit(cpuPtr.ac[twoAcc1Word.acs], 0) {
			log.Fatal("ERROR: Indirect 32-bit BIT pointers not yet supported")
		}
		wordAddr = dg.PhysAddrT(cpuPtr.ac[twoAcc1Word.acs])
	}
	offset := dg.PhysAddrT(cpuPtr.ac[twoAcc1Word.acd]) >> 4
	wordAddr += offset // add unsigned offset
	bitNum = uint(cpuPtr.ac[twoAcc1Word.acd] & 0x000f)
	return wordAddr, bitNum
}
