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
	var bpp, pitch = s.Bpp(), uint(s.Surface.Pitch)
	
	s.setValueAt(y * pitch + x*bpp, value)
}

func (s *SDLSurface) setPixel(x, y uint, color RGBA) {
	s.setValue(x, y, color.value32())
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

func (display *SDLScreen) render(screen, oldScreen_orNil *Screen) {
	unscaledDisplay.newFrame()
	unscaledDisplay.render(screen, oldScreen_orNil)
	
	surface := display.ScreenSurface
	pixels := &unscaledDisplay.pixels
	for _,r := range *unscaledDisplay.changedRegions {
		end_x := int(r.X) + int(r.W)
		end_y := int(r.Y) + int(r.H)
		
		for y:=int(r.Y) ; y<end_y ; y++ {
			wy := TotalScreenWidth*y
			for x:=int(r.X) ; x<end_x ; x++ {
				surface.setValue(uint(x), uint(y), palette[pixels[wy+x]])
			}
		}
	}
	
	if unscaledDisplay.border_orNil != nil {
		SDL_renderBorder(display.ScreenSurface.Surface, unscaledDisplay.changedRegions, /*scale*/1, *unscaledDisplay.border_orNil)
	}
	
	SDL_updateRects(display.ScreenSurface.Surface, unscaledDisplay.changedRegions, /*scale*/1)
	
	unscaledDisplay.releaseMemory()
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
	unscaledDisplay.newFrame()
	unscaledDisplay.render(screen, oldScreen_orNil)
	
	pixels := &unscaledDisplay.pixels
	for _,r := range *unscaledDisplay.changedRegions {
		end_x := int(r.X) + int(r.W)
		end_y := int(r.Y) + int(r.H)
		
		for y:=int(r.Y) ; y<end_y ; y++ {
			wy := TotalScreenWidth*y
			for x:=int(r.X) ; x<end_x ; x++ {
				display.setPixel_2x(uint(x), uint(y), palette[pixels[wy+x]])
			}
		}
	}
	
	if unscaledDisplay.border_orNil != nil {
		SDL_renderBorder(display.ScreenSurface.Surface, unscaledDisplay.changedRegions, /*scale*/2, *unscaledDisplay.border_orNil)
	}
	
	SDL_updateRects(display.ScreenSurface.Surface, unscaledDisplay.changedRegions, /*scale*/2)
	
	unscaledDisplay.releaseMemory()
}

func (display *SDLScreen2x) setPixel_2x(x, y uint, color uint32) {
	scaleX := x<<1
	scaleY := y<<1

	surface := display.ScreenSurface
	surface.setValue(scaleX  , scaleY  , color)
	surface.setValue(scaleX+1, scaleY+1, color)
	surface.setValue(scaleX  , scaleY+1, color)
	surface.setValue(scaleX+1, scaleY  , color)
}





func SDL_updateRects(surface *sdl.Surface, surfaceChanges *ListOfRects, scale uint) {
	if scale == 1 {
		surface.UpdateRects(*surfaceChanges)
	} else {
		scaledRects := make([]sdl.Rect, len(*surfaceChanges))
		
		for i,r := range *surfaceChanges {
			scaledRects[i] = sdl.Rect{ int16(scale)*r.X , int16(scale)*r.Y , uint16(scale)*r.W , uint16(scale)*r.H }
		}
		
		surface.UpdateRects(scaledRects)
	}
}

func SDL_renderBorder(surface *sdl.Surface, surfaceChanges *ListOfRects, scale uint, color byte) {
	s := scale
	c := palette[color]
	
	const W = ScreenWidth
	const H = ScreenHeight
	const BW = ScreenBorderX
	const BH = ScreenBorderY
	const TW = TotalScreenWidth

	surface.FillRect( &sdl.Rect{int16(s*0)     , int16(s*0)     , uint16(s*TW), uint16(s*BH)}, c )
	surface.FillRect( &sdl.Rect{int16(s*0)     , int16(s*(BH+H)), uint16(s*TW), uint16(s*BH)}, c )
	surface.FillRect( &sdl.Rect{int16(s*0)     , int16(s*BH)    , uint16(s*BW), uint16(s*H )}, c )
	surface.FillRect( &sdl.Rect{int16(s*(BW+W)), int16(s*BH)    , uint16(s*BW), uint16(s*H )}, c )
	
	// This is NOT a typo, the scale is actually 1 here
	surfaceChanges.addBorder(/*scale*/1)
}





type ListOfRects []sdl.Rect

func newListOfRects() *ListOfRects {
	l := new(ListOfRects)
	*l = make([]sdl.Rect,0,8)
	return l
}

func (l *ListOfRects) addRect(rect sdl.Rect) {
	slice := *l
	
	len_slice := len(slice)
	if len_slice == cap(slice) {
		// Double the capacity (assumes non-zero initial capacity)
        newSlice := make([]sdl.Rect, len_slice, 2*cap(slice))
		copy(newSlice, slice)
		slice = newSlice
	}
	
	slice = slice[0:len_slice+1]
	slice[len_slice] = rect
	
	*l = slice
}

func (l *ListOfRects) add(x,y int, w,h uint) {
	l.addRect( sdl.Rect{int16(x), int16(y), uint16(w), uint16(h)} )
}

func (l *ListOfRects) addBorder(scale uint) {
	s := scale
	
	const W = ScreenWidth
	const H = ScreenHeight
	const BW = ScreenBorderX
	const BH = ScreenBorderY
	const TW = TotalScreenWidth

	l.add( int(s*0)     , int(s*0)     , s*TW, s*BH )
	l.add( int(s*0)     , int(s*(BH+H)), s*TW, s*BH )
	l.add( int(s*0)     , int(s*BH)    , s*BW, s*H  )
	l.add( int(s*(BW+W)), int(s*BH)    , s*BW, s*H  )
}





type UnscaledDisplay struct {
	pixels         [TotalScreenWidth*TotalScreenHeight]byte
	changedRegions *ListOfRects
	border_orNil   *byte	// Valid in case the whole border has a single color
}

var unscaledDisplay = newUnscaledDisplay()

func newUnscaledDisplay() *UnscaledDisplay {
	return &UnscaledDisplay{changedRegions:newListOfRects(), border_orNil:nil}
}

func (disp *UnscaledDisplay) newFrame() {
	disp.changedRegions = newListOfRects()
	disp.border_orNil = nil
}

func (disp *UnscaledDisplay) releaseMemory() {
	disp.changedRegions = nil
	disp.border_orNil = nil
}

func (disp *UnscaledDisplay) renderBorder(screen *Screen, oldScreen_orNil *Screen) {
	if (oldScreen_orNil == nil) || (screen.border != oldScreen_orNil.border) || (oldScreen_orNil.borderEvents != nil) {
		var border byte = screen.border
		disp.border_orNil = &border
		
		disp.changedRegions.addBorder(/*scale*/1)
	}
}

// FIXME: This shouldn't be a public type
type SimplifiedBorderEvent struct {
	tstate uint
	color byte
}

func (disp *UnscaledDisplay) scanlineFill(minx,maxx,y uint, color byte) {
	wy := TotalScreenWidth*y
	
	if (y < ScreenBorderY) || (y >= TotalScreenHeight-ScreenBorderY) {
		for x:=minx ; x<=maxx ; x++ {
			disp.pixels[wy+x] = color
		}
	} else {
		for x:=minx ; x<ScreenBorderX ; x++ {
			disp.pixels[wy+x] = color
		}
		for x:=uint(TotalScreenWidth-ScreenBorderX) ; x<=maxx; x++ {
			disp.pixels[wy+x] = color
		}
	}
}

// Render border in the interval [start,end)
func (disp *UnscaledDisplay) renderBorderBetweenTwoEvents(start *SimplifiedBorderEvent, end *SimplifiedBorderEvent) {
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
		disp.scanlineFill(uint(start_x), uint(end_x), uint(y), color);
	} else {
		// Top scanline (start_y)
		disp.scanlineFill(uint(start_x), TotalScreenWidth-1, uint(start_y), color);
		
		// Scanlines (start_y+1) ... (end_y-1)
		for y:=start_y+1 ; y<end_y ; y++ {
			disp.scanlineFill(0, TotalScreenWidth-1, uint(y), color);
		}
		
		// Bottom scanline (end_y)
		disp.scanlineFill(0, uint(end_x), uint(end_y), color);
	}
}

