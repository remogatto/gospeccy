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
	"fmt"
	"⚛perf"
	"os"
	"syscall"
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


var opcodesMap [1536]func(z80 *Z80, tempaddr uint16)

const SHIFT_0xCB = 256
const SHIFT_0xED = 512
const SHIFT_0xDD = 768
const SHIFT_0xDDCB = 1024
const SHIFT_0xFDCB = 1024
const SHIFT_0xFD = 1280

func shift0xcb(opcode byte) int {
	return SHIFT_0xCB + int(opcode)
}

func shift0xed(opcode byte) int {
	return SHIFT_0xED + int(opcode)
}

func shift0xdd(opcode byte) int {
	return SHIFT_0xDD + int(opcode)
}

func shift0xddcb(opcode byte) int {
	return SHIFT_0xDDCB + int(opcode)
}

func shift0xfdcb(opcode byte) int {
	return SHIFT_0xFDCB + int(opcode)
}

func shift0xfd(opcode byte) int {
	return SHIFT_0xFD + int(opcode)
}


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

	// Number of tstates since the beginning of the last frame.
	// The value of this variable is usually smaller than TStatesPerFrame,
	// but in some unlikely circumstances it may be >= than that.
	tstates uint

	halted bool

	interruptsEnabledAt int

	memory MemoryAccessor

	ports PortAccessor

	speccy *Spectrum48k

	rzxInstructionsOffset int

	eventNextEvent uint

	LogEvents bool

	z80_instructionCounter     uint64 // Number of Z80 instructions executed
	z80_instructionsMeasured   uint64 // Number of Z80 instrs that can be related to 'hostCpu_instructionCounter'
	hostCpu_instructionCounter uint64
	perfCounter_hostCpuInstr   *perf.PerfCounter // Can be nil (if creating the counter fails)
}

func NewZ80(memory MemoryAccessor, ports PortAccessor) *Z80 {
	z80 := &Z80{memory: memory, ports: ports}

	z80.bc = register16{&z80.b, &z80.c}
	z80.bc_ = register16{&z80.b_, &z80.c_}
	z80.hl = register16{&z80.h, &z80.l}
	z80.hl_ = register16{&z80.h_, &z80.l_}
	z80.af = register16{&z80.a, &z80.f}
	z80.de = register16{&z80.d, &z80.e}
	z80.ix = register16{&z80.ixh, &z80.ixl}
	z80.iy = register16{&z80.iyh, &z80.iyl}
	z80.de_ = register16{&z80.d_, &z80.e_}

	z80.perfCounter_hostCpuInstr = perf.NewCounter_Instructions( /*user*/ true, /*kernel*/ false)

	return z80
}

func (z80 *Z80) init(speccy *Spectrum48k) {
	z80.speccy = speccy
}

func (z80 *Z80) close() {
	if z80.perfCounter_hostCpuInstr != nil {
		z80.perfCounter_hostCpuInstr.Close()
		z80.perfCounter_hostCpuInstr = nil
	}
}

// Returns the average number of host-CPU instructions required to execute one Z80 instruction.
// Returns zero if this information is not available.
func (z80 *Z80) GetEmulationEfficiency() uint {
	var eff uint
	if z80.z80_instructionsMeasured > 0 {
		eff = uint(z80.hostCpu_instructionCounter / z80.z80_instructionsMeasured)
	} else {
		eff = 0
	}
	return eff
}

func (z80 *Z80) reset() {
	z80.a, z80.f, z80.b, z80.c, z80.d, z80.e, z80.h, z80.l = 0, 0, 0, 0, 0, 0, 0, 0
	z80.a_, z80.f_, z80.b_, z80.c_, z80.d_, z80.e_, z80.h_, z80.l_ = 0, 0, 0, 0, 0, 0, 0, 0
	z80.ixh, z80.ixl, z80.iyh, z80.iyl = 0, 0, 0, 0

	z80.sp, z80.i, z80.r, z80.r7, z80.pc, z80.iff1, z80.iff2, z80.im = 0, 0, 0, 0, 0, 0, 0, 0

	z80.tstates = 0

	z80.halted = false
	z80.interruptsEnabledAt = 0
}


