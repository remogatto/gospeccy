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
	"sdl"
	"unsafe"
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

func (s *SDLSurface) setValueAt(offset uint, value uint32) {
	// FIXME: Need locking?
	var pixel = uintptr(unsafe.Pointer(s.Surface.Pixels))

	pixel += uintptr(offset)

	var p =(*uint32)(unsafe.Pointer(pixel))
	
	*p = value
}

func (s *SDLSurface) setPixelAt(offset uint, color RGBA) {
	s.setValueAt(offset, color.value32())
}

func (s *SDLSurface) setValue(x, y uint, value uint32) {
	var bpp, pitch = uint(s.Surface.Format.BytesPerPixel), uint(s.Surface.Pitch)
	
	s.setValueAt(y * pitch + x*bpp, value)
}

func (s *SDLSurface) setPixel(x, y uint, color RGBA) {
	var bpp, pitch = uint(s.Surface.Format.BytesPerPixel), uint(s.Surface.Pitch)
	
	s.setPixelAt(y * pitch + x*bpp, color)
}

type SDLScreen struct {
	// Channel for receiving screen data
	screenChannel  chan *Screen
	
	// The whole screen, borders included
	ScreenSurface SDLSurface
}

type _ScreenRenderer interface {
	render(screen, oldScreen_orNil *Screen)
}

func NewSDLScreen(app *Application, screenSurface *sdl.Surface) *SDLScreen {
	SDL_screen := &SDLScreen{ make(chan *Screen), SDLSurface{ screenSurface } }
	
	go screenRenderLoop(app.NewEventLoop(), SDL_screen.screenChannel, SDL_screen)
	
	return SDL_screen
}

// Implement DisplayChannel
func (display *SDLScreen) getScreenChannel() chan *Screen {
	return display.screenChannel
}

func screenRenderLoop(evtLoop *EventLoop, screenChannel chan *Screen, renderer _ScreenRenderer) {
	var screen    *Screen = nil
	var oldScreen *Screen = nil
	for {
		select {
			case <-evtLoop.Pause:
				evtLoop.Pause <- 0

			case <-evtLoop.Terminate:
				// Terminate this Go routine
				if evtLoop.App.Verbose { println("screen render loop: exit") }
				evtLoop.Terminate <- 0
				return
				
			case screen = <- screenChannel:
				renderer.render(screen, oldScreen)
				oldScreen = screen
				screen = nil
		}
	}
}

func renderBorder(surface *sdl.Surface, scale uint, screen *Screen, oldScreen_orNil *Screen) {
	borderValue := screen.border.value32()
	
	if (oldScreen_orNil == nil) || (borderValue != oldScreen_orNil.border.value32()) ||
	   (screen.borderEvents != nil) || (oldScreen_orNil.borderEvents != nil) {
		s := scale
		
		const W = ScreenWidth
		const H = ScreenHeight
		const BW = ScreenBorderX
		const BH = ScreenBorderY
		const TW = TotalScreenWidth

		surface.FillRect( &sdl.Rect{int16(s*0)     , int16(s*0)     , uint16(s*TW), uint16(s*BH)}, borderValue )
		surface.FillRect( &sdl.Rect{int16(s*0)     , int16(s*(BH+H)), uint16(s*TW), uint16(s*BH)}, borderValue )
		surface.FillRect( &sdl.Rect{int16(s*0)     , int16(s*BH)    , uint16(s*BW), uint16(s*H )}, borderValue )
		surface.FillRect( &sdl.Rect{int16(s*(BW+W)), int16(s*BH)    , uint16(s*BW), uint16(s*H )}, borderValue )
		
		if screen.borderEvents == nil {
			updateBorder(surface, scale)
		}
	}
}

func updateBorder(surface *sdl.Surface, scale uint) {
	s := scale
	
	const W = ScreenWidth
	const H = ScreenHeight
	const BW = ScreenBorderX
	const BH = ScreenBorderY
	const TW = TotalScreenWidth

	surface.UpdateRect( int32(s*0)     , int32(s*0)     , uint32(s*TW), uint32(s*BH) )
	surface.UpdateRect( int32(s*0)     , int32(s*(BH+H)), uint32(s*TW), uint32(s*BH) )
	surface.UpdateRect( int32(s*0)     , int32(s*BH)    , uint32(s*BW), uint32(s*H ) )
	surface.UpdateRect( int32(s*(BW+W)), int32(s*BH)    , uint32(s*BW), uint32(s*H ) )
}

// (This shouldn't be a public type)
type SimplifiedBorderEvent struct {
	tstate uint
	color RGBA
}

