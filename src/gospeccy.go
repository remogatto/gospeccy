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
	"spectrum"
	"spectrum/formats"
	"⚛sdl"
	"⚛sdl/ttf"
	"fmt"
	"flag"
	"os"
	"os/signal"
	"runtime"
	"clingon"
)

var (
	// The application instance
	app *spectrum.Application

	// The speccy instance
	speccy *spectrum.Spectrum48k

	// The CLI renderer
	cliRenderer *clingon.SDLRenderer

	// The CLI
	cli *clingon.Console

	// The font used by the CLI
	font *ttf.Font

	// The application renderer
	r *SDLRenderer

	// CLI animations
	slideDown, slideUp *clingon.Animation
)

type SDLSurfaceAccessor interface {
	UpdatedRectsCh() <-chan []sdl.Rect
	GetSurface() *sdl.Surface
}

type SDLRenderer struct {
	app *spectrum.Application
	consoleX, consoleY                            int16
	consoleW, consoleH, consoleH_2, width, height uint16
	appSurface                                    *sdl.Surface
	speccySurface, cliSurface                     SDLSurfaceAccessor
	fps                                           float
	toggling                                      bool
	animationDuration                             int64
}

func NewSDLRenderer(app *spectrum.Application) *SDLRenderer {
	return &SDLRenderer{
	app: app,
	}
}

func (r *SDLRenderer) loop() {
	evtLoop := app.NewEventLoop()
	for {
		select {
		case <-evtLoop.Pause:
			cli.Pause(true)
			evtLoop.Pause <- 0

		case <-evtLoop.Terminate:
			// Terminate this Go routine
			if app.Verbose {
				app.PrintfMsg("SDL renderer event loop: exit")
			}
			evtLoop.Terminate <- 0
			return
		case value := <-slideDown.ValueCh():
			r.consoleY = int16(r.consoleH_2) + int16(value)
			r.render(nil, nil)
		case value := <-slideUp.ValueCh():
			r.consoleY = int16(r.consoleH_2) + int16(value)
			r.render(nil, nil)
		case <-slideDown.FinishedCh():
			r.toggling = false
		case <-slideUp.FinishedCh():
			r.toggling = false
		case cliRects := <-r.cliSurface.UpdatedRectsCh():
			r.render(nil, cliRects)
		case speccyRects := <-r.speccySurface.UpdatedRectsCh():
			r.render(speccyRects, nil)
		}
	}
}

func (r *SDLRenderer) render(speccyRects, cliRects []sdl.Rect) {
	if !r.toggling {
		if cli.Paused() {
			for _, rect := range speccyRects {
				r.appSurface.Blit(&rect, r.speccySurface.GetSurface(), &rect)
				r.appSurface.UpdateRect(int32(rect.X), int32(rect.Y), uint32(rect.W), uint32(rect.H))
			}
		} else {
			for _, rect := range speccyRects {
				x, y, w, h := rect.X, rect.Y - int16(r.consoleY), rect.W, rect.H
				r.appSurface.Blit(&rect, r.speccySurface.GetSurface(), &rect)
				r.appSurface.Blit(&sdl.Rect{x, y + int16(r.consoleY), 0, 0}, r.cliSurface.GetSurface(), &sdl.Rect{x, y, w, h})
			}
			for _, rect := range cliRects {
				x, y, w, h := rect.X, rect.Y + int16(r.consoleY), rect.W, rect.H
				r.appSurface.Blit(&sdl.Rect{x, y, 0, 0}, r.speccySurface.GetSurface(), &sdl.Rect{x, y, w, h})
				r.appSurface.Blit(&sdl.Rect{rect.X, rect.Y + int16(r.consoleY), 0, 0}, r.cliSurface.GetSurface(), &rect)
			}
			r.appSurface.Flip()
		}
	} else {
		r.appSurface.Blit(nil, r.speccySurface.GetSurface(), nil)
		r.appSurface.Blit(&sdl.Rect{0, int16(r.consoleY), 0, 0}, r.cliSurface.GetSurface(), nil)
		r.appSurface.Flip()
	}
}

