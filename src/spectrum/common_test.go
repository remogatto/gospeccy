package spectrum

import (
	"time"
	"os"
	"⚛sdl"
)

var (
	speccy *Spectrum48k
	app *Application
	romLoaded bool
	romLoadedCh chan bool = make(chan bool)
)

func emulatorLoop(evtLoop *EventLoop, speccy *Spectrum48k) {
	romLoaded = false
	app := evtLoop.App()

	fps := <-speccy.FPS
	ticker := time.NewTicker(int64(1e9 / fps))

	// Render the 1st frame (the 2nd frame will be rendered after 1/FPS seconds)
	{
		completionTime := make(chan int64)
		speccy.CommandChannel <- Cmd_RenderFrame{completionTime}

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
			Drain(ticker)
			evtLoop.Pause <- 0

		case <-evtLoop.Terminate:
			// Terminate this Go routine
			if app.Verbose {
				app.PrintfMsg("emulator loop: exit")
			}
			evtLoop.Terminate <- 0
			return

		case <-ticker.C:
			speccy.CommandChannel <- Cmd_RenderFrame{}
			
			if speccy.Cpu.PC() == 0x10ac && !romLoaded {
				romLoadedCh <- true
				romLoaded = true
			}

		case FPS_new := <-speccy.FPS:
			if (FPS_new != fps) && (FPS_new > 0) {
				if app.Verbose {
					app.PrintfMsg("setting FPS to %f", FPS_new)
				}
				ticker.Stop()
				Drain(ticker)
				ticker = time.NewTicker(int64(1e9 / FPS_new))
				fps = FPS_new
			}
		}
	}
}

func StartFullEmulation() {
	var err os.Error

	app = NewApplication()
	speccy, err = NewSpectrum48k(app, "testdata/48.rom")

	if err != nil {
		panic(err)
	}

	if sdl.Init(sdl.INIT_VIDEO|sdl.INIT_AUDIO) != 0 {
		app.PrintfMsg("%s", sdl.GetError())
		app.RequestExit()
		<-app.HasTerminated
		sdl.Quit()
	}

	sdl.WM_SetCaption("GoSpeccy - ZX Spectrum Emulator - Test mode", "")

	speccy.CommandChannel <- Cmd_AddDisplay{ NewSDLScreen(app) }

	audio, err := NewSDLAudio(app)

	if err == nil {
		speccy.CommandChannel <- Cmd_AddAudioReceiver{ audio }
	} else {
		app.PrintfMsg("%s", err)
	}

	go emulatorLoop(app.NewEventLoop(), speccy)

	<-romLoadedCh
}
