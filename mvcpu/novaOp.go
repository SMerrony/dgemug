// novaOp.go

// Copyright (C) 2017  Steve Merrony

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

func novaOp(cpu *CPUT, iPtr *decodedInstrT) bool {

	var (
		shifter          dg.WordT
		wideShifter      dg.DwordT
		tmpAcS, tmpAcD   dg.WordT
		savedCry, tmpCry bool
		pcInc            dg.PhysAddrT
		novaTwoAccMultOp novaTwoAccMultOpT
	)

	novaTwoAccMultOp = iPtr.variant.(novaTwoAccMultOpT)

	tmpAcS = memory.DwordGetLowerWord(cpu.ac[novaTwoAccMultOp.acs])
	tmpAcD = memory.DwordGetLowerWord(cpu.ac[novaTwoAccMultOp.acd])
	savedCry = cpu.carry

	// Preset Carry if required
	switch novaTwoAccMultOp.c {
	case 'Z': // zero
		cpu.carry = false
	case 'O': // One
		cpu.carry = true
	case 'C': // Complement
		cpu.carry = !cpu.carry
	}

	// perform the operation
	switch iPtr.ix {
	case instrADC:
		wideShifter = dg.DwordT(tmpAcD) + dg.DwordT(^tmpAcS)
		shifter = memory.DwordGetLowerWord(wideShifter)
		if wideShifter > 65535 {
			cpu.carry = !cpu.carry
		}

	case instrADD: // unsigned
		wideShifter = dg.DwordT(tmpAcD) + dg.DwordT(tmpAcS)
		shifter = memory.DwordGetLowerWord(wideShifter)
		if wideShifter > 65535 {
			cpu.carry = !cpu.carry
		}

	case instrAND:
		shifter = tmpAcD & tmpAcS

	case instrCOM:
		shifter = ^tmpAcS

	case instrINC:
		shifter = tmpAcS + 1
		if tmpAcS == 0xffff {
			cpu.carry = !cpu.carry
		}

	case instrMOV:
		shifter = tmpAcS

	case instrNEG:
		shifter = dg.WordT(-int16(tmpAcS))
		if tmpAcS == 0 {
			cpu.carry = !cpu.carry
		}

	case instrSUB:
		shifter = tmpAcD - tmpAcS
		if tmpAcS <= tmpAcD {
			cpu.carry = !cpu.carry
		}

	default:
		log.Fatalf("ERROR: NOVA_MEMREF instruction <%s> not yet implemented\n", iPtr.mnemonic)
	}

	// shift if required
	switch novaTwoAccMultOp.sh {
	case 'L':
		tmpCry = cpu.carry
		cpu.carry = memory.TestWbit(shifter, 0)
		shifter <<= 1
		if tmpCry {
			shifter |= 0x0001
		}
	case 'R':
		tmpCry = cpu.carry
		cpu.carry = memory.TestWbit(shifter, 15)
		shifter >>= 1
		if tmpCry {
			shifter |= 0x8000
		}
	case 'S':
		shifter = memory.SwapBytes(shifter)
	}

	// Skip?
	switch novaTwoAccMultOp.skip {
	case noSkip:
		pcInc = 1
	case skpSkip:
		pcInc = 2
	case szcSkip:
		if !cpu.carry {
			pcInc = 2
		} else {
			pcInc = 1
		}
	case sncSkip:
		if cpu.carry {
			pcInc = 2
		} else {
			pcInc = 1
		}
	case szrSkip:
		if shifter == 0 {
			pcInc = 2
		} else {
			pcInc = 1
		}
	case snrSkip:
		if shifter != 0 {
			pcInc = 2
		} else {
			pcInc = 1
		}
	case sezSkip:
		if !cpu.carry || shifter == 0 {
			pcInc = 2
		} else {
			pcInc = 1
		}
	case sbnSkip:
		if cpu.carry && shifter != 0 {
			pcInc = 2
		} else {
			pcInc = 1
		}
	default:
		log.Fatalln("ERROR: Invalid skip in novaOp()")
	}

	// No-Load?
	if novaTwoAccMultOp.nl == '#' {
		// don't load the result from the shifter, restore the Carry flag
		cpu.carry = savedCry
	} else {
		cpu.ac[novaTwoAccMultOp.acd] = dg.DwordT(shifter) & 0x0000ffff
	}

	cpu.pc += pcInc
	return true
}
