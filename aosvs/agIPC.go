// +build virtual !physical

// agIPC.go - 'Agent' Portion of IPC System Call Emulation

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
	"github.com/SMerrony/dgemug/dg"
	"github.com/SMerrony/dgemug/logging"
)

type agCreateIPCReqT struct {
	PID         int
	filename    string
	localPortNo int
	ACL         string
}

func agCreateIPC(req agCreateIPCReqT) (errCode dg.WordT) {
	path := perProcessData[req.PID].virtualRoot + "/" + req.filename
	if _, found := agIPCs[path]; found {
		logging.DebugPrint(logging.ScLog, "\t?CREATE called for extant IPC file %s\n", path)
		errCode = ernae
	} else {
		var agIPC agIPCT
		agIPC.ownerPID = req.PID
		agIPC.localPortNo = req.localPortNo
		agIPC.name = req.filename
		agIPC.globalPortNo = req.PID<<16 | req.localPortNo // TODO should really include ring #
		agIPCs[path] = &agIPC
		logging.DebugPrint(logging.ScLog, "\t?CREATEd virtual IPC file %s\n", path)
	}
	return errCode
}

type agIlkupReqT struct {
	PID      int
	filename string
}
type agIlkupRespT struct {
	globalPortNo int
	ipcType      int
	errCode      int
}

func agIlkup(req agIlkupReqT) (resp agIlkupRespT) {
	path := perProcessData[req.PID].virtualRoot + "/" + req.filename
	agIPC, found := agIPCs[path]
	logging.DebugPrint(logging.ScLog, "\tChecking for virtual IPC %s\n", path)
	if !found {
		resp.errCode = erfde
		logging.DebugPrint(logging.ScLog, "\tIPC Lookup failed\n")
	} else {
		resp.globalPortNo = agIPC.globalPortNo
		resp.ipcType = fipc // TODO ???
		logging.DebugPrint(logging.ScLog, "\tIPC Lookup succeeded: %v\n", agIPC)
	}
	return resp
}
