package z80

type MemoryAccessor interface {
	ReadByte(address uint16) byte
	ReadByteInternal(address uint16) byte

	WriteByte(address uint16, value byte)
	WriteByteInternal(address uint16, value byte)

	ContendRead(address uint16, time uint)
	ContendReadNoMreq(address uint16, time uint)
	ContendReadNoMreq_loop(address uint16, time uint, count uint)

	ContendWriteNoMreq(address uint16, time uint)
	ContendWriteNoMreq_loop(address uint16, time uint, count uint)

	Read(address uint16) byte
	Write(address uint16, value byte, protectROM bool)
	Data() *[0x10000]byte
}
