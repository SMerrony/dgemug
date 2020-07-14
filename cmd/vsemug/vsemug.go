// vsmemug project main src

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

// N.B. Build with "-tags virtual"

package main

import (
	"flag"
	"log"
	"net"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"

	"github.com/SMerrony/dgemug/aosvs"
	"github.com/SMerrony/dgemug/logging"
	"github.com/SMerrony/dgemug/memory"
	"github.com/SMerrony/dgemug/mvcpu"
)

// we need the instructionDefinitions.go file to have been generated in mvcpu
//go:generate dginstr -action=makego -cputype=mv -csv=../dginstr/dginstrs.csv -go=../../mvcpu/instructionDefinitions.go

// program options - Change arg slicing in main if these are changed
var (
	argsFlag        = flag.String("args", "", "arguments to pass to program (surround multiple args with double-quotes)")
	consoleAddrFlag = flag.String("consoleaddr", "localhost:10001", "network interface/port for @CONSOLE for 1st process, others will be assigned sequentially")
	prFlag          = flag.String("pr", "", "program to run at startup")
)

func main() {

	debugLogging := true // SLOWS execution dramatically

	flag.Parse()
	log.Printf("INFO: Waiting for terminal connection to %s\n", *consoleAddrFlag)
	l, err := net.Listen("tcp", *consoleAddrFlag)
	if err != nil {
		log.Println("ERROR: Could not listen on @CONSOLE port: ", err.Error())
		os.Exit(1)
	}
	defer l.Close()

	conn, err := l.Accept()
	if err != nil {
		log.Println("ERROR: Could not accept on @CONSOLE port: ", err.Error())
		os.Exit(1)
	}
	conn.Write([]byte("\n *** Welcome to the VSemuG AOS/VS Emulator ***" + "\n"))
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			exitNicely(conn, " *** VSemuG Internal Panic ***")
		}
	}()

	if *prFlag == "" {
		exitNicely(conn, "Please supply an initial PR file to run"+"\n")
	}

	// DEBUGGING ONLY... (exits in x minutes)
	go stopper(5, conn)

	memory.MemInit()
	mvcpu.InstructionsInit()
	args := make([]string, 1)
	// Stripping path as slashes will confuse AOS/VS argument parsing
	// We are taking the virtual root from the path of the PR file for now
	args[0] = filepath.Base(*prFlag)
	vRoot := filepath.Dir(*prFlag)
	if *argsFlag != "" {
		args = append(args, strings.Fields(*argsFlag)...)
	}

	agentChan := aosvs.StartAgent(conn) // start the pseudo-Agent which will serialise syscalls in the process's tasks

	err = aosvs.CreateProcess(args, vRoot, *prFlag, 7, conn, agentChan, debugLogging) // TODO - Eventually this should be a call to ?PROC
	if err != nil {
		exitNicely(conn, err.Error())
	}

	ppd := aosvs.PerProcessData[5]
	ppd.ActiveTasksWg.Wait()

	exitNicely(conn, "")
}

func exitNicely(con net.Conn, msg string) {
	con.Write([]byte(msg))
	con.Write([]byte("\n *** Exiting Emulator ***\n"))
	log.Println(msg)
	logging.DebugLogsDump("logs/")
	os.Exit(1)
}

// stopper will panic 5 minutes after it is started - for desparated debugging only!
func stopper(mins int, con net.Conn) {
	time.Sleep(time.Minute * time.Duration(mins))
	exitNicely(con, "STOPPING: Stopper func has timed-out")
}
