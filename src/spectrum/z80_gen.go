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

package spectrum

//
// Automatically generated file -- DO NOT EDIT
//

// Generated getters and INC/DEC functions for 8bit registers

func (z80 *Z80) A() byte {
	return z80.a
}

func (z80 *Z80) incA() {
	z80.a++
	z80.f = (z80.f & FLAG_C) | (ternOpB(z80.a == 0x80, FLAG_V, 0)) | (ternOpB((z80.a&0x0f) != 0, 0, FLAG_H)) | sz53Table[z80.a]
}

func (z80 *Z80) decA() {
	z80.f = (z80.f & FLAG_C) | (ternOpB(z80.a&0x0f != 0, 0, FLAG_H)) | FLAG_N
	z80.a--
	z80.f |= (ternOpB(z80.a == 0x7f, FLAG_V, 0)) | sz53Table[z80.a]
}

func (z80 *Z80) B() byte {
	return z80.b
}

func (z80 *Z80) incB() {
	z80.b++
	z80.f = (z80.f & FLAG_C) | (ternOpB(z80.b == 0x80, FLAG_V, 0)) | (ternOpB((z80.b&0x0f) != 0, 0, FLAG_H)) | sz53Table[z80.b]
}

func (z80 *Z80) decB() {
	z80.f = (z80.f & FLAG_C) | (ternOpB(z80.b&0x0f != 0, 0, FLAG_H)) | FLAG_N
	z80.b--
	z80.f |= (ternOpB(z80.b == 0x7f, FLAG_V, 0)) | sz53Table[z80.b]
}

func (z80 *Z80) C() byte {
	return z80.c
}

func (z80 *Z80) incC() {
	z80.c++
	z80.f = (z80.f & FLAG_C) | (ternOpB(z80.c == 0x80, FLAG_V, 0)) | (ternOpB((z80.c&0x0f) != 0, 0, FLAG_H)) | sz53Table[z80.c]
}

func (z80 *Z80) decC() {
	z80.f = (z80.f & FLAG_C) | (ternOpB(z80.c&0x0f != 0, 0, FLAG_H)) | FLAG_N
	z80.c--
	z80.f |= (ternOpB(z80.c == 0x7f, FLAG_V, 0)) | sz53Table[z80.c]
}

func (z80 *Z80) D() byte {
	return z80.d
}

func (z80 *Z80) incD() {
	z80.d++
	z80.f = (z80.f & FLAG_C) | (ternOpB(z80.d == 0x80, FLAG_V, 0)) | (ternOpB((z80.d&0x0f) != 0, 0, FLAG_H)) | sz53Table[z80.d]
}

func (z80 *Z80) decD() {
	z80.f = (z80.f & FLAG_C) | (ternOpB(z80.d&0x0f != 0, 0, FLAG_H)) | FLAG_N
	z80.d--
	z80.f |= (ternOpB(z80.d == 0x7f, FLAG_V, 0)) | sz53Table[z80.d]
}

func (z80 *Z80) E() byte {
	return z80.e
}

func (z80 *Z80) incE() {
	z80.e++
	z80.f = (z80.f & FLAG_C) | (ternOpB(z80.e == 0x80, FLAG_V, 0)) | (ternOpB((z80.e&0x0f) != 0, 0, FLAG_H)) | sz53Table[z80.e]
}

func (z80 *Z80) decE() {
	z80.f = (z80.f & FLAG_C) | (ternOpB(z80.e&0x0f != 0, 0, FLAG_H)) | FLAG_N
	z80.e--
	z80.f |= (ternOpB(z80.e == 0x7f, FLAG_V, 0)) | sz53Table[z80.e]
}

func (z80 *Z80) H() byte {
	return z80.h
}

func (z80 *Z80) incH() {
	z80.h++
	z80.f = (z80.f & FLAG_C) | (ternOpB(z80.h == 0x80, FLAG_V, 0)) | (ternOpB((z80.h&0x0f) != 0, 0, FLAG_H)) | sz53Table[z80.h]
}

func (z80 *Z80) decH() {
	z80.f = (z80.f & FLAG_C) | (ternOpB(z80.h&0x0f != 0, 0, FLAG_H)) | FLAG_N
	z80.h--
	z80.f |= (ternOpB(z80.h == 0x7f, FLAG_V, 0)) | sz53Table[z80.h]
}

func (z80 *Z80) L() byte {
	return z80.l
}

func (z80 *Z80) incL() {
	z80.l++
	z80.f = (z80.f & FLAG_C) | (ternOpB(z80.l == 0x80, FLAG_V, 0)) | (ternOpB((z80.l&0x0f) != 0, 0, FLAG_H)) | sz53Table[z80.l]
}

func (z80 *Z80) decL() {
	z80.f = (z80.f & FLAG_C) | (ternOpB(z80.l&0x0f != 0, 0, FLAG_H)) | FLAG_N
	z80.l--
	z80.f |= (ternOpB(z80.l == 0x7f, FLAG_V, 0)) | sz53Table[z80.l]
}

func (z80 *Z80) IXL() byte {
	return z80.ixl
}

func (z80 *Z80) incIXL() {
	z80.ixl++
	z80.f = (z80.f & FLAG_C) | (ternOpB(z80.ixl == 0x80, FLAG_V, 0)) | (ternOpB((z80.ixl&0x0f) != 0, 0, FLAG_H)) | sz53Table[z80.ixl]
}

