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
	"fmt"
	"os"
	"⚛sdl"
	"time"
	"unsafe"
)

func init() {
	const expectedVersion = "⚛SDL bindings 1.0"
	actualVersion := sdl.GoSdlVersion()
	if actualVersion != expectedVersion {
		fmt.Fprintf(os.Stderr, "Invalid SDL bindings version: expected \"%s\", got \"%s\"\n",
					expectedVersion, actualVersion)
		os.Exit(1)
	}
}


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

// ==============================
// Screen render loop (goroutine)
// ==============================

func screenRenderLoop(evtLoop *EventLoop, screenChannel <-chan *DisplayData, renderer screen_renderer_t) {
	var screen *DisplayData
	for {
		select {
		case <-evtLoop.Pause:
			evtLoop.Pause <- 0

		case <-evtLoop.Terminate:
			// Terminate this Go routine
			if evtLoop.App().Verbose {
				evtLoop.App().PrintfMsg("screen render loop: exit")
			}
			evtLoop.Terminate <- 0
			return

		case screen = <-screenChannel:
		}
		if screen != nil {
			renderer.render(screen)
		} else {
			evtLoop.Delete()
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

	app *Application
}

type screen_renderer_t interface {
	render(screen *DisplayData)
}

func NewSDLScreen(app *Application) *SDLScreen {
	SDL_screen := &SDLScreen{
		screenChannel:   make(chan *DisplayData),
		screenSurface:   SDLSurface{nil},
		unscaledDisplay: newUnscaledDisplay(),
		app:             app,
	}

	go screenRenderLoop(app.NewEventLoop(), SDL_screen.screenChannel, SDL_screen)

	return SDL_screen
}

func (display *SDLScreen) GetSurface() *sdl.Surface {
	return display.screenSurface.surface
}

// Implement DisplayReceiver
func (display *SDLScreen) getDisplayDataChannel() chan<- *DisplayData {
	return display.screenChannel
}
func (display *SDLScreen) close() {
	display.screenChannel <- nil
}

// Implement screen_renderer_t
func (display *SDLScreen) render(screen *DisplayData) {
	unscaledDisplay := display.unscaledDisplay
	unscaledDisplay.newFrame()
	unscaledDisplay.render(screen)

	if display.screenSurface.surface == nil {
		var sdlMode uint32 = 0

		surface := sdl.SetVideoMode(TotalScreenWidth, TotalScreenHeight, 32, sdlMode)
		if surface == nil {
			display.app.PrintfMsg("%s", sdl.GetError())
			display.app.RequestExit()
			return
		}

		display.screenSurface.surface = surface
	}

	surface := display.screenSurface
	bpp     := surface.Bpp()
	pixels  := &unscaledDisplay.pixels

	for _, r := range *unscaledDisplay.changedRegions {
		end_x := uint(r.X) + uint(r.W)
		end_y := uint(r.Y) + uint(r.H)

		for y := uint(r.Y); y < end_y; y++ {
			wy := TotalScreenWidth * y
			addr := surface.addrXY(uint(r.X), y)
			for x := uint(r.X); x < end_x; x++ {
				*(*uint32)(unsafe.Pointer(addr)) = palette[pixels[wy+x]]
				addr += uintptr(bpp)
			}
		}
	}

	if screen.completionTime_orNil != nil {
		screen.completionTime_orNil <- time.Nanoseconds()
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
	frontendSurface *sdl.Surface

	unscaledDisplay *UnscaledDisplay
	
	app *Application
}

func NewSDLScreen2x(app *Application, fullscreen bool) *SDLScreen2x {
	SDL_screen := &SDLScreen2x{
		screenChannel:   make(chan *DisplayData),
		fullscreen:      fullscreen,
		screenSurface:   SDLSurface{nil},
		unscaledDisplay: newUnscaledDisplay(),
		app:             app,
	}

	go screenRenderLoop(app.NewEventLoop(), SDL_screen.screenChannel, SDL_screen)

	return SDL_screen
}

func (display *SDLScreen2x) GetSurface() *sdl.Surface {
	return display.screenSurface.surface
}

// Implement DisplayReceiver
func (display *SDLScreen2x) getDisplayDataChannel() chan<- *DisplayData {
	return display.screenChannel
}

func (display *SDLScreen2x) close() {
	display.screenChannel <- nil
}

// Implement screen_renderer_t
func (display *SDLScreen2x) render(screen *DisplayData) {
	unscaledDisplay := display.unscaledDisplay
	unscaledDisplay.newFrame()
	unscaledDisplay.render(screen)

	if display.screenSurface.surface == nil {
//		var sdlMode uint32
		// if display.fullscreen {
		// 	sdlMode = sdl.FULLSCREEN
		// } else {
		// 	sdlMode = 0
		// }

		surface := sdl.CreateRGBSurface(sdl.SWSURFACE, 2*TotalScreenWidth, 2*TotalScreenHeight, 32, 0, 0 ,0, 0)
		if surface == nil {
			display.app.PrintfMsg("%s", sdl.GetError())
			display.app.RequestExit()
			return
		}

		if display.fullscreen {
			sdl.ShowCursor(sdl.DISABLE)
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

				// Fill a 2x2 rectangle
				*(*uint32)(unsafe.Pointer(addr)) = color
				*(*uint32)(unsafe.Pointer(addr+bpp)) = color
				*(*uint32)(unsafe.Pointer(addr+pitch)) = color
				*(*uint32)(unsafe.Pointer(addr+pitch+bpp)) = color

				addr += bpp2
			}
		}
	}

	if screen.completionTime_orNil != nil {
		screen.completionTime_orNil <- time.Nanoseconds()
	}
 
	SDL_updateRects(surface.surface, unscaledDisplay.changedRegions, /*scale*/ 2)
	unscaledDisplay.releaseMemory()
}


// ==============
// Misc functions
// ==============

func SDL_updateRects(surface *sdl.Surface, surfaceChanges *ListOfRects, scale uint) {
	// Implementation note:
	//   This function does NOT make use of 'surface.UpdateRects',
	//   although in theory that would be much more efficient than 'surface.UpdateRect'.
	//   The reason is that using multiple rectangles makes SDL send *multiple*
	//   messages the X server, thus causing serious visual artifacting.
	//   For example, the Overscan demo by Busy Soft looks really bad
	//   if 'surface.UpdateRects' is used.
	//
	//   However, the current implementation should be more efficient
	//   than simply calling 'surface.Flip'.

	n := len(*surfaceChanges)
	if n == 0 {
		return
	}

	minx := int((*surfaceChanges)[0].X)
	miny := int((*surfaceChanges)[0].Y)
	maxx := minx + int((*surfaceChanges)[0].W)
	maxy := miny + int((*surfaceChanges)[0].H)

	for i := 1; i < n; i++ {
		rect_i := &(*surfaceChanges)[i]
		minx_i := int(rect_i.X)
		miny_i := int(rect_i.Y)
		maxx_i := minx_i + int(rect_i.W)
		maxy_i := miny_i + int(rect_i.H)

		if minx_i < minx {
			minx = minx_i
		}
		if miny_i < miny {
			miny = miny_i
		}
		if maxx_i > maxx {
			maxx = maxx_i
		}
		if maxy_i > maxy {
			maxy = maxy_i
		}
	}

	x := int32(scale) * int32(minx)
	y := int32(scale) * int32(miny)
	w := uint32(scale) * uint32(maxx-minx)
	h := uint32(scale) * uint32(maxy-miny)
	surface.UpdateRect(x, y, w, h)
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

	l.add( int(s*0)     , int(s*0)     , s*TW, s*BH )	// Top
	l.add( int(s*0)     , int(s*(BH+H)), s*TW, s*BH )	// Bottom
	l.add( int(s*0)     , int(s*BH)    , s*BW, s*H  )	// Left
	l.add( int(s*(BW+W)), int(s*BH)    , s*BW, s*H  )	// Right
}


// ===============
// UnscaledDisplay
// ===============

type UnscaledDisplay struct {
	pixels         [TotalScreenWidth * TotalScreenHeight]byte
	changedRegions *ListOfRects

	// This is the border which was rendered to 'pixels'
	border_orNil *BorderEvent
}

func newUnscaledDisplay() *UnscaledDisplay {
	return &UnscaledDisplay{changedRegions: newListOfRects()}
}

func (disp *UnscaledDisplay) newFrame() {
	disp.changedRegions = newListOfRects()
}

func (disp *UnscaledDisplay) releaseMemory() {
	disp.changedRegions = nil
}

// Set pixels from (minx,y) to (maxx,y). Both bounds are inclusive.
func (disp *UnscaledDisplay) scanlineFill(minx, maxx, y uint, color byte) {
	wy := TotalScreenWidth * y
	pixels := &disp.pixels

	if !(minx <= maxx) {
		return
	}

	if maxx >= TotalScreenWidth {
		maxx = TotalScreenWidth - 1
	}

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
func (disp *UnscaledDisplay) renderBorderBetweenTwoEvents(start *simplifiedBorderEvent_t, end *simplifiedBorderEvent_t) {
	assert(start.tstate < end.tstate)

	if start.tstate < DISPLAY_START {
		start.tstate = DISPLAY_START
	}
	if end.tstate-1 < DISPLAY_START {
		return
	}
	if start.tstate >= DISPLAY_START+TotalScreenHeight*TSTATES_PER_LINE {
		return
	}

	start_y := (start.tstate - DISPLAY_START) / TSTATES_PER_LINE
	end_y   := (end.tstate-1 - DISPLAY_START) / TSTATES_PER_LINE

	start_x := (start.tstate - DISPLAY_START) % TSTATES_PER_LINE
	end_x   := (end.tstate-1 - DISPLAY_START) % TSTATES_PER_LINE

	start_x = (start_x << PIXELS_PER_TSTATE_LOG2) &^ 7
	end_x = (end_x << PIXELS_PER_TSTATE_LOG2) &^ 7
	end_x += 7

	// Clip to visible screen area
	{
		if end_y >= TotalScreenHeight {
			end_x = TotalScreenWidth - 1
			end_y = TotalScreenHeight - 1
		}
	}

	// Fill scanlines from (start_x,start_y) to (end_x,end_y)
	color := start.color
	if start_y == end_y {
		y := start_y
		disp.scanlineFill(start_x, end_x, y, color)
	} else {
		// Top scanline (start_y)
		disp.scanlineFill(start_x, TotalScreenWidth-1, start_y, color)

		// Scanlines (start_y+1) ... (end_y-1)
		for y := (start_y + 1); y < end_y; y++ {
			disp.scanlineFill(0, TotalScreenWidth-1, y, color)
		}

		// Bottom scanline (end_y)
		disp.scanlineFill(0, end_x, end_y, color)
	}
}

func (disp *UnscaledDisplay) renderBorder(lastEvent_orNil *BorderEvent) {
	if !disp.border_orNil.Equals(lastEvent_orNil) {
		if lastEvent_orNil != nil {
			lastEvent := lastEvent_orNil
			assert(lastEvent.tstate == TStatesPerFrame)

			// Put the events in an array, sorted by T-state value in ascending order
			var events []simplifiedBorderEvent_t
			{
				events_array := &simplifiedBorderEvent_array_t{}
				EventListToArray_Ascending(lastEvent, events_array, nil)
				events = events_array.events
			}

			numEvents := len(events)

			for i := 0; i < numEvents-1; i++ {
				disp.renderBorderBetweenTwoEvents(&events[i], &events[i+1])
			}

			disp.changedRegions.addBorder( /*scale*/ 1)
		}

		disp.border_orNil = lastEvent_orNil
	}
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

func (disp *UnscaledDisplay) render(screen *DisplayData) {
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

	disp.renderBorder(screen.borderEvents_orNil)
}


// =======================
// Simplified border-event
// =======================

type simplifiedBorderEvent_t struct {
	tstate uint
	color  byte
}

type simplifiedBorderEvent_array_t struct {
	events []simplifiedBorderEvent_t
}

func (a *simplifiedBorderEvent_array_t) Init(n int) {
	a.events = make([]simplifiedBorderEvent_t, n)
}

func (a *simplifiedBorderEvent_array_t) Set(i int, _e Event) {
	e := _e.(*BorderEvent)
	a.events[i] = simplifiedBorderEvent_t{e.tstate, e.color}
}