// Initializes state from data in SNA format.
// Returns nil on success.
func (z80 *Z80) loadSna(data []byte) os.Error {
	if len(data) != 49179 {
		return os.NewError(fmt.Sprintf("snapshot has invalid size"))
	}

	// Populate registers
	z80.i = data[0]
	z80.l_ = data[1]
	z80.h_ = data[2]
	z80.e_ = data[3]
	z80.d_ = data[4]
	z80.c_ = data[5]
	z80.b_ = data[6]
	z80.f_ = data[7]
	z80.a_ = data[8]
	z80.l = data[9]
	z80.h = data[10]
	z80.e = data[11]
	z80.d = data[12]
	z80.c = data[13]
	z80.b = data[14]
	z80.iyl = data[15]
	z80.iyh = data[16]
	z80.ixl = data[17]
	z80.ixh = data[18]

	z80.iff1 = uint16(ternOpB((data[19]&0x04) != 0, 1, 0))
	z80.iff2 = z80.iff1

	var r = uint16(data[20])

	z80.r = r & 0x7f
	z80.r7 = r & 0x80

	z80.f = data[21]
	z80.a = data[22]
	z80.sp = uint16(data[23]) | (uint16(data[24]) << 8)
	z80.im = uint16(data[25])

	// Border color
	z80.writePort(0xfe, data[26]&0x07)

	// Populate memory
	for i := 0; i < 0xc000; i++ {
		z80.memory.Data()[i+0x4000] = data[i+27]
	}

	// Send a RETN
	z80.iff1 = z80.iff2
	z80.ret()

	if z80.iff1 != 0 {
		z80.tstates = InterruptLength
	} else {
		z80.tstates = 0
	}

	return nil
}

func (z80 *Z80) saveSna() ([]byte, os.Error) {

	data := make([]byte, 49179)

	// Save registers
	data[0] = z80.i
	data[1] = z80.l_
	data[2] = z80.h_
	data[3] = z80.e_
	data[4] = z80.d_
	data[5] = z80.c_
	data[6] = z80.b_
	data[7] = z80.f_
	data[8] = z80.a_
	data[9] = z80.l
	data[10] = z80.h
	data[11] = z80.e
	data[12] = z80.d
	data[13] = z80.c
	data[14] = z80.b
	data[15] = z80.iyl
	data[16] = z80.iyh
	data[17] = z80.ixl
	data[18] = z80.ixh

	if z80.iff1 != 0 {
		data[19] = (1 << 2)
	} else {
		data[19] = (0 << 2)
	}

	data[20] = byte(z80.r&0x7f) | byte(z80.r7&0x80)

	data[21] = z80.f
	data[22] = z80.a

	sp_afterSimulatedPushPC := z80.sp - 2
	if (sp_afterSimulatedPushPC < 0x4000) || (sp_afterSimulatedPushPC > 0xfffe) {
		// We would be saving the PC to ROM or outside of memory
		return nil, os.NewError(fmt.Sprintf("failed to simulate a RETN"))
	}

	data[23] = byte(sp_afterSimulatedPushPC & 0xff)
	data[24] = byte(sp_afterSimulatedPushPC >> 8)
	data[25] = byte(z80.im)

	// Border color
	data[26] = z80.speccy.ula.getBorderColor() & 0x07

	// Memory
	for i := 0; i < 0xc000; i++ {
		data[i+27] = z80.memory.Data()[i+0x4000]
	}

	// Push PC
	pch, pcl := splitWord(z80.pc)
	data[(sp_afterSimulatedPushPC-0x4000+0)+27] = pcl
	data[(sp_afterSimulatedPushPC-0x4000+1)+27] = pch

	return data, nil
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

		z80.tstates += InterruptLength
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
	z80.f |= ternOpB(*value == 0x7f, FLAG_V, 0) | sz53Table[*value]
}

func (z80 *Z80) inc(value *byte) {
	*value++
	z80.f = (z80.f & FLAG_C) | ternOpB(*value == 0x80, FLAG_V, 0) | ternOpB((*value&0x0f) != 0, 0, FLAG_H) | sz53Table[(*value)]
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
		sz53Table[z80.a]
}

