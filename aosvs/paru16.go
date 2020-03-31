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
