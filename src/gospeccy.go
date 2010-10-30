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
	"spectrum/console"
	"⚛sdl"
	"fmt"
	"flag"
	"os"
	"os/signal"
	"runtime"
	"time"
)

// A Go routine for processing SDL events.
func sdlEventLoop(evtLoop *spectrum.EventLoop, speccy *spectrum.Spectrum48k, verboseKeyboard bool) {
	app := evtLoop.App()

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

func emulatorLoop(evtLoop *spectrum.EventLoop, speccy *spectrum.Spectrum48k) {
	app := evtLoop.App()

	fps := <-speccy.FPS
	ticker := time.NewTicker(int64(1e9 / fps))

	// Render the 1st frame (the 2nd frame will be rendered after 1/FPS seconds)
	{
		completionTime := make(chan int64)
		speccy.CommandChannel <- spectrum.Cmd_RenderFrame{completionTime}

		go func() {
			start := app.CreationTime
			end := <-completionTime
			if app.Verbose {
				app.PrintfMsg("first frame latency: %d ms", (end-start)/1e6)
			}
		}()
	}

	for {
		select {
		case <-evtLoop.Pause:
			ticker.Stop()
			spectrum.Drain(ticker)
			evtLoop.Pause <- 0

		case <-evtLoop.Terminate:
			// Terminate this Go routine
			if app.Verbose {
				app.PrintfMsg("emulator loop: exit")
			}
			evtLoop.Terminate <- 0
			return

		case <-ticker.C:
			//app.PrintfMsg("%d", time.Nanoseconds()/1e6)
			speccy.CommandChannel <- spectrum.Cmd_RenderFrame{}

		case FPS_new := <-speccy.FPS:
			if (FPS_new != fps) && (FPS_new > 0) {
				if app.Verbose {
					app.PrintfMsg("setting FPS to %f", FPS_new)
				}
				ticker.Stop()
				spectrum.Drain(ticker)
				ticker = time.NewTicker(int64(1e9 / FPS_new))
				fps = FPS_new
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

func main() {
	// Use at least two OS threads. This helps to prevent sound buffer underflows
	// in case SDL rendering is consuming too much CPU.
	if runtime.GOMAXPROCS(-1) < 2 {
		runtime.GOMAXPROCS(2)
	}

	help := flag.Bool("help", false, "Show usage")
	scale2x := flag.Bool("2x", false, "2x display scaler")
	fullscreen := flag.Bool("fullscreen", false, "Fullscreen (enable 2x scaler by default)")
	fps := flag.Float("fps", spectrum.DefaultFPS, "Frames per second")
	sound := flag.Bool("sound", true, "Enable or disable sound")
	verbose := flag.Bool("verbose", false, "Enable debugging messages")
	verboseKeyboard := flag.Bool("verbose-keyboard", false, "Enable debugging messages (keyboard events)")

	{
		flag.Usage = func() {
			fmt.Fprintf(os.Stderr, "GoSpeccy - A ZX Spectrum 48k Emulator written in GO\n\n")
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

	app := spectrum.NewApplication()
	app.Verbose = *verbose

	// Install SIGTERM handler
	{
		handler := handler_SIGTERM{app}
		spectrum.InstallSignalHandler(&handler)
	}

	// Create new emulator core
	speccy, err := spectrum.NewSpectrum48k(app, spectrum.SystemRomPath("48.rom"))
	if err != nil {
		app.PrintfMsg("%s", err)
		app.RequestExit()
		goto quit
	}

	if sdl.Init(sdl.INIT_VIDEO|sdl.INIT_AUDIO) != 0 {
		app.PrintfMsg("%s", sdl.GetError())
		app.RequestExit()
		goto quit
	}

	sdl.WM_SetCaption("GoSpeccy - ZX Spectrum Emulator", "")

	// Run startup scripts. The startup scripts may create a display/audio receiver.
	{
		console.Init(app, speccy)

		if app.TerminationInProgress() || closed(app.HasTerminated) {
			goto quit
		}
	}

	// Load snapshot (if any)
	if flag.Arg(0) != "" {
		file := flag.Arg(0)
		path := spectrum.SnaPath(file)

		snapshot, err := formats.ReadSnapshot(path)
		if err != nil {
			app.PrintfMsg("%s", err)
			app.RequestExit()
			goto quit
		}

		errChan := make(chan os.Error)
		speccy.CommandChannel <- spectrum.Cmd_LoadSnapshot{file, snapshot, errChan}
		err = <-errChan
		if err != nil {
			app.PrintfMsg("%s", err)
			app.RequestExit()
			goto quit
		}
	}

	{
		n := make(chan uint)

		// Setup the display
		speccy.CommandChannel <- spectrum.Cmd_GetNumDisplayReceivers{n}
		if <-n == 0 {
			if *fullscreen {
				*scale2x = true
			}

			if *scale2x {
				speccy.CommandChannel <- spectrum.Cmd_AddDisplay{spectrum.NewSDLScreen2x(app, *fullscreen)}
			} else {
				speccy.CommandChannel <- spectrum.Cmd_AddDisplay{spectrum.NewSDLScreen(app)}
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

	// Begin speccy emulation
	go sdlEventLoop(app.NewEventLoop(), speccy, *verboseKeyboard)
	go emulatorLoop(app.NewEventLoop(), speccy)
	speccy.CommandChannel <- spectrum.Cmd_SetFPS{*fps}

	// Start the console goroutine.
	go console.Run(true)

quit:
	<-app.HasTerminated
	sdl.Quit()
}
