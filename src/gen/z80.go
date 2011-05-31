// z80.go: generate Go code for Z80 opcodes
//
// Based on z80.pl by Philip Kendall
// Based on modified z80.pl by Andrea Fazzi
//
// Copyright (c) 1999-2006 Philip Kendall <philip-fuse@shadowmagic.org.uk>
// Copyright (c) 2010 Andrea Fazzi <andrea.fazzi@alcacoop.it>
// Copyright (c) 2011 ⚛ <0xe2.0x9a.0x9b@gmail.com> 
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along
// with this program; if not, write to the Free Software Foundation, Inc.,
// 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// The status of which flags relates to which condition

// These conditions involve !( F & FLAG_<whatever> )
var not = map[string]bool{
	"NC": true,
	"NZ": true,
	"P":  true,
	"PO": true,
}

// Use F & FLAG_<whatever>
var flag = map[string]string{
	"C":  "C",
	"NC": "C",
	"PE": "P",
	"PO": "P",
	"M":  "S",
	"P":  "S",
	"Z":  "Z",
	"NZ": "Z",
}

// Returns whether 's' matches 'pattern'
func matches(s, pattern string) bool {
	return regexp.MustCompile(pattern).MatchString(s)
}

// Return the lowercased version of 's'
func lc(s string) string {
	return strings.ToLower(s)
}

// Joins the strings in the 'stringList', prints them on to standard output
// and sends a new line to standard output
func ln(stringList ...string) {
	for _, s := range stringList {
		fmt.Print(s)
	}
	fmt.Println()
}

// Returns 'if_true' or 'if_false' depending on the value of 'cond'.
// The purpose of this function is to reduce the number of code lines.
func _if(cond bool, if_true, if_false string) string {
	if cond {
		return if_true
	}
	return if_false
}

// Generalised opcode routines

func arithmetic_logical(opcode, arg1, arg2 string) {
	if arg2 == "" {
		arg2 = arg1
		arg1 = "A"
	}

	if len(arg1) == 1 {
		if len(arg2) == 1 || matches(arg2, "^REGISTER[HL]$") {
			lc_opcode, lc_arg2 := lc(opcode), lc(arg2)
			ln("z80.", lc_opcode, "(z80.", lc_arg2, ")")
		} else if arg2 == "(REGISTER+dd)" {
			lc_opcode := lc(opcode)
			ln("var offset byte = z80.memory.readByte( z80.pc )")
			ln("z80.memory.contendReadNoMreq_loop( z80.pc, 1, 5 )")
			ln("z80.pc++")
			ln("var bytetemp byte = z80.memory.readByte(z80.REGISTER() + uint16(signExtend(offset)))")
			ln("z80.", lc_opcode, "(bytetemp)")
		} else {
			register := _if(arg2 == "(HL)", "HL", "PC")
			increment := _if(register == "PC", "z80.pc++", "")
			lc_opcode := lc(opcode)
			ln("var bytetemp byte = z80.memory.readByte(z80.", register, "())")
			ln(increment)
			ln("z80.", lc_opcode, "(bytetemp)")
		}
	} else if opcode == "ADD" {
		lc_opcode := lc(opcode)
		lc_arg1 := lc(arg1)
		ln("z80.memory.contendReadNoMreq_loop( z80.IR(), 1, 7 )")
		ln("z80.", lc_opcode, "16(z80.", lc_arg1, ", z80.", arg2, "())")
	} else if (arg1 == "HL") && len(arg2) == 2 {
		lc_opcode := lc(opcode)
		ln("z80.memory.contendReadNoMreq_loop( z80.IR(), 1, 7 )")
		ln("z80.", lc_opcode, "16(z80.", arg2, "())")
	}
}

func call_jp(opcode, condition, offset string) {
	lc_opcode := lc(opcode)
	if offset == "" {
		ln("z80.", lc_opcode, "()")
	} else {
		var condition_string string
		if not[condition] {
			condition_string = "(z80.f & FLAG_" + flag[condition] + ") == 0"
		} else {
			condition_string = "(z80.f & FLAG_" + flag[condition] + ") != 0"
		}
		ln("if ", condition_string, "{")
		ln("  z80.", lc_opcode, "()")
		ln("} else {")
		ln("  z80.memory.contendRead(z80.pc, 3); z80.memory.contendRead( z80.pc + 1, 3 ); z80.pc += 2;")
		ln("}")
	}
}

func cpi_cpd(opcode string) {
	modifier := _if(opcode == "CPI", "inc", "dec")

	ln("var value byte = z80.memory.readByte( z80.HL() )")
	ln("var bytetemp byte = z80.a - value")
	ln("var lookup byte = ((z80.a & 0x08 ) >> 3 ) | (((value) & 0x08 ) >> 2 ) | ((bytetemp & 0x08 ) >> 1)")

	ln("z80.memory.contendReadNoMreq_loop( z80.HL(), 1, 5 )")
	ln("z80.", modifier, "HL(); z80.decBC()")
	ln("z80.f = (z80.f & FLAG_C) | ternOpB(z80.BC() != 0, FLAG_V | FLAG_N, FLAG_N) | halfcarrySubTable[lookup] | ternOpB(bytetemp != 0, 0, FLAG_Z) | (bytetemp & FLAG_S )")
	ln("if((z80.f & FLAG_H) != 0) { bytetemp-- }")
	ln("z80.f |= (bytetemp & FLAG_3) | ternOpB((bytetemp & 0x02) != 0, FLAG_5, 0)")
}

