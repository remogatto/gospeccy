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

const (
	ScreenWidth = 256
	ScreenHeight = 192

	ScreenWidth_Attr = ScreenWidth/8	// =32
	ScreenHeight_Attr = ScreenHeight/8	// =24

	ScreenBorderX = 32
	ScreenBorderY = 24

	// Screen dimensions, including the border
	TotalScreenWidth = ScreenWidth + ScreenBorderX * 2
	TotalScreenHeight = ScreenHeight + ScreenBorderY * 2
)

// Spectrum 48k video timings
const (
	TSTATES_PER_PIXEL = 2

	// Horizontal
	LINE_SCREEN       = ScreenWidth/TSTATES_PER_PIXEL	// 128 T states of screen
	LINE_RIGHT_BORDER = 24		// 24 T states of right border
	LINE_RETRACE      = 48		// 48 T states of horizontal retrace
	LINE_LEFT_BORDER  = 24		// 24 T states of left border
	
	TSTATES_PER_LINE  = (LINE_RIGHT_BORDER + LINE_SCREEN + LINE_LEFT_BORDER + LINE_RETRACE)	// 224 T states
	
	FIRST_SCREEN_BYTE = 14336	// T states before the first byte of the screen (16384) is displayed
	
	// Vertical
	LINES_TOP         = 64
	LINES_SCREEN      = ScreenHeight
	LINES_BOTTOM      = 56
	BORDER_TOP        = ScreenBorderY
	BORDER_BOTTOM     = ScreenBorderY
	
	// The T-state which corresponds to pixel (0,0) on the (SDL) surface.
	// That pixel belongs to the border.
	DISPLAY_START     = ( FIRST_SCREEN_BYTE - TSTATES_PER_LINE*BORDER_TOP - ScreenBorderX/TSTATES_PER_PIXEL )
)



type RGBA struct {
	R,G,B,A byte
}

func (color RGBA) value32() uint32 {
	return (uint32(color.A) << 24) | (uint32(color.R) << 16) | (uint32(color.G) << 8) | uint32(color.B)
}

var palette [16]uint32 = [16]uint32{
	RGBA{0  , 0  , 0  , 255}.value32(),
	RGBA{0  , 0  , 192, 255}.value32(),
	RGBA{192, 0  , 0  , 255}.value32(),
	RGBA{192, 0  , 192, 255}.value32(),
	RGBA{0  , 192, 0  , 255}.value32(),
	RGBA{0  , 192, 192, 255}.value32(),
	RGBA{192, 192, 0  , 255}.value32(),
	RGBA{192, 192, 192, 255}.value32(),
	RGBA{0  , 0  , 0  , 255}.value32(),
	RGBA{0  , 0  , 255, 255}.value32(),
	RGBA{255, 0  , 0  , 255}.value32(),
	RGBA{255, 0  , 255, 255}.value32(),
	RGBA{0  , 255, 0  , 255}.value32(),
	RGBA{0  , 255, 255, 255}.value32(),
	RGBA{255, 255, 0  , 255}.value32(),
	RGBA{255, 255, 255, 255}.value32()}

type DisplayChannel interface {
	getScreenChannel() chan *Screen
}

