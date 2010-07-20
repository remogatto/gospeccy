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

type RGBA struct {
	R,G,B,A byte
}

var palette [16]RGBA = [16]RGBA{
	RGBA{0  , 0  , 0  , 255},
	RGBA{0  , 0  , 192, 255},
	RGBA{192, 0  , 0  , 255},
	RGBA{192, 0  , 192, 255},
	RGBA{0  , 192, 0  , 255},
	RGBA{0  , 192, 192, 255},
	RGBA{192, 192, 0  , 255},
	RGBA{192, 192, 192, 255},
	RGBA{0  , 0  , 0  , 255},
	RGBA{0  , 0  , 255, 255},
	RGBA{255, 0  , 0  , 255},
	RGBA{255, 0  , 255, 255},
	RGBA{0  , 255, 0  , 255},
	RGBA{0  , 255, 255, 255},
	RGBA{255, 255, 0  , 255},
	RGBA{255, 255, 255, 255}}

func (color RGBA) value32() uint32 {
	return (uint32(color.A) << 24) | (uint32(color.R) << 16) | (uint32(color.G) << 8) | uint32(color.B)
}

type SurfaceAccessor interface {
	Width() uint
	Height() uint
	SizeInBytes() uint
	Bpp() uint

	setValueAt(offset uint, value uint32)
	setPixelAt(offset uint, color RGBA)
	
	setValue(x, y uint, value uint32)
	setPixel(x, y uint, color RGBA)
}

type DisplayChannel interface {
	getScreenChannel() chan *Screen
}