func cpir_cpdr(opcode string) {
	modifier := _if(opcode == "CPIR", "inc", "dec")

	ln("var value byte = z80.memory.readByte( z80.HL() )")
	ln("var bytetemp byte = z80.a - value")
	ln("var lookup byte = ((z80.a & 0x08) >> 3) | (((value) & 0x08) >> 2) | ((bytetemp & 0x08) >> 1)")

	ln("z80.memory.contendReadNoMreq_loop( z80.HL(), 1, 5 )")
	ln("z80.decBC()")
	ln("z80.f = ( z80.f & FLAG_C ) | ( ternOpB(z80.BC() != 0, ( FLAG_V | FLAG_N ),FLAG_N)) | halfcarrySubTable[lookup] | ( ternOpB(bytetemp != 0, 0, FLAG_Z )) | ( bytetemp & FLAG_S )")
	ln("if((z80.f & FLAG_H) != 0) {")
	ln("  bytetemp--")
	ln("}")
	ln("z80.f |= ( bytetemp & FLAG_3 ) | ternOpB((bytetemp & 0x02) != 0, FLAG_5, 0)")
	ln("if( ( z80.f & ( FLAG_V | FLAG_Z ) ) == FLAG_V ) {")
	ln("  z80.memory.contendReadNoMreq_loop( z80.HL(), 1, 5 )")
	ln("  z80.pc -= 2")
	ln("}")
	ln("z80.", modifier, "HL()")
}

func inc_dec(opcode, arg string) {
	modifier := _if(opcode == "INC", "inc", "dec")

	if len(arg) == 1 || matches(arg, "^REGISTER[HL]$") {
		ln("z80.", lc(opcode), arg, "()")
	} else if len(arg) == 2 || arg == "REGISTER" {
		ln("z80.memory.contendReadNoMreq_loop( z80.IR(), 1, 2 )")
		ln("z80.", modifier, arg, "()")
	} else if arg == "(HL)" {
		ln("{")
		ln("  var bytetemp byte = z80.memory.readByte( z80.HL() )")
		ln("  z80.memory.contendReadNoMreq( z80.HL(), 1 )")
		ln("  z80.", lc(opcode), "(&bytetemp)")
		ln("  z80.memory.writeByte(z80.HL(), bytetemp)")
		ln("}")
	} else if arg == "(REGISTER+dd)" {
		ln("var offset byte = z80.memory.readByte( z80.pc )")
		ln("z80.memory.contendReadNoMreq_loop( z80.pc, 1, 5 )")
		ln("z80.pc++")
		ln("var wordtemp uint16 = z80.REGISTER() + uint16(signExtend(offset))")
		ln("var bytetemp byte = z80.memory.readByte( wordtemp )")
		ln("z80.memory.contendReadNoMreq( wordtemp, 1 )")
		ln("z80.", lc(opcode), "(&bytetemp)")
		ln("z80.memory.writeByte(wordtemp,bytetemp)")
	}

}

func ini_ind(opcode string) {
	modifier := _if(opcode == "INI", "inc", "dec")
	operation := _if(opcode == "INI", "+", "-")

	ln("z80.memory.contendReadNoMreq( z80.IR(), 1 );")
	ln("var initemp byte = z80.readPort(z80.BC());")
	ln("z80.memory.writeByte( z80.HL(), initemp );")
	ln()
	ln("z80.b--; z80.", modifier, "HL()")
	ln("var initemp2 byte = initemp + z80.c ", operation, " 1;")
	ln("z80.f = ternOpB((initemp & 0x80) != 0, FLAG_N, 0) |")
	ln("        ternOpB(initemp2 < initemp, FLAG_H | FLAG_C, 0) |")
	ln("        ternOpB(parityTable[(initemp2 & 0x07) ^ z80.b] != 0, FLAG_P, 0 ) |")
	ln("        sz53Table[z80.b]")
}

func inir_indr(opcode string) {
	modifier := _if(opcode == "INIR", "inc", "dec")
	operation := _if(opcode == "INIR", "+", "-")

	ln("z80.memory.contendReadNoMreq( z80.IR(), 1 );")
	ln("var initemp byte = z80.readPort(z80.BC());")
	ln("z80.memory.writeByte( z80.HL(), initemp );")
	ln()
	ln("z80.b--;")
	ln("var initemp2 byte = initemp + z80.c ", operation, " 1;")
	ln("z80.f = ternOpB(initemp & 0x80 != 0, FLAG_N, 0) |")
	ln("        ternOpB(initemp2 < initemp, FLAG_H | FLAG_C, 0 ) |")
	ln("        ternOpB(parityTable[ ( initemp2 & 0x07 ) ^ z80.b ] != 0, FLAG_P, 0) |")
	ln("        sz53Table[z80.b];")
	ln()
	ln("if( z80.b != 0 ) {")
	ln("  z80.memory.contendWriteNoMreq_loop( z80.HL(), 1, 5 )")
	ln("  z80.pc -= 2")
	ln("}")
	ln("z80.", modifier, "HL()")
}


func ldi_ldd(opcode string) {
	modifier := _if(opcode == "LDI", "inc", "dec")

	ln("var bytetemp byte = z80.memory.readByte( z80.HL() )")
	ln("z80.decBC()")
	ln("z80.memory.writeByte(z80.DE(), bytetemp);")
	ln("z80.memory.contendWriteNoMreq_loop( z80.DE(), 1, 2 )")
	ln("z80.", modifier, "DE(); z80.", modifier, "HL();")
	ln("bytetemp += z80.a;")
	ln("z80.f = ( z80.f & ( FLAG_C | FLAG_Z | FLAG_S ) ) |")
	ln("        ternOpB(z80.BC() != 0, FLAG_V, 0) |")
	ln("        ( bytetemp & FLAG_3 ) |")
	ln("        ternOpB((bytetemp & 0x02) != 0, FLAG_5, 0)")
}

