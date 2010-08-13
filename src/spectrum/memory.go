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

	getBorder() byte
	setBorder(borderColor byte)

	frame_begin()
	sendScreenToDisplay(display *DisplayInfo, borderEvents *BorderEvent)
	reset()

	At(addr uint) byte
        Data() []byte
}

type ula_byte_t struct {
	valid bool
	value uint8
}

type ula_attr_t struct {
	valid  bool
	value  uint8
	tstate uint
}

type Memory struct {
	data [0x10000]byte

	// Frame number
	frame uint

	borderColor byte

	// Screen bitmap data read by ULA, if they differ from data in memory at the end of a frame.
	// Spectrum y-coordinate.
	ula_bitmap [BytesPerLine * ScreenHeight]ula_byte_t

	// Screen attributes read by ULA, if they differ from data in memory at the end of a frame.
	// Linear y-coordinate.
	ula_attr [BytesPerLine * ScreenHeight]ula_attr_t

	// Whether the 8x8 rectangular screen area was modified during the current frame
	dirtyScreen [ScreenWidth_Attr * ScreenHeight_Attr]bool

	z80 *Z80
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
	if (address >= SCREEN_BASE_ADDR) && (address < ATTR_BASE_ADDR) {
		if memory.data[address] != b {
			memory.screenBitmapWrite(address)

			rel_addr := address - SCREEN_BASE_ADDR
			screenline_start_tstate := screenline_start_tstates[rel_addr>>BytesPerLine_log2]
			x, _ := screenAddr_to_xy(address)
			screen_tstate := screenline_start_tstate + uint(x>>PIXELS_PER_TSTATE_LOG2)
			if memory.z80.tstates > screen_tstate {
				// Remember the value read by ULA
				memory.ula_bitmap[rel_addr] = ula_byte_t{true, memory.data[address]}
			}
		}
	} else if (address >= ATTR_BASE_ADDR) && (address < 0x5b00) {
		memory.screenAttrWrite(address)

		attr_x := (address & 0x001f)
		attr_y := (address - ATTR_BASE_ADDR) / ScreenWidth_Attr

		x := uint(8 * attr_x)
		y := uint(8 * attr_y)
		for i := 0; i < 8; i++ {
			screenline_start_tstate := FIRST_SCREEN_BYTE + y*TSTATES_PER_LINE
			screen_tstate := screenline_start_tstate + (x >> PIXELS_PER_TSTATE_LOG2)
			if memory.z80.tstates > screen_tstate {
				ofs := (y << BytesPerLine_log2) + uint(attr_x)
				ula_attr := &memory.ula_attr[ofs]
				if !ula_attr.valid || (screen_tstate > ula_attr.tstate) {
					*ula_attr = ula_attr_t{true, memory.data[address], memory.z80.tstates}
				}
			}
			y++
		}
	}

	memory.data[address] = b
}

func (memory *Memory) screenBitmapWrite(address uint16) {
	// address: [0 1 0 y7 y6 y2 y1 y0 / y5 y4 y3 x4 x3 x2 x1 x0]
	var attr_x = (address & 0x001f)
	var attr_y = (((address & 0x0700) >> 8) | ((address & 0x00e0) >> 2) | ((address & 0x1800) >> 5)) / 8
	memory.dirtyScreen[attr_y*ScreenWidth_Attr+attr_x] = true
}

func (memory *Memory) screenAttrWrite(address uint16) {
	memory.dirtyScreen[address-ATTR_BASE_ADDR] = true
}

func (memory *Memory) getDirtyScreen() []bool {
	return &memory.dirtyScreen
}

// This function is called at the beginning of each frame
func (memory *Memory) frame_begin() {
	memory.frame++
	if memory.frame == 1 {
		// The very first frame --> repaint the whole screen
		for i := 0; i < ScreenWidth_Attr*ScreenHeight_Attr; i++ {
			memory.dirtyScreen[i] = true
		}
	} else {
		for i := 0; i < ScreenWidth_Attr*ScreenHeight_Attr; i++ {
			memory.dirtyScreen[i] = false
		}
	}

	ula_bitmap := &memory.ula_bitmap
	for ofs := uint16(0); ofs < BytesPerLine*ScreenHeight; ofs++ {
		if ula_bitmap[ofs].valid {
			memory.screenBitmapWrite(SCREEN_BASE_ADDR + ofs)
		}

		ula_bitmap[ofs].valid = false
	}

	ula_attr := &memory.ula_attr
	for ofs := uint16(0); ofs < BytesPerLine*ScreenHeight; ofs++ {
		if ula_attr[ofs].valid {
			linearY := (ofs >> BytesPerLine_log2)
			attr_y := (linearY >> 3)
			attr_x := (ofs & 0x001f)
			memory.screenAttrWrite(ATTR_BASE_ADDR + (attr_y << BytesPerLine_log2) + attr_x)
		}

		ula_attr[ofs].valid = false
	}
}

