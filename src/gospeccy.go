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

package main

import (
	"io/ioutil"
	"spectrum"
	"spectrum/formats"
	"spectrum/interpreter"
	"⚛sdl"
	"⚛sdl/ttf"
	"fmt"
	"flag"
	"os"
	"os/signal"
	"runtime"
	"clingon"
)

const DEFAULT_JOYSTICK_ID = 0

var (
	// The application instance
	app *spectrum.Application

	// The speccy instance
	speccy *spectrum.Spectrum48k

	// The CLI
	cli *clingon.Console

	// The application renderer
	r *SDLRenderer

	joystick *sdl.Joystick
)

type SDLSurfaceAccessor interface {
	UpdatedRectsCh() <-chan []sdl.Rect
	GetSurface() *sdl.Surface
}

type cmd_newSurface struct {
	surface SDLSurfaceAccessor
	done    chan bool
}

type cmd_newCliSurface struct {
	surface_orNil *clingon.SDLRenderer
	done          chan bool
}

const (
	HIDE = iota
	SHOW
)

type cmd_newSlider struct {
	anim        *clingon.Animation
	targetState int // Either HIDE or SHOW
}

type SDLRenderer struct {
	app                           *spectrum.Application
	scale2x, fullscreen           bool
	consoleY                      int16
	width, height                 int
	appSurface, speccySurface     SDLSurfaceAccessor
	cliSurface_orNil              *clingon.SDLRenderer
	toggling                      bool
	appSurfaceCh, speccySurfaceCh chan cmd_newSurface
	cliSurfaceCh                  chan cmd_newCliSurface
	evtLoop                       *spectrum.EventLoop

	sliderCh chan cmd_newSlider
}

type wrapSurface struct {
	surface *sdl.Surface
}

func (s *wrapSurface) GetSurface() *sdl.Surface {
	return s.surface
}

func (s *wrapSurface) UpdatedRectsCh() <-chan []sdl.Rect {
	return nil
}

func width(scale2x, fullscreen bool) int {
	if fullscreen {
		scale2x = true
	}
	if scale2x {
		return spectrum.TotalScreenWidth * 2
	}
	return spectrum.TotalScreenWidth
}

func height(scale2x, fullscreen bool) int {
	if fullscreen {
		scale2x = true
	}
	if scale2x {
		return spectrum.TotalScreenHeight * 2
	}
	return spectrum.TotalScreenHeight
}

func newAppSurface(scale2x, fullscreen bool) SDLSurfaceAccessor {
	var sdlMode uint32
	if fullscreen {
		scale2x = true
		sdlMode = sdl.FULLSCREEN
		sdl.ShowCursor(sdl.DISABLE)
	} else {
		sdl.ShowCursor(sdl.ENABLE)
		sdlMode = sdl.SWSURFACE
	}
	surface := sdl.SetVideoMode(int(width(scale2x, fullscreen)), int(height(scale2x, fullscreen)), 32, sdlMode)
	if app.Verbose {
		app.PrintfMsg("video surface resolution: %dx%d", surface.W, surface.H)
	}
	return &wrapSurface{surface}
}

func newSpeccySurface(app *spectrum.Application, scale2x, fullscreen bool) SDLSurfaceAccessor {
	var speccySurface SDLSurfaceAccessor
	if fullscreen {
		scale2x = true
	}
	if scale2x {
		sdlScreen := spectrum.NewSDLScreen2x(app)
		speccy.CommandChannel <- spectrum.Cmd_AddDisplay{sdlScreen}
		speccySurface = sdlScreen
	} else {
		sdlScreen := spectrum.NewSDLScreen(app)
		speccy.CommandChannel <- spectrum.Cmd_AddDisplay{sdlScreen}
		speccySurface = sdlScreen
	}
	return speccySurface
}

func newCLISurface(scale2x, fullscreen bool) *clingon.SDLRenderer {
	cliSurface := clingon.NewSDLRenderer(
		sdl.CreateRGBSurface(
			sdl.SRCALPHA,
			width(scale2x, fullscreen),
			height(scale2x, fullscreen)/2, 32, 0, 0, 0, 0),
		newFont(scale2x, fullscreen),
	)
	cliSurface.GetSurface().SetAlpha(sdl.SRCALPHA, 0xdd)
	return cliSurface
}

