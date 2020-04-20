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
	"io"
	"log"
	"os"
	"strings"

	"github.com/SMerrony/dgemug/dg"
	"github.com/SMerrony/dgemug/logging"
)

type agCloseReqT struct {
	chanNo int
}
type agCloseRespT struct {
	errCode dg.WordT
}

func agFileClose(req agCloseReqT) (resp agCloseRespT) {
	if req.chanNo == 0 {
		resp.errCode = 0
	} else {
		logging.DebugPrint(logging.ScLog, "\tChannel # %d\n", req.chanNo)
		agChan, isOpen := agChannels[req.chanNo]
		if isOpen {
			if agChan.isConsole {
				logging.DebugPrint(logging.ScLog, "\tIgnoring ?CLOSE on console channel\n")
			} else {
				agChan.file.Close()
				delete(agChannels, req.chanNo)
				logging.DebugPrint(logging.ScLog, "\tFile closed\n")
			}
		} else {
			logging.DebugPrint(logging.ScLog, "\tFILE WAS NOT OPEN\n")
			resp.errCode = eracu
		}
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
	// TODO currently returning same channel for these common generic files, they might need separate ones...
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
		logging.DebugPrint(logging.ScLog, "\tAttempting to Open file: %s\n", perProcessData[req.PID].virtualRoot+"/"+req.path)
		fp, err = os.OpenFile(perProcessData[req.PID].virtualRoot+"/"+req.path, flags, 0755)
	} else {
		fp, err = os.OpenFile(req.path, flags, 0755)
	}
	if err != nil {
		resp.ac0 = erfad
		return resp
	}
	agChan.file = fp
	newChan := len(agChannels)
	agChannels[newChan] = &agChan
	resp.channelNo = dg.WordT(newChan)
	return resp
}

type agReadReqT struct {
	chanNo   int
	specs    dg.WordT
	length   int
	readLine bool
}
type agReadRespT struct {
	ac0  dg.WordT
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
			switch {
			case req.specs&ipst != 0:
				log.Panic("Absolute positining NYI")
			}
			buf := make([]byte, req.length)
			n, err := agChannels[req.chanNo].file.Read(buf)
			if n == 0 && err == io.EOF {
				resp.ac0 = ereof
			} else {
				resp.data = buf
			}
		}
	} else {
		log.Panic("ERROR: attempt to ?READ from unopened file")
	}
	logging.DebugPrint(logging.ScLog, "?READ returning <%v>\n", resp.data)
	return resp
}

type agRecreateReqT struct {
	PID         int
	aosFilename string
}
type agRecreateRespT struct {
	ok      bool
	errCode dg.DwordT
}

func agFileRecreate(req agRecreateReqT) (resp agRecreateRespT) {
	filename := strings.ReplaceAll(req.aosFilename, ":", "/") // convert any : to /
	if filename[0] == '@' {
		log.Panicf("ERROR: ?RECREATE in :PER not yet implemented (file: %s)", filename)
	}
	if filename[0] != '/' {
		filename = perProcessData[req.PID].virtualRoot + "/" + filename
		logging.DebugPrint(logging.ScLog, "\tResolved %s to %s\n", req.aosFilename, filename)
	}
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		resp.errCode = erfde
		resp.ok = false
	} else {
		os.Truncate(filename, 0)
		resp.ok = true
	}
	return resp
}

type agSharedOpenReqT struct {
	PID      int
	filename string
	readonly bool
}
type agSharedOpenRespT struct {
	ac0       dg.DwordT
	channelNo dg.DwordT
}

func agSharedOpen(req agSharedOpenReqT) (resp agSharedOpenRespT) {
	var (
		fp     *os.File
		flags  int
		err    error
		agChan agChannelT
	)
	agChan.path = req.filename
	if req.readonly {
		flags = os.O_RDONLY
	} else {
		flags = os.O_RDWR
	}
	if req.filename[0] != ':' && perProcessData[req.PID].virtualRoot != "" {
		logging.DebugPrint(logging.ScLog, "\tAttempting to SOpen file: %s\n", perProcessData[req.PID].virtualRoot+"/"+req.filename)
		fp, err = os.OpenFile(perProcessData[req.PID].virtualRoot+"/"+req.filename, flags, 0755)
	} else {
		fp, err = os.OpenFile(req.filename, flags, 0755)
	}
	if err != nil {
		resp.ac0 = erfad // TODO add more errors here
		return resp
	}
	agChan.file = fp
	newChan := len(agChannels)
	agChannels[newChan] = &agChan
	resp.channelNo = dg.DwordT(newChan)
	logging.DebugPrint(logging.ScLog, "\tReturning channel: %d.\n", newChan)
	return resp
}

type agSharedReadReqT struct {
	chanNo   int
	length   int
	startPos int64
}
type agSharedReadRespT struct {
	ac0  dg.DwordT
	data []byte
}

func agSharedRead(req agSharedReadReqT) (resp agSharedReadRespT) {
	_, isOpen := agChannels[req.chanNo]
	if isOpen {
		buf := make([]byte, req.length)
		logging.DebugPrint(logging.ScLog, "\tAttempting to Seek to byte: %d. on channel: %d.\n", req.startPos, req.chanNo)
		_, err := agChannels[req.chanNo].file.Seek(req.startPos, 0)
		if err != nil {
			log.Panicf("ERROR: ?SPAGE positioning failed: %v", err)
		}
		n, err := agChannels[req.chanNo].file.Read(buf)
		if n == 0 && err == io.EOF {
			// It looks as if we should create pages here - can't find it in the docs though...
			resp.data = make([]byte, req.length)
			logging.DebugPrint(logging.ScLog, "\tRead no bytes from channnel #%d., returning %d. empty bytes\n", req.chanNo, req.length)
		} else {
			resp.data = buf
			logging.DebugPrint(logging.ScLog, "\tRead %d. bytes from channnel #%d.\n", n, req.chanNo)
		}
	} else {
		log.Panic("ERROR: attempt to ?SPAGE from unopened file")
	}
	return resp
}

type agWriteReqT struct {
	channel    int
	isExtended bool
	isAbsolute bool
	recLen     int16
	bytes      []byte
	position   int32
}
type agWriteRespT struct {
	bytesTxfrd dg.WordT
}

func agFileWrite(req agWriteReqT) (resp agWriteRespT) {
	logging.DebugPrint(logging.ScLog, "----- Chan: %d., Extended: %v, Posn: %#x, Len: %d.\n", req.channel, req.isExtended, req.position, req.recLen)

	agChan, isOpen := agChannels[req.channel]
	if isOpen {
		if agChan.isConsole {
			resp.bytesTxfrd = dg.WordT(agWriteToUserConsole(agChan, req.bytes))
		}
	} else {
		log.Panic("ERROR: attempt to ?WRITE to unopened file")
	}
	return resp
}

func agWriteToUserConsole(agChan *agChannelT, b []byte) (n int) {
	n, err := agChan.rwc.Write(b)
	if err != nil {
		log.Panic("ERROR: Could not write to @CONSOLE")
	}
	logging.DebugPrint(logging.ScLog, "-----  wrote %d., bytes <%v> to @CONSOLE\n", n, b)
	return n
}
