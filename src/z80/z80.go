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

var OpcodesMap [1536]func(z80 *Z80)

const SHIFT_0xCB = 256
const SHIFT_0xED = 512
const SHIFT_0xDD = 768
const SHIFT_0xDDCB = 1024
const SHIFT_0xFDCB = 1024
const SHIFT_0xFD = 1280

type register16 struct {
	high, low *byte
}

func (r register16) inc() {
	temp := r.get() + 1
	*r.high = byte(temp >> 8)
	*r.low = byte(temp & 0xff)
}

func (r register16) dec() {
	temp := r.get() - 1
	*r.high = byte(temp >> 8)
	*r.low = byte(temp & 0xff)
}

func (r register16) set(value uint16) {
	*r.high, *r.low = splitWord(value)
}

func (r register16) get() uint16 {
	return joinBytes(*r.high, *r.low)
}

type Z80 struct {
	A, F, B, C, D, E, H, L         byte
	A_, F_, B_, C_, D_, E_, H_, L_ byte
	IXH, IXL, IYH, IYL             byte
	I, IFF1, IFF2, IM              byte

	// The highest bit (bit 7) of the R register
	R7 byte

	// The low 7 bits of the R register. 16 bits long so it can
	// also act as an RZX instruction counter.
	R uint16

	sp, pc uint16

	bc, bc_, hl, hl_, af, de, de_, ix, iy register16

	EventNextEvent uint

	// Number of tstates since the beginning of the last frame.
	// The value of this variable is usually smaller than TStatesPerFrame,
	// but in some unlikely circumstances it may be >= than that.
	Tstates uint

	Halted bool

	// Needed when executing opcodes prefixed by 0xCB
	tempaddr uint16

	interruptsEnabledAt int

	memory          MemoryAccessor
	ports           PortAccessor

	rzxInstructionsOffset int
}

func NewZ80(memory MemoryAccessor, ports PortAccessor) *Z80 {
	z80 := &Z80{memory: memory, ports: ports}
	z80.bc = register16{&z80.B, &z80.C}
	z80.bc_ = register16{&z80.B_, &z80.C_}
	z80.hl = register16{&z80.H, &z80.L}
	z80.hl_ = register16{&z80.H_, &z80.L_}
	z80.af = register16{&z80.A, &z80.F}
	z80.de = register16{&z80.D, &z80.E}
	z80.ix = register16{&z80.IXH, &z80.IXL}
	z80.iy = register16{&z80.IYH, &z80.IYL}
	z80.de_ = register16{&z80.D_, &z80.E_}
	return z80
}

func (z80 *Z80) Reset() {
	z80.A, z80.F, z80.B, z80.C, z80.D, z80.E, z80.H, z80.L = 0, 0, 0, 0, 0, 0, 0, 0
	z80.A_, z80.F_, z80.B_, z80.C_, z80.D_, z80.E_, z80.H_, z80.L_ = 0, 0, 0, 0, 0, 0, 0, 0
	z80.IXH, z80.IXL, z80.IYH, z80.IYL = 0, 0, 0, 0

	z80.sp, z80.I, z80.R, z80.R7, z80.pc, z80.IFF1, z80.IFF2, z80.IM = 0, 0, 0, 0, 0, 0, 0, 0

	z80.Tstates = 0

	z80.Halted = false
	z80.interruptsEnabledAt = 0
}

func splitWord(word uint16) (byte, byte) {
	return byte(word >> 8), byte(word & 0xff)
}
func joinBytes(h, l byte) uint16 {
	return uint16(l) | (uint16(h) << 8)
}

