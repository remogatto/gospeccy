package test

import (
	"io/ioutil"
	"⚛sdl"
	"⚛sdl/ttf"
	"prettytest"
	"clingon"
	"spectrum"
	"spectrum/formats"
	"spectrum/interpreter"
	"spectrum/output"
)

var (
	font        *ttf.Font
	speccy      *spectrum.Spectrum48k
	console     *clingon.Console
	cliRenderer *clingon.SDLRenderer
	app         *spectrum.Application
	r           *renderer
)

type SDLSurfaceAccessor interface {
	UpdatedRectsCh() <-chan []sdl.Rect
	GetSurface() *sdl.Surface
}

type renderer struct {
	app              *spectrum.Application
	appSurface       *sdl.Surface
	speccySurface    SDLSurfaceAccessor
	cliSurface_orNil *clingon.SDLRenderer
	width, height    int
	consoleY         int16
}

func newRenderer(app *spectrum.Application, speccySurface SDLSurfaceAccessor, cliSurface_orNil *clingon.SDLRenderer) *renderer {
	width := spectrum.TotalScreenWidth * 2
	height := spectrum.TotalScreenHeight * 2
	r := &renderer{
		app:              app,
		appSurface:       sdl.SetVideoMode(width, height, 32, 0),
		speccySurface:    speccySurface,
		cliSurface_orNil: cliSurface_orNil,
		width:            width,
		height:           height,
	}
	if cliSurface_orNil != nil {
		go r.loopWithCLI(app.NewEventLoop())
	} else {
		go r.loop(app.NewEventLoop())
	}
	return r
}

// Implement 'interpreter.UserInterfaceSettings'
func (r *renderer) ResizeVideo(scale2x, fullscreen bool) {
	// Empty
}

// Implement 'interpreter.UserInterfaceSettings'
func (r *renderer) ShowPaintedRegions(enable bool) {
	// Empty
}

// Implement 'interpreter.UserInterfaceSettings'
func (r *renderer) EnableSound(enable bool) {
	// Empty
}

func (r *renderer) render(speccyRects, cliRects []sdl.Rect) {
	for _, rect := range speccyRects {
		x, y, w, h := rect.X, rect.Y-int16(r.consoleY), rect.W, rect.H
		r.appSurface.Blit(&rect, r.speccySurface.GetSurface(), &rect)
		if r.cliSurface_orNil != nil {
			r.appSurface.Blit(&sdl.Rect{x, rect.Y, 0, 0}, r.cliSurface_orNil.GetSurface(), &sdl.Rect{x, y, w, h})
		}
	}
	for _, rect := range cliRects {
		x, y, w, h := rect.X, rect.Y+int16(r.consoleY), rect.W, rect.H
		r.appSurface.Blit(&sdl.Rect{x, y, 0, 0}, r.speccySurface.GetSurface(), &sdl.Rect{x, y, w, h})
		r.appSurface.Blit(&sdl.Rect{rect.X, rect.Y + int16(r.consoleY), 0, 0}, r.cliSurface_orNil.GetSurface(), &rect)
	}
	r.appSurface.Flip()
}

// Synchronously destroy the CLI renderer
func (r *renderer) destroyCliRenderer() {
	if r.cliSurface_orNil != nil {
		cliSurface := r.cliSurface_orNil
		r.cliSurface_orNil = nil

		console.SetRenderer(nil)

		go func() {
			for r := range cliSurface.UpdatedRectsCh() {
				if r == nil {
					//if app.Verbose {
					//	app.PrintfMsg("command-line renderer: end of the stream of update-rectangles")
					//}
					break
				}
			}
		}()

		done := make(chan bool)
		cliSurface.EventCh() <- clingon.Cmd_Terminate{done}
		<-done

		cliSurface.GetSurface().Free()
		// Note: 'font' is a global variable, it cannot be freed
	}
}

