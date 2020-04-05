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
	agentFileClose = iota
	agentFileOpen
	agentFileRead
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

type agChannelT struct {
	path        string
	isConsole   bool
	read, write bool
	rwc         io.ReadWriteCloser
}

const consoleChan = 0

var agChannels = map[dg.WordT]*agChannelT{ // TODO should probably not be at this scope
	consoleChan: {"@CONSOLE", true, true, true, nil}, // @CONSOLE is always available
}

var invocationArgs []string

// StartAgent fires of the pseudo-agent Goroutine and returns its msg channel
func StartAgent(conn net.Conn, args []string) chan AgentReqT {
	invocationArgs = args
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
		case agentFileClose:
			request.result = doAgentFileClose(request.reqParms.(agCloseReqT))
		case agentFileOpen:
			request.result = doAgentFileOpen(request.reqParms.(agOpenReqT))
		case agentFileRead:
			request.result = doAgentFileRead(request.reqParms.(agReadReqT))
		case agentFileWrite:
			request.result = doAgentFileWrite(request.reqParms.(agWriteReqT))
		case agentGetChars:
			request.result = doAgentGetChars(request.reqParms.(agGchrReqT))
		case agentGetMessage:
			request.result = doAgentGetMessage(request.reqParms.(agGtMesReqT))

		default:
			log.Panicf("ERROR: Agent received unknown request type %d\n", request.action)
		}
		agentChan <- request
	}
}

type agCloseReqT struct {
	chanNo dg.WordT
}
type agCloseRespT struct {
	errCode dg.WordT
}

func doAgentFileClose(req agCloseReqT) (resp agCloseRespT) {
	if req.chanNo == 0 {
		resp.errCode = 0
	} else {
		log.Panicf("ERROR: real file closing not yet implemented")
	}
	return resp
}

type agOpenReqT struct {
	path string
	mode dg.WordT
}
type agOpenRespT struct {
	channelNo dg.WordT
}

func doAgentFileOpen(req agOpenReqT) (resp agOpenRespT) {
	log.Printf("DEBUG: Agent received File Open request for %s\n", req.path)
	if req.path == "@CONSOLE" || req.path == "@OUTPUT" || req.path == "@INPUT" {
		resp.channelNo = 0
	} else {
		log.Panicf("ERROR: real file opening not yet implemented")
	}
	return resp
}

type agReadReqT struct {
	chanNo   dg.WordT
	length   int
	readLine bool
}
type agReadRespT struct {
	data []byte
}

func doAgentFileRead(req agReadReqT) (resp agReadRespT) {
	agChan, isOpen := agChannels[req.chanNo]
	if isOpen {
		if agChan.isConsole {
			if !req.readLine {
				log.Panic("ERROR: Fixed-length Input from @CONSOLE not yet implemented")
			}
			buff := make([]byte, 0)
			for {
				oneByte := make([]byte, 1, 1)
				l, err := agChan.rwc.Read(oneByte)
				if err != nil {
					log.Panic("ERROR: Could not read from @CONSOLE")
				}
				if l == 0 {
					log.Panic("ERROR: ?READ got 0 bytes from @CONSOLE")
				}
				if oneByte[0] == dg.ASCIINL || oneByte[0] == '\n' || oneByte[0] == '\r' {
					break
				}
				// TODO DELete
				buff = append(buff, oneByte[0])
			}
			resp.data = buff
		} else {
			log.Panicf("ERROR: real file reading not yet implemented")
		}
	} else {
		log.Panic("ERROR: attempt to ?READ from unopened file")
	}
	log.Printf("?READ returning <%v>\n", resp.data)
	return resp
}

type agWriteReqT struct {
	chanNo dg.WordT
	bytes  []byte
}
type agWriteRespT struct {
	bytesTxfrd dg.WordT
}

func doAgentFileWrite(req agWriteReqT) (resp agWriteRespT) {
	agChan, isOpen := agChannels[req.chanNo]
	if isOpen {
		if agChan.isConsole {
			n, err := agChan.rwc.Write(req.bytes)
			if err != nil {
				log.Panic("ERROR: Could not write to @CONSOLE")
			}
			resp.bytesTxfrd = dg.WordT(n)
		}
	} else {
		log.Panic("ERROR: attempt to ?WRITE to unopened file")
	}
	return resp
}

type agGchrReqT struct {
	getDefaults bool // otherwise get current
	useChan     bool // otherwise use name
	devChan     dg.WordT
	devName     string
}
type agGchrRespT struct {
	words [3]dg.WordT
}

func doAgentGetChars(req agGchrReqT) (resp agGchrRespT) {

	return resp
}

type agGtMesReqT struct {
	greq dg.WordT
	gnum dg.WordT
	gsw  dg.DwordT
}
type agGtMesRespT struct {
	ac0, ac1 dg.DwordT
	result   string
}

func doAgentGetMessage(req agGtMesReqT) (resp agGtMesRespT) {
	switch req.greq {
	case gmes: // get entire message
		first := true
		for _, arg := range invocationArgs {
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
		for _, arg := range invocationArgs {
			if first {
				first = false
			} else {
				resp.result += ","
			}
			resp.result += strings.ToUpper(arg)
		}
		resp.ac1 = dg.DwordT(len(resp.result))
	case gcnt:
		resp.ac0 = dg.DwordT(len(invocationArgs) - 1)
	case garg: // get the nth arg - special handing for integers
		if int(req.gnum) > len(invocationArgs)-1 {
			log.Panicf("ERROR: ?GTMES attempted to retrieve non-extant argument no. %d.", req.gnum)
		}
		argS := invocationArgs[int(req.gnum)]
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
