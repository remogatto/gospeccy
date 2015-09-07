/*
 * Copyright: ⚛ <0xe2.0x9a.0x9b@gmail.com> 2011
 *
 * The contents of this file can be used freely,
 * except for usages in immoral contexts.
 */

// +build linux freebsd

package sdl_output

import (
	"github.com/scottferg/Go-SDL/sdl"
	"github.com/remogatto/gospeccy/src/spectrum"
	"math/rand"
	"unsafe"
)

// Composes multiple SDL surfaces into a single surface
type SDLSurfaceComposer struct {
	// The surfaces to compose.
	//
	// The order of surfaces in this array defines the compositing order.
	// The first surface will visually appear at the bottom,
	// while the last surface will visually appear at the top.
	inputs []*input_surface_t

	// The surface where to put the composited image
	output_orNil *sdl.Surface

	commandChannel chan interface{}

	showPaintedRegions bool
}

type input_surface_t struct {
	surface        *sdl.Surface
	updatedRectsCh <-chan []sdl.Rect
	forwarderLoop  *spectrum.EventLoop
	x, y           int
}

// Creates a new composer, and starts its command-loop in a goroutine
func NewSDLSurfaceComposer(app *spectrum.Application) *SDLSurfaceComposer {
	composer := &SDLSurfaceComposer{
		inputs:             make([]*input_surface_t, 0),
		output_orNil:       nil,
		commandChannel:     make(chan interface{}),
		showPaintedRegions: false,
	}

	go composer.commandLoop(app)

	return composer
}

// Enqueues a command that will append the specified surface
// to the end of [the list of input surfaces of 'composer'].
//
// The order of surfaces in the mentioned list defines the compositing order.
// The first surface will visually appear at the bottom,
// while the last surface will visually appear at the top.
func (composer *SDLSurfaceComposer) AddInputSurface(surface *sdl.Surface, x, y int, updatedRectsCh <-chan []sdl.Rect) {
	composer.commandChannel <- cmd_add{surface, x, y, updatedRectsCh}
}

// Enqueues a command that will remove the specified surface
// from [the list of input surfaces of 'composer'].
// The returned channel will receive a single value when the command completes.
func (composer *SDLSurfaceComposer) RemoveInputSurface(surface *sdl.Surface) <-chan byte {
	done := make(chan byte)
	composer.commandChannel <- cmd_remove{surface, done}
	return done
}

// Enqueues a command that will clear [the list of input surfaces of 'composer'].
// The returned channel will receive a single value when the command completes.
func (composer *SDLSurfaceComposer) RemoveAllInputSurfaces() <-chan byte {
	done := make(chan byte)
	composer.commandChannel <- cmd_removeAll{done}
	return done
}

// Enqueues a command that will move the specified surface to a new position
func (composer *SDLSurfaceComposer) SetPosition(surface *sdl.Surface, x, y int) {
	composer.commandChannel <- cmd_setPosition{surface, x, y}
}

// Enqueues a command that will replace the output surface.
// The returned channel will receive a single value when the command completes.
func (composer *SDLSurfaceComposer) ReplaceOutputSurface(surface_orNil *sdl.Surface) <-chan byte {
	done := make(chan byte)
	composer.commandChannel <- cmd_replaceOutputSurface{surface_orNil, done}
	return done
}

// Enqueues a command that will set the "show painted regions" flag.
// If the flag is set to true, then each paint operation is covered
// by a semi-transparent randomly colored rectangle.
func (composer *SDLSurfaceComposer) ShowPaintedRegions(enable bool) {
	composer.commandChannel <- cmd_showPaintedRegions{enable}
}

type cmd_add struct {
	surface        *sdl.Surface
	x, y           int
	updatedRectsCh <-chan []sdl.Rect
}

type cmd_remove struct {
	surface *sdl.Surface
	done    chan<- byte
}

type cmd_removeAll struct {
	done chan<- byte
}

type cmd_setPosition struct {
	surface *sdl.Surface
	x, y    int
}

type cmd_replaceOutputSurface struct {
	surface_orNil *sdl.Surface
	done          chan<- byte
}

type cmd_showPaintedRegions struct {
	enable bool
}

type cmd_update struct {
	surface *input_surface_t
	rects   []sdl.Rect
}

