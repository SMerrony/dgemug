// scFileIO.go - File I/O System Call Emulation

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
	"log"
	"strings"

	"github.com/SMerrony/dgemug/dg"
	"github.com/SMerrony/dgemug/logging"
	"github.com/SMerrony/dgemug/memory"
)

func scClose(p syscallParmsT) bool {
	pktAddr := dg.PhysAddrT(p.cpu.GetAc(2))
	channel := memory.ReadWord(pktAddr + ich)
	var creq = agCloseReqT{int(channel)}
	var areq = AgentReqT{agentFileClose, creq, nil}
	p.agentChan <- areq
	areq = <-p.agentChan
	if areq.result.(agCloseRespT).errCode != 0 {
		p.cpu.SetAc(0, dg.DwordT(areq.result.(agCloseRespT).errCode))
		return false
	}
	return true
}

func scGchr(p syscallParmsT) bool {
	var areq AgentReqT
	var gchrReq agGchrReqT
	if memory.TestDwbit(p.cpu.GetAc(1), 0) {
		// AC0 should contain a channel #
		if memory.TestDwbit(p.cpu.GetAc(1), 1) {
			// get default chars
			gchrReq = agGchrReqT{p.PID, true, true, dg.WordT(p.cpu.GetAc(0)), ""}
		} else {
			// get current chars
			gchrReq = agGchrReqT{p.PID, false, true, dg.WordT(p.cpu.GetAc(0)), ""}
		}
	} else {
		// AC0 should contain BP to device name
		bpPathname := p.cpu.GetAc(0)
		path := strings.ToUpper(readString(bpPathname, p.ringMask))
		if memory.TestDwbit(p.cpu.GetAc(1), 1) {
			// get default chars
			gchrReq = agGchrReqT{p.PID, true, false, 0, path}
		} else {
			// get current chars
			gchrReq = agGchrReqT{p.PID, false, false, 0, path}
		}
	}
	areq.action = agentGetChars
	areq.reqParms = gchrReq
	p.agentChan <- areq
	areq = <-p.agentChan
	wrAddr := dg.PhysAddrT(p.cpu.GetAc(2)) | p.ringMask
	memory.WriteWord(dg.PhysAddrT(wrAddr), areq.result.(agGchrRespT).words[0])
	memory.WriteWord(dg.PhysAddrT(wrAddr+1), areq.result.(agGchrRespT).words[1])
	memory.WriteWord(dg.PhysAddrT(wrAddr+2), areq.result.(agGchrRespT).words[2])
	return true
}

func scGechr(p syscallParmsT) bool {

	return true
}

func scSchr(p syscallParmsT) bool {

	return true
}

func scGopen(p syscallParmsT) bool {
	filename := readString(p.cpu.GetAc(0), p.cpu.GetPC())
	logging.DebugPrint(logging.ScLog, "----- Filename: %s\n", filename)
	panic("NYI")
	return true
}

func scOpen(p syscallParmsT) bool {
	pktAddr := dg.PhysAddrT(p.cpu.GetAc(2)) | (p.ringMask)
	fileSpec := memory.ReadWord(pktAddr + isti)
	fileType := memory.ReadWord(pktAddr + isto)
	// blockSize := memory.ReadWord(pktAddr+imrs)
	recLen := int(int16(memory.ReadWord(pktAddr + ircl)))
	bpPathname := memory.ReadDWord(pktAddr + ifnp)
	path := strings.ToUpper(readString(bpPathname, p.ringMask))
	logging.DebugPrint(logging.ScLog, "?OPEN Pathname: %s, Type: %#x, FileSpec: %#x, RecLen: %d.\n", path, fileType, fileSpec, recLen)
	var areq AgentReqT
	var openReq = agOpenReqT{p.PID, path, fileSpec, recLen}
	areq.action = agentFileOpen
	areq.reqParms = openReq
	p.agentChan <- areq
	areq = <-p.agentChan
	if areq.result.(agOpenRespT).ac0 != 0 {
		p.cpu.SetAc(0, areq.result.(agOpenRespT).ac0)
		return false
	}
	memory.WriteWord(pktAddr+ich, areq.result.(agOpenRespT).channelNo)
	logging.DebugPrint(logging.ScLog, "----- Returned channel # %d\n", areq.result.(agOpenRespT).channelNo)
	return true
}

