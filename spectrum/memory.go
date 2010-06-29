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

	set(addr uint, value byte)

	renderScreen()

	At(addr uint) byte
	Data() []byte
}

type Memory struct {
	Display DisplayAccessor
	data    [0x10000]byte
}

func (memory *Memory) readByteInternal(addr uint16) byte {
	return memory.data[addr]
}

func (memory *Memory) readByte(addr uint16) byte {
	tstates += 3
	return memory.readByteInternal(addr)
}

func (memory *Memory) writeByte(address uint16, b byte) {
	tstates += 3
	memory.writeByteInternal(address, b)
}

func (memory *Memory) writeByteInternal(address uint16, b byte) {

	switch {

	case (address > 0x4000) && (address < 0x5800):
		memory.drawScreenByte(address, b)

	case (address >= 0x5800) && (address < 0x5b00):
		memory.drawAttrByte(address, b)

	}

	memory.data[address] = b
}

func (memory *Memory) renderScreen() {
	var address uint16

	flashFrame = (flashFrame + 1) & 0x1f

	if (flashFrame & 0x0f) == 0 {
		// Need to redraw flashing attributes on this frame
		for address = 0x5800; address < 0x5b00; address++ {
			if (memory.data[address] & 0x80) != 0 {
				memory.drawAttrByte(address, memory.data[address])
			}
		}
	} else {
		memory.Display.flush()
	}
}

func (memory *Memory) contendRead(addr uint16, time uint) {
	tstates += time
}

func (memory *Memory) contendReadNoMreq(address uint16, time uint) {
	memory.contendRead(address, time)
}

func (memory *Memory) contendWriteNoMreq(address uint16, time uint) {
	tstates += time
}

func (memory *Memory) set(address uint, value byte) {
	memory.data[address] = value
}

// Draw a line of 8 pixels on the speccy display. The starting point
// of that line is encoded in address. The value byte specifies which
// pixels to turn on. Pixels turned on are painted with ink
// value. Pixels turned off are painted with paper value.
func (memory *Memory) drawScreenByte(address uint16, value byte) {

	// FIXME: Need to use proper type for x, y etc getting rid of
	// all these type conversions. I think uint16 is not needed at
	// all, byte type should be enough.

	// 0 1 0 y7 y6 y2 y1 y0 / y5 y4 y3 x4 x3 x2 x1 x0 
	var x uint16 = (address & 0x001f) // counted in characters
	var y uint16 = ((address & 0x0700) >> 8) | ((address & 0x00e0) >> 2) | ((address & 0x1800) >> 5) // counted in pixels
	var attributeByte byte = memory.data[0x5800|((y & 0xf8) << 2) | x]

	var ink [3]byte
	var paper [3]byte

	if ((attributeByte & 0x80) != 0) && ((flashFrame & 0x10) != 0) {
		// invert flashing attributes
		ink = palette[(attributeByte & 0x78) >> 3]
		paper = palette[((attributeByte & 0x40) >> 3)|(attributeByte & 0x07)]
	} else {
		ink = palette[((attributeByte & 0x40)>>3)|(attributeByte & 0x07)]
		paper = palette[(attributeByte & 0x78)>>3]
	}

	var p byte

	startx := x * 8

	for p = 0; p < 8; p++ {
		
		if (value & (1 << (7 - p))) != 0 {
			memory.Display.setPixel(uint(uint16(startx) + uint16(p)), uint(y), ink)
		} else {
			memory.Display.setPixel(uint(uint16(startx) + uint16(p)), uint(y), paper)
		}

	}
}

func (memory *Memory) drawAttrByte(address uint16, value byte) {
	/* 0 1 0 1 1 0 y4 y3 / y2 y1 y0 x4 x3 x2 x1 x0 */
	var x0 uint16 = (address & 0x001f)      /* counted in characters */
	var y0 uint16 = (address & 0x03e0) >> 2 /* counted in pixels */

	var ink [3]byte
	var paper [3]byte

	if ((value & 0x80) != 0) && (flashFrame&0x10) != 0 {
		/* invert flashing attributes */
		ink = palette[(value&0x78)>>3]
		paper = palette[((value&0x40)>>3)|(value&0x07)]
	} else {
		ink = palette[((value&0x40)>>3)|(value&0x07)]
		paper = palette[(value&0x78)>>3]
	}

	var y uint16
	
	startx := x0 * 8
	starty := y0

	for y = 0; y < 8; y++ {

//		var pixelAddress uint = ((uint(y0) | uint(y)) << 10) | (uint(x0) << 5)
		var screenByte byte = memory.data[0x4000|((y0 & 0xc0) << 5) |( y << 8) | ((y0 & 0x38) << 2) | x0]

		var p byte

		for p = 0; p < 8; p++ {
//			fmt.Printf("x %d y %d\n", uint(uint16(startx) + uint16(p)), uint(starty + y))
			if (screenByte & (1 << (7-p))) != 0 {
				memory.Display.setPixel(uint(uint16(startx) + uint16(p)), uint(starty + y), ink)
				// memory.Display.setPixelAt(uint(pixelAddress), ink)
				// pixelAddress += 4
			} else {
				memory.Display.setPixel(uint(uint16(startx) + uint16(p)), uint(starty + y), paper)
				// memory.Display.setPixelAt(uint(pixelAddress), paper)
				// pixelAddress += 4
			}
		}
	}
}

func (memory *Memory) At(address uint) byte {
	return memory.data[address]
}

func (memory *Memory) Data() []byte {
	return []byte(memory.Data())
}
