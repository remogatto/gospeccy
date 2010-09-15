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

	// Whether to discern between [data read by ULA] and [data in memory at the end of a frame].
	// If the value is 'false', then fields 'bitmap' and 'attr' will contain no information.
	// The default value is 'true'.
	accurateEmulation bool

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
	return &ULA{accurateEmulation: true}
}

func (ula *ULA) init(speccy *Spectrum48k) {
	ula.speccy = speccy
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


func (ula *ULA) setEmulationAccuracy(accurateEmulation bool) {
	ula.accurateEmulation = accurateEmulation
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
func (ula *ULA) screenBitmapWrite(address uint16, oldValue byte, newValue byte) {
	if oldValue != newValue {
		ula.screenBitmapTouch(address)

		if ula.accurateEmulation {
			rel_addr := address - SCREEN_BASE_ADDR
			ula_lineStart_tstate := screenline_start_tstates[rel_addr>>BytesPerLine_log2]
			x, _ := screenAddr_to_xy(address)
			ula_tstate := ula_lineStart_tstate + uint(x>>PIXELS_PER_TSTATE_LOG2)
			if ula_tstate <= ula.speccy.Cpu.tstates {
				// Remember the value read by ULA
				ula.bitmap[rel_addr] = ula_byte_t{true, oldValue}
			}
		}
	}
}

// Handle a write to an address in range (ATTR_BASE_ADDR ... ATTR_BASE_ADDR+0x300-1)
func (ula *ULA) screenAttrWrite(address uint16, oldValue byte, newValue byte) {
	if oldValue != newValue {
		ula.screenAttrTouch(address)

		if ula.accurateEmulation {
			speccy := ula.speccy

			attr_x := uint(address & 0x001f)
			attr_y := uint((address - ATTR_BASE_ADDR) >> ScreenWidth_Attr_log2)

			x := 8 * attr_x
			y := 8 * attr_y

			ofs := (y << BytesPerLine_log2) + attr_x
			ula_tstate := FIRST_SCREEN_BYTE + y*TSTATES_PER_LINE + (x >> PIXELS_PER_TSTATE_LOG2)

			for i := 0; i < 8; i++ {
				if ula_tstate <= speccy.Cpu.tstates {
					ula_attr := &ula.attr[ofs]
					if !ula_attr.valid || (ula_tstate > ula_attr.tstate) {
						*ula_attr = ula_attr_t{true, oldValue, speccy.Cpu.tstates}
					}
					ofs += BytesPerLine
					ula_tstate += TSTATES_PER_LINE
				} else {
					break
				}
			}
		}
	}
}

func (ula *ULA) prepare(display *DisplayInfo) *DisplayData {
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
		var memory_data *[0x10000]byte = ula.speccy.Memory.Data()
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
		borderEvents := ula.speccy.Ports.getBorderEvents()
		if (borderEvents != nil) && (borderEvents.previous_orNil == nil) {
			// Only the one event which was added there at the start of the frame - ignore it
			screen.borderEvents = nil
		} else {
			screen.borderEvents = borderEvents
		}
	}

	return &screen
}

func (ula *ULA) sendScreenToDisplay(display *DisplayInfo, completionTime_orNil chan<- int64) {
	displayData := ula.prepare(display)
	displayData.completionTime_orNil = completionTime_orNil

	displayChannel := display.displayReceiver.getDisplayDataChannel()
	nonBlockingSend := displayChannel <- displayData

	if nonBlockingSend {
		if display.lastFrame == nil {
			display.lastFrame = new(uint)
		}
		*display.lastFrame = ula.frame
	} else {
		// Throw away the frame since the send would block.
		// This allows the CPU emulation to proceed when the next tick arrives,
		// instead of waiting for the display backend to receive the previous frame.
		// Note that 'display.lastFrame' is NOT updated.
		display.numMissedFrames++
	}
}
