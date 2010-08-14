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


// ==========
// SDLSurface
// ==========

type SDLSurface struct {
	surface *sdl.Surface
}

func (s SDLSurface) Width() uint {
	return uint(s.surface.W)
}

func (s SDLSurface) Height() uint {
	return uint(s.surface.H)
}

func (s SDLSurface) Bpp() uint {
	return uint(s.surface.Format.BytesPerPixel)
}

func (s SDLSurface) Pitch() uint {
	return uint(s.surface.Pitch)
}

// Return the address of pixel at (x,y)
func (s SDLSurface) addrXY(x, y uint) uintptr {
	pixels := uintptr(unsafe.Pointer(s.surface.Pixels))
	offset := uintptr(y*s.Pitch() + x*s.Bpp())

	return uintptr(unsafe.Pointer(pixels + offset))
}

func (s SDLSurface) setPixel(addr uintptr, color uint32) {
	*(*uint32)(unsafe.Pointer(addr)) = color
}


// ==============================
// Screen render loop (goroutine)
// ==============================

func screenRenderLoop(evtLoop *EventLoop, screenChannel chan *DisplayData, renderer screen_renderer_t) {
	var screen *DisplayData = nil
	var oldScreen *DisplayData = nil
	for {
		select {
		case <-evtLoop.Pause:
			evtLoop.Pause <- 0

		case <-evtLoop.Terminate:
			// Terminate this Go routine
			if evtLoop.App().Verbose {
				println("screen render loop: exit")
			}
			evtLoop.Terminate <- 0
			return

		case screen = <-screenChannel:
			if screen != nil {
				renderer.render(screen, oldScreen)
				oldScreen = screen
				screen = nil
			} else {
				evtLoop.Delete()
			}
		}
	}
}


// =========
// SDLScreen
// =========

type SDLScreen struct {
	// Channel for receiving display changes
	screenChannel chan *DisplayData

	// The whole screen, borders included.
	// Initially nil.
	screenSurface SDLSurface

	unscaledDisplay *UnscaledDisplay
}

type screen_renderer_t interface {
	render(screen, oldScreen_orNil *DisplayData)
}

func NewSDLScreen(app *Application) *SDLScreen {
	SDL_screen := &SDLScreen{make(chan *DisplayData), SDLSurface{nil}, newUnscaledDisplay()}

	go screenRenderLoop(app.NewEventLoop(), SDL_screen.screenChannel, SDL_screen)

	return SDL_screen
}

// Implement DisplayReceiver
func (display *SDLScreen) getDisplayDataChannel() chan *DisplayData {
	return display.screenChannel
}
func (display *SDLScreen) close() {
	display.screenChannel <- nil
}

// Implement screen_renderer_t
func (display *SDLScreen) render(screen, oldScreen_orNil *DisplayData) {
	unscaledDisplay := display.unscaledDisplay
	unscaledDisplay.newFrame()
	unscaledDisplay.render(screen, oldScreen_orNil)

	if display.screenSurface.surface == nil {
		var sdlMode uint32 = 0

		surface := sdl.SetVideoMode(TotalScreenWidth, TotalScreenHeight, 32, sdlMode)
		if surface == nil {
			panic(sdl.GetError())
		}

		display.screenSurface.surface = surface
	}

	surface := display.screenSurface
	bpp     := surface.Bpp()
	pixels  := &unscaledDisplay.pixels

	for _,r := range *unscaledDisplay.changedRegions {
		end_x := uint(r.X) + uint(r.W)
		end_y := uint(r.Y) + uint(r.H)

		for y := uint(r.Y); y < end_y; y++ {
			wy := TotalScreenWidth * y
			addr := surface.addrXY(uint(r.X), y)
			for x := uint(r.X); x < end_x; x++ {
				surface.setPixel(addr, palette[pixels[wy+x]])
				addr += uintptr(bpp)
			}
		}
	}

	if unscaledDisplay.border_orNil != nil {
		SDL_renderBorder(surface.surface, unscaledDisplay.changedRegions, /*scale*/ 1, *unscaledDisplay.border_orNil)
	}

	SDL_updateRects(surface.surface, unscaledDisplay.changedRegions, /*scale*/ 1)

	unscaledDisplay.releaseMemory()
}


// ===========
// SDLScreen2x
// ===========

type SDLScreen2x struct {
	// Channel for receiving display changes
	screenChannel chan *DisplayData

	fullscreen bool

	// The whole screen, borders included.
	// Initially nil.
	screenSurface SDLSurface

	unscaledDisplay *UnscaledDisplay
}

func NewSDLScreen2x(app *Application, fullscreen bool) *SDLScreen2x {
	SDL_screen := &SDLScreen2x{make(chan *DisplayData), fullscreen, SDLSurface{nil}, newUnscaledDisplay()}

	go screenRenderLoop(app.NewEventLoop(), SDL_screen.screenChannel, SDL_screen)

	return SDL_screen
}