func ldir_lddr(opcode string) {
	modifier := _if(opcode == "LDIR", "inc", "dec")

	ln("var bytetemp byte = z80.memory.readByte( z80.HL() )")
	ln("z80.memory.writeByte(z80.DE(), bytetemp);")
	ln("z80.memory.contendWriteNoMreq_loop(z80.DE(), 1, 2)")
	ln("z80.decBC()")
	ln("bytetemp += z80.a;")
	ln("z80.f = (z80.f & ( FLAG_C | FLAG_Z | FLAG_S )) | ternOpB(z80.BC() != 0, FLAG_V, 0 ) | (bytetemp & FLAG_3) | ternOpB((bytetemp & 0x02 != 0), FLAG_5, 0 )")
	ln("if(z80.BC() != 0) {")
	ln("  z80.memory.contendWriteNoMreq_loop( z80.DE(), 1, 5 )")
	ln("  z80.pc -= 2")
	ln("}")
	ln("z80.", modifier, "HL(); z80.", modifier, "DE()")
}

func otir_otdr(opcode string) {
	modifier := _if(opcode == "OTIR", "inc", "dec")

	ln("z80.memory.contendReadNoMreq( z80.IR(), 1 );")
	ln("var outitemp byte = z80.memory.readByte( z80.HL() );")
	ln("z80.b--;	/* This does happen first, despite what the specs say */")
	ln("z80.writePort(z80.BC(), outitemp);")
	ln()
	ln("z80.", modifier, "HL()")
	ln("var outitemp2 byte = outitemp + z80.l;")
	ln("z80.f = ternOpB((outitemp & 0x80) != 0, FLAG_N, 0 ) |")
	ln("    ternOpB(outitemp2 < outitemp, FLAG_H | FLAG_C, 0) |")
	ln("    ternOpB(parityTable[ ( outitemp2 & 0x07 ) ^ z80.b ] != 0, FLAG_P, 0 ) |")
	ln("    sz53Table[z80.b]")
	ln()
	ln("if( z80.b != 0 ) {")
	ln("  z80.memory.contendReadNoMreq_loop( z80.BC(), 1, 5 )")
	ln("  z80.pc -= 2")
	ln("}")
}

func outi_outd(opcode string) {
	modifier := _if(opcode == "OUTI", "inc", "dec")

	ln("z80.memory.contendReadNoMreq( z80.IR(), 1 )")
	ln("var outitemp byte = z80.memory.readByte( z80.HL() )")
	ln("z80.b--;	/* This does happen first, despite what the specs say */")
	ln("z80.writePort(z80.BC(), outitemp)")
	ln()
	ln("z80.", modifier, "HL()")
	ln("var outitemp2 byte = outitemp + z80.l")
	ln("z80.f = ternOpB((outitemp & 0x80) != 0, FLAG_N, 0) |")
	ln("        ternOpB(outitemp2 < outitemp, FLAG_H | FLAG_C, 0) |")
	ln("        ternOpB(parityTable[ ( outitemp2 & 0x07 ) ^ z80.b ] != 0, FLAG_P, 0 ) |")
	ln("        sz53Table[z80.b]")
}

func push_pop(opcode, regpair string) {
	var high, low string

	if regpair == "REGISTER" {
		high, low = "REGISTERH", "REGISTERL"
	} else {
		high, low = regpair[0:0+1], regpair[1:1+1]
	}

	lc_opcode := lc(opcode)
	lc_low := lc(low)
	lc_high := lc(high)
	if lc_opcode == "pop" {
		ln("z80.", lc_low, ", z80.", lc_high, " = z80.", lc_opcode, "16()")
	} else {
		ln("z80.", lc_opcode, "16(z80.", lc_low, ", z80.", lc_high, ")")
	}
}

func res_set_hexmask(opcode string, bit uint) string {
	mask := 1 << bit
	if opcode == "RES" {
		mask = 0xff - mask
	}

	return fmt.Sprintf("0x%02x", mask)
}

func res_set(opcode string, bit uint, register string) {
	operator := _if(opcode == "RES", "&", "|")
	hex_mask := res_set_hexmask(opcode, bit)
	lc_register := lc(register)

	if len(register) == 1 {
		ln("z80.", lc_register, " ", operator, "= ", hex_mask)
	} else if register == "(HL)" {
		ln("var bytetemp byte = z80.memory.readByte( z80.HL() )")
		ln("z80.memory.contendReadNoMreq( z80.HL(), 1 )")
		ln("z80.memory.writeByte( z80.HL(), bytetemp ", operator, " ", hex_mask, " )")
	} else if register == "(REGISTER+dd)" {
		ln("var bytetemp byte = z80.memory.readByte( z80.tempaddr )")
		ln("z80.memory.contendReadNoMreq( z80.tempaddr, 1 )")
		ln("z80.memory.writeByte( z80.tempaddr, bytetemp ", operator, " ", hex_mask, " )")
	} else {
		panic("invalid register: " + register)
	}
}

