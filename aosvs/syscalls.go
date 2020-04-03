// syscalls.go - map of AOS/VS system calls

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

package aosvs

import (
	"bytes"
	"log"

	"github.com/SMerrony/dgemug/memory"

	"github.com/SMerrony/dgemug/dg"
	"github.com/SMerrony/dgemug/mvcpu"
)

type syscallDescT struct {
	name        string                                 // AOS/VS System Call Name
	alias       string                                 // AOS/VS 4-char Name Alias
	syscallType int                                    // groupings as per Table 2-1 in Sys Call Dict
	fn          func(*mvcpu.CPUT, chan AgentReqT) bool // implementation
	fn16        func(*mvcpu.CPUT, chan AgentReqT) bool // 16-bit implementation - may be the same as fn
}

// System Call Types as per Chap 2 of Sys Call Dictionary
const (
	scMemory = iota
	scProcess
	scFileManage
	scFileIO
	scDebugging
	scWindowing
	scMultitasking
	scIPC
	scConnection
	scMultiproc
	scClass
	scSystem
	scUserDev
	scBisync
	sc16Bit
)

const scReturn = 0310 // We need special access to this call number - it is handled differently

var syscalls = map[dg.WordT]syscallDescT{
	0:    {"?CREATE", "?CREA", scFileManage, nil, nil},
	1:    {"?DELETE", "?DELE", scFileManage, nil, nil},
	3:    {"?MEM", "?MEM", scMemory, scMem, scMem},
	014:  {"?MEMI", "?MEMI", scMemory, scMemi, scMemi},
	027:  {"?ILKUP", "?ILKU", scIPC, scIlkup, nil},
	036:  {"?GTOD", "?GTOD", scSystem, scGtod, scGtod},
	041:  {"?GDAY", "?GDAY", scSystem, scGday, scGday},
	070:  {"?PRIPR", "?PRIP", scProcess, scDummy, scDummy},
	072:  {"?GUNM", "?GUNM", scProcess, scGunm, nil},
	074:  {"?GHRZ", "?GHRZ", scSystem, scGhrz, scGhrz},
	0127: {"?DADID", "?DADI", scProcess, scDadid, scDadid},
	0157: {"?SINFO", "?SINF", scSystem, scInfo, nil},
	0166: {"?DACL", "?DACL", scFileManage, scDacl, nil},
	0263: {"?WDELAY", "?WDEL", scMultitasking, scWdelay, nil},
	0265: {"?LEFE", "?LEFE", scUserDev, scLefe, scLefe},
	0300: {"?OPEN", "?OPEN", scFileIO, scOpen, scOpen16},
	0301: {"?CLOSE", "?CLOS", scFileIO, scClose, nil},
	0302: {"?READ", "?READ", scFileIO, nil, scRead16},
	0303: {"?WRITE", "?WRIT", scFileIO, scWrite, scWrite16},
	0312: {"?GCHR", "?GCHR", scFileIO, scGchr, scGchr},
	0330: {"?EXEC", "?EXEC", scSystem, scExec, nil},
	0307: {"?GTMES", "?GTME", scSystem, scGtmes, scGtmes16},
	0415: {"?GECHR", "?GECH", scFileIO, scGechr, nil},
	0500: {"?TASK", "?TASK", scMultitasking, nil, nil},
	0503: {"?PRI", "?PRI", scMultitasking, scDummy, scDummy},
	0505: {"?KILAD", "?KILA", scMultitasking, scDummy, nil},
	0527: {"?DRSCH", "?DRSC", scMultitasking, scDummy, scDummy}, // Suspend all other tasks
	0542: {"?IFPU", "?IFPU", scMultitasking, scIfpu, scIfpu},
	0573: {"?SYSPRV", "?SYSP", scProcess, scSysprv, nil},
	0576: {"?XPSTAT", "?XPST", scProcess, scXpstat, nil},
}

