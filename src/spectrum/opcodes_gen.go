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
	opcodesMap[0x00] = func(z80 *Z80) {
	}
	/* LD BC,nnnn */
	opcodesMap[0x01] = func(z80 *Z80) {
		b1 := z80.memory.readByte(z80.pc)
		z80.pc++
		b2 := z80.memory.readByte(z80.pc)
		z80.pc++
		z80.setBC(joinBytes(b2, b1))
	}
	/* LD (BC),A */
	opcodesMap[0x02] = func(z80 *Z80) {
		z80.memory.writeByte(z80.BC(), z80.a)
	}
	/* INC BC */
	opcodesMap[0x03] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 2)
		z80.incBC()
	}
	/* INC B */
	opcodesMap[0x04] = func(z80 *Z80) {
		z80.incB()
	}
	/* DEC B */
	opcodesMap[0x05] = func(z80 *Z80) {
		z80.decB()
	}
	/* LD B,nn */
	opcodesMap[0x06] = func(z80 *Z80) {
		z80.b = z80.memory.readByte(z80.pc)
		z80.pc++
	}
	/* RLCA */
	opcodesMap[0x07] = func(z80 *Z80) {
		z80.a = (z80.a << 1) | (z80.a >> 7)
		z80.f = (z80.f & (FLAG_P | FLAG_Z | FLAG_S)) |
			(z80.a & (FLAG_C | FLAG_3 | FLAG_5))
	}
	/* EX AF,AF' */
	opcodesMap[0x08] = func(z80 *Z80) {
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
	opcodesMap[0x09] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 7)
		z80.add16(z80.hl, z80.BC())
	}
	/* LD A,(BC) */
	opcodesMap[0x0a] = func(z80 *Z80) {
		z80.a = z80.memory.readByte(z80.BC())
	}
	/* DEC BC */
	opcodesMap[0x0b] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 2)
		z80.decBC()
	}
	/* INC C */
	opcodesMap[0x0c] = func(z80 *Z80) {
		z80.incC()
	}
	/* DEC C */
	opcodesMap[0x0d] = func(z80 *Z80) {
		z80.decC()
	}
	/* LD C,nn */
	opcodesMap[0x0e] = func(z80 *Z80) {
		z80.c = z80.memory.readByte(z80.pc)
		z80.pc++
	}
	/* RRCA */
	opcodesMap[0x0f] = func(z80 *Z80) {
		z80.f = (z80.f & (FLAG_P | FLAG_Z | FLAG_S)) | (z80.a & FLAG_C)
		z80.a = (z80.a >> 1) | (z80.a << 7)
		z80.f |= (z80.a & (FLAG_3 | FLAG_5))
	}
	/* DJNZ offset */
	opcodesMap[0x10] = func(z80 *Z80) {
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
	opcodesMap[0x11] = func(z80 *Z80) {
		b1 := z80.memory.readByte(z80.pc)
		z80.pc++
		b2 := z80.memory.readByte(z80.pc)
		z80.pc++
		z80.setDE(joinBytes(b2, b1))
	}
	/* LD (DE),A */
	opcodesMap[0x12] = func(z80 *Z80) {
		z80.memory.writeByte(z80.DE(), z80.a)
	}
	/* INC DE */
	opcodesMap[0x13] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 2)
		z80.incDE()
	}
	/* INC D */
	opcodesMap[0x14] = func(z80 *Z80) {
		z80.incD()
	}
	/* DEC D */
	opcodesMap[0x15] = func(z80 *Z80) {
		z80.decD()
	}
	/* LD D,nn */
	opcodesMap[0x16] = func(z80 *Z80) {
		z80.d = z80.memory.readByte(z80.pc)
		z80.pc++
	}
	/* RLA */
	opcodesMap[0x17] = func(z80 *Z80) {
		var bytetemp byte = z80.a
		z80.a = (z80.a << 1) | (z80.f & FLAG_C)
		z80.f = (z80.f & (FLAG_P | FLAG_Z | FLAG_S)) | (z80.a & (FLAG_3 | FLAG_5)) | (bytetemp >> 7)
	}
	/* JR offset */
	opcodesMap[0x18] = func(z80 *Z80) {
		z80.jr()
		z80.pc++
	}
	/* ADD HL,DE */
	opcodesMap[0x19] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 7)
		z80.add16(z80.hl, z80.DE())
	}
	/* LD A,(DE) */
	opcodesMap[0x1a] = func(z80 *Z80) {
		z80.a = z80.memory.readByte(z80.DE())
	}
	/* DEC DE */
	opcodesMap[0x1b] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 2)
		z80.decDE()
	}
	/* INC E */
	opcodesMap[0x1c] = func(z80 *Z80) {
		z80.incE()
	}
	/* DEC E */
	opcodesMap[0x1d] = func(z80 *Z80) {
		z80.decE()
	}
	/* LD E,nn */
	opcodesMap[0x1e] = func(z80 *Z80) {
		z80.e = z80.memory.readByte(z80.pc)
		z80.pc++
	}
	/* RRA */
	opcodesMap[0x1f] = func(z80 *Z80) {
		var bytetemp byte = z80.a
		z80.a = (z80.a >> 1) | (z80.f << 7)
		z80.f = (z80.f & (FLAG_P | FLAG_Z | FLAG_S)) | (z80.a & (FLAG_3 | FLAG_5)) | (bytetemp & FLAG_C)
	}
	/* JR NZ,offset */
	opcodesMap[0x20] = func(z80 *Z80) {
		if (z80.f & FLAG_Z) == 0 {
			z80.jr()
		} else {
			z80.memory.contendRead(z80.pc, 3)
		}
		z80.pc++
	}
	/* LD HL,nnnn */
	opcodesMap[0x21] = func(z80 *Z80) {
		b1 := z80.memory.readByte(z80.pc)
		z80.pc++
		b2 := z80.memory.readByte(z80.pc)
		z80.pc++
		z80.setHL(joinBytes(b2, b1))
	}
	/* LD (nnnn),HL */
	opcodesMap[0x22] = func(z80 *Z80) {
		z80.ld16nnrr(z80.l, z80.h)
		// break
	}
	/* INC HL */
	opcodesMap[0x23] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 2)
		z80.incHL()
	}
	/* INC H */
	opcodesMap[0x24] = func(z80 *Z80) {
		z80.incH()
	}
	/* DEC H */
	opcodesMap[0x25] = func(z80 *Z80) {
		z80.decH()
	}
	/* LD H,nn */
	opcodesMap[0x26] = func(z80 *Z80) {
		z80.h = z80.memory.readByte(z80.pc)
		z80.pc++
	}
	/* DAA */
	opcodesMap[0x27] = func(z80 *Z80) {
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
	opcodesMap[0x28] = func(z80 *Z80) {
		if (z80.f & FLAG_Z) != 0 {
			z80.jr()
		} else {
			z80.memory.contendRead(z80.pc, 3)
		}
		z80.pc++
	}
	/* ADD HL,HL */
	opcodesMap[0x29] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 7)
		z80.add16(z80.hl, z80.HL())
	}
	/* LD HL,(nnnn) */
	opcodesMap[0x2a] = func(z80 *Z80) {
		z80.ld16rrnn(&z80.l, &z80.h)
		// break
	}
	/* DEC HL */
	opcodesMap[0x2b] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 2)
		z80.decHL()
	}
	/* INC L */
	opcodesMap[0x2c] = func(z80 *Z80) {
		z80.incL()
	}
	/* DEC L */
	opcodesMap[0x2d] = func(z80 *Z80) {
		z80.decL()
	}
	/* LD L,nn */
	opcodesMap[0x2e] = func(z80 *Z80) {
		z80.l = z80.memory.readByte(z80.pc)
		z80.pc++
	}
	/* CPL */
	opcodesMap[0x2f] = func(z80 *Z80) {
		z80.a ^= 0xff
		z80.f = (z80.f & (FLAG_C | FLAG_P | FLAG_Z | FLAG_S)) |
			(z80.a & (FLAG_3 | FLAG_5)) |
			(FLAG_N | FLAG_H)
	}
	/* JR NC,offset */
	opcodesMap[0x30] = func(z80 *Z80) {
		if (z80.f & FLAG_C) == 0 {
			z80.jr()
		} else {
			z80.memory.contendRead(z80.pc, 3)
		}
		z80.pc++
	}
	/* LD SP,nnnn */
	opcodesMap[0x31] = func(z80 *Z80) {
		b1 := z80.memory.readByte(z80.pc)
		z80.pc++
		b2 := z80.memory.readByte(z80.pc)
		z80.pc++
		z80.setSP(joinBytes(b2, b1))
	}
	/* LD (nnnn),A */
	opcodesMap[0x32] = func(z80 *Z80) {
		var wordtemp uint16 = uint16(z80.memory.readByte(z80.pc))
		z80.pc++
		wordtemp |= uint16(z80.memory.readByte(z80.pc)) << 8
		z80.pc++
		z80.memory.writeByte(wordtemp, z80.a)
	}
	/* INC SP */
	opcodesMap[0x33] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 2)
		z80.incSP()
	}
	/* INC (HL) */
	opcodesMap[0x34] = func(z80 *Z80) {
		{
			var bytetemp byte = z80.memory.readByte(z80.HL())
			z80.memory.contendReadNoMreq(z80.HL(), 1)
			z80.inc(&bytetemp)
			z80.memory.writeByte(z80.HL(), bytetemp)
		}
	}
	/* DEC (HL) */
	opcodesMap[0x35] = func(z80 *Z80) {
		{
			var bytetemp byte = z80.memory.readByte(z80.HL())
			z80.memory.contendReadNoMreq(z80.HL(), 1)
			z80.dec(&bytetemp)
			z80.memory.writeByte(z80.HL(), bytetemp)
		}
	}
	/* LD (HL),nn */
	opcodesMap[0x36] = func(z80 *Z80) {
		z80.memory.writeByte(z80.HL(), z80.memory.readByte(z80.pc))
		z80.pc++
	}
	/* SCF */
	opcodesMap[0x37] = func(z80 *Z80) {
		z80.f = (z80.f & (FLAG_P | FLAG_Z | FLAG_S)) |
			(z80.a & (FLAG_3 | FLAG_5)) |
			FLAG_C
	}
	/* JR C,offset */
	opcodesMap[0x38] = func(z80 *Z80) {
		if (z80.f & FLAG_C) != 0 {
			z80.jr()
		} else {
			z80.memory.contendRead(z80.pc, 3)
		}
		z80.pc++
	}
	/* ADD HL,SP */
	opcodesMap[0x39] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 7)
		z80.add16(z80.hl, z80.SP())
	}
	/* LD A,(nnnn) */
	opcodesMap[0x3a] = func(z80 *Z80) {
		var wordtemp uint16 = uint16(z80.memory.readByte(z80.pc))
		z80.pc++
		wordtemp |= uint16(z80.memory.readByte(z80.pc)) << 8
		z80.pc++
		z80.a = z80.memory.readByte(wordtemp)
	}
	/* DEC SP */
	opcodesMap[0x3b] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 2)
		z80.decSP()
	}
	/* INC A */
	opcodesMap[0x3c] = func(z80 *Z80) {
		z80.incA()
	}
	/* DEC A */
	opcodesMap[0x3d] = func(z80 *Z80) {
		z80.decA()
	}
	/* LD A,nn */
	opcodesMap[0x3e] = func(z80 *Z80) {
		z80.a = z80.memory.readByte(z80.pc)
		z80.pc++
	}
	/* CCF */
	opcodesMap[0x3f] = func(z80 *Z80) {
		z80.f = (z80.f & (FLAG_P | FLAG_Z | FLAG_S)) |
			ternOpB((z80.f&FLAG_C) != 0, FLAG_H, FLAG_C) |
			(z80.a & (FLAG_3 | FLAG_5))
	}
	/* LD B,B */
	opcodesMap[0x40] = func(z80 *Z80) {
	}
	/* LD B,C */
	opcodesMap[0x41] = func(z80 *Z80) {
		z80.b = z80.c
	}
	/* LD B,D */
	opcodesMap[0x42] = func(z80 *Z80) {
		z80.b = z80.d
	}
	/* LD B,E */
	opcodesMap[0x43] = func(z80 *Z80) {
		z80.b = z80.e
	}
	/* LD B,H */
	opcodesMap[0x44] = func(z80 *Z80) {
		z80.b = z80.h
	}
	/* LD B,L */
	opcodesMap[0x45] = func(z80 *Z80) {
		z80.b = z80.l
	}
	/* LD B,(HL) */
	opcodesMap[0x46] = func(z80 *Z80) {
		z80.b = z80.memory.readByte(z80.HL())
	}
	/* LD B,A */
	opcodesMap[0x47] = func(z80 *Z80) {
		z80.b = z80.a
	}
	/* LD C,B */
	opcodesMap[0x48] = func(z80 *Z80) {
		z80.c = z80.b
	}
	/* LD C,C */
	opcodesMap[0x49] = func(z80 *Z80) {
	}
	/* LD C,D */
	opcodesMap[0x4a] = func(z80 *Z80) {
		z80.c = z80.d
	}
	/* LD C,E */
	opcodesMap[0x4b] = func(z80 *Z80) {
		z80.c = z80.e
	}
	/* LD C,H */
	opcodesMap[0x4c] = func(z80 *Z80) {
		z80.c = z80.h
	}
	/* LD C,L */
	opcodesMap[0x4d] = func(z80 *Z80) {
		z80.c = z80.l
	}
	/* LD C,(HL) */
	opcodesMap[0x4e] = func(z80 *Z80) {
		z80.c = z80.memory.readByte(z80.HL())
	}
	/* LD C,A */
	opcodesMap[0x4f] = func(z80 *Z80) {
		z80.c = z80.a
	}
	/* LD D,B */
	opcodesMap[0x50] = func(z80 *Z80) {
		z80.d = z80.b
	}
	/* LD D,C */
	opcodesMap[0x51] = func(z80 *Z80) {
		z80.d = z80.c
	}
	/* LD D,D */
	opcodesMap[0x52] = func(z80 *Z80) {
	}
	/* LD D,E */
	opcodesMap[0x53] = func(z80 *Z80) {
		z80.d = z80.e
	}
	/* LD D,H */
	opcodesMap[0x54] = func(z80 *Z80) {
		z80.d = z80.h
	}
	/* LD D,L */
	opcodesMap[0x55] = func(z80 *Z80) {
		z80.d = z80.l
	}
	/* LD D,(HL) */
	opcodesMap[0x56] = func(z80 *Z80) {
		z80.d = z80.memory.readByte(z80.HL())
	}
	/* LD D,A */
	opcodesMap[0x57] = func(z80 *Z80) {
		z80.d = z80.a
	}
	/* LD E,B */
	opcodesMap[0x58] = func(z80 *Z80) {
		z80.e = z80.b
	}
	/* LD E,C */
	opcodesMap[0x59] = func(z80 *Z80) {
		z80.e = z80.c
	}
	/* LD E,D */
	opcodesMap[0x5a] = func(z80 *Z80) {
		z80.e = z80.d
	}
	/* LD E,E */
	opcodesMap[0x5b] = func(z80 *Z80) {
	}
	/* LD E,H */
	opcodesMap[0x5c] = func(z80 *Z80) {
		z80.e = z80.h
	}
	/* LD E,L */
	opcodesMap[0x5d] = func(z80 *Z80) {
		z80.e = z80.l
	}
	/* LD E,(HL) */
	opcodesMap[0x5e] = func(z80 *Z80) {
		z80.e = z80.memory.readByte(z80.HL())
	}
	/* LD E,A */
	opcodesMap[0x5f] = func(z80 *Z80) {
		z80.e = z80.a
	}
	/* LD H,B */
	opcodesMap[0x60] = func(z80 *Z80) {
		z80.h = z80.b
	}
	/* LD H,C */
	opcodesMap[0x61] = func(z80 *Z80) {
		z80.h = z80.c
	}
	/* LD H,D */
	opcodesMap[0x62] = func(z80 *Z80) {
		z80.h = z80.d
	}
	/* LD H,E */
	opcodesMap[0x63] = func(z80 *Z80) {
		z80.h = z80.e
	}
	/* LD H,H */
	opcodesMap[0x64] = func(z80 *Z80) {
	}
	/* LD H,L */
	opcodesMap[0x65] = func(z80 *Z80) {
		z80.h = z80.l
	}
	/* LD H,(HL) */
	opcodesMap[0x66] = func(z80 *Z80) {
		z80.h = z80.memory.readByte(z80.HL())
	}
	/* LD H,A */
	opcodesMap[0x67] = func(z80 *Z80) {
		z80.h = z80.a
	}
	/* LD L,B */
	opcodesMap[0x68] = func(z80 *Z80) {
		z80.l = z80.b
	}
	/* LD L,C */
	opcodesMap[0x69] = func(z80 *Z80) {
		z80.l = z80.c
	}
	/* LD L,D */
	opcodesMap[0x6a] = func(z80 *Z80) {
		z80.l = z80.d
	}
	/* LD L,E */
	opcodesMap[0x6b] = func(z80 *Z80) {
		z80.l = z80.e
	}
	/* LD L,H */
	opcodesMap[0x6c] = func(z80 *Z80) {
		z80.l = z80.h
	}
	/* LD L,L */
	opcodesMap[0x6d] = func(z80 *Z80) {
	}
	/* LD L,(HL) */
	opcodesMap[0x6e] = func(z80 *Z80) {
		z80.l = z80.memory.readByte(z80.HL())
	}
	/* LD L,A */
	opcodesMap[0x6f] = func(z80 *Z80) {
		z80.l = z80.a
	}
	/* LD (HL),B */
	opcodesMap[0x70] = func(z80 *Z80) {
		z80.memory.writeByte(z80.HL(), z80.b)
	}
	/* LD (HL),C */
	opcodesMap[0x71] = func(z80 *Z80) {
		z80.memory.writeByte(z80.HL(), z80.c)
	}
	/* LD (HL),D */
	opcodesMap[0x72] = func(z80 *Z80) {
		z80.memory.writeByte(z80.HL(), z80.d)
	}
	/* LD (HL),E */
	opcodesMap[0x73] = func(z80 *Z80) {
		z80.memory.writeByte(z80.HL(), z80.e)
	}
	/* LD (HL),H */
	opcodesMap[0x74] = func(z80 *Z80) {
		z80.memory.writeByte(z80.HL(), z80.h)
	}
	/* LD (HL),L */
	opcodesMap[0x75] = func(z80 *Z80) {
		z80.memory.writeByte(z80.HL(), z80.l)
	}
	/* HALT */
	opcodesMap[0x76] = func(z80 *Z80) {
		z80.halted = true
		z80.pc--
		return
	}
	/* LD (HL),A */
	opcodesMap[0x77] = func(z80 *Z80) {
		z80.memory.writeByte(z80.HL(), z80.a)
	}
	/* LD A,B */
	opcodesMap[0x78] = func(z80 *Z80) {
		z80.a = z80.b
	}
	/* LD A,C */
	opcodesMap[0x79] = func(z80 *Z80) {
		z80.a = z80.c
	}
	/* LD A,D */
	opcodesMap[0x7a] = func(z80 *Z80) {
		z80.a = z80.d
	}
	/* LD A,E */
	opcodesMap[0x7b] = func(z80 *Z80) {
		z80.a = z80.e
	}
	/* LD A,H */
	opcodesMap[0x7c] = func(z80 *Z80) {
		z80.a = z80.h
	}
	/* LD A,L */
	opcodesMap[0x7d] = func(z80 *Z80) {
		z80.a = z80.l
	}
	/* LD A,(HL) */
	opcodesMap[0x7e] = func(z80 *Z80) {
		z80.a = z80.memory.readByte(z80.HL())
	}
	/* LD A,A */
	opcodesMap[0x7f] = func(z80 *Z80) {
	}
	/* ADD A,B */
	opcodesMap[0x80] = func(z80 *Z80) {
		z80.add(z80.b)
	}
	/* ADD A,C */
	opcodesMap[0x81] = func(z80 *Z80) {
		z80.add(z80.c)
	}
	/* ADD A,D */
	opcodesMap[0x82] = func(z80 *Z80) {
		z80.add(z80.d)
	}
	/* ADD A,E */
	opcodesMap[0x83] = func(z80 *Z80) {
		z80.add(z80.e)
	}
	/* ADD A,H */
	opcodesMap[0x84] = func(z80 *Z80) {
		z80.add(z80.h)
	}
	/* ADD A,L */
	opcodesMap[0x85] = func(z80 *Z80) {
		z80.add(z80.l)
	}
	/* ADD A,(HL) */
	opcodesMap[0x86] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.HL())

		z80.add(bytetemp)
	}
	/* ADD A,A */
	opcodesMap[0x87] = func(z80 *Z80) {
		z80.add(z80.a)
	}
	/* ADC A,B */
	opcodesMap[0x88] = func(z80 *Z80) {
		z80.adc(z80.b)
	}
	/* ADC A,C */
	opcodesMap[0x89] = func(z80 *Z80) {
		z80.adc(z80.c)
	}
	/* ADC A,D */
	opcodesMap[0x8a] = func(z80 *Z80) {
		z80.adc(z80.d)
	}
	/* ADC A,E */
	opcodesMap[0x8b] = func(z80 *Z80) {
		z80.adc(z80.e)
	}
	/* ADC A,H */
	opcodesMap[0x8c] = func(z80 *Z80) {
		z80.adc(z80.h)
	}
	/* ADC A,L */
	opcodesMap[0x8d] = func(z80 *Z80) {
		z80.adc(z80.l)
	}
	/* ADC A,(HL) */
	opcodesMap[0x8e] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.HL())

		z80.adc(bytetemp)
	}
	/* ADC A,A */
	opcodesMap[0x8f] = func(z80 *Z80) {
		z80.adc(z80.a)
	}
	/* SUB A,B */
	opcodesMap[0x90] = func(z80 *Z80) {
		z80.sub(z80.b)
	}
	/* SUB A,C */
	opcodesMap[0x91] = func(z80 *Z80) {
		z80.sub(z80.c)
	}
	/* SUB A,D */
	opcodesMap[0x92] = func(z80 *Z80) {
		z80.sub(z80.d)
	}
	/* SUB A,E */
	opcodesMap[0x93] = func(z80 *Z80) {
		z80.sub(z80.e)
	}
	/* SUB A,H */
	opcodesMap[0x94] = func(z80 *Z80) {
		z80.sub(z80.h)
	}
	/* SUB A,L */
	opcodesMap[0x95] = func(z80 *Z80) {
		z80.sub(z80.l)
	}
	/* SUB A,(HL) */
	opcodesMap[0x96] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.HL())

		z80.sub(bytetemp)
	}
	/* SUB A,A */
	opcodesMap[0x97] = func(z80 *Z80) {
		z80.sub(z80.a)
	}
	/* SBC A,B */
	opcodesMap[0x98] = func(z80 *Z80) {
		z80.sbc(z80.b)
	}
	/* SBC A,C */
	opcodesMap[0x99] = func(z80 *Z80) {
		z80.sbc(z80.c)
	}
	/* SBC A,D */
	opcodesMap[0x9a] = func(z80 *Z80) {
		z80.sbc(z80.d)
	}
	/* SBC A,E */
	opcodesMap[0x9b] = func(z80 *Z80) {
		z80.sbc(z80.e)
	}
	/* SBC A,H */
	opcodesMap[0x9c] = func(z80 *Z80) {
		z80.sbc(z80.h)
	}
	/* SBC A,L */
	opcodesMap[0x9d] = func(z80 *Z80) {
		z80.sbc(z80.l)
	}
	/* SBC A,(HL) */
	opcodesMap[0x9e] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.HL())

		z80.sbc(bytetemp)
	}
	/* SBC A,A */
	opcodesMap[0x9f] = func(z80 *Z80) {
		z80.sbc(z80.a)
	}
	/* AND A,B */
	opcodesMap[0xa0] = func(z80 *Z80) {
		z80.and(z80.b)
	}
	/* AND A,C */
	opcodesMap[0xa1] = func(z80 *Z80) {
		z80.and(z80.c)
	}
	/* AND A,D */
	opcodesMap[0xa2] = func(z80 *Z80) {
		z80.and(z80.d)
	}
	/* AND A,E */
	opcodesMap[0xa3] = func(z80 *Z80) {
		z80.and(z80.e)
	}
	/* AND A,H */
	opcodesMap[0xa4] = func(z80 *Z80) {
		z80.and(z80.h)
	}
	/* AND A,L */
	opcodesMap[0xa5] = func(z80 *Z80) {
		z80.and(z80.l)
	}
	/* AND A,(HL) */
	opcodesMap[0xa6] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.HL())

		z80.and(bytetemp)
	}
	/* AND A,A */
	opcodesMap[0xa7] = func(z80 *Z80) {
		z80.and(z80.a)
	}
	/* XOR A,B */
	opcodesMap[0xa8] = func(z80 *Z80) {
		z80.xor(z80.b)
	}
	/* XOR A,C */
	opcodesMap[0xa9] = func(z80 *Z80) {
		z80.xor(z80.c)
	}
	/* XOR A,D */
	opcodesMap[0xaa] = func(z80 *Z80) {
		z80.xor(z80.d)
	}
	/* XOR A,E */
	opcodesMap[0xab] = func(z80 *Z80) {
		z80.xor(z80.e)
	}
	/* XOR A,H */
	opcodesMap[0xac] = func(z80 *Z80) {
		z80.xor(z80.h)
	}
	/* XOR A,L */
	opcodesMap[0xad] = func(z80 *Z80) {
		z80.xor(z80.l)
	}
	/* XOR A,(HL) */
	opcodesMap[0xae] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.HL())

		z80.xor(bytetemp)
	}
	/* XOR A,A */
	opcodesMap[0xaf] = func(z80 *Z80) {
		z80.xor(z80.a)
	}
	/* OR A,B */
	opcodesMap[0xb0] = func(z80 *Z80) {
		z80.or(z80.b)
	}
	/* OR A,C */
	opcodesMap[0xb1] = func(z80 *Z80) {
		z80.or(z80.c)
	}
	/* OR A,D */
	opcodesMap[0xb2] = func(z80 *Z80) {
		z80.or(z80.d)
	}
	/* OR A,E */
	opcodesMap[0xb3] = func(z80 *Z80) {
		z80.or(z80.e)
	}
	/* OR A,H */
	opcodesMap[0xb4] = func(z80 *Z80) {
		z80.or(z80.h)
	}
	/* OR A,L */
	opcodesMap[0xb5] = func(z80 *Z80) {
		z80.or(z80.l)
	}
	/* OR A,(HL) */
	opcodesMap[0xb6] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.HL())

		z80.or(bytetemp)
	}
	/* OR A,A */
	opcodesMap[0xb7] = func(z80 *Z80) {
		z80.or(z80.a)
	}
	/* CP B */
	opcodesMap[0xb8] = func(z80 *Z80) {
		z80.cp(z80.b)
	}
	/* CP C */
	opcodesMap[0xb9] = func(z80 *Z80) {
		z80.cp(z80.c)
	}
	/* CP D */
	opcodesMap[0xba] = func(z80 *Z80) {
		z80.cp(z80.d)
	}
	/* CP E */
	opcodesMap[0xbb] = func(z80 *Z80) {
		z80.cp(z80.e)
	}
	/* CP H */
	opcodesMap[0xbc] = func(z80 *Z80) {
		z80.cp(z80.h)
	}
	/* CP L */
	opcodesMap[0xbd] = func(z80 *Z80) {
		z80.cp(z80.l)
	}
	/* CP (HL) */
	opcodesMap[0xbe] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.HL())

		z80.cp(bytetemp)
	}
	/* CP A */
	opcodesMap[0xbf] = func(z80 *Z80) {
		z80.cp(z80.a)
	}
	/* RET NZ */
	opcodesMap[0xc0] = func(z80 *Z80) {
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
	opcodesMap[0xc1] = func(z80 *Z80) {
		z80.c, z80.b = z80.pop16()
	}
	/* JP NZ,nnnn */
	opcodesMap[0xc2] = func(z80 *Z80) {
		if (z80.f & FLAG_Z) == 0 {
			z80.jp()
		} else {
			z80.memory.contendRead(z80.pc, 3)
			z80.memory.contendRead(z80.pc+1, 3)
			z80.pc += 2
		}
	}
	/* JP nnnn */
	opcodesMap[0xc3] = func(z80 *Z80) {
		z80.jp()
	}
	/* CALL NZ,nnnn */
	opcodesMap[0xc4] = func(z80 *Z80) {
		if (z80.f & FLAG_Z) == 0 {
			z80.call()
		} else {
			z80.memory.contendRead(z80.pc, 3)
			z80.memory.contendRead(z80.pc+1, 3)
			z80.pc += 2
		}
	}
	/* PUSH BC */
	opcodesMap[0xc5] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq(z80.IR(), 1)
		z80.push16(z80.c, z80.b)
	}
	/* ADD A,nn */
	opcodesMap[0xc6] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.PC())
		z80.pc++
		z80.add(bytetemp)
	}
	/* RST 00 */
	opcodesMap[0xc7] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq(z80.IR(), 1)
		z80.rst(0x00)
	}
	/* RET Z */
	opcodesMap[0xc8] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq(z80.IR(), 1)
		if (z80.f & FLAG_Z) != 0 {
			z80.ret()
		}
	}
	/* RET */
	opcodesMap[0xc9] = func(z80 *Z80) {
		z80.ret()
	}
	/* JP Z,nnnn */
	opcodesMap[0xca] = func(z80 *Z80) {
		if (z80.f & FLAG_Z) != 0 {
			z80.jp()
		} else {
			z80.memory.contendRead(z80.pc, 3)
			z80.memory.contendRead(z80.pc+1, 3)
			z80.pc += 2
		}
	}
	/* shift CB */
	opcodesMap[0xcb] = func(z80 *Z80) {
	}
	/* CALL Z,nnnn */
	opcodesMap[0xcc] = func(z80 *Z80) {
		if (z80.f & FLAG_Z) != 0 {
			z80.call()
		} else {
			z80.memory.contendRead(z80.pc, 3)
			z80.memory.contendRead(z80.pc+1, 3)
			z80.pc += 2
		}
	}
	/* CALL nnnn */
	opcodesMap[0xcd] = func(z80 *Z80) {
		z80.call()
	}
	/* ADC A,nn */
	opcodesMap[0xce] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.PC())
		z80.pc++
		z80.adc(bytetemp)
	}
	/* RST 8 */
	opcodesMap[0xcf] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq(z80.IR(), 1)
		z80.rst(0x8)
	}
	/* RET NC */
	opcodesMap[0xd0] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq(z80.IR(), 1)
		if !((z80.f & FLAG_C) != 0) {
			z80.ret()
		}
	}
	/* POP DE */
	opcodesMap[0xd1] = func(z80 *Z80) {
		z80.e, z80.d = z80.pop16()
	}
	/* JP NC,nnnn */
	opcodesMap[0xd2] = func(z80 *Z80) {
		if (z80.f & FLAG_C) == 0 {
			z80.jp()
		} else {
			z80.memory.contendRead(z80.pc, 3)
			z80.memory.contendRead(z80.pc+1, 3)
			z80.pc += 2
		}
	}
	/* OUT (nn),A */
	opcodesMap[0xd3] = func(z80 *Z80) {
		var outtemp uint16 = uint16(z80.memory.readByte(z80.pc)) + (uint16(z80.a) << 8)
		z80.pc++
		z80.writePort(outtemp, z80.a)
	}
	/* CALL NC,nnnn */
	opcodesMap[0xd4] = func(z80 *Z80) {
		if (z80.f & FLAG_C) == 0 {
			z80.call()
		} else {
			z80.memory.contendRead(z80.pc, 3)
			z80.memory.contendRead(z80.pc+1, 3)
			z80.pc += 2
		}
	}
	/* PUSH DE */
	opcodesMap[0xd5] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq(z80.IR(), 1)
		z80.push16(z80.e, z80.d)
	}
	/* SUB nn */
	opcodesMap[0xd6] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.PC())
		z80.pc++
		z80.sub(bytetemp)
	}
	/* RST 10 */
	opcodesMap[0xd7] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq(z80.IR(), 1)
		z80.rst(0x10)
	}
	/* RET C */
	opcodesMap[0xd8] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq(z80.IR(), 1)
		if (z80.f & FLAG_C) != 0 {
			z80.ret()
		}
	}
	/* EXX */
	opcodesMap[0xd9] = func(z80 *Z80) {
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
	opcodesMap[0xda] = func(z80 *Z80) {
		if (z80.f & FLAG_C) != 0 {
			z80.jp()
		} else {
			z80.memory.contendRead(z80.pc, 3)
			z80.memory.contendRead(z80.pc+1, 3)
			z80.pc += 2
		}
	}
	/* IN A,(nn) */
	opcodesMap[0xdb] = func(z80 *Z80) {
		var intemp uint16 = uint16(z80.memory.readByte(z80.pc)) + (uint16(z80.a) << 8)
		z80.pc++
		z80.a = z80.readPort(intemp)
	}
	/* CALL C,nnnn */
	opcodesMap[0xdc] = func(z80 *Z80) {
		if (z80.f & FLAG_C) != 0 {
			z80.call()
		} else {
			z80.memory.contendRead(z80.pc, 3)
			z80.memory.contendRead(z80.pc+1, 3)
			z80.pc += 2
		}
	}
	/* shift DD */
	opcodesMap[0xdd] = func(z80 *Z80) {
	}
	/* SBC A,nn */
	opcodesMap[0xde] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.PC())
		z80.pc++
		z80.sbc(bytetemp)
	}
	/* RST 18 */
	opcodesMap[0xdf] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq(z80.IR(), 1)
		z80.rst(0x18)
	}
	/* RET PO */
	opcodesMap[0xe0] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq(z80.IR(), 1)
		if !((z80.f & FLAG_P) != 0) {
			z80.ret()
		}
	}
	/* POP HL */
	opcodesMap[0xe1] = func(z80 *Z80) {
		z80.l, z80.h = z80.pop16()
	}
	/* JP PO,nnnn */
	opcodesMap[0xe2] = func(z80 *Z80) {
		if (z80.f & FLAG_P) == 0 {
			z80.jp()
		} else {
			z80.memory.contendRead(z80.pc, 3)
			z80.memory.contendRead(z80.pc+1, 3)
			z80.pc += 2
		}
	}
	/* EX (SP),HL */
	opcodesMap[0xe3] = func(z80 *Z80) {
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
	opcodesMap[0xe4] = func(z80 *Z80) {
		if (z80.f & FLAG_P) == 0 {
			z80.call()
		} else {
			z80.memory.contendRead(z80.pc, 3)
			z80.memory.contendRead(z80.pc+1, 3)
			z80.pc += 2
		}
	}
	/* PUSH HL */
	opcodesMap[0xe5] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq(z80.IR(), 1)
		z80.push16(z80.l, z80.h)
	}
	/* AND nn */
	opcodesMap[0xe6] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.PC())
		z80.pc++
		z80.and(bytetemp)
	}
	/* RST 20 */
	opcodesMap[0xe7] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq(z80.IR(), 1)
		z80.rst(0x20)
	}
	/* RET PE */
	opcodesMap[0xe8] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq(z80.IR(), 1)
		if (z80.f & FLAG_P) != 0 {
			z80.ret()
		}
	}
	/* JP HL */
	opcodesMap[0xe9] = func(z80 *Z80) {
		z80.pc = z80.HL() /* NB: NOT INDIRECT! */
	}
	/* JP PE,nnnn */
	opcodesMap[0xea] = func(z80 *Z80) {
		if (z80.f & FLAG_P) != 0 {
			z80.jp()
		} else {
			z80.memory.contendRead(z80.pc, 3)
			z80.memory.contendRead(z80.pc+1, 3)
			z80.pc += 2
		}
	}
	/* EX DE,HL */
	opcodesMap[0xeb] = func(z80 *Z80) {
		var wordtemp uint16 = z80.DE()
		z80.setDE(z80.HL())
		z80.setHL(wordtemp)
	}
	/* CALL PE,nnnn */
	opcodesMap[0xec] = func(z80 *Z80) {
		if (z80.f & FLAG_P) != 0 {
			z80.call()
		} else {
			z80.memory.contendRead(z80.pc, 3)
			z80.memory.contendRead(z80.pc+1, 3)
			z80.pc += 2
		}
	}
	/* shift ED */
	opcodesMap[0xed] = func(z80 *Z80) {
	}
	/* XOR A,nn */
	opcodesMap[0xee] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.PC())
		z80.pc++
		z80.xor(bytetemp)
	}
	/* RST 28 */
	opcodesMap[0xef] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq(z80.IR(), 1)
		z80.rst(0x28)
	}
	/* RET P */
	opcodesMap[0xf0] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq(z80.IR(), 1)
		if !((z80.f & FLAG_S) != 0) {
			z80.ret()
		}
	}
	/* POP AF */
	opcodesMap[0xf1] = func(z80 *Z80) {
		z80.f, z80.a = z80.pop16()
	}
	/* JP P,nnnn */
	opcodesMap[0xf2] = func(z80 *Z80) {
		if (z80.f & FLAG_S) == 0 {
			z80.jp()
		} else {
			z80.memory.contendRead(z80.pc, 3)
			z80.memory.contendRead(z80.pc+1, 3)
			z80.pc += 2
		}
	}
	/* DI */
	opcodesMap[0xf3] = func(z80 *Z80) {
		z80.iff1, z80.iff2 = 0, 0
	}
	/* CALL P,nnnn */
	opcodesMap[0xf4] = func(z80 *Z80) {
		if (z80.f & FLAG_S) == 0 {
			z80.call()
		} else {
			z80.memory.contendRead(z80.pc, 3)
			z80.memory.contendRead(z80.pc+1, 3)
			z80.pc += 2
		}
	}
	/* PUSH AF */
	opcodesMap[0xf5] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq(z80.IR(), 1)
		z80.push16(z80.f, z80.a)
	}
	/* OR nn */
	opcodesMap[0xf6] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.PC())
		z80.pc++
		z80.or(bytetemp)
	}
	/* RST 30 */
	opcodesMap[0xf7] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq(z80.IR(), 1)
		z80.rst(0x30)
	}
	/* RET M */
	opcodesMap[0xf8] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq(z80.IR(), 1)
		if (z80.f & FLAG_S) != 0 {
			z80.ret()
		}
	}
	/* LD SP,HL */
	opcodesMap[0xf9] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 2)
		z80.sp = z80.HL()
	}
	/* JP M,nnnn */
	opcodesMap[0xfa] = func(z80 *Z80) {
		if (z80.f & FLAG_S) != 0 {
			z80.jp()
		} else {
			z80.memory.contendRead(z80.pc, 3)
			z80.memory.contendRead(z80.pc+1, 3)
			z80.pc += 2
		}
	}
	/* EI */
	opcodesMap[0xfb] = func(z80 *Z80) {
		/* Interrupts are not accepted immediately after an EI, but are
		   accepted after the next instruction */
		z80.iff1, z80.iff2 = 1, 1
		z80.interruptsEnabledAt = int(z80.tstates)
		// eventAdd(z80.tstates + 1, z80InterruptEvent)
	}
	/* CALL M,nnnn */
	opcodesMap[0xfc] = func(z80 *Z80) {
		if (z80.f & FLAG_S) != 0 {
			z80.call()
		} else {
			z80.memory.contendRead(z80.pc, 3)
			z80.memory.contendRead(z80.pc+1, 3)
			z80.pc += 2
		}
	}
	/* shift FD */
	opcodesMap[0xfd] = func(z80 *Z80) {
	}
	/* CP nn */
	opcodesMap[0xfe] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.PC())
		z80.pc++
		z80.cp(bytetemp)
	}
	/* RST 38 */
	opcodesMap[0xff] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq(z80.IR(), 1)
		z80.rst(0x38)
	}

	// END of non shifted opcodes

	// BEGIN of 0xcb shifted opcodes

	/* RLC B */
	opcodesMap[shift0xcb(0x00)] = func(z80 *Z80) {
		z80.b = z80.rlc(z80.b)
	}
	/* RLC C */
	opcodesMap[shift0xcb(0x01)] = func(z80 *Z80) {
		z80.c = z80.rlc(z80.c)
	}
	/* RLC D */
	opcodesMap[shift0xcb(0x02)] = func(z80 *Z80) {
		z80.d = z80.rlc(z80.d)
	}
	/* RLC E */
	opcodesMap[shift0xcb(0x03)] = func(z80 *Z80) {
		z80.e = z80.rlc(z80.e)
	}
	/* RLC H */
	opcodesMap[shift0xcb(0x04)] = func(z80 *Z80) {
		z80.h = z80.rlc(z80.h)
	}
	/* RLC L */
	opcodesMap[shift0xcb(0x05)] = func(z80 *Z80) {
		z80.l = z80.rlc(z80.l)
	}
	/* RLC (HL) */
	opcodesMap[shift0xcb(0x06)] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.HL())
		z80.memory.contendReadNoMreq(z80.HL(), 1)
		bytetemp = z80.rlc(bytetemp)
		z80.memory.writeByte(z80.HL(), bytetemp)
	}
	/* RLC A */
	opcodesMap[shift0xcb(0x07)] = func(z80 *Z80) {
		z80.a = z80.rlc(z80.a)
	}
	/* RRC B */
	opcodesMap[shift0xcb(0x08)] = func(z80 *Z80) {
		z80.b = z80.rrc(z80.b)
	}
	/* RRC C */
	opcodesMap[shift0xcb(0x09)] = func(z80 *Z80) {
		z80.c = z80.rrc(z80.c)
	}
	/* RRC D */
	opcodesMap[shift0xcb(0x0a)] = func(z80 *Z80) {
		z80.d = z80.rrc(z80.d)
	}
	/* RRC E */
	opcodesMap[shift0xcb(0x0b)] = func(z80 *Z80) {
		z80.e = z80.rrc(z80.e)
	}
	/* RRC H */
	opcodesMap[shift0xcb(0x0c)] = func(z80 *Z80) {
		z80.h = z80.rrc(z80.h)
	}
	/* RRC L */
	opcodesMap[shift0xcb(0x0d)] = func(z80 *Z80) {
		z80.l = z80.rrc(z80.l)
	}
	/* RRC (HL) */
	opcodesMap[shift0xcb(0x0e)] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.HL())
		z80.memory.contendReadNoMreq(z80.HL(), 1)
		bytetemp = z80.rrc(bytetemp)
		z80.memory.writeByte(z80.HL(), bytetemp)
	}
	/* RRC A */
	opcodesMap[shift0xcb(0x0f)] = func(z80 *Z80) {
		z80.a = z80.rrc(z80.a)
	}
	/* RL B */
	opcodesMap[shift0xcb(0x10)] = func(z80 *Z80) {
		z80.b = z80.rl(z80.b)
	}
	/* RL C */
	opcodesMap[shift0xcb(0x11)] = func(z80 *Z80) {
		z80.c = z80.rl(z80.c)
	}
	/* RL D */
	opcodesMap[shift0xcb(0x12)] = func(z80 *Z80) {
		z80.d = z80.rl(z80.d)
	}
	/* RL E */
	opcodesMap[shift0xcb(0x13)] = func(z80 *Z80) {
		z80.e = z80.rl(z80.e)
	}
	/* RL H */
	opcodesMap[shift0xcb(0x14)] = func(z80 *Z80) {
		z80.h = z80.rl(z80.h)
	}
	/* RL L */
	opcodesMap[shift0xcb(0x15)] = func(z80 *Z80) {
		z80.l = z80.rl(z80.l)
	}
	/* RL (HL) */
	opcodesMap[shift0xcb(0x16)] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.HL())
		z80.memory.contendReadNoMreq(z80.HL(), 1)
		bytetemp = z80.rl(bytetemp)
		z80.memory.writeByte(z80.HL(), bytetemp)
	}
	/* RL A */
	opcodesMap[shift0xcb(0x17)] = func(z80 *Z80) {
		z80.a = z80.rl(z80.a)
	}
	/* RR B */
	opcodesMap[shift0xcb(0x18)] = func(z80 *Z80) {
		z80.b = z80.rr(z80.b)
	}
	/* RR C */
	opcodesMap[shift0xcb(0x19)] = func(z80 *Z80) {
		z80.c = z80.rr(z80.c)
	}
	/* RR D */
	opcodesMap[shift0xcb(0x1a)] = func(z80 *Z80) {
		z80.d = z80.rr(z80.d)
	}
	/* RR E */
	opcodesMap[shift0xcb(0x1b)] = func(z80 *Z80) {
		z80.e = z80.rr(z80.e)
	}
	/* RR H */
	opcodesMap[shift0xcb(0x1c)] = func(z80 *Z80) {
		z80.h = z80.rr(z80.h)
	}
	/* RR L */
	opcodesMap[shift0xcb(0x1d)] = func(z80 *Z80) {
		z80.l = z80.rr(z80.l)
	}
	/* RR (HL) */
	opcodesMap[shift0xcb(0x1e)] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.HL())
		z80.memory.contendReadNoMreq(z80.HL(), 1)
		bytetemp = z80.rr(bytetemp)
		z80.memory.writeByte(z80.HL(), bytetemp)
	}
	/* RR A */
	opcodesMap[shift0xcb(0x1f)] = func(z80 *Z80) {
		z80.a = z80.rr(z80.a)
	}
	/* SLA B */
	opcodesMap[shift0xcb(0x20)] = func(z80 *Z80) {
		z80.b = z80.sla(z80.b)
	}
	/* SLA C */
	opcodesMap[shift0xcb(0x21)] = func(z80 *Z80) {
		z80.c = z80.sla(z80.c)
	}
	/* SLA D */
	opcodesMap[shift0xcb(0x22)] = func(z80 *Z80) {
		z80.d = z80.sla(z80.d)
	}
	/* SLA E */
	opcodesMap[shift0xcb(0x23)] = func(z80 *Z80) {
		z80.e = z80.sla(z80.e)
	}
	/* SLA H */
	opcodesMap[shift0xcb(0x24)] = func(z80 *Z80) {
		z80.h = z80.sla(z80.h)
	}
	/* SLA L */
	opcodesMap[shift0xcb(0x25)] = func(z80 *Z80) {
		z80.l = z80.sla(z80.l)
	}
	/* SLA (HL) */
	opcodesMap[shift0xcb(0x26)] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.HL())
		z80.memory.contendReadNoMreq(z80.HL(), 1)
		bytetemp = z80.sla(bytetemp)
		z80.memory.writeByte(z80.HL(), bytetemp)
	}
	/* SLA A */
	opcodesMap[shift0xcb(0x27)] = func(z80 *Z80) {
		z80.a = z80.sla(z80.a)
	}
	/* SRA B */
	opcodesMap[shift0xcb(0x28)] = func(z80 *Z80) {
		z80.b = z80.sra(z80.b)
	}
	/* SRA C */
	opcodesMap[shift0xcb(0x29)] = func(z80 *Z80) {
		z80.c = z80.sra(z80.c)
	}
	/* SRA D */
	opcodesMap[shift0xcb(0x2a)] = func(z80 *Z80) {
		z80.d = z80.sra(z80.d)
	}
	/* SRA E */
	opcodesMap[shift0xcb(0x2b)] = func(z80 *Z80) {
		z80.e = z80.sra(z80.e)
	}
	/* SRA H */
	opcodesMap[shift0xcb(0x2c)] = func(z80 *Z80) {
		z80.h = z80.sra(z80.h)
	}
	/* SRA L */
	opcodesMap[shift0xcb(0x2d)] = func(z80 *Z80) {
		z80.l = z80.sra(z80.l)
	}
	/* SRA (HL) */
	opcodesMap[shift0xcb(0x2e)] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.HL())
		z80.memory.contendReadNoMreq(z80.HL(), 1)
		bytetemp = z80.sra(bytetemp)
		z80.memory.writeByte(z80.HL(), bytetemp)
	}
	/* SRA A */
	opcodesMap[shift0xcb(0x2f)] = func(z80 *Z80) {
		z80.a = z80.sra(z80.a)
	}
	/* SLL B */
	opcodesMap[shift0xcb(0x30)] = func(z80 *Z80) {
		z80.b = z80.sll(z80.b)
	}
	/* SLL C */
	opcodesMap[shift0xcb(0x31)] = func(z80 *Z80) {
		z80.c = z80.sll(z80.c)
	}
	/* SLL D */
	opcodesMap[shift0xcb(0x32)] = func(z80 *Z80) {
		z80.d = z80.sll(z80.d)
	}
	/* SLL E */
	opcodesMap[shift0xcb(0x33)] = func(z80 *Z80) {
		z80.e = z80.sll(z80.e)
	}
	/* SLL H */
	opcodesMap[shift0xcb(0x34)] = func(z80 *Z80) {
		z80.h = z80.sll(z80.h)
	}
	/* SLL L */
	opcodesMap[shift0xcb(0x35)] = func(z80 *Z80) {
		z80.l = z80.sll(z80.l)
	}
	/* SLL (HL) */
	opcodesMap[shift0xcb(0x36)] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.HL())
		z80.memory.contendReadNoMreq(z80.HL(), 1)
		bytetemp = z80.sll(bytetemp)
		z80.memory.writeByte(z80.HL(), bytetemp)
	}
	/* SLL A */
	opcodesMap[shift0xcb(0x37)] = func(z80 *Z80) {
		z80.a = z80.sll(z80.a)
	}
	/* SRL B */
	opcodesMap[shift0xcb(0x38)] = func(z80 *Z80) {
		z80.b = z80.srl(z80.b)
	}
	/* SRL C */
	opcodesMap[shift0xcb(0x39)] = func(z80 *Z80) {
		z80.c = z80.srl(z80.c)
	}
	/* SRL D */
	opcodesMap[shift0xcb(0x3a)] = func(z80 *Z80) {
		z80.d = z80.srl(z80.d)
	}
	/* SRL E */
	opcodesMap[shift0xcb(0x3b)] = func(z80 *Z80) {
		z80.e = z80.srl(z80.e)
	}
	/* SRL H */
	opcodesMap[shift0xcb(0x3c)] = func(z80 *Z80) {
		z80.h = z80.srl(z80.h)
	}
	/* SRL L */
	opcodesMap[shift0xcb(0x3d)] = func(z80 *Z80) {
		z80.l = z80.srl(z80.l)
	}
	/* SRL (HL) */
	opcodesMap[shift0xcb(0x3e)] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.HL())
		z80.memory.contendReadNoMreq(z80.HL(), 1)
		bytetemp = z80.srl(bytetemp)
		z80.memory.writeByte(z80.HL(), bytetemp)
	}
	/* SRL A */
	opcodesMap[shift0xcb(0x3f)] = func(z80 *Z80) {
		z80.a = z80.srl(z80.a)
	}
	/* BIT 0,B */
	opcodesMap[shift0xcb(0x40)] = func(z80 *Z80) {
		z80.bit(0, z80.b)
	}
	/* BIT 0,C */
	opcodesMap[shift0xcb(0x41)] = func(z80 *Z80) {
		z80.bit(0, z80.c)
	}
	/* BIT 0,D */
	opcodesMap[shift0xcb(0x42)] = func(z80 *Z80) {
		z80.bit(0, z80.d)
	}
	/* BIT 0,E */
	opcodesMap[shift0xcb(0x43)] = func(z80 *Z80) {
		z80.bit(0, z80.e)
	}
	/* BIT 0,H */
	opcodesMap[shift0xcb(0x44)] = func(z80 *Z80) {
		z80.bit(0, z80.h)
	}
	/* BIT 0,L */
	opcodesMap[shift0xcb(0x45)] = func(z80 *Z80) {
		z80.bit(0, z80.l)
	}
	/* BIT 0,(HL) */
	opcodesMap[shift0xcb(0x46)] = func(z80 *Z80) {
		bytetemp := z80.memory.readByte(z80.HL())
		z80.memory.contendReadNoMreq(z80.HL(), 1)
		z80.bit(0, bytetemp)
	}
	/* BIT 0,A */
	opcodesMap[shift0xcb(0x47)] = func(z80 *Z80) {
		z80.bit(0, z80.a)
	}
	/* BIT 1,B */
	opcodesMap[shift0xcb(0x48)] = func(z80 *Z80) {
		z80.bit(1, z80.b)
	}
	/* BIT 1,C */
	opcodesMap[shift0xcb(0x49)] = func(z80 *Z80) {
		z80.bit(1, z80.c)
	}
	/* BIT 1,D */
	opcodesMap[shift0xcb(0x4a)] = func(z80 *Z80) {
		z80.bit(1, z80.d)
	}
	/* BIT 1,E */
	opcodesMap[shift0xcb(0x4b)] = func(z80 *Z80) {
		z80.bit(1, z80.e)
	}
	/* BIT 1,H */
	opcodesMap[shift0xcb(0x4c)] = func(z80 *Z80) {
		z80.bit(1, z80.h)
	}
	/* BIT 1,L */
	opcodesMap[shift0xcb(0x4d)] = func(z80 *Z80) {
		z80.bit(1, z80.l)
	}
	/* BIT 1,(HL) */
	opcodesMap[shift0xcb(0x4e)] = func(z80 *Z80) {
		bytetemp := z80.memory.readByte(z80.HL())
		z80.memory.contendReadNoMreq(z80.HL(), 1)
		z80.bit(1, bytetemp)
	}
	/* BIT 1,A */
	opcodesMap[shift0xcb(0x4f)] = func(z80 *Z80) {
		z80.bit(1, z80.a)
	}
	/* BIT 2,B */
	opcodesMap[shift0xcb(0x50)] = func(z80 *Z80) {
		z80.bit(2, z80.b)
	}
	/* BIT 2,C */
	opcodesMap[shift0xcb(0x51)] = func(z80 *Z80) {
		z80.bit(2, z80.c)
	}
	/* BIT 2,D */
	opcodesMap[shift0xcb(0x52)] = func(z80 *Z80) {
		z80.bit(2, z80.d)
	}
	/* BIT 2,E */
	opcodesMap[shift0xcb(0x53)] = func(z80 *Z80) {
		z80.bit(2, z80.e)
	}
	/* BIT 2,H */
	opcodesMap[shift0xcb(0x54)] = func(z80 *Z80) {
		z80.bit(2, z80.h)
	}
	/* BIT 2,L */
	opcodesMap[shift0xcb(0x55)] = func(z80 *Z80) {
		z80.bit(2, z80.l)
	}
	/* BIT 2,(HL) */
	opcodesMap[shift0xcb(0x56)] = func(z80 *Z80) {
		bytetemp := z80.memory.readByte(z80.HL())
		z80.memory.contendReadNoMreq(z80.HL(), 1)
		z80.bit(2, bytetemp)
	}
	/* BIT 2,A */
	opcodesMap[shift0xcb(0x57)] = func(z80 *Z80) {
		z80.bit(2, z80.a)
	}
	/* BIT 3,B */
	opcodesMap[shift0xcb(0x58)] = func(z80 *Z80) {
		z80.bit(3, z80.b)
	}
	/* BIT 3,C */
	opcodesMap[shift0xcb(0x59)] = func(z80 *Z80) {
		z80.bit(3, z80.c)
	}
	/* BIT 3,D */
	opcodesMap[shift0xcb(0x5a)] = func(z80 *Z80) {
		z80.bit(3, z80.d)
	}
	/* BIT 3,E */
	opcodesMap[shift0xcb(0x5b)] = func(z80 *Z80) {
		z80.bit(3, z80.e)
	}
	/* BIT 3,H */
	opcodesMap[shift0xcb(0x5c)] = func(z80 *Z80) {
		z80.bit(3, z80.h)
	}
	/* BIT 3,L */
	opcodesMap[shift0xcb(0x5d)] = func(z80 *Z80) {
		z80.bit(3, z80.l)
	}
	/* BIT 3,(HL) */
	opcodesMap[shift0xcb(0x5e)] = func(z80 *Z80) {
		bytetemp := z80.memory.readByte(z80.HL())
		z80.memory.contendReadNoMreq(z80.HL(), 1)
		z80.bit(3, bytetemp)
	}
	/* BIT 3,A */
	opcodesMap[shift0xcb(0x5f)] = func(z80 *Z80) {
		z80.bit(3, z80.a)
	}
	/* BIT 4,B */
	opcodesMap[shift0xcb(0x60)] = func(z80 *Z80) {
		z80.bit(4, z80.b)
	}
	/* BIT 4,C */
	opcodesMap[shift0xcb(0x61)] = func(z80 *Z80) {
		z80.bit(4, z80.c)
	}
	/* BIT 4,D */
	opcodesMap[shift0xcb(0x62)] = func(z80 *Z80) {
		z80.bit(4, z80.d)
	}
	/* BIT 4,E */
	opcodesMap[shift0xcb(0x63)] = func(z80 *Z80) {
		z80.bit(4, z80.e)
	}
	/* BIT 4,H */
	opcodesMap[shift0xcb(0x64)] = func(z80 *Z80) {
		z80.bit(4, z80.h)
	}
	/* BIT 4,L */
	opcodesMap[shift0xcb(0x65)] = func(z80 *Z80) {
		z80.bit(4, z80.l)
	}
	/* BIT 4,(HL) */
	opcodesMap[shift0xcb(0x66)] = func(z80 *Z80) {
		bytetemp := z80.memory.readByte(z80.HL())
		z80.memory.contendReadNoMreq(z80.HL(), 1)
		z80.bit(4, bytetemp)
	}
	/* BIT 4,A */
	opcodesMap[shift0xcb(0x67)] = func(z80 *Z80) {
		z80.bit(4, z80.a)
	}
	/* BIT 5,B */
	opcodesMap[shift0xcb(0x68)] = func(z80 *Z80) {
		z80.bit(5, z80.b)
	}
	/* BIT 5,C */
	opcodesMap[shift0xcb(0x69)] = func(z80 *Z80) {
		z80.bit(5, z80.c)
	}
	/* BIT 5,D */
	opcodesMap[shift0xcb(0x6a)] = func(z80 *Z80) {
		z80.bit(5, z80.d)
	}
	/* BIT 5,E */
	opcodesMap[shift0xcb(0x6b)] = func(z80 *Z80) {
		z80.bit(5, z80.e)
	}
	/* BIT 5,H */
	opcodesMap[shift0xcb(0x6c)] = func(z80 *Z80) {
		z80.bit(5, z80.h)
	}
	/* BIT 5,L */
	opcodesMap[shift0xcb(0x6d)] = func(z80 *Z80) {
		z80.bit(5, z80.l)
	}
	/* BIT 5,(HL) */
	opcodesMap[shift0xcb(0x6e)] = func(z80 *Z80) {
		bytetemp := z80.memory.readByte(z80.HL())
		z80.memory.contendReadNoMreq(z80.HL(), 1)
		z80.bit(5, bytetemp)
	}
	/* BIT 5,A */
	opcodesMap[shift0xcb(0x6f)] = func(z80 *Z80) {
		z80.bit(5, z80.a)
	}
	/* BIT 6,B */
	opcodesMap[shift0xcb(0x70)] = func(z80 *Z80) {
		z80.bit(6, z80.b)
	}
	/* BIT 6,C */
	opcodesMap[shift0xcb(0x71)] = func(z80 *Z80) {
		z80.bit(6, z80.c)
	}
	/* BIT 6,D */
	opcodesMap[shift0xcb(0x72)] = func(z80 *Z80) {
		z80.bit(6, z80.d)
	}
	/* BIT 6,E */
	opcodesMap[shift0xcb(0x73)] = func(z80 *Z80) {
		z80.bit(6, z80.e)
	}
	/* BIT 6,H */
	opcodesMap[shift0xcb(0x74)] = func(z80 *Z80) {
		z80.bit(6, z80.h)
	}
	/* BIT 6,L */
	opcodesMap[shift0xcb(0x75)] = func(z80 *Z80) {
		z80.bit(6, z80.l)
	}
	/* BIT 6,(HL) */
	opcodesMap[shift0xcb(0x76)] = func(z80 *Z80) {
		bytetemp := z80.memory.readByte(z80.HL())
		z80.memory.contendReadNoMreq(z80.HL(), 1)
		z80.bit(6, bytetemp)
	}
	/* BIT 6,A */
	opcodesMap[shift0xcb(0x77)] = func(z80 *Z80) {
		z80.bit(6, z80.a)
	}
	/* BIT 7,B */
	opcodesMap[shift0xcb(0x78)] = func(z80 *Z80) {
		z80.bit(7, z80.b)
	}
	/* BIT 7,C */
	opcodesMap[shift0xcb(0x79)] = func(z80 *Z80) {
		z80.bit(7, z80.c)
	}
	/* BIT 7,D */
	opcodesMap[shift0xcb(0x7a)] = func(z80 *Z80) {
		z80.bit(7, z80.d)
	}
	/* BIT 7,E */
	opcodesMap[shift0xcb(0x7b)] = func(z80 *Z80) {
		z80.bit(7, z80.e)
	}
	/* BIT 7,H */
	opcodesMap[shift0xcb(0x7c)] = func(z80 *Z80) {
		z80.bit(7, z80.h)
	}
	/* BIT 7,L */
	opcodesMap[shift0xcb(0x7d)] = func(z80 *Z80) {
		z80.bit(7, z80.l)
	}
	/* BIT 7,(HL) */
	opcodesMap[shift0xcb(0x7e)] = func(z80 *Z80) {
		bytetemp := z80.memory.readByte(z80.HL())
		z80.memory.contendReadNoMreq(z80.HL(), 1)
		z80.bit(7, bytetemp)
	}
	/* BIT 7,A */
	opcodesMap[shift0xcb(0x7f)] = func(z80 *Z80) {
		z80.bit(7, z80.a)
	}
	/* RES 0,B */
	opcodesMap[shift0xcb(0x80)] = func(z80 *Z80) {
		z80.b &= 0xfe
	}
	/* RES 0,C */
	opcodesMap[shift0xcb(0x81)] = func(z80 *Z80) {
		z80.c &= 0xfe
	}
	/* RES 0,D */
	opcodesMap[shift0xcb(0x82)] = func(z80 *Z80) {
		z80.d &= 0xfe
	}
	/* RES 0,E */
	opcodesMap[shift0xcb(0x83)] = func(z80 *Z80) {
		z80.e &= 0xfe
	}
	/* RES 0,H */
	opcodesMap[shift0xcb(0x84)] = func(z80 *Z80) {
		z80.h &= 0xfe
	}
	/* RES 0,L */
	opcodesMap[shift0xcb(0x85)] = func(z80 *Z80) {
		z80.l &= 0xfe
	}
	/* RES 0,(HL) */
	opcodesMap[shift0xcb(0x86)] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.HL())
		z80.memory.contendReadNoMreq(z80.HL(), 1)
		z80.memory.writeByte(z80.HL(), bytetemp&0xfe)
	}
	/* RES 0,A */
	opcodesMap[shift0xcb(0x87)] = func(z80 *Z80) {
		z80.a &= 0xfe
	}
	/* RES 1,B */
	opcodesMap[shift0xcb(0x88)] = func(z80 *Z80) {
		z80.b &= 0xfd
	}
	/* RES 1,C */
	opcodesMap[shift0xcb(0x89)] = func(z80 *Z80) {
		z80.c &= 0xfd
	}
	/* RES 1,D */
	opcodesMap[shift0xcb(0x8a)] = func(z80 *Z80) {
		z80.d &= 0xfd
	}
	/* RES 1,E */
	opcodesMap[shift0xcb(0x8b)] = func(z80 *Z80) {
		z80.e &= 0xfd
	}
	/* RES 1,H */
	opcodesMap[shift0xcb(0x8c)] = func(z80 *Z80) {
		z80.h &= 0xfd
	}
	/* RES 1,L */
	opcodesMap[shift0xcb(0x8d)] = func(z80 *Z80) {
		z80.l &= 0xfd
	}
	/* RES 1,(HL) */
	opcodesMap[shift0xcb(0x8e)] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.HL())
		z80.memory.contendReadNoMreq(z80.HL(), 1)
		z80.memory.writeByte(z80.HL(), bytetemp&0xfd)
	}
	/* RES 1,A */
	opcodesMap[shift0xcb(0x8f)] = func(z80 *Z80) {
		z80.a &= 0xfd
	}
	/* RES 2,B */
	opcodesMap[shift0xcb(0x90)] = func(z80 *Z80) {
		z80.b &= 0xfb
	}
	/* RES 2,C */
	opcodesMap[shift0xcb(0x91)] = func(z80 *Z80) {
		z80.c &= 0xfb
	}
	/* RES 2,D */
	opcodesMap[shift0xcb(0x92)] = func(z80 *Z80) {
		z80.d &= 0xfb
	}
	/* RES 2,E */
	opcodesMap[shift0xcb(0x93)] = func(z80 *Z80) {
		z80.e &= 0xfb
	}
	/* RES 2,H */
	opcodesMap[shift0xcb(0x94)] = func(z80 *Z80) {
		z80.h &= 0xfb
	}
	/* RES 2,L */
	opcodesMap[shift0xcb(0x95)] = func(z80 *Z80) {
		z80.l &= 0xfb
	}
	/* RES 2,(HL) */
	opcodesMap[shift0xcb(0x96)] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.HL())
		z80.memory.contendReadNoMreq(z80.HL(), 1)
		z80.memory.writeByte(z80.HL(), bytetemp&0xfb)
	}
	/* RES 2,A */
	opcodesMap[shift0xcb(0x97)] = func(z80 *Z80) {
		z80.a &= 0xfb
	}
	/* RES 3,B */
	opcodesMap[shift0xcb(0x98)] = func(z80 *Z80) {
		z80.b &= 0xf7
	}
	/* RES 3,C */
	opcodesMap[shift0xcb(0x99)] = func(z80 *Z80) {
		z80.c &= 0xf7
	}
	/* RES 3,D */
	opcodesMap[shift0xcb(0x9a)] = func(z80 *Z80) {
		z80.d &= 0xf7
	}
	/* RES 3,E */
	opcodesMap[shift0xcb(0x9b)] = func(z80 *Z80) {
		z80.e &= 0xf7
	}
	/* RES 3,H */
	opcodesMap[shift0xcb(0x9c)] = func(z80 *Z80) {
		z80.h &= 0xf7
	}
	/* RES 3,L */
	opcodesMap[shift0xcb(0x9d)] = func(z80 *Z80) {
		z80.l &= 0xf7
	}
	/* RES 3,(HL) */
	opcodesMap[shift0xcb(0x9e)] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.HL())
		z80.memory.contendReadNoMreq(z80.HL(), 1)
		z80.memory.writeByte(z80.HL(), bytetemp&0xf7)
	}
	/* RES 3,A */
	opcodesMap[shift0xcb(0x9f)] = func(z80 *Z80) {
		z80.a &= 0xf7
	}
	/* RES 4,B */
	opcodesMap[shift0xcb(0xa0)] = func(z80 *Z80) {
		z80.b &= 0xef
	}
	/* RES 4,C */
	opcodesMap[shift0xcb(0xa1)] = func(z80 *Z80) {
		z80.c &= 0xef
	}
	/* RES 4,D */
	opcodesMap[shift0xcb(0xa2)] = func(z80 *Z80) {
		z80.d &= 0xef
	}
	/* RES 4,E */
	opcodesMap[shift0xcb(0xa3)] = func(z80 *Z80) {
		z80.e &= 0xef
	}
	/* RES 4,H */
	opcodesMap[shift0xcb(0xa4)] = func(z80 *Z80) {
		z80.h &= 0xef
	}
	/* RES 4,L */
	opcodesMap[shift0xcb(0xa5)] = func(z80 *Z80) {
		z80.l &= 0xef
	}
	/* RES 4,(HL) */
	opcodesMap[shift0xcb(0xa6)] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.HL())
		z80.memory.contendReadNoMreq(z80.HL(), 1)
		z80.memory.writeByte(z80.HL(), bytetemp&0xef)
	}
	/* RES 4,A */
	opcodesMap[shift0xcb(0xa7)] = func(z80 *Z80) {
		z80.a &= 0xef
	}
	/* RES 5,B */
	opcodesMap[shift0xcb(0xa8)] = func(z80 *Z80) {
		z80.b &= 0xdf
	}
	/* RES 5,C */
	opcodesMap[shift0xcb(0xa9)] = func(z80 *Z80) {
		z80.c &= 0xdf
	}
	/* RES 5,D */
	opcodesMap[shift0xcb(0xaa)] = func(z80 *Z80) {
		z80.d &= 0xdf
	}
	/* RES 5,E */
	opcodesMap[shift0xcb(0xab)] = func(z80 *Z80) {
		z80.e &= 0xdf
	}
	/* RES 5,H */
	opcodesMap[shift0xcb(0xac)] = func(z80 *Z80) {
		z80.h &= 0xdf
	}
	/* RES 5,L */
	opcodesMap[shift0xcb(0xad)] = func(z80 *Z80) {
		z80.l &= 0xdf
	}
	/* RES 5,(HL) */
	opcodesMap[shift0xcb(0xae)] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.HL())
		z80.memory.contendReadNoMreq(z80.HL(), 1)
		z80.memory.writeByte(z80.HL(), bytetemp&0xdf)
	}
	/* RES 5,A */
	opcodesMap[shift0xcb(0xaf)] = func(z80 *Z80) {
		z80.a &= 0xdf
	}
	/* RES 6,B */
	opcodesMap[shift0xcb(0xb0)] = func(z80 *Z80) {
		z80.b &= 0xbf
	}
	/* RES 6,C */
	opcodesMap[shift0xcb(0xb1)] = func(z80 *Z80) {
		z80.c &= 0xbf
	}
	/* RES 6,D */
	opcodesMap[shift0xcb(0xb2)] = func(z80 *Z80) {
		z80.d &= 0xbf
	}
	/* RES 6,E */
	opcodesMap[shift0xcb(0xb3)] = func(z80 *Z80) {
		z80.e &= 0xbf
	}
	/* RES 6,H */
	opcodesMap[shift0xcb(0xb4)] = func(z80 *Z80) {
		z80.h &= 0xbf
	}
	/* RES 6,L */
	opcodesMap[shift0xcb(0xb5)] = func(z80 *Z80) {
		z80.l &= 0xbf
	}
	/* RES 6,(HL) */
	opcodesMap[shift0xcb(0xb6)] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.HL())
		z80.memory.contendReadNoMreq(z80.HL(), 1)
		z80.memory.writeByte(z80.HL(), bytetemp&0xbf)
	}
	/* RES 6,A */
	opcodesMap[shift0xcb(0xb7)] = func(z80 *Z80) {
		z80.a &= 0xbf
	}
	/* RES 7,B */
	opcodesMap[shift0xcb(0xb8)] = func(z80 *Z80) {
		z80.b &= 0x7f
	}
	/* RES 7,C */
	opcodesMap[shift0xcb(0xb9)] = func(z80 *Z80) {
		z80.c &= 0x7f
	}
	/* RES 7,D */
	opcodesMap[shift0xcb(0xba)] = func(z80 *Z80) {
		z80.d &= 0x7f
	}
	/* RES 7,E */
	opcodesMap[shift0xcb(0xbb)] = func(z80 *Z80) {
		z80.e &= 0x7f
	}
	/* RES 7,H */
	opcodesMap[shift0xcb(0xbc)] = func(z80 *Z80) {
		z80.h &= 0x7f
	}
	/* RES 7,L */
	opcodesMap[shift0xcb(0xbd)] = func(z80 *Z80) {
		z80.l &= 0x7f
	}
	/* RES 7,(HL) */
	opcodesMap[shift0xcb(0xbe)] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.HL())
		z80.memory.contendReadNoMreq(z80.HL(), 1)
		z80.memory.writeByte(z80.HL(), bytetemp&0x7f)
	}
	/* RES 7,A */
	opcodesMap[shift0xcb(0xbf)] = func(z80 *Z80) {
		z80.a &= 0x7f
	}
	/* SET 0,B */
	opcodesMap[shift0xcb(0xc0)] = func(z80 *Z80) {
		z80.b |= 0x01
	}
	/* SET 0,C */
	opcodesMap[shift0xcb(0xc1)] = func(z80 *Z80) {
		z80.c |= 0x01
	}
	/* SET 0,D */
	opcodesMap[shift0xcb(0xc2)] = func(z80 *Z80) {
		z80.d |= 0x01
	}
	/* SET 0,E */
	opcodesMap[shift0xcb(0xc3)] = func(z80 *Z80) {
		z80.e |= 0x01
	}
	/* SET 0,H */
	opcodesMap[shift0xcb(0xc4)] = func(z80 *Z80) {
		z80.h |= 0x01
	}
	/* SET 0,L */
	opcodesMap[shift0xcb(0xc5)] = func(z80 *Z80) {
		z80.l |= 0x01
	}
	/* SET 0,(HL) */
	opcodesMap[shift0xcb(0xc6)] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.HL())
		z80.memory.contendReadNoMreq(z80.HL(), 1)
		z80.memory.writeByte(z80.HL(), bytetemp|0x01)
	}
	/* SET 0,A */
	opcodesMap[shift0xcb(0xc7)] = func(z80 *Z80) {
		z80.a |= 0x01
	}
	/* SET 1,B */
	opcodesMap[shift0xcb(0xc8)] = func(z80 *Z80) {
		z80.b |= 0x02
	}
	/* SET 1,C */
	opcodesMap[shift0xcb(0xc9)] = func(z80 *Z80) {
		z80.c |= 0x02
	}
	/* SET 1,D */
	opcodesMap[shift0xcb(0xca)] = func(z80 *Z80) {
		z80.d |= 0x02
	}
	/* SET 1,E */
	opcodesMap[shift0xcb(0xcb)] = func(z80 *Z80) {
		z80.e |= 0x02
	}
	/* SET 1,H */
	opcodesMap[shift0xcb(0xcc)] = func(z80 *Z80) {
		z80.h |= 0x02
	}
	/* SET 1,L */
	opcodesMap[shift0xcb(0xcd)] = func(z80 *Z80) {
		z80.l |= 0x02
	}
	/* SET 1,(HL) */
	opcodesMap[shift0xcb(0xce)] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.HL())
		z80.memory.contendReadNoMreq(z80.HL(), 1)
		z80.memory.writeByte(z80.HL(), bytetemp|0x02)
	}
	/* SET 1,A */
	opcodesMap[shift0xcb(0xcf)] = func(z80 *Z80) {
		z80.a |= 0x02
	}
	/* SET 2,B */
	opcodesMap[shift0xcb(0xd0)] = func(z80 *Z80) {
		z80.b |= 0x04
	}
	/* SET 2,C */
	opcodesMap[shift0xcb(0xd1)] = func(z80 *Z80) {
		z80.c |= 0x04
	}
	/* SET 2,D */
	opcodesMap[shift0xcb(0xd2)] = func(z80 *Z80) {
		z80.d |= 0x04
	}
	/* SET 2,E */
	opcodesMap[shift0xcb(0xd3)] = func(z80 *Z80) {
		z80.e |= 0x04
	}
	/* SET 2,H */
	opcodesMap[shift0xcb(0xd4)] = func(z80 *Z80) {
		z80.h |= 0x04
	}
	/* SET 2,L */
	opcodesMap[shift0xcb(0xd5)] = func(z80 *Z80) {
		z80.l |= 0x04
	}
	/* SET 2,(HL) */
	opcodesMap[shift0xcb(0xd6)] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.HL())
		z80.memory.contendReadNoMreq(z80.HL(), 1)
		z80.memory.writeByte(z80.HL(), bytetemp|0x04)
	}
	/* SET 2,A */
	opcodesMap[shift0xcb(0xd7)] = func(z80 *Z80) {
		z80.a |= 0x04
	}
	/* SET 3,B */
	opcodesMap[shift0xcb(0xd8)] = func(z80 *Z80) {
		z80.b |= 0x08
	}
	/* SET 3,C */
	opcodesMap[shift0xcb(0xd9)] = func(z80 *Z80) {
		z80.c |= 0x08
	}
	/* SET 3,D */
	opcodesMap[shift0xcb(0xda)] = func(z80 *Z80) {
		z80.d |= 0x08
	}
	/* SET 3,E */
	opcodesMap[shift0xcb(0xdb)] = func(z80 *Z80) {
		z80.e |= 0x08
	}
	/* SET 3,H */
	opcodesMap[shift0xcb(0xdc)] = func(z80 *Z80) {
		z80.h |= 0x08
	}
	/* SET 3,L */
	opcodesMap[shift0xcb(0xdd)] = func(z80 *Z80) {
		z80.l |= 0x08
	}
	/* SET 3,(HL) */
	opcodesMap[shift0xcb(0xde)] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.HL())
		z80.memory.contendReadNoMreq(z80.HL(), 1)
		z80.memory.writeByte(z80.HL(), bytetemp|0x08)
	}
	/* SET 3,A */
	opcodesMap[shift0xcb(0xdf)] = func(z80 *Z80) {
		z80.a |= 0x08
	}
	/* SET 4,B */
	opcodesMap[shift0xcb(0xe0)] = func(z80 *Z80) {
		z80.b |= 0x10
	}
	/* SET 4,C */
	opcodesMap[shift0xcb(0xe1)] = func(z80 *Z80) {
		z80.c |= 0x10
	}
	/* SET 4,D */
	opcodesMap[shift0xcb(0xe2)] = func(z80 *Z80) {
		z80.d |= 0x10
	}
	/* SET 4,E */
	opcodesMap[shift0xcb(0xe3)] = func(z80 *Z80) {
		z80.e |= 0x10
	}
	/* SET 4,H */
	opcodesMap[shift0xcb(0xe4)] = func(z80 *Z80) {
		z80.h |= 0x10
	}
	/* SET 4,L */
	opcodesMap[shift0xcb(0xe5)] = func(z80 *Z80) {
		z80.l |= 0x10
	}
	/* SET 4,(HL) */
	opcodesMap[shift0xcb(0xe6)] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.HL())
		z80.memory.contendReadNoMreq(z80.HL(), 1)
		z80.memory.writeByte(z80.HL(), bytetemp|0x10)
	}
	/* SET 4,A */
	opcodesMap[shift0xcb(0xe7)] = func(z80 *Z80) {
		z80.a |= 0x10
	}
	/* SET 5,B */
	opcodesMap[shift0xcb(0xe8)] = func(z80 *Z80) {
		z80.b |= 0x20
	}
	/* SET 5,C */
	opcodesMap[shift0xcb(0xe9)] = func(z80 *Z80) {
		z80.c |= 0x20
	}
	/* SET 5,D */
	opcodesMap[shift0xcb(0xea)] = func(z80 *Z80) {
		z80.d |= 0x20
	}
	/* SET 5,E */
	opcodesMap[shift0xcb(0xeb)] = func(z80 *Z80) {
		z80.e |= 0x20
	}
	/* SET 5,H */
	opcodesMap[shift0xcb(0xec)] = func(z80 *Z80) {
		z80.h |= 0x20
	}
	/* SET 5,L */
	opcodesMap[shift0xcb(0xed)] = func(z80 *Z80) {
		z80.l |= 0x20
	}
	/* SET 5,(HL) */
	opcodesMap[shift0xcb(0xee)] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.HL())
		z80.memory.contendReadNoMreq(z80.HL(), 1)
		z80.memory.writeByte(z80.HL(), bytetemp|0x20)
	}
	/* SET 5,A */
	opcodesMap[shift0xcb(0xef)] = func(z80 *Z80) {
		z80.a |= 0x20
	}
	/* SET 6,B */
	opcodesMap[shift0xcb(0xf0)] = func(z80 *Z80) {
		z80.b |= 0x40
	}
	/* SET 6,C */
	opcodesMap[shift0xcb(0xf1)] = func(z80 *Z80) {
		z80.c |= 0x40
	}
	/* SET 6,D */
	opcodesMap[shift0xcb(0xf2)] = func(z80 *Z80) {
		z80.d |= 0x40
	}
	/* SET 6,E */
	opcodesMap[shift0xcb(0xf3)] = func(z80 *Z80) {
		z80.e |= 0x40
	}
	/* SET 6,H */
	opcodesMap[shift0xcb(0xf4)] = func(z80 *Z80) {
		z80.h |= 0x40
	}
	/* SET 6,L */
	opcodesMap[shift0xcb(0xf5)] = func(z80 *Z80) {
		z80.l |= 0x40
	}
	/* SET 6,(HL) */
	opcodesMap[shift0xcb(0xf6)] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.HL())
		z80.memory.contendReadNoMreq(z80.HL(), 1)
		z80.memory.writeByte(z80.HL(), bytetemp|0x40)
	}
	/* SET 6,A */
	opcodesMap[shift0xcb(0xf7)] = func(z80 *Z80) {
		z80.a |= 0x40
	}
	/* SET 7,B */
	opcodesMap[shift0xcb(0xf8)] = func(z80 *Z80) {
		z80.b |= 0x80
	}
	/* SET 7,C */
	opcodesMap[shift0xcb(0xf9)] = func(z80 *Z80) {
		z80.c |= 0x80
	}
	/* SET 7,D */
	opcodesMap[shift0xcb(0xfa)] = func(z80 *Z80) {
		z80.d |= 0x80
	}
	/* SET 7,E */
	opcodesMap[shift0xcb(0xfb)] = func(z80 *Z80) {
		z80.e |= 0x80
	}
	/* SET 7,H */
	opcodesMap[shift0xcb(0xfc)] = func(z80 *Z80) {
		z80.h |= 0x80
	}
	/* SET 7,L */
	opcodesMap[shift0xcb(0xfd)] = func(z80 *Z80) {
		z80.l |= 0x80
	}
	/* SET 7,(HL) */
	opcodesMap[shift0xcb(0xfe)] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.HL())
		z80.memory.contendReadNoMreq(z80.HL(), 1)
		z80.memory.writeByte(z80.HL(), bytetemp|0x80)
	}
	/* SET 7,A */
	opcodesMap[shift0xcb(0xff)] = func(z80 *Z80) {
		z80.a |= 0x80
	}

	// END of 0xcb shifted opcodes

	// BEGIN of 0xed shifted opcodes

	/* IN B,(C) */
	opcodesMap[shift0xed(0x40)] = func(z80 *Z80) {
		z80.in(&z80.b, z80.BC())
	}
	/* OUT (C),B */
	opcodesMap[shift0xed(0x41)] = func(z80 *Z80) {
		z80.writePort(z80.BC(), z80.b)
	}
	/* SBC HL,BC */
	opcodesMap[shift0xed(0x42)] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 7)
		z80.sbc16(z80.BC())
	}
	/* LD (nnnn),BC */
	opcodesMap[shift0xed(0x43)] = func(z80 *Z80) {
		z80.ld16nnrr(z80.c, z80.b)
		// break
	}
	/* NEG */
	opcodesMap[shift0xed(0x7c)] = func(z80 *Z80) {
		bytetemp := z80.a
		z80.a = 0
		z80.sub(bytetemp)
	}
	// Fallthrough cases
	opcodesMap[shift0xed(0x44)] = opcodesMap[shift0xed(0x7c)]
	opcodesMap[shift0xed(0x4c)] = opcodesMap[shift0xed(0x7c)]
	opcodesMap[shift0xed(0x54)] = opcodesMap[shift0xed(0x7c)]
	opcodesMap[shift0xed(0x5c)] = opcodesMap[shift0xed(0x7c)]
	opcodesMap[shift0xed(0x64)] = opcodesMap[shift0xed(0x7c)]
	opcodesMap[shift0xed(0x6c)] = opcodesMap[shift0xed(0x7c)]
	opcodesMap[shift0xed(0x74)] = opcodesMap[shift0xed(0x7c)]
	/* RETN */
	opcodesMap[shift0xed(0x7d)] = func(z80 *Z80) {
		z80.iff1 = z80.iff2
		z80.ret()
	}
	// Fallthrough cases
	opcodesMap[shift0xed(0x45)] = opcodesMap[shift0xed(0x7d)]
	opcodesMap[shift0xed(0x4d)] = opcodesMap[shift0xed(0x7d)]
	opcodesMap[shift0xed(0x55)] = opcodesMap[shift0xed(0x7d)]
	opcodesMap[shift0xed(0x5d)] = opcodesMap[shift0xed(0x7d)]
	opcodesMap[shift0xed(0x65)] = opcodesMap[shift0xed(0x7d)]
	opcodesMap[shift0xed(0x6d)] = opcodesMap[shift0xed(0x7d)]
	opcodesMap[shift0xed(0x75)] = opcodesMap[shift0xed(0x7d)]
	/* IM 0 */
	opcodesMap[shift0xed(0x6e)] = func(z80 *Z80) {
		z80.im = 0
	}
	// Fallthrough cases
	opcodesMap[shift0xed(0x46)] = opcodesMap[shift0xed(0x6e)]
	opcodesMap[shift0xed(0x4e)] = opcodesMap[shift0xed(0x6e)]
	opcodesMap[shift0xed(0x66)] = opcodesMap[shift0xed(0x6e)]
	/* LD I,A */
	opcodesMap[shift0xed(0x47)] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq(z80.IR(), 1)
		z80.i = z80.a
	}
	/* IN C,(C) */
	opcodesMap[shift0xed(0x48)] = func(z80 *Z80) {
		z80.in(&z80.c, z80.BC())
	}
	/* OUT (C),C */
	opcodesMap[shift0xed(0x49)] = func(z80 *Z80) {
		z80.writePort(z80.BC(), z80.c)
	}
	/* ADC HL,BC */
	opcodesMap[shift0xed(0x4a)] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 7)
		z80.adc16(z80.BC())
	}
	/* LD BC,(nnnn) */
	opcodesMap[shift0xed(0x4b)] = func(z80 *Z80) {
		z80.ld16rrnn(&z80.c, &z80.b)
		// break
	}
	/* LD R,A */
	opcodesMap[shift0xed(0x4f)] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq(z80.IR(), 1)
		/* Keep the RZX instruction counter right */
		z80.rzxInstructionsOffset += (int(z80.r) - int(z80.a))
		z80.r, z80.r7 = uint16(z80.a), z80.a
	}
	/* IN D,(C) */
	opcodesMap[shift0xed(0x50)] = func(z80 *Z80) {
		z80.in(&z80.d, z80.BC())
	}
	/* OUT (C),D */
	opcodesMap[shift0xed(0x51)] = func(z80 *Z80) {
		z80.writePort(z80.BC(), z80.d)
	}
	/* SBC HL,DE */
	opcodesMap[shift0xed(0x52)] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 7)
		z80.sbc16(z80.DE())
	}
	/* LD (nnnn),DE */
	opcodesMap[shift0xed(0x53)] = func(z80 *Z80) {
		z80.ld16nnrr(z80.e, z80.d)
		// break
	}
	/* IM 1 */
	opcodesMap[shift0xed(0x76)] = func(z80 *Z80) {
		z80.im = 1
	}
	// Fallthrough cases
	opcodesMap[shift0xed(0x56)] = opcodesMap[shift0xed(0x76)]
	/* LD A,I */
	opcodesMap[shift0xed(0x57)] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq(z80.IR(), 1)
		z80.a = z80.i
		z80.f = (z80.f & FLAG_C) | sz53Table[z80.a] | ternOpB(z80.iff2 != 0, FLAG_V, 0)
	}
	/* IN E,(C) */
	opcodesMap[shift0xed(0x58)] = func(z80 *Z80) {
		z80.in(&z80.e, z80.BC())
	}
	/* OUT (C),E */
	opcodesMap[shift0xed(0x59)] = func(z80 *Z80) {
		z80.writePort(z80.BC(), z80.e)
	}
	/* ADC HL,DE */
	opcodesMap[shift0xed(0x5a)] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 7)
		z80.adc16(z80.DE())
	}
	/* LD DE,(nnnn) */
	opcodesMap[shift0xed(0x5b)] = func(z80 *Z80) {
		z80.ld16rrnn(&z80.e, &z80.d)
		// break
	}
	/* IM 2 */
	opcodesMap[shift0xed(0x7e)] = func(z80 *Z80) {
		z80.im = 2
	}
	// Fallthrough cases
	opcodesMap[shift0xed(0x5e)] = opcodesMap[shift0xed(0x7e)]
	/* LD A,R */
	opcodesMap[shift0xed(0x5f)] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq(z80.IR(), 1)
		z80.a = byte(z80.r&0x7f) | (z80.r7 & 0x80)
		z80.f = (z80.f & FLAG_C) | sz53Table[z80.a] | ternOpB(z80.iff2 != 0, FLAG_V, 0)
	}
	/* IN H,(C) */
	opcodesMap[shift0xed(0x60)] = func(z80 *Z80) {
		z80.in(&z80.h, z80.BC())
	}
	/* OUT (C),H */
	opcodesMap[shift0xed(0x61)] = func(z80 *Z80) {
		z80.writePort(z80.BC(), z80.h)
	}
	/* SBC HL,HL */
	opcodesMap[shift0xed(0x62)] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 7)
		z80.sbc16(z80.HL())
	}
	/* LD (nnnn),HL */
	opcodesMap[shift0xed(0x63)] = func(z80 *Z80) {
		z80.ld16nnrr(z80.l, z80.h)
		// break
	}
	/* RRD */
	opcodesMap[shift0xed(0x67)] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.HL())
		z80.memory.contendReadNoMreq_loop(z80.HL(), 1, 4)
		z80.memory.writeByte(z80.HL(), (z80.a<<4)|(bytetemp>>4))
		z80.a = (z80.a & 0xf0) | (bytetemp & 0x0f)
		z80.f = (z80.f & FLAG_C) | sz53pTable[z80.a]
	}
	/* IN L,(C) */
	opcodesMap[shift0xed(0x68)] = func(z80 *Z80) {
		z80.in(&z80.l, z80.BC())
	}
	/* OUT (C),L */
	opcodesMap[shift0xed(0x69)] = func(z80 *Z80) {
		z80.writePort(z80.BC(), z80.l)
	}
	/* ADC HL,HL */
	opcodesMap[shift0xed(0x6a)] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 7)
		z80.adc16(z80.HL())
	}
	/* LD HL,(nnnn) */
	opcodesMap[shift0xed(0x6b)] = func(z80 *Z80) {
		z80.ld16rrnn(&z80.l, &z80.h)
		// break
	}
	/* RLD */
	opcodesMap[shift0xed(0x6f)] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.HL())
		z80.memory.contendReadNoMreq_loop(z80.HL(), 1, 4)
		z80.memory.writeByte(z80.HL(), (bytetemp<<4)|(z80.a&0x0f))
		z80.a = (z80.a & 0xf0) | (bytetemp >> 4)
		z80.f = (z80.f & FLAG_C) | sz53pTable[z80.a]
	}
	/* IN F,(C) */
	opcodesMap[shift0xed(0x70)] = func(z80 *Z80) {
		var bytetemp byte
		z80.in(&bytetemp, z80.BC())
	}
	/* OUT (C),0 */
	opcodesMap[shift0xed(0x71)] = func(z80 *Z80) {
		z80.writePort(z80.BC(), 0)
	}
	/* SBC HL,SP */
	opcodesMap[shift0xed(0x72)] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 7)
		z80.sbc16(z80.SP())
	}
	/* LD (nnnn),SP */
	opcodesMap[shift0xed(0x73)] = func(z80 *Z80) {
		sph, spl := splitWord(z80.sp)
		z80.ld16nnrr(spl, sph)
		// break
	}
	/* IN A,(C) */
	opcodesMap[shift0xed(0x78)] = func(z80 *Z80) {
		z80.in(&z80.a, z80.BC())
	}
	/* OUT (C),A */
	opcodesMap[shift0xed(0x79)] = func(z80 *Z80) {
		z80.writePort(z80.BC(), z80.a)
	}
	/* ADC HL,SP */
	opcodesMap[shift0xed(0x7a)] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 7)
		z80.adc16(z80.SP())
	}
	/* LD SP,(nnnn) */
	opcodesMap[shift0xed(0x7b)] = func(z80 *Z80) {
		sph, spl := splitWord(z80.sp)
		z80.ld16rrnn(&spl, &sph)
		z80.sp = joinBytes(sph, spl)
		// break
	}
	/* LDI */
	opcodesMap[shift0xed(0xa0)] = func(z80 *Z80) {
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
	opcodesMap[shift0xed(0xa1)] = func(z80 *Z80) {
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
	opcodesMap[shift0xed(0xa2)] = func(z80 *Z80) {
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
	opcodesMap[shift0xed(0xa3)] = func(z80 *Z80) {
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
	opcodesMap[shift0xed(0xa8)] = func(z80 *Z80) {
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
	opcodesMap[shift0xed(0xa9)] = func(z80 *Z80) {
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
	opcodesMap[shift0xed(0xaa)] = func(z80 *Z80) {
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
	opcodesMap[shift0xed(0xab)] = func(z80 *Z80) {
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
	opcodesMap[shift0xed(0xb0)] = func(z80 *Z80) {
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
	opcodesMap[shift0xed(0xb1)] = func(z80 *Z80) {
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
	opcodesMap[shift0xed(0xb2)] = func(z80 *Z80) {
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
	opcodesMap[shift0xed(0xb3)] = func(z80 *Z80) {
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
	opcodesMap[shift0xed(0xb8)] = func(z80 *Z80) {
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
	opcodesMap[shift0xed(0xb9)] = func(z80 *Z80) {
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
	opcodesMap[shift0xed(0xba)] = func(z80 *Z80) {
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
	opcodesMap[shift0xed(0xbb)] = func(z80 *Z80) {
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
	opcodesMap[shift0xed(0xfb)] = func(z80 *Z80) {
		z80.sltTrap(int16(z80.HL()), z80.a)
	}

	// END of 0xed shifted opcodes

	// BEGIN of 0xdd shifted opcodes

	/* ADD ix,BC */
	opcodesMap[shift0xdd(0x09)] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 7)
		z80.add16(z80.ix, z80.BC())
	}
	/* ADD ix,DE */
	opcodesMap[shift0xdd(0x19)] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 7)
		z80.add16(z80.ix, z80.DE())
	}
	/* LD ix,nnnn */
	opcodesMap[shift0xdd(0x21)] = func(z80 *Z80) {
		b1 := z80.memory.readByte(z80.pc)
		z80.pc++
		b2 := z80.memory.readByte(z80.pc)
		z80.pc++
		z80.setIX(joinBytes(b2, b1))
	}
	/* LD (nnnn),ix */
	opcodesMap[shift0xdd(0x22)] = func(z80 *Z80) {
		z80.ld16nnrr(z80.ixl, z80.ixh)
		// break
	}
	/* INC ix */
	opcodesMap[shift0xdd(0x23)] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 2)
		z80.incIX()
	}
	/* INC z80.IXH() */
	opcodesMap[shift0xdd(0x24)] = func(z80 *Z80) {
		z80.incIXH()
	}
	/* DEC z80.IXH() */
	opcodesMap[shift0xdd(0x25)] = func(z80 *Z80) {
		z80.decIXH()
	}
	/* LD z80.IXH(),nn */
	opcodesMap[shift0xdd(0x26)] = func(z80 *Z80) {
		z80.ixh = z80.memory.readByte(z80.pc)
		z80.pc++
	}
	/* ADD ix,ix */
	opcodesMap[shift0xdd(0x29)] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 7)
		z80.add16(z80.ix, z80.IX())
	}
	/* LD ix,(nnnn) */
	opcodesMap[shift0xdd(0x2a)] = func(z80 *Z80) {
		z80.ld16rrnn(&z80.ixl, &z80.ixh)
		// break
	}
	/* DEC ix */
	opcodesMap[shift0xdd(0x2b)] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 2)
		z80.decIX()
	}
	/* INC z80.IXL() */
	opcodesMap[shift0xdd(0x2c)] = func(z80 *Z80) {
		z80.incIXL()
	}
	/* DEC z80.IXL() */
	opcodesMap[shift0xdd(0x2d)] = func(z80 *Z80) {
		z80.decIXL()
	}
	/* LD z80.IXL(),nn */
	opcodesMap[shift0xdd(0x2e)] = func(z80 *Z80) {
		z80.ixl = z80.memory.readByte(z80.pc)
		z80.pc++
	}
	/* INC (ix+dd) */
	opcodesMap[shift0xdd(0x34)] = func(z80 *Z80) {
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
	opcodesMap[shift0xdd(0x35)] = func(z80 *Z80) {
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
	opcodesMap[shift0xdd(0x36)] = func(z80 *Z80) {
		offset := z80.memory.readByte(z80.pc)
		z80.pc++
		value := z80.memory.readByte(z80.pc)
		z80.memory.contendReadNoMreq_loop(z80.pc, 1, 2)
		z80.pc++
		z80.memory.writeByte(z80.IX()+uint16(signExtend(offset)), value)
	}
	/* ADD ix,SP */
	opcodesMap[shift0xdd(0x39)] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 7)
		z80.add16(z80.ix, z80.SP())
	}
	/* LD B,z80.IXH() */
	opcodesMap[shift0xdd(0x44)] = func(z80 *Z80) {
		z80.b = z80.ixh
	}
	/* LD B,z80.IXL() */
	opcodesMap[shift0xdd(0x45)] = func(z80 *Z80) {
		z80.b = z80.ixl
	}
	/* LD B,(ix+dd) */
	opcodesMap[shift0xdd(0x46)] = func(z80 *Z80) {
		var offset byte = z80.memory.readByte(z80.pc)
		z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
		z80.pc++
		z80.b = z80.memory.readByte(z80.IX() + uint16(signExtend(offset)))
	}
	/* LD C,z80.IXH() */
	opcodesMap[shift0xdd(0x4c)] = func(z80 *Z80) {
		z80.c = z80.ixh
	}
	/* LD C,z80.IXL() */
	opcodesMap[shift0xdd(0x4d)] = func(z80 *Z80) {
		z80.c = z80.ixl
	}
	/* LD C,(ix+dd) */
	opcodesMap[shift0xdd(0x4e)] = func(z80 *Z80) {
		var offset byte = z80.memory.readByte(z80.pc)
		z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
		z80.pc++
		z80.c = z80.memory.readByte(z80.IX() + uint16(signExtend(offset)))
	}
	/* LD D,z80.IXH() */
	opcodesMap[shift0xdd(0x54)] = func(z80 *Z80) {
		z80.d = z80.ixh
	}
	/* LD D,z80.IXL() */
	opcodesMap[shift0xdd(0x55)] = func(z80 *Z80) {
		z80.d = z80.ixl
	}
	/* LD D,(ix+dd) */
	opcodesMap[shift0xdd(0x56)] = func(z80 *Z80) {
		var offset byte = z80.memory.readByte(z80.pc)
		z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
		z80.pc++
		z80.d = z80.memory.readByte(z80.IX() + uint16(signExtend(offset)))
	}
	/* LD E,z80.IXH() */
	opcodesMap[shift0xdd(0x5c)] = func(z80 *Z80) {
		z80.e = z80.ixh
	}
	/* LD E,z80.IXL() */
	opcodesMap[shift0xdd(0x5d)] = func(z80 *Z80) {
		z80.e = z80.ixl
	}
	/* LD E,(ix+dd) */
	opcodesMap[shift0xdd(0x5e)] = func(z80 *Z80) {
		var offset byte = z80.memory.readByte(z80.pc)
		z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
		z80.pc++
		z80.e = z80.memory.readByte(z80.IX() + uint16(signExtend(offset)))
	}
	/* LD z80.IXH(),B */
	opcodesMap[shift0xdd(0x60)] = func(z80 *Z80) {
		z80.ixh = z80.b
	}
	/* LD z80.IXH(),C */
	opcodesMap[shift0xdd(0x61)] = func(z80 *Z80) {
		z80.ixh = z80.c
	}
	/* LD z80.IXH(),D */
	opcodesMap[shift0xdd(0x62)] = func(z80 *Z80) {
		z80.ixh = z80.d
	}
	/* LD z80.IXH(),E */
	opcodesMap[shift0xdd(0x63)] = func(z80 *Z80) {
		z80.ixh = z80.e
	}
	/* LD z80.IXH(),z80.IXH() */
	opcodesMap[shift0xdd(0x64)] = func(z80 *Z80) {
	}
	/* LD z80.IXH(),z80.IXL() */
	opcodesMap[shift0xdd(0x65)] = func(z80 *Z80) {
		z80.ixh = z80.ixl
	}
	/* LD H,(ix+dd) */
	opcodesMap[shift0xdd(0x66)] = func(z80 *Z80) {
		var offset byte = z80.memory.readByte(z80.pc)
		z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
		z80.pc++
		z80.h = z80.memory.readByte(z80.IX() + uint16(signExtend(offset)))
	}
	/* LD z80.IXH(),A */
	opcodesMap[shift0xdd(0x67)] = func(z80 *Z80) {
		z80.ixh = z80.a
	}
	/* LD z80.IXL(),B */
	opcodesMap[shift0xdd(0x68)] = func(z80 *Z80) {
		z80.ixl = z80.b
	}
	/* LD z80.IXL(),C */
	opcodesMap[shift0xdd(0x69)] = func(z80 *Z80) {
		z80.ixl = z80.c
	}
	/* LD z80.IXL(),D */
	opcodesMap[shift0xdd(0x6a)] = func(z80 *Z80) {
		z80.ixl = z80.d
	}
	/* LD z80.IXL(),E */
	opcodesMap[shift0xdd(0x6b)] = func(z80 *Z80) {
		z80.ixl = z80.e
	}
	/* LD z80.IXL(),z80.IXH() */
	opcodesMap[shift0xdd(0x6c)] = func(z80 *Z80) {
		z80.ixl = z80.ixh
	}
	/* LD z80.IXL(),z80.IXL() */
	opcodesMap[shift0xdd(0x6d)] = func(z80 *Z80) {
	}
	/* LD L,(ix+dd) */
	opcodesMap[shift0xdd(0x6e)] = func(z80 *Z80) {
		var offset byte = z80.memory.readByte(z80.pc)
		z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
		z80.pc++
		z80.l = z80.memory.readByte(z80.IX() + uint16(signExtend(offset)))
	}
	/* LD z80.IXL(),A */
	opcodesMap[shift0xdd(0x6f)] = func(z80 *Z80) {
		z80.ixl = z80.a
	}
	/* LD (ix+dd),B */
	opcodesMap[shift0xdd(0x70)] = func(z80 *Z80) {
		offset := z80.memory.readByte(z80.pc)
		z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
		z80.pc++
		z80.memory.writeByte(z80.IX()+uint16(signExtend(offset)), z80.b)
	}
	/* LD (ix+dd),C */
	opcodesMap[shift0xdd(0x71)] = func(z80 *Z80) {
		offset := z80.memory.readByte(z80.pc)
		z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
		z80.pc++
		z80.memory.writeByte(z80.IX()+uint16(signExtend(offset)), z80.c)
	}
	/* LD (ix+dd),D */
	opcodesMap[shift0xdd(0x72)] = func(z80 *Z80) {
		offset := z80.memory.readByte(z80.pc)
		z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
		z80.pc++
		z80.memory.writeByte(z80.IX()+uint16(signExtend(offset)), z80.d)
	}
	/* LD (ix+dd),E */
	opcodesMap[shift0xdd(0x73)] = func(z80 *Z80) {
		offset := z80.memory.readByte(z80.pc)
		z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
		z80.pc++
		z80.memory.writeByte(z80.IX()+uint16(signExtend(offset)), z80.e)
	}
	/* LD (ix+dd),H */
	opcodesMap[shift0xdd(0x74)] = func(z80 *Z80) {
		offset := z80.memory.readByte(z80.pc)
		z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
		z80.pc++
		z80.memory.writeByte(z80.IX()+uint16(signExtend(offset)), z80.h)
	}
	/* LD (ix+dd),L */
	opcodesMap[shift0xdd(0x75)] = func(z80 *Z80) {
		offset := z80.memory.readByte(z80.pc)
		z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
		z80.pc++
		z80.memory.writeByte(z80.IX()+uint16(signExtend(offset)), z80.l)
	}
	/* LD (ix+dd),A */
	opcodesMap[shift0xdd(0x77)] = func(z80 *Z80) {
		offset := z80.memory.readByte(z80.pc)
		z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
		z80.pc++
		z80.memory.writeByte(z80.IX()+uint16(signExtend(offset)), z80.a)
	}
	/* LD A,z80.IXH() */
	opcodesMap[shift0xdd(0x7c)] = func(z80 *Z80) {
		z80.a = z80.ixh
	}
	/* LD A,z80.IXL() */
	opcodesMap[shift0xdd(0x7d)] = func(z80 *Z80) {
		z80.a = z80.ixl
	}
	/* LD A,(ix+dd) */
	opcodesMap[shift0xdd(0x7e)] = func(z80 *Z80) {
		var offset byte = z80.memory.readByte(z80.pc)
		z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
		z80.pc++
		z80.a = z80.memory.readByte(z80.IX() + uint16(signExtend(offset)))
	}
	/* ADD A,z80.IXH() */
	opcodesMap[shift0xdd(0x84)] = func(z80 *Z80) {
		z80.add(z80.ixh)
	}
	/* ADD A,z80.IXL() */
	opcodesMap[shift0xdd(0x85)] = func(z80 *Z80) {
		z80.add(z80.ixl)
	}
	/* ADD A,(ix+dd) */
	opcodesMap[shift0xdd(0x86)] = func(z80 *Z80) {
		var offset byte = z80.memory.readByte(z80.pc)
		z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
		z80.pc++
		var bytetemp byte = z80.memory.readByte(z80.IX() + uint16(signExtend(offset)))
		z80.add(bytetemp)
	}
	/* ADC A,z80.IXH() */
	opcodesMap[shift0xdd(0x8c)] = func(z80 *Z80) {
		z80.adc(z80.ixh)
	}
	/* ADC A,z80.IXL() */
	opcodesMap[shift0xdd(0x8d)] = func(z80 *Z80) {
		z80.adc(z80.ixl)
	}
	/* ADC A,(ix+dd) */
	opcodesMap[shift0xdd(0x8e)] = func(z80 *Z80) {
		var offset byte = z80.memory.readByte(z80.pc)
		z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
		z80.pc++
		var bytetemp byte = z80.memory.readByte(z80.IX() + uint16(signExtend(offset)))
		z80.adc(bytetemp)
	}
	/* SUB A,z80.IXH() */
	opcodesMap[shift0xdd(0x94)] = func(z80 *Z80) {
		z80.sub(z80.ixh)
	}
	/* SUB A,z80.IXL() */
	opcodesMap[shift0xdd(0x95)] = func(z80 *Z80) {
		z80.sub(z80.ixl)
	}
	/* SUB A,(ix+dd) */
	opcodesMap[shift0xdd(0x96)] = func(z80 *Z80) {
		var offset byte = z80.memory.readByte(z80.pc)
		z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
		z80.pc++
		var bytetemp byte = z80.memory.readByte(z80.IX() + uint16(signExtend(offset)))
		z80.sub(bytetemp)
	}
	/* SBC A,z80.IXH() */
	opcodesMap[shift0xdd(0x9c)] = func(z80 *Z80) {
		z80.sbc(z80.ixh)
	}
	/* SBC A,z80.IXL() */
	opcodesMap[shift0xdd(0x9d)] = func(z80 *Z80) {
		z80.sbc(z80.ixl)
	}
	/* SBC A,(ix+dd) */
	opcodesMap[shift0xdd(0x9e)] = func(z80 *Z80) {
		var offset byte = z80.memory.readByte(z80.pc)
		z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
		z80.pc++
		var bytetemp byte = z80.memory.readByte(z80.IX() + uint16(signExtend(offset)))
		z80.sbc(bytetemp)
	}
	/* AND A,z80.IXH() */
	opcodesMap[shift0xdd(0xa4)] = func(z80 *Z80) {
		z80.and(z80.ixh)
	}
	/* AND A,z80.IXL() */
	opcodesMap[shift0xdd(0xa5)] = func(z80 *Z80) {
		z80.and(z80.ixl)
	}
	/* AND A,(ix+dd) */
	opcodesMap[shift0xdd(0xa6)] = func(z80 *Z80) {
		var offset byte = z80.memory.readByte(z80.pc)
		z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
		z80.pc++
		var bytetemp byte = z80.memory.readByte(z80.IX() + uint16(signExtend(offset)))
		z80.and(bytetemp)
	}
	/* XOR A,z80.IXH() */
	opcodesMap[shift0xdd(0xac)] = func(z80 *Z80) {
		z80.xor(z80.ixh)
	}
	/* XOR A,z80.IXL() */
	opcodesMap[shift0xdd(0xad)] = func(z80 *Z80) {
		z80.xor(z80.ixl)
	}
	/* XOR A,(ix+dd) */
	opcodesMap[shift0xdd(0xae)] = func(z80 *Z80) {
		var offset byte = z80.memory.readByte(z80.pc)
		z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
		z80.pc++
		var bytetemp byte = z80.memory.readByte(z80.IX() + uint16(signExtend(offset)))
		z80.xor(bytetemp)
	}
	/* OR A,z80.IXH() */
	opcodesMap[shift0xdd(0xb4)] = func(z80 *Z80) {
		z80.or(z80.ixh)
	}
	/* OR A,z80.IXL() */
	opcodesMap[shift0xdd(0xb5)] = func(z80 *Z80) {
		z80.or(z80.ixl)
	}
	/* OR A,(ix+dd) */
	opcodesMap[shift0xdd(0xb6)] = func(z80 *Z80) {
		var offset byte = z80.memory.readByte(z80.pc)
		z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
		z80.pc++
		var bytetemp byte = z80.memory.readByte(z80.IX() + uint16(signExtend(offset)))
		z80.or(bytetemp)
	}
	/* CP A,z80.IXH() */
	opcodesMap[shift0xdd(0xbc)] = func(z80 *Z80) {
		z80.cp(z80.ixh)
	}
	/* CP A,z80.IXL() */
	opcodesMap[shift0xdd(0xbd)] = func(z80 *Z80) {
		z80.cp(z80.ixl)
	}
	/* CP A,(ix+dd) */
	opcodesMap[shift0xdd(0xbe)] = func(z80 *Z80) {
		var offset byte = z80.memory.readByte(z80.pc)
		z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
		z80.pc++
		var bytetemp byte = z80.memory.readByte(z80.IX() + uint16(signExtend(offset)))
		z80.cp(bytetemp)
	}
	/* shift DDFDCB */
	opcodesMap[shift0xdd(0xcb)] = func(z80 *Z80) {
	}
	/* POP ix */
	opcodesMap[shift0xdd(0xe1)] = func(z80 *Z80) {
		z80.ixl, z80.ixh = z80.pop16()
	}
	/* EX (SP),ix */
	opcodesMap[shift0xdd(0xe3)] = func(z80 *Z80) {
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
	opcodesMap[shift0xdd(0xe5)] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq(z80.IR(), 1)
		z80.push16(z80.ixl, z80.ixh)
	}
	/* JP ix */
	opcodesMap[shift0xdd(0xe9)] = func(z80 *Z80) {
		z80.pc = z80.IX() /* NB: NOT INDIRECT! */
	}
	/* LD SP,ix */
	opcodesMap[shift0xdd(0xf9)] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 2)
		z80.sp = z80.IX()
	}

	// END of 0xdd shifted opcodes

	// BEGIN of 0xfd shifted opcodes

	/* ADD iy,BC */
	opcodesMap[shift0xfd(0x09)] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 7)
		z80.add16(z80.iy, z80.BC())
	}
	/* ADD iy,DE */
	opcodesMap[shift0xfd(0x19)] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 7)
		z80.add16(z80.iy, z80.DE())
	}
	/* LD iy,nnnn */
	opcodesMap[shift0xfd(0x21)] = func(z80 *Z80) {
		b1 := z80.memory.readByte(z80.pc)
		z80.pc++
		b2 := z80.memory.readByte(z80.pc)
		z80.pc++
		z80.setIY(joinBytes(b2, b1))
	}
	/* LD (nnnn),iy */
	opcodesMap[shift0xfd(0x22)] = func(z80 *Z80) {
		z80.ld16nnrr(z80.iyl, z80.iyh)
		// break
	}
	/* INC iy */
	opcodesMap[shift0xfd(0x23)] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 2)
		z80.incIY()
	}
	/* INC z80.IYH() */
	opcodesMap[shift0xfd(0x24)] = func(z80 *Z80) {
		z80.incIYH()
	}
	/* DEC z80.IYH() */
	opcodesMap[shift0xfd(0x25)] = func(z80 *Z80) {
		z80.decIYH()
	}
	/* LD z80.IYH(),nn */
	opcodesMap[shift0xfd(0x26)] = func(z80 *Z80) {
		z80.iyh = z80.memory.readByte(z80.pc)
		z80.pc++
	}
	/* ADD iy,iy */
	opcodesMap[shift0xfd(0x29)] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 7)
		z80.add16(z80.iy, z80.IY())
	}
	/* LD iy,(nnnn) */
	opcodesMap[shift0xfd(0x2a)] = func(z80 *Z80) {
		z80.ld16rrnn(&z80.iyl, &z80.iyh)
		// break
	}
	/* DEC iy */
	opcodesMap[shift0xfd(0x2b)] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 2)
		z80.decIY()
	}
	/* INC z80.IYL() */
	opcodesMap[shift0xfd(0x2c)] = func(z80 *Z80) {
		z80.incIYL()
	}
	/* DEC z80.IYL() */
	opcodesMap[shift0xfd(0x2d)] = func(z80 *Z80) {
		z80.decIYL()
	}
	/* LD z80.IYL(),nn */
	opcodesMap[shift0xfd(0x2e)] = func(z80 *Z80) {
		z80.iyl = z80.memory.readByte(z80.pc)
		z80.pc++
	}
	/* INC (iy+dd) */
	opcodesMap[shift0xfd(0x34)] = func(z80 *Z80) {
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
	opcodesMap[shift0xfd(0x35)] = func(z80 *Z80) {
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
	opcodesMap[shift0xfd(0x36)] = func(z80 *Z80) {
		offset := z80.memory.readByte(z80.pc)
		z80.pc++
		value := z80.memory.readByte(z80.pc)
		z80.memory.contendReadNoMreq_loop(z80.pc, 1, 2)
		z80.pc++
		z80.memory.writeByte(z80.IY()+uint16(signExtend(offset)), value)
	}
	/* ADD iy,SP */
	opcodesMap[shift0xfd(0x39)] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 7)
		z80.add16(z80.iy, z80.SP())
	}
	/* LD B,z80.IYH() */
	opcodesMap[shift0xfd(0x44)] = func(z80 *Z80) {
		z80.b = z80.iyh
	}
	/* LD B,z80.IYL() */
	opcodesMap[shift0xfd(0x45)] = func(z80 *Z80) {
		z80.b = z80.iyl
	}
	/* LD B,(iy+dd) */
	opcodesMap[shift0xfd(0x46)] = func(z80 *Z80) {
		var offset byte = z80.memory.readByte(z80.pc)
		z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
		z80.pc++
		z80.b = z80.memory.readByte(z80.IY() + uint16(signExtend(offset)))
	}
	/* LD C,z80.IYH() */
	opcodesMap[shift0xfd(0x4c)] = func(z80 *Z80) {
		z80.c = z80.iyh
	}
	/* LD C,z80.IYL() */
	opcodesMap[shift0xfd(0x4d)] = func(z80 *Z80) {
		z80.c = z80.iyl
	}
	/* LD C,(iy+dd) */
	opcodesMap[shift0xfd(0x4e)] = func(z80 *Z80) {
		var offset byte = z80.memory.readByte(z80.pc)
		z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
		z80.pc++
		z80.c = z80.memory.readByte(z80.IY() + uint16(signExtend(offset)))
	}
	/* LD D,z80.IYH() */
	opcodesMap[shift0xfd(0x54)] = func(z80 *Z80) {
		z80.d = z80.iyh
	}
	/* LD D,z80.IYL() */
	opcodesMap[shift0xfd(0x55)] = func(z80 *Z80) {
		z80.d = z80.iyl
	}
	/* LD D,(iy+dd) */
	opcodesMap[shift0xfd(0x56)] = func(z80 *Z80) {
		var offset byte = z80.memory.readByte(z80.pc)
		z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
		z80.pc++
		z80.d = z80.memory.readByte(z80.IY() + uint16(signExtend(offset)))
	}
	/* LD E,z80.IYH() */
	opcodesMap[shift0xfd(0x5c)] = func(z80 *Z80) {
		z80.e = z80.iyh
	}
	/* LD E,z80.IYL() */
	opcodesMap[shift0xfd(0x5d)] = func(z80 *Z80) {
		z80.e = z80.iyl
	}
	/* LD E,(iy+dd) */
	opcodesMap[shift0xfd(0x5e)] = func(z80 *Z80) {
		var offset byte = z80.memory.readByte(z80.pc)
		z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
		z80.pc++
		z80.e = z80.memory.readByte(z80.IY() + uint16(signExtend(offset)))
	}
	/* LD z80.IYH(),B */
	opcodesMap[shift0xfd(0x60)] = func(z80 *Z80) {
		z80.iyh = z80.b
	}
	/* LD z80.IYH(),C */
	opcodesMap[shift0xfd(0x61)] = func(z80 *Z80) {
		z80.iyh = z80.c
	}
	/* LD z80.IYH(),D */
	opcodesMap[shift0xfd(0x62)] = func(z80 *Z80) {
		z80.iyh = z80.d
	}
	/* LD z80.IYH(),E */
	opcodesMap[shift0xfd(0x63)] = func(z80 *Z80) {
		z80.iyh = z80.e
	}
	/* LD z80.IYH(),z80.IYH() */
	opcodesMap[shift0xfd(0x64)] = func(z80 *Z80) {
	}
	/* LD z80.IYH(),z80.IYL() */
	opcodesMap[shift0xfd(0x65)] = func(z80 *Z80) {
		z80.iyh = z80.iyl
	}
	/* LD H,(iy+dd) */
	opcodesMap[shift0xfd(0x66)] = func(z80 *Z80) {
		var offset byte = z80.memory.readByte(z80.pc)
		z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
		z80.pc++
		z80.h = z80.memory.readByte(z80.IY() + uint16(signExtend(offset)))
	}
	/* LD z80.IYH(),A */
	opcodesMap[shift0xfd(0x67)] = func(z80 *Z80) {
		z80.iyh = z80.a
	}
	/* LD z80.IYL(),B */
	opcodesMap[shift0xfd(0x68)] = func(z80 *Z80) {
		z80.iyl = z80.b
	}
	/* LD z80.IYL(),C */
	opcodesMap[shift0xfd(0x69)] = func(z80 *Z80) {
		z80.iyl = z80.c
	}
	/* LD z80.IYL(),D */
	opcodesMap[shift0xfd(0x6a)] = func(z80 *Z80) {
		z80.iyl = z80.d
	}
	/* LD z80.IYL(),E */
	opcodesMap[shift0xfd(0x6b)] = func(z80 *Z80) {
		z80.iyl = z80.e
	}
	/* LD z80.IYL(),z80.IYH() */
	opcodesMap[shift0xfd(0x6c)] = func(z80 *Z80) {
		z80.iyl = z80.iyh
	}
	/* LD z80.IYL(),z80.IYL() */
	opcodesMap[shift0xfd(0x6d)] = func(z80 *Z80) {
	}
	/* LD L,(iy+dd) */
	opcodesMap[shift0xfd(0x6e)] = func(z80 *Z80) {
		var offset byte = z80.memory.readByte(z80.pc)
		z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
		z80.pc++
		z80.l = z80.memory.readByte(z80.IY() + uint16(signExtend(offset)))
	}
	/* LD z80.IYL(),A */
	opcodesMap[shift0xfd(0x6f)] = func(z80 *Z80) {
		z80.iyl = z80.a
	}
	/* LD (iy+dd),B */
	opcodesMap[shift0xfd(0x70)] = func(z80 *Z80) {
		offset := z80.memory.readByte(z80.pc)
		z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
		z80.pc++
		z80.memory.writeByte(z80.IY()+uint16(signExtend(offset)), z80.b)
	}
	/* LD (iy+dd),C */
	opcodesMap[shift0xfd(0x71)] = func(z80 *Z80) {
		offset := z80.memory.readByte(z80.pc)
		z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
		z80.pc++
		z80.memory.writeByte(z80.IY()+uint16(signExtend(offset)), z80.c)
	}
	/* LD (iy+dd),D */
	opcodesMap[shift0xfd(0x72)] = func(z80 *Z80) {
		offset := z80.memory.readByte(z80.pc)
		z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
		z80.pc++
		z80.memory.writeByte(z80.IY()+uint16(signExtend(offset)), z80.d)
	}
	/* LD (iy+dd),E */
	opcodesMap[shift0xfd(0x73)] = func(z80 *Z80) {
		offset := z80.memory.readByte(z80.pc)
		z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
		z80.pc++
		z80.memory.writeByte(z80.IY()+uint16(signExtend(offset)), z80.e)
	}
	/* LD (iy+dd),H */
	opcodesMap[shift0xfd(0x74)] = func(z80 *Z80) {
		offset := z80.memory.readByte(z80.pc)
		z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
		z80.pc++
		z80.memory.writeByte(z80.IY()+uint16(signExtend(offset)), z80.h)
	}
	/* LD (iy+dd),L */
	opcodesMap[shift0xfd(0x75)] = func(z80 *Z80) {
		offset := z80.memory.readByte(z80.pc)
		z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
		z80.pc++
		z80.memory.writeByte(z80.IY()+uint16(signExtend(offset)), z80.l)
	}
	/* LD (iy+dd),A */
	opcodesMap[shift0xfd(0x77)] = func(z80 *Z80) {
		offset := z80.memory.readByte(z80.pc)
		z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
		z80.pc++
		z80.memory.writeByte(z80.IY()+uint16(signExtend(offset)), z80.a)
	}
	/* LD A,z80.IYH() */
	opcodesMap[shift0xfd(0x7c)] = func(z80 *Z80) {
		z80.a = z80.iyh
	}
	/* LD A,z80.IYL() */
	opcodesMap[shift0xfd(0x7d)] = func(z80 *Z80) {
		z80.a = z80.iyl
	}
	/* LD A,(iy+dd) */
	opcodesMap[shift0xfd(0x7e)] = func(z80 *Z80) {
		var offset byte = z80.memory.readByte(z80.pc)
		z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
		z80.pc++
		z80.a = z80.memory.readByte(z80.IY() + uint16(signExtend(offset)))
	}
	/* ADD A,z80.IYH() */
	opcodesMap[shift0xfd(0x84)] = func(z80 *Z80) {
		z80.add(z80.iyh)
	}
	/* ADD A,z80.IYL() */
	opcodesMap[shift0xfd(0x85)] = func(z80 *Z80) {
		z80.add(z80.iyl)
	}
	/* ADD A,(iy+dd) */
	opcodesMap[shift0xfd(0x86)] = func(z80 *Z80) {
		var offset byte = z80.memory.readByte(z80.pc)
		z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
		z80.pc++
		var bytetemp byte = z80.memory.readByte(z80.IY() + uint16(signExtend(offset)))
		z80.add(bytetemp)
	}
	/* ADC A,z80.IYH() */
	opcodesMap[shift0xfd(0x8c)] = func(z80 *Z80) {
		z80.adc(z80.iyh)
	}
	/* ADC A,z80.IYL() */
	opcodesMap[shift0xfd(0x8d)] = func(z80 *Z80) {
		z80.adc(z80.iyl)
	}
	/* ADC A,(iy+dd) */
	opcodesMap[shift0xfd(0x8e)] = func(z80 *Z80) {
		var offset byte = z80.memory.readByte(z80.pc)
		z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
		z80.pc++
		var bytetemp byte = z80.memory.readByte(z80.IY() + uint16(signExtend(offset)))
		z80.adc(bytetemp)
	}
	/* SUB A,z80.IYH() */
	opcodesMap[shift0xfd(0x94)] = func(z80 *Z80) {
		z80.sub(z80.iyh)
	}
	/* SUB A,z80.IYL() */
	opcodesMap[shift0xfd(0x95)] = func(z80 *Z80) {
		z80.sub(z80.iyl)
	}
	/* SUB A,(iy+dd) */
	opcodesMap[shift0xfd(0x96)] = func(z80 *Z80) {
		var offset byte = z80.memory.readByte(z80.pc)
		z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
		z80.pc++
		var bytetemp byte = z80.memory.readByte(z80.IY() + uint16(signExtend(offset)))
		z80.sub(bytetemp)
	}
	/* SBC A,z80.IYH() */
	opcodesMap[shift0xfd(0x9c)] = func(z80 *Z80) {
		z80.sbc(z80.iyh)
	}
	/* SBC A,z80.IYL() */
	opcodesMap[shift0xfd(0x9d)] = func(z80 *Z80) {
		z80.sbc(z80.iyl)
	}
	/* SBC A,(iy+dd) */
	opcodesMap[shift0xfd(0x9e)] = func(z80 *Z80) {
		var offset byte = z80.memory.readByte(z80.pc)
		z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
		z80.pc++
		var bytetemp byte = z80.memory.readByte(z80.IY() + uint16(signExtend(offset)))
		z80.sbc(bytetemp)
	}
	/* AND A,z80.IYH() */
	opcodesMap[shift0xfd(0xa4)] = func(z80 *Z80) {
		z80.and(z80.iyh)
	}
	/* AND A,z80.IYL() */
	opcodesMap[shift0xfd(0xa5)] = func(z80 *Z80) {
		z80.and(z80.iyl)
	}
	/* AND A,(iy+dd) */
	opcodesMap[shift0xfd(0xa6)] = func(z80 *Z80) {
		var offset byte = z80.memory.readByte(z80.pc)
		z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
		z80.pc++
		var bytetemp byte = z80.memory.readByte(z80.IY() + uint16(signExtend(offset)))
		z80.and(bytetemp)
	}
	/* XOR A,z80.IYH() */
	opcodesMap[shift0xfd(0xac)] = func(z80 *Z80) {
		z80.xor(z80.iyh)
	}
	/* XOR A,z80.IYL() */
	opcodesMap[shift0xfd(0xad)] = func(z80 *Z80) {
		z80.xor(z80.iyl)
	}
	/* XOR A,(iy+dd) */
	opcodesMap[shift0xfd(0xae)] = func(z80 *Z80) {
		var offset byte = z80.memory.readByte(z80.pc)
		z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
		z80.pc++
		var bytetemp byte = z80.memory.readByte(z80.IY() + uint16(signExtend(offset)))
		z80.xor(bytetemp)
	}
	/* OR A,z80.IYH() */
	opcodesMap[shift0xfd(0xb4)] = func(z80 *Z80) {
		z80.or(z80.iyh)
	}
	/* OR A,z80.IYL() */
	opcodesMap[shift0xfd(0xb5)] = func(z80 *Z80) {
		z80.or(z80.iyl)
	}
	/* OR A,(iy+dd) */
	opcodesMap[shift0xfd(0xb6)] = func(z80 *Z80) {
		var offset byte = z80.memory.readByte(z80.pc)
		z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
		z80.pc++
		var bytetemp byte = z80.memory.readByte(z80.IY() + uint16(signExtend(offset)))
		z80.or(bytetemp)
	}
	/* CP A,z80.IYH() */
	opcodesMap[shift0xfd(0xbc)] = func(z80 *Z80) {
		z80.cp(z80.iyh)
	}
	/* CP A,z80.IYL() */
	opcodesMap[shift0xfd(0xbd)] = func(z80 *Z80) {
		z80.cp(z80.iyl)
	}
	/* CP A,(iy+dd) */
	opcodesMap[shift0xfd(0xbe)] = func(z80 *Z80) {
		var offset byte = z80.memory.readByte(z80.pc)
		z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
		z80.pc++
		var bytetemp byte = z80.memory.readByte(z80.IY() + uint16(signExtend(offset)))
		z80.cp(bytetemp)
	}
	/* shift DDFDCB */
	opcodesMap[shift0xfd(0xcb)] = func(z80 *Z80) {
	}
	/* POP iy */
	opcodesMap[shift0xfd(0xe1)] = func(z80 *Z80) {
		z80.iyl, z80.iyh = z80.pop16()
	}
	/* EX (SP),iy */
	opcodesMap[shift0xfd(0xe3)] = func(z80 *Z80) {
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
	opcodesMap[shift0xfd(0xe5)] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq(z80.IR(), 1)
		z80.push16(z80.iyl, z80.iyh)
	}
	/* JP iy */
	opcodesMap[shift0xfd(0xe9)] = func(z80 *Z80) {
		z80.pc = z80.IY() /* NB: NOT INDIRECT! */
	}
	/* LD SP,iy */
	opcodesMap[shift0xfd(0xf9)] = func(z80 *Z80) {
		z80.memory.contendReadNoMreq_loop(z80.IR(), 1, 2)
		z80.sp = z80.IY()
	}

	// END of 0xfd shifted opcodes

	// BEGIN of 0xddfdcb shifted opcodes

	/* LD B,RLC (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x00)] = func(z80 *Z80) {
		z80.b = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.b = z80.rlc(z80.b)
		z80.memory.writeByte(z80.tempaddr, z80.b)
	}
	/* LD C,RLC (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x01)] = func(z80 *Z80) {
		z80.c = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.c = z80.rlc(z80.c)
		z80.memory.writeByte(z80.tempaddr, z80.c)
	}
	/* LD D,RLC (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x02)] = func(z80 *Z80) {
		z80.d = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.d = z80.rlc(z80.d)
		z80.memory.writeByte(z80.tempaddr, z80.d)
	}
	/* LD E,RLC (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x03)] = func(z80 *Z80) {
		z80.e = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.e = z80.rlc(z80.e)
		z80.memory.writeByte(z80.tempaddr, z80.e)
	}
	/* LD H,RLC (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x04)] = func(z80 *Z80) {
		z80.h = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.h = z80.rlc(z80.h)
		z80.memory.writeByte(z80.tempaddr, z80.h)
	}
	/* LD L,RLC (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x05)] = func(z80 *Z80) {
		z80.l = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.l = z80.rlc(z80.l)
		z80.memory.writeByte(z80.tempaddr, z80.l)
	}
	/* RLC (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x06)] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		bytetemp = z80.rlc(bytetemp)
		z80.memory.writeByte(z80.tempaddr, bytetemp)
	}
	/* LD A,RLC (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x07)] = func(z80 *Z80) {
		z80.a = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.a = z80.rlc(z80.a)
		z80.memory.writeByte(z80.tempaddr, z80.a)
	}
	/* LD B,RRC (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x08)] = func(z80 *Z80) {
		z80.b = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.b = z80.rrc(z80.b)
		z80.memory.writeByte(z80.tempaddr, z80.b)
	}
	/* LD C,RRC (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x09)] = func(z80 *Z80) {
		z80.c = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.c = z80.rrc(z80.c)
		z80.memory.writeByte(z80.tempaddr, z80.c)
	}
	/* LD D,RRC (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x0a)] = func(z80 *Z80) {
		z80.d = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.d = z80.rrc(z80.d)
		z80.memory.writeByte(z80.tempaddr, z80.d)
	}
	/* LD E,RRC (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x0b)] = func(z80 *Z80) {
		z80.e = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.e = z80.rrc(z80.e)
		z80.memory.writeByte(z80.tempaddr, z80.e)
	}
	/* LD H,RRC (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x0c)] = func(z80 *Z80) {
		z80.h = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.h = z80.rrc(z80.h)
		z80.memory.writeByte(z80.tempaddr, z80.h)
	}
	/* LD L,RRC (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x0d)] = func(z80 *Z80) {
		z80.l = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.l = z80.rrc(z80.l)
		z80.memory.writeByte(z80.tempaddr, z80.l)
	}
	/* RRC (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x0e)] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		bytetemp = z80.rrc(bytetemp)
		z80.memory.writeByte(z80.tempaddr, bytetemp)
	}
	/* LD A,RRC (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x0f)] = func(z80 *Z80) {
		z80.a = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.a = z80.rrc(z80.a)
		z80.memory.writeByte(z80.tempaddr, z80.a)
	}
	/* LD B,RL (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x10)] = func(z80 *Z80) {
		z80.b = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.b = z80.rl(z80.b)
		z80.memory.writeByte(z80.tempaddr, z80.b)
	}
	/* LD C,RL (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x11)] = func(z80 *Z80) {
		z80.c = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.c = z80.rl(z80.c)
		z80.memory.writeByte(z80.tempaddr, z80.c)
	}
	/* LD D,RL (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x12)] = func(z80 *Z80) {
		z80.d = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.d = z80.rl(z80.d)
		z80.memory.writeByte(z80.tempaddr, z80.d)
	}
	/* LD E,RL (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x13)] = func(z80 *Z80) {
		z80.e = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.e = z80.rl(z80.e)
		z80.memory.writeByte(z80.tempaddr, z80.e)
	}
	/* LD H,RL (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x14)] = func(z80 *Z80) {
		z80.h = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.h = z80.rl(z80.h)
		z80.memory.writeByte(z80.tempaddr, z80.h)
	}
	/* LD L,RL (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x15)] = func(z80 *Z80) {
		z80.l = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.l = z80.rl(z80.l)
		z80.memory.writeByte(z80.tempaddr, z80.l)
	}
	/* RL (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x16)] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		bytetemp = z80.rl(bytetemp)
		z80.memory.writeByte(z80.tempaddr, bytetemp)
	}
	/* LD A,RL (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x17)] = func(z80 *Z80) {
		z80.a = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.a = z80.rl(z80.a)
		z80.memory.writeByte(z80.tempaddr, z80.a)
	}
	/* LD B,RR (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x18)] = func(z80 *Z80) {
		z80.b = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.b = z80.rr(z80.b)
		z80.memory.writeByte(z80.tempaddr, z80.b)
	}
	/* LD C,RR (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x19)] = func(z80 *Z80) {
		z80.c = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.c = z80.rr(z80.c)
		z80.memory.writeByte(z80.tempaddr, z80.c)
	}
	/* LD D,RR (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x1a)] = func(z80 *Z80) {
		z80.d = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.d = z80.rr(z80.d)
		z80.memory.writeByte(z80.tempaddr, z80.d)
	}
	/* LD E,RR (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x1b)] = func(z80 *Z80) {
		z80.e = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.e = z80.rr(z80.e)
		z80.memory.writeByte(z80.tempaddr, z80.e)
	}
	/* LD H,RR (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x1c)] = func(z80 *Z80) {
		z80.h = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.h = z80.rr(z80.h)
		z80.memory.writeByte(z80.tempaddr, z80.h)
	}
	/* LD L,RR (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x1d)] = func(z80 *Z80) {
		z80.l = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.l = z80.rr(z80.l)
		z80.memory.writeByte(z80.tempaddr, z80.l)
	}
	/* RR (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x1e)] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		bytetemp = z80.rr(bytetemp)
		z80.memory.writeByte(z80.tempaddr, bytetemp)
	}
	/* LD A,RR (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x1f)] = func(z80 *Z80) {
		z80.a = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.a = z80.rr(z80.a)
		z80.memory.writeByte(z80.tempaddr, z80.a)
	}
	/* LD B,SLA (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x20)] = func(z80 *Z80) {
		z80.b = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.b = z80.sla(z80.b)
		z80.memory.writeByte(z80.tempaddr, z80.b)
	}
	/* LD C,SLA (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x21)] = func(z80 *Z80) {
		z80.c = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.c = z80.sla(z80.c)
		z80.memory.writeByte(z80.tempaddr, z80.c)
	}
	/* LD D,SLA (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x22)] = func(z80 *Z80) {
		z80.d = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.d = z80.sla(z80.d)
		z80.memory.writeByte(z80.tempaddr, z80.d)
	}
	/* LD E,SLA (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x23)] = func(z80 *Z80) {
		z80.e = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.e = z80.sla(z80.e)
		z80.memory.writeByte(z80.tempaddr, z80.e)
	}
	/* LD H,SLA (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x24)] = func(z80 *Z80) {
		z80.h = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.h = z80.sla(z80.h)
		z80.memory.writeByte(z80.tempaddr, z80.h)
	}
	/* LD L,SLA (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x25)] = func(z80 *Z80) {
		z80.l = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.l = z80.sla(z80.l)
		z80.memory.writeByte(z80.tempaddr, z80.l)
	}
	/* SLA (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x26)] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		bytetemp = z80.sla(bytetemp)
		z80.memory.writeByte(z80.tempaddr, bytetemp)
	}
	/* LD A,SLA (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x27)] = func(z80 *Z80) {
		z80.a = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.a = z80.sla(z80.a)
		z80.memory.writeByte(z80.tempaddr, z80.a)
	}
	/* LD B,SRA (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x28)] = func(z80 *Z80) {
		z80.b = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.b = z80.sra(z80.b)
		z80.memory.writeByte(z80.tempaddr, z80.b)
	}
	/* LD C,SRA (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x29)] = func(z80 *Z80) {
		z80.c = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.c = z80.sra(z80.c)
		z80.memory.writeByte(z80.tempaddr, z80.c)
	}
	/* LD D,SRA (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x2a)] = func(z80 *Z80) {
		z80.d = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.d = z80.sra(z80.d)
		z80.memory.writeByte(z80.tempaddr, z80.d)
	}
	/* LD E,SRA (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x2b)] = func(z80 *Z80) {
		z80.e = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.e = z80.sra(z80.e)
		z80.memory.writeByte(z80.tempaddr, z80.e)
	}
	/* LD H,SRA (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x2c)] = func(z80 *Z80) {
		z80.h = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.h = z80.sra(z80.h)
		z80.memory.writeByte(z80.tempaddr, z80.h)
	}
	/* LD L,SRA (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x2d)] = func(z80 *Z80) {
		z80.l = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.l = z80.sra(z80.l)
		z80.memory.writeByte(z80.tempaddr, z80.l)
	}
	/* SRA (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x2e)] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		bytetemp = z80.sra(bytetemp)
		z80.memory.writeByte(z80.tempaddr, bytetemp)
	}
	/* LD A,SRA (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x2f)] = func(z80 *Z80) {
		z80.a = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.a = z80.sra(z80.a)
		z80.memory.writeByte(z80.tempaddr, z80.a)
	}
	/* LD B,SLL (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x30)] = func(z80 *Z80) {
		z80.b = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.b = z80.sll(z80.b)
		z80.memory.writeByte(z80.tempaddr, z80.b)
	}
	/* LD C,SLL (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x31)] = func(z80 *Z80) {
		z80.c = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.c = z80.sll(z80.c)
		z80.memory.writeByte(z80.tempaddr, z80.c)
	}
	/* LD D,SLL (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x32)] = func(z80 *Z80) {
		z80.d = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.d = z80.sll(z80.d)
		z80.memory.writeByte(z80.tempaddr, z80.d)
	}
	/* LD E,SLL (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x33)] = func(z80 *Z80) {
		z80.e = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.e = z80.sll(z80.e)
		z80.memory.writeByte(z80.tempaddr, z80.e)
	}
	/* LD H,SLL (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x34)] = func(z80 *Z80) {
		z80.h = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.h = z80.sll(z80.h)
		z80.memory.writeByte(z80.tempaddr, z80.h)
	}
	/* LD L,SLL (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x35)] = func(z80 *Z80) {
		z80.l = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.l = z80.sll(z80.l)
		z80.memory.writeByte(z80.tempaddr, z80.l)
	}
	/* SLL (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x36)] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		bytetemp = z80.sll(bytetemp)
		z80.memory.writeByte(z80.tempaddr, bytetemp)
	}
	/* LD A,SLL (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x37)] = func(z80 *Z80) {
		z80.a = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.a = z80.sll(z80.a)
		z80.memory.writeByte(z80.tempaddr, z80.a)
	}
	/* LD B,SRL (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x38)] = func(z80 *Z80) {
		z80.b = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.b = z80.srl(z80.b)
		z80.memory.writeByte(z80.tempaddr, z80.b)
	}
	/* LD C,SRL (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x39)] = func(z80 *Z80) {
		z80.c = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.c = z80.srl(z80.c)
		z80.memory.writeByte(z80.tempaddr, z80.c)
	}
	/* LD D,SRL (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x3a)] = func(z80 *Z80) {
		z80.d = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.d = z80.srl(z80.d)
		z80.memory.writeByte(z80.tempaddr, z80.d)
	}
	/* LD E,SRL (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x3b)] = func(z80 *Z80) {
		z80.e = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.e = z80.srl(z80.e)
		z80.memory.writeByte(z80.tempaddr, z80.e)
	}
	/* LD H,SRL (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x3c)] = func(z80 *Z80) {
		z80.h = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.h = z80.srl(z80.h)
		z80.memory.writeByte(z80.tempaddr, z80.h)
	}
	/* LD L,SRL (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x3d)] = func(z80 *Z80) {
		z80.l = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.l = z80.srl(z80.l)
		z80.memory.writeByte(z80.tempaddr, z80.l)
	}
	/* SRL (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x3e)] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		bytetemp = z80.srl(bytetemp)
		z80.memory.writeByte(z80.tempaddr, bytetemp)
	}
	/* LD A,SRL (REGISTER+dd) */
	opcodesMap[shift0xddcb(0x3f)] = func(z80 *Z80) {
		z80.a = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.a = z80.srl(z80.a)
		z80.memory.writeByte(z80.tempaddr, z80.a)
	}
	/* BIT 0,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0x47)] = func(z80 *Z80) {
		bytetemp := z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.biti(0, bytetemp, z80.tempaddr)
	}
	// Fallthrough cases
	opcodesMap[shift0xddcb(0x40)] = opcodesMap[shift0xddcb(0x47)]
	opcodesMap[shift0xddcb(0x41)] = opcodesMap[shift0xddcb(0x47)]
	opcodesMap[shift0xddcb(0x42)] = opcodesMap[shift0xddcb(0x47)]
	opcodesMap[shift0xddcb(0x43)] = opcodesMap[shift0xddcb(0x47)]
	opcodesMap[shift0xddcb(0x44)] = opcodesMap[shift0xddcb(0x47)]
	opcodesMap[shift0xddcb(0x45)] = opcodesMap[shift0xddcb(0x47)]
	opcodesMap[shift0xddcb(0x46)] = opcodesMap[shift0xddcb(0x47)]
	/* BIT 1,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0x4f)] = func(z80 *Z80) {
		bytetemp := z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.biti(1, bytetemp, z80.tempaddr)
	}
	// Fallthrough cases
	opcodesMap[shift0xddcb(0x48)] = opcodesMap[shift0xddcb(0x4f)]
	opcodesMap[shift0xddcb(0x49)] = opcodesMap[shift0xddcb(0x4f)]
	opcodesMap[shift0xddcb(0x4a)] = opcodesMap[shift0xddcb(0x4f)]
	opcodesMap[shift0xddcb(0x4b)] = opcodesMap[shift0xddcb(0x4f)]
	opcodesMap[shift0xddcb(0x4c)] = opcodesMap[shift0xddcb(0x4f)]
	opcodesMap[shift0xddcb(0x4d)] = opcodesMap[shift0xddcb(0x4f)]
	opcodesMap[shift0xddcb(0x4e)] = opcodesMap[shift0xddcb(0x4f)]
	/* BIT 2,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0x57)] = func(z80 *Z80) {
		bytetemp := z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.biti(2, bytetemp, z80.tempaddr)
	}
	// Fallthrough cases
	opcodesMap[shift0xddcb(0x50)] = opcodesMap[shift0xddcb(0x57)]
	opcodesMap[shift0xddcb(0x51)] = opcodesMap[shift0xddcb(0x57)]
	opcodesMap[shift0xddcb(0x52)] = opcodesMap[shift0xddcb(0x57)]
	opcodesMap[shift0xddcb(0x53)] = opcodesMap[shift0xddcb(0x57)]
	opcodesMap[shift0xddcb(0x54)] = opcodesMap[shift0xddcb(0x57)]
	opcodesMap[shift0xddcb(0x55)] = opcodesMap[shift0xddcb(0x57)]
	opcodesMap[shift0xddcb(0x56)] = opcodesMap[shift0xddcb(0x57)]
	/* BIT 3,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0x5f)] = func(z80 *Z80) {
		bytetemp := z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.biti(3, bytetemp, z80.tempaddr)
	}
	// Fallthrough cases
	opcodesMap[shift0xddcb(0x58)] = opcodesMap[shift0xddcb(0x5f)]
	opcodesMap[shift0xddcb(0x59)] = opcodesMap[shift0xddcb(0x5f)]
	opcodesMap[shift0xddcb(0x5a)] = opcodesMap[shift0xddcb(0x5f)]
	opcodesMap[shift0xddcb(0x5b)] = opcodesMap[shift0xddcb(0x5f)]
	opcodesMap[shift0xddcb(0x5c)] = opcodesMap[shift0xddcb(0x5f)]
	opcodesMap[shift0xddcb(0x5d)] = opcodesMap[shift0xddcb(0x5f)]
	opcodesMap[shift0xddcb(0x5e)] = opcodesMap[shift0xddcb(0x5f)]
	/* BIT 4,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0x67)] = func(z80 *Z80) {
		bytetemp := z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.biti(4, bytetemp, z80.tempaddr)
	}
	// Fallthrough cases
	opcodesMap[shift0xddcb(0x60)] = opcodesMap[shift0xddcb(0x67)]
	opcodesMap[shift0xddcb(0x61)] = opcodesMap[shift0xddcb(0x67)]
	opcodesMap[shift0xddcb(0x62)] = opcodesMap[shift0xddcb(0x67)]
	opcodesMap[shift0xddcb(0x63)] = opcodesMap[shift0xddcb(0x67)]
	opcodesMap[shift0xddcb(0x64)] = opcodesMap[shift0xddcb(0x67)]
	opcodesMap[shift0xddcb(0x65)] = opcodesMap[shift0xddcb(0x67)]
	opcodesMap[shift0xddcb(0x66)] = opcodesMap[shift0xddcb(0x67)]
	/* BIT 5,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0x6f)] = func(z80 *Z80) {
		bytetemp := z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.biti(5, bytetemp, z80.tempaddr)
	}
	// Fallthrough cases
	opcodesMap[shift0xddcb(0x68)] = opcodesMap[shift0xddcb(0x6f)]
	opcodesMap[shift0xddcb(0x69)] = opcodesMap[shift0xddcb(0x6f)]
	opcodesMap[shift0xddcb(0x6a)] = opcodesMap[shift0xddcb(0x6f)]
	opcodesMap[shift0xddcb(0x6b)] = opcodesMap[shift0xddcb(0x6f)]
	opcodesMap[shift0xddcb(0x6c)] = opcodesMap[shift0xddcb(0x6f)]
	opcodesMap[shift0xddcb(0x6d)] = opcodesMap[shift0xddcb(0x6f)]
	opcodesMap[shift0xddcb(0x6e)] = opcodesMap[shift0xddcb(0x6f)]
	/* BIT 6,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0x77)] = func(z80 *Z80) {
		bytetemp := z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.biti(6, bytetemp, z80.tempaddr)
	}
	// Fallthrough cases
	opcodesMap[shift0xddcb(0x70)] = opcodesMap[shift0xddcb(0x77)]
	opcodesMap[shift0xddcb(0x71)] = opcodesMap[shift0xddcb(0x77)]
	opcodesMap[shift0xddcb(0x72)] = opcodesMap[shift0xddcb(0x77)]
	opcodesMap[shift0xddcb(0x73)] = opcodesMap[shift0xddcb(0x77)]
	opcodesMap[shift0xddcb(0x74)] = opcodesMap[shift0xddcb(0x77)]
	opcodesMap[shift0xddcb(0x75)] = opcodesMap[shift0xddcb(0x77)]
	opcodesMap[shift0xddcb(0x76)] = opcodesMap[shift0xddcb(0x77)]
	/* BIT 7,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0x7f)] = func(z80 *Z80) {
		bytetemp := z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.biti(7, bytetemp, z80.tempaddr)
	}
	// Fallthrough cases
	opcodesMap[shift0xddcb(0x78)] = opcodesMap[shift0xddcb(0x7f)]
	opcodesMap[shift0xddcb(0x79)] = opcodesMap[shift0xddcb(0x7f)]
	opcodesMap[shift0xddcb(0x7a)] = opcodesMap[shift0xddcb(0x7f)]
	opcodesMap[shift0xddcb(0x7b)] = opcodesMap[shift0xddcb(0x7f)]
	opcodesMap[shift0xddcb(0x7c)] = opcodesMap[shift0xddcb(0x7f)]
	opcodesMap[shift0xddcb(0x7d)] = opcodesMap[shift0xddcb(0x7f)]
	opcodesMap[shift0xddcb(0x7e)] = opcodesMap[shift0xddcb(0x7f)]
	/* LD B,RES 0,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0x80)] = func(z80 *Z80) {
		z80.b = z80.memory.readByte(z80.tempaddr) & 0xfe
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.b)
	}
	/* LD C,RES 0,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0x81)] = func(z80 *Z80) {
		z80.c = z80.memory.readByte(z80.tempaddr) & 0xfe
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.c)
	}
	/* LD D,RES 0,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0x82)] = func(z80 *Z80) {
		z80.d = z80.memory.readByte(z80.tempaddr) & 0xfe
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.d)
	}
	/* LD E,RES 0,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0x83)] = func(z80 *Z80) {
		z80.e = z80.memory.readByte(z80.tempaddr) & 0xfe
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.e)
	}
	/* LD H,RES 0,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0x84)] = func(z80 *Z80) {
		z80.h = z80.memory.readByte(z80.tempaddr) & 0xfe
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.h)
	}
	/* LD L,RES 0,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0x85)] = func(z80 *Z80) {
		z80.l = z80.memory.readByte(z80.tempaddr) & 0xfe
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.l)
	}
	/* RES 0,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0x86)] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, bytetemp&0xfe)
	}
	/* LD A,RES 0,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0x87)] = func(z80 *Z80) {
		z80.a = z80.memory.readByte(z80.tempaddr) & 0xfe
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.a)
	}
	/* LD B,RES 1,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0x88)] = func(z80 *Z80) {
		z80.b = z80.memory.readByte(z80.tempaddr) & 0xfd
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.b)
	}
	/* LD C,RES 1,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0x89)] = func(z80 *Z80) {
		z80.c = z80.memory.readByte(z80.tempaddr) & 0xfd
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.c)
	}
	/* LD D,RES 1,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0x8a)] = func(z80 *Z80) {
		z80.d = z80.memory.readByte(z80.tempaddr) & 0xfd
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.d)
	}
	/* LD E,RES 1,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0x8b)] = func(z80 *Z80) {
		z80.e = z80.memory.readByte(z80.tempaddr) & 0xfd
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.e)
	}
	/* LD H,RES 1,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0x8c)] = func(z80 *Z80) {
		z80.h = z80.memory.readByte(z80.tempaddr) & 0xfd
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.h)
	}
	/* LD L,RES 1,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0x8d)] = func(z80 *Z80) {
		z80.l = z80.memory.readByte(z80.tempaddr) & 0xfd
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.l)
	}
	/* RES 1,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0x8e)] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, bytetemp&0xfd)
	}
	/* LD A,RES 1,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0x8f)] = func(z80 *Z80) {
		z80.a = z80.memory.readByte(z80.tempaddr) & 0xfd
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.a)
	}
	/* LD B,RES 2,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0x90)] = func(z80 *Z80) {
		z80.b = z80.memory.readByte(z80.tempaddr) & 0xfb
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.b)
	}
	/* LD C,RES 2,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0x91)] = func(z80 *Z80) {
		z80.c = z80.memory.readByte(z80.tempaddr) & 0xfb
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.c)
	}
	/* LD D,RES 2,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0x92)] = func(z80 *Z80) {
		z80.d = z80.memory.readByte(z80.tempaddr) & 0xfb
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.d)
	}
	/* LD E,RES 2,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0x93)] = func(z80 *Z80) {
		z80.e = z80.memory.readByte(z80.tempaddr) & 0xfb
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.e)
	}
	/* LD H,RES 2,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0x94)] = func(z80 *Z80) {
		z80.h = z80.memory.readByte(z80.tempaddr) & 0xfb
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.h)
	}
	/* LD L,RES 2,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0x95)] = func(z80 *Z80) {
		z80.l = z80.memory.readByte(z80.tempaddr) & 0xfb
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.l)
	}
	/* RES 2,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0x96)] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, bytetemp&0xfb)
	}
	/* LD A,RES 2,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0x97)] = func(z80 *Z80) {
		z80.a = z80.memory.readByte(z80.tempaddr) & 0xfb
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.a)
	}
	/* LD B,RES 3,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0x98)] = func(z80 *Z80) {
		z80.b = z80.memory.readByte(z80.tempaddr) & 0xf7
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.b)
	}
	/* LD C,RES 3,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0x99)] = func(z80 *Z80) {
		z80.c = z80.memory.readByte(z80.tempaddr) & 0xf7
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.c)
	}
	/* LD D,RES 3,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0x9a)] = func(z80 *Z80) {
		z80.d = z80.memory.readByte(z80.tempaddr) & 0xf7
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.d)
	}
	/* LD E,RES 3,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0x9b)] = func(z80 *Z80) {
		z80.e = z80.memory.readByte(z80.tempaddr) & 0xf7
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.e)
	}
	/* LD H,RES 3,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0x9c)] = func(z80 *Z80) {
		z80.h = z80.memory.readByte(z80.tempaddr) & 0xf7
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.h)
	}
	/* LD L,RES 3,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0x9d)] = func(z80 *Z80) {
		z80.l = z80.memory.readByte(z80.tempaddr) & 0xf7
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.l)
	}
	/* RES 3,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0x9e)] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, bytetemp&0xf7)
	}
	/* LD A,RES 3,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0x9f)] = func(z80 *Z80) {
		z80.a = z80.memory.readByte(z80.tempaddr) & 0xf7
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.a)
	}
	/* LD B,RES 4,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xa0)] = func(z80 *Z80) {
		z80.b = z80.memory.readByte(z80.tempaddr) & 0xef
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.b)
	}
	/* LD C,RES 4,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xa1)] = func(z80 *Z80) {
		z80.c = z80.memory.readByte(z80.tempaddr) & 0xef
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.c)
	}
	/* LD D,RES 4,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xa2)] = func(z80 *Z80) {
		z80.d = z80.memory.readByte(z80.tempaddr) & 0xef
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.d)
	}
	/* LD E,RES 4,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xa3)] = func(z80 *Z80) {
		z80.e = z80.memory.readByte(z80.tempaddr) & 0xef
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.e)
	}
	/* LD H,RES 4,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xa4)] = func(z80 *Z80) {
		z80.h = z80.memory.readByte(z80.tempaddr) & 0xef
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.h)
	}
	/* LD L,RES 4,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xa5)] = func(z80 *Z80) {
		z80.l = z80.memory.readByte(z80.tempaddr) & 0xef
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.l)
	}
	/* RES 4,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xa6)] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, bytetemp&0xef)
	}
	/* LD A,RES 4,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xa7)] = func(z80 *Z80) {
		z80.a = z80.memory.readByte(z80.tempaddr) & 0xef
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.a)
	}
	/* LD B,RES 5,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xa8)] = func(z80 *Z80) {
		z80.b = z80.memory.readByte(z80.tempaddr) & 0xdf
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.b)
	}
	/* LD C,RES 5,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xa9)] = func(z80 *Z80) {
		z80.c = z80.memory.readByte(z80.tempaddr) & 0xdf
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.c)
	}
	/* LD D,RES 5,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xaa)] = func(z80 *Z80) {
		z80.d = z80.memory.readByte(z80.tempaddr) & 0xdf
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.d)
	}
	/* LD E,RES 5,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xab)] = func(z80 *Z80) {
		z80.e = z80.memory.readByte(z80.tempaddr) & 0xdf
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.e)
	}
	/* LD H,RES 5,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xac)] = func(z80 *Z80) {
		z80.h = z80.memory.readByte(z80.tempaddr) & 0xdf
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.h)
	}
	/* LD L,RES 5,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xad)] = func(z80 *Z80) {
		z80.l = z80.memory.readByte(z80.tempaddr) & 0xdf
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.l)
	}
	/* RES 5,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xae)] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, bytetemp&0xdf)
	}
	/* LD A,RES 5,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xaf)] = func(z80 *Z80) {
		z80.a = z80.memory.readByte(z80.tempaddr) & 0xdf
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.a)
	}
	/* LD B,RES 6,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xb0)] = func(z80 *Z80) {
		z80.b = z80.memory.readByte(z80.tempaddr) & 0xbf
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.b)
	}
	/* LD C,RES 6,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xb1)] = func(z80 *Z80) {
		z80.c = z80.memory.readByte(z80.tempaddr) & 0xbf
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.c)
	}
	/* LD D,RES 6,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xb2)] = func(z80 *Z80) {
		z80.d = z80.memory.readByte(z80.tempaddr) & 0xbf
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.d)
	}
	/* LD E,RES 6,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xb3)] = func(z80 *Z80) {
		z80.e = z80.memory.readByte(z80.tempaddr) & 0xbf
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.e)
	}
	/* LD H,RES 6,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xb4)] = func(z80 *Z80) {
		z80.h = z80.memory.readByte(z80.tempaddr) & 0xbf
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.h)
	}
	/* LD L,RES 6,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xb5)] = func(z80 *Z80) {
		z80.l = z80.memory.readByte(z80.tempaddr) & 0xbf
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.l)
	}
	/* RES 6,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xb6)] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, bytetemp&0xbf)
	}
	/* LD A,RES 6,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xb7)] = func(z80 *Z80) {
		z80.a = z80.memory.readByte(z80.tempaddr) & 0xbf
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.a)
	}
	/* LD B,RES 7,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xb8)] = func(z80 *Z80) {
		z80.b = z80.memory.readByte(z80.tempaddr) & 0x7f
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.b)
	}
	/* LD C,RES 7,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xb9)] = func(z80 *Z80) {
		z80.c = z80.memory.readByte(z80.tempaddr) & 0x7f
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.c)
	}
	/* LD D,RES 7,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xba)] = func(z80 *Z80) {
		z80.d = z80.memory.readByte(z80.tempaddr) & 0x7f
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.d)
	}
	/* LD E,RES 7,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xbb)] = func(z80 *Z80) {
		z80.e = z80.memory.readByte(z80.tempaddr) & 0x7f
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.e)
	}
	/* LD H,RES 7,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xbc)] = func(z80 *Z80) {
		z80.h = z80.memory.readByte(z80.tempaddr) & 0x7f
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.h)
	}
	/* LD L,RES 7,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xbd)] = func(z80 *Z80) {
		z80.l = z80.memory.readByte(z80.tempaddr) & 0x7f
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.l)
	}
	/* RES 7,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xbe)] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, bytetemp&0x7f)
	}
	/* LD A,RES 7,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xbf)] = func(z80 *Z80) {
		z80.a = z80.memory.readByte(z80.tempaddr) & 0x7f
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.a)
	}
	/* LD B,SET 0,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xc0)] = func(z80 *Z80) {
		z80.b = z80.memory.readByte(z80.tempaddr) | 0x01
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.b)
	}
	/* LD C,SET 0,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xc1)] = func(z80 *Z80) {
		z80.c = z80.memory.readByte(z80.tempaddr) | 0x01
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.c)
	}
	/* LD D,SET 0,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xc2)] = func(z80 *Z80) {
		z80.d = z80.memory.readByte(z80.tempaddr) | 0x01
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.d)
	}
	/* LD E,SET 0,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xc3)] = func(z80 *Z80) {
		z80.e = z80.memory.readByte(z80.tempaddr) | 0x01
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.e)
	}
	/* LD H,SET 0,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xc4)] = func(z80 *Z80) {
		z80.h = z80.memory.readByte(z80.tempaddr) | 0x01
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.h)
	}
	/* LD L,SET 0,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xc5)] = func(z80 *Z80) {
		z80.l = z80.memory.readByte(z80.tempaddr) | 0x01
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.l)
	}
	/* SET 0,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xc6)] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, bytetemp|0x01)
	}
	/* LD A,SET 0,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xc7)] = func(z80 *Z80) {
		z80.a = z80.memory.readByte(z80.tempaddr) | 0x01
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.a)
	}
	/* LD B,SET 1,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xc8)] = func(z80 *Z80) {
		z80.b = z80.memory.readByte(z80.tempaddr) | 0x02
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.b)
	}
	/* LD C,SET 1,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xc9)] = func(z80 *Z80) {
		z80.c = z80.memory.readByte(z80.tempaddr) | 0x02
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.c)
	}
	/* LD D,SET 1,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xca)] = func(z80 *Z80) {
		z80.d = z80.memory.readByte(z80.tempaddr) | 0x02
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.d)
	}
	/* LD E,SET 1,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xcb)] = func(z80 *Z80) {
		z80.e = z80.memory.readByte(z80.tempaddr) | 0x02
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.e)
	}
	/* LD H,SET 1,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xcc)] = func(z80 *Z80) {
		z80.h = z80.memory.readByte(z80.tempaddr) | 0x02
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.h)
	}
	/* LD L,SET 1,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xcd)] = func(z80 *Z80) {
		z80.l = z80.memory.readByte(z80.tempaddr) | 0x02
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.l)
	}
	/* SET 1,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xce)] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, bytetemp|0x02)
	}
	/* LD A,SET 1,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xcf)] = func(z80 *Z80) {
		z80.a = z80.memory.readByte(z80.tempaddr) | 0x02
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.a)
	}
	/* LD B,SET 2,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xd0)] = func(z80 *Z80) {
		z80.b = z80.memory.readByte(z80.tempaddr) | 0x04
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.b)
	}
	/* LD C,SET 2,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xd1)] = func(z80 *Z80) {
		z80.c = z80.memory.readByte(z80.tempaddr) | 0x04
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.c)
	}
	/* LD D,SET 2,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xd2)] = func(z80 *Z80) {
		z80.d = z80.memory.readByte(z80.tempaddr) | 0x04
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.d)
	}
	/* LD E,SET 2,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xd3)] = func(z80 *Z80) {
		z80.e = z80.memory.readByte(z80.tempaddr) | 0x04
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.e)
	}
	/* LD H,SET 2,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xd4)] = func(z80 *Z80) {
		z80.h = z80.memory.readByte(z80.tempaddr) | 0x04
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.h)
	}
	/* LD L,SET 2,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xd5)] = func(z80 *Z80) {
		z80.l = z80.memory.readByte(z80.tempaddr) | 0x04
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.l)
	}
	/* SET 2,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xd6)] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, bytetemp|0x04)
	}
	/* LD A,SET 2,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xd7)] = func(z80 *Z80) {
		z80.a = z80.memory.readByte(z80.tempaddr) | 0x04
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.a)
	}
	/* LD B,SET 3,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xd8)] = func(z80 *Z80) {
		z80.b = z80.memory.readByte(z80.tempaddr) | 0x08
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.b)
	}
	/* LD C,SET 3,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xd9)] = func(z80 *Z80) {
		z80.c = z80.memory.readByte(z80.tempaddr) | 0x08
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.c)
	}
	/* LD D,SET 3,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xda)] = func(z80 *Z80) {
		z80.d = z80.memory.readByte(z80.tempaddr) | 0x08
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.d)
	}
	/* LD E,SET 3,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xdb)] = func(z80 *Z80) {
		z80.e = z80.memory.readByte(z80.tempaddr) | 0x08
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.e)
	}
	/* LD H,SET 3,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xdc)] = func(z80 *Z80) {
		z80.h = z80.memory.readByte(z80.tempaddr) | 0x08
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.h)
	}
	/* LD L,SET 3,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xdd)] = func(z80 *Z80) {
		z80.l = z80.memory.readByte(z80.tempaddr) | 0x08
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.l)
	}
	/* SET 3,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xde)] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, bytetemp|0x08)
	}
	/* LD A,SET 3,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xdf)] = func(z80 *Z80) {
		z80.a = z80.memory.readByte(z80.tempaddr) | 0x08
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.a)
	}
	/* LD B,SET 4,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xe0)] = func(z80 *Z80) {
		z80.b = z80.memory.readByte(z80.tempaddr) | 0x10
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.b)
	}
	/* LD C,SET 4,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xe1)] = func(z80 *Z80) {
		z80.c = z80.memory.readByte(z80.tempaddr) | 0x10
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.c)
	}
	/* LD D,SET 4,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xe2)] = func(z80 *Z80) {
		z80.d = z80.memory.readByte(z80.tempaddr) | 0x10
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.d)
	}
	/* LD E,SET 4,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xe3)] = func(z80 *Z80) {
		z80.e = z80.memory.readByte(z80.tempaddr) | 0x10
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.e)
	}
	/* LD H,SET 4,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xe4)] = func(z80 *Z80) {
		z80.h = z80.memory.readByte(z80.tempaddr) | 0x10
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.h)
	}
	/* LD L,SET 4,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xe5)] = func(z80 *Z80) {
		z80.l = z80.memory.readByte(z80.tempaddr) | 0x10
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.l)
	}
	/* SET 4,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xe6)] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, bytetemp|0x10)
	}
	/* LD A,SET 4,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xe7)] = func(z80 *Z80) {
		z80.a = z80.memory.readByte(z80.tempaddr) | 0x10
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.a)
	}
	/* LD B,SET 5,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xe8)] = func(z80 *Z80) {
		z80.b = z80.memory.readByte(z80.tempaddr) | 0x20
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.b)
	}
	/* LD C,SET 5,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xe9)] = func(z80 *Z80) {
		z80.c = z80.memory.readByte(z80.tempaddr) | 0x20
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.c)
	}
	/* LD D,SET 5,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xea)] = func(z80 *Z80) {
		z80.d = z80.memory.readByte(z80.tempaddr) | 0x20
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.d)
	}
	/* LD E,SET 5,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xeb)] = func(z80 *Z80) {
		z80.e = z80.memory.readByte(z80.tempaddr) | 0x20
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.e)
	}
	/* LD H,SET 5,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xec)] = func(z80 *Z80) {
		z80.h = z80.memory.readByte(z80.tempaddr) | 0x20
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.h)
	}
	/* LD L,SET 5,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xed)] = func(z80 *Z80) {
		z80.l = z80.memory.readByte(z80.tempaddr) | 0x20
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.l)
	}
	/* SET 5,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xee)] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, bytetemp|0x20)
	}
	/* LD A,SET 5,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xef)] = func(z80 *Z80) {
		z80.a = z80.memory.readByte(z80.tempaddr) | 0x20
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.a)
	}
	/* LD B,SET 6,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xf0)] = func(z80 *Z80) {
		z80.b = z80.memory.readByte(z80.tempaddr) | 0x40
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.b)
	}
	/* LD C,SET 6,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xf1)] = func(z80 *Z80) {
		z80.c = z80.memory.readByte(z80.tempaddr) | 0x40
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.c)
	}
	/* LD D,SET 6,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xf2)] = func(z80 *Z80) {
		z80.d = z80.memory.readByte(z80.tempaddr) | 0x40
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.d)
	}
	/* LD E,SET 6,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xf3)] = func(z80 *Z80) {
		z80.e = z80.memory.readByte(z80.tempaddr) | 0x40
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.e)
	}
	/* LD H,SET 6,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xf4)] = func(z80 *Z80) {
		z80.h = z80.memory.readByte(z80.tempaddr) | 0x40
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.h)
	}
	/* LD L,SET 6,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xf5)] = func(z80 *Z80) {
		z80.l = z80.memory.readByte(z80.tempaddr) | 0x40
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.l)
	}
	/* SET 6,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xf6)] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, bytetemp|0x40)
	}
	/* LD A,SET 6,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xf7)] = func(z80 *Z80) {
		z80.a = z80.memory.readByte(z80.tempaddr) | 0x40
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.a)
	}
	/* LD B,SET 7,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xf8)] = func(z80 *Z80) {
		z80.b = z80.memory.readByte(z80.tempaddr) | 0x80
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.b)
	}
	/* LD C,SET 7,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xf9)] = func(z80 *Z80) {
		z80.c = z80.memory.readByte(z80.tempaddr) | 0x80
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.c)
	}
	/* LD D,SET 7,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xfa)] = func(z80 *Z80) {
		z80.d = z80.memory.readByte(z80.tempaddr) | 0x80
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.d)
	}
	/* LD E,SET 7,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xfb)] = func(z80 *Z80) {
		z80.e = z80.memory.readByte(z80.tempaddr) | 0x80
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.e)
	}
	/* LD H,SET 7,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xfc)] = func(z80 *Z80) {
		z80.h = z80.memory.readByte(z80.tempaddr) | 0x80
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.h)
	}
	/* LD L,SET 7,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xfd)] = func(z80 *Z80) {
		z80.l = z80.memory.readByte(z80.tempaddr) | 0x80
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.l)
	}
	/* SET 7,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xfe)] = func(z80 *Z80) {
		var bytetemp byte = z80.memory.readByte(z80.tempaddr)
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, bytetemp|0x80)
	}
	/* LD A,SET 7,(REGISTER+dd) */
	opcodesMap[shift0xddcb(0xff)] = func(z80 *Z80) {
		z80.a = z80.memory.readByte(z80.tempaddr) | 0x80
		z80.memory.contendReadNoMreq(z80.tempaddr, 1)
		z80.memory.writeByte(z80.tempaddr, z80.a)
	}

	// END of 0xddfdcb shifted opcodes

}
