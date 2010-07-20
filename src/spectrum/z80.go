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

import (
	"io/ioutil"
	"container/vector"
	"fmt"
	"os"
)

/* The flags */

const FLAG_C = 0x01
const FLAG_N = 0x02
const FLAG_P = 0x04
const FLAG_V = FLAG_P
const FLAG_3 = 0x08
const FLAG_H = 0x10
const FLAG_5 = 0x20
const FLAG_Z = 0x40
const FLAG_S = 0x80

var z80interruptEvent int

/* Whether a half carry occurred or not can be determined by looking at
   the 3rd bit of the two arguments and the result; these are hashed
   into this table in the form r12, where r is the 3rd bit of the
   result, 1 is the 3rd bit of the 1st argument and 2 is the
   third bit of the 2nd argument; the tables differ for add and subtract
   operations */
var halfcarryAddTable = []byte{0, FLAG_H, FLAG_H, FLAG_H, 0, 0, 0, FLAG_H}
var halfcarrySubTable = []byte{0, 0, FLAG_H, 0, FLAG_H, 0, FLAG_H, FLAG_H}

/* Similarly, overflow can be determined by looking at the 7th bits; again
   the hash into this table is r12 */
var overflowAddTable = []byte{0, 0, 0, FLAG_V, FLAG_V, 0, 0, 0}
var overflowSubTable = []byte{0, FLAG_V, 0, 0, 0, 0, FLAG_V, 0}

var rzxInstructionsOffset int

type register16 struct {
	high, low *byte
}

func (r *register16) inc() {
	temp := r.get() + 1
	*r.high = byte(temp >> 8)
	*r.low = byte(temp & 0xff)
}

func (r *register16) dec() {
	temp := r.get() - 1
	*r.high = byte(temp >> 8)
	*r.low = byte(temp & 0xff)
}

func (r *register16) set(value uint16) {
	*r.high, *r.low = splitWord(value)
}

func (r *register16) get() uint16 {
	return joinBytes(*r.high, *r.low)
}

type Z80 struct {
	a, f, b, c, d, e, h, i, l      byte
	a_, f_, b_, c_, d_, e_, h_, l_ byte
	ixh, ixl, iyh, iyl             byte

	sp, r, r7, pc, iff1, iff2, im uint16

	bc, bc_, hl, hl_, af, de, de_, ix, iy register16

	sz53Table, sz53pTable, parityTable [0x100]byte

	// Number of tstates since the beginning of the last frame
	tstates uint

	halted bool

	interruptsEnabledAt int

	memory MemoryAccessor

	PortAccessor

	LogEvents bool
}

var initialMemory [0x10000]byte
var eventNextEvent uint

func NewZ80(memory MemoryAccessor, port PortAccessor) *Z80 {
	z80 := &Z80{memory: memory, PortAccessor: port}

	z80.bc = register16{&z80.b, &z80.c}
	z80.bc_ = register16{&z80.b_, &z80.c_}
	z80.hl = register16{&z80.h, &z80.l}
	z80.hl_ = register16{&z80.h_, &z80.l_}
	z80.af = register16{&z80.a, &z80.f}
	z80.de = register16{&z80.d, &z80.e}
	z80.ix = register16{&z80.ixh, &z80.ixl}
	z80.iy = register16{&z80.iyh, &z80.iyl}
	z80.de_ = register16{&z80.d_, &z80.e_}

	z80.initTables()

	return z80
}

func (z80 *Z80) DumpRegisters(out *vector.StringVector) {
	out.Push(fmt.Sprintf("%02x%02x %02x%02x %02x%02x %02x%02x %02x%02x %02x%02x %02x%02x %02x%02x %02x%02x %02x%02x %04x %04x\n",
		z80.a, z80.f, z80.b, z80.c, z80.d, z80.e, z80.h, z80.l, z80.a_, z80.f_, z80.b_, z80.c_, z80.d_, z80.e_, z80.h_, z80.l_, z80.ixh, z80.ixl, z80.iyh, z80.iyl, z80.sp, z80.pc))
	out.Push(fmt.Sprintf("%02x %02x %d %d %d %d %d\n", z80.i, (z80.r7&0x80)|(z80.r&0x7f),
		z80.iff1, z80.iff2, z80.im, z80.halted, z80.tstates))
}

func (z80 *Z80) DumpMemory(out *vector.StringVector) {
	var i uint
	for i = 0; i < 0x10000; i++ {
		if z80.memory.At(i) == initialMemory[i] {
			continue
		}

		line := fmt.Sprintf("%04x ", i)

		for (i < 0x10000) && (z80.memory.At(i) != initialMemory[i]) {
			line += fmt.Sprintf("%02x ", z80.memory.At(i))
			i++
		}

		line += fmt.Sprintf("-1\n")

		out.Push(line)
	}
}

func (z80 *Z80) Reset() {
	z80.a, z80.f, z80.b, z80.c, z80.d, z80.e, z80.h, z80.l = 0, 0, 0, 0, 0, 0, 0, 0
	z80.a_, z80.f_, z80.b_, z80.c_, z80.d_, z80.e_, z80.h_, z80.l_ = 0, 0, 0, 0, 0, 0, 0, 0
	z80.ixh, z80.ixl, z80.iyh, z80.iyl = 0, 0, 0, 0

	z80.sp, z80.i, z80.r, z80.r7, z80.pc, z80.iff1, z80.iff2, z80.im = 0, 0, 0, 0, 0, 0, 0, 0

	z80.tstates = 0

	z80.halted = false

	for i := 0; i < 0x10000; i++ {
		z80.memory.set(uint16(i), 0)
	}
}

// Initialize state from the snapshot defined by the specified filename.
// Returns nil on success.
func (z80 *Z80) LoadSna(filename string) os.Error {
	bytes, err := ioutil.ReadFile(filename)

	if err != nil {
		return err
	} else {
		if len(bytes) != 49179 {
			return os.NewError(fmt.Sprintf("snapshot \"%s\" has invalid size", filename))
		}

		z80.tstates = 0

		// Populate registers
		z80.i = bytes[0]
		z80.l_ = bytes[1]
		z80.h_ = bytes[2]
		z80.e_ = bytes[3]
		z80.d_ = bytes[4]
		z80.c_ = bytes[5]
		z80.b_ = bytes[6]
		z80.f_ = bytes[7]
		z80.a_ = bytes[8]
		z80.l = bytes[9]
		z80.h = bytes[10]
		z80.e = bytes[11]
		z80.d = bytes[12]
		z80.c = bytes[13]
		z80.b = bytes[14]
		z80.iyl = bytes[15]
		z80.iyh = bytes[16]
		z80.ixl = bytes[17]
		z80.ixh = bytes[18]

		z80.iff1 = uint16(ternOpB((bytes[19]&0x04) != 0, 1, 0))
		z80.iff2 = z80.iff1

		var r = uint16(bytes[20])

		z80.r = r & 0x7f
		z80.r7 = r & 0x80

		z80.f = bytes[21]
		z80.a = bytes[22]
		z80.sp = uint16(bytes[23]) | (uint16(bytes[24]) << 8)
		z80.im = uint16(bytes[25])

		// Border color
		z80.writePort(0xfe, bytes[26]&0x07)

		// Populate memory
		var i uint16
		for i = 0; i < 0xc000; i++ {
			z80.memory.set(uint16(i+0x4000), bytes[i+27])
		}

		// Set attribute bytes to force repaint of whole screen
		for i = 0x5800; i < 0x5b00; i++ {
			z80.memory.writeByte(i, z80.memory.At(uint(i)))
		}

		// Send a RETN
		z80.iff1 = z80.iff2
		z80.ret()
	}

	return nil
}

func splitWord(word uint16) (byte, byte) {
	return byte(word >> 8), byte(word & 0xff)
}
func joinBytes(h, l byte) uint16 {
	return uint16(l) | (uint16(h) << 8)
}

/* Process a z80 maskable interrupt */
func (z80 *Z80) interrupt() {
	if z80.iff1 != 0 {

		if z80.halted {
			z80.pc++
			z80.halted = false
		}

		z80.iff1, z80.iff2 = 0, 0

		pch, pcl := splitWord(z80.pc)

		z80.sp--

		z80.memory.writeByte(z80.sp, pch)

		z80.sp--

		z80.memory.writeByte(z80.sp, pcl)

		z80.r = (z80.r + 1) & 0x7f

		switch z80.im {
		case 0:
			z80.pc = 0x0038
			break
		case 1:
			z80.pc = 0x0038
			break
		case 2:
			var inttemp uint16 = (0x100 * uint16(z80.i)) + 0xff
			pcl := z80.memory.readByte(inttemp)
			inttemp++
			pch := z80.memory.readByte(inttemp)
			z80.pc = joinBytes(pch, pcl)
			break
		default:
			panic("Unknown interrupt mode")
		}
	}
}

func ternOpB(cond bool, ret1, ret2 byte) byte {
	if cond {
		return ret1
	}
	return ret2
}

func signExtend(v byte) int8 {
	if v < 128 {
		return int8(v)
	}
	return int8(int(v) - 256)
}

func (z80 *Z80) tapeSaveTrap() int {
	panic("tapeSaveTrap() should never be called")
}

func (z80 *Z80) tapeLoadTrap() int {
	/* Should never be called */
	panic("tapeLoadTrap() should never be called")
}

func (z80 *Z80) initTables() {

	var i int16
	var j, k byte
	var parity byte

	for i = 0; i < 0x100; i++ {
		z80.sz53Table[i] = byte(i) & (0x08 | 0x20 | 0x80)
		j = byte(i)
		parity = 0
		for k = 0; k < 8; k++ {
			parity ^= j & 1
			j >>= 1
		}
		z80.parityTable[i] = ternOpB(parity != 0, 0, 0x04)
		z80.sz53pTable[i] = z80.sz53Table[i] | z80.parityTable[i]
	}

	z80.sz53Table[0] |= 0x40
	z80.sz53pTable[0] |= 0x40

}

func (z80 *Z80) jp() {
	var jptemp uint16 = z80.pc
	pcl := z80.memory.readByte(jptemp)
	jptemp++
	pch := z80.memory.readByte(jptemp)
	z80.pc = joinBytes(pch, pcl)
}

func (z80 *Z80) dec(value *byte) {
	z80.f = (z80.f & FLAG_C) | ternOpB((*value&0x0f) != 0, 0, FLAG_H) | FLAG_N
	*value--
	z80.f |= ternOpB(*value == 0x7f, FLAG_V, 0) | z80.sz53Table[*value]
}

func (z80 *Z80) inc(value *byte) {
	*value++
	z80.f = (z80.f & FLAG_C) | ternOpB(*value == 0x80, FLAG_V, 0) | ternOpB((*value&0x0f) != 0, 0, FLAG_H) | z80.sz53Table[(*value)]
}

func (z80 *Z80) jr() {
	var jrtemp int8 = signExtend(z80.memory.readByte(z80.pc))
	z80.memory.contendReadNoMreq(z80.pc, 1)
	z80.memory.contendReadNoMreq(z80.pc, 1)
	z80.memory.contendReadNoMreq(z80.pc, 1)
	z80.memory.contendReadNoMreq(z80.pc, 1)
	z80.memory.contendReadNoMreq(z80.pc, 1)
	z80.pc += uint16(jrtemp)
}

func (z80 *Z80) ld16nnrr(regl, regh byte) {
	var ldtemp uint16

	ldtemp = uint16(z80.memory.readByte(z80.pc))
	z80.pc++
	ldtemp |= uint16(z80.memory.readByte(z80.pc)) << 8
	z80.pc++
	z80.memory.writeByte(ldtemp, regl)
	ldtemp++
	z80.memory.writeByte(ldtemp, regh)
}

func (z80 *Z80) ld16rrnn(regl, regh *byte) {
	var ldtemp uint16

	ldtemp = uint16(z80.memory.readByte(z80.pc))
	z80.pc++
	ldtemp |= uint16(z80.memory.readByte(z80.pc)) << 8
	z80.pc++
	*regl = z80.memory.readByte(ldtemp)
	ldtemp++
	*regh = z80.memory.readByte(ldtemp)
}

func (z80 *Z80) sub(value byte) {
	var subtemp uint16 = uint16(z80.a) - uint16(value)
	var lookup byte = ((z80.a & 0x88) >> 3) | ((value & 0x88) >> 2) | byte((subtemp&0x88)>>1)
	z80.a = byte(subtemp)
	z80.f = ternOpB(subtemp&0x100 != 0, FLAG_C, 0) | FLAG_N |
		halfcarrySubTable[lookup&0x07] | overflowSubTable[lookup>>4] |
		z80.sz53Table[z80.a]
}

func (z80 *Z80) and(value byte) {
	z80.a &= value
	z80.f = FLAG_H | z80.sz53pTable[z80.a]
}

func (z80 *Z80) adc(value byte) {
	var adctemp uint16 = uint16(z80.a) + uint16(value) + (uint16(z80.f) & FLAG_C)
	var lookup byte = byte(((uint16(z80.a) & 0x88) >> 3) | ((uint16(value) & 0x88) >> 2) | ((uint16(adctemp) & 0x88) >> 1))

	z80.a = byte(adctemp)

	z80.f = ternOpB((adctemp&0x100) != 0, FLAG_C, 0) | halfcarryAddTable[lookup&0x07] | overflowAddTable[lookup>>4] | z80.sz53Table[z80.a]
}

func (z80 *Z80) adc16(value uint16) {
	var add16temp uint = uint(z80.HL()) + uint(value) + (uint(z80.f) & FLAG_C)
	var lookup byte = byte(((uint(z80.HL()) & 0x8800) >> 11) | ((uint(value) & 0x8800) >> 10) | (add16temp&0x8800)>>9)

	z80.setHL(uint16(add16temp))

	z80.f = ternOpB((uint(add16temp)&0x10000) != 0, FLAG_C, 0) | overflowAddTable[lookup>>4] | (z80.h & (FLAG_3 | FLAG_5 | FLAG_S)) | halfcarryAddTable[lookup&0x07] | ternOpB(z80.HL() != 0, 0, FLAG_Z)
}

func (z80 *Z80) add16(value1 register16, value2 uint16) {
	var add16temp uint = uint(value1.get()) + uint(value2)
	var lookup byte = byte(((value1.get() & 0x0800) >> 11) | ((value2 & 0x0800) >> 10) | (uint16(add16temp)&0x0800)>>9)

	value1.set(uint16(add16temp))

	z80.f = (z80.f & (FLAG_V | FLAG_Z | FLAG_S)) | ternOpB((add16temp&0x10000) != 0, FLAG_C, 0) | (byte(add16temp>>8) & (FLAG_3 | FLAG_5)) | halfcarryAddTable[lookup]
}

func (z80 *Z80) add(value byte) {
	var addtemp uint = uint(z80.a) + uint(value)
	var lookup byte = ((z80.a & 0x88) >> 3) | ((value & 0x88) >> 2) | byte((addtemp&0x88)>>1)
	z80.a = byte(addtemp)
	z80.f = ternOpB(addtemp&0x100 != 0, FLAG_C, 0) | halfcarryAddTable[lookup&0x07] | overflowAddTable[lookup>>4] | z80.sz53Table[z80.a]
}

func (z80 *Z80) or(value byte) {
	z80.a |= value
	z80.f = z80.sz53pTable[z80.a]
}

func (z80 *Z80) pop16(regl, regh *byte) {
	*regl = z80.memory.readByte(z80.sp)
	z80.sp++
	*regh = z80.memory.readByte(z80.sp)
	z80.sp++
}

func (z80 *Z80) push16(regl, regh byte) {
	z80.sp--
	z80.memory.writeByte(z80.sp, regh)
	z80.sp--
	z80.memory.writeByte(z80.sp, regl)
}

func (z80 *Z80) ret() {
	pch, pcl := splitWord(z80.pc)
	z80.pop16(&pcl, &pch)
	z80.pc = joinBytes(pch, pcl)
}

func (z80 *Z80) rl(value *byte) {
	rltemp := *value
	*value = (*value << 1) | (z80.f & FLAG_C)
	z80.f = (rltemp >> 7) | z80.sz53pTable[*value]
}

func (z80 *Z80) rlc(value *byte) {
	*value = (*value << 1) | (*value >> 7)
	z80.f = (*value & FLAG_C) | z80.sz53pTable[*value]
}

func (z80 *Z80) rr(value *byte) {
	rrtemp := *value
	*value = (*value >> 1) | (z80.f << 7)
	z80.f = (rrtemp & FLAG_C) | z80.sz53pTable[*value]
}

func (z80 *Z80) rrc(value *byte) {
	z80.f = *value & FLAG_C
	*value = (*value >> 1) | (*value << 7)
	z80.f |= z80.sz53pTable[*value]
}

func (z80 *Z80) rst(value byte) {
	pch, pcl := splitWord(z80.pc)
	z80.push16(pcl, pch)
	z80.pc = uint16(value)
}

func (z80 *Z80) sbc(value byte) {
	var sbctemp uint16 = uint16(z80.a) - uint16(value) - (uint16(z80.f) & FLAG_C)
	var lookup byte = ((z80.a & 0x88) >> 3) | ((value & 0x88) >> 2) | byte((sbctemp&0x88)>>1)
	z80.a = byte(sbctemp)
	z80.f = ternOpB((sbctemp&0x100) != 0, FLAG_C, 0) | FLAG_N | halfcarrySubTable[lookup&0x07] | overflowSubTable[lookup>>4] | z80.sz53Table[z80.a]
}

func (z80 *Z80) sbc16(value uint16) {
	var sub16temp uint = uint(z80.HL()) - uint(value) - (uint(z80.f) & FLAG_C)
	var lookup byte = byte(((z80.HL() & 0x8800) >> 11) | ((uint16(value) & 0x8800) >> 10) | ((uint16(sub16temp) & 0x8800) >> 9))

	z80.setHL(uint16(sub16temp))

	z80.f = ternOpB((sub16temp&0x10000) != 0, FLAG_C, 0) | FLAG_N | overflowSubTable[lookup>>4] | (z80.h & (FLAG_3 | FLAG_5 | FLAG_S)) | halfcarrySubTable[lookup&0x07] | ternOpB(z80.HL() != 0, 0, FLAG_Z)
}

func (z80 *Z80) sla(value *byte) {
	z80.f = *value >> 7
	*value <<= 1
	z80.f |= z80.sz53pTable[*value]
}

func (z80 *Z80) sll(value *byte) {
	z80.f = *value >> 7
	*value = (*value << 1) | 0x01
	z80.f |= z80.sz53pTable[(*value)]
}

func (z80 *Z80) sra(value *byte) {
	z80.f = *value & FLAG_C
	*value = (*value & 0x80) | (*value >> 1)
	z80.f |= z80.sz53pTable[*value]
}

func (z80 *Z80) srl(value *byte) {
	z80.f = *value & FLAG_C
	*value >>= 1
	z80.f |= z80.sz53pTable[*value]
}

func (z80 *Z80) xor(value byte) {
	z80.a ^= value
	z80.f = z80.sz53pTable[z80.a]
}

func (z80 *Z80) bit(bit, value byte) {
	z80.f = (z80.f & FLAG_C) | FLAG_H | (value & (FLAG_3 | FLAG_5))
	if value&(0x01<<bit) == 0 {
		z80.f |= FLAG_P | FLAG_Z
	}
	if bit == 7 && (value&0x80) != 0 {
		z80.f |= FLAG_S
	}
}

func (z80 *Z80) biti(bit, value byte, address uint16) {
	z80.f = (z80.f & FLAG_C) | FLAG_H | (byte(address>>8) & (FLAG_3 | FLAG_5))
	if value&(0x01<<bit) == 0 {
		z80.f |= FLAG_P | FLAG_Z
	}
	if (bit == 7) && (value&0x80) != 0 {
		z80.f |= FLAG_S
	}
}

func (z80 *Z80) call() {
	var calltempl, calltemph byte
	calltempl = z80.memory.readByte(z80.pc)
	z80.pc++
	calltemph = z80.memory.readByte(z80.pc)
	z80.memory.contendReadNoMreq(z80.pc, 1)
	z80.pc++
	pch, pcl := splitWord(z80.pc)
	z80.push16(pcl, pch)
	z80.pc = joinBytes(calltemph, calltempl)
}

func (z80 *Z80) cp(value byte) {
	var cptemp uint16 = uint16(z80.a) - uint16(value)
	var lookup byte = ((z80.a & 0x88) >> 3) | ((value & 0x88) >> 2) | byte((cptemp&0x88)>>1)
	z80.f = ternOpB((cptemp&0x100) != 0, FLAG_C, ternOpB(cptemp != 0, 0, FLAG_Z)) | FLAG_N | halfcarrySubTable[lookup&0x07] | overflowSubTable[lookup>>4] | (value & (FLAG_3 | FLAG_5)) | byte(cptemp&FLAG_S)
}

func (z80 *Z80) in(reg *byte, port uint16) {
	*reg = z80.readPort(port)
	z80.f = (z80.f & FLAG_C) | z80.sz53pTable[*reg]
}

// Generated getters and INC/DEC functions for 8bit registers


func (z80 *Z80) A() byte {
	return z80.a
}

func (z80 *Z80) incA() {
	z80.a++
	z80.f = (z80.f & FLAG_C) | (ternOpB(z80.a == 0x80, FLAG_V, 0)) | (ternOpB((z80.a&0x0f) != 0, 0, FLAG_H)) | z80.sz53Table[z80.a]
}

func (z80 *Z80) decA() {
	z80.f = (z80.f & FLAG_C) | (ternOpB(z80.a&0x0f != 0, 0, FLAG_H)) | FLAG_N
	z80.a--
	z80.f |= (ternOpB(z80.a == 0x7f, FLAG_V, 0)) | z80.sz53Table[z80.a]

}


func (z80 *Z80) B() byte {
	return z80.b
}

func (z80 *Z80) incB() {
	z80.b++
	z80.f = (z80.f & FLAG_C) | (ternOpB(z80.b == 0x80, FLAG_V, 0)) | (ternOpB((z80.b&0x0f) != 0, 0, FLAG_H)) | z80.sz53Table[z80.b]
}

func (z80 *Z80) decB() {
	z80.f = (z80.f & FLAG_C) | (ternOpB(z80.b&0x0f != 0, 0, FLAG_H)) | FLAG_N
	z80.b--
	z80.f |= (ternOpB(z80.b == 0x7f, FLAG_V, 0)) | z80.sz53Table[z80.b]

}


func (z80 *Z80) C() byte {
	return z80.c
}

func (z80 *Z80) incC() {
	z80.c++
	z80.f = (z80.f & FLAG_C) | (ternOpB(z80.c == 0x80, FLAG_V, 0)) | (ternOpB((z80.c&0x0f) != 0, 0, FLAG_H)) | z80.sz53Table[z80.c]
}

func (z80 *Z80) decC() {
	z80.f = (z80.f & FLAG_C) | (ternOpB(z80.c&0x0f != 0, 0, FLAG_H)) | FLAG_N
	z80.c--
	z80.f |= (ternOpB(z80.c == 0x7f, FLAG_V, 0)) | z80.sz53Table[z80.c]

}


func (z80 *Z80) D() byte {
	return z80.d
}

func (z80 *Z80) incD() {
	z80.d++
	z80.f = (z80.f & FLAG_C) | (ternOpB(z80.d == 0x80, FLAG_V, 0)) | (ternOpB((z80.d&0x0f) != 0, 0, FLAG_H)) | z80.sz53Table[z80.d]
}

func (z80 *Z80) decD() {
	z80.f = (z80.f & FLAG_C) | (ternOpB(z80.d&0x0f != 0, 0, FLAG_H)) | FLAG_N
	z80.d--
	z80.f |= (ternOpB(z80.d == 0x7f, FLAG_V, 0)) | z80.sz53Table[z80.d]

}


func (z80 *Z80) E() byte {
	return z80.e
}

func (z80 *Z80) incE() {
	z80.e++
	z80.f = (z80.f & FLAG_C) | (ternOpB(z80.e == 0x80, FLAG_V, 0)) | (ternOpB((z80.e&0x0f) != 0, 0, FLAG_H)) | z80.sz53Table[z80.e]
}

func (z80 *Z80) decE() {
	z80.f = (z80.f & FLAG_C) | (ternOpB(z80.e&0x0f != 0, 0, FLAG_H)) | FLAG_N
	z80.e--
	z80.f |= (ternOpB(z80.e == 0x7f, FLAG_V, 0)) | z80.sz53Table[z80.e]

}


func (z80 *Z80) H() byte {
	return z80.h
}

func (z80 *Z80) incH() {
	z80.h++
	z80.f = (z80.f & FLAG_C) | (ternOpB(z80.h == 0x80, FLAG_V, 0)) | (ternOpB((z80.h&0x0f) != 0, 0, FLAG_H)) | z80.sz53Table[z80.h]
}

func (z80 *Z80) decH() {
	z80.f = (z80.f & FLAG_C) | (ternOpB(z80.h&0x0f != 0, 0, FLAG_H)) | FLAG_N
	z80.h--
	z80.f |= (ternOpB(z80.h == 0x7f, FLAG_V, 0)) | z80.sz53Table[z80.h]

}


func (z80 *Z80) L() byte {
	return z80.l
}

func (z80 *Z80) incL() {
	z80.l++
	z80.f = (z80.f & FLAG_C) | (ternOpB(z80.l == 0x80, FLAG_V, 0)) | (ternOpB((z80.l&0x0f) != 0, 0, FLAG_H)) | z80.sz53Table[z80.l]
}

func (z80 *Z80) decL() {
	z80.f = (z80.f & FLAG_C) | (ternOpB(z80.l&0x0f != 0, 0, FLAG_H)) | FLAG_N
	z80.l--
	z80.f |= (ternOpB(z80.l == 0x7f, FLAG_V, 0)) | z80.sz53Table[z80.l]

}


func (z80 *Z80) IXL() byte {
	return z80.ixl
}

func (z80 *Z80) incIXL() {
	z80.ixl++
	z80.f = (z80.f & FLAG_C) | (ternOpB(z80.ixl == 0x80, FLAG_V, 0)) | (ternOpB((z80.ixl&0x0f) != 0, 0, FLAG_H)) | z80.sz53Table[z80.ixl]
}

func (z80 *Z80) decIXL() {
	z80.f = (z80.f & FLAG_C) | (ternOpB(z80.ixl&0x0f != 0, 0, FLAG_H)) | FLAG_N
	z80.ixl--
	z80.f |= (ternOpB(z80.ixl == 0x7f, FLAG_V, 0)) | z80.sz53Table[z80.ixl]

}


func (z80 *Z80) IXH() byte {
	return z80.ixh
}