// Implement DisplayReceiver
func (display *SDLScreen2x) getDisplayDataChannel() chan *DisplayData {
	return display.screenChannel
}

func (display *SDLScreen2x) close() {
	display.screenChannel <- nil
}

// Implement screen_renderer_t
func (display *SDLScreen2x) render(screen, oldScreen_orNil *DisplayData) {
	unscaledDisplay := display.unscaledDisplay
	unscaledDisplay.newFrame()
	unscaledDisplay.render(screen, oldScreen_orNil)

	if display.screenSurface.surface == nil {
		var sdlMode uint32
		if display.fullscreen {
			sdlMode = sdl.FULLSCREEN
		} else {
			sdlMode = 0
		}

		surface := sdl.SetVideoMode(2*TotalScreenWidth, 2*TotalScreenHeight, 32, sdlMode)
		if surface == nil {
			panic(sdl.GetError())
		}

		display.screenSurface.surface = surface
	}

	surface := display.screenSurface
	bpp     := uintptr(surface.Bpp())
	bpp2    := 2 * bpp
	pitch   := uintptr(surface.Pitch())
	pixels  := &unscaledDisplay.pixels

	for _, r := range *unscaledDisplay.changedRegions {
		end_x := uint(r.X) + uint(r.W)
		end_y := uint(r.Y) + uint(r.H)

		for y := uint(r.Y); y < end_y; y++ {
			addr := surface.addrXY(2*uint(r.X), 2*y)
			wy := TotalScreenWidth * y

			for x := uint(r.X); x < end_x; x++ {
				color := palette[pixels[wy+x]]

				surface.setPixel(addr+0, color)
				surface.setPixel(addr+bpp, color)
				surface.setPixel(addr+pitch+0, color)
				surface.setPixel(addr+pitch+bpp, color)

				addr += bpp2
			}
		}
	}

	if unscaledDisplay.border_orNil != nil {
		SDL_renderBorder(surface.surface, unscaledDisplay.changedRegions, /*scale*/ 2, *unscaledDisplay.border_orNil)
	}

	SDL_updateRects(surface.surface, unscaledDisplay.changedRegions, /*scale*/ 2)

	unscaledDisplay.releaseMemory()
}


// ==============
// Misc functions
// ==============