func newFont(scale2x, fullscreen bool) *ttf.Font {
	var font *ttf.Font
	if fullscreen {
		scale2x = true
	}
	if scale2x {
		font = ttf.OpenFont(spectrum.FontPath("VeraMono.ttf"), 12)
	} else {
		font = ttf.OpenFont(spectrum.FontPath("VeraMono.ttf"), 10)
	}
	if font == nil {
		panic(sdl.GetError())
	}
	return font
}

func NewSDLRenderer(app *spectrum.Application, scale2x, fullscreen bool) *SDLRenderer {
	width := width(scale2x, fullscreen)
	height := height(scale2x, fullscreen)
	r := &SDLRenderer{
		app:              app,
		scale2x:          scale2x,
		fullscreen:       fullscreen,
		appSurfaceCh:     make(chan cmd_newSurface),
		speccySurfaceCh:  make(chan cmd_newSurface),
		cliSurfaceCh:     make(chan cmd_newCliSurface),
		appSurface:       newAppSurface(scale2x, fullscreen),
		speccySurface:    newSpeccySurface(app, scale2x, fullscreen),
		cliSurface_orNil: nil,
		width:            width,
		height:           height,
		consoleY:         int16(height),
		sliderCh:         make(chan cmd_newSlider),
	}
	go r.loop()
	return r
}

func (r *SDLRenderer) Resize(app *spectrum.Application, scale2x, fullscreen bool) {
	if r.scale2x != scale2x {
		if scale2x {
			// 1x --> 2x
			y := int16(r.height) - r.consoleY
			r.consoleY = int16(2*r.height) - 2*y
		} else {
			// 2x --> 1x
			y := int16(r.height) - r.consoleY
			r.consoleY = int16(r.height/2) - y/2
		}
	}

	r.width = width(scale2x, fullscreen)
	r.height = height(scale2x, fullscreen)
	r.scale2x = scale2x
	r.fullscreen = fullscreen

	done := make(chan bool)
	r.appSurfaceCh <- cmd_newSurface{newAppSurface(scale2x, fullscreen), done}
	<-done

	done = make(chan bool)
	r.speccySurfaceCh <- cmd_newSurface{newSpeccySurface(app, scale2x, fullscreen), done}
	<-done

	if r.cliSurface_orNil != nil {
		done = make(chan bool)
		r.cliSurfaceCh <- cmd_newCliSurface{newCLISurface(scale2x, fullscreen), done}
		<-done
	}
}

// Synchronously destroy the CLI renderer
func (r *SDLRenderer) destroyCliRenderer() {
	if r.cliSurface_orNil != nil {
		cliSurface := r.cliSurface_orNil
		r.cliSurface_orNil = nil

		cli.SetRenderer(nil)

		go func() {
			for r := range cliSurface.UpdatedRectsCh() {
				if r == nil {
					break
				}
			}
		}()

		done := make(chan bool)
		cliSurface.EventCh() <- clingon.Cmd_Terminate{done}
		<-done

		cliSurface.GetSurface().Free()
		cliSurface.Font.Close()
	}
}