// syscall redirects System Call according to the syscalls map
func syscall(callID dg.WordT, agent chan AgentReqT, cpu *mvcpu.CPUT) (ok bool) {
	call, defined := syscalls[callID]
	if !defined {
		log.Panicf("ERROR: System call No. %#o not yet defined at PC=%#x", callID, cpu.GetPC())
	}
	if call.fn == nil {
		log.Panicf("ERROR: System call No. %s not yet implemented at PC=%#x", call.name, cpu.GetPC())
	}
	if cpu.GetDebugLogging() {
		log.Printf("%s System Call...\n", call.name)
	}
	return call.fn(cpu, agent)
}

// syscall16 redirects a 16-bit System Call according to the syscalls map
func syscall16(callID dg.WordT, agent chan AgentReqT, cpu *mvcpu.CPUT) (ok bool) {
	call, defined := syscalls[callID]
	if !defined {
		log.Panicf("ERROR: System call No. %#o not yet defined at PC=%#x", callID, cpu.GetPC())
	}
	if call.fn16 == nil {
		log.Panicf("ERROR: 16-bit System call No. %s not yet implemented at PC=%#x", call.name, cpu.GetPC())
	}
	if cpu.GetDebugLogging() {
		log.Printf("%s System Call (16-bit)...\n", call.name)
	}
	return call.fn16(cpu, agent)
}

// readPacket just loads a chunk of memory into a slice of words
// TODO maybe this should be in ram_virtual.go as 'ReadWords' for efficiency?
func readPacket(addr dg.PhysAddrT, pktLen int) (pkt []dg.WordT) {
	pkt = make([]dg.WordT, pktLen, pktLen)
	for w := range pkt {
		pkt[w] = memory.ReadWord(addr + dg.PhysAddrT(w))
	}
	return pkt
}

// readBytes reads characters from memory up to the first NUL from the given doubleword byte address
func readBytes(bpAddr dg.DwordT, pc dg.PhysAddrT) []byte {
	buff := bytes.NewBufferString("")
	lobyte := (bpAddr & 0x0001) == 1
	wdAddr := dg.PhysAddrT(bpAddr>>1) | (pc & 0x7000_0000)
	c := memory.ReadByte(wdAddr, lobyte)
	for c != 0 {
		buff.WriteByte(byte(c))
		if lobyte {
			wdAddr++
		}
		lobyte = !lobyte
		c = memory.ReadByte(wdAddr, lobyte)
	}
	return buff.Bytes()
}

// readString reads characters from memory up to the first NUL from the given doubleword byte address
func readString(bpAddr dg.DwordT, pc dg.PhysAddrT) string {
	buff := bytes.NewBufferString("")
	lobyte := (bpAddr & 0x0001) == 1
	wdAddr := dg.PhysAddrT(bpAddr>>1) | (pc & 0x7000_0000)
	c := memory.ReadByte(wdAddr, lobyte)
	for c != 0 {
		buff.WriteByte(byte(c))
		if lobyte {
			wdAddr++
		}
		lobyte = !lobyte
		c = memory.ReadByte(wdAddr, lobyte)
	}
	return buff.String()
}

// writeBytes writes the whole byte array into memory at the given doubleword byte address
func writeBytes(bpAddr dg.DwordT, pc dg.PhysAddrT, arr []byte) {
	lobyte := (bpAddr & 0x0001) == 1
	wdAddr := dg.PhysAddrT(bpAddr>>1) | (pc & 0x7000_0000)
	for c := 0; c < len(arr); c++ {
		memory.WriteByte(wdAddr, lobyte, dg.ByteT(arr[c]))
		if lobyte {
			wdAddr++
		}
		lobyte = !lobyte
	}
}

// scDummy is a stub func for sys calls we are ignoring for now
func scDummy(cpu *mvcpu.CPUT, agentChan chan AgentReqT) bool {
	return true
}
