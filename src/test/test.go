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
	r renderer
 )

type renderer struct {
	appSurface *sdl.Surface
	speccySurface *spectrum.SDLScreen2x
	width, height int
}

func (r *renderer) render(speccyRects []sdl.Rect) {
	for _, rect := range speccyRects {
		r.appSurface.Blit(&rect, r.speccySurface.GetSurface(), &rect)
		r.appSurface.UpdateRect(int32(rect.X), int32(rect.Y), uint32(rect.W), uint32(rect.H))
	}
}

type testSuite struct { prettytest.Suite }

func (t *testSuite) beforeAll() {
	if sdl.Init(sdl.INIT_VIDEO|sdl.INIT_AUDIO) != 0 {
		app.PrintfMsg("%s", sdl.GetError())
		app.RequestExit()
		<-app.HasTerminated
		sdl.Quit()
	}

	sdl.WM_SetCaption("GoSpeccy - ZX Spectrum Emulator - Test mode", "")

	r.width = spectrum.TotalScreenWidth*2
	r.height = spectrum.TotalScreenHeight*2
	r.appSurface = sdl.SetVideoMode(r.width, r.height, 32, 0)
}

func (t *testSuite) afterAll() {
	sdl.Quit()
}

func (t *testSuite) before() {
	StartFullEmulation()
}

func (t *testSuite) after() {
	app.RequestExit()
	<-app.HasTerminated
}

func StartFullEmulation() {
	var (
		err os.Error
		speccyRects []sdl.Rect
	)

	app = spectrum.NewApplication()

	speccy, err = spectrum.NewSpectrum48k(app, "testdata/48.rom")
	speccy.TapeDrive().NotifyLoadComplete = true

	if err != nil {
		panic(err)
	}

	sdlScreen := spectrum.NewSDLScreen2x(app)
	speccy.CommandChannel <- spectrum.Cmd_AddDisplay{sdlScreen}
	r.speccySurface = sdlScreen

	audio, err := spectrum.NewSDLAudio(app)

	if err == nil {
		speccy.CommandChannel <- spectrum.Cmd_AddAudioReceiver{audio}
	} else {
		app.PrintfMsg("%s", err)
	}

	go speccy.EmulatorLoop()

	ticker := time.NewTicker(1e9/int64(60))
	evtLoop := app.NewEventLoop()

	go func() {
		for {
			select {
			case <-evtLoop.Pause:
				evtLoop.Pause <- 0

			case <-evtLoop.Terminate:
				close(r.speccySurface.UpdatedRectsCh())
				ticker.Stop()
				evtLoop.Terminate <- 0
				return

			case speccyRects = <-r.speccySurface.UpdatedRectsCh():

			case <-ticker.C: r.render(speccyRects)
			}
		}
	}()

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
