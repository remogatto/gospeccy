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
	"testing"
	"io/ioutil"
	"fmt"
	"strings"
	"strconv"
	"os"
	"bufio"
	"container/vector"
)

var events *vector.StringVector = new(vector.StringVector)

type testMemory [0x10000]byte

func (memory *testMemory) readByteInternal(addr uint16) byte {
	events.Push(fmt.Sprintf("%5d MR %04x %02x\n", tstates, addr, memory[addr]))
	return memory[addr]
}

func (memory *testMemory) readByte(addr uint16) byte {
	events.Push(fmt.Sprintf("%5d MC %04x\n", tstates, addr))
	tstates += 3
	return memory.readByteInternal(addr)
}

func (memory *testMemory) writeByte(address uint16, b byte) {
	events.Push(fmt.Sprintf("%5d MC %04x\n", tstates, address))
	tstates += 3
	memory.writeByteInternal(address, b)
}

func (memory *testMemory) writeByteInternal(address uint16, b byte) {
	events.Push(fmt.Sprintf("%5d MW %04x %02x\n", tstates, address, b))
	memory[address] = b
}

func (memory *testMemory) contendRead(addr uint16, time uint) {
	events.Push(fmt.Sprintf("%5d MC %04x\n", tstates, addr))
	tstates += time
}

func (memory *testMemory) contendReadNoMreq(address uint16, time uint) {
	memory.contendRead(address, time)
}

func (memory *testMemory) contendWriteNoMreq(address uint16, time uint) {
	events.Push(fmt.Sprintf("%5d MC %04x\n", tstates, address))
	tstates += time
}

func (memory *testMemory) renderScreen() {

}

func (memory *testMemory) At(address uint) byte {
	return memory[address]
}

func (memory *testMemory) set(address uint, value byte) {
	memory[address] = value
}

func (memory *testMemory) Data() []byte {
	return []byte(memory)
}

type testPort struct{}

func (p *testPort) readPort(port uint16) byte {
	var r byte = byte(port >> 8)
	p.contendPortPreio(port)

	events.Push(fmt.Sprintf("%5d PR %04x %02x\n", tstates, port, r))

	p.contendPortPostio(port)
	return r
}

func (p *testPort) writePort(port uint16, b byte) {
	p.contendPortPreio(port)

	events.Push(fmt.Sprintf("%5d PW %04x %02x\n", tstates, port, b))

	p.contendPortPostio(port)
}

func (p *testPort) contendPortPreio(port uint16) {
	if (port & 0xc000) == 0x4000 {
		events.Push(fmt.Sprintf("%5d PC %04x\n", tstates, port))

	}
	tstates++
}

func (p *testPort) contendPortPostio(port uint16) {
	if (port & 0x0001) != 0 {
		if (port & 0xc000) == 0x4000 {

			events.Push(fmt.Sprintf("%5d PC %04x\n", tstates, port))
			tstates++
			events.Push(fmt.Sprintf("%5d PC %04x\n", tstates, port))
			tstates++
			events.Push(fmt.Sprintf("%5d PC %04x\n", tstates, port))
			tstates++

		} else {
			tstates += 3
		}

	} else {
		events.Push(fmt.Sprintf("%5d PC %04x\n", tstates, port))
		tstates += 3

	}

}

var maxLines = 20000