/* Process a z80 maskable interrupt */
func (z80 *Z80) Interrupt() {
	if z80.IFF1 != 0 {
		if z80.Halted {
			z80.pc++
			z80.Halted = false
		}

		z80.Tstates += 7

		z80.R = (z80.R + 1) & 0x7f
		z80.IFF1, z80.IFF2 = 0, 0

		// push PC
		{
			pch, pcl := splitWord(z80.pc)

			z80.sp--
			z80.memory.WriteByte(z80.sp, pch)
			z80.sp--
			z80.memory.WriteByte(z80.sp, pcl)
		}

		switch z80.IM {
		case 0, 1:
			z80.pc = 0x0038

		case 2:
			var inttemp uint16 = (uint16(z80.I) << 8) | 0xff
			pcl := z80.memory.ReadByte(inttemp)
			inttemp++
			pch := z80.memory.ReadByte(inttemp)
			z80.pc = joinBytes(pch, pcl)

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

func signExtend(v byte) int16 {
	return int16(int8(v))
}

func (z80 *Z80) jp() {
	var jptemp uint16 = z80.pc
	pcl := z80.memory.ReadByte(jptemp)
	jptemp++
	pch := z80.memory.ReadByte(jptemp)
	z80.pc = joinBytes(pch, pcl)
}

func (z80 *Z80) dec(value *byte) {
	z80.F = (z80.F & FLAG_C) | ternOpB((*value&0x0f) != 0, 0, FLAG_H) | FLAG_N
	*value--
	z80.F |= ternOpB(*value == 0x7f, FLAG_V, 0) | sz53Table[*value]
}

func (z80 *Z80) inc(value *byte) {
	*value++
	z80.F = (z80.F & FLAG_C) | ternOpB(*value == 0x80, FLAG_V, 0) | ternOpB((*value&0x0f) != 0, 0, FLAG_H) | sz53Table[(*value)]
}

func (z80 *Z80) jr() {
	var jrtemp int16 = signExtend(z80.memory.ReadByte(z80.pc))
	z80.memory.ContendReadNoMreq_loop(z80.pc, 1, 5)
	z80.pc += uint16(jrtemp)
}

func (z80 *Z80) ld16nnrr(regl, regh byte) {
	var ldtemp uint16

	ldtemp = uint16(z80.memory.ReadByte(z80.pc))
	z80.pc++
	ldtemp |= uint16(z80.memory.ReadByte(z80.pc)) << 8
	z80.pc++
	z80.memory.WriteByte(ldtemp, regl)
	ldtemp++
	z80.memory.WriteByte(ldtemp, regh)
}

func (z80 *Z80) ld16rrnn(regl, regh *byte) {
	var ldtemp uint16

	ldtemp = uint16(z80.memory.ReadByte(z80.pc))
	z80.pc++
	ldtemp |= uint16(z80.memory.ReadByte(z80.pc)) << 8
	z80.pc++
	*regl = z80.memory.ReadByte(ldtemp)
	ldtemp++
	*regh = z80.memory.ReadByte(ldtemp)
}

func (z80 *Z80) sub(value byte) {
	var subtemp uint16 = uint16(z80.A) - uint16(value)
	var lookup byte = ((z80.A & 0x88) >> 3) | ((value & 0x88) >> 2) | byte((subtemp&0x88)>>1)
	z80.A = byte(subtemp)
	z80.F = ternOpB(subtemp&0x100 != 0, FLAG_C, 0) | FLAG_N |
		halfcarrySubTable[lookup&0x07] | overflowSubTable[lookup>>4] |
		sz53Table[z80.A]
}

func (z80 *Z80) and(value byte) {
	z80.A &= value
	z80.F = FLAG_H | sz53pTable[z80.A]
}

func (z80 *Z80) adc(value byte) {
	var adctemp uint16 = uint16(z80.A) + uint16(value) + (uint16(z80.F) & FLAG_C)
	var lookup byte = byte(((uint16(z80.A) & 0x88) >> 3) | ((uint16(value) & 0x88) >> 2) | ((uint16(adctemp) & 0x88) >> 1))

	z80.A = byte(adctemp)

	z80.F = ternOpB((adctemp&0x100) != 0, FLAG_C, 0) | halfcarryAddTable[lookup&0x07] | overflowAddTable[lookup>>4] | sz53Table[z80.A]
}

func (z80 *Z80) adc16(value uint16) {
	var add16temp uint = uint(z80.HL()) + uint(value) + (uint(z80.F) & FLAG_C)
	var lookup byte = byte(((uint(z80.HL()) & 0x8800) >> 11) | ((uint(value) & 0x8800) >> 10) | (add16temp&0x8800)>>9)

	z80.SetHL(uint16(add16temp))

	z80.F = ternOpB((uint(add16temp)&0x10000) != 0, FLAG_C, 0) | overflowAddTable[lookup>>4] | (z80.H & (FLAG_3 | FLAG_5 | FLAG_S)) | halfcarryAddTable[lookup&0x07] | ternOpB(z80.HL() != 0, 0, FLAG_Z)
}

func (z80 *Z80) add16(value1 register16, value2 uint16) {
	var add16temp uint = uint(value1.get()) + uint(value2)
	var lookup byte = byte(((value1.get() & 0x0800) >> 11) | ((value2 & 0x0800) >> 10) | (uint16(add16temp)&0x0800)>>9)

	value1.set(uint16(add16temp))

	z80.F = (z80.F & (FLAG_V | FLAG_Z | FLAG_S)) | ternOpB((add16temp&0x10000) != 0, FLAG_C, 0) | (byte(add16temp>>8) & (FLAG_3 | FLAG_5)) | halfcarryAddTable[lookup]
}

func (z80 *Z80) add(value byte) {
	var addtemp uint = uint(z80.A) + uint(value)
	var lookup byte = ((z80.A & 0x88) >> 3) | ((value & 0x88) >> 2) | byte((addtemp&0x88)>>1)
	z80.A = byte(addtemp)
	z80.F = ternOpB(addtemp&0x100 != 0, FLAG_C, 0) | halfcarryAddTable[lookup&0x07] | overflowAddTable[lookup>>4] | sz53Table[z80.A]
}

func (z80 *Z80) or(value byte) {
	z80.A |= value
	z80.F = sz53pTable[z80.A]
}

func (z80 *Z80) pop16() (regl, regh byte) {
	regl = z80.memory.ReadByte(z80.sp)
	z80.sp++
	regh = z80.memory.ReadByte(z80.sp)
	z80.sp++
	return
}

func (z80 *Z80) push16(regl, regh byte) {
	z80.sp--
	z80.memory.WriteByte(z80.sp, regh)
	z80.sp--
	z80.memory.WriteByte(z80.sp, regl)
}

func (z80 *Z80) ret() {
	pcl, pch := z80.pop16()
	z80.pc = joinBytes(pch, pcl)
}

func (z80 *Z80) rl(value byte) byte {
	rltemp := value
	value = (value << 1) | (z80.F & FLAG_C)
	z80.F = (rltemp >> 7) | sz53pTable[value]
	return value
}

func (z80 *Z80) rlc(value byte) byte {
	value = (value << 1) | (value >> 7)
	z80.F = (value & FLAG_C) | sz53pTable[value]
	return value
}

func (z80 *Z80) rr(value byte) byte {
	rrtemp := value
	value = (value >> 1) | (z80.F << 7)
	z80.F = (rrtemp & FLAG_C) | sz53pTable[value]
	return value
}

func (z80 *Z80) rrc(value byte) byte {
	z80.F = value & FLAG_C
	value = (value >> 1) | (value << 7)
	z80.F |= sz53pTable[value]
	return value
}

func (z80 *Z80) rst(value byte) {
	pch, pcl := splitWord(z80.pc)
	z80.push16(pcl, pch)
	z80.pc = uint16(value)
}

func (z80 *Z80) sbc(value byte) {
	var sbctemp uint16 = uint16(z80.A) - uint16(value) - (uint16(z80.F) & FLAG_C)
	var lookup byte = ((z80.A & 0x88) >> 3) | ((value & 0x88) >> 2) | byte((sbctemp&0x88)>>1)
	z80.A = byte(sbctemp)
	z80.F = ternOpB((sbctemp&0x100) != 0, FLAG_C, 0) | FLAG_N | halfcarrySubTable[lookup&0x07] | overflowSubTable[lookup>>4] | sz53Table[z80.A]
}

func (z80 *Z80) sbc16(value uint16) {
	var sub16temp uint = uint(z80.HL()) - uint(value) - (uint(z80.F) & FLAG_C)
	var lookup byte = byte(((z80.HL() & 0x8800) >> 11) | ((uint16(value) & 0x8800) >> 10) | ((uint16(sub16temp) & 0x8800) >> 9))

	z80.SetHL(uint16(sub16temp))

	z80.F = ternOpB((sub16temp&0x10000) != 0, FLAG_C, 0) | FLAG_N | overflowSubTable[lookup>>4] | (z80.H & (FLAG_3 | FLAG_5 | FLAG_S)) | halfcarrySubTable[lookup&0x07] | ternOpB(z80.HL() != 0, 0, FLAG_Z)
}

func (z80 *Z80) sla(value byte) byte {
	z80.F = value >> 7
	value <<= 1
	z80.F |= sz53pTable[value]
	return value
}

func (z80 *Z80) sll(value byte) byte {
	z80.F = value >> 7
	value = (value << 1) | 0x01
	z80.F |= sz53pTable[(value)]
	return value
}

func (z80 *Z80) sra(value byte) byte {
	z80.F = value & FLAG_C
	value = (value & 0x80) | (value >> 1)
	z80.F |= sz53pTable[value]
	return value
}

func (z80 *Z80) srl(value byte) byte {
	z80.F = value & FLAG_C
	value >>= 1
	z80.F |= sz53pTable[value]
	return value
}

func (z80 *Z80) xor(value byte) {
	z80.A ^= value
	z80.F = sz53pTable[z80.A]
}

func (z80 *Z80) bit(bit, value byte) {
	z80.F = (z80.F & FLAG_C) | FLAG_H | (value & (FLAG_3 | FLAG_5))
	if value&(0x01<<bit) == 0 {
		z80.F |= FLAG_P | FLAG_Z
	}
	if bit == 7 && (value&0x80) != 0 {
		z80.F |= FLAG_S
	}
}

func (z80 *Z80) biti(bit, value byte, address uint16) {
	z80.F = (z80.F & FLAG_C) | FLAG_H | (byte(address>>8) & (FLAG_3 | FLAG_5))
	if value&(0x01<<bit) == 0 {
		z80.F |= FLAG_P | FLAG_Z
	}
	if (bit == 7) && (value&0x80) != 0 {
		z80.F |= FLAG_S
	}
}

func (z80 *Z80) call() {
	var calltempl, calltemph byte
	calltempl = z80.memory.ReadByte(z80.pc)
	z80.pc++
	calltemph = z80.memory.ReadByte(z80.pc)
	z80.memory.ContendReadNoMreq(z80.pc, 1)
	z80.pc++
	pch, pcl := splitWord(z80.pc)
	z80.push16(pcl, pch)
	z80.pc = joinBytes(calltemph, calltempl)
}

func (z80 *Z80) cp(value byte) {
	var cptemp uint16 = uint16(z80.A) - uint16(value)
	var lookup byte = ((z80.A & 0x88) >> 3) | ((value & 0x88) >> 2) | byte((cptemp&0x88)>>1)
	z80.F = ternOpB((cptemp&0x100) != 0, FLAG_C, ternOpB(cptemp != 0, 0, FLAG_Z)) | FLAG_N | halfcarrySubTable[lookup&0x07] | overflowSubTable[lookup>>4] | (value & (FLAG_3 | FLAG_5)) | byte(cptemp&FLAG_S)
}

func (z80 *Z80) in(reg *byte, port uint16) {
	*reg = z80.readPort(port)
	z80.F = (z80.F & FLAG_C) | sz53pTable[*reg]
}

func (z80 *Z80) readPort(address uint16) byte {
	return z80.ports.ReadPort(address)
}

func (z80 *Z80) writePort(address uint16, b byte) {
	z80.ports.WritePort(address, b)
}

// The following functions can not be generated as they need special treatments

func (z80 *Z80) PC() uint16 {
	return z80.pc
}

func (z80 *Z80) SetPC(value uint16) {
	z80.pc = value
}

func (z80 *Z80) IncPC(value uint16) {
	z80.pc += value
}

func (z80 *Z80) DecPC(value uint16) {
	z80.pc -= value
}

func (z80 *Z80) SP() uint16 {
	return z80.sp
}

func (z80 *Z80) SetSP(value uint16) {
	z80.sp = value
}

func (z80 *Z80) IncSP() {
	z80.sp++
}

func (z80 *Z80) DecSP() {
	z80.sp--
}

func (z80 *Z80) IR() uint16 {
	var ir uint16
	ir |= uint16(z80.I) << 8
	ir |= uint16(z80.R7&0x80) | (z80.R & 0x7f)
	return ir
}

func (z80 *Z80) sltTrap(address int16, level byte) int {
	// Dummy implementation
	return 0
}

func (z80 *Z80) doOpcodes() {
	// Main instruction emulation loop
	{
		for (z80.Tstates < z80.EventNextEvent) && !z80.Halted {
			z80.memory.ContendRead(z80.pc, 4)
			opcode := z80.memory.ReadByteInternal(z80.pc)
			z80.R = (z80.R + 1) & 0x7f
			z80.pc++
			OpcodesMap[opcode](z80)

		}

		if z80.Halted {
			// Repeat emulating the HALT instruction until 'z80.eventNextEvent'
			for z80.Tstates < z80.EventNextEvent {
				z80.memory.ContendRead(z80.pc, 4)
				z80.R = (z80.R + 1) & 0x7f
			}
		}
	}
}

func invalidOpcode(z80 *Z80) {
	panic("invalid opcode")
}

func opcode_cb(z80 *Z80) {
	z80.memory.ContendRead(z80.pc, 4)
	var opcode2 byte = z80.memory.ReadByteInternal(z80.pc)
	z80.pc++
	z80.R++
	OpcodesMap[SHIFT_0xCB+int(opcode2)](z80)
}

func opcode_ed(z80 *Z80) {
	z80.memory.ContendRead(z80.pc, 4)
	var opcode2 byte = z80.memory.ReadByteInternal(z80.pc)
	z80.pc++
	z80.R++

	if f := OpcodesMap[SHIFT_0xED+int(opcode2)]; f != nil {
		f(z80)
	} else {
		invalidOpcode(z80)
	}
}

func opcode_dd(z80 *Z80) {
	z80.memory.ContendRead(z80.pc, 4)
	var opcode2 byte = z80.memory.ReadByteInternal(z80.pc)
	z80.pc++
	z80.R++

	switch opcode2 {
	case 0xcb:
		z80.memory.ContendRead(z80.pc, 3)
		z80.tempaddr = z80.IX() + uint16(signExtend(z80.memory.ReadByteInternal(z80.pc)))
		z80.pc++
		z80.memory.ContendRead(z80.pc, 3)
		var opcode3 byte = z80.memory.ReadByteInternal(z80.pc)
		z80.memory.ContendReadNoMreq_loop(z80.pc, 1, 2)
		z80.pc++
		OpcodesMap[SHIFT_0xDDCB+int(opcode3)](z80)
	default:
		if f := OpcodesMap[SHIFT_0xDD+int(opcode2)]; f != nil {
			f(z80)
		} else {
			/* Instruction did not involve H or L */
			OpcodesMap[opcode2](z80)
		}
	}
}

func opcode_fd(z80 *Z80) {
	z80.memory.ContendRead(z80.pc, 4)
	var opcode2 byte = z80.memory.ReadByteInternal(z80.pc)
	z80.pc++
	z80.R++

	switch opcode2 {
	case 0xcb:
		z80.memory.ContendRead(z80.pc, 3)
		z80.tempaddr = z80.IY() + uint16(signExtend(z80.memory.ReadByteInternal(z80.pc)))
		z80.pc++
		z80.memory.ContendRead(z80.pc, 3)
		var opcode3 byte = z80.memory.ReadByteInternal(z80.pc)
		z80.memory.ContendReadNoMreq_loop(z80.pc, 1, 2)
		z80.pc++

		OpcodesMap[SHIFT_0xFDCB+int(opcode3)](z80)

	default:
		if f := OpcodesMap[SHIFT_0xFD+int(opcode2)]; f != nil {
			f(z80)
		} else {
			/* Instruction did not involve H or L */
			OpcodesMap[opcode2](z80)
		}
	}
}

func init() {
	initOpcodes()
	OpcodesMap[0xcb] = opcode_cb
	OpcodesMap[0xdd] = opcode_dd
	OpcodesMap[0xed] = opcode_ed
	OpcodesMap[0xfd] = opcode_fd
}
