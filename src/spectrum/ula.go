package spectrum


type ula_byte_t struct {
	valid bool
	value uint8
}

type ula_attr_t struct {
	valid  bool
	value  uint8
	tstate uint
}

type ULA struct {
	// Frame number
	frame uint

	borderColor byte

	// Screen bitmap data read by ULA, if they differ from data in memory at the end of a frame.
	// Spectrum y-coordinate.
	bitmap [BytesPerLine * ScreenHeight]ula_byte_t

	// Screen attributes read by ULA, if they differ from data in memory at the end of a frame.
	// Linear y-coordinate.
	attr [BytesPerLine * ScreenHeight]ula_attr_t

	// Whether the 8x8 rectangular screen area was modified during the current frame
	dirtyScreen [ScreenWidth_Attr * ScreenHeight_Attr]bool

	speccy *Spectrum48k
}


func NewULA() *ULA {
	return &ULA{}
}

func (ula *ULA) reset() {
	ula.frame = 0
}


func (ula *ULA) getBorderColor() byte {
	return ula.borderColor
}

func (ula *ULA) setBorderColor(borderColor byte) {
	ula.borderColor = borderColor
}


// This function is called at the beginning of each frame
func (ula *ULA) frame_begin() {
	ula.frame++
	if ula.frame == 1 {
		// The very first frame --> repaint the whole screen
		for i := 0; i < ScreenWidth_Attr*ScreenHeight_Attr; i++ {
			ula.dirtyScreen[i] = true
		}
	} else {
		for i := 0; i < ScreenWidth_Attr*ScreenHeight_Attr; i++ {
			ula.dirtyScreen[i] = false
		}
	}

	bitmap := &ula.bitmap
	for ofs := uint16(0); ofs < BytesPerLine*ScreenHeight; ofs++ {
		if bitmap[ofs].valid {
			ula.screenBitmapTouch(SCREEN_BASE_ADDR + ofs)
		}

		bitmap[ofs].valid = false
	}

	attr := &ula.attr
	for ofs := uint16(0); ofs < BytesPerLine*ScreenHeight; ofs++ {
		if attr[ofs].valid {
			linearY := (ofs >> BytesPerLine_log2)
			attr_y := (linearY >> 3)
			attr_x := (ofs & 0x001f)
			ula.screenAttrTouch(ATTR_BASE_ADDR + (attr_y << BytesPerLine_log2) + attr_x)
		}

		attr[ofs].valid = false
	}
}

func (ula *ULA) screenBitmapTouch(address uint16) {
	// address: [0 1 0 y7 y6 y2 y1 y0 / y5 y4 y3 x4 x3 x2 x1 x0]
	var attr_x = (address & 0x001f)
	var attr_y = (((address & 0x0700) >> 8) | ((address & 0x00e0) >> 2) | ((address & 0x1800) >> 5)) / 8

	ula.dirtyScreen[attr_y*ScreenWidth_Attr+attr_x] = true
}

func (ula *ULA) screenAttrTouch(address uint16) {
	ula.dirtyScreen[address-ATTR_BASE_ADDR] = true
}

// Handle a write to an address in range (SCREEN_BASE_ADDR ... SCREEN_BASE_ADDR+0x1800-1)
func (ula *ULA) screenBitmapWrite(address uint16, b byte) {
	ula.screenBitmapTouch(address)

	rel_addr := address - SCREEN_BASE_ADDR
	screenline_start_tstate := screenline_start_tstates[rel_addr>>BytesPerLine_log2]
	x, _ := screenAddr_to_xy(address)
	screen_tstate := screenline_start_tstate + uint(x>>PIXELS_PER_TSTATE_LOG2)
	if ula.speccy.Cpu.tstates > screen_tstate {
		// Remember the value read by ULA
		ula.bitmap[rel_addr] = ula_byte_t{true, ula.speccy.Memory.data[address]}
	}
}

// Handle a write to an address in range (ATTR_BASE_ADDR ... ATTR_BASE_ADDR+0x300-1)
func (ula *ULA) screenAttrWrite(address uint16, b byte) {
	ula.screenAttrTouch(address)

	speccy := ula.speccy

	attr_x := (address & 0x001f)
	attr_y := (address - ATTR_BASE_ADDR) / ScreenWidth_Attr

	x := uint(8 * attr_x)
	y := uint(8 * attr_y)
	for i := 0; i < 8; i++ {
		screenline_start_tstate := FIRST_SCREEN_BYTE + y*TSTATES_PER_LINE
		screen_tstate := screenline_start_tstate + (x >> PIXELS_PER_TSTATE_LOG2)
		if speccy.Cpu.tstates > screen_tstate {
			ofs := (y << BytesPerLine_log2) + uint(attr_x)
			ula_attr := &ula.attr[ofs]
			if !ula_attr.valid || (screen_tstate > ula_attr.tstate) {
				*ula_attr = ula_attr_t{true, speccy.Memory.data[address], speccy.Cpu.tstates}
			}
		}
		y++
	}
}

func (ula *ULA) sendScreenToDisplay(display *DisplayInfo) {
	sendDiffOnly := false
	if (display.lastFrame != nil) && (*display.lastFrame == ula.frame-1) {
		sendDiffOnly = true
	}

	var screen DisplayData
	{
		flash := (ula.frame & 0x10) != 0
		flash_previous := ((ula.frame - 1) & 0x10) != 0
		flash_diff := (flash != flash_previous)

		// screen.dirty
		if sendDiffOnly {
			screen.dirty = ula.dirtyScreen
		} else {
			for i := 0; i < ScreenWidth_Attr*ScreenHeight_Attr; i++ {
				screen.dirty[i] = true
			}
		}

		// Fill screen.bitmap & screen.attr, but only the dirty regions.
		memory_data := &ula.speccy.Memory.data
		ula_bitmap := &ula.bitmap
		ula_attr := &ula.attr
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
						if !ula_attr[linearY_ofs].valid {
							attr = memory_data[ATTR_BASE_ADDR+attr_ofs]
						} else {
							attr = ula_attr[linearY_ofs].value
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
						if !ula_bitmap[screen_addr-SCREEN_BASE_ADDR].valid {
							screen_bitmap[linearY_ofs] = memory_data[screen_addr]
						} else {
							screen_bitmap[linearY_ofs] = ula_bitmap[screen_addr-SCREEN_BASE_ADDR].value
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
						if !ula_attr[linearY_ofs].valid {
							attr = memory_data[ATTR_BASE_ADDR+attr_ofs]
						} else {
							attr = ula_attr[linearY_ofs].value
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
		screen.border = ula.borderColor

		// screen.borderEvents
		borderEvents := ula.speccy.Ports.borderEvents
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
	*display.lastFrame = ula.frame

	display.displayReceiver.getDisplayDataChannel() <- &screen
}