func rotate_shift(opcode, register string) {
	lc_opcode := lc(opcode)
	lc_register := lc(register)

	if len(register) == 1 {
		ln("z80.", lc_register, " = z80.", lc_opcode, "(z80.", lc_register, ")")
	} else if register == "(HL)" {
		ln("var bytetemp byte = z80.memory.readByte(z80.HL())")
		ln("z80.memory.contendReadNoMreq( z80.HL(), 1 )")
		ln("bytetemp = z80.", lc_opcode, "(bytetemp)")
		ln("z80.memory.writeByte(z80.HL(),bytetemp)")
	} else if register == "(REGISTER+dd)" {
		ln("var bytetemp byte = z80.memory.readByte(z80.tempaddr)")
		ln("z80.memory.contendReadNoMreq( z80.tempaddr, 1 )")
		ln("bytetemp = z80.", lc_opcode, "(bytetemp)")
		ln("z80.memory.writeByte(z80.tempaddr, bytetemp)")
	} else {
		panic("invalid register: " + register)
	}
}

// A type on which we define methods corresponding to opcodes.
// The methods are found at run-time by using Go's reflection capabilities.
type Opcode byte

// Individual opcode routines

func (Opcode) ADC(a, b string) { arithmetic_logical("ADC", a, b) }

func (Opcode) ADD(a, b string) { arithmetic_logical("ADD", a, b) }

func (Opcode) AND(a, b string) { arithmetic_logical("AND", a, b) }

func (Opcode) BIT(bit, register string) {
	lc_register := lc(register)
	if len(register) == 1 {
		ln("z80.bit(", bit, ", z80.", lc_register, ")")
	} else if register == "(REGISTER+dd)" {
		ln("bytetemp := z80.memory.readByte( z80.tempaddr )")
		ln("z80.memory.contendReadNoMreq( z80.tempaddr, 1 )")
		ln("z80.biti(", bit, ", bytetemp, z80.tempaddr)")
	} else {
		ln("bytetemp := z80.memory.readByte( z80.HL() )")
		ln("z80.memory.contendReadNoMreq( z80.HL(), 1 )")
		ln("z80.bit(", bit, ", bytetemp)")
	}
}

func (Opcode) CALL(a, b string) { call_jp("CALL", a, b) }

func (Opcode) CCF() {
	ln("z80.f = ( z80.f & ( FLAG_P | FLAG_Z | FLAG_S ) ) |")
	ln("        ternOpB( ( z80.f & FLAG_C ) != 0, FLAG_H, FLAG_C ) |")
	ln("        ( z80.a & ( FLAG_3 | FLAG_5 ) )")
}

func (Opcode) CP(a, b string) { arithmetic_logical("CP", a, b) }

func (Opcode) CPD() { cpi_cpd("CPD") }

func (Opcode) CPDR() { cpir_cpdr("CPDR") }

func (Opcode) CPI() { cpi_cpd("CPI") }

func (Opcode) CPIR() { cpir_cpdr("CPIR") }

func (Opcode) CPL() {
	ln("z80.a ^= 0xff")
	ln("z80.f = ( z80.f & ( FLAG_C | FLAG_P | FLAG_Z | FLAG_S ) ) |")
	ln("        ( z80.a & ( FLAG_3 | FLAG_5 ) ) | ")
	ln("        ( FLAG_N | FLAG_H )")
}

func (Opcode) DAA() {
	ln("var add, carry byte = 0, ( z80.f & FLAG_C )")
	ln("if( ( (z80.f & FLAG_H ) != 0) || ( ( z80.a & 0x0f ) > 9 ) ) { add = 6 }")
	ln("if( (carry != 0) || ( z80.a > 0x99 ) ) { add |= 0x60 }")
	ln("if( z80.a > 0x99 ) { carry = FLAG_C }")
	ln("if( (z80.f & FLAG_N) != 0 ) {")
	ln("  z80.sub(add)")
	ln("} else {")
	ln("  z80.add(add)")
	ln("}")
	ln("var temp byte = byte(int(z80.f) & ^(FLAG_C | FLAG_P)) | carry | parityTable[z80.a]")
	ln("z80.f = temp")
}

func (Opcode) DEC(a string) { inc_dec("DEC", a) }

func (Opcode) DI() { ln("z80.iff1, z80.iff2 = 0, 0") }

func (Opcode) DJNZ() {
	ln("z80.memory.contendReadNoMreq(z80.IR(), 1)")
	ln("z80.b--")
	ln("if(z80.b != 0) {")
	ln("  z80.jr()")
	ln("} else {")
	ln("  z80.memory.contendRead( z80.pc, 3 )")
	ln("}")
	ln("z80.pc++")
}

func (Opcode) EI() {
	ln("/* Interrupts are not accepted immediately after an EI, but are")
	ln("   accepted after the next instruction */")
	ln("z80.iff1, z80.iff2 = 1, 1")
	ln("z80.interruptsEnabledAt = int(z80.tstates)")
	ln("// eventAdd(z80.tstates + 1, z80InterruptEvent)")
}

