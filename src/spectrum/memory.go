package spectrum

import "github.com/remogatto/z80"

type Memory struct {
	data [0x10000]byte
	speccy *Spectrum48k
}

func NewMemory() *Memory {
	return &Memory{}
}

func (memory *Memory) init(speccy *Spectrum48k) {
	memory.speccy = speccy
}

func (memory *Memory) reset() {
	for i := 0; i < 0x10000; i++ {
		memory.data[i] = 0
	}
}

func (memory *Memory) ReadByteInternal(address uint16) byte {
	return memory.data[address]
}

func (memory *Memory) WriteByteInternal(address uint16, b byte) {
	if (address >= SCREEN_BASE_ADDR) && (address < ATTR_BASE_ADDR) {
		memory.speccy.ula.screenBitmapWrite(address, memory.data[address], b)
	} else if (address >= ATTR_BASE_ADDR) && (address < 0x5b00) {
		memory.speccy.ula.screenAttrWrite(address, memory.data[address], b)
	}

	if address >= 0x4000 {
		memory.data[address] = b
	}
}

func (memory *Memory) ReadByte(address uint16) byte {
	contendMemory(memory.speccy.Cpu, address, 3)
	return memory.ReadByteInternal(address)
}

func (memory *Memory) WriteByte(address uint16, b byte) {
	contendMemory(memory.speccy.Cpu, address, 3)
	memory.WriteByteInternal(address, b)
}

func contendMemory(z80 *z80.Z80, address uint16, time int) {
	tstates_p := &z80.Tstates
	tstates := *tstates_p

	if (address & 0xc000) == 0x4000 {
		tstates += int(delay_table[tstates])
	}

	tstates += time

	*tstates_p = tstates
}

// Equivalent to executing "contendMemory(z80, address, time)" count times
func contendMemory_loop(z80 *z80.Z80, address uint16, time int, count uint) {
	tstates_p := &z80.Tstates
	tstates := *tstates_p

	if (address & 0xc000) == 0x4000 {
		for i := uint(0); i < count; i++ {
			tstates += int(delay_table[tstates])
			tstates += time
		}
	} else {
		tstates += time * int(count)
	}

	*tstates_p = tstates
}

func (memory *Memory) ContendRead(address uint16, time int) {
	contendMemory(memory.speccy.Cpu, address, time)
}

func (memory *Memory) ContendReadNoMreq(address uint16, time int) {
	contendMemory(memory.speccy.Cpu, address, time)
}

func (memory *Memory) ContendReadNoMreq_loop(address uint16, time int, count uint) {
	contendMemory_loop(memory.speccy.Cpu, address, time, count)
}

func (memory *Memory) ContendWriteNoMreq(address uint16, time int) {
	contendMemory(memory.speccy.Cpu, address, time)
}

func (memory *Memory) ContendWriteNoMreq_loop(address uint16, time int, count uint) {
	contendMemory_loop(memory.speccy.Cpu, address, time, count)
}

func (memory *Memory) Read(address uint16) byte {
	return memory.data[address]
}

func (memory *Memory) Write(address uint16, value byte, protectROM bool) {
	if (address >= 0x4000) || !protectROM {
		memory.data[address] = value
	}
}

func (memory *Memory) Data() []byte {
	return memory.data[:]
}

// Number of T-states to delay, for each possible T-state within a frame.
// The array is extended at the end - this covers the case when the emulator
// begins to execute an instruction at Tstate=(TStatesPerFrame-1). Such an
// instruction will finish at (TStatesPerFrame-1+4) or later.
var delay_table [TStatesPerFrame + 100]byte

// Initialize 'delay_table' at program startup
func init() {
	// Note: The language automatically initialized all values
	//       of the 'delay_table' array to zeroes. So, we only
	//       have to modify the non-zero elements.

	tstate := FIRST_SCREEN_BYTE - 1
	for y := 0; y < ScreenHeight; y++ {
		for x := 0; x < ScreenWidth; x += 16 {
			tstate_x := x / PIXELS_PER_TSTATE
			delay_table[tstate+tstate_x+0] = 6
			delay_table[tstate+tstate_x+1] = 5
			delay_table[tstate+tstate_x+2] = 4
			delay_table[tstate+tstate_x+3] = 3
			delay_table[tstate+tstate_x+4] = 2
			delay_table[tstate+tstate_x+5] = 1
		}
		tstate += TSTATES_PER_LINE
	}
}