func (z80 *Z80) and(value byte) {
	z80.a &= value
	z80.f = FLAG_H | sz53pTable[z80.a]
}

func (z80 *Z80) adc(value byte) {
	var adctemp uint16 = uint16(z80.a) + uint16(value) + (uint16(z80.f) & FLAG_C)
	var lookup byte = byte(((uint16(z80.a) & 0x88) >> 3) | ((uint16(value) & 0x88) >> 2) | ((uint16(adctemp) & 0x88) >> 1))

	z80.a = byte(adctemp)

	z80.f = ternOpB((adctemp&0x100) != 0, FLAG_C, 0) | halfcarryAddTable[lookup&0x07] | overflowAddTable[lookup>>4] | sz53Table[z80.a]
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
	z80.f = ternOpB(addtemp&0x100 != 0, FLAG_C, 0) | halfcarryAddTable[lookup&0x07] | overflowAddTable[lookup>>4] | sz53Table[z80.a]
}

func (z80 *Z80) or(value byte) {
	z80.a |= value
	z80.f = sz53pTable[z80.a]
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
	z80.f = (rltemp >> 7) | sz53pTable[*value]
}

func (z80 *Z80) rlc(value *byte) {
	*value = (*value << 1) | (*value >> 7)
	z80.f = (*value & FLAG_C) | sz53pTable[*value]
}

func (z80 *Z80) rr(value *byte) {
	rrtemp := *value
	*value = (*value >> 1) | (z80.f << 7)
	z80.f = (rrtemp & FLAG_C) | sz53pTable[*value]
}

func (z80 *Z80) rrc(value *byte) {
	z80.f = *value & FLAG_C
	*value = (*value >> 1) | (*value << 7)
	z80.f |= sz53pTable[*value]
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
	z80.f = ternOpB((sbctemp&0x100) != 0, FLAG_C, 0) | FLAG_N | halfcarrySubTable[lookup&0x07] | overflowSubTable[lookup>>4] | sz53Table[z80.a]
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
	z80.f |= sz53pTable[*value]
}

func (z80 *Z80) sll(value *byte) {
	z80.f = *value >> 7
	*value = (*value << 1) | 0x01
	z80.f |= sz53pTable[(*value)]
}

func (z80 *Z80) sra(value *byte) {
	z80.f = *value & FLAG_C
	*value = (*value & 0x80) | (*value >> 1)
	z80.f |= sz53pTable[*value]
}

func (z80 *Z80) srl(value *byte) {
	z80.f = *value & FLAG_C
	*value >>= 1
	z80.f |= sz53pTable[*value]
}

func (z80 *Z80) xor(value byte) {
	z80.a ^= value
	z80.f = sz53pTable[z80.a]
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
	z80.f = (z80.f & FLAG_C) | sz53pTable[*reg]
}

func (z80 *Z80) readPort(address uint16) byte {
	return z80.ports.readPort(address)
}