func (memory *Memory) reset() {
	for i := 0; i < 0x10000; i++ {
		memory.set(uint16(i), 0)
	}

	memory.frame = 0
}

func (memory *Memory) sendScreenToDisplay(display *DisplayInfo, borderEvents *BorderEvent) {
	sendDiffOnly := false
	if (display.lastFrame != nil) && (*display.lastFrame == memory.frame-1) {
		sendDiffOnly = true
	}

	var screen DisplayData
	{
		flash := (memory.frame & 0x10) != 0
		flash_previous := ((memory.frame - 1) & 0x10) != 0
		flash_diff := (flash != flash_previous)

		// screen.dirty
		if sendDiffOnly {
			screen.dirty = memory.dirtyScreen
		} else {
			for i := 0; i < ScreenWidth_Attr*ScreenHeight_Attr; i++ {
				screen.dirty[i] = true
			}
		}

		// Fill screen.bitmap & screen.attr, but only the dirty regions.
		memory_data := &memory.data
		memory_ulaBitmap := &memory.ula_bitmap
		memory_ulaAttr := &memory.ula_attr
		screen_dirty := &screen.dirty
		screen_bitmap := &screen.bitmap
		screen_attr := &screen.attr
		for attr_y := uint(0); attr_y < ScreenHeight_Attr; attr_y++ {
			attr_y8 := 8 * attr_y

			for attr_x := uint(0); attr_x < ScreenWidth_Attr; attr_x++ {
				attr_ofs := attr_y*ScreenWidth_Attr + attr_x

				// Make sure to send all changed flashing pixels to the DisplayReceiver
				if flash_diff {
					linearY_ofs := (attr_y8 << BytesPerLine_log2) + attr_x

					for y := 0; y < 8; y++ {
						var attr byte
						if !memory_ulaAttr[linearY_ofs].valid {
							attr = memory_data[ATTR_BASE_ADDR+attr_ofs]
						} else {
							attr = memory_ulaAttr[linearY_ofs].value
						}

						if (attr & 0x80) != 0 {
							screen_dirty[attr_ofs] = true
							break
						}

						linearY_ofs += BytesPerLine
					}
				}

				if !screen_dirty[attr_ofs] {
					continue
				}

				// screen.bitmap
				{
					screen_addr := xy_to_screenAddr(uint8(8*attr_x), uint8(attr_y8))
					linearY_ofs := (attr_y8 << BytesPerLine_log2) + attr_x

					for y := 0; y < 8; y++ {
						if !memory_ulaBitmap[screen_addr-SCREEN_BASE_ADDR].valid {
							screen_bitmap[linearY_ofs] = memory_data[screen_addr]
						} else {
							screen_bitmap[linearY_ofs] = memory_ulaBitmap[screen_addr-SCREEN_BASE_ADDR].value
						}

						screen_addr += 8 * BytesPerLine
						linearY_ofs += BytesPerLine
					}
				}

				// screen.attr
				{
					linearY_ofs := (attr_y8 << BytesPerLine_log2) + attr_x

					for y := 0; y < 8; y++ {
						var attr byte
						if !memory_ulaAttr[linearY_ofs].valid {
							attr = memory_data[ATTR_BASE_ADDR+attr_ofs]
						} else {
							attr = memory_ulaAttr[linearY_ofs].value
						}

						ink := ((attr & 0x40) >> 3) | (attr & 0x07)
						paper := (attr & 0x78) >> 3

						if flash && ((attr & 0x80) != 0) {
							/* invert flashing attributes */
							ink, paper = paper, ink
						}

						screen_attr[linearY_ofs] = attr_4bit((ink << 4) | paper)

						linearY_ofs += BytesPerLine
					}
				}
			}
		}

		// screen.border
		screen.border = memory.borderColor

		// screen.borderEvents
		if (borderEvents != nil) && (borderEvents.previous_orNil == nil) {
			// Only the one event which was added there at the start of the frame - ignore it
			screen.borderEvents = nil
		} else {
			screen.borderEvents = borderEvents
		}
	}

	if display.lastFrame == nil {
		display.lastFrame = new(uint)
	}
	*display.lastFrame = memory.frame

	display.displayReceiver.getDisplayDataChannel() <- &screen
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