// The composer's command loop.
// This function runs in a separate goroutine.
func (composer *SDLSurfaceComposer) commandLoop(app *spectrum.Application) {
	evtLoop := app.NewEventLoop()

	shutdown.Add(1)
	for {
		select {
		case <-evtLoop.Pause:
			evtLoop.Pause <- 0

		case <-evtLoop.Terminate:
			// Terminate this goroutine
			if app.Verbose {
				app.PrintfMsg("surface compositing loop: exit")
			}
			evtLoop.Terminate <- 0
			shutdown.Done()
			return

		case untyped_cmd := <-composer.commandChannel:
			switch cmd := untyped_cmd.(type) {
			case cmd_add:
				composer.add(app, cmd.surface, cmd.x, cmd.y, cmd.updatedRectsCh)

			case cmd_remove:
				composer.remove(cmd.surface, cmd.done)

			case cmd_removeAll:
				composer.removeAll(cmd.done)

			case cmd_setPosition:
				composer.setPosition(cmd.surface, cmd.x, cmd.y)

			case cmd_replaceOutputSurface:
				composer.output_orNil = cmd.surface_orNil
				cmd.done <- 0

			case cmd_showPaintedRegions:
				composer.showPaintedRegions = cmd.enable
				composer.repaintTheWholeOutputSurface()

			case cmd_update:
				composer.performCompositing(cmd.surface.x, cmd.surface.y, cmd.rects)
			}
		}
	}
}

// Receives changed rectangles from the 'updatedRectsCh' channel,
// and sends them to the 'commandChannel' as instances of 'cmd_update'.
// This function runs in a separate goroutine.
func (composer *SDLSurfaceComposer) forwarderLoop(s *input_surface_t) {
	evtLoop := s.forwarderLoop
	updatedRectsCh_orNil := s.updatedRectsCh

	shutdown.Add(1)
	for {
		select {
		case <-evtLoop.Pause:
			updatedRectsCh_orNil = nil
			go func() {
				for rect := range s.updatedRectsCh {
					if rect == nil {
						break
					}
				}
			}()
			evtLoop.Pause <- 0

		case <-evtLoop.Terminate:
			// Terminate this function
			if evtLoop.App().Verbose {
				evtLoop.App().PrintfMsg("surface compositing: a forwarder loop: exit")
			}
			evtLoop.Terminate <- 0
			shutdown.Done()
			return

		case rects := <-updatedRectsCh_orNil:
			if rects != nil {
				composer.commandChannel <- cmd_update{s, rects}
			}
		}
	}
}

func (composer *SDLSurfaceComposer) indexOf(surface *sdl.Surface) int {
	n := len(composer.inputs)
	for i := 0; i < n; i++ {
		if composer.inputs[i].surface == surface {
			return i
		}
	}

	panic("no such surface")
}

func (composer *SDLSurfaceComposer) add(app *spectrum.Application, surface *sdl.Surface, x, y int, updatedRectsCh <-chan []sdl.Rect) {
	newInput := &input_surface_t{
		surface:        surface,
		updatedRectsCh: updatedRectsCh,
		forwarderLoop:  app.NewEventLoop(),
		x:              x,
		y:              y,
	}
	composer.inputs = append(composer.inputs, newInput)

	updateRect := sdl.Rect{
		X: int16(0),
		Y: int16(0),
		W: uint16(newInput.surface.W),
		H: uint16(newInput.surface.H),
	}

	composer.performCompositing(x, y, []sdl.Rect{updateRect})

	go composer.forwarderLoop(newInput)
}

func (composer *SDLSurfaceComposer) remove(surface *sdl.Surface, done chan<- byte) {
	i := composer.indexOf(surface)
	input := composer.inputs[i]

	// Remove the i-th element from 'composer.inputs'
	copy(composer.inputs[i:], composer.inputs[i+1:])
	composer.inputs = composer.inputs[0 : len(composer.inputs)-1]

	updateRect := sdl.Rect{
		X: int16(0),
		Y: int16(0),
		W: uint16(input.surface.W),
		H: uint16(input.surface.H),
	}
	composer.performCompositing(input.x, input.y, []sdl.Rect{updateRect})

	go func() {
		deleted := input.forwarderLoop.Delete()
		<-deleted

		done <- 0
	}()
}

func (composer *SDLSurfaceComposer) removeAll(done chan<- byte) {
	oldInputs := composer.inputs
	composer.inputs = make([]*input_surface_t, 0)

	composer.repaintTheWholeOutputSurface()

	go func() {
		for _, oldInput := range oldInputs {
			deleted := oldInput.forwarderLoop.Delete()
			<-deleted
		}

		done <- 0
	}()
}

func (composer *SDLSurfaceComposer) setPosition(surface *sdl.Surface, newX, newY int) {
	input := composer.inputs[composer.indexOf(surface)]

	oldX, oldY := input.x, input.y

	if (oldX != newX) || (oldY != newY) {
		input.x, input.y = newX, newY

		// Compute the rectangle which needs to repainted

		var x, y int
		var absDeltaX, absDeltaY uint
		if oldX <= newX {
			x = oldX
			absDeltaX = uint(newX - oldX)
		} else {
			x = newX
			absDeltaX = uint(oldX - newX)
		}
		if oldY <= newY {
			y = oldY
			absDeltaY = uint(newY - oldY)
		} else {
			y = newY
			absDeltaY = uint(oldY - newY)
		}

		updateRect := sdl.Rect{
			X: int16(x - newX),
			Y: int16(y - newY),
			W: uint16(uint(input.surface.W) + absDeltaX),
			H: uint16(uint(input.surface.H) + absDeltaY),
		}

		// Repaint
		composer.performCompositing(input.x, input.y, []sdl.Rect{updateRect})
	}
}

