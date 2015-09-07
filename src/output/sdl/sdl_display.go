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

// +build linux freebsd

package sdl_output

import (
	"fmt"
	"github.com/scottferg/Go-SDL/sdl"
	"github.com/scottferg/Go-SDL/ttf"
	"github.com/remogatto/gospeccy/src/spectrum"
	"os"
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

func init() {
	const expectedVersion = "⚛SDL TTF bindings 1.0"
	actualVersion := ttf.GoSdlVersion()
	if actualVersion != expectedVersion {
		fmt.Fprintf(os.Stderr, "Invalid SDL font bindings version: expected \"%s\", got \"%s\"\n",
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
	pixels := uintptr(s.surface.Pixels)
	offset := uintptr(y*s.Pitch() + x*s.Bpp())

	return pixels + offset
}

func newSDLSurface(app *spectrum.Application, w, h int) *SDLSurface {
	surface := sdl.CreateRGBSurface(sdl.SWSURFACE, w, h, 32, 0, 0, 0, 0)
	if surface == nil {
		app.PrintfMsg("%s", sdl.GetError())
		app.RequestExit()
		return nil
	}
	return &SDLSurface{surface}
}

// Create an SDL surface suitable for a 2x scaled screen
func NewSDLSurface2x(app *spectrum.Application) *SDLSurface {
	return newSDLSurface(app, 2*spectrum.TotalScreenWidth, 2*spectrum.TotalScreenHeight)
}

// Create an SDL surface suitable for an unscaled screen
func NewSDLSurface(app *spectrum.Application) *SDLSurface {
	return newSDLSurface(app, spectrum.TotalScreenWidth, spectrum.TotalScreenHeight)
}

// ==============================
// Screen render loop (goroutine)
// ==============================

func screenRenderLoop(evtLoop *spectrum.EventLoop, screenChannel <-chan *spectrum.DisplayData, renderer screen_renderer_t) {
	terminating := false

	shutdown.Add(1)
	for {
		select {
		case <-evtLoop.Pause:
			terminating = true
			evtLoop.Pause <- 0

		case <-evtLoop.Terminate:
			// Terminate this Go routine
			if evtLoop.App().Verbose {
				evtLoop.App().PrintfMsg("screen render loop: exit")
			}
			evtLoop.Terminate <- 0
			shutdown.Done()
			return

		case screen := <-screenChannel:
			if screen != nil {
				if !terminating {
					renderer.render(screen)
				}
			} else {
				done := evtLoop.Delete()
				go func() { <-done }()
			}
		}

	}
}

// =========
// SDLScreen
// =========

type SDLScreen struct {
	// Channel for receiving display changes
	screenChannel chan *spectrum.DisplayData

	// The whole screen, borders included.
	// Initially nil.
	screenSurface *SDLSurface

	unscaledDisplay *UnscaledDisplay

	updatedRectsCh chan []sdl.Rect

	app *spectrum.Application
}

type screen_renderer_t interface {
	render(screen *spectrum.DisplayData)
}

func NewSDLScreen(app *spectrum.Application) *SDLScreen {
	SDL_screen := &SDLScreen{
		screenChannel:   make(chan *spectrum.DisplayData),
		screenSurface:   NewSDLSurface(app),
		unscaledDisplay: newUnscaledDisplay(),
		updatedRectsCh:  make(chan []sdl.Rect),
		app:             app,
	}

	go screenRenderLoop(app.NewEventLoop(), SDL_screen.screenChannel, SDL_screen)

	return SDL_screen
}

func (display *SDLScreen) UpdatedRectsCh() <-chan []sdl.Rect {
	return display.updatedRectsCh
}

func (display *SDLScreen) GetSurface() *sdl.Surface {
	return display.screenSurface.surface
}

// Implement DisplayReceiver
func (display *SDLScreen) GetDisplayDataChannel() chan<- *spectrum.DisplayData {
	return display.screenChannel
}

func (display *SDLScreen) Close() {
	display.screenChannel <- nil
}

// Implement screen_renderer_t
func (display *SDLScreen) render(screen *spectrum.DisplayData) {
	unscaledDisplay := display.unscaledDisplay
	unscaledDisplay.newFrame()
	unscaledDisplay.render(screen)

	surface := display.screenSurface
	bpp := surface.Bpp()
	pixels := &unscaledDisplay.pixels

	surface.surface.Lock()
	for _, r := range *unscaledDisplay.changedRegions {
		end_x := uint(r.X) + uint(r.W)
		end_y := uint(r.Y) + uint(r.H)

		for y := uint(r.Y); y < end_y; y++ {
			wy := spectrum.TotalScreenWidth * y
			addr := surface.addrXY(uint(r.X), y)
			for x := uint(r.X); x < end_x; x++ {
				*(*uint32)(unsafe.Pointer(addr)) = spectrum.Palette[pixels[wy+x]]
				addr += uintptr(bpp)
			}
		}
	}
	surface.surface.Unlock()

	if screen.CompletionTime_orNil != nil {
		screen.CompletionTime_orNil <- time.Now()
	}

	SDL_updateRects(surface.surface, unscaledDisplay.changedRegions, 1 /*scale*/, display.updatedRectsCh)
	unscaledDisplay.releaseMemory()

}

// ===========
// SDLScreen2x
// ===========

type SDLScreen2x struct {
	// Channel for receiving display changes
	screenChannel chan *spectrum.DisplayData

	// The whole screen, borders included.
	// Initially nil.
	screenSurface *SDLSurface

	unscaledDisplay *UnscaledDisplay

	updatedRectsCh chan []sdl.Rect

	app *spectrum.Application
}

func NewSDLScreen2x(app *spectrum.Application) *SDLScreen2x {
	SDL_screen := &SDLScreen2x{
		screenChannel:   make(chan *spectrum.DisplayData),
		screenSurface:   NewSDLSurface2x(app),
		unscaledDisplay: newUnscaledDisplay(),
		updatedRectsCh:  make(chan []sdl.Rect),
		app:             app,
	}

	go screenRenderLoop(app.NewEventLoop(), SDL_screen.screenChannel, SDL_screen)

	return SDL_screen
}

func (display *SDLScreen2x) UpdatedRectsCh() <-chan []sdl.Rect {
	return display.updatedRectsCh
}

func (display *SDLScreen2x) GetSurface() *sdl.Surface {
	return display.screenSurface.surface
}

// Implement DisplayReceiver
func (display *SDLScreen2x) GetDisplayDataChannel() chan<- *spectrum.DisplayData {
	return display.screenChannel
}

func (display *SDLScreen2x) Close() {
	display.screenChannel <- nil
}

// Implement screen_renderer_t
func (display *SDLScreen2x) render(screen *spectrum.DisplayData) {
	unscaledDisplay := display.unscaledDisplay
	unscaledDisplay.newFrame()
	unscaledDisplay.render(screen)

	surface := display.screenSurface
	bpp := uintptr(surface.Bpp())
	bpp2 := 2 * bpp
	pitch := uintptr(surface.Pitch())
	pixels := &unscaledDisplay.pixels

	surface.surface.Lock()
	for _, r := range *unscaledDisplay.changedRegions {
		end_x := uint(r.X) + uint(r.W)
		end_y := uint(r.Y) + uint(r.H)

		for y := uint(r.Y); y < end_y; y++ {
			addr := surface.addrXY(2*uint(r.X), 2*y)
			wy := spectrum.TotalScreenWidth * y

			for x := uint(r.X); x < end_x; x++ {
				color := spectrum.Palette[pixels[wy+x]]

				// Fill a 2x2 rectangle
				*(*uint32)(unsafe.Pointer(addr)) = color
				*(*uint32)(unsafe.Pointer(addr + bpp)) = color
				*(*uint32)(unsafe.Pointer(addr + pitch)) = color
				*(*uint32)(unsafe.Pointer(addr + pitch + bpp)) = color

				addr += bpp2
			}
		}
	}
	surface.surface.Unlock()

	if screen.CompletionTime_orNil != nil {
		screen.CompletionTime_orNil <- time.Now()
	}

	SDL_updateRects(surface.surface, unscaledDisplay.changedRegions, 2 /*scale*/, display.updatedRectsCh)
	unscaledDisplay.releaseMemory()
}

// ==============
// Misc functions
// ==============

func SDL_updateRects(surface *sdl.Surface, surfaceChanges *ListOfRects, scale uint, updatedRectsCh chan []sdl.Rect) {
	// Implementation note:
	//   This function does NOT make use of 'surface.UpdateRects',
	//   although in theory that would be much more efficient than 'surface.UpdateRect'.
	//   The reason is that using multiple rectangles makes SDL send *multiple*
	//   messages to the X server, thus causing serious visual artifacting.
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

	updatedRectsCh <- []sdl.Rect{{int16(x), int16(y), uint16(w), uint16(h)}}
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

	const W = spectrum.ScreenWidth
	const H = spectrum.ScreenHeight
	const BW = spectrum.ScreenBorderX
	const BH = spectrum.ScreenBorderY
	const TW = spectrum.TotalScreenWidth

	l.add(int(s*0), int(s*0), s*TW, s*BH)      // Top
	l.add(int(s*0), int(s*(BH+H)), s*TW, s*BH) // Bottom
	l.add(int(s*0), int(s*BH), s*BW, s*H)      // Left
	l.add(int(s*(BW+W)), int(s*BH), s*BW, s*H) // Right
}

// ===============
// UnscaledDisplay
// ===============

type UnscaledDisplay struct {
	pixels         [spectrum.TotalScreenWidth * spectrum.TotalScreenHeight]byte
	changedRegions *ListOfRects

	// This is the border which was rendered to 'pixels'
	border []spectrum.BorderEvent
}

func newUnscaledDisplay() *UnscaledDisplay {
	return &UnscaledDisplay{
		changedRegions: newListOfRects(),
		border:         nil,
	}
}

func (disp *UnscaledDisplay) newFrame() {
	disp.changedRegions = newListOfRects()
}

func (disp *UnscaledDisplay) releaseMemory() {
	disp.changedRegions = nil
}

// Set pixels from (minx,y) to (maxx,y). Both bounds are inclusive.
func (disp *UnscaledDisplay) scanlineFill(minx, maxx, y int, color byte) {
	wy := spectrum.TotalScreenWidth * y
	pixels := &disp.pixels

	if !(minx <= maxx) {
		return
	}

	if maxx >= spectrum.TotalScreenWidth {
		maxx = spectrum.TotalScreenWidth - 1
	}

	if (y < spectrum.ScreenBorderY) || (y >= spectrum.TotalScreenHeight-spectrum.ScreenBorderY) {
		for x := minx; x <= maxx; x++ {
			pixels[wy+x] = color
		}
	} else {
		for x := minx; x < spectrum.ScreenBorderX; x++ {
			pixels[wy+x] = color
		}
		for x := spectrum.TotalScreenWidth - spectrum.ScreenBorderX; x <= maxx; x++ {
			pixels[wy+x] = color
		}
	}
}

// Render border in the interval [start,end)
func (disp *UnscaledDisplay) renderBorderBetweenTwoEvents(start spectrum.BorderEvent, end spectrum.BorderEvent) {
	spectrum.Assert(start.TState < end.TState)

	const DISPLAY_START = spectrum.DISPLAY_START
	const TSTATES_PER_LINE = spectrum.TSTATES_PER_LINE

	if start.TState < DISPLAY_START {
		start.TState = DISPLAY_START
	}
	if end.TState-1 < DISPLAY_START {
		return
	}
	if start.TState >= DISPLAY_START+spectrum.TotalScreenHeight*TSTATES_PER_LINE {
		return
	}

	start_y := (start.TState - DISPLAY_START) / TSTATES_PER_LINE
	end_y := (end.TState - 1 - DISPLAY_START) / TSTATES_PER_LINE

	start_x := (start.TState - DISPLAY_START) % TSTATES_PER_LINE
	end_x := (end.TState - 1 - DISPLAY_START) % TSTATES_PER_LINE

	start_x = (start_x << spectrum.PIXELS_PER_TSTATE_LOG2) &^ 7
	end_x = (end_x << spectrum.PIXELS_PER_TSTATE_LOG2) &^ 7
	end_x += 7

	// Clip to visible screen area
	{
		if end_y >= spectrum.TotalScreenHeight {
			end_x = spectrum.TotalScreenWidth - 1
			end_y = spectrum.TotalScreenHeight - 1
		}
	}

	// Fill scanlines from (start_x,start_y) to (end_x,end_y)
	color := start.Color
	if start_y == end_y {
		y := start_y
		disp.scanlineFill(start_x, end_x, y, color)
	} else {
		// Top scanline (start_y)
		disp.scanlineFill(start_x, spectrum.TotalScreenWidth-1, start_y, color)

		// Scanlines (start_y+1) ... (end_y-1)
		for y := (start_y + 1); y < end_y; y++ {
			disp.scanlineFill(0, spectrum.TotalScreenWidth-1, y, color)
		}

		// Bottom scanline (end_y)
		disp.scanlineFill(0, end_x, end_y, color)
	}
}

func (disp *UnscaledDisplay) renderBorder(events []spectrum.BorderEvent) {
	if !spectrum.SameBorderEvents(disp.border, events) {
		if len(events) > 0 {
			firstEvent := &events[0]
			spectrum.Assert(firstEvent.TState == 0)

			lastEvent := &events[len(events)-1]
			spectrum.Assert(lastEvent.TState == spectrum.TStatesPerFrame)

			numEvents := len(events)

			for i := 0; i < numEvents-1; i++ {
				disp.renderBorderBetweenTwoEvents(events[i], events[i+1])
			}

			disp.changedRegions.addBorder( /*scale*/ 1)
		}

		disp.border = events
	}
}

// Table for extracting the numeric value of individual bits in an 8-bit number
var bitmap_unpack_table [1 << 8][8]uint

func init() {
	for a := uint(0); a < (1 << 8); a++ {
		bitmap_unpack_table_a := &bitmap_unpack_table[a]
		bitmap_unpack_table_a[0] = (a >> 7) & 1
		bitmap_unpack_table_a[1] = (a >> 6) & 1
		bitmap_unpack_table_a[2] = (a >> 5) & 1
		bitmap_unpack_table_a[3] = (a >> 4) & 1
		bitmap_unpack_table_a[4] = (a >> 3) & 1
		bitmap_unpack_table_a[5] = (a >> 2) & 1
		bitmap_unpack_table_a[6] = (a >> 1) & 1
		bitmap_unpack_table_a[7] = (a >> 0) & 1
	}
}

func (disp *UnscaledDisplay) render(screen *spectrum.DisplayData) {
	const X0 = spectrum.ScreenBorderX
	const Y0 = spectrum.ScreenBorderY

	screen_dirty := &screen.Dirty
	screen_attr := &screen.Attr
	screen_bitmap := &screen.Bitmap

	pixels := &disp.pixels

	var attr_x, attr_y uint
	for attr_y = 0; attr_y < spectrum.ScreenHeight_Attr; attr_y++ {
		dst_Y0 := Y0 + 8*attr_y
		attr_wy := spectrum.ScreenWidth_Attr * attr_y

		for attr_x = 0; attr_x < spectrum.ScreenWidth_Attr; attr_x++ {
			if screen_dirty[attr_wy+attr_x] {
				dst_X0 := X0 + 8*attr_x

				var y uint = 0
				var src_ofs uint = ((8 * attr_y) << spectrum.BytesPerLine_log2) + attr_x
				var dst_ofs uint = spectrum.TotalScreenWidth*(dst_Y0+y) + dst_X0
				for y < 8 {
					// Paper is in the lower 4 bits, ink is in the higher 4 bits
					var paperInk spectrum.Attr_4bit = screen_attr[src_ofs]
					paperInk_array := [2]uint8{uint8(paperInk) & 0xf, (uint8(paperInk) >> 4) & 0xf}

					var value byte = screen_bitmap[src_ofs]
					var unpacked_value *[8]uint = &bitmap_unpack_table[value]

					for x := 0; x < 8; x++ {
						color := paperInk_array[unpacked_value[x]]
						pixels[dst_ofs+uint(x)] = color
					}

					y += 1
					src_ofs += spectrum.BytesPerLine
					dst_ofs += spectrum.TotalScreenWidth
				}

				disp.changedRegions.add(int(dst_X0), int(dst_Y0), 8, 8)
			}
		}
	}

	disp.renderBorder(screen.BorderEvents)
}
