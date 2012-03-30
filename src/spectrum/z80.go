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
	"github.com/0xe2-0x9a-0x9b/Go-PerfEvents"
	"github.com/remogatto/gospeccy/src/formats"
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

var opcodesMap [1536]func(z80 *Z80)

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
	a, f, b, c, d, e, h, l         byte
	a_, f_, b_, c_, d_, e_, h_, l_ byte
	ixh, ixl, iyh, iyl             byte
	i, iff1, iff2, im              byte

	// The highest bit (bit 7) of the R register
	r7 byte

	// The low 7 bits of the R register. 16 bits long so it can also act as an RZX instruction counter.
	r uint16

	sp, pc uint16

	bc, bc_, hl, hl_, af, de, de_, ix, iy register16

	// Needed when executing opcodes prefixed by 0xCB
	tempaddr uint16

	// Number of tstates since the beginning of the last frame.
	// The value of this variable is usually smaller than TStatesPerFrame,
	// but in some unlikely circumstances it may be >= than that.
	tstates uint

	halted bool

	interruptsEnabledAt int

	memory          MemoryAccessor
	ports           PortAccessor
	ula             *ULA
	tapeDrive_orNil *TapeDrive

	rzxInstructionsOffset int

	eventNextEvent uint

	z80_instructionCounter     uint64 // Number of Z80 instructions executed
	z80_instructionsMeasured   uint64 // Number of Z80 instrs that can be related to 'hostCpu_instructionCounter'
	hostCpu_instructionCounter uint64
	perfCounter_hostCpuInstr   *perf.Counter // Can be nil (if creating the counter fails)

	readFromTape bool

	// The value is non-zero if a couple of the most recent frames
	// executed instructions which appeared to be reading from the tape
	shouldPlayTheTape int
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

	z80.perfCounter_hostCpuInstr, _ = perf.NewCounter_Instructions(true /*user*/, false /*kernel*/)

	return z80
}