// A Go routine for processing SDL events.
func sdlEventLoop(evtLoop *spectrum.EventLoop, speccy *spectrum.Spectrum48k, verboseKeyboard bool) {
	app = evtLoop.App()

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

			case sdl.KeyboardEvent:
				keyName := sdl.GetKeyName(sdl.Key(e.Keysym.Sym))

				if verboseKeyboard {
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
					if app.Verbose {
						app.PrintfMsg("f10 key -> toggle console")
					}
					if !r.toggling {
						cli.Pause(!cli.Paused())
						if cli.Paused() {
							slideDown.Start()
						} else {
							slideUp.Start()
						}
					}
					r.toggling = true
				} else {
					if !cli.Paused() {
						if (keyName == "page up") && (e.Type == sdl.KEYDOWN) {
							cliRenderer.ScrollCh() <- clingon.SCROLL_UP
						} else if (keyName == "page down") && (e.Type == sdl.KEYDOWN) {
							cliRenderer.ScrollCh() <- clingon.SCROLL_DOWN
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

func initSDLRenderer() {
	r = NewSDLRenderer(app)
}

// Create new emulator core
func initEmulationCore(acceleratedLoad bool) (err os.Error) {
	speccy, err = spectrum.NewSpectrum48k(app, spectrum.SystemRomPath("48.rom"))
	if acceleratedLoad {
		speccy.TapeDrive().AcceleratedLoad = true
	}
	return
}

func initSDLSubSystems() os.Error {
	if sdl.Init(sdl.INIT_VIDEO|sdl.INIT_AUDIO) != 0 {
		return os.NewError(sdl.GetError())
	}

	if ttf.Init() != 0 {
		return os.NewError(sdl.GetError())
	}
	sdl.WM_SetCaption("GoSpeccy - ZX Spectrum Emulator", "")
	sdl.EnableUNICODE(1)
	return nil
}

func initDisplay(scale2x, fullscreen bool) {
	var sdlMode uint32
	if fullscreen {
		sdlMode = sdl.FULLSCREEN
		scale2x = true
		sdl.ShowCursor(sdl.DISABLE)
	} else {
		sdlMode = 0
	}
	if scale2x {
		r.width = spectrum.TotalScreenWidth * 2
		r.height = spectrum.TotalScreenHeight * 2
		r.consoleH_2 = spectrum.TotalScreenHeight

		sdlScreen := spectrum.NewSDLScreen2x(app)
		speccy.CommandChannel <- spectrum.Cmd_AddDisplay{sdlScreen}
		r.speccySurface = sdlScreen
		initFont(12)
	} else {
		r.width = spectrum.TotalScreenWidth
		r.height = spectrum.TotalScreenHeight
		r.consoleH_2 = r.height / 2

		sdlScreen := spectrum.NewSDLScreen(app)
		speccy.CommandChannel <- spectrum.Cmd_AddDisplay{sdlScreen}
		r.speccySurface = sdlScreen
		initFont(10)
	}

	r.fps = spectrum.DefaultFPS
	r.consoleY = int16(r.height)
	r.appSurface = sdl.SetVideoMode(int(r.width), int(r.height), 32, sdlMode)
	slidingDistance := uint16(r.appSurface.H) - r.consoleH_2
	slideDown = clingon.NewSlideDownAnimation(500*1e6, float64(slidingDistance))
	slideUp = clingon.NewSlideUpAnimation(500*1e6, float64(slidingDistance))
}

func initFont(fontSize int) {
	font = ttf.OpenFont("font/VeraMono.ttf", fontSize)

	if font == nil {
		panic(sdl.GetError())
	}
}

func initCLI() {
	// Initialize CLI
	initInterpreter()
	cliRenderer = clingon.NewSDLRenderer(sdl.CreateRGBSurface(sdl.SRCALPHA, int(r.width), int(r.consoleH_2), 32, 0, 0, 0, 0), font)
	cliRenderer.GetSurface().SetAlpha(sdl.SRCALPHA, 0xdd)
	cli = clingon.NewConsole(cliRenderer, &i)
	cli.Print(`
Welcome to the GoSpeccy Command Line Interface (CLI)
----------------------------------------------------
Available keys:
* F10 toggle/untoggle the CLI
* PagUp/PagDown for scrolling
`)
	cli.SetPrompt("gospeccy> ")
	r.cliSurface = cliRenderer
}

func main() {
	// Use at least two OS threads. This helps to prevent sound
	// buffer underflows in case SDL rendering is consuming too
	// much CPU.
	if runtime.GOMAXPROCS(-1) < 2 {
		runtime.GOMAXPROCS(2)
	}

	// Handle options
	help := flag.Bool("help", false, "Show usage")
	scale2x := flag.Bool("2x", false, "2x display scaler")
	fullscreen := flag.Bool("fullscreen", false, "Fullscreen (enable 2x scaler by default)")
	fps := flag.Float("fps", spectrum.DefaultFPS, "Frames per second")
	sound := flag.Bool("sound", true, "Enable or disable sound")
	acceleratedLoad := flag.Bool("accelerated-load", false, "Enable or disable accelerated tapes loading")
	verbose := flag.Bool("verbose", false, "Enable debugging messages")
	verboseKeyboard := flag.Bool("verbose-keyboard", false, "Enable debugging messages (keyboard events)")

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
	initSDLRenderer()

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
			initDisplay(*scale2x, *fullscreen)
		}

		// Run startup scripts. The startup scripts may create a display/audio receiver.
		{
			fmt.Println("Hint: Press F10 to invoke the built-in console.")
			fmt.Println("      Input an empty line in the console to display available commands.")
			initCLI()
			if app.TerminationInProgress() || closed(app.HasTerminated) {
				goto quit
			}
		}

		// Setup the audio
		speccy.CommandChannel <- spectrum.Cmd_GetNumAudioReceivers{n}
		if *sound && (<-n == 0) {
			audio, err := spectrum.NewSDLAudio(app)
			if err == nil {
				speccy.CommandChannel <- spectrum.Cmd_AddAudioReceiver{audio}
			} else {
				app.PrintfMsg("%s", err)
			}
		}

		close(n)
	}

	// Start the SDL event loop
	go sdlEventLoop(app.NewEventLoop(), speccy, *verboseKeyboard)

	// Begin speccy emulation
	go speccy.EmulatorLoop()

	// Begin the application render loop
	go r.loop()

	// Set the FPS
	speccy.CommandChannel <- spectrum.Cmd_SetFPS{*fps}

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
			romLoaded := make(chan bool, 1)
			speccy.CommandChannel <- spectrum.Cmd_Reset{romLoaded}
			<-romLoaded
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

	// Drain systemROMLoaded channel
	<-speccy.ROMLoaded()

quit:
	<-app.HasTerminated
	sdl.Quit()
}