func scOpen16(p syscallParmsT) bool {
	pktAddr := dg.PhysAddrT(p.cpu.GetAc(2)) | (p.ringMask)
	fileSpec := memory.ReadWord(pktAddr + isti16)
	fileType := memory.ReadWord(pktAddr + isto16)
	// blockSize := memory.ReadWord(pktAddr+imrs16)
	recLen := int(int16(memory.ReadWord(pktAddr + ircl16)))
	bpPathname := dg.DwordT(memory.ReadWord(pktAddr + ifnp16))
	path := strings.ToUpper(readString(bpPathname, p.ringMask))
	logging.DebugPrint(logging.ScLog, "?OPEN (16-bit) Pathname: %s, Type: %#x, FileSpec: %#x, RecLen: %d.\n", path, fileType, fileSpec, recLen)
	var areq AgentReqT
	var openReq = agOpenReqT{p.PID, path, fileSpec, recLen}
	areq.action = agentFileOpen
	areq.reqParms = openReq
	p.agentChan <- areq
	areq = <-p.agentChan
	if areq.result.(agOpenRespT).ac0 != 0 {
		p.cpu.SetAc(0, areq.result.(agOpenRespT).ac0)
		return false
	}
	memory.WriteWord(pktAddr+ich16, areq.result.(agOpenRespT).channelNo)
	return true
}

func scRead(p syscallParmsT) bool {
	pktAddr := dg.PhysAddrT(p.cpu.GetAc(2))
	channel := int(memory.ReadWord(pktAddr + ich))
	specs := memory.ReadWord(pktAddr + isti)
	length := int(int16(memory.ReadWord(pktAddr + ircl)))
	dest := memory.ReadDWord(pktAddr + ibad)
	readLine := (memory.ReadWord(pktAddr+isti) & ibin) != 0
	if specs&ipkl != 0 {
		if memory.TestDwbit(memory.ReadDWord(pktAddr+etsp), 0) {
			smPktAddr := dg.PhysAddrT(memory.ReadDWord(pktAddr+etsp) & 0x7fff_ffff)
			memory.WriteDWord(pktAddr+etsp, dg.DwordT(smPktAddr))
			flagWd := memory.ReadWord(smPktAddr)
			if flagWd != 0 {
				logging.DebugPrint(logging.ScLog, "\t?ESSE: %v\n", flagWd&esse != 0)
				logging.DebugPrint(logging.ScLog, "\t?ESRD: %v\n", flagWd&esrd != 0)
				logging.DebugPrint(logging.ScLog, "\t?ESNR: %v\n", flagWd&esnr != 0)
				logging.DebugPrint(logging.ScLog, "\t?ESED: %v\n", flagWd&esed != 0)
				logging.DebugPrint(logging.ScLog, "\t?ESCP: %v\n", flagWd&escp != 0)
				logging.DebugPrint(logging.ScLog, "\t?ESDD: %v\n", flagWd&esdd != 0)
				logging.DebugPrint(logging.ScLog, "\t?ESRP: %v\n", flagWd&esrp != 0)
				logging.DebugPrint(logging.ScLog, "\t?ESNE: %v\n", flagWd&esne != 0)
				logging.DebugPrint(logging.ScLog, "\t?ESGT: %v\n", flagWd&esgt != 0)
				logging.DebugPrint(logging.ScLog, "\t?ESBE: %v\n", flagWd&esbe != 0)
				logging.DebugPrint(logging.ScLog, "\t?ESPE: %v\n", flagWd&espe != 0)
				logging.DebugPrint(logging.ScLog, "\t?ESEP: %#x\n", memory.ReadWord(smPktAddr+1))
				logging.DebugPrint(logging.ScLog, "\t?ESCR: %#x\n", memory.ReadWord(smPktAddr+2))
				logging.DebugPrint(logging.ScLog, "\tExtended packet ignored...\n")
			} else {
				logging.DebugPrint(logging.ScLog, "\tFlag word is zero - ignoring\n")
			}
		}
	}
	logging.DebugPrint(logging.ScLog, "?READ (32-bit) Channel: %#x, Specs: %#x, Bytes: %#x, Dest: %#x, Line Mode: %v\n", channel, specs, length, dest, readLine)
	var readReq = agReadReqT{channel, specs, length, readLine}
	var areq = AgentReqT{agentFileRead, readReq, nil}
	p.agentChan <- areq
	areq = <-p.agentChan
	resp := areq.result.(agReadRespT)
	memory.WriteWord(pktAddr+irlr, dg.WordT(len(resp.data)))
	writeBytes(dest, p.ringMask, resp.data)
	if resp.ac0 != 0 {
		p.cpu.SetAc(0, dg.DwordT(resp.ac0))
	}
	return true
}

