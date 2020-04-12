// paru16.go - Go version of parts of AOS/VS PARU.16.SR definitions file

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

// N.B. To avoid namespace collisions, all these names have "16" appended
//      No need to duplicate symbols here if they are the same as those in
//      PARU.32.SR - eg. many option flags are identical.

package aosvs

const (
	//  GENERAL USER I/O PACKET
	//
	//        USED FOR ?OPEN/?READ/?WRITE/?CLOSE
	//
	ich16  = 0          // CHANNEL NUMBER
	isti16 = ich16 + 1  // STATUS WORD (IN)
	isto16 = isti16 + 1 // RIGHT=FILE TYPE, LEFT=RESERVED
	ibad16 = isto16 + 1 // BYTE POINTER TO BUFFER
	ires16 = ibad16 + 1 // RESERVED
	iflg16 = ires16     // WORD OF FLAGS
	ircl16 = ires16 + 1 // RECORD LENGTH
	irlr16 = ircl16 + 1 // RECORD LENGTH (RETURNED)
	irnh16 = irlr16 + 1 // RECORD NUMBER (HIGH)
	irnl16 = irnh16 + 1 // RECORD NUMBER (LOW)
	ifnp16 = irnl16 + 1 // BYTE POINTER TO FILE NAME
	imrs16 = ifnp16 + 1 // PHYSICAL RECORD SIZE - 1 (BYTES)
	idel16 = imrs16 + 1 // DELIMITER TABLE ADDRESS

	iblt16 = idel16 + 1 // PACKET LENGTH
)

const (
	// PACKET TO GET INITIAL MESSAGE (?GTMES)
	//
	greq16 = 0          // REQUEST TYPE (SEE BELOW)
	gnum16 = greq16 + 1 // ARGUMENT NUMBER
	gsw16  = gnum16 + 1 // BYTE PTR TO POSSIBLE SWITCH
	gres16 = gsw16 + 1  // BYTE PTR TO AREA TO RECEIVE SWITCH
	gtln16 = gres16 + 1 // PACKET LENGTH
)

// This is just here as a sanity check; the locations DO correspond to those in paru32
const (
	// UST.16
	// USER STATUS TABLE (UST) TEMPLATE
	ust16   = 400         // START OF USER STATUS AREA
	ustez16 = 0           // EXTENDED VARIABLE  WORD COUNT
	ustes16 = ustez16 + 1 // EXTENDED VARIABLE PAGE 0 START
	ustss16 = ustes16 + 2 // SYMBOLS START
	// NOTE THAT THE 16. BIT USER IS BEING
	// POINTED TO THE LOWER HALF OF A 32. BIT
	// BUCKET FOR THINGS LIKE USTSS, ETC.  THAT'S
	// WHY USTSS = USTES + 2, EVEN THO USTES IS 16.
	// BITS LONG
	ustse16 = ustss16 + 2 // SYMBOLS END
	ustda16 = ustse16 + 2 // DEB ADDR OR -1
	//.DUSR  USTSL=  ustda+1 // SHARED LIBRARY LIST POINTER
	ustrv16 = ustda16 + 1  // REVISION OF PROGRAM
	usttc16 = ustrv16 + 2  // NUMBER OF TASKS (1 TO 32.)
	ustbl16 = usttc16 + 2  // # IMPURE BLKS
	ustod16 = ustbl16 + 1  // OVLY DIRECTORY ADDR
	ustst16 = ustod16 + 2  // SHARED STARTING BLK #
	ustit16 = ustst16 + 2  // INTERRUPT ADDRESS
	ustsz16 = ustit16 + 2  // SHARED SIZE IN BLKS
	ustpr16 = ustsz16 + 1  // PROGRAM TYPE (16 OR 32 BIT)
	ustsh16 = ustpr16 + 5  // PHYSICAL STARTING PAGE OF SHARED AREA IN .PR
	usten16 = ustpr16 + 21 // LAST ENTRY VISIBLE TO USER
)