func (Opcode) EX(arg1, arg2 string) {
	if (arg1 == "AF") && (arg2 == "AF'") {
		ln("/* Tape saving trap: note this traps the EX AF,AF' at #04d0, not")
		ln("   #04d1 as PC has already been incremented */")
		ln("/* 0x76 - Timex 2068 save routine in EXROM */")
		ln("if( z80.pc == 0x04d1 || z80.pc == 0x0077 ) {")
		ln("  if( z80.tapeSaveTrap() == 0 ) { /*break*/ }")
		ln("}")
		ln()
		ln("var olda, oldf = z80.a, z80.f")
		ln("z80.a = z80.a_; z80.f = z80.f_")
		ln("z80.a_ = olda; z80.f_ = oldf")
	} else if (arg1 == "(SP)") && (arg2 == "HL" || arg2 == "REGISTER") {
		var high, low string

		if arg2 == "HL" {
			high, low = "H", "L"
		} else {
			high, low = "REGISTERH", "REGISTERL"
		}
		lc_low := lc(low)
		lc_high := lc(high)

		ln("var bytetempl = z80.memory.readByte( z80.SP() )")
		ln("var bytetemph = z80.memory.readByte( z80.SP() + 1 )")
		ln("z80.memory.contendReadNoMreq( z80.SP() + 1, 1 )")
		ln("z80.memory.writeByte( z80.SP() + 1, z80.", lc_high, " )")
		ln("z80.memory.writeByte( z80.SP(),     z80.", lc_low, "  )")
		ln("z80.memory.contendWriteNoMreq_loop( z80.SP(), 1, 2 )")
		ln("z80.", lc_low, " = bytetempl")
		ln("z80.", lc_high, " = bytetemph")
	} else if (arg1 == "DE") && (arg2 == "HL") {
		ln("var wordtemp uint16 = z80.DE()")
		ln("z80.setDE(z80.HL())")
		ln("z80.setHL(wordtemp)")
	} else {
		panic("invalid args: " + arg1 + ", " + arg2)
	}
}

func (Opcode) EXX() {
	ln("var wordtemp uint16 = z80.BC()")
	ln("z80.setBC(z80.BC_())")
	ln("z80.setBC_(wordtemp)")
	ln()
	ln("wordtemp = z80.DE()")
	ln("z80.setDE(z80.DE_())")
	ln("z80.setDE_(wordtemp)")
	ln()
	ln("wordtemp = z80.HL()")
	ln("z80.setHL(z80.HL_())")
	ln("z80.setHL_(wordtemp)")
}

func (Opcode) HALT() {
	ln("z80.halted = true")
	ln("z80.pc--")
	ln("return")
}

func (Opcode) IM(mode string) {
	ln("z80.im = ", mode)
}

func (Opcode) IN(register, port string) {
	if (register == "A") && (port == "(nn)") {
		ln("var intemp uint16 = uint16(z80.memory.readByte(z80.pc)) + (uint16(z80.a) << 8 )")
		ln("z80.pc++")
		ln("z80.a = z80.readPort(intemp)")
	} else if (register == "F") && (port == "(C)") {
		ln("var bytetemp byte")
		ln("z80.in(&bytetemp, z80.BC())")
	} else if (len(register) == 1) && (port == "(C)") {
		lc_register := lc(register)
		ln("z80.in(&z80.", lc_register, ", z80.BC())")
	} else {
		panic("invalid args: " + register + ", " + port)
	}
}

func (Opcode) INC(a string) { inc_dec("INC", a) }

func (Opcode) IND() { ini_ind("IND") }

func (Opcode) INDR() { inir_indr("INDR") }

func (Opcode) INI() { ini_ind("INI") }

func (Opcode) INIR() { inir_indr("INIR") }

func (Opcode) JP(condition, offset string) {
	if (condition == "HL") || (condition == "REGISTER") {
		ln("z80.pc = z80.", condition, "()\t\t/* NB: NOT INDIRECT! */")
	} else {
		call_jp("JP", condition, offset)
	}
}

func (Opcode) JR(condition, offset string) {
	if offset == "" {
		offset = condition
		condition = ""
	}

	if condition == "" {
		ln("z80.jr()")
	} else {
		var condition_string string
		if not[condition] {
			condition_string = "(z80.f & FLAG_" + flag[condition] + ") == 0"
		} else {
			condition_string = "(z80.f & FLAG_" + flag[condition] + ") != 0"
		}
		ln("if ", condition_string, " {")
		ln("  z80.jr()")
		ln("} else {")
		ln("  z80.memory.contendRead( z80.pc, 3 )")
		ln("}")
	}

	ln("z80.pc++")
}

