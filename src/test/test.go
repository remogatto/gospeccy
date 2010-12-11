package test

import (
	"time"
	"os"
	"io/ioutil"
	"⚛sdl"
	"spectrum"
	"spectrum/formats"
	"prettytest"
 )

 var (
	 speccy *spectrum.Spectrum48k
	 app    *spectrum.Application
 )

type testSuite struct { prettytest.Suite }

func (t *testSuite) beforeAll() {
	StartFullEmulation()
}

func (t *testSuite) afterAll() {
	app.RequestExit()
	<-app.HasTerminated
}

func (t *testSuite) before() {
	StartFullEmulation()
}

func (t *testSuite) after() {
	app.RequestExit()
	<-app.HasTerminated
}

func emulatorLoop(evtLoop *spectrum.EventLoop, speccy *spectrum.Spectrum48k) {
	app = evtLoop.App()

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

			close(speccy.ROMLoaded())

			evtLoop.Pause <- 0

		case <-evtLoop.Terminate:
			// Terminate this Go routine
			if app.Verbose {
				app.PrintfMsg("emulator loop: exit")
			}
			evtLoop.Terminate <- 0
			return

		case <-ticker.C:
			speccy.CommandChannel <- spectrum.Cmd_RenderFrame{}
			speccy.CommandChannel <- spectrum.Cmd_CheckSystemROMLoaded{}

		case FPS_new := <-speccy.FPS:
			if FPS_new != fps && FPS_new > 0 {
				if app.Verbose {
					app.PrintfMsg("setting FPS to %f", FPS_new)
				}
				ticker.Stop()
				spectrum.Drain(ticker)
				// ticker = time.NewTicker(int64(1e9 / FPS_new))
				ticker = time.NewTicker(1)
				fps = FPS_new
			}

		}
	}
}

func StartFullEmulation() {
	var err os.Error

	app = spectrum.NewApplication()

	speccy, err = spectrum.NewSpectrum48k(app, "testdata/48.rom")
	speccy.TapeDrive().NotifyLoadComplete = true

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

	speccy.CommandChannel <- spectrum.Cmd_AddDisplay{spectrum.NewSDLScreen(app)}

	audio, err := spectrum.NewSDLAudio(app)

	if err == nil {
		speccy.CommandChannel <- spectrum.Cmd_AddAudioReceiver{audio}
	} else {
		app.PrintfMsg("%s", err)
	}

	go speccy.EmulatorLoop()

	<-speccy.ROMLoaded()
}

func loadSnapshot(filename string) formats.Snapshot {
	_snapshot, err := formats.ReadProgram(filename)

	if err != nil {
		panic(err)
	}

	var snapshot formats.Snapshot
	var ok bool
	if snapshot, ok = _snapshot.(formats.Snapshot); !ok {
		panic("invalid type")
	}

	return snapshot
}

func assertScreenEqual(expected, actual formats.Snapshot) bool {
	// Doesn't compare screen attributes
	for address, actualValue := range actual.Memory()[:0x1800] {
		if expected.Memory()[address] != actualValue {
			return false
		}
	}
	return true
}

func assertStateEqual(expected, actual formats.Snapshot) bool {
	// Compare memory ignoring PC value
	for address, actualValue := range actual.Memory() {
		if expected.Memory()[address] != actualValue && address != 0xbf4a && address != 0xbf4b {
			return false
		}
	}

	// FIXME: Should compare also CPU and ULA states

	return true
}

func stateEqualTo(filename string) bool {
	return assertStateEqual(loadSnapshot(filename), speccy.Cpu.MakeSnapshot())
}

func screenEqualTo(filename string) bool {
	return assertScreenEqual(loadSnapshot(filename), speccy.Cpu.MakeSnapshot())
}

func saveSnapshot(filename string) {
	fullSnapshot := speccy.Cpu.MakeSnapshot()

	data, err := fullSnapshot.EncodeSNA()

	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(filename, data, 0600)

	if err != nil {
		panic(err)
	}
}
