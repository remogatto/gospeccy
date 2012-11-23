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

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"testing"
)

var (
	events        []string
	initialMemory map[uint16]byte = make(map[uint16]byte)
	dirtyMemory   map[uint16]bool = make(map[uint16]bool)
)

func (z80 *Z80) DumpRegisters(out *[]string) {
	var halted byte

	if z80.Halted {
		halted = 1
	} else {
		halted = 0
	}

	*out = append(*out, fmt.Sprintf("%02x%02x %02x%02x %02x%02x %02x%02x %02x%02x %02x%02x %02x%02x %02x%02x %02x%02x %02x%02x %04x %04x\n",
		z80.A, z80.F, z80.B, z80.C, z80.D, z80.E, z80.H, z80.L, z80.A_, z80.F_, z80.B_, z80.C_, z80.D_, z80.E_, z80.H_, z80.L_, z80.IXH, z80.IXL, z80.IYH, z80.IYL, z80.sp, z80.pc))
	*out = append(*out, fmt.Sprintf("%02x %02x %d %d %d %d %d\n", z80.I, (z80.R7&0x80)|byte(z80.R&0x7f),
		z80.IFF1, z80.IFF2, z80.IM, halted, z80.Tstates))
}

func (memory *testMemory) DumpMemory(out *[]string) {
	var addresses sort.IntSlice = make([]int, 0, 32)

	for address, _ := range memory.data_map {
		addresses = append(addresses, int(address))
	}

	sort.Ints(addresses)

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
		*out = append(*out, line)
	}
}

func contendMemory(z80 *Z80, address uint16, time uint) {
	tstates_p := &z80.Tstates
	tstates := *tstates_p

	tstates += time

	*tstates_p = tstates
}

func contendPort(z80 *Z80, time uint) {
	tstates_p := &z80.Tstates
	*tstates_p += time
}

type testMemory struct {
	data_array [0x10000]byte
	data_map   map[uint16]byte
	z80        *Z80
}

func (memory *testMemory) ReadByteInternal(addr uint16) byte {
	events = append(events, fmt.Sprintf("%5d MR %04x %02x\n", memory.z80.Tstates, addr, memory.data_array[addr]))
	return memory.data_array[addr]
}

func (memory *testMemory) WriteByteInternal(address uint16, b byte) {
	events = append(events, fmt.Sprintf("%5d MW %04x %02x\n", memory.z80.Tstates, address, b))
	memory.data_array[address] = b
	memory.data_map[address] = b
	if b == 0 {
		dirtyMemory[address] = true
	}
}

func (memory *testMemory) ReadByte(addr uint16) byte {
	events = append(events, fmt.Sprintf("%5d MC %04x\n", memory.z80.Tstates, addr))
	contendMemory(memory.z80, addr, 3)
	return memory.ReadByteInternal(addr)
}

func (memory *testMemory) WriteByte(address uint16, b byte) {
	events = append(events, fmt.Sprintf("%5d MC %04x\n", memory.z80.Tstates, address))
	contendMemory(memory.z80, address, 3)
	memory.WriteByteInternal(address, b)
}

func (memory *testMemory) ContendRead(address uint16, time uint) {
	events = append(events, fmt.Sprintf("%5d MC %04x\n", memory.z80.Tstates, address))
	contendMemory(memory.z80, address, time)
}

func (memory *testMemory) ContendReadNoMreq(address uint16, time uint) {
	memory.ContendRead(address, time)
}

func (memory *testMemory) ContendReadNoMreq_loop(address uint16, time uint, count uint) {
	for i := uint(0); i < count; i++ {
		memory.ContendReadNoMreq(address, time)
	}
}

func (memory *testMemory) ContendWriteNoMreq(address uint16, time uint) {
	events = append(events, fmt.Sprintf("%5d MC %04x\n", memory.z80.Tstates, address))
	contendMemory(memory.z80, address, time)
}

func (memory *testMemory) ContendWriteNoMreq_loop(address uint16, time uint, count uint) {
	for i := uint(0); i < count; i++ {
		memory.ContendWriteNoMreq(address, time)
	}
}

func (memory *testMemory) Read(address uint16) byte {
	return memory.data_array[address]
}

