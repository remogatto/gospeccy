/*

Copyright (c) 2010 Andrea Fazzi

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
"Software"), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

*/

//
// Automatically generated file -- DO NOT EDIT
//

package spectrum

func initOpcodes() {
	// BEGIN of non shifted opcodes
	/* NOP */
	opcodesMap[0x00] = instr__NOP
	/* LD BC,nnnn */
	opcodesMap[0x01] = instr__LD_BC_NNNN
	/* LD (BC),A */
	opcodesMap[0x02] = instr__LD_iBC_A
	/* INC BC */
	opcodesMap[0x03] = instr__INC_BC
	/* INC B */
	opcodesMap[0x04] = instr__INC_B
	/* DEC B */
	opcodesMap[0x05] = instr__DEC_B
	/* LD B,nn */
	opcodesMap[0x06] = instr__LD_B_NN
	/* RLCA */
	opcodesMap[0x07] = instr__RLCA
	/* EX AF,AF' */
	opcodesMap[0x08] = instr__EX_AF_AF
	/* ADD HL,BC */
	opcodesMap[0x09] = instr__ADD_HL_BC
	/* LD A,(BC) */
	opcodesMap[0x0a] = instr__LD_A_iBC
	/* DEC BC */
	opcodesMap[0x0b] = instr__DEC_BC
	/* INC C */
	opcodesMap[0x0c] = instr__INC_C
	/* DEC C */
	opcodesMap[0x0d] = instr__DEC_C
	/* LD C,nn */
	opcodesMap[0x0e] = instr__LD_C_NN
	/* RRCA */
	opcodesMap[0x0f] = instr__RRCA
	/* DJNZ offset */
	opcodesMap[0x10] = instr__DJNZ_OFFSET
	/* LD DE,nnnn */
	opcodesMap[0x11] = instr__LD_DE_NNNN
	/* LD (DE),A */
	opcodesMap[0x12] = instr__LD_iDE_A
	/* INC DE */
	opcodesMap[0x13] = instr__INC_DE
	/* INC D */
	opcodesMap[0x14] = instr__INC_D
	/* DEC D */
	opcodesMap[0x15] = instr__DEC_D
	/* LD D,nn */
	opcodesMap[0x16] = instr__LD_D_NN
	/* RLA */
	opcodesMap[0x17] = instr__RLA
	/* JR offset */
	opcodesMap[0x18] = instr__JR_OFFSET
	/* ADD HL,DE */
	opcodesMap[0x19] = instr__ADD_HL_DE
	/* LD A,(DE) */
	opcodesMap[0x1a] = instr__LD_A_iDE
	/* DEC DE */
	opcodesMap[0x1b] = instr__DEC_DE
	/* INC E */
	opcodesMap[0x1c] = instr__INC_E
	/* DEC E */
	opcodesMap[0x1d] = instr__DEC_E
	/* LD E,nn */
	opcodesMap[0x1e] = instr__LD_E_NN
	/* RRA */
	opcodesMap[0x1f] = instr__RRA
	/* JR NZ,offset */
	opcodesMap[0x20] = instr__JR_NZ_OFFSET
	/* LD HL,nnnn */
	opcodesMap[0x21] = instr__LD_HL_NNNN
	/* LD (nnnn),HL */
	opcodesMap[0x22] = instr__LD_iNNNN_HL
	/* INC HL */
	opcodesMap[0x23] = instr__INC_HL
	/* INC H */
	opcodesMap[0x24] = instr__INC_H
	/* DEC H */
	opcodesMap[0x25] = instr__DEC_H
	/* LD H,nn */
	opcodesMap[0x26] = instr__LD_H_NN
	/* DAA */
	opcodesMap[0x27] = instr__DAA
	/* JR Z,offset */
	opcodesMap[0x28] = instr__JR_Z_OFFSET
	/* ADD HL,HL */
	opcodesMap[0x29] = instr__ADD_HL_HL
	/* LD HL,(nnnn) */
	opcodesMap[0x2a] = instr__LD_HL_iNNNN
	/* DEC HL */
	opcodesMap[0x2b] = instr__DEC_HL
	/* INC L */
	opcodesMap[0x2c] = instr__INC_L
	/* DEC L */
	opcodesMap[0x2d] = instr__DEC_L
	/* LD L,nn */
	opcodesMap[0x2e] = instr__LD_L_NN
	/* CPL */
	opcodesMap[0x2f] = instr__CPL
	/* JR NC,offset */
	opcodesMap[0x30] = instr__JR_NC_OFFSET
	/* LD SP,nnnn */
	opcodesMap[0x31] = instr__LD_SP_NNNN
	/* LD (nnnn),A */
	opcodesMap[0x32] = instr__LD_iNNNN_A
	/* INC SP */
	opcodesMap[0x33] = instr__INC_SP
	/* INC (HL) */
	opcodesMap[0x34] = instr__INC_iHL
	/* DEC (HL) */
	opcodesMap[0x35] = instr__DEC_iHL
	/* LD (HL),nn */
	opcodesMap[0x36] = instr__LD_iHL_NN
	/* SCF */
	opcodesMap[0x37] = instr__SCF
	/* JR C,offset */
	opcodesMap[0x38] = instr__JR_C_OFFSET
	/* ADD HL,SP */
	opcodesMap[0x39] = instr__ADD_HL_SP
	/* LD A,(nnnn) */
	opcodesMap[0x3a] = instr__LD_A_iNNNN
	/* DEC SP */
	opcodesMap[0x3b] = instr__DEC_SP
	/* INC A */
	opcodesMap[0x3c] = instr__INC_A
	/* DEC A */
	opcodesMap[0x3d] = instr__DEC_A
	/* LD A,nn */
	opcodesMap[0x3e] = instr__LD_A_NN
	/* CCF */
	opcodesMap[0x3f] = instr__CCF
	/* LD B,B */
	opcodesMap[0x40] = instr__LD_B_B
	/* LD B,C */
	opcodesMap[0x41] = instr__LD_B_C
	/* LD B,D */
	opcodesMap[0x42] = instr__LD_B_D
	/* LD B,E */
	opcodesMap[0x43] = instr__LD_B_E
	/* LD B,H */
	opcodesMap[0x44] = instr__LD_B_H
	/* LD B,L */
	opcodesMap[0x45] = instr__LD_B_L
	/* LD B,(HL) */
	opcodesMap[0x46] = instr__LD_B_iHL
	/* LD B,A */
	opcodesMap[0x47] = instr__LD_B_A
	/* LD C,B */
	opcodesMap[0x48] = instr__LD_C_B
	/* LD C,C */
	opcodesMap[0x49] = instr__LD_C_C
	/* LD C,D */
	opcodesMap[0x4a] = instr__LD_C_D
	/* LD C,E */
	opcodesMap[0x4b] = instr__LD_C_E
	/* LD C,H */
	opcodesMap[0x4c] = instr__LD_C_H
	/* LD C,L */
	opcodesMap[0x4d] = instr__LD_C_L
	/* LD C,(HL) */
	opcodesMap[0x4e] = instr__LD_C_iHL
	/* LD C,A */
	opcodesMap[0x4f] = instr__LD_C_A
	/* LD D,B */
	opcodesMap[0x50] = instr__LD_D_B
	/* LD D,C */
	opcodesMap[0x51] = instr__LD_D_C
	/* LD D,D */
	opcodesMap[0x52] = instr__LD_D_D
	/* LD D,E */
	opcodesMap[0x53] = instr__LD_D_E
	/* LD D,H */
	opcodesMap[0x54] = instr__LD_D_H
	/* LD D,L */
	opcodesMap[0x55] = instr__LD_D_L
	/* LD D,(HL) */
	opcodesMap[0x56] = instr__LD_D_iHL
	/* LD D,A */
	opcodesMap[0x57] = instr__LD_D_A
	/* LD E,B */
	opcodesMap[0x58] = instr__LD_E_B
	/* LD E,C */
	opcodesMap[0x59] = instr__LD_E_C
	/* LD E,D */
	opcodesMap[0x5a] = instr__LD_E_D
	/* LD E,E */
	opcodesMap[0x5b] = instr__LD_E_E
	/* LD E,H */
	opcodesMap[0x5c] = instr__LD_E_H
	/* LD E,L */
	opcodesMap[0x5d] = instr__LD_E_L
	/* LD E,(HL) */
	opcodesMap[0x5e] = instr__LD_E_iHL
	/* LD E,A */
	opcodesMap[0x5f] = instr__LD_E_A
	/* LD H,B */
	opcodesMap[0x60] = instr__LD_H_B
	/* LD H,C */
	opcodesMap[0x61] = instr__LD_H_C
	/* LD H,D */
	opcodesMap[0x62] = instr__LD_H_D
	/* LD H,E */
	opcodesMap[0x63] = instr__LD_H_E
	/* LD H,H */
	opcodesMap[0x64] = instr__LD_H_H
	/* LD H,L */
	opcodesMap[0x65] = instr__LD_H_L
	/* LD H,(HL) */
	opcodesMap[0x66] = instr__LD_H_iHL
	/* LD H,A */
	opcodesMap[0x67] = instr__LD_H_A
	/* LD L,B */
	opcodesMap[0x68] = instr__LD_L_B
	/* LD L,C */
	opcodesMap[0x69] = instr__LD_L_C
	/* LD L,D */
	opcodesMap[0x6a] = instr__LD_L_D
	/* LD L,E */
	opcodesMap[0x6b] = instr__LD_L_E
	/* LD L,H */
	opcodesMap[0x6c] = instr__LD_L_H
	/* LD L,L */
	opcodesMap[0x6d] = instr__LD_L_L
	/* LD L,(HL) */
	opcodesMap[0x6e] = instr__LD_L_iHL
	/* LD L,A */
	opcodesMap[0x6f] = instr__LD_L_A
	/* LD (HL),B */
	opcodesMap[0x70] = instr__LD_iHL_B
	/* LD (HL),C */
	opcodesMap[0x71] = instr__LD_iHL_C
	/* LD (HL),D */
	opcodesMap[0x72] = instr__LD_iHL_D
	/* LD (HL),E */
	opcodesMap[0x73] = instr__LD_iHL_E
	/* LD (HL),H */
	opcodesMap[0x74] = instr__LD_iHL_H
	/* LD (HL),L */
	opcodesMap[0x75] = instr__LD_iHL_L
	/* HALT */
	opcodesMap[0x76] = instr__HALT
	/* LD (HL),A */
	opcodesMap[0x77] = instr__LD_iHL_A
	/* LD A,B */
	opcodesMap[0x78] = instr__LD_A_B
	/* LD A,C */
	opcodesMap[0x79] = instr__LD_A_C
	/* LD A,D */
	opcodesMap[0x7a] = instr__LD_A_D
	/* LD A,E */
	opcodesMap[0x7b] = instr__LD_A_E
	/* LD A,H */
	opcodesMap[0x7c] = instr__LD_A_H
	/* LD A,L */
	opcodesMap[0x7d] = instr__LD_A_L
	/* LD A,(HL) */
	opcodesMap[0x7e] = instr__LD_A_iHL
	/* LD A,A */
	opcodesMap[0x7f] = instr__LD_A_A
	/* ADD A,B */
	opcodesMap[0x80] = instr__ADD_A_B
	/* ADD A,C */
	opcodesMap[0x81] = instr__ADD_A_C
	/* ADD A,D */
	opcodesMap[0x82] = instr__ADD_A_D
	/* ADD A,E */
	opcodesMap[0x83] = instr__ADD_A_E
	/* ADD A,H */
	opcodesMap[0x84] = instr__ADD_A_H
	/* ADD A,L */
	opcodesMap[0x85] = instr__ADD_A_L
	/* ADD A,(HL) */
	opcodesMap[0x86] = instr__ADD_A_iHL
	/* ADD A,A */
	opcodesMap[0x87] = instr__ADD_A_A
	/* ADC A,B */
	opcodesMap[0x88] = instr__ADC_A_B
	/* ADC A,C */
	opcodesMap[0x89] = instr__ADC_A_C
	/* ADC A,D */
	opcodesMap[0x8a] = instr__ADC_A_D
	/* ADC A,E */
	opcodesMap[0x8b] = instr__ADC_A_E
	/* ADC A,H */
	opcodesMap[0x8c] = instr__ADC_A_H
	/* ADC A,L */
	opcodesMap[0x8d] = instr__ADC_A_L
	/* ADC A,(HL) */
	opcodesMap[0x8e] = instr__ADC_A_iHL
	/* ADC A,A */
	opcodesMap[0x8f] = instr__ADC_A_A
	/* SUB A,B */
	opcodesMap[0x90] = instr__SUB_A_B
	/* SUB A,C */
	opcodesMap[0x91] = instr__SUB_A_C
	/* SUB A,D */
	opcodesMap[0x92] = instr__SUB_A_D
	/* SUB A,E */
	opcodesMap[0x93] = instr__SUB_A_E
	/* SUB A,H */
	opcodesMap[0x94] = instr__SUB_A_H
	/* SUB A,L */
	opcodesMap[0x95] = instr__SUB_A_L
	/* SUB A,(HL) */
	opcodesMap[0x96] = instr__SUB_A_iHL
	/* SUB A,A */
	opcodesMap[0x97] = instr__SUB_A_A
	/* SBC A,B */
	opcodesMap[0x98] = instr__SBC_A_B
	/* SBC A,C */
	opcodesMap[0x99] = instr__SBC_A_C
	/* SBC A,D */
	opcodesMap[0x9a] = instr__SBC_A_D
	/* SBC A,E */
	opcodesMap[0x9b] = instr__SBC_A_E
	/* SBC A,H */
	opcodesMap[0x9c] = instr__SBC_A_H
	/* SBC A,L */
	opcodesMap[0x9d] = instr__SBC_A_L
	/* SBC A,(HL) */
	opcodesMap[0x9e] = instr__SBC_A_iHL
	/* SBC A,A */
	opcodesMap[0x9f] = instr__SBC_A_A
	/* AND A,B */
	opcodesMap[0xa0] = instr__AND_A_B
	/* AND A,C */
	opcodesMap[0xa1] = instr__AND_A_C
	/* AND A,D */
	opcodesMap[0xa2] = instr__AND_A_D
	/* AND A,E */
	opcodesMap[0xa3] = instr__AND_A_E
	/* AND A,H */
	opcodesMap[0xa4] = instr__AND_A_H
	/* AND A,L */
	opcodesMap[0xa5] = instr__AND_A_L
	/* AND A,(HL) */
	opcodesMap[0xa6] = instr__AND_A_iHL
	/* AND A,A */
	opcodesMap[0xa7] = instr__AND_A_A
	/* XOR A,B */
	opcodesMap[0xa8] = instr__XOR_A_B
	/* XOR A,C */
	opcodesMap[0xa9] = instr__XOR_A_C
	/* XOR A,D */
	opcodesMap[0xaa] = instr__XOR_A_D
	/* XOR A,E */
	opcodesMap[0xab] = instr__XOR_A_E
	/* XOR A,H */
	opcodesMap[0xac] = instr__XOR_A_H
	/* XOR A,L */
	opcodesMap[0xad] = instr__XOR_A_L
	/* XOR A,(HL) */
	opcodesMap[0xae] = instr__XOR_A_iHL
	/* XOR A,A */
	opcodesMap[0xaf] = instr__XOR_A_A
	/* OR A,B */
	opcodesMap[0xb0] = instr__OR_A_B
	/* OR A,C */
	opcodesMap[0xb1] = instr__OR_A_C
	/* OR A,D */
	opcodesMap[0xb2] = instr__OR_A_D
	/* OR A,E */
	opcodesMap[0xb3] = instr__OR_A_E
	/* OR A,H */
	opcodesMap[0xb4] = instr__OR_A_H
	/* OR A,L */
	opcodesMap[0xb5] = instr__OR_A_L
	/* OR A,(HL) */
	opcodesMap[0xb6] = instr__OR_A_iHL
	/* OR A,A */
	opcodesMap[0xb7] = instr__OR_A_A
	/* CP B */
	opcodesMap[0xb8] = instr__CP_B
	/* CP C */
	opcodesMap[0xb9] = instr__CP_C
	/* CP D */
	opcodesMap[0xba] = instr__CP_D
	/* CP E */
	opcodesMap[0xbb] = instr__CP_E
	/* CP H */
	opcodesMap[0xbc] = instr__CP_H
	/* CP L */
	opcodesMap[0xbd] = instr__CP_L
	/* CP (HL) */
	opcodesMap[0xbe] = instr__CP_iHL
	/* CP A */
	opcodesMap[0xbf] = instr__CP_A
	/* RET NZ */
	opcodesMap[0xc0] = instr__RET_NZ
	/* POP BC */
	opcodesMap[0xc1] = instr__POP_BC
	/* JP NZ,nnnn */
	opcodesMap[0xc2] = instr__JP_NZ_NNNN
	/* JP nnnn */
	opcodesMap[0xc3] = instr__JP_NNNN
	/* CALL NZ,nnnn */
	opcodesMap[0xc4] = instr__CALL_NZ_NNNN
	/* PUSH BC */
	opcodesMap[0xc5] = instr__PUSH_BC
	/* ADD A,nn */
	opcodesMap[0xc6] = instr__ADD_A_NN
	/* RST 00 */
	opcodesMap[0xc7] = instr__RST_00
	/* RET Z */
	opcodesMap[0xc8] = instr__RET_Z
	/* RET */
	opcodesMap[0xc9] = instr__RET
	/* JP Z,nnnn */
	opcodesMap[0xca] = instr__JP_Z_NNNN
	/* shift CB */
	opcodesMap[0xcb] = instr__SHIFT_CB
	/* CALL Z,nnnn */
	opcodesMap[0xcc] = instr__CALL_Z_NNNN
	/* CALL nnnn */
	opcodesMap[0xcd] = instr__CALL_NNNN
	/* ADC A,nn */
	opcodesMap[0xce] = instr__ADC_A_NN
	/* RST 8 */
	opcodesMap[0xcf] = instr__RST_8
	/* RET NC */
	opcodesMap[0xd0] = instr__RET_NC
	/* POP DE */
	opcodesMap[0xd1] = instr__POP_DE
	/* JP NC,nnnn */
	opcodesMap[0xd2] = instr__JP_NC_NNNN
	/* OUT (nn),A */
	opcodesMap[0xd3] = instr__OUT_iNN_A
	/* CALL NC,nnnn */
	opcodesMap[0xd4] = instr__CALL_NC_NNNN
	/* PUSH DE */
	opcodesMap[0xd5] = instr__PUSH_DE
	/* SUB nn */
	opcodesMap[0xd6] = instr__SUB_NN
	/* RST 10 */
	opcodesMap[0xd7] = instr__RST_10
	/* RET C */
	opcodesMap[0xd8] = instr__RET_C
	/* EXX */
	opcodesMap[0xd9] = instr__EXX
	/* JP C,nnnn */
	opcodesMap[0xda] = instr__JP_C_NNNN
	/* IN A,(nn) */
	opcodesMap[0xdb] = instr__IN_A_iNN
	/* CALL C,nnnn */
	opcodesMap[0xdc] = instr__CALL_C_NNNN
	/* shift DD */
	opcodesMap[0xdd] = instr__SHIFT_DD
	/* SBC A,nn */
	opcodesMap[0xde] = instr__SBC_A_NN
	/* RST 18 */
	opcodesMap[0xdf] = instr__RST_18
	/* RET PO */
	opcodesMap[0xe0] = instr__RET_PO
	/* POP HL */
	opcodesMap[0xe1] = instr__POP_HL
	/* JP PO,nnnn */
	opcodesMap[0xe2] = instr__JP_PO_NNNN
	/* EX (SP),HL */
	opcodesMap[0xe3] = instr__EX_iSP_HL
	/* CALL PO,nnnn */
	opcodesMap[0xe4] = instr__CALL_PO_NNNN
	/* PUSH HL */
	opcodesMap[0xe5] = instr__PUSH_HL
	/* AND nn */
	opcodesMap[0xe6] = instr__AND_NN
	/* RST 20 */
	opcodesMap[0xe7] = instr__RST_20
	/* RET PE */
	opcodesMap[0xe8] = instr__RET_PE
	/* JP HL */
	opcodesMap[0xe9] = instr__JP_HL
	/* JP PE,nnnn */
	opcodesMap[0xea] = instr__JP_PE_NNNN
	/* EX DE,HL */
	opcodesMap[0xeb] = instr__EX_DE_HL
	/* CALL PE,nnnn */
	opcodesMap[0xec] = instr__CALL_PE_NNNN
	/* shift ED */
	opcodesMap[0xed] = instr__SHIFT_ED
	/* XOR A,nn */
	opcodesMap[0xee] = instr__XOR_A_NN
	/* RST 28 */
	opcodesMap[0xef] = instr__RST_28
	/* RET P */
	opcodesMap[0xf0] = instr__RET_P
	/* POP AF */
	opcodesMap[0xf1] = instr__POP_AF
	/* JP P,nnnn */
	opcodesMap[0xf2] = instr__JP_P_NNNN
	/* DI */
	opcodesMap[0xf3] = instr__DI
	/* CALL P,nnnn */
	opcodesMap[0xf4] = instr__CALL_P_NNNN
	/* PUSH AF */
	opcodesMap[0xf5] = instr__PUSH_AF
	/* OR nn */
	opcodesMap[0xf6] = instr__OR_NN
	/* RST 30 */
	opcodesMap[0xf7] = instr__RST_30
	/* RET M */
	opcodesMap[0xf8] = instr__RET_M
	/* LD SP,HL */
	opcodesMap[0xf9] = instr__LD_SP_HL
	/* JP M,nnnn */
	opcodesMap[0xfa] = instr__JP_M_NNNN
	/* EI */
	opcodesMap[0xfb] = instr__EI
	/* CALL M,nnnn */
	opcodesMap[0xfc] = instr__CALL_M_NNNN
	/* shift FD */
	opcodesMap[0xfd] = instr__SHIFT_FD
	/* CP nn */
	opcodesMap[0xfe] = instr__CP_NN
	/* RST 38 */
	opcodesMap[0xff] = instr__RST_38

	// END of non shifted opcodes

	// BEGIN of 0xcb shifted opcodes
	/* RLC B */
	opcodesMap[SHIFT_0xCB+0x00] = instrCB__RLC_B
	/* RLC C */
	opcodesMap[SHIFT_0xCB+0x01] = instrCB__RLC_C
	/* RLC D */
	opcodesMap[SHIFT_0xCB+0x02] = instrCB__RLC_D
	/* RLC E */
	opcodesMap[SHIFT_0xCB+0x03] = instrCB__RLC_E
	/* RLC H */
	opcodesMap[SHIFT_0xCB+0x04] = instrCB__RLC_H
	/* RLC L */
	opcodesMap[SHIFT_0xCB+0x05] = instrCB__RLC_L
	/* RLC (HL) */
	opcodesMap[SHIFT_0xCB+0x06] = instrCB__RLC_iHL
	/* RLC A */
	opcodesMap[SHIFT_0xCB+0x07] = instrCB__RLC_A
	/* RRC B */
	opcodesMap[SHIFT_0xCB+0x08] = instrCB__RRC_B
	/* RRC C */
	opcodesMap[SHIFT_0xCB+0x09] = instrCB__RRC_C
	/* RRC D */
	opcodesMap[SHIFT_0xCB+0x0a] = instrCB__RRC_D
	/* RRC E */
	opcodesMap[SHIFT_0xCB+0x0b] = instrCB__RRC_E
	/* RRC H */
	opcodesMap[SHIFT_0xCB+0x0c] = instrCB__RRC_H
	/* RRC L */
	opcodesMap[SHIFT_0xCB+0x0d] = instrCB__RRC_L
	/* RRC (HL) */
	opcodesMap[SHIFT_0xCB+0x0e] = instrCB__RRC_iHL
	/* RRC A */
	opcodesMap[SHIFT_0xCB+0x0f] = instrCB__RRC_A
	/* RL B */
	opcodesMap[SHIFT_0xCB+0x10] = instrCB__RL_B
	/* RL C */
	opcodesMap[SHIFT_0xCB+0x11] = instrCB__RL_C
	/* RL D */
	opcodesMap[SHIFT_0xCB+0x12] = instrCB__RL_D
	/* RL E */
	opcodesMap[SHIFT_0xCB+0x13] = instrCB__RL_E
	/* RL H */
	opcodesMap[SHIFT_0xCB+0x14] = instrCB__RL_H
	/* RL L */
	opcodesMap[SHIFT_0xCB+0x15] = instrCB__RL_L
	/* RL (HL) */
	opcodesMap[SHIFT_0xCB+0x16] = instrCB__RL_iHL
	/* RL A */
	opcodesMap[SHIFT_0xCB+0x17] = instrCB__RL_A
	/* RR B */
	opcodesMap[SHIFT_0xCB+0x18] = instrCB__RR_B
	/* RR C */
	opcodesMap[SHIFT_0xCB+0x19] = instrCB__RR_C
	/* RR D */
	opcodesMap[SHIFT_0xCB+0x1a] = instrCB__RR_D
	/* RR E */
	opcodesMap[SHIFT_0xCB+0x1b] = instrCB__RR_E
	/* RR H */
	opcodesMap[SHIFT_0xCB+0x1c] = instrCB__RR_H
	/* RR L */
	opcodesMap[SHIFT_0xCB+0x1d] = instrCB__RR_L
	/* RR (HL) */
	opcodesMap[SHIFT_0xCB+0x1e] = instrCB__RR_iHL
	/* RR A */
	opcodesMap[SHIFT_0xCB+0x1f] = instrCB__RR_A
	/* SLA B */
	opcodesMap[SHIFT_0xCB+0x20] = instrCB__SLA_B
	/* SLA C */
	opcodesMap[SHIFT_0xCB+0x21] = instrCB__SLA_C
	/* SLA D */
	opcodesMap[SHIFT_0xCB+0x22] = instrCB__SLA_D
	/* SLA E */
	opcodesMap[SHIFT_0xCB+0x23] = instrCB__SLA_E
	/* SLA H */
	opcodesMap[SHIFT_0xCB+0x24] = instrCB__SLA_H
	/* SLA L */
	opcodesMap[SHIFT_0xCB+0x25] = instrCB__SLA_L
	/* SLA (HL) */
	opcodesMap[SHIFT_0xCB+0x26] = instrCB__SLA_iHL
	/* SLA A */
	opcodesMap[SHIFT_0xCB+0x27] = instrCB__SLA_A
	/* SRA B */
	opcodesMap[SHIFT_0xCB+0x28] = instrCB__SRA_B
	/* SRA C */
	opcodesMap[SHIFT_0xCB+0x29] = instrCB__SRA_C
	/* SRA D */
	opcodesMap[SHIFT_0xCB+0x2a] = instrCB__SRA_D
	/* SRA E */
	opcodesMap[SHIFT_0xCB+0x2b] = instrCB__SRA_E
	/* SRA H */
	opcodesMap[SHIFT_0xCB+0x2c] = instrCB__SRA_H
	/* SRA L */
	opcodesMap[SHIFT_0xCB+0x2d] = instrCB__SRA_L
	/* SRA (HL) */
	opcodesMap[SHIFT_0xCB+0x2e] = instrCB__SRA_iHL
	/* SRA A */
	opcodesMap[SHIFT_0xCB+0x2f] = instrCB__SRA_A
	/* SLL B */
	opcodesMap[SHIFT_0xCB+0x30] = instrCB__SLL_B
	/* SLL C */
	opcodesMap[SHIFT_0xCB+0x31] = instrCB__SLL_C
	/* SLL D */
	opcodesMap[SHIFT_0xCB+0x32] = instrCB__SLL_D
	/* SLL E */
	opcodesMap[SHIFT_0xCB+0x33] = instrCB__SLL_E
	/* SLL H */
	opcodesMap[SHIFT_0xCB+0x34] = instrCB__SLL_H
	/* SLL L */
	opcodesMap[SHIFT_0xCB+0x35] = instrCB__SLL_L
	/* SLL (HL) */
	opcodesMap[SHIFT_0xCB+0x36] = instrCB__SLL_iHL
	/* SLL A */
	opcodesMap[SHIFT_0xCB+0x37] = instrCB__SLL_A
	/* SRL B */
	opcodesMap[SHIFT_0xCB+0x38] = instrCB__SRL_B
	/* SRL C */
	opcodesMap[SHIFT_0xCB+0x39] = instrCB__SRL_C
	/* SRL D */
	opcodesMap[SHIFT_0xCB+0x3a] = instrCB__SRL_D
	/* SRL E */
	opcodesMap[SHIFT_0xCB+0x3b] = instrCB__SRL_E
	/* SRL H */
	opcodesMap[SHIFT_0xCB+0x3c] = instrCB__SRL_H
	/* SRL L */
	opcodesMap[SHIFT_0xCB+0x3d] = instrCB__SRL_L
	/* SRL (HL) */
	opcodesMap[SHIFT_0xCB+0x3e] = instrCB__SRL_iHL
	/* SRL A */
	opcodesMap[SHIFT_0xCB+0x3f] = instrCB__SRL_A
	/* BIT 0,B */
	opcodesMap[SHIFT_0xCB+0x40] = instrCB__BIT_0_B
	/* BIT 0,C */
	opcodesMap[SHIFT_0xCB+0x41] = instrCB__BIT_0_C
	/* BIT 0,D */
	opcodesMap[SHIFT_0xCB+0x42] = instrCB__BIT_0_D
	/* BIT 0,E */
	opcodesMap[SHIFT_0xCB+0x43] = instrCB__BIT_0_E
	/* BIT 0,H */
	opcodesMap[SHIFT_0xCB+0x44] = instrCB__BIT_0_H
	/* BIT 0,L */
	opcodesMap[SHIFT_0xCB+0x45] = instrCB__BIT_0_L
	/* BIT 0,(HL) */
	opcodesMap[SHIFT_0xCB+0x46] = instrCB__BIT_0_iHL
	/* BIT 0,A */
	opcodesMap[SHIFT_0xCB+0x47] = instrCB__BIT_0_A
	/* BIT 1,B */
	opcodesMap[SHIFT_0xCB+0x48] = instrCB__BIT_1_B
	/* BIT 1,C */
	opcodesMap[SHIFT_0xCB+0x49] = instrCB__BIT_1_C
	/* BIT 1,D */
	opcodesMap[SHIFT_0xCB+0x4a] = instrCB__BIT_1_D
	/* BIT 1,E */
	opcodesMap[SHIFT_0xCB+0x4b] = instrCB__BIT_1_E
	/* BIT 1,H */
	opcodesMap[SHIFT_0xCB+0x4c] = instrCB__BIT_1_H
	/* BIT 1,L */
	opcodesMap[SHIFT_0xCB+0x4d] = instrCB__BIT_1_L
	/* BIT 1,(HL) */
	opcodesMap[SHIFT_0xCB+0x4e] = instrCB__BIT_1_iHL
	/* BIT 1,A */
	opcodesMap[SHIFT_0xCB+0x4f] = instrCB__BIT_1_A
	/* BIT 2,B */
	opcodesMap[SHIFT_0xCB+0x50] = instrCB__BIT_2_B
	/* BIT 2,C */
	opcodesMap[SHIFT_0xCB+0x51] = instrCB__BIT_2_C
	/* BIT 2,D */
	opcodesMap[SHIFT_0xCB+0x52] = instrCB__BIT_2_D
	/* BIT 2,E */
	opcodesMap[SHIFT_0xCB+0x53] = instrCB__BIT_2_E
	/* BIT 2,H */
	opcodesMap[SHIFT_0xCB+0x54] = instrCB__BIT_2_H
	/* BIT 2,L */
	opcodesMap[SHIFT_0xCB+0x55] = instrCB__BIT_2_L
	/* BIT 2,(HL) */
	opcodesMap[SHIFT_0xCB+0x56] = instrCB__BIT_2_iHL
	/* BIT 2,A */
	opcodesMap[SHIFT_0xCB+0x57] = instrCB__BIT_2_A
	/* BIT 3,B */
	opcodesMap[SHIFT_0xCB+0x58] = instrCB__BIT_3_B
	/* BIT 3,C */
	opcodesMap[SHIFT_0xCB+0x59] = instrCB__BIT_3_C
	/* BIT 3,D */
	opcodesMap[SHIFT_0xCB+0x5a] = instrCB__BIT_3_D
	/* BIT 3,E */
	opcodesMap[SHIFT_0xCB+0x5b] = instrCB__BIT_3_E
	/* BIT 3,H */
	opcodesMap[SHIFT_0xCB+0x5c] = instrCB__BIT_3_H
	/* BIT 3,L */
	opcodesMap[SHIFT_0xCB+0x5d] = instrCB__BIT_3_L
	/* BIT 3,(HL) */
	opcodesMap[SHIFT_0xCB+0x5e] = instrCB__BIT_3_iHL
	/* BIT 3,A */
	opcodesMap[SHIFT_0xCB+0x5f] = instrCB__BIT_3_A
	/* BIT 4,B */
	opcodesMap[SHIFT_0xCB+0x60] = instrCB__BIT_4_B
	/* BIT 4,C */
	opcodesMap[SHIFT_0xCB+0x61] = instrCB__BIT_4_C
	/* BIT 4,D */
	opcodesMap[SHIFT_0xCB+0x62] = instrCB__BIT_4_D
	/* BIT 4,E */
	opcodesMap[SHIFT_0xCB+0x63] = instrCB__BIT_4_E
	/* BIT 4,H */
	opcodesMap[SHIFT_0xCB+0x64] = instrCB__BIT_4_H
	/* BIT 4,L */
	opcodesMap[SHIFT_0xCB+0x65] = instrCB__BIT_4_L
	/* BIT 4,(HL) */
	opcodesMap[SHIFT_0xCB+0x66] = instrCB__BIT_4_iHL
	/* BIT 4,A */
	opcodesMap[SHIFT_0xCB+0x67] = instrCB__BIT_4_A
	/* BIT 5,B */
	opcodesMap[SHIFT_0xCB+0x68] = instrCB__BIT_5_B
	/* BIT 5,C */
	opcodesMap[SHIFT_0xCB+0x69] = instrCB__BIT_5_C
	/* BIT 5,D */
	opcodesMap[SHIFT_0xCB+0x6a] = instrCB__BIT_5_D
	/* BIT 5,E */
	opcodesMap[SHIFT_0xCB+0x6b] = instrCB__BIT_5_E
	/* BIT 5,H */
	opcodesMap[SHIFT_0xCB+0x6c] = instrCB__BIT_5_H
	/* BIT 5,L */
	opcodesMap[SHIFT_0xCB+0x6d] = instrCB__BIT_5_L
	/* BIT 5,(HL) */
	opcodesMap[SHIFT_0xCB+0x6e] = instrCB__BIT_5_iHL
	/* BIT 5,A */
	opcodesMap[SHIFT_0xCB+0x6f] = instrCB__BIT_5_A
	/* BIT 6,B */
	opcodesMap[SHIFT_0xCB+0x70] = instrCB__BIT_6_B
	/* BIT 6,C */
	opcodesMap[SHIFT_0xCB+0x71] = instrCB__BIT_6_C
	/* BIT 6,D */
	opcodesMap[SHIFT_0xCB+0x72] = instrCB__BIT_6_D
	/* BIT 6,E */
	opcodesMap[SHIFT_0xCB+0x73] = instrCB__BIT_6_E
	/* BIT 6,H */
	opcodesMap[SHIFT_0xCB+0x74] = instrCB__BIT_6_H
	/* BIT 6,L */
	opcodesMap[SHIFT_0xCB+0x75] = instrCB__BIT_6_L
	/* BIT 6,(HL) */
	opcodesMap[SHIFT_0xCB+0x76] = instrCB__BIT_6_iHL
	/* BIT 6,A */
	opcodesMap[SHIFT_0xCB+0x77] = instrCB__BIT_6_A
	/* BIT 7,B */
	opcodesMap[SHIFT_0xCB+0x78] = instrCB__BIT_7_B
	/* BIT 7,C */
	opcodesMap[SHIFT_0xCB+0x79] = instrCB__BIT_7_C
	/* BIT 7,D */
	opcodesMap[SHIFT_0xCB+0x7a] = instrCB__BIT_7_D
	/* BIT 7,E */
	opcodesMap[SHIFT_0xCB+0x7b] = instrCB__BIT_7_E
	/* BIT 7,H */
	opcodesMap[SHIFT_0xCB+0x7c] = instrCB__BIT_7_H
	/* BIT 7,L */
	opcodesMap[SHIFT_0xCB+0x7d] = instrCB__BIT_7_L
	/* BIT 7,(HL) */
	opcodesMap[SHIFT_0xCB+0x7e] = instrCB__BIT_7_iHL
	/* BIT 7,A */
	opcodesMap[SHIFT_0xCB+0x7f] = instrCB__BIT_7_A
	/* RES 0,B */
	opcodesMap[SHIFT_0xCB+0x80] = instrCB__RES_0_B
	/* RES 0,C */
	opcodesMap[SHIFT_0xCB+0x81] = instrCB__RES_0_C
	/* RES 0,D */
	opcodesMap[SHIFT_0xCB+0x82] = instrCB__RES_0_D
	/* RES 0,E */
	opcodesMap[SHIFT_0xCB+0x83] = instrCB__RES_0_E
	/* RES 0,H */
	opcodesMap[SHIFT_0xCB+0x84] = instrCB__RES_0_H
	/* RES 0,L */
	opcodesMap[SHIFT_0xCB+0x85] = instrCB__RES_0_L
	/* RES 0,(HL) */
	opcodesMap[SHIFT_0xCB+0x86] = instrCB__RES_0_iHL
	/* RES 0,A */
	opcodesMap[SHIFT_0xCB+0x87] = instrCB__RES_0_A
	/* RES 1,B */
	opcodesMap[SHIFT_0xCB+0x88] = instrCB__RES_1_B
	/* RES 1,C */
	opcodesMap[SHIFT_0xCB+0x89] = instrCB__RES_1_C
	/* RES 1,D */
	opcodesMap[SHIFT_0xCB+0x8a] = instrCB__RES_1_D
	/* RES 1,E */
	opcodesMap[SHIFT_0xCB+0x8b] = instrCB__RES_1_E
	/* RES 1,H */
	opcodesMap[SHIFT_0xCB+0x8c] = instrCB__RES_1_H
	/* RES 1,L */
	opcodesMap[SHIFT_0xCB+0x8d] = instrCB__RES_1_L
	/* RES 1,(HL) */
	opcodesMap[SHIFT_0xCB+0x8e] = instrCB__RES_1_iHL
	/* RES 1,A */
	opcodesMap[SHIFT_0xCB+0x8f] = instrCB__RES_1_A
	/* RES 2,B */
	opcodesMap[SHIFT_0xCB+0x90] = instrCB__RES_2_B
	/* RES 2,C */
	opcodesMap[SHIFT_0xCB+0x91] = instrCB__RES_2_C
	/* RES 2,D */
	opcodesMap[SHIFT_0xCB+0x92] = instrCB__RES_2_D
	/* RES 2,E */
	opcodesMap[SHIFT_0xCB+0x93] = instrCB__RES_2_E
	/* RES 2,H */
	opcodesMap[SHIFT_0xCB+0x94] = instrCB__RES_2_H
	/* RES 2,L */
	opcodesMap[SHIFT_0xCB+0x95] = instrCB__RES_2_L
	/* RES 2,(HL) */
	opcodesMap[SHIFT_0xCB+0x96] = instrCB__RES_2_iHL
	/* RES 2,A */
	opcodesMap[SHIFT_0xCB+0x97] = instrCB__RES_2_A
	/* RES 3,B */
	opcodesMap[SHIFT_0xCB+0x98] = instrCB__RES_3_B
	/* RES 3,C */
	opcodesMap[SHIFT_0xCB+0x99] = instrCB__RES_3_C
	/* RES 3,D */
	opcodesMap[SHIFT_0xCB+0x9a] = instrCB__RES_3_D
	/* RES 3,E */
	opcodesMap[SHIFT_0xCB+0x9b] = instrCB__RES_3_E
	/* RES 3,H */
	opcodesMap[SHIFT_0xCB+0x9c] = instrCB__RES_3_H
	/* RES 3,L */
	opcodesMap[SHIFT_0xCB+0x9d] = instrCB__RES_3_L
	/* RES 3,(HL) */
	opcodesMap[SHIFT_0xCB+0x9e] = instrCB__RES_3_iHL
	/* RES 3,A */
	opcodesMap[SHIFT_0xCB+0x9f] = instrCB__RES_3_A
	/* RES 4,B */
	opcodesMap[SHIFT_0xCB+0xa0] = instrCB__RES_4_B
	/* RES 4,C */
	opcodesMap[SHIFT_0xCB+0xa1] = instrCB__RES_4_C
	/* RES 4,D */
	opcodesMap[SHIFT_0xCB+0xa2] = instrCB__RES_4_D
	/* RES 4,E */
	opcodesMap[SHIFT_0xCB+0xa3] = instrCB__RES_4_E
	/* RES 4,H */
	opcodesMap[SHIFT_0xCB+0xa4] = instrCB__RES_4_H
	/* RES 4,L */
	opcodesMap[SHIFT_0xCB+0xa5] = instrCB__RES_4_L
	/* RES 4,(HL) */
	opcodesMap[SHIFT_0xCB+0xa6] = instrCB__RES_4_iHL
	/* RES 4,A */
	opcodesMap[SHIFT_0xCB+0xa7] = instrCB__RES_4_A
	/* RES 5,B */
	opcodesMap[SHIFT_0xCB+0xa8] = instrCB__RES_5_B
	/* RES 5,C */
	opcodesMap[SHIFT_0xCB+0xa9] = instrCB__RES_5_C
	/* RES 5,D */
	opcodesMap[SHIFT_0xCB+0xaa] = instrCB__RES_5_D
	/* RES 5,E */
	opcodesMap[SHIFT_0xCB+0xab] = instrCB__RES_5_E
	/* RES 5,H */
	opcodesMap[SHIFT_0xCB+0xac] = instrCB__RES_5_H
	/* RES 5,L */
	opcodesMap[SHIFT_0xCB+0xad] = instrCB__RES_5_L
	/* RES 5,(HL) */
	opcodesMap[SHIFT_0xCB+0xae] = instrCB__RES_5_iHL
	/* RES 5,A */
	opcodesMap[SHIFT_0xCB+0xaf] = instrCB__RES_5_A
	/* RES 6,B */
	opcodesMap[SHIFT_0xCB+0xb0] = instrCB__RES_6_B
	/* RES 6,C */
	opcodesMap[SHIFT_0xCB+0xb1] = instrCB__RES_6_C
	/* RES 6,D */
	opcodesMap[SHIFT_0xCB+0xb2] = instrCB__RES_6_D
	/* RES 6,E */
	opcodesMap[SHIFT_0xCB+0xb3] = instrCB__RES_6_E
	/* RES 6,H */
	opcodesMap[SHIFT_0xCB+0xb4] = instrCB__RES_6_H
	/* RES 6,L */
	opcodesMap[SHIFT_0xCB+0xb5] = instrCB__RES_6_L
	/* RES 6,(HL) */
	opcodesMap[SHIFT_0xCB+0xb6] = instrCB__RES_6_iHL
	/* RES 6,A */
	opcodesMap[SHIFT_0xCB+0xb7] = instrCB__RES_6_A
	/* RES 7,B */
	opcodesMap[SHIFT_0xCB+0xb8] = instrCB__RES_7_B
	/* RES 7,C */
	opcodesMap[SHIFT_0xCB+0xb9] = instrCB__RES_7_C
	/* RES 7,D */
	opcodesMap[SHIFT_0xCB+0xba] = instrCB__RES_7_D
	/* RES 7,E */
	opcodesMap[SHIFT_0xCB+0xbb] = instrCB__RES_7_E
	/* RES 7,H */
	opcodesMap[SHIFT_0xCB+0xbc] = instrCB__RES_7_H
	/* RES 7,L */
	opcodesMap[SHIFT_0xCB+0xbd] = instrCB__RES_7_L
	/* RES 7,(HL) */
	opcodesMap[SHIFT_0xCB+0xbe] = instrCB__RES_7_iHL
	/* RES 7,A */
	opcodesMap[SHIFT_0xCB+0xbf] = instrCB__RES_7_A
	/* SET 0,B */
	opcodesMap[SHIFT_0xCB+0xc0] = instrCB__SET_0_B
	/* SET 0,C */
	opcodesMap[SHIFT_0xCB+0xc1] = instrCB__SET_0_C
	/* SET 0,D */
	opcodesMap[SHIFT_0xCB+0xc2] = instrCB__SET_0_D
	/* SET 0,E */
	opcodesMap[SHIFT_0xCB+0xc3] = instrCB__SET_0_E
	/* SET 0,H */
	opcodesMap[SHIFT_0xCB+0xc4] = instrCB__SET_0_H
	/* SET 0,L */
	opcodesMap[SHIFT_0xCB+0xc5] = instrCB__SET_0_L
	/* SET 0,(HL) */
	opcodesMap[SHIFT_0xCB+0xc6] = instrCB__SET_0_iHL
	/* SET 0,A */
	opcodesMap[SHIFT_0xCB+0xc7] = instrCB__SET_0_A
	/* SET 1,B */
	opcodesMap[SHIFT_0xCB+0xc8] = instrCB__SET_1_B
	/* SET 1,C */
	opcodesMap[SHIFT_0xCB+0xc9] = instrCB__SET_1_C
	/* SET 1,D */
	opcodesMap[SHIFT_0xCB+0xca] = instrCB__SET_1_D
	/* SET 1,E */
	opcodesMap[SHIFT_0xCB+0xcb] = instrCB__SET_1_E
	/* SET 1,H */
	opcodesMap[SHIFT_0xCB+0xcc] = instrCB__SET_1_H
	/* SET 1,L */
	opcodesMap[SHIFT_0xCB+0xcd] = instrCB__SET_1_L
	/* SET 1,(HL) */
	opcodesMap[SHIFT_0xCB+0xce] = instrCB__SET_1_iHL
	/* SET 1,A */
	opcodesMap[SHIFT_0xCB+0xcf] = instrCB__SET_1_A
	/* SET 2,B */
	opcodesMap[SHIFT_0xCB+0xd0] = instrCB__SET_2_B
	/* SET 2,C */
	opcodesMap[SHIFT_0xCB+0xd1] = instrCB__SET_2_C
	/* SET 2,D */
	opcodesMap[SHIFT_0xCB+0xd2] = instrCB__SET_2_D
	/* SET 2,E */
	opcodesMap[SHIFT_0xCB+0xd3] = instrCB__SET_2_E
	/* SET 2,H */
	opcodesMap[SHIFT_0xCB+0xd4] = instrCB__SET_2_H
	/* SET 2,L */
	opcodesMap[SHIFT_0xCB+0xd5] = instrCB__SET_2_L
	/* SET 2,(HL) */
	opcodesMap[SHIFT_0xCB+0xd6] = instrCB__SET_2_iHL
	/* SET 2,A */
	opcodesMap[SHIFT_0xCB+0xd7] = instrCB__SET_2_A
	/* SET 3,B */
	opcodesMap[SHIFT_0xCB+0xd8] = instrCB__SET_3_B
	/* SET 3,C */
	opcodesMap[SHIFT_0xCB+0xd9] = instrCB__SET_3_C
	/* SET 3,D */
	opcodesMap[SHIFT_0xCB+0xda] = instrCB__SET_3_D
	/* SET 3,E */
	opcodesMap[SHIFT_0xCB+0xdb] = instrCB__SET_3_E
	/* SET 3,H */
	opcodesMap[SHIFT_0xCB+0xdc] = instrCB__SET_3_H
	/* SET 3,L */
	opcodesMap[SHIFT_0xCB+0xdd] = instrCB__SET_3_L
	/* SET 3,(HL) */
	opcodesMap[SHIFT_0xCB+0xde] = instrCB__SET_3_iHL
	/* SET 3,A */
	opcodesMap[SHIFT_0xCB+0xdf] = instrCB__SET_3_A
	/* SET 4,B */
	opcodesMap[SHIFT_0xCB+0xe0] = instrCB__SET_4_B
	/* SET 4,C */
	opcodesMap[SHIFT_0xCB+0xe1] = instrCB__SET_4_C
	/* SET 4,D */
	opcodesMap[SHIFT_0xCB+0xe2] = instrCB__SET_4_D
	/* SET 4,E */
	opcodesMap[SHIFT_0xCB+0xe3] = instrCB__SET_4_E
	/* SET 4,H */
	opcodesMap[SHIFT_0xCB+0xe4] = instrCB__SET_4_H
	/* SET 4,L */
	opcodesMap[SHIFT_0xCB+0xe5] = instrCB__SET_4_L
	/* SET 4,(HL) */
	opcodesMap[SHIFT_0xCB+0xe6] = instrCB__SET_4_iHL
	/* SET 4,A */
	opcodesMap[SHIFT_0xCB+0xe7] = instrCB__SET_4_A
	/* SET 5,B */
	opcodesMap[SHIFT_0xCB+0xe8] = instrCB__SET_5_B
	/* SET 5,C */
	opcodesMap[SHIFT_0xCB+0xe9] = instrCB__SET_5_C
	/* SET 5,D */
	opcodesMap[SHIFT_0xCB+0xea] = instrCB__SET_5_D
	/* SET 5,E */
	opcodesMap[SHIFT_0xCB+0xeb] = instrCB__SET_5_E
	/* SET 5,H */
	opcodesMap[SHIFT_0xCB+0xec] = instrCB__SET_5_H
	/* SET 5,L */
	opcodesMap[SHIFT_0xCB+0xed] = instrCB__SET_5_L
	/* SET 5,(HL) */
	opcodesMap[SHIFT_0xCB+0xee] = instrCB__SET_5_iHL
	/* SET 5,A */
	opcodesMap[SHIFT_0xCB+0xef] = instrCB__SET_5_A
	/* SET 6,B */
	opcodesMap[SHIFT_0xCB+0xf0] = instrCB__SET_6_B
	/* SET 6,C */
	opcodesMap[SHIFT_0xCB+0xf1] = instrCB__SET_6_C
	/* SET 6,D */
	opcodesMap[SHIFT_0xCB+0xf2] = instrCB__SET_6_D
	/* SET 6,E */
	opcodesMap[SHIFT_0xCB+0xf3] = instrCB__SET_6_E
	/* SET 6,H */
	opcodesMap[SHIFT_0xCB+0xf4] = instrCB__SET_6_H
	/* SET 6,L */
	opcodesMap[SHIFT_0xCB+0xf5] = instrCB__SET_6_L
	/* SET 6,(HL) */
	opcodesMap[SHIFT_0xCB+0xf6] = instrCB__SET_6_iHL
	/* SET 6,A */
	opcodesMap[SHIFT_0xCB+0xf7] = instrCB__SET_6_A
	/* SET 7,B */
	opcodesMap[SHIFT_0xCB+0xf8] = instrCB__SET_7_B
	/* SET 7,C */
	opcodesMap[SHIFT_0xCB+0xf9] = instrCB__SET_7_C
	/* SET 7,D */
	opcodesMap[SHIFT_0xCB+0xfa] = instrCB__SET_7_D
	/* SET 7,E */
	opcodesMap[SHIFT_0xCB+0xfb] = instrCB__SET_7_E
	/* SET 7,H */
	opcodesMap[SHIFT_0xCB+0xfc] = instrCB__SET_7_H
	/* SET 7,L */
	opcodesMap[SHIFT_0xCB+0xfd] = instrCB__SET_7_L
	/* SET 7,(HL) */
	opcodesMap[SHIFT_0xCB+0xfe] = instrCB__SET_7_iHL
	/* SET 7,A */
	opcodesMap[SHIFT_0xCB+0xff] = instrCB__SET_7_A

	// END of 0xcb shifted opcodes

	// BEGIN of 0xed shifted opcodes
	/* IN B,(C) */
	opcodesMap[SHIFT_0xED+0x40] = instrED__IN_B_iC
	/* OUT (C),B */
	opcodesMap[SHIFT_0xED+0x41] = instrED__OUT_iC_B
	/* SBC HL,BC */
	opcodesMap[SHIFT_0xED+0x42] = instrED__SBC_HL_BC
	/* LD (nnnn),BC */
	opcodesMap[SHIFT_0xED+0x43] = instrED__LD_iNNNN_BC
	/* NEG */
	opcodesMap[SHIFT_0xED+0x7c] = instrED__NEG
	// Fallthrough cases
	opcodesMap[SHIFT_0xED+0x44] = opcodesMap[SHIFT_0xED+0x7c]
	opcodesMap[SHIFT_0xED+0x4c] = opcodesMap[SHIFT_0xED+0x7c]
	opcodesMap[SHIFT_0xED+0x54] = opcodesMap[SHIFT_0xED+0x7c]
	opcodesMap[SHIFT_0xED+0x5c] = opcodesMap[SHIFT_0xED+0x7c]
	opcodesMap[SHIFT_0xED+0x64] = opcodesMap[SHIFT_0xED+0x7c]
	opcodesMap[SHIFT_0xED+0x6c] = opcodesMap[SHIFT_0xED+0x7c]
	opcodesMap[SHIFT_0xED+0x74] = opcodesMap[SHIFT_0xED+0x7c]
	/* RETN */
	opcodesMap[SHIFT_0xED+0x7d] = instrED__RETN
	// Fallthrough cases
	opcodesMap[SHIFT_0xED+0x45] = opcodesMap[SHIFT_0xED+0x7d]
	opcodesMap[SHIFT_0xED+0x4d] = opcodesMap[SHIFT_0xED+0x7d]
	opcodesMap[SHIFT_0xED+0x55] = opcodesMap[SHIFT_0xED+0x7d]
	opcodesMap[SHIFT_0xED+0x5d] = opcodesMap[SHIFT_0xED+0x7d]
	opcodesMap[SHIFT_0xED+0x65] = opcodesMap[SHIFT_0xED+0x7d]
	opcodesMap[SHIFT_0xED+0x6d] = opcodesMap[SHIFT_0xED+0x7d]
	opcodesMap[SHIFT_0xED+0x75] = opcodesMap[SHIFT_0xED+0x7d]
	/* IM 0 */
	opcodesMap[SHIFT_0xED+0x6e] = instrED__IM_0
	// Fallthrough cases
	opcodesMap[SHIFT_0xED+0x46] = opcodesMap[SHIFT_0xED+0x6e]
	opcodesMap[SHIFT_0xED+0x4e] = opcodesMap[SHIFT_0xED+0x6e]
	opcodesMap[SHIFT_0xED+0x66] = opcodesMap[SHIFT_0xED+0x6e]
	/* LD I,A */
	opcodesMap[SHIFT_0xED+0x47] = instrED__LD_I_A
	/* IN C,(C) */
	opcodesMap[SHIFT_0xED+0x48] = instrED__IN_C_iC
	/* OUT (C),C */
	opcodesMap[SHIFT_0xED+0x49] = instrED__OUT_iC_C
	/* ADC HL,BC */
	opcodesMap[SHIFT_0xED+0x4a] = instrED__ADC_HL_BC
	/* LD BC,(nnnn) */
	opcodesMap[SHIFT_0xED+0x4b] = instrED__LD_BC_iNNNN
	/* LD R,A */
	opcodesMap[SHIFT_0xED+0x4f] = instrED__LD_R_A
	/* IN D,(C) */
	opcodesMap[SHIFT_0xED+0x50] = instrED__IN_D_iC
	/* OUT (C),D */
	opcodesMap[SHIFT_0xED+0x51] = instrED__OUT_iC_D
	/* SBC HL,DE */
	opcodesMap[SHIFT_0xED+0x52] = instrED__SBC_HL_DE
	/* LD (nnnn),DE */
	opcodesMap[SHIFT_0xED+0x53] = instrED__LD_iNNNN_DE
	/* IM 1 */
	opcodesMap[SHIFT_0xED+0x76] = instrED__IM_1
	// Fallthrough cases
	opcodesMap[SHIFT_0xED+0x56] = opcodesMap[SHIFT_0xED+0x76]
	/* LD A,I */
	opcodesMap[SHIFT_0xED+0x57] = instrED__LD_A_I
	/* IN E,(C) */
	opcodesMap[SHIFT_0xED+0x58] = instrED__IN_E_iC
	/* OUT (C),E */
	opcodesMap[SHIFT_0xED+0x59] = instrED__OUT_iC_E
	/* ADC HL,DE */
	opcodesMap[SHIFT_0xED+0x5a] = instrED__ADC_HL_DE
	/* LD DE,(nnnn) */
	opcodesMap[SHIFT_0xED+0x5b] = instrED__LD_DE_iNNNN
	/* IM 2 */
	opcodesMap[SHIFT_0xED+0x7e] = instrED__IM_2
	// Fallthrough cases
	opcodesMap[SHIFT_0xED+0x5e] = opcodesMap[SHIFT_0xED+0x7e]
	/* LD A,R */
	opcodesMap[SHIFT_0xED+0x5f] = instrED__LD_A_R
	/* IN H,(C) */
	opcodesMap[SHIFT_0xED+0x60] = instrED__IN_H_iC
	/* OUT (C),H */
	opcodesMap[SHIFT_0xED+0x61] = instrED__OUT_iC_H
	/* SBC HL,HL */
	opcodesMap[SHIFT_0xED+0x62] = instrED__SBC_HL_HL
	/* LD (nnnn),HL */
	opcodesMap[SHIFT_0xED+0x63] = instrED__LD_iNNNN_HL
	/* RRD */
	opcodesMap[SHIFT_0xED+0x67] = instrED__RRD
	/* IN L,(C) */
	opcodesMap[SHIFT_0xED+0x68] = instrED__IN_L_iC
	/* OUT (C),L */
	opcodesMap[SHIFT_0xED+0x69] = instrED__OUT_iC_L
	/* ADC HL,HL */
	opcodesMap[SHIFT_0xED+0x6a] = instrED__ADC_HL_HL
	/* LD HL,(nnnn) */
	opcodesMap[SHIFT_0xED+0x6b] = instrED__LD_HL_iNNNN
	/* RLD */
	opcodesMap[SHIFT_0xED+0x6f] = instrED__RLD
	/* IN F,(C) */
	opcodesMap[SHIFT_0xED+0x70] = instrED__IN_F_iC
	/* OUT (C),0 */
	opcodesMap[SHIFT_0xED+0x71] = instrED__OUT_iC_0
	/* SBC HL,SP */
	opcodesMap[SHIFT_0xED+0x72] = instrED__SBC_HL_SP
	/* LD (nnnn),SP */
	opcodesMap[SHIFT_0xED+0x73] = instrED__LD_iNNNN_SP
	/* IN A,(C) */
	opcodesMap[SHIFT_0xED+0x78] = instrED__IN_A_iC
	/* OUT (C),A */
	opcodesMap[SHIFT_0xED+0x79] = instrED__OUT_iC_A
	/* ADC HL,SP */
	opcodesMap[SHIFT_0xED+0x7a] = instrED__ADC_HL_SP
	/* LD SP,(nnnn) */
	opcodesMap[SHIFT_0xED+0x7b] = instrED__LD_SP_iNNNN
	/* LDI */
	opcodesMap[SHIFT_0xED+0xa0] = instrED__LDI
	/* CPI */
	opcodesMap[SHIFT_0xED+0xa1] = instrED__CPI
	/* INI */
	opcodesMap[SHIFT_0xED+0xa2] = instrED__INI
	/* OUTI */
	opcodesMap[SHIFT_0xED+0xa3] = instrED__OUTI
	/* LDD */
	opcodesMap[SHIFT_0xED+0xa8] = instrED__LDD
	/* CPD */
	opcodesMap[SHIFT_0xED+0xa9] = instrED__CPD
	/* IND */
	opcodesMap[SHIFT_0xED+0xaa] = instrED__IND
	/* OUTD */
	opcodesMap[SHIFT_0xED+0xab] = instrED__OUTD
	/* LDIR */
	opcodesMap[SHIFT_0xED+0xb0] = instrED__LDIR
	/* CPIR */
	opcodesMap[SHIFT_0xED+0xb1] = instrED__CPIR
	/* INIR */
	opcodesMap[SHIFT_0xED+0xb2] = instrED__INIR
	/* OTIR */
	opcodesMap[SHIFT_0xED+0xb3] = instrED__OTIR
	/* LDDR */
	opcodesMap[SHIFT_0xED+0xb8] = instrED__LDDR
	/* CPDR */
	opcodesMap[SHIFT_0xED+0xb9] = instrED__CPDR
	/* INDR */
	opcodesMap[SHIFT_0xED+0xba] = instrED__INDR
	/* OTDR */
	opcodesMap[SHIFT_0xED+0xbb] = instrED__OTDR
	/* slttrap */
	opcodesMap[SHIFT_0xED+0xfb] = instrED__SLTTRAP

	// END of 0xed shifted opcodes

	// BEGIN of 0xdd shifted opcodes
	/* ADD REGISTER,BC */
	opcodesMap[SHIFT_0xDD+0x09] = instrDD__ADD_REG_BC
	/* ADD REGISTER,DE */
	opcodesMap[SHIFT_0xDD+0x19] = instrDD__ADD_REG_DE
	/* LD REGISTER,nnnn */
	opcodesMap[SHIFT_0xDD+0x21] = instrDD__LD_REG_NNNN
	/* LD (nnnn),REGISTER */
	opcodesMap[SHIFT_0xDD+0x22] = instrDD__LD_iNNNN_REG
	/* INC REGISTER */
	opcodesMap[SHIFT_0xDD+0x23] = instrDD__INC_REG
	/* INC REGISTERH */
	opcodesMap[SHIFT_0xDD+0x24] = instrDD__INC_REGH
	/* DEC REGISTERH */
	opcodesMap[SHIFT_0xDD+0x25] = instrDD__DEC_REGH
	/* LD REGISTERH,nn */
	opcodesMap[SHIFT_0xDD+0x26] = instrDD__LD_REGH_NN
	/* ADD REGISTER,REGISTER */
	opcodesMap[SHIFT_0xDD+0x29] = instrDD__ADD_REG_REG
	/* LD REGISTER,(nnnn) */
	opcodesMap[SHIFT_0xDD+0x2a] = instrDD__LD_REG_iNNNN
	/* DEC REGISTER */
	opcodesMap[SHIFT_0xDD+0x2b] = instrDD__DEC_REG
	/* INC REGISTERL */
	opcodesMap[SHIFT_0xDD+0x2c] = instrDD__INC_REGL
	/* DEC REGISTERL */
	opcodesMap[SHIFT_0xDD+0x2d] = instrDD__DEC_REGL
	/* LD REGISTERL,nn */
	opcodesMap[SHIFT_0xDD+0x2e] = instrDD__LD_REGL_NN
	/* INC (REGISTER+dd) */
	opcodesMap[SHIFT_0xDD+0x34] = instrDD__INC_iREGpDD
	/* DEC (REGISTER+dd) */
	opcodesMap[SHIFT_0xDD+0x35] = instrDD__DEC_iREGpDD
	/* LD (REGISTER+dd),nn */
	opcodesMap[SHIFT_0xDD+0x36] = instrDD__LD_iREGpDD_NN
	/* ADD REGISTER,SP */
	opcodesMap[SHIFT_0xDD+0x39] = instrDD__ADD_REG_SP
	/* LD B,REGISTERH */
	opcodesMap[SHIFT_0xDD+0x44] = instrDD__LD_B_REGH
	/* LD B,REGISTERL */
	opcodesMap[SHIFT_0xDD+0x45] = instrDD__LD_B_REGL
	/* LD B,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDD+0x46] = instrDD__LD_B_iREGpDD
	/* LD C,REGISTERH */
	opcodesMap[SHIFT_0xDD+0x4c] = instrDD__LD_C_REGH
	/* LD C,REGISTERL */
	opcodesMap[SHIFT_0xDD+0x4d] = instrDD__LD_C_REGL
	/* LD C,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDD+0x4e] = instrDD__LD_C_iREGpDD
	/* LD D,REGISTERH */
	opcodesMap[SHIFT_0xDD+0x54] = instrDD__LD_D_REGH
	/* LD D,REGISTERL */
	opcodesMap[SHIFT_0xDD+0x55] = instrDD__LD_D_REGL
	/* LD D,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDD+0x56] = instrDD__LD_D_iREGpDD
	/* LD E,REGISTERH */
	opcodesMap[SHIFT_0xDD+0x5c] = instrDD__LD_E_REGH
	/* LD E,REGISTERL */
	opcodesMap[SHIFT_0xDD+0x5d] = instrDD__LD_E_REGL
	/* LD E,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDD+0x5e] = instrDD__LD_E_iREGpDD
	/* LD REGISTERH,B */
	opcodesMap[SHIFT_0xDD+0x60] = instrDD__LD_REGH_B
	/* LD REGISTERH,C */
	opcodesMap[SHIFT_0xDD+0x61] = instrDD__LD_REGH_C
	/* LD REGISTERH,D */
	opcodesMap[SHIFT_0xDD+0x62] = instrDD__LD_REGH_D
	/* LD REGISTERH,E */
	opcodesMap[SHIFT_0xDD+0x63] = instrDD__LD_REGH_E
	/* LD REGISTERH,REGISTERH */
	opcodesMap[SHIFT_0xDD+0x64] = instrDD__LD_REGH_REGH
	/* LD REGISTERH,REGISTERL */
	opcodesMap[SHIFT_0xDD+0x65] = instrDD__LD_REGH_REGL
	/* LD H,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDD+0x66] = instrDD__LD_H_iREGpDD
	/* LD REGISTERH,A */
	opcodesMap[SHIFT_0xDD+0x67] = instrDD__LD_REGH_A
	/* LD REGISTERL,B */
	opcodesMap[SHIFT_0xDD+0x68] = instrDD__LD_REGL_B
	/* LD REGISTERL,C */
	opcodesMap[SHIFT_0xDD+0x69] = instrDD__LD_REGL_C
	/* LD REGISTERL,D */
	opcodesMap[SHIFT_0xDD+0x6a] = instrDD__LD_REGL_D
	/* LD REGISTERL,E */
	opcodesMap[SHIFT_0xDD+0x6b] = instrDD__LD_REGL_E
	/* LD REGISTERL,REGISTERH */
	opcodesMap[SHIFT_0xDD+0x6c] = instrDD__LD_REGL_REGH
	/* LD REGISTERL,REGISTERL */
	opcodesMap[SHIFT_0xDD+0x6d] = instrDD__LD_REGL_REGL
	/* LD L,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDD+0x6e] = instrDD__LD_L_iREGpDD
	/* LD REGISTERL,A */
	opcodesMap[SHIFT_0xDD+0x6f] = instrDD__LD_REGL_A
	/* LD (REGISTER+dd),B */
	opcodesMap[SHIFT_0xDD+0x70] = instrDD__LD_iREGpDD_B
	/* LD (REGISTER+dd),C */
	opcodesMap[SHIFT_0xDD+0x71] = instrDD__LD_iREGpDD_C
	/* LD (REGISTER+dd),D */
	opcodesMap[SHIFT_0xDD+0x72] = instrDD__LD_iREGpDD_D
	/* LD (REGISTER+dd),E */
	opcodesMap[SHIFT_0xDD+0x73] = instrDD__LD_iREGpDD_E
	/* LD (REGISTER+dd),H */
	opcodesMap[SHIFT_0xDD+0x74] = instrDD__LD_iREGpDD_H
	/* LD (REGISTER+dd),L */
	opcodesMap[SHIFT_0xDD+0x75] = instrDD__LD_iREGpDD_L
	/* LD (REGISTER+dd),A */
	opcodesMap[SHIFT_0xDD+0x77] = instrDD__LD_iREGpDD_A
	/* LD A,REGISTERH */
	opcodesMap[SHIFT_0xDD+0x7c] = instrDD__LD_A_REGH
	/* LD A,REGISTERL */
	opcodesMap[SHIFT_0xDD+0x7d] = instrDD__LD_A_REGL
	/* LD A,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDD+0x7e] = instrDD__LD_A_iREGpDD
	/* ADD A,REGISTERH */
	opcodesMap[SHIFT_0xDD+0x84] = instrDD__ADD_A_REGH
	/* ADD A,REGISTERL */
	opcodesMap[SHIFT_0xDD+0x85] = instrDD__ADD_A_REGL
	/* ADD A,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDD+0x86] = instrDD__ADD_A_iREGpDD
	/* ADC A,REGISTERH */
	opcodesMap[SHIFT_0xDD+0x8c] = instrDD__ADC_A_REGH
	/* ADC A,REGISTERL */
	opcodesMap[SHIFT_0xDD+0x8d] = instrDD__ADC_A_REGL
	/* ADC A,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDD+0x8e] = instrDD__ADC_A_iREGpDD
	/* SUB A,REGISTERH */
	opcodesMap[SHIFT_0xDD+0x94] = instrDD__SUB_A_REGH
	/* SUB A,REGISTERL */
	opcodesMap[SHIFT_0xDD+0x95] = instrDD__SUB_A_REGL
	/* SUB A,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDD+0x96] = instrDD__SUB_A_iREGpDD
	/* SBC A,REGISTERH */
	opcodesMap[SHIFT_0xDD+0x9c] = instrDD__SBC_A_REGH
	/* SBC A,REGISTERL */
	opcodesMap[SHIFT_0xDD+0x9d] = instrDD__SBC_A_REGL
	/* SBC A,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDD+0x9e] = instrDD__SBC_A_iREGpDD
	/* AND A,REGISTERH */
	opcodesMap[SHIFT_0xDD+0xa4] = instrDD__AND_A_REGH
	/* AND A,REGISTERL */
	opcodesMap[SHIFT_0xDD+0xa5] = instrDD__AND_A_REGL
	/* AND A,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDD+0xa6] = instrDD__AND_A_iREGpDD
	/* XOR A,REGISTERH */
	opcodesMap[SHIFT_0xDD+0xac] = instrDD__XOR_A_REGH
	/* XOR A,REGISTERL */
	opcodesMap[SHIFT_0xDD+0xad] = instrDD__XOR_A_REGL
	/* XOR A,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDD+0xae] = instrDD__XOR_A_iREGpDD
	/* OR A,REGISTERH */
	opcodesMap[SHIFT_0xDD+0xb4] = instrDD__OR_A_REGH
	/* OR A,REGISTERL */
	opcodesMap[SHIFT_0xDD+0xb5] = instrDD__OR_A_REGL
	/* OR A,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDD+0xb6] = instrDD__OR_A_iREGpDD
	/* CP A,REGISTERH */
	opcodesMap[SHIFT_0xDD+0xbc] = instrDD__CP_A_REGH
	/* CP A,REGISTERL */
	opcodesMap[SHIFT_0xDD+0xbd] = instrDD__CP_A_REGL
	/* CP A,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDD+0xbe] = instrDD__CP_A_iREGpDD
	/* shift DDFDCB */
	opcodesMap[SHIFT_0xDD+0xcb] = instrDD__SHIFT_DDFDCB
	/* POP REGISTER */
	opcodesMap[SHIFT_0xDD+0xe1] = instrDD__POP_REG
	/* EX (SP),REGISTER */
	opcodesMap[SHIFT_0xDD+0xe3] = instrDD__EX_iSP_REG
	/* PUSH REGISTER */
	opcodesMap[SHIFT_0xDD+0xe5] = instrDD__PUSH_REG
	/* JP REGISTER */
	opcodesMap[SHIFT_0xDD+0xe9] = instrDD__JP_REG
	/* LD SP,REGISTER */
	opcodesMap[SHIFT_0xDD+0xf9] = instrDD__LD_SP_REG

	// END of 0xdd shifted opcodes

	// BEGIN of 0xfd shifted opcodes
	/* ADD REGISTER,BC */
	opcodesMap[SHIFT_0xFD+0x09] = instrFD__ADD_REG_BC
	/* ADD REGISTER,DE */
	opcodesMap[SHIFT_0xFD+0x19] = instrFD__ADD_REG_DE
	/* LD REGISTER,nnnn */
	opcodesMap[SHIFT_0xFD+0x21] = instrFD__LD_REG_NNNN
	/* LD (nnnn),REGISTER */
	opcodesMap[SHIFT_0xFD+0x22] = instrFD__LD_iNNNN_REG
	/* INC REGISTER */
	opcodesMap[SHIFT_0xFD+0x23] = instrFD__INC_REG
	/* INC REGISTERH */
	opcodesMap[SHIFT_0xFD+0x24] = instrFD__INC_REGH
	/* DEC REGISTERH */
	opcodesMap[SHIFT_0xFD+0x25] = instrFD__DEC_REGH
	/* LD REGISTERH,nn */
	opcodesMap[SHIFT_0xFD+0x26] = instrFD__LD_REGH_NN
	/* ADD REGISTER,REGISTER */
	opcodesMap[SHIFT_0xFD+0x29] = instrFD__ADD_REG_REG
	/* LD REGISTER,(nnnn) */
	opcodesMap[SHIFT_0xFD+0x2a] = instrFD__LD_REG_iNNNN
	/* DEC REGISTER */
	opcodesMap[SHIFT_0xFD+0x2b] = instrFD__DEC_REG
	/* INC REGISTERL */
	opcodesMap[SHIFT_0xFD+0x2c] = instrFD__INC_REGL
	/* DEC REGISTERL */
	opcodesMap[SHIFT_0xFD+0x2d] = instrFD__DEC_REGL
	/* LD REGISTERL,nn */
	opcodesMap[SHIFT_0xFD+0x2e] = instrFD__LD_REGL_NN
	/* INC (REGISTER+dd) */
	opcodesMap[SHIFT_0xFD+0x34] = instrFD__INC_iREGpDD
	/* DEC (REGISTER+dd) */
	opcodesMap[SHIFT_0xFD+0x35] = instrFD__DEC_iREGpDD
	/* LD (REGISTER+dd),nn */
	opcodesMap[SHIFT_0xFD+0x36] = instrFD__LD_iREGpDD_NN
	/* ADD REGISTER,SP */
	opcodesMap[SHIFT_0xFD+0x39] = instrFD__ADD_REG_SP
	/* LD B,REGISTERH */
	opcodesMap[SHIFT_0xFD+0x44] = instrFD__LD_B_REGH
	/* LD B,REGISTERL */
	opcodesMap[SHIFT_0xFD+0x45] = instrFD__LD_B_REGL
	/* LD B,(REGISTER+dd) */
	opcodesMap[SHIFT_0xFD+0x46] = instrFD__LD_B_iREGpDD
	/* LD C,REGISTERH */
	opcodesMap[SHIFT_0xFD+0x4c] = instrFD__LD_C_REGH
	/* LD C,REGISTERL */
	opcodesMap[SHIFT_0xFD+0x4d] = instrFD__LD_C_REGL
	/* LD C,(REGISTER+dd) */
	opcodesMap[SHIFT_0xFD+0x4e] = instrFD__LD_C_iREGpDD
	/* LD D,REGISTERH */
	opcodesMap[SHIFT_0xFD+0x54] = instrFD__LD_D_REGH
	/* LD D,REGISTERL */
	opcodesMap[SHIFT_0xFD+0x55] = instrFD__LD_D_REGL
	/* LD D,(REGISTER+dd) */
	opcodesMap[SHIFT_0xFD+0x56] = instrFD__LD_D_iREGpDD
	/* LD E,REGISTERH */
	opcodesMap[SHIFT_0xFD+0x5c] = instrFD__LD_E_REGH
	/* LD E,REGISTERL */
	opcodesMap[SHIFT_0xFD+0x5d] = instrFD__LD_E_REGL
	/* LD E,(REGISTER+dd) */
	opcodesMap[SHIFT_0xFD+0x5e] = instrFD__LD_E_iREGpDD
	/* LD REGISTERH,B */
	opcodesMap[SHIFT_0xFD+0x60] = instrFD__LD_REGH_B
	/* LD REGISTERH,C */
	opcodesMap[SHIFT_0xFD+0x61] = instrFD__LD_REGH_C
	/* LD REGISTERH,D */
	opcodesMap[SHIFT_0xFD+0x62] = instrFD__LD_REGH_D
	/* LD REGISTERH,E */
	opcodesMap[SHIFT_0xFD+0x63] = instrFD__LD_REGH_E
	/* LD REGISTERH,REGISTERH */
	opcodesMap[SHIFT_0xFD+0x64] = instrFD__LD_REGH_REGH
	/* LD REGISTERH,REGISTERL */
	opcodesMap[SHIFT_0xFD+0x65] = instrFD__LD_REGH_REGL
	/* LD H,(REGISTER+dd) */
	opcodesMap[SHIFT_0xFD+0x66] = instrFD__LD_H_iREGpDD
	/* LD REGISTERH,A */
	opcodesMap[SHIFT_0xFD+0x67] = instrFD__LD_REGH_A
	/* LD REGISTERL,B */
	opcodesMap[SHIFT_0xFD+0x68] = instrFD__LD_REGL_B
	/* LD REGISTERL,C */
	opcodesMap[SHIFT_0xFD+0x69] = instrFD__LD_REGL_C
	/* LD REGISTERL,D */
	opcodesMap[SHIFT_0xFD+0x6a] = instrFD__LD_REGL_D
	/* LD REGISTERL,E */
	opcodesMap[SHIFT_0xFD+0x6b] = instrFD__LD_REGL_E
	/* LD REGISTERL,REGISTERH */
	opcodesMap[SHIFT_0xFD+0x6c] = instrFD__LD_REGL_REGH
	/* LD REGISTERL,REGISTERL */
	opcodesMap[SHIFT_0xFD+0x6d] = instrFD__LD_REGL_REGL
	/* LD L,(REGISTER+dd) */
	opcodesMap[SHIFT_0xFD+0x6e] = instrFD__LD_L_iREGpDD
	/* LD REGISTERL,A */
	opcodesMap[SHIFT_0xFD+0x6f] = instrFD__LD_REGL_A
	/* LD (REGISTER+dd),B */
	opcodesMap[SHIFT_0xFD+0x70] = instrFD__LD_iREGpDD_B
	/* LD (REGISTER+dd),C */
	opcodesMap[SHIFT_0xFD+0x71] = instrFD__LD_iREGpDD_C
	/* LD (REGISTER+dd),D */
	opcodesMap[SHIFT_0xFD+0x72] = instrFD__LD_iREGpDD_D
	/* LD (REGISTER+dd),E */
	opcodesMap[SHIFT_0xFD+0x73] = instrFD__LD_iREGpDD_E
	/* LD (REGISTER+dd),H */
	opcodesMap[SHIFT_0xFD+0x74] = instrFD__LD_iREGpDD_H
	/* LD (REGISTER+dd),L */
	opcodesMap[SHIFT_0xFD+0x75] = instrFD__LD_iREGpDD_L
	/* LD (REGISTER+dd),A */
	opcodesMap[SHIFT_0xFD+0x77] = instrFD__LD_iREGpDD_A
	/* LD A,REGISTERH */
	opcodesMap[SHIFT_0xFD+0x7c] = instrFD__LD_A_REGH
	/* LD A,REGISTERL */
	opcodesMap[SHIFT_0xFD+0x7d] = instrFD__LD_A_REGL
	/* LD A,(REGISTER+dd) */
	opcodesMap[SHIFT_0xFD+0x7e] = instrFD__LD_A_iREGpDD
	/* ADD A,REGISTERH */
	opcodesMap[SHIFT_0xFD+0x84] = instrFD__ADD_A_REGH
	/* ADD A,REGISTERL */
	opcodesMap[SHIFT_0xFD+0x85] = instrFD__ADD_A_REGL
	/* ADD A,(REGISTER+dd) */
	opcodesMap[SHIFT_0xFD+0x86] = instrFD__ADD_A_iREGpDD
	/* ADC A,REGISTERH */
	opcodesMap[SHIFT_0xFD+0x8c] = instrFD__ADC_A_REGH
	/* ADC A,REGISTERL */
	opcodesMap[SHIFT_0xFD+0x8d] = instrFD__ADC_A_REGL
	/* ADC A,(REGISTER+dd) */
	opcodesMap[SHIFT_0xFD+0x8e] = instrFD__ADC_A_iREGpDD
	/* SUB A,REGISTERH */
	opcodesMap[SHIFT_0xFD+0x94] = instrFD__SUB_A_REGH
	/* SUB A,REGISTERL */
	opcodesMap[SHIFT_0xFD+0x95] = instrFD__SUB_A_REGL
	/* SUB A,(REGISTER+dd) */
	opcodesMap[SHIFT_0xFD+0x96] = instrFD__SUB_A_iREGpDD
	/* SBC A,REGISTERH */
	opcodesMap[SHIFT_0xFD+0x9c] = instrFD__SBC_A_REGH
	/* SBC A,REGISTERL */
	opcodesMap[SHIFT_0xFD+0x9d] = instrFD__SBC_A_REGL
	/* SBC A,(REGISTER+dd) */
	opcodesMap[SHIFT_0xFD+0x9e] = instrFD__SBC_A_iREGpDD
	/* AND A,REGISTERH */
	opcodesMap[SHIFT_0xFD+0xa4] = instrFD__AND_A_REGH
	/* AND A,REGISTERL */
	opcodesMap[SHIFT_0xFD+0xa5] = instrFD__AND_A_REGL
	/* AND A,(REGISTER+dd) */
	opcodesMap[SHIFT_0xFD+0xa6] = instrFD__AND_A_iREGpDD
	/* XOR A,REGISTERH */
	opcodesMap[SHIFT_0xFD+0xac] = instrFD__XOR_A_REGH
	/* XOR A,REGISTERL */
	opcodesMap[SHIFT_0xFD+0xad] = instrFD__XOR_A_REGL
	/* XOR A,(REGISTER+dd) */
	opcodesMap[SHIFT_0xFD+0xae] = instrFD__XOR_A_iREGpDD
	/* OR A,REGISTERH */
	opcodesMap[SHIFT_0xFD+0xb4] = instrFD__OR_A_REGH
	/* OR A,REGISTERL */
	opcodesMap[SHIFT_0xFD+0xb5] = instrFD__OR_A_REGL
	/* OR A,(REGISTER+dd) */
	opcodesMap[SHIFT_0xFD+0xb6] = instrFD__OR_A_iREGpDD
	/* CP A,REGISTERH */
	opcodesMap[SHIFT_0xFD+0xbc] = instrFD__CP_A_REGH
	/* CP A,REGISTERL */
	opcodesMap[SHIFT_0xFD+0xbd] = instrFD__CP_A_REGL
	/* CP A,(REGISTER+dd) */
	opcodesMap[SHIFT_0xFD+0xbe] = instrFD__CP_A_iREGpDD
	/* shift DDFDCB */
	opcodesMap[SHIFT_0xFD+0xcb] = instrFD__SHIFT_DDFDCB
	/* POP REGISTER */
	opcodesMap[SHIFT_0xFD+0xe1] = instrFD__POP_REG
	/* EX (SP),REGISTER */
	opcodesMap[SHIFT_0xFD+0xe3] = instrFD__EX_iSP_REG
	/* PUSH REGISTER */
	opcodesMap[SHIFT_0xFD+0xe5] = instrFD__PUSH_REG
	/* JP REGISTER */
	opcodesMap[SHIFT_0xFD+0xe9] = instrFD__JP_REG
	/* LD SP,REGISTER */
	opcodesMap[SHIFT_0xFD+0xf9] = instrFD__LD_SP_REG

	// END of 0xfd shifted opcodes

	// BEGIN of 0xddfdcb shifted opcodes
	/* LD B,RLC (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x00] = instrDDCB__LD_B_RLC_iREGpDD
	/* LD C,RLC (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x01] = instrDDCB__LD_C_RLC_iREGpDD
	/* LD D,RLC (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x02] = instrDDCB__LD_D_RLC_iREGpDD
	/* LD E,RLC (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x03] = instrDDCB__LD_E_RLC_iREGpDD
	/* LD H,RLC (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x04] = instrDDCB__LD_H_RLC_iREGpDD
	/* LD L,RLC (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x05] = instrDDCB__LD_L_RLC_iREGpDD
	/* RLC (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x06] = instrDDCB__RLC_iREGpDD
	/* LD A,RLC (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x07] = instrDDCB__LD_A_RLC_iREGpDD
	/* LD B,RRC (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x08] = instrDDCB__LD_B_RRC_iREGpDD
	/* LD C,RRC (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x09] = instrDDCB__LD_C_RRC_iREGpDD
	/* LD D,RRC (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x0a] = instrDDCB__LD_D_RRC_iREGpDD
	/* LD E,RRC (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x0b] = instrDDCB__LD_E_RRC_iREGpDD
	/* LD H,RRC (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x0c] = instrDDCB__LD_H_RRC_iREGpDD
	/* LD L,RRC (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x0d] = instrDDCB__LD_L_RRC_iREGpDD
	/* RRC (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x0e] = instrDDCB__RRC_iREGpDD
	/* LD A,RRC (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x0f] = instrDDCB__LD_A_RRC_iREGpDD
	/* LD B,RL (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x10] = instrDDCB__LD_B_RL_iREGpDD
	/* LD C,RL (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x11] = instrDDCB__LD_C_RL_iREGpDD
	/* LD D,RL (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x12] = instrDDCB__LD_D_RL_iREGpDD
	/* LD E,RL (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x13] = instrDDCB__LD_E_RL_iREGpDD
	/* LD H,RL (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x14] = instrDDCB__LD_H_RL_iREGpDD
	/* LD L,RL (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x15] = instrDDCB__LD_L_RL_iREGpDD
	/* RL (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x16] = instrDDCB__RL_iREGpDD
	/* LD A,RL (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x17] = instrDDCB__LD_A_RL_iREGpDD
	/* LD B,RR (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x18] = instrDDCB__LD_B_RR_iREGpDD
	/* LD C,RR (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x19] = instrDDCB__LD_C_RR_iREGpDD
	/* LD D,RR (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x1a] = instrDDCB__LD_D_RR_iREGpDD
	/* LD E,RR (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x1b] = instrDDCB__LD_E_RR_iREGpDD
	/* LD H,RR (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x1c] = instrDDCB__LD_H_RR_iREGpDD
	/* LD L,RR (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x1d] = instrDDCB__LD_L_RR_iREGpDD
	/* RR (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x1e] = instrDDCB__RR_iREGpDD
	/* LD A,RR (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x1f] = instrDDCB__LD_A_RR_iREGpDD
	/* LD B,SLA (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x20] = instrDDCB__LD_B_SLA_iREGpDD
	/* LD C,SLA (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x21] = instrDDCB__LD_C_SLA_iREGpDD
	/* LD D,SLA (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x22] = instrDDCB__LD_D_SLA_iREGpDD
	/* LD E,SLA (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x23] = instrDDCB__LD_E_SLA_iREGpDD
	/* LD H,SLA (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x24] = instrDDCB__LD_H_SLA_iREGpDD
	/* LD L,SLA (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x25] = instrDDCB__LD_L_SLA_iREGpDD
	/* SLA (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x26] = instrDDCB__SLA_iREGpDD
	/* LD A,SLA (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x27] = instrDDCB__LD_A_SLA_iREGpDD
	/* LD B,SRA (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x28] = instrDDCB__LD_B_SRA_iREGpDD
	/* LD C,SRA (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x29] = instrDDCB__LD_C_SRA_iREGpDD
	/* LD D,SRA (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x2a] = instrDDCB__LD_D_SRA_iREGpDD
	/* LD E,SRA (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x2b] = instrDDCB__LD_E_SRA_iREGpDD
	/* LD H,SRA (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x2c] = instrDDCB__LD_H_SRA_iREGpDD
	/* LD L,SRA (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x2d] = instrDDCB__LD_L_SRA_iREGpDD
	/* SRA (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x2e] = instrDDCB__SRA_iREGpDD
	/* LD A,SRA (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x2f] = instrDDCB__LD_A_SRA_iREGpDD
	/* LD B,SLL (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x30] = instrDDCB__LD_B_SLL_iREGpDD
	/* LD C,SLL (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x31] = instrDDCB__LD_C_SLL_iREGpDD
	/* LD D,SLL (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x32] = instrDDCB__LD_D_SLL_iREGpDD
	/* LD E,SLL (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x33] = instrDDCB__LD_E_SLL_iREGpDD
	/* LD H,SLL (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x34] = instrDDCB__LD_H_SLL_iREGpDD
	/* LD L,SLL (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x35] = instrDDCB__LD_L_SLL_iREGpDD
	/* SLL (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x36] = instrDDCB__SLL_iREGpDD
	/* LD A,SLL (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x37] = instrDDCB__LD_A_SLL_iREGpDD
	/* LD B,SRL (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x38] = instrDDCB__LD_B_SRL_iREGpDD
	/* LD C,SRL (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x39] = instrDDCB__LD_C_SRL_iREGpDD
	/* LD D,SRL (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x3a] = instrDDCB__LD_D_SRL_iREGpDD
	/* LD E,SRL (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x3b] = instrDDCB__LD_E_SRL_iREGpDD
	/* LD H,SRL (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x3c] = instrDDCB__LD_H_SRL_iREGpDD
	/* LD L,SRL (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x3d] = instrDDCB__LD_L_SRL_iREGpDD
	/* SRL (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x3e] = instrDDCB__SRL_iREGpDD
	/* LD A,SRL (REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x3f] = instrDDCB__LD_A_SRL_iREGpDD
	/* BIT 0,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x47] = instrDDCB__BIT_0_iREGpDD
	// Fallthrough cases
	opcodesMap[SHIFT_0xDDCB+0x40] = opcodesMap[SHIFT_0xDDCB+0x47]
	opcodesMap[SHIFT_0xDDCB+0x41] = opcodesMap[SHIFT_0xDDCB+0x47]
	opcodesMap[SHIFT_0xDDCB+0x42] = opcodesMap[SHIFT_0xDDCB+0x47]
	opcodesMap[SHIFT_0xDDCB+0x43] = opcodesMap[SHIFT_0xDDCB+0x47]
	opcodesMap[SHIFT_0xDDCB+0x44] = opcodesMap[SHIFT_0xDDCB+0x47]
	opcodesMap[SHIFT_0xDDCB+0x45] = opcodesMap[SHIFT_0xDDCB+0x47]
	opcodesMap[SHIFT_0xDDCB+0x46] = opcodesMap[SHIFT_0xDDCB+0x47]
	/* BIT 1,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x4f] = instrDDCB__BIT_1_iREGpDD
	// Fallthrough cases
	opcodesMap[SHIFT_0xDDCB+0x48] = opcodesMap[SHIFT_0xDDCB+0x4f]
	opcodesMap[SHIFT_0xDDCB+0x49] = opcodesMap[SHIFT_0xDDCB+0x4f]
	opcodesMap[SHIFT_0xDDCB+0x4a] = opcodesMap[SHIFT_0xDDCB+0x4f]
	opcodesMap[SHIFT_0xDDCB+0x4b] = opcodesMap[SHIFT_0xDDCB+0x4f]
	opcodesMap[SHIFT_0xDDCB+0x4c] = opcodesMap[SHIFT_0xDDCB+0x4f]
	opcodesMap[SHIFT_0xDDCB+0x4d] = opcodesMap[SHIFT_0xDDCB+0x4f]
	opcodesMap[SHIFT_0xDDCB+0x4e] = opcodesMap[SHIFT_0xDDCB+0x4f]
	/* BIT 2,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x57] = instrDDCB__BIT_2_iREGpDD
	// Fallthrough cases
	opcodesMap[SHIFT_0xDDCB+0x50] = opcodesMap[SHIFT_0xDDCB+0x57]
	opcodesMap[SHIFT_0xDDCB+0x51] = opcodesMap[SHIFT_0xDDCB+0x57]
	opcodesMap[SHIFT_0xDDCB+0x52] = opcodesMap[SHIFT_0xDDCB+0x57]
	opcodesMap[SHIFT_0xDDCB+0x53] = opcodesMap[SHIFT_0xDDCB+0x57]
	opcodesMap[SHIFT_0xDDCB+0x54] = opcodesMap[SHIFT_0xDDCB+0x57]
	opcodesMap[SHIFT_0xDDCB+0x55] = opcodesMap[SHIFT_0xDDCB+0x57]
	opcodesMap[SHIFT_0xDDCB+0x56] = opcodesMap[SHIFT_0xDDCB+0x57]
	/* BIT 3,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x5f] = instrDDCB__BIT_3_iREGpDD
	// Fallthrough cases
	opcodesMap[SHIFT_0xDDCB+0x58] = opcodesMap[SHIFT_0xDDCB+0x5f]
	opcodesMap[SHIFT_0xDDCB+0x59] = opcodesMap[SHIFT_0xDDCB+0x5f]
	opcodesMap[SHIFT_0xDDCB+0x5a] = opcodesMap[SHIFT_0xDDCB+0x5f]
	opcodesMap[SHIFT_0xDDCB+0x5b] = opcodesMap[SHIFT_0xDDCB+0x5f]
	opcodesMap[SHIFT_0xDDCB+0x5c] = opcodesMap[SHIFT_0xDDCB+0x5f]
	opcodesMap[SHIFT_0xDDCB+0x5d] = opcodesMap[SHIFT_0xDDCB+0x5f]
	opcodesMap[SHIFT_0xDDCB+0x5e] = opcodesMap[SHIFT_0xDDCB+0x5f]
	/* BIT 4,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x67] = instrDDCB__BIT_4_iREGpDD
	// Fallthrough cases
	opcodesMap[SHIFT_0xDDCB+0x60] = opcodesMap[SHIFT_0xDDCB+0x67]
	opcodesMap[SHIFT_0xDDCB+0x61] = opcodesMap[SHIFT_0xDDCB+0x67]
	opcodesMap[SHIFT_0xDDCB+0x62] = opcodesMap[SHIFT_0xDDCB+0x67]
	opcodesMap[SHIFT_0xDDCB+0x63] = opcodesMap[SHIFT_0xDDCB+0x67]
	opcodesMap[SHIFT_0xDDCB+0x64] = opcodesMap[SHIFT_0xDDCB+0x67]
	opcodesMap[SHIFT_0xDDCB+0x65] = opcodesMap[SHIFT_0xDDCB+0x67]
	opcodesMap[SHIFT_0xDDCB+0x66] = opcodesMap[SHIFT_0xDDCB+0x67]
	/* BIT 5,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x6f] = instrDDCB__BIT_5_iREGpDD
	// Fallthrough cases
	opcodesMap[SHIFT_0xDDCB+0x68] = opcodesMap[SHIFT_0xDDCB+0x6f]
	opcodesMap[SHIFT_0xDDCB+0x69] = opcodesMap[SHIFT_0xDDCB+0x6f]
	opcodesMap[SHIFT_0xDDCB+0x6a] = opcodesMap[SHIFT_0xDDCB+0x6f]
	opcodesMap[SHIFT_0xDDCB+0x6b] = opcodesMap[SHIFT_0xDDCB+0x6f]
	opcodesMap[SHIFT_0xDDCB+0x6c] = opcodesMap[SHIFT_0xDDCB+0x6f]
	opcodesMap[SHIFT_0xDDCB+0x6d] = opcodesMap[SHIFT_0xDDCB+0x6f]
	opcodesMap[SHIFT_0xDDCB+0x6e] = opcodesMap[SHIFT_0xDDCB+0x6f]
	/* BIT 6,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x77] = instrDDCB__BIT_6_iREGpDD
	// Fallthrough cases
	opcodesMap[SHIFT_0xDDCB+0x70] = opcodesMap[SHIFT_0xDDCB+0x77]
	opcodesMap[SHIFT_0xDDCB+0x71] = opcodesMap[SHIFT_0xDDCB+0x77]
	opcodesMap[SHIFT_0xDDCB+0x72] = opcodesMap[SHIFT_0xDDCB+0x77]
	opcodesMap[SHIFT_0xDDCB+0x73] = opcodesMap[SHIFT_0xDDCB+0x77]
	opcodesMap[SHIFT_0xDDCB+0x74] = opcodesMap[SHIFT_0xDDCB+0x77]
	opcodesMap[SHIFT_0xDDCB+0x75] = opcodesMap[SHIFT_0xDDCB+0x77]
	opcodesMap[SHIFT_0xDDCB+0x76] = opcodesMap[SHIFT_0xDDCB+0x77]
	/* BIT 7,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x7f] = instrDDCB__BIT_7_iREGpDD
	// Fallthrough cases
	opcodesMap[SHIFT_0xDDCB+0x78] = opcodesMap[SHIFT_0xDDCB+0x7f]
	opcodesMap[SHIFT_0xDDCB+0x79] = opcodesMap[SHIFT_0xDDCB+0x7f]
	opcodesMap[SHIFT_0xDDCB+0x7a] = opcodesMap[SHIFT_0xDDCB+0x7f]
	opcodesMap[SHIFT_0xDDCB+0x7b] = opcodesMap[SHIFT_0xDDCB+0x7f]
	opcodesMap[SHIFT_0xDDCB+0x7c] = opcodesMap[SHIFT_0xDDCB+0x7f]
	opcodesMap[SHIFT_0xDDCB+0x7d] = opcodesMap[SHIFT_0xDDCB+0x7f]
	opcodesMap[SHIFT_0xDDCB+0x7e] = opcodesMap[SHIFT_0xDDCB+0x7f]
	/* LD B,RES 0,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x80] = instrDDCB__LD_B_RES_0_iREGpDD
	/* LD C,RES 0,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x81] = instrDDCB__LD_C_RES_0_iREGpDD
	/* LD D,RES 0,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x82] = instrDDCB__LD_D_RES_0_iREGpDD
	/* LD E,RES 0,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x83] = instrDDCB__LD_E_RES_0_iREGpDD
	/* LD H,RES 0,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x84] = instrDDCB__LD_H_RES_0_iREGpDD
	/* LD L,RES 0,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x85] = instrDDCB__LD_L_RES_0_iREGpDD
	/* RES 0,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x86] = instrDDCB__RES_0_iREGpDD
	/* LD A,RES 0,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x87] = instrDDCB__LD_A_RES_0_iREGpDD
	/* LD B,RES 1,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x88] = instrDDCB__LD_B_RES_1_iREGpDD
	/* LD C,RES 1,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x89] = instrDDCB__LD_C_RES_1_iREGpDD
	/* LD D,RES 1,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x8a] = instrDDCB__LD_D_RES_1_iREGpDD
	/* LD E,RES 1,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x8b] = instrDDCB__LD_E_RES_1_iREGpDD
	/* LD H,RES 1,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x8c] = instrDDCB__LD_H_RES_1_iREGpDD
	/* LD L,RES 1,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x8d] = instrDDCB__LD_L_RES_1_iREGpDD
	/* RES 1,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x8e] = instrDDCB__RES_1_iREGpDD
	/* LD A,RES 1,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x8f] = instrDDCB__LD_A_RES_1_iREGpDD
	/* LD B,RES 2,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x90] = instrDDCB__LD_B_RES_2_iREGpDD
	/* LD C,RES 2,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x91] = instrDDCB__LD_C_RES_2_iREGpDD
	/* LD D,RES 2,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x92] = instrDDCB__LD_D_RES_2_iREGpDD
	/* LD E,RES 2,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x93] = instrDDCB__LD_E_RES_2_iREGpDD
	/* LD H,RES 2,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x94] = instrDDCB__LD_H_RES_2_iREGpDD
	/* LD L,RES 2,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x95] = instrDDCB__LD_L_RES_2_iREGpDD
	/* RES 2,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x96] = instrDDCB__RES_2_iREGpDD
	/* LD A,RES 2,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x97] = instrDDCB__LD_A_RES_2_iREGpDD
	/* LD B,RES 3,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x98] = instrDDCB__LD_B_RES_3_iREGpDD
	/* LD C,RES 3,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x99] = instrDDCB__LD_C_RES_3_iREGpDD
	/* LD D,RES 3,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x9a] = instrDDCB__LD_D_RES_3_iREGpDD
	/* LD E,RES 3,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x9b] = instrDDCB__LD_E_RES_3_iREGpDD
	/* LD H,RES 3,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x9c] = instrDDCB__LD_H_RES_3_iREGpDD
	/* LD L,RES 3,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x9d] = instrDDCB__LD_L_RES_3_iREGpDD
	/* RES 3,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x9e] = instrDDCB__RES_3_iREGpDD
	/* LD A,RES 3,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0x9f] = instrDDCB__LD_A_RES_3_iREGpDD
	/* LD B,RES 4,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xa0] = instrDDCB__LD_B_RES_4_iREGpDD
	/* LD C,RES 4,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xa1] = instrDDCB__LD_C_RES_4_iREGpDD
	/* LD D,RES 4,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xa2] = instrDDCB__LD_D_RES_4_iREGpDD
	/* LD E,RES 4,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xa3] = instrDDCB__LD_E_RES_4_iREGpDD
	/* LD H,RES 4,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xa4] = instrDDCB__LD_H_RES_4_iREGpDD
	/* LD L,RES 4,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xa5] = instrDDCB__LD_L_RES_4_iREGpDD
	/* RES 4,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xa6] = instrDDCB__RES_4_iREGpDD
	/* LD A,RES 4,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xa7] = instrDDCB__LD_A_RES_4_iREGpDD
	/* LD B,RES 5,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xa8] = instrDDCB__LD_B_RES_5_iREGpDD
	/* LD C,RES 5,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xa9] = instrDDCB__LD_C_RES_5_iREGpDD
	/* LD D,RES 5,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xaa] = instrDDCB__LD_D_RES_5_iREGpDD
	/* LD E,RES 5,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xab] = instrDDCB__LD_E_RES_5_iREGpDD
	/* LD H,RES 5,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xac] = instrDDCB__LD_H_RES_5_iREGpDD
	/* LD L,RES 5,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xad] = instrDDCB__LD_L_RES_5_iREGpDD
	/* RES 5,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xae] = instrDDCB__RES_5_iREGpDD
	/* LD A,RES 5,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xaf] = instrDDCB__LD_A_RES_5_iREGpDD
	/* LD B,RES 6,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xb0] = instrDDCB__LD_B_RES_6_iREGpDD
	/* LD C,RES 6,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xb1] = instrDDCB__LD_C_RES_6_iREGpDD
	/* LD D,RES 6,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xb2] = instrDDCB__LD_D_RES_6_iREGpDD
	/* LD E,RES 6,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xb3] = instrDDCB__LD_E_RES_6_iREGpDD
	/* LD H,RES 6,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xb4] = instrDDCB__LD_H_RES_6_iREGpDD
	/* LD L,RES 6,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xb5] = instrDDCB__LD_L_RES_6_iREGpDD
	/* RES 6,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xb6] = instrDDCB__RES_6_iREGpDD
	/* LD A,RES 6,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xb7] = instrDDCB__LD_A_RES_6_iREGpDD
	/* LD B,RES 7,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xb8] = instrDDCB__LD_B_RES_7_iREGpDD
	/* LD C,RES 7,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xb9] = instrDDCB__LD_C_RES_7_iREGpDD
	/* LD D,RES 7,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xba] = instrDDCB__LD_D_RES_7_iREGpDD
	/* LD E,RES 7,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xbb] = instrDDCB__LD_E_RES_7_iREGpDD
	/* LD H,RES 7,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xbc] = instrDDCB__LD_H_RES_7_iREGpDD
	/* LD L,RES 7,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xbd] = instrDDCB__LD_L_RES_7_iREGpDD
	/* RES 7,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xbe] = instrDDCB__RES_7_iREGpDD
	/* LD A,RES 7,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xbf] = instrDDCB__LD_A_RES_7_iREGpDD
	/* LD B,SET 0,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xc0] = instrDDCB__LD_B_SET_0_iREGpDD
	/* LD C,SET 0,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xc1] = instrDDCB__LD_C_SET_0_iREGpDD
	/* LD D,SET 0,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xc2] = instrDDCB__LD_D_SET_0_iREGpDD
	/* LD E,SET 0,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xc3] = instrDDCB__LD_E_SET_0_iREGpDD
	/* LD H,SET 0,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xc4] = instrDDCB__LD_H_SET_0_iREGpDD
	/* LD L,SET 0,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xc5] = instrDDCB__LD_L_SET_0_iREGpDD
	/* SET 0,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xc6] = instrDDCB__SET_0_iREGpDD
	/* LD A,SET 0,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xc7] = instrDDCB__LD_A_SET_0_iREGpDD
	/* LD B,SET 1,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xc8] = instrDDCB__LD_B_SET_1_iREGpDD
	/* LD C,SET 1,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xc9] = instrDDCB__LD_C_SET_1_iREGpDD
	/* LD D,SET 1,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xca] = instrDDCB__LD_D_SET_1_iREGpDD
	/* LD E,SET 1,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xcb] = instrDDCB__LD_E_SET_1_iREGpDD
	/* LD H,SET 1,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xcc] = instrDDCB__LD_H_SET_1_iREGpDD
	/* LD L,SET 1,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xcd] = instrDDCB__LD_L_SET_1_iREGpDD
	/* SET 1,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xce] = instrDDCB__SET_1_iREGpDD
	/* LD A,SET 1,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xcf] = instrDDCB__LD_A_SET_1_iREGpDD
	/* LD B,SET 2,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xd0] = instrDDCB__LD_B_SET_2_iREGpDD
	/* LD C,SET 2,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xd1] = instrDDCB__LD_C_SET_2_iREGpDD
	/* LD D,SET 2,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xd2] = instrDDCB__LD_D_SET_2_iREGpDD
	/* LD E,SET 2,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xd3] = instrDDCB__LD_E_SET_2_iREGpDD
	/* LD H,SET 2,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xd4] = instrDDCB__LD_H_SET_2_iREGpDD
	/* LD L,SET 2,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xd5] = instrDDCB__LD_L_SET_2_iREGpDD
	/* SET 2,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xd6] = instrDDCB__SET_2_iREGpDD
	/* LD A,SET 2,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xd7] = instrDDCB__LD_A_SET_2_iREGpDD
	/* LD B,SET 3,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xd8] = instrDDCB__LD_B_SET_3_iREGpDD
	/* LD C,SET 3,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xd9] = instrDDCB__LD_C_SET_3_iREGpDD
	/* LD D,SET 3,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xda] = instrDDCB__LD_D_SET_3_iREGpDD
	/* LD E,SET 3,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xdb] = instrDDCB__LD_E_SET_3_iREGpDD
	/* LD H,SET 3,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xdc] = instrDDCB__LD_H_SET_3_iREGpDD
	/* LD L,SET 3,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xdd] = instrDDCB__LD_L_SET_3_iREGpDD
	/* SET 3,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xde] = instrDDCB__SET_3_iREGpDD
	/* LD A,SET 3,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xdf] = instrDDCB__LD_A_SET_3_iREGpDD
	/* LD B,SET 4,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xe0] = instrDDCB__LD_B_SET_4_iREGpDD
	/* LD C,SET 4,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xe1] = instrDDCB__LD_C_SET_4_iREGpDD
	/* LD D,SET 4,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xe2] = instrDDCB__LD_D_SET_4_iREGpDD
	/* LD E,SET 4,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xe3] = instrDDCB__LD_E_SET_4_iREGpDD
	/* LD H,SET 4,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xe4] = instrDDCB__LD_H_SET_4_iREGpDD
	/* LD L,SET 4,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xe5] = instrDDCB__LD_L_SET_4_iREGpDD
	/* SET 4,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xe6] = instrDDCB__SET_4_iREGpDD
	/* LD A,SET 4,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xe7] = instrDDCB__LD_A_SET_4_iREGpDD
	/* LD B,SET 5,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xe8] = instrDDCB__LD_B_SET_5_iREGpDD
	/* LD C,SET 5,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xe9] = instrDDCB__LD_C_SET_5_iREGpDD
	/* LD D,SET 5,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xea] = instrDDCB__LD_D_SET_5_iREGpDD
	/* LD E,SET 5,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xeb] = instrDDCB__LD_E_SET_5_iREGpDD
	/* LD H,SET 5,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xec] = instrDDCB__LD_H_SET_5_iREGpDD
	/* LD L,SET 5,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xed] = instrDDCB__LD_L_SET_5_iREGpDD
	/* SET 5,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xee] = instrDDCB__SET_5_iREGpDD
	/* LD A,SET 5,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xef] = instrDDCB__LD_A_SET_5_iREGpDD
	/* LD B,SET 6,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xf0] = instrDDCB__LD_B_SET_6_iREGpDD
	/* LD C,SET 6,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xf1] = instrDDCB__LD_C_SET_6_iREGpDD
	/* LD D,SET 6,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xf2] = instrDDCB__LD_D_SET_6_iREGpDD
	/* LD E,SET 6,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xf3] = instrDDCB__LD_E_SET_6_iREGpDD
	/* LD H,SET 6,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xf4] = instrDDCB__LD_H_SET_6_iREGpDD
	/* LD L,SET 6,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xf5] = instrDDCB__LD_L_SET_6_iREGpDD
	/* SET 6,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xf6] = instrDDCB__SET_6_iREGpDD
	/* LD A,SET 6,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xf7] = instrDDCB__LD_A_SET_6_iREGpDD
	/* LD B,SET 7,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xf8] = instrDDCB__LD_B_SET_7_iREGpDD
	/* LD C,SET 7,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xf9] = instrDDCB__LD_C_SET_7_iREGpDD
	/* LD D,SET 7,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xfa] = instrDDCB__LD_D_SET_7_iREGpDD
	/* LD E,SET 7,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xfb] = instrDDCB__LD_E_SET_7_iREGpDD
	/* LD H,SET 7,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xfc] = instrDDCB__LD_H_SET_7_iREGpDD
	/* LD L,SET 7,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xfd] = instrDDCB__LD_L_SET_7_iREGpDD
	/* SET 7,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xfe] = instrDDCB__SET_7_iREGpDD
	/* LD A,SET 7,(REGISTER+dd) */
	opcodesMap[SHIFT_0xDDCB+0xff] = instrDDCB__LD_A_SET_7_iREGpDD

	// END of 0xddfdcb shifted opcodes
}

