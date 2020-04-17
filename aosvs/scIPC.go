// scIPC.go - Inter-Process Communication System Call Emulation

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
	"strings"

	"github.com/SMerrony/dgemug/dg"
)

func scIlkup(p syscallParmsT) bool {
	bpPathname := p.cpu.GetAc(0)
	path := strings.ToUpper(readString(bpPathname, p.cpu.GetPC()))
	agIlkupReq := agIlkupReqT{p.PID, path}
	areq := AgentReqT{agentIlkup, agIlkupReq, nil}
	p.agentChan <- areq
	areq = <-p.agentChan
	resp := areq.result.(agIlkupRespT)
	if resp.errCode != 0 {
		p.cpu.SetAc(0, dg.DwordT(resp.errCode))
		return false
	}
	p.cpu.SetAc(1, dg.DwordT(resp.globalPortNo))
	p.cpu.SetAc(2, dg.DwordT(resp.ipcType))
	// // trap SWAT checking...
	// if path == "?10.SWAT.IPC" {
	// 	p.cpu.SetAc(0, erfde)
	// 	return false
	// }
	// log.Panicf("?ILKUP of %s not yet implemented", path)
	return true
}