func scanlineFill(surface SDLSurface, scale uint, minx,maxx,y uint, color RGBA) {
	ys := y*scale
	
	if (y < ScreenBorderY) || (y >= TotalScreenHeight-ScreenBorderY) {
		for x:=minx ; x<=maxx ; x++ {
			xs := x*scale
			
			for yy:=uint(0) ; yy<scale; yy++ {
				for xx:=uint(0) ; xx<scale; xx++ {
					surface.setPixel(xs+xx, ys+yy, color)
				}
			}
		}
	} else {
		for x:=minx ; x<ScreenBorderX ; x++ {
			xs := x*scale
			
			for yy:=uint(0) ; yy<scale; yy++ {
				for xx:=uint(0) ; xx<scale; xx++ {
					surface.setPixel(xs+xx, ys+yy, color)
				}
			}
		}
		for x:=uint(TotalScreenWidth-ScreenBorderX) ; x<=maxx; x++ {
			xs := x*scale
			
			for yy:=uint(0) ; yy<scale; yy++ {
				for xx:=uint(0) ; xx<scale; xx++ {
					surface.setPixel(xs+xx, ys+yy, color)
				}
			}
		}
	}
}

// Render border in the interval [start,end)
func renderBorderBetweenTwoEvents(surface SDLSurface, scale uint, start *SimplifiedBorderEvent, end *SimplifiedBorderEvent) {
	// Spectrum 48k video timings
	const (
		TSTATES_PER_PIXEL = 2
	
		// Horizontal
		LINE_SCREEN       = 128		// 128 T states of screen
		LINE_RIGHT_BORDER = 24 		// 24 T states of right border
		LINE_RETRACE      = 48		// 48 T states of horizontal retrace
		LINE_LEFT_BORDER  = 24		// 24 T states of left border
		
		TSTATES_PER_LINE  = (LINE_RIGHT_BORDER + LINE_SCREEN + LINE_LEFT_BORDER + LINE_RETRACE)
		
		FIRST_BYTE        = 14336	// T states before the first byte of the screen (16384) is displayed
		
		// Vertical
		LINES_TOP         = 64
		LINES_SCREEN      = 192
		LINES_BOTTOM      = 56
		BORDER_TOP        = ScreenBorderY
		BORDER_BOTTOM     = ScreenBorderY
		
		// The T-state which corresponds to pixel (0,0) on the SDL surface
		DISPLAY_START     = ( FIRST_BYTE - TSTATES_PER_LINE*BORDER_TOP - ScreenBorderX/TSTATES_PER_PIXEL )
	)
	
	start_y := (int(start.tstate) - DISPLAY_START) / TSTATES_PER_LINE
	end_y   := (int(end.tstate)-1 - DISPLAY_START) / TSTATES_PER_LINE
	
	start_x := (int(start.tstate) - DISPLAY_START) - start_y*TSTATES_PER_LINE
	end_x   := (int(end.tstate)-1 - DISPLAY_START) - end_y*TSTATES_PER_LINE
	
	// Clip to visible screen area
	{
		if start_y < 0 {
			start_x = 0
			start_y = 0
		}
		if end_y >= TotalScreenHeight {
			end_x = TotalScreenWidth-1
			end_y = TotalScreenHeight-1
		}
		if start_x < 0 {
			start_x = 0
		}
		if end_x < 0 {
			end_x = 0
		}
		
		if end_y < 0 { return }
		if !(start_y <= end_y) { return }
		if (start_y == end_y) && !(start_x <= end_x) { return }
	}
	
	// Fill scanlines from (start_x,start_y) to (end_x,end_y)
	color := start.color
	if start_y == end_y {
		y := start_y
		scanlineFill(surface, scale, uint(start_x), uint(end_x), uint(y), color);
	} else {
		// Top scanline (start_y)
		scanlineFill(surface, scale, uint(start_x), TotalScreenWidth-1, uint(start_y), color);
		
		// Scanlines (start_y+1) ... (end_y-1)
		for y:=start_y+1 ; y<end_y ; y++ {
			scanlineFill(surface, scale, 0, TotalScreenWidth-1, uint(y), color);
		}
		
		// Bottom scanline (end_y)
		scanlineFill(surface, scale, 0, uint(end_x), uint(end_y), color);
	}
}

func renderBorderEvents(surface SDLSurface, scale uint, lastEvent_orNil *BorderEvent) {
	// Determine the number of border-events
	numEvents := 0
	for e:=lastEvent_orNil ; e!=nil ; e=e.previous_orNil {
		numEvents++
	}
	
	// Create the array 'events' with the events sorted by T-state value in *ascending* order
	events := make([]SimplifiedBorderEvent, numEvents+1)
	{
		i := numEvents-1
		for e:=lastEvent_orNil ; e!=nil ; e=e.previous_orNil {
			events[i] = SimplifiedBorderEvent{e.tstate, e.color}
			i--
		}
		// At this point: 'i' should equal to -1
		
		// The [border color from the last event] lasts until the end of the frame
		if lastEvent_orNil != nil {
			events[numEvents] = SimplifiedBorderEvent{TStatesPerFrame, lastEvent_orNil.color}
		}
	}
	
	// Note; If 'lastEvent_orNil' is nil, then 'event[numEvents]' is also nil. But this is OK.
	
	for i:=0 ; i < numEvents ; i++ {
		renderBorderBetweenTwoEvents(surface, scale, &events[i], &events[i+1])
	}
	
	updateBorder(surface.Surface, scale)
}

