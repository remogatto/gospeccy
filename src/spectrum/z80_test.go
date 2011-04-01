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
	"bufio"
	"container/vector"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"spectrum/formats"
	"strconv"
	"strings"
	"testing"
)

var (
	events        vector.StringVector
	initialMemory map[uint16]byte = make(map[uint16]byte)
	dirtyMemory   map[uint16]bool = make(map[uint16]bool)
)

func (z80 *Z80) DumpRegisters(out *vector.StringVector) {
	var halted byte

	if z80.halted {
		halted = 1
	} else {
		halted = 0
	}

	out.Push(fmt.Sprintf("%02x%02x %02x%02x %02x%02x %02x%02x %02x%02x %02x%02x %02x%02x %02x%02x %02x%02x %02x%02x %04x %04x\n",
		z80.a, z80.f, z80.b, z80.c, z80.d, z80.e, z80.h, z80.l, z80.a_, z80.f_, z80.b_, z80.c_, z80.d_, z80.e_, z80.h_, z80.l_, z80.ixh, z80.ixl, z80.iyh, z80.iyl, z80.sp, z80.pc))
	out.Push(fmt.Sprintf("%02x %02x %d %d %d %d %d\n", z80.i, (z80.r7&0x80)|byte(z80.r&0x7f),
		z80.iff1, z80.iff2, z80.im, halted, z80.tstates))
}

func (memory *testMemory) DumpMemory(out *vector.StringVector) {
	var addresses sort.IntArray = make([]int, 0, 32)

	for address, _ := range memory.data_map {
		addresses = append(addresses, int(address))
	}

	sort.SortInts(addresses)

	for i := 0; i < len(addresses); i++ {
		addr := uint16(addresses[i])

		if memory.Read(addr) == initialMemory[addr] {
			continue
		}

		line := fmt.Sprintf("%04x ", addr)

		for (memory.Read(addr) != initialMemory[addr]) || dirtyMemory[addr] {
			line += fmt.Sprintf("%02x ", memory.Read(addr))
			addr++
			i++

			if i >= len(addresses) {
				break
			}
			if addresses[i] != int(addr) {
				i--
				break
			}
		}

		line += fmt.Sprintf("-1\n")
		out.Push(line)
	}
}

type testMemory struct {
	data_array [0x10000]byte
	data_map   map[uint16]byte
	z80        *Z80
}

func (memory *testMemory) readByteInternal(addr uint16) byte {
	events.Push(fmt.Sprintf("%5d MR %04x %02x\n", memory.z80.tstates, addr, memory.data_array[addr]))
	return memory.data_array[addr]
}

func (memory *testMemory) writeByteInternal(address uint16, b byte) {
	events.Push(fmt.Sprintf("%5d MW %04x %02x\n", memory.z80.tstates, address, b))

	// Note: ROM is not protected from writes

	memory.data_array[address] = b
	memory.data_map[address] = b

	if b == 0 {
		dirtyMemory[address] = true
	}
}

func (memory *testMemory) readByte(addr uint16) byte {
	events.Push(fmt.Sprintf("%5d MC %04x\n", memory.z80.tstates, addr))
	contendMemory(memory.z80, addr, 3)
	return memory.readByteInternal(addr)
}

func (memory *testMemory) writeByte(address uint16, b byte) {
	events.Push(fmt.Sprintf("%5d MC %04x\n", memory.z80.tstates, address))
	contendMemory(memory.z80, address, 3)
	memory.writeByteInternal(address, b)
}

func (memory *testMemory) contendRead(address uint16, time uint) {
	events.Push(fmt.Sprintf("%5d MC %04x\n", memory.z80.tstates, address))
	contendMemory(memory.z80, address, time)
}

func (memory *testMemory) contendReadNoMreq(address uint16, time uint) {
	memory.contendRead(address, time)
}

func (memory *testMemory) contendReadNoMreq_loop(address uint16, time uint, count uint) {
	for i := uint(0); i < count; i++ {
		memory.contendReadNoMreq(address, time)
	}
}

func (memory *testMemory) contendWriteNoMreq(address uint16, time uint) {
	events.Push(fmt.Sprintf("%5d MC %04x\n", memory.z80.tstates, address))
	contendMemory(memory.z80, address, time)
}

func (memory *testMemory) contendWriteNoMreq_loop(address uint16, time uint, count uint) {
	for i := uint(0); i < count; i++ {
		memory.contendWriteNoMreq(address, time)
	}
}

func (memory *testMemory) Read(address uint16) byte {
	return memory.data_array[address]
}

func (memory *testMemory) Write(address uint16, value byte, protectROM bool) {
	// 'protectROM' is ignored
	memory.data_array[address] = value
	memory.data_map[address] = value
}

func (memory *testMemory) Data() *[0x10000]byte {
	return &memory.data_array
}

func (memory *testMemory) reset() {
	for address, _ := range memory.data_map {
		memory.data_array[address] = 0
	}

	memory.data_map = make(map[uint16]byte)
}


type testPort struct {
	z80 *Z80
}

func (p *testPort) readPortInternal(address uint16, contend bool) byte {
	if contend {
		p.contendPortPreio(address)
	}

	var r byte = byte(address >> 8)
	events.Push(fmt.Sprintf("%5d PR %04x %02x\n", p.z80.tstates, address, r))

	if contend {
		p.contendPortPostio(address)
	}

	return r
}

func (p *testPort) readPort(port uint16) byte {
	return p.readPortInternal(port, true)
}

func (p *testPort) writePortInternal(address uint16, b byte, contend bool) {
	if contend {
		p.contendPortPreio(address)
	}

	events.Push(fmt.Sprintf("%5d PW %04x %02x\n", p.z80.tstates, address, b))

	if contend {
		p.contendPortPostio(address)
	}
}

