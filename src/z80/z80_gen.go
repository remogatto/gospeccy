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

package z80

//
// Automatically generated file -- DO NOT EDIT
//

// Generated INC/DEC functions for 8bit registers

func (z80 *Z80) incA() {
	z80.A++
	z80.F = (z80.F & FLAG_C) | (ternOpB(z80.A == 0x80, FLAG_V, 0)) | (ternOpB((z80.A&0x0f) != 0, 0, FLAG_H)) | sz53Table[z80.A]
}

func (z80 *Z80) decA() {
	z80.F = (z80.F & FLAG_C) | (ternOpB(z80.A&0x0f != 0, 0, FLAG_H)) | FLAG_N
	z80.A--
	z80.F |= (ternOpB(z80.A == 0x7f, FLAG_V, 0)) | sz53Table[z80.A]
}

func (z80 *Z80) incB() {
	z80.B++
	z80.F = (z80.F & FLAG_C) | (ternOpB(z80.B == 0x80, FLAG_V, 0)) | (ternOpB((z80.B&0x0f) != 0, 0, FLAG_H)) | sz53Table[z80.B]
}

func (z80 *Z80) decB() {
	z80.F = (z80.F & FLAG_C) | (ternOpB(z80.B&0x0f != 0, 0, FLAG_H)) | FLAG_N
	z80.B--
	z80.F |= (ternOpB(z80.B == 0x7f, FLAG_V, 0)) | sz53Table[z80.B]
}

func (z80 *Z80) incC() {
	z80.C++
	z80.F = (z80.F & FLAG_C) | (ternOpB(z80.C == 0x80, FLAG_V, 0)) | (ternOpB((z80.C&0x0f) != 0, 0, FLAG_H)) | sz53Table[z80.C]
}

func (z80 *Z80) decC() {
	z80.F = (z80.F & FLAG_C) | (ternOpB(z80.C&0x0f != 0, 0, FLAG_H)) | FLAG_N
	z80.C--
	z80.F |= (ternOpB(z80.C == 0x7f, FLAG_V, 0)) | sz53Table[z80.C]
}

func (z80 *Z80) incD() {
	z80.D++
	z80.F = (z80.F & FLAG_C) | (ternOpB(z80.D == 0x80, FLAG_V, 0)) | (ternOpB((z80.D&0x0f) != 0, 0, FLAG_H)) | sz53Table[z80.D]
}

func (z80 *Z80) decD() {
	z80.F = (z80.F & FLAG_C) | (ternOpB(z80.D&0x0f != 0, 0, FLAG_H)) | FLAG_N
	z80.D--
	z80.F |= (ternOpB(z80.D == 0x7f, FLAG_V, 0)) | sz53Table[z80.D]
}

func (z80 *Z80) incE() {
	z80.E++
	z80.F = (z80.F & FLAG_C) | (ternOpB(z80.E == 0x80, FLAG_V, 0)) | (ternOpB((z80.E&0x0f) != 0, 0, FLAG_H)) | sz53Table[z80.E]
}

func (z80 *Z80) decE() {
	z80.F = (z80.F & FLAG_C) | (ternOpB(z80.E&0x0f != 0, 0, FLAG_H)) | FLAG_N
	z80.E--
	z80.F |= (ternOpB(z80.E == 0x7f, FLAG_V, 0)) | sz53Table[z80.E]
}

func (z80 *Z80) incF() {
	z80.F++
	z80.F = (z80.F & FLAG_C) | (ternOpB(z80.F == 0x80, FLAG_V, 0)) | (ternOpB((z80.F&0x0f) != 0, 0, FLAG_H)) | sz53Table[z80.F]
}

func (z80 *Z80) decF() {
	z80.F = (z80.F & FLAG_C) | (ternOpB(z80.F&0x0f != 0, 0, FLAG_H)) | FLAG_N
	z80.F--
	z80.F |= (ternOpB(z80.F == 0x7f, FLAG_V, 0)) | sz53Table[z80.F]
}

func (z80 *Z80) incH() {
	z80.H++
	z80.F = (z80.F & FLAG_C) | (ternOpB(z80.H == 0x80, FLAG_V, 0)) | (ternOpB((z80.H&0x0f) != 0, 0, FLAG_H)) | sz53Table[z80.H]
}