func (composer *SDLSurfaceComposer) repaintTheWholeOutputSurface() {
	if composer.output_orNil != nil {
		updateRect := sdl.Rect{
			X: int16(0),
			Y: int16(0),
			W: uint16(composer.output_orNil.W),
			H: uint16(composer.output_orNil.H),
		}
		composer.performCompositing(0, 0, []sdl.Rect{updateRect})
	}
}

// Used to generate colors when 'showPaintedRegions' is true
var rnd *rand.Rand = rand.New(rand.NewSource(0))

// The surface compositing function.
//
// rects: The list of changes
//
// ofsX, ofsY: The translation to be applied to each element of 'rects'.
//             After the translation, the position of each rectangle is relative
//             to the coordinate system of the output surface.
func (composer *SDLSurfaceComposer) performCompositing(ofsX, ofsY int, rects []sdl.Rect) {
	if composer.output_orNil != nil {
		output := composer.output_orNil

		updateRects := make([]sdl.Rect, 0)

		for inputIndex, input := range composer.inputs {
			for _, rect := range rects {
				out_rect := rect
				out_rect.X += int16(ofsX)
				out_rect.Y += int16(ofsY)

				in_rect := out_rect
				in_rect.X -= int16(input.x)
				in_rect.Y -= int16(input.y)

				if inputIndex == 0 {
					updateRects = append(updateRects, clip(out_rect, output))
				}

				output.Blit(&out_rect, input.surface, &in_rect)
			}
		}

		if composer.showPaintedRegions {
			R := rnd.Float32()
			G := rnd.Float32() * (1.0 - R)
			B := rnd.Float32() * (1.0 - R - G)
			color := (uint32(R*0xFF) << 16) | (uint32(G*0xFF) << 8) | (uint32(B*0xFF) << 0)
			const alpha = 0x80

			for _, updateRect := range updateRects {
				fillRect(&SDLSurface{output}, updateRect, color, alpha)
			}
		}

		output.UpdateRects(updateRects)
	}
}

// Fills a rectangle with the specified RGB color and alpha value.
//
// The value of the alpha channel is from 0 to 256 (including 256).
// The value 256 maps to an alpha value of 1.0.
func fillRect(surface *SDLSurface, r sdl.Rect, RGB uint32, A uint32) {
	R := uint32((RGB >> 16) & 0xFF)
	G := uint32((RGB >> 8) & 0xFF)
	B := uint32((RGB >> 0) & 0xFF)

	RA := R * A
	GA := G * A
	BA := B * A
	A100 := 0x100 - A

	bpp := surface.Bpp()
	end_x := uint(r.X) + uint(r.W)
	end_y := uint(r.Y) + uint(r.H)

	for y := uint(r.Y); y < end_y; y++ {
		addr := surface.addrXY(uint(r.X), y)
		for x := uint(r.X); x < end_x; x++ {
			var c uint32 = *(*uint32)(unsafe.Pointer(addr))

			cR := (c >> 16) & 0xFF
			cG := (c >> 8) & 0xFF
			cB := (c >> 0) & 0xFF

			cR = ((cR * A100) + RA) >> 8
			cG = ((cG * A100) + GA) >> 8
			cB = ((cB * A100) + BA) >> 8

			c = (cR << 16) | (cG << 8) | cB
			*(*uint32)(unsafe.Pointer(addr)) = c

			addr += uintptr(bpp)
		}
	}
}

// Clips 'rect' to the dimensions of 'surface'.
// Returns the clipped rectangle.
func clip(rect sdl.Rect, surface *sdl.Surface) sdl.Rect {
	if (rect.X >= int16(surface.W)) || (rect.Y >= int16(surface.H)) {
		return sdl.Rect{}
	}

	clippedRect := rect

	if clippedRect.X < 0 {
		w := int16(clippedRect.W) - (-clippedRect.X)
		if w <= 0 {
			return sdl.Rect{}
		}

		clippedRect.X = 0
		clippedRect.W = uint16(w)
	}
	if clippedRect.Y < 0 {
		h := int16(clippedRect.H) - (-clippedRect.Y)
		if h <= 0 {
			return sdl.Rect{}
		}

		clippedRect.Y = 0
		clippedRect.H = uint16(h)
	}
	if (clippedRect.X + int16(clippedRect.W)) > int16(surface.W) {
		clippedRect.W = uint16(int16(surface.W) - clippedRect.X)
	}
	if (clippedRect.Y + int16(clippedRect.H)) > int16(surface.H) {
		clippedRect.H = uint16(int16(surface.H) - clippedRect.Y)
	}

	return clippedRect
}