func (z80 *Z80) writePort(address uint16, b byte) {
	z80.ports.writePort(address, b)
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
	ttid_start := syscall.Gettid()

	var hostCpu_instrCount_start uint64
	if z80.perfCounter_hostCpuInstr != nil {
		hostCpu_instrCount_start, _ = z80.perfCounter_hostCpuInstr.Read()
	} else {
		hostCpu_instrCount_start = 0
	}

	var z80_localInstructionCounter uint = 0

	for (z80.tstates < z80.eventNextEvent) && !z80.halted {
		z80.memory.contendRead(z80.pc, 4)
		opcode := z80.memory.readByteInternal(z80.pc)

		z80.r = (z80.r + 1) & 0x7f
		z80.pc++

		z80_localInstructionCounter++

		opcodesMap[opcode](z80, 0)
	}

	// Update emulation efficiency counters
	{
		ttid_end := syscall.Gettid()

		var hostCpu_instrCount_end uint64
		if z80.perfCounter_hostCpuInstr != nil {
			hostCpu_instrCount_end, _ = z80.perfCounter_hostCpuInstr.Read()
		} else {
			hostCpu_instrCount_end = 0
		}

		z80.z80_instructionCounter += uint64(z80_localInstructionCounter)

		/*if z80_localInstructionCounter > 0 {
			println( z80_localInstructionCounter, hostCpu_instrCount_start, hostCpu_instrCount_end,
					hostCpu_instrCount_end-hostCpu_instrCount_start,
					(hostCpu_instrCount_end - hostCpu_instrCount_start) / uint64(z80_localInstructionCounter) )
		}*/

		if (ttid_start == ttid_end) &&
			(z80_localInstructionCounter > 0) &&
			(hostCpu_instrCount_start > 0) &&
			(hostCpu_instrCount_end > 0) &&
			(hostCpu_instrCount_end > hostCpu_instrCount_start) {

			avg := uint((hostCpu_instrCount_end - hostCpu_instrCount_start) / uint64(z80_localInstructionCounter))

			// It may happen that the measured values are invalid.
			// The primary cause of this is that the Go runtime
			// can move a goroutine to a different OS thread,
			// without notifying us when it does so.
			// The majority of these cases is detected by (ttid_start == ttid_end) constraint.
			eff := z80.GetEmulationEfficiency()
			bogusMeasurement := (avg < eff/4) || ((eff > 0) && (avg > eff*4))

			if !bogusMeasurement {
				z80.z80_instructionsMeasured += uint64(z80_localInstructionCounter)
				z80.hostCpu_instructionCounter += (hostCpu_instrCount_end - hostCpu_instrCount_start)
			}
		}
	}
}


func invalidOpcode(z80 *Z80, tempaddr uint16) {
	panic("invalid opcode")
}

func opcode_cb(z80 *Z80, tempaddr uint16) {
	var opcode2 byte
	z80.memory.contendRead(z80.pc, 4)
	opcode2 = z80.memory.readByteInternal(z80.pc)
	z80.pc++
	z80.r++
	opcodesMap[SHIFT_0xCB+int(opcode2)](z80, 0)
}

func opcode_ed(z80 *Z80, tempaddr uint16) {
	var opcode2 byte
	z80.memory.contendRead(z80.pc, 4)
	opcode2 = z80.memory.readByteInternal(z80.pc)
	z80.pc++
	z80.r++

	if f := opcodesMap[SHIFT_0xED+int(opcode2)]; f != nil {
		f(z80, 0)
	} else {
		invalidOpcode(z80, 0)
	}
}

func opcode_dd(z80 *Z80, tempaddr uint16) {
	var opcode2 byte
	z80.memory.contendRead(z80.pc, 4)
	opcode2 = z80.memory.readByteInternal(z80.pc)
	z80.pc++
	z80.r++

	switch opcode2 {
	case 0xcb:
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
		opcodesMap[SHIFT_0xDDCB+int(opcode3)](z80, tempaddr)
	default:
		if f := opcodesMap[SHIFT_0xDD+int(opcode2)]; f != nil {
			f(z80, 0)
		} else {
			/* Instruction did not involve H or L */
			opcodesMap[opcode2](z80, 0)
		}
	}
}

func opcode_fd(z80 *Z80, tempaddr uint16) {
	var opcode2 byte
	z80.memory.contendRead(z80.pc, 4)
	opcode2 = z80.memory.readByteInternal(z80.pc)
	z80.pc++
	z80.r++

	switch opcode2 {
	case 0xcb:
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

		opcodesMap[SHIFT_0xFDCB+int(opcode3)](z80, tempaddr)

	default:
		if f := opcodesMap[SHIFT_0xFD+int(opcode2)]; f != nil {
			f(z80, 0)
		} else {
			/* Instruction did not involve H or L */
			opcodesMap[opcode2](z80, 0)
		}
	}
}


func init() {
	initOpcodes()
	opcodesMap[0xcb] = opcode_cb
	opcodesMap[0xdd] = opcode_dd
	opcodesMap[0xed] = opcode_ed
	opcodesMap[0xfd] = opcode_fd
}
