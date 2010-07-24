package spectrum

import (
	"testing"
	"sdl"
	"io/ioutil"
	"image"
	"image/png"
	"os"
	"strings"
	"unsafe"
)

func (s *SDLSurface) At(x, y int) image.Color {
	var bpp = int(s.Surface.Format.BytesPerPixel)

	var pixel = uintptr(unsafe.Pointer(s.Surface.Pixels))

	pixel += uintptr(y*int(s.Surface.Pitch) + x*bpp)

	var color = *((*uint32)(unsafe.Pointer(pixel)))

	var r uint8
	var g uint8
	var b uint8
	var a uint8

	sdl.GetRGBA(color, s.Surface.Format, &r, &g, &b, &a)

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

func readOutputImage(filename string) image.Image {
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

func readInputImage(filename string) *Display {
	display := &Display{}
	display.memory, _ = ioutil.ReadFile(filename)
	return display
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
			if colorsAreNotEqual(got.ScreenSurface.At(x, y), expected.At(x, y)) {
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
	borderColor RGBA
	flash       bool
	diffImage   image.Image
}

func (r *RenderTest) renderInputImage() bool {
	renderedScreen := &SDLScreen{nil, SDLSurface{newSurface()}}

	expectedImage := readOutputImage(r.out)
	inputImage := readInputImage(r.in)

	inputImage.borderColor = r.borderColor

	if r.flash {
		inputImage.flashFrame = 0x10
		inputImage.prepare()
	}

	displayData := inputImage.prepare()

	renderedScreen.render(displayData, nil)

	if diffImage := imagesAreNotEqual(renderedScreen, expectedImage); diffImage != nil {
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
	RenderTest{in: "testdata/initial.scr", out: "testdata/initial.png", borderColor: RGBA{192, 192, 192, 255}},
	RenderTest{in: "testdata/flash.scr", out: "testdata/flash_0.png", borderColor: RGBA{192, 192, 192, 255}},
	RenderTest{in: "testdata/flash.scr", out: "testdata/flash_1.png", borderColor: RGBA{192, 192, 192, 255}, flash: true},
}

func TestSDLRenderer(t *testing.T) {

	initSDL()

	for _, r := range RenderTests {
		if notEqual := r.renderInputImage(); notEqual {
			r.reportError(t)
		}
	}

	sdl.Quit()

}

func BenchmarkRender(b *testing.B) {
	renderedScreen := &SDLScreen{nil, SDLSurface{newSurface()}}

	inputImage := readInputImage("testdata/initial.scr")
	inputImage.borderColor = RGBA{192, 192, 192, 255}

	displayData := inputImage.prepare()

	for i := 0; i < b.N; i++ {
		renderedScreen.render(displayData, nil)
	}

}

func BenchmarkRenderWithoutChanges(b *testing.B) {
	var oldDisplayData *DisplayData = nil
	renderedScreen := &SDLScreen{nil, SDLSurface{newSurface()}}

	inputImage := readInputImage("testdata/initial.scr")
	inputImage.borderColor = RGBA{192, 192, 192, 255}

	displayData := inputImage.prepare()

	for i := 0; i < b.N; i++ {
		renderedScreen.render(displayData, oldDisplayData)
		oldDisplayData = displayData
	}

}