func (display *SDLScreen) render(screen, oldScreen_orNil *Screen) {
	const X0 = ScreenBorderX
	const Y0 = ScreenBorderY
	
	var attr_x, attr_y uint
	for attr_y = 0; attr_y < ScreenHeight_Attr; attr_y++ {
		for attr_x = 0; attr_x < ScreenWidth_Attr; attr_x++ {
			attr_ofs := (0x20*attr_y)+attr_x
			
			ink_paper := screen.attr[attr_ofs]
			
			changed_attr := false
			if oldScreen_orNil != nil {
				if !equals(oldScreen_orNil.attr[attr_ofs], ink_paper) {
					changed_attr = true
				}
			} else {
				changed_attr = true
			}
			
			srcBaseAddr := (0x800*(attr_y>>3) + 0x20*(attr_y&7)) + attr_x
			
			dst_X0 := X0+8*attr_x
			dst_Y0 := Y0+8*attr_y
			
			var y    uint = 0
			var y100 uint = 0
			for y < 8 {
				var value byte = screen.bitmap[srcBaseAddr + y100];
				
				if !changed_attr && (value == oldScreen_orNil.bitmap[srcBaseAddr + y100]) {
					y    += 1
					y100 += 0x100
					continue
				}
				
				for x := 7; x >= 0; x-- {
					color := ink_paper[value&1]
					display.ScreenSurface.setPixel(dst_X0+uint(x), dst_Y0+y, color)
					value = value >> 1
				}
				
				y    += 1
				y100 += 0x100
			}
		}
	}
	
	screenSurface := display.ScreenSurface.Surface
	screenSurface.UpdateRect(X0, Y0, ScreenWidth, ScreenHeight)
	
	renderBorder(screenSurface, /*scale*/1, screen, oldScreen_orNil)
	
	if screen.borderEvents != nil {
		renderBorderEvents(display.ScreenSurface, /*scale*/1, screen.borderEvents)
	}
}

type SDLScreen2x struct {
	// Channel for receiving screen data
	screenChannel  chan *Screen
	
	// The whole screen, borders included
	ScreenSurface SDLSurface
}

func NewSDLScreen2x(app *Application, screenSurface *sdl.Surface) *SDLScreen2x {
	SDL_screen := &SDLScreen2x{ make(chan *Screen), SDLSurface{ screenSurface } }
	
	go screenRenderLoop(app.NewEventLoop(), SDL_screen.screenChannel, SDL_screen)
	
	return SDL_screen
}

// Implement DisplayChannel
func (display *SDLScreen2x) getScreenChannel() chan *Screen {
	return display.screenChannel
}

func (display *SDLScreen2x) render(screen, oldScreen_orNil *Screen) {
	const X0 = ScreenBorderX
	const Y0 = ScreenBorderY
	
	var attr_x, attr_y uint
	for attr_y = 0; attr_y < ScreenHeight_Attr; attr_y++ {
		for attr_x = 0; attr_x < ScreenWidth_Attr; attr_x++ {
			attr_ofs := (0x20*attr_y)+attr_x
			
			ink_paper := screen.attr[attr_ofs]
			
			changed_attr := false
			if oldScreen_orNil != nil {
				if !equals(oldScreen_orNil.attr[attr_ofs], ink_paper) {
					changed_attr = true
				}
			} else {
				changed_attr = true
			}
			
			srcBaseAddr := (0x800*(attr_y>>3) + 0x20*(attr_y&7)) + attr_x
			
			dst_X0 := X0+8*attr_x
			dst_Y0 := Y0+8*attr_y
			
			var y    uint = 0
			var y100 uint = 0
			for y < 8 {
				var value byte = screen.bitmap[srcBaseAddr + y100];
				
				if !changed_attr && (value == oldScreen_orNil.bitmap[srcBaseAddr + y100]) {
					y    += 1
					y100 += 0x100
					continue
				}
				
				for x := 7; x >= 0; x-- {
					color := ink_paper[value&1]
					display.setPixel_2x(dst_X0+uint(x), dst_Y0+y, color)
					value = value >> 1
				}
				
				y    += 1
				y100 += 0x100
			}
		}
	}
	
	screenSurface := display.ScreenSurface.Surface
	screenSurface.UpdateRect(2*X0, 2*Y0, 2*ScreenWidth, 2*ScreenHeight)
	
	renderBorder(screenSurface, /*scale*/2, screen, oldScreen_orNil)
	
	if screen.borderEvents != nil {
		renderBorderEvents(display.ScreenSurface, /*scale*/2, screen.borderEvents)
	}
}

func (display *SDLScreen2x) setPixel_2x(x, y uint, color RGBA) {
	scaleX := x*2
	scaleY := y*2

	surface := display.ScreenSurface
	surface.setPixel(scaleX  , scaleY  , color)
	surface.setPixel(scaleX+1, scaleY+1, color)
	surface.setPixel(scaleX  , scaleY+1, color)
	surface.setPixel(scaleX+1, scaleY  , color)
}

