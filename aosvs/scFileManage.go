// scFileManage.go - File Management System Call Emulation

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

func scCreate(p syscallParmsT) bool {
	pktAddr := dg.PhysAddrT(p.cpu.GetAc(2)) | (p.ringMask)
	bpFilename := p.cpu.GetAc(0)
	filename := strings.ToUpper(readString(bpFilename, p.ringMask))
	fileType := memory.ReadWord(pktAddr+cftyp) & 0x00ff
	switch fileType {
	case fipc:
		localPortNo := int(memory.ReadWord(pktAddr + cpor))
		// not handling ?CTIM
		var acl string
		bpACL := memory.ReadDWord(pktAddr + cacp)
		switch bpACL {
		case 0xffff_ffff:
			acl = "[DEFACL]"
		case 0:
			acl = ""
		default:
			acl = readString(bpACL, p.ringMask)
		}
		logging.DebugPrint(logging.ScLog, "----- IPC File: %s Local Port #: %d, ACL: %s\n", filename, localPortNo, acl)
		crIPCreq := agCreateIPCReqT{p.PID, filename, localPortNo, acl}
		areq := AgentReqT{agentCreateIPC, crIPCreq, nil}
		p.agentChan <- areq
		areq = <-p.agentChan
	default:
		log.Panicf("ERROR: ?CREATE called for %s - not yet implemented", filename)
	}

	return true
}

func scDacl(p syscallParmsT) bool {
	// TODO this is all faked...
	switch dg.WordT(p.cpu.GetAc(0)) { // make 16-bit safe
	case 0xffff:
		log.Panic("ERROR: Setting new DefACL not yet implemented in ?DACL")
	case 0:
		defacl := []byte("XYZZY")
		defacl = append(defacl, 0)
		defacl = append(defacl, faco+facw+faca+facw+face)
		defacl = append(defacl, 0)
		memory.WriteBytesBA(defacl, p.cpu.GetAc(1))
	case 1:
		log.Panic("ERROR: Turning off DefACL not yet implemented in ?DACL")
	}
	return true
}

func scGname(p syscallParmsT) bool {
	bpFilename := p.cpu.GetAc(0)
	filename := strings.ToUpper(readString(bpFilename, p.ringMask))
	if filename[0] == '@' {
		filename = ":PER:" + filename[1:] // convert @ to :PER
	}
	filename = strings.ReplaceAll(filename, "/", ":") // convert any / to :
	bpPathname := p.cpu.GetAc(1)
	writeBytes(bpPathname, p.ringMask, []byte(filename))
	p.cpu.SetAc(2, dg.DwordT(len(filename)))
	logging.DebugPrint(logging.ScLog, "?GNAME returning %s for %s\n", readString(bpFilename, p.ringMask), filename)
	return true
}

func scRecreate(p syscallParmsT) bool {
	bpFilename := p.cpu.GetAc(0)
	filename := strings.ToUpper(readString(bpFilename, p.ringMask))
	var recReq = agRecreateReqT{PID: p.PID, aosFilename: filename}
	var areq = AgentReqT{agentFileRecreate, recReq, nil}
	p.agentChan <- areq
	areq = <-p.agentChan
	resp := areq.result.(agRecreateRespT)
	if !resp.ok {
		p.cpu.SetAc(0, resp.errCode)
		return false
	}
	return true
}