func (memory *testMemory) Write(address uint16, value byte, protectROM bool) {
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

func (p *testPort) ReadPortInternal(address uint16, contend bool) byte {
	if contend {
		p.ContendPortPreio(address)
	}

	var r byte = byte(address >> 8)
	events = append(events, fmt.Sprintf("%5d PR %04x %02x\n", p.z80.Tstates, address, r))

	if contend {
		p.ContendPortPostio(address)
	}
	return r
}

func (p *testPort) ReadPort(port uint16) byte {
	return p.ReadPortInternal(port, true)
}

func (p *testPort) WritePortInternal(address uint16, b byte, contend bool) {
	if contend {
		p.ContendPortPreio(address)
	}

	events = append(events, fmt.Sprintf("%5d PW %04x %02x\n", p.z80.Tstates, address, b))

	if contend {
		p.ContendPortPostio(address)
	}
}

func (p *testPort) WritePort(port uint16, b byte) {
	p.WritePortInternal(port, b, true)
}

func (p *testPort) ContendPortPreio(port uint16) {
	if (port & 0xc000) == 0x4000 {
		events = append(events, fmt.Sprintf("%5d PC %04x\n", p.z80.Tstates, port))
	}
	p.z80.Tstates += 1
}

func (p *testPort) ContendPortPostio(port uint16) {
	if (port & 0x0001) == 1 {
		if (port & 0xc000) == 0x4000 {
			for i := 0; i < 3; i++ {
				events = append(events, fmt.Sprintf("%5d PC %04x\n", p.z80.Tstates, port))
				contendPort(p.z80, 1)
			}
		} else {
			p.z80.Tstates += 3
		}

	} else {
		events = append(events, fmt.Sprintf("%5d PC %04x\n", p.z80.Tstates, port))
		contendPort(p.z80, 3)
	}
}

const maxLines = 20000

func TestDoOpcodes(t *testing.T) {

	var memory testMemory
	var port testPort
	memory.data_map = make(map[uint16]byte)

	// Instantiate a Z80 processor
	z80 := NewZ80(&memory, &port)
//	ula := NewULA()
//	z80.init(ula, nil /*tapeDrive_orNil*/)
//	ula.init(z80, &memory, &port)

	memory.z80 = z80
	port.z80 = z80

	// Read the "tests.in" file

	bytes, err := ioutil.ReadFile("testdata/tests.in")

	if err != nil {
		fmt.Println("Error reading tests.in")
	} else {
		content := string(bytes)
		lines := strings.Split(content, "\n")

		currLine := 0

		for (currLine < len(lines)-1) && currLine < maxLines {

			// Skip all blank lines and consume the first non-blank line
			if lines[currLine] == "" {
				currLine++
				continue
			}

			currOp := lines[currLine]

			currLine++

			mainRegs := strings.Split(lines[currLine], " ")

			// Fill registers

			af, _ := strconv.ParseUint(mainRegs[0], 16, 0)
			z80.A, z80.F = byte(int16(af)>>8), byte(uint16(af)&0xff)

			bc, _ := strconv.ParseUint(mainRegs[1], 16, 0)
			z80.B, z80.C = byte(int16(bc)>>8), byte(uint16(bc)&0xff)

			de, _ := strconv.ParseUint(mainRegs[2], 16, 0)
			z80.D, z80.E = byte(int16(de)>>8), byte(uint16(de)&0xff)

			hl, _ := strconv.ParseUint(mainRegs[3], 16, 0)
			z80.H, z80.L = byte(int16(hl)>>8), byte(uint16(hl)&0xff)

			af_, _ := strconv.ParseUint(mainRegs[4], 16, 0)
			z80.A_, z80.F_ = byte(int16(af_)>>8), byte(uint16(af_)&0xff)

			bc_, _ := strconv.ParseUint(mainRegs[5], 16, 0)
			z80.B_, z80.C_ = byte(int16(bc_)>>8), byte(uint16(bc_)&0xff)

			de_, _ := strconv.ParseUint(mainRegs[6], 16, 0)
			z80.D_, z80.E_ = byte(int16(de_)>>8), byte(uint16(de_)&0xff)

			hl_, _ := strconv.ParseUint(mainRegs[7], 16, 0)
			z80.H_, z80.L_ = byte(int16(hl_)>>8), byte(uint16(hl_)&0xff)

			ix, _ := strconv.ParseUint(mainRegs[8], 16, 0)
			z80.IXH, z80.IXL = byte(int16(ix)>>8), byte(uint16(ix)&0xff)

			iy, _ := strconv.ParseUint(mainRegs[9], 16, 0)
			z80.IYH, z80.IYL = byte(int16(iy)>>8), byte(uint16(iy)&0xff)

			sp, _ := strconv.ParseUint(mainRegs[10], 16, 0)
			z80.sp = uint16(sp)

			pc, _ := strconv.ParseUint(mainRegs[11], 16, 0)
			z80.pc = uint16(pc)

			currLine++

			otherRegs := strings.Split(lines[currLine], " ")

			i, _ := strconv.ParseUint(otherRegs[0], 16, 0)
			z80.I = byte(i)

			r, _ := strconv.ParseUint(otherRegs[1], 16, 0)
			z80.R, z80.R7 = uint16(r), byte(r)

			iff1, _ := strconv.ParseUint(otherRegs[2], 16, 0)
			z80.IFF1 = byte(iff1)

			iff2, _ := strconv.ParseUint(otherRegs[3], 16, 0)
			z80.IFF2 = byte(iff2)

			im, _ := strconv.ParseUint(otherRegs[4], 16, 0)
			z80.IM = byte(im)

			halted, _ := strconv.ParseUint(otherRegs[5], 10, 0)

			if halted != 0 {
				z80.Halted = true
			} else {
				z80.Halted = false
			}

			// Should set event_next_event and tstates

			event, _ := strconv.ParseUint(otherRegs[len(otherRegs)-1], 10, 0)

			z80.EventNextEvent = uint(event)

			// Fill memory

			currLine++

			for lines[currLine] != "-1" {
				memWrites := strings.Split(lines[currLine], " ")
				addr, _ := strconv.ParseUint(memWrites[0], 16, 0)
				for i := 1; i < (len(memWrites)); i++ {
					byte := memWrites[i]
					if byte != "-1" {
						value, _ := strconv.ParseUint(byte, 16, 0)
						z80.memory.Write(uint16(addr), uint8(value), false /*protectROM*/)
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
			events = append(events, currOp+"\n")
			z80.doOpcodes()

			// dump registers and memory
			z80.DumpRegisters(&events)
			memory.DumpMemory(&events)

			events = append(events, "\n")

			currLine++

			z80.Reset()
			memory.reset()
			dirtyMemory = make(map[uint16]bool)
		}
	}

	// Read the "tests.expected" file

	if file, err := os.Open("testdata/tests.expected"); err != nil {
		t.Fatalf("Error opening tests.expected\n")
	} else {
		var nextIsTestDescription bool = false
		var testDescription string
		currLine := 0
		passed := 0
		buf := bufio.NewReader(file)
		for {
			l, err := buf.ReadString('\n') // parse line-by-line
			if err == io.EOF || currLine >= maxLines {
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

				if currLine >= len(events) {
					t.Errorf("** No events at line %d **", currLine)
				} else {
					if l != events[currLine] {
						// diff with expected
						fmt.Printf("F(%s) ", testDescription)
						t.Errorf("\nTest 0x%s failed at line %d\nEXPECTED: %sGOT:      %s\n", testDescription, currLine+1, l, events[currLine])
					} else {
						passed++
					}
				}
			}

			currLine++
		}

	}

}

// func BenchmarkZ80(b *testing.B) {
// 	b.StopTimer()

// 	rom, err := ReadROM("testdata/48.rom")
// 	if err != nil {
// 		panic(err)
// 	}

// 	app := NewApplication()
// 	speccy := NewSpectrum48k(app, *rom)

// 	snapshot, err := formats.ReadProgram("testdata/fire.z80")
// 	if err != nil {
// 		panic(err)
// 	}

// 	err = speccy.loadSnapshot(snapshot.(formats.Snapshot))
// 	if err != nil {
// 		panic(err)
// 	}

// 	b.StartTimer()

// 	for i := 0; i < b.N; i++ {
// 		speccy.CommandChannel <- Cmd_RenderFrame{CompletionTime_orNil: nil}
// 		//speccy.renderFrame(nil /*completionTime_orNil*/)
// 	}

// 	app.RequestExit()
// 	<-app.HasTerminated
// }