func (r *renderer) loopWithCLI(evtLoop *spectrum.EventLoop) {
	var cliSurface_updatedRectsCh_orNil <-chan []sdl.Rect = r.cliSurface_orNil.UpdatedRectsCh()

	go func() {
		for {
			select {
			case <-evtLoop.Pause:
				r.destroyCliRenderer()
				evtLoop.Pause <- 0

			case <-evtLoop.Terminate:
				evtLoop.Terminate <- 0
				return

			case speccyRects := <-r.speccySurface.UpdatedRectsCh():
				r.render(speccyRects, nil)

			case cliRects := <-cliSurface_updatedRectsCh_orNil:
				r.render(nil, cliRects)
			}
		}
	}()
}

func (r *renderer) loop(evtLoop *spectrum.EventLoop) {
	go func() {
		for {
			select {
			case <-evtLoop.Pause:
				evtLoop.Pause <- 0

			case <-evtLoop.Terminate:
				evtLoop.Terminate <- 0
				return

			case speccyRects := <-r.speccySurface.UpdatedRectsCh():
				r.render(speccyRects, nil)
			}
		}
	}()
}

type testSuite struct {
	prettytest.Suite
}

func (t *testSuite) beforeAll() {
	if sdl.Init(sdl.INIT_VIDEO|sdl.INIT_AUDIO) != 0 {
		app.PrintfMsg("%s", sdl.GetError())
		app.RequestExit()
		<-app.HasTerminated
		sdl.Quit()
	}
	sdl.WM_SetCaption("GoSpeccy - ZX Spectrum Emulator - Test mode", "")
}

func (t *testSuite) afterAll() {
	sdl.Quit()
}

func (t *testSuite) before() {
	StartFullEmulation(false)
}

func (t *testSuite) after() {
	app.RequestExit()
	<-app.HasTerminated
}

type cliTestSuite struct {
	prettytest.Suite
	t *testSuite
}

func (t *cliTestSuite) beforeAll() {
	t.t.beforeAll()
	if ttf.Init() != 0 {
		panic(sdl.GetError())
	}
	sdl.EnableUNICODE(1)
	font = ttf.OpenFont("testdata/VeraMono.ttf", 12)
	if font == nil {
		panic(sdl.GetError())
	}
}

func (t *cliTestSuite) afterAll() {
	font.Close()
	t.t.afterAll()
}

func (t *cliTestSuite) after() {
	t.t.after()
}

func (t *cliTestSuite) before() {
	StartFullEmulation(true)
}

func StartFullEmulation(cli bool) {
	rom, err := spectrum.ReadROM("testdata/48.rom")
	if err != nil {
		panic(err)
	}

	app = spectrum.NewApplication()
	speccy = spectrum.NewSpectrum48k(app, *rom)
	speccy.TapeDrive().NotifyLoadComplete = true
	sdlScreen := output.NewSDLScreen2x(app)
	speccy.CommandChannel <- spectrum.Cmd_AddDisplay{sdlScreen}
	if !cli {
		r = newRenderer(app, sdlScreen, nil)
	} else {
		width := spectrum.TotalScreenWidth * 2
		height := spectrum.TotalScreenHeight * 2
		cliRenderer := clingon.NewSDLRenderer(sdl.CreateRGBSurface(sdl.SRCALPHA, int(width), int(height/2), 32, 0, 0, 0, 0), font)
		cliRenderer.GetSurface().SetAlpha(sdl.SRCALPHA, 0xdd)
		r = newRenderer(app, sdlScreen, cliRenderer)
		r.consoleY = int16(r.height / 2)
		interpreter.IgnoreStartupScript = true
		interpreter.Init(app, "", speccy, r)
		console = clingon.NewConsole(&interpreter.Interpreter{})
		console.SetRenderer(cliRenderer)
		console.SetPrompt("gospeccy> ")

		console.Print(`
Welcome to the GoSpeccy CLI Testing Mode
----------------------------------------
`)
	}
	audio, err := output.NewSDLAudio(app)
	if err == nil {
		speccy.CommandChannel <- spectrum.Cmd_AddAudioReceiver{audio}
	} else {
		app.PrintfMsg("%s", err)
	}
	go speccy.EmulatorLoop()
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
