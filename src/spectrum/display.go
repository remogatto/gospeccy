/*

Copyright (c) 2010 Andrea Fazzi

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
"Software"), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

*/

package spectrum

import (
	"time"
)

const (
	ScreenWidth  = 256
	ScreenHeight = 192

	BytesPerLine      = ScreenWidth / 8 // =32
	BytesPerLine_log2 = 5               // =log2(BytesPerLine)

	ScreenWidth_Attr      = ScreenWidth / 8  // =32
	ScreenWidth_Attr_log2 = 5                // =log2(ScreenWidth_Attr)
	ScreenHeight_Attr     = ScreenHeight / 8 // =24

	ScreenBorderX = 32
	ScreenBorderY = 32

	// Screen dimensions, including the border
	TotalScreenWidth  = ScreenWidth + ScreenBorderX*2
	TotalScreenHeight = ScreenHeight + ScreenBorderY*2

	SCREEN_BASE_ADDR = 0x4000
	ATTR_BASE_ADDR   = 0x5800
)

// Spectrum 48k video timings
const (
	PIXELS_PER_TSTATE      = 2 // The number of screen pixels painted per T-state
	PIXELS_PER_TSTATE_LOG2 = 1 // = Log2(PIXELS_PER_TSTATE)

	// Horizontal
	LINE_SCREEN       = ScreenWidth / PIXELS_PER_TSTATE // 128 T states of screen
	LINE_RIGHT_BORDER = 24                              // 24 T states of right border
	LINE_RETRACE      = 48                              // 48 T states of horizontal retrace
	LINE_LEFT_BORDER  = 24                              // 24 T states of left border

	TSTATES_PER_LINE = (LINE_RIGHT_BORDER + LINE_SCREEN + LINE_LEFT_BORDER + LINE_RETRACE) // 224 T states

	FIRST_SCREEN_BYTE = 14336 // T-state when the first byte of the screen (16384) is displayed

	// Vertical
	LINES_TOP     = 64
	LINES_SCREEN  = ScreenHeight
	LINES_BOTTOM  = 56
	BORDER_TOP    = ScreenBorderY
	BORDER_BOTTOM = ScreenBorderY

	// The T-state which corresponds to pixel (0,0) on the host-machine display.
	// That pixel belongs to the border.
	DISPLAY_START            = (FIRST_SCREEN_BYTE - TSTATES_PER_LINE*BORDER_TOP - ScreenBorderX/PIXELS_PER_TSTATE + BORDER_TSTATE_ADJUSTMENT)
	BORDER_TSTATE_ADJUSTMENT = 2
)

type RGBA struct {
	R, G, B, A byte
}

func (color RGBA) value32() uint32 {
	return (uint32(color.A) << 24) | (uint32(color.R) << 16) | (uint32(color.G) << 8) | uint32(color.B)
}

var Palette [16]uint32 = [16]uint32{
	RGBA{000, 000, 000, 255}.value32(),
	RGBA{000, 000, 192, 255}.value32(),
	RGBA{192, 000, 000, 255}.value32(),
	RGBA{192, 000, 192, 255}.value32(),
	RGBA{000, 192, 000, 255}.value32(),
	RGBA{000, 192, 192, 255}.value32(),
	RGBA{192, 192, 000, 255}.value32(),
	RGBA{192, 192, 192, 255}.value32(),
	RGBA{000, 000, 000, 255}.value32(),
	RGBA{000, 000, 255, 255}.value32(),
	RGBA{255, 000, 000, 255}.value32(),
	RGBA{255, 000, 255, 255}.value32(),
	RGBA{000, 255, 000, 255}.value32(),
	RGBA{000, 255, 255, 255}.value32(),
	RGBA{255, 255, 000, 255}.value32(),
	RGBA{255, 255, 255, 255}.value32(),
}

func screenAddr_to_xy(screenAddr uint16) (x, y uint8) {
	// address: [0 1 0 y7 y6 y2 y1 y0 / y5 y4 y3 x4 x3 x2 x1 x0]
	x = uint8((screenAddr & 0x001f) << 3)
	y = uint8(((screenAddr & 0x0700) >> 8) | ((screenAddr & 0x00e0) >> 2) | ((screenAddr & 0x1800) >> 5))
	return
}

func screenAddr_to_attrXY(screenAddr uint16) (attr_x, attr_y uint8) {
	// address: [0 1 0 y7 y6 y2 y1 y0 / y5 y4 y3 x4 x3 x2 x1 x0]
	attr_x = uint8(screenAddr & 0x001f)
	attr_y = uint8((((screenAddr & 0x0700) >> 8) | ((screenAddr & 0x00e0) >> 2) | ((screenAddr & 0x1800) >> 5)) >> 3)
	return
}

func xy_to_screenAddr(x, y uint8) uint16 {
	yy := uint(y)
	addr_y := SCREEN_BASE_ADDR | 0x800*(yy>>6) | BytesPerLine*((yy&0x38)>>3) | ((yy & 0x07) << 8)
	return uint16(addr_y | uint(x>>3))
}

// The lower 4 bits define the paper, the higher 4 bits define the ink.
// Note that the paper is in the *lower* half.
// There is no flash bit.
type Attr_4bit byte

// This is the primary structure for sending display changes
// from the Z80 CPU emulation core to a rendering backend.
// The data is already preprocessed, to make the rendering-backend's code simpler and faster.
//
// The content of 'bitmap' and 'attr' corresponding to non-dirty regions is unspecified.
type DisplayData struct {
	Bitmap [BytesPerLine * ScreenHeight]byte          // Linear y-coordinate
	Attr   [BytesPerLine * ScreenHeight]Attr_4bit     // Linear y-coordinate
	Dirty  [ScreenWidth_Attr * ScreenHeight_Attr]bool // The 8x8 rectangular region was modified, either the bitmap or the attr

	BorderEvents []BorderEvent

	// From structure Cmd_RenderFrame
	CompletionTime_orNil chan<- time.Time
}

// Interface to a rendering backend awaiting display changes
type DisplayReceiver interface {
	GetDisplayDataChannel() chan<- *DisplayData

	// Closes the display associated with this DisplayReceiver
	Close()
}

// Let 'addr' be in range 0x4000 ... 0x5800-1.
// Then 'screenline_start_tstates[(addr-0x4000)/BytesPerLine]' is the T-state when the Spectrum
// starts painting the screenline containing 'addr'.
var screenline_start_tstates [ScreenHeight]int

func init() {
	for y := uint8(0); y < ScreenHeight; y++ {
		addr := xy_to_screenAddr(0, y)
		screenline_start_tstates[(addr-SCREEN_BASE_ADDR)/BytesPerLine] = FIRST_SCREEN_BYTE + int(y)*TSTATES_PER_LINE
	}
}

func init() {
	// Some sanity checks
	Assert(ScreenBorderX <= LINE_RIGHT_BORDER*PIXELS_PER_TSTATE)
	Assert(ScreenBorderY <= LINE_RIGHT_BORDER*PIXELS_PER_TSTATE)
	Assert(ScreenBorderY <= LINES_TOP)
	Assert(ScreenBorderY <= LINES_BOTTOM)
}