func (p *testPort) writePort(port uint16, b byte) {
	p.writePortInternal(port, b, true)
}

func (p *testPort) contendPortPreio(port uint16) {
	if (port & 0xc000) == 0x4000 {
		events.Push(fmt.Sprintf("%5d PC %04x\n", p.z80.tstates, port))
	}

	if (port & 0xc000) == 0x4000 {
		contendPort(p.z80, 1)
	} else {
		p.z80.tstates += 1
	}
}

func (p *testPort) contendPortPostio(port uint16) {
	if (port & 0x0001) == 1 {
		if (port & 0xc000) == 0x4000 {
			for i := 0; i < 3; i++ {
				events.Push(fmt.Sprintf("%5d PC %04x\n", p.z80.tstates, port))
				contendPort(p.z80, 1)
			}
		} else {
			p.z80.tstates += 3
		}

	} else {
		events.Push(fmt.Sprintf("%5d PC %04x\n", p.z80.tstates, port))
		contendPort(p.z80, 3)
	}
}

func (p *testPort) frame_begin() {
}

func (p *testPort) frame_end() FrameStatusOfPorts {
	return FrameStatusOfPorts{
		shouldPlayTheTape: false,
	}
}

func (p *testPort) getBorderEvents_orNil() *BorderEvent {
	return nil
}

func (p *testPort) getBeeperEvents_orNil() *BeeperEvent {
	return nil
}

func (p *testPort) reset() {
}


const maxLines = 20000

func TestDoOpcodes(t *testing.T) {

	var memory testMemory
	var port testPort
	memory.data_map = make(map[uint16]byte)

	// Instantiate a Z80 processor
	z80 := NewZ80(&memory, &port)
	ula := NewULA()
	z80.init(ula, /*tapeDrive_orNil*/ nil)
	ula.init(z80, &memory, &port)

	memory.z80 = z80
	port.z80 = z80

	// Read the "tests.in" file

	bytes, err := ioutil.ReadFile("testdata/tests.in")

	if err != nil {
		fmt.Println("Error reading tests.in")
	} else {
		content := string(bytes)
		lines := strings.Split(content, "\n", -1)

		currLine := 0

		for (currLine < len(lines)-1) && currLine < maxLines {

			// Skip all blank lines and consume the first non-blank line
			if lines[currLine] == "" {
				currLine++
				continue
			}

			currOp := lines[currLine]

			currLine++

			mainRegs := strings.Split(lines[currLine], " ", -1)

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

			otherRegs := strings.Split(lines[currLine], " ", -1)

			i, _ := strconv.Btoui64(otherRegs[0], 16)
			z80.i = byte(i)

			r, _ := strconv.Btoui64(otherRegs[1], 16)
			z80.r, z80.r7 = uint16(r), byte(r)

			iff1, _ := strconv.Btoui64(otherRegs[2], 16)
			z80.iff1 = byte(iff1)

			iff2, _ := strconv.Btoui64(otherRegs[3], 16)
			z80.iff2 = byte(iff2)

			im, _ := strconv.Btoui64(otherRegs[4], 16)
			z80.im = byte(im)

			halted, _ := strconv.Btoui64(otherRegs[5], 10)

			if halted != 0 {
				z80.halted = true
			} else {
				z80.halted = false
			}

			// Should set event_next_event and tstates

			event, _ := strconv.Btoui64(otherRegs[len(otherRegs)-1], 10)

			z80.eventNextEvent = uint(event)

			// Fill memory

			currLine++

			for lines[currLine] != "-1" {
				memWrites := strings.Split(lines[currLine], " ", -1)
				addr, _ := strconv.Btoui64(memWrites[0], 16)
				for i := 1; i < (len(memWrites)); i++ {
					byte := memWrites[i]
					if byte != "-1" {
						value, _ := strconv.Btoui64(byte, 16)
						z80.memory.Write(uint16(addr), uint8(value), /*protectROM*/ false)
						addr++
					}
				}
				currLine++
			}

			// Take a picture of the initial memory
			initialMemory = make(map[uint16]byte)
			for address, value := range memory.data_map {
				initialMemory[address] = value
			}

			// doOpcodes
			events.Push(currOp + "\n")
			z80.doOpcodes()

			// dump registers and memory
			z80.DumpRegisters(&events)
			memory.DumpMemory(&events)

			events.Push("\n")

			currLine++

			z80.reset()
			memory.reset()
			dirtyMemory = make(map[uint16]bool)
		}
	}

	// Read the "tests.expected" file

	if file, err := os.Open("testdata/tests.expected", os.O_RDONLY, 0); err != nil {
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
						fmt.Printf("F(%s) ", testDescription)
						t.Errorf("\nTest 0x%s failed at line %d\nEXPECTED: %sGOT:      %s\n", testDescription, currLine+1, l, events.At(currLine))
					} else {
						passed++
					}
				}
			}

			currLine++
		}

	}

}


func BenchmarkZ80(b *testing.B) {
	b.StopTimer()

	rom, err := ReadROM("testdata/48.rom")
	if err != nil {
		panic(err)
	}

	app := NewApplication()
	speccy := NewSpectrum48k(app, *rom)

	snapshot, err := formats.ReadProgram("testdata/fire.z80")
	if err != nil {
		panic(err)
	}

	err = speccy.loadSnapshot(snapshot.(formats.Snapshot))
	if err != nil {
		panic(err)
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		speccy.CommandChannel <- Cmd_RenderFrame{CompletionTime_orNil: nil}
		//speccy.renderFrame(/*completionTime_orNil*/ nil)
	}

	app.RequestExit()
	<-app.HasTerminated
}
