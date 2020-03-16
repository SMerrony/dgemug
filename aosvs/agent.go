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

	"github.com/SMerrony/dgemug/dg"
)

const (
	agentFileOpen = iota
	agentFileWrite
	agentFileClose
)

// AgentReqT is the type of messages passed to and from the pseudo-agent
type AgentReqT struct {
	action   int
	reqParms interface{}
	result   interface{}
}

type agCloseReqT struct {
	chanNo dg.WordT
}
type agCloseRespT struct {
	errCode dg.WordT
}

type agOpenReqT struct {
	path string
	mode dg.WordT
}
type agOpenRespT struct {
	channelNo dg.WordT
}

type agWriteReqT struct {
	chanNo dg.WordT
	bytes  []byte
}
type agWriteRespT struct {
	bytesTxfrd dg.WordT
}

type agChannelT struct {
	path        string
	isConsole   bool
	read, write bool
	rwc         io.ReadWriteCloser
}

var agChannels = map[dg.WordT]*agChannelT{ // TODO should probably not be at this scope
	0: {"@CONSOLE", true, true, true, nil}, // @CONSOLE is always available
}

// StartAgent fires of the pseudo-agent Goroutine and returns its msg channel
func StartAgent(conn net.Conn) chan AgentReqT {
	agentChan := make(chan AgentReqT) // unbuffered to serialise requests
	agChannels[0].rwc = conn

	go agentHandler(agentChan)
	return agentChan
}

func agentHandler(agentChan chan AgentReqT) {
	for {
		request := <-agentChan
		switch request.action {
		case agentFileClose:
			request.result = agentFileCloser(request.reqParms.(agCloseReqT))
		case agentFileOpen:
			request.result = agentFileOpener(request.reqParms.(agOpenReqT))
		case agentFileWrite:
			request.result = agentFileWriter(request.reqParms.(agWriteReqT))

		default:
			log.Fatalf("ERROR: Agent received unknown request type %d\n", request.action)
		}
		agentChan <- request
	}
}

func agentFileCloser(req agCloseReqT) (resp agCloseRespT) {
	if req.chanNo == 0 {
		resp.errCode = 0
	} else {
		log.Fatalf("ERROR: real file opening not yet implemented")
	}
	return resp
}

func agentFileOpener(req agOpenReqT) (resp agOpenRespT) {
	log.Printf("DEBUG: Agent received File Open request for %s\n", req.path)
	if req.path == "@CONSOLE" {
		resp.channelNo = 0
	} else {
		log.Fatalf("ERROR: real file opening not yet implemented")
	}
	return resp
}

func agentFileWriter(req agWriteReqT) (resp agWriteRespT) {
	agChan, isOpen := agChannels[req.chanNo]
	if isOpen {
		if agChan.isConsole {
			n, err := agChan.rwc.Write(req.bytes)
			if err != nil {
				log.Fatal("ERROR: Could not write to @CONSOLE")
			}
			resp.bytesTxfrd = dg.WordT(n)
		}
	} else {
		log.Fatal("ERROR: attempt to ?WRITE to unopened file")
	}
	return resp
}