func (z80 *Z80) incIXH() {
	z80.ixh++
	z80.f = (z80.f & FLAG_C) | (ternOpB(z80.ixh == 0x80, FLAG_V, 0)) | (ternOpB((z80.ixh&0x0f) != 0, 0, FLAG_H)) | z80.sz53Table[z80.ixh]
}

func (z80 *Z80) decIXH() {
	z80.f = (z80.f & FLAG_C) | (ternOpB(z80.ixh&0x0f != 0, 0, FLAG_H)) | FLAG_N
	z80.ixh--
	z80.f |= (ternOpB(z80.ixh == 0x7f, FLAG_V, 0)) | z80.sz53Table[z80.ixh]

}


func (z80 *Z80) IYL() byte {
	return z80.iyl
}

func (z80 *Z80) incIYL() {
	z80.iyl++
	z80.f = (z80.f & FLAG_C) | (ternOpB(z80.iyl == 0x80, FLAG_V, 0)) | (ternOpB((z80.iyl&0x0f) != 0, 0, FLAG_H)) | z80.sz53Table[z80.iyl]
}

func (z80 *Z80) decIYL() {
	z80.f = (z80.f & FLAG_C) | (ternOpB(z80.iyl&0x0f != 0, 0, FLAG_H)) | FLAG_N
	z80.iyl--
	z80.f |= (ternOpB(z80.iyl == 0x7f, FLAG_V, 0)) | z80.sz53Table[z80.iyl]

}


func (z80 *Z80) IYH() byte {
	return z80.iyh
}

func (z80 *Z80) incIYH() {
	z80.iyh++
	z80.f = (z80.f & FLAG_C) | (ternOpB(z80.iyh == 0x80, FLAG_V, 0)) | (ternOpB((z80.iyh&0x0f) != 0, 0, FLAG_H)) | z80.sz53Table[z80.iyh]
}

func (z80 *Z80) decIYH() {
	z80.f = (z80.f & FLAG_C) | (ternOpB(z80.iyh&0x0f != 0, 0, FLAG_H)) | FLAG_N
	z80.iyh--
	z80.f |= (ternOpB(z80.iyh == 0x7f, FLAG_V, 0)) | z80.sz53Table[z80.iyh]

}


// Generated getters/setters and INC/DEC functions for 16bit registers


func (z80 *Z80) AF() uint16 {
	return z80.af.get()
}

func (z80 *Z80) setAF(value uint16) {
	z80.af.set(value)
}

func (z80 *Z80) decAF() {
	z80.af.dec()
}

func (z80 *Z80) incAF() {
	z80.af.inc()
}

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


// The following functions can not be generated as they need special treatments

func (z80 *Z80) PC() uint16 {
	return z80.pc
}

func (z80 *Z80) SP() uint16 {
	return z80.sp
}

func (z80 *Z80) setSP(value uint16) {
	z80.sp = value
}

func (z80 *Z80) incSP() {
	z80.sp++
}

func (z80 *Z80) decSP() {
	z80.sp--
}

func (z80 *Z80) IR() uint16 {
	return uint16(uint16(z80.i)<<8 | (z80.r7 & 0x80) | (z80.r & 0x7f))
}

func (z80 *Z80) sltTrap(address int16, level byte) int {
	return 0
}