func TestDoOpcodes(t *testing.T) {

	var memory testMemory
	var port testPort

	// Instantiate a Z80 processor
	z80 := NewZ80(&memory, &port)

	// Read the test.in file

	bytes, err := ioutil.ReadFile("tests.in")

	if err != nil {
		fmt.Println("Error reading tests.in")
	} else {
		content := string(bytes)
		lines := strings.Split(content, "\n", 0)

		currLine := 0

		for (currLine < len(lines)-1) && currLine < maxLines {

			// Skip all blank lines and consume the first non-blank line
			if lines[currLine] == "" {
				currLine++
				continue
			}

			currOp := lines[currLine]

			currLine++

			mainRegs := strings.Split(lines[currLine], " ", 0)

			// Fill registers

			af, _ := strconv.Btoui64(mainRegs[0], 16)
			z80.a, z80.f = byte(int16(af)>>8), byte(int16(af)&0xff)

			bc, _ := strconv.Btoui64(mainRegs[1], 16)
			z80.b, z80.c = byte(int16(bc)>>8), byte(int16(bc)&0xff)

			de, _ := strconv.Btoui64(mainRegs[2], 16)
			z80.d, z80.e = byte(int16(de)>>8), byte(int16(de)&0xff)

			hl, _ := strconv.Btoui64(mainRegs[3], 16)
			z80.h, z80.l = byte(int16(hl)>>8), byte(int16(hl)&0xff)

			af_, _ := strconv.Btoui64(mainRegs[4], 16)
			z80.a_, z80.f_ = byte(int16(af_)>>8), byte(int16(af_)&0xff)

			bc_, _ := strconv.Btoui64(mainRegs[5], 16)
			z80.b_, z80.c_ = byte(int16(bc_)>>8), byte(int16(bc_)&0xff)

			de_, _ := strconv.Btoui64(mainRegs[6], 16)
			z80.d_, z80.e_ = byte(int16(de_)>>8), byte(int16(de_)&0xff)

			hl_, _ := strconv.Btoui64(mainRegs[7], 16)
			z80.h_, z80.l_ = byte(int16(hl_)>>8), byte(int16(hl_)&0xff)

			ix, _ := strconv.Btoui64(mainRegs[8], 16)
			z80.ixh, z80.ixl = byte(int16(ix)>>8), byte(int16(ix)&0xff)

			iy, _ := strconv.Btoui64(mainRegs[9], 16)
			z80.iyh, z80.iyl = byte(int16(iy)>>8), byte(int16(iy)&0xff)

			sp, _ := strconv.Btoui64(mainRegs[10], 16)
			z80.sp = uint16(sp)

			pc, _ := strconv.Btoui64(mainRegs[11], 16)
			z80.pc = uint16(pc)

			currLine++

			otherRegs := strings.Split(lines[currLine], " ", 0)

			i, _ := strconv.Btoui64(otherRegs[0], 16)
			z80.i = byte(i)

			r, _ := strconv.Btoui64(otherRegs[1], 16)
			z80.r, z80.r7 = uint16(r), uint16(r)

			iff1, _ := strconv.Btoui64(otherRegs[2], 16)
			z80.iff1 = uint16(iff1)

			iff2, _ := strconv.Btoui64(otherRegs[3], 16)
			z80.iff2 = uint16(iff2)

			im, _ := strconv.Btoui64(otherRegs[4], 16)
			z80.im = uint16(im)

			halted, _ := strconv.Btoui64(otherRegs[5], 10)
			z80.halted = byte(halted)

			// Should set event_next_event and tstates

			event, _ := strconv.Btoui64(otherRegs[len(otherRegs)-1], 10)
			eventNextEvent = uint(event)

			// Fill memory

			currLine++

			for lines[currLine] != "-1" {
				memWrites := strings.Split(lines[currLine], " ", 0)
				addr, _ := strconv.Btoui64(memWrites[0], 16)
				for i := 1; i < (len(memWrites)); i++ {
					byte := memWrites[i]
					if byte != "-1" {
						value, _ := strconv.Btoui64(byte, 16)
						z80.memory.set(uint(addr), uint8(value))
						addr++
					}
				}
				currLine++
			}

			// Take a picture of the initial memory

			for i, val := range z80.memory.Data() {
				initialMemory[i] = val
			}

			// doOpcodes

			z80.LogEvents = true

			events.Push(currOp + "\n")

			z80.doOpcodes()

			// dump registers and memory and save the

			z80.DumpRegisters(events)
			z80.DumpMemory(events)

			events.Push("\n")

			currLine++

			z80.Reset()
		}
	}


	// Read the tests.expected file

	if file, err := os.Open("tests.expected", os.O_RDONLY, 0); err != nil {
		t.Fatalf("Error opening tests.expected\n")
	} else {
		var nextIsTestDescription bool = false
		var testDescription string
		currLine := 0
		passed := 0
		buf := bufio.NewReader(file)
		for {
			l, err := buf.ReadString('\n') // parse line-by-line
			if err == os.EOF || currLine >= maxLines {
				break
			} else if err != nil {
				t.Fatalf("Error reading file\n")
			} else {

				if nextIsTestDescription {
					testDescription = strings.Trim(l, "\n")
					nextIsTestDescription = false
				}

				if l == "\n" {
					nextIsTestDescription = true
				}

				if currLine >= events.Len() {
					t.Errorf("** No events at line %d **", currLine)
				} else {
					if l != events.At(currLine) {
						// diff with expected
						fmt.Print("F")
						t.Errorf("\nTest 0x%s failed at line %d\nEXPECTED: %sGOT:      %s\n", testDescription, currLine+1, l, events.At(currLine))
					} else {
						passed++
						fmt.Print(".")
					}
				}
			}

			currLine++

		}

	}

}