func (z80 *Z80) decIXL() {
	z80.f = (z80.f & FLAG_C) | (ternOpB(z80.ixl&0x0f != 0, 0, FLAG_H)) | FLAG_N
	z80.ixl--
	z80.f |= (ternOpB(z80.ixl == 0x7f, FLAG_V, 0)) | sz53Table[z80.ixl]
}

func (z80 *Z80) IXH() byte {
	return z80.ixh
}

func (z80 *Z80) incIXH() {
	z80.ixh++
	z80.f = (z80.f & FLAG_C) | (ternOpB(z80.ixh == 0x80, FLAG_V, 0)) | (ternOpB((z80.ixh&0x0f) != 0, 0, FLAG_H)) | sz53Table[z80.ixh]
}

func (z80 *Z80) decIXH() {
	z80.f = (z80.f & FLAG_C) | (ternOpB(z80.ixh&0x0f != 0, 0, FLAG_H)) | FLAG_N
	z80.ixh--
	z80.f |= (ternOpB(z80.ixh == 0x7f, FLAG_V, 0)) | sz53Table[z80.ixh]
}

func (z80 *Z80) IYL() byte {
	return z80.iyl
}

func (z80 *Z80) incIYL() {
	z80.iyl++
	z80.f = (z80.f & FLAG_C) | (ternOpB(z80.iyl == 0x80, FLAG_V, 0)) | (ternOpB((z80.iyl&0x0f) != 0, 0, FLAG_H)) | sz53Table[z80.iyl]
}

func (z80 *Z80) decIYL() {
	z80.f = (z80.f & FLAG_C) | (ternOpB(z80.iyl&0x0f != 0, 0, FLAG_H)) | FLAG_N
	z80.iyl--
	z80.f |= (ternOpB(z80.iyl == 0x7f, FLAG_V, 0)) | sz53Table[z80.iyl]
}

func (z80 *Z80) IYH() byte {
	return z80.iyh
}

func (z80 *Z80) incIYH() {
	z80.iyh++
	z80.f = (z80.f & FLAG_C) | (ternOpB(z80.iyh == 0x80, FLAG_V, 0)) | (ternOpB((z80.iyh&0x0f) != 0, 0, FLAG_H)) | sz53Table[z80.iyh]
}

func (z80 *Z80) decIYH() {
	z80.f = (z80.f & FLAG_C) | (ternOpB(z80.iyh&0x0f != 0, 0, FLAG_H)) | FLAG_N
	z80.iyh--
	z80.f |= (ternOpB(z80.iyh == 0x7f, FLAG_V, 0)) | sz53Table[z80.iyh]
}

// Generated getters/setters and INC/DEC functions for 16bit registers

func (z80 *Z80) BC() uint16 {
	return z80.bc.get()
}

func (z80 *Z80) setBC(value uint16) {
	z80.bc.set(value)
}

func (z80 *Z80) decBC() {
	z80.bc.dec()
}

func (z80 *Z80) incBC() {
	z80.bc.inc()
}

func (z80 *Z80) DE() uint16 {
	return z80.de.get()
}

func (z80 *Z80) setDE(value uint16) {
	z80.de.set(value)
}

func (z80 *Z80) decDE() {
	z80.de.dec()
}

func (z80 *Z80) incDE() {
	z80.de.inc()
}

func (z80 *Z80) HL() uint16 {
	return z80.hl.get()
}

func (z80 *Z80) setHL(value uint16) {
	z80.hl.set(value)
}

func (z80 *Z80) decHL() {
	z80.hl.dec()
}

func (z80 *Z80) incHL() {
	z80.hl.inc()
}

func (z80 *Z80) BC_() uint16 {
	return z80.bc_.get()
}

func (z80 *Z80) setBC_(value uint16) {
	z80.bc_.set(value)
}

func (z80 *Z80) decBC_() {
	z80.bc_.dec()
}

func (z80 *Z80) incBC_() {
	z80.bc_.inc()
}

func (z80 *Z80) DE_() uint16 {
	return z80.de_.get()
}

func (z80 *Z80) setDE_(value uint16) {
	z80.de_.set(value)
}

func (z80 *Z80) decDE_() {
	z80.de_.dec()
}

func (z80 *Z80) incDE_() {
	z80.de_.inc()
}

func (z80 *Z80) HL_() uint16 {
	return z80.hl_.get()
}

func (z80 *Z80) setHL_(value uint16) {
	z80.hl_.set(value)
}

func (z80 *Z80) decHL_() {
	z80.hl_.dec()
}

func (z80 *Z80) incHL_() {
	z80.hl_.inc()
}

func (z80 *Z80) IX() uint16 {
	return z80.ix.get()
}

func (z80 *Z80) setIX(value uint16) {
	z80.ix.set(value)
}

func (z80 *Z80) decIX() {
	z80.ix.dec()
}

func (z80 *Z80) incIX() {
	z80.ix.inc()
}

func (z80 *Z80) IY() uint16 {
	return z80.iy.get()
}

func (z80 *Z80) setIY(value uint16) {
	z80.iy.set(value)
}

func (z80 *Z80) decIY() {
	z80.iy.dec()
}

func (z80 *Z80) incIY() {
	z80.iy.inc()
}
