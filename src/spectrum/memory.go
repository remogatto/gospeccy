package spectrum

type MemoryAccessor interface {
	readByte(address uint16) byte
	readByteInternal(address uint16) byte

	writeByte(address uint16, value byte)
	writeByteInternal(address uint16, value byte)

	contendRead(addr uint16, time uint)
	contendReadNoMreq(addr uint16, time uint)

	contendWriteNoMreq(addr uint16, time uint)

	set(addr uint16, value byte)

	At(addr uint) byte
	Data() []byte

	getDirtyScreen() []bool

	frame_begin()
}

type Memory struct {
	data        [0x10000]byte
	dirtyScreen [ScreenWidth_Attr * ScreenHeight_Attr]bool
	z80         *Z80
}

func NewMemory() *Memory {
	return &Memory{}
}

func (memory *Memory) readByteInternal(addr uint16) byte {
	return memory.data[addr]
}

func (memory *Memory) readByte(addr uint16) byte {
	memory.contend(addr, 3)
	return memory.readByteInternal(addr)
}

func (memory *Memory) writeByte(addr uint16, b byte) {
	memory.contend(addr, 3)
	memory.writeByteInternal(addr, b)
}

func (memory *Memory) writeByteInternal(address uint16, b byte) {
	if (address >= 0x4000) && (address < 0x5800) && (memory.data[address] != b) {
		memory.screenBitmapWrite(address)
	}

	memory.data[address] = b
}

func (memory *Memory) screenBitmapWrite(address uint16) {
	// address: [0 1 0 y7 y6 y2 y1 y0 / y5 y4 y3 x4 x3 x2 x1 x0]
	var attr_x = (address & 0x001f)
	var attr_y = (((address & 0x0700) >> 8) | ((address & 0x00e0) >> 2) | ((address & 0x1800) >> 5)) / 8

	memory.dirtyScreen[attr_y*ScreenWidth_Attr+attr_x] = true
}

func (memory *Memory) getDirtyScreen() []bool {
	return &memory.dirtyScreen
}

func (memory *Memory) frame_begin() {
	for i := 0; i < (ScreenWidth_Attr * ScreenHeight_Attr); i++ {
		memory.dirtyScreen[i] = false
	}
}

func (memory *Memory) contend(address uint16, time uint) {
	tstates_p := &memory.z80.tstates

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

func (memory *Memory) set(address uint16, value byte) {
	memory.data[address] = value
}

func (memory *Memory) At(address uint) byte {
	return memory.data[address]
}

func (memory *Memory) Data() []byte {
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
			tstate_x := x / TSTATES_PER_PIXEL
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
