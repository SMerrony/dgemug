// agFileIO.go - 'Agent' Portion of File I/O System Call Emulation

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
	"os"

	"github.com/SMerrony/dgemug/dg"
)

type agCloseReqT struct {
	chanNo dg.WordT
}
type agCloseRespT struct {
	errCode dg.WordT
}

func agFileClose(req agCloseReqT) (resp agCloseRespT) {
	if req.chanNo == 0 {
		resp.errCode = 0
	} else {
		log.Panicf("ERROR: real file closing not yet implemented")
	}
	return resp
}

type agOpenReqT struct {
	PID  int
	path string
	mode dg.WordT
}
type agOpenRespT struct {
	channelNo dg.WordT
	ac0       dg.DwordT
}

func agFileOpen(req agOpenReqT) (resp agOpenRespT) {
	resp.ac0 = 0
	log.Printf("DEBUG: Agent received File Open request for %s\n", req.path)
	if req.path == "@CONSOLE" || req.path == "@OUTPUT" || req.path == "@INPUT" {
		resp.channelNo = consoleChan
		return resp
	}
	var (
		fp     *os.File
		flags  int
		err    error
		agChan agChannelT
	)
	agChan.path = req.path
	// parse creation options
	switch {
	case (req.mode&ofcr != 0) && (req.mode&ofce != 0):
		// delete, then create before open - i.e. truncate
		flags |= os.O_TRUNC
	case req.mode&ofcr != 0:
		// create new file before open
		flags |= os.O_CREATE | os.O_EXCL
	case req.mode&ofce != 0:
		// create if it doesn't already exist
		flags |= os.O_CREATE
	}
	// parse R/W options
	switch {
	case req.mode&ofin != 0 && req.mode&ofot == 0:
		flags |= os.O_RDONLY
		agChan.read = true
	case req.mode&ofot != 0 && req.mode&ofin == 0:
		flags |= os.O_WRONLY
		agChan.write = true
	case req.mode&ofio != 0:
		flags |= os.O_RDWR
	}
	// append?
	if req.mode&apnd != 0 {
		flags |= os.O_APPEND
	}
	if req.path[0] != ':' && perProcessData[req.PID].virtualRoot != "" {
		log.Printf("DEBUG: Attempting to Open file: %s\n", perProcessData[req.PID].virtualRoot+"/"+req.path)
		fp, err = os.OpenFile(perProcessData[req.PID].virtualRoot+"/"+req.path, flags, 0755)
	} else {
		fp, err = os.OpenFile(req.path, flags, 0755)
	}
	if err != nil {
		resp.ac0 = erfad
		return resp
	}
	agChan.file = fp
	newChan := len(agChannels) + 1
	agChannels[newChan] = &agChan

	return resp
}

type agReadReqT struct {
	chanNo   int
	length   int
	readLine bool
}
type agReadRespT struct {
	data []byte
}

func agFileRead(req agReadReqT) (resp agReadRespT) {
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
	chanNo int
	bytes  []byte
}
type agWriteRespT struct {
	bytesTxfrd dg.WordT
}

func agFileWrite(req agWriteReqT) (resp agWriteRespT) {
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
