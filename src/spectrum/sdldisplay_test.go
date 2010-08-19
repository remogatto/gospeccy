package spectrum

import (
	"testing"
	"⚛sdl"
	"io/ioutil"
	"image"
	"image/png"
	"os"
	"strings"
	"unsafe"
)

func (s *SDLSurface) At(x, y int) image.Color {
	var bpp = int(s.surface.Format.BytesPerPixel)

	var pixel = uintptr(unsafe.Pointer(s.surface.Pixels))

	pixel += uintptr(y*int(s.surface.Pitch) + x*bpp)

	var color = *((*uint32)(unsafe.Pointer(pixel)))

	var r uint8
	var g uint8
	var b uint8
	var a uint8

	sdl.GetRGBA(color, s.surface.Format, &r, &g, &b, &a)

	return image.RGBAColor{uint8(r), uint8(g), uint8(b), uint8(a)}
}

func initSDL() {
	if sdl.Init(sdl.INIT_VIDEO) != 0 {
		panic(sdl.GetError())
	}
}

func newSurface() *sdl.Surface {
	return sdl.SetVideoMode(TotalScreenWidth, TotalScreenHeight, 32, 0)
}

func loadExpectedImage(filename string) image.Image {
	var file *os.File
	var err os.Error
	var image image.Image

	if file, err = os.Open(filename, os.O_RDONLY, 0); err != nil {
		panic(err)
	}

	if image, err = png.Decode(file); err != nil {
		panic(err)
	}
	return image
}

func loadScreen(filename string, speccy *Spectrum48k) *ULA {
	memory := NewMemory()
	memory.init(speccy)
	ula := NewULA()
	ula.init(speccy)
	image, _ := ioutil.ReadFile(filename)
	for offset, value := range image {
		speccy.Memory.Write(0x4000+uint16(offset), value)
	}
	return ula
}

func colorsAreNotEqual(got, expected image.Color) bool {
	got_r, got_g, got_b, got_a := got.RGBA()
	expected_r, expected_g, expected_b, expected_a := expected.RGBA()
	if (got_r != expected_r) || (got_g != expected_g) || (got_b != expected_b) || (got_a != expected_a) {
		return true
	}
	return false
}

func imagesAreNotEqual(got *SDLScreen, expected image.Image) image.Image {
	diff := false
	diffImage := image.NewRGBA(TotalScreenWidth, TotalScreenHeight)

	for y := 0; y < TotalScreenHeight; y++ {
		for x := 0; x < TotalScreenWidth; x++ {
			if colorsAreNotEqual(got.screenSurface.At(x, y), expected.At(x, y)) {
				diff = true
				diffImage.Set(x, y, image.Red)
			}
		}
	}

	if diff {
		return diffImage
	}

	return nil

}

type RenderTest struct {
	in, out     string
	borderColor byte
	flash       bool
	diffImage   image.Image
}

func (r *RenderTest) renderScreen(speccy *Spectrum48k) bool {
	renderedSDLScreen := &SDLScreen{nil, SDLSurface{newSurface()}, newUnscaledDisplay()}

	speccy.addDisplay(renderedSDLScreen)

	expectedImage := loadExpectedImage(r.out)
	inputScreen := loadScreen(r.in, speccy)

	inputScreen.borderColor = r.borderColor

	if r.flash {
		inputScreen.frame = 0x10
		inputScreen.prepare(speccy.displays.At(0).(*DisplayInfo))
	}

	displayData := inputScreen.prepare(speccy.displays.At(0).(*DisplayInfo))

	renderedSDLScreen.render(displayData, nil)

	if diffImage := imagesAreNotEqual(renderedSDLScreen, expectedImage); diffImage != nil {
		r.diffImage = diffImage
		return true
	}

	return false
}

func (r *RenderTest) getDiffFn() string {
	return strings.TrimRight(r.out, ".png") + "_diff.png"
}

func (r *RenderTest) reportError(t *testing.T) {
	t.Errorf("Expected image %s is not equal to the rendered one! Check %s\n", r.out, r.getDiffFn())

	if file, err := os.Open(r.getDiffFn(), os.O_CREATE|os.O_WRONLY, 0666); err != nil {
		panic(err)
	} else {
		if err := png.Encode(file, r.diffImage); err != nil {
			panic(err)
		}
	}
}

var RenderTests = []RenderTest{
	RenderTest{in: "testdata/initial.scr", out: "testdata/initial.png", borderColor: 7},
	RenderTest{in: "testdata/flash.scr", out: "testdata/flash_0.png", borderColor: 7},
	RenderTest{in: "testdata/flash.scr", out: "testdata/flash_1.png", borderColor: 7, flash: true},
}

func TestSDLRenderer(t *testing.T) {

	initSDL()

	romPath := "testdata/48.rom"

	app := NewApplication()

	if speccy, err := NewSpectrum48k(app, romPath); err != nil {
		panic(err)
	} else {
		for _, r := range RenderTests {
			if notEqual := r.renderScreen(speccy); notEqual {
				r.reportError(t)
			}
		}
	}

	sdl.Quit()

}

func BenchmarkRender(b *testing.B) {

	b.StopTimer()

	initSDL()

	romPath = "testdata/48.rom"

	app := NewApplication()
	surface := newSurface()

	const numFrames = 1000

	var (
		frames    [numFrames]DisplayData
		prevFrame *DisplayData = nil
	)

	sdlScreen := &SDLScreen{make(chan *DisplayData), SDLSurface{surface}, newUnscaledDisplay()}

	if speccy, err := NewSpectrum48k(app, romPath); err != nil {
		panic(err)
	} else {
		speccy.addDisplay(sdlScreen)
		speccy.loadSna("testdata/fire.sna")

		go func() {
			for i := 0; i < numFrames; i++ {
				speccy.renderFrame()
			}
		}()

		for i := 0; i < numFrames; i++ {
			frames[i] = *<-sdlScreen.getDisplayDataChannel()
		}

		var j int

		b.StartTimer()
		for i := 0; i < b.N; i++ {
			j %= numFrames
			sdlScreen.render(&frames[j], prevFrame)
			prevFrame = &frames[j]
			j++
		}

	}

	sdl.Quit()
}