func (Opcode) LD(dest, src string) {
	if (len(dest) == 1) || matches(dest, "^REGISTER[HL]$") {
		if (len(src) == 1) || matches(src, "^REGISTER[HL]$") {
			if (dest == "R") && (src == "A") {
				ln("z80.memory.contendReadNoMreq( z80.IR(), 1 )")
				ln("/* Keep the RZX instruction counter right */")
				ln("z80.rzxInstructionsOffset += ( int(z80.r) - int(z80.a))")
				ln("z80.r, z80.r7 = uint16(z80.a), z80.a")
			} else if (dest == "A") && (src == "R") {
				ln("z80.memory.contendReadNoMreq( z80.IR(), 1 )")
				ln("z80.a = byte(z80.r&0x7f) | (z80.r7 & 0x80)")
				ln("z80.f = ( z80.f & FLAG_C ) | sz53Table[z80.a] | ternOpB(z80.iff2 != 0, FLAG_V, 0)")
			} else {
				if (src == "I") || (dest == "I") {
					ln("z80.memory.contendReadNoMreq( z80.IR(), 1 )")
				}
				lc_dest, lc_src := lc(dest), lc(src)
				if dest != src {
					ln("z80.", lc_dest, " = z80.", lc_src)
				}
				if (dest == "A") && (src == "I") {
					ln("z80.f = ( z80.f & FLAG_C ) | sz53Table[z80.a] | ternOpB(z80.iff2 != 0, FLAG_V, 0)")
				}
			}
		} else if src == "nn" {
			ln("z80.", lc(dest), " = z80.memory.readByte(z80.pc)")
			ln("z80.pc++")
		} else if matches(src, "^\\(..\\)$") {
			register := src[1 : 1+2]
			ln("z80.", lc(dest), " = z80.memory.readByte(z80.", register, "())")
		} else if src == "(nnnn)" {
			ln("var wordtemp uint16 = uint16(z80.memory.readByte(z80.pc))")
			ln("z80.pc++")
			ln("wordtemp |= uint16(z80.memory.readByte(z80.pc)) << 8")
			ln("z80.pc++")
			ln("z80.a = z80.memory.readByte(wordtemp)")
		} else if src == "(REGISTER+dd)" {
			ln("var offset byte = z80.memory.readByte( z80.pc )")
			ln("z80.memory.contendReadNoMreq_loop( z80.pc, 1, 5 )")
			ln("z80.pc++")
			ln("z80.", lc(dest), " = z80.memory.readByte(z80.REGISTER() + uint16(signExtend(offset)))")
		} else {
			panic("invalid src: " + src)
		}
	} else if (len(dest) == 2) || (dest == "REGISTER") {
		var high, low string

		if (dest == "SP") || (dest == "REGISTER") {
			high, low = (dest + "H"), (dest + "L")
		} else {
			high, low = dest[0:0+1], dest[1:1+1]
		}

		if src == "nnnn" {
			ln("b1 := z80.memory.readByte(z80.pc)")
			ln("z80.pc++")
			ln("b2 := z80.memory.readByte(z80.pc)")
			ln("z80.pc++")
			ln("z80.set", high, low, "(joinBytes(b2, b1))")
		} else if (src == "HL") || (src == "REGISTER") {
			ln("z80.memory.contendReadNoMreq_loop( z80.IR(), 1, 2 )")
			ln("z80.sp = z80.", src, "()")
		} else if src == "(nnnn)" {
			if low == "SPL" {
				ln("sph, spl := splitWord(z80.sp)\nz80.ld16rrnn(&spl, &sph)\nz80.sp = joinBytes(sph, spl)\n // break")
			} else {
				ln("z80.ld16rrnn(&z80.", lc(low), ", &z80.", lc(high), ")\n // break")
			}
		} else {
			panic("invalid src: " + src)
		}
	} else if matches(dest, "^\\(..\\)$") {
		register := dest[1 : 1+2]

		if len(src) == 1 {
			ln("z80.memory.writeByte(z80.", register, "(),z80.", lc(src), ")")
		} else if src == "nn" {
			ln("z80.memory.writeByte(z80.", register, "(),z80.memory.readByte(z80.pc))")
			ln("z80.pc++")
		} else {
			panic("invalid src: " + src)
		}
	} else if dest == "(nnnn)" {
		if src == "A" {
			ln("var wordtemp uint16 = uint16(z80.memory.readByte(z80.pc))")
			ln("z80.pc++")
			ln("wordtemp |= uint16(z80.memory.readByte(z80.pc)) << 8")
			ln("z80.pc++")
			ln("z80.memory.writeByte(wordtemp, z80.a)")
		} else if (len(src) == 2) || (src == "REGISTER") {
			var high, low string

			if (src == "SP") || (src == "REGISTER") {
				high, low = (src + "H"), (src + "L")
			} else {
				high, low = src[0:0+1], src[1:1+1]
			}
			if low == "SPL" {
				ln("sph, spl := splitWord(z80.sp)\nz80.ld16nnrr(spl, sph)\n // break")
			} else {
				ln("z80.ld16nnrr(z80.", lc(low), ", z80.", lc(high), ")\n // break")
			}
		} else {
			panic("invalid src: " + src)
		}
	} else if dest == "(REGISTER+dd)" {
		if len(src) == 1 {
			ln("offset := z80.memory.readByte( z80.pc )")
			ln("z80.memory.contendReadNoMreq_loop( z80.pc, 1, 5 )")
			ln("z80.pc++")
			ln("z80.memory.writeByte(z80.REGISTER() + uint16(signExtend(offset)), z80.", lc(src), " )")
		} else if src == "nn" {
			ln("offset := z80.memory.readByte( z80.pc )")
			ln("z80.pc++")
			ln("value := z80.memory.readByte( z80.pc )")
			ln("z80.memory.contendReadNoMreq_loop( z80.pc, 1, 2 )")
			ln("z80.pc++")
			ln("z80.memory.writeByte(z80.REGISTER() + uint16(signExtend(offset)), value )")
		} else {
			panic("invalid src: " + src)
		}
	} else {
		panic("invalid dest: " + dest)
	}
}

func (Opcode) LDD() { ldi_ldd("LDD") }

func (Opcode) LDDR() { ldir_lddr("LDDR") }

func (Opcode) LDI() { ldi_ldd("LDI") }

func (Opcode) LDIR() { ldir_lddr("LDIR") }

func (Opcode) NEG() {
	ln("bytetemp := z80.a")
	ln("z80.a = 0")
	ln("z80.sub(bytetemp)")
}

func (Opcode) NOP() {}

func (Opcode) OR(a, b string) { arithmetic_logical("OR", a, b) }

func (Opcode) OTDR() { otir_otdr("OTDR") }

func (Opcode) OTIR() { otir_otdr("OTIR") }

func (Opcode) OUT(port, register string) {
	if (port == "(nn)") && (register == "A") {
		ln("var outtemp uint16 = uint16(z80.memory.readByte(z80.pc)) + (uint16(z80.a) << 8)")
		ln("z80.pc++")
		ln("z80.writePort(outtemp, z80.a)")
	} else if (port == "(C)") && (len(register) == 1) {
		if register == "0" {
			ln("z80.writePort(z80.BC(), ", register, ")")
		} else {
			lc_register := lc(register)
			ln("z80.writePort(z80.BC(), z80.", lc_register, ")")
		}
	} else {
		panic("invalid args: " + port + ", " + register)
	}
}


