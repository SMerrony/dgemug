// debugLogs.go
// Copyright ©2018-2020  Steve Merrony

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

package logging

import (
	"fmt"
	"os"
	"sync"
)

const (
	numLogs          = 7
	numDebugLogLines = 100000 // each circular buffer contains this many lines
	logPerms         = 0644
)

// Log IDs
const (
	DebugLog = iota // DebugLog is the general-purpose log
	MtLog           // MtLog is for the type 6026 MT tape module
	DkpLog          // DkpLog is for the type 4231a Moving-Head Disk
	DpfLog          // DpfLog is for the type 6061 DPF disk module
	DskpLog         // DskpLog is for the type 6239 DSKP disk module
	MapLog          // MapLog is for BMC/DCH-related logging
	ScLog           // ScLog is for System Call logging in the VS emulator
)

var logs = map[int]string{
	DebugLog: "debug.log",
	MtLog:    "mt_debug.log",
	DkpLog:   "dkp_debug.log",
	DpfLog:   "dpf_debug.log",
	DskpLog:  "dskp_debug.log",
	MapLog:   "bmcdch_debug.log",
	ScLog:    "syscall_debug.log",
}

var (
	// N.B. I tried using strings.Builder for the logs with Go 1.10, it seemed to use c.1000x more heap
	logArr    [numLogs][numDebugLogLines]string // the stored log messages
	logArrMu  [numLogs]sync.Mutex
	firstLine [numLogs]int // pointer to the first line of each log
	lastLine  [numLogs]int // pointer to the last line of each log
)

// DebugLogsDump can be called to dump out each of the non-empty debug logs to text files
// @dir should be empty, or a /-terminated subdirectory to receive the logs
func DebugLogsDump(dir string) {

	var (
		debugDumpFile *os.File
	)

	for l := range logArr {
		logArrMu[l].Lock()
		if firstLine[l] != lastLine[l] { // ignore unused or empty logs
			debugDumpFile, _ = os.OpenFile(dir+logs[l], os.O_WRONLY|os.O_CREATE|os.O_TRUNC, logPerms)
			debugDumpFile.WriteString(">>> Dumping Debug Log\n\n")
			thisLine := firstLine[l]
			for thisLine != lastLine[l] {
				debugDumpFile.WriteString(logArr[l][thisLine])
				thisLine++
				if thisLine == numDebugLogLines {
					thisLine = 0
				}
			}
			debugDumpFile.WriteString(logArr[l][thisLine])
			debugDumpFile.WriteString("\n>>> Debug Log Ends\n")
			debugDumpFile.Close()
		}
		logArrMu[l].Unlock()
	}
}

// DebugPrint doesn't print anything!  It stores the log message
// in array-backed circular arrays
// for printing when debugLogsDump() is invoked.
// This func can be called very often, KISS...
func DebugPrint(log int, aFmt string, msg ...interface{}) {

	logArrMu[log].Lock()

	lastLine[log]++

	// end of log array?
	if lastLine[log] == numDebugLogLines {
		lastLine[log] = 0 // wrap-around
	}

	// has the tail hit the head of the circular buffer?
	if lastLine[log] == firstLine[log] {
		firstLine[log]++ // advance the head pointer
		if firstLine[log] == numDebugLogLines {
			firstLine[log] = 0 // but reset if at limit
		}
	}

	// sprintf the given message to tail of the specified log
	logArr[log][lastLine[log]] = fmt.Sprintf(aFmt, msg...)
	logArrMu[log].Unlock()
}