func (z80 *Z80) init(ula *ULA, tapeDrive_orNil *TapeDrive) {
	z80.ula = ula
	z80.tapeDrive_orNil = tapeDrive_orNil
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

// Initializes state from the specified snapshot.
// Returns nil on success.
func (z80 *Z80) loadSnapshot(s formats.Snapshot) error {
	cpu := s.CpuState()
	ula := s.UlaState()
	mem := s.Memory()

	// Populate registers
	z80.a = cpu.A
	z80.f = cpu.F
	z80.b = cpu.B
	z80.c = cpu.C
	z80.d = cpu.D
	z80.e = cpu.E
	z80.h = cpu.H
	z80.l = cpu.L
	z80.a_ = cpu.A_
	z80.f_ = cpu.F_
	z80.b_ = cpu.B_
	z80.c_ = cpu.C_
	z80.d_ = cpu.D_
	z80.e_ = cpu.E_
	z80.h_ = cpu.H_
	z80.l_ = cpu.L_
	z80.ixl = byte(cpu.IX & 0xff)
	z80.ixh = byte(cpu.IX >> 8)
	z80.iyl = byte(cpu.IY & 0xff)
	z80.iyh = byte(cpu.IY >> 8)

	z80.i = cpu.I
	z80.iff1 = cpu.IFF1
	z80.iff2 = cpu.IFF2
	z80.im = cpu.IM

	z80.r = uint16(cpu.R & 0x7f)
	z80.r7 = cpu.R & 0x80

	z80.pc = cpu.PC
	z80.sp = cpu.SP

	// Border color
	z80.ports.writePortInternal(0xfe, ula.Border&0x07, false /*contend*/)

	// Populate memory
	copy(z80.memory.Data()[0x4000:], mem[:])

	z80.tstates = cpu.Tstate

	return nil
}

func (z80 *Z80) MakeSnapshot() *formats.FullSnapshot {
	var s formats.FullSnapshot

	// Save registers
	s.Cpu.A = z80.a
	s.Cpu.F = z80.f
	s.Cpu.B = z80.b
	s.Cpu.C = z80.c
	s.Cpu.D = z80.d
	s.Cpu.E = z80.e
	s.Cpu.H = z80.h
	s.Cpu.L = z80.l
	s.Cpu.A_ = z80.a_
	s.Cpu.F_ = z80.f_
	s.Cpu.B_ = z80.b_
	s.Cpu.C_ = z80.c_
	s.Cpu.D_ = z80.d_
	s.Cpu.E_ = z80.e_
	s.Cpu.H_ = z80.h_
	s.Cpu.L_ = z80.l_
	s.Cpu.IX = uint16(z80.ixl) | (uint16(z80.ixh) << 8)
	s.Cpu.IY = uint16(z80.iyl) | (uint16(z80.iyh) << 8)

	s.Cpu.I = z80.i
	s.Cpu.IFF1 = z80.iff1
	s.Cpu.IFF2 = z80.iff2
	s.Cpu.IM = z80.im

	s.Cpu.R = byte(z80.r&0x7f) | (z80.r7 & 0x80)

	s.Cpu.SP = z80.sp
	s.Cpu.PC = z80.pc

	// Border color
	s.Ula.Border = z80.ula.getBorderColor() & 0x07

	// Memory
	copy(s.Mem[:], z80.memory.Data()[0x4000:])

	return &s
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

		z80.tstates += 7

		z80.r = (z80.r + 1) & 0x7f
		z80.iff1, z80.iff2 = 0, 0

		// push PC
		{
			pch, pcl := splitWord(z80.pc)

			z80.sp--
			z80.memory.writeByte(z80.sp, pch)
			z80.sp--
			z80.memory.writeByte(z80.sp, pcl)
		}

		switch z80.im {
		case 0, 1:
			z80.pc = 0x0038

		case 2:
			var inttemp uint16 = (uint16(z80.i) << 8) | 0xff
			pcl := z80.memory.readByte(inttemp)
			inttemp++
			pch := z80.memory.readByte(inttemp)
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

func (z80 *Z80) tapeSaveTrap() int {
	// Dummy implementation
	return 0
}

func (z80 *Z80) tapeLoadTrap() int {
	// Dummy implementation
	return 0
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
	var jrtemp int16 = signExtend(z80.memory.readByte(z80.pc))
	z80.memory.contendReadNoMreq_loop(z80.pc, 1, 5)
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

func (z80 *Z80) pop16() (regl, regh byte) {
	regl = z80.memory.readByte(z80.sp)
	z80.sp++
	regh = z80.memory.readByte(z80.sp)
	z80.sp++
	return
}

func (z80 *Z80) push16(regl, regh byte) {
	z80.sp--
	z80.memory.writeByte(z80.sp, regh)
	z80.sp--
	z80.memory.writeByte(z80.sp, regl)
}

func (z80 *Z80) ret() {
	pcl, pch := z80.pop16()
	z80.pc = joinBytes(pch, pcl)
}

func (z80 *Z80) rl(value byte) byte {
	rltemp := value
	value = (value << 1) | (z80.f & FLAG_C)
	z80.f = (rltemp >> 7) | sz53pTable[value]
	return value
}

func (z80 *Z80) rlc(value byte) byte {
	value = (value << 1) | (value >> 7)
	z80.f = (value & FLAG_C) | sz53pTable[value]
	return value
}

func (z80 *Z80) rr(value byte) byte {
	rrtemp := value
	value = (value >> 1) | (z80.f << 7)
	z80.f = (rrtemp & FLAG_C) | sz53pTable[value]
	return value
}

func (z80 *Z80) rrc(value byte) byte {
	z80.f = value & FLAG_C
	value = (value >> 1) | (value << 7)
	z80.f |= sz53pTable[value]
	return value
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

func (z80 *Z80) sla(value byte) byte {
	z80.f = value >> 7
	value <<= 1
	z80.f |= sz53pTable[value]
	return value
}

func (z80 *Z80) sll(value byte) byte {
	z80.f = value >> 7
	value = (value << 1) | 0x01
	z80.f |= sz53pTable[(value)]
	return value
}

func (z80 *Z80) sra(value byte) byte {
	z80.f = value & FLAG_C
	value = (value & 0x80) | (value >> 1)
	z80.f |= sz53pTable[value]
	return value
}

func (z80 *Z80) srl(value byte) byte {
	z80.f = value & FLAG_C
	value >>= 1
	z80.f |= sz53pTable[value]
	return value
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
	var ir uint16
	ir |= uint16(z80.i) << 8
	ir |= uint16(z80.r7&0x80) | (z80.r & 0x7f)
	return ir
}

func (z80 *Z80) sltTrap(address int16, level byte) int {
	// Dummy implementation
	return 0
}

func (z80 *Z80) doOpcodes() {
	var ttid_start int
	if z80.perfCounter_hostCpuInstr != nil {
		ttid_start = z80.perfCounter_hostCpuInstr.Gettid()
	} else {
		ttid_start = -1
	}

	var hostCpu_instrCount_start uint64 = 0
	var hostCpu_instrCount_startErr error = nil
	if z80.perfCounter_hostCpuInstr != nil {
		hostCpu_instrCount_start, hostCpu_instrCount_startErr = z80.perfCounter_hostCpuInstr.Read()
	}

	var z80_localInstructionCounter uint = 0

	// Main instruction emulation loop
	{
		var readFromTape bool = (z80.readFromTape && (z80.shouldPlayTheTape > 0) && (z80.tapeDrive_orNil != nil))

		if z80.tapeDrive_orNil != nil && z80.tapeDrive_orNil.NotifyLoadComplete && z80.tapeDrive_orNil.notifyCpuLoadCompleted {
			z80.tapeDrive_orNil.notifyCpuLoadCompleted = false
			z80.tapeDrive_orNil.loadComplete <- true
		}

		if !readFromTape {
			if z80.tapeDrive_orNil != nil {
				z80.tapeDrive_orNil.decelerate()
			}
		}

		for (z80.tstates < z80.eventNextEvent) && !z80.halted {
			z80.memory.contendRead(z80.pc, 4)
			opcode := z80.memory.readByteInternal(z80.pc)

			z80.r = (z80.r + 1) & 0x7f
			z80.pc++

			z80_localInstructionCounter++

			opcodesMap[opcode](z80)

			if readFromTape {
				endOfBlock := z80.tapeDrive_orNil.doPlay()
				if endOfBlock {
					readFromTape = false
					z80.shouldPlayTheTape = 0
					z80.tapeDrive_orNil.decelerate()
				}
			}
		}

		if z80.halted {
			z80.shouldPlayTheTape = 0
			if z80.tapeDrive_orNil != nil {
				z80.tapeDrive_orNil.decelerate()
			}

			// Repeat emulating the HALT instruction until 'z80.eventNextEvent'
			for z80.tstates < z80.eventNextEvent {
				z80.memory.contendRead(z80.pc, 4)

				z80.r = (z80.r + 1) & 0x7f
				z80_localInstructionCounter++
			}
		}
	}

	// Update emulation efficiency counters
	if z80.perfCounter_hostCpuInstr != nil {
		ttid_end := z80.perfCounter_hostCpuInstr.Gettid()

		var hostCpu_instrCount_end uint64
		var hostCpu_instrCount_endErr error
		hostCpu_instrCount_end, hostCpu_instrCount_endErr = z80.perfCounter_hostCpuInstr.Read()

		z80.z80_instructionCounter += uint64(z80_localInstructionCounter)

		/*if z80_localInstructionCounter > 0 {
			println( z80_localInstructionCounter, hostCpu_instrCount_start, hostCpu_instrCount_end,
					hostCpu_instrCount_end-hostCpu_instrCount_start,
					(hostCpu_instrCount_end - hostCpu_instrCount_start) / uint64(z80_localInstructionCounter) )
		}*/

		if (hostCpu_instrCount_startErr == nil) &&
			(hostCpu_instrCount_endErr == nil) &&
			(ttid_start == ttid_end) &&
			(z80_localInstructionCounter > 0) &&
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

func invalidOpcode(z80 *Z80) {
	panic("invalid opcode")
}

func opcode_cb(z80 *Z80) {
	z80.memory.contendRead(z80.pc, 4)
	var opcode2 byte = z80.memory.readByteInternal(z80.pc)
	z80.pc++
	z80.r++
	opcodesMap[SHIFT_0xCB+int(opcode2)](z80)
}

func opcode_ed(z80 *Z80) {
	z80.memory.contendRead(z80.pc, 4)
	var opcode2 byte = z80.memory.readByteInternal(z80.pc)
	z80.pc++
	z80.r++

	if f := opcodesMap[SHIFT_0xED+int(opcode2)]; f != nil {
		f(z80)
	} else {
		invalidOpcode(z80)
	}
}

func opcode_dd(z80 *Z80) {
	z80.memory.contendRead(z80.pc, 4)
	var opcode2 byte = z80.memory.readByteInternal(z80.pc)
	z80.pc++
	z80.r++

	switch opcode2 {
	case 0xcb:
		z80.memory.contendRead(z80.pc, 3)
		z80.tempaddr = z80.IX() + uint16(signExtend(z80.memory.readByteInternal(z80.pc)))
		z80.pc++
		z80.memory.contendRead(z80.pc, 3)
		var opcode3 byte = z80.memory.readByteInternal(z80.pc)
		z80.memory.contendReadNoMreq_loop(z80.pc, 1, 2)
		z80.pc++
		opcodesMap[SHIFT_0xDDCB+int(opcode3)](z80)
	default:
		if f := opcodesMap[SHIFT_0xDD+int(opcode2)]; f != nil {
			f(z80)
		} else {
			/* Instruction did not involve H or L */
			opcodesMap[opcode2](z80)
		}
	}
}

func opcode_fd(z80 *Z80) {
	z80.memory.contendRead(z80.pc, 4)
	var opcode2 byte = z80.memory.readByteInternal(z80.pc)
	z80.pc++
	z80.r++

	switch opcode2 {
	case 0xcb:
		z80.memory.contendRead(z80.pc, 3)
		z80.tempaddr = z80.IY() + uint16(signExtend(z80.memory.readByteInternal(z80.pc)))
		z80.pc++
		z80.memory.contendRead(z80.pc, 3)
		var opcode3 byte = z80.memory.readByteInternal(z80.pc)
		z80.memory.contendReadNoMreq_loop(z80.pc, 1, 2)
		z80.pc++

		opcodesMap[SHIFT_0xFDCB+int(opcode3)](z80)

	default:
		if f := opcodesMap[SHIFT_0xFD+int(opcode2)]; f != nil {
			f(z80)
		} else {
			/* Instruction did not involve H or L */
			opcodesMap[opcode2](z80)
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