func (Opcode) OUTD() { outi_outd("OUTD") }

func (Opcode) OUTI() { outi_outd("OUTI") }

func (Opcode) POP(a string) { push_pop("POP", a) }

func (Opcode) PUSH(regpair string) {
	ln("z80.memory.contendReadNoMreq( z80.IR(), 1 )")
	push_pop("PUSH", regpair)
}

func (Opcode) RES(bit, register string) {
	bitNum, err := strconv.Atoui(bit)
	if err != nil {
		panic(err.String())
	}
	res_set("RES", bitNum, register)
}

func (Opcode) RET(condition string) {
	if condition == "" {
		ln("z80.ret()")
	} else {
		ln("z80.memory.contendReadNoMreq( z80.IR(), 1 )")

		if condition == "NZ" {
			ln("if( z80.pc==0x056c || z80.pc == 0x0112 ) {")
			ln("  if(z80.tapeLoadTrap() == 0) { /*break*/ }")
			ln("}")
		}

		if not[condition] {
			ln("if(!((z80.f & FLAG_", flag[condition], ") != 0)) { z80.ret() }")
		} else {
			ln("if((z80.f & FLAG_", flag[condition], ") != 0) { z80.ret() }")
		}
	}
}

func (Opcode) RETN() {
	ln("z80.iff1 = z80.iff2")
	ln("z80.ret()")
}

func (Opcode) RL(a string) { rotate_shift("RL", a) }

func (Opcode) RLC(a string) { rotate_shift("RLC", a) }

func (Opcode) RLCA() {
	ln("z80.a = ( z80.a << 1 ) | ( z80.a >> 7 )")
	ln("z80.f = ( z80.f & ( FLAG_P | FLAG_Z | FLAG_S ) ) |")
	ln("        ( z80.a & ( FLAG_C | FLAG_3 | FLAG_5 ) )")
}

func (Opcode) RLA() {
	ln("var bytetemp byte = z80.a")
	ln("z80.a = ( z80.a << 1 ) | ( z80.f & FLAG_C )")
	ln("z80.f = ( z80.f & ( FLAG_P | FLAG_Z | FLAG_S ) ) | ( z80.a & ( FLAG_3 | FLAG_5 ) ) | ( bytetemp >> 7 )")
}

func (Opcode) RLD() {
	ln("var bytetemp byte = z80.memory.readByte( z80.HL() )")
	ln("z80.memory.contendReadNoMreq_loop( z80.HL(), 1, 4 )")
	ln("z80.memory.writeByte(z80.HL(), (bytetemp << 4 ) | ( z80.a & 0x0f ) )")
	ln("z80.a = ( z80.a & 0xf0 ) | ( bytetemp >> 4 )")
	ln("z80.f = ( z80.f & FLAG_C ) | sz53pTable[z80.a]")
}

func (Opcode) RR(a string) { rotate_shift("RR", a) }

func (Opcode) RRA(a string) {
	ln("var bytetemp byte = z80.a")
	ln("z80.a = ( z80.a >> 1 ) | ( z80.f << 7 )")
	ln("z80.f = ( z80.f & ( FLAG_P | FLAG_Z | FLAG_S ) ) | ( z80.a & ( FLAG_3 | FLAG_5 ) ) | ( bytetemp & FLAG_C )")
}

func (Opcode) RRC(a string) { rotate_shift("RRC", a) }

func (Opcode) RRCA() {
	ln("z80.f = ( z80.f & ( FLAG_P | FLAG_Z | FLAG_S ) ) | ( z80.a & FLAG_C )")
	ln("z80.a = ( z80.a >> 1) | ( z80.a << 7 )")
	ln("z80.f |= ( z80.a & ( FLAG_3 | FLAG_5 ) )")
}

func (Opcode) RRD() {
	ln("var bytetemp byte = z80.memory.readByte( z80.HL() )")
	ln("z80.memory.contendReadNoMreq_loop( z80.HL(), 1, 4 )")
	ln("z80.memory.writeByte(z80.HL(),  ( z80.a << 4 ) | ( bytetemp >> 4 ) )")
	ln("z80.a = ( z80.a & 0xf0 ) | ( bytetemp & 0x0f )")
	ln("z80.f = ( z80.f & FLAG_C ) | sz53pTable[z80.a]")
}

func (Opcode) RST(value string) {
	ln("z80.memory.contendReadNoMreq( z80.IR(), 1 )")
	ln("z80.rst(0x", value, ")")
}

func (Opcode) SBC(a, b string) { arithmetic_logical("SBC", a, b) }

func (Opcode) SCF() {
	ln("z80.f = ( z80.f & ( FLAG_P | FLAG_Z | FLAG_S ) ) |")
	ln("        ( z80.a & ( FLAG_3 | FLAG_5          ) ) |")
	ln("        FLAG_C")
}

func (Opcode) SET(bit, register string) {
	bitNum, err := strconv.Atoui(bit)
	if err != nil {
		panic(err.String())
	}
	res_set("SET", bitNum, register)
}

func (Opcode) SLA(a string) { rotate_shift("SLA", a) }

func (Opcode) SLL(a string) { rotate_shift("SLL", a) }

func (Opcode) SRA(a string) { rotate_shift("SRA", a) }

func (Opcode) SRL(a string) { rotate_shift("SRL", a) }

