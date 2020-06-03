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

// pseudo-Agent function calls...
const (
	agentAllocatePID = iota
	agentAllocateTID
	agentCreateIPC
	agentFileClose
	agentFileOpen
	agentFileRead
	agentFileRecreate
	agentFileWrite
	agentGetChars
	agentGetMessage
	agentIlkup
	agentSharedOpen
	agentSharedRead
)

// AgentReqT is the type of messages passed to and from the pseudo-agent
type AgentReqT struct {
	action   int
	reqParms interface{}
	result   interface{}
}

const (
	maxPID = 255
)

type perProcessDataT struct {
	invocationArgs []string
	virtualRoot    string
	sixteenBit     bool
	name           string
	conn           io.ReadWriteCloser // stream I/O port for proc's CONSOLE
	tidsInUse      [maxTasksPerProc]bool
}

type perTaskDataT struct {
	priority dg.WordT
	tid      dg.WordT
	startPC  dg.PhysAddrT
	initAC2  dg.DwordT
	wsb      dg.PhysAddrT
	wsfh     dg.PhysAddrT
	wssz     dg.DwordT
}

// agChannelT holds status of a file opened by the Agent for a user proc
type agChannelT struct {
	openerPID    int
	path         string
	isConsole    bool
	read, write  bool
	forShared    bool     // indicated this has been ?SOPENed
	recordLength int      // default I/O record length set at ?OPEN time
	conn         net.Conn // stream I/O
	file         *os.File // file I/O
}

type agIPCT struct {
	ownerPID     int
	name         string
	localPortNo  int
	globalPortNo int
	spool        chan []byte
}

// fixed channel #s for @CONSOLE, @INPUT and @OUTPUT
const (
	consoleChan = 0
	inputChan   = 1
	outputChan  = 2
)

var (
	pidInUse       [maxPID]bool
	perProcessData = map[int]perProcessDataT{}
	// uniqueTIDs     []uint16
	agChannels = map[int]*agChannelT{
		consoleChan: {path: "@CONSOLE", isConsole: true, read: true, write: true, forShared: false, recordLength: -1, conn: nil, file: nil},
		inputChan:   {path: "@INPUT", isConsole: true, read: true, write: true, forShared: false, recordLength: -1, conn: nil, file: nil},
		outputChan:  {path: "@OUTPUT", isConsole: true, read: true, write: true, forShared: false, recordLength: -1, conn: nil, file: nil},
	}
	agIPCs = map[string]*agIPCT{} // key is unique pathname
)

// StartAgent fires of the pseudo-agent Goroutine and returns its msg channel
func StartAgent(conn net.Conn) chan AgentReqT {
	// fake some in-use PIDs so they are not used
	pidInUse[0], pidInUse[1], pidInUse[2], pidInUse[3], pidInUse[4] = true, true, true, true, true
	agentChan := make(chan AgentReqT) // unbuffered to serialise requests
	agChannels[consoleChan].conn = conn
	agChannels[inputChan].conn = conn
	agChannels[outputChan].conn = conn

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
		case agentAllocateTID:
			request.result = agAllocateTID(request.reqParms.(agAllocateTIDReqT))
		case agentCreateIPC:
			request.result = agCreateIPC(request.reqParms.(agCreateIPCReqT))
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
		case agentIlkup:
			request.result = agIlkup(request.reqParms.(agIlkupReqT))
		case agentSharedOpen:
			request.result = agSharedOpen(request.reqParms.(agSharedOpenReqT))
		case agentSharedRead:
			request.result = agSharedRead(request.reqParms.(agSharedReadReqT))
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

func getNextFreeTID(PID int) (TID uint8, ok bool) {
	ppd := perProcessData[PID]
	for t := 1; t < maxTasksPerProc; t++ { // Zero TID is invalid
		if !ppd.tidsInUse[t] {
			ppd.tidsInUse[t] = true
			return uint8(t), true
		}
	}
	return 0, false // all TIDs in use
}

type agAllocatePIDReqT struct {
	invocationArgs []string
	virtualRoot    string
	sixteenBit     bool
	name           string
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
	perProcessData[resp.PID] = perProcessDataT{
		invocationArgs: req.invocationArgs,
		virtualRoot:    req.virtualRoot,
		sixteenBit:     req.sixteenBit,
		name:           req.name,
	}
	logging.DebugPrint(logging.ScLog, "AGENT assigned PID %d  Name: %s Args: %v\n", resp.PID, req.name, req.invocationArgs)
	if req.sixteenBit {
		logging.DebugPrint(logging.ScLog, "----- 16-bit program type\n")
	} else {
		logging.DebugPrint(logging.ScLog, "----- 32-bit program type\n")
	}
	return resp
}

type agAllocateTIDReqT struct {
	PID int
}
type agAllocateTIDRespT struct {
	uniqueTID   uint16
	tsw         dg.WordT
	standardTID uint8
	priority    dg.WordT
}

func agAllocateTID(req agAllocateTIDReqT) (resp agAllocateTIDRespT) {
	var ok bool
	resp.standardTID, ok = getNextFreeTID(req.PID)
	if !ok {
		return resp
	}
	resp.uniqueTID = uint16(req.PID)<<8 | uint16(resp.standardTID)
	resp.tsw = 0      // TODO
	resp.priority = 0 // TODO
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
	case gtsw:
		// TODO faked for now
		logging.DebugPrint(logging.ScLog, "WARNING: Faking empty ?GTSW response to ?GTMES system call\n")
		resp.ac0 = 0xffff_ffff
		resp.ac1 = 0
	case gsws:
		// TODO faked for now
		logging.DebugPrint(logging.ScLog, "WARNING: Faking empty ?GSWS response to ?GTMES system call\n")
		resp.ac0 = 0
		resp.ac1 = 0
	default:
		log.Panicf("ERROR: ?GTMES request type %#x not yet supported\n", req.greq)
	}
	logging.DebugPrint(logging.ScLog, "?GTMES returning %s\n", resp.result)
	return resp
}

// type agUidstatReqT struct {
// 	PID, TID int
// }
// type agUidstatRespT struct {
// 	uniqueTID   uint16
// 	tsw         dg.WordT
// 	standardTID uint8
// 	priority    dg.WordT
// }
// func agUidstat(req agUidstatReqT) (resp agUidstatRespT) {

// }