func SDL_updateRects(surface *sdl.Surface, surfaceChanges *ListOfRects, scale uint) {
	//for _,r := range *surfaceChanges {
	//	println("(", r.X, r.Y, r.W, r.H, ")")
	//}

	if scale == 1 {
		surface.UpdateRects(*surfaceChanges)
	} else {
		scaledRects := make([]sdl.Rect, len(*surfaceChanges))

		for i, r := range *surfaceChanges {
			scaledRects[i] = sdl.Rect{int16(scale) * r.X, int16(scale) * r.Y, uint16(scale) * r.W, uint16(scale) * r.H}
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

	// This is NOT a typo, the scale is actually 1 here.
	// The rectangles will be scaled later.
	surfaceChanges.addBorder( /*scale*/ 1)
}


// ===========
// ListOfRects
// ===========

type ListOfRects []sdl.Rect

func newListOfRects() *ListOfRects {
	l := new(ListOfRects)
	*l = make([]sdl.Rect, 0, 8)
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

	slice = slice[0 : len_slice+1]
	slice[len_slice] = rect

	*l = slice
}

func (l *ListOfRects) add(x, y int, w, h uint) {
	l.addRect(sdl.Rect{int16(x), int16(y), uint16(w), uint16(h)})
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


// ===============
// UnspacedDisplay
// ===============

type UnscaledDisplay struct {
	pixels         [TotalScreenWidth * TotalScreenHeight]byte
	changedRegions *ListOfRects
	border_orNil   *byte // Valid in case the whole border has a single color
}

func newUnscaledDisplay() *UnscaledDisplay {
	return &UnscaledDisplay{changedRegions: newListOfRects(), border_orNil: nil}
}

func (disp *UnscaledDisplay) newFrame() {
	disp.changedRegions = newListOfRects()
	disp.border_orNil = nil
}

func (disp *UnscaledDisplay) releaseMemory() {
	disp.changedRegions = nil
	disp.border_orNil = nil
}

func (disp *UnscaledDisplay) renderBorder(screen, oldScreen_orNil *DisplayData) {
	if (oldScreen_orNil == nil) || (screen.border != oldScreen_orNil.border) || (oldScreen_orNil.borderEvents != nil) {
		var border byte = screen.border
		disp.border_orNil = &border

		disp.changedRegions.addBorder( /*scale*/ 1)
	}
}

// FIXME: This shouldn't be a public type
type SimplifiedBorderEvent struct {
	tstate uint
	color  byte
}

func (disp *UnscaledDisplay) scanlineFill(minx, maxx, y uint, color byte) {
	wy := TotalScreenWidth * y
	pixels := &disp.pixels

	if (y < ScreenBorderY) || (y >= TotalScreenHeight-ScreenBorderY) {
		for x := minx; x <= maxx; x++ {
			pixels[wy+x] = color
		}
	} else {
		for x := minx; x < ScreenBorderX; x++ {
			pixels[wy+x] = color
		}
		for x := uint(TotalScreenWidth - ScreenBorderX); x <= maxx; x++ {
			pixels[wy+x] = color
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
			end_x = TotalScreenWidth - 1
			end_y = TotalScreenHeight - 1
		}
		if start_x < 0 {
			start_x = 0
		}
		if end_x < 0 {
			end_x = 0
		}

		if end_y < 0 {
			return
		}
		if !(start_y <= end_y) {
			return
		}
		if (start_y == end_y) && !(start_x <= end_x) {
			return
		}
	}

	// Fill scanlines from (start_x,start_y) to (end_x,end_y)
	color := start.color
	if start_y == end_y {
		y := start_y
		disp.scanlineFill(uint(start_x), uint(end_x), uint(y), color)
	} else {
		// Top scanline (start_y)
		disp.scanlineFill(uint(start_x), TotalScreenWidth-1, uint(start_y), color)

		// Scanlines (start_y+1) ... (end_y-1)
		for y := (start_y + 1); y < end_y; y++ {
			disp.scanlineFill(0, TotalScreenWidth-1, uint(y), color)
		}

		// Bottom scanline (end_y)
		disp.scanlineFill(0, uint(end_x), uint(end_y), color)
	}
}

func (disp *UnscaledDisplay) renderBorderEvents(lastEvent_orNil *BorderEvent) {
	// Determine the number of border-events
	numEvents := 0
	for e := lastEvent_orNil; e != nil; e = e.previous_orNil {
		numEvents++
	}

	// Create an array called 'events' and initialize it with
	// the events sorted by T-state value in *ascending* order
	events := make([]SimplifiedBorderEvent, numEvents+1)
	{
		i := numEvents - 1
		for e := lastEvent_orNil; e != nil; e = e.previous_orNil {
			events[i] = SimplifiedBorderEvent{e.tstate, e.color}
			i--
		}
		// At this point: 'i' should equal to -1

		// The [border color from the last event] lasts until the end of the frame
		if lastEvent_orNil != nil {
			events[numEvents] = SimplifiedBorderEvent{TStatesPerFrame, lastEvent_orNil.color}
		}
	}

	// Note: If 'lastEvent_orNil' is nil, then 'event[numEvents]' is also nil. But this is OK.

	for i := 0; i < numEvents; i++ {
		disp.renderBorderBetweenTwoEvents(&events[i], &events[i+1])
	}

	disp.changedRegions.addBorder( /*scale*/ 1)
}

// Table for extracting the numeric value of individual bits in an 8-bit number
var bitmap_unpack_table [1 << 8][8]uint

func init() {
	for a := uint(0); a < (1 << 8); a++ {
		bitmap_unpack_table[a][0] = (a >> 7) & 1
		bitmap_unpack_table[a][1] = (a >> 6) & 1
		bitmap_unpack_table[a][2] = (a >> 5) & 1
		bitmap_unpack_table[a][3] = (a >> 4) & 1
		bitmap_unpack_table[a][4] = (a >> 3) & 1
		bitmap_unpack_table[a][5] = (a >> 2) & 1
		bitmap_unpack_table[a][6] = (a >> 1) & 1
		bitmap_unpack_table[a][7] = (a >> 0) & 1
	}
}

func (disp *UnscaledDisplay) render(screen, oldScreen_orNil *DisplayData) {
	const X0 = ScreenBorderX
	const Y0 = ScreenBorderY

	screen_dirty := &screen.dirty
	screen_attr := &screen.attr
	screen_bitmap := &screen.bitmap

	pixels := &disp.pixels

	var attr_x, attr_y uint
	for attr_y = 0; attr_y < ScreenHeight_Attr; attr_y++ {
		dst_Y0 := Y0 + 8*attr_y
		attr_wy := ScreenWidth_Attr * attr_y

		for attr_x = 0; attr_x < ScreenWidth_Attr; attr_x++ {
			if screen_dirty[attr_wy+attr_x] {
				dst_X0 := X0 + 8*attr_x

				var y       uint = 0
				var src_ofs uint = ((8 * attr_y) << BytesPerLine_log2) + attr_x
				var dst_ofs uint = TotalScreenWidth*(dst_Y0+y) + dst_X0
				for y < 8 {
					// Paper is in the lower 4 bits, ink is in the higher 4 bits
					var paperInk attr_4bit = screen_attr[src_ofs]
					paperInk_array := [2]uint8{uint8(paperInk) & 0xf, (uint8(paperInk) >> 4) & 0xf}

					var value          byte     = screen_bitmap[src_ofs]
					var unpacked_value *[8]uint = &bitmap_unpack_table[value]

					for x := 0; x < 8; x++ {
						color := paperInk_array[unpacked_value[x]]
						pixels[dst_ofs+uint(x)] = color
					}

					y += 1
					src_ofs += BytesPerLine
					dst_ofs += TotalScreenWidth
				}

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