func (Opcode) SUB(a, b string) { arithmetic_logical("SUB", a, b) }

func (Opcode) XOR(a, b string) { arithmetic_logical("XOR", a, b) }

func (Opcode) SLTTRAP() {
	ln("z80.sltTrap(int16(z80.HL()), z80.a)")
}

// Description of each file
var description = map[string]string{
	"opcodes_cb.dat":     "z80_cb.c: Z80 CBxx opcodes",
	"opcodes_ddfd.dat":   "z80_ddfd.c Z80 {DD,FD}xx opcodes",
	"opcodes_ddfdcb.dat": "z80_ddfdcb.c Z80 {DD,FD}CBxx opcodes",
	"opcodes_ed.dat":     "z80_ed.c: Z80 CBxx opcodes",
	"opcodes_base.dat":   "opcodes_base.c: unshifted Z80 opcodes",
}

var funcTable = make(map[string]reflect.Value)
func init() {
	var opcode Opcode
	var opcodeType reflect.Type = reflect.TypeOf(opcode)
	var opcodeValue reflect.Value = reflect.ValueOf(opcode)

	numMethods := opcodeType.NumMethod()
	for i:=0; i<numMethods; i++ {
		funcName := opcodeType.Method(i).Name
		funcValue := opcodeValue.Method(i)
		funcTable[funcName] = funcValue
	}
}

// Main program
func main() {
	data_file := os.Args[1]

	var data []byte
	var err os.Error
	data, err = ioutil.ReadFile(data_file)
	if err != nil {
		panic(err.String())
	}

	lines := strings.Split(string(data), "\n", -1)

	var fallthrough_cases []string

	for _, line := range lines {
		// Remove comments
		if strings.Contains(line, "#") {
			line = line[0:strings.Index(line, "#")]
		}

		line = strings.TrimSpace(line)

		// Skip blank lines
		if len(line) == 0 {
			continue
		}

		var l []string = strings.Split(line, " ", -1)

		var number, opcode, arguments, extra string
		number = l[0]
		if len(l) >= 2 {
			opcode = l[1]
		}
		if len(l) >= 3 {
			arguments = l[2]
		}
		if len(l) >= 4 {
			extra = l[3]
		}

		var args []string
		if arguments != "" {
			args = strings.Split(arguments, ",", -1)
		}

		var shift_op string
		if data_file == "opcodes_cb.dat" {
			shift_op = "shift0xcb(" + number + ")"
		} else if data_file == "opcodes_ed.dat" {
			shift_op = "shift0xed(" + number + ")"
		} else if data_file == "opcodes_ddfd.dat" {
			shift_op = "shift0xdd(" + number + ")"
		} else if data_file == "opcodes_ddfdcb.dat" {
			shift_op = "shift0xddcb(" + number + ")"
		} else {
			shift_op = number
		}

		if opcode != "" {
			comment := "/* " + opcode
			if arguments != "" {
				comment += " " + arguments
			}
			if extra != "" {
				comment += " " + extra
			}
			comment += " */"
			ln(comment)
		} else {
			fallthrough_cases = append(fallthrough_cases, shift_op)
			continue
		}

		ln("opcodesMap[", shift_op, "] = func (z80 *Z80) {")

		// Handle the undocumented rotate-shift-or-bit and store-in-register opcodes specially
		if extra != "" {
			register, opcode2 := args[0], args[1]
			lc_register, lc_opcode2 := lc(register), lc(opcode2)

			if (opcode2 == "RES") || (opcode2 == "SET") {
				bit := strings.Split(extra, ",", -1)[0]
				bitNum, err2 := strconv.Atoui(bit)
				if err2 != nil {
					panic("invalid bit number: " + bit)
				}

				operator := _if(opcode2 == "RES", "&", "|")
				hexmask := res_set_hexmask(opcode2, bitNum)

				ln("  z80.", lc_register, " = z80.memory.readByte(z80.tempaddr) ", operator, " ", hexmask)
				ln("  z80.memory.contendReadNoMreq(z80.tempaddr, 1)")
				ln("  z80.memory.writeByte(z80.tempaddr, z80.", lc_register, ")")
				ln("}")
			} else {
				ln("  z80.", lc_register, " = z80.memory.readByte(z80.tempaddr)")
				ln("  z80.memory.contendReadNoMreq( z80.tempaddr, 1 )")
				ln("  z80.", lc_register, " = z80.", lc_opcode2, "(z80.", lc_register, ")")
				ln("  z80.memory.writeByte(z80.tempaddr, z80.", lc_register, ")")
				ln("}")
			}
			continue
		}

		if fn := funcTable[strings.ToUpper(opcode)]; fn.IsValid() {
			reflect_args := make([]reflect.Value, len(args))
			for i, arg := range args {
				reflect_args[i] = reflect.ValueOf(arg)
			}

			// Missing arguments are substituted with ""
			fn_numArgs := fn.Type().NumIn()
			for len(reflect_args) < fn_numArgs {
				reflect_args = append(reflect_args, reflect.ValueOf(""))
			}

			//fmt.Printf("%s   %#v\n", opcode, args)
			//println(fn.String())

			// Call the method. Excessive arguments are simply ignored.
			fn.Call(reflect_args[0:fn_numArgs])
		}

		ln("}")

		if len(fallthrough_cases) > 0 {
			ln("// Fallthrough cases")
			for _, fallthrough_case := range fallthrough_cases {
				ln("opcodesMap[", fallthrough_case, "] = opcodesMap[", shift_op, "]")
			}
			fallthrough_cases = fallthrough_cases[0:0]
		}
	}
}