func (r *SDLRenderer) loop() {
	var cliSurface_updatedRectsCh_orNil <-chan []sdl.Rect = nil

	var slider_orNil *clingon.Animation = nil
	var sliderTargetState int = -1
	var sliderValueCh_orNil <-chan float64 = nil
	var sliderFinishedCh_orNil <-chan bool = nil

	r.evtLoop = r.app.NewEventLoop()
	for {
		select {
		case <-r.evtLoop.Pause:
			if slider_orNil != nil {
				slider_orNil.Terminate()
				<-sliderFinishedCh_orNil

				slider_orNil = nil
				sliderTargetState = -1
				sliderValueCh_orNil = nil
				sliderFinishedCh_orNil = nil
			}

			r.destroyCliRenderer()
			cliSurface_updatedRectsCh_orNil = nil

			r.evtLoop.Pause <- 0

		case <-r.evtLoop.Terminate:
			// Terminate this Go routine
			if app.Verbose {
				app.PrintfMsg("frontend SDL renderer event loop: exit")
			}
			r.evtLoop.Terminate <- 0
			return

		case cmd := <-r.sliderCh:
			slider_orNil = cmd.anim
			sliderTargetState = cmd.targetState
			sliderValueCh_orNil = cmd.anim.ValueCh()
			sliderFinishedCh_orNil = cmd.anim.FinishedCh()

		case value := <-sliderValueCh_orNil:
			if sliderTargetState == HIDE {
				r.consoleY = int16(float64(r.height/2) + value*float64(r.height/2))
			} else {
				r.consoleY = int16(float64(r.height) - value*float64(r.height/2))
			}
			r.blitAll()

		case <-sliderFinishedCh_orNil:
			if sliderTargetState == HIDE {
				r.destroyCliRenderer()
				cliSurface_updatedRectsCh_orNil = nil

				r.blitAll()
			}

			slider_orNil = nil
			sliderTargetState = -1
			sliderValueCh_orNil = nil
			sliderFinishedCh_orNil = nil

			r.toggling = false

		case cmd := <-r.cliSurfaceCh:
			r.destroyCliRenderer()
			cliSurface_updatedRectsCh_orNil = nil

			r.cliSurface_orNil = cmd.surface_orNil
			cli.SetRenderer(cmd.surface_orNil)

			if cmd.surface_orNil != nil {
				cliSurface_updatedRectsCh_orNil = cmd.surface_orNil.UpdatedRectsCh()
			} else {
				cliSurface_updatedRectsCh_orNil = nil
			}

			cmd.done <- true
			r.blitAll()

		case cmd := <-r.speccySurfaceCh:
			r.speccySurface.GetSurface().Free()
			r.speccySurface = cmd.surface
			cmd.done <- true
			r.blitAll()

		case cmd := <-r.appSurfaceCh:
			r.appSurface.GetSurface().Free()
			r.appSurface = cmd.surface
			cmd.done <- true
			r.blitAll()

		case cliRects := <-cliSurface_updatedRectsCh_orNil:
			r.render(nil, cliRects)

		case speccyRects := <-r.speccySurface.UpdatedRectsCh():
			r.render(speccyRects, nil)
		}
	}
}

func (r *SDLRenderer) blitAll() {
	appSdlSurface := r.appSurface.GetSurface()

	appSdlSurface.Blit(nil, r.speccySurface.GetSurface(), nil)
	if r.cliSurface_orNil != nil {
		cliSurface := r.cliSurface_orNil
		appSdlSurface.Blit(&sdl.Rect{0, int16(r.consoleY), 0, 0}, cliSurface.GetSurface(), nil)
	}

	appSdlSurface.Flip()
}

func (r *SDLRenderer) render(speccyRects, cliRects []sdl.Rect) {
	appSdlSurface := r.appSurface.GetSurface()

	var cliSdlSurface_orNil *sdl.Surface = nil
	if r.cliSurface_orNil != nil {
		cliSdlSurface_orNil = r.cliSurface_orNil.GetSurface()
	}

	for _, rect := range speccyRects {
		// 'x' and 'y' are relative positions in respect to 'appSdlSurface'
		x, y, w, h := rect.X, rect.Y, rect.W, rect.H
		appSdlSurface.Blit(&rect, r.speccySurface.GetSurface(), &rect)

		if cliSdlSurface_orNil != nil {
			cy := int16(r.consoleY)
			appSdlSurface.Blit(&sdl.Rect{x, y, 0, 0}, cliSdlSurface_orNil, &sdl.Rect{x, y - cy, w, h})
		}
	}

	for _, rect := range cliRects {
		// 'x' and 'y' are relative positions in respect to 'appSdlSurface'
		x, y, w, h := rect.X, rect.Y+int16(r.consoleY), rect.W, rect.H
		appSdlSurface.Blit(&sdl.Rect{x, y, 0, 0}, r.speccySurface.GetSurface(), &sdl.Rect{x, y, w, h})

		if cliSdlSurface_orNil != nil {
			cy := int16(r.consoleY)
			appSdlSurface.Blit(&sdl.Rect{x, y, 0, 0}, cliSdlSurface_orNil, &sdl.Rect{x, y - cy, w, h})
		}
	}

	appSdlSurface.Flip()
}

func initCLI() {
	cli = clingon.NewConsole(&interpreter.Interpreter{})
	cli.Print(`
GoSpeccy Command Line Interface (CLI)
-------------------------------------
Available keys:
* F10 toggle/untoggle the CLI
* Up/Down for history browsing
* PageUp/PageDown for scrolling
`)
	cli.SetPrompt("gospeccy> ")
}

