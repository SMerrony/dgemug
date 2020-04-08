// Code generated by dginstr.go; DO NOT EDIT.

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

package mvcpu

// Instruction Types
const (
	NOVA_MEMREF = iota
	NOVA_OP
	NOVA_IO
	NOVA_MATH
	NOVA_PC
	ECLIPSE_FPU
	ECLIPSE_MEMREF
	ECLIPSE_OP
	ECLIPSE_PC
	ECLIPSE_STACK
	EAGLE_DECIMAL
	EAGLE_IO
	EAGLE_PC
	EAGLE_OP
	EAGLE_MEMREF
	EAGLE_STACK
	EAGLE_FPU
)

// Instruction Formats
const (
	DERR_FMT = iota
	IMM_MODE_2_WORD_FMT
	IMM_ONEACC_FMT
	IO_FLAGS_DEV_FMT
	IO_TEST_DEV_FMT
	LNDO_4_WORD_FMT
	NOACC_MODE_2_WORD_FMT
	NOACC_MODE_3_WORD_FMT
	NOACC_MODE_IMM_IND_3_WORD_FMT
	NOACC_MODE_IND_2_WORD_E_FMT
	NOACC_MODE_IND_2_WORD_X_FMT
	NOACC_MODE_IND_3_WORD_FMT
	NOACC_MODE_IND_3_WORD_XCALL_FMT
	NOACC_MODE_IND_4_WORD_FMT
	NOVA_DATA_IO_FMT
	NOVA_NOACC_EFF_ADDR_FMT
	NOVA_ONEACC_EFF_ADDR_FMT
	NOVA_TWOACC_MULT_OP_FMT
	ONEACC_IMM_2_WORD_FMT
	ONEACC_IMMWD_2_WORD_FMT
	ONEACC_IMM_3_WORD_FMT
	ONEACC_IMMDWD_3_WORD_FMT
	ONEACC_MODE_2_WORD_E_FMT
	ONEACC_MODE_2_WORD_X_B_FMT
	ONEACC_MODE_3_WORD_FMT
	ONEACC_MODE_IND_2_WORD_E_FMT
	ONEACC_MODE_IND_2_WORD_X_FMT
	ONEACC_MODE_IND_3_WORD_FMT
	ONEACC_1_WORD_FMT
	UNIQUE_1_WORD_FMT
	UNIQUE_2_WORD_FMT
	SPLIT_8BIT_DISP_FMT
	THREE_WORD_DO_FMT
	TWOACC_1_WORD_FMT
	TWOACC_IMM_2_WORD_FMT
	WIDE_DEC_SPECIAL_FMT
	WSKB_FMT
)

// Instruction Mnemonic Consts
const (
	instrADC = iota
	instrADD
	instrADDI
	instrADI
	instrANC
	instrAND
	instrANDI
	instrBAM
	instrBKPT
	instrBLM
	instrBTO
	instrBTZ
	instrCIO
	instrCIOI
	instrCLM
	instrCMP
	instrCMT
	instrCMV
	instrCOB
	instrCOM
	instrCRYTC
	instrCRYTO
	instrCRYTZ
	instrCTR
	instrCVWN
	instrDAD
	instrDEQUE
	instrDERR
	instrDHXL
	instrDHXR
	instrDIA
	instrDIB
	instrDIC
	instrDIV
	instrDIVS
	instrDIVX
	instrDLSH
	instrDOA
	instrDOB
	instrDOC
	instrDSB
	instrDSPA
	instrDSZ
	instrDSZTS
	instrECLID
	instrEDIT
	instrEDSZ
	instrEISZ
	instrEJMP
	instrEJSR
	instrELDA
	instrELDB
	instrELEF
	instrENQH
	instrENQT
	instrESTA
	instrESTB
	instrFAD
	instrFAS
	instrFCLE
	instrFCMP
	instrFFAS
	instrFLAS
	instrFLDS
	instrFNEG
	instrFNS
	instrFPOP
	instrFPSH
	instrFSA
	instrFSEQ
	instrFSGE
	instrFSGT
	instrFSLE
	instrFSLT
	instrFSNE
	instrFSS
	instrFSST
	instrFSTS
	instrFTD
	instrFTE
	instrFXTD
	instrFXTE
	instrHALT
	instrHLV
	instrHXL
	instrHXR
	instrINC
	instrINTA
	instrINTDS
	instrINTEN
	instrIOR
	instrIORI
	instrIORST
	instrISZ
	instrISZTS
	instrJMP
	instrJSR
	instrLCALL
	instrLCPID
	instrLDA
	instrLDAFP
	instrLDASB
	instrLDASL
	instrLDASP
	instrLDATS
	instrLDB
	instrLDSP
	instrLEF
	instrLFDMS
	instrLJMP
	instrLJSR
	instrLLDB
	instrLLEF
	instrLLEFB
	instrLMRF
	instrLNADD
	instrLNADI
	instrLNDIV
	instrLNDO
	instrLNDSZ
	instrLNISZ
	instrLNLDA
	instrLNMUL
	instrLNSBI
	instrLNSTA
	instrLNSUB
	instrLOB
	instrLPEF
	instrLPEFB
	instrLPHY
	instrLPSHJ
	instrLPSR
	instrLRB
	instrLSH
	instrLSTB
	instrLWADD
	instrLWDO
	instrLWDSZ
	instrLWISZ
	instrLWLDA
	instrLWSTA
	instrLWSUB
	instrMOV
	instrMSP
	instrMUL
	instrMULS
	instrNADD
	instrNADDI
	instrNADI
	instrNCLID
	instrNEG
	instrNIO
	instrNLDAI
	instrNMUL
	instrNSALA
	instrNSANA
	instrNSBI
	instrNSUB
	instrPIO
	instrPOP
	instrPOPB
	instrPOPJ
	instrPRTSEL
	instrPSH
	instrPSHJ
	instrPSHR
	instrREADS
	instrRSTR
	instrRTN
	instrSAVE
	instrSBI
	instrSEX
	instrSGE
	instrSGT
	instrSKP
	instrSNB
	instrSNOVR
	instrSPSR
	instrSPTE
	instrSSPT
	instrSTA
	instrSTAFP
	instrSTASB
	instrSTASL
	instrSTASP
	instrSTATS
	instrSTB
	instrSUB
	instrSZB
	instrSZBO
	instrWADC
	instrWADD
	instrWADDI
	instrWADI
	instrWANC
	instrWAND
	instrWANDI
	instrWASH
	instrWASHI
	instrWBLM
	instrWBR
	instrWBTO
	instrWBTZ
	instrWCLM
	instrWCMP
	instrWCMV
	instrWCOM
	instrWCST
	instrWCTR
	instrWDecOp
	instrWDIV
	instrWDIVS
	instrWFLAD
	instrWFPOP
	instrWFPSH
	instrWHLV
	instrWINC
	instrWIOR
	instrWIORI
	instrWLDAI
	instrWLDB
	instrWLMP
	instrWLSH
	instrWLSHI
	instrWLSI
	instrWMESS
	instrWMOV
	instrWMOVR
	instrWMSP
	instrWMUL
	instrWMULS
	instrWNADI
	instrWNEG
	instrWPOP
	instrWPOPB
	instrWPOPJ
	instrWPSH
	instrWRTN
	instrWSAVR
	instrWSAVS
	instrWSBI
	instrWSEQ
	instrWSEQI
	instrWSGE
	instrWSGT
	instrWSGTI
	instrWSKBO
	instrWSKBZ
	instrWSLE
	instrWSLEI
	instrWSLT
	instrWSNB
	instrWSNE
	instrWSNEI
	instrWSSVR
	instrWSSVS
	instrWSTB
	instrWSTI
	instrWSUB
	instrWSZB
	instrWSZBO
	instrWUSGT
	instrWUGTI
	instrWXCH
	instrWXORI
	instrXCALL
	instrXCH
	instrXCT
	instrXFLDS
	instrXJMP
	instrXJSR
	instrXLDB
	instrXLEF
	instrXLEFB
	instrXNADD
	instrXNADI
	instrXNDO
	instrXNDSZ
	instrXNISZ
	instrXNLDA
	instrXNSBI
	instrXNSTA
	instrXNSUB
	instrXOR
	instrXORI
	instrXPEF
	instrXPEFB
	instrXPSHJ
	instrXSTB
	instrXWADD
	instrXWADI
	instrXWDSZ
	instrXWISZ
	instrXWLDA
	instrXWMUL
	instrXWSBI
	instrXWSTA
	instrXWSUB
	instrZEX
)

