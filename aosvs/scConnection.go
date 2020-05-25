// scConnection.go - 'Connection Management'-related System Call Emulation

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

import "github.com/SMerrony/dgemug/logging"

func scCon(p syscallParmsT) bool {
	ac1 := p.cpu.GetAc(1)
	if ac1&mcpid != 0 {
		// ac0 is a b.p. to a process name
		serverName := readString(p.cpu.GetAc(0), p.cpu.GetPC())
		logging.DebugPrint(logging.ScLog, "----- Faking connection to proc name %s\n", serverName)
	} else {
		// ac0 is a PID
		serverPID := int(p.cpu.GetAc(0))
		logging.DebugPrint(logging.ScLog, "----- Faking connection to PID %d\n", serverPID)
	}

	return true
}