// A Go routine for processing SDL events.
func sdlEventLoop(evtLoop *spectrum.EventLoop, speccy *spectrum.Spectrum48k, verboseInput bool) {
	app = evtLoop.App()

	consoleIsVisible := false

	for {
		select {
		case <-evtLoop.Pause:
			evtLoop.Pause <- 0

		case <-evtLoop.Terminate:
			// Terminate this Go routine
			if app.Verbose {
				app.PrintfMsg("SDL event loop: exit")
			}
			evtLoop.Terminate <- 0
			return

		case event := <-sdl.Events:
			switch e := event.(type) {
			case sdl.QuitEvent:
				if app.Verbose {
					app.PrintfMsg("SDL quit -> request[exit the application]")
				}
				app.RequestExit()

			case sdl.JoyAxisEvent:
				if verboseInput {
					app.PrintfMsg("[Joystick] Axis: %d, Value: %d", e.Axis, e.Value)
				}
				if e.Axis == 0 {
					if e.Value > 0 {
						speccy.Joystick.KempstonDown(spectrum.KEMPSTON_RIGHT)
					} else if e.Value < 0 {
						speccy.Joystick.KempstonDown(spectrum.KEMPSTON_LEFT)
					} else {
						speccy.Joystick.KempstonUp(spectrum.KEMPSTON_RIGHT)
						speccy.Joystick.KempstonUp(spectrum.KEMPSTON_LEFT)
					}
				} else if e.Axis == 1 {
					if e.Value > 0 {
						speccy.Joystick.KempstonDown(spectrum.KEMPSTON_UP)
					} else if e.Value < 0 {
						speccy.Joystick.KempstonDown(spectrum.KEMPSTON_DOWN)
					} else {
						speccy.Joystick.KempstonUp(spectrum.KEMPSTON_UP)
						speccy.Joystick.KempstonUp(spectrum.KEMPSTON_DOWN)
					}
				}

			case sdl.JoyButtonEvent:
				if verboseInput {
					app.PrintfMsg("[Joystick] Button: %d, State: %d", e.Button, e.State)
				}
				if e.Button == 0 {
					if e.State > 0 {
						speccy.Joystick.KempstonDown(spectrum.KEMPSTON_FIRE)
					} else {
						speccy.Joystick.KempstonUp(spectrum.KEMPSTON_FIRE)
					}
				}

			case sdl.KeyboardEvent:
				keyName := sdl.GetKeyName(sdl.Key(e.Keysym.Sym))

				if verboseInput {
					app.PrintfMsg("\n")
					app.PrintfMsg("%v: %v", e.Keysym.Sym, ": ", keyName)

					app.PrintfMsg("%04x ", e.Type)

					for i := 0; i < len(e.Pad0); i++ {
						app.PrintfMsg("%02x ", e.Pad0[i])
					}
					app.PrintfMsg("\n")

					app.PrintfMsg("Type: %02x Which: %02x State: %02x Pad: %02x\n", e.Type, e.Which, e.State, e.Pad0[0])
					app.PrintfMsg("Scancode: %02x Sym: %08x Mod: %04x Unicode: %04x\n", e.Keysym.Scancode, e.Keysym.Sym, e.Keysym.Mod, e.Keysym.Unicode)
				}

				if (keyName == "escape") && (e.Type == sdl.KEYDOWN) {
					if app.Verbose {
						app.PrintfMsg("escape key -> request[exit the application]")
					}
					app.RequestExit()

				} else if (keyName == "f10") && (e.Type == sdl.KEYDOWN) {
					//if app.Verbose {
					//	app.PrintfMsg("f10 key -> toggle console")
					//}
					if !r.toggling {
						r.toggling = true

						if r.cliSurface_orNil == nil {
							done := make(chan bool)
							r.cliSurfaceCh <- cmd_newCliSurface{newCLISurface(r.scale2x, r.fullscreen), done}
							<-done
						}

						anim := clingon.NewSliderAnimation(0.500, 1.0)

						var targetState int
						if consoleIsVisible {
							targetState = HIDE
						} else {
							targetState = SHOW
						}

						r.sliderCh <- cmd_newSlider{anim, targetState}
						anim.Start()

						consoleIsVisible = !consoleIsVisible
					}
				} else {
					if r.cliSurface_orNil != nil {
						cliSurface := r.cliSurface_orNil

						if (keyName == "page up") && (e.Type == sdl.KEYDOWN) {
							cliSurface.EventCh() <- clingon.Cmd_Scroll{clingon.SCROLL_UP}
						} else if (keyName == "page down") && (e.Type == sdl.KEYDOWN) {
							cliSurface.EventCh() <- clingon.Cmd_Scroll{clingon.SCROLL_DOWN}
						} else if (keyName == "up") && (e.Type == sdl.KEYDOWN) {
							cli.PutReadline(clingon.HISTORY_PREV)
						} else if (keyName == "down") && (e.Type == sdl.KEYDOWN) {
							cli.PutReadline(clingon.HISTORY_NEXT)
						} else if (keyName == "left") && (e.Type == sdl.KEYDOWN) {
							cli.PutReadline(clingon.CURSOR_LEFT)
						} else if (keyName == "right") && (e.Type == sdl.KEYDOWN) {
							cli.PutReadline(clingon.CURSOR_RIGHT)
						} else {
							unicode := e.Keysym.Unicode
							if unicode > 0 {
								cli.PutUnicode(unicode)
							}
						}
					} else {
						sequence, haveMapping := spectrum.SDL_KeyMap[keyName]

						if haveMapping {
							switch e.Type {
							case sdl.KEYDOWN:
								// Normal order
								for i := 0; i < len(sequence); i++ {
									speccy.Keyboard.KeyDown(sequence[i])
								}
							case sdl.KEYUP:
								// Reverse order
								for i := len(sequence) - 1; i >= 0; i-- {
									speccy.Keyboard.KeyUp(sequence[i])
								}
							}
						}
					}
				}
			}
		}
	}
}

