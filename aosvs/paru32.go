// paru32.go - Go version of parts of AOS/VS PARU.32.SR definitions file

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
	"github.com/SMerrony/dgemug/dg"
)

const (
	// System Constants

	// System Implementation Numbers
	saos   = 1    // AOS SYSTEM
	savs   = 2    // AOS/VS SYSTEM
	smps   = 3    // MP/OS SYSTEM
	sr32   = 4    // RT/32 SYSTEM
	sdvs   = 5    // AOS/DVS SYSTEM
	savsii = 6    // AOS/VS II SYSTEM
	sin    = savs // THIS IS AOS/VS
	// 6-7 RESERVED

	l32b = 0 // 32 BIT USER FLAG
	l16b = 1 // 16 BIT USER FLAG

	dfll = 136  // DEFAULT MAXIMUM LINE LENGTH
	mxpl = 256  // MAX PATHNAME LENGTH (BYTES)
	mxpn = 16   // MAX PROCESS NAME LENGTH (BYTES)
	mxun = 16   // MAX USERNAME LENGTH (BYTES)
	mxhn = 32   // MAX HOST NAME LENGTH (BYTES)
	mxul = 1024 // BUFFER FOR USER COMMANDS
	mxfs = mxpl // MAX LENGTH OF fedfunc STRING

	mxfp  = mxhn + mxun + mxpn // MAX FULL PROCESS NAME
	mxfn  = 32                 // MAX FILENAME LENGTH (BYTES)
	mxacl = 256                // MAX ACL LENGTH (BYTES)
	// ALL OF THE ABOVE INCLUDE THE
	// TRAILING NULL BYTE

	lowpid = 255     // HIGHEST AOS/VS LEGAL LOW PID
	hipid  = 4095    // HIGHEST AOS/VS LEGAL PID
	mxpid  = 64      // HIGHEST AOS LEGAL PID(NOT USED AOS/VS)
	vspids = 32767   // HIGHEST AOS/VS LEGAL PID (FOR INFOS USE)
	mxipc  = 1024    // MAX INTERHOST IPC LENGTH
	mxio   = 1024    // MAX INTERHOST I/O XFER LENGTH
	mxpsl  = 8       // MAX # PATHNAMES IN A SEARCHLIST
	midsc  = -9      // RESERVED
	hmsk   = 0177000 // RESERVED

	// ENTRY TYPE RANGES

	smin = 0        // SYSTEM MINIMUM
	smax = 63       // SYSTEM MAXIMUM
	dmin = smax + 1 // DGC MINIMUM
	dmax = 127      // DGC MAXIMUM
	umin = dmax + 1 // USER MINIMUM
	umax = 255      // USER MAXIMUM

	//       LOCATIONS DEFINED IN PHYSICAL BLOCK 0 OF THE .PR FILE

	preswsz = 0377 // PAGE SIZE OF EXTENSIBLE SWAPFILE REQUESTED
	// FOR RUNNING THIS PROGRAM.  IF 0 ( DEFAULT )
	// DEFAULT SWAPFILE SIZE IS USED.

	//       USER STATUS TABLE (UST) TEMPLATE

	ust = 0400 // START OF USER STATUS AREA 256.

	ustez = 0         // EXTENDED VARIABLE  WORD COUNT 256.  = 0
	ustes = ustez + 1 // EXTENDED VARIABLE PAGE 0 START 257. = 1
	ustss = ustes + 1 // SYMBOLS START 258.                  = 2
	ustse = ustss + 2 // SYMBOLS END 260.                    = 4
	ustda = ustse + 2 // DEB ADDR OR -1 262.                 = 6
	ustrv = ustda + 2 // REVISION OF PROGRAM 264.          = 010
	usttc = ustrv + 2 // NUMBER OF TASKS (1 TO 32.) 266.   = 012
	ustbl = usttc + 1 // # IMPURE BLKS 267.                = 013
	ustst = ustbl + 3 // SHARED STARTING BLK # 270.        = 016
	// USTST IS USTBL+3 BECAUSE THE 16. BIT USER'S
	// USTOD IS HIDDEN UNDERNEATH
	ustit = ustst + 2  // INTERRUPT ADDRESS 272.                = 020
	ustsz = ustit + 2  // SHARED SIZE IN BLKS 274.              = 022
	ustpr = ustsz + 2  // PROGRAM FILE TYPE (16 OR 32 BIT) 276. = 024
	ustsh = ustpr + 5  // PHYSICAL STARTING PAGE OF SHARED AREA IN .PR 281. = 031
	usten = ustpr + 21 // END OF USER UST
	ustpl = usten + 6  // PROGRAM LOCALITY

	// 		// USTPR FLAGS
	//   ust16=  1B0     // 16 BIT PROGRAM TYPE
	//   ust32=  0B0     // 32 BIT PROGRAM TYPE

	//   ustpa=  0B15    // PID SIZE TYPE 'SMALLPID' (<256)
	//   ustpb=  2B15    // PID SIZE TYPE 'HYBRID' (<256)
	//   ustpc=  3B15    // PID SIZE TYPE 'ANYPID' (>256)
	//
	// // TASK STATUS BITS (RETURNED BY tidstat CALL)

	//   tspn=  1B0     // TASK PENDED
	//   tssg=  1B1     // WAITING FOR .XMTW/.REC
	//   tssp=  1B2     // SUSPENDED
	//   tsrc=  1B3     // WAITING FOR TRCON
	//   tsov=  1B4     // WAITING FOR OVERLAY
	//   tswp=  1B5     // UNPEND TASK VIA WDPOP
	//   tsgs=  1B6     // TASK PENDED DUE TO AGENT SYNCHRONIZATION
	//   tsab=  1B7     // PENDED AWAITING AGENT ABORT PROCESSING
	//   tstl=  1B8     // PENDED AWAITING tunlock FROM ANOTHER TASK
	//   tsyg=  1B9     // TASK HAS BEEN signled (NOT A PEND BIT)
	//   tsdr=  1B10    // TASK PENDED FROM drsch
	//   tslk=  1B11    // TASK PENDED ON A flock REQUEST
	//   tsxr=  1B12    // TASK PENDED ON XMT OR REC
	//   twsg=  1B13    // TASK wtsignl PENDED
	//   tsut=  1B14    // AWAITING RETURN FROM USER utsk CODE
	//   tsuk=  1B15    // AWAITING RETURN FROM USER ukil CODE

	// PACKET FOR TASK DEFINITION (task)

	dlnk   = 0          // NON-ZERO = SHORT PACKET, ZERO = EXTENDED
	dlnl   = dlnk + 1   // LOWER PORTION OF dlnk
	dlnkb  = dlnl + 1   // BACKWARDS LINK (UPPER PORTION)
	dlnkbl = dlnkb + 1  // BACKWARDS LINK (LOWER PORTION)
	dpri   = dlnkbl + 1 // PRIORITY, ZERO TO USE CALLER'S
	did    = dpri + 1   // I.D., ZERO FOR NONE
	dpc    = did + 1    // STARTING ADDRESS OR RESOURCE ENTRY
	dpcl   = dpc + 1    // LOWER PORTION OF dpc
	dac2   = dpcl + 1   // INITIAL AC2 CONTENTS
	dcl2   = dac2 + 1   // LOWER PORTION OF dac2
	dstb   = dcl2 + 1   // STACK BASE, MINUS ONE FOR NO STACK
	dstl   = dstb + 1   // LOWER PORTION OF dstb
	dsflt  = dstl + 1   // STACK FAULT ROUTINE ADDR OR -1 IF SAME AS CURRENT
	dssz   = dsflt + 1  // STACK SIZE, IGNORED IF NO STACK
	dssl   = dssz + 1   // LOWER PORTION OF dssz
	dflgs  = dssl + 1   // FLAGS
	dfl0   = 1 << 15    //1B0     // RESERVED FOR SYSTEM
	dflrc  = 1 << 14    //1B1     // RESOURCE CALL TASK
	dfl15  = 1          //1B15    // RESERVED FOR SYSTEM
	dres   = dflgs + 1  // RESERVED FOR SYSTEM
	dnum   = dres + 1   // NUMBER OF TASKS TO CREATE

	dslth = dnum + 1 // LENGTH OF SHORT PACKET

	dsh  = dnum + 1 // STARTING HOUR, -1 IF IMMEDIATE
	dsms = dsh + 1  // STARTING SECOND IN HOUR, IGNORED IF IMMEDIATE
	dcc  = dsms + 1 // NUMBER OF TIMES TO CREATE TASK(S)
	dci  = dcc + 1  // CREATION INCREMENT  IN SECONDS

	dxlth = dci + 1 // LENGTH OF EXTENDED PACKET

	// BIT POINTER TO TASK DEF BITS

	dfbrc = dflgs*16 + 1 // RESOURCE CALL

	//  GENERAL USER I/O PACKET
	//
	//        USED FOR ?OPEN/?READ/?WRITE/?CLOSE
	//
	ich  dg.PhysAddrT = 0             // CHANNEL NUMBER
	isti              = ich + 1       // STATUS WORD (IN)
	isto              = isti + 1      // RIGHT=FILE TYPE, LEFT=RESERVED
	imrs              = isto + 1      // PHYSICAL RECORD SIZE - 1 (BYTES)
	ibad              = imrs + 1      // BYTE POINTER TO BUFFER
	ibal              = ibad + 1      // LOW ORDER BITS OF ?IBAD
	ires              = ibal + 1      // RESERVED
	ircl              = ires + 1      // RECORD LENGTH
	irlr              = ircl + 1      // RECORD LENGTH (RETURNED)
	irnw              = irlr + 1      // RESERVED
	irnh              = irnw + 1      // RECORD NUMBER (HIGH)
	irnl              = irnh + 1      // RECORD NUMBER (LOW)
	ifnp              = irnl + 1      // BYTE POINTER TO FILE NAME
	ifnl              = ifnp + 1      // LOW ORDER BITS OF ?IFNP
	idel              = ifnl + 1      // DELIMITER TABLE ADDRESS
	idll              = idel + 1      // LOWER BITS OF ?IDEL
	iosz int          = int(idll) + 1 // LENGTH OF STANDARD I/O PACKET

	//  ?ISTI FLAGS: BIT DEFINITIONS
	iplb = 0  // PACKET LENGTH BIT (0 => SHORT PACKET)
	icfb = 1  // CHANGE FORMAT BIT (0 => DEFAULT)
	icdm = 1  // DUMP MODE BIT (ON ?CLOSE ONLY)
	iptb = 2  // POSITIONING TYPE (0 => RELATIVE)
	ibib = 3  // BINARY I/O
	ifob = 4  // FORCE OUTPUT
	ioex = 5  // EXCLUSIVE OPEN
	iips = 6  // IPC NO WAIT BIT
	pdlm = 7  // PRIORITY REQUEST
	apbt = 8  // OPEN FILE FOR APPENDING
	of1b = 9  // OPEN TYPE BIT 1
	of2b = 10 // OPEN TYPE BIT 2
	opib = 11 // OPEN FOR INPUT
	opob = 12 // OPEN FOR OUTPUT
	rf1b = 13 // RECORD FORMAT BIT 1
	rf2b = 14 // RECORD FORMAT BIT 2
	rf3b = 15 // RECORD FORMAT BIT 3

	//  ?ISTI FLAGS: MASK DEFINITIONS
	ipkl = 0x8000 >> iplb // EXTENDED PACKET (IF SET)
	icrf = 0x8000 >> icfb // CHANGE RECORD FORMAT (IF SET)
	cdmp = 0x8000 >> icdm // SET DUMP BIT (ONLY ON ?CLOSE)
	ipst = 0x8000 >> iptb // RECORD POSITIONING TYPE (1 - ABSOLUTE)
	ibin = 0x8000 >> ibib // BINARY I/O
	ifop = 0x8000 >> ifob // FORCE OUTPUT
	iexo = 0x8000 >> ioex // EXCLUSIVE OPEN
	iipc = 0x8000 >> iips // IPC NO WAIT BIT
	pdel = 0x8000 >> pdlm // PRIORITY OPEN-I/O
	apnd = 0x8000 >> apbt // OPEN FILE FOR APPENDING
	ofcr = 0x8000 >> of1b // ATTEMPT CREATE BEFORE OPEN
	ofce = 0x8000 >> of2b // CORRECT ERROR ON CREATE OR OPEN
	ofin = 0x8000 >> opib // OPEN FOR INPUT
	ofot = 0x8000 >> opob // OPEN FOR OUTPUT
	ofio = ofin + ofot    // OPEN FOR INPUT AND OUTPUT
)

// FLAGS FOR RETURN TO CLI (?RETURN)
const (
	Rfcf = 1 << 7 // 1B0             // CLI FORMAT
	Rfwa = 1 << 5 // 1B2             // WARNING (SEVERITY=1)
	Rfer = 2 << 5 // 2B2             // ERROR   (SEVERITY=2)
	Rfab = 3 << 5 // 3B2             // ABORT   (SEVERITY=3)
	Rfec = 1 << 4 // 1B3             // ERROR CODE FLAG. IF SET, AC0 CONTAINS ERROR CODE
)
