// process.go - abstraction of an AOS/VS process

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
	"errors"
	"fmt"
	"io/ioutil"
	"net"

	"github.com/SMerrony/dgemug/dg"
)

// ProcessT represents an AOS/VS Process which will contain one or more Tasks
type ProcessT struct {
	pid             int
	name            string
	programFileName string
	tasks           []taskT
	console         net.Conn
}

// CreateProcess creates, but does not start, an emulated AOS/VS Process
func CreateProcess(pid int, prName string, con net.Conn) (proc *ProcessT, err error) {
	userProgramBytes, err := ioutil.ReadFile(prName)
	if err != nil {
		return nil, fmt.Errorf("could not read PR file %s", prName)
	}
	if len(userProgramBytes) == 0 {
		return nil, errors.New("empty PR file supplied")
	}

	userProgramWords := make([]dg.WordT, len(userProgramBytes)/2)
	for word := 0; word < len(userProgramBytes)/2; word++ {
		userProgramWords[word] = dg.WordT(userProgramBytes[word*2])<<8 | dg.WordT(userProgramBytes[word*2+1])
	}
	proc.pid = pid
	proc.programFileName = prName
	proc.console = con
	return proc, nil
}

// Run launches a new AOS/VS process and its initial Task
func (proc *ProcessT) Run() {

}