func (z80 *Z80) decH() {
	z80.F = (z80.F & FLAG_C) | (ternOpB(z80.H&0x0f != 0, 0, FLAG_H)) | FLAG_N
	z80.H--
	z80.F |= (ternOpB(z80.H == 0x7f, FLAG_V, 0)) | sz53Table[z80.H]
}

func (z80 *Z80) incI() {
	z80.I++
	z80.F = (z80.F & FLAG_C) | (ternOpB(z80.I == 0x80, FLAG_V, 0)) | (ternOpB((z80.I&0x0f) != 0, 0, FLAG_H)) | sz53Table[z80.I]
}

func (z80 *Z80) decI() {
	z80.F = (z80.F & FLAG_C) | (ternOpB(z80.I&0x0f != 0, 0, FLAG_H)) | FLAG_N
	z80.I--
	z80.F |= (ternOpB(z80.I == 0x7f, FLAG_V, 0)) | sz53Table[z80.I]
}

func (z80 *Z80) incL() {
	z80.L++
	z80.F = (z80.F & FLAG_C) | (ternOpB(z80.L == 0x80, FLAG_V, 0)) | (ternOpB((z80.L&0x0f) != 0, 0, FLAG_H)) | sz53Table[z80.L]
}

func (z80 *Z80) decL() {
	z80.F = (z80.F & FLAG_C) | (ternOpB(z80.L&0x0f != 0, 0, FLAG_H)) | FLAG_N
	z80.L--
	z80.F |= (ternOpB(z80.L == 0x7f, FLAG_V, 0)) | sz53Table[z80.L]
}

func (z80 *Z80) incR7() {
	z80.R7++
	z80.F = (z80.F & FLAG_C) | (ternOpB(z80.R7 == 0x80, FLAG_V, 0)) | (ternOpB((z80.R7&0x0f) != 0, 0, FLAG_H)) | sz53Table[z80.R7]
}

func (z80 *Z80) decR7() {
	z80.F = (z80.F & FLAG_C) | (ternOpB(z80.R7&0x0f != 0, 0, FLAG_H)) | FLAG_N
	z80.R7--
	z80.F |= (ternOpB(z80.R7 == 0x7f, FLAG_V, 0)) | sz53Table[z80.R7]
}

func (z80 *Z80) incA_() {
	z80.A_++
	z80.F = (z80.F & FLAG_C) | (ternOpB(z80.A_ == 0x80, FLAG_V, 0)) | (ternOpB((z80.A_&0x0f) != 0, 0, FLAG_H)) | sz53Table[z80.A_]
}

func (z80 *Z80) decA_() {
	z80.F = (z80.F & FLAG_C) | (ternOpB(z80.A_&0x0f != 0, 0, FLAG_H)) | FLAG_N
	z80.A_--
	z80.F |= (ternOpB(z80.A_ == 0x7f, FLAG_V, 0)) | sz53Table[z80.A_]
}

func (z80 *Z80) incB_() {
	z80.B_++
	z80.F = (z80.F & FLAG_C) | (ternOpB(z80.B_ == 0x80, FLAG_V, 0)) | (ternOpB((z80.B_&0x0f) != 0, 0, FLAG_H)) | sz53Table[z80.B_]
}

func (z80 *Z80) decB_() {
	z80.F = (z80.F & FLAG_C) | (ternOpB(z80.B_&0x0f != 0, 0, FLAG_H)) | FLAG_N
	z80.B_--
	z80.F |= (ternOpB(z80.B_ == 0x7f, FLAG_V, 0)) | sz53Table[z80.B_]
}

func (z80 *Z80) incC_() {
	z80.C_++
	z80.F = (z80.F & FLAG_C) | (ternOpB(z80.C_ == 0x80, FLAG_V, 0)) | (ternOpB((z80.C_&0x0f) != 0, 0, FLAG_H)) | sz53Table[z80.C_]
}

func (z80 *Z80) decC_() {
	z80.F = (z80.F & FLAG_C) | (ternOpB(z80.C_&0x0f != 0, 0, FLAG_H)) | FLAG_N
	z80.C_--
	z80.F |= (ternOpB(z80.C_ == 0x7f, FLAG_V, 0)) | sz53Table[z80.C_]
}