func scRead16(p syscallParmsT) bool {
	pktAddr := dg.PhysAddrT(p.cpu.GetAc(2)) | p.ringMask
	channel := int(memory.ReadWord(pktAddr + ich16))
	specs := memory.ReadWord(pktAddr + isti16)
	length := int(int16(memory.ReadWord(pktAddr + ircl16)))
	dest := dg.DwordT(memory.ReadWord(pktAddr + ibad16))
	readLine := (memory.ReadWord(pktAddr+isti16) & ibin) == 0
	if specs&ipkl != 0 {
		log.Panic("ERROR: ?READ (16-bit) extended packet not yet implemented")
	}
	logging.DebugPrint(logging.ScLog, "?READ (16-bit) Channel: %#x, Specs: %#x, Bytes: %#x, Dest: %#x, Line Mode: %v\n", channel, specs, length, dest, readLine)
	var readReq = agReadReqT{channel, specs, length, readLine}
	var areq = AgentReqT{agentFileRead, readReq, nil}
	p.agentChan <- areq
	areq = <-p.agentChan
	resp := areq.result.(agReadRespT)
	memory.WriteWord(pktAddr+irlr16, dg.WordT(len(resp.data)))
	writeBytes(dest, p.ringMask, resp.data)
	return true
}

func scSend(p syscallParmsT) bool {
	msgLen := int(p.cpu.GetAc(2) & 0x00ff)
	msg := memory.ReadBytes(p.cpu.GetAc(1), p.cpu.GetPC(), msgLen)
	// flag := memory.GetDwbits(p.cpu.GetAc(2), 22, 2)
	// switch flag {
	// case 0x00: // AC0 contains a PID
	// case 0x01: // AC0 is a bp to a proc name
	// case 0x02: // AC0 is a bp to a console name
	// case 0x03: // undefined
	// }
	agWrite := agWriteReqT{0, false, false, true, int(int16(msgLen)), msg, 0}
	areq := AgentReqT{agentFileWrite, agWrite, nil}
	logging.DebugPrint(logging.ScLog, "\tWriting <%s> to @CONSOLE\n", string(msg))
	p.agentChan <- areq
	areq = <-p.agentChan
	return true
}

func scWrite(p syscallParmsT) bool {
	pkt := dg.PhysAddrT(p.cpu.GetAc(2))
	channel := int(memory.ReadWord(pkt + ich))
	specsWd := memory.ReadWord(pkt + isti)
	extendedPkt := specsWd&ipkl != 0
	absPositioning := specsWd&ipst != 0
	dataSens := specsWd&rtds != 0
	recLen := int(int16(memory.ReadWord(pkt + ircl)))
	byteslice := memory.ReadBytes(memory.ReadDWord(pkt+ibad), p.ringMask, int(recLen))
	if dataSens {
		// this is done here to avoid passing excess data into the Agent
		maxLen := recLen
		if recLen == -1 {
			maxLen = agChannels[channel].recordLength
		}
		var tooLong bool
		byteslice, tooLong = getDataSensitivePortion(byteslice, maxLen)
		if tooLong {
			p.cpu.SetAc(0, dg.DwordT(erltl))
			return false
		}
	}
	position := int32(memory.ReadDWord(pkt + irnh))
	var writeReq = agWriteReqT{channel, extendedPkt, absPositioning, dataSens, recLen, byteslice, position}
	var areq = AgentReqT{agentFileWrite, writeReq, nil}
	p.agentChan <- areq
	areq = <-p.agentChan
	if areq.result.(agWriteRespT).errCode != 0 {
		p.cpu.SetAc(0, dg.DwordT(areq.result.(agWriteRespT).errCode))
		return false
	}
	memory.WriteWord(pkt+irlr, areq.result.(agWriteRespT).bytesTxfrd)
	return true
}

func scWrite16(p syscallParmsT) bool {
	pkt := dg.PhysAddrT(p.cpu.GetAc(2)) | p.ringMask
	channel := int(memory.ReadWord(pkt + ich16))
	specsWd := memory.ReadWord(pkt + isti16)
	extendedPkt := specsWd&ipkl != 0
	absPositioning := specsWd&ipst != 0
	dataSens := specsWd&rtds != 0
	recLen := int(int16(memory.ReadWord(pkt + ircl16)))
	byteslice := memory.ReadBytes(dg.DwordT(memory.ReadWord(pkt+ibad16)), p.ringMask, int(recLen))
	if dataSens {
		maxLen := recLen
		if recLen == -1 {
			maxLen = agChannels[channel].recordLength
		}
		var tooLong bool
		byteslice, tooLong = getDataSensitivePortion(byteslice, maxLen)
		if tooLong {
			p.cpu.SetAc(0, dg.DwordT(erltl))
			return false
		}
	}
	position := int32(memory.ReadDWord(pkt + irnh16))
	var writeReq = agWriteReqT{channel, extendedPkt, absPositioning, specsWd&rtds != 0, recLen, byteslice, position}
	var areq = AgentReqT{agentFileWrite, writeReq, nil}
	p.agentChan <- areq
	areq = <-p.agentChan
	memory.WriteWord(pkt+irlr16, areq.result.(agWriteRespT).bytesTxfrd)
	return true
}
