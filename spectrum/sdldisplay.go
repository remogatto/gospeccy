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
