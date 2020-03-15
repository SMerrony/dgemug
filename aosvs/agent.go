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
	"log"

	"github.com/SMerrony/dgemug/dg"
)

const (
	agentFileOpen = iota
	agentFileClose
)

// AgentReqT is the type of messages passed to and from the pseudo-agent
type AgentReqT struct {
	action   int
	reqParms interface{}
	result   interface{}
}

type agOpenReqT struct {
	path string
	mode dg.WordT
}
type agOpenRespT struct {
	channelNo dg.WordT
}

type agChannelT struct {
	channelNo   dg.WordT
	isConsole   bool
	read, write bool
}

var agChannels = map[string]agChannelT{ // TODO should probably not be at this scope
	"@CONSOLE": {0, true, true, true}, // @CONSOLE is always available
}

// StartAgent fires of the pseudo-agent Goroutine and returns its msg channel
func StartAgent() chan AgentReqT {
	agentChan := make(chan AgentReqT) // unbuffered to serialise requests
	go agentHandler(agentChan)
	return agentChan
}

func agentHandler(agentChan chan AgentReqT) {
	for {
		request := <-agentChan
		switch request.action {
		case agentFileOpen:
			request.result = agentFileOpener(request.reqParms.(agOpenReqT))
		case agentFileClose:

		default:
			log.Fatalf("ERROR: Agent received unknown request type %d\n", request.action)
		}
		agentChan <- request
	}
}

func agentFileOpener(req agOpenReqT) (resp agOpenRespT) {
	log.Printf("DEBUG: Agent received File Open request for %s\n", req.path)
	agChan, isOpen := agChannels[req.path]
	if isOpen {
		resp.channelNo = agChan.channelNo
	} else {
		log.Fatalf("ERROR: real file opening not yet implemented")
	}
	return resp
}