func (z80 *Z80) incD_() {
	z80.D_++
	z80.F = (z80.F & FLAG_C) | (ternOpB(z80.D_ == 0x80, FLAG_V, 0)) | (ternOpB((z80.D_&0x0f) != 0, 0, FLAG_H)) | sz53Table[z80.D_]
}

func (z80 *Z80) decD_() {
	z80.F = (z80.F & FLAG_C) | (ternOpB(z80.D_&0x0f != 0, 0, FLAG_H)) | FLAG_N
	z80.D_--
	z80.F |= (ternOpB(z80.D_ == 0x7f, FLAG_V, 0)) | sz53Table[z80.D_]
}

func (z80 *Z80) incE_() {
	z80.E_++
	z80.F = (z80.F & FLAG_C) | (ternOpB(z80.E_ == 0x80, FLAG_V, 0)) | (ternOpB((z80.E_&0x0f) != 0, 0, FLAG_H)) | sz53Table[z80.E_]
}

func (z80 *Z80) decE_() {
	z80.F = (z80.F & FLAG_C) | (ternOpB(z80.E_&0x0f != 0, 0, FLAG_H)) | FLAG_N
	z80.E_--
	z80.F |= (ternOpB(z80.E_ == 0x7f, FLAG_V, 0)) | sz53Table[z80.E_]
}

func (z80 *Z80) incF_() {
	z80.F_++
	z80.F = (z80.F & FLAG_C) | (ternOpB(z80.F_ == 0x80, FLAG_V, 0)) | (ternOpB((z80.F_&0x0f) != 0, 0, FLAG_H)) | sz53Table[z80.F_]
}

func (z80 *Z80) decF_() {
	z80.F = (z80.F & FLAG_C) | (ternOpB(z80.F_&0x0f != 0, 0, FLAG_H)) | FLAG_N
	z80.F_--
	z80.F |= (ternOpB(z80.F_ == 0x7f, FLAG_V, 0)) | sz53Table[z80.F_]
}

func (z80 *Z80) incH_() {
	z80.H_++
	z80.F = (z80.F & FLAG_C) | (ternOpB(z80.H_ == 0x80, FLAG_V, 0)) | (ternOpB((z80.H_&0x0f) != 0, 0, FLAG_H)) | sz53Table[z80.H_]
}

func (z80 *Z80) decH_() {
	z80.F = (z80.F & FLAG_C) | (ternOpB(z80.H_&0x0f != 0, 0, FLAG_H)) | FLAG_N
	z80.H_--
	z80.F |= (ternOpB(z80.H_ == 0x7f, FLAG_V, 0)) | sz53Table[z80.H_]
}

func (z80 *Z80) incL_() {
	z80.L_++
	z80.F = (z80.F & FLAG_C) | (ternOpB(z80.L_ == 0x80, FLAG_V, 0)) | (ternOpB((z80.L_&0x0f) != 0, 0, FLAG_H)) | sz53Table[z80.L_]
}

func (z80 *Z80) decL_() {
	z80.F = (z80.F & FLAG_C) | (ternOpB(z80.L_&0x0f != 0, 0, FLAG_H)) | FLAG_N
	z80.L_--
	z80.F |= (ternOpB(z80.L_ == 0x7f, FLAG_V, 0)) | sz53Table[z80.L_]
}

func (z80 *Z80) incIXL() {
	z80.IXL++
	z80.F = (z80.F & FLAG_C) | (ternOpB(z80.IXL == 0x80, FLAG_V, 0)) | (ternOpB((z80.IXL&0x0f) != 0, 0, FLAG_H)) | sz53Table[z80.IXL]
}

func (z80 *Z80) decIXL() {
	z80.F = (z80.F & FLAG_C) | (ternOpB(z80.IXL&0x0f != 0, 0, FLAG_H)) | FLAG_N
	z80.IXL--
	z80.F |= (ternOpB(z80.IXL == 0x7f, FLAG_V, 0)) | sz53Table[z80.IXL]
}

func (z80 *Z80) incIXH() {
	z80.IXH++
	z80.F = (z80.F & FLAG_C) | (ternOpB(z80.IXH == 0x80, FLAG_V, 0)) | (ternOpB((z80.IXH&0x0f) != 0, 0, FLAG_H)) | sz53Table[z80.IXH]
}