type handler_SIGTERM struct {
	app *spectrum.Application
}

func (h *handler_SIGTERM) HandleSignal(s signal.Signal) {
	switch ss := s.(type) {
	case signal.UnixSignal:
		switch ss {
		case signal.SIGTERM, signal.SIGINT:
			if h.app.Verbose {
				h.app.PrintfMsg("%v", ss)
			}

			h.app.RequestExit()
		}
	}
}

func initApplication(verbose bool) {
	app = spectrum.NewApplication()
	app.Verbose = verbose
}

// Create new emulator core
func initEmulationCore(acceleratedLoad bool) os.Error {
	var rom [0x4000]byte
	{
		romPath := spectrum.SystemRomPath("48.rom")

		rom48k, err := ioutil.ReadFile(romPath)
		if err != nil {
			return err
		}
		if len(rom48k) != 0x4000 {
			return os.NewError(fmt.Sprintf("ROM file \"%s\" has an invalid size", romPath))
		}

		copy(rom[:], rom48k)
	}

	speccy = spectrum.NewSpectrum48k(app, rom)
	if acceleratedLoad {
		speccy.TapeDrive().AcceleratedLoad = true
	}

	return nil
}

func initSDLSubSystems() os.Error {
	if sdl.Init(sdl.INIT_VIDEO|sdl.INIT_AUDIO|sdl.INIT_JOYSTICK) != 0 {
		return os.NewError(sdl.GetError())
	}
	if ttf.Init() != 0 {
		return os.NewError(sdl.GetError())
	}
	if sdl.NumJoysticks() > 0 {
		// Open joystick
		joystick = sdl.JoystickOpen(DEFAULT_JOYSTICK_ID)
		if joystick != nil {
			if app.Verbose {
				app.PrintfMsg("Opened Joystick %d", DEFAULT_JOYSTICK_ID)
				app.PrintfMsg("Name: %s", sdl.JoystickName(DEFAULT_JOYSTICK_ID))
				app.PrintfMsg("Number of Axes: %d", joystick.NumAxes())
				app.PrintfMsg("Number of Buttons: %d", joystick.NumButtons())
				app.PrintfMsg("Number of Balls: %d", joystick.NumBalls())
			}
		} else {
			return os.NewError("Couldn't open Joystick!")
		}
	}
	sdl.WM_SetCaption("GoSpeccy - ZX Spectrum Emulator", "")
	sdl.EnableUNICODE(1)
	return nil
}

