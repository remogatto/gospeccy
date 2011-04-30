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
	"spectrum/interpreter"
	"spectrum/output"
	"⚛sdl"
	"⚛sdl/ttf"
	"fmt"
	"flag"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"clingon"
)

const DEFAULT_JOYSTICK_ID = 0

var (
	// The speccy instance
	speccy *spectrum.Spectrum48k

	// The CLI
	cli *clingon.Console

	// The application renderer
	r *SDLRenderer

	joystick *sdl.Joystick

	composer *output.SDLSurfaceComposer
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

	sliderCh chan cmd_newSlider

	audio     bool
	audioFreq uint
	hqAudio   bool
}

// Passed to interpreter.Init()
type DummyRenderer struct {
	scale2x    *bool
	fullscreen *bool

	audio     *bool
	audioFreq *uint
	hqAudio   *bool
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

func newAppSurface(app *spectrum.Application, scale2x, fullscreen bool) SDLSurfaceAccessor {
	var sdlMode uint32
	if fullscreen {
		scale2x = true
		sdlMode = sdl.FULLSCREEN
		sdl.ShowCursor(sdl.DISABLE)
	} else {
		sdl.ShowCursor(sdl.ENABLE)
		sdlMode = sdl.SWSURFACE
	}

	<-composer.ReplaceOutputSurface(nil)

	surface := sdl.SetVideoMode(int(width(scale2x, fullscreen)), int(height(scale2x, fullscreen)), 32, sdlMode)
	if app.Verbose {
		app.PrintfMsg("video surface resolution: %dx%d", surface.W, surface.H)
	}

	<-composer.ReplaceOutputSurface(surface)

	return &wrapSurface{surface}
}

func newSpeccySurface(app *spectrum.Application, scale2x, fullscreen bool) SDLSurfaceAccessor {
	var speccySurface SDLSurfaceAccessor
	if fullscreen {
		scale2x = true
	}
	if scale2x {
		sdlScreen := output.NewSDLScreen2x(app)
		speccy.CommandChannel <- spectrum.Cmd_AddDisplay{sdlScreen}
		speccySurface = sdlScreen
	} else {
		sdlScreen := output.NewSDLScreen(app)
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

func NewSDLRenderer(app *spectrum.Application, scale2x, fullscreen bool, audio, hqAudio bool, audioFreq uint) *SDLRenderer {
	width := width(scale2x, fullscreen)
	height := height(scale2x, fullscreen)
	r := &SDLRenderer{
		app:              app,
		scale2x:          scale2x,
		fullscreen:       fullscreen,
		appSurfaceCh:     make(chan cmd_newSurface),
		speccySurfaceCh:  make(chan cmd_newSurface),
		cliSurfaceCh:     make(chan cmd_newCliSurface),
		appSurface:       newAppSurface(app, scale2x, fullscreen),
		speccySurface:    newSpeccySurface(app, scale2x, fullscreen),
		cliSurface_orNil: nil,
		width:            width,
		height:           height,
		consoleY:         int16(height),
		sliderCh:         make(chan cmd_newSlider),
		audio:            audio,
		audioFreq:        audioFreq,
		hqAudio:          hqAudio,
	}

	composer.AddInputSurface(r.speccySurface.GetSurface(), 0, 0, r.speccySurface.UpdatedRectsCh())

	go r.loop()
	return r
}

func (r *SDLRenderer) ResizeVideo(scale2x, fullscreen bool) {
	finished := make(chan byte)
	speccy.CommandChannel <- spectrum.Cmd_CloseAllDisplays{finished}
	<-finished

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
	r.appSurfaceCh <- cmd_newSurface{newAppSurface(r.app, scale2x, fullscreen), done}
	<-done

	r.speccySurfaceCh <- cmd_newSurface{newSpeccySurface(r.app, scale2x, fullscreen), done}
	<-done

	if r.cliSurface_orNil != nil {
		r.cliSurfaceCh <- cmd_newCliSurface{newCLISurface(scale2x, fullscreen), done}
		<-done
	}
}

func (r *SDLRenderer) ShowPaintedRegions(enable bool) {
	composer.ShowPaintedRegions(enable)
}

func (r *SDLRenderer) setAudioParameters(enable, hqAudio bool, freq uint) {
	r.audio = enable
	r.hqAudio = hqAudio
	r.audioFreq = freq

	finished := make(chan byte)
	speccy.CommandChannel <- spectrum.Cmd_CloseAllAudioReceivers{finished}
	<-finished

	if enable {
		audio, err := output.NewSDLAudio(r.app, freq, hqAudio)
		if err == nil {
			finished := make(chan byte)
			speccy.CommandChannel <- spectrum.Cmd_CloseAllAudioReceivers{finished}
			<-finished

			speccy.CommandChannel <- spectrum.Cmd_AddAudioReceiver{audio}
		} else {
			r.app.PrintfMsg("%s", err)
			return
		}
	}
}

func (r *SDLRenderer) EnableAudio(enable bool) {
	r.setAudioParameters(enable, r.hqAudio, r.audioFreq)
}

func (r *SDLRenderer) SetAudioFreq(freq uint) {
	if r.audioFreq != freq {
		r.setAudioParameters(r.audio, r.hqAudio, freq)
	}
}

func (r *SDLRenderer) SetAudioQuality(hqAudio bool) {
	if r.hqAudio != hqAudio {
		r.setAudioParameters(r.audio, hqAudio, r.audioFreq)
	}
}

func (r *DummyRenderer) ResizeVideo(scale2x, fullscreen bool) {
	// Overwrite the command-line settings
	*r.scale2x = scale2x
	*r.fullscreen = fullscreen
}

func (r *DummyRenderer) ShowPaintedRegions(enable bool) {
	composer.ShowPaintedRegions(enable)
}

func (r *DummyRenderer) EnableAudio(enable bool) {
	// Overwrite the command-line settings
	*r.audio = enable
}

func (r *DummyRenderer) SetAudioFreq(freq uint) {
	// Overwrite the command-line settings
	*r.audioFreq = freq
}

func (r *DummyRenderer) SetAudioQuality(hqAudio bool) {
	// Overwrite the command-line settings
	*r.hqAudio = hqAudio
}

// Synchronously destroy the CLI renderer
func (r *SDLRenderer) destroyCliRenderer() {
	if r.cliSurface_orNil != nil {
		cliSurface := r.cliSurface_orNil
		r.cliSurface_orNil = nil

		cli.SetRenderer(nil)

		<-composer.RemoveInputSurface(cliSurface.GetSurface())

		done := make(chan bool)
		cliSurface.EventCh() <- clingon.Cmd_Terminate{done}
		<-done

		cliSurface.GetSurface().Free()
		cliSurface.Font.Close()
	}
}

func (r *SDLRenderer) loop() {
	var slider_orNil *clingon.Animation = nil
	var sliderTargetState int = -1
	var sliderValueCh_orNil <-chan float64 = nil
	var sliderFinishedCh_orNil <-chan bool = nil

	evtLoop := r.app.NewEventLoop()

	for {
		select {
		case <-evtLoop.Pause:
			if slider_orNil != nil {
				slider_orNil.Terminate()
				<-sliderFinishedCh_orNil

				slider_orNil = nil
				sliderTargetState = -1
				sliderValueCh_orNil = nil
				sliderFinishedCh_orNil = nil
			}

			r.destroyCliRenderer()

			evtLoop.Pause <- 0

		case <-evtLoop.Terminate:
			// Terminate this Go routine
			if r.app.Verbose {
				r.app.PrintfMsg("frontend SDL renderer event loop: exit")
			}
			evtLoop.Terminate <- 0
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
			composer.SetPosition(r.cliSurface_orNil.GetSurface(), 0, int(r.consoleY))

		case <-sliderFinishedCh_orNil:
			if sliderTargetState == HIDE {
				r.destroyCliRenderer()
			}

			slider_orNil = nil
			sliderTargetState = -1
			sliderValueCh_orNil = nil
			sliderFinishedCh_orNil = nil

			r.toggling = false

		case cmd := <-r.cliSurfaceCh:
			r.destroyCliRenderer()

			r.cliSurface_orNil = cmd.surface_orNil
			cli.SetRenderer(cmd.surface_orNil)

			if r.cliSurface_orNil != nil {
				composer.AddInputSurface(r.cliSurface_orNil.GetSurface(), 0, int(r.consoleY), r.cliSurface_orNil.UpdatedRectsCh())
			}

			cmd.done <- true

		case cmd := <-r.speccySurfaceCh:
			<-composer.RemoveAllInputSurfaces()

			r.speccySurface.GetSurface().Free()
			r.speccySurface = cmd.surface

			composer.AddInputSurface(r.speccySurface.GetSurface(), 0, 0, r.speccySurface.UpdatedRectsCh())
			if r.cliSurface_orNil != nil {
				composer.AddInputSurface(r.cliSurface_orNil.GetSurface(), 0, int(r.consoleY), r.cliSurface_orNil.UpdatedRectsCh())
			}
			cmd.done <- true

		case cmd := <-r.appSurfaceCh:
			<-composer.ReplaceOutputSurface(nil)

			r.appSurface.GetSurface().Free()
			r.appSurface = cmd.surface

			<-composer.ReplaceOutputSurface(r.appSurface.GetSurface())

			cmd.done <- true
		}
	}
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
func sdlEventLoop(app *spectrum.Application, speccy *spectrum.Spectrum48k, verboseInput bool) {
	evtLoop := app.NewEventLoop()

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
					app.PrintfMsg("%v: %v", e.Keysym.Sym, keyName)
					app.PrintfMsg("Type: %02x Which: %02x State: %02x\n", e.Type, e.Which, e.State)
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

func createApplication(verbose bool) *spectrum.Application {
	app := spectrum.NewApplication()
	app.Verbose = verbose
	return app
}

// Create new emulator core
func initEmulationCore(app *spectrum.Application, acceleratedLoad bool) os.Error {
	romPath := spectrum.SystemRomPath("48.rom")
	rom, err := spectrum.ReadROM(romPath)
	if err != nil {
		return err
	}

	speccy = spectrum.NewSpectrum48k(app, *rom)
	if acceleratedLoad {
		speccy.TapeDrive().AcceleratedLoad = true
	}

	return nil
}

func initSDLSubSystems(app *spectrum.Application) os.Error {
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
	audio := flag.Bool("audio", true, "Enable or disable audio")
	audioFreq := flag.Uint("audio-freq", output.PLAYBACK_FREQUENCY, "Audio playback frequency (units: Hz)")
	hqAudio := flag.Bool("audio-hq", true, "Enable or disable higher-quality audio")
	acceleratedLoad := flag.Bool("accelerated-load", false, "Enable or disable accelerated tapes loading")
	showPaintedRegions := flag.Bool("show-paint", false, "Show painted display regions")
	verbose := flag.Bool("verbose", false, "Enable debugging messages")
	verboseInput := flag.Bool("verbose-input", false, "Enable debugging messages (input device events)")
	cpuProfile := flag.String("hostcpu-profile", "", "Write host-CPU profile to the specified file (for 'pprof')")

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

	// Start host-CPU profiling (if enabled).
	// The setup code is based on the contents of Go's file "src/pkg/testing/testing.go".
	var pprof_file *os.File
	if *cpuProfile != "" {
		var err os.Error

		pprof_file, err = os.Create(*cpuProfile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
			return
		}

		err = pprof.StartCPUProfile(pprof_file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to start host-CPU profiling: %s", err)
			pprof_file.Close()
			return
		}
	}

	app := createApplication(*verbose)

	composer = output.NewSDLSurfaceComposer(app)
	composer.ShowPaintedRegions(*showPaintedRegions)

	// Use at least 2 OS threads.
	// This helps to prevent audio buffer underflows
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

	SDL_initialized := false

	if err := initEmulationCore(app, *acceleratedLoad); err != nil {
		app.PrintfMsg("%s", err)
		app.RequestExit()
		goto quit
	}

	// Run startup scripts.
	// The startup scripts may change the display settings or enable/disable the audio.
	// They may also terminate the program.
	{
		dummyRenderer := DummyRenderer{
			scale2x:    scale2x,
			fullscreen: fullscreen,
			audio:      audio,
			audioFreq:  audioFreq,
			hqAudio:    hqAudio,
		}
		interpreter.Init(app, flag.Arg(0), speccy, &dummyRenderer)

		if app.TerminationInProgress() || app.Terminated() {
			goto quit
		}
	}

	// Optional: Read and categorize the contents
	//           of the file specified on the command-line
	var program_orNil interface{} = nil
	if flag.Arg(0) != "" {
		file := flag.Arg(0)
		path := spectrum.ProgramPath(file)

		var err os.Error
		program_orNil, err = formats.ReadProgram(path)
		if err != nil {
			app.PrintfMsg("read %s: %s", file, err)
			app.RequestExit()
			goto quit
		}
	}

	// SDL subsystems init
	if err := initSDLSubSystems(app); err != nil {
		app.PrintfMsg("%s", err)
		app.RequestExit()
		goto quit
	} else {
		SDL_initialized = true
	}

	{
		n := make(chan uint)

		// Setup the display

		speccy.CommandChannel <- spectrum.Cmd_GetNumDisplayReceivers{n}
		if <-n == 0 {
			r = NewSDLRenderer(app, *scale2x, *fullscreen, *audio, *hqAudio, *audioFreq)
			interpreter.SetUI(r)
		}

		initCLI()

		// Setup the audio
		speccy.CommandChannel <- spectrum.Cmd_GetNumAudioReceivers{n}
		numAudioReceivers := <-n
		if *audio && (numAudioReceivers == 0) {
			audio, err := output.NewSDLAudio(app, *audioFreq, *hqAudio)
			if err == nil {
				speccy.CommandChannel <- spectrum.Cmd_AddAudioReceiver{audio}
			} else {
				app.PrintfMsg("%s", err)
			}
		}
	}

	// Start the SDL event loop
	go sdlEventLoop(app, speccy, *verboseInput)

	// Begin speccy emulation
	go speccy.EmulatorLoop()

	// Set the FPS
	speccy.CommandChannel <- spectrum.Cmd_SetFPS{float32(*fps), nil}

	// Optional: Load the program specified on the command-line
	if program_orNil != nil {
		program := program_orNil
		file := flag.Arg(0)

		if _, isTAP := program.(*formats.TAP); isTAP {
			romLoaded := make(chan (<-chan bool))
			speccy.CommandChannel <- spectrum.Cmd_Reset{romLoaded}
			<-(<-romLoaded)
		}

		errChan := make(chan os.Error)
		speccy.CommandChannel <- spectrum.Cmd_Load{file, program, errChan}
		err := <-errChan
		if err != nil {
			app.PrintfMsg("%s", err)
			app.RequestExit()
			goto quit
		}
	}

	hint := ""
	hint += "Hint: Press F10 to invoke the built-in console.\n"
	hint += "      Input an empty line in the console to display available commands.\n"
	fmt.Print(hint)

quit:
	<-app.HasTerminated
	if SDL_initialized {
		sdl.Quit()
	}

	if app.Verbose {
		app.PrintfMsg("GC: %d garbage collections, %f ms total pause time",
			runtime.MemStats.NumGC, float64(runtime.MemStats.PauseTotalNs)/1e6)
	}

	// Stop host-CPU profiling
	if *cpuProfile != "" {
		pprof.StopCPUProfile() // flushes profile to disk
	}
}
