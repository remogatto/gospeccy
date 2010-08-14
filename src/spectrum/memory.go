package spectrum

type MemoryAccessor interface {
	readByte(address uint16) byte
	readByteInternal(address uint16) byte

	writeByte(address uint16, value byte)
	writeByteInternal(address uint16, value byte)

	contendRead(addr uint16, time uint)
	contendReadNoMreq(addr uint16, time uint)

	contendWriteNoMreq(addr uint16, time uint)

	Read(addr uint16) byte
	Write(addr uint16, value byte)
	Data() *[0x10000]byte
	
	reset()
}

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


func (memory *Memory) readByteInternal(addr uint16) byte {
	return memory.data[addr]
}

func (memory *Memory) writeByteInternal(address uint16, b byte) {
	if (address >= SCREEN_BASE_ADDR) && (address < ATTR_BASE_ADDR) {
		if memory.data[address] != b {
			memory.speccy.ula.screenBitmapWrite(address, b)
		}
	} else if (address >= ATTR_BASE_ADDR) && (address < 0x5b00) {
		memory.speccy.ula.screenAttrWrite(address, b)
	}

	memory.data[address] = b
}

func (memory *Memory) readByte(addr uint16) byte {
	memory.contend(addr, 3)
	return memory.readByteInternal(addr)
}

func (memory *Memory) writeByte(addr uint16, b byte) {
	memory.contend(addr, 3)
	memory.writeByteInternal(addr, b)
}

func (memory *Memory) contend(address uint16, time uint) {
	tstates_p := &memory.speccy.Cpu.tstates

	if (address >= 0x4000) && (address <= 0x7fff) {
		*tstates_p += uint(delay_table[*tstates_p])
	}

	*tstates_p += time
}

func (memory *Memory) contendRead(address uint16, time uint) {
	memory.contend(address, time)
}

func (memory *Memory) contendReadNoMreq(address uint16, time uint) {
	memory.contend(address, time)
}

func (memory *Memory) contendWriteNoMreq(address uint16, time uint) {
	memory.contend(address, time)
}

func (memory *Memory) Read(address uint16) byte {
	return memory.data[address]
}

func (memory *Memory) Write(address uint16, value byte) {
	memory.data[address] = value
}

func (memory *Memory) Data() *[0x10000]byte {
	return &memory.data
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