func (z80 *Z80) doOpcodes() {
	for z80.tstates < eventNextEvent {

		z80.memory.contendRead(z80.pc, 4)

		opcode := z80.memory.readByteInternal(z80.pc)

	EndOpcode:
		z80.r = (z80.r + 1) & 0x7f
		z80.pc++

		switch opcode {

		/* opcodes_base.c: unshifted Z80 opcodes
		   Copyright (c) 1999-2003 Philip Kendall

		   This program is free software; you can redistribute it and/or modify
		   it under the terms of the GNU General Public License as published by
		   the Free Software Foundation; either version 2 of the License, or
		   (at your option) any later version.

		   This program is distributed in the hope that it will be useful,
		   but WITHOUT ANY WARRANTY; without even the implied warranty of
		   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
		   GNU General Public License for more details.

		   You should have received a copy of the GNU General Public License along
		   with this program; if not, write to the Free Software Foundation, Inc.,
		   51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.

		   Author contact information:

		   E-mail: philip-fuse@shadowmagic.org.uk

		*/

		/* NB: this file is autogenerated by './z80.pl' from 'opcodes_base.dat',
		   and included in 'z80_ops.c' */

		case 0x00: /* NOP */
			break
		case 0x01: /* LD BC,nnnn */
			b1 := z80.memory.readByte(z80.pc)
			z80.pc++
			b2 := z80.memory.readByte(z80.pc)
			z80.pc++
			z80.setBC(joinBytes(b2, b1))
			break
		case 0x02: /* LD (BC),A */
			z80.memory.writeByte(z80.BC(), z80.a)
			break
		case 0x03: /* INC BC */
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.incBC()
			break
		case 0x04: /* INC B */
			z80.incB()
			break
		case 0x05: /* DEC B */
			z80.decB()
			break
		case 0x06: /* LD B,nn */
			z80.b = z80.memory.readByte(z80.pc)
			z80.pc++
			break
		case 0x07: /* RLCA */
			z80.a = (z80.a << 1) | (z80.a >> 7)
			z80.f = (z80.f & (FLAG_P | FLAG_Z | FLAG_S)) |
				(z80.a & (FLAG_C | FLAG_3 | FLAG_5))
			break
		case 0x08: /* EX AF,AF' */
			/* Tape saving trap: note this traps the EX AF,AF' at #04d0, not
			#04d1 as PC has already been incremented */
			/* 0x76 - Timex 2068 save routine in EXROM */
			if z80.pc == 0x04d1 || z80.pc == 0x0077 {
				if z80.tapeSaveTrap() == 0 {
					break
				}
			}

			var olda, oldf = z80.a, z80.f
			z80.a = z80.a_
			z80.f = z80.f_
			z80.a_ = olda
			z80.f_ = oldf
			break
		case 0x09: /* ADD HL,BC */
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.add16(z80.hl, z80.BC())
			break
		case 0x0a: /* LD A,(BC) */
			z80.a = z80.memory.readByte(z80.BC())
			break
		case 0x0b: /* DEC BC */
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.decBC()
			break
		case 0x0c: /* INC C */
			z80.incC()
			break
		case 0x0d: /* DEC C */
			z80.decC()
			break
		case 0x0e: /* LD C,nn */
			z80.c = z80.memory.readByte(z80.pc)
			z80.pc++
			break
		case 0x0f: /* RRCA */
			z80.f = (z80.f & (FLAG_P | FLAG_Z | FLAG_S)) | (z80.a & FLAG_C)
			z80.a = (z80.a >> 1) | (z80.a << 7)
			z80.f |= (z80.a & (FLAG_3 | FLAG_5))
			break
		case 0x10: /* DJNZ offset */
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.b--
			if z80.b != 0 {
				z80.jr()
			} else {
				z80.memory.contendRead(z80.pc, 3)
			}
			z80.pc++
			break
		case 0x11: /* LD DE,nnnn */
			b1 := z80.memory.readByte(z80.pc)
			z80.pc++
			b2 := z80.memory.readByte(z80.pc)
			z80.pc++
			z80.setDE(joinBytes(b2, b1))
			break
		case 0x12: /* LD (DE),A */
			z80.memory.writeByte(z80.DE(), z80.a)
			break
		case 0x13: /* INC DE */
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.incDE()
			break
		case 0x14: /* INC D */
			z80.incD()
			break
		case 0x15: /* DEC D */
			z80.decD()
			break
		case 0x16: /* LD D,nn */
			z80.d = z80.memory.readByte(z80.pc)
			z80.pc++
			break
		case 0x17: /* RLA */
			var bytetemp byte = z80.a
			z80.a = (z80.a << 1) | (z80.f & FLAG_C)
			z80.f = (z80.f & (FLAG_P | FLAG_Z | FLAG_S)) | (z80.a & (FLAG_3 | FLAG_5)) | (bytetemp >> 7)
			break
		case 0x18: /* JR offset */
			z80.jr()
			z80.pc++
			break
		case 0x19: /* ADD HL,DE */
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.add16(z80.hl, z80.DE())
			break
		case 0x1a: /* LD A,(DE) */
			z80.a = z80.memory.readByte(z80.DE())
			break
		case 0x1b: /* DEC DE */
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.decDE()
			break
		case 0x1c: /* INC E */
			z80.incE()
			break
		case 0x1d: /* DEC E */
			z80.decE()
			break
		case 0x1e: /* LD E,nn */
			z80.e = z80.memory.readByte(z80.pc)
			z80.pc++
			break
		case 0x1f: /* RRA */
			var bytetemp byte = z80.a
			z80.a = (z80.a >> 1) | (z80.f << 7)
			z80.f = (z80.f & (FLAG_P | FLAG_Z | FLAG_S)) | (z80.a & (FLAG_3 | FLAG_5)) | (bytetemp & FLAG_C)
			break
		case 0x20: /* JR NZ,offset */
			if (z80.f & FLAG_Z) == 0 {
				z80.jr()
			} else {
				z80.memory.contendRead(z80.pc, 3)
			}
			z80.pc++
			break
		case 0x21: /* LD HL,nnnn */
			b1 := z80.memory.readByte(z80.pc)
			z80.pc++
			b2 := z80.memory.readByte(z80.pc)
			z80.pc++
			z80.setHL(joinBytes(b2, b1))
			break
		case 0x22: /* LD (nnnn),HL */
			z80.ld16nnrr(z80.l, z80.h)
			break
			break
		case 0x23: /* INC HL */
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.incHL()
			break
		case 0x24: /* INC H */
			z80.incH()
			break
		case 0x25: /* DEC H */
			z80.decH()
			break
		case 0x26: /* LD H,nn */
			z80.h = z80.memory.readByte(z80.pc)
			z80.pc++
			break
		case 0x27: /* DAA */
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
			var temp int = (int(z80.f) & ^(FLAG_C | FLAG_P)) | int(carry) | int(z80.parityTable[z80.a])
			z80.f = byte(temp)
			break
		case 0x28: /* JR Z,offset */
			if (z80.f & FLAG_Z) != 0 {
				z80.jr()
			} else {
				z80.memory.contendRead(z80.pc, 3)
			}
			z80.pc++
			break
		case 0x29: /* ADD HL,HL */
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.add16(z80.hl, z80.HL())
			break
		case 0x2a: /* LD HL,(nnnn) */
			z80.ld16rrnn(&z80.l, &z80.h)
			break
			break
		case 0x2b: /* DEC HL */
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.decHL()
			break
		case 0x2c: /* INC L */
			z80.incL()
			break
		case 0x2d: /* DEC L */
			z80.decL()
			break
		case 0x2e: /* LD L,nn */
			z80.l = z80.memory.readByte(z80.pc)
			z80.pc++
			break
		case 0x2f: /* CPL */
			z80.a ^= 0xff
			z80.f = (z80.f & (FLAG_C | FLAG_P | FLAG_Z | FLAG_S)) |
				(z80.a & (FLAG_3 | FLAG_5)) | (FLAG_N | FLAG_H)
			break
		case 0x30: /* JR NC,offset */
			if (z80.f & FLAG_C) == 0 {
				z80.jr()
			} else {
				z80.memory.contendRead(z80.pc, 3)
			}
			z80.pc++
			break
		case 0x31: /* LD SP,nnnn */
			b1 := z80.memory.readByte(z80.pc)
			z80.pc++
			b2 := z80.memory.readByte(z80.pc)
			z80.pc++
			z80.setSP(joinBytes(b2, b1))
			break
		case 0x32: /* LD (nnnn),A */
			var wordtemp uint16 = uint16(z80.memory.readByte(z80.pc))
			z80.pc++
			wordtemp |= uint16(z80.memory.readByte(z80.pc)) << 8
			z80.pc++
			z80.memory.writeByte(wordtemp, z80.a)
			break
		case 0x33: /* INC SP */
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.incSP()
			break
		case 0x34: /* INC (HL) */
			{
				var bytetemp byte = z80.memory.readByte(z80.HL())
				z80.memory.contendReadNoMreq(z80.HL(), 1)
				z80.inc(&bytetemp)
				z80.memory.writeByte(z80.HL(), bytetemp)
			}
			break
		case 0x35: /* DEC (HL) */
			{
				var bytetemp byte = z80.memory.readByte(z80.HL())
				z80.memory.contendReadNoMreq(z80.HL(), 1)
				z80.dec(&bytetemp)
				z80.memory.writeByte(z80.HL(), bytetemp)
			}
			break
		case 0x36: /* LD (HL),nn */
			z80.memory.writeByte(z80.HL(), z80.memory.readByte(z80.pc))
			z80.pc++
			break
		case 0x37: /* SCF */
			z80.f = (z80.f & (FLAG_P | FLAG_Z | FLAG_S)) |
				(z80.a & (FLAG_3 | FLAG_5)) |
				FLAG_C
			break
		case 0x38: /* JR C,offset */
			if (z80.f & FLAG_C) != 0 {
				z80.jr()
			} else {
				z80.memory.contendRead(z80.pc, 3)
			}
			z80.pc++
			break
		case 0x39: /* ADD HL,SP */
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.add16(z80.hl, z80.SP())
			break
		case 0x3a: /* LD A,(nnnn) */
			var wordtemp uint16
			wordtemp = uint16(z80.memory.readByte(z80.pc))
			z80.pc++
			wordtemp |= uint16(z80.memory.readByte(z80.pc)) << 8
			z80.pc++
			z80.a = z80.memory.readByte(wordtemp)
			break
		case 0x3b: /* DEC SP */
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.decSP()
			break
		case 0x3c: /* INC A */
			z80.incA()
			break
		case 0x3d: /* DEC A */
			z80.decA()
			break
		case 0x3e: /* LD A,nn */
			z80.a = z80.memory.readByte(z80.pc)
			z80.pc++
			break
		case 0x3f: /* CCF */
			z80.f = (z80.f & (FLAG_P | FLAG_Z | FLAG_S)) |
				ternOpB((z80.f&FLAG_C) != 0, FLAG_H, FLAG_C) | (z80.a & (FLAG_3 | FLAG_5))
			break
		case 0x40: /* LD B,B */
			break
		case 0x41: /* LD B,C */
			z80.b = z80.c
			break
		case 0x42: /* LD B,D */
			z80.b = z80.d
			break
		case 0x43: /* LD B,E */
			z80.b = z80.e
			break
		case 0x44: /* LD B,H */
			z80.b = z80.h
			break
		case 0x45: /* LD B,L */
			z80.b = z80.l
			break
		case 0x46: /* LD B,(HL) */
			z80.b = z80.memory.readByte(z80.HL())
			break
		case 0x47: /* LD B,A */
			z80.b = z80.a
			break
		case 0x48: /* LD C,B */
			z80.c = z80.b
			break
		case 0x49: /* LD C,C */
			break
		case 0x4a: /* LD C,D */
			z80.c = z80.d
			break
		case 0x4b: /* LD C,E */
			z80.c = z80.e
			break
		case 0x4c: /* LD C,H */
			z80.c = z80.h
			break
		case 0x4d: /* LD C,L */
			z80.c = z80.l
			break
		case 0x4e: /* LD C,(HL) */
			z80.c = z80.memory.readByte(z80.HL())
			break
		case 0x4f: /* LD C,A */
			z80.c = z80.a
			break
		case 0x50: /* LD D,B */
			z80.d = z80.b
			break
		case 0x51: /* LD D,C */
			z80.d = z80.c
			break
		case 0x52: /* LD D,D */
			break
		case 0x53: /* LD D,E */
			z80.d = z80.e
			break
		case 0x54: /* LD D,H */
			z80.d = z80.h
			break
		case 0x55: /* LD D,L */
			z80.d = z80.l
			break
		case 0x56: /* LD D,(HL) */
			z80.d = z80.memory.readByte(z80.HL())
			break
		case 0x57: /* LD D,A */
			z80.d = z80.a
			break
		case 0x58: /* LD E,B */
			z80.e = z80.b
			break
		case 0x59: /* LD E,C */
			z80.e = z80.c
			break
		case 0x5a: /* LD E,D */
			z80.e = z80.d
			break
		case 0x5b: /* LD E,E */
			break
		case 0x5c: /* LD E,H */
			z80.e = z80.h
			break
		case 0x5d: /* LD E,L */
			z80.e = z80.l
			break
		case 0x5e: /* LD E,(HL) */
			z80.e = z80.memory.readByte(z80.HL())
			break
		case 0x5f: /* LD E,A */
			z80.e = z80.a
			break
		case 0x60: /* LD H,B */
			z80.h = z80.b
			break
		case 0x61: /* LD H,C */
			z80.h = z80.c
			break
		case 0x62: /* LD H,D */
			z80.h = z80.d
			break
		case 0x63: /* LD H,E */
			z80.h = z80.e
			break
		case 0x64: /* LD H,H */
			break
		case 0x65: /* LD H,L */
			z80.h = z80.l
			break
		case 0x66: /* LD H,(HL) */
			z80.h = z80.memory.readByte(z80.HL())
			break
		case 0x67: /* LD H,A */
			z80.h = z80.a
			break
		case 0x68: /* LD L,B */
			z80.l = z80.b
			break
		case 0x69: /* LD L,C */
			z80.l = z80.c
			break
		case 0x6a: /* LD L,D */
			z80.l = z80.d
			break
		case 0x6b: /* LD L,E */
			z80.l = z80.e
			break
		case 0x6c: /* LD L,H */
			z80.l = z80.h
			break
		case 0x6d: /* LD L,L */
			break
		case 0x6e: /* LD L,(HL) */
			z80.l = z80.memory.readByte(z80.HL())
			break
		case 0x6f: /* LD L,A */
			z80.l = z80.a
			break
		case 0x70: /* LD (HL),B */
			z80.memory.writeByte(z80.HL(), z80.b)
			break
		case 0x71: /* LD (HL),C */
			z80.memory.writeByte(z80.HL(), z80.c)
			break
		case 0x72: /* LD (HL),D */
			z80.memory.writeByte(z80.HL(), z80.d)
			break
		case 0x73: /* LD (HL),E */
			z80.memory.writeByte(z80.HL(), z80.e)
			break
		case 0x74: /* LD (HL),H */
			z80.memory.writeByte(z80.HL(), z80.h)
			break
		case 0x75: /* LD (HL),L */
			z80.memory.writeByte(z80.HL(), z80.l)
			break
		case 0x76: /* HALT */
			z80.halted = true
			z80.pc--
			break
		case 0x77: /* LD (HL),A */
			z80.memory.writeByte(z80.HL(), z80.a)
			break
		case 0x78: /* LD A,B */
			z80.a = z80.b
			break
		case 0x79: /* LD A,C */
			z80.a = z80.c
			break
		case 0x7a: /* LD A,D */
			z80.a = z80.d
			break
		case 0x7b: /* LD A,E */
			z80.a = z80.e
			break
		case 0x7c: /* LD A,H */
			z80.a = z80.h
			break
		case 0x7d: /* LD A,L */
			z80.a = z80.l
			break
		case 0x7e: /* LD A,(HL) */
			z80.a = z80.memory.readByte(z80.HL())
			break
		case 0x7f: /* LD A,A */
			break
		case 0x80: /* ADD A,B */
			z80.add(z80.b)
			break
		case 0x81: /* ADD A,C */
			z80.add(z80.c)
			break
		case 0x82: /* ADD A,D */
			z80.add(z80.d)
			break
		case 0x83: /* ADD A,E */
			z80.add(z80.e)
			break
		case 0x84: /* ADD A,H */
			z80.add(z80.h)
			break
		case 0x85: /* ADD A,L */
			z80.add(z80.l)
			break
		case 0x86: /* ADD A,(HL) */
			{
				var bytetemp byte = z80.memory.readByte(z80.HL())

				z80.add(bytetemp)
			}
			break
		case 0x87: /* ADD A,A */
			z80.add(z80.a)
			break
		case 0x88: /* ADC A,B */
			z80.adc(z80.b)
			break
		case 0x89: /* ADC A,C */
			z80.adc(z80.c)
			break
		case 0x8a: /* ADC A,D */
			z80.adc(z80.d)
			break
		case 0x8b: /* ADC A,E */
			z80.adc(z80.e)
			break
		case 0x8c: /* ADC A,H */
			z80.adc(z80.h)
			break
		case 0x8d: /* ADC A,L */
			z80.adc(z80.l)
			break
		case 0x8e: /* ADC A,(HL) */
			{
				var bytetemp byte = z80.memory.readByte(z80.HL())

				z80.adc(bytetemp)
			}
			break
		case 0x8f: /* ADC A,A */
			z80.adc(z80.a)
			break
		case 0x90: /* SUB A,B */
			z80.sub(z80.b)
			break
		case 0x91: /* SUB A,C */
			z80.sub(z80.c)
			break
		case 0x92: /* SUB A,D */
			z80.sub(z80.d)
			break
		case 0x93: /* SUB A,E */
			z80.sub(z80.e)
			break
		case 0x94: /* SUB A,H */
			z80.sub(z80.h)
			break
		case 0x95: /* SUB A,L */
			z80.sub(z80.l)
			break
		case 0x96: /* SUB A,(HL) */
			{
				var bytetemp byte = z80.memory.readByte(z80.HL())

				z80.sub(bytetemp)
			}
			break
		case 0x97: /* SUB A,A */
			z80.sub(z80.a)
			break
		case 0x98: /* SBC A,B */
			z80.sbc(z80.b)
			break
		case 0x99: /* SBC A,C */
			z80.sbc(z80.c)
			break
		case 0x9a: /* SBC A,D */
			z80.sbc(z80.d)
			break
		case 0x9b: /* SBC A,E */
			z80.sbc(z80.e)
			break
		case 0x9c: /* SBC A,H */
			z80.sbc(z80.h)
			break
		case 0x9d: /* SBC A,L */
			z80.sbc(z80.l)
			break
		case 0x9e: /* SBC A,(HL) */
			{
				var bytetemp byte = z80.memory.readByte(z80.HL())

				z80.sbc(bytetemp)
			}
			break
		case 0x9f: /* SBC A,A */
			z80.sbc(z80.a)
			break
		case 0xa0: /* AND A,B */
			z80.and(z80.b)
			break
		case 0xa1: /* AND A,C */
			z80.and(z80.c)
			break
		case 0xa2: /* AND A,D */
			z80.and(z80.d)
			break
		case 0xa3: /* AND A,E */
			z80.and(z80.e)
			break
		case 0xa4: /* AND A,H */
			z80.and(z80.h)
			break
		case 0xa5: /* AND A,L */
			z80.and(z80.l)
			break
		case 0xa6: /* AND A,(HL) */
			{
				var bytetemp byte = z80.memory.readByte(z80.HL())

				z80.and(bytetemp)
			}
			break
		case 0xa7: /* AND A,A */
			z80.and(z80.a)
			break
		case 0xa8: /* XOR A,B */
			z80.xor(z80.b)
			break
		case 0xa9: /* XOR A,C */
			z80.xor(z80.c)
			break
		case 0xaa: /* XOR A,D */
			z80.xor(z80.d)
			break
		case 0xab: /* XOR A,E */
			z80.xor(z80.e)
			break
		case 0xac: /* XOR A,H */
			z80.xor(z80.h)
			break
		case 0xad: /* XOR A,L */
			z80.xor(z80.l)
			break
		case 0xae: /* XOR A,(HL) */
			{
				var bytetemp byte = z80.memory.readByte(z80.HL())

				z80.xor(bytetemp)
			}
			break
		case 0xaf: /* XOR A,A */
			z80.xor(z80.a)
			break
		case 0xb0: /* OR A,B */
			z80.or(z80.b)
			break
		case 0xb1: /* OR A,C */
			z80.or(z80.c)
			break
		case 0xb2: /* OR A,D */
			z80.or(z80.d)
			break
		case 0xb3: /* OR A,E */
			z80.or(z80.e)
			break
		case 0xb4: /* OR A,H */
			z80.or(z80.h)
			break
		case 0xb5: /* OR A,L */
			z80.or(z80.l)
			break
		case 0xb6: /* OR A,(HL) */
			{
				var bytetemp byte = z80.memory.readByte(z80.HL())

				z80.or(bytetemp)
			}
			break
		case 0xb7: /* OR A,A */
			z80.or(z80.a)
			break
		case 0xb8: /* CP B */
			z80.cp(z80.b)
			break
		case 0xb9: /* CP C */
			z80.cp(z80.c)
			break
		case 0xba: /* CP D */
			z80.cp(z80.d)
			break
		case 0xbb: /* CP E */
			z80.cp(z80.e)
			break
		case 0xbc: /* CP H */
			z80.cp(z80.h)
			break
		case 0xbd: /* CP L */
			z80.cp(z80.l)
			break
		case 0xbe: /* CP (HL) */
			{
				var bytetemp byte = z80.memory.readByte(z80.HL())

				z80.cp(bytetemp)
			}
			break
		case 0xbf: /* CP A */
			z80.cp(z80.a)
			break
		case 0xc0: /* RET NZ */
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			if z80.pc == 0x056c || z80.pc == 0x0112 {
				if z80.tapeLoadTrap() == 0 {
					break
				}
			}
			if !((z80.f & FLAG_Z) != 0) {
				z80.ret()
			}
			break
		case 0xc1: /* POP BC */
			z80.pop16(&z80.c, &z80.b)
			break
		case 0xc2: /* JP NZ,nnnn */
			if (z80.f & FLAG_Z) == 0 {
				z80.jp()
			} else {
				z80.memory.contendRead(z80.pc, 3)
				z80.memory.contendRead(z80.pc+1, 3)
				z80.pc += 2
			}
			break
		case 0xc3: /* JP nnnn */
			z80.jp()
			break
		case 0xc4: /* CALL NZ,nnnn */
			if (z80.f & FLAG_Z) == 0 {
				z80.call()
			} else {
				z80.memory.contendRead(z80.pc, 3)
				z80.memory.contendRead(z80.pc+1, 3)
				z80.pc += 2
			}
			break
		case 0xc5: /* PUSH BC */
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.push16(z80.c, z80.b)
			break
		case 0xc6: /* ADD A,nn */
			{
				var bytetemp byte = z80.memory.readByte(z80.PC())
				z80.pc++
				z80.add(bytetemp)
			}
			break
		case 0xc7: /* RST 00 */
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.rst(0x00)
			break
		case 0xc8: /* RET Z */
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			if (z80.f & FLAG_Z) != 0 {
				z80.ret()
			}
			break
		case 0xc9: /* RET */
			z80.ret()
			break
		case 0xca: /* JP Z,nnnn */
			if (z80.f & FLAG_Z) != 0 {
				z80.jp()
			} else {
				z80.memory.contendRead(z80.pc, 3)
				z80.memory.contendRead(z80.pc+1, 3)
				z80.pc += 2
			}
			break
		case 0xcb: /* shift CB */
			{
				var opcode2 byte
				z80.memory.contendRead(z80.pc, 4)
				opcode2 = z80.memory.readByteInternal(z80.pc)
				z80.pc++
				z80.r++

				switch opcode2 {
				/* z80_cb.c: Z80 CBxx opcodes
				   Copyright (c) 1999-2003 Philip Kendall

				   This program is free software; you can redistribute it and/or modify
				   it under the terms of the GNU General Public License as published by
				   the Free Software Foundation; either version 2 of the License, or
				   (at your option) any later version.

				   This program is distributed in the hope that it will be useful,
				   but WITHOUT ANY WARRANTY; without even the implied warranty of
				   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
				   GNU General Public License for more details.

				   You should have received a copy of the GNU General Public License along
				   with this program; if not, write to the Free Software Foundation, Inc.,
				   51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.

				   Author contact information:

				   E-mail: philip-fuse@shadowmagic.org.uk

				*/

				/* NB: this file is autogenerated by './z80.pl' from 'opcodes_cb.dat',
				   and included in 'z80_ops.c' */

				case 0x00: /* RLC B */
					z80.rlc(&z80.b)
					break
				case 0x01: /* RLC C */
					z80.rlc(&z80.c)
					break
				case 0x02: /* RLC D */
					z80.rlc(&z80.d)
					break
				case 0x03: /* RLC E */
					z80.rlc(&z80.e)
					break
				case 0x04: /* RLC H */
					z80.rlc(&z80.h)
					break
				case 0x05: /* RLC L */
					z80.rlc(&z80.l)
					break
				case 0x06: /* RLC (HL) */
					var bytetemp byte = z80.memory.readByte(z80.HL())
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.rlc(&bytetemp)
					z80.memory.writeByte(z80.HL(), bytetemp)
					break
				case 0x07: /* RLC A */
					z80.rlc(&z80.a)
					break
				case 0x08: /* RRC B */
					z80.rrc(&z80.b)
					break
				case 0x09: /* RRC C */
					z80.rrc(&z80.c)
					break
				case 0x0a: /* RRC D */
					z80.rrc(&z80.d)
					break
				case 0x0b: /* RRC E */
					z80.rrc(&z80.e)
					break
				case 0x0c: /* RRC H */
					z80.rrc(&z80.h)
					break
				case 0x0d: /* RRC L */
					z80.rrc(&z80.l)
					break
				case 0x0e: /* RRC (HL) */
					var bytetemp byte = z80.memory.readByte(z80.HL())
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.rrc(&bytetemp)
					z80.memory.writeByte(z80.HL(), bytetemp)
					break
				case 0x0f: /* RRC A */
					z80.rrc(&z80.a)
					break
				case 0x10: /* RL B */
					z80.rl(&z80.b)
					break
				case 0x11: /* RL C */
					z80.rl(&z80.c)
					break
				case 0x12: /* RL D */
					z80.rl(&z80.d)
					break
				case 0x13: /* RL E */
					z80.rl(&z80.e)
					break
				case 0x14: /* RL H */
					z80.rl(&z80.h)
					break
				case 0x15: /* RL L */
					z80.rl(&z80.l)
					break
				case 0x16: /* RL (HL) */
					var bytetemp byte = z80.memory.readByte(z80.HL())
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.rl(&bytetemp)
					z80.memory.writeByte(z80.HL(), bytetemp)
					break
				case 0x17: /* RL A */
					z80.rl(&z80.a)
					break
				case 0x18: /* RR B */
					z80.rr(&z80.b)
					break
				case 0x19: /* RR C */
					z80.rr(&z80.c)
					break
				case 0x1a: /* RR D */
					z80.rr(&z80.d)
					break
				case 0x1b: /* RR E */
					z80.rr(&z80.e)
					break
				case 0x1c: /* RR H */
					z80.rr(&z80.h)
					break
				case 0x1d: /* RR L */
					z80.rr(&z80.l)
					break
				case 0x1e: /* RR (HL) */
					var bytetemp byte = z80.memory.readByte(z80.HL())
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.rr(&bytetemp)
					z80.memory.writeByte(z80.HL(), bytetemp)
					break
				case 0x1f: /* RR A */
					z80.rr(&z80.a)
					break
				case 0x20: /* SLA B */
					z80.sla(&z80.b)
					break
				case 0x21: /* SLA C */
					z80.sla(&z80.c)
					break
				case 0x22: /* SLA D */
					z80.sla(&z80.d)
					break
				case 0x23: /* SLA E */
					z80.sla(&z80.e)
					break
				case 0x24: /* SLA H */
					z80.sla(&z80.h)
					break
				case 0x25: /* SLA L */
					z80.sla(&z80.l)
					break
				case 0x26: /* SLA (HL) */
					var bytetemp byte = z80.memory.readByte(z80.HL())
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.sla(&bytetemp)
					z80.memory.writeByte(z80.HL(), bytetemp)
					break
				case 0x27: /* SLA A */
					z80.sla(&z80.a)
					break
				case 0x28: /* SRA B */
					z80.sra(&z80.b)
					break
				case 0x29: /* SRA C */
					z80.sra(&z80.c)
					break
				case 0x2a: /* SRA D */
					z80.sra(&z80.d)
					break
				case 0x2b: /* SRA E */
					z80.sra(&z80.e)
					break
				case 0x2c: /* SRA H */
					z80.sra(&z80.h)
					break
				case 0x2d: /* SRA L */
					z80.sra(&z80.l)
					break
				case 0x2e: /* SRA (HL) */
					var bytetemp byte = z80.memory.readByte(z80.HL())
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.sra(&bytetemp)
					z80.memory.writeByte(z80.HL(), bytetemp)
					break
				case 0x2f: /* SRA A */
					z80.sra(&z80.a)
					break
				case 0x30: /* SLL B */
					z80.sll(&z80.b)
					break
				case 0x31: /* SLL C */
					z80.sll(&z80.c)
					break
				case 0x32: /* SLL D */
					z80.sll(&z80.d)
					break
				case 0x33: /* SLL E */
					z80.sll(&z80.e)
					break
				case 0x34: /* SLL H */
					z80.sll(&z80.h)
					break
				case 0x35: /* SLL L */
					z80.sll(&z80.l)
					break
				case 0x36: /* SLL (HL) */
					var bytetemp byte = z80.memory.readByte(z80.HL())
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.sll(&bytetemp)
					z80.memory.writeByte(z80.HL(), bytetemp)
					break
				case 0x37: /* SLL A */
					z80.sll(&z80.a)
					break
				case 0x38: /* SRL B */
					z80.srl(&z80.b)
					break
				case 0x39: /* SRL C */
					z80.srl(&z80.c)
					break
				case 0x3a: /* SRL D */
					z80.srl(&z80.d)
					break
				case 0x3b: /* SRL E */
					z80.srl(&z80.e)
					break
				case 0x3c: /* SRL H */
					z80.srl(&z80.h)
					break
				case 0x3d: /* SRL L */
					z80.srl(&z80.l)
					break
				case 0x3e: /* SRL (HL) */
					var bytetemp byte = z80.memory.readByte(z80.HL())
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.srl(&bytetemp)
					z80.memory.writeByte(z80.HL(), bytetemp)
					break
				case 0x3f: /* SRL A */
					z80.srl(&z80.a)
					break
				case 0x40: /* BIT 0,B */
					z80.bit(0, z80.b)
					break
				case 0x41: /* BIT 0,C */
					z80.bit(0, z80.c)
					break
				case 0x42: /* BIT 0,D */
					z80.bit(0, z80.d)
					break
				case 0x43: /* BIT 0,E */
					z80.bit(0, z80.e)
					break
				case 0x44: /* BIT 0,H */
					z80.bit(0, z80.h)
					break
				case 0x45: /* BIT 0,L */
					z80.bit(0, z80.l)
					break
				case 0x46: /* BIT 0,(HL) */
					bytetemp := z80.memory.readByte(z80.HL())
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.bit(0, bytetemp)
					break
				case 0x47: /* BIT 0,A */
					z80.bit(0, z80.a)
					break
				case 0x48: /* BIT 1,B */
					z80.bit(1, z80.b)
					break
				case 0x49: /* BIT 1,C */
					z80.bit(1, z80.c)
					break
				case 0x4a: /* BIT 1,D */
					z80.bit(1, z80.d)
					break
				case 0x4b: /* BIT 1,E */
					z80.bit(1, z80.e)
					break
				case 0x4c: /* BIT 1,H */
					z80.bit(1, z80.h)
					break
				case 0x4d: /* BIT 1,L */
					z80.bit(1, z80.l)
					break
				case 0x4e: /* BIT 1,(HL) */
					bytetemp := z80.memory.readByte(z80.HL())
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.bit(1, bytetemp)
					break
				case 0x4f: /* BIT 1,A */
					z80.bit(1, z80.a)
					break
				case 0x50: /* BIT 2,B */
					z80.bit(2, z80.b)
					break
				case 0x51: /* BIT 2,C */
					z80.bit(2, z80.c)
					break
				case 0x52: /* BIT 2,D */
					z80.bit(2, z80.d)
					break
				case 0x53: /* BIT 2,E */
					z80.bit(2, z80.e)
					break
				case 0x54: /* BIT 2,H */
					z80.bit(2, z80.h)
					break
				case 0x55: /* BIT 2,L */
					z80.bit(2, z80.l)
					break
				case 0x56: /* BIT 2,(HL) */
					bytetemp := z80.memory.readByte(z80.HL())
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.bit(2, bytetemp)
					break
				case 0x57: /* BIT 2,A */
					z80.bit(2, z80.a)
					break
				case 0x58: /* BIT 3,B */
					z80.bit(3, z80.b)
					break
				case 0x59: /* BIT 3,C */
					z80.bit(3, z80.c)
					break
				case 0x5a: /* BIT 3,D */
					z80.bit(3, z80.d)
					break
				case 0x5b: /* BIT 3,E */
					z80.bit(3, z80.e)
					break
				case 0x5c: /* BIT 3,H */
					z80.bit(3, z80.h)
					break
				case 0x5d: /* BIT 3,L */
					z80.bit(3, z80.l)
					break
				case 0x5e: /* BIT 3,(HL) */
					bytetemp := z80.memory.readByte(z80.HL())
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.bit(3, bytetemp)
					break
				case 0x5f: /* BIT 3,A */
					z80.bit(3, z80.a)
					break
				case 0x60: /* BIT 4,B */
					z80.bit(4, z80.b)
					break
				case 0x61: /* BIT 4,C */
					z80.bit(4, z80.c)
					break
				case 0x62: /* BIT 4,D */
					z80.bit(4, z80.d)
					break
				case 0x63: /* BIT 4,E */
					z80.bit(4, z80.e)
					break
				case 0x64: /* BIT 4,H */
					z80.bit(4, z80.h)
					break
				case 0x65: /* BIT 4,L */
					z80.bit(4, z80.l)
					break
				case 0x66: /* BIT 4,(HL) */
					bytetemp := z80.memory.readByte(z80.HL())
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.bit(4, bytetemp)
					break
				case 0x67: /* BIT 4,A */
					z80.bit(4, z80.a)
					break
				case 0x68: /* BIT 5,B */
					z80.bit(5, z80.b)
					break
				case 0x69: /* BIT 5,C */
					z80.bit(5, z80.c)
					break
				case 0x6a: /* BIT 5,D */
					z80.bit(5, z80.d)
					break
				case 0x6b: /* BIT 5,E */
					z80.bit(5, z80.e)
					break
				case 0x6c: /* BIT 5,H */
					z80.bit(5, z80.h)
					break
				case 0x6d: /* BIT 5,L */
					z80.bit(5, z80.l)
					break
				case 0x6e: /* BIT 5,(HL) */
					bytetemp := z80.memory.readByte(z80.HL())
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.bit(5, bytetemp)
					break
				case 0x6f: /* BIT 5,A */
					z80.bit(5, z80.a)
					break
				case 0x70: /* BIT 6,B */
					z80.bit(6, z80.b)
					break
				case 0x71: /* BIT 6,C */
					z80.bit(6, z80.c)
					break
				case 0x72: /* BIT 6,D */
					z80.bit(6, z80.d)
					break
				case 0x73: /* BIT 6,E */
					z80.bit(6, z80.e)
					break
				case 0x74: /* BIT 6,H */
					z80.bit(6, z80.h)
					break
				case 0x75: /* BIT 6,L */
					z80.bit(6, z80.l)
					break
				case 0x76: /* BIT 6,(HL) */
					bytetemp := z80.memory.readByte(z80.HL())
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.bit(6, bytetemp)
					break
				case 0x77: /* BIT 6,A */
					z80.bit(6, z80.a)
					break
				case 0x78: /* BIT 7,B */
					z80.bit(7, z80.b)
					break
				case 0x79: /* BIT 7,C */
					z80.bit(7, z80.c)
					break
				case 0x7a: /* BIT 7,D */
					z80.bit(7, z80.d)
					break
				case 0x7b: /* BIT 7,E */
					z80.bit(7, z80.e)
					break
				case 0x7c: /* BIT 7,H */
					z80.bit(7, z80.h)
					break
				case 0x7d: /* BIT 7,L */
					z80.bit(7, z80.l)
					break
				case 0x7e: /* BIT 7,(HL) */
					bytetemp := z80.memory.readByte(z80.HL())
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.bit(7, bytetemp)
					break
				case 0x7f: /* BIT 7,A */
					z80.bit(7, z80.a)
					break
				case 0x80: /* RES 0,B */
					z80.b &= 0xfe
					break
				case 0x81: /* RES 0,C */
					z80.c &= 0xfe
					break
				case 0x82: /* RES 0,D */
					z80.d &= 0xfe
					break
				case 0x83: /* RES 0,E */
					z80.e &= 0xfe
					break
				case 0x84: /* RES 0,H */
					z80.h &= 0xfe
					break
				case 0x85: /* RES 0,L */
					z80.l &= 0xfe
					break
				case 0x86: /* RES 0,(HL) */
					var bytetemp byte = z80.memory.readByte(z80.HL())
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.memory.writeByte(z80.HL(), bytetemp&0xfe)
					break
				case 0x87: /* RES 0,A */
					z80.a &= 0xfe
					break
				case 0x88: /* RES 1,B */
					z80.b &= 0xfd
					break
				case 0x89: /* RES 1,C */
					z80.c &= 0xfd
					break
				case 0x8a: /* RES 1,D */
					z80.d &= 0xfd
					break
				case 0x8b: /* RES 1,E */
					z80.e &= 0xfd
					break
				case 0x8c: /* RES 1,H */
					z80.h &= 0xfd
					break
				case 0x8d: /* RES 1,L */
					z80.l &= 0xfd
					break
				case 0x8e: /* RES 1,(HL) */
					var bytetemp byte = z80.memory.readByte(z80.HL())
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.memory.writeByte(z80.HL(), bytetemp&0xfd)
					break
				case 0x8f: /* RES 1,A */
					z80.a &= 0xfd
					break
				case 0x90: /* RES 2,B */
					z80.b &= 0xfb
					break
				case 0x91: /* RES 2,C */
					z80.c &= 0xfb
					break
				case 0x92: /* RES 2,D */
					z80.d &= 0xfb
					break
				case 0x93: /* RES 2,E */
					z80.e &= 0xfb
					break
				case 0x94: /* RES 2,H */
					z80.h &= 0xfb
					break
				case 0x95: /* RES 2,L */
					z80.l &= 0xfb
					break
				case 0x96: /* RES 2,(HL) */
					var bytetemp byte = z80.memory.readByte(z80.HL())
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.memory.writeByte(z80.HL(), bytetemp&0xfb)
					break
				case 0x97: /* RES 2,A */
					z80.a &= 0xfb
					break
				case 0x98: /* RES 3,B */
					z80.b &= 0xf7
					break
				case 0x99: /* RES 3,C */
					z80.c &= 0xf7
					break
				case 0x9a: /* RES 3,D */
					z80.d &= 0xf7
					break
				case 0x9b: /* RES 3,E */
					z80.e &= 0xf7
					break
				case 0x9c: /* RES 3,H */
					z80.h &= 0xf7
					break
				case 0x9d: /* RES 3,L */
					z80.l &= 0xf7
					break
				case 0x9e: /* RES 3,(HL) */
					var bytetemp byte = z80.memory.readByte(z80.HL())
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.memory.writeByte(z80.HL(), bytetemp&0xf7)
					break
				case 0x9f: /* RES 3,A */
					z80.a &= 0xf7
					break
				case 0xa0: /* RES 4,B */
					z80.b &= 0xef
					break
				case 0xa1: /* RES 4,C */
					z80.c &= 0xef
					break
				case 0xa2: /* RES 4,D */
					z80.d &= 0xef
					break
				case 0xa3: /* RES 4,E */
					z80.e &= 0xef
					break
				case 0xa4: /* RES 4,H */
					z80.h &= 0xef
					break
				case 0xa5: /* RES 4,L */
					z80.l &= 0xef
					break
				case 0xa6: /* RES 4,(HL) */
					var bytetemp byte = z80.memory.readByte(z80.HL())
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.memory.writeByte(z80.HL(), bytetemp&0xef)
					break
				case 0xa7: /* RES 4,A */
					z80.a &= 0xef
					break
				case 0xa8: /* RES 5,B */
					z80.b &= 0xdf
					break
				case 0xa9: /* RES 5,C */
					z80.c &= 0xdf
					break
				case 0xaa: /* RES 5,D */
					z80.d &= 0xdf
					break
				case 0xab: /* RES 5,E */
					z80.e &= 0xdf
					break
				case 0xac: /* RES 5,H */
					z80.h &= 0xdf
					break
				case 0xad: /* RES 5,L */
					z80.l &= 0xdf
					break
				case 0xae: /* RES 5,(HL) */
					var bytetemp byte = z80.memory.readByte(z80.HL())
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.memory.writeByte(z80.HL(), bytetemp&0xdf)
					break
				case 0xaf: /* RES 5,A */
					z80.a &= 0xdf
					break
				case 0xb0: /* RES 6,B */
					z80.b &= 0xbf
					break
				case 0xb1: /* RES 6,C */
					z80.c &= 0xbf
					break
				case 0xb2: /* RES 6,D */
					z80.d &= 0xbf
					break
				case 0xb3: /* RES 6,E */
					z80.e &= 0xbf
					break
				case 0xb4: /* RES 6,H */
					z80.h &= 0xbf
					break
				case 0xb5: /* RES 6,L */
					z80.l &= 0xbf
					break
				case 0xb6: /* RES 6,(HL) */
					var bytetemp byte = z80.memory.readByte(z80.HL())
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.memory.writeByte(z80.HL(), bytetemp&0xbf)
					break
				case 0xb7: /* RES 6,A */
					z80.a &= 0xbf
					break
				case 0xb8: /* RES 7,B */
					z80.b &= 0x7f
					break
				case 0xb9: /* RES 7,C */
					z80.c &= 0x7f
					break
				case 0xba: /* RES 7,D */
					z80.d &= 0x7f
					break
				case 0xbb: /* RES 7,E */
					z80.e &= 0x7f
					break
				case 0xbc: /* RES 7,H */
					z80.h &= 0x7f
					break
				case 0xbd: /* RES 7,L */
					z80.l &= 0x7f
					break
				case 0xbe: /* RES 7,(HL) */
					var bytetemp byte = z80.memory.readByte(z80.HL())
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.memory.writeByte(z80.HL(), bytetemp&0x7f)
					break
				case 0xbf: /* RES 7,A */
					z80.a &= 0x7f
					break
				case 0xc0: /* SET 0,B */
					z80.b |= 0x01
					break
				case 0xc1: /* SET 0,C */
					z80.c |= 0x01
					break
				case 0xc2: /* SET 0,D */
					z80.d |= 0x01
					break
				case 0xc3: /* SET 0,E */
					z80.e |= 0x01
					break
				case 0xc4: /* SET 0,H */
					z80.h |= 0x01
					break
				case 0xc5: /* SET 0,L */
					z80.l |= 0x01
					break
				case 0xc6: /* SET 0,(HL) */
					var bytetemp byte = z80.memory.readByte(z80.HL())
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.memory.writeByte(z80.HL(), bytetemp|0x01)
					break
				case 0xc7: /* SET 0,A */
					z80.a |= 0x01
					break
				case 0xc8: /* SET 1,B */
					z80.b |= 0x02
					break
				case 0xc9: /* SET 1,C */
					z80.c |= 0x02
					break
				case 0xca: /* SET 1,D */
					z80.d |= 0x02
					break
				case 0xcb: /* SET 1,E */
					z80.e |= 0x02
					break
				case 0xcc: /* SET 1,H */
					z80.h |= 0x02
					break
				case 0xcd: /* SET 1,L */
					z80.l |= 0x02
					break
				case 0xce: /* SET 1,(HL) */
					var bytetemp byte = z80.memory.readByte(z80.HL())
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.memory.writeByte(z80.HL(), bytetemp|0x02)
					break
				case 0xcf: /* SET 1,A */
					z80.a |= 0x02
					break
				case 0xd0: /* SET 2,B */
					z80.b |= 0x04
					break
				case 0xd1: /* SET 2,C */
					z80.c |= 0x04
					break
				case 0xd2: /* SET 2,D */
					z80.d |= 0x04
					break
				case 0xd3: /* SET 2,E */
					z80.e |= 0x04
					break
				case 0xd4: /* SET 2,H */
					z80.h |= 0x04
					break
				case 0xd5: /* SET 2,L */
					z80.l |= 0x04
					break
				case 0xd6: /* SET 2,(HL) */
					var bytetemp byte = z80.memory.readByte(z80.HL())
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.memory.writeByte(z80.HL(), bytetemp|0x04)
					break
				case 0xd7: /* SET 2,A */
					z80.a |= 0x04
					break
				case 0xd8: /* SET 3,B */
					z80.b |= 0x08
					break
				case 0xd9: /* SET 3,C */
					z80.c |= 0x08
					break
				case 0xda: /* SET 3,D */
					z80.d |= 0x08
					break
				case 0xdb: /* SET 3,E */
					z80.e |= 0x08
					break
				case 0xdc: /* SET 3,H */
					z80.h |= 0x08
					break
				case 0xdd: /* SET 3,L */
					z80.l |= 0x08
					break
				case 0xde: /* SET 3,(HL) */
					var bytetemp byte = z80.memory.readByte(z80.HL())
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.memory.writeByte(z80.HL(), bytetemp|0x08)
					break
				case 0xdf: /* SET 3,A */
					z80.a |= 0x08
					break
				case 0xe0: /* SET 4,B */
					z80.b |= 0x10
					break
				case 0xe1: /* SET 4,C */
					z80.c |= 0x10
					break
				case 0xe2: /* SET 4,D */
					z80.d |= 0x10
					break
				case 0xe3: /* SET 4,E */
					z80.e |= 0x10
					break
				case 0xe4: /* SET 4,H */
					z80.h |= 0x10
					break
				case 0xe5: /* SET 4,L */
					z80.l |= 0x10
					break
				case 0xe6: /* SET 4,(HL) */
					var bytetemp byte = z80.memory.readByte(z80.HL())
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.memory.writeByte(z80.HL(), bytetemp|0x10)
					break
				case 0xe7: /* SET 4,A */
					z80.a |= 0x10
					break
				case 0xe8: /* SET 5,B */
					z80.b |= 0x20
					break
				case 0xe9: /* SET 5,C */
					z80.c |= 0x20
					break
				case 0xea: /* SET 5,D */
					z80.d |= 0x20
					break
				case 0xeb: /* SET 5,E */
					z80.e |= 0x20
					break
				case 0xec: /* SET 5,H */
					z80.h |= 0x20
					break
				case 0xed: /* SET 5,L */
					z80.l |= 0x20
					break
				case 0xee: /* SET 5,(HL) */
					var bytetemp byte = z80.memory.readByte(z80.HL())
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.memory.writeByte(z80.HL(), bytetemp|0x20)
					break
				case 0xef: /* SET 5,A */
					z80.a |= 0x20
					break
				case 0xf0: /* SET 6,B */
					z80.b |= 0x40
					break
				case 0xf1: /* SET 6,C */
					z80.c |= 0x40
					break
				case 0xf2: /* SET 6,D */
					z80.d |= 0x40
					break
				case 0xf3: /* SET 6,E */
					z80.e |= 0x40
					break
				case 0xf4: /* SET 6,H */
					z80.h |= 0x40
					break
				case 0xf5: /* SET 6,L */
					z80.l |= 0x40
					break
				case 0xf6: /* SET 6,(HL) */
					var bytetemp byte = z80.memory.readByte(z80.HL())
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.memory.writeByte(z80.HL(), bytetemp|0x40)
					break
				case 0xf7: /* SET 6,A */
					z80.a |= 0x40
					break
				case 0xf8: /* SET 7,B */
					z80.b |= 0x80
					break
				case 0xf9: /* SET 7,C */
					z80.c |= 0x80
					break
				case 0xfa: /* SET 7,D */
					z80.d |= 0x80
					break
				case 0xfb: /* SET 7,E */
					z80.e |= 0x80
					break
				case 0xfc: /* SET 7,H */
					z80.h |= 0x80
					break
				case 0xfd: /* SET 7,L */
					z80.l |= 0x80
					break
				case 0xfe: /* SET 7,(HL) */
					var bytetemp byte = z80.memory.readByte(z80.HL())
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.memory.writeByte(z80.HL(), bytetemp|0x80)
					break
				case 0xff: /* SET 7,A */
					z80.a |= 0x80
					break

				}
			}
			break
		case 0xcc: /* CALL Z,nnnn */
			if (z80.f & FLAG_Z) != 0 {
				z80.call()
			} else {
				z80.memory.contendRead(z80.pc, 3)
				z80.memory.contendRead(z80.pc+1, 3)
				z80.pc += 2
			}
			break
		case 0xcd: /* CALL nnnn */
			z80.call()
			break
		case 0xce: /* ADC A,nn */
			{
				var bytetemp byte = z80.memory.readByte(z80.PC())
				z80.pc++
				z80.adc(bytetemp)
			}
			break
		case 0xcf: /* RST 8 */
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.rst(0x08)
			break
		case 0xd0: /* RET NC */
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			if !((z80.f & FLAG_C) != 0) {
				z80.ret()
			}
			break
		case 0xd1: /* POP DE */
			z80.pop16(&z80.e, &z80.d)
			break
		case 0xd2: /* JP NC,nnnn */
			if (z80.f & FLAG_C) == 0 {
				z80.jp()
			} else {
				z80.memory.contendRead(z80.pc, 3)
				z80.memory.contendRead(z80.pc+1, 3)
				z80.pc += 2
			}
			break
		case 0xd3: /* OUT (nn),A */
			var outtemp uint16
			outtemp = uint16(z80.memory.readByte(z80.pc)) + (uint16(z80.a) << 8)
			z80.pc++
			z80.writePort(outtemp, z80.a)
			break
		case 0xd4: /* CALL NC,nnnn */
			if (z80.f & FLAG_C) == 0 {
				z80.call()
			} else {
				z80.memory.contendRead(z80.pc, 3)
				z80.memory.contendRead(z80.pc+1, 3)
				z80.pc += 2
			}
			break
		case 0xd5: /* PUSH DE */
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.push16(z80.e, z80.d)
			break
		case 0xd6: /* SUB nn */
			{
				var bytetemp byte = z80.memory.readByte(z80.PC())
				z80.pc++
				z80.sub(bytetemp)
			}
			break
		case 0xd7: /* RST 10 */
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.rst(0x10)
			break
		case 0xd8: /* RET C */
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			if (z80.f & FLAG_C) != 0 {
				z80.ret()
			}
			break
		case 0xd9: /* EXX */
			var wordtemp uint16
			wordtemp = z80.BC()
			z80.setBC(z80.BC_())
			z80.setBC_(wordtemp)

			wordtemp = z80.DE()
			z80.setDE(z80.DE_())
			z80.setDE_(wordtemp)

			wordtemp = z80.HL()
			z80.setHL(z80.HL_())
			z80.setHL_(wordtemp)
			break
		case 0xda: /* JP C,nnnn */
			if (z80.f & FLAG_C) != 0 {
				z80.jp()
			} else {
				z80.memory.contendRead(z80.pc, 3)
				z80.memory.contendRead(z80.pc+1, 3)
				z80.pc += 2
			}
			break
		case 0xdb: /* IN A,(nn) */
			var intemp uint16
			intemp = uint16(z80.memory.readByte(z80.pc)) + (uint16(z80.a) << 8)
			z80.pc++
			z80.a = z80.readPort(intemp)
			break
		case 0xdc: /* CALL C,nnnn */
			if (z80.f & FLAG_C) != 0 {
				z80.call()
			} else {
				z80.memory.contendRead(z80.pc, 3)
				z80.memory.contendRead(z80.pc+1, 3)
				z80.pc += 2
			}
			break
		case 0xdd: /* shift DD */
			{
				var opcode2 byte
				z80.memory.contendRead(z80.pc, 4)
				opcode2 = z80.memory.readByteInternal(z80.pc)
				z80.pc++
				z80.r++

				switch opcode2 {

				/* z80_ddfd.c Z80 {DD,FD}xx opcodes
				   Copyright (c) 1999-2003 Philip Kendall

				   This program is free software; you can redistribute it and/or modify
				   it under the terms of the GNU General Public License as published by
				   the Free Software Foundation; either version 2 of the License, or
				   (at your option) any later version.

				   This program is distributed in the hope that it will be useful,
				   but WITHOUT ANY WARRANTY; without even the implied warranty of
				   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
				   GNU General Public License for more details.

				   You should have received a copy of the GNU General Public License along
				   with this program; if not, write to the Free Software Foundation, Inc.,
				   51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.

				   Author contact information:

				   E-mail: philip-fuse@shadowmagic.org.uk

				*/

				/* NB: this file is autogenerated by './z80.pl' from 'opcodes_ddfd.dat',
				   and included in 'z80_ops.c' */

				case 0x09: /* ADD ix,BC */
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.add16(z80.ix, z80.BC())
					break
				case 0x19: /* ADD ix,DE */
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.add16(z80.ix, z80.DE())
					break
				case 0x21: /* LD ix,nnnn */
					b1 := z80.memory.readByte(z80.pc)
					z80.pc++
					b2 := z80.memory.readByte(z80.pc)
					z80.pc++
					z80.setIX(joinBytes(b2, b1))
					break
				case 0x22: /* LD (nnnn),ix */
					z80.ld16nnrr(z80.ixl, z80.ixh)
					break
					break
				case 0x23: /* INC ix */
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.incIX()
					break
				case 0x24: /* INC z80.IXH() */
					z80.incIXH()
					break
				case 0x25: /* DEC z80.IXH() */
					z80.decIXH()
					break
				case 0x26: /* LD z80.IXH(),nn */
					z80.ixh = z80.memory.readByte(z80.pc)
					z80.pc++
					break
				case 0x29: /* ADD ix,ix */
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.add16(z80.ix, z80.IX())
					break
				case 0x2a: /* LD ix,(nnnn) */
					z80.ld16rrnn(&z80.ixl, &z80.ixh)
					break
					break
				case 0x2b: /* DEC ix */
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.decIX()
					break
				case 0x2c: /* INC z80.IXL() */
					z80.incIXL()
					break
				case 0x2d: /* DEC z80.IXL() */
					z80.decIXL()
					break
				case 0x2e: /* LD z80.IXL(),nn */
					z80.ixl = z80.memory.readByte(z80.pc)
					z80.pc++
					break
				case 0x34: /* INC (ix+dd) */
					var offset, bytetemp byte
					var wordtemp uint16
					offset = z80.memory.readByte(z80.pc)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.pc++
					wordtemp = uint16(int(z80.IX()) + int(signExtend(offset)))
					bytetemp = z80.memory.readByte(wordtemp)
					z80.memory.contendReadNoMreq(wordtemp, 1)
					z80.inc(&bytetemp)
					z80.memory.writeByte(wordtemp, bytetemp)
					break
				case 0x35: /* DEC (ix+dd) */
					var offset, bytetemp byte
					var wordtemp uint16
					offset = z80.memory.readByte(z80.pc)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.pc++
					wordtemp = uint16(int(z80.IX()) + int(signExtend(offset)))
					bytetemp = z80.memory.readByte(wordtemp)
					z80.memory.contendReadNoMreq(wordtemp, 1)
					z80.dec(&bytetemp)
					z80.memory.writeByte(wordtemp, bytetemp)
					break
				case 0x36: /* LD (ix+dd),nn */
					offset := z80.memory.readByte(z80.pc)
					z80.pc++
					value := z80.memory.readByte(z80.pc)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.pc++
					z80.memory.writeByte(uint16(int(z80.IX())+int(signExtend(offset))), value)
					break
				case 0x39: /* ADD ix,SP */
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.add16(z80.ix, z80.SP())
					break
				case 0x44: /* LD B,z80.IXH() */
					z80.b = z80.ixh
					break
				case 0x45: /* LD B,z80.IXL() */
					z80.b = z80.ixl
					break
				case 0x46: /* LD B,(ix+dd) */
					var offset byte
					offset = z80.memory.readByte(z80.pc)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.pc++
					z80.b = z80.memory.readByte(uint16(int(z80.IX()) + int(signExtend(offset))))
					break
				case 0x4c: /* LD C,z80.IXH() */
					z80.c = z80.ixh
					break
				case 0x4d: /* LD C,z80.IXL() */
					z80.c = z80.ixl
					break
				case 0x4e: /* LD C,(ix+dd) */
					var offset byte
					offset = z80.memory.readByte(z80.pc)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.pc++
					z80.c = z80.memory.readByte(uint16(int(z80.IX()) + int(signExtend(offset))))
					break
				case 0x54: /* LD D,z80.IXH() */
					z80.d = z80.ixh
					break
				case 0x55: /* LD D,z80.IXL() */
					z80.d = z80.ixl
					break
				case 0x56: /* LD D,(ix+dd) */
					var offset byte
					offset = z80.memory.readByte(z80.pc)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.pc++
					z80.d = z80.memory.readByte(uint16(int(z80.IX()) + int(signExtend(offset))))
					break
				case 0x5c: /* LD E,z80.IXH() */
					z80.e = z80.ixh
					break
				case 0x5d: /* LD E,z80.IXL() */
					z80.e = z80.ixl
					break
				case 0x5e: /* LD E,(ix+dd) */
					var offset byte
					offset = z80.memory.readByte(z80.pc)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.pc++
					z80.e = z80.memory.readByte(uint16(int(z80.IX()) + int(signExtend(offset))))
					break
				case 0x60: /* LD z80.IXH(),B */
					z80.ixh = z80.b
					break
				case 0x61: /* LD z80.IXH(),C */
					z80.ixh = z80.c
					break
				case 0x62: /* LD z80.IXH(),D */
					z80.ixh = z80.d
					break
				case 0x63: /* LD z80.IXH(),E */
					z80.ixh = z80.e
					break
				case 0x64: /* LD z80.IXH(),z80.IXH() */
					break
				case 0x65: /* LD z80.IXH(),z80.IXL() */
					z80.ixh = z80.ixl
					break
				case 0x66: /* LD H,(ix+dd) */
					var offset byte
					offset = z80.memory.readByte(z80.pc)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.pc++
					z80.h = z80.memory.readByte(uint16(int(z80.IX()) + int(signExtend(offset))))
					break
				case 0x67: /* LD z80.IXH(),A */
					z80.ixh = z80.a
					break
				case 0x68: /* LD z80.IXL(),B */
					z80.ixl = z80.b
					break
				case 0x69: /* LD z80.IXL(),C */
					z80.ixl = z80.c
					break
				case 0x6a: /* LD z80.IXL(),D */
					z80.ixl = z80.d
					break
				case 0x6b: /* LD z80.IXL(),E */
					z80.ixl = z80.e
					break
				case 0x6c: /* LD z80.IXL(),z80.IXH() */
					z80.ixl = z80.ixh
					break
				case 0x6d: /* LD z80.IXL(),z80.IXL() */
					break
				case 0x6e: /* LD L,(ix+dd) */
					var offset byte
					offset = z80.memory.readByte(z80.pc)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.pc++
					z80.l = z80.memory.readByte(uint16(int(z80.IX()) + int(signExtend(offset))))
					break
				case 0x6f: /* LD z80.IXL(),A */
					z80.ixl = z80.a
					break
				case 0x70: /* LD (ix+dd),B */
					offset := z80.memory.readByte(z80.pc)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.pc++
					z80.memory.writeByte(uint16(int(z80.IX())+int(signExtend(offset))), z80.b)
					break
				case 0x71: /* LD (ix+dd),C */
					offset := z80.memory.readByte(z80.pc)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.pc++
					z80.memory.writeByte(uint16(int(z80.IX())+int(signExtend(offset))), z80.c)
					break
				case 0x72: /* LD (ix+dd),D */
					offset := z80.memory.readByte(z80.pc)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.pc++
					z80.memory.writeByte(uint16(int(z80.IX())+int(signExtend(offset))), z80.d)
					break
				case 0x73: /* LD (ix+dd),E */
					offset := z80.memory.readByte(z80.pc)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.pc++
					z80.memory.writeByte(uint16(int(z80.IX())+int(signExtend(offset))), z80.e)
					break
				case 0x74: /* LD (ix+dd),H */
					offset := z80.memory.readByte(z80.pc)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.pc++
					z80.memory.writeByte(uint16(int(z80.IX())+int(signExtend(offset))), z80.h)
					break
				case 0x75: /* LD (ix+dd),L */
					offset := z80.memory.readByte(z80.pc)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.pc++
					z80.memory.writeByte(uint16(int(z80.IX())+int(signExtend(offset))), z80.l)
					break
				case 0x77: /* LD (ix+dd),A */
					offset := z80.memory.readByte(z80.pc)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.pc++
					z80.memory.writeByte(uint16(int(z80.IX())+int(signExtend(offset))), z80.a)
					break
				case 0x7c: /* LD A,z80.IXH() */
					z80.a = z80.ixh
					break
				case 0x7d: /* LD A,z80.IXL() */
					z80.a = z80.ixl
					break
				case 0x7e: /* LD A,(ix+dd) */
					var offset byte
					offset = z80.memory.readByte(z80.pc)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.pc++
					z80.a = z80.memory.readByte(uint16(int(z80.IX()) + int(signExtend(offset))))
					break
				case 0x84: /* ADD A,z80.IXH() */
					z80.add(z80.ixh)
					break
				case 0x85: /* ADD A,z80.IXL() */
					z80.add(z80.ixl)
					break
				case 0x86: /* ADD A,(ix+dd) */

					var offset, bytetemp byte
					offset = z80.memory.readByte(z80.pc)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.pc++
					bytetemp = z80.memory.readByte(uint16(int(z80.IX()) + int(signExtend(offset))))
					z80.add(bytetemp)

					break
				case 0x8c: /* ADC A,z80.IXH() */
					z80.adc(z80.ixh)
					break
				case 0x8d: /* ADC A,z80.IXL() */
					z80.adc(z80.ixl)
					break
				case 0x8e: /* ADC A,(ix+dd) */

					var offset, bytetemp byte
					offset = z80.memory.readByte(z80.pc)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.pc++
					bytetemp = z80.memory.readByte(uint16(int(z80.IX()) + int(signExtend(offset))))
					z80.adc(bytetemp)

					break
				case 0x94: /* SUB A,z80.IXH() */
					z80.sub(z80.ixh)
					break
				case 0x95: /* SUB A,z80.IXL() */
					z80.sub(z80.ixl)
					break
				case 0x96: /* SUB A,(ix+dd) */

					var offset, bytetemp byte
					offset = z80.memory.readByte(z80.pc)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.pc++
					bytetemp = z80.memory.readByte(uint16(int(z80.IX()) + int(signExtend(offset))))
					z80.sub(bytetemp)

					break
				case 0x9c: /* SBC A,z80.IXH() */
					z80.sbc(z80.ixh)
					break
				case 0x9d: /* SBC A,z80.IXL() */
					z80.sbc(z80.ixl)
					break
				case 0x9e: /* SBC A,(ix+dd) */

					var offset, bytetemp byte
					offset = z80.memory.readByte(z80.pc)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.pc++
					bytetemp = z80.memory.readByte(uint16(int(z80.IX()) + int(signExtend(offset))))
					z80.sbc(bytetemp)

					break
				case 0xa4: /* AND A,z80.IXH() */
					z80.and(z80.ixh)
					break
				case 0xa5: /* AND A,z80.IXL() */
					z80.and(z80.ixl)
					break
				case 0xa6: /* AND A,(ix+dd) */

					var offset, bytetemp byte
					offset = z80.memory.readByte(z80.pc)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.pc++
					bytetemp = z80.memory.readByte(uint16(int(z80.IX()) + int(signExtend(offset))))
					z80.and(bytetemp)

					break
				case 0xac: /* XOR A,z80.IXH() */
					z80.xor(z80.ixh)
					break
				case 0xad: /* XOR A,z80.IXL() */
					z80.xor(z80.ixl)
					break
				case 0xae: /* XOR A,(ix+dd) */

					var offset, bytetemp byte
					offset = z80.memory.readByte(z80.pc)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.pc++
					bytetemp = z80.memory.readByte(uint16(int(z80.IX()) + int(signExtend(offset))))
					z80.xor(bytetemp)

					break
				case 0xb4: /* OR A,z80.IXH() */
					z80.or(z80.ixh)
					break
				case 0xb5: /* OR A,z80.IXL() */
					z80.or(z80.ixl)
					break
				case 0xb6: /* OR A,(ix+dd) */

					var offset, bytetemp byte
					offset = z80.memory.readByte(z80.pc)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.pc++
					bytetemp = z80.memory.readByte(uint16(int(z80.IX()) + int(signExtend(offset))))
					z80.or(bytetemp)

					break
				case 0xbc: /* CP A,z80.IXH() */
					z80.cp(z80.ixh)
					break
				case 0xbd: /* CP A,z80.IXL() */
					z80.cp(z80.ixl)
					break
				case 0xbe: /* CP A,(ix+dd) */

					var offset, bytetemp byte
					offset = z80.memory.readByte(z80.pc)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.pc++
					bytetemp = z80.memory.readByte(uint16(int(z80.IX()) + int(signExtend(offset))))
					z80.cp(bytetemp)

					break
				case 0xcb: /* shift DDFDCB */

					var tempaddr uint16
					var opcode3 byte
					z80.memory.contendRead(z80.pc, 3)
					tempaddr = uint16(int(z80.IX()) + int(signExtend(z80.memory.readByteInternal(z80.pc))))
					z80.pc++
					z80.memory.contendRead(z80.pc, 3)
					opcode3 = z80.memory.readByteInternal(z80.pc)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.pc++

					switch opcode3 {
					/* z80_ddfdcb.c Z80 {DD,FD}CBxx opcodes
					   Copyright (c) 1999-2003 Philip Kendall

					   This program is free software; you can redistribute it and/or modify
					   it under the terms of the GNU General Public License as published by
					   the Free Software Foundation; either version 2 of the License, or
					   (at your option) any later version.

					   This program is distributed in the hope that it will be useful,
					   but WITHOUT ANY WARRANTY; without even the implied warranty of
					   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
					   GNU General Public License for more details.

					   You should have received a copy of the GNU General Public License along
					   with this program; if not, write to the Free Software Foundation, Inc.,
					   51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.

					   Author contact information:

					   E-mail: philip-fuse@shadowmagic.org.uk

					*/

					/* NB: this file is autogenerated by './z80.pl' from 'opcodes_ddfdcb.dat',
					   and included in 'z80_ops.c' */

					case 0x00: /* LD B,RLC (REGISTER+dd) */
						z80.b = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rlc(&z80.b)
						z80.memory.writeByte(tempaddr, z80.b)
						break
					case 0x01: /* LD C,RLC (REGISTER+dd) */
						z80.c = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rlc(&z80.c)
						z80.memory.writeByte(tempaddr, z80.c)
						break
					case 0x02: /* LD D,RLC (REGISTER+dd) */
						z80.d = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rlc(&z80.d)
						z80.memory.writeByte(tempaddr, z80.d)
						break
					case 0x03: /* LD E,RLC (REGISTER+dd) */
						z80.e = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rlc(&z80.e)
						z80.memory.writeByte(tempaddr, z80.e)
						break
					case 0x04: /* LD H,RLC (REGISTER+dd) */
						z80.h = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rlc(&z80.h)
						z80.memory.writeByte(tempaddr, z80.h)
						break
					case 0x05: /* LD L,RLC (REGISTER+dd) */
						z80.l = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rlc(&z80.l)
						z80.memory.writeByte(tempaddr, z80.l)
						break
					case 0x06: /* RLC (REGISTER+dd) */
						var bytetemp byte = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rlc(&bytetemp)
						z80.memory.writeByte(tempaddr, bytetemp)
						break
					case 0x07: /* LD A,RLC (REGISTER+dd) */
						z80.a = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rlc(&z80.a)
						z80.memory.writeByte(tempaddr, z80.a)
						break
					case 0x08: /* LD B,RRC (REGISTER+dd) */
						z80.b = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rrc(&z80.b)
						z80.memory.writeByte(tempaddr, z80.b)
						break
					case 0x09: /* LD C,RRC (REGISTER+dd) */
						z80.c = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rrc(&z80.c)
						z80.memory.writeByte(tempaddr, z80.c)
						break
					case 0x0a: /* LD D,RRC (REGISTER+dd) */
						z80.d = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rrc(&z80.d)
						z80.memory.writeByte(tempaddr, z80.d)
						break
					case 0x0b: /* LD E,RRC (REGISTER+dd) */
						z80.e = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rrc(&z80.e)
						z80.memory.writeByte(tempaddr, z80.e)
						break
					case 0x0c: /* LD H,RRC (REGISTER+dd) */
						z80.h = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rrc(&z80.h)
						z80.memory.writeByte(tempaddr, z80.h)
						break
					case 0x0d: /* LD L,RRC (REGISTER+dd) */
						z80.l = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rrc(&z80.l)
						z80.memory.writeByte(tempaddr, z80.l)
						break
					case 0x0e: /* RRC (REGISTER+dd) */
						var bytetemp byte = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rrc(&bytetemp)
						z80.memory.writeByte(tempaddr, bytetemp)
						break
					case 0x0f: /* LD A,RRC (REGISTER+dd) */
						z80.a = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rrc(&z80.a)
						z80.memory.writeByte(tempaddr, z80.a)
						break
					case 0x10: /* LD B,RL (REGISTER+dd) */
						z80.b = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rl(&z80.b)
						z80.memory.writeByte(tempaddr, z80.b)
						break
					case 0x11: /* LD C,RL (REGISTER+dd) */
						z80.c = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rl(&z80.c)
						z80.memory.writeByte(tempaddr, z80.c)
						break
					case 0x12: /* LD D,RL (REGISTER+dd) */
						z80.d = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rl(&z80.d)
						z80.memory.writeByte(tempaddr, z80.d)
						break
					case 0x13: /* LD E,RL (REGISTER+dd) */
						z80.e = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rl(&z80.e)
						z80.memory.writeByte(tempaddr, z80.e)
						break
					case 0x14: /* LD H,RL (REGISTER+dd) */
						z80.h = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rl(&z80.h)
						z80.memory.writeByte(tempaddr, z80.h)
						break
					case 0x15: /* LD L,RL (REGISTER+dd) */
						z80.l = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rl(&z80.l)
						z80.memory.writeByte(tempaddr, z80.l)
						break
					case 0x16: /* RL (REGISTER+dd) */
						var bytetemp byte = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rl(&bytetemp)
						z80.memory.writeByte(tempaddr, bytetemp)
						break
					case 0x17: /* LD A,RL (REGISTER+dd) */
						z80.a = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rl(&z80.a)
						z80.memory.writeByte(tempaddr, z80.a)
						break
					case 0x18: /* LD B,RR (REGISTER+dd) */
						z80.b = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rr(&z80.b)
						z80.memory.writeByte(tempaddr, z80.b)
						break
					case 0x19: /* LD C,RR (REGISTER+dd) */
						z80.c = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rr(&z80.c)
						z80.memory.writeByte(tempaddr, z80.c)
						break
					case 0x1a: /* LD D,RR (REGISTER+dd) */
						z80.d = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rr(&z80.d)
						z80.memory.writeByte(tempaddr, z80.d)
						break
					case 0x1b: /* LD E,RR (REGISTER+dd) */
						z80.e = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rr(&z80.e)
						z80.memory.writeByte(tempaddr, z80.e)
						break
					case 0x1c: /* LD H,RR (REGISTER+dd) */
						z80.h = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rr(&z80.h)
						z80.memory.writeByte(tempaddr, z80.h)
						break
					case 0x1d: /* LD L,RR (REGISTER+dd) */
						z80.l = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rr(&z80.l)
						z80.memory.writeByte(tempaddr, z80.l)
						break
					case 0x1e: /* RR (REGISTER+dd) */
						var bytetemp byte = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rr(&bytetemp)
						z80.memory.writeByte(tempaddr, bytetemp)
						break
					case 0x1f: /* LD A,RR (REGISTER+dd) */
						z80.a = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rr(&z80.a)
						z80.memory.writeByte(tempaddr, z80.a)
						break
					case 0x20: /* LD B,SLA (REGISTER+dd) */
						z80.b = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.sla(&z80.b)
						z80.memory.writeByte(tempaddr, z80.b)
						break
					case 0x21: /* LD C,SLA (REGISTER+dd) */
						z80.c = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.sla(&z80.c)
						z80.memory.writeByte(tempaddr, z80.c)
						break
					case 0x22: /* LD D,SLA (REGISTER+dd) */
						z80.d = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.sla(&z80.d)
						z80.memory.writeByte(tempaddr, z80.d)
						break
					case 0x23: /* LD E,SLA (REGISTER+dd) */
						z80.e = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.sla(&z80.e)
						z80.memory.writeByte(tempaddr, z80.e)
						break
					case 0x24: /* LD H,SLA (REGISTER+dd) */
						z80.h = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.sla(&z80.h)
						z80.memory.writeByte(tempaddr, z80.h)
						break
					case 0x25: /* LD L,SLA (REGISTER+dd) */
						z80.l = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.sla(&z80.l)
						z80.memory.writeByte(tempaddr, z80.l)
						break
					case 0x26: /* SLA (REGISTER+dd) */
						var bytetemp byte = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.sla(&bytetemp)
						z80.memory.writeByte(tempaddr, bytetemp)
						break
					case 0x27: /* LD A,SLA (REGISTER+dd) */
						z80.a = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.sla(&z80.a)
						z80.memory.writeByte(tempaddr, z80.a)
						break
					case 0x28: /* LD B,SRA (REGISTER+dd) */
						z80.b = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.sra(&z80.b)
						z80.memory.writeByte(tempaddr, z80.b)
						break
					case 0x29: /* LD C,SRA (REGISTER+dd) */
						z80.c = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.sra(&z80.c)
						z80.memory.writeByte(tempaddr, z80.c)
						break
					case 0x2a: /* LD D,SRA (REGISTER+dd) */
						z80.d = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.sra(&z80.d)
						z80.memory.writeByte(tempaddr, z80.d)
						break
					case 0x2b: /* LD E,SRA (REGISTER+dd) */
						z80.e = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.sra(&z80.e)
						z80.memory.writeByte(tempaddr, z80.e)
						break
					case 0x2c: /* LD H,SRA (REGISTER+dd) */
						z80.h = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.sra(&z80.h)
						z80.memory.writeByte(tempaddr, z80.h)
						break
					case 0x2d: /* LD L,SRA (REGISTER+dd) */
						z80.l = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.sra(&z80.l)
						z80.memory.writeByte(tempaddr, z80.l)
						break
					case 0x2e: /* SRA (REGISTER+dd) */
						var bytetemp byte = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.sra(&bytetemp)
						z80.memory.writeByte(tempaddr, bytetemp)
						break
					case 0x2f: /* LD A,SRA (REGISTER+dd) */
						z80.a = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.sra(&z80.a)
						z80.memory.writeByte(tempaddr, z80.a)
						break
					case 0x30: /* LD B,SLL (REGISTER+dd) */
						z80.b = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.sll(&z80.b)
						z80.memory.writeByte(tempaddr, z80.b)
						break
					case 0x31: /* LD C,SLL (REGISTER+dd) */
						z80.c = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.sll(&z80.c)
						z80.memory.writeByte(tempaddr, z80.c)
						break
					case 0x32: /* LD D,SLL (REGISTER+dd) */
						z80.d = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.sll(&z80.d)
						z80.memory.writeByte(tempaddr, z80.d)
						break
					case 0x33: /* LD E,SLL (REGISTER+dd) */
						z80.e = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.sll(&z80.e)
						z80.memory.writeByte(tempaddr, z80.e)
						break
					case 0x34: /* LD H,SLL (REGISTER+dd) */
						z80.h = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.sll(&z80.h)
						z80.memory.writeByte(tempaddr, z80.h)
						break
					case 0x35: /* LD L,SLL (REGISTER+dd) */
						z80.l = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.sll(&z80.l)
						z80.memory.writeByte(tempaddr, z80.l)
						break
					case 0x36: /* SLL (REGISTER+dd) */
						var bytetemp byte = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.sll(&bytetemp)
						z80.memory.writeByte(tempaddr, bytetemp)
						break
					case 0x37: /* LD A,SLL (REGISTER+dd) */
						z80.a = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.sll(&z80.a)
						z80.memory.writeByte(tempaddr, z80.a)
						break
					case 0x38: /* LD B,SRL (REGISTER+dd) */
						z80.b = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.srl(&z80.b)
						z80.memory.writeByte(tempaddr, z80.b)
						break
					case 0x39: /* LD C,SRL (REGISTER+dd) */
						z80.c = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.srl(&z80.c)
						z80.memory.writeByte(tempaddr, z80.c)
						break
					case 0x3a: /* LD D,SRL (REGISTER+dd) */
						z80.d = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.srl(&z80.d)
						z80.memory.writeByte(tempaddr, z80.d)
						break
					case 0x3b: /* LD E,SRL (REGISTER+dd) */
						z80.e = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.srl(&z80.e)
						z80.memory.writeByte(tempaddr, z80.e)
						break
					case 0x3c: /* LD H,SRL (REGISTER+dd) */
						z80.h = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.srl(&z80.h)
						z80.memory.writeByte(tempaddr, z80.h)
						break
					case 0x3d: /* LD L,SRL (REGISTER+dd) */
						z80.l = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.srl(&z80.l)
						z80.memory.writeByte(tempaddr, z80.l)
						break
					case 0x3e: /* SRL (REGISTER+dd) */
						var bytetemp byte = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.srl(&bytetemp)
						z80.memory.writeByte(tempaddr, bytetemp)
						break
					case 0x3f: /* LD A,SRL (REGISTER+dd) */
						z80.a = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.srl(&z80.a)
						z80.memory.writeByte(tempaddr, z80.a)
						break
					case 0x40:
						fallthrough
					case 0x41:
						fallthrough
					case 0x42:
						fallthrough
					case 0x43:
						fallthrough
					case 0x44:
						fallthrough
					case 0x45:
						fallthrough
					case 0x46:
						fallthrough
					case 0x47: /* BIT 0,(REGISTER+dd) */
						bytetemp := z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.biti(0, bytetemp, tempaddr)
						break
					case 0x48:
						fallthrough
					case 0x49:
						fallthrough
					case 0x4a:
						fallthrough
					case 0x4b:
						fallthrough
					case 0x4c:
						fallthrough
					case 0x4d:
						fallthrough
					case 0x4e:
						fallthrough
					case 0x4f: /* BIT 1,(REGISTER+dd) */
						bytetemp := z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.biti(1, bytetemp, tempaddr)
						break
					case 0x50:
						fallthrough
					case 0x51:
						fallthrough
					case 0x52:
						fallthrough
					case 0x53:
						fallthrough
					case 0x54:
						fallthrough
					case 0x55:
						fallthrough
					case 0x56:
						fallthrough
					case 0x57: /* BIT 2,(REGISTER+dd) */
						bytetemp := z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.biti(2, bytetemp, tempaddr)
						break
					case 0x58:
						fallthrough
					case 0x59:
						fallthrough
					case 0x5a:
						fallthrough
					case 0x5b:
						fallthrough
					case 0x5c:
						fallthrough
					case 0x5d:
						fallthrough
					case 0x5e:
						fallthrough
					case 0x5f: /* BIT 3,(REGISTER+dd) */
						bytetemp := z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.biti(3, bytetemp, tempaddr)
						break
					case 0x60:
						fallthrough
					case 0x61:
						fallthrough
					case 0x62:
						fallthrough
					case 0x63:
						fallthrough
					case 0x64:
						fallthrough
					case 0x65:
						fallthrough
					case 0x66:
						fallthrough
					case 0x67: /* BIT 4,(REGISTER+dd) */
						bytetemp := z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.biti(4, bytetemp, tempaddr)
						break
					case 0x68:
						fallthrough
					case 0x69:
						fallthrough
					case 0x6a:
						fallthrough
					case 0x6b:
						fallthrough
					case 0x6c:
						fallthrough
					case 0x6d:
						fallthrough
					case 0x6e:
						fallthrough
					case 0x6f: /* BIT 5,(REGISTER+dd) */
						bytetemp := z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.biti(5, bytetemp, tempaddr)
						break
					case 0x70:
						fallthrough
					case 0x71:
						fallthrough
					case 0x72:
						fallthrough
					case 0x73:
						fallthrough
					case 0x74:
						fallthrough
					case 0x75:
						fallthrough
					case 0x76:
						fallthrough
					case 0x77: /* BIT 6,(REGISTER+dd) */
						bytetemp := z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.biti(6, bytetemp, tempaddr)
						break
					case 0x78:
						fallthrough
					case 0x79:
						fallthrough
					case 0x7a:
						fallthrough
					case 0x7b:
						fallthrough
					case 0x7c:
						fallthrough
					case 0x7d:
						fallthrough
					case 0x7e:
						fallthrough
					case 0x7f: /* BIT 7,(REGISTER+dd) */
						bytetemp := z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.biti(7, bytetemp, tempaddr)
						break
					case 0x80: /* LD B,RES 0,(REGISTER+dd) */
						z80.b = z80.memory.readByte(tempaddr) & 0xfe
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.b)
						break
					case 0x81: /* LD C,RES 0,(REGISTER+dd) */
						z80.c = z80.memory.readByte(tempaddr) & 0xfe
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.c)
						break
					case 0x82: /* LD D,RES 0,(REGISTER+dd) */
						z80.d = z80.memory.readByte(tempaddr) & 0xfe
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.d)
						break
					case 0x83: /* LD E,RES 0,(REGISTER+dd) */
						z80.e = z80.memory.readByte(tempaddr) & 0xfe
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.e)
						break
					case 0x84: /* LD H,RES 0,(REGISTER+dd) */
						z80.h = z80.memory.readByte(tempaddr) & 0xfe
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.h)
						break
					case 0x85: /* LD L,RES 0,(REGISTER+dd) */
						z80.l = z80.memory.readByte(tempaddr) & 0xfe
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.l)
						break
					case 0x86: /* RES 0,(REGISTER+dd) */

						var bytetemp byte
						bytetemp = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, bytetemp&0xfe)

						break
					case 0x87: /* LD A,RES 0,(REGISTER+dd) */
						z80.a = z80.memory.readByte(tempaddr) & 0xfe
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.a)
						break
					case 0x88: /* LD B,RES 1,(REGISTER+dd) */
						z80.b = z80.memory.readByte(tempaddr) & 0xfd
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.b)
						break
					case 0x89: /* LD C,RES 1,(REGISTER+dd) */
						z80.c = z80.memory.readByte(tempaddr) & 0xfd
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.c)
						break
					case 0x8a: /* LD D,RES 1,(REGISTER+dd) */
						z80.d = z80.memory.readByte(tempaddr) & 0xfd
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.d)
						break
					case 0x8b: /* LD E,RES 1,(REGISTER+dd) */
						z80.e = z80.memory.readByte(tempaddr) & 0xfd
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.e)
						break
					case 0x8c: /* LD H,RES 1,(REGISTER+dd) */
						z80.h = z80.memory.readByte(tempaddr) & 0xfd
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.h)
						break
					case 0x8d: /* LD L,RES 1,(REGISTER+dd) */
						z80.l = z80.memory.readByte(tempaddr) & 0xfd
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.l)
						break
					case 0x8e: /* RES 1,(REGISTER+dd) */

						var bytetemp byte
						bytetemp = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, bytetemp&0xfd)

						break
					case 0x8f: /* LD A,RES 1,(REGISTER+dd) */
						z80.a = z80.memory.readByte(tempaddr) & 0xfd
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.a)
						break
					case 0x90: /* LD B,RES 2,(REGISTER+dd) */
						z80.b = z80.memory.readByte(tempaddr) & 0xfb
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.b)
						break
					case 0x91: /* LD C,RES 2,(REGISTER+dd) */
						z80.c = z80.memory.readByte(tempaddr) & 0xfb
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.c)
						break
					case 0x92: /* LD D,RES 2,(REGISTER+dd) */
						z80.d = z80.memory.readByte(tempaddr) & 0xfb
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.d)
						break
					case 0x93: /* LD E,RES 2,(REGISTER+dd) */
						z80.e = z80.memory.readByte(tempaddr) & 0xfb
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.e)
						break
					case 0x94: /* LD H,RES 2,(REGISTER+dd) */
						z80.h = z80.memory.readByte(tempaddr) & 0xfb
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.h)
						break
					case 0x95: /* LD L,RES 2,(REGISTER+dd) */
						z80.l = z80.memory.readByte(tempaddr) & 0xfb
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.l)
						break
					case 0x96: /* RES 2,(REGISTER+dd) */

						var bytetemp byte
						bytetemp = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, bytetemp&0xfb)

						break
					case 0x97: /* LD A,RES 2,(REGISTER+dd) */
						z80.a = z80.memory.readByte(tempaddr) & 0xfb
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.a)
						break
					case 0x98: /* LD B,RES 3,(REGISTER+dd) */
						z80.b = z80.memory.readByte(tempaddr) & 0xf7
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.b)
						break
					case 0x99: /* LD C,RES 3,(REGISTER+dd) */
						z80.c = z80.memory.readByte(tempaddr) & 0xf7
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.c)
						break
					case 0x9a: /* LD D,RES 3,(REGISTER+dd) */
						z80.d = z80.memory.readByte(tempaddr) & 0xf7
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.d)
						break
					case 0x9b: /* LD E,RES 3,(REGISTER+dd) */
						z80.e = z80.memory.readByte(tempaddr) & 0xf7
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.e)
						break
					case 0x9c: /* LD H,RES 3,(REGISTER+dd) */
						z80.h = z80.memory.readByte(tempaddr) & 0xf7
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.h)
						break
					case 0x9d: /* LD L,RES 3,(REGISTER+dd) */
						z80.l = z80.memory.readByte(tempaddr) & 0xf7
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.l)
						break
					case 0x9e: /* RES 3,(REGISTER+dd) */

						var bytetemp byte
						bytetemp = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, bytetemp&0xf7)

						break
					case 0x9f: /* LD A,RES 3,(REGISTER+dd) */
						z80.a = z80.memory.readByte(tempaddr) & 0xf7
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.a)
						break
					case 0xa0: /* LD B,RES 4,(REGISTER+dd) */
						z80.b = z80.memory.readByte(tempaddr) & 0xef
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.b)
						break
					case 0xa1: /* LD C,RES 4,(REGISTER+dd) */
						z80.c = z80.memory.readByte(tempaddr) & 0xef
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.c)
						break
					case 0xa2: /* LD D,RES 4,(REGISTER+dd) */
						z80.d = z80.memory.readByte(tempaddr) & 0xef
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.d)
						break
					case 0xa3: /* LD E,RES 4,(REGISTER+dd) */
						z80.e = z80.memory.readByte(tempaddr) & 0xef
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.e)
						break
					case 0xa4: /* LD H,RES 4,(REGISTER+dd) */
						z80.h = z80.memory.readByte(tempaddr) & 0xef
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.h)
						break
					case 0xa5: /* LD L,RES 4,(REGISTER+dd) */
						z80.l = z80.memory.readByte(tempaddr) & 0xef
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.l)
						break
					case 0xa6: /* RES 4,(REGISTER+dd) */

						var bytetemp byte
						bytetemp = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, bytetemp&0xef)

						break
					case 0xa7: /* LD A,RES 4,(REGISTER+dd) */
						z80.a = z80.memory.readByte(tempaddr) & 0xef
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.a)
						break
					case 0xa8: /* LD B,RES 5,(REGISTER+dd) */
						z80.b = z80.memory.readByte(tempaddr) & 0xdf
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.b)
						break
					case 0xa9: /* LD C,RES 5,(REGISTER+dd) */
						z80.c = z80.memory.readByte(tempaddr) & 0xdf
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.c)
						break
					case 0xaa: /* LD D,RES 5,(REGISTER+dd) */
						z80.d = z80.memory.readByte(tempaddr) & 0xdf
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.d)
						break
					case 0xab: /* LD E,RES 5,(REGISTER+dd) */
						z80.e = z80.memory.readByte(tempaddr) & 0xdf
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.e)
						break
					case 0xac: /* LD H,RES 5,(REGISTER+dd) */
						z80.h = z80.memory.readByte(tempaddr) & 0xdf
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.h)
						break
					case 0xad: /* LD L,RES 5,(REGISTER+dd) */
						z80.l = z80.memory.readByte(tempaddr) & 0xdf
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.l)
						break
					case 0xae: /* RES 5,(REGISTER+dd) */

						var bytetemp byte
						bytetemp = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, bytetemp&0xdf)

						break
					case 0xaf: /* LD A,RES 5,(REGISTER+dd) */
						z80.a = z80.memory.readByte(tempaddr) & 0xdf
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.a)
						break
					case 0xb0: /* LD B,RES 6,(REGISTER+dd) */
						z80.b = z80.memory.readByte(tempaddr) & 0xbf
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.b)
						break
					case 0xb1: /* LD C,RES 6,(REGISTER+dd) */
						z80.c = z80.memory.readByte(tempaddr) & 0xbf
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.c)
						break
					case 0xb2: /* LD D,RES 6,(REGISTER+dd) */
						z80.d = z80.memory.readByte(tempaddr) & 0xbf
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.d)
						break
					case 0xb3: /* LD E,RES 6,(REGISTER+dd) */
						z80.e = z80.memory.readByte(tempaddr) & 0xbf
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.e)
						break
					case 0xb4: /* LD H,RES 6,(REGISTER+dd) */
						z80.h = z80.memory.readByte(tempaddr) & 0xbf
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.h)
						break
					case 0xb5: /* LD L,RES 6,(REGISTER+dd) */
						z80.l = z80.memory.readByte(tempaddr) & 0xbf
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.l)
						break
					case 0xb6: /* RES 6,(REGISTER+dd) */

						var bytetemp byte
						bytetemp = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, bytetemp&0xbf)

						break
					case 0xb7: /* LD A,RES 6,(REGISTER+dd) */
						z80.a = z80.memory.readByte(tempaddr) & 0xbf
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.a)
						break
					case 0xb8: /* LD B,RES 7,(REGISTER+dd) */
						z80.b = z80.memory.readByte(tempaddr) & 0x7f
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.b)
						break
					case 0xb9: /* LD C,RES 7,(REGISTER+dd) */
						z80.c = z80.memory.readByte(tempaddr) & 0x7f
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.c)
						break
					case 0xba: /* LD D,RES 7,(REGISTER+dd) */
						z80.d = z80.memory.readByte(tempaddr) & 0x7f
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.d)
						break
					case 0xbb: /* LD E,RES 7,(REGISTER+dd) */
						z80.e = z80.memory.readByte(tempaddr) & 0x7f
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.e)
						break
					case 0xbc: /* LD H,RES 7,(REGISTER+dd) */
						z80.h = z80.memory.readByte(tempaddr) & 0x7f
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.h)
						break
					case 0xbd: /* LD L,RES 7,(REGISTER+dd) */
						z80.l = z80.memory.readByte(tempaddr) & 0x7f
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.l)
						break
					case 0xbe: /* RES 7,(REGISTER+dd) */

						var bytetemp byte
						bytetemp = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, bytetemp&0x7f)

						break
					case 0xbf: /* LD A,RES 7,(REGISTER+dd) */
						z80.a = z80.memory.readByte(tempaddr) & 0x7f
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.a)
						break
					case 0xc0: /* LD B,SET 0,(REGISTER+dd) */
						z80.b = z80.memory.readByte(tempaddr) | 0x01
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.b)
						break
					case 0xc1: /* LD C,SET 0,(REGISTER+dd) */
						z80.c = z80.memory.readByte(tempaddr) | 0x01
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.c)
						break
					case 0xc2: /* LD D,SET 0,(REGISTER+dd) */
						z80.d = z80.memory.readByte(tempaddr) | 0x01
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.d)
						break
					case 0xc3: /* LD E,SET 0,(REGISTER+dd) */
						z80.e = z80.memory.readByte(tempaddr) | 0x01
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.e)
						break
					case 0xc4: /* LD H,SET 0,(REGISTER+dd) */
						z80.h = z80.memory.readByte(tempaddr) | 0x01
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.h)
						break
					case 0xc5: /* LD L,SET 0,(REGISTER+dd) */
						z80.l = z80.memory.readByte(tempaddr) | 0x01
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.l)
						break
					case 0xc6: /* SET 0,(REGISTER+dd) */

						var bytetemp byte
						bytetemp = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, bytetemp|0x01)

						break
					case 0xc7: /* LD A,SET 0,(REGISTER+dd) */
						z80.a = z80.memory.readByte(tempaddr) | 0x01
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.a)
						break
					case 0xc8: /* LD B,SET 1,(REGISTER+dd) */
						z80.b = z80.memory.readByte(tempaddr) | 0x02
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.b)
						break
					case 0xc9: /* LD C,SET 1,(REGISTER+dd) */
						z80.c = z80.memory.readByte(tempaddr) | 0x02
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.c)
						break
					case 0xca: /* LD D,SET 1,(REGISTER+dd) */
						z80.d = z80.memory.readByte(tempaddr) | 0x02
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.d)
						break
					case 0xcb: /* LD E,SET 1,(REGISTER+dd) */
						z80.e = z80.memory.readByte(tempaddr) | 0x02
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.e)
						break
					case 0xcc: /* LD H,SET 1,(REGISTER+dd) */
						z80.h = z80.memory.readByte(tempaddr) | 0x02
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.h)
						break
					case 0xcd: /* LD L,SET 1,(REGISTER+dd) */
						z80.l = z80.memory.readByte(tempaddr) | 0x02
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.l)
						break
					case 0xce: /* SET 1,(REGISTER+dd) */

						var bytetemp byte
						bytetemp = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, bytetemp|0x02)

						break
					case 0xcf: /* LD A,SET 1,(REGISTER+dd) */
						z80.a = z80.memory.readByte(tempaddr) | 0x02
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.a)
						break
					case 0xd0: /* LD B,SET 2,(REGISTER+dd) */
						z80.b = z80.memory.readByte(tempaddr) | 0x04
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.b)
						break
					case 0xd1: /* LD C,SET 2,(REGISTER+dd) */
						z80.c = z80.memory.readByte(tempaddr) | 0x04
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.c)
						break
					case 0xd2: /* LD D,SET 2,(REGISTER+dd) */
						z80.d = z80.memory.readByte(tempaddr) | 0x04
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.d)
						break
					case 0xd3: /* LD E,SET 2,(REGISTER+dd) */
						z80.e = z80.memory.readByte(tempaddr) | 0x04
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.e)
						break
					case 0xd4: /* LD H,SET 2,(REGISTER+dd) */
						z80.h = z80.memory.readByte(tempaddr) | 0x04
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.h)
						break
					case 0xd5: /* LD L,SET 2,(REGISTER+dd) */
						z80.l = z80.memory.readByte(tempaddr) | 0x04
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.l)
						break
					case 0xd6: /* SET 2,(REGISTER+dd) */

						var bytetemp byte
						bytetemp = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, bytetemp|0x04)

						break
					case 0xd7: /* LD A,SET 2,(REGISTER+dd) */
						z80.a = z80.memory.readByte(tempaddr) | 0x04
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.a)
						break
					case 0xd8: /* LD B,SET 3,(REGISTER+dd) */
						z80.b = z80.memory.readByte(tempaddr) | 0x08
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.b)
						break
					case 0xd9: /* LD C,SET 3,(REGISTER+dd) */
						z80.c = z80.memory.readByte(tempaddr) | 0x08
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.c)
						break
					case 0xda: /* LD D,SET 3,(REGISTER+dd) */
						z80.d = z80.memory.readByte(tempaddr) | 0x08
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.d)
						break
					case 0xdb: /* LD E,SET 3,(REGISTER+dd) */
						z80.e = z80.memory.readByte(tempaddr) | 0x08
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.e)
						break
					case 0xdc: /* LD H,SET 3,(REGISTER+dd) */
						z80.h = z80.memory.readByte(tempaddr) | 0x08
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.h)
						break
					case 0xdd: /* LD L,SET 3,(REGISTER+dd) */
						z80.l = z80.memory.readByte(tempaddr) | 0x08
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.l)
						break
					case 0xde: /* SET 3,(REGISTER+dd) */

						var bytetemp byte
						bytetemp = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, bytetemp|0x08)

						break
					case 0xdf: /* LD A,SET 3,(REGISTER+dd) */
						z80.a = z80.memory.readByte(tempaddr) | 0x08
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.a)
						break
					case 0xe0: /* LD B,SET 4,(REGISTER+dd) */
						z80.b = z80.memory.readByte(tempaddr) | 0x10
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.b)
						break
					case 0xe1: /* LD C,SET 4,(REGISTER+dd) */
						z80.c = z80.memory.readByte(tempaddr) | 0x10
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.c)
						break
					case 0xe2: /* LD D,SET 4,(REGISTER+dd) */
						z80.d = z80.memory.readByte(tempaddr) | 0x10
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.d)
						break
					case 0xe3: /* LD E,SET 4,(REGISTER+dd) */
						z80.e = z80.memory.readByte(tempaddr) | 0x10
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.e)
						break
					case 0xe4: /* LD H,SET 4,(REGISTER+dd) */
						z80.h = z80.memory.readByte(tempaddr) | 0x10
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.h)
						break
					case 0xe5: /* LD L,SET 4,(REGISTER+dd) */
						z80.l = z80.memory.readByte(tempaddr) | 0x10
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.l)
						break
					case 0xe6: /* SET 4,(REGISTER+dd) */

						var bytetemp byte
						bytetemp = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, bytetemp|0x10)

						break
					case 0xe7: /* LD A,SET 4,(REGISTER+dd) */
						z80.a = z80.memory.readByte(tempaddr) | 0x10
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.a)
						break
					case 0xe8: /* LD B,SET 5,(REGISTER+dd) */
						z80.b = z80.memory.readByte(tempaddr) | 0x20
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.b)
						break
					case 0xe9: /* LD C,SET 5,(REGISTER+dd) */
						z80.c = z80.memory.readByte(tempaddr) | 0x20
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.c)
						break
					case 0xea: /* LD D,SET 5,(REGISTER+dd) */
						z80.d = z80.memory.readByte(tempaddr) | 0x20
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.d)
						break
					case 0xeb: /* LD E,SET 5,(REGISTER+dd) */
						z80.e = z80.memory.readByte(tempaddr) | 0x20
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.e)
						break
					case 0xec: /* LD H,SET 5,(REGISTER+dd) */
						z80.h = z80.memory.readByte(tempaddr) | 0x20
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.h)
						break
					case 0xed: /* LD L,SET 5,(REGISTER+dd) */
						z80.l = z80.memory.readByte(tempaddr) | 0x20
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.l)
						break
					case 0xee: /* SET 5,(REGISTER+dd) */

						var bytetemp byte
						bytetemp = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, bytetemp|0x20)

						break
					case 0xef: /* LD A,SET 5,(REGISTER+dd) */
						z80.a = z80.memory.readByte(tempaddr) | 0x20
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.a)
						break
					case 0xf0: /* LD B,SET 6,(REGISTER+dd) */
						z80.b = z80.memory.readByte(tempaddr) | 0x40
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.b)
						break
					case 0xf1: /* LD C,SET 6,(REGISTER+dd) */
						z80.c = z80.memory.readByte(tempaddr) | 0x40
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.c)
						break
					case 0xf2: /* LD D,SET 6,(REGISTER+dd) */
						z80.d = z80.memory.readByte(tempaddr) | 0x40
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.d)
						break
					case 0xf3: /* LD E,SET 6,(REGISTER+dd) */
						z80.e = z80.memory.readByte(tempaddr) | 0x40
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.e)
						break
					case 0xf4: /* LD H,SET 6,(REGISTER+dd) */
						z80.h = z80.memory.readByte(tempaddr) | 0x40
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.h)
						break
					case 0xf5: /* LD L,SET 6,(REGISTER+dd) */
						z80.l = z80.memory.readByte(tempaddr) | 0x40
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.l)
						break
					case 0xf6: /* SET 6,(REGISTER+dd) */

						var bytetemp byte
						bytetemp = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, bytetemp|0x40)

						break
					case 0xf7: /* LD A,SET 6,(REGISTER+dd) */
						z80.a = z80.memory.readByte(tempaddr) | 0x40
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.a)
						break
					case 0xf8: /* LD B,SET 7,(REGISTER+dd) */
						z80.b = z80.memory.readByte(tempaddr) | 0x80
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.b)
						break
					case 0xf9: /* LD C,SET 7,(REGISTER+dd) */
						z80.c = z80.memory.readByte(tempaddr) | 0x80
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.c)
						break
					case 0xfa: /* LD D,SET 7,(REGISTER+dd) */
						z80.d = z80.memory.readByte(tempaddr) | 0x80
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.d)
						break
					case 0xfb: /* LD E,SET 7,(REGISTER+dd) */
						z80.e = z80.memory.readByte(tempaddr) | 0x80
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.e)
						break
					case 0xfc: /* LD H,SET 7,(REGISTER+dd) */
						z80.h = z80.memory.readByte(tempaddr) | 0x80
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.h)
						break
					case 0xfd: /* LD L,SET 7,(REGISTER+dd) */
						z80.l = z80.memory.readByte(tempaddr) | 0x80
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.l)
						break
					case 0xfe: /* SET 7,(REGISTER+dd) */

						var bytetemp byte
						bytetemp = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, bytetemp|0x80)

						break
					case 0xff: /* LD A,SET 7,(REGISTER+dd) */
						z80.a = z80.memory.readByte(tempaddr) | 0x80
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.a)
						break

					}

					break
				case 0xe1: /* POP ix */
					z80.pop16(&z80.ixl, &z80.ixh)
					break
				case 0xe3: /* EX (SP),ix */
					var bytetempl, bytetemph byte
					bytetempl = z80.memory.readByte(z80.SP())
					bytetemph = z80.memory.readByte(z80.SP() + 1)
					z80.memory.contendReadNoMreq(z80.SP()+1, 1)
					z80.memory.writeByte(z80.SP()+1, z80.ixh)
					z80.memory.writeByte(z80.SP(), z80.ixl)
					z80.memory.contendWriteNoMreq(z80.SP(), 1)
					z80.memory.contendWriteNoMreq(z80.SP(), 1)
					z80.ixl = bytetempl
					z80.ixh = bytetemph
					break
				case 0xe5: /* PUSH ix */
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.push16(z80.ixl, z80.ixh)
					break
				case 0xe9: /* JP ix */
					z80.pc = z80.IX() /* NB: NOT INDIRECT! */
					break
				case 0xf9: /* LD SP,ix */
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.sp = z80.IX()
					break
				default: /* Instruction did not involve H or L, so backtrack
					   one instruction and parse again */
					z80.pc--
					z80.r--
					opcode = opcode2

					goto EndOpcode

				}
			}
			break
		case 0xde: /* SBC A,nn */
			{
				var bytetemp byte = z80.memory.readByte(z80.PC())
				z80.pc++
				z80.sbc(bytetemp)
			}
			break
		case 0xdf: /* RST 18 */
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.rst(0x18)
			break
		case 0xe0: /* RET PO */
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			if !((z80.f & FLAG_P) != 0) {
				z80.ret()
			}
			break
		case 0xe1: /* POP HL */
			z80.pop16(&z80.l, &z80.h)
			break
		case 0xe2: /* JP PO,nnnn */
			if (z80.f & FLAG_P) == 0 {
				z80.jp()
			} else {
				z80.memory.contendRead(z80.pc, 3)
				z80.memory.contendRead(z80.pc+1, 3)
				z80.pc += 2
			}
			break
		case 0xe3: /* EX (SP),HL */
			var bytetempl, bytetemph byte
			bytetempl = z80.memory.readByte(z80.SP())
			bytetemph = z80.memory.readByte(z80.SP() + 1)
			z80.memory.contendReadNoMreq(z80.SP()+1, 1)
			z80.memory.writeByte(z80.SP()+1, z80.h)
			z80.memory.writeByte(z80.SP(), z80.l)
			z80.memory.contendWriteNoMreq(z80.SP(), 1)
			z80.memory.contendWriteNoMreq(z80.SP(), 1)
			z80.l = bytetempl
			z80.h = bytetemph
			break
		case 0xe4: /* CALL PO,nnnn */
			if (z80.f & FLAG_P) == 0 {
				z80.call()
			} else {
				z80.memory.contendRead(z80.pc, 3)
				z80.memory.contendRead(z80.pc+1, 3)
				z80.pc += 2
			}
			break
		case 0xe5: /* PUSH HL */
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.push16(z80.l, z80.h)
			break
		case 0xe6: /* AND nn */
			{
				var bytetemp byte = z80.memory.readByte(z80.PC())
				z80.pc++
				z80.and(bytetemp)
			}
			break
		case 0xe7: /* RST 20 */
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.rst(0x20)
			break
		case 0xe8: /* RET PE */
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			if (z80.f & FLAG_P) != 0 {
				z80.ret()
			}
			break
		case 0xe9: /* JP HL */
			z80.pc = z80.HL() /* NB: NOT INDIRECT! */
			break
		case 0xea: /* JP PE,nnnn */
			if (z80.f & FLAG_P) != 0 {
				z80.jp()
			} else {
				z80.memory.contendRead(z80.pc, 3)
				z80.memory.contendRead(z80.pc+1, 3)
				z80.pc += 2
			}
			break
		case 0xeb: /* EX DE,HL */
			var wordtemp uint16 = z80.DE()
			z80.setDE(z80.HL())
			z80.setHL(wordtemp)
			break
		case 0xec: /* CALL PE,nnnn */
			if (z80.f & FLAG_P) != 0 {
				z80.call()
			} else {
				z80.memory.contendRead(z80.pc, 3)
				z80.memory.contendRead(z80.pc+1, 3)
				z80.pc += 2
			}
			break
		case 0xed: /* shift ED */
			{
				var opcode2 byte
				z80.memory.contendRead(z80.pc, 4)
				opcode2 = z80.memory.readByteInternal(z80.pc)
				z80.pc++
				z80.r++

				switch opcode2 {
				/* z80_ed.c: Z80 CBxx opcodes
				   Copyright (c) 1999-2003 Philip Kendall

				   This program is free software; you can redistribute it and/or modify
				   it under the terms of the GNU General Public License as published by
				   the Free Software Foundation; either version 2 of the License, or
				   (at your option) any later version.

				   This program is distributed in the hope that it will be useful,
				   but WITHOUT ANY WARRANTY; without even the implied warranty of
				   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
				   GNU General Public License for more details.

				   You should have received a copy of the GNU General Public License along
				   with this program; if not, write to the Free Software Foundation, Inc.,
				   51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.

				   Author contact information:

				   E-mail: philip-fuse@shadowmagic.org.uk

				*/

				/* NB: this file is autogenerated by './z80.pl' from 'opcodes_ed.dat',
				   and included in 'z80_ops.c' */

				case 0x40: /* IN B,(C) */
					z80.in(&z80.b, z80.BC())
					break
				case 0x41: /* OUT (C),B */
					z80.writePort(z80.BC(), z80.b)
					break
				case 0x42: /* SBC HL,BC */
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.sbc16(z80.BC())
					break
				case 0x43: /* LD (nnnn),BC */
					z80.ld16nnrr(z80.c, z80.b)
					break
					break
				case 0x44:
					fallthrough
				case 0x4c:
					fallthrough
				case 0x54:
					fallthrough
				case 0x5c:
					fallthrough
				case 0x64:
					fallthrough
				case 0x6c:
					fallthrough
				case 0x74:
					fallthrough
				case 0x7c: /* NEG */
					bytetemp := z80.a
					z80.a = 0
					z80.sub(bytetemp)
					break
				case 0x45:
					fallthrough
				case 0x4d:
					fallthrough
				case 0x55:
					fallthrough
				case 0x5d:
					fallthrough
				case 0x65:
					fallthrough
				case 0x6d:
					fallthrough
				case 0x75:
					fallthrough
				case 0x7d: /* RETN */
					z80.iff1 = z80.iff2
					z80.ret()
					break
				case 0x46:
					fallthrough
				case 0x4e:
					fallthrough
				case 0x66:
					fallthrough
				case 0x6e: /* IM 0 */
					z80.im = 0
					break
				case 0x47: /* LD I,A */
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.i = z80.a
					break
				case 0x48: /* IN C,(C) */
					z80.in(&z80.c, z80.BC())
					break
				case 0x49: /* OUT (C),C */
					z80.writePort(z80.BC(), z80.c)
					break
				case 0x4a: /* ADC HL,BC */
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.adc16(z80.BC())
					break
				case 0x4b: /* LD BC,(nnnn) */
					z80.ld16rrnn(&z80.c, &z80.b)
					break
					break
				case 0x4f: /* LD R,A */
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					/* Keep the RZX instruction counter right */
					rzxInstructionsOffset += (int(z80.r) - int(z80.a))
					z80.r, z80.r7 = uint16(z80.a), uint16(z80.a)
					break
				case 0x50: /* IN D,(C) */
					z80.in(&z80.d, z80.BC())
					break
				case 0x51: /* OUT (C),D */
					z80.writePort(z80.BC(), z80.d)
					break
				case 0x52: /* SBC HL,DE */
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.sbc16(z80.DE())
					break
				case 0x53: /* LD (nnnn),DE */
					z80.ld16nnrr(z80.e, z80.d)
					break
					break
				case 0x56:
					fallthrough
				case 0x76: /* IM 1 */
					z80.im = 1
					break
				case 0x57: /* LD A,I */
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.a = z80.i
					z80.f = (z80.f & FLAG_C) | z80.sz53Table[z80.a] | ternOpB(z80.iff2 != 0, FLAG_V, 0)
					break
				case 0x58: /* IN E,(C) */
					z80.in(&z80.e, z80.BC())
					break
				case 0x59: /* OUT (C),E */
					z80.writePort(z80.BC(), z80.e)
					break
				case 0x5a: /* ADC HL,DE */
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.adc16(z80.DE())
					break
				case 0x5b: /* LD DE,(nnnn) */
					z80.ld16rrnn(&z80.e, &z80.d)
					break
					break
				case 0x5e:
					fallthrough
				case 0x7e: /* IM 2 */
					z80.im = 2
					break
				case 0x5f: /* LD A,R */
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.a = byte((z80.r & 0x7f) | (z80.r7 & 0x80))
					z80.f = (z80.f & FLAG_C) | z80.sz53Table[z80.a] | ternOpB(z80.iff2 != 0, FLAG_V, 0)
					break
				case 0x60: /* IN H,(C) */
					z80.in(&z80.h, z80.BC())
					break
				case 0x61: /* OUT (C),H */
					z80.writePort(z80.BC(), z80.h)
					break
				case 0x62: /* SBC HL,HL */
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.sbc16(z80.HL())
					break
				case 0x63: /* LD (nnnn),HL */
					z80.ld16nnrr(z80.l, z80.h)
					break
					break
				case 0x67: /* RRD */
					var bytetemp byte = z80.memory.readByte(z80.HL())
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.memory.writeByte(z80.HL(), (z80.a<<4)|(bytetemp>>4))
					z80.a = (z80.a & 0xf0) | (bytetemp & 0x0f)
					z80.f = (z80.f & FLAG_C) | z80.sz53pTable[z80.a]
					break
				case 0x68: /* IN L,(C) */
					z80.in(&z80.l, z80.BC())
					break
				case 0x69: /* OUT (C),L */
					z80.writePort(z80.BC(), z80.l)
					break
				case 0x6a: /* ADC HL,HL */
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.adc16(z80.HL())
					break
				case 0x6b: /* LD HL,(nnnn) */
					z80.ld16rrnn(&z80.l, &z80.h)
					break
					break
				case 0x6f: /* RLD */
					var bytetemp byte = z80.memory.readByte(z80.HL())
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.memory.writeByte(z80.HL(), (bytetemp<<4)|(z80.a&0x0f))
					z80.a = (z80.a & 0xf0) | (bytetemp >> 4)
					z80.f = (z80.f & FLAG_C) | z80.sz53pTable[z80.a]
					break
				case 0x70: /* IN F,(C) */
					var bytetemp byte
					z80.in(&bytetemp, z80.BC())
					break
				case 0x71: /* OUT (C),0 */
					z80.writePort(z80.BC(), 0)
					break
				case 0x72: /* SBC HL,SP */
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.sbc16(z80.SP())
					break
				case 0x73: /* LD (nnnn),SP */
					sph, spl := splitWord(z80.sp)
					z80.ld16nnrr(spl, sph)
					break
					break
				case 0x78: /* IN A,(C) */
					z80.in(&z80.a, z80.BC())
					break
				case 0x79: /* OUT (C),A */
					z80.writePort(z80.BC(), z80.a)
					break
				case 0x7a: /* ADC HL,SP */
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.adc16(z80.SP())
					break
				case 0x7b: /* LD SP,(nnnn) */
					sph, spl := splitWord(z80.sp)
					z80.ld16rrnn(&spl, &sph)
					z80.sp = joinBytes(sph, spl)
					break
					break
				case 0xa0: /* LDI */
					var bytetemp byte = z80.memory.readByte(z80.HL())
					z80.decBC()
					z80.memory.writeByte(z80.DE(), bytetemp)
					z80.memory.contendWriteNoMreq(z80.DE(), 1)
					z80.memory.contendWriteNoMreq(z80.DE(), 1)
					z80.incDE()
					z80.incHL()
					bytetemp += z80.a
					z80.f = (z80.f & (FLAG_C | FLAG_Z | FLAG_S)) | ternOpB(z80.BC() != 0, FLAG_V, 0) |
						(bytetemp & FLAG_3) | ternOpB((bytetemp&0x02) != 0, FLAG_5, 0)
					break
				case 0xa1: /* CPI */

					var value, bytetemp, lookup byte

					value = z80.memory.readByte(z80.HL())
					bytetemp = z80.a - value
					lookup = ((z80.a & 0x08) >> 3) | (((value) & 0x08) >> 2) | ((bytetemp & 0x08) >> 1)

					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.incHL()
					z80.decBC()
					z80.f = (z80.f & FLAG_C) | ternOpB(z80.BC() != 0, FLAG_V|FLAG_N, FLAG_N) | halfcarrySubTable[lookup] | ternOpB(bytetemp != 0, 0, FLAG_Z) | (bytetemp & FLAG_S)
					if (z80.f & FLAG_H) != 0 {
						bytetemp--
					}
					z80.f |= (bytetemp & FLAG_3) | ternOpB((bytetemp&0x02) != 0, FLAG_5, 0)

					break
				case 0xa2: /* INI */
					var initemp, initemp2 byte

					z80.memory.contendReadNoMreq(z80.IR(), 1)
					initemp = z80.readPort(z80.BC())
					z80.memory.writeByte(z80.HL(), initemp)

					z80.b--
					z80.incHL()
					initemp2 = initemp + z80.c + 1
					z80.f = ternOpB((initemp&0x80) != 0, FLAG_N, 0) | ternOpB(initemp2 < initemp, FLAG_H|FLAG_C, 0) | ternOpB(z80.parityTable[(initemp2&0x07)^z80.b] != 0, FLAG_P, 0) | z80.sz53Table[z80.b]
					break
				case 0xa3: /* OUTI */
					var outitemp, outitemp2 byte

					z80.memory.contendReadNoMreq(z80.IR(), 1)
					outitemp = z80.memory.readByte(z80.HL())
					z80.b-- /* This does happen first, despite what the specs say */
					z80.writePort(z80.BC(), outitemp)

					z80.incHL()
					outitemp2 = outitemp + z80.l
					z80.f = ternOpB((outitemp&0x80) != 0, FLAG_N, 0) |
						ternOpB(outitemp2 < outitemp, FLAG_H|FLAG_C, 0) |
						ternOpB(z80.parityTable[(outitemp2&0x07)^z80.b] != 0, FLAG_P, 0) |
						z80.sz53Table[z80.b]
					break
				case 0xa8: /* LDD */
					var bytetemp byte = z80.memory.readByte(z80.HL())
					z80.decBC()
					z80.memory.writeByte(z80.DE(), bytetemp)
					z80.memory.contendWriteNoMreq(z80.DE(), 1)
					z80.memory.contendWriteNoMreq(z80.DE(), 1)
					z80.decDE()
					z80.decHL()
					bytetemp += z80.a
					z80.f = (z80.f & (FLAG_C | FLAG_Z | FLAG_S)) | ternOpB(z80.BC() != 0, FLAG_V, 0) |
						(bytetemp & FLAG_3) | ternOpB((bytetemp&0x02) != 0, FLAG_5, 0)
					break
				case 0xa9: /* CPD */

					var value, bytetemp, lookup byte

					value = z80.memory.readByte(z80.HL())
					bytetemp = z80.a - value
					lookup = ((z80.a & 0x08) >> 3) | (((value) & 0x08) >> 2) | ((bytetemp & 0x08) >> 1)

					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.decHL()
					z80.decBC()
					z80.f = (z80.f & FLAG_C) | ternOpB(z80.BC() != 0, FLAG_V|FLAG_N, FLAG_N) | halfcarrySubTable[lookup] | ternOpB(bytetemp != 0, 0, FLAG_Z) | (bytetemp & FLAG_S)
					if (z80.f & FLAG_H) != 0 {
						bytetemp--
					}
					z80.f |= (bytetemp & FLAG_3) | ternOpB((bytetemp&0x02) != 0, FLAG_5, 0)

					break
				case 0xaa: /* IND */
					var initemp, initemp2 byte

					z80.memory.contendReadNoMreq(z80.IR(), 1)
					initemp = z80.readPort(z80.BC())
					z80.memory.writeByte(z80.HL(), initemp)

					z80.b--
					z80.decHL()
					initemp2 = initemp + z80.c - 1
					z80.f = ternOpB((initemp&0x80) != 0, FLAG_N, 0) | ternOpB(initemp2 < initemp, FLAG_H|FLAG_C, 0) | ternOpB(z80.parityTable[(initemp2&0x07)^z80.b] != 0, FLAG_P, 0) | z80.sz53Table[z80.b]
					break
				case 0xab: /* OUTD */
					var outitemp, outitemp2 byte

					z80.memory.contendReadNoMreq(z80.IR(), 1)
					outitemp = z80.memory.readByte(z80.HL())
					z80.b-- /* This does happen first, despite what the specs say */
					z80.writePort(z80.BC(), outitemp)

					z80.decHL()
					outitemp2 = outitemp + z80.l
					z80.f = ternOpB((outitemp&0x80) != 0, FLAG_N, 0) |
						ternOpB(outitemp2 < outitemp, FLAG_H|FLAG_C, 0) |
						ternOpB(z80.parityTable[(outitemp2&0x07)^z80.b] != 0, FLAG_P, 0) |
						z80.sz53Table[z80.b]
					break
				case 0xb0: /* LDIR */
					var bytetemp byte = z80.memory.readByte(z80.HL())
					z80.memory.writeByte(z80.DE(), bytetemp)
					z80.memory.contendWriteNoMreq(z80.DE(), 1)
					z80.memory.contendWriteNoMreq(z80.DE(), 1)
					z80.decBC()
					bytetemp += z80.a
					z80.f = (z80.f & (FLAG_C | FLAG_Z | FLAG_S)) | ternOpB(z80.BC() != 0, FLAG_V, 0) | (bytetemp & FLAG_3) | ternOpB((bytetemp&0x02 != 0), FLAG_5, 0)
					if z80.BC() != 0 {
						z80.memory.contendWriteNoMreq(z80.DE(), 1)
						z80.memory.contendWriteNoMreq(z80.DE(), 1)
						z80.memory.contendWriteNoMreq(z80.DE(), 1)
						z80.memory.contendWriteNoMreq(z80.DE(), 1)
						z80.memory.contendWriteNoMreq(z80.DE(), 1)
						z80.pc -= 2
					}
					z80.incHL()
					z80.incDE()
					break
				case 0xb1: /* CPIR */
					var value, bytetemp, lookup byte

					value = z80.memory.readByte(z80.HL())
					bytetemp = z80.a - value
					lookup = ((z80.a & 0x08) >> 3) | (((value) & 0x08) >> 2) | ((bytetemp & 0x08) >> 1)

					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.decBC()
					z80.f = (z80.f & FLAG_C) | (ternOpB(z80.BC() != 0, (FLAG_V | FLAG_N), FLAG_N)) | halfcarrySubTable[lookup] | (ternOpB(bytetemp != 0, 0, FLAG_Z)) | (bytetemp & FLAG_S)
					if (z80.f & FLAG_H) != 0 {
						bytetemp--
					}
					z80.f |= (bytetemp & FLAG_3) | ternOpB((bytetemp&0x02) != 0, FLAG_5, 0)
					if (z80.f & (FLAG_V | FLAG_Z)) == FLAG_V {
						z80.memory.contendReadNoMreq(z80.HL(), 1)
						z80.memory.contendReadNoMreq(z80.HL(), 1)
						z80.memory.contendReadNoMreq(z80.HL(), 1)
						z80.memory.contendReadNoMreq(z80.HL(), 1)
						z80.memory.contendReadNoMreq(z80.HL(), 1)
						z80.pc -= 2
					}
					z80.incHL()
					break
				case 0xb2: /* INIR */
					var initemp, initemp2 byte

					z80.memory.contendReadNoMreq(z80.IR(), 1)
					initemp = z80.readPort(z80.BC())
					z80.memory.writeByte(z80.HL(), initemp)

					z80.b--
					initemp2 = initemp + z80.c + 1
					z80.f = ternOpB(initemp&0x80 != 0, FLAG_N, 0) |
						ternOpB(initemp2 < initemp, FLAG_H|FLAG_C, 0) |
						ternOpB(z80.parityTable[(initemp2&0x07)^z80.b] != 0, FLAG_P, 0) |
						z80.sz53Table[z80.b]

					if z80.b != 0 {
						z80.memory.contendWriteNoMreq(z80.HL(), 1)
						z80.memory.contendWriteNoMreq(z80.HL(), 1)
						z80.memory.contendWriteNoMreq(z80.HL(), 1)
						z80.memory.contendWriteNoMreq(z80.HL(), 1)
						z80.memory.contendWriteNoMreq(z80.HL(), 1)
						z80.pc -= 2
					}
					z80.incHL()
					break
				case 0xb3: /* OTIR */
					var outitemp, outitemp2 byte

					z80.memory.contendReadNoMreq(z80.IR(), 1)
					outitemp = z80.memory.readByte(z80.HL())
					z80.b-- /* This does happen first, despite what the specs say */
					z80.writePort(z80.BC(), outitemp)

					z80.incHL()
					outitemp2 = outitemp + z80.l
					z80.f = ternOpB((outitemp&0x80) != 0, FLAG_N, 0) |
						ternOpB(outitemp2 < outitemp, FLAG_H|FLAG_C, 0) |
						ternOpB(z80.parityTable[(outitemp2&0x07)^z80.b] != 0, FLAG_P, 0) |
						z80.sz53Table[z80.b]

					if z80.b != 0 {
						z80.memory.contendReadNoMreq(z80.BC(), 1)
						z80.memory.contendReadNoMreq(z80.BC(), 1)
						z80.memory.contendReadNoMreq(z80.BC(), 1)
						z80.memory.contendReadNoMreq(z80.BC(), 1)
						z80.memory.contendReadNoMreq(z80.BC(), 1)
						z80.pc -= 2
					}
					break
				case 0xb8: /* LDDR */
					var bytetemp byte = z80.memory.readByte(z80.HL())
					z80.memory.writeByte(z80.DE(), bytetemp)
					z80.memory.contendWriteNoMreq(z80.DE(), 1)
					z80.memory.contendWriteNoMreq(z80.DE(), 1)
					z80.decBC()
					bytetemp += z80.a
					z80.f = (z80.f & (FLAG_C | FLAG_Z | FLAG_S)) | ternOpB(z80.BC() != 0, FLAG_V, 0) | (bytetemp & FLAG_3) | ternOpB((bytetemp&0x02 != 0), FLAG_5, 0)
					if z80.BC() != 0 {
						z80.memory.contendWriteNoMreq(z80.DE(), 1)
						z80.memory.contendWriteNoMreq(z80.DE(), 1)
						z80.memory.contendWriteNoMreq(z80.DE(), 1)
						z80.memory.contendWriteNoMreq(z80.DE(), 1)
						z80.memory.contendWriteNoMreq(z80.DE(), 1)
						z80.pc -= 2
					}
					z80.decHL()
					z80.decDE()
					break
				case 0xb9: /* CPDR */
					var value, bytetemp, lookup byte

					value = z80.memory.readByte(z80.HL())
					bytetemp = z80.a - value
					lookup = ((z80.a & 0x08) >> 3) | (((value) & 0x08) >> 2) | ((bytetemp & 0x08) >> 1)

					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.memory.contendReadNoMreq(z80.HL(), 1)
					z80.decBC()
					z80.f = (z80.f & FLAG_C) | (ternOpB(z80.BC() != 0, (FLAG_V | FLAG_N), FLAG_N)) | halfcarrySubTable[lookup] | (ternOpB(bytetemp != 0, 0, FLAG_Z)) | (bytetemp & FLAG_S)
					if (z80.f & FLAG_H) != 0 {
						bytetemp--
					}
					z80.f |= (bytetemp & FLAG_3) | ternOpB((bytetemp&0x02) != 0, FLAG_5, 0)
					if (z80.f & (FLAG_V | FLAG_Z)) == FLAG_V {
						z80.memory.contendReadNoMreq(z80.HL(), 1)
						z80.memory.contendReadNoMreq(z80.HL(), 1)
						z80.memory.contendReadNoMreq(z80.HL(), 1)
						z80.memory.contendReadNoMreq(z80.HL(), 1)
						z80.memory.contendReadNoMreq(z80.HL(), 1)
						z80.pc -= 2
					}
					z80.decHL()
					break
				case 0xba: /* INDR */
					var initemp, initemp2 byte

					z80.memory.contendReadNoMreq(z80.IR(), 1)
					initemp = z80.readPort(z80.BC())
					z80.memory.writeByte(z80.HL(), initemp)

					z80.b--
					initemp2 = initemp + z80.c - 1
					z80.f = ternOpB(initemp&0x80 != 0, FLAG_N, 0) |
						ternOpB(initemp2 < initemp, FLAG_H|FLAG_C, 0) |
						ternOpB(z80.parityTable[(initemp2&0x07)^z80.b] != 0, FLAG_P, 0) |
						z80.sz53Table[z80.b]

					if z80.b != 0 {
						z80.memory.contendWriteNoMreq(z80.HL(), 1)
						z80.memory.contendWriteNoMreq(z80.HL(), 1)
						z80.memory.contendWriteNoMreq(z80.HL(), 1)
						z80.memory.contendWriteNoMreq(z80.HL(), 1)
						z80.memory.contendWriteNoMreq(z80.HL(), 1)
						z80.pc -= 2
					}
					z80.decHL()
					break
				case 0xbb: /* OTDR */
					var outitemp, outitemp2 byte

					z80.memory.contendReadNoMreq(z80.IR(), 1)
					outitemp = z80.memory.readByte(z80.HL())
					z80.b-- /* This does happen first, despite what the specs say */
					z80.writePort(z80.BC(), outitemp)

					z80.decHL()
					outitemp2 = outitemp + z80.l
					z80.f = ternOpB((outitemp&0x80) != 0, FLAG_N, 0) |
						ternOpB(outitemp2 < outitemp, FLAG_H|FLAG_C, 0) |
						ternOpB(z80.parityTable[(outitemp2&0x07)^z80.b] != 0, FLAG_P, 0) |
						z80.sz53Table[z80.b]

					if z80.b != 0 {
						z80.memory.contendReadNoMreq(z80.BC(), 1)
						z80.memory.contendReadNoMreq(z80.BC(), 1)
						z80.memory.contendReadNoMreq(z80.BC(), 1)
						z80.memory.contendReadNoMreq(z80.BC(), 1)
						z80.memory.contendReadNoMreq(z80.BC(), 1)
						z80.pc -= 2
					}
					break
				case 0xfb: /* slttrap */
					z80.sltTrap(int16(z80.HL()), z80.a)
					break
				default: /* All other opcodes are NOPD */
					break

				}
			}
			break
		case 0xee: /* XOR A,nn */
			{
				var bytetemp byte = z80.memory.readByte(z80.PC())
				z80.pc++
				z80.xor(bytetemp)
			}
			break
		case 0xef: /* RST 28 */
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.rst(0x28)
			break
		case 0xf0: /* RET P */
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			if !((z80.f & FLAG_S) != 0) {
				z80.ret()
			}
			break
		case 0xf1: /* POP AF */
			z80.pop16(&z80.f, &z80.a)
			break
		case 0xf2: /* JP P,nnnn */
			if (z80.f & FLAG_S) == 0 {
				z80.jp()
			} else {
				z80.memory.contendRead(z80.pc, 3)
				z80.memory.contendRead(z80.pc+1, 3)
				z80.pc += 2
			}
			break
		case 0xf3: /* DI */
			z80.iff1, z80.iff2 = 0, 0
			break
		case 0xf4: /* CALL P,nnnn */
			if (z80.f & FLAG_S) == 0 {
				z80.call()
			} else {
				z80.memory.contendRead(z80.pc, 3)
				z80.memory.contendRead(z80.pc+1, 3)
				z80.pc += 2
			}
			break
		case 0xf5: /* PUSH AF */
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.push16(z80.f, z80.a)
			break
		case 0xf6: /* OR nn */
			{
				var bytetemp byte = z80.memory.readByte(z80.PC())
				z80.pc++
				z80.or(bytetemp)
			}
			break
		case 0xf7: /* RST 30 */
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.rst(0x30)
			break
		case 0xf8: /* RET M */
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			if (z80.f & FLAG_S) != 0 {
				z80.ret()
			}
			break
		case 0xf9: /* LD SP,HL */
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.sp = z80.HL()
			break
		case 0xfa: /* JP M,nnnn */
			if (z80.f & FLAG_S) != 0 {
				z80.jp()
			} else {
				z80.memory.contendRead(z80.pc, 3)
				z80.memory.contendRead(z80.pc+1, 3)
				z80.pc += 2
			}
			break
		case 0xfb: /* EI */
			/* Interrupts are not accepted immediately after an EI, but are
			accepted after the next instruction */
			z80.iff1, z80.iff2 = 1, 1
			z80.interruptsEnabledAt = int(z80.tstates)
			// eventAdd(z80.tstates + 1, z80InterruptEvent)
			break
		case 0xfc: /* CALL M,nnnn */
			if (z80.f & FLAG_S) != 0 {
				z80.call()
			} else {
				z80.memory.contendRead(z80.pc, 3)
				z80.memory.contendRead(z80.pc+1, 3)
				z80.pc += 2
			}
			break
		case 0xfd: /* shift FD */
			{
				var opcode2 byte
				z80.memory.contendRead(z80.pc, 4)
				opcode2 = z80.memory.readByteInternal(z80.pc)
				z80.pc++
				z80.r++

				switch opcode2 {

				/* z80_ddfd.c Z80 {DD,FD}xx opcodes
				   Copyright (c) 1999-2003 Philip Kendall

				   This program is free software; you can redistribute it and/or modify
				   it under the terms of the GNU General Public License as published by
				   the Free Software Foundation; either version 2 of the License, or
				   (at your option) any later version.

				   This program is distributed in the hope that it will be useful,
				   but WITHOUT ANY WARRANTY; without even the implied warranty of
				   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
				   GNU General Public License for more details.

				   You should have received a copy of the GNU General Public License along
				   with this program; if not, write to the Free Software Foundation, Inc.,
				   51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.

				   Author contact information:

				   E-mail: philip-fuse@shadowmagic.org.uk

				*/

				/* NB: this file is autogenerated by './z80.pl' from 'opcodes_ddfd.dat',
				   and included in 'z80_ops.c' */

				case 0x09: /* ADD iy,BC */
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.add16(z80.iy, z80.BC())
					break
				case 0x19: /* ADD iy,DE */
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.add16(z80.iy, z80.DE())
					break
				case 0x21: /* LD iy,nnnn */
					b1 := z80.memory.readByte(z80.pc)
					z80.pc++
					b2 := z80.memory.readByte(z80.pc)
					z80.pc++
					z80.setIY(joinBytes(b2, b1))
					break
				case 0x22: /* LD (nnnn),iy */
					z80.ld16nnrr(z80.iyl, z80.iyh)
					break
					break
				case 0x23: /* INC iy */
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.incIY()
					break
				case 0x24: /* INC z80.IYL() */
					z80.incIYH()
					break
				case 0x25: /* DEC z80.IYL() */
					z80.decIYH()
					break
				case 0x26: /* LD z80.IYL(),nn */
					z80.iyh = z80.memory.readByte(z80.pc)
					z80.pc++
					break
				case 0x29: /* ADD iy,iy */
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.add16(z80.iy, z80.IY())
					break
				case 0x2a: /* LD iy,(nnnn) */
					z80.ld16rrnn(&z80.iyl, &z80.iyh)
					break
					break
				case 0x2b: /* DEC iy */
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.decIY()
					break
				case 0x2c: /* INC z80.IYH() */
					z80.incIYL()
					break
				case 0x2d: /* DEC z80.IYH() */
					z80.decIYL()
					break
				case 0x2e: /* LD z80.IYH(),nn */
					z80.iyl = z80.memory.readByte(z80.pc)
					z80.pc++
					break
				case 0x34: /* INC (iy+dd) */
					var offset, bytetemp byte
					var wordtemp uint16
					offset = z80.memory.readByte(z80.pc)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.pc++
					wordtemp = uint16(int(z80.IY()) + int(signExtend(offset)))
					bytetemp = z80.memory.readByte(wordtemp)
					z80.memory.contendReadNoMreq(wordtemp, 1)
					z80.inc(&bytetemp)
					z80.memory.writeByte(wordtemp, bytetemp)
					break
				case 0x35: /* DEC (iy+dd) */
					var offset, bytetemp byte
					var wordtemp uint16
					offset = z80.memory.readByte(z80.pc)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.pc++
					wordtemp = uint16(int(z80.IY()) + int(signExtend(offset)))
					bytetemp = z80.memory.readByte(wordtemp)
					z80.memory.contendReadNoMreq(wordtemp, 1)
					z80.dec(&bytetemp)
					z80.memory.writeByte(wordtemp, bytetemp)
					break
				case 0x36: /* LD (iy+dd),nn */
					offset := z80.memory.readByte(z80.pc)
					z80.pc++
					value := z80.memory.readByte(z80.pc)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.pc++
					z80.memory.writeByte(uint16(int(z80.IY())+int(signExtend(offset))), value)
					break
				case 0x39: /* ADD iy,SP */
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.add16(z80.iy, z80.SP())
					break
				case 0x44: /* LD B,z80.IYL() */
					z80.b = z80.iyh
					break
				case 0x45: /* LD B,z80.IYH() */
					z80.b = z80.iyl
					break
				case 0x46: /* LD B,(iy+dd) */
					var offset byte
					offset = z80.memory.readByte(z80.pc)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.pc++
					z80.b = z80.memory.readByte(uint16(int(z80.IY()) + int(signExtend(offset))))
					break
				case 0x4c: /* LD C,z80.IYL() */
					z80.c = z80.iyh
					break
				case 0x4d: /* LD C,z80.IYH() */
					z80.c = z80.iyl
					break
				case 0x4e: /* LD C,(iy+dd) */
					var offset byte
					offset = z80.memory.readByte(z80.pc)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.pc++
					z80.c = z80.memory.readByte(uint16(int(z80.IY()) + int(signExtend(offset))))
					break
				case 0x54: /* LD D,z80.IYL() */
					z80.d = z80.iyh
					break
				case 0x55: /* LD D,z80.IYH() */
					z80.d = z80.iyl
					break
				case 0x56: /* LD D,(iy+dd) */
					var offset byte
					offset = z80.memory.readByte(z80.pc)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.pc++
					z80.d = z80.memory.readByte(uint16(int(z80.IY()) + int(signExtend(offset))))
					break
				case 0x5c: /* LD E,z80.IYL() */
					z80.e = z80.iyh
					break
				case 0x5d: /* LD E,z80.IYH() */
					z80.e = z80.iyl
					break
				case 0x5e: /* LD E,(iy+dd) */
					var offset byte
					offset = z80.memory.readByte(z80.pc)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.pc++
					z80.e = z80.memory.readByte(uint16(int(z80.IY()) + int(signExtend(offset))))
					break
				case 0x60: /* LD z80.IYL(),B */
					z80.iyh = z80.b
					break
				case 0x61: /* LD z80.IYL(),C */
					z80.iyh = z80.c
					break
				case 0x62: /* LD z80.IYL(),D */
					z80.iyh = z80.d
					break
				case 0x63: /* LD z80.IYL(),E */
					z80.iyh = z80.e
					break
				case 0x64: /* LD z80.IYL(),z80.IYL() */
					break
				case 0x65: /* LD z80.IYL(),z80.IYH() */
					z80.iyh = z80.iyl
					break
				case 0x66: /* LD H,(iy+dd) */
					var offset byte
					offset = z80.memory.readByte(z80.pc)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.pc++
					z80.h = z80.memory.readByte(uint16(int(z80.IY()) + int(signExtend(offset))))
					break
				case 0x67: /* LD z80.IYL(),A */
					z80.iyh = z80.a
					break
				case 0x68: /* LD z80.IYH(),B */
					z80.iyl = z80.b
					break
				case 0x69: /* LD z80.IYH(),C */
					z80.iyl = z80.c
					break
				case 0x6a: /* LD z80.IYH(),D */
					z80.iyl = z80.d
					break
				case 0x6b: /* LD z80.IYH(),E */
					z80.iyl = z80.e
					break
				case 0x6c: /* LD z80.IYH(),z80.IYL() */
					z80.iyl = z80.iyh
					break
				case 0x6d: /* LD z80.IYH(),z80.IYH() */
					break
				case 0x6e: /* LD L,(iy+dd) */
					var offset byte
					offset = z80.memory.readByte(z80.pc)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.pc++
					z80.l = z80.memory.readByte(uint16(int(z80.IY()) + int(signExtend(offset))))
					break
				case 0x6f: /* LD z80.IYH(),A */
					z80.iyl = z80.a
					break
				case 0x70: /* LD (iy+dd),B */
					offset := z80.memory.readByte(z80.pc)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.pc++
					z80.memory.writeByte(uint16(int(z80.IY())+int(signExtend(offset))), z80.b)
					break
				case 0x71: /* LD (iy+dd),C */
					offset := z80.memory.readByte(z80.pc)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.pc++
					z80.memory.writeByte(uint16(int(z80.IY())+int(signExtend(offset))), z80.c)
					break
				case 0x72: /* LD (iy+dd),D */
					offset := z80.memory.readByte(z80.pc)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.pc++
					z80.memory.writeByte(uint16(int(z80.IY())+int(signExtend(offset))), z80.d)
					break
				case 0x73: /* LD (iy+dd),E */
					offset := z80.memory.readByte(z80.pc)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.pc++
					z80.memory.writeByte(uint16(int(z80.IY())+int(signExtend(offset))), z80.e)
					break
				case 0x74: /* LD (iy+dd),H */
					offset := z80.memory.readByte(z80.pc)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.pc++
					z80.memory.writeByte(uint16(int(z80.IY())+int(signExtend(offset))), z80.h)
					break
				case 0x75: /* LD (iy+dd),L */
					offset := z80.memory.readByte(z80.pc)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.pc++
					z80.memory.writeByte(uint16(int(z80.IY())+int(signExtend(offset))), z80.l)
					break
				case 0x77: /* LD (iy+dd),A */
					offset := z80.memory.readByte(z80.pc)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.pc++
					z80.memory.writeByte(uint16(int(z80.IY())+int(signExtend(offset))), z80.a)
					break
				case 0x7c: /* LD A,z80.IYL() */
					z80.a = z80.iyh
					break
				case 0x7d: /* LD A,z80.IYH() */
					z80.a = z80.iyl
					break
				case 0x7e: /* LD A,(iy+dd) */
					var offset byte
					offset = z80.memory.readByte(z80.pc)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.pc++
					z80.a = z80.memory.readByte(uint16(int(z80.IY()) + int(signExtend(offset))))
					break
				case 0x84: /* ADD A,z80.IYL() */
					z80.add(z80.iyh)
					break
				case 0x85: /* ADD A,z80.IYH() */
					z80.add(z80.iyl)
					break
				case 0x86: /* ADD A,(iy+dd) */

					var offset, bytetemp byte
					offset = z80.memory.readByte(z80.pc)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.pc++
					bytetemp = z80.memory.readByte(uint16(int(z80.IY()) + int(signExtend(offset))))
					z80.add(bytetemp)

					break
				case 0x8c: /* ADC A,z80.IYL() */
					z80.adc(z80.iyh)
					break
				case 0x8d: /* ADC A,z80.IYH() */
					z80.adc(z80.iyl)
					break
				case 0x8e: /* ADC A,(iy+dd) */

					var offset, bytetemp byte
					offset = z80.memory.readByte(z80.pc)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.pc++
					bytetemp = z80.memory.readByte(uint16(int(z80.IY()) + int(signExtend(offset))))
					z80.adc(bytetemp)

					break
				case 0x94: /* SUB A,z80.IYL() */
					z80.sub(z80.iyh)
					break
				case 0x95: /* SUB A,z80.IYH() */
					z80.sub(z80.iyl)
					break
				case 0x96: /* SUB A,(iy+dd) */

					var offset, bytetemp byte
					offset = z80.memory.readByte(z80.pc)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.pc++
					bytetemp = z80.memory.readByte(uint16(int(z80.IY()) + int(signExtend(offset))))
					z80.sub(bytetemp)

					break
				case 0x9c: /* SBC A,z80.IYL() */
					z80.sbc(z80.iyh)
					break
				case 0x9d: /* SBC A,z80.IYH() */
					z80.sbc(z80.iyl)
					break
				case 0x9e: /* SBC A,(iy+dd) */

					var offset, bytetemp byte
					offset = z80.memory.readByte(z80.pc)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.pc++
					bytetemp = z80.memory.readByte(uint16(int(z80.IY()) + int(signExtend(offset))))
					z80.sbc(bytetemp)

					break
				case 0xa4: /* AND A,z80.IYL() */
					z80.and(z80.iyh)
					break
				case 0xa5: /* AND A,z80.IYH() */
					z80.and(z80.iyl)
					break
				case 0xa6: /* AND A,(iy+dd) */

					var offset, bytetemp byte
					offset = z80.memory.readByte(z80.pc)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.pc++
					bytetemp = z80.memory.readByte(uint16(int(z80.IY()) + int(signExtend(offset))))
					z80.and(bytetemp)

					break
				case 0xac: /* XOR A,z80.IYL() */
					z80.xor(z80.iyh)
					break
				case 0xad: /* XOR A,z80.IYH() */
					z80.xor(z80.iyl)
					break
				case 0xae: /* XOR A,(iy+dd) */

					var offset, bytetemp byte
					offset = z80.memory.readByte(z80.pc)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.pc++
					bytetemp = z80.memory.readByte(uint16(int(z80.IY()) + int(signExtend(offset))))
					z80.xor(bytetemp)

					break
				case 0xb4: /* OR A,z80.IYL() */
					z80.or(z80.iyh)
					break
				case 0xb5: /* OR A,z80.IYH() */
					z80.or(z80.iyl)
					break
				case 0xb6: /* OR A,(iy+dd) */

					var offset, bytetemp byte
					offset = z80.memory.readByte(z80.pc)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.pc++
					bytetemp = z80.memory.readByte(uint16(int(z80.IY()) + int(signExtend(offset))))
					z80.or(bytetemp)

					break
				case 0xbc: /* CP A,z80.IYL() */
					z80.cp(z80.iyh)
					break
				case 0xbd: /* CP A,z80.IYH() */
					z80.cp(z80.iyl)
					break
				case 0xbe: /* CP A,(iy+dd) */

					var offset, bytetemp byte
					offset = z80.memory.readByte(z80.pc)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.pc++
					bytetemp = z80.memory.readByte(uint16(int(z80.IY()) + int(signExtend(offset))))
					z80.cp(bytetemp)

					break
				case 0xcb: /* shift DDFDCB */

					var tempaddr uint16
					var opcode3 byte
					z80.memory.contendRead(z80.pc, 3)
					tempaddr = uint16(int(z80.IY()) + int(signExtend(z80.memory.readByteInternal(z80.pc))))
					z80.pc++
					z80.memory.contendRead(z80.pc, 3)
					opcode3 = z80.memory.readByteInternal(z80.pc)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.memory.contendReadNoMreq(z80.pc, 1)
					z80.pc++

					switch opcode3 {
					/* z80_ddfdcb.c Z80 {DD,FD}CBxx opcodes
					   Copyright (c) 1999-2003 Philip Kendall

					   This program is free software; you can redistribute it and/or modify
					   it under the terms of the GNU General Public License as published by
					   the Free Software Foundation; either version 2 of the License, or
					   (at your option) any later version.

					   This program is distributed in the hope that it will be useful,
					   but WITHOUT ANY WARRANTY; without even the implied warranty of
					   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
					   GNU General Public License for more details.

					   You should have received a copy of the GNU General Public License along
					   with this program; if not, write to the Free Software Foundation, Inc.,
					   51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.

					   Author contact information:

					   E-mail: philip-fuse@shadowmagic.org.uk

					*/

					/* NB: this file is autogenerated by './z80.pl' from 'opcodes_ddfdcb.dat',
					   and included in 'z80_ops.c' */

					case 0x00: /* LD B,RLC (REGISTER+dd) */
						z80.b = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rlc(&z80.b)
						z80.memory.writeByte(tempaddr, z80.b)
						break
					case 0x01: /* LD C,RLC (REGISTER+dd) */
						z80.c = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rlc(&z80.c)
						z80.memory.writeByte(tempaddr, z80.c)
						break
					case 0x02: /* LD D,RLC (REGISTER+dd) */
						z80.d = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rlc(&z80.d)
						z80.memory.writeByte(tempaddr, z80.d)
						break
					case 0x03: /* LD E,RLC (REGISTER+dd) */
						z80.e = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rlc(&z80.e)
						z80.memory.writeByte(tempaddr, z80.e)
						break
					case 0x04: /* LD H,RLC (REGISTER+dd) */
						z80.h = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rlc(&z80.h)
						z80.memory.writeByte(tempaddr, z80.h)
						break
					case 0x05: /* LD L,RLC (REGISTER+dd) */
						z80.l = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rlc(&z80.l)
						z80.memory.writeByte(tempaddr, z80.l)
						break
					case 0x06: /* RLC (REGISTER+dd) */
						var bytetemp byte = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rlc(&bytetemp)
						z80.memory.writeByte(tempaddr, bytetemp)
						break
					case 0x07: /* LD A,RLC (REGISTER+dd) */
						z80.a = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rlc(&z80.a)
						z80.memory.writeByte(tempaddr, z80.a)
						break
					case 0x08: /* LD B,RRC (REGISTER+dd) */
						z80.b = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rrc(&z80.b)
						z80.memory.writeByte(tempaddr, z80.b)
						break
					case 0x09: /* LD C,RRC (REGISTER+dd) */
						z80.c = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rrc(&z80.c)
						z80.memory.writeByte(tempaddr, z80.c)
						break
					case 0x0a: /* LD D,RRC (REGISTER+dd) */
						z80.d = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rrc(&z80.d)
						z80.memory.writeByte(tempaddr, z80.d)
						break
					case 0x0b: /* LD E,RRC (REGISTER+dd) */
						z80.e = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rrc(&z80.e)
						z80.memory.writeByte(tempaddr, z80.e)
						break
					case 0x0c: /* LD H,RRC (REGISTER+dd) */
						z80.h = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rrc(&z80.h)
						z80.memory.writeByte(tempaddr, z80.h)
						break
					case 0x0d: /* LD L,RRC (REGISTER+dd) */
						z80.l = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rrc(&z80.l)
						z80.memory.writeByte(tempaddr, z80.l)
						break
					case 0x0e: /* RRC (REGISTER+dd) */
						var bytetemp byte = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rrc(&bytetemp)
						z80.memory.writeByte(tempaddr, bytetemp)
						break
					case 0x0f: /* LD A,RRC (REGISTER+dd) */
						z80.a = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rrc(&z80.a)
						z80.memory.writeByte(tempaddr, z80.a)
						break
					case 0x10: /* LD B,RL (REGISTER+dd) */
						z80.b = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rl(&z80.b)
						z80.memory.writeByte(tempaddr, z80.b)
						break
					case 0x11: /* LD C,RL (REGISTER+dd) */
						z80.c = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rl(&z80.c)
						z80.memory.writeByte(tempaddr, z80.c)
						break
					case 0x12: /* LD D,RL (REGISTER+dd) */
						z80.d = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rl(&z80.d)
						z80.memory.writeByte(tempaddr, z80.d)
						break
					case 0x13: /* LD E,RL (REGISTER+dd) */
						z80.e = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rl(&z80.e)
						z80.memory.writeByte(tempaddr, z80.e)
						break
					case 0x14: /* LD H,RL (REGISTER+dd) */
						z80.h = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rl(&z80.h)
						z80.memory.writeByte(tempaddr, z80.h)
						break
					case 0x15: /* LD L,RL (REGISTER+dd) */
						z80.l = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rl(&z80.l)
						z80.memory.writeByte(tempaddr, z80.l)
						break
					case 0x16: /* RL (REGISTER+dd) */
						var bytetemp byte = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rl(&bytetemp)
						z80.memory.writeByte(tempaddr, bytetemp)
						break
					case 0x17: /* LD A,RL (REGISTER+dd) */
						z80.a = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rl(&z80.a)
						z80.memory.writeByte(tempaddr, z80.a)
						break
					case 0x18: /* LD B,RR (REGISTER+dd) */
						z80.b = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rr(&z80.b)
						z80.memory.writeByte(tempaddr, z80.b)
						break
					case 0x19: /* LD C,RR (REGISTER+dd) */
						z80.c = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rr(&z80.c)
						z80.memory.writeByte(tempaddr, z80.c)
						break
					case 0x1a: /* LD D,RR (REGISTER+dd) */
						z80.d = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rr(&z80.d)
						z80.memory.writeByte(tempaddr, z80.d)
						break
					case 0x1b: /* LD E,RR (REGISTER+dd) */
						z80.e = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rr(&z80.e)
						z80.memory.writeByte(tempaddr, z80.e)
						break
					case 0x1c: /* LD H,RR (REGISTER+dd) */
						z80.h = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rr(&z80.h)
						z80.memory.writeByte(tempaddr, z80.h)
						break
					case 0x1d: /* LD L,RR (REGISTER+dd) */
						z80.l = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rr(&z80.l)
						z80.memory.writeByte(tempaddr, z80.l)
						break
					case 0x1e: /* RR (REGISTER+dd) */
						var bytetemp byte = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rr(&bytetemp)
						z80.memory.writeByte(tempaddr, bytetemp)
						break
					case 0x1f: /* LD A,RR (REGISTER+dd) */
						z80.a = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.rr(&z80.a)
						z80.memory.writeByte(tempaddr, z80.a)
						break
					case 0x20: /* LD B,SLA (REGISTER+dd) */
						z80.b = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.sla(&z80.b)
						z80.memory.writeByte(tempaddr, z80.b)
						break
					case 0x21: /* LD C,SLA (REGISTER+dd) */
						z80.c = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.sla(&z80.c)
						z80.memory.writeByte(tempaddr, z80.c)
						break
					case 0x22: /* LD D,SLA (REGISTER+dd) */
						z80.d = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.sla(&z80.d)
						z80.memory.writeByte(tempaddr, z80.d)
						break
					case 0x23: /* LD E,SLA (REGISTER+dd) */
						z80.e = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.sla(&z80.e)
						z80.memory.writeByte(tempaddr, z80.e)
						break
					case 0x24: /* LD H,SLA (REGISTER+dd) */
						z80.h = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.sla(&z80.h)
						z80.memory.writeByte(tempaddr, z80.h)
						break
					case 0x25: /* LD L,SLA (REGISTER+dd) */
						z80.l = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.sla(&z80.l)
						z80.memory.writeByte(tempaddr, z80.l)
						break
					case 0x26: /* SLA (REGISTER+dd) */
						var bytetemp byte = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.sla(&bytetemp)
						z80.memory.writeByte(tempaddr, bytetemp)
						break
					case 0x27: /* LD A,SLA (REGISTER+dd) */
						z80.a = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.sla(&z80.a)
						z80.memory.writeByte(tempaddr, z80.a)
						break
					case 0x28: /* LD B,SRA (REGISTER+dd) */
						z80.b = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.sra(&z80.b)
						z80.memory.writeByte(tempaddr, z80.b)
						break
					case 0x29: /* LD C,SRA (REGISTER+dd) */
						z80.c = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.sra(&z80.c)
						z80.memory.writeByte(tempaddr, z80.c)
						break
					case 0x2a: /* LD D,SRA (REGISTER+dd) */
						z80.d = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.sra(&z80.d)
						z80.memory.writeByte(tempaddr, z80.d)
						break
					case 0x2b: /* LD E,SRA (REGISTER+dd) */
						z80.e = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.sra(&z80.e)
						z80.memory.writeByte(tempaddr, z80.e)
						break
					case 0x2c: /* LD H,SRA (REGISTER+dd) */
						z80.h = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.sra(&z80.h)
						z80.memory.writeByte(tempaddr, z80.h)
						break
					case 0x2d: /* LD L,SRA (REGISTER+dd) */
						z80.l = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.sra(&z80.l)
						z80.memory.writeByte(tempaddr, z80.l)
						break
					case 0x2e: /* SRA (REGISTER+dd) */
						var bytetemp byte = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.sra(&bytetemp)
						z80.memory.writeByte(tempaddr, bytetemp)
						break
					case 0x2f: /* LD A,SRA (REGISTER+dd) */
						z80.a = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.sra(&z80.a)
						z80.memory.writeByte(tempaddr, z80.a)
						break
					case 0x30: /* LD B,SLL (REGISTER+dd) */
						z80.b = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.sll(&z80.b)
						z80.memory.writeByte(tempaddr, z80.b)
						break
					case 0x31: /* LD C,SLL (REGISTER+dd) */
						z80.c = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.sll(&z80.c)
						z80.memory.writeByte(tempaddr, z80.c)
						break
					case 0x32: /* LD D,SLL (REGISTER+dd) */
						z80.d = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.sll(&z80.d)
						z80.memory.writeByte(tempaddr, z80.d)
						break
					case 0x33: /* LD E,SLL (REGISTER+dd) */
						z80.e = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.sll(&z80.e)
						z80.memory.writeByte(tempaddr, z80.e)
						break
					case 0x34: /* LD H,SLL (REGISTER+dd) */
						z80.h = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.sll(&z80.h)
						z80.memory.writeByte(tempaddr, z80.h)
						break
					case 0x35: /* LD L,SLL (REGISTER+dd) */
						z80.l = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.sll(&z80.l)
						z80.memory.writeByte(tempaddr, z80.l)
						break
					case 0x36: /* SLL (REGISTER+dd) */
						var bytetemp byte = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.sll(&bytetemp)
						z80.memory.writeByte(tempaddr, bytetemp)
						break
					case 0x37: /* LD A,SLL (REGISTER+dd) */
						z80.a = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.sll(&z80.a)
						z80.memory.writeByte(tempaddr, z80.a)
						break
					case 0x38: /* LD B,SRL (REGISTER+dd) */
						z80.b = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.srl(&z80.b)
						z80.memory.writeByte(tempaddr, z80.b)
						break
					case 0x39: /* LD C,SRL (REGISTER+dd) */
						z80.c = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.srl(&z80.c)
						z80.memory.writeByte(tempaddr, z80.c)
						break
					case 0x3a: /* LD D,SRL (REGISTER+dd) */
						z80.d = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.srl(&z80.d)
						z80.memory.writeByte(tempaddr, z80.d)
						break
					case 0x3b: /* LD E,SRL (REGISTER+dd) */
						z80.e = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.srl(&z80.e)
						z80.memory.writeByte(tempaddr, z80.e)
						break
					case 0x3c: /* LD H,SRL (REGISTER+dd) */
						z80.h = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.srl(&z80.h)
						z80.memory.writeByte(tempaddr, z80.h)
						break
					case 0x3d: /* LD L,SRL (REGISTER+dd) */
						z80.l = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.srl(&z80.l)
						z80.memory.writeByte(tempaddr, z80.l)
						break
					case 0x3e: /* SRL (REGISTER+dd) */
						var bytetemp byte = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.srl(&bytetemp)
						z80.memory.writeByte(tempaddr, bytetemp)
						break
					case 0x3f: /* LD A,SRL (REGISTER+dd) */
						z80.a = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.srl(&z80.a)
						z80.memory.writeByte(tempaddr, z80.a)
						break
					case 0x40:
						fallthrough
					case 0x41:
						fallthrough
					case 0x42:
						fallthrough
					case 0x43:
						fallthrough
					case 0x44:
						fallthrough
					case 0x45:
						fallthrough
					case 0x46:
						fallthrough
					case 0x47: /* BIT 0,(REGISTER+dd) */
						bytetemp := z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.biti(0, bytetemp, tempaddr)
						break
					case 0x48:
						fallthrough
					case 0x49:
						fallthrough
					case 0x4a:
						fallthrough
					case 0x4b:
						fallthrough
					case 0x4c:
						fallthrough
					case 0x4d:
						fallthrough
					case 0x4e:
						fallthrough
					case 0x4f: /* BIT 1,(REGISTER+dd) */
						bytetemp := z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.biti(1, bytetemp, tempaddr)
						break
					case 0x50:
						fallthrough
					case 0x51:
						fallthrough
					case 0x52:
						fallthrough
					case 0x53:
						fallthrough
					case 0x54:
						fallthrough
					case 0x55:
						fallthrough
					case 0x56:
						fallthrough
					case 0x57: /* BIT 2,(REGISTER+dd) */
						bytetemp := z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.biti(2, bytetemp, tempaddr)
						break
					case 0x58:
						fallthrough
					case 0x59:
						fallthrough
					case 0x5a:
						fallthrough
					case 0x5b:
						fallthrough
					case 0x5c:
						fallthrough
					case 0x5d:
						fallthrough
					case 0x5e:
						fallthrough
					case 0x5f: /* BIT 3,(REGISTER+dd) */
						bytetemp := z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.biti(3, bytetemp, tempaddr)
						break
					case 0x60:
						fallthrough
					case 0x61:
						fallthrough
					case 0x62:
						fallthrough
					case 0x63:
						fallthrough
					case 0x64:
						fallthrough
					case 0x65:
						fallthrough
					case 0x66:
						fallthrough
					case 0x67: /* BIT 4,(REGISTER+dd) */
						bytetemp := z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.biti(4, bytetemp, tempaddr)
						break
					case 0x68:
						fallthrough
					case 0x69:
						fallthrough
					case 0x6a:
						fallthrough
					case 0x6b:
						fallthrough
					case 0x6c:
						fallthrough
					case 0x6d:
						fallthrough
					case 0x6e:
						fallthrough
					case 0x6f: /* BIT 5,(REGISTER+dd) */
						bytetemp := z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.biti(5, bytetemp, tempaddr)
						break
					case 0x70:
						fallthrough
					case 0x71:
						fallthrough
					case 0x72:
						fallthrough
					case 0x73:
						fallthrough
					case 0x74:
						fallthrough
					case 0x75:
						fallthrough
					case 0x76:
						fallthrough
					case 0x77: /* BIT 6,(REGISTER+dd) */
						bytetemp := z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.biti(6, bytetemp, tempaddr)
						break
					case 0x78:
						fallthrough
					case 0x79:
						fallthrough
					case 0x7a:
						fallthrough
					case 0x7b:
						fallthrough
					case 0x7c:
						fallthrough
					case 0x7d:
						fallthrough
					case 0x7e:
						fallthrough
					case 0x7f: /* BIT 7,(REGISTER+dd) */
						bytetemp := z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.biti(7, bytetemp, tempaddr)
						break
					case 0x80: /* LD B,RES 0,(REGISTER+dd) */
						z80.b = z80.memory.readByte(tempaddr) & 0xfe
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.b)
						break
					case 0x81: /* LD C,RES 0,(REGISTER+dd) */
						z80.c = z80.memory.readByte(tempaddr) & 0xfe
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.c)
						break
					case 0x82: /* LD D,RES 0,(REGISTER+dd) */
						z80.d = z80.memory.readByte(tempaddr) & 0xfe
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.d)
						break
					case 0x83: /* LD E,RES 0,(REGISTER+dd) */
						z80.e = z80.memory.readByte(tempaddr) & 0xfe
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.e)
						break
					case 0x84: /* LD H,RES 0,(REGISTER+dd) */
						z80.h = z80.memory.readByte(tempaddr) & 0xfe
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.h)
						break
					case 0x85: /* LD L,RES 0,(REGISTER+dd) */
						z80.l = z80.memory.readByte(tempaddr) & 0xfe
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.l)
						break
					case 0x86: /* RES 0,(REGISTER+dd) */

						var bytetemp byte
						bytetemp = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, bytetemp&0xfe)

						break
					case 0x87: /* LD A,RES 0,(REGISTER+dd) */
						z80.a = z80.memory.readByte(tempaddr) & 0xfe
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.a)
						break
					case 0x88: /* LD B,RES 1,(REGISTER+dd) */
						z80.b = z80.memory.readByte(tempaddr) & 0xfd
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.b)
						break
					case 0x89: /* LD C,RES 1,(REGISTER+dd) */
						z80.c = z80.memory.readByte(tempaddr) & 0xfd
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.c)
						break
					case 0x8a: /* LD D,RES 1,(REGISTER+dd) */
						z80.d = z80.memory.readByte(tempaddr) & 0xfd
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.d)
						break
					case 0x8b: /* LD E,RES 1,(REGISTER+dd) */
						z80.e = z80.memory.readByte(tempaddr) & 0xfd
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.e)
						break
					case 0x8c: /* LD H,RES 1,(REGISTER+dd) */
						z80.h = z80.memory.readByte(tempaddr) & 0xfd
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.h)
						break
					case 0x8d: /* LD L,RES 1,(REGISTER+dd) */
						z80.l = z80.memory.readByte(tempaddr) & 0xfd
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.l)
						break
					case 0x8e: /* RES 1,(REGISTER+dd) */

						var bytetemp byte
						bytetemp = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, bytetemp&0xfd)

						break
					case 0x8f: /* LD A,RES 1,(REGISTER+dd) */
						z80.a = z80.memory.readByte(tempaddr) & 0xfd
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.a)
						break
					case 0x90: /* LD B,RES 2,(REGISTER+dd) */
						z80.b = z80.memory.readByte(tempaddr) & 0xfb
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.b)
						break
					case 0x91: /* LD C,RES 2,(REGISTER+dd) */
						z80.c = z80.memory.readByte(tempaddr) & 0xfb
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.c)
						break
					case 0x92: /* LD D,RES 2,(REGISTER+dd) */
						z80.d = z80.memory.readByte(tempaddr) & 0xfb
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.d)
						break
					case 0x93: /* LD E,RES 2,(REGISTER+dd) */
						z80.e = z80.memory.readByte(tempaddr) & 0xfb
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.e)
						break
					case 0x94: /* LD H,RES 2,(REGISTER+dd) */
						z80.h = z80.memory.readByte(tempaddr) & 0xfb
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.h)
						break
					case 0x95: /* LD L,RES 2,(REGISTER+dd) */
						z80.l = z80.memory.readByte(tempaddr) & 0xfb
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.l)
						break
					case 0x96: /* RES 2,(REGISTER+dd) */

						var bytetemp byte
						bytetemp = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, bytetemp&0xfb)

						break
					case 0x97: /* LD A,RES 2,(REGISTER+dd) */
						z80.a = z80.memory.readByte(tempaddr) & 0xfb
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.a)
						break
					case 0x98: /* LD B,RES 3,(REGISTER+dd) */
						z80.b = z80.memory.readByte(tempaddr) & 0xf7
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.b)
						break
					case 0x99: /* LD C,RES 3,(REGISTER+dd) */
						z80.c = z80.memory.readByte(tempaddr) & 0xf7
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.c)
						break
					case 0x9a: /* LD D,RES 3,(REGISTER+dd) */
						z80.d = z80.memory.readByte(tempaddr) & 0xf7
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.d)
						break
					case 0x9b: /* LD E,RES 3,(REGISTER+dd) */
						z80.e = z80.memory.readByte(tempaddr) & 0xf7
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.e)
						break
					case 0x9c: /* LD H,RES 3,(REGISTER+dd) */
						z80.h = z80.memory.readByte(tempaddr) & 0xf7
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.h)
						break
					case 0x9d: /* LD L,RES 3,(REGISTER+dd) */
						z80.l = z80.memory.readByte(tempaddr) & 0xf7
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.l)
						break
					case 0x9e: /* RES 3,(REGISTER+dd) */

						var bytetemp byte
						bytetemp = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, bytetemp&0xf7)

						break
					case 0x9f: /* LD A,RES 3,(REGISTER+dd) */
						z80.a = z80.memory.readByte(tempaddr) & 0xf7
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.a)
						break
					case 0xa0: /* LD B,RES 4,(REGISTER+dd) */
						z80.b = z80.memory.readByte(tempaddr) & 0xef
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.b)
						break
					case 0xa1: /* LD C,RES 4,(REGISTER+dd) */
						z80.c = z80.memory.readByte(tempaddr) & 0xef
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.c)
						break
					case 0xa2: /* LD D,RES 4,(REGISTER+dd) */
						z80.d = z80.memory.readByte(tempaddr) & 0xef
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.d)
						break
					case 0xa3: /* LD E,RES 4,(REGISTER+dd) */
						z80.e = z80.memory.readByte(tempaddr) & 0xef
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.e)
						break
					case 0xa4: /* LD H,RES 4,(REGISTER+dd) */
						z80.h = z80.memory.readByte(tempaddr) & 0xef
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.h)
						break
					case 0xa5: /* LD L,RES 4,(REGISTER+dd) */
						z80.l = z80.memory.readByte(tempaddr) & 0xef
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.l)
						break
					case 0xa6: /* RES 4,(REGISTER+dd) */

						var bytetemp byte
						bytetemp = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, bytetemp&0xef)

						break
					case 0xa7: /* LD A,RES 4,(REGISTER+dd) */
						z80.a = z80.memory.readByte(tempaddr) & 0xef
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.a)
						break
					case 0xa8: /* LD B,RES 5,(REGISTER+dd) */
						z80.b = z80.memory.readByte(tempaddr) & 0xdf
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.b)
						break
					case 0xa9: /* LD C,RES 5,(REGISTER+dd) */
						z80.c = z80.memory.readByte(tempaddr) & 0xdf
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.c)
						break
					case 0xaa: /* LD D,RES 5,(REGISTER+dd) */
						z80.d = z80.memory.readByte(tempaddr) & 0xdf
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.d)
						break
					case 0xab: /* LD E,RES 5,(REGISTER+dd) */
						z80.e = z80.memory.readByte(tempaddr) & 0xdf
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.e)
						break
					case 0xac: /* LD H,RES 5,(REGISTER+dd) */
						z80.h = z80.memory.readByte(tempaddr) & 0xdf
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.h)
						break
					case 0xad: /* LD L,RES 5,(REGISTER+dd) */
						z80.l = z80.memory.readByte(tempaddr) & 0xdf
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.l)
						break
					case 0xae: /* RES 5,(REGISTER+dd) */

						var bytetemp byte
						bytetemp = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, bytetemp&0xdf)

						break
					case 0xaf: /* LD A,RES 5,(REGISTER+dd) */
						z80.a = z80.memory.readByte(tempaddr) & 0xdf
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.a)
						break
					case 0xb0: /* LD B,RES 6,(REGISTER+dd) */
						z80.b = z80.memory.readByte(tempaddr) & 0xbf
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.b)
						break
					case 0xb1: /* LD C,RES 6,(REGISTER+dd) */
						z80.c = z80.memory.readByte(tempaddr) & 0xbf
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.c)
						break
					case 0xb2: /* LD D,RES 6,(REGISTER+dd) */
						z80.d = z80.memory.readByte(tempaddr) & 0xbf
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.d)
						break
					case 0xb3: /* LD E,RES 6,(REGISTER+dd) */
						z80.e = z80.memory.readByte(tempaddr) & 0xbf
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.e)
						break
					case 0xb4: /* LD H,RES 6,(REGISTER+dd) */
						z80.h = z80.memory.readByte(tempaddr) & 0xbf
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.h)
						break
					case 0xb5: /* LD L,RES 6,(REGISTER+dd) */
						z80.l = z80.memory.readByte(tempaddr) & 0xbf
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.l)
						break
					case 0xb6: /* RES 6,(REGISTER+dd) */

						var bytetemp byte
						bytetemp = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, bytetemp&0xbf)

						break
					case 0xb7: /* LD A,RES 6,(REGISTER+dd) */
						z80.a = z80.memory.readByte(tempaddr) & 0xbf
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.a)
						break
					case 0xb8: /* LD B,RES 7,(REGISTER+dd) */
						z80.b = z80.memory.readByte(tempaddr) & 0x7f
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.b)
						break
					case 0xb9: /* LD C,RES 7,(REGISTER+dd) */
						z80.c = z80.memory.readByte(tempaddr) & 0x7f
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.c)
						break
					case 0xba: /* LD D,RES 7,(REGISTER+dd) */
						z80.d = z80.memory.readByte(tempaddr) & 0x7f
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.d)
						break
					case 0xbb: /* LD E,RES 7,(REGISTER+dd) */
						z80.e = z80.memory.readByte(tempaddr) & 0x7f
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.e)
						break
					case 0xbc: /* LD H,RES 7,(REGISTER+dd) */
						z80.h = z80.memory.readByte(tempaddr) & 0x7f
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.h)
						break
					case 0xbd: /* LD L,RES 7,(REGISTER+dd) */
						z80.l = z80.memory.readByte(tempaddr) & 0x7f
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.l)
						break
					case 0xbe: /* RES 7,(REGISTER+dd) */

						var bytetemp byte
						bytetemp = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, bytetemp&0x7f)

						break
					case 0xbf: /* LD A,RES 7,(REGISTER+dd) */
						z80.a = z80.memory.readByte(tempaddr) & 0x7f
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.a)
						break
					case 0xc0: /* LD B,SET 0,(REGISTER+dd) */
						z80.b = z80.memory.readByte(tempaddr) | 0x01
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.b)
						break
					case 0xc1: /* LD C,SET 0,(REGISTER+dd) */
						z80.c = z80.memory.readByte(tempaddr) | 0x01
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.c)
						break
					case 0xc2: /* LD D,SET 0,(REGISTER+dd) */
						z80.d = z80.memory.readByte(tempaddr) | 0x01
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.d)
						break
					case 0xc3: /* LD E,SET 0,(REGISTER+dd) */
						z80.e = z80.memory.readByte(tempaddr) | 0x01
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.e)
						break
					case 0xc4: /* LD H,SET 0,(REGISTER+dd) */
						z80.h = z80.memory.readByte(tempaddr) | 0x01
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.h)
						break
					case 0xc5: /* LD L,SET 0,(REGISTER+dd) */
						z80.l = z80.memory.readByte(tempaddr) | 0x01
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.l)
						break
					case 0xc6: /* SET 0,(REGISTER+dd) */

						var bytetemp byte
						bytetemp = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, bytetemp|0x01)

						break
					case 0xc7: /* LD A,SET 0,(REGISTER+dd) */
						z80.a = z80.memory.readByte(tempaddr) | 0x01
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.a)
						break
					case 0xc8: /* LD B,SET 1,(REGISTER+dd) */
						z80.b = z80.memory.readByte(tempaddr) | 0x02
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.b)
						break
					case 0xc9: /* LD C,SET 1,(REGISTER+dd) */
						z80.c = z80.memory.readByte(tempaddr) | 0x02
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.c)
						break
					case 0xca: /* LD D,SET 1,(REGISTER+dd) */
						z80.d = z80.memory.readByte(tempaddr) | 0x02
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.d)
						break
					case 0xcb: /* LD E,SET 1,(REGISTER+dd) */
						z80.e = z80.memory.readByte(tempaddr) | 0x02
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.e)
						break
					case 0xcc: /* LD H,SET 1,(REGISTER+dd) */
						z80.h = z80.memory.readByte(tempaddr) | 0x02
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.h)
						break
					case 0xcd: /* LD L,SET 1,(REGISTER+dd) */
						z80.l = z80.memory.readByte(tempaddr) | 0x02
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.l)
						break
					case 0xce: /* SET 1,(REGISTER+dd) */

						var bytetemp byte
						bytetemp = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, bytetemp|0x02)

						break
					case 0xcf: /* LD A,SET 1,(REGISTER+dd) */
						z80.a = z80.memory.readByte(tempaddr) | 0x02
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.a)
						break
					case 0xd0: /* LD B,SET 2,(REGISTER+dd) */
						z80.b = z80.memory.readByte(tempaddr) | 0x04
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.b)
						break
					case 0xd1: /* LD C,SET 2,(REGISTER+dd) */
						z80.c = z80.memory.readByte(tempaddr) | 0x04
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.c)
						break
					case 0xd2: /* LD D,SET 2,(REGISTER+dd) */
						z80.d = z80.memory.readByte(tempaddr) | 0x04
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.d)
						break
					case 0xd3: /* LD E,SET 2,(REGISTER+dd) */
						z80.e = z80.memory.readByte(tempaddr) | 0x04
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.e)
						break
					case 0xd4: /* LD H,SET 2,(REGISTER+dd) */
						z80.h = z80.memory.readByte(tempaddr) | 0x04
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.h)
						break
					case 0xd5: /* LD L,SET 2,(REGISTER+dd) */
						z80.l = z80.memory.readByte(tempaddr) | 0x04
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.l)
						break
					case 0xd6: /* SET 2,(REGISTER+dd) */

						var bytetemp byte
						bytetemp = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, bytetemp|0x04)

						break
					case 0xd7: /* LD A,SET 2,(REGISTER+dd) */
						z80.a = z80.memory.readByte(tempaddr) | 0x04
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.a)
						break
					case 0xd8: /* LD B,SET 3,(REGISTER+dd) */
						z80.b = z80.memory.readByte(tempaddr) | 0x08
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.b)
						break
					case 0xd9: /* LD C,SET 3,(REGISTER+dd) */
						z80.c = z80.memory.readByte(tempaddr) | 0x08
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.c)
						break
					case 0xda: /* LD D,SET 3,(REGISTER+dd) */
						z80.d = z80.memory.readByte(tempaddr) | 0x08
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.d)
						break
					case 0xdb: /* LD E,SET 3,(REGISTER+dd) */
						z80.e = z80.memory.readByte(tempaddr) | 0x08
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.e)
						break
					case 0xdc: /* LD H,SET 3,(REGISTER+dd) */
						z80.h = z80.memory.readByte(tempaddr) | 0x08
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.h)
						break
					case 0xdd: /* LD L,SET 3,(REGISTER+dd) */
						z80.l = z80.memory.readByte(tempaddr) | 0x08
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.l)
						break
					case 0xde: /* SET 3,(REGISTER+dd) */

						var bytetemp byte
						bytetemp = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, bytetemp|0x08)

						break
					case 0xdf: /* LD A,SET 3,(REGISTER+dd) */
						z80.a = z80.memory.readByte(tempaddr) | 0x08
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.a)
						break
					case 0xe0: /* LD B,SET 4,(REGISTER+dd) */
						z80.b = z80.memory.readByte(tempaddr) | 0x10
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.b)
						break
					case 0xe1: /* LD C,SET 4,(REGISTER+dd) */
						z80.c = z80.memory.readByte(tempaddr) | 0x10
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.c)
						break
					case 0xe2: /* LD D,SET 4,(REGISTER+dd) */
						z80.d = z80.memory.readByte(tempaddr) | 0x10
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.d)
						break
					case 0xe3: /* LD E,SET 4,(REGISTER+dd) */
						z80.e = z80.memory.readByte(tempaddr) | 0x10
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.e)
						break
					case 0xe4: /* LD H,SET 4,(REGISTER+dd) */
						z80.h = z80.memory.readByte(tempaddr) | 0x10
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.h)
						break
					case 0xe5: /* LD L,SET 4,(REGISTER+dd) */
						z80.l = z80.memory.readByte(tempaddr) | 0x10
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.l)
						break
					case 0xe6: /* SET 4,(REGISTER+dd) */

						var bytetemp byte
						bytetemp = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, bytetemp|0x10)

						break
					case 0xe7: /* LD A,SET 4,(REGISTER+dd) */
						z80.a = z80.memory.readByte(tempaddr) | 0x10
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.a)
						break
					case 0xe8: /* LD B,SET 5,(REGISTER+dd) */
						z80.b = z80.memory.readByte(tempaddr) | 0x20
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.b)
						break
					case 0xe9: /* LD C,SET 5,(REGISTER+dd) */
						z80.c = z80.memory.readByte(tempaddr) | 0x20
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.c)
						break
					case 0xea: /* LD D,SET 5,(REGISTER+dd) */
						z80.d = z80.memory.readByte(tempaddr) | 0x20
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.d)
						break
					case 0xeb: /* LD E,SET 5,(REGISTER+dd) */
						z80.e = z80.memory.readByte(tempaddr) | 0x20
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.e)
						break
					case 0xec: /* LD H,SET 5,(REGISTER+dd) */
						z80.h = z80.memory.readByte(tempaddr) | 0x20
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.h)
						break
					case 0xed: /* LD L,SET 5,(REGISTER+dd) */
						z80.l = z80.memory.readByte(tempaddr) | 0x20
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.l)
						break
					case 0xee: /* SET 5,(REGISTER+dd) */

						var bytetemp byte
						bytetemp = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, bytetemp|0x20)

						break
					case 0xef: /* LD A,SET 5,(REGISTER+dd) */
						z80.a = z80.memory.readByte(tempaddr) | 0x20
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.a)
						break
					case 0xf0: /* LD B,SET 6,(REGISTER+dd) */
						z80.b = z80.memory.readByte(tempaddr) | 0x40
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.b)
						break
					case 0xf1: /* LD C,SET 6,(REGISTER+dd) */
						z80.c = z80.memory.readByte(tempaddr) | 0x40
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.c)
						break
					case 0xf2: /* LD D,SET 6,(REGISTER+dd) */
						z80.d = z80.memory.readByte(tempaddr) | 0x40
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.d)
						break
					case 0xf3: /* LD E,SET 6,(REGISTER+dd) */
						z80.e = z80.memory.readByte(tempaddr) | 0x40
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.e)
						break
					case 0xf4: /* LD H,SET 6,(REGISTER+dd) */
						z80.h = z80.memory.readByte(tempaddr) | 0x40
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.h)
						break
					case 0xf5: /* LD L,SET 6,(REGISTER+dd) */
						z80.l = z80.memory.readByte(tempaddr) | 0x40
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.l)
						break
					case 0xf6: /* SET 6,(REGISTER+dd) */

						var bytetemp byte
						bytetemp = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, bytetemp|0x40)

						break
					case 0xf7: /* LD A,SET 6,(REGISTER+dd) */
						z80.a = z80.memory.readByte(tempaddr) | 0x40
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.a)
						break
					case 0xf8: /* LD B,SET 7,(REGISTER+dd) */
						z80.b = z80.memory.readByte(tempaddr) | 0x80
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.b)
						break
					case 0xf9: /* LD C,SET 7,(REGISTER+dd) */
						z80.c = z80.memory.readByte(tempaddr) | 0x80
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.c)
						break
					case 0xfa: /* LD D,SET 7,(REGISTER+dd) */
						z80.d = z80.memory.readByte(tempaddr) | 0x80
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.d)
						break
					case 0xfb: /* LD E,SET 7,(REGISTER+dd) */
						z80.e = z80.memory.readByte(tempaddr) | 0x80
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.e)
						break
					case 0xfc: /* LD H,SET 7,(REGISTER+dd) */
						z80.h = z80.memory.readByte(tempaddr) | 0x80
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.h)
						break
					case 0xfd: /* LD L,SET 7,(REGISTER+dd) */
						z80.l = z80.memory.readByte(tempaddr) | 0x80
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.l)
						break
					case 0xfe: /* SET 7,(REGISTER+dd) */

						var bytetemp byte
						bytetemp = z80.memory.readByte(tempaddr)
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, bytetemp|0x80)

						break
					case 0xff: /* LD A,SET 7,(REGISTER+dd) */
						z80.a = z80.memory.readByte(tempaddr) | 0x80
						z80.memory.contendReadNoMreq(tempaddr, 1)
						z80.memory.writeByte(tempaddr, z80.a)
						break

					}

					break
				case 0xe1: /* POP iy */
					z80.pop16(&z80.iyl, &z80.iyh)
					break
				case 0xe3: /* EX (SP),iy */
					var bytetempl, bytetemph byte
					bytetempl = z80.memory.readByte(z80.SP())
					bytetemph = z80.memory.readByte(z80.SP() + 1)
					z80.memory.contendReadNoMreq(z80.SP()+1, 1)
					z80.memory.writeByte(z80.SP()+1, z80.iyh)
					z80.memory.writeByte(z80.SP(), z80.iyl)
					z80.memory.contendWriteNoMreq(z80.SP(), 1)
					z80.memory.contendWriteNoMreq(z80.SP(), 1)
					z80.iyl = bytetempl
					z80.iyh = bytetemph
					break
				case 0xe5: /* PUSH iy */
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.push16(z80.iyl, z80.iyh)
					break
				case 0xe9: /* JP iy */
					z80.pc = z80.IY() /* NB: NOT INDIRECT! */
					break
				case 0xf9: /* LD SP,iy */
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.memory.contendReadNoMreq(z80.IR(), 1)
					z80.sp = z80.IY()
					break
				default: /* Instruction did not involve H or L, so backtrack
					   one instruction and parse again */
					z80.pc--
					z80.r--
					opcode = opcode2

					goto EndOpcode

				}
			}
			break
		case 0xfe: /* CP nn */
			{
				var bytetemp byte = z80.memory.readByte(z80.PC())
				z80.pc++
				z80.cp(bytetemp)
			}
			break
		case 0xff: /* RST 38 */
			z80.memory.contendReadNoMreq(z80.IR(), 1)
			z80.rst(0x38)
			break

		}
	}
}