func (disp *UnscaledDisplay) renderBorderEvents(lastEvent_orNil *BorderEvent) {
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
		disp.renderBorderBetweenTwoEvents(&events[i], &events[i+1])
	}
	
	disp.changedRegions.addBorder(/*scale*/1)
}

func (disp *UnscaledDisplay) render(screen, oldScreen_orNil *Screen) {
	const X0 = ScreenBorderX
	const Y0 = ScreenBorderY
	
	var attr_x, attr_y uint
	for attr_y = 0; attr_y < ScreenHeight_Attr; attr_y++ {
		srcBaseAddr_y := (0x800*(attr_y>>3) + ScreenWidth_Attr*(attr_y&7))
		dst_Y0        := Y0+8*attr_y
		
		for attr_x = 0; attr_x < ScreenWidth_Attr; attr_x++ {
			attr_ofs := (ScreenWidth_Attr*attr_y)+attr_x
			
			ink_paper := screen.attr[attr_ofs]
			
			changed_attr := false
			if oldScreen_orNil != nil {
				if !equals(oldScreen_orNil.attr[attr_ofs], ink_paper) {
					changed_attr = true
				}
			} else {
				changed_attr = true
			}
			
			dst_X0 := X0+8*attr_x
			
			changed_bitmap := false
			if changed_attr || screen.dirty[attr_ofs] {
				srcBaseAddr  := srcBaseAddr_y + attr_x
				single_color := (ink_paper[0] == ink_paper[1])
				
				var y      uint = 0
				var y100   uint = 0
				var dst_wy uint = TotalScreenWidth*(dst_Y0+y)
				for y < 8 {
					var value byte = screen.bitmap[srcBaseAddr + y100];
					
					if changed_attr || ((value != oldScreen_orNil.bitmap[srcBaseAddr + y100] && !single_color)) {
						changed_bitmap = true
						
						for x := 7; x >= 0; x-- {
							color := ink_paper[value&1]
							disp.pixels[dst_wy+(dst_X0+uint(x))] = color
							value = value >> 1
						}
					}
					
					y      += 1
					y100   += 0x100
					dst_wy += TotalScreenWidth
				}
			}
			
			if changed_attr || changed_bitmap {
				disp.changedRegions.add(int(dst_X0), int(dst_Y0), 8, 8)
			}
		}
	}
	
	if screen.borderEvents != nil {
		disp.renderBorderEvents(screen.borderEvents)
	} else {
		disp.renderBorder(screen, oldScreen_orNil)
	}
}


