package spectrum

import (
	"unsafe"
	"sdl"
//	"image"
//	"fmt"
)

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

type DisplayAccessor interface {
	setPixel(address uint, color [3]byte)
	setBorderColor(color [3]byte)
	flip()
}

type SDLDisplay struct {
	// The whole screen borders included
	ScreenSurface *sdl.Surface

	// The drawable display without borders
	DisplaySurface *sdl.Surface

	borderColor [3]byte
}

func NewSDLDisplay(screenSurface *sdl.Surface) *SDLDisplay {
	displaySurface := sdl.CreateRGBSurface(sdl.SWSURFACE, 256, 192, 32, 0, 0, 0, 0)
	return &SDLDisplay{ ScreenSurface: screenSurface, DisplaySurface: displaySurface }
}

func (display *SDLDisplay) flip() {
	color := (uint32(display.borderColor[0]) << 16) | (uint32(display.borderColor[1]) << 8) | uint32(display.borderColor[2])

	display.ScreenSurface.FillRect(nil, color)
	display.ScreenSurface.Blit(&sdl.Rect{32, 24, 320 - 32, 240 - 24}, display.DisplaySurface, &sdl.Rect{0, 0, 256, 192})
	display.ScreenSurface.Flip()
}

func (display *SDLDisplay) setBorderColor(color [3]byte) {
	display.borderColor = color
}

func (display *SDLDisplay) setPixel(address uint, color [3]byte) {
	// FIXME: Need locking?
	var pixel = uintptr(unsafe.Pointer(display.DisplaySurface.Pixels))

	pixel += uintptr(address)

//	var p = (*image.RGBAColor)(unsafe.Pointer(pixel))
	var p =(*uint32)(unsafe.Pointer(pixel))
	*p = (uint32(color[0]) << 16) | (uint32(color[1]) << 8) | uint32(color[2])
	// p.R = uint8(color[0])
	// p.G = uint8(color[1])
	// p.B = uint8(color[2])

//	fmt.Printf("R:%d G:%d B:%d\n", color[0], color[1], color[2])

	// Alpha value set to max opacity
	// p.A = uint8(255)
}
