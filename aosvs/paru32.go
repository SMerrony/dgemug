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
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENTIN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package aosvs

import (
	"github.com/SMerrony/dgemug/dg"
)

const (
	// SYSTEM ERROR CODES
	ericm = 01   // ILLEGAL SYSTEM COMMAND
	erfno = 02   // CHANNEL NOT OPEN
	eropr = 03   // CHANNEL ALREADY OPEN
	ersal = 04   // SHARED I/O REQ NOT MAP SLOT ALIGNED
	ermem = 05   // INSUFFICIENT MEMORY AVAILABLE
	eradr = 06   // ILLEGAL STARTING ADDRESS
	erovn = 07   // ILLEGAL OVERLAY NUMBER
	ertim = 010  // ILLEGAL TIME ARGUMENT
	ernot = 011  // NO TASK CONTROL BLOCK AVAILABLE
	erxmt = 012  // SIGNAL TO ADDRESS ALREADY IN USE
	erqts = 013  // ERROR IN QTASK REQUEST
	ertid = 014  // TASK I.D. ERROR
	erdch = 015  // DATA CHANNEL MAP FULL
	ermpr = 016  // SYSTEM CALL PARAMETER ADDRESS ERROR
	erabt = 017  // TASK NOT FOUND FOR ABORT
	erirb = 020  // INSUFFICIENT ROOM IN BUFFER
	erspc = 021  // FILE SPACE EXHAUSTED
	ersft = 022  // USER STACK FAULT
	erdde = 023  // DIRECTORY DOES NOT EXIST
	erifc = 024  // ILLEGAL FILENAME CHARACTER
	erfde = 025  // FILE DOES NOT EXIST
	ernae = 026  // FILE NAME ALREADY EXISTS
	ernad = 027  // NON-DIRECTORY ARGUMENT IN PATHNAME
	ereof = 030  // END OF FILE
	erdid = 031  // DIRECTORY DELETE ERROR
	erwad = 032  // WRITE ACCESS DENIED
	errad = 033  // READ ACCESS DENIED
	erawd = 034  // APPEND AND/OR WRITE ACCESS DENIED
	ernmc = 035  // NO CHANNELS AVAILABLE
	ersrl = 036  // RELEASE OF NON-ACTIVE SHARED SLOT
	erprp = 037  // ILLEGAL PRIORITY
	erbmx = 040  // ILLEGAL MAX SIZE ON PROCESS CREATE
	erpty = 041  // ILLEGAL PROCESS TYPE
	ercon = 042  // CONSOLE DEVICE SPECIFICATION ERROR
	ernsw = 043  // SWAP FILE SPACE EXHAUSTED
	eribs = 044  // DEVICE ALREADY IN SYSTEM
	erdnm = 045  // ILLEGAL DEVICE CODE
	ershp = 046  // ERROR ON SHARED PARTITION SET
	errmp = 047  // ERROR ON REMAP CALL
	ergsg = 050  // ILLEGAL AGENT GATE CALL
	erprn = 051  // NUMBER OF PROCESSES EXCEEDS MAX
	ernef = 052  // IPC MESSAGE EXCEEDS BUFFER LENGTH
	erivp = 053  // INVALID PORT NUMBER
	ernms = 054  // NO MATCHING SEND
	ernor = 055  // NO OUTSTANDING RECEIVE
	eriop = 056  // ILLEGAL ORIGIN PORT
	eridp = 057  // ILLEGAL DESTINATION PORT
	ersen = 060  // INVALID SHARED LIBRARY REFERENCE
	erirl = 061  // ILLEGAL RECORD LENGTH SPECIFIED(=0)
	erarc = 062  // ATTEMPT TO RELEASE CONSOLE DEVICE
	erdai = 063  // DEVICE ALREADY IN USE
	eraru = 064  // ATTEMPT TO RELEASE UNASSIGNED DEVICE
	eracu = 065  // ATTEMPT TO CLOSE UNOPEN CHANNEL/DEVICE
	eritc = 066  // I/O TERMINATED BY CLOSE
	erltl = 067  // LINE TOO LONG
	erpar = 070  // PARITY ERROR
	erexc = 071  // RESDENT PROC TRIED TO PUSH (.EXEC)
	erndr = 072  // NOT A DIRECTORY
	ernsa = 073  // SHARED I/O REQUEST NOT TO SHARED AREA
	ersnm = 074  // ATTEMPT TO CREATE > MAX # SONS
	erfil = 075  // FILE READ ERROR
	erdto = 076  // DEVICE TIMEOUT
	eriot = 077  // WRONG TYPE I/O FOR OPEN TYPE
	erftl = 0100 // FILENAME TOO LONG
	erbof = 0101 // POSITIONING BEFORE BEGINNING OF FILE
	erprv = 0102 // CALLER NOT PRIVILEGED FOR THIS ACTION
	ersim = 0103 // SIMULTANEOUS REQUESTS ON SAME CHANNEL
	erift = 0104 // ILLEGAL FILE TYPE
	ernrd = 0105 // INSUFFICIENT ROOM IN DIRECTORY
	erilo = 0106 // ILLEGAL OPEN
	erprh = 0107 // ATTEMPT TO ACCESS PROC NOT IN HIERARCHY
	erblr = 0110 // ATTEMPT TO BLOCK UNBLOCKABLE PROC
	erpre = 0111 // INVALID SYSTEM CALL PARAMETER
	erges = 0112 // ATTEMPT TO START MULTIPLE AGENTS
	erciu = 0113 // CHANNEL IN USE
	ericb = 0114 // INSUFFICIENT CONTIGUOUS DISK BLOCKS
	ersto = 0115 // STACK OVERFLOW
	eribm = 0116 // INCONSISTENT BIT MAP DATA
	erbsz = 0117 // ILLEGAL BLOCK SIZE FOR DEVICE
	erxmz = 0120 // ATTEMPT TO XMT ILLEGAL MESSAGE
	erpuf = 0121 // PHYSICAL UNIT FAILURE
	erpwl = 0122 // PHYSICAL WRITE LOCK
	eruol = 0123 // PHYSICAL UNIT OFFLINE
	erioo = 0124 // ILLEGAL OPEN OPTION FOR FILE TYPE
	erndv = 0125 // TOO MANY OR TOO FEW DEVICE NAMES
	ermis = 0126 // DISK AND FILE SYS REV #'S DON'T MATCH
	eridd = 0127 // INCONSISTENT DIB DATA
	erild = 0130 // INCONSISTENT LD
	eridu = 0131 // INCOMPLETE LD
	eridt = 0132 // ILLEGAL DEVICE NAME TYPE
	erpdf = 0133 // ERROR IN PROCESS UST DEFINITION
	erviu = 0134 // LD IN USE, CANNOT RELEASE
	ersre = 0135 // SEARCH LIST RESOLUTION ERROR
	ercgf = 0136 // CAN'T GET IPC DATA FROM FATHER
	erilb = 0137 // ILLEGAL LIBRARY NUMBER GIVEN
	errfm = 0140 // ILLEGAL RECORD FORMAT
	erarg = 0141 // TOO MANY OR TOO FEW ARGUMENTS TO PMGR
	erigm = 0142 // ILLEGAL gtmes PARAMETERS
	ericl = 0143 // ILLEGAL CLI MESSAGE
	ermrd = 0144 // MESSAGE RECEIVE DISABLED
	ernac = 0145 // NOT A CONSOLE DEVICE
	ermil = 0146 // ATTEMPT TO EXCEED MAX INDEX LEVEL
	ericn = 0147 // ILLEGAL CHANNEL
	ernrr = 0150 // NO RECEIVER WAITING
	ersrr = 0151 // SHORT RECEIVE REQUEST
	ertin = 0152 // TRANSMITTER INOPERATIVE
	erunm = 0153 // ILLEGAL USER NAME
	eriln = 0154 // ILLEGAL LINK #
	erdpe = 0155 // DISK POSITIONING ERROR
	ertxt = 0156 // MSG TEXT LONGER THAN SPEC'D.
	erstr = 0157 // SHORT TRANSMISSION
	erhis = 0160 // ERROR ON HISTOGRAM INIT/DELETE
	erirv = 0161 // ILLEGAL RETRY VALUE
	erass = 0162 // ASSIGN ERROR - ALREADY YOUR DEVICE
	erpet = 0163 // MAG TAPE REQ PAST LOGICAL END OF TAPE
	ersts = 0164 // STACK TOO SMALL (task)
	ertmt = 0165 // TOO MANY TASKS REQUESTED (task)
	ersoc = 0166 // SPOOLER OPEN RETRY COUNT EXCEEDED
	eracl = 0167 // ILLEGAL ACL
	erwpb = 0170 // stmap BUFFER INVALID OR WRITE PROTECTED
	erinp = 0171 // IPC FILE NOT OPENED BY ANOTHER PROC
	erfpu = 0172 // FPU HARDWARE NOT INSTALLED
	erpnm = 0173 // ILLEGAL PROCESS NAME
	erpnu = 0174 // PROCESS NAME ALREADY IN USE
	erdct = 0175 // DISCONNECT ERROR (MODEM CONTROLLED)
	eripr = 0176 // NONBLOCKING PROC REQUEST ERROR
	ersni = 0177 // SYSTEM NOT INSTALLED
	erlvl = 0200 // MAX DIRECTORY TREE DEPTH EXCEEDED
	erroo = 0201 // RELEASING OUT-OF-USE OVERLAY
	errdl = 0202 // RESOURCE DEADLOCK
	ereo1 = 0203 // FILE IS OPEN, CAN'T EXCLUSIVE OPEN
	ereo2 = 0204 // FILE IS EXCLUSIVE OPENED, CAN'T OPEN
	eripd = 0205 // INIT PRIVILEGE DENIED
	ermim = 0206 // MULTIPLE imsg CALLS TO SAME DCT
	erlnk = 0207 // ILLEGAL LINK
	eridf = 0210 // ILLEGAL DUMP FORMAT
	erxna = 0211 // EXEC NOT AVAILABLE (MOUNT, ETC.)
	erxuf = 0212 // EXEC REQUEST FUNCTION UNKNOWN
	ereso = 0213 // ONLY EXEC'S SONS CAN DO THAT
	errbo = 0214 // REFUSED BY OPERATOR
	erwmt = 0215 // VOLUME NOT MOUNTED
	erisv = 0216 // ILLEGAL SWITCH VALUE (>65K DECIMAL)

	// THE NEXT FOUR ERROR CODES MUST BE CONTIGUOUSLY NUMBERED

	erifn = 0217 // INPUT FILE  DOES NOT EXIST
	erofn = 0220 // OUTPUT FILE  DOES NOT EXIST
	erlfn = 0221 // LIST FILE  DOES NOT EXIST
	erdfn = 0222 // DATA FILE DOES NOT EXIST

	ergfe = 0223 // RECURSIVE GENERIC FILE OPEN FAILURE
	ernmw = 0224 // NO MESSAGE WAITING
	ernud = 0225 // USER DATA AREA DOES NOT EXIST
	erdvc = 0226 // ILLEGAL DEVICE TYPE FROM VSGEN
	errst = 0227 // RESTART OF SYSTEM CALL
	erfur = 0230 // PROBABLY FATAL HARDWARE RUNTIME ERROR
	ercft = 0231 // USER COMMERCIAL STACK FAULT
	erfft = 0232 // USER FLOATING POINT STACK FAULT
	eruae = 0233 // USER DATA AREA ALREADY EXISTS
	eriso = 0234 // ILLEGAL SCREEN_EDIT REQUEST (PMGR)
	erddh = 0235 // "send" DESTINATION DEVICE HELD BY "^S"
	erovr = 0236 // DATA OVERRUN ERROR
	ercpd = 0237 // CONTROL POINT DIRECTORY MAX SIZE EXCEEDED
	ernsd = 0240 // SYS OR BOOT DISK NOT PART OF MASTER LD
	erusy = 0241 // UNIVERSAL SYSTEM, YOU CAN'T DO THAT
	eread = 0242 // EXECUTE ACCESS DENIED
	erfix = 0243 // CAN'T INIT LD, RUN FIXUP ON IT
	erfad = 0244 // FILE ACCESS DENIED
	erdad = 0245 // DIRECTORY ACCESS DENIED
	eriad = 0246 // ATTEMPT TO DEFINE > 1 SPECIAL PROC
	erind = 0247 // NO SPECIAL PROCESS IS DEFINED
	erpro = 0250 // ATTEMPT TO ISSUE MCA REQUEST WITH
	// DIRECT I/O IN PROGRESS
	erdio = 0251 // ATTEMPT TO ISSUE MCA DIRECT I/O WITH
	// OUTSTANDING REQUESTS
	erltk = 0252 // LAST TASK WAS KILLED
	erlrf = 0253 // RESOURCE LOAD OR RELEASE FAILURE
	ernnl = 0254 // ZERO LENGTH FILENAME SPECIFIED

	// FOLLOWING ARE AOS/VS ONLY ERROR CODES
	// THEY ARE IN A SEPARATE GROUP FROM THE AOS ERROR CODES
	// THEY ARE IN GROUP 77, HENCE CODES ARE 77*1000+X

	erxxx = 077 * 01000 // FIRST CODE FOR AOS/VS
	// HARDWARE POINTER VALIDATION ERRORS

	ervwp = erxxx     // INVALID ADDRESS PASSED AS SYSTEM CALL ARGUMENT
	ervbp = ervwp + 1 // INVALID BYTE POINTER PASSED AS SYS CALL ARGUMENT
	erdpt = ervbp + 1 // DIFFERENT TYPE PROCESS(32/16 BIT) WITHOUT PRIVILEGE
	erral = erdpt + 1 // RING ALREADY LOADED
	errni = erral + 1 // RING NUMBER INVALID
	errtb = errni + 1 // RING TOO BIG
	erwsm = errtb + 1 // SET WKG SET MIN, NOT PRIVILEGED
	ertne = erwsm + 1 // PMGR- TRACING NOT ENABLED
	ertae = ertne + 1 // PMGR- TRACING ALREADY ENABLED
	ernuf = ertae + 1 // PMGR- TRACING FILE NOT A USER DATA FILE
	errna = ernuf + 1 // PMGR- REQUESTOR NOT TRACING AUTHORIZED
	erpnl = errna + 1 // PMGR- PATHNAME LENGTH AREA
	ersnf = erpnl + 1 // SYMBOL NOT FOUND IN .ST FILE
	ersnr = ersnf + 1 // SOURCE NOT RESIDENT ON lmap
	erdnr = ersnr + 1 // DESTINATION NOT RESIDENT ON lmap
	eribp = erdnr + 1 // BKPT SEEN IN USER PROGRAM WHEN DEBUG NOT INIT'ED
	erbst = eribp + 1 // BAD SYMBOL TABLE FORMAT SEEN (gtna CALL)
	erpdo = erbst + 1 // PAGE FILE DIRECTORY OVERFLOW
	ermwt = erpdo + 1 // MORE THAN ONE WS TRACE DEFINED ON TARGET
	erhwt = ermwt + 1 // BOTH TRACE AND HISTOGRAM CALLED, OR > 1 TRACE
	erdtc = erhwt + 1 // DIFFERENT TYPE CHAIN
	erwst = erdtc + 1 // NO WS TRACE DEFINED ON THIS TARGET
	erwss = erwst + 1 // INVALID WKG SET MAX/MIN
	erwsb = erwss + 1 // INVALID WORKING SET TRACE BUFFER
	erwsf = erwsb + 1 // WORKING SET NOT SWAPPABLE
	erawm = erwsf + 1 // TRYING TO WIRE MORE PAGES THAN WS MAX
	ertpw = erawm + 1 // TOO MANY PAGES WIRED
	eracc = ertpw + 1 // ACCESS DENIED ON valad
	errnl = eracc + 1 // RING NOT LOADED
	ertal = errnl + 1 // TOO MANY ARGUMENTS ON LCALL
	erxbl = ertal + 1 // ixit FROM BASE LEVEL
	erppr = erxbl + 1 // PMGR PANIC REQUESTED BY NON-PMGR PROCESS
	ersci = erppr + 1 // SYSTEM CALL AT INTERRUPT LEVEL
	ernip = ersci + 1 // PMGR -- NOT AN IAC PMGR
	ernid = ernip + 1 // PMGR -- NOT AN IAC-DRIVEN DEVICE
	ersgo = ernid + 1 // signl ALREADY OUTSTANDING
	erufr = ersgo + 1 // UNKNOWN REQUEST FUNCTION
	erifs = erufr + 1 // ILLEGAL FED STRING
	era1o = erifs + 1 // ATTEMPT TO 1ST OPEN AN OPEN FILE
	erifi = era1o + 1 // INVALID PROTECTED FILE ID
	erapu = erifi + 1 // ATTEMPT TO PASS UNHELD ACCESS PRIVILEGES
	ernbk = erapu + 1 // NO BREAKFILE ENABLED FOR THIS RING
	ercds = ernbk + 1 // PMGR: MODEM DISCONNECT IN PROGRESS - CAN'T OPEN
	ertnf = ercds + 1 // TASK IS NOT FAULTING
	ernmt = ertnf + 1 // MAP TARGET DOES NOT EXIST
	ermte = ernmt + 1 // MAP TARGET (ALREADY) MAPPED ELSEWHERE
	ermsi = ermte + 1 // MAP SPECIFICATION ILLEGAL FOR TARGET
	errau = ermsi + 1 // MAP REGION ALREADY IN USE
	erjai = errau + 1 // JP ALREADY INITIALIZED
	erjni = erjai + 1 // JP NOT INITIALIZED
	erlne = erjni + 1 // LP DOES NOT EXIST
	erlai = erlne + 1 // LP ALREADY EXISTS
	erljp = erlai + 1 // ATTEMPT TO RELEASE LAST JP ATTACHED TO AN LP
	erijp = erljp + 1 // INVALID JPID
	erilp = erijp + 1 // INVALID LPID
	erjst = erilp + 1 // JP RUNNING ONE OR MORE SYSTEM TASKS
	erjaa = erjst + 1 // JP ALREADY ATTACHED TO LP
	erjna = erjaa + 1 // JP NOT ATTACHED TO LP
	ermlp = erjna + 1 // ATTEMPT TO EXCEED MAXIMUM LP COUNT
	erjpa = ermlp + 1 // CANNOT DELETE LP WITH JP ATTACHED
	eriti = erjpa + 1 // INVALID TIME INTERVAL
	erici = eriti + 1 // INVALID CLASS ID
	ercpc = erici + 1 // INVALID CLASS PERCENTAGE
	erihl = ercpc + 1 // INVALID HIERARCHICAL LEVEL
	erclu = erihl + 1 // CLASS IN USE
	erimp = erclu + 1 // ILLEGAL BIT MAP
	erilv = erimp + 1 // ILLEGAL LOCALITY VALUE
	erlp0 = erilv + 1 // CANNOT DELETE LP 0
	ernmp = erlp0 + 1 // NOT A MULTI-PROCESSOR SYSTEM
	ercne = ernmp + 1 // CLASS DOES NOT EXIST
	erhlp = ercne + 1 // ILLEGAL HIERARCHY LEVEL / PERCENTAGE PAIR
	ericd = erhlp + 1 // ILLEGAL FUNCTION CODE
	erjps = ericd + 1 // JP IS IN A BAD STATE
	ercmm = erjps + 1 // MICROCODE IS INCOMPATIBLE WITH CURRENT SYSTEM
	ermcr = ercmm + 1 // INCORRECT MICROCODE REVISION
	ermff = ermcr + 1 // MICROCODE FILE FORMAT ERROR
	erucp = ermff + 1 // INVALID CPU MODEL NUMBER
	ercso = erucp + 1 // CLASS SCHEDULING IS ENABLED
	erhlt = ercso + 1 // NON-SEQUENTIAL HIERARCHY LEVELS DESIGNATED
	erpor = erhlt + 1 // PID IS OUT OF RANGE FOR THIS PROCESS
	erpno = erpor + 1 // PROCESS NOT AN OPERATOR
	erbce = erpno + 1 // MAX BLOCK COUNT EXCEEDED
	erdeb = erbce + 1 // DAEMON ERROR IN ERROR BUFFER
	erdrf = erdeb + 1 // DAEMON RESOURCE FAILURE
	erlas = erdrf + 1 // LOG ALREADY STARTED
	ercl0 = erlas + 1 // CANNOT DELETE CLASS 0
	erpvm = ercl0 + 1 // UNKNOWN PRIVILEGE MODE
	erpvx = erpvm + 1 // PRIVILEGE HELD EXCLUSIVELY BY OTHER PROCESS
	erpvo = erpvx + 1 // PRIVILEGE CANNOT BE HELD EXCLUSIVELY
	erpvp = erpvo + 1 // OTHER PROCESSES USING PRIVILEGE
	erwcp = erpvp + 1 // WORKING SET CHANGE ONLY PARTLY DONE
	erfrd = erwcp + 1 // FAULT RECURSION DEPTH EXCEEDED
	eralp = erfrd + 1 // PROCESS'S CLASS NOT SCHEDULABLE ON AN ACTIVE LP
	ercll = eralp + 1 // INVALID CELL COUNT
	ernml = ercll + 1 // NO MICROCODE LOADED IN THIS JP
	eregn = ernml + 1 // END OF GET NEXT SEQUENCE
	erctd = eregn + 1 // CORRUPTED TASK CONTROL BLOCK DATA DETECTED

	eriwr = erctd + 1 // INVALID WINDOW REFERENCE.
	erwnn = eriwr + 1 // MAXIMUM NUMBER OF WINDOWS EXCEEDED.
	erwmd = erwnn + 1 // WINDOW MARKED FOR DELETION.
	erigp = erwmd + 1 // INVALID GRAPHICS PARAMETER.
	eripp = erigp + 1 // INVALID POINTER DEVICE PARAMETER.
	erivs = eripp + 1 // INVALID VIEW OR SCAN PORT.
	eriwo = erivs + 1 // INVALID WINDOWING OPERATION.
	eriwp = eriwo + 1 // INVALID WINDOWING PARAMETER.
	erade = eriwp + 1 // ASSOCIATION DOES NOT EXIST.
	eruwe = erade + 1 // UNKNOWN WINDOWING SUBSYSTEM ERROR
	ernsp = eruwe + 1 // HARDWARE/MICROCODE DOES NOT SUPPORT PIXEL MAPS
	erifl = ernsp + 1 // IAC FAILURE
	ertmo = erifl + 1 // TOO MANY OPENS ON THIS DEVICE.

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
	// FOR RUNNING THIS PROGRAM IF 0 ( DEFAULT )
	// DEFAULT SWAPFILE SIZE IS USED.

	//       USER STATUS TABLE (UST) TEMPLATE
	ust   = 0400      // START OF USER STATUS AREA 256.
	ustez = 0         // EXTENDED VARIABLE  WORD COUNT 256 = 0
	ustes = ustez + 1 // EXTENDED VARIABLE PAGE 0 START 257= 1
	ustss = ustes + 1 // SYMBOLS START 258                 = 2
	ustse = ustss + 2 // SYMBOLS END 260                   = 4
	ustda = ustse + 2 // DEB ADDR OR -1 262                = 6
	ustrv = ustda + 2 // REVISION OF PROGRAM 264         = 010
	usttc = ustrv + 2 // NUMBER OF TASKS (1 TO 32.) 266  = 012
	ustbl = usttc + 1 // # IMPURE BLKS 267               = 013
	ustst = ustbl + 3 // SHARED STARTING BLK # 270       = 016
	// USTST IS USTBL+3 BECAUSE THE 16BIT USER'S
	// USTOD IS HIDDEN UNDERNEATH
	ustit = ustst + 2  // INTERRUPT ADDRESS 272               = 020
	ustsz = ustit + 2  // SHARED SIZE IN BLKS 274             = 022
	ustpr = ustsz + 2  // PROGRAM FILE TYPE (16 OR 32 BIT) 276= 024
	ustsh = ustpr + 5  // PHYSICAL STARTING PAGE OF SHARED AREA IN .PR 281= 031
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
	dslth  = dnum + 1   // LENGTH OF SHORT PACKET

	dsh   = dnum + 1 // STARTING HOUR, -1 IF IMMEDIATE
	dsms  = dsh + 1  // STARTING SECOND IN HOUR, IGNORED IF IMMEDIATE
	dcc   = dsms + 1 // NUMBER OF TIMES TO CREATE TASK(S)
	dci   = dcc + 1  // CREATION INCREMENT  IN SECONDS
	dxlth = dci + 1  // LENGTH OF EXTENDED PACKET

	// BIT POINTER TO TASK DEF BITS
	dfbrc = dflgs*16 + 1 // RESOURCE CALL

	//  GENERAL USER I/O PACKET
	//
	//        USED FOR open/read/write/close
	//
	ich  dg.PhysAddrT = 0             // CHANNEL NUMBER
	isti              = ich + 1       // STATUS WORD (IN)
	isto              = isti + 1      // RIGHT=FILE TYPE, LEFT=RESERVED
	imrs              = isto + 1      // PHYSICAL RECORD SIZE - 1 (BYTES)
	ibad              = imrs + 1      // BYTE POINTER TO BUFFER
	ibal              = ibad + 1      // LOW ORDER BITS OF ibad
	ires              = ibal + 1      // RESERVED
	ircl              = ires + 1      // RECORD LENGTH
	irlr              = ircl + 1      // RECORD LENGTH (RETURNED)
	irnw              = irlr + 1      // RESERVED
	irnh              = irnw + 1      // RECORD NUMBER (HIGH)
	irnl              = irnh + 1      // RECORD NUMBER (LOW)
	ifnp              = irnl + 1      // BYTE POINTER TO FILE NAME
	ifnl              = ifnp + 1      // LOW ORDER BITS OF ifnp
	idel              = ifnl + 1      // DELIMITER TABLE ADDRESS
	idll              = idel + 1      // LOWER BITS OF idel
	iosz int          = int(idll) + 1 // LENGTH OF STANDARD I/O PACKET

	//  isti FLAGS: BIT DEFINITIONS
	iplb = 0  // PACKET LENGTH BIT (0 => SHORT PACKET)
	icfb = 1  // CHANGE FORMAT BIT (0 => DEFAULT)
	icdm = 1  // DUMP MODE BIT (ON close ONLY)
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

	//  isti FLAGS: MASK DEFINITIONS
	ipkl = 0x8000 >> iplb // EXTENDED PACKET (IF SET)
	icrf = 0x8000 >> icfb // CHANGE RECORD FORMAT (IF SET)
	cdmp = 0x8000 >> icdm // SET DUMP BIT (ONLY ON close)
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

// FLAGS FOR RETURN TO CLI (return)
const (
	Rfcf = 1 << 7 // 1B0             // CLI FORMAT
	Rfwa = 1 << 5 // 1B2             // WARNING (SEVERITY=1)
	Rfer = 2 << 5 // 2B2             // ERROR   (SEVERITY=2)
	Rfab = 3 << 5 // 3B2             // ABORT   (SEVERITY=3)
	Rfec = 1 << 4 // 1B3             // ERROR CODE FLAGIF SET, AC0 CONTAINS ERROR CODE
)

const (
	// PACKET TO GET INITIAL MESSAGE (gtmes)
	//
	greq dg.PhysAddrT = 0        // REQUEST TYPE (SEE BELOW)
	gnum              = greq + 1 // ARGUMENT NUMBER
	gsw               = gnum + 1 // BYTE PTR TO POSSIBLE SWITCH
	gsw1              = gsw + 1  // LOWER PORTION OF gsw
	gres              = gsw1 + 1 // BYTE PTR TO AREA TO RECEIVE
	grel              = gres + 1 // LOWER PORTION OF gres
	// SWITCH
	gtln = grel + 1 // PACKET LENGTH

	// REQUEST TYPES (greq)
	gmes dg.WordT = 0        // GET ENTIRE MESSAGE
	gcmd          = gmes + 1 // GET CLI COMMAND
	gcnt          = gcmd + 1 // GET ARGUMENT COUNT
	garg          = gcnt + 1 // GET ARGUMENT
	gtsw          = garg + 1 // TEST SWITCH
	gsws          = gtsw + 1 // TEST SWITCHES
	gdlc          = 1 << 15  //1B0             // DISABLE LOWER TO UPPERCASE CONVERSION

	// FLAGS RETURNED ON gflg TYPE CALLS
	gfcf = 1 << 15 // 1B0             // CLI FORMAT

	// BY CONVENTION, PROGRAMS CALLABLE FROM EXEC USE BITS 1 & 2
	// IF gfcf IS 0.
	gfex = 1 << 14 //1B1             // FROM EXEC IF ON

	//IF gfex IS ON, gfxb GIVES JOB'S BATCH/INTERACTIVE STATUS
	gfxb = 1 << 13 //1B2             // ON=BATCH, OFF=INTERACTIVE
	// IN ADDITION, IF CLI IS INVOKED WITH gfcf 0, BOTH gfxb & gfex
	// EQUAL TO ZERO => EXECUTE COMMAND PASSED IN MESSAGE AND RETURN.

)

const (
	// PACKET TO GET SYSTEM INFORMATION (sinfo)
	sirn = 0        // SYSTEM REV, LEFT BYTE=MAJOR,RIGHT BYTE=MINOR
	sirs = sirn + 1 // RESERVED
	simm = sirs + 1 // LENGTH OF PHYSICAL MEMORY (HPAGE)
	siml = simm + 1 // LOWER PORTION OF simm
	siln = siml + 1 // BYTE POINTER TO RECEIVE MASTER LDU NAME
	sill = siln + 1 // LOWER PORTION OF siln
	siid = sill + 1 // BYTE POINTER TO RECEIVE SYSTEM IDENTIFIER
	siil = siid + 1 // LOWER PORTION OF siid
	sipl = siil + 1 // UNEXTENDED PACKET LENGTH
	sios = siil + 1 // BYTE POINTER TO EXECUTING OP SYS PATHNAME
	siol = sios + 1 // LOWER PORTION OF sios
	ssin = siol + 1 // SYSTEM IMPLEMENTATION NUMBER (savs FOR AOSVS)

	siex = ssin + 6 // EXTENDED PACKET LENGTH (INCLUDE 3 DOUBLE
	//        WORDS FOR FUTURE EXPANSIONS)

)

const (
	// USER PACKET DEFINITION FOR     xpstat     SYSTEM CALL

	// RETURNS INTERESTING ITEMS FROM PROCESS TABLE, (AND MORE)
	// THIS CALL SHOULD BE USED IN PLACE OF THE OLD pstat

	xpsid  = (0576 * 0200000) + (0 * 0400) + 0 // xpstat PACKET ID(IE.(xpsp)
	xpsid1 = (0576 * 0200000) + (0 * 0400) + 1 // xpstat NEW PACKET ID

	xpsp  = 0        // SUB PACKET IDENTIFIER DWORD
	xpsf  = xpsp + 2 // SUB FUNCTION ID
	xpfp  = xpsf + 1 // PROCESS ID OF TARGET PID'S FATHER
	xpnr  = xpfp + 1 // # OF TASKS SUSPENDED ON irec
	xpsns = xpnr + 1 // # OF TASKS BLOCKED AWAITING SYSTEM STACKS

	// SEE pstat PACKET FOR THE DEFINITIONS OF THE BITS IN THE FOLLOWING WORDS

	xpsw  = xpsns + 1 // PROCESS STATUS WORD
	xpsqf = xpsw + 1  // PRIORITY QUEUE FACTOR
	xpfl  = xpsqf + 1 // FIRST  PROCESS FLAG WORD
	xpf2  = xpfl + 1  // SECOND PROCESS FLAG WORD
	xpf3  = xpf2 + 1  // THIRD  PROCESS FLAG WORD
	xpf4  = xpf3 + 1  // FOURTH PROCESS FLAG WORD
	xpf5  = xpf4 + 1  // FIFTH  PROCESS FLAG WORD

	xppr  = xpf5 + 1  // PROCESS PRIORITY
	xpcw  = xppr + 1  // CURRENT WORKING SET SIZE IN PAGES (DWORD)
	xpr1  = xpcw + 2  // RESERVED FOR FUTURE USE
	xppv  = xpr1 + 1  // PROCESS PRIVILEGE BITS (SEE pstat DEFINITIONS)
	xpex  = xppv + 1  // TIME SLICE EXPONENT
	xppd  = xpex + 1  // PID OR VPID OF THE TARGET PROCESS
	xprh  = xppd + 1  // # SECONDS ELAPSED SINCE PROCESS CREATION
	xpch  = xprh + 2  // MILLISECONDS OF CPU TIME
	xpcpl = xpch + 2  // MAXIMUM CPU TIME ALLOWED
	xpph  = xpcpl + 2 // PAGE USAGE OVER CPU TIME (PAGES/SEC)
	xpmx  = xpph + 2  // MAXIMUM LOGICAL PAGES FOR RING 7
	xpws  = xpmx + 2  // MAXIMUM WORKING SET SIZE
	xpwm  = xpws + 2  // MINIMUM WORKING SET SIZE
	xpfa  = xpwm + 2  // NUMBER OF PAGE FAULTS SINCE PROCESS CREATION

	xpdis = xpfa + 2  // WORD POINTER TO MEMORY DESCRIPTORS ARRAY BUFFER
	xpdbs = xpdis + 2 // AVAILABLE BUFFER SIZE FOR MEMORY DESCRIPTORS
	//  (MUST BE AT LEAST 7*pdesln)
	xpdrs = xpdbs + 1 // ACTUAL NUMBER OF WORDS RETURNED IN BUFFER

	xpih  = xpdrs + 1 // # OF BLOCKS READ/WRITTEN
	xplfa = xpih + 2  // # OF PAGE FAULTS NOT REQUIRING I/O
	xpsl  = xplfa + 2 // # OF SUB-SLICES LEFT
	xpcpu = xpsl + 1  // CURRENT CPU NUMBER FOR THIS PROCESS

	xpll  = xpcpu + 1 // LEGAL LOCALITIES
	xpulc = xpll + 1  // CURRENT USER LOCALITY
	xppl  = xpulc + 1 // PROGRAM LOCALITY
	xpcid = xppl + 1  // CLASS ID OF THE PROCESS

	xppg  = xpcid + 1 // BYTE POINTER TO PROCESS GROUP NAME
	xpgbs = xppg + 2  // AVAILABLE BUFFER SIZE FOR PROCESS GROUP NAME
	xpgrs = xpgbs + 1 // RETURNED PROCESS GROUP NAME SIZE IN BYTES

	xpun  = xpgrs + 1 // BYTE POINTER TO USER NAME
	xpnbs = xpun + 2  // AVAILABLE BUFFER SIZE FOR USER NAME
	xpnrs = xpnbs + 1 // RETURNED USER NAME SIZE IN BYTES

	xppu  = xpnrs + 1 // BYTE POINTER TO PROC'D USER NAME
	xpubs = xppu + 2  // AVAILABLE BUFFER SIZE FOR PROC'D USER NAME
	xpurs = xpubs + 1 // RETURNED PROC'D USER NAME SIZE IN BYTES

	xpupd = xpurs + 1 // 128 BIT UNIQUE PROCESS ID (UPID)
	// (128BITS = EIGHT WORDS)

	xpglt  = xpurs + 1  // RESERVED FOR GROUP ACLS
	xpgab  = xpglt + 2  // RESERVED FOR GROUP ACLS
	xpgrb  = xpgab + 1  // RESERVED FOR GORUP ACLS
	xpuqsh = xpgrb + 1  // UNIQUE SHARED PAGES
	xprs1  = xpuqsh + 2 // RESERVED FOR FUTURE USE

	xplth = xprs1 + 2 // PACKET LENGTH

//               END OF   xpstat    PACKET DEFINITION

)

const (
	//  PERIPHERAL DEVICE CHARACTERISTICS

	//        The following parameters are for the characteristic packet offsets

	ch1  = 0  // word 1 (offset 0)
	ch2  = 1  // word 2 (offset 1)
	ch3  = 2  // word 3 (offset 2)
	ch4  = 3  // word 4 (offset 3)
	ch5  = 4  // word 5 (offset 4)
	ch6  = 5  // word 6 (offset 5)
	ch7  = 6  // word 7 (offset 6)
	ch8  = 7  // word 8 (offset 7)
	ch9  = 8  // word 9 (offset 8)
	ch10 = 9  // word 10 (offset 9)
	ch11 = 10 // word 11 (offset 10)
	ch12 = 11 // word 12 (offset 11)
	ch13 = 12 // word 13 (offset 12)
	ch14 = 13 // word 14 (offset 13)
	ch15 = 14 // word 15 (offset 14)

	//        Packet length parameters

	clmin = 3  //  MIN LENGTH OF CHARACTERISTICS PACKET
	clmax = 15 //  MAX LENGTH OF CHARACTERISTICS PACKET
	bmlth = 20 //  LENGTH OF INQUIRE PACKET

	//        ch1 - offset 0

	cst  = 0  // SIMULATE TABS
	csff = 1  // SIMULATE FORM FEEDS
	cepi = 2  // REQUIRE EVEN PARITY ON INPUT
	c8bt = 3  // ALLOW 8 DATA BITS/CHARACTER
	cspo = 4  // SET PARITY ON OUTPUT (EVEN ONLY)
	craf = 5  // SEND RUBOUTS AFTER FORM FEEDS
	crat = 6  // SEND RUBOUTS AFTER TABS
	crac = 7  // SEND RUBOUTS AFTER CR AND NL
	cnas = 8  // NON ANSI STANDARD DEVICE
	cott = 9  // CONVERT ESC CHARACTER (FOR OLD TTY'S)
	ceol = 10 // DO NOT AUTO CR/LF AT END OF LINE
	cuco = 11 // OUTPUT UPPER CASE ONLY DEVICE
	cmri = 12 // MONITOR RING INDICATOR ON MODEM CONTROL LINE
	cff  = 13 // FORM FEED ON OPEN
	//        THE FOLLOWING TWO BITS MUST NOT BE MOVED :
	ceb0 = 14 // ECHO MODE BIT 0
	ceb1 = 15 // ECHO MODE BIT 1

	//        ECHO MODES :
	//        0=      NO ECHO
	//        1=      STRAIGHT ECHO
	//        2=      ECHO CONTROL CHARS AS ^B ^F (ETC.), ESC AS $
	//        3=      (RESERVED FOR FUTURE USE)

	ceos = 0x8000 >> ceb1 // 0x8000 >> ceb1       // STRAIGHT ECHO BIT MASK
	ceoc = 0x8000 >> ceb0 // 0x8000 >> ceb0       // CNTRL SPECIAL ECHO BIT MASK

	//        ch2 - offset 1

	culc = 0 // INPUT UPPER/LOWER CASE DEVICE
	cpm  = 1 // DEVICE IS IN PAGE MODE
	cnrm = 2 // DISABLE MESSAGE RECEPTION
	cmod = 3 // DEVICE ON MODEM INTERFACE

	//        THE FOLLOWING FOUR BITS MUST NOT BE MOVED :
	cdt0 = 4 // DEVICE TYPE BIT 0
	cdt1 = 5 // DEVICE TYPE BIT 1
	cdt2 = 6 // DEVICE TYPE BIT 2
	cdt3 = 7 // DEVICE TYPE BIT 3

	cto  = 8  // DEVICE TIME-OUTS ENABLED
	ctsp = 9  // CRA- NO TRAILING BLANK SUPPRESSION
	cpbn = 10 // CRA- PACKED FORMATE ON BINARY READ
	cesc = 11 // ESC CHARACTER PRODUCES INTERRUPT
	cwrp = 12 // HARDWARE WRAPS AROUND ON LINE TOO LONG
	cfkt = 13 // FUNCTION KEYS ARE INPUT DELIMITERS
	cnnl = 14 // CRA- NO NEW-LINE CHARACTERS APPENDED
	//                15    // BIT 15 USED IN PARU.16.SR FOR TRA/TPA

	//        DEFINE DEVICE TYPE MASK.

	dtype = 0x8000>>cdt0 + 0x8000>>cdt1 + 0x8000>>cdt2 + 0x8000>>cdt3

	tty   = 0                                          // 4010A CONSOLE DEVICE TYPE
	crt1  = 0x8000 >> cdt3                             // 4010I CONSOLE DEVICE TYPE
	crt2  = 0x8000 >> cdt2                             // 6012  CONSOLE DEVICE TYPE
	crt3  = 0x8000>>cdt2 + 0x8000>>cdt3                // 605X CONSOLE DEVICE TYPE
	crt4  = 0x8000 >> cdt1                             // ANOTHER CONSOLE DEVICE TYPE
	crt5  = 0x8000>>cdt1 + 0x8000>>cdt3                // PSEUDO 6012 DEVICE
	crt6  = 0x8000>>cdt1 + 0x8000>>cdt2                // 6130 CONSOLE DEVICE TYPE
	crt7  = 0x8000>>cdt1 + 0x8000>>cdt2 + 0x8000>>cdt3 // USER DEFINED DEVICE
	crt8  = 0x8000 >> cdt0                             // USER DEFINED DEVICE
	crt9  = 0x8000>>cdt0 + 0x8000>>cdt3                // USER DEFINED DEVICE
	crt10 = 0x8000>>cdt0 + 0x8000>>cdt2                // USER DEFINED DEVICE
	crt11 = 0x8000>>cdt0 + 0x8000>>cdt2 + 0x8000>>cdt3 // USER DEFINED DEVICE
	crt12 = 0x8000>>cdt0 + 0x8000>>cdt1                // USER DEFINED DEVICE
	crt13 = 0x8000>>cdt0 + 0x8000>>cdt1 + 0x8000>>cdt3 // USER DEFINED DEVICE
	crt14 = 0x8000>>cdt0 + 0x8000>>cdt1 + 0x8000>>cdt2 // USER DEFINED DEVICE
	crt15 = 0x8000>>cdt0 + 0x8000>>cdt1 + 0x8000>>cdt2 + 0x8000>>cdt3

	//        ch3 - offset 2
	//
	//        HIGH BYTE IS LPP (LINES PER PAGE)
	//        LOW  BYTE IS CPL (CHARACTERS PER LINE)

	cpgsz = ch3 // Page size

	//        ch4 - offset 3

	cval = 0 // INDICATES THAT THE CONTENTS OF THIS
	// OFFSET ARE VALID(USED ON RETURN
	// FROM gechr.)  IN GENERAL, cval= 1
	// FOR AN IAC SYSTEM, AND cval OTHERWISE.
	br0bit = 1      // BAUD RATE FIELD (BIT 0)
	ctck   = 2      // INTERNAL TRANSMITER CLOCK
	crck   = 3      // INTERNAL RECIEVER CLOCK
	br1bit = 4      // BAUD RATE FIELD (BIT 1)
	br2bit = 5      // BAUD RATE FIELD (BIT 2)
	br3bit = 6      // BAUD RATE FIELD (BIT 3)
	br4bit = 7      // BAUD RATE FIELD (BIT 4)
	cst0   = 8      // STOP BIT 0
	cst1   = 9      // STOP BIT 1
	cpty   = 10     // ODD/EVEN PARITY
	cpen   = 11     // PARITY DISABLED/ENABLED
	clt0   = 12     // DATA LENGTH BITS
	clt1   = 13     // DATA LENGTH BITS
	brfct  = 14     // BAUD RATE FACTOR 16X
	hrdflc = 15     // HARDWARE FLOW CONTROL (CTS)
	chofc  = hrdflc // HARDWARE OUTPUT FLOW CONTROL

	//        SPLIT BAUD RATE VALUES:

	csben = 0x8000>>ctck + 0x8000>>brfct                // ENABLE SPLIT BAUD
	csbds = 0x8000>>ctck + 0x8000>>crck + 0x8000>>brfct // DISABLE SPLIT BAUD

	//        STOP BIT FIELD VALUES ARE:

	csmsk = 0x8000>>cst0 + 0x8000>>cst1 // STOP BIT FIELD MASK

	// cs10=  0bcst0+0x8000 >> cst1         // 1 STOP BIT
	// cs15=  0x8000 >> cst0+0bcst1         // 1.5 STOP BITS
	// cs20=  0x8000 >> cst0+0x8000 >> cst1         // 2 STOP BITS

	//        PARITY BIT FIELD VALUES ARE:

	// cpmsk= 0x8000 >> cpen+0x8000 >> cpty         // PARITY FIELD MASK

	// cpr0=  0bcpen                 // DISABLE PARITY CHECKING
	// cpr1=  1bcpen+0bcpty         // ENABLE ODD  PARITY
	// cpr2=  1bcpen+1bcpty         // ENABLE EVEN PARITY

	// //        BAUD RATES ARE:

	// brmsk= 0x8000 >> br0bt)!17B(br4bit)        // BAUD RATE MASK

	// cr50=  0B(br0bit)+0.B(br4bit)        // 50
	// cr75=  0B(br0bit)+1.B(br4bit)        // 75
	// cr110= 0B(br0bit)+2.B(br4bit)        // 110
	// cr134= 0B(br0bit)+3.B(br4bit)        // 134.5
	// cr150= 0B(br0bit)+4.B(br4bit)        // 150
	// cr300= 0B(br0bit)+5.B(br4bit)        // 300
	// cr600= 0B(br0bit)+6.B(br4bit)        // 600
	// cr12h= 0B(br0bit)+7.B(br4bit)        // 1200
	// cr18h= 0B(br0bit)+8.B(br4bit)        // 1800
	// cr20h= 0B(br0bit)+9.B(br4bit)        // 2000
	// cr24h= 0B(br0bit)+10.B(br4bit)       // 2400
	// cr36h= 0B(br0bit)+11.B(br4bit)       // 3600
	// cr48h= 0B(br0bit)+12.B(br4bit)       // 4800
	// cr72h= 0B(br0bit)+13.B(br4bit)       // 7200
	// cr96h= 0B(br0bit)+14.B(br4bit)       // 9600
	// cr19k= 0B(br0bit)+15.B(br4bit)       // 19200

	// cr45=  0x8000 >> br0bt)+0.B(br4bit)        // 45.5
	// cr38k= 0x8000 >> br0bt)+1.B(br4bit)        // 38400
	//                            2- 15           //  - RESERVED

	// //        DATA LENGTH FIELD VALUES ARE:

	// clmsk= 1bclt0+1bclt1         // DATA LENGTH FIELD MASK

	// cln5=  0bclt0+0bclt1         // 5 BITS
	// cln6=  0bclt0+1bclt1         // 6 BITS
	// cln7=  1bclt0+0bclt1         // 7 BITS
	// cln8=  1bclt0+1bclt1         // 8 BITS

	//        ch5 - offset 4

	shco    = 0  // SHARED CONSOLE OWNERSHIP CHARACTERISTIC
	xofc    = 1  // XON XOFF OUTPUT FLOW CONTROL
	xifc    = 2  // XON XOFF INPUT  FLOW CONTROL
	c16b    = 3  // Enable double byte handling (16 bit characters)
	bmdev   = 4  // BITMAP DEVICE
	trpe    = 5  // TERMINATE READ ON POINTER EVENT
	cwin    = 6  // WINDOW CHARACTERISTIC
	cacc    = 7  // ENFORCE ACCESS CONTROL
	cctd    = 8  // PORT IS IN A CONTENDED ENVIRONMENT (PBX, TERMSERVER)
	csrds   = 9  // SUPRESS RECEIVER DISABLE
	cxlt    = 10 // TRANSLATE (ANSI TERMINAL)
	cabd    = 11 // [1] DO AUTOBAUD MATCH IF SET
	callout = 12 // CALL OUT (PBX SUPPORT)
	cbk0    = 13 // BREAK FUNCTION BIT 0
	cbk1    = 14 // BREAK FUNCTION BIT 1
	cbk2    = 15 // BREAK FUNCTION BIT 2

	// // BREAK FUNCTION FIELD DEFINITION:

	// cbkm=  1bcbk0+1bcbk1+1bcbk2 // MASK

	// cbbm=  0B(cbk2)               // BREAK BINARY MODE
	// cbds=  0x8000 >> cbk2               // FORCE DISCONNECT
	// cbca=  2B(cbk2)               // SEND ^C^A SEQUENCE
	// cbcb=  3B(cbk2)               // SEND ^C^B SEQUENCE
	// cbcf=  4B(cbk2)               // SEND ^C^F SEQUENCE
	//                5B(cbk2)               //  - RESERVED
	//                6B(cbk2)               //  - RESERVED
	//                7B(cbk2)               //  - RESERVED

	//        ch6 - offset 5
	//        (MODEM ENHANCEMENTS)

	cmdop = ch6 // Modem options

	cdmc  = 0        // RESERVED
	cmdua = cdmc + 1 // DIRECT USER ACCESS TO MODEM
	// (DON'T PEND FIRST WRITE)
	chdpx = cmdua + 1 // HALF DUPLEX
	csmcd = chdpx + 1 // SUPPRESS MONITORING CD
	// (FOR MODEM CONNECTION)
	crtscd = csmcd + 1 // ON HALF DUPLEX, DON'T RAISE
	// RTS UNTIL CD DROPS
	chifc = crtscd + 1 // HARDWARE INPUT FLOW CONTROL

	//        ch7 - offset 6
	ctcc = ch7 // Time (in msec) to wait for CD on a modem
	// connect

	//        ch8 - offset 7
	ctcd = ch8 // Time (in msec) to wait for CD if it drops

	//        ch9 - offset 8
	ctdw = ch9 // Time (in msec) to wait after connection
	// before allowing I/O

	//        ch10 - offset 9
	cthc = ch10 // Time (in msec) to wait after disconnect
	// for modem to settle

	//        ch11 - offset 10
	ctlt = ch11 // Time (in msec) to wait before turning
	// the line around (from XMIT to REC) for
	// half duplex

	//        ch12 - offset 11
	//        (Console Type)
	//
	//        HIGH BYTE IS RESERVED (=0)
	//        LOW  BYTE IS CONSOLE TYPE

	cctype = ch12 // Console type

	//        Mask for accessing just console type

	cctypmsk = 377 // mask for just console type

	//        These are the current values for console types

	cdcc = 0        // Direct Connect
	clnc = cdcc + 1 // Term Server
	ctnc = clnc + 1 // TELNET Consoles
	cpdc = ctnc + 1 // PAD Consoles
	cvrc = cpdc + 1 // Virtual (SVTA-like) Consoles
	cpxc = cvrc + 1 // PBX Consoles (PIM)
	cpcc = cpxc + 1 // PC/TS Consoles
	cbmc = cpcc + 1 // Bitmapped (Windowing) Console
	ctpc = cbmc + 1 // T1 Primary Rate Console(IIC)

	//        ch13 - offset 12
	//        (Language Front-end Processor)

	clfp = ch13 // LFP options

	ckg0 = 0 // G1-G0 double-byte handling
	ckhw = 1 // Kanji half-wide characters
	cnlx = 2 // Native language translation

//        DEVICE TYPES : (FOR RUBOUT ECHO & CURSOR CONTROLS)
//
//        PIBC2   CHARACTERS TO :
//        DEVICE  MODEL   MOVE    MOVE    ERASE   RUBOUT
//        TYPE :  # :     LEFT:   RIGHT:  LINE:   ECHO:
//
//        0       4010A   (NONE)  (NONE)  (NONE)  SHIFT O
//        0       6040    (NONE)  (NONE)  (NONE)  SHIFT O
//        1       4010I   ^Z      ^Y      ^K      ^Z,SPACE,^Z
//        2       6012    ^Y      ^X      ^K      ^Y,SPACE,^Y
//        3       6052    ^Y      ^X      ^K      ^Y,SPACE,^Y
//        4       ----    ESC,D   ESC,C   ESC,K   ESC,D,SPACE,ESC,D
//        5       ----
//        6       6130    ^Y      ^X      ^K      ^Z,SPACE,^Z
//        7-15  (FOR FUTURE EXPANSION)
)

const (
	// ACCESS CONTROL SPECIFICATIONS

	// ACCESS PRIVILEGE BIT POSITIONS
	faob = 11 // OWNER ACCESS
	fawb = 12 // WRITE ACCESS
	faab = 13 // APPEND ACCESS
	farb = 14 // READ ACCESS
	faeb = 15 // EXECUTE ACCESS

	// DEFINITIONS FOR USE WITH WSKBN/WSKBZ EAGLE INSTRUCTIONS
	faobb = 16 + faob // OWNER ACCESS
	fawbb = 16 + fawb // WRITE ACCESS
	faabb = 16 + faab // APPEND ACCESS
	farbb = 16 + farb // READ ACCESS
	faebb = 16 + faeb // EXECUTE ACCESS

	// ACCESS PRIVILEGE MASKS
	faco = 0x8000 >> faob // OWNER
	facw = 0x8000 >> fawb // WRITE
	faca = 0x8000 >> faab // APPEND
	facr = 0x8000 >> farb // READ
	face = 0x8000 >> faeb // EXECUTE
)

// SYSTEM RECORD I/O PACKET FOR ALL DISK AND MAG. TAPEAND MCA REQUESTS FROM EITHER THE AGENT OR USER CONTEXTS. USED FOR rdb/wrb, prdb/PWRB, spage AND allocate
const (
	psti = 0        // RECORD COUNT (RIGHT), STATUS IN (LEFT)
	psto = psti + 1 // RESERVED (LEFT) PRIORITY (RIGHT)
	pcad = psto + 1 // WORD ADDRESS FOR DATA
	pcdl = pcad + 1 // LOW ORDER PORTION OF pcad
	prnh = pcdl + 1 // RECORD NUMBER (HIGH) LINK # (MCA)
	prnl = prnh + 1 // RECORD NUMBER (LOW)  RETRY COUNT (MCA)
	prcl = prnl + 1 // MAX LENGTH OF EACH RECORD (MAG TAPE)
	//                 BYTE COUNT IN LAST BLOCK (DISK WRITES)
	//                 BYTE COUNT (MCA)
	pres = prcl + 1 // RESERVED WORD
	pblt = pres + 1 // PACKET SIZE
)

// SYSTEM ENTRY TYPES
const (
	// MISC
	flnk = smin     // LINK
	fsdf = flnk + 1 // SYSTEM DATA FILE
	fmtf = fsdf + 1 // MAG TAPE FILE
	fgfn = fmtf + 1 // GENERIC FILE NAME

	// DIRECTORIES (DO NOT CHANGE THEIR ORDER)
	fdir = 10.      // DISK DIRECTORY
	fldu = fdir + 1 // LD ROOT DIRECTORY
	fcpd = fldu + 1 // CONTROL POINT DIRECTORY
	fmtv = fcpd + 1 // MAG TAPE VOLUME
	fmdr = fmtv + 1 // RESERVED FOR RT32(MEM DIRS), NOT LEGAL FOR AOS
	fgnr = fmdr + 1 // RESERVED FOR RT32, NOT LEGAL FOR AOS
	ldir = fdir     // LOW DIR TYPE
	hdir = fgnr     // HIGH DIR TYPE
	lcpd = fldu     // LOW CONTROL POINT DIR TYPE
	hcpd = fcpd     // HIGH CONTROL POINT DIR TYPE

	// UNITS
	fdku = 20.      // DISK UNIT
	fmcu = fdku + 1 // MULTIPROCESSOR COMMUNICATIONS UNIT
	fmtu = fmcu + 1 // MAG TAPE UNIT
	flpu = fmtu + 1 // DATA CHANNEL LINE PRINTER
	flpd = flpu + 1 // DATA CHANNEL LP2 UNIT
	flpe = flpd + 1 // DATA CHANNEL LINE PRINTER (LASER)
	fpgn = flpe + 1 // RESERVED FOR RT32(PROCESS GROUP)
	fltu = flpu + 1 // LABELLED MAG TAPE UNIT
	// ***** NO LONGER USED *****
	lunt = fdku // LOW UNIT TYPE
	hunt = fpgn // HIGH UNIT TYPE

	// IPC ENTRY TYPES
	fipc = 30. // IPC PORT ENTRY

	// DGC ENTRY TYPES
	fudf = dmin     // USER DATA FILE
	fprg = fudf + 1 // PROGRAM FILE
	fupf = fprg + 1 // USER PROFILE FILE
	fstf = fupf + 1 // SYMBOL TABLE FILE
	ftxt = fstf + 1 // TEXT FILE
	flog = ftxt + 1 // SYSTEM LOG FILE (ACCOUNTING FILE)
	fncc = flog + 1 // FORTRAN CARRIAGE CONTROL FILE
	flcc = fncc + 1 // FORTRAN CARRIAGE CONTROL FILE
	ffcc = flcc + 1 // FORTRAN CARRIAGE CONTROL FILE
	focc = ffcc + 1 // FORTRAN CARRIAGE CONTROL FILE
	fprv = focc + 1 // AOS/VS PROGRAM FILE
	fwrd = fprv + 1 // WORD PROCESSING
	fafi = fwrd + 1 // APL FILE
	faws = fafi + 1 // APL WORKSPACE FILE
	fbci = faws + 1 // BASIC CORE IMAGE FILE
	fdcf = fbci + 1 // DEVICE CONFIGURATION FILE (NETWORKING)
	flcf = fdcf + 1 // LINK CONFIGURATION FILE (NETWORKING)
	flug = flcf + 1 // LOGICAL UNIT GROUP FILE (SNA)
	frtl = flug + 1 // AOS/RT32 RESERVED FILE TYPE RANGE (LO)
	frth = frtl + 4 // AOS/RT32 RESERVED FILE TYPE RANGE (HI)
	funx = frth + 1 // VS/UNIX FILE
	fbbs = funx + 1 // BUSINESS BASIC SYMBOL FILE
	fvlf = fbbs + 1 // BUSINESS BASIC VOLUME LABEL FILE
	fdbf = fvlf + 1 // BUSINESS BASIC DATA BASE FILE

	// CEO FILE TYPES
	fgkm = fdbf + 1 // DG GRAPHICS KERNAL METAFILE
	fvdm = fgkm + 1 // VIRTUAL DEVICE METAFILE
	fnap = fvdm + 1 // NAPLPS STANDARD GRAPH FILE
	ftrv = fnap + 1 // TRENDVIEW COMMAND FILE
	fspd = ftrv + 1 // SPREADSHEET FILE
	fqry = fspd + 1 // PRESENT QUERY MACRO
	fdtb = fqry + 1 // PHD DATA TABLE
	ffmt = fdtb + 1 // PHD FORMAT FILE
	fwpt = ffmt + 1 // TEXT INTERCHANGE FORMAT
	fdif = fwpt + 1 // DATA INTERCHANGE FORMAT
	fvif = fdif + 1 // VOICE IMAGE FILE
	fimg = fvif + 1 // FACSIMILE IMAGE
	fprf = fimg + 1 // PRINT READY FILE

	// MORE DGC ENTRY TYPES
	fpip = fprf + 1 // PIPE FILE
	fttx = fpip + 1 // TELETEX FILE
	fdxf = fttx + 1 // RESERVED FOR DXA
	fdxr = fdxf + 1 // RESERVED FOR DXA
	fcwp = fdxr + 1 // CEO WORD PROCESSOR FILE
	fcwt = fcwp     // CEOwrite WORD PROCESSOR FILE
	frpt = fcwp + 1 // PHD REPORT FILE
)

//  PACKET FOR DIRECTORY ENTRY CREATION (?CREATE)
const (
	cftyp = 0        // ENTRY TYPE (RH) AND RECORD FORMAT (LH)
	cpor  = 1        // PORT NUMBER (IPC TYPES ONLY)
	chfs  = 1        // HASH FRAME SIZE (DIRECTORY TYPES ONLY)
	chid  = 1        // HOST ID (frem TYPE FILES ONLY )
	ccps  = 1        // FILE CONTROL PARAMETER (OTHERS)
	ctim  = 2        // POINTER TO TIME BLOCK
	ctil  = ctim + 1 // LOWER PORTION OF ctim
	cacp  = ctil + 1 // POINTER TO INITIAL ACL
	cacl  = cacp + 1 // LOWER PORTION OF cacp
	cmsh  = cacl + 1 // MAX SPACE ALLOCATED (fcpd)
	cmsl  = cmsh + 1 // MAX SPACE ALLOCATED (LOW)
	cdeh  = cacl + 1 // RESERVED
	cdel  = cdeh + 1 // FILE ELEMENT SIZE
	cmil  = cdel + 1 // MAXIMUM INDEX LEVEL DEPTH
	cmrs  = cmil + 1 // RESERVED
	clth  = cmrs + 1 // LENGTH OF THE PARAMETER BLOCK
)