func (z80 *Z80) decIXH() {
	z80.F = (z80.F & FLAG_C) | (ternOpB(z80.IXH&0x0f != 0, 0, FLAG_H)) | FLAG_N
	z80.IXH--
	z80.F |= (ternOpB(z80.IXH == 0x7f, FLAG_V, 0)) | sz53Table[z80.IXH]
}

func (z80 *Z80) incIYL() {
	z80.IYL++
	z80.F = (z80.F & FLAG_C) | (ternOpB(z80.IYL == 0x80, FLAG_V, 0)) | (ternOpB((z80.IYL&0x0f) != 0, 0, FLAG_H)) | sz53Table[z80.IYL]
}

func (z80 *Z80) decIYL() {
	z80.F = (z80.F & FLAG_C) | (ternOpB(z80.IYL&0x0f != 0, 0, FLAG_H)) | FLAG_N
	z80.IYL--
	z80.F |= (ternOpB(z80.IYL == 0x7f, FLAG_V, 0)) | sz53Table[z80.IYL]
}

func (z80 *Z80) incIYH() {
	z80.IYH++
	z80.F = (z80.F & FLAG_C) | (ternOpB(z80.IYH == 0x80, FLAG_V, 0)) | (ternOpB((z80.IYH&0x0f) != 0, 0, FLAG_H)) | sz53Table[z80.IYH]
}

func (z80 *Z80) decIYH() {
	z80.F = (z80.F & FLAG_C) | (ternOpB(z80.IYH&0x0f != 0, 0, FLAG_H)) | FLAG_N
	z80.IYH--
	z80.F |= (ternOpB(z80.IYH == 0x7f, FLAG_V, 0)) | sz53Table[z80.IYH]
}

// Generated getters/setters and INC/DEC functions for 16bit registers

func (z80 *Z80) BC() uint16 {
	return z80.bc.get()
}

func (z80 *Z80) SetBC(value uint16) {
	z80.bc.set(value)
}

func (z80 *Z80) DecBC() {
	z80.bc.dec()
}

func (z80 *Z80) IncBC() {
	z80.bc.inc()
}

func (z80 *Z80) DE() uint16 {
	return z80.de.get()
}

func (z80 *Z80) SetDE(value uint16) {
	z80.de.set(value)
}

func (z80 *Z80) DecDE() {
	z80.de.dec()
}

func (z80 *Z80) IncDE() {
	z80.de.inc()
}

func (z80 *Z80) HL() uint16 {
	return z80.hl.get()
}

func (z80 *Z80) SetHL(value uint16) {
	z80.hl.set(value)
}

func (z80 *Z80) DecHL() {
	z80.hl.dec()
}

func (z80 *Z80) IncHL() {
	z80.hl.inc()
}

func (z80 *Z80) BC_() uint16 {
	return z80.bc_.get()
}

func (z80 *Z80) SetBC_(value uint16) {
	z80.bc_.set(value)
}

func (z80 *Z80) DecBC_() {
	z80.bc_.dec()
}

func (z80 *Z80) IncBC_() {
	z80.bc_.inc()
}

func (z80 *Z80) DE_() uint16 {
	return z80.de_.get()
}

func (z80 *Z80) SetDE_(value uint16) {
	z80.de_.set(value)
}

func (z80 *Z80) DecDE_() {
	z80.de_.dec()
}

func (z80 *Z80) IncDE_() {
	z80.de_.inc()
}

func (z80 *Z80) HL_() uint16 {
	return z80.hl_.get()
}

func (z80 *Z80) SetHL_(value uint16) {
	z80.hl_.set(value)
}

func (z80 *Z80) DecHL_() {
	z80.hl_.dec()
}

func (z80 *Z80) IncHL_() {
	z80.hl_.inc()
}

func (z80 *Z80) IX() uint16 {
	return z80.ix.get()
}

func (z80 *Z80) SetIX(value uint16) {
	z80.ix.set(value)
}

func (z80 *Z80) DecIX() {
	z80.ix.dec()
}

func (z80 *Z80) IncIX() {
	z80.ix.inc()
}

func (z80 *Z80) IY() uint16 {
	return z80.iy.get()
}

func (z80 *Z80) SetIY(value uint16) {
	z80.iy.set(value)
}

func (z80 *Z80) DecIY() {
	z80.iy.dec()
}

func (z80 *Z80) IncIY() {
	z80.iy.inc()
}
