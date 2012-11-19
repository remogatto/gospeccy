package z80

type MemoryAccessor interface {
	readByte(address uint16) byte
	readByteInternal(address uint16) byte

	writeByte(address uint16, value byte)
	writeByteInternal(address uint16, value byte)

	contendRead(address uint16, time uint)
	contendReadNoMreq(address uint16, time uint)
	contendReadNoMreq_loop(address uint16, time uint, count uint)

	contendWriteNoMreq(address uint16, time uint)
	contendWriteNoMreq_loop(address uint16, time uint, count uint)

	Read(address uint16) byte
	Write(address uint16, value byte, protectROM bool)
	Data() *[0x10000]byte
}