// InstructionsInit initialises the instruction characterstics for each instruction(
func InstructionsInit() {
	instructionSet[instrADC] = instrChars{"ADC", 0x8400, 0x8700, 1, NOVA_TWOACC_MULT_OP_FMT, NOVA_OP, 0}
	instructionSet[instrADD] = instrChars{"ADD", 0x8600, 0x8700, 1, NOVA_TWOACC_MULT_OP_FMT, NOVA_OP, 0}
	instructionSet[instrADDI] = instrChars{"ADDI", 0xe7f8, 0xe7ff, 2, ONEACC_IMM_2_WORD_FMT, ECLIPSE_OP, 0}
	instructionSet[instrADI] = instrChars{"ADI", 0x8008, 0x87ff, 1, IMM_ONEACC_FMT, ECLIPSE_OP, 0}
	instructionSet[instrANC] = instrChars{"ANC", 0x8188, 0x87ff, 1, TWOACC_1_WORD_FMT, ECLIPSE_OP, 0}
	instructionSet[instrAND] = instrChars{"AND", 0x8700, 0x8700, 1, NOVA_TWOACC_MULT_OP_FMT, NOVA_OP, 0}
	instructionSet[instrANDI] = instrChars{"ANDI", 0xc7f8, 0xe7ff, 2, ONEACC_IMMWD_2_WORD_FMT, ECLIPSE_OP, 0}
	instructionSet[instrBAM] = instrChars{"BAM", 0x97c8, 0xffff, 1, UNIQUE_1_WORD_FMT, ECLIPSE_MEMREF, 0}
	instructionSet[instrBKPT] = instrChars{"BKPT", 0xc789, 0xffff, 1, UNIQUE_1_WORD_FMT, EAGLE_PC, 0}
	instructionSet[instrBLM] = instrChars{"BLM", 0xb7c8, 0xffff, 1, UNIQUE_1_WORD_FMT, ECLIPSE_MEMREF, 0}
	instructionSet[instrBTO] = instrChars{"BTO", 0x8408, 0x87ff, 1, TWOACC_1_WORD_FMT, ECLIPSE_MEMREF, 0}
	instructionSet[instrBTZ] = instrChars{"BTZ", 0x8448, 0x87ff, 1, TWOACC_1_WORD_FMT, ECLIPSE_MEMREF, 0}
	instructionSet[instrCIO] = instrChars{"CIO", 0x85e9, 0x87ff, 1, TWOACC_1_WORD_FMT, EAGLE_IO, 0}
	instructionSet[instrCIOI] = instrChars{"CIOI", 0x85f9, 0x87ff, 2, TWOACC_IMM_2_WORD_FMT, EAGLE_IO, 0}
	instructionSet[instrCLM] = instrChars{"CLM", 0x84f8, 0x87ff, 1, TWOACC_1_WORD_FMT, ECLIPSE_PC, 0}
	instructionSet[instrCMP] = instrChars{"CMP", 0xdfa8, 0xffff, 1, UNIQUE_1_WORD_FMT, ECLIPSE_MEMREF, 0}
	instructionSet[instrCMT] = instrChars{"CMT", 0xefa8, 0xffff, 1, UNIQUE_1_WORD_FMT, ECLIPSE_MEMREF, 0}
	instructionSet[instrCMV] = instrChars{"CMV", 0xd7a8, 0xffff, 1, UNIQUE_1_WORD_FMT, ECLIPSE_MEMREF, 0}
	instructionSet[instrCOB] = instrChars{"COB", 0x8588, 0x87ff, 1, TWOACC_1_WORD_FMT, ECLIPSE_OP, 0}
	instructionSet[instrCOM] = instrChars{"COM", 0x8000, 0x8700, 1, NOVA_TWOACC_MULT_OP_FMT, NOVA_OP, 0}
	instructionSet[instrCRYTC] = instrChars{"CRYTC", 0xa7e9, 0xffff, 1, UNIQUE_1_WORD_FMT, EAGLE_OP, 0}
	instructionSet[instrCRYTO] = instrChars{"CRYTO", 0xa7c9, 0xffff, 1, UNIQUE_1_WORD_FMT, EAGLE_OP, 0}
	instructionSet[instrCRYTZ] = instrChars{"CRYTZ", 0xa7d9, 0xffff, 1, UNIQUE_1_WORD_FMT, EAGLE_OP, 0}
	instructionSet[instrCTR] = instrChars{"CTR", 0xe7a8, 0xffff, 1, UNIQUE_1_WORD_FMT, ECLIPSE_OP, 0}
	instructionSet[instrCVWN] = instrChars{"CVWN", 0xe669, 0xe7ff, 1, ONEACC_1_WORD_FMT, EAGLE_OP, 0}
	instructionSet[instrDAD] = instrChars{"DAD", 0x8088, 0x87ff, 1, TWOACC_1_WORD_FMT, ECLIPSE_OP, 0}
	instructionSet[instrDEQUE] = instrChars{"DEQUE", 0xe7c9, 0xffff, 1, UNIQUE_1_WORD_FMT, EAGLE_OP, 0}
	instructionSet[instrDERR] = instrChars{"DERR", 0x8f09, 0x8fcf, 1, DERR_FMT, EAGLE_PC, 0}
	instructionSet[instrDHXL] = instrChars{"DHXL", 0x8388, 0x87ff, 1, IMM_ONEACC_FMT, ECLIPSE_OP, 0}
	instructionSet[instrDHXR] = instrChars{"DHXR", 0x83c8, 0x87ff, 1, IMM_ONEACC_FMT, ECLIPSE_OP, 0}
	instructionSet[instrDIA] = instrChars{"DIA", 0x6100, 0xe700, 1, NOVA_DATA_IO_FMT, NOVA_IO, 0}
	instructionSet[instrDIB] = instrChars{"DIB", 0x6300, 0xe700, 1, NOVA_DATA_IO_FMT, NOVA_IO, 0}
	instructionSet[instrDIC] = instrChars{"DIC", 0x6500, 0xe700, 1, NOVA_DATA_IO_FMT, NOVA_IO, 0}
	instructionSet[instrDIV] = instrChars{"DIV", 0xd7c8, 0xffff, 1, UNIQUE_1_WORD_FMT, NOVA_MATH, 0}
	instructionSet[instrDIVS] = instrChars{"DIVS", 0xdfc8, 0xffff, 1, UNIQUE_1_WORD_FMT, ECLIPSE_OP, 0}
	instructionSet[instrDIVX] = instrChars{"DIVX", 0xbfc8, 0xffff, 1, UNIQUE_1_WORD_FMT, ECLIPSE_OP, 0}
	instructionSet[instrDLSH] = instrChars{"DLSH", 0x82c8, 0x87ff, 1, TWOACC_1_WORD_FMT, ECLIPSE_OP, 0}
	instructionSet[instrDOA] = instrChars{"DOA", 0x6200, 0xe700, 1, NOVA_DATA_IO_FMT, NOVA_IO, 0}
	instructionSet[instrDOB] = instrChars{"DOB", 0x6400, 0xe700, 1, NOVA_DATA_IO_FMT, NOVA_IO, 0}
	instructionSet[instrDOC] = instrChars{"DOC", 0x6600, 0xe700, 1, NOVA_DATA_IO_FMT, NOVA_IO, 0}
	instructionSet[instrDSB] = instrChars{"DSB", 0x80c8, 0x87ff, 1, TWOACC_1_WORD_FMT, ECLIPSE_OP, 0}
	instructionSet[instrDSPA] = instrChars{"DSPA", 0xc478, 0xe4ff, 2, ONEACC_MODE_IND_2_WORD_E_FMT, ECLIPSE_PC, 1}
	instructionSet[instrDSZ] = instrChars{"DSZ", 0x1800, 0xf800, 1, NOVA_NOACC_EFF_ADDR_FMT, NOVA_MEMREF, 0}
	instructionSet[instrDSZTS] = instrChars{"DSZTS", 0xc7d9, 0xffff, 1, UNIQUE_1_WORD_FMT, EAGLE_PC, 0}
	instructionSet[instrECLID] = instrChars{"ECLID", 0xffc8, 0xffff, 1, UNIQUE_1_WORD_FMT, EAGLE_IO, 0}
	instructionSet[instrEDIT] = instrChars{"EDIT", 0xf7a8, 0xffff, 1, UNIQUE_1_WORD_FMT, ECLIPSE_OP, 0}
	instructionSet[instrEDSZ] = instrChars{"EDSZ", 0x9c38, 0xfcff, 2, NOACC_MODE_IND_2_WORD_E_FMT, ECLIPSE_PC, 1}
	instructionSet[instrEISZ] = instrChars{"EISZ", 0x9438, 0xfcff, 2, NOACC_MODE_IND_2_WORD_E_FMT, ECLIPSE_PC, 1}
	instructionSet[instrEJMP] = instrChars{"EJMP", 0x8438, 0xfcff, 2, NOACC_MODE_IND_2_WORD_E_FMT, ECLIPSE_PC, 1}
	instructionSet[instrEJSR] = instrChars{"EJSR", 0x8c38, 0xfcff, 2, NOACC_MODE_IND_2_WORD_E_FMT, ECLIPSE_PC, 1}
	instructionSet[instrELDA] = instrChars{"ELDA", 0xa438, 0xe4ff, 2, ONEACC_MODE_IND_2_WORD_E_FMT, ECLIPSE_MEMREF, 1}
	instructionSet[instrELDB] = instrChars{"ELDB", 0x8478, 0xe4ff, 2, ONEACC_MODE_IND_2_WORD_E_FMT, ECLIPSE_MEMREF, 1}
	instructionSet[instrELEF] = instrChars{"ELEF", 0xe438, 0xe4ff, 2, ONEACC_MODE_IND_2_WORD_E_FMT, ECLIPSE_MEMREF, 1}
	instructionSet[instrENQH] = instrChars{"ENQH", 0xc7e9, 0xffff, 1, UNIQUE_1_WORD_FMT, EAGLE_OP, 0}
	instructionSet[instrENQT] = instrChars{"ENQT", 0xc7f9, 0xffff, 1, UNIQUE_1_WORD_FMT, EAGLE_OP, 0}
	instructionSet[instrESTA] = instrChars{"ESTA", 0xc438, 0xe4ff, 2, ONEACC_MODE_IND_2_WORD_E_FMT, ECLIPSE_MEMREF, 1}
	instructionSet[instrESTB] = instrChars{"ESTB", 0xa478, 0xe4ff, 2, ONEACC_MODE_2_WORD_E_FMT, ECLIPSE_OP, 1}
	instructionSet[instrFAD] = instrChars{"FAD", 0x8068, 0x87ff, 1, TWOACC_1_WORD_FMT, ECLIPSE_FPU, 0}
	instructionSet[instrFAS] = instrChars{"FAS", 0x8028, 0x87ff, 1, TWOACC_1_WORD_FMT, ECLIPSE_FPU, 0}
	instructionSet[instrFCLE] = instrChars{"FCLE", 0xd6e8, 0xffff, 1, UNIQUE_1_WORD_FMT, ECLIPSE_FPU, 0}
	instructionSet[instrFCMP] = instrChars{"FCMP", 0x8728, 0x87ff, 1, TWOACC_1_WORD_FMT, ECLIPSE_FPU, 0}
	instructionSet[instrFFAS] = instrChars{"FFAS", 0x85a8, 0x87ff, 1, TWOACC_1_WORD_FMT, ECLIPSE_FPU, 0}
	instructionSet[instrFLAS] = instrChars{"FLAS", 0x8528, 0x87ff, 1, TWOACC_1_WORD_FMT, ECLIPSE_FPU, 0}
	instructionSet[instrFLDS] = instrChars{"FLDS", 0x8428, 0x87ff, 2, ONEACC_MODE_IND_2_WORD_X_FMT, ECLIPSE_FPU, 1}
	instructionSet[instrFNEG] = instrChars{"FNEG", 0xe628, 0xe7ff, 1, ONEACC_1_WORD_FMT, ECLIPSE_FPU, 0}
	instructionSet[instrFNS] = instrChars{"FNS", 0x86a8, 0xffff, 1, UNIQUE_1_WORD_FMT, ECLIPSE_PC, 0}
	instructionSet[instrFPOP] = instrChars{"FPOP", 0xeee8, 0xffff, 1, UNIQUE_1_WORD_FMT, ECLIPSE_STACK, 0}
	instructionSet[instrFPSH] = instrChars{"FPSH", 0xe6e8, 0xffff, 1, UNIQUE_1_WORD_FMT, ECLIPSE_STACK, 0}
	instructionSet[instrFSA] = instrChars{"FSA", 0x8ea8, 0xffff, 1, UNIQUE_1_WORD_FMT, ECLIPSE_PC, 0}
	instructionSet[instrFSEQ] = instrChars{"FSEQ", 0x96a8, 0xffff, 1, UNIQUE_1_WORD_FMT, ECLIPSE_FPU, 0}
	instructionSet[instrFSGE] = instrChars{"FSGE", 0xaea8, 0xffff, 1, UNIQUE_1_WORD_FMT, ECLIPSE_FPU, 0}
	instructionSet[instrFSGT] = instrChars{"FSGT", 0xbea8, 0xffff, 1, UNIQUE_1_WORD_FMT, ECLIPSE_FPU, 0}
	instructionSet[instrFSLE] = instrChars{"FSLE", 0xb6a8, 0xffff, 1, UNIQUE_1_WORD_FMT, ECLIPSE_FPU, 0}
	instructionSet[instrFSLT] = instrChars{"FSLT", 0xa6a8, 0xffff, 1, UNIQUE_1_WORD_FMT, ECLIPSE_FPU, 0}
	instructionSet[instrFSNE] = instrChars{"FSNE", 0x9ea8, 0xffff, 1, UNIQUE_1_WORD_FMT, ECLIPSE_FPU, 0}
	instructionSet[instrFSS] = instrChars{"FSS", 0x80a8, 0x87ff, 1, TWOACC_1_WORD_FMT, ECLIPSE_FPU, 0}
	instructionSet[instrFSST] = instrChars{"FSST", 0x86e8, 0xe7ff, 2, NOACC_MODE_IND_2_WORD_X_FMT, ECLIPSE_FPU, 0}
	instructionSet[instrFSTS] = instrChars{"FSTS", 0x84a8, 0x87ff, 2, ONEACC_MODE_IND_2_WORD_X_FMT, ECLIPSE_FPU, 0}
	instructionSet[instrFTD] = instrChars{"FTD", 0xcee8, 0xffff, 1, UNIQUE_1_WORD_FMT, ECLIPSE_FPU, 0}
	instructionSet[instrFTE] = instrChars{"FTE", 0xc6e8, 0xffff, 1, UNIQUE_1_WORD_FMT, ECLIPSE_FPU, 0}
	instructionSet[instrFXTD] = instrChars{"FXTD", 0xa779, 0xffff, 1, UNIQUE_1_WORD_FMT, ECLIPSE_OP, 0}
	instructionSet[instrFXTE] = instrChars{"FXTE", 0xc749, 0xffff, 1, UNIQUE_1_WORD_FMT, ECLIPSE_OP, 0}
	instructionSet[instrHALT] = instrChars{"HALT", 0x647f, 0xffff, 1, UNIQUE_1_WORD_FMT, NOVA_IO, 0}
	instructionSet[instrHLV] = instrChars{"HLV", 0xc6f8, 0xe7ff, 1, ONEACC_1_WORD_FMT, ECLIPSE_OP, 0}
	instructionSet[instrHXL] = instrChars{"HXL", 0x8308, 0x87ff, 1, IMM_ONEACC_FMT, ECLIPSE_OP, 0}
	instructionSet[instrHXR] = instrChars{"HXR", 0x8348, 0x87ff, 1, IMM_ONEACC_FMT, ECLIPSE_OP, 0}
	instructionSet[instrINC] = instrChars{"INC", 0x8300, 0x8700, 1, NOVA_TWOACC_MULT_OP_FMT, NOVA_OP, 0}
	instructionSet[instrINTA] = instrChars{"INTA", 0x633f, 0xe7ff, 1, ONEACC_1_WORD_FMT, NOVA_IO, 0}
	instructionSet[instrINTDS] = instrChars{"INTDS", 0x60bf, 0xffff, 1, UNIQUE_1_WORD_FMT, NOVA_IO, 0}
	instructionSet[instrINTEN] = instrChars{"INTEN", 0x607f, 0xffff, 1, UNIQUE_1_WORD_FMT, NOVA_IO, 0}
	instructionSet[instrIOR] = instrChars{"IOR", 0x8108, 0x87ff, 1, TWOACC_1_WORD_FMT, ECLIPSE_OP, 0}
	instructionSet[instrIORI] = instrChars{"IORI", 0x87f8, 0xe7ff, 2, ONEACC_IMMWD_2_WORD_FMT, ECLIPSE_OP, 0}
	instructionSet[instrIORST] = instrChars{"IORST", 0x653f, 0xe73f, 1, ONEACC_1_WORD_FMT, NOVA_IO, 0}
	instructionSet[instrISZ] = instrChars{"ISZ", 0x1000, 0xf800, 1, NOVA_NOACC_EFF_ADDR_FMT, NOVA_MEMREF, 0}
	instructionSet[instrISZTS] = instrChars{"ISZTS", 0xc7c9, 0xffff, 1, UNIQUE_1_WORD_FMT, EAGLE_PC, 0}
	instructionSet[instrJMP] = instrChars{"JMP", 0x0, 0xf800, 1, NOVA_NOACC_EFF_ADDR_FMT, NOVA_PC, 0}
	instructionSet[instrJSR] = instrChars{"JSR", 0x800, 0xf800, 1, NOVA_NOACC_EFF_ADDR_FMT, NOVA_PC, 0}
	instructionSet[instrLCALL] = instrChars{"LCALL", 0xa6c9, 0xe7ff, 4, NOACC_MODE_IND_4_WORD_FMT, EAGLE_PC, 1}
	instructionSet[instrLCPID] = instrChars{"LCPID", 0x8759, 0xffff, 1, UNIQUE_1_WORD_FMT, EAGLE_IO, 0}
	instructionSet[instrLDA] = instrChars{"LDA", 0x2000, 0xe000, 1, NOVA_ONEACC_EFF_ADDR_FMT, NOVA_MEMREF, 0}
	instructionSet[instrLDAFP] = instrChars{"LDAFP", 0xc669, 0xe7ff, 1, ONEACC_1_WORD_FMT, EAGLE_STACK, 0}
	instructionSet[instrLDASB] = instrChars{"LDASB", 0xc649, 0xe7ff, 1, ONEACC_1_WORD_FMT, EAGLE_STACK, 0}
	instructionSet[instrLDASL] = instrChars{"LDASL", 0xa669, 0xe7ff, 1, ONEACC_1_WORD_FMT, EAGLE_STACK, 0}
	instructionSet[instrLDASP] = instrChars{"LDASP", 0xa649, 0xe7ff, 1, ONEACC_1_WORD_FMT, EAGLE_STACK, 0}
	instructionSet[instrLDATS] = instrChars{"LDATS", 0x8649, 0xe7ff, 1, ONEACC_1_WORD_FMT, EAGLE_STACK, 0}
	instructionSet[instrLDB] = instrChars{"LDB", 0x85c8, 0x87ff, 1, TWOACC_1_WORD_FMT, ECLIPSE_MEMREF, 0}
	instructionSet[instrLDSP] = instrChars{"LDSP", 0x8519, 0x87ff, 3, ONEACC_MODE_IND_3_WORD_FMT, EAGLE_PC, 1}
	instructionSet[instrLEF] = instrChars{"LEF", 0x6000, 0xe000, 1, NOVA_ONEACC_EFF_ADDR_FMT, ECLIPSE_MEMREF, 0}
	instructionSet[instrLFDMS] = instrChars{"LFDMS", 0x81e9, 0x87ff, 3, ONEACC_MODE_IND_3_WORD_FMT, EAGLE_FPU, 1}
	instructionSet[instrLJMP] = instrChars{"LJMP", 0xa6d9, 0xe7ff, 3, NOACC_MODE_IND_3_WORD_FMT, EAGLE_PC, 1}
	instructionSet[instrLJSR] = instrChars{"LJSR", 0xa6e9, 0xe7ff, 3, NOACC_MODE_IND_3_WORD_FMT, EAGLE_PC, 1}
	instructionSet[instrLLDB] = instrChars{"LLDB", 0x84c9, 0x87ff, 3, ONEACC_MODE_3_WORD_FMT, EAGLE_MEMREF, 1}
	instructionSet[instrLLEF] = instrChars{"LLEF", 0x83e9, 0x87ff, 3, ONEACC_MODE_IND_3_WORD_FMT, EAGLE_MEMREF, 1}
	instructionSet[instrLLEFB] = instrChars{"LLEFB", 0x84e9, 0x87ff, 3, ONEACC_MODE_3_WORD_FMT, EAGLE_MEMREF, 1}
	instructionSet[instrLMRF] = instrChars{"LMRF", 0x87c9, 0xffff, 1, UNIQUE_1_WORD_FMT, EAGLE_OP, 0}
	instructionSet[instrLNADD] = instrChars{"LNADD", 0x8218, 0x87ff, 3, ONEACC_MODE_IND_3_WORD_FMT, EAGLE_OP, 1}
	instructionSet[instrLNADI] = instrChars{"LNADI", 0x8618, 0x87ff, 3, NOACC_MODE_IMM_IND_3_WORD_FMT, EAGLE_MEMREF, 1}
	instructionSet[instrLNDIV] = instrChars{"LNDIV", 0x82d8, 0x87ff, 3, ONEACC_MODE_IND_3_WORD_FMT, EAGLE_OP, 1}
	instructionSet[instrLNDO] = instrChars{"LNDO", 0x8698, 0x87ff, 4, LNDO_4_WORD_FMT, EAGLE_PC, 1}
	instructionSet[instrLNDSZ] = instrChars{"LNDSZ", 0x86d9, 0xe7ff, 3, NOACC_MODE_IND_3_WORD_FMT, EAGLE_PC, 1}
	instructionSet[instrLNISZ] = instrChars{"LNISZ", 0x86c9, 0xe7ff, 3, NOACC_MODE_IND_3_WORD_FMT, EAGLE_PC, 1}
	instructionSet[instrLNLDA] = instrChars{"LNLDA", 0x83c9, 0x87ff, 3, ONEACC_MODE_IND_3_WORD_FMT, EAGLE_MEMREF, 1}
	instructionSet[instrLNMUL] = instrChars{"LNMUL", 0x8298, 0x87ff, 3, ONEACC_MODE_IND_3_WORD_FMT, EAGLE_MEMREF, 1}
	instructionSet[instrLNSBI] = instrChars{"LNSBI", 0x8658, 0x87ff, 3, NOACC_MODE_IMM_IND_3_WORD_FMT, EAGLE_MEMREF, 1}
	instructionSet[instrLNSTA] = instrChars{"LNSTA", 0x83d9, 0x87ff, 3, ONEACC_MODE_IND_3_WORD_FMT, EAGLE_MEMREF, 1}
	instructionSet[instrLNSUB] = instrChars{"LNSUB", 0x8258, 0x87ff, 3, ONEACC_MODE_IND_3_WORD_FMT, EAGLE_MEMREF, 1}
	instructionSet[instrLOB] = instrChars{"LOB", 0x8508, 0x87ff, 1, TWOACC_1_WORD_FMT, ECLIPSE_OP, 0}
	instructionSet[instrLPEF] = instrChars{"LPEF", 0xa6f9, 0xe7ff, 3, NOACC_MODE_IND_3_WORD_FMT, EAGLE_STACK, 1}
	instructionSet[instrLPEFB] = instrChars{"LPEFB", 0xc6f9, 0xe7ff, 3, NOACC_MODE_3_WORD_FMT, EAGLE_STACK, 1}
	instructionSet[instrLPHY] = instrChars{"LPHY", 0x87e9, 0xffff, 1, UNIQUE_1_WORD_FMT, EAGLE_OP, 0}
	instructionSet[instrLPSHJ] = instrChars{"LPSHJ", 0xC6C9, 0xE7FF, 3, NOACC_MODE_IND_3_WORD_FMT, EAGLE_PC, 1}
	instructionSet[instrLPSR] = instrChars{"LPSR", 0xa799, 0xffff, 1, UNIQUE_1_WORD_FMT, EAGLE_OP, 0}
	instructionSet[instrLRB] = instrChars{"LRB", 0x8548, 0x87ff, 1, TWOACC_1_WORD_FMT, ECLIPSE_OP, 0}
	instructionSet[instrLSH] = instrChars{"LSH", 0x8288, 0x87ff, 1, TWOACC_1_WORD_FMT, ECLIPSE_OP, 0}
	instructionSet[instrLSTB] = instrChars{"LSTB", 0x84d9, 0x87ff, 3, ONEACC_MODE_3_WORD_FMT, EAGLE_MEMREF, 1}
	instructionSet[instrLWADD] = instrChars{"LWADD", 0x8318, 0x87ff, 3, ONEACC_MODE_IND_3_WORD_FMT, EAGLE_MEMREF, 1}
	instructionSet[instrLWDO] = instrChars{"LWDO", 0x8798, 0x87ff, 4, LNDO_4_WORD_FMT, EAGLE_PC, 1}
	instructionSet[instrLWDSZ] = instrChars{"LWDSZ", 0x86f9, 0xe7ff, 3, NOACC_MODE_IND_3_WORD_FMT, EAGLE_PC, 1}
	instructionSet[instrLWISZ] = instrChars{"LWISZ", 0x86e9, 0xe7ff, 3, NOACC_MODE_IND_3_WORD_FMT, EAGLE_PC, 1}
	instructionSet[instrLWLDA] = instrChars{"LWLDA", 0x83f9, 0x87ff, 3, ONEACC_MODE_IND_3_WORD_FMT, EAGLE_MEMREF, 1}
	instructionSet[instrLWSTA] = instrChars{"LWSTA", 0x84f9, 0x87ff, 3, ONEACC_MODE_IND_3_WORD_FMT, EAGLE_MEMREF, 1}
	instructionSet[instrLWSUB] = instrChars{"LWSUB", 0x8358, 0x87ff, 3, ONEACC_MODE_IND_3_WORD_FMT, EAGLE_MEMREF, 1}
	instructionSet[instrMOV] = instrChars{"MOV", 0x8200, 0x8700, 1, NOVA_TWOACC_MULT_OP_FMT, NOVA_OP, 0}
	instructionSet[instrMSP] = instrChars{"MSP", 0x86f8, 0xe7ff, 1, ONEACC_1_WORD_FMT, ECLIPSE_STACK, 0}
	instructionSet[instrMUL] = instrChars{"MUL", 0xc7c8, 0xffff, 1, UNIQUE_1_WORD_FMT, NOVA_MATH, 0}
	instructionSet[instrMULS] = instrChars{"MULS", 0xcfc8, 0xffff, 1, UNIQUE_1_WORD_FMT, ECLIPSE_OP, 0}
	instructionSet[instrNADD] = instrChars{"NADD", 0x8049, 0x87ff, 1, TWOACC_1_WORD_FMT, EAGLE_OP, 0}
	instructionSet[instrNADDI] = instrChars{"NADDI", 0xc639, 0xe7ff, 2, ONEACC_IMM_2_WORD_FMT, EAGLE_OP, 0}
	instructionSet[instrNADI] = instrChars{"NADI", 0x8599, 0x87ff, 1, IMM_ONEACC_FMT, EAGLE_OP, 0}
	instructionSet[instrNCLID] = instrChars{"NCLID", 0x683f, 0xffff, 1, UNIQUE_1_WORD_FMT, EAGLE_IO, 0}
	instructionSet[instrNEG] = instrChars{"NEG", 0x8100, 0x8700, 1, NOVA_TWOACC_MULT_OP_FMT, NOVA_OP, 0}
	instructionSet[instrNIO] = instrChars{"NIO", 0x6000, 0xff00, 1, IO_FLAGS_DEV_FMT, NOVA_IO, 0}
	instructionSet[instrNLDAI] = instrChars{"NLDAI", 0xc629, 0xe7ff, 2, ONEACC_IMM_2_WORD_FMT, EAGLE_OP, 0}
	instructionSet[instrNMUL] = instrChars{"NMUL", 0x8069, 0x87ff, 1, TWOACC_1_WORD_FMT, EAGLE_OP, 0}
	instructionSet[instrNSALA] = instrChars{"NSALA", 0xe609, 0xe7ff, 2, ONEACC_IMM_2_WORD_FMT, EAGLE_PC, 0}
	instructionSet[instrNSANA] = instrChars{"NSANA", 0xe629, 0xe7ff, 2, ONEACC_IMM_2_WORD_FMT, EAGLE_PC, 0}
	instructionSet[instrNSBI] = instrChars{"NSBI", 0x85a9, 0x87ff, 1, IMM_ONEACC_FMT, EAGLE_OP, 0}
	instructionSet[instrNSUB] = instrChars{"NSUB", 0x8059, 0x87ff, 1, TWOACC_1_WORD_FMT, EAGLE_OP, 0}
	instructionSet[instrPIO] = instrChars{"PIO", 0x85d9, 0x87ff, 1, TWOACC_1_WORD_FMT, EAGLE_IO, 0}
	instructionSet[instrPOP] = instrChars{"POP", 0x8688, 0x87ff, 1, TWOACC_1_WORD_FMT, ECLIPSE_STACK, 0}
	instructionSet[instrPOPB] = instrChars{"POPB", 0x8fc8, 0xffff, 1, UNIQUE_1_WORD_FMT, ECLIPSE_STACK, 0}
	instructionSet[instrPOPJ] = instrChars{"POPJ", 0x9fc8, 0xffff, 1, UNIQUE_1_WORD_FMT, ECLIPSE_STACK, 0}
	instructionSet[instrPRTSEL] = instrChars{"PRTSEL", 0x783f, 0xffff, 1, UNIQUE_1_WORD_FMT, EAGLE_IO, 0}
	instructionSet[instrPSH] = instrChars{"PSH", 0x8648, 0x87ff, 1, TWOACC_1_WORD_FMT, ECLIPSE_STACK, 0}
	instructionSet[instrPSHJ] = instrChars{"PSHJ", 0x84b8, 0xfcff, 2, NOACC_MODE_IND_2_WORD_E_FMT, ECLIPSE_STACK, 1}
	instructionSet[instrPSHR] = instrChars{"PSHR", 0x87c8, 0xffff, 1, UNIQUE_1_WORD_FMT, ECLIPSE_STACK, 0}
	instructionSet[instrREADS] = instrChars{"READS", 0x613f, 0xe7ff, 1, ONEACC_1_WORD_FMT, EAGLE_IO, 0}
	instructionSet[instrRSTR] = instrChars{"RSTR", 0xefc8, 0xffff, 1, UNIQUE_1_WORD_FMT, ECLIPSE_STACK, 0}
	instructionSet[instrRTN] = instrChars{"RTN", 0xafc8, 0xffff, 1, UNIQUE_1_WORD_FMT, ECLIPSE_STACK, 0}
	instructionSet[instrSAVE] = instrChars{"SAVE", 0xe7c8, 0xffff, 2, UNIQUE_2_WORD_FMT, ECLIPSE_STACK, 0}
	instructionSet[instrSBI] = instrChars{"SBI", 0x8048, 0x87ff, 1, IMM_ONEACC_FMT, ECLIPSE_OP, 0}
	instructionSet[instrSEX] = instrChars{"SEX", 0x8349, 0x87ff, 1, TWOACC_1_WORD_FMT, EAGLE_OP, 0}
	instructionSet[instrSGE] = instrChars{"SGE", 0x8248, 0x87ff, 1, TWOACC_1_WORD_FMT, ECLIPSE_PC, 0}
	instructionSet[instrSGT] = instrChars{"SGT", 0x8208, 0x87ff, 1, TWOACC_1_WORD_FMT, ECLIPSE_PC, 0}
	instructionSet[instrSKP] = instrChars{"SKP", 0x6700, 0xff00, 1, IO_TEST_DEV_FMT, NOVA_IO, 0}
	instructionSet[instrSNB] = instrChars{"SNB", 0x85f8, 0x87ff, 1, TWOACC_1_WORD_FMT, ECLIPSE_PC, 0}
	instructionSet[instrSNOVR] = instrChars{"SNOVR", 0xa7b9, 0xffff, 1, UNIQUE_1_WORD_FMT, EAGLE_PC, 0}
	instructionSet[instrSPSR] = instrChars{"SPSR", 0xa7a9, 0xffff, 1, UNIQUE_1_WORD_FMT, EAGLE_OP, 0}
	instructionSet[instrSPTE] = instrChars{"SPTE", 0xe729, 0xffff, 1, UNIQUE_1_WORD_FMT, EAGLE_OP, 0}
	instructionSet[instrSSPT] = instrChars{"SSPT", 0xe7d9, 0xffff, 1, UNIQUE_1_WORD_FMT, EAGLE_OP, 0}
	instructionSet[instrSTA] = instrChars{"STA", 0x4000, 0xe000, 1, NOVA_ONEACC_EFF_ADDR_FMT, NOVA_MEMREF, 0}
	instructionSet[instrSTAFP] = instrChars{"STAFP", 0xc679, 0xe7ff, 1, ONEACC_1_WORD_FMT, EAGLE_STACK, 0}
	instructionSet[instrSTASB] = instrChars{"STASB", 0xc659, 0xe7ff, 1, ONEACC_1_WORD_FMT, EAGLE_STACK, 0}
	instructionSet[instrSTASL] = instrChars{"STASL", 0xa679, 0xe7ff, 1, ONEACC_1_WORD_FMT, EAGLE_STACK, 0}
	instructionSet[instrSTASP] = instrChars{"STASP", 0xa659, 0xe7ff, 1, ONEACC_1_WORD_FMT, EAGLE_STACK, 0}
	instructionSet[instrSTATS] = instrChars{"STATS", 0x8659, 0xe7ff, 1, ONEACC_1_WORD_FMT, EAGLE_STACK, 0}
	instructionSet[instrSTB] = instrChars{"STB", 0x8608, 0x87ff, 1, TWOACC_1_WORD_FMT, ECLIPSE_MEMREF, 0}
	instructionSet[instrSUB] = instrChars{"SUB", 0x8500, 0x8700, 1, NOVA_TWOACC_MULT_OP_FMT, NOVA_OP, 0}
	instructionSet[instrSZB] = instrChars{"SZB", 0x8488, 0x87ff, 1, TWOACC_1_WORD_FMT, ECLIPSE_PC, 0}
	instructionSet[instrSZBO] = instrChars{"SZBO", 0x84c8, 0x87ff, 1, TWOACC_1_WORD_FMT, ECLIPSE_PC, 0}
	instructionSet[instrWADC] = instrChars{"WADC", 0x8249, 0x87ff, 1, TWOACC_1_WORD_FMT, EAGLE_OP, 0}
	instructionSet[instrWADD] = instrChars{"WADD", 0x8149, 0x87ff, 1, TWOACC_1_WORD_FMT, EAGLE_OP, 0}
	instructionSet[instrWADDI] = instrChars{"WADDI", 0x8689, 0xe7ff, 3, ONEACC_IMM_3_WORD_FMT, EAGLE_OP, 0}
	instructionSet[instrWADI] = instrChars{"WADI", 0x84b9, 0x87ff, 1, IMM_ONEACC_FMT, EAGLE_OP, 0}
	instructionSet[instrWANC] = instrChars{"WANC", 0x8549, 0x87ff, 1, TWOACC_1_WORD_FMT, EAGLE_OP, 0}
	instructionSet[instrWAND] = instrChars{"WAND", 0x8449, 0x87ff, 1, TWOACC_1_WORD_FMT, EAGLE_OP, 0}
	instructionSet[instrWANDI] = instrChars{"WANDI", 0x8699, 0xe7ff, 3, ONEACC_IMMDWD_3_WORD_FMT, EAGLE_OP, 0}
	instructionSet[instrWASH] = instrChars{"WASH", 0x8279, 0x87ff, 1, TWOACC_1_WORD_FMT, EAGLE_OP, 0}
	instructionSet[instrWASHI] = instrChars{"WASHI", 0xc6a9, 0xe7ff, 2, ONEACC_IMM_2_WORD_FMT, EAGLE_OP, 0}
	instructionSet[instrWBLM] = instrChars{"WBLM", 0xe749, 0xffff, 1, UNIQUE_1_WORD_FMT, EAGLE_MEMREF, 0}
	instructionSet[instrWBR] = instrChars{"WBR", 0x8038, 0x843f, 1, SPLIT_8BIT_DISP_FMT, EAGLE_PC, 0}
	instructionSet[instrWBTO] = instrChars{"WBTO", 0x8299, 0x87ff, 1, TWOACC_1_WORD_FMT, EAGLE_MEMREF, 0}
	instructionSet[instrWBTZ] = instrChars{"WBTZ", 0x82a9, 0x87ff, 1, TWOACC_1_WORD_FMT, EAGLE_MEMREF, 0}
	instructionSet[instrWCLM] = instrChars{"WCLM", 0x8569, 0x87ff, 1, TWOACC_1_WORD_FMT, EAGLE_PC, 0}
	instructionSet[instrWCMP] = instrChars{"WCMP", 0xa759, 0xffff, 1, UNIQUE_1_WORD_FMT, EAGLE_MEMREF, 0}
	instructionSet[instrWCMV] = instrChars{"WCMV", 0x8779, 0xFFFF, 1, UNIQUE_1_WORD_FMT, EAGLE_MEMREF, 0}
	instructionSet[instrWCOM] = instrChars{"WCOM", 0x8459, 0x87ff, 1, TWOACC_1_WORD_FMT, EAGLE_OP, 0}
	instructionSet[instrWCST] = instrChars{"WCST", 0xe709, 0xffff, 1, UNIQUE_1_WORD_FMT, EAGLE_MEMREF, 0}
	instructionSet[instrWCTR] = instrChars{"WCTR", 0x8769, 0xffff, 1, UNIQUE_1_WORD_FMT, EAGLE_MEMREF, 0}
	instructionSet[instrWDecOp] = instrChars{"WDecOp", 0x8719, 0xffff, 2, WIDE_DEC_SPECIAL_FMT, EAGLE_DECIMAL, 1}
	instructionSet[instrWDIV] = instrChars{"WDIV", 0x8179, 0x87ff, 1, TWOACC_1_WORD_FMT, EAGLE_OP, 0}
	instructionSet[instrWDIVS] = instrChars{"WDIVS", 0xe769, 0xffff, 1, UNIQUE_1_WORD_FMT, EAGLE_OP, 0}
	instructionSet[instrWFLAD] = instrChars{"WFLAD", 0x84a9, 0x87ff, 1, TWOACC_1_WORD_FMT, EAGLE_FPU, 0}
	instructionSet[instrWFPOP] = instrChars{"WFPOP", 0xa789, 0xffff, 1, UNIQUE_1_WORD_FMT, EAGLE_STACK, 0}
	instructionSet[instrWFPSH] = instrChars{"WFPSH", 0x87b9, 0xffff, 1, UNIQUE_1_WORD_FMT, EAGLE_STACK, 0}
	instructionSet[instrWHLV] = instrChars{"WHLV", 0xe659, 0xe7ff, 1, ONEACC_1_WORD_FMT, EAGLE_OP, 0}
	instructionSet[instrWINC] = instrChars{"WINC", 0x8259, 0x87ff, 1, TWOACC_1_WORD_FMT, EAGLE_OP, 0}
	instructionSet[instrWIOR] = instrChars{"WIOR", 0x8469, 0x87ff, 1, TWOACC_1_WORD_FMT, EAGLE_OP, 0}
	instructionSet[instrWIORI] = instrChars{"WIORI", 0x86a9, 0xe7ff, 3, ONEACC_IMMDWD_3_WORD_FMT, EAGLE_OP, 0}
	instructionSet[instrWLDAI] = instrChars{"WLDAI", 0xc689, 0xe7ff, 3, ONEACC_IMMDWD_3_WORD_FMT, EAGLE_OP, 0}
	instructionSet[instrWLDB] = instrChars{"WLDB", 0x8529, 0x87ff, 1, TWOACC_1_WORD_FMT, EAGLE_MEMREF, 0}
	instructionSet[instrWLMP] = instrChars{"WLMP", 0xa7f9, 0xffff, 1, UNIQUE_1_WORD_FMT, EAGLE_IO, 0}
	instructionSet[instrWLSH] = instrChars{"WLSH", 0x8559, 0x87ff, 1, TWOACC_1_WORD_FMT, EAGLE_OP, 0}
	instructionSet[instrWLSHI] = instrChars{"WLSHI", 0xe6d9, 0xe7ff, 2, ONEACC_IMM_2_WORD_FMT, EAGLE_OP, 0}
	instructionSet[instrWLSI] = instrChars{"WLSI", 0x85b9, 0x87ff, 1, IMM_ONEACC_FMT, EAGLE_OP, 0}
	instructionSet[instrWMESS] = instrChars{"WMESS", 0xe719, 0xffff, 1, UNIQUE_1_WORD_FMT, EAGLE_PC, 0}
	instructionSet[instrWMOV] = instrChars{"WMOV", 0x8379, 0x87ff, 1, TWOACC_1_WORD_FMT, EAGLE_OP, 0}
	instructionSet[instrWMOVR] = instrChars{"WMOVR", 0xe699, 0xe7ff, 1, ONEACC_1_WORD_FMT, EAGLE_OP, 0}
	instructionSet[instrWMSP] = instrChars{"WMSP", 0xe649, 0xe7ff, 1, ONEACC_1_WORD_FMT, EAGLE_STACK, 0}
	instructionSet[instrWMUL] = instrChars{"WMUL", 0x8169, 0x87ff, 1, TWOACC_1_WORD_FMT, EAGLE_OP, 0}
	instructionSet[instrWMULS] = instrChars{"WMULS", 0xe759, 0xffff, 1, UNIQUE_1_WORD_FMT, EAGLE_OP, 0}
	instructionSet[instrWNADI] = instrChars{"WNADI", 0xe6f9, 0xe7ff, 2, ONEACC_IMM_2_WORD_FMT, EAGLE_OP, 0}
	instructionSet[instrWNEG] = instrChars{"WNEG", 0x8269, 0x87ff, 1, TWOACC_1_WORD_FMT, EAGLE_OP, 0}
	instructionSet[instrWPOP] = instrChars{"WPOP", 0x8089, 0x87ff, 1, TWOACC_1_WORD_FMT, EAGLE_STACK, 0}
	instructionSet[instrWPOPB] = instrChars{"WPOPB", 0xe779, 0xFFFF, 1, UNIQUE_1_WORD_FMT, EAGLE_PC, 0}
	instructionSet[instrWPOPJ] = instrChars{"WPOPJ", 0x8789, 0xFFFF, 1, UNIQUE_1_WORD_FMT, EAGLE_PC, 0}
	instructionSet[instrWPSH] = instrChars{"WPSH", 0x8579, 0x87ff, 1, TWOACC_1_WORD_FMT, EAGLE_STACK, 0}
	instructionSet[instrWRTN] = instrChars{"WRTN", 0x87a9, 0xffff, 1, UNIQUE_1_WORD_FMT, EAGLE_PC, 0}
	instructionSet[instrWSAVR] = instrChars{"WSAVR", 0xA729, 0xFFFF, 2, UNIQUE_2_WORD_FMT, EAGLE_STACK, 0}
	instructionSet[instrWSAVS] = instrChars{"WSAVS", 0xA739, 0xFFFF, 2, UNIQUE_2_WORD_FMT, EAGLE_STACK, 0}
	instructionSet[instrWSBI] = instrChars{"WSBI", 0x8589, 0x87ff, 1, IMM_ONEACC_FMT, EAGLE_OP, 0}
	instructionSet[instrWSEQ] = instrChars{"WSEQ", 0x80b9, 0x87ff, 1, TWOACC_1_WORD_FMT, EAGLE_PC, 0}
	instructionSet[instrWSEQI] = instrChars{"WSEQI", 0xe6c9, 0xe7ff, 2, ONEACC_IMM_2_WORD_FMT, EAGLE_PC, 0}
	instructionSet[instrWSGE] = instrChars{"WSGE", 0x8199, 0x87ff, 1, TWOACC_1_WORD_FMT, EAGLE_PC, 0}
	instructionSet[instrWSGT] = instrChars{"WSGT", 0x81b9, 0x87ff, 1, TWOACC_1_WORD_FMT, EAGLE_PC, 0}
	instructionSet[instrWSGTI] = instrChars{"WSGTI", 0xe689, 0xe7ff, 2, ONEACC_IMM_2_WORD_FMT, EAGLE_PC, 0}
	instructionSet[instrWSKBO] = instrChars{"WSKBO", 0x8f49, 0x8fcf, 1, WSKB_FMT, EAGLE_PC, 0}
	instructionSet[instrWSKBZ] = instrChars{"WSKBZ", 0x8f89, 0x8fcf, 1, WSKB_FMT, EAGLE_PC, 0}
	instructionSet[instrWSLE] = instrChars{"WSLE", 0x81a9, 0x87ff, 1, TWOACC_1_WORD_FMT, EAGLE_PC, 0}
	instructionSet[instrWSLEI] = instrChars{"WSLEI", 0xe6a9, 0xe7ff, 2, ONEACC_IMM_2_WORD_FMT, EAGLE_PC, 0}
	instructionSet[instrWSLT] = instrChars{"WSLT", 0x8289, 0x87ff, 1, TWOACC_1_WORD_FMT, EAGLE_PC, 0}
	instructionSet[instrWSNB] = instrChars{"WSNB", 0x8389, 0x87ff, 1, TWOACC_1_WORD_FMT, EAGLE_PC, 0}
	instructionSet[instrWSNE] = instrChars{"WSNE", 0x8189, 0x87ff, 1, TWOACC_1_WORD_FMT, EAGLE_PC, 0}
	instructionSet[instrWSNEI] = instrChars{"WSNEI", 0xe6e9, 0xe7ff, 2, ONEACC_IMM_2_WORD_FMT, EAGLE_PC, 0}
	instructionSet[instrWSSVR] = instrChars{"WSSVR", 0x8729, 0xffff, 2, UNIQUE_2_WORD_FMT, EAGLE_STACK, 0}
	instructionSet[instrWSSVS] = instrChars{"WSSVS", 0x8739, 0xffff, 2, UNIQUE_2_WORD_FMT, EAGLE_STACK, 0}
	instructionSet[instrWSTB] = instrChars{"WSTB", 0x8539, 0x87ff, 1, TWOACC_1_WORD_FMT, EAGLE_MEMREF, 0}
	instructionSet[instrWSTI] = instrChars{"WSTI", 0xe6b9, 0xe7ff, 1, ONEACC_1_WORD_FMT, EAGLE_FPU, 0}
	instructionSet[instrWSUB] = instrChars{"WSUB", 0x8159, 0x87ff, 1, TWOACC_1_WORD_FMT, EAGLE_OP, 0}
	instructionSet[instrWSZB] = instrChars{"WSZB", 0x82b9, 0x87ff, 1, TWOACC_1_WORD_FMT, EAGLE_PC, 0}
	instructionSet[instrWSZBO] = instrChars{"WSZBO", 0x8399, 0x87ff, 1, TWOACC_1_WORD_FMT, EAGLE_PC, 0}
	instructionSet[instrWUSGT] = instrChars{"WUSGT", 0x80a9, 0x87ff, 1, TWOACC_1_WORD_FMT, EAGLE_PC, 0}
	instructionSet[instrWUGTI] = instrChars{"WUGTI", 0xc699, 0xe7ff, 3, ONEACC_IMM_3_WORD_FMT, EAGLE_PC, 0}
	instructionSet[instrWXCH] = instrChars{"WXCH", 0x8369, 0x87ff, 1, TWOACC_1_WORD_FMT, EAGLE_OP, 0}
	instructionSet[instrWXORI] = instrChars{"WXORI", 0x86b9, 0xe7ff, 3, ONEACC_IMM_3_WORD_FMT, EAGLE_OP, 0}
	instructionSet[instrXCALL] = instrChars{"XCALL", 0x8609, 0xe7ff, 3, NOACC_MODE_IND_3_WORD_XCALL_FMT, EAGLE_PC, 1}
	instructionSet[instrXCH] = instrChars{"XCH", 0x81c8, 0x87ff, 1, TWOACC_1_WORD_FMT, ECLIPSE_OP, 0}
	instructionSet[instrXCT] = instrChars{"XCT", 0xa6f8, 0xe7ff, 1, ONEACC_1_WORD_FMT, ECLIPSE_OP, 0}
	instructionSet[instrXFLDS] = instrChars{"XFLDS", 0x8209, 0x87ff, 2, ONEACC_MODE_IND_2_WORD_X_FMT, EAGLE_FPU, 0}
	instructionSet[instrXJMP] = instrChars{"XJMP", 0xc609, 0xe7ff, 2, NOACC_MODE_IND_2_WORD_X_FMT, EAGLE_PC, 1}
	instructionSet[instrXJSR] = instrChars{"XJSR", 0xc619, 0xe7ff, 2, NOACC_MODE_IND_2_WORD_X_FMT, EAGLE_PC, 1}
	instructionSet[instrXLDB] = instrChars{"XLDB", 0x8419, 0x87ff, 2, ONEACC_MODE_2_WORD_X_B_FMT, EAGLE_MEMREF, 1}
	instructionSet[instrXLEF] = instrChars{"XLEF", 0x8409, 0x87ff, 2, ONEACC_MODE_IND_2_WORD_X_FMT, EAGLE_MEMREF, 1}
	instructionSet[instrXLEFB] = instrChars{"XLEFB", 0x8439, 0x87ff, 2, ONEACC_MODE_2_WORD_X_B_FMT, EAGLE_MEMREF, 1}
	instructionSet[instrXNADD] = instrChars{"XNADD", 0x8018, 0x87ff, 2, ONEACC_MODE_IND_2_WORD_X_FMT, EAGLE_MEMREF, 1}
	instructionSet[instrXNADI] = instrChars{"XNADI", 0x8418, 0x87ff, 2, IMM_MODE_2_WORD_FMT, EAGLE_MEMREF, 1}
	instructionSet[instrXNDO] = instrChars{"XNDO", 0x8498, 0x87ff, 3, THREE_WORD_DO_FMT, EAGLE_PC, 1}
	instructionSet[instrXNDSZ] = instrChars{"XNDSZ", 0xa609, 0xe7ff, 2, NOACC_MODE_IND_2_WORD_X_FMT, EAGLE_PC, 1}
	instructionSet[instrXNISZ] = instrChars{"XNISZ", 0x8639, 0xe7ff, 2, NOACC_MODE_IND_2_WORD_X_FMT, EAGLE_PC, 1}
	instructionSet[instrXNLDA] = instrChars{"XNLDA", 0x8329, 0x87ff, 2, ONEACC_MODE_IND_2_WORD_X_FMT, EAGLE_MEMREF, 1}
	instructionSet[instrXNSBI] = instrChars{"XNSBI", 0x8458, 0x87ff, 2, IMM_MODE_2_WORD_FMT, EAGLE_MEMREF, 1}
	instructionSet[instrXNSTA] = instrChars{"XNSTA", 0x8339, 0x87ff, 2, ONEACC_MODE_IND_2_WORD_X_FMT, EAGLE_MEMREF, 1}
	instructionSet[instrXNSUB] = instrChars{"XNSUB", 0x8058, 0x87ff, 2, ONEACC_MODE_IND_2_WORD_X_FMT, EAGLE_MEMREF, 1}
	instructionSet[instrXOR] = instrChars{"XOR", 0x8148, 0x87ff, 1, TWOACC_1_WORD_FMT, ECLIPSE_OP, 0}
	instructionSet[instrXORI] = instrChars{"XORI", 0xa7f8, 0xe7ff, 2, ONEACC_IMM_2_WORD_FMT, ECLIPSE_OP, 0}
	instructionSet[instrXPEF] = instrChars{"XPEF", 0x8629, 0xe7ff, 2, NOACC_MODE_IND_2_WORD_X_FMT, EAGLE_STACK, 1}
	instructionSet[instrXPEFB] = instrChars{"XPEFB", 0xa629, 0xe7ff, 2, NOACC_MODE_2_WORD_FMT, EAGLE_STACK, 1}
	instructionSet[instrXPSHJ] = instrChars{"XPSHJ", 0x8619, 0xe7ff, 2, IMM_MODE_2_WORD_FMT, EAGLE_STACK, 1}
	instructionSet[instrXSTB] = instrChars{"XSTB", 0x8429, 0x87ff, 2, ONEACC_MODE_2_WORD_X_B_FMT, EAGLE_MEMREF, 1}
	instructionSet[instrXWADD] = instrChars{"XWADD", 0x8118, 0x87ff, 2, ONEACC_MODE_IND_2_WORD_X_FMT, EAGLE_MEMREF, 1}
	instructionSet[instrXWADI] = instrChars{"XWADI", 0x8518, 0x87ff, 2, IMM_MODE_2_WORD_FMT, EAGLE_MEMREF, 1}
	instructionSet[instrXWDSZ] = instrChars{"XWDSZ", 0xA639, 0xe7FF, 2, NOACC_MODE_IND_2_WORD_X_FMT, EAGLE_PC, 1}
	instructionSet[instrXWISZ] = instrChars{"XWISZ", 0xa619, 0xe7ff, 2, NOACC_MODE_IND_2_WORD_X_FMT, EAGLE_PC, 1}
	instructionSet[instrXWLDA] = instrChars{"XWLDA", 0x8309, 0x87ff, 2, ONEACC_MODE_IND_2_WORD_X_FMT, EAGLE_MEMREF, 1}
	instructionSet[instrXWMUL] = instrChars{"XWMUL", 0x8198, 0x87ff, 2, ONEACC_MODE_IND_2_WORD_X_FMT, EAGLE_MEMREF, 1}
	instructionSet[instrXWSBI] = instrChars{"XWSBI", 0x8558, 0x87ff, 2, IMM_MODE_2_WORD_FMT, EAGLE_MEMREF, 1}
	instructionSet[instrXWSTA] = instrChars{"XWSTA", 0x8319, 0x87ff, 2, ONEACC_MODE_IND_2_WORD_X_FMT, EAGLE_MEMREF, 1}
	instructionSet[instrXWSUB] = instrChars{"XWSUB", 0x8158, 0x87ff, 2, ONEACC_MODE_IND_2_WORD_X_FMT, EAGLE_MEMREF, 1}
	instructionSet[instrZEX] = instrChars{"ZEX", 0x8359, 0x87ff, 1, TWOACC_1_WORD_FMT, EAGLE_OP, 0}
}
