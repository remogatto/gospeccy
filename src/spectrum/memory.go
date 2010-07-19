package spectrum

// import "fmt"

var flashFrame byte

type MemoryAccessor interface {
	readByte(address uint16) byte
	readByteInternal(address uint16) byte

	writeByte(address uint16, value byte)
	writeByteInternal(address uint16, value byte)

	contendRead(addr uint16, time uint)
	contendReadNoMreq(addr uint16, time uint)

	contendWriteNoMreq(addr uint16, time uint)

	set(addr uint16, value byte)
	
	getBorder() RGBA
	setBorder(borderColor RGBA)

	sendScreenToDisplay(display DisplayChannel, borderEvents *BorderEvent)

	At(addr uint) byte
	Data() []byte
}

type PaperInk [2]RGBA

func equals(a,b PaperInk) bool {
	return	(a[0].R == b[0].R) &&
			(a[0].G == b[0].G) &&
			(a[0].B == b[0].B) &&
			(a[0].A == b[0].A) &&
			(a[1].R == b[1].R) &&
			(a[1].G == b[1].G) &&
			(a[1].B == b[1].B) &&
			(a[1].A == b[1].A)
}

type Screen struct {
	bitmap       [ScreenWidth/8*ScreenHeight] byte
	attr         [ScreenWidth_Attr*ScreenHeight_Attr] PaperInk
	border       RGBA
	borderEvents *BorderEvent	// Might be nil
	flash        bool
}

type Memory struct {
	data [0x10000]byte
	borderColor RGBA
	z80 *Z80
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

func (memory *Memory) getBorder() RGBA {
	return memory.borderColor
}

func (memory *Memory) setBorder(borderColor RGBA) {
	memory.borderColor = borderColor
}

func (memory *Memory) sendScreenToDisplay(display DisplayChannel, borderEvents *BorderEvent) {
	flashFrame = (flashFrame + 1) & 0x1f
	
	var screen Screen
	{
		// screen.bitmap
		for ofs := 0; ofs < ScreenWidth/8*ScreenHeight; ofs++ {
			screen.bitmap[ofs] = memory.data[0x4000+ofs]
		}
		
		// screen.flash
		flash := (flashFrame & 0x10 != 0)
		screen.flash = flash
		
		// screen.attr
		for attr_ofs := 0; attr_ofs < ScreenWidth_Attr*ScreenHeight_Attr; attr_ofs++ {
			attr := memory.data[0x5800+attr_ofs]
				
			var ink RGBA
			var paper RGBA
				
			if flash && ((attr & 0x80) != 0) {
				/* invert flashing attributes */
				ink = palette[(attr&0x78)>>3]
				paper = palette[((attr&0x40)>>3)|(attr&0x07)]
			} else {
				ink = palette[((attr&0x40)>>3)|(attr&0x07)]
				paper = palette[(attr&0x78)>>3]
			}
			
			screen.attr[attr_ofs] = PaperInk{paper, ink}
		}
		
		// screen.border
		screen.border = memory.borderColor
		
		// screen.borderEvents
		if (borderEvents != nil) && (borderEvents.previous_orNil==nil) {
			// Only the one event which was added there at the start of the frame - ignore it
			screen.borderEvents = nil
		} else {
			screen.borderEvents = borderEvents
		}
	}
	
	display.getScreenChannel() <- &screen
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
