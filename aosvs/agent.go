// agent.go - provides some agent-like serveices

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
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/SMerrony/dgemug/dg"
	"github.com/SMerrony/dgemug/logging"
)

const (
	agentAllocatePID = iota
	agentFileClose
	agentFileOpen
	agentFileRead
	agentFileRecreate
	agentFileWrite
	agentGetChars
	agentGetMessage
)

// AgentReqT is the type of messages passed to and from the pseudo-agent
type AgentReqT struct {
	action   int
	reqParms interface{}
	result   interface{}
}

const maxPID = 255

type perProcessDataT struct {
	invocationArgs []string
	virtualRoot    string
}

type agChannelT struct {
	path        string
	isConsole   bool
	read, write bool
	rwc         io.ReadWriteCloser // stream I/O
	file        *os.File           // file I/O
}

const consoleChan = 0

var (
	pidInUse       [maxPID]bool
	perProcessData = map[int]perProcessDataT{}
	agChannels     = map[int]*agChannelT{
		consoleChan: {path: "@CONSOLE", isConsole: true, read: true, write: true, rwc: nil, file: nil}, // @CONSOLE is always available
	}
)

// StartAgent fires of the pseudo-agent Goroutine and returns its msg channel
func StartAgent(conn net.Conn) chan AgentReqT {
	// fake some in-use PIDs so they are not used
	pidInUse[0], pidInUse[1], pidInUse[2], pidInUse[3], pidInUse[4] = true, true, true, true, true
	agentChan := make(chan AgentReqT) // unbuffered to serialise requests
	agChannels[0].rwc = conn

	go agentHandler(agentChan)
	return agentChan
}

func agentHandler(agentChan chan AgentReqT) {
	defer func() {
		if r := recover(); r != nil {
			logging.DebugLogsDump("logs/")
			os.Exit(1)
		}
	}()
	for {
		request := <-agentChan
		switch request.action {
		case agentAllocatePID:
			request.result = agAllocatePID(request.reqParms.(agAllocatePIDReqT))
		case agentFileClose:
			request.result = agFileClose(request.reqParms.(agCloseReqT))
		case agentFileOpen:
			request.result = agFileOpen(request.reqParms.(agOpenReqT))
		case agentFileRead:
			request.result = agFileRead(request.reqParms.(agReadReqT))
		case agentFileRecreate:
			request.result = agFileRecreate(request.reqParms.(agRecreateReqT))
		case agentFileWrite:
			request.result = agFileWrite(request.reqParms.(agWriteReqT))
		case agentGetChars:
			request.result = agGetChars(request.reqParms.(agGchrReqT))
		case agentGetMessage:
			request.result = agGetMessage(request.reqParms.(agGtMesReqT))

		default:
			log.Panicf("ERROR: Agent received unknown request type %d\n", request.action)
		}
		agentChan <- request
	}
}

func getNextFreePID() (pid int, ok bool) {
	for p := 1; p < maxPID; p++ {
		if !pidInUse[p] {
			pidInUse[p] = true
			return p, true
		}
	}
	return 0, false // all PIDs in use
}

type agAllocatePIDReqT struct {
	invocationArgs []string
	virtualRoot    string
}
type agAllocatePIDRespT struct {
	PID int
	ok  bool
}

func agAllocatePID(req agAllocatePIDReqT) (resp agAllocatePIDRespT) {
	resp.PID, resp.ok = getNextFreePID()
	if !resp.ok {
		return resp
	}
	perProcessData[resp.PID] = perProcessDataT{invocationArgs: req.invocationArgs, virtualRoot: req.virtualRoot}
	log.Printf("DEBUG: AGENT assigned PID %d  Args: %v\n", resp.PID, req.invocationArgs)
	return resp
}

type agGchrReqT struct {
	PID         int
	getDefaults bool // otherwise get current
	useChan     bool // otherwise use name
	devChan     dg.WordT
	devName     string
}
type agGchrRespT struct {
	words [3]dg.WordT
}

func agGetChars(req agGchrReqT) (resp agGchrRespT) {

	return resp
}

type agGtMesReqT struct {
	PID  int
	greq dg.WordT
	gnum dg.WordT
	gsw  dg.DwordT
}
type agGtMesRespT struct {
	ac0, ac1 dg.DwordT
	result   string
}

func agGetMessage(req agGtMesReqT) (resp agGtMesRespT) {
	switch req.greq {
	case gmes: // get entire message
		first := true
		for _, arg := range perProcessData[req.PID].invocationArgs {
			if first {
				first = false
			} else {
				resp.result += " "
			}
			resp.result += strings.ToUpper(arg)
		}
		resp.ac0 = gfcf
		resp.ac1 = dg.DwordT(len(resp.result)) >> 1 // words not bytes
	case gcmd: // get a parsed version of the command line
		first := true
		for _, arg := range perProcessData[req.PID].invocationArgs {
			if first {
				first = false
			} else {
				resp.result += ","
			}
			resp.result += strings.ToUpper(arg)
		}
		resp.ac1 = dg.DwordT(len(resp.result))
	case gcnt:
		resp.ac0 = dg.DwordT(len(perProcessData[req.PID].invocationArgs) - 1)
	case garg: // get the nth arg - special handing for integers
		if int(req.gnum) > len(perProcessData[req.PID].invocationArgs)-1 {
			log.Panicf("ERROR: ?GTMES attempted to retrieve non-extant argument no. %d.", req.gnum)
		}
		argS := perProcessData[req.PID].invocationArgs[int(req.gnum)]
		i, err := strconv.ParseInt(argS, 10, 16)
		if err == nil { // integer-only case
			resp.ac1 = dg.DwordT(i)
			resp.ac0 = dg.DwordT(len(argS))
		} else {
			resp.result = argS + "\x00"
			resp.ac0 = dg.DwordT(len(argS))
		}
	// case gtsw:
	case gsws:
		// TODO faked for now
		log.Println("WARNING: Faking empty ?GSWS response to ?GTMES system call")
		resp.ac0 = 0
		resp.ac1 = 0
	default:
		log.Panicf("ERROR: ?GTMES request type %#x not yet supported\n", req.greq)
	}
	return resp
}
