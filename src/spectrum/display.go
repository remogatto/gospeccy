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
	ScreenWidth  = 256
	ScreenHeight = 192

	ScreenWidth_Attr  = ScreenWidth / 8  // =32
	ScreenHeight_Attr = ScreenHeight / 8 // =24

	ScreenBorderX = 32
	ScreenBorderY = 24

	// Screen dimensions, including the border
	TotalScreenWidth  = ScreenWidth + ScreenBorderX*2
	TotalScreenHeight = ScreenHeight + ScreenBorderY*2
)

type RGBA struct {
	R, G, B, A byte
}

var palette [16]RGBA = [16]RGBA{
	RGBA{0, 0, 0, 255},
	RGBA{0, 0, 192, 255},
	RGBA{192, 0, 0, 255},
	RGBA{192, 0, 192, 255},
	RGBA{0, 192, 0, 255},
	RGBA{0, 192, 192, 255},
	RGBA{192, 192, 0, 255},
	RGBA{192, 192, 192, 255},
	RGBA{0, 0, 0, 255},
	RGBA{0, 0, 255, 255},
	RGBA{255, 0, 0, 255},
	RGBA{255, 0, 255, 255},
	RGBA{0, 255, 0, 255},
	RGBA{0, 255, 255, 255},
	RGBA{255, 255, 0, 255},
	RGBA{255, 255, 255, 255}}

func (color RGBA) value32() uint32 {
	return (uint32(color.A) << 24) | (uint32(color.R) << 16) | (uint32(color.G) << 8) | uint32(color.B)
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

type Display struct {
	// Shared VRAM
	memory []byte
	borderColor RGBA
}

type DecodedDisplay struct {
	borderColor RGBA
	borderEvents *BorderEvent // Might be nil
	flash        bool

	bitmap       [ScreenWidth/8*ScreenHeight] byte
	attr         [ScreenWidth_Attr*ScreenHeight_Attr] PaperInk	
}

type DisplayChannel interface {
	getScreenChannel() chan *DecodedDisplay
}

func NewDisplay(systemMemory []byte) *Display {
	return &Display{ memory: systemMemory[0x4000:0x5b00] }
}

func (display *Display) getBorderColor() RGBA { return display.borderColor }
func (display *Display) setBorderColor(color RGBA) { display.borderColor = color }

func (display *Display) decode() *DecodedDisplay {
	
	var decodedDisplay DecodedDisplay

	flashFrame = (flashFrame + 1) & 0x1f
	
	// screen.bitmap
	for ofs := 0; ofs < ScreenWidth/8*ScreenHeight; ofs++ {
		decodedDisplay.bitmap[ofs] = display.memory[ofs]
	}
	
	// screen.flash
	flash := (flashFrame & 0x10 != 0)
	decodedDisplay.flash = flash
	
	// screen.attr
	for attr_ofs := 0; attr_ofs < ScreenWidth_Attr*ScreenHeight_Attr; attr_ofs++ {
		attr := display.memory[(0x5800-0x4000)+attr_ofs]
		
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
		
		decodedDisplay.attr[attr_ofs] = PaperInk{paper, ink}
	}

	decodedDisplay.borderColor = display.getBorderColor()

	return &decodedDisplay

}

func (display *Display) sendDisplayData(ch DisplayChannel, borderEvents *BorderEvent) {

	decodedDisplay := display.decode()

	// screen.borderEvents
	if (borderEvents != nil) && (borderEvents.previous_orNil == nil) {
		// Only the one event which was added there at the start of the frame - ignore it
		decodedDisplay.borderEvents = nil
	} else {
		decodedDisplay.borderEvents = borderEvents
	}
	
	ch.getScreenChannel() <- decodedDisplay

}
