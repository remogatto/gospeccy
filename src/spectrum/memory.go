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
}

type Memory struct {
	data [0x10000]byte
	z80  *Z80
}

func NewMemory() *Memory {
	return &Memory{}
}

func (memory *Memory) readByteInternal(addr uint16) byte {
	return memory.data[addr]
}

func (memory *Memory) readByte(addr uint16) byte {
	memory.z80.tstates += 3
	return memory.readByteInternal(addr)
}

func (memory *Memory) writeByte(address uint16, b byte) {
	memory.z80.tstates += 3
	memory.writeByteInternal(address, b)
}

func (memory *Memory) writeByteInternal(address uint16, b byte) {
	memory.data[address] = b
}

func (memory *Memory) contendRead(addr uint16, time uint) {
	memory.z80.tstates += time
}

func (memory *Memory) contendReadNoMreq(address uint16, time uint) {
	memory.contendRead(address, time)
}

func (memory *Memory) contendWriteNoMreq(address uint16, time uint) {
	memory.z80.tstates += time
}

func (memory *Memory) set(address uint16, value byte) {
	memory.data[address] = value
}

func (memory *Memory) At(address uint) byte {
	return memory.data[address]
}

func (memory *Memory) Data() []byte {
	return []byte(memory.Data())
}