func main() {
	// Handle options
	help := flag.Bool("help", false, "Show usage")
	scale2x := flag.Bool("2x", false, "2x display scaler")
	fullscreen := flag.Bool("fullscreen", false, "Fullscreen (enable 2x scaler by default)")
	fps := flag.Float64("fps", spectrum.DefaultFPS, "Frames per second")
	sound := flag.Bool("sound", true, "Enable or disable sound")
	acceleratedLoad := flag.Bool("accelerated-load", false, "Enable or disable accelerated tapes loading")
	verbose := flag.Bool("verbose", false, "Enable debugging messages")
	verboseInput := flag.Bool("verbose-input", false, "Enable debugging messages (input device events)")

	{
		flag.Usage = func() {
			fmt.Fprintf(os.Stderr, "GoSpeccy - A ZX Spectrum 48k Emulator written in Go\n\n")
			fmt.Fprintf(os.Stderr, "Usage:\n\n")
			fmt.Fprintf(os.Stderr, "\tgospeccy [options] [image.sna]\n\n")
			fmt.Fprintf(os.Stderr, "Options are:\n\n")
			flag.PrintDefaults()
		}

		flag.Parse()

		if *help == true {
			flag.Usage()
			return
		}
	}

	initApplication(*verbose)

	// Use at least 2 OS threads.
	// This helps to prevent sound buffer underflows
	// in case SDL rendering is consuming too much CPU.
	if (os.Getenv("GOMAXPROCS") == "") && (runtime.GOMAXPROCS(-1) < 2) {
		runtime.GOMAXPROCS(2)
	}
	if app.Verbose {
		app.PrintfMsg("using %d OS threads", runtime.GOMAXPROCS(-1))
	}

	// Install SIGTERM handler
	{
		handler := handler_SIGTERM{app}
		spectrum.InstallSignalHandler(&handler)
	}

	if err := initEmulationCore(*acceleratedLoad); err != nil {
		app.PrintfMsg("%s", err)
		app.RequestExit()
		goto quit
	}

	// SDL subsystems init
	if err := initSDLSubSystems(); err != nil {
		app.PrintfMsg("%s", err)
		app.RequestExit()
		goto quit
	}

	{
		n := make(chan uint)

		// Setup the display

		speccy.CommandChannel <- spectrum.Cmd_GetNumDisplayReceivers{n}
		if <-n == 0 {
			r = NewSDLRenderer(app, *scale2x, *fullscreen)
		}

		initCLI()

		// Run startup scripts. The startup scripts may create a display/audio receiver.
		{
			fmt.Println("Hint: Press F10 to invoke the built-in console.")
			fmt.Println("      Input an empty line in the console to display available commands.")
			if app.TerminationInProgress() || closed(app.HasTerminated) {
				goto quit
			}
		}

		// Setup the audio
		speccy.CommandChannel <- spectrum.Cmd_GetNumAudioReceivers{n}
		numAudioReceivers := <-n
		if *sound && (numAudioReceivers == 0) {
			audio, err := spectrum.NewSDLAudio(app)
			if err == nil {
				speccy.CommandChannel <- spectrum.Cmd_AddAudioReceiver{audio}
			} else {
				app.PrintfMsg("%s", err)
			}
		}
	}

	// Start the SDL event loop
	go sdlEventLoop(app.NewEventLoop(), speccy, *verboseInput)

	// Begin speccy emulation
	go speccy.EmulatorLoop()

	// Set the FPS
	speccy.CommandChannel <- spectrum.Cmd_SetFPS{float32(*fps), nil}

	interpreter.Init(app, flag.Arg(0), speccy, r)

	// Process command line argument. Load the given program (if any)
	if flag.Arg(0) != "" {
		file := flag.Arg(0)

		path := spectrum.ProgramPath(file)

		program, err := formats.ReadProgram(path)
		if err != nil {
			app.PrintfMsg("%s", err)
			app.RequestExit()
			goto quit
		}

		if _, isTAP := program.(*formats.TAP); isTAP {
			romLoaded := make(chan (<-chan bool))
			speccy.CommandChannel <- spectrum.Cmd_Reset{romLoaded}
			<-(<-romLoaded)
		}

		errChan := make(chan os.Error)
		speccy.CommandChannel <- spectrum.Cmd_Load{file, program, errChan}
		err = <-errChan
		if err != nil {
			app.PrintfMsg("%s", err)
			app.RequestExit()
			goto quit
		}
	}

quit:
	<-app.HasTerminated
	sdl.Quit()

	if app.Verbose {
		app.PrintfMsg("GC: %d garbage collections, %f ms total pause time",
			runtime.MemStats.NumGC, float64(runtime.MemStats.PauseTotalNs)/1e6)
	}
}