/* NOP */
func instr__NOP(z80 *Z80) {
}
/* LD BC,nnnn */
func instr__LD_BC_NNNN(z80 *Z80) {
	b1 := z80.memory.readByte(z80.pc)
	z80.pc++
	b2 := z80.memory.readByte(z80.pc)
	z80.pc++
	z80.setBC(joinBytes(b2, b1))
}
/* LD (BC),A */
func instr__LD_iBC_A(z80 *Z80) {
	z80.memory.writeByte(z80.BC(), z80.a)
}
/* INC BC */
func instr__INC_BC(z80 *Z80) {
	z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 2)
	z80.incBC()
}
/* INC B */
func instr__INC_B(z80 *Z80) {
	z80.incB()
}
/* DEC B */
func instr__DEC_B(z80 *Z80) {
	z80.decB()
}
/* LD B,nn */
func instr__LD_B_NN(z80 *Z80) {
	z80.b = z80.memory.readByte(z80.pc)
	z80.pc++
}
/* RLCA */
func instr__RLCA(z80 *Z80) {
	z80.a = (z80.a << 1) | (z80.a >> 7)
	z80.f = (z80.f & (FLAG_P | FLAG_Z | FLAG_S)) |
		(z80.a & (FLAG_C | FLAG_3 | FLAG_5))
}
/* EX AF,AF' */
func instr__EX_AF_AF(z80 *Z80) {
	/* Tape saving trap: note this traps the EX AF,AF' at #04d0, not
	   #04d1 as PC has already been incremented */
	/* 0x76 - Timex 2068 save routine in EXROM */
	if z80.pc == 0x04d1 || z80.pc == 0x0077 {
		if z80.tapeSaveTrap() == 0 { /*break*/
		}
	}

	var olda, oldf = z80.a, z80.f
	z80.a = z80.a_
	z80.f = z80.f_
	z80.a_ = olda
	z80.f_ = oldf
}
/* ADD HL,BC */
func instr__ADD_HL_BC(z80 *Z80) {
	z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 7)
	z80.add16(z80.hl, z80.BC())
}
/* LD A,(BC) */
func instr__LD_A_iBC(z80 *Z80) {
	z80.a = z80.memory.readByte(z80.BC())
}
/* DEC BC */
func instr__DEC_BC(z80 *Z80) {
	z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 2)
	z80.decBC()
}
/* INC C */
func instr__INC_C(z80 *Z80) {
	z80.incC()
}
/* DEC C */
func instr__DEC_C(z80 *Z80) {
	z80.decC()
}
/* LD C,nn */
func instr__LD_C_NN(z80 *Z80) {
	z80.c = z80.memory.readByte(z80.pc)
	z80.pc++
}
/* RRCA */
func instr__RRCA(z80 *Z80) {
	z80.f = (z80.f & (FLAG_P | FLAG_Z | FLAG_S)) | (z80.a & FLAG_C)
	z80.a = (z80.a >> 1) | (z80.a << 7)
	z80.f |= (z80.a & (FLAG_3 | FLAG_5))
}
/* DJNZ offset */
func instr__DJNZ_OFFSET(z80 *Z80) {
	z80.memory.contendReadNoMreq(z80.IR(), 1)
	z80.b--
	if z80.b != 0 {
		z80.jr()
	} else {
		z80.memory.contendRead(z80.pc, 3)
	}
	z80.pc++
}
/* LD DE,nnnn */
func instr__LD_DE_NNNN(z80 *Z80) {
	b1 := z80.memory.readByte(z80.pc)
	z80.pc++
	b2 := z80.memory.readByte(z80.pc)
	z80.pc++
	z80.setDE(joinBytes(b2, b1))
}
/* LD (DE),A */
func instr__LD_iDE_A(z80 *Z80) {
	z80.memory.writeByte(z80.DE(), z80.a)
}
/* INC DE */
func instr__INC_DE(z80 *Z80) {
	z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 2)
	z80.incDE()
}
/* INC D */
func instr__INC_D(z80 *Z80) {
	z80.incD()
}
/* DEC D */
func instr__DEC_D(z80 *Z80) {
	z80.decD()
}
/* LD D,nn */
func instr__LD_D_NN(z80 *Z80) {
	z80.d = z80.memory.readByte(z80.pc)
	z80.pc++
}
/* RLA */
func instr__RLA(z80 *Z80) {
	var bytetemp byte = z80.a
	z80.a = (z80.a << 1) | (z80.f & FLAG_C)
	z80.f = (z80.f & (FLAG_P | FLAG_Z | FLAG_S)) | (z80.a & (FLAG_3 | FLAG_5)) | (bytetemp >> 7)
}
/* JR offset */
func instr__JR_OFFSET(z80 *Z80) {
	z80.jr()
	z80.pc++
}
/* ADD HL,DE */
func instr__ADD_HL_DE(z80 *Z80) {
	z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 7)
	z80.add16(z80.hl, z80.DE())
}
/* LD A,(DE) */
func instr__LD_A_iDE(z80 *Z80) {
	z80.a = z80.memory.readByte(z80.DE())
}
/* DEC DE */
func instr__DEC_DE(z80 *Z80) {
	z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 2)
	z80.decDE()
}
/* INC E */
func instr__INC_E(z80 *Z80) {
	z80.incE()
}
/* DEC E */
func instr__DEC_E(z80 *Z80) {
	z80.decE()
}
/* LD E,nn */
func instr__LD_E_NN(z80 *Z80) {
	z80.e = z80.memory.readByte(z80.pc)
	z80.pc++
}
/* RRA */
func instr__RRA(z80 *Z80) {
	var bytetemp byte = z80.a
	z80.a = (z80.a >> 1) | (z80.f << 7)
	z80.f = (z80.f & (FLAG_P | FLAG_Z | FLAG_S)) | (z80.a & (FLAG_3 | FLAG_5)) | (bytetemp & FLAG_C)
}
/* JR NZ,offset */
func instr__JR_NZ_OFFSET(z80 *Z80) {
	if (z80.f & FLAG_Z) == 0 {
		z80.jr()
	} else {
		z80.memory.contendRead(z80.pc, 3)
	}
	z80.pc++
}
/* LD HL,nnnn */
func instr__LD_HL_NNNN(z80 *Z80) {
	b1 := z80.memory.readByte(z80.pc)
	z80.pc++
	b2 := z80.memory.readByte(z80.pc)
	z80.pc++
	z80.setHL(joinBytes(b2, b1))
}
/* LD (nnnn),HL */
func instr__LD_iNNNN_HL(z80 *Z80) {
	z80.ld16nnrr(z80.l, z80.h)
	// break
}
/* INC HL */
func instr__INC_HL(z80 *Z80) {
	z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 2)
	z80.incHL()
}
/* INC H */
func instr__INC_H(z80 *Z80) {
	z80.incH()
}
/* DEC H */
func instr__DEC_H(z80 *Z80) {
	z80.decH()
}
/* LD H,nn */
func instr__LD_H_NN(z80 *Z80) {
	z80.h = z80.memory.readByte(z80.pc)
	z80.pc++
}
/* DAA */
func instr__DAA(z80 *Z80) {
	var add, carry byte = 0, (z80.f & FLAG_C)
	if ((z80.f & FLAG_H) != 0) || ((z80.a & 0x0f) > 9) {
		add = 6
	}
	if (carry != 0) || (z80.a > 0x99) {
		add |= 0x60
	}
	if z80.a > 0x99 {
		carry = FLAG_C
	}
	if (z80.f & FLAG_N) != 0 {
		z80.sub(add)
	} else {
		z80.add(add)
	}
	var temp byte = byte(int(z80.f) & ^(FLAG_C|FLAG_P)) | carry | parityTable[z80.a]
	z80.f = temp
}
/* JR Z,offset */
func instr__JR_Z_OFFSET(z80 *Z80) {
	if (z80.f & FLAG_Z) != 0 {
		z80.jr()
	} else {
		z80.memory.contendRead(z80.pc, 3)
	}
	z80.pc++
}
/* ADD HL,HL */
func instr__ADD_HL_HL(z80 *Z80) {
	z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 7)
	z80.add16(z80.hl, z80.HL())
}
/* LD HL,(nnnn) */
func instr__LD_HL_iNNNN(z80 *Z80) {
	z80.ld16rrnn(&z80.l, &z80.h)
	// break
}
/* DEC HL */
func instr__DEC_HL(z80 *Z80) {
	z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 2)
	z80.decHL()
}
/* INC L */
func instr__INC_L(z80 *Z80) {
	z80.incL()
}
/* DEC L */
func instr__DEC_L(z80 *Z80) {
	z80.decL()
}
/* LD L,nn */
func instr__LD_L_NN(z80 *Z80) {
	z80.l = z80.memory.readByte(z80.pc)
	z80.pc++
}
/* CPL */
func instr__CPL(z80 *Z80) {
	z80.a ^= 0xff
	z80.f = (z80.f & (FLAG_C | FLAG_P | FLAG_Z | FLAG_S)) |
		(z80.a & (FLAG_3 | FLAG_5)) |
		(FLAG_N | FLAG_H)
}
/* JR NC,offset */
func instr__JR_NC_OFFSET(z80 *Z80) {
	if (z80.f & FLAG_C) == 0 {
		z80.jr()
	} else {
		z80.memory.contendRead(z80.pc, 3)
	}
	z80.pc++
}
/* LD SP,nnnn */
func instr__LD_SP_NNNN(z80 *Z80) {
	b1 := z80.memory.readByte(z80.pc)
	z80.pc++
	b2 := z80.memory.readByte(z80.pc)
	z80.pc++
	z80.setSP(joinBytes(b2, b1))
}
/* LD (nnnn),A */
func instr__LD_iNNNN_A(z80 *Z80) {
	var wordtemp uint16 = uint16(z80.memory.readByte(z80.pc))
	z80.pc++
	wordtemp |= uint16(z80.memory.readByte(z80.pc)) << 8
	z80.pc++
	z80.memory.writeByte(wordtemp, z80.a)
}
/* INC SP */
func instr__INC_SP(z80 *Z80) {
	z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 2)
	z80.incSP()
}
/* INC (HL) */
func instr__INC_iHL(z80 *Z80) {
	{
		var bytetemp byte = z80.memory.readByte(z80.HL())
		z80.memory.contendReadNoMreq(z80.HL(), 1)
		z80.inc(&bytetemp)
		z80.memory.writeByte(z80.HL(), bytetemp)
	}
}
/* DEC (HL) */
func instr__DEC_iHL(z80 *Z80) {
	{
		var bytetemp byte = z80.memory.readByte(z80.HL())
		z80.memory.contendReadNoMreq(z80.HL(), 1)
		z80.dec(&bytetemp)
		z80.memory.writeByte(z80.HL(), bytetemp)
	}
}
/* LD (HL),nn */
func instr__LD_iHL_NN(z80 *Z80) {
	z80.memory.writeByte(z80.HL(), z80.memory.readByte(z80.pc))
	z80.pc++
}
/* SCF */
func instr__SCF(z80 *Z80) {
	z80.f = (z80.f & (FLAG_P | FLAG_Z | FLAG_S)) |
		(z80.a & (FLAG_3 | FLAG_5)) |
		FLAG_C
}
/* JR C,offset */
func instr__JR_C_OFFSET(z80 *Z80) {
	if (z80.f & FLAG_C) != 0 {
		z80.jr()
	} else {
		z80.memory.contendRead(z80.pc, 3)
	}
	z80.pc++
}
/* ADD HL,SP */
func instr__ADD_HL_SP(z80 *Z80) {
	z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 7)
	z80.add16(z80.hl, z80.SP())
}
/* LD A,(nnnn) */
func instr__LD_A_iNNNN(z80 *Z80) {
	var wordtemp uint16 = uint16(z80.memory.readByte(z80.pc))
	z80.pc++
	wordtemp |= uint16(z80.memory.readByte(z80.pc)) << 8
	z80.pc++
	z80.a = z80.memory.readByte(wordtemp)
}
/* DEC SP */
func instr__DEC_SP(z80 *Z80) {
	z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 2)
	z80.decSP()
}
/* INC A */
func instr__INC_A(z80 *Z80) {
	z80.incA()
}
/* DEC A */
func instr__DEC_A(z80 *Z80) {
	z80.decA()
}
/* LD A,nn */
func instr__LD_A_NN(z80 *Z80) {
	z80.a = z80.memory.readByte(z80.pc)
	z80.pc++
}
/* CCF */
func instr__CCF(z80 *Z80) {
	z80.f = (z80.f & (FLAG_P | FLAG_Z | FLAG_S)) |
		ternOpB((z80.f&FLAG_C) != 0, FLAG_H, FLAG_C) |
		(z80.a & (FLAG_3 | FLAG_5))
}
/* LD B,B */
func instr__LD_B_B(z80 *Z80) {
}
/* LD B,C */
func instr__LD_B_C(z80 *Z80) {
	z80.b = z80.c
}
/* LD B,D */
func instr__LD_B_D(z80 *Z80) {
	z80.b = z80.d
}
/* LD B,E */
func instr__LD_B_E(z80 *Z80) {
	z80.b = z80.e
}
/* LD B,H */
func instr__LD_B_H(z80 *Z80) {
	z80.b = z80.h
}
/* LD B,L */
func instr__LD_B_L(z80 *Z80) {
	z80.b = z80.l
}
/* LD B,(HL) */
func instr__LD_B_iHL(z80 *Z80) {
	z80.b = z80.memory.readByte(z80.HL())
}
/* LD B,A */
func instr__LD_B_A(z80 *Z80) {
	z80.b = z80.a
}
/* LD C,B */
func instr__LD_C_B(z80 *Z80) {
	z80.c = z80.b
}
/* LD C,C */
func instr__LD_C_C(z80 *Z80) {
}
/* LD C,D */
func instr__LD_C_D(z80 *Z80) {
	z80.c = z80.d
}
/* LD C,E */
func instr__LD_C_E(z80 *Z80) {
	z80.c = z80.e
}
/* LD C,H */
func instr__LD_C_H(z80 *Z80) {
	z80.c = z80.h
}
/* LD C,L */
func instr__LD_C_L(z80 *Z80) {
	z80.c = z80.l
}
/* LD C,(HL) */
func instr__LD_C_iHL(z80 *Z80) {
	z80.c = z80.memory.readByte(z80.HL())
}
/* LD C,A */
func instr__LD_C_A(z80 *Z80) {
	z80.c = z80.a
}
/* LD D,B */
func instr__LD_D_B(z80 *Z80) {
	z80.d = z80.b
}
/* LD D,C */
func instr__LD_D_C(z80 *Z80) {
	z80.d = z80.c
}
/* LD D,D */
func instr__LD_D_D(z80 *Z80) {
}
/* LD D,E */
func instr__LD_D_E(z80 *Z80) {
	z80.d = z80.e
}
/* LD D,H */
func instr__LD_D_H(z80 *Z80) {
	z80.d = z80.h
}
/* LD D,L */
func instr__LD_D_L(z80 *Z80) {
	z80.d = z80.l
}
/* LD D,(HL) */
func instr__LD_D_iHL(z80 *Z80) {
	z80.d = z80.memory.readByte(z80.HL())
}
/* LD D,A */
func instr__LD_D_A(z80 *Z80) {
	z80.d = z80.a
}
/* LD E,B */
func instr__LD_E_B(z80 *Z80) {
	z80.e = z80.b
}
/* LD E,C */
func instr__LD_E_C(z80 *Z80) {
	z80.e = z80.c
}
/* LD E,D */
func instr__LD_E_D(z80 *Z80) {
	z80.e = z80.d
}
/* LD E,E */
func instr__LD_E_E(z80 *Z80) {
}
/* LD E,H */
func instr__LD_E_H(z80 *Z80) {
	z80.e = z80.h
}
/* LD E,L */
func instr__LD_E_L(z80 *Z80) {
	z80.e = z80.l
}
/* LD E,(HL) */
func instr__LD_E_iHL(z80 *Z80) {
	z80.e = z80.memory.readByte(z80.HL())
}
/* LD E,A */
func instr__LD_E_A(z80 *Z80) {
	z80.e = z80.a
}
/* LD H,B */
func instr__LD_H_B(z80 *Z80) {
	z80.h = z80.b
}
/* LD H,C */
func instr__LD_H_C(z80 *Z80) {
	z80.h = z80.c
}
/* LD H,D */
func instr__LD_H_D(z80 *Z80) {
	z80.h = z80.d
}
/* LD H,E */
func instr__LD_H_E(z80 *Z80) {
	z80.h = z80.e
}
/* LD H,H */
func instr__LD_H_H(z80 *Z80) {
}
/* LD H,L */
func instr__LD_H_L(z80 *Z80) {
	z80.h = z80.l
}
/* LD H,(HL) */
func instr__LD_H_iHL(z80 *Z80) {
	z80.h = z80.memory.readByte(z80.HL())
}
/* LD H,A */
func instr__LD_H_A(z80 *Z80) {
	z80.h = z80.a
}
/* LD L,B */
func instr__LD_L_B(z80 *Z80) {
	z80.l = z80.b
}
/* LD L,C */
func instr__LD_L_C(z80 *Z80) {
	z80.l = z80.c
}
/* LD L,D */
func instr__LD_L_D(z80 *Z80) {
	z80.l = z80.d
}
/* LD L,E */
func instr__LD_L_E(z80 *Z80) {
	z80.l = z80.e
}
/* LD L,H */
func instr__LD_L_H(z80 *Z80) {
	z80.l = z80.h
}
/* LD L,L */
func instr__LD_L_L(z80 *Z80) {
}
/* LD L,(HL) */
func instr__LD_L_iHL(z80 *Z80) {
	z80.l = z80.memory.readByte(z80.HL())
}
/* LD L,A */
func instr__LD_L_A(z80 *Z80) {
	z80.l = z80.a
}
/* LD (HL),B */
func instr__LD_iHL_B(z80 *Z80) {
	z80.memory.writeByte(z80.HL(), z80.b)
}
/* LD (HL),C */
func instr__LD_iHL_C(z80 *Z80) {
	z80.memory.writeByte(z80.HL(), z80.c)
}
/* LD (HL),D */
func instr__LD_iHL_D(z80 *Z80) {
	z80.memory.writeByte(z80.HL(), z80.d)
}
/* LD (HL),E */
func instr__LD_iHL_E(z80 *Z80) {
	z80.memory.writeByte(z80.HL(), z80.e)
}
/* LD (HL),H */
func instr__LD_iHL_H(z80 *Z80) {
	z80.memory.writeByte(z80.HL(), z80.h)
}
/* LD (HL),L */
func instr__LD_iHL_L(z80 *Z80) {
	z80.memory.writeByte(z80.HL(), z80.l)
}
/* HALT */
func instr__HALT(z80 *Z80) {
	z80.halted = true
	z80.pc--
	return
}
/* LD (HL),A */
func instr__LD_iHL_A(z80 *Z80) {
	z80.memory.writeByte(z80.HL(), z80.a)
}
/* LD A,B */
func instr__LD_A_B(z80 *Z80) {
	z80.a = z80.b
}
/* LD A,C */
func instr__LD_A_C(z80 *Z80) {
	z80.a = z80.c
}
/* LD A,D */
func instr__LD_A_D(z80 *Z80) {
	z80.a = z80.d
}
/* LD A,E */
func instr__LD_A_E(z80 *Z80) {
	z80.a = z80.e
}
/* LD A,H */
func instr__LD_A_H(z80 *Z80) {
	z80.a = z80.h
}
/* LD A,L */
func instr__LD_A_L(z80 *Z80) {
	z80.a = z80.l
}
/* LD A,(HL) */
func instr__LD_A_iHL(z80 *Z80) {
	z80.a = z80.memory.readByte(z80.HL())
}
/* LD A,A */
func instr__LD_A_A(z80 *Z80) {
}
/* ADD A,B */
func instr__ADD_A_B(z80 *Z80) {
	z80.add(z80.b)
}
/* ADD A,C */
func instr__ADD_A_C(z80 *Z80) {
	z80.add(z80.c)
}
/* ADD A,D */
func instr__ADD_A_D(z80 *Z80) {
	z80.add(z80.d)
}
/* ADD A,E */
func instr__ADD_A_E(z80 *Z80) {
	z80.add(z80.e)
}
/* ADD A,H */
func instr__ADD_A_H(z80 *Z80) {
	z80.add(z80.h)
}
/* ADD A,L */
func instr__ADD_A_L(z80 *Z80) {
	z80.add(z80.l)
}
/* ADD A,(HL) */
func instr__ADD_A_iHL(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.HL())

	z80.add(bytetemp)
}
/* ADD A,A */
func instr__ADD_A_A(z80 *Z80) {
	z80.add(z80.a)
}
/* ADC A,B */
func instr__ADC_A_B(z80 *Z80) {
	z80.adc(z80.b)
}
/* ADC A,C */
func instr__ADC_A_C(z80 *Z80) {
	z80.adc(z80.c)
}
/* ADC A,D */
func instr__ADC_A_D(z80 *Z80) {
	z80.adc(z80.d)
}
/* ADC A,E */
func instr__ADC_A_E(z80 *Z80) {
	z80.adc(z80.e)
}
/* ADC A,H */
func instr__ADC_A_H(z80 *Z80) {
	z80.adc(z80.h)
}
/* ADC A,L */
func instr__ADC_A_L(z80 *Z80) {
	z80.adc(z80.l)
}
/* ADC A,(HL) */
func instr__ADC_A_iHL(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.HL())

	z80.adc(bytetemp)
}
/* ADC A,A */
func instr__ADC_A_A(z80 *Z80) {
	z80.adc(z80.a)
}
/* SUB A,B */
func instr__SUB_A_B(z80 *Z80) {
	z80.sub(z80.b)
}
/* SUB A,C */
func instr__SUB_A_C(z80 *Z80) {
	z80.sub(z80.c)
}
/* SUB A,D */
func instr__SUB_A_D(z80 *Z80) {
	z80.sub(z80.d)
}
/* SUB A,E */
func instr__SUB_A_E(z80 *Z80) {
	z80.sub(z80.e)
}
/* SUB A,H */
func instr__SUB_A_H(z80 *Z80) {
	z80.sub(z80.h)
}
/* SUB A,L */
func instr__SUB_A_L(z80 *Z80) {
	z80.sub(z80.l)
}
/* SUB A,(HL) */
func instr__SUB_A_iHL(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.HL())

	z80.sub(bytetemp)
}
/* SUB A,A */
func instr__SUB_A_A(z80 *Z80) {
	z80.sub(z80.a)
}
/* SBC A,B */
func instr__SBC_A_B(z80 *Z80) {
	z80.sbc(z80.b)
}
/* SBC A,C */
func instr__SBC_A_C(z80 *Z80) {
	z80.sbc(z80.c)
}
/* SBC A,D */
func instr__SBC_A_D(z80 *Z80) {
	z80.sbc(z80.d)
}
/* SBC A,E */
func instr__SBC_A_E(z80 *Z80) {
	z80.sbc(z80.e)
}
/* SBC A,H */
func instr__SBC_A_H(z80 *Z80) {
	z80.sbc(z80.h)
}
/* SBC A,L */
func instr__SBC_A_L(z80 *Z80) {
	z80.sbc(z80.l)
}
/* SBC A,(HL) */
func instr__SBC_A_iHL(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.HL())

	z80.sbc(bytetemp)
}
/* SBC A,A */
func instr__SBC_A_A(z80 *Z80) {
	z80.sbc(z80.a)
}
/* AND A,B */
func instr__AND_A_B(z80 *Z80) {
	z80.and(z80.b)
}
/* AND A,C */
func instr__AND_A_C(z80 *Z80) {
	z80.and(z80.c)
}
/* AND A,D */
func instr__AND_A_D(z80 *Z80) {
	z80.and(z80.d)
}
/* AND A,E */
func instr__AND_A_E(z80 *Z80) {
	z80.and(z80.e)
}
/* AND A,H */
func instr__AND_A_H(z80 *Z80) {
	z80.and(z80.h)
}
/* AND A,L */
func instr__AND_A_L(z80 *Z80) {
	z80.and(z80.l)
}
/* AND A,(HL) */
func instr__AND_A_iHL(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.HL())

	z80.and(bytetemp)
}
/* AND A,A */
func instr__AND_A_A(z80 *Z80) {
	z80.and(z80.a)
}
/* XOR A,B */
func instr__XOR_A_B(z80 *Z80) {
	z80.xor(z80.b)
}
/* XOR A,C */
func instr__XOR_A_C(z80 *Z80) {
	z80.xor(z80.c)
}
/* XOR A,D */
func instr__XOR_A_D(z80 *Z80) {
	z80.xor(z80.d)
}
/* XOR A,E */
func instr__XOR_A_E(z80 *Z80) {
	z80.xor(z80.e)
}
/* XOR A,H */
func instr__XOR_A_H(z80 *Z80) {
	z80.xor(z80.h)
}
/* XOR A,L */
func instr__XOR_A_L(z80 *Z80) {
	z80.xor(z80.l)
}
/* XOR A,(HL) */
func instr__XOR_A_iHL(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.HL())

	z80.xor(bytetemp)
}
/* XOR A,A */
func instr__XOR_A_A(z80 *Z80) {
	z80.xor(z80.a)
}
/* OR A,B */
func instr__OR_A_B(z80 *Z80) {
	z80.or(z80.b)
}
/* OR A,C */
func instr__OR_A_C(z80 *Z80) {
	z80.or(z80.c)
}
/* OR A,D */
func instr__OR_A_D(z80 *Z80) {
	z80.or(z80.d)
}
/* OR A,E */
func instr__OR_A_E(z80 *Z80) {
	z80.or(z80.e)
}
/* OR A,H */
func instr__OR_A_H(z80 *Z80) {
	z80.or(z80.h)
}
/* OR A,L */
func instr__OR_A_L(z80 *Z80) {
	z80.or(z80.l)
}
/* OR A,(HL) */
func instr__OR_A_iHL(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.HL())

	z80.or(bytetemp)
}
/* OR A,A */
func instr__OR_A_A(z80 *Z80) {
	z80.or(z80.a)
}
/* CP B */
func instr__CP_B(z80 *Z80) {
	z80.cp(z80.b)
}
/* CP C */
func instr__CP_C(z80 *Z80) {
	z80.cp(z80.c)
}
/* CP D */
func instr__CP_D(z80 *Z80) {
	z80.cp(z80.d)
}
/* CP E */
func instr__CP_E(z80 *Z80) {
	z80.cp(z80.e)
}
/* CP H */
func instr__CP_H(z80 *Z80) {
	z80.cp(z80.h)
}
/* CP L */
func instr__CP_L(z80 *Z80) {
	z80.cp(z80.l)
}
/* CP (HL) */
func instr__CP_iHL(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.HL())

	z80.cp(bytetemp)
}
/* CP A */
func instr__CP_A(z80 *Z80) {
	z80.cp(z80.a)
}
/* RET NZ */
func instr__RET_NZ(z80 *Z80) {
	z80.memory.contendReadNoMreq(z80.IR(), 1)
	if z80.pc == 0x056c || z80.pc == 0x0112 {
		if z80.tapeLoadTrap() == 0 { /*break*/
		}
	}
	if !((z80.f & FLAG_Z) != 0) {
		z80.ret()
	}
}
/* POP BC */
func instr__POP_BC(z80 *Z80) {
	z80.c, z80.b = z80.pop16()
}
/* JP NZ,nnnn */
func instr__JP_NZ_NNNN(z80 *Z80) {
	if (z80.f & FLAG_Z) == 0 {
		z80.jp()
	} else {
		z80.memory.contendRead(z80.pc, 3)
		z80.memory.contendRead(z80.pc+1, 3)
		z80.pc += 2
	}
}
/* JP nnnn */
func instr__JP_NNNN(z80 *Z80) {
	z80.jp()
}
/* CALL NZ,nnnn */
func instr__CALL_NZ_NNNN(z80 *Z80) {
	if (z80.f & FLAG_Z) == 0 {
		z80.call()
	} else {
		z80.memory.contendRead(z80.pc, 3)
		z80.memory.contendRead(z80.pc+1, 3)
		z80.pc += 2
	}
}
/* PUSH BC */
func instr__PUSH_BC(z80 *Z80) {
	z80.memory.contendReadNoMreq(z80.IR(), 1)
	z80.push16(z80.c, z80.b)
}
/* ADD A,nn */
func instr__ADD_A_NN(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.PC())
	z80.pc++
	z80.add(bytetemp)
}
/* RST 00 */
func instr__RST_00(z80 *Z80) {
	z80.memory.contendReadNoMreq(z80.IR(), 1)
	z80.rst(0x00)
}
/* RET Z */
func instr__RET_Z(z80 *Z80) {
	z80.memory.contendReadNoMreq(z80.IR(), 1)
	if (z80.f & FLAG_Z) != 0 {
		z80.ret()
	}
}
/* RET */
func instr__RET(z80 *Z80) {
	z80.ret()
}
/* JP Z,nnnn */
func instr__JP_Z_NNNN(z80 *Z80) {
	if (z80.f & FLAG_Z) != 0 {
		z80.jp()
	} else {
		z80.memory.contendRead(z80.pc, 3)
		z80.memory.contendRead(z80.pc+1, 3)
		z80.pc += 2
	}
}
/* shift CB */
func instr__SHIFT_CB(z80 *Z80) {
}
/* CALL Z,nnnn */
func instr__CALL_Z_NNNN(z80 *Z80) {
	if (z80.f & FLAG_Z) != 0 {
		z80.call()
	} else {
		z80.memory.contendRead(z80.pc, 3)
		z80.memory.contendRead(z80.pc+1, 3)
		z80.pc += 2
	}
}
/* CALL nnnn */
func instr__CALL_NNNN(z80 *Z80) {
	z80.call()
}
/* ADC A,nn */
func instr__ADC_A_NN(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.PC())
	z80.pc++
	z80.adc(bytetemp)
}
/* RST 8 */
func instr__RST_8(z80 *Z80) {
	z80.memory.contendReadNoMreq(z80.IR(), 1)
	z80.rst(0x8)
}
/* RET NC */
func instr__RET_NC(z80 *Z80) {
	z80.memory.contendReadNoMreq(z80.IR(), 1)
	if !((z80.f & FLAG_C) != 0) {
		z80.ret()
	}
}
/* POP DE */
func instr__POP_DE(z80 *Z80) {
	z80.e, z80.d = z80.pop16()
}
/* JP NC,nnnn */
func instr__JP_NC_NNNN(z80 *Z80) {
	if (z80.f & FLAG_C) == 0 {
		z80.jp()
	} else {
		z80.memory.contendRead(z80.pc, 3)
		z80.memory.contendRead(z80.pc+1, 3)
		z80.pc += 2
	}
}
/* OUT (nn),A */
func instr__OUT_iNN_A(z80 *Z80) {
	var outtemp uint16 = uint16(z80.memory.readByte(z80.pc)) + (uint16(z80.a) << 8)
	z80.pc++
	z80.writePort(outtemp, z80.a)
}
/* CALL NC,nnnn */
func instr__CALL_NC_NNNN(z80 *Z80) {
	if (z80.f & FLAG_C) == 0 {
		z80.call()
	} else {
		z80.memory.contendRead(z80.pc, 3)
		z80.memory.contendRead(z80.pc+1, 3)
		z80.pc += 2
	}
}
/* PUSH DE */
func instr__PUSH_DE(z80 *Z80) {
	z80.memory.contendReadNoMreq(z80.IR(), 1)
	z80.push16(z80.e, z80.d)
}
/* SUB nn */
func instr__SUB_NN(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.PC())
	z80.pc++
	z80.sub(bytetemp)
}
/* RST 10 */
func instr__RST_10(z80 *Z80) {
	z80.memory.contendReadNoMreq(z80.IR(), 1)
	z80.rst(0x10)
}
/* RET C */
func instr__RET_C(z80 *Z80) {
	z80.memory.contendReadNoMreq(z80.IR(), 1)
	if (z80.f & FLAG_C) != 0 {
		z80.ret()
	}
}
/* EXX */
func instr__EXX(z80 *Z80) {
	var wordtemp uint16 = z80.BC()
	z80.setBC(z80.BC_())
	z80.setBC_(wordtemp)

	wordtemp = z80.DE()
	z80.setDE(z80.DE_())
	z80.setDE_(wordtemp)

	wordtemp = z80.HL()
	z80.setHL(z80.HL_())
	z80.setHL_(wordtemp)
}
/* JP C,nnnn */
func instr__JP_C_NNNN(z80 *Z80) {
	if (z80.f & FLAG_C) != 0 {
		z80.jp()
	} else {
		z80.memory.contendRead(z80.pc, 3)
		z80.memory.contendRead(z80.pc+1, 3)
		z80.pc += 2
	}
}
/* IN A,(nn) */
func instr__IN_A_iNN(z80 *Z80) {
	var intemp uint16 = uint16(z80.memory.readByte(z80.pc)) + (uint16(z80.a) << 8)
	z80.pc++
	z80.a = z80.readPort(intemp)
}
/* CALL C,nnnn */
func instr__CALL_C_NNNN(z80 *Z80) {
	if (z80.f & FLAG_C) != 0 {
		z80.call()
	} else {
		z80.memory.contendRead(z80.pc, 3)
		z80.memory.contendRead(z80.pc+1, 3)
		z80.pc += 2
	}
}
/* shift DD */
func instr__SHIFT_DD(z80 *Z80) {
}
/* SBC A,nn */
func instr__SBC_A_NN(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.PC())
	z80.pc++
	z80.sbc(bytetemp)
}
/* RST 18 */
func instr__RST_18(z80 *Z80) {
	z80.memory.contendReadNoMreq(z80.IR(), 1)
	z80.rst(0x18)
}
/* RET PO */
func instr__RET_PO(z80 *Z80) {
	z80.memory.contendReadNoMreq(z80.IR(), 1)
	if !((z80.f & FLAG_P) != 0) {
		z80.ret()
	}
}
/* POP HL */
func instr__POP_HL(z80 *Z80) {
	z80.l, z80.h = z80.pop16()
}
/* JP PO,nnnn */
func instr__JP_PO_NNNN(z80 *Z80) {
	if (z80.f & FLAG_P) == 0 {
		z80.jp()
	} else {
		z80.memory.contendRead(z80.pc, 3)
		z80.memory.contendRead(z80.pc+1, 3)
		z80.pc += 2
	}
}
/* EX (SP),HL */
func instr__EX_iSP_HL(z80 *Z80) {
	var bytetempl = z80.memory.readByte(z80.SP())
	var bytetemph = z80.memory.readByte(z80.SP() + 1)
	z80.memory.contendReadNoMreq(z80.SP()+1, 1)
	z80.memory.writeByte(z80.SP()+1, z80.h)
	z80.memory.writeByte(z80.SP(), z80.l)
	z80.memory.contendWriteNoMreq_loop(z80.SP(), 1, 2)
	z80.l = bytetempl
	z80.h = bytetemph
}
/* CALL PO,nnnn */
func instr__CALL_PO_NNNN(z80 *Z80) {
	if (z80.f & FLAG_P) == 0 {
		z80.call()
	} else {
		z80.memory.contendRead(z80.pc, 3)
		z80.memory.contendRead(z80.pc+1, 3)
		z80.pc += 2
	}
}
/* PUSH HL */
func instr__PUSH_HL(z80 *Z80) {
	z80.memory.contendReadNoMreq(z80.IR(), 1)
	z80.push16(z80.l, z80.h)
}
/* AND nn */
func instr__AND_NN(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.PC())
	z80.pc++
	z80.and(bytetemp)
}
/* RST 20 */
func instr__RST_20(z80 *Z80) {
	z80.memory.contendReadNoMreq(z80.IR(), 1)
	z80.rst(0x20)
}
/* RET PE */
func instr__RET_PE(z80 *Z80) {
	z80.memory.contendReadNoMreq(z80.IR(), 1)
	if (z80.f & FLAG_P) != 0 {
		z80.ret()
	}
}
/* JP HL */
func instr__JP_HL(z80 *Z80) {
	z80.pc = z80.HL() /* NB: NOT INDIRECT! */
}
/* JP PE,nnnn */
func instr__JP_PE_NNNN(z80 *Z80) {
	if (z80.f & FLAG_P) != 0 {
		z80.jp()
	} else {
		z80.memory.contendRead(z80.pc, 3)
		z80.memory.contendRead(z80.pc+1, 3)
		z80.pc += 2
	}
}
/* EX DE,HL */
func instr__EX_DE_HL(z80 *Z80) {
	var wordtemp uint16 = z80.DE()
	z80.setDE(z80.HL())
	z80.setHL(wordtemp)
}
/* CALL PE,nnnn */
func instr__CALL_PE_NNNN(z80 *Z80) {
	if (z80.f & FLAG_P) != 0 {
		z80.call()
	} else {
		z80.memory.contendRead(z80.pc, 3)
		z80.memory.contendRead(z80.pc+1, 3)
		z80.pc += 2
	}
}
/* shift ED */
func instr__SHIFT_ED(z80 *Z80) {
}
/* XOR A,nn */
func instr__XOR_A_NN(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.PC())
	z80.pc++
	z80.xor(bytetemp)
}
/* RST 28 */
func instr__RST_28(z80 *Z80) {
	z80.memory.contendReadNoMreq(z80.IR(), 1)
	z80.rst(0x28)
}
/* RET P */
func instr__RET_P(z80 *Z80) {
	z80.memory.contendReadNoMreq(z80.IR(), 1)
	if !((z80.f & FLAG_S) != 0) {
		z80.ret()
	}
}
/* POP AF */
func instr__POP_AF(z80 *Z80) {
	z80.f, z80.a = z80.pop16()
}
/* JP P,nnnn */
func instr__JP_P_NNNN(z80 *Z80) {
	if (z80.f & FLAG_S) == 0 {
		z80.jp()
	} else {
		z80.memory.contendRead(z80.pc, 3)
		z80.memory.contendRead(z80.pc+1, 3)
		z80.pc += 2
	}
}
/* DI */
func instr__DI(z80 *Z80) {
	z80.iff1, z80.iff2 = 0, 0
}
/* CALL P,nnnn */
func instr__CALL_P_NNNN(z80 *Z80) {
	if (z80.f & FLAG_S) == 0 {
		z80.call()
	} else {
		z80.memory.contendRead(z80.pc, 3)
		z80.memory.contendRead(z80.pc+1, 3)
		z80.pc += 2
	}
}
/* PUSH AF */
func instr__PUSH_AF(z80 *Z80) {
	z80.memory.contendReadNoMreq(z80.IR(), 1)
	z80.push16(z80.f, z80.a)
}
/* OR nn */
func instr__OR_NN(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.PC())
	z80.pc++
	z80.or(bytetemp)
}
/* RST 30 */
func instr__RST_30(z80 *Z80) {
	z80.memory.contendReadNoMreq(z80.IR(), 1)
	z80.rst(0x30)
}
/* RET M */
func instr__RET_M(z80 *Z80) {
	z80.memory.contendReadNoMreq(z80.IR(), 1)
	if (z80.f & FLAG_S) != 0 {
		z80.ret()
	}
}
/* LD SP,HL */
func instr__LD_SP_HL(z80 *Z80) {
	z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 2)
	z80.sp = z80.HL()
}
/* JP M,nnnn */
func instr__JP_M_NNNN(z80 *Z80) {
	if (z80.f & FLAG_S) != 0 {
		z80.jp()
	} else {
		z80.memory.contendRead(z80.pc, 3)
		z80.memory.contendRead(z80.pc+1, 3)
		z80.pc += 2
	}
}
/* EI */
func instr__EI(z80 *Z80) {
	/* Interrupts are not accepted immediately after an EI, but are
	   accepted after the next instruction */
	z80.iff1, z80.iff2 = 1, 1
	z80.interruptsEnabledAt = int(z80.tstates)
	// eventAdd(z80.tstates + 1, z80InterruptEvent)
}
/* CALL M,nnnn */
func instr__CALL_M_NNNN(z80 *Z80) {
	if (z80.f & FLAG_S) != 0 {
		z80.call()
	} else {
		z80.memory.contendRead(z80.pc, 3)
		z80.memory.contendRead(z80.pc+1, 3)
		z80.pc += 2
	}
}
/* shift FD */
func instr__SHIFT_FD(z80 *Z80) {
}
/* CP nn */
func instr__CP_NN(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.PC())
	z80.pc++
	z80.cp(bytetemp)
}
/* RST 38 */
func instr__RST_38(z80 *Z80) {
	z80.memory.contendReadNoMreq(z80.IR(), 1)
	z80.rst(0x38)
}
/* RLC B */
func instrCB__RLC_B(z80 *Z80) {
	z80.b = z80.rlc(z80.b)
}
/* RLC C */
func instrCB__RLC_C(z80 *Z80) {
	z80.c = z80.rlc(z80.c)
}
/* RLC D */
func instrCB__RLC_D(z80 *Z80) {
	z80.d = z80.rlc(z80.d)
}
/* RLC E */
func instrCB__RLC_E(z80 *Z80) {
	z80.e = z80.rlc(z80.e)
}
/* RLC H */
func instrCB__RLC_H(z80 *Z80) {
	z80.h = z80.rlc(z80.h)
}
/* RLC L */
func instrCB__RLC_L(z80 *Z80) {
	z80.l = z80.rlc(z80.l)
}
/* RLC (HL) */
func instrCB__RLC_iHL(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.HL())
	z80.memory.contendReadNoMreq(z80.HL(), 1)
	bytetemp = z80.rlc(bytetemp)
	z80.memory.writeByte(z80.HL(), bytetemp)
}
/* RLC A */
func instrCB__RLC_A(z80 *Z80) {
	z80.a = z80.rlc(z80.a)
}
/* RRC B */
func instrCB__RRC_B(z80 *Z80) {
	z80.b = z80.rrc(z80.b)
}
/* RRC C */
func instrCB__RRC_C(z80 *Z80) {
	z80.c = z80.rrc(z80.c)
}
/* RRC D */
func instrCB__RRC_D(z80 *Z80) {
	z80.d = z80.rrc(z80.d)
}
/* RRC E */
func instrCB__RRC_E(z80 *Z80) {
	z80.e = z80.rrc(z80.e)
}
/* RRC H */
func instrCB__RRC_H(z80 *Z80) {
	z80.h = z80.rrc(z80.h)
}
/* RRC L */
func instrCB__RRC_L(z80 *Z80) {
	z80.l = z80.rrc(z80.l)
}
/* RRC (HL) */
func instrCB__RRC_iHL(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.HL())
	z80.memory.contendReadNoMreq(z80.HL(), 1)
	bytetemp = z80.rrc(bytetemp)
	z80.memory.writeByte(z80.HL(), bytetemp)
}
/* RRC A */
func instrCB__RRC_A(z80 *Z80) {
	z80.a = z80.rrc(z80.a)
}
/* RL B */
func instrCB__RL_B(z80 *Z80) {
	z80.b = z80.rl(z80.b)
}
/* RL C */
func instrCB__RL_C(z80 *Z80) {
	z80.c = z80.rl(z80.c)
}
/* RL D */
func instrCB__RL_D(z80 *Z80) {
	z80.d = z80.rl(z80.d)
}
/* RL E */
func instrCB__RL_E(z80 *Z80) {
	z80.e = z80.rl(z80.e)
}
/* RL H */
func instrCB__RL_H(z80 *Z80) {
	z80.h = z80.rl(z80.h)
}
/* RL L */
func instrCB__RL_L(z80 *Z80) {
	z80.l = z80.rl(z80.l)
}
/* RL (HL) */
func instrCB__RL_iHL(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.HL())
	z80.memory.contendReadNoMreq(z80.HL(), 1)
	bytetemp = z80.rl(bytetemp)
	z80.memory.writeByte(z80.HL(), bytetemp)
}
/* RL A */
func instrCB__RL_A(z80 *Z80) {
	z80.a = z80.rl(z80.a)
}
/* RR B */
func instrCB__RR_B(z80 *Z80) {
	z80.b = z80.rr(z80.b)
}
/* RR C */
func instrCB__RR_C(z80 *Z80) {
	z80.c = z80.rr(z80.c)
}
/* RR D */
func instrCB__RR_D(z80 *Z80) {
	z80.d = z80.rr(z80.d)
}
/* RR E */
func instrCB__RR_E(z80 *Z80) {
	z80.e = z80.rr(z80.e)
}
/* RR H */
func instrCB__RR_H(z80 *Z80) {
	z80.h = z80.rr(z80.h)
}
/* RR L */
func instrCB__RR_L(z80 *Z80) {
	z80.l = z80.rr(z80.l)
}
/* RR (HL) */
func instrCB__RR_iHL(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.HL())
	z80.memory.contendReadNoMreq(z80.HL(), 1)
	bytetemp = z80.rr(bytetemp)
	z80.memory.writeByte(z80.HL(), bytetemp)
}
/* RR A */
func instrCB__RR_A(z80 *Z80) {
	z80.a = z80.rr(z80.a)
}
/* SLA B */
func instrCB__SLA_B(z80 *Z80) {
	z80.b = z80.sla(z80.b)
}
/* SLA C */
func instrCB__SLA_C(z80 *Z80) {
	z80.c = z80.sla(z80.c)
}
/* SLA D */
func instrCB__SLA_D(z80 *Z80) {
	z80.d = z80.sla(z80.d)
}
/* SLA E */
func instrCB__SLA_E(z80 *Z80) {
	z80.e = z80.sla(z80.e)
}
/* SLA H */
func instrCB__SLA_H(z80 *Z80) {
	z80.h = z80.sla(z80.h)
}
/* SLA L */
func instrCB__SLA_L(z80 *Z80) {
	z80.l = z80.sla(z80.l)
}
/* SLA (HL) */
func instrCB__SLA_iHL(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.HL())
	z80.memory.contendReadNoMreq(z80.HL(), 1)
	bytetemp = z80.sla(bytetemp)
	z80.memory.writeByte(z80.HL(), bytetemp)
}
/* SLA A */
func instrCB__SLA_A(z80 *Z80) {
	z80.a = z80.sla(z80.a)
}
/* SRA B */
func instrCB__SRA_B(z80 *Z80) {
	z80.b = z80.sra(z80.b)
}
/* SRA C */
func instrCB__SRA_C(z80 *Z80) {
	z80.c = z80.sra(z80.c)
}
/* SRA D */
func instrCB__SRA_D(z80 *Z80) {
	z80.d = z80.sra(z80.d)
}
/* SRA E */
func instrCB__SRA_E(z80 *Z80) {
	z80.e = z80.sra(z80.e)
}
/* SRA H */
func instrCB__SRA_H(z80 *Z80) {
	z80.h = z80.sra(z80.h)
}
/* SRA L */
func instrCB__SRA_L(z80 *Z80) {
	z80.l = z80.sra(z80.l)
}
/* SRA (HL) */
func instrCB__SRA_iHL(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.HL())
	z80.memory.contendReadNoMreq(z80.HL(), 1)
	bytetemp = z80.sra(bytetemp)
	z80.memory.writeByte(z80.HL(), bytetemp)
}
/* SRA A */
func instrCB__SRA_A(z80 *Z80) {
	z80.a = z80.sra(z80.a)
}
/* SLL B */
func instrCB__SLL_B(z80 *Z80) {
	z80.b = z80.sll(z80.b)
}
/* SLL C */
func instrCB__SLL_C(z80 *Z80) {
	z80.c = z80.sll(z80.c)
}
/* SLL D */
func instrCB__SLL_D(z80 *Z80) {
	z80.d = z80.sll(z80.d)
}
/* SLL E */
func instrCB__SLL_E(z80 *Z80) {
	z80.e = z80.sll(z80.e)
}
/* SLL H */
func instrCB__SLL_H(z80 *Z80) {
	z80.h = z80.sll(z80.h)
}
/* SLL L */
func instrCB__SLL_L(z80 *Z80) {
	z80.l = z80.sll(z80.l)
}
/* SLL (HL) */
func instrCB__SLL_iHL(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.HL())
	z80.memory.contendReadNoMreq(z80.HL(), 1)
	bytetemp = z80.sll(bytetemp)
	z80.memory.writeByte(z80.HL(), bytetemp)
}
/* SLL A */
func instrCB__SLL_A(z80 *Z80) {
	z80.a = z80.sll(z80.a)
}
/* SRL B */
func instrCB__SRL_B(z80 *Z80) {
	z80.b = z80.srl(z80.b)
}
/* SRL C */
func instrCB__SRL_C(z80 *Z80) {
	z80.c = z80.srl(z80.c)
}
/* SRL D */
func instrCB__SRL_D(z80 *Z80) {
	z80.d = z80.srl(z80.d)
}
/* SRL E */
func instrCB__SRL_E(z80 *Z80) {
	z80.e = z80.srl(z80.e)
}
/* SRL H */
func instrCB__SRL_H(z80 *Z80) {
	z80.h = z80.srl(z80.h)
}
/* SRL L */
func instrCB__SRL_L(z80 *Z80) {
	z80.l = z80.srl(z80.l)
}
/* SRL (HL) */
func instrCB__SRL_iHL(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.HL())
	z80.memory.contendReadNoMreq(z80.HL(), 1)
	bytetemp = z80.srl(bytetemp)
	z80.memory.writeByte(z80.HL(), bytetemp)
}
/* SRL A */
func instrCB__SRL_A(z80 *Z80) {
	z80.a = z80.srl(z80.a)
}
/* BIT 0,B */
func instrCB__BIT_0_B(z80 *Z80) {
	z80.bit(0, z80.b)
}
/* BIT 0,C */
func instrCB__BIT_0_C(z80 *Z80) {
	z80.bit(0, z80.c)
}
/* BIT 0,D */
func instrCB__BIT_0_D(z80 *Z80) {
	z80.bit(0, z80.d)
}
/* BIT 0,E */
func instrCB__BIT_0_E(z80 *Z80) {
	z80.bit(0, z80.e)
}
/* BIT 0,H */
func instrCB__BIT_0_H(z80 *Z80) {
	z80.bit(0, z80.h)
}
/* BIT 0,L */
func instrCB__BIT_0_L(z80 *Z80) {
	z80.bit(0, z80.l)
}
/* BIT 0,(HL) */
func instrCB__BIT_0_iHL(z80 *Z80) {
	bytetemp := z80.memory.readByte(z80.HL())
	z80.memory.contendReadNoMreq(z80.HL(), 1)
	z80.bit(0, bytetemp)
}
/* BIT 0,A */
func instrCB__BIT_0_A(z80 *Z80) {
	z80.bit(0, z80.a)
}
/* BIT 1,B */
func instrCB__BIT_1_B(z80 *Z80) {
	z80.bit(1, z80.b)
}
/* BIT 1,C */
func instrCB__BIT_1_C(z80 *Z80) {
	z80.bit(1, z80.c)
}
/* BIT 1,D */
func instrCB__BIT_1_D(z80 *Z80) {
	z80.bit(1, z80.d)
}
/* BIT 1,E */
func instrCB__BIT_1_E(z80 *Z80) {
	z80.bit(1, z80.e)
}
/* BIT 1,H */
func instrCB__BIT_1_H(z80 *Z80) {
	z80.bit(1, z80.h)
}
/* BIT 1,L */
func instrCB__BIT_1_L(z80 *Z80) {
	z80.bit(1, z80.l)
}
/* BIT 1,(HL) */
func instrCB__BIT_1_iHL(z80 *Z80) {
	bytetemp := z80.memory.readByte(z80.HL())
	z80.memory.contendReadNoMreq(z80.HL(), 1)
	z80.bit(1, bytetemp)
}
/* BIT 1,A */
func instrCB__BIT_1_A(z80 *Z80) {
	z80.bit(1, z80.a)
}
/* BIT 2,B */
func instrCB__BIT_2_B(z80 *Z80) {
	z80.bit(2, z80.b)
}
/* BIT 2,C */
func instrCB__BIT_2_C(z80 *Z80) {
	z80.bit(2, z80.c)
}
/* BIT 2,D */
func instrCB__BIT_2_D(z80 *Z80) {
	z80.bit(2, z80.d)
}
/* BIT 2,E */
func instrCB__BIT_2_E(z80 *Z80) {
	z80.bit(2, z80.e)
}
/* BIT 2,H */
func instrCB__BIT_2_H(z80 *Z80) {
	z80.bit(2, z80.h)
}
/* BIT 2,L */
func instrCB__BIT_2_L(z80 *Z80) {
	z80.bit(2, z80.l)
}
/* BIT 2,(HL) */
func instrCB__BIT_2_iHL(z80 *Z80) {
	bytetemp := z80.memory.readByte(z80.HL())
	z80.memory.contendReadNoMreq(z80.HL(), 1)
	z80.bit(2, bytetemp)
}
/* BIT 2,A */
func instrCB__BIT_2_A(z80 *Z80) {
	z80.bit(2, z80.a)
}
/* BIT 3,B */
func instrCB__BIT_3_B(z80 *Z80) {
	z80.bit(3, z80.b)
}
/* BIT 3,C */
func instrCB__BIT_3_C(z80 *Z80) {
	z80.bit(3, z80.c)
}
/* BIT 3,D */
func instrCB__BIT_3_D(z80 *Z80) {
	z80.bit(3, z80.d)
}
/* BIT 3,E */
func instrCB__BIT_3_E(z80 *Z80) {
	z80.bit(3, z80.e)
}
/* BIT 3,H */
func instrCB__BIT_3_H(z80 *Z80) {
	z80.bit(3, z80.h)
}
/* BIT 3,L */
func instrCB__BIT_3_L(z80 *Z80) {
	z80.bit(3, z80.l)
}
/* BIT 3,(HL) */
func instrCB__BIT_3_iHL(z80 *Z80) {
	bytetemp := z80.memory.readByte(z80.HL())
	z80.memory.contendReadNoMreq(z80.HL(), 1)
	z80.bit(3, bytetemp)
}
/* BIT 3,A */
func instrCB__BIT_3_A(z80 *Z80) {
	z80.bit(3, z80.a)
}
/* BIT 4,B */
func instrCB__BIT_4_B(z80 *Z80) {
	z80.bit(4, z80.b)
}
/* BIT 4,C */
func instrCB__BIT_4_C(z80 *Z80) {
	z80.bit(4, z80.c)
}
/* BIT 4,D */
func instrCB__BIT_4_D(z80 *Z80) {
	z80.bit(4, z80.d)
}
/* BIT 4,E */
func instrCB__BIT_4_E(z80 *Z80) {
	z80.bit(4, z80.e)
}
/* BIT 4,H */
func instrCB__BIT_4_H(z80 *Z80) {
	z80.bit(4, z80.h)
}
/* BIT 4,L */
func instrCB__BIT_4_L(z80 *Z80) {
	z80.bit(4, z80.l)
}
/* BIT 4,(HL) */
func instrCB__BIT_4_iHL(z80 *Z80) {
	bytetemp := z80.memory.readByte(z80.HL())
	z80.memory.contendReadNoMreq(z80.HL(), 1)
	z80.bit(4, bytetemp)
}
/* BIT 4,A */
func instrCB__BIT_4_A(z80 *Z80) {
	z80.bit(4, z80.a)
}
/* BIT 5,B */
func instrCB__BIT_5_B(z80 *Z80) {
	z80.bit(5, z80.b)
}
/* BIT 5,C */
func instrCB__BIT_5_C(z80 *Z80) {
	z80.bit(5, z80.c)
}
/* BIT 5,D */
func instrCB__BIT_5_D(z80 *Z80) {
	z80.bit(5, z80.d)
}
/* BIT 5,E */
func instrCB__BIT_5_E(z80 *Z80) {
	z80.bit(5, z80.e)
}
/* BIT 5,H */
func instrCB__BIT_5_H(z80 *Z80) {
	z80.bit(5, z80.h)
}
/* BIT 5,L */
func instrCB__BIT_5_L(z80 *Z80) {
	z80.bit(5, z80.l)
}
/* BIT 5,(HL) */
func instrCB__BIT_5_iHL(z80 *Z80) {
	bytetemp := z80.memory.readByte(z80.HL())
	z80.memory.contendReadNoMreq(z80.HL(), 1)
	z80.bit(5, bytetemp)
}
/* BIT 5,A */
func instrCB__BIT_5_A(z80 *Z80) {
	z80.bit(5, z80.a)
}
/* BIT 6,B */
func instrCB__BIT_6_B(z80 *Z80) {
	z80.bit(6, z80.b)
}
/* BIT 6,C */
func instrCB__BIT_6_C(z80 *Z80) {
	z80.bit(6, z80.c)
}
/* BIT 6,D */
func instrCB__BIT_6_D(z80 *Z80) {
	z80.bit(6, z80.d)
}
/* BIT 6,E */
func instrCB__BIT_6_E(z80 *Z80) {
	z80.bit(6, z80.e)
}
/* BIT 6,H */
func instrCB__BIT_6_H(z80 *Z80) {
	z80.bit(6, z80.h)
}
/* BIT 6,L */
func instrCB__BIT_6_L(z80 *Z80) {
	z80.bit(6, z80.l)
}
/* BIT 6,(HL) */
func instrCB__BIT_6_iHL(z80 *Z80) {
	bytetemp := z80.memory.readByte(z80.HL())
	z80.memory.contendReadNoMreq(z80.HL(), 1)
	z80.bit(6, bytetemp)
}
/* BIT 6,A */
func instrCB__BIT_6_A(z80 *Z80) {
	z80.bit(6, z80.a)
}
/* BIT 7,B */
func instrCB__BIT_7_B(z80 *Z80) {
	z80.bit(7, z80.b)
}
/* BIT 7,C */
func instrCB__BIT_7_C(z80 *Z80) {
	z80.bit(7, z80.c)
}
/* BIT 7,D */
func instrCB__BIT_7_D(z80 *Z80) {
	z80.bit(7, z80.d)
}
/* BIT 7,E */
func instrCB__BIT_7_E(z80 *Z80) {
	z80.bit(7, z80.e)
}
/* BIT 7,H */
func instrCB__BIT_7_H(z80 *Z80) {
	z80.bit(7, z80.h)
}
/* BIT 7,L */
func instrCB__BIT_7_L(z80 *Z80) {
	z80.bit(7, z80.l)
}
/* BIT 7,(HL) */
func instrCB__BIT_7_iHL(z80 *Z80) {
	bytetemp := z80.memory.readByte(z80.HL())
	z80.memory.contendReadNoMreq(z80.HL(), 1)
	z80.bit(7, bytetemp)
}
/* BIT 7,A */
func instrCB__BIT_7_A(z80 *Z80) {
	z80.bit(7, z80.a)
}
/* RES 0,B */
func instrCB__RES_0_B(z80 *Z80) {
	z80.b &= 0xfe
}
/* RES 0,C */
func instrCB__RES_0_C(z80 *Z80) {
	z80.c &= 0xfe
}
/* RES 0,D */
func instrCB__RES_0_D(z80 *Z80) {
	z80.d &= 0xfe
}
/* RES 0,E */
func instrCB__RES_0_E(z80 *Z80) {
	z80.e &= 0xfe
}
/* RES 0,H */
func instrCB__RES_0_H(z80 *Z80) {
	z80.h &= 0xfe
}
/* RES 0,L */
func instrCB__RES_0_L(z80 *Z80) {
	z80.l &= 0xfe
}
/* RES 0,(HL) */
func instrCB__RES_0_iHL(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.HL())
	z80.memory.contendReadNoMreq(z80.HL(), 1)
	z80.memory.writeByte(z80.HL(), bytetemp&0xfe)
}
/* RES 0,A */
func instrCB__RES_0_A(z80 *Z80) {
	z80.a &= 0xfe
}
/* RES 1,B */
func instrCB__RES_1_B(z80 *Z80) {
	z80.b &= 0xfd
}
/* RES 1,C */
func instrCB__RES_1_C(z80 *Z80) {
	z80.c &= 0xfd
}
/* RES 1,D */
func instrCB__RES_1_D(z80 *Z80) {
	z80.d &= 0xfd
}
/* RES 1,E */
func instrCB__RES_1_E(z80 *Z80) {
	z80.e &= 0xfd
}
/* RES 1,H */
func instrCB__RES_1_H(z80 *Z80) {
	z80.h &= 0xfd
}
/* RES 1,L */
func instrCB__RES_1_L(z80 *Z80) {
	z80.l &= 0xfd
}
/* RES 1,(HL) */
func instrCB__RES_1_iHL(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.HL())
	z80.memory.contendReadNoMreq(z80.HL(), 1)
	z80.memory.writeByte(z80.HL(), bytetemp&0xfd)
}
/* RES 1,A */
func instrCB__RES_1_A(z80 *Z80) {
	z80.a &= 0xfd
}
/* RES 2,B */
func instrCB__RES_2_B(z80 *Z80) {
	z80.b &= 0xfb
}
/* RES 2,C */
func instrCB__RES_2_C(z80 *Z80) {
	z80.c &= 0xfb
}
/* RES 2,D */
func instrCB__RES_2_D(z80 *Z80) {
	z80.d &= 0xfb
}
/* RES 2,E */
func instrCB__RES_2_E(z80 *Z80) {
	z80.e &= 0xfb
}
/* RES 2,H */
func instrCB__RES_2_H(z80 *Z80) {
	z80.h &= 0xfb
}
/* RES 2,L */
func instrCB__RES_2_L(z80 *Z80) {
	z80.l &= 0xfb
}
/* RES 2,(HL) */
func instrCB__RES_2_iHL(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.HL())
	z80.memory.contendReadNoMreq(z80.HL(), 1)
	z80.memory.writeByte(z80.HL(), bytetemp&0xfb)
}
/* RES 2,A */
func instrCB__RES_2_A(z80 *Z80) {
	z80.a &= 0xfb
}
/* RES 3,B */
func instrCB__RES_3_B(z80 *Z80) {
	z80.b &= 0xf7
}
/* RES 3,C */
func instrCB__RES_3_C(z80 *Z80) {
	z80.c &= 0xf7
}
/* RES 3,D */
func instrCB__RES_3_D(z80 *Z80) {
	z80.d &= 0xf7
}
/* RES 3,E */
func instrCB__RES_3_E(z80 *Z80) {
	z80.e &= 0xf7
}
/* RES 3,H */
func instrCB__RES_3_H(z80 *Z80) {
	z80.h &= 0xf7
}
/* RES 3,L */
func instrCB__RES_3_L(z80 *Z80) {
	z80.l &= 0xf7
}
/* RES 3,(HL) */
func instrCB__RES_3_iHL(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.HL())
	z80.memory.contendReadNoMreq(z80.HL(), 1)
	z80.memory.writeByte(z80.HL(), bytetemp&0xf7)
}
/* RES 3,A */
func instrCB__RES_3_A(z80 *Z80) {
	z80.a &= 0xf7
}
/* RES 4,B */
func instrCB__RES_4_B(z80 *Z80) {
	z80.b &= 0xef
}
/* RES 4,C */
func instrCB__RES_4_C(z80 *Z80) {
	z80.c &= 0xef
}
/* RES 4,D */
func instrCB__RES_4_D(z80 *Z80) {
	z80.d &= 0xef
}
/* RES 4,E */
func instrCB__RES_4_E(z80 *Z80) {
	z80.e &= 0xef
}
/* RES 4,H */
func instrCB__RES_4_H(z80 *Z80) {
	z80.h &= 0xef
}
/* RES 4,L */
func instrCB__RES_4_L(z80 *Z80) {
	z80.l &= 0xef
}
/* RES 4,(HL) */
func instrCB__RES_4_iHL(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.HL())
	z80.memory.contendReadNoMreq(z80.HL(), 1)
	z80.memory.writeByte(z80.HL(), bytetemp&0xef)
}
/* RES 4,A */
func instrCB__RES_4_A(z80 *Z80) {
	z80.a &= 0xef
}
/* RES 5,B */
func instrCB__RES_5_B(z80 *Z80) {
	z80.b &= 0xdf
}
/* RES 5,C */
func instrCB__RES_5_C(z80 *Z80) {
	z80.c &= 0xdf
}
/* RES 5,D */
func instrCB__RES_5_D(z80 *Z80) {
	z80.d &= 0xdf
}
/* RES 5,E */
func instrCB__RES_5_E(z80 *Z80) {
	z80.e &= 0xdf
}
/* RES 5,H */
func instrCB__RES_5_H(z80 *Z80) {
	z80.h &= 0xdf
}
/* RES 5,L */
func instrCB__RES_5_L(z80 *Z80) {
	z80.l &= 0xdf
}
/* RES 5,(HL) */
func instrCB__RES_5_iHL(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.HL())
	z80.memory.contendReadNoMreq(z80.HL(), 1)
	z80.memory.writeByte(z80.HL(), bytetemp&0xdf)
}
/* RES 5,A */
func instrCB__RES_5_A(z80 *Z80) {
	z80.a &= 0xdf
}
/* RES 6,B */
func instrCB__RES_6_B(z80 *Z80) {
	z80.b &= 0xbf
}
/* RES 6,C */
func instrCB__RES_6_C(z80 *Z80) {
	z80.c &= 0xbf
}
/* RES 6,D */
func instrCB__RES_6_D(z80 *Z80) {
	z80.d &= 0xbf
}
/* RES 6,E */
func instrCB__RES_6_E(z80 *Z80) {
	z80.e &= 0xbf
}
/* RES 6,H */
func instrCB__RES_6_H(z80 *Z80) {
	z80.h &= 0xbf
}
/* RES 6,L */
func instrCB__RES_6_L(z80 *Z80) {
	z80.l &= 0xbf
}
/* RES 6,(HL) */
func instrCB__RES_6_iHL(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.HL())
	z80.memory.contendReadNoMreq(z80.HL(), 1)
	z80.memory.writeByte(z80.HL(), bytetemp&0xbf)
}
/* RES 6,A */
func instrCB__RES_6_A(z80 *Z80) {
	z80.a &= 0xbf
}
/* RES 7,B */
func instrCB__RES_7_B(z80 *Z80) {
	z80.b &= 0x7f
}
/* RES 7,C */
func instrCB__RES_7_C(z80 *Z80) {
	z80.c &= 0x7f
}
/* RES 7,D */
func instrCB__RES_7_D(z80 *Z80) {
	z80.d &= 0x7f
}
/* RES 7,E */
func instrCB__RES_7_E(z80 *Z80) {
	z80.e &= 0x7f
}
/* RES 7,H */
func instrCB__RES_7_H(z80 *Z80) {
	z80.h &= 0x7f
}
/* RES 7,L */
func instrCB__RES_7_L(z80 *Z80) {
	z80.l &= 0x7f
}
/* RES 7,(HL) */
func instrCB__RES_7_iHL(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.HL())
	z80.memory.contendReadNoMreq(z80.HL(), 1)
	z80.memory.writeByte(z80.HL(), bytetemp&0x7f)
}
/* RES 7,A */
func instrCB__RES_7_A(z80 *Z80) {
	z80.a &= 0x7f
}
/* SET 0,B */
func instrCB__SET_0_B(z80 *Z80) {
	z80.b |= 0x01
}
/* SET 0,C */
func instrCB__SET_0_C(z80 *Z80) {
	z80.c |= 0x01
}
/* SET 0,D */
func instrCB__SET_0_D(z80 *Z80) {
	z80.d |= 0x01
}
/* SET 0,E */
func instrCB__SET_0_E(z80 *Z80) {
	z80.e |= 0x01
}
/* SET 0,H */
func instrCB__SET_0_H(z80 *Z80) {
	z80.h |= 0x01
}
/* SET 0,L */
func instrCB__SET_0_L(z80 *Z80) {
	z80.l |= 0x01
}
/* SET 0,(HL) */
func instrCB__SET_0_iHL(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.HL())
	z80.memory.contendReadNoMreq(z80.HL(), 1)
	z80.memory.writeByte(z80.HL(), bytetemp|0x01)
}
/* SET 0,A */
func instrCB__SET_0_A(z80 *Z80) {
	z80.a |= 0x01
}
/* SET 1,B */
func instrCB__SET_1_B(z80 *Z80) {
	z80.b |= 0x02
}
/* SET 1,C */
func instrCB__SET_1_C(z80 *Z80) {
	z80.c |= 0x02
}
/* SET 1,D */
func instrCB__SET_1_D(z80 *Z80) {
	z80.d |= 0x02
}
/* SET 1,E */
func instrCB__SET_1_E(z80 *Z80) {
	z80.e |= 0x02
}
/* SET 1,H */
func instrCB__SET_1_H(z80 *Z80) {
	z80.h |= 0x02
}
/* SET 1,L */
func instrCB__SET_1_L(z80 *Z80) {
	z80.l |= 0x02
}
/* SET 1,(HL) */
func instrCB__SET_1_iHL(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.HL())
	z80.memory.contendReadNoMreq(z80.HL(), 1)
	z80.memory.writeByte(z80.HL(), bytetemp|0x02)
}
/* SET 1,A */
func instrCB__SET_1_A(z80 *Z80) {
	z80.a |= 0x02
}
/* SET 2,B */
func instrCB__SET_2_B(z80 *Z80) {
	z80.b |= 0x04
}
/* SET 2,C */
func instrCB__SET_2_C(z80 *Z80) {
	z80.c |= 0x04
}
/* SET 2,D */
func instrCB__SET_2_D(z80 *Z80) {
	z80.d |= 0x04
}
/* SET 2,E */
func instrCB__SET_2_E(z80 *Z80) {
	z80.e |= 0x04
}
/* SET 2,H */
func instrCB__SET_2_H(z80 *Z80) {
	z80.h |= 0x04
}
/* SET 2,L */
func instrCB__SET_2_L(z80 *Z80) {
	z80.l |= 0x04
}
/* SET 2,(HL) */
func instrCB__SET_2_iHL(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.HL())
	z80.memory.contendReadNoMreq(z80.HL(), 1)
	z80.memory.writeByte(z80.HL(), bytetemp|0x04)
}
/* SET 2,A */
func instrCB__SET_2_A(z80 *Z80) {
	z80.a |= 0x04
}
/* SET 3,B */
func instrCB__SET_3_B(z80 *Z80) {
	z80.b |= 0x08
}
/* SET 3,C */
func instrCB__SET_3_C(z80 *Z80) {
	z80.c |= 0x08
}
/* SET 3,D */
func instrCB__SET_3_D(z80 *Z80) {
	z80.d |= 0x08
}
/* SET 3,E */
func instrCB__SET_3_E(z80 *Z80) {
	z80.e |= 0x08
}
/* SET 3,H */
func instrCB__SET_3_H(z80 *Z80) {
	z80.h |= 0x08
}
/* SET 3,L */
func instrCB__SET_3_L(z80 *Z80) {
	z80.l |= 0x08
}
/* SET 3,(HL) */
func instrCB__SET_3_iHL(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.HL())
	z80.memory.contendReadNoMreq(z80.HL(), 1)
	z80.memory.writeByte(z80.HL(), bytetemp|0x08)
}
/* SET 3,A */
func instrCB__SET_3_A(z80 *Z80) {
	z80.a |= 0x08
}
/* SET 4,B */
func instrCB__SET_4_B(z80 *Z80) {
	z80.b |= 0x10
}
/* SET 4,C */
func instrCB__SET_4_C(z80 *Z80) {
	z80.c |= 0x10
}
/* SET 4,D */
func instrCB__SET_4_D(z80 *Z80) {
	z80.d |= 0x10
}
/* SET 4,E */
func instrCB__SET_4_E(z80 *Z80) {
	z80.e |= 0x10
}
/* SET 4,H */
func instrCB__SET_4_H(z80 *Z80) {
	z80.h |= 0x10
}
/* SET 4,L */
func instrCB__SET_4_L(z80 *Z80) {
	z80.l |= 0x10
}
/* SET 4,(HL) */
func instrCB__SET_4_iHL(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.HL())
	z80.memory.contendReadNoMreq(z80.HL(), 1)
	z80.memory.writeByte(z80.HL(), bytetemp|0x10)
}
/* SET 4,A */
func instrCB__SET_4_A(z80 *Z80) {
	z80.a |= 0x10
}
/* SET 5,B */
func instrCB__SET_5_B(z80 *Z80) {
	z80.b |= 0x20
}
/* SET 5,C */
func instrCB__SET_5_C(z80 *Z80) {
	z80.c |= 0x20
}
/* SET 5,D */
func instrCB__SET_5_D(z80 *Z80) {
	z80.d |= 0x20
}
/* SET 5,E */
func instrCB__SET_5_E(z80 *Z80) {
	z80.e |= 0x20
}
/* SET 5,H */
func instrCB__SET_5_H(z80 *Z80) {
	z80.h |= 0x20
}
/* SET 5,L */
func instrCB__SET_5_L(z80 *Z80) {
	z80.l |= 0x20
}
/* SET 5,(HL) */
func instrCB__SET_5_iHL(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.HL())
	z80.memory.contendReadNoMreq(z80.HL(), 1)
	z80.memory.writeByte(z80.HL(), bytetemp|0x20)
}
/* SET 5,A */
func instrCB__SET_5_A(z80 *Z80) {
	z80.a |= 0x20
}
/* SET 6,B */
func instrCB__SET_6_B(z80 *Z80) {
	z80.b |= 0x40
}
/* SET 6,C */
func instrCB__SET_6_C(z80 *Z80) {
	z80.c |= 0x40
}
/* SET 6,D */
func instrCB__SET_6_D(z80 *Z80) {
	z80.d |= 0x40
}
/* SET 6,E */
func instrCB__SET_6_E(z80 *Z80) {
	z80.e |= 0x40
}
/* SET 6,H */
func instrCB__SET_6_H(z80 *Z80) {
	z80.h |= 0x40
}
/* SET 6,L */
func instrCB__SET_6_L(z80 *Z80) {
	z80.l |= 0x40
}
/* SET 6,(HL) */
func instrCB__SET_6_iHL(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.HL())
	z80.memory.contendReadNoMreq(z80.HL(), 1)
	z80.memory.writeByte(z80.HL(), bytetemp|0x40)
}
/* SET 6,A */
func instrCB__SET_6_A(z80 *Z80) {
	z80.a |= 0x40
}
/* SET 7,B */
func instrCB__SET_7_B(z80 *Z80) {
	z80.b |= 0x80
}
/* SET 7,C */
func instrCB__SET_7_C(z80 *Z80) {
	z80.c |= 0x80
}
/* SET 7,D */
func instrCB__SET_7_D(z80 *Z80) {
	z80.d |= 0x80
}
/* SET 7,E */
func instrCB__SET_7_E(z80 *Z80) {
	z80.e |= 0x80
}
/* SET 7,H */
func instrCB__SET_7_H(z80 *Z80) {
	z80.h |= 0x80
}
/* SET 7,L */
func instrCB__SET_7_L(z80 *Z80) {
	z80.l |= 0x80
}
/* SET 7,(HL) */
func instrCB__SET_7_iHL(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.HL())
	z80.memory.contendReadNoMreq(z80.HL(), 1)
	z80.memory.writeByte(z80.HL(), bytetemp|0x80)
}
/* SET 7,A */
func instrCB__SET_7_A(z80 *Z80) {
	z80.a |= 0x80
}
/* IN B,(C) */
func instrED__IN_B_iC(z80 *Z80) {
	z80.in(&z80.b, z80.BC())
}
/* OUT (C),B */
func instrED__OUT_iC_B(z80 *Z80) {
	z80.writePort(z80.BC(), z80.b)
}
/* SBC HL,BC */
func instrED__SBC_HL_BC(z80 *Z80) {
	z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 7)
	z80.sbc16(z80.BC())
}
/* LD (nnnn),BC */
func instrED__LD_iNNNN_BC(z80 *Z80) {
	z80.ld16nnrr(z80.c, z80.b)
	// break
}
/* NEG */
func instrED__NEG(z80 *Z80) {
	bytetemp := z80.a
	z80.a = 0
	z80.sub(bytetemp)
}
/* RETN */
func instrED__RETN(z80 *Z80) {
	z80.iff1 = z80.iff2
	z80.ret()
}
/* IM 0 */
func instrED__IM_0(z80 *Z80) {
	z80.im = 0
}
/* LD I,A */
func instrED__LD_I_A(z80 *Z80) {
	z80.memory.contendReadNoMreq(z80.IR(), 1)
	z80.i = z80.a
}
/* IN C,(C) */
func instrED__IN_C_iC(z80 *Z80) {
	z80.in(&z80.c, z80.BC())
}
/* OUT (C),C */
func instrED__OUT_iC_C(z80 *Z80) {
	z80.writePort(z80.BC(), z80.c)
}
/* ADC HL,BC */
func instrED__ADC_HL_BC(z80 *Z80) {
	z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 7)
	z80.adc16(z80.BC())
}
/* LD BC,(nnnn) */
func instrED__LD_BC_iNNNN(z80 *Z80) {
	z80.ld16rrnn(&z80.c, &z80.b)
	// break
}
/* LD R,A */
func instrED__LD_R_A(z80 *Z80) {
	z80.memory.contendReadNoMreq(z80.IR(), 1)
	/* Keep the RZX instruction counter right */
	z80.rzxInstructionsOffset += (int(z80.r) - int(z80.a))
	z80.r, z80.r7 = uint16(z80.a), z80.a
}
/* IN D,(C) */
func instrED__IN_D_iC(z80 *Z80) {
	z80.in(&z80.d, z80.BC())
}
/* OUT (C),D */
func instrED__OUT_iC_D(z80 *Z80) {
	z80.writePort(z80.BC(), z80.d)
}
/* SBC HL,DE */
func instrED__SBC_HL_DE(z80 *Z80) {
	z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 7)
	z80.sbc16(z80.DE())
}
/* LD (nnnn),DE */
func instrED__LD_iNNNN_DE(z80 *Z80) {
	z80.ld16nnrr(z80.e, z80.d)
	// break
}
/* IM 1 */
func instrED__IM_1(z80 *Z80) {
	z80.im = 1
}
/* LD A,I */
func instrED__LD_A_I(z80 *Z80) {
	z80.memory.contendReadNoMreq(z80.IR(), 1)
	z80.a = z80.i
	z80.f = (z80.f & FLAG_C) | sz53Table[z80.a] | ternOpB(z80.iff2 != 0, FLAG_V, 0)
}
/* IN E,(C) */
func instrED__IN_E_iC(z80 *Z80) {
	z80.in(&z80.e, z80.BC())
}
/* OUT (C),E */
func instrED__OUT_iC_E(z80 *Z80) {
	z80.writePort(z80.BC(), z80.e)
}
/* ADC HL,DE */
func instrED__ADC_HL_DE(z80 *Z80) {
	z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 7)
	z80.adc16(z80.DE())
}
/* LD DE,(nnnn) */
func instrED__LD_DE_iNNNN(z80 *Z80) {
	z80.ld16rrnn(&z80.e, &z80.d)
	// break
}
/* IM 2 */
func instrED__IM_2(z80 *Z80) {
	z80.im = 2
}
/* LD A,R */
func instrED__LD_A_R(z80 *Z80) {
	z80.memory.contendReadNoMreq(z80.IR(), 1)
	z80.a = byte(z80.r&0x7f) | (z80.r7 & 0x80)
	z80.f = (z80.f & FLAG_C) | sz53Table[z80.a] | ternOpB(z80.iff2 != 0, FLAG_V, 0)
}
/* IN H,(C) */
func instrED__IN_H_iC(z80 *Z80) {
	z80.in(&z80.h, z80.BC())
}
/* OUT (C),H */
func instrED__OUT_iC_H(z80 *Z80) {
	z80.writePort(z80.BC(), z80.h)
}
/* SBC HL,HL */
func instrED__SBC_HL_HL(z80 *Z80) {
	z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 7)
	z80.sbc16(z80.HL())
}
/* LD (nnnn),HL */
func instrED__LD_iNNNN_HL(z80 *Z80) {
	z80.ld16nnrr(z80.l, z80.h)
	// break
}
/* RRD */
func instrED__RRD(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.HL())
	z80.memory.contendReadNoMreq_loop(z80.HL(), 1, 4)
	z80.memory.writeByte(z80.HL(), (z80.a<<4)|(bytetemp>>4))
	z80.a = (z80.a & 0xf0) | (bytetemp & 0x0f)
	z80.f = (z80.f & FLAG_C) | sz53pTable[z80.a]
}
/* IN L,(C) */
func instrED__IN_L_iC(z80 *Z80) {
	z80.in(&z80.l, z80.BC())
}
/* OUT (C),L */
func instrED__OUT_iC_L(z80 *Z80) {
	z80.writePort(z80.BC(), z80.l)
}
/* ADC HL,HL */
func instrED__ADC_HL_HL(z80 *Z80) {
	z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 7)
	z80.adc16(z80.HL())
}
/* LD HL,(nnnn) */
func instrED__LD_HL_iNNNN(z80 *Z80) {
	z80.ld16rrnn(&z80.l, &z80.h)
	// break
}
/* RLD */
func instrED__RLD(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.HL())
	z80.memory.contendReadNoMreq_loop(z80.HL(), 1, 4)
	z80.memory.writeByte(z80.HL(), (bytetemp<<4)|(z80.a&0x0f))
	z80.a = (z80.a & 0xf0) | (bytetemp >> 4)
	z80.f = (z80.f & FLAG_C) | sz53pTable[z80.a]
}
/* IN F,(C) */
func instrED__IN_F_iC(z80 *Z80) {
	var bytetemp byte
	z80.in(&bytetemp, z80.BC())
}
/* OUT (C),0 */
func instrED__OUT_iC_0(z80 *Z80) {
	z80.writePort(z80.BC(), 0)
}
/* SBC HL,SP */
func instrED__SBC_HL_SP(z80 *Z80) {
	z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 7)
	z80.sbc16(z80.SP())
}
/* LD (nnnn),SP */
func instrED__LD_iNNNN_SP(z80 *Z80) {
	sph, spl := splitWord(z80.sp)
	z80.ld16nnrr(spl, sph)
	// break
}
/* IN A,(C) */
func instrED__IN_A_iC(z80 *Z80) {
	z80.in(&z80.a, z80.BC())
}
/* OUT (C),A */
func instrED__OUT_iC_A(z80 *Z80) {
	z80.writePort(z80.BC(), z80.a)
}
/* ADC HL,SP */
func instrED__ADC_HL_SP(z80 *Z80) {
	z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 7)
	z80.adc16(z80.SP())
}
/* LD SP,(nnnn) */
func instrED__LD_SP_iNNNN(z80 *Z80) {
	sph, spl := splitWord(z80.sp)
	z80.ld16rrnn(&spl, &sph)
	z80.sp = joinBytes(sph, spl)
	// break
}
/* LDI */
func instrED__LDI(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.HL())
	z80.decBC()
	z80.memory.writeByte(z80.DE(), bytetemp)
	z80.memory.contendWriteNoMreq_loop(z80.DE(), 1, 2)
	z80.incDE()
	z80.incHL()
	bytetemp += z80.a
	z80.f = (z80.f & (FLAG_C | FLAG_Z | FLAG_S)) |
		ternOpB(z80.BC() != 0, FLAG_V, 0) |
		(bytetemp & FLAG_3) |
		ternOpB((bytetemp&0x02) != 0, FLAG_5, 0)
}
/* CPI */
func instrED__CPI(z80 *Z80) {
	var value byte = z80.memory.readByte(z80.HL())
	var bytetemp byte = z80.a - value
	var lookup byte = ((z80.a & 0x08) >> 3) | (((value) & 0x08) >> 2) | ((bytetemp & 0x08) >> 1)
	z80.memory.contendReadNoMreq_loop(z80.HL(), 1, 5)
	z80.incHL()
	z80.decBC()
	z80.f = (z80.f & FLAG_C) | ternOpB(z80.BC() != 0, FLAG_V|FLAG_N, FLAG_N) | halfcarrySubTable[lookup] | ternOpB(bytetemp != 0, 0, FLAG_Z) | (bytetemp & FLAG_S)
	if (z80.f & FLAG_H) != 0 {
		bytetemp--
	}
	z80.f |= (bytetemp & FLAG_3) | ternOpB((bytetemp&0x02) != 0, FLAG_5, 0)
}
/* INI */
func instrED__INI(z80 *Z80) {
	z80.memory.contendReadNoMreq(z80.IR(), 1)
	var initemp byte = z80.readPort(z80.BC())
	z80.memory.writeByte(z80.HL(), initemp)

	z80.b--
	z80.incHL()
	var initemp2 byte = initemp + z80.c + 1
	z80.f = ternOpB((initemp&0x80) != 0, FLAG_N, 0) |
		ternOpB(initemp2 < initemp, FLAG_H|FLAG_C, 0) |
		ternOpB(parityTable[(initemp2&0x07)^z80.b] != 0, FLAG_P, 0) |
		sz53Table[z80.b]
}
/* OUTI */
func instrED__OUTI(z80 *Z80) {
	z80.memory.contendReadNoMreq(z80.IR(), 1)
	var outitemp byte = z80.memory.readByte(z80.HL())
	z80.b-- /* This does happen first, despite what the specs say */
	z80.writePort(z80.BC(), outitemp)

	z80.incHL()
	var outitemp2 byte = outitemp + z80.l
	z80.f = ternOpB((outitemp&0x80) != 0, FLAG_N, 0) |
		ternOpB(outitemp2 < outitemp, FLAG_H|FLAG_C, 0) |
		ternOpB(parityTable[(outitemp2&0x07)^z80.b] != 0, FLAG_P, 0) |
		sz53Table[z80.b]
}
/* LDD */
func instrED__LDD(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.HL())
	z80.decBC()
	z80.memory.writeByte(z80.DE(), bytetemp)
	z80.memory.contendWriteNoMreq_loop(z80.DE(), 1, 2)
	z80.decDE()
	z80.decHL()
	bytetemp += z80.a
	z80.f = (z80.f & (FLAG_C | FLAG_Z | FLAG_S)) |
		ternOpB(z80.BC() != 0, FLAG_V, 0) |
		(bytetemp & FLAG_3) |
		ternOpB((bytetemp&0x02) != 0, FLAG_5, 0)
}
/* CPD */
func instrED__CPD(z80 *Z80) {
	var value byte = z80.memory.readByte(z80.HL())
	var bytetemp byte = z80.a - value
	var lookup byte = ((z80.a & 0x08) >> 3) | (((value) & 0x08) >> 2) | ((bytetemp & 0x08) >> 1)
	z80.memory.contendReadNoMreq_loop(z80.HL(), 1, 5)
	z80.decHL()
	z80.decBC()
	z80.f = (z80.f & FLAG_C) | ternOpB(z80.BC() != 0, FLAG_V|FLAG_N, FLAG_N) | halfcarrySubTable[lookup] | ternOpB(bytetemp != 0, 0, FLAG_Z) | (bytetemp & FLAG_S)
	if (z80.f & FLAG_H) != 0 {
		bytetemp--
	}
	z80.f |= (bytetemp & FLAG_3) | ternOpB((bytetemp&0x02) != 0, FLAG_5, 0)
}
/* IND */
func instrED__IND(z80 *Z80) {
	z80.memory.contendReadNoMreq(z80.IR(), 1)
	var initemp byte = z80.readPort(z80.BC())
	z80.memory.writeByte(z80.HL(), initemp)

	z80.b--
	z80.decHL()
	var initemp2 byte = initemp + z80.c - 1
	z80.f = ternOpB((initemp&0x80) != 0, FLAG_N, 0) |
		ternOpB(initemp2 < initemp, FLAG_H|FLAG_C, 0) |
		ternOpB(parityTable[(initemp2&0x07)^z80.b] != 0, FLAG_P, 0) |
		sz53Table[z80.b]
}
/* OUTD */
func instrED__OUTD(z80 *Z80) {
	z80.memory.contendReadNoMreq(z80.IR(), 1)
	var outitemp byte = z80.memory.readByte(z80.HL())
	z80.b-- /* This does happen first, despite what the specs say */
	z80.writePort(z80.BC(), outitemp)

	z80.decHL()
	var outitemp2 byte = outitemp + z80.l
	z80.f = ternOpB((outitemp&0x80) != 0, FLAG_N, 0) |
		ternOpB(outitemp2 < outitemp, FLAG_H|FLAG_C, 0) |
		ternOpB(parityTable[(outitemp2&0x07)^z80.b] != 0, FLAG_P, 0) |
		sz53Table[z80.b]
}
/* LDIR */
func instrED__LDIR(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.HL())
	z80.memory.writeByte(z80.DE(), bytetemp)
	z80.memory.contendWriteNoMreq_loop(z80.DE(), 1, 2)
	z80.decBC()
	bytetemp += z80.a
	z80.f = (z80.f & (FLAG_C | FLAG_Z | FLAG_S)) | ternOpB(z80.BC() != 0, FLAG_V, 0) | (bytetemp & FLAG_3) | ternOpB((bytetemp&0x02 != 0), FLAG_5, 0)
	if z80.BC() != 0 {
		z80.memory.contendWriteNoMreq_loop(z80.DE(), 1, 5)
		z80.pc -= 2
	}
	z80.incHL()
	z80.incDE()
}
/* CPIR */
func instrED__CPIR(z80 *Z80) {
	var value byte = z80.memory.readByte(z80.HL())
	var bytetemp byte = z80.a - value
	var lookup byte = ((z80.a & 0x08) >> 3) | (((value) & 0x08) >> 2) | ((bytetemp & 0x08) >> 1)
	z80.memory.contendReadNoMreq_loop(z80.HL(), 1, 5)
	z80.decBC()
	z80.f = (z80.f & FLAG_C) | (ternOpB(z80.BC() != 0, (FLAG_V | FLAG_N), FLAG_N)) | halfcarrySubTable[lookup] | (ternOpB(bytetemp != 0, 0, FLAG_Z)) | (bytetemp & FLAG_S)
	if (z80.f & FLAG_H) != 0 {
		bytetemp--
	}
	z80.f |= (bytetemp & FLAG_3) | ternOpB((bytetemp&0x02) != 0, FLAG_5, 0)
	if (z80.f & (FLAG_V | FLAG_Z)) == FLAG_V {
		z80.memory.contendReadNoMreq_loop(z80.HL(), 1, 5)
		z80.pc -= 2
	}
	z80.incHL()
}
/* INIR */
func instrED__INIR(z80 *Z80) {
	z80.memory.contendReadNoMreq(z80.IR(), 1)
	var initemp byte = z80.readPort(z80.BC())
	z80.memory.writeByte(z80.HL(), initemp)

	z80.b--
	var initemp2 byte = initemp + z80.c + 1
	z80.f = ternOpB(initemp&0x80 != 0, FLAG_N, 0) |
		ternOpB(initemp2 < initemp, FLAG_H|FLAG_C, 0) |
		ternOpB(parityTable[(initemp2&0x07)^z80.b] != 0, FLAG_P, 0) |
		sz53Table[z80.b]

	if z80.b != 0 {
		z80.memory.contendWriteNoMreq_loop(z80.HL(), 1, 5)
		z80.pc -= 2
	}
	z80.incHL()
}
/* OTIR */
func instrED__OTIR(z80 *Z80) {
	z80.memory.contendReadNoMreq(z80.IR(), 1)
	var outitemp byte = z80.memory.readByte(z80.HL())
	z80.b-- /* This does happen first, despite what the specs say */
	z80.writePort(z80.BC(), outitemp)

	z80.incHL()
	var outitemp2 byte = outitemp + z80.l
	z80.f = ternOpB((outitemp&0x80) != 0, FLAG_N, 0) |
		ternOpB(outitemp2 < outitemp, FLAG_H|FLAG_C, 0) |
		ternOpB(parityTable[(outitemp2&0x07)^z80.b] != 0, FLAG_P, 0) |
		sz53Table[z80.b]

	if z80.b != 0 {
		z80.memory.contendReadNoMreq_loop(z80.BC(), 1, 5)
		z80.pc -= 2
	}
}
/* LDDR */
func instrED__LDDR(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.HL())
	z80.memory.writeByte(z80.DE(), bytetemp)
	z80.memory.contendWriteNoMreq_loop(z80.DE(), 1, 2)
	z80.decBC()
	bytetemp += z80.a
	z80.f = (z80.f & (FLAG_C | FLAG_Z | FLAG_S)) | ternOpB(z80.BC() != 0, FLAG_V, 0) | (bytetemp & FLAG_3) | ternOpB((bytetemp&0x02 != 0), FLAG_5, 0)
	if z80.BC() != 0 {
		z80.memory.contendWriteNoMreq_loop(z80.DE(), 1, 5)
		z80.pc -= 2
	}
	z80.decHL()
	z80.decDE()
}
/* CPDR */
func instrED__CPDR(z80 *Z80) {
	var value byte = z80.memory.readByte(z80.HL())
	var bytetemp byte = z80.a - value
	var lookup byte = ((z80.a & 0x08) >> 3) | (((value) & 0x08) >> 2) | ((bytetemp & 0x08) >> 1)
	z80.memory.contendReadNoMreq_loop(z80.HL(), 1, 5)
	z80.decBC()
	z80.f = (z80.f & FLAG_C) | (ternOpB(z80.BC() != 0, (FLAG_V | FLAG_N), FLAG_N)) | halfcarrySubTable[lookup] | (ternOpB(bytetemp != 0, 0, FLAG_Z)) | (bytetemp & FLAG_S)
	if (z80.f & FLAG_H) != 0 {
		bytetemp--
	}
	z80.f |= (bytetemp & FLAG_3) | ternOpB((bytetemp&0x02) != 0, FLAG_5, 0)
	if (z80.f & (FLAG_V | FLAG_Z)) == FLAG_V {
		z80.memory.contendReadNoMreq_loop(z80.HL(), 1, 5)
		z80.pc -= 2
	}
	z80.decHL()
}
/* INDR */
func instrED__INDR(z80 *Z80) {
	z80.memory.contendReadNoMreq(z80.IR(), 1)
	var initemp byte = z80.readPort(z80.BC())
	z80.memory.writeByte(z80.HL(), initemp)

	z80.b--
	var initemp2 byte = initemp + z80.c - 1
	z80.f = ternOpB(initemp&0x80 != 0, FLAG_N, 0) |
		ternOpB(initemp2 < initemp, FLAG_H|FLAG_C, 0) |
		ternOpB(parityTable[(initemp2&0x07)^z80.b] != 0, FLAG_P, 0) |
		sz53Table[z80.b]

	if z80.b != 0 {
		z80.memory.contendWriteNoMreq_loop(z80.HL(), 1, 5)
		z80.pc -= 2
	}
	z80.decHL()
}
/* OTDR */
func instrED__OTDR(z80 *Z80) {
	z80.memory.contendReadNoMreq(z80.IR(), 1)
	var outitemp byte = z80.memory.readByte(z80.HL())
	z80.b-- /* This does happen first, despite what the specs say */
	z80.writePort(z80.BC(), outitemp)

	z80.decHL()
	var outitemp2 byte = outitemp + z80.l
	z80.f = ternOpB((outitemp&0x80) != 0, FLAG_N, 0) |
		ternOpB(outitemp2 < outitemp, FLAG_H|FLAG_C, 0) |
		ternOpB(parityTable[(outitemp2&0x07)^z80.b] != 0, FLAG_P, 0) |
		sz53Table[z80.b]

	if z80.b != 0 {
		z80.memory.contendReadNoMreq_loop(z80.BC(), 1, 5)
		z80.pc -= 2
	}
}
/* slttrap */
func instrED__SLTTRAP(z80 *Z80) {
	z80.sltTrap(int16(z80.HL()), z80.a)
}
/* ADD ix,BC */
func instrDD__ADD_REG_BC(z80 *Z80) {
	z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 7)
	z80.add16(z80.ix, z80.BC())
}
/* ADD ix,DE */
func instrDD__ADD_REG_DE(z80 *Z80) {
	z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 7)
	z80.add16(z80.ix, z80.DE())
}
/* LD ix,nnnn */
func instrDD__LD_REG_NNNN(z80 *Z80) {
	b1 := z80.memory.readByte(z80.pc)
	z80.pc++
	b2 := z80.memory.readByte(z80.pc)
	z80.pc++
	z80.setIX(joinBytes(b2, b1))
}
/* LD (nnnn),ix */
func instrDD__LD_iNNNN_REG(z80 *Z80) {
	z80.ld16nnrr(z80.ixl, z80.ixh)
	// break
}
/* INC ix */
func instrDD__INC_REG(z80 *Z80) {
	z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 2)
	z80.incIX()
}
/* INC z80.IXH() */
func instrDD__INC_REGH(z80 *Z80) {
	z80.incIXH()
}
/* DEC z80.IXH() */
func instrDD__DEC_REGH(z80 *Z80) {
	z80.decIXH()
}
/* LD z80.IXH(),nn */
func instrDD__LD_REGH_NN(z80 *Z80) {
	z80.ixh = z80.memory.readByte(z80.pc)
	z80.pc++
}
/* ADD ix,ix */
func instrDD__ADD_REG_REG(z80 *Z80) {
	z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 7)
	z80.add16(z80.ix, z80.IX())
}
/* LD ix,(nnnn) */
func instrDD__LD_REG_iNNNN(z80 *Z80) {
	z80.ld16rrnn(&z80.ixl, &z80.ixh)
	// break
}
/* DEC ix */
func instrDD__DEC_REG(z80 *Z80) {
	z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 2)
	z80.decIX()
}
/* INC z80.IXL() */
func instrDD__INC_REGL(z80 *Z80) {
	z80.incIXL()
}
/* DEC z80.IXL() */
func instrDD__DEC_REGL(z80 *Z80) {
	z80.decIXL()
}
/* LD z80.IXL(),nn */
func instrDD__LD_REGL_NN(z80 *Z80) {
	z80.ixl = z80.memory.readByte(z80.pc)
	z80.pc++
}
/* INC (ix+dd) */
func instrDD__INC_iREGpDD(z80 *Z80) {
	var offset byte = z80.memory.readByte(z80.pc)
	z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
	z80.pc++
	var wordtemp uint16 = z80.IX() + uint16(signExtend(offset))
	var bytetemp byte = z80.memory.readByte(wordtemp)
	z80.memory.contendReadNoMreq(wordtemp, 1)
	z80.inc(&bytetemp)
	z80.memory.writeByte(wordtemp, bytetemp)
}
/* DEC (ix+dd) */
func instrDD__DEC_iREGpDD(z80 *Z80) {
	var offset byte = z80.memory.readByte(z80.pc)
	z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
	z80.pc++
	var wordtemp uint16 = z80.IX() + uint16(signExtend(offset))
	var bytetemp byte = z80.memory.readByte(wordtemp)
	z80.memory.contendReadNoMreq(wordtemp, 1)
	z80.dec(&bytetemp)
	z80.memory.writeByte(wordtemp, bytetemp)
}
/* LD (ix+dd),nn */
func instrDD__LD_iREGpDD_NN(z80 *Z80) {
	offset := z80.memory.readByte(z80.pc)
	z80.pc++
	value := z80.memory.readByte(z80.pc)
	z80.memory.contendReadNoMreq_loop(z80.pc, 1, 2)
	z80.pc++
	z80.memory.writeByte(z80.IX()+uint16(signExtend(offset)), value)
}
/* ADD ix,SP */
func instrDD__ADD_REG_SP(z80 *Z80) {
	z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 7)
	z80.add16(z80.ix, z80.SP())
}
/* LD B,z80.IXH() */
func instrDD__LD_B_REGH(z80 *Z80) {
	z80.b = z80.ixh
}
/* LD B,z80.IXL() */
func instrDD__LD_B_REGL(z80 *Z80) {
	z80.b = z80.ixl
}
/* LD B,(ix+dd) */
func instrDD__LD_B_iREGpDD(z80 *Z80) {
	var offset byte = z80.memory.readByte(z80.pc)
	z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
	z80.pc++
	z80.b = z80.memory.readByte(z80.IX() + uint16(signExtend(offset)))
}
/* LD C,z80.IXH() */
func instrDD__LD_C_REGH(z80 *Z80) {
	z80.c = z80.ixh
}
/* LD C,z80.IXL() */
func instrDD__LD_C_REGL(z80 *Z80) {
	z80.c = z80.ixl
}
/* LD C,(ix+dd) */
func instrDD__LD_C_iREGpDD(z80 *Z80) {
	var offset byte = z80.memory.readByte(z80.pc)
	z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
	z80.pc++
	z80.c = z80.memory.readByte(z80.IX() + uint16(signExtend(offset)))
}
/* LD D,z80.IXH() */
func instrDD__LD_D_REGH(z80 *Z80) {
	z80.d = z80.ixh
}
/* LD D,z80.IXL() */
func instrDD__LD_D_REGL(z80 *Z80) {
	z80.d = z80.ixl
}
/* LD D,(ix+dd) */
func instrDD__LD_D_iREGpDD(z80 *Z80) {
	var offset byte = z80.memory.readByte(z80.pc)
	z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
	z80.pc++
	z80.d = z80.memory.readByte(z80.IX() + uint16(signExtend(offset)))
}
/* LD E,z80.IXH() */
func instrDD__LD_E_REGH(z80 *Z80) {
	z80.e = z80.ixh
}
/* LD E,z80.IXL() */
func instrDD__LD_E_REGL(z80 *Z80) {
	z80.e = z80.ixl
}
/* LD E,(ix+dd) */
func instrDD__LD_E_iREGpDD(z80 *Z80) {
	var offset byte = z80.memory.readByte(z80.pc)
	z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
	z80.pc++
	z80.e = z80.memory.readByte(z80.IX() + uint16(signExtend(offset)))
}
/* LD z80.IXH(),B */
func instrDD__LD_REGH_B(z80 *Z80) {
	z80.ixh = z80.b
}
/* LD z80.IXH(),C */
func instrDD__LD_REGH_C(z80 *Z80) {
	z80.ixh = z80.c
}
/* LD z80.IXH(),D */
func instrDD__LD_REGH_D(z80 *Z80) {
	z80.ixh = z80.d
}
/* LD z80.IXH(),E */
func instrDD__LD_REGH_E(z80 *Z80) {
	z80.ixh = z80.e
}
/* LD z80.IXH(),z80.IXH() */
func instrDD__LD_REGH_REGH(z80 *Z80) {
}
/* LD z80.IXH(),z80.IXL() */
func instrDD__LD_REGH_REGL(z80 *Z80) {
	z80.ixh = z80.ixl
}
/* LD H,(ix+dd) */
func instrDD__LD_H_iREGpDD(z80 *Z80) {
	var offset byte = z80.memory.readByte(z80.pc)
	z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
	z80.pc++
	z80.h = z80.memory.readByte(z80.IX() + uint16(signExtend(offset)))
}
/* LD z80.IXH(),A */
func instrDD__LD_REGH_A(z80 *Z80) {
	z80.ixh = z80.a
}
/* LD z80.IXL(),B */
func instrDD__LD_REGL_B(z80 *Z80) {
	z80.ixl = z80.b
}
/* LD z80.IXL(),C */
func instrDD__LD_REGL_C(z80 *Z80) {
	z80.ixl = z80.c
}
/* LD z80.IXL(),D */
func instrDD__LD_REGL_D(z80 *Z80) {
	z80.ixl = z80.d
}
/* LD z80.IXL(),E */
func instrDD__LD_REGL_E(z80 *Z80) {
	z80.ixl = z80.e
}
/* LD z80.IXL(),z80.IXH() */
func instrDD__LD_REGL_REGH(z80 *Z80) {
	z80.ixl = z80.ixh
}
/* LD z80.IXL(),z80.IXL() */
func instrDD__LD_REGL_REGL(z80 *Z80) {
}
/* LD L,(ix+dd) */
func instrDD__LD_L_iREGpDD(z80 *Z80) {
	var offset byte = z80.memory.readByte(z80.pc)
	z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
	z80.pc++
	z80.l = z80.memory.readByte(z80.IX() + uint16(signExtend(offset)))
}
/* LD z80.IXL(),A */
func instrDD__LD_REGL_A(z80 *Z80) {
	z80.ixl = z80.a
}
/* LD (ix+dd),B */
func instrDD__LD_iREGpDD_B(z80 *Z80) {
	offset := z80.memory.readByte(z80.pc)
	z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
	z80.pc++
	z80.memory.writeByte(z80.IX()+uint16(signExtend(offset)), z80.b)
}
/* LD (ix+dd),C */
func instrDD__LD_iREGpDD_C(z80 *Z80) {
	offset := z80.memory.readByte(z80.pc)
	z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
	z80.pc++
	z80.memory.writeByte(z80.IX()+uint16(signExtend(offset)), z80.c)
}
/* LD (ix+dd),D */
func instrDD__LD_iREGpDD_D(z80 *Z80) {
	offset := z80.memory.readByte(z80.pc)
	z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
	z80.pc++
	z80.memory.writeByte(z80.IX()+uint16(signExtend(offset)), z80.d)
}
/* LD (ix+dd),E */
func instrDD__LD_iREGpDD_E(z80 *Z80) {
	offset := z80.memory.readByte(z80.pc)
	z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
	z80.pc++
	z80.memory.writeByte(z80.IX()+uint16(signExtend(offset)), z80.e)
}
/* LD (ix+dd),H */
func instrDD__LD_iREGpDD_H(z80 *Z80) {
	offset := z80.memory.readByte(z80.pc)
	z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
	z80.pc++
	z80.memory.writeByte(z80.IX()+uint16(signExtend(offset)), z80.h)
}
/* LD (ix+dd),L */
func instrDD__LD_iREGpDD_L(z80 *Z80) {
	offset := z80.memory.readByte(z80.pc)
	z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
	z80.pc++
	z80.memory.writeByte(z80.IX()+uint16(signExtend(offset)), z80.l)
}
/* LD (ix+dd),A */
func instrDD__LD_iREGpDD_A(z80 *Z80) {
	offset := z80.memory.readByte(z80.pc)
	z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
	z80.pc++
	z80.memory.writeByte(z80.IX()+uint16(signExtend(offset)), z80.a)
}
/* LD A,z80.IXH() */
func instrDD__LD_A_REGH(z80 *Z80) {
	z80.a = z80.ixh
}
/* LD A,z80.IXL() */
func instrDD__LD_A_REGL(z80 *Z80) {
	z80.a = z80.ixl
}
/* LD A,(ix+dd) */
func instrDD__LD_A_iREGpDD(z80 *Z80) {
	var offset byte = z80.memory.readByte(z80.pc)
	z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
	z80.pc++
	z80.a = z80.memory.readByte(z80.IX() + uint16(signExtend(offset)))
}
/* ADD A,z80.IXH() */
func instrDD__ADD_A_REGH(z80 *Z80) {
	z80.add(z80.ixh)
}
/* ADD A,z80.IXL() */
func instrDD__ADD_A_REGL(z80 *Z80) {
	z80.add(z80.ixl)
}
/* ADD A,(ix+dd) */
func instrDD__ADD_A_iREGpDD(z80 *Z80) {
	var offset byte = z80.memory.readByte(z80.pc)
	z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
	z80.pc++
	var bytetemp byte = z80.memory.readByte(z80.IX() + uint16(signExtend(offset)))
	z80.add(bytetemp)
}
/* ADC A,z80.IXH() */
func instrDD__ADC_A_REGH(z80 *Z80) {
	z80.adc(z80.ixh)
}
/* ADC A,z80.IXL() */
func instrDD__ADC_A_REGL(z80 *Z80) {
	z80.adc(z80.ixl)
}
/* ADC A,(ix+dd) */
func instrDD__ADC_A_iREGpDD(z80 *Z80) {
	var offset byte = z80.memory.readByte(z80.pc)
	z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
	z80.pc++
	var bytetemp byte = z80.memory.readByte(z80.IX() + uint16(signExtend(offset)))
	z80.adc(bytetemp)
}
/* SUB A,z80.IXH() */
func instrDD__SUB_A_REGH(z80 *Z80) {
	z80.sub(z80.ixh)
}
/* SUB A,z80.IXL() */
func instrDD__SUB_A_REGL(z80 *Z80) {
	z80.sub(z80.ixl)
}
/* SUB A,(ix+dd) */
func instrDD__SUB_A_iREGpDD(z80 *Z80) {
	var offset byte = z80.memory.readByte(z80.pc)
	z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
	z80.pc++
	var bytetemp byte = z80.memory.readByte(z80.IX() + uint16(signExtend(offset)))
	z80.sub(bytetemp)
}
/* SBC A,z80.IXH() */
func instrDD__SBC_A_REGH(z80 *Z80) {
	z80.sbc(z80.ixh)
}
/* SBC A,z80.IXL() */
func instrDD__SBC_A_REGL(z80 *Z80) {
	z80.sbc(z80.ixl)
}
/* SBC A,(ix+dd) */
func instrDD__SBC_A_iREGpDD(z80 *Z80) {
	var offset byte = z80.memory.readByte(z80.pc)
	z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
	z80.pc++
	var bytetemp byte = z80.memory.readByte(z80.IX() + uint16(signExtend(offset)))
	z80.sbc(bytetemp)
}
/* AND A,z80.IXH() */
func instrDD__AND_A_REGH(z80 *Z80) {
	z80.and(z80.ixh)
}
/* AND A,z80.IXL() */
func instrDD__AND_A_REGL(z80 *Z80) {
	z80.and(z80.ixl)
}
/* AND A,(ix+dd) */
func instrDD__AND_A_iREGpDD(z80 *Z80) {
	var offset byte = z80.memory.readByte(z80.pc)
	z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
	z80.pc++
	var bytetemp byte = z80.memory.readByte(z80.IX() + uint16(signExtend(offset)))
	z80.and(bytetemp)
}
/* XOR A,z80.IXH() */
func instrDD__XOR_A_REGH(z80 *Z80) {
	z80.xor(z80.ixh)
}
/* XOR A,z80.IXL() */
func instrDD__XOR_A_REGL(z80 *Z80) {
	z80.xor(z80.ixl)
}
/* XOR A,(ix+dd) */
func instrDD__XOR_A_iREGpDD(z80 *Z80) {
	var offset byte = z80.memory.readByte(z80.pc)
	z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
	z80.pc++
	var bytetemp byte = z80.memory.readByte(z80.IX() + uint16(signExtend(offset)))
	z80.xor(bytetemp)
}
/* OR A,z80.IXH() */
func instrDD__OR_A_REGH(z80 *Z80) {
	z80.or(z80.ixh)
}
/* OR A,z80.IXL() */
func instrDD__OR_A_REGL(z80 *Z80) {
	z80.or(z80.ixl)
}
/* OR A,(ix+dd) */
func instrDD__OR_A_iREGpDD(z80 *Z80) {
	var offset byte = z80.memory.readByte(z80.pc)
	z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
	z80.pc++
	var bytetemp byte = z80.memory.readByte(z80.IX() + uint16(signExtend(offset)))
	z80.or(bytetemp)
}
/* CP A,z80.IXH() */
func instrDD__CP_A_REGH(z80 *Z80) {
	z80.cp(z80.ixh)
}
/* CP A,z80.IXL() */
func instrDD__CP_A_REGL(z80 *Z80) {
	z80.cp(z80.ixl)
}
/* CP A,(ix+dd) */
func instrDD__CP_A_iREGpDD(z80 *Z80) {
	var offset byte = z80.memory.readByte(z80.pc)
	z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
	z80.pc++
	var bytetemp byte = z80.memory.readByte(z80.IX() + uint16(signExtend(offset)))
	z80.cp(bytetemp)
}
/* shift DDFDCB */
func instrDD__SHIFT_DDFDCB(z80 *Z80) {
}
/* POP ix */
func instrDD__POP_REG(z80 *Z80) {
	z80.ixl, z80.ixh = z80.pop16()
}
/* EX (SP),ix */
func instrDD__EX_iSP_REG(z80 *Z80) {
	var bytetempl = z80.memory.readByte(z80.SP())
	var bytetemph = z80.memory.readByte(z80.SP() + 1)
	z80.memory.contendReadNoMreq(z80.SP()+1, 1)
	z80.memory.writeByte(z80.SP()+1, z80.ixh)
	z80.memory.writeByte(z80.SP(), z80.ixl)
	z80.memory.contendWriteNoMreq_loop(z80.SP(), 1, 2)
	z80.ixl = bytetempl
	z80.ixh = bytetemph
}
/* PUSH ix */
func instrDD__PUSH_REG(z80 *Z80) {
	z80.memory.contendReadNoMreq(z80.IR(), 1)
	z80.push16(z80.ixl, z80.ixh)
}
/* JP ix */
func instrDD__JP_REG(z80 *Z80) {
	z80.pc = z80.IX() /* NB: NOT INDIRECT! */
}
/* LD SP,ix */
func instrDD__LD_SP_REG(z80 *Z80) {
	z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 2)
	z80.sp = z80.IX()
}
/* ADD iy,BC */
func instrFD__ADD_REG_BC(z80 *Z80) {
	z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 7)
	z80.add16(z80.iy, z80.BC())
}
/* ADD iy,DE */
func instrFD__ADD_REG_DE(z80 *Z80) {
	z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 7)
	z80.add16(z80.iy, z80.DE())
}
/* LD iy,nnnn */
func instrFD__LD_REG_NNNN(z80 *Z80) {
	b1 := z80.memory.readByte(z80.pc)
	z80.pc++
	b2 := z80.memory.readByte(z80.pc)
	z80.pc++
	z80.setIY(joinBytes(b2, b1))
}
/* LD (nnnn),iy */
func instrFD__LD_iNNNN_REG(z80 *Z80) {
	z80.ld16nnrr(z80.iyl, z80.iyh)
	// break
}
/* INC iy */
func instrFD__INC_REG(z80 *Z80) {
	z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 2)
	z80.incIY()
}
/* INC z80.IYH() */
func instrFD__INC_REGH(z80 *Z80) {
	z80.incIYH()
}
/* DEC z80.IYH() */
func instrFD__DEC_REGH(z80 *Z80) {
	z80.decIYH()
}
/* LD z80.IYH(),nn */
func instrFD__LD_REGH_NN(z80 *Z80) {
	z80.iyh = z80.memory.readByte(z80.pc)
	z80.pc++
}
/* ADD iy,iy */
func instrFD__ADD_REG_REG(z80 *Z80) {
	z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 7)
	z80.add16(z80.iy, z80.IY())
}
/* LD iy,(nnnn) */
func instrFD__LD_REG_iNNNN(z80 *Z80) {
	z80.ld16rrnn(&z80.iyl, &z80.iyh)
	// break
}
/* DEC iy */
func instrFD__DEC_REG(z80 *Z80) {
	z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 2)
	z80.decIY()
}
/* INC z80.IYL() */
func instrFD__INC_REGL(z80 *Z80) {
	z80.incIYL()
}
/* DEC z80.IYL() */
func instrFD__DEC_REGL(z80 *Z80) {
	z80.decIYL()
}
/* LD z80.IYL(),nn */
func instrFD__LD_REGL_NN(z80 *Z80) {
	z80.iyl = z80.memory.readByte(z80.pc)
	z80.pc++
}
/* INC (iy+dd) */
func instrFD__INC_iREGpDD(z80 *Z80) {
	var offset byte = z80.memory.readByte(z80.pc)
	z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
	z80.pc++
	var wordtemp uint16 = z80.IY() + uint16(signExtend(offset))
	var bytetemp byte = z80.memory.readByte(wordtemp)
	z80.memory.contendReadNoMreq(wordtemp, 1)
	z80.inc(&bytetemp)
	z80.memory.writeByte(wordtemp, bytetemp)
}
/* DEC (iy+dd) */
func instrFD__DEC_iREGpDD(z80 *Z80) {
	var offset byte = z80.memory.readByte(z80.pc)
	z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
	z80.pc++
	var wordtemp uint16 = z80.IY() + uint16(signExtend(offset))
	var bytetemp byte = z80.memory.readByte(wordtemp)
	z80.memory.contendReadNoMreq(wordtemp, 1)
	z80.dec(&bytetemp)
	z80.memory.writeByte(wordtemp, bytetemp)
}
/* LD (iy+dd),nn */
func instrFD__LD_iREGpDD_NN(z80 *Z80) {
	offset := z80.memory.readByte(z80.pc)
	z80.pc++
	value := z80.memory.readByte(z80.pc)
	z80.memory.contendReadNoMreq_loop(z80.pc, 1, 2)
	z80.pc++
	z80.memory.writeByte(z80.IY()+uint16(signExtend(offset)), value)
}
/* ADD iy,SP */
func instrFD__ADD_REG_SP(z80 *Z80) {
	z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 7)
	z80.add16(z80.iy, z80.SP())
}
/* LD B,z80.IYH() */
func instrFD__LD_B_REGH(z80 *Z80) {
	z80.b = z80.iyh
}
/* LD B,z80.IYL() */
func instrFD__LD_B_REGL(z80 *Z80) {
	z80.b = z80.iyl
}
/* LD B,(iy+dd) */
func instrFD__LD_B_iREGpDD(z80 *Z80) {
	var offset byte = z80.memory.readByte(z80.pc)
	z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
	z80.pc++
	z80.b = z80.memory.readByte(z80.IY() + uint16(signExtend(offset)))
}
/* LD C,z80.IYH() */
func instrFD__LD_C_REGH(z80 *Z80) {
	z80.c = z80.iyh
}
/* LD C,z80.IYL() */
func instrFD__LD_C_REGL(z80 *Z80) {
	z80.c = z80.iyl
}
/* LD C,(iy+dd) */
func instrFD__LD_C_iREGpDD(z80 *Z80) {
	var offset byte = z80.memory.readByte(z80.pc)
	z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
	z80.pc++
	z80.c = z80.memory.readByte(z80.IY() + uint16(signExtend(offset)))
}
/* LD D,z80.IYH() */
func instrFD__LD_D_REGH(z80 *Z80) {
	z80.d = z80.iyh
}
/* LD D,z80.IYL() */
func instrFD__LD_D_REGL(z80 *Z80) {
	z80.d = z80.iyl
}
/* LD D,(iy+dd) */
func instrFD__LD_D_iREGpDD(z80 *Z80) {
	var offset byte = z80.memory.readByte(z80.pc)
	z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
	z80.pc++
	z80.d = z80.memory.readByte(z80.IY() + uint16(signExtend(offset)))
}
/* LD E,z80.IYH() */
func instrFD__LD_E_REGH(z80 *Z80) {
	z80.e = z80.iyh
}
/* LD E,z80.IYL() */
func instrFD__LD_E_REGL(z80 *Z80) {
	z80.e = z80.iyl
}
/* LD E,(iy+dd) */
func instrFD__LD_E_iREGpDD(z80 *Z80) {
	var offset byte = z80.memory.readByte(z80.pc)
	z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
	z80.pc++
	z80.e = z80.memory.readByte(z80.IY() + uint16(signExtend(offset)))
}
/* LD z80.IYH(),B */
func instrFD__LD_REGH_B(z80 *Z80) {
	z80.iyh = z80.b
}
/* LD z80.IYH(),C */
func instrFD__LD_REGH_C(z80 *Z80) {
	z80.iyh = z80.c
}
/* LD z80.IYH(),D */
func instrFD__LD_REGH_D(z80 *Z80) {
	z80.iyh = z80.d
}
/* LD z80.IYH(),E */
func instrFD__LD_REGH_E(z80 *Z80) {
	z80.iyh = z80.e
}
/* LD z80.IYH(),z80.IYH() */
func instrFD__LD_REGH_REGH(z80 *Z80) {
}
/* LD z80.IYH(),z80.IYL() */
func instrFD__LD_REGH_REGL(z80 *Z80) {
	z80.iyh = z80.iyl
}
/* LD H,(iy+dd) */
func instrFD__LD_H_iREGpDD(z80 *Z80) {
	var offset byte = z80.memory.readByte(z80.pc)
	z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
	z80.pc++
	z80.h = z80.memory.readByte(z80.IY() + uint16(signExtend(offset)))
}
/* LD z80.IYH(),A */
func instrFD__LD_REGH_A(z80 *Z80) {
	z80.iyh = z80.a
}
/* LD z80.IYL(),B */
func instrFD__LD_REGL_B(z80 *Z80) {
	z80.iyl = z80.b
}
/* LD z80.IYL(),C */
func instrFD__LD_REGL_C(z80 *Z80) {
	z80.iyl = z80.c
}
/* LD z80.IYL(),D */
func instrFD__LD_REGL_D(z80 *Z80) {
	z80.iyl = z80.d
}
/* LD z80.IYL(),E */
func instrFD__LD_REGL_E(z80 *Z80) {
	z80.iyl = z80.e
}
/* LD z80.IYL(),z80.IYH() */
func instrFD__LD_REGL_REGH(z80 *Z80) {
	z80.iyl = z80.iyh
}
/* LD z80.IYL(),z80.IYL() */
func instrFD__LD_REGL_REGL(z80 *Z80) {
}
/* LD L,(iy+dd) */
func instrFD__LD_L_iREGpDD(z80 *Z80) {
	var offset byte = z80.memory.readByte(z80.pc)
	z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
	z80.pc++
	z80.l = z80.memory.readByte(z80.IY() + uint16(signExtend(offset)))
}
/* LD z80.IYL(),A */
func instrFD__LD_REGL_A(z80 *Z80) {
	z80.iyl = z80.a
}
/* LD (iy+dd),B */
func instrFD__LD_iREGpDD_B(z80 *Z80) {
	offset := z80.memory.readByte(z80.pc)
	z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
	z80.pc++
	z80.memory.writeByte(z80.IY()+uint16(signExtend(offset)), z80.b)
}
/* LD (iy+dd),C */
func instrFD__LD_iREGpDD_C(z80 *Z80) {
	offset := z80.memory.readByte(z80.pc)
	z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
	z80.pc++
	z80.memory.writeByte(z80.IY()+uint16(signExtend(offset)), z80.c)
}
/* LD (iy+dd),D */
func instrFD__LD_iREGpDD_D(z80 *Z80) {
	offset := z80.memory.readByte(z80.pc)
	z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
	z80.pc++
	z80.memory.writeByte(z80.IY()+uint16(signExtend(offset)), z80.d)
}
/* LD (iy+dd),E */
func instrFD__LD_iREGpDD_E(z80 *Z80) {
	offset := z80.memory.readByte(z80.pc)
	z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
	z80.pc++
	z80.memory.writeByte(z80.IY()+uint16(signExtend(offset)), z80.e)
}
/* LD (iy+dd),H */
func instrFD__LD_iREGpDD_H(z80 *Z80) {
	offset := z80.memory.readByte(z80.pc)
	z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
	z80.pc++
	z80.memory.writeByte(z80.IY()+uint16(signExtend(offset)), z80.h)
}
/* LD (iy+dd),L */
func instrFD__LD_iREGpDD_L(z80 *Z80) {
	offset := z80.memory.readByte(z80.pc)
	z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
	z80.pc++
	z80.memory.writeByte(z80.IY()+uint16(signExtend(offset)), z80.l)
}
/* LD (iy+dd),A */
func instrFD__LD_iREGpDD_A(z80 *Z80) {
	offset := z80.memory.readByte(z80.pc)
	z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
	z80.pc++
	z80.memory.writeByte(z80.IY()+uint16(signExtend(offset)), z80.a)
}
/* LD A,z80.IYH() */
func instrFD__LD_A_REGH(z80 *Z80) {
	z80.a = z80.iyh
}
/* LD A,z80.IYL() */
func instrFD__LD_A_REGL(z80 *Z80) {
	z80.a = z80.iyl
}
/* LD A,(iy+dd) */
func instrFD__LD_A_iREGpDD(z80 *Z80) {
	var offset byte = z80.memory.readByte(z80.pc)
	z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
	z80.pc++
	z80.a = z80.memory.readByte(z80.IY() + uint16(signExtend(offset)))
}
/* ADD A,z80.IYH() */
func instrFD__ADD_A_REGH(z80 *Z80) {
	z80.add(z80.iyh)
}
/* ADD A,z80.IYL() */
func instrFD__ADD_A_REGL(z80 *Z80) {
	z80.add(z80.iyl)
}
/* ADD A,(iy+dd) */
func instrFD__ADD_A_iREGpDD(z80 *Z80) {
	var offset byte = z80.memory.readByte(z80.pc)
	z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
	z80.pc++
	var bytetemp byte = z80.memory.readByte(z80.IY() + uint16(signExtend(offset)))
	z80.add(bytetemp)
}
/* ADC A,z80.IYH() */
func instrFD__ADC_A_REGH(z80 *Z80) {
	z80.adc(z80.iyh)
}
/* ADC A,z80.IYL() */
func instrFD__ADC_A_REGL(z80 *Z80) {
	z80.adc(z80.iyl)
}
/* ADC A,(iy+dd) */
func instrFD__ADC_A_iREGpDD(z80 *Z80) {
	var offset byte = z80.memory.readByte(z80.pc)
	z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
	z80.pc++
	var bytetemp byte = z80.memory.readByte(z80.IY() + uint16(signExtend(offset)))
	z80.adc(bytetemp)
}
/* SUB A,z80.IYH() */
func instrFD__SUB_A_REGH(z80 *Z80) {
	z80.sub(z80.iyh)
}
/* SUB A,z80.IYL() */
func instrFD__SUB_A_REGL(z80 *Z80) {
	z80.sub(z80.iyl)
}
/* SUB A,(iy+dd) */
func instrFD__SUB_A_iREGpDD(z80 *Z80) {
	var offset byte = z80.memory.readByte(z80.pc)
	z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
	z80.pc++
	var bytetemp byte = z80.memory.readByte(z80.IY() + uint16(signExtend(offset)))
	z80.sub(bytetemp)
}
/* SBC A,z80.IYH() */
func instrFD__SBC_A_REGH(z80 *Z80) {
	z80.sbc(z80.iyh)
}
/* SBC A,z80.IYL() */
func instrFD__SBC_A_REGL(z80 *Z80) {
	z80.sbc(z80.iyl)
}
/* SBC A,(iy+dd) */
func instrFD__SBC_A_iREGpDD(z80 *Z80) {
	var offset byte = z80.memory.readByte(z80.pc)
	z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
	z80.pc++
	var bytetemp byte = z80.memory.readByte(z80.IY() + uint16(signExtend(offset)))
	z80.sbc(bytetemp)
}
/* AND A,z80.IYH() */
func instrFD__AND_A_REGH(z80 *Z80) {
	z80.and(z80.iyh)
}
/* AND A,z80.IYL() */
func instrFD__AND_A_REGL(z80 *Z80) {
	z80.and(z80.iyl)
}
/* AND A,(iy+dd) */
func instrFD__AND_A_iREGpDD(z80 *Z80) {
	var offset byte = z80.memory.readByte(z80.pc)
	z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
	z80.pc++
	var bytetemp byte = z80.memory.readByte(z80.IY() + uint16(signExtend(offset)))
	z80.and(bytetemp)
}
/* XOR A,z80.IYH() */
func instrFD__XOR_A_REGH(z80 *Z80) {
	z80.xor(z80.iyh)
}
/* XOR A,z80.IYL() */
func instrFD__XOR_A_REGL(z80 *Z80) {
	z80.xor(z80.iyl)
}
/* XOR A,(iy+dd) */
func instrFD__XOR_A_iREGpDD(z80 *Z80) {
	var offset byte = z80.memory.readByte(z80.pc)
	z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
	z80.pc++
	var bytetemp byte = z80.memory.readByte(z80.IY() + uint16(signExtend(offset)))
	z80.xor(bytetemp)
}
/* OR A,z80.IYH() */
func instrFD__OR_A_REGH(z80 *Z80) {
	z80.or(z80.iyh)
}
/* OR A,z80.IYL() */
func instrFD__OR_A_REGL(z80 *Z80) {
	z80.or(z80.iyl)
}
/* OR A,(iy+dd) */
func instrFD__OR_A_iREGpDD(z80 *Z80) {
	var offset byte = z80.memory.readByte(z80.pc)
	z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
	z80.pc++
	var bytetemp byte = z80.memory.readByte(z80.IY() + uint16(signExtend(offset)))
	z80.or(bytetemp)
}
/* CP A,z80.IYH() */
func instrFD__CP_A_REGH(z80 *Z80) {
	z80.cp(z80.iyh)
}
/* CP A,z80.IYL() */
func instrFD__CP_A_REGL(z80 *Z80) {
	z80.cp(z80.iyl)
}
/* CP A,(iy+dd) */
func instrFD__CP_A_iREGpDD(z80 *Z80) {
	var offset byte = z80.memory.readByte(z80.pc)
	z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
	z80.pc++
	var bytetemp byte = z80.memory.readByte(z80.IY() + uint16(signExtend(offset)))
	z80.cp(bytetemp)
}
/* shift DDFDCB */
func instrFD__SHIFT_DDFDCB(z80 *Z80) {
}
/* POP iy */
func instrFD__POP_REG(z80 *Z80) {
	z80.iyl, z80.iyh = z80.pop16()
}
/* EX (SP),iy */
func instrFD__EX_iSP_REG(z80 *Z80) {
	var bytetempl = z80.memory.readByte(z80.SP())
	var bytetemph = z80.memory.readByte(z80.SP() + 1)
	z80.memory.contendReadNoMreq(z80.SP()+1, 1)
	z80.memory.writeByte(z80.SP()+1, z80.iyh)
	z80.memory.writeByte(z80.SP(), z80.iyl)
	z80.memory.contendWriteNoMreq_loop(z80.SP(), 1, 2)
	z80.iyl = bytetempl
	z80.iyh = bytetemph
}
/* PUSH iy */
func instrFD__PUSH_REG(z80 *Z80) {
	z80.memory.contendReadNoMreq(z80.IR(), 1)
	z80.push16(z80.iyl, z80.iyh)
}
/* JP iy */
func instrFD__JP_REG(z80 *Z80) {
	z80.pc = z80.IY() /* NB: NOT INDIRECT! */
}
/* LD SP,iy */
func instrFD__LD_SP_REG(z80 *Z80) {
	z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 2)
	z80.sp = z80.IY()
}
/* LD B,RLC (REGISTER+dd) */
func instrDDCB__LD_B_RLC_iREGpDD(z80 *Z80) {
	z80.b = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.b = z80.rlc(z80.b)
	z80.memory.writeByte(z80.tempaddr, z80.b)
}
/* LD C,RLC (REGISTER+dd) */
func instrDDCB__LD_C_RLC_iREGpDD(z80 *Z80) {
	z80.c = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.c = z80.rlc(z80.c)
	z80.memory.writeByte(z80.tempaddr, z80.c)
}
/* LD D,RLC (REGISTER+dd) */
func instrDDCB__LD_D_RLC_iREGpDD(z80 *Z80) {
	z80.d = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.d = z80.rlc(z80.d)
	z80.memory.writeByte(z80.tempaddr, z80.d)
}
/* LD E,RLC (REGISTER+dd) */
func instrDDCB__LD_E_RLC_iREGpDD(z80 *Z80) {
	z80.e = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.e = z80.rlc(z80.e)
	z80.memory.writeByte(z80.tempaddr, z80.e)
}
/* LD H,RLC (REGISTER+dd) */
func instrDDCB__LD_H_RLC_iREGpDD(z80 *Z80) {
	z80.h = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.h = z80.rlc(z80.h)
	z80.memory.writeByte(z80.tempaddr, z80.h)
}
/* LD L,RLC (REGISTER+dd) */
func instrDDCB__LD_L_RLC_iREGpDD(z80 *Z80) {
	z80.l = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.l = z80.rlc(z80.l)
	z80.memory.writeByte(z80.tempaddr, z80.l)
}
/* RLC (REGISTER+dd) */
func instrDDCB__RLC_iREGpDD(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	bytetemp = z80.rlc(bytetemp)
	z80.memory.writeByte(z80.tempaddr, bytetemp)
}
/* LD A,RLC (REGISTER+dd) */
func instrDDCB__LD_A_RLC_iREGpDD(z80 *Z80) {
	z80.a = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.a = z80.rlc(z80.a)
	z80.memory.writeByte(z80.tempaddr, z80.a)
}
/* LD B,RRC (REGISTER+dd) */
func instrDDCB__LD_B_RRC_iREGpDD(z80 *Z80) {
	z80.b = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.b = z80.rrc(z80.b)
	z80.memory.writeByte(z80.tempaddr, z80.b)
}
/* LD C,RRC (REGISTER+dd) */
func instrDDCB__LD_C_RRC_iREGpDD(z80 *Z80) {
	z80.c = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.c = z80.rrc(z80.c)
	z80.memory.writeByte(z80.tempaddr, z80.c)
}
/* LD D,RRC (REGISTER+dd) */
func instrDDCB__LD_D_RRC_iREGpDD(z80 *Z80) {
	z80.d = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.d = z80.rrc(z80.d)
	z80.memory.writeByte(z80.tempaddr, z80.d)
}
/* LD E,RRC (REGISTER+dd) */
func instrDDCB__LD_E_RRC_iREGpDD(z80 *Z80) {
	z80.e = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.e = z80.rrc(z80.e)
	z80.memory.writeByte(z80.tempaddr, z80.e)
}
/* LD H,RRC (REGISTER+dd) */
func instrDDCB__LD_H_RRC_iREGpDD(z80 *Z80) {
	z80.h = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.h = z80.rrc(z80.h)
	z80.memory.writeByte(z80.tempaddr, z80.h)
}
/* LD L,RRC (REGISTER+dd) */
func instrDDCB__LD_L_RRC_iREGpDD(z80 *Z80) {
	z80.l = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.l = z80.rrc(z80.l)
	z80.memory.writeByte(z80.tempaddr, z80.l)
}
/* RRC (REGISTER+dd) */
func instrDDCB__RRC_iREGpDD(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	bytetemp = z80.rrc(bytetemp)
	z80.memory.writeByte(z80.tempaddr, bytetemp)
}
/* LD A,RRC (REGISTER+dd) */
func instrDDCB__LD_A_RRC_iREGpDD(z80 *Z80) {
	z80.a = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.a = z80.rrc(z80.a)
	z80.memory.writeByte(z80.tempaddr, z80.a)
}
/* LD B,RL (REGISTER+dd) */
func instrDDCB__LD_B_RL_iREGpDD(z80 *Z80) {
	z80.b = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.b = z80.rl(z80.b)
	z80.memory.writeByte(z80.tempaddr, z80.b)
}
/* LD C,RL (REGISTER+dd) */
func instrDDCB__LD_C_RL_iREGpDD(z80 *Z80) {
	z80.c = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.c = z80.rl(z80.c)
	z80.memory.writeByte(z80.tempaddr, z80.c)
}
/* LD D,RL (REGISTER+dd) */
func instrDDCB__LD_D_RL_iREGpDD(z80 *Z80) {
	z80.d = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.d = z80.rl(z80.d)
	z80.memory.writeByte(z80.tempaddr, z80.d)
}
/* LD E,RL (REGISTER+dd) */
func instrDDCB__LD_E_RL_iREGpDD(z80 *Z80) {
	z80.e = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.e = z80.rl(z80.e)
	z80.memory.writeByte(z80.tempaddr, z80.e)
}
/* LD H,RL (REGISTER+dd) */
func instrDDCB__LD_H_RL_iREGpDD(z80 *Z80) {
	z80.h = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.h = z80.rl(z80.h)
	z80.memory.writeByte(z80.tempaddr, z80.h)
}
/* LD L,RL (REGISTER+dd) */
func instrDDCB__LD_L_RL_iREGpDD(z80 *Z80) {
	z80.l = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.l = z80.rl(z80.l)
	z80.memory.writeByte(z80.tempaddr, z80.l)
}
/* RL (REGISTER+dd) */
func instrDDCB__RL_iREGpDD(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	bytetemp = z80.rl(bytetemp)
	z80.memory.writeByte(z80.tempaddr, bytetemp)
}
/* LD A,RL (REGISTER+dd) */
func instrDDCB__LD_A_RL_iREGpDD(z80 *Z80) {
	z80.a = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.a = z80.rl(z80.a)
	z80.memory.writeByte(z80.tempaddr, z80.a)
}
/* LD B,RR (REGISTER+dd) */
func instrDDCB__LD_B_RR_iREGpDD(z80 *Z80) {
	z80.b = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.b = z80.rr(z80.b)
	z80.memory.writeByte(z80.tempaddr, z80.b)
}
/* LD C,RR (REGISTER+dd) */
func instrDDCB__LD_C_RR_iREGpDD(z80 *Z80) {
	z80.c = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.c = z80.rr(z80.c)
	z80.memory.writeByte(z80.tempaddr, z80.c)
}
/* LD D,RR (REGISTER+dd) */
func instrDDCB__LD_D_RR_iREGpDD(z80 *Z80) {
	z80.d = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.d = z80.rr(z80.d)
	z80.memory.writeByte(z80.tempaddr, z80.d)
}
/* LD E,RR (REGISTER+dd) */
func instrDDCB__LD_E_RR_iREGpDD(z80 *Z80) {
	z80.e = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.e = z80.rr(z80.e)
	z80.memory.writeByte(z80.tempaddr, z80.e)
}
/* LD H,RR (REGISTER+dd) */
func instrDDCB__LD_H_RR_iREGpDD(z80 *Z80) {
	z80.h = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.h = z80.rr(z80.h)
	z80.memory.writeByte(z80.tempaddr, z80.h)
}
/* LD L,RR (REGISTER+dd) */
func instrDDCB__LD_L_RR_iREGpDD(z80 *Z80) {
	z80.l = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.l = z80.rr(z80.l)
	z80.memory.writeByte(z80.tempaddr, z80.l)
}
/* RR (REGISTER+dd) */
func instrDDCB__RR_iREGpDD(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	bytetemp = z80.rr(bytetemp)
	z80.memory.writeByte(z80.tempaddr, bytetemp)
}
/* LD A,RR (REGISTER+dd) */
func instrDDCB__LD_A_RR_iREGpDD(z80 *Z80) {
	z80.a = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.a = z80.rr(z80.a)
	z80.memory.writeByte(z80.tempaddr, z80.a)
}
/* LD B,SLA (REGISTER+dd) */
func instrDDCB__LD_B_SLA_iREGpDD(z80 *Z80) {
	z80.b = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.b = z80.sla(z80.b)
	z80.memory.writeByte(z80.tempaddr, z80.b)
}
/* LD C,SLA (REGISTER+dd) */
func instrDDCB__LD_C_SLA_iREGpDD(z80 *Z80) {
	z80.c = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.c = z80.sla(z80.c)
	z80.memory.writeByte(z80.tempaddr, z80.c)
}
/* LD D,SLA (REGISTER+dd) */
func instrDDCB__LD_D_SLA_iREGpDD(z80 *Z80) {
	z80.d = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.d = z80.sla(z80.d)
	z80.memory.writeByte(z80.tempaddr, z80.d)
}
/* LD E,SLA (REGISTER+dd) */
func instrDDCB__LD_E_SLA_iREGpDD(z80 *Z80) {
	z80.e = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.e = z80.sla(z80.e)
	z80.memory.writeByte(z80.tempaddr, z80.e)
}
/* LD H,SLA (REGISTER+dd) */
func instrDDCB__LD_H_SLA_iREGpDD(z80 *Z80) {
	z80.h = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.h = z80.sla(z80.h)
	z80.memory.writeByte(z80.tempaddr, z80.h)
}
/* LD L,SLA (REGISTER+dd) */
func instrDDCB__LD_L_SLA_iREGpDD(z80 *Z80) {
	z80.l = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.l = z80.sla(z80.l)
	z80.memory.writeByte(z80.tempaddr, z80.l)
}
/* SLA (REGISTER+dd) */
func instrDDCB__SLA_iREGpDD(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	bytetemp = z80.sla(bytetemp)
	z80.memory.writeByte(z80.tempaddr, bytetemp)
}
/* LD A,SLA (REGISTER+dd) */
func instrDDCB__LD_A_SLA_iREGpDD(z80 *Z80) {
	z80.a = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.a = z80.sla(z80.a)
	z80.memory.writeByte(z80.tempaddr, z80.a)
}
/* LD B,SRA (REGISTER+dd) */
func instrDDCB__LD_B_SRA_iREGpDD(z80 *Z80) {
	z80.b = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.b = z80.sra(z80.b)
	z80.memory.writeByte(z80.tempaddr, z80.b)
}
/* LD C,SRA (REGISTER+dd) */
func instrDDCB__LD_C_SRA_iREGpDD(z80 *Z80) {
	z80.c = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.c = z80.sra(z80.c)
	z80.memory.writeByte(z80.tempaddr, z80.c)
}
/* LD D,SRA (REGISTER+dd) */
func instrDDCB__LD_D_SRA_iREGpDD(z80 *Z80) {
	z80.d = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.d = z80.sra(z80.d)
	z80.memory.writeByte(z80.tempaddr, z80.d)
}
/* LD E,SRA (REGISTER+dd) */
func instrDDCB__LD_E_SRA_iREGpDD(z80 *Z80) {
	z80.e = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.e = z80.sra(z80.e)
	z80.memory.writeByte(z80.tempaddr, z80.e)
}
/* LD H,SRA (REGISTER+dd) */
func instrDDCB__LD_H_SRA_iREGpDD(z80 *Z80) {
	z80.h = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.h = z80.sra(z80.h)
	z80.memory.writeByte(z80.tempaddr, z80.h)
}
/* LD L,SRA (REGISTER+dd) */
func instrDDCB__LD_L_SRA_iREGpDD(z80 *Z80) {
	z80.l = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.l = z80.sra(z80.l)
	z80.memory.writeByte(z80.tempaddr, z80.l)
}
/* SRA (REGISTER+dd) */
func instrDDCB__SRA_iREGpDD(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	bytetemp = z80.sra(bytetemp)
	z80.memory.writeByte(z80.tempaddr, bytetemp)
}
/* LD A,SRA (REGISTER+dd) */
func instrDDCB__LD_A_SRA_iREGpDD(z80 *Z80) {
	z80.a = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.a = z80.sra(z80.a)
	z80.memory.writeByte(z80.tempaddr, z80.a)
}
/* LD B,SLL (REGISTER+dd) */
func instrDDCB__LD_B_SLL_iREGpDD(z80 *Z80) {
	z80.b = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.b = z80.sll(z80.b)
	z80.memory.writeByte(z80.tempaddr, z80.b)
}
/* LD C,SLL (REGISTER+dd) */
func instrDDCB__LD_C_SLL_iREGpDD(z80 *Z80) {
	z80.c = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.c = z80.sll(z80.c)
	z80.memory.writeByte(z80.tempaddr, z80.c)
}
/* LD D,SLL (REGISTER+dd) */
func instrDDCB__LD_D_SLL_iREGpDD(z80 *Z80) {
	z80.d = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.d = z80.sll(z80.d)
	z80.memory.writeByte(z80.tempaddr, z80.d)
}
/* LD E,SLL (REGISTER+dd) */
func instrDDCB__LD_E_SLL_iREGpDD(z80 *Z80) {
	z80.e = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.e = z80.sll(z80.e)
	z80.memory.writeByte(z80.tempaddr, z80.e)
}
/* LD H,SLL (REGISTER+dd) */
func instrDDCB__LD_H_SLL_iREGpDD(z80 *Z80) {
	z80.h = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.h = z80.sll(z80.h)
	z80.memory.writeByte(z80.tempaddr, z80.h)
}
/* LD L,SLL (REGISTER+dd) */
func instrDDCB__LD_L_SLL_iREGpDD(z80 *Z80) {
	z80.l = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.l = z80.sll(z80.l)
	z80.memory.writeByte(z80.tempaddr, z80.l)
}
/* SLL (REGISTER+dd) */
func instrDDCB__SLL_iREGpDD(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	bytetemp = z80.sll(bytetemp)
	z80.memory.writeByte(z80.tempaddr, bytetemp)
}
/* LD A,SLL (REGISTER+dd) */
func instrDDCB__LD_A_SLL_iREGpDD(z80 *Z80) {
	z80.a = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.a = z80.sll(z80.a)
	z80.memory.writeByte(z80.tempaddr, z80.a)
}
/* LD B,SRL (REGISTER+dd) */
func instrDDCB__LD_B_SRL_iREGpDD(z80 *Z80) {
	z80.b = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.b = z80.srl(z80.b)
	z80.memory.writeByte(z80.tempaddr, z80.b)
}
/* LD C,SRL (REGISTER+dd) */
func instrDDCB__LD_C_SRL_iREGpDD(z80 *Z80) {
	z80.c = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.c = z80.srl(z80.c)
	z80.memory.writeByte(z80.tempaddr, z80.c)
}
/* LD D,SRL (REGISTER+dd) */
func instrDDCB__LD_D_SRL_iREGpDD(z80 *Z80) {
	z80.d = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.d = z80.srl(z80.d)
	z80.memory.writeByte(z80.tempaddr, z80.d)
}
/* LD E,SRL (REGISTER+dd) */
func instrDDCB__LD_E_SRL_iREGpDD(z80 *Z80) {
	z80.e = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.e = z80.srl(z80.e)
	z80.memory.writeByte(z80.tempaddr, z80.e)
}
/* LD H,SRL (REGISTER+dd) */
func instrDDCB__LD_H_SRL_iREGpDD(z80 *Z80) {
	z80.h = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.h = z80.srl(z80.h)
	z80.memory.writeByte(z80.tempaddr, z80.h)
}
/* LD L,SRL (REGISTER+dd) */
func instrDDCB__LD_L_SRL_iREGpDD(z80 *Z80) {
	z80.l = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.l = z80.srl(z80.l)
	z80.memory.writeByte(z80.tempaddr, z80.l)
}
/* SRL (REGISTER+dd) */
func instrDDCB__SRL_iREGpDD(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	bytetemp = z80.srl(bytetemp)
	z80.memory.writeByte(z80.tempaddr, bytetemp)
}
/* LD A,SRL (REGISTER+dd) */
func instrDDCB__LD_A_SRL_iREGpDD(z80 *Z80) {
	z80.a = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.a = z80.srl(z80.a)
	z80.memory.writeByte(z80.tempaddr, z80.a)
}
/* BIT 0,(REGISTER+dd) */
func instrDDCB__BIT_0_iREGpDD(z80 *Z80) {
	bytetemp := z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.biti(0, bytetemp, z80.tempaddr)
}
/* BIT 1,(REGISTER+dd) */
func instrDDCB__BIT_1_iREGpDD(z80 *Z80) {
	bytetemp := z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.biti(1, bytetemp, z80.tempaddr)
}
/* BIT 2,(REGISTER+dd) */
func instrDDCB__BIT_2_iREGpDD(z80 *Z80) {
	bytetemp := z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.biti(2, bytetemp, z80.tempaddr)
}
/* BIT 3,(REGISTER+dd) */
func instrDDCB__BIT_3_iREGpDD(z80 *Z80) {
	bytetemp := z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.biti(3, bytetemp, z80.tempaddr)
}
/* BIT 4,(REGISTER+dd) */
func instrDDCB__BIT_4_iREGpDD(z80 *Z80) {
	bytetemp := z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.biti(4, bytetemp, z80.tempaddr)
}
/* BIT 5,(REGISTER+dd) */
func instrDDCB__BIT_5_iREGpDD(z80 *Z80) {
	bytetemp := z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.biti(5, bytetemp, z80.tempaddr)
}
/* BIT 6,(REGISTER+dd) */
func instrDDCB__BIT_6_iREGpDD(z80 *Z80) {
	bytetemp := z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.biti(6, bytetemp, z80.tempaddr)
}
/* BIT 7,(REGISTER+dd) */
func instrDDCB__BIT_7_iREGpDD(z80 *Z80) {
	bytetemp := z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.biti(7, bytetemp, z80.tempaddr)
}
/* LD B,RES 0,(REGISTER+dd) */
func instrDDCB__LD_B_RES_0_iREGpDD(z80 *Z80) {
	z80.b = z80.memory.readByte(z80.tempaddr) & 0xfe
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.b)
}
/* LD C,RES 0,(REGISTER+dd) */
func instrDDCB__LD_C_RES_0_iREGpDD(z80 *Z80) {
	z80.c = z80.memory.readByte(z80.tempaddr) & 0xfe
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.c)
}
/* LD D,RES 0,(REGISTER+dd) */
func instrDDCB__LD_D_RES_0_iREGpDD(z80 *Z80) {
	z80.d = z80.memory.readByte(z80.tempaddr) & 0xfe
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.d)
}
/* LD E,RES 0,(REGISTER+dd) */
func instrDDCB__LD_E_RES_0_iREGpDD(z80 *Z80) {
	z80.e = z80.memory.readByte(z80.tempaddr) & 0xfe
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.e)
}
/* LD H,RES 0,(REGISTER+dd) */
func instrDDCB__LD_H_RES_0_iREGpDD(z80 *Z80) {
	z80.h = z80.memory.readByte(z80.tempaddr) & 0xfe
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.h)
}
/* LD L,RES 0,(REGISTER+dd) */
func instrDDCB__LD_L_RES_0_iREGpDD(z80 *Z80) {
	z80.l = z80.memory.readByte(z80.tempaddr) & 0xfe
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.l)
}
/* RES 0,(REGISTER+dd) */
func instrDDCB__RES_0_iREGpDD(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, bytetemp&0xfe)
}
/* LD A,RES 0,(REGISTER+dd) */
func instrDDCB__LD_A_RES_0_iREGpDD(z80 *Z80) {
	z80.a = z80.memory.readByte(z80.tempaddr) & 0xfe
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.a)
}
/* LD B,RES 1,(REGISTER+dd) */
func instrDDCB__LD_B_RES_1_iREGpDD(z80 *Z80) {
	z80.b = z80.memory.readByte(z80.tempaddr) & 0xfd
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.b)
}
/* LD C,RES 1,(REGISTER+dd) */
func instrDDCB__LD_C_RES_1_iREGpDD(z80 *Z80) {
	z80.c = z80.memory.readByte(z80.tempaddr) & 0xfd
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.c)
}
/* LD D,RES 1,(REGISTER+dd) */
func instrDDCB__LD_D_RES_1_iREGpDD(z80 *Z80) {
	z80.d = z80.memory.readByte(z80.tempaddr) & 0xfd
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.d)
}
/* LD E,RES 1,(REGISTER+dd) */
func instrDDCB__LD_E_RES_1_iREGpDD(z80 *Z80) {
	z80.e = z80.memory.readByte(z80.tempaddr) & 0xfd
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.e)
}
/* LD H,RES 1,(REGISTER+dd) */
func instrDDCB__LD_H_RES_1_iREGpDD(z80 *Z80) {
	z80.h = z80.memory.readByte(z80.tempaddr) & 0xfd
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.h)
}
/* LD L,RES 1,(REGISTER+dd) */
func instrDDCB__LD_L_RES_1_iREGpDD(z80 *Z80) {
	z80.l = z80.memory.readByte(z80.tempaddr) & 0xfd
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.l)
}
/* RES 1,(REGISTER+dd) */
func instrDDCB__RES_1_iREGpDD(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, bytetemp&0xfd)
}
/* LD A,RES 1,(REGISTER+dd) */
func instrDDCB__LD_A_RES_1_iREGpDD(z80 *Z80) {
	z80.a = z80.memory.readByte(z80.tempaddr) & 0xfd
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.a)
}
/* LD B,RES 2,(REGISTER+dd) */
func instrDDCB__LD_B_RES_2_iREGpDD(z80 *Z80) {
	z80.b = z80.memory.readByte(z80.tempaddr) & 0xfb
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.b)
}
/* LD C,RES 2,(REGISTER+dd) */
func instrDDCB__LD_C_RES_2_iREGpDD(z80 *Z80) {
	z80.c = z80.memory.readByte(z80.tempaddr) & 0xfb
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.c)
}
/* LD D,RES 2,(REGISTER+dd) */
func instrDDCB__LD_D_RES_2_iREGpDD(z80 *Z80) {
	z80.d = z80.memory.readByte(z80.tempaddr) & 0xfb
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.d)
}
/* LD E,RES 2,(REGISTER+dd) */
func instrDDCB__LD_E_RES_2_iREGpDD(z80 *Z80) {
	z80.e = z80.memory.readByte(z80.tempaddr) & 0xfb
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.e)
}
/* LD H,RES 2,(REGISTER+dd) */
func instrDDCB__LD_H_RES_2_iREGpDD(z80 *Z80) {
	z80.h = z80.memory.readByte(z80.tempaddr) & 0xfb
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.h)
}
/* LD L,RES 2,(REGISTER+dd) */
func instrDDCB__LD_L_RES_2_iREGpDD(z80 *Z80) {
	z80.l = z80.memory.readByte(z80.tempaddr) & 0xfb
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.l)
}
/* RES 2,(REGISTER+dd) */
func instrDDCB__RES_2_iREGpDD(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, bytetemp&0xfb)
}
/* LD A,RES 2,(REGISTER+dd) */
func instrDDCB__LD_A_RES_2_iREGpDD(z80 *Z80) {
	z80.a = z80.memory.readByte(z80.tempaddr) & 0xfb
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.a)
}
/* LD B,RES 3,(REGISTER+dd) */
func instrDDCB__LD_B_RES_3_iREGpDD(z80 *Z80) {
	z80.b = z80.memory.readByte(z80.tempaddr) & 0xf7
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.b)
}
/* LD C,RES 3,(REGISTER+dd) */
func instrDDCB__LD_C_RES_3_iREGpDD(z80 *Z80) {
	z80.c = z80.memory.readByte(z80.tempaddr) & 0xf7
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.c)
}
/* LD D,RES 3,(REGISTER+dd) */
func instrDDCB__LD_D_RES_3_iREGpDD(z80 *Z80) {
	z80.d = z80.memory.readByte(z80.tempaddr) & 0xf7
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.d)
}
/* LD E,RES 3,(REGISTER+dd) */
func instrDDCB__LD_E_RES_3_iREGpDD(z80 *Z80) {
	z80.e = z80.memory.readByte(z80.tempaddr) & 0xf7
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.e)
}
/* LD H,RES 3,(REGISTER+dd) */
func instrDDCB__LD_H_RES_3_iREGpDD(z80 *Z80) {
	z80.h = z80.memory.readByte(z80.tempaddr) & 0xf7
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.h)
}
/* LD L,RES 3,(REGISTER+dd) */
func instrDDCB__LD_L_RES_3_iREGpDD(z80 *Z80) {
	z80.l = z80.memory.readByte(z80.tempaddr) & 0xf7
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.l)
}
/* RES 3,(REGISTER+dd) */
func instrDDCB__RES_3_iREGpDD(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, bytetemp&0xf7)
}
/* LD A,RES 3,(REGISTER+dd) */
func instrDDCB__LD_A_RES_3_iREGpDD(z80 *Z80) {
	z80.a = z80.memory.readByte(z80.tempaddr) & 0xf7
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.a)
}
/* LD B,RES 4,(REGISTER+dd) */
func instrDDCB__LD_B_RES_4_iREGpDD(z80 *Z80) {
	z80.b = z80.memory.readByte(z80.tempaddr) & 0xef
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.b)
}
/* LD C,RES 4,(REGISTER+dd) */
func instrDDCB__LD_C_RES_4_iREGpDD(z80 *Z80) {
	z80.c = z80.memory.readByte(z80.tempaddr) & 0xef
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.c)
}
/* LD D,RES 4,(REGISTER+dd) */
func instrDDCB__LD_D_RES_4_iREGpDD(z80 *Z80) {
	z80.d = z80.memory.readByte(z80.tempaddr) & 0xef
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.d)
}
/* LD E,RES 4,(REGISTER+dd) */
func instrDDCB__LD_E_RES_4_iREGpDD(z80 *Z80) {
	z80.e = z80.memory.readByte(z80.tempaddr) & 0xef
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.e)
}
/* LD H,RES 4,(REGISTER+dd) */
func instrDDCB__LD_H_RES_4_iREGpDD(z80 *Z80) {
	z80.h = z80.memory.readByte(z80.tempaddr) & 0xef
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.h)
}
/* LD L,RES 4,(REGISTER+dd) */
func instrDDCB__LD_L_RES_4_iREGpDD(z80 *Z80) {
	z80.l = z80.memory.readByte(z80.tempaddr) & 0xef
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.l)
}
/* RES 4,(REGISTER+dd) */
func instrDDCB__RES_4_iREGpDD(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, bytetemp&0xef)
}
/* LD A,RES 4,(REGISTER+dd) */
func instrDDCB__LD_A_RES_4_iREGpDD(z80 *Z80) {
	z80.a = z80.memory.readByte(z80.tempaddr) & 0xef
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.a)
}
/* LD B,RES 5,(REGISTER+dd) */
func instrDDCB__LD_B_RES_5_iREGpDD(z80 *Z80) {
	z80.b = z80.memory.readByte(z80.tempaddr) & 0xdf
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.b)
}
/* LD C,RES 5,(REGISTER+dd) */
func instrDDCB__LD_C_RES_5_iREGpDD(z80 *Z80) {
	z80.c = z80.memory.readByte(z80.tempaddr) & 0xdf
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.c)
}
/* LD D,RES 5,(REGISTER+dd) */
func instrDDCB__LD_D_RES_5_iREGpDD(z80 *Z80) {
	z80.d = z80.memory.readByte(z80.tempaddr) & 0xdf
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.d)
}
/* LD E,RES 5,(REGISTER+dd) */
func instrDDCB__LD_E_RES_5_iREGpDD(z80 *Z80) {
	z80.e = z80.memory.readByte(z80.tempaddr) & 0xdf
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.e)
}
/* LD H,RES 5,(REGISTER+dd) */
func instrDDCB__LD_H_RES_5_iREGpDD(z80 *Z80) {
	z80.h = z80.memory.readByte(z80.tempaddr) & 0xdf
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.h)
}
/* LD L,RES 5,(REGISTER+dd) */
func instrDDCB__LD_L_RES_5_iREGpDD(z80 *Z80) {
	z80.l = z80.memory.readByte(z80.tempaddr) & 0xdf
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.l)
}
/* RES 5,(REGISTER+dd) */
func instrDDCB__RES_5_iREGpDD(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, bytetemp&0xdf)
}
/* LD A,RES 5,(REGISTER+dd) */
func instrDDCB__LD_A_RES_5_iREGpDD(z80 *Z80) {
	z80.a = z80.memory.readByte(z80.tempaddr) & 0xdf
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.a)
}
/* LD B,RES 6,(REGISTER+dd) */
func instrDDCB__LD_B_RES_6_iREGpDD(z80 *Z80) {
	z80.b = z80.memory.readByte(z80.tempaddr) & 0xbf
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.b)
}
/* LD C,RES 6,(REGISTER+dd) */
func instrDDCB__LD_C_RES_6_iREGpDD(z80 *Z80) {
	z80.c = z80.memory.readByte(z80.tempaddr) & 0xbf
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.c)
}
/* LD D,RES 6,(REGISTER+dd) */
func instrDDCB__LD_D_RES_6_iREGpDD(z80 *Z80) {
	z80.d = z80.memory.readByte(z80.tempaddr) & 0xbf
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.d)
}
/* LD E,RES 6,(REGISTER+dd) */
func instrDDCB__LD_E_RES_6_iREGpDD(z80 *Z80) {
	z80.e = z80.memory.readByte(z80.tempaddr) & 0xbf
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.e)
}
/* LD H,RES 6,(REGISTER+dd) */
func instrDDCB__LD_H_RES_6_iREGpDD(z80 *Z80) {
	z80.h = z80.memory.readByte(z80.tempaddr) & 0xbf
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.h)
}
/* LD L,RES 6,(REGISTER+dd) */
func instrDDCB__LD_L_RES_6_iREGpDD(z80 *Z80) {
	z80.l = z80.memory.readByte(z80.tempaddr) & 0xbf
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.l)
}
/* RES 6,(REGISTER+dd) */
func instrDDCB__RES_6_iREGpDD(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, bytetemp&0xbf)
}
/* LD A,RES 6,(REGISTER+dd) */
func instrDDCB__LD_A_RES_6_iREGpDD(z80 *Z80) {
	z80.a = z80.memory.readByte(z80.tempaddr) & 0xbf
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.a)
}
/* LD B,RES 7,(REGISTER+dd) */
func instrDDCB__LD_B_RES_7_iREGpDD(z80 *Z80) {
	z80.b = z80.memory.readByte(z80.tempaddr) & 0x7f
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.b)
}
/* LD C,RES 7,(REGISTER+dd) */
func instrDDCB__LD_C_RES_7_iREGpDD(z80 *Z80) {
	z80.c = z80.memory.readByte(z80.tempaddr) & 0x7f
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.c)
}
/* LD D,RES 7,(REGISTER+dd) */
func instrDDCB__LD_D_RES_7_iREGpDD(z80 *Z80) {
	z80.d = z80.memory.readByte(z80.tempaddr) & 0x7f
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.d)
}
/* LD E,RES 7,(REGISTER+dd) */
func instrDDCB__LD_E_RES_7_iREGpDD(z80 *Z80) {
	z80.e = z80.memory.readByte(z80.tempaddr) & 0x7f
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.e)
}
/* LD H,RES 7,(REGISTER+dd) */
func instrDDCB__LD_H_RES_7_iREGpDD(z80 *Z80) {
	z80.h = z80.memory.readByte(z80.tempaddr) & 0x7f
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.h)
}
/* LD L,RES 7,(REGISTER+dd) */
func instrDDCB__LD_L_RES_7_iREGpDD(z80 *Z80) {
	z80.l = z80.memory.readByte(z80.tempaddr) & 0x7f
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.l)
}
/* RES 7,(REGISTER+dd) */
func instrDDCB__RES_7_iREGpDD(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, bytetemp&0x7f)
}
/* LD A,RES 7,(REGISTER+dd) */
func instrDDCB__LD_A_RES_7_iREGpDD(z80 *Z80) {
	z80.a = z80.memory.readByte(z80.tempaddr) & 0x7f
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.a)
}
/* LD B,SET 0,(REGISTER+dd) */
func instrDDCB__LD_B_SET_0_iREGpDD(z80 *Z80) {
	z80.b = z80.memory.readByte(z80.tempaddr) | 0x01
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.b)
}
/* LD C,SET 0,(REGISTER+dd) */
func instrDDCB__LD_C_SET_0_iREGpDD(z80 *Z80) {
	z80.c = z80.memory.readByte(z80.tempaddr) | 0x01
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.c)
}
/* LD D,SET 0,(REGISTER+dd) */
func instrDDCB__LD_D_SET_0_iREGpDD(z80 *Z80) {
	z80.d = z80.memory.readByte(z80.tempaddr) | 0x01
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.d)
}
/* LD E,SET 0,(REGISTER+dd) */
func instrDDCB__LD_E_SET_0_iREGpDD(z80 *Z80) {
	z80.e = z80.memory.readByte(z80.tempaddr) | 0x01
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.e)
}
/* LD H,SET 0,(REGISTER+dd) */
func instrDDCB__LD_H_SET_0_iREGpDD(z80 *Z80) {
	z80.h = z80.memory.readByte(z80.tempaddr) | 0x01
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.h)
}
/* LD L,SET 0,(REGISTER+dd) */
func instrDDCB__LD_L_SET_0_iREGpDD(z80 *Z80) {
	z80.l = z80.memory.readByte(z80.tempaddr) | 0x01
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.l)
}
/* SET 0,(REGISTER+dd) */
func instrDDCB__SET_0_iREGpDD(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, bytetemp|0x01)
}
/* LD A,SET 0,(REGISTER+dd) */
func instrDDCB__LD_A_SET_0_iREGpDD(z80 *Z80) {
	z80.a = z80.memory.readByte(z80.tempaddr) | 0x01
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.a)
}
/* LD B,SET 1,(REGISTER+dd) */
func instrDDCB__LD_B_SET_1_iREGpDD(z80 *Z80) {
	z80.b = z80.memory.readByte(z80.tempaddr) | 0x02
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.b)
}
/* LD C,SET 1,(REGISTER+dd) */
func instrDDCB__LD_C_SET_1_iREGpDD(z80 *Z80) {
	z80.c = z80.memory.readByte(z80.tempaddr) | 0x02
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.c)
}
/* LD D,SET 1,(REGISTER+dd) */
func instrDDCB__LD_D_SET_1_iREGpDD(z80 *Z80) {
	z80.d = z80.memory.readByte(z80.tempaddr) | 0x02
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.d)
}
/* LD E,SET 1,(REGISTER+dd) */
func instrDDCB__LD_E_SET_1_iREGpDD(z80 *Z80) {
	z80.e = z80.memory.readByte(z80.tempaddr) | 0x02
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.e)
}
/* LD H,SET 1,(REGISTER+dd) */
func instrDDCB__LD_H_SET_1_iREGpDD(z80 *Z80) {
	z80.h = z80.memory.readByte(z80.tempaddr) | 0x02
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.h)
}
/* LD L,SET 1,(REGISTER+dd) */
func instrDDCB__LD_L_SET_1_iREGpDD(z80 *Z80) {
	z80.l = z80.memory.readByte(z80.tempaddr) | 0x02
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.l)
}
/* SET 1,(REGISTER+dd) */
func instrDDCB__SET_1_iREGpDD(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, bytetemp|0x02)
}
/* LD A,SET 1,(REGISTER+dd) */
func instrDDCB__LD_A_SET_1_iREGpDD(z80 *Z80) {
	z80.a = z80.memory.readByte(z80.tempaddr) | 0x02
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.a)
}
/* LD B,SET 2,(REGISTER+dd) */
func instrDDCB__LD_B_SET_2_iREGpDD(z80 *Z80) {
	z80.b = z80.memory.readByte(z80.tempaddr) | 0x04
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.b)
}
/* LD C,SET 2,(REGISTER+dd) */
func instrDDCB__LD_C_SET_2_iREGpDD(z80 *Z80) {
	z80.c = z80.memory.readByte(z80.tempaddr) | 0x04
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.c)
}
/* LD D,SET 2,(REGISTER+dd) */
func instrDDCB__LD_D_SET_2_iREGpDD(z80 *Z80) {
	z80.d = z80.memory.readByte(z80.tempaddr) | 0x04
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.d)
}
/* LD E,SET 2,(REGISTER+dd) */
func instrDDCB__LD_E_SET_2_iREGpDD(z80 *Z80) {
	z80.e = z80.memory.readByte(z80.tempaddr) | 0x04
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.e)
}
/* LD H,SET 2,(REGISTER+dd) */
func instrDDCB__LD_H_SET_2_iREGpDD(z80 *Z80) {
	z80.h = z80.memory.readByte(z80.tempaddr) | 0x04
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.h)
}
/* LD L,SET 2,(REGISTER+dd) */
func instrDDCB__LD_L_SET_2_iREGpDD(z80 *Z80) {
	z80.l = z80.memory.readByte(z80.tempaddr) | 0x04
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.l)
}
/* SET 2,(REGISTER+dd) */
func instrDDCB__SET_2_iREGpDD(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, bytetemp|0x04)
}
/* LD A,SET 2,(REGISTER+dd) */
func instrDDCB__LD_A_SET_2_iREGpDD(z80 *Z80) {
	z80.a = z80.memory.readByte(z80.tempaddr) | 0x04
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.a)
}
/* LD B,SET 3,(REGISTER+dd) */
func instrDDCB__LD_B_SET_3_iREGpDD(z80 *Z80) {
	z80.b = z80.memory.readByte(z80.tempaddr) | 0x08
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.b)
}
/* LD C,SET 3,(REGISTER+dd) */
func instrDDCB__LD_C_SET_3_iREGpDD(z80 *Z80) {
	z80.c = z80.memory.readByte(z80.tempaddr) | 0x08
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.c)
}
/* LD D,SET 3,(REGISTER+dd) */
func instrDDCB__LD_D_SET_3_iREGpDD(z80 *Z80) {
	z80.d = z80.memory.readByte(z80.tempaddr) | 0x08
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.d)
}
/* LD E,SET 3,(REGISTER+dd) */
func instrDDCB__LD_E_SET_3_iREGpDD(z80 *Z80) {
	z80.e = z80.memory.readByte(z80.tempaddr) | 0x08
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.e)
}
/* LD H,SET 3,(REGISTER+dd) */
func instrDDCB__LD_H_SET_3_iREGpDD(z80 *Z80) {
	z80.h = z80.memory.readByte(z80.tempaddr) | 0x08
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.h)
}
/* LD L,SET 3,(REGISTER+dd) */
func instrDDCB__LD_L_SET_3_iREGpDD(z80 *Z80) {
	z80.l = z80.memory.readByte(z80.tempaddr) | 0x08
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.l)
}
/* SET 3,(REGISTER+dd) */
func instrDDCB__SET_3_iREGpDD(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, bytetemp|0x08)
}
/* LD A,SET 3,(REGISTER+dd) */
func instrDDCB__LD_A_SET_3_iREGpDD(z80 *Z80) {
	z80.a = z80.memory.readByte(z80.tempaddr) | 0x08
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.a)
}
/* LD B,SET 4,(REGISTER+dd) */
func instrDDCB__LD_B_SET_4_iREGpDD(z80 *Z80) {
	z80.b = z80.memory.readByte(z80.tempaddr) | 0x10
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.b)
}
/* LD C,SET 4,(REGISTER+dd) */
func instrDDCB__LD_C_SET_4_iREGpDD(z80 *Z80) {
	z80.c = z80.memory.readByte(z80.tempaddr) | 0x10
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.c)
}
/* LD D,SET 4,(REGISTER+dd) */
func instrDDCB__LD_D_SET_4_iREGpDD(z80 *Z80) {
	z80.d = z80.memory.readByte(z80.tempaddr) | 0x10
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.d)
}
/* LD E,SET 4,(REGISTER+dd) */
func instrDDCB__LD_E_SET_4_iREGpDD(z80 *Z80) {
	z80.e = z80.memory.readByte(z80.tempaddr) | 0x10
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.e)
}
/* LD H,SET 4,(REGISTER+dd) */
func instrDDCB__LD_H_SET_4_iREGpDD(z80 *Z80) {
	z80.h = z80.memory.readByte(z80.tempaddr) | 0x10
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.h)
}
/* LD L,SET 4,(REGISTER+dd) */
func instrDDCB__LD_L_SET_4_iREGpDD(z80 *Z80) {
	z80.l = z80.memory.readByte(z80.tempaddr) | 0x10
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.l)
}
/* SET 4,(REGISTER+dd) */
func instrDDCB__SET_4_iREGpDD(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, bytetemp|0x10)
}
/* LD A,SET 4,(REGISTER+dd) */
func instrDDCB__LD_A_SET_4_iREGpDD(z80 *Z80) {
	z80.a = z80.memory.readByte(z80.tempaddr) | 0x10
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.a)
}
/* LD B,SET 5,(REGISTER+dd) */
func instrDDCB__LD_B_SET_5_iREGpDD(z80 *Z80) {
	z80.b = z80.memory.readByte(z80.tempaddr) | 0x20
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.b)
}
/* LD C,SET 5,(REGISTER+dd) */
func instrDDCB__LD_C_SET_5_iREGpDD(z80 *Z80) {
	z80.c = z80.memory.readByte(z80.tempaddr) | 0x20
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.c)
}
/* LD D,SET 5,(REGISTER+dd) */
func instrDDCB__LD_D_SET_5_iREGpDD(z80 *Z80) {
	z80.d = z80.memory.readByte(z80.tempaddr) | 0x20
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.d)
}
/* LD E,SET 5,(REGISTER+dd) */
func instrDDCB__LD_E_SET_5_iREGpDD(z80 *Z80) {
	z80.e = z80.memory.readByte(z80.tempaddr) | 0x20
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.e)
}
/* LD H,SET 5,(REGISTER+dd) */
func instrDDCB__LD_H_SET_5_iREGpDD(z80 *Z80) {
	z80.h = z80.memory.readByte(z80.tempaddr) | 0x20
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.h)
}
/* LD L,SET 5,(REGISTER+dd) */
func instrDDCB__LD_L_SET_5_iREGpDD(z80 *Z80) {
	z80.l = z80.memory.readByte(z80.tempaddr) | 0x20
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.l)
}
/* SET 5,(REGISTER+dd) */
func instrDDCB__SET_5_iREGpDD(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, bytetemp|0x20)
}
/* LD A,SET 5,(REGISTER+dd) */
func instrDDCB__LD_A_SET_5_iREGpDD(z80 *Z80) {
	z80.a = z80.memory.readByte(z80.tempaddr) | 0x20
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.a)
}
/* LD B,SET 6,(REGISTER+dd) */
func instrDDCB__LD_B_SET_6_iREGpDD(z80 *Z80) {
	z80.b = z80.memory.readByte(z80.tempaddr) | 0x40
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.b)
}
/* LD C,SET 6,(REGISTER+dd) */
func instrDDCB__LD_C_SET_6_iREGpDD(z80 *Z80) {
	z80.c = z80.memory.readByte(z80.tempaddr) | 0x40
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.c)
}
/* LD D,SET 6,(REGISTER+dd) */
func instrDDCB__LD_D_SET_6_iREGpDD(z80 *Z80) {
	z80.d = z80.memory.readByte(z80.tempaddr) | 0x40
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.d)
}
/* LD E,SET 6,(REGISTER+dd) */
func instrDDCB__LD_E_SET_6_iREGpDD(z80 *Z80) {
	z80.e = z80.memory.readByte(z80.tempaddr) | 0x40
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.e)
}
/* LD H,SET 6,(REGISTER+dd) */
func instrDDCB__LD_H_SET_6_iREGpDD(z80 *Z80) {
	z80.h = z80.memory.readByte(z80.tempaddr) | 0x40
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.h)
}
/* LD L,SET 6,(REGISTER+dd) */
func instrDDCB__LD_L_SET_6_iREGpDD(z80 *Z80) {
	z80.l = z80.memory.readByte(z80.tempaddr) | 0x40
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.l)
}
/* SET 6,(REGISTER+dd) */
func instrDDCB__SET_6_iREGpDD(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, bytetemp|0x40)
}
/* LD A,SET 6,(REGISTER+dd) */
func instrDDCB__LD_A_SET_6_iREGpDD(z80 *Z80) {
	z80.a = z80.memory.readByte(z80.tempaddr) | 0x40
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.a)
}
/* LD B,SET 7,(REGISTER+dd) */
func instrDDCB__LD_B_SET_7_iREGpDD(z80 *Z80) {
	z80.b = z80.memory.readByte(z80.tempaddr) | 0x80
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.b)
}
/* LD C,SET 7,(REGISTER+dd) */
func instrDDCB__LD_C_SET_7_iREGpDD(z80 *Z80) {
	z80.c = z80.memory.readByte(z80.tempaddr) | 0x80
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.c)
}
/* LD D,SET 7,(REGISTER+dd) */
func instrDDCB__LD_D_SET_7_iREGpDD(z80 *Z80) {
	z80.d = z80.memory.readByte(z80.tempaddr) | 0x80
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.d)
}
/* LD E,SET 7,(REGISTER+dd) */
func instrDDCB__LD_E_SET_7_iREGpDD(z80 *Z80) {
	z80.e = z80.memory.readByte(z80.tempaddr) | 0x80
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.e)
}
/* LD H,SET 7,(REGISTER+dd) */
func instrDDCB__LD_H_SET_7_iREGpDD(z80 *Z80) {
	z80.h = z80.memory.readByte(z80.tempaddr) | 0x80
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.h)
}
/* LD L,SET 7,(REGISTER+dd) */
func instrDDCB__LD_L_SET_7_iREGpDD(z80 *Z80) {
	z80.l = z80.memory.readByte(z80.tempaddr) | 0x80
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.l)
}
/* SET 7,(REGISTER+dd) */
func instrDDCB__SET_7_iREGpDD(z80 *Z80) {
	var bytetemp byte = z80.memory.readByte(z80.tempaddr)
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, bytetemp|0x80)
}
/* LD A,SET 7,(REGISTER+dd) */
func instrDDCB__LD_A_SET_7_iREGpDD(z80 *Z80) {
	z80.a = z80.memory.readByte(z80.tempaddr) | 0x80
	z80.memory.contendReadNoMreq(z80.tempaddr, 1)
	z80.memory.writeByte(z80.tempaddr, z80.a)
}
