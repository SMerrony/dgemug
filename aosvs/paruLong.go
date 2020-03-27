// paruLog.go - Go version of parts of AOS/VS PARULONG.SR 32-bit definitions file

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

const (
	// =======================================================================
	//                        ?SYSPRV PACKET
	// =======================================================================
	sysprvPktPktID  = 0                  // Packet ID
	sysprvPktFunc   = sysprvPktPktID + 2 // Function Code
	sysprvPktFlags  = sysprvPktFunc + 1  // Flags word
	sysprvPktMode   = sysprvPktFlags + 1 // Mode Value
	sysprvPktSubpkt = sysprvPktMode + 1  // W(subpacket)
	// <not used>

	sysprvPktLen = sysprvPktSubpkt + 2 // Length of packet

	// FUNCTION VALUES (sysprvPktFunc)
	sysprvGet       = 1 // Get privilege status
	sysprvEnter     = 2 // Enter a privileged mode
	sysprvEnterExcl = 3 // Enter a priv mode exclusively
	sysprvLeave     = 4 // Leave a privileged mode

	sysprvFuncMin = sysprvGet   // Minimum legal function
	sysprvFuncMax = sysprvLeave // Maximum legal function

	// MODE VALUES (sysprvPktMode)
	sysprvSuser    = 1 // Superuser mode
	sysprvSprocess = 2 // Superprocess mode
	sysprvSysmgr   = 3 // System manager mode

	sysprvModeMin = sysprvSuser  // Minimum legal mode
	sysprvModeMax = sysprvSysmgr // Maximum legal mode

	// FLAG WORD BITS (BITS IN sysprvPktFlags)
	sysprvFlagsCaller      = 0 // Caller has privelege on
	sysprvFlagsCallersExcl = 1 // Caller has priv on exclusively
	sysprvFlagsOthers      = 2 // Others have privilege on
	sysprvFlagsOthersExcl  = 3 // Others have priv on exclusively

	// FLAG WORD BIT POINTERS (BITS IN sysprvPktFlags)
	sysprvBPFlagsCaller     = (sysprvPktFlags * 16.) + sysprvFlagsCaller
	sysprvBPFlagsCallerExcl = (sysprvPktFlags * 16.) + sysprvFlagsCallersExcl
	sysprvBPFlagsOthers     = (sysprvPktFlags * 16.) + sysprvFlagsOthers
	sysprvBPFlagsOthersExcl = (sysprvPktFlags * 16.) + sysprvFlagsOthersExcl
)
