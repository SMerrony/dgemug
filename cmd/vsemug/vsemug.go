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
	"strconv"

	"github.com/SMerrony/dgemug/aosvs"
	"github.com/SMerrony/dgemug/logging"
	"github.com/SMerrony/dgemug/memory"
	"github.com/SMerrony/dgemug/mvcpu"
)

// we need the instructionDefinitions.go file to have been generated in mvcpu
//go:generate dginstr -action=makego -cputype=mv -csv=../dginstr/dginstrs.csv -go=../../mvcpu/instructionDefinitions.go

// program options - Change arg slicing in main if these are changed
var (
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

	memory.MemInit()
	mvcpu.InstructionsInit()
	args := make([]string, 1)
	// N.B. Change this if program invocation flags are changed...
	// Stripping path as slashes will confuse AOS/VS argument parsing
	// We are taking the virtual root from the path of the PR file for now
	var vRoot string
	if *consoleAddrFlag == "" {
		args[0] = filepath.Base(os.Args[4])
		vRoot = filepath.Dir(os.Args[4])
		args = append(args, os.Args[5:]...)
	} else {
		args[0] = filepath.Base(os.Args[2])
		vRoot = filepath.Dir(os.Args[2])
		args = append(args, os.Args[3:]...)
	}
	agentChan := aosvs.StartAgent(conn) // start the pseudo-Agent which will serialise syscalls in the process's tasks

	proc, err := aosvs.CreateProcess(args, vRoot, *prFlag, 7, conn, agentChan, debugLogging)
	if err != nil {
		exitNicely(conn, err.Error())
	}

	errorCode, termMessage, flags := proc.Run()
	switch flags {
	case aosvs.Rfwa:
		termMessage = "WARNING: " + termMessage
	case aosvs.Rfer:
		termMessage = "ERROR: " + termMessage
	case aosvs.Rfab:
		termMessage = "ABORT: " + termMessage
	}
	if flags&aosvs.Rfec != 0 {
		termMessage += "\nError Code: " + strconv.Itoa(int(errorCode))
	}
	exitNicely(conn, termMessage)
}

func exitNicely(con net.Conn, msg string) {
	con.Write([]byte(msg))
	con.Write([]byte("\n *** Exiting Emulator ***\n"))
	log.Println(msg)
	logging.DebugLogsDump("logs/")
	os.Exit(1)
}
