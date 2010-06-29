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
	"unsafe"
	"sdl"
)

type SDLSurface struct {
	Surface *sdl.Surface
}

func (s *SDLSurface) Width() uint {
	return uint(s.Surface.W)
}

func (s *SDLSurface) Height() uint {
	return uint(s.Surface.H)
}

func (s *SDLSurface) Bpp() uint {
	return uint(s.Surface.Format.BytesPerPixel)
}

func (s *SDLSurface) SizeInBytes() uint {
	return s.Width() * s.Height() * s.Bpp()
}

func (s *SDLSurface) setPixelAt(address uint, color [3]byte) {
	// FIXME: Need setting alpha channel to max opacity!
	s.setValueAt(address, (uint32(color[0]) << 16) | (uint32(color[1]) << 8) | uint32(color[2]))
}

func (s *SDLSurface) setPixel(x, y uint, color [3]byte) {
	var bpp, pitch = uint(s.Surface.Format.BytesPerPixel), uint(s.Surface.Pitch)
	
	// FIXME: Need setting alpha channel to max opacity!
	s.setValueAt(y * pitch + x*bpp, (uint32(color[0]) << 16) | (uint32(color[1]) << 8) | uint32(color[2]))
}

func (s *SDLSurface) setPixelValue(x, y uint, value uint32) {
	var bpp, pitch = uint(s.Surface.Format.BytesPerPixel), uint(s.Surface.Pitch)
	
	// FIXME: Need setting alpha channel to max opacity!
	s.setValueAt(y * pitch + x*bpp, value)
}

func (s *SDLSurface) setValueAt(id uint, value uint32) {
	// FIXME: Need locking?
	var pixel = uintptr(unsafe.Pointer(s.Surface.Pixels))

	pixel += uintptr(id)

	var p =(*uint32)(unsafe.Pointer(pixel))
	
	*p = value
}

func (s *SDLSurface) getValueAt(id uint) uint32 {
	// FIXME: Need locking?
	var pixel = uintptr(unsafe.Pointer(s.Surface.Pixels))

	pixel += uintptr(id)

	var p =(*uint32)(unsafe.Pointer(pixel))
	
	return *p
}

type SDLDisplay struct {
	// The whole screen borders included
	ScreenSurface *SDLSurface

	// The drawable display without borders
	DisplaySurface *SDLSurface

	borderColor [3]byte
}

func NewSDLDisplay(screenSurface *sdl.Surface) *SDLDisplay {

	// Here below we create the internal DisplaySurface i.e. the
	// drawable display area without borders
	displaySurface := &SDLSurface { sdl.CreateRGBSurface(sdl.SWSURFACE, 256, 192, 32, 0, 0, 0, 0) }

	// Literal initialization of the SDLDisplay object (the whole
	// screen: drawable area + borders)
	return &SDLDisplay{ ScreenSurface: &SDLSurface{ screenSurface }, DisplaySurface: displaySurface }
}

func (display *SDLDisplay) flush() {
	color := (uint32(display.borderColor[0]) << 16) | (uint32(display.borderColor[1]) << 8) | uint32(display.borderColor[2])

	display.ScreenSurface.Surface.FillRect(nil, color)
	display.ScreenSurface.Surface.Blit(&sdl.Rect{32, 24, 320 - 32, 240 - 24}, display.DisplaySurface.Surface, &sdl.Rect{0, 0, 256, 192})
}

func (display *SDLDisplay) setBorderColor(color [3]byte) {
	display.borderColor = color
}

func (display *SDLDisplay) setPixelAt(address uint, color [3]byte) {
	display.DisplaySurface.setPixelAt(address, color)
}

func (display *SDLDisplay) setPixel(x, y uint, color [3]byte) {
	display.DisplaySurface.setPixel(x, y, color)
}

// Experimental SDLDoubledDisplay
type SDLDoubledDisplay struct {
	SDLDisplay
}

func NewSDLDoubledDisplay(screenSurface *sdl.Surface) *SDLDoubledDisplay {
	// Here below we create the internal DisplaySurface i.e. the
	// drawable display area without borders
	displaySurface := &SDLSurface { sdl.CreateRGBSurface(sdl.SWSURFACE, 512, 384, 32, 0, 0, 0, 0) }

	// Literal initialization of the SDLDisplay object (the whole
	// screen: drawable area + borders)
	return &SDLDoubledDisplay{ SDLDisplay{ &SDLSurface{ screenSurface }, displaySurface, [3]byte { 0, 0, 0 } } }
}

func (display *SDLDoubledDisplay) flush() {
	color := (uint32(display.borderColor[0]) << 16) | (uint32(display.borderColor[1]) << 8) | uint32(display.borderColor[2])

	display.ScreenSurface.Surface.FillRect(nil, color)
	display.ScreenSurface.Surface.Blit(&sdl.Rect{32 * 2, 24 * 2, 640 - 32*2, 480 - 24*2}, display.DisplaySurface.Surface, &sdl.Rect{0, 0, 512, 384 })
}

func (display *SDLDoubledDisplay) setBorderColor(color [3]byte) {
	display.borderColor = color
}

func (display *SDLDoubledDisplay) setPixel(x, y uint, color [3]byte) {
	var (
		scaleX uint = x * 2
		scaleY uint = y * 2
		scaleXInc uint = scaleX + 1
 		scaleYInc uint = scaleY + 1
	)

	display.DisplaySurface.setPixel(scaleX, scaleY, color)
	display.DisplaySurface.setPixel(scaleXInc, scaleYInc, color)
	display.DisplaySurface.setPixel(scaleX, scaleYInc, color)
	display.DisplaySurface.setPixel(scaleXInc, scaleY, color)
}

