// Copyright ©2017-2020  Steve Merrony

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

package dg

// ASCII and DASHER special characters
const (
	ASCIIBEL = 007
	ASCIIBS  = 010
	ASCIITAB = 011
	ASCIINL  = 012
	ASCIIFF  = 014
	ASCIICR  = 015
	ASCIIESC = 033
	ASCIISPC = 040

	DasherERASEEOL        = 013
	DasherERASEPAGE       = 014
	DasherCURSORLEFT      = 031
	DasherWRITEWINDOWADDR = 020 //followed by col then row
	DasherDIMON           = 034
	DasherDIMOFF          = 035
	DasherUNDERLINE       = 024
	DasherNORMAL          = 025
	DasherDELETE          = 0177
)
