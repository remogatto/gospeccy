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

var palette [16][3]byte = [16][3]byte{
	[3]byte{0, 0, 0},
	[3]byte{0, 0, 192},
	[3]byte{192, 0, 0},
	[3]byte{192, 0, 192},
	[3]byte{0, 192, 0},
	[3]byte{0, 192, 192},
	[3]byte{192, 192, 0},
	[3]byte{192, 192, 192},
	[3]byte{0, 0, 0},
	[3]byte{0, 0, 255},
	[3]byte{255, 0, 0},
	[3]byte{255, 0, 255},
	[3]byte{0, 255, 0},
	[3]byte{0, 255, 255},
	[3]byte{255, 255, 0},
	[3]byte{255, 255, 255}}

type SurfaceAccessor interface {
	Width() uint
	Height() uint
	SizeInBytes() uint
	Bpp() uint

	getValueAt(id uint) uint32
	setValueAt(id uint, value uint32)

	setPixelAt(address uint, color [3]byte)
	setPixelValue(x, y uint, value uint32)
	setPixel(x, y uint, color [3]byte)
}


type DisplayAccessor interface {

	// FIXME: ZX Spectrum display coords should be of byte size!
	setPixel(x, y uint, color [3]byte)

	setPixelAt(address uint, color [3]byte)
	setBorderColor(color [3]byte)
	flush()
}

