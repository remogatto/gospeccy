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
	"sdl"
	"sdl/ttf"
	"fmt"
	"flag"
	"os"
	"time"
	"io/ioutil"
)

const frameDuration = 1e9 / 50 // 50 frames a second

var (
	applicationSurface *sdl.Surface
	port *spectrum.Port
	memory *spectrum.Memory
	speccy *spectrum.Spectrum48k

	sdlMode uint32
	lastFrame, delay int64
	systemROM []byte
	showFPS bool
)

func gospeccyDir() string {
	return os.Getenv("HOME") + "/.gospeccy/"
}

func fontFilename(filename string) string {
	return gospeccyDir() + "font/" + filename
}

func romFilename(filename string) string {
	return gospeccyDir() + "rom/" + filename
}

func snaFilename(filename string) string {
	return gospeccyDir() + "sna/" + filename
}

func loadSystemROM(filename string) []byte {
	rom, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return rom
}

func printKeyInfo(e *sdl.KeyboardEvent) {
	println("")
	println(e.Keysym.Sym, ": ", sdl.GetKeyName(sdl.Key(e.Keysym.Sym)))

	fmt.Printf("Type: %02x Which: %02x State: %02x Pad: %02x\n", e.Type, e.Which, e.State, e.Pad0[0])
	fmt.Printf("Scancode: %02x Sym: %08x Mod: %04x Unicode: %04x\n", e.Keysym.Scancode, e.Keysym.Sym, e.Keysym.Mod, e.Keysym.Unicode)
}

func printFPS(appSurface *sdl.Surface, font *ttf.Font, delay int64) {
	white := sdl.Color{255, 255, 255, 0}
	text := ttf.RenderText_Blended(font, fmt.Sprintf("FPS %d", int(1/(float(delay)/1e9))), white)
	appSurface.Blit(&sdl.Rect{10, 10, 0, 0}, text, nil)
}

// Big loop block. Need a bit of refactoring I guess :)
func run() {

	help := flag.Bool("help", false, "Show usage")
	verbose := flag.Bool("verbose", false, "Print a lot of stupid messages")
	scale := flag.Bool("doubled", false, "Double size display")
	fullscreen:= flag.Bool("fullscreen", false, "Fullscreen (enable double size by default)")
	rom := flag.String("rom", romFilename("48.rom"), "Start with the given system rom")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "GoSpeccy - A simple ZX Spectrum 48k Emulator written in GO\n\n")
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

	// Load system ROM
	if *rom != "" {
		systemROM = loadSystemROM(*rom)
	} else {
		systemROM = loadSystemROM(romFilename("48.rom"))
	}

	if sdl.Init(sdl.INIT_EVERYTHING) != 0 {
		panic(sdl.GetError())
	}

	// Initialize system memory, ports and display objects. Memory
	// and port objects depends by a display object that, in turn,
	// implements the DisplayAccessor interface. System memory,
	// ports and display are related because, on the speccy,
	// display memory is shared with system memory and ports
	// control the display borders. Using Go's interfaces you
	// could easily implement new display backend (e.g. exp/draw).
	
	if *fullscreen {
		sdlMode = sdl.FULLSCREEN
		*scale = true
	} else {
		sdlMode = 0
	}

	if *scale {
		display := spectrum.NewSDLDoubledScreen(sdl.SetVideoMode(640, 480, 32, sdlMode))
		speccy = spectrum.NewSpectrum48k(systemROM, display)
		applicationSurface = display.ScreenSurface.Surface
	} else {
		display := spectrum.NewSDLScreen(sdl.SetVideoMode(320, 240, 32, sdlMode))
		speccy = spectrum.NewSpectrum48k(systemROM, display)
		applicationSurface = display.ScreenSurface.Surface
	}

	if applicationSurface == nil {
		panic(sdl.GetError())
	}

	// Try to load the snapshot file from the current dir. If not
	// found fall back to $HOME/.gospeccy/sna/filename
	if flag.Arg(0) != "" {
		if speccy.LoadSna(flag.Arg(0)) != nil {
			err := speccy.LoadSna(snaFilename(flag.Arg(0)))
			if err != nil {
				panic(err)
			}
		}
	}

	sdl.WM_SetCaption("GoSpeccy - ZX Spectrum Emulator", "")

	// Initialize font library and load fonts

	if ttf.Init() != 0 {
		panic(sdl.GetError())
	}

	font := ttf.OpenFont(fontFilename("Fontin Sans.otf"), 14)

	if font == nil {
		panic(sdl.GetError())
	}

	running := true

	for running {

		e := &sdl.Event{}

		for e.Poll() {
			switch e.Type {
			case sdl.QUIT:
				running = false
				break
			case sdl.KEYDOWN, sdl.KEYUP:

				k := e.Keyboard()

				if *verbose {
					printKeyInfo(k)
				}

				switch k.Keysym.Sym {

				case 27: // ESC
					running = false

				case 291: // f10
					if k.State != 0 {
						showFPS = !showFPS
					}

				default:
					if k.State != 0 {
						speccy.KeyDown(uint(k.Keysym.Sym))
					} else {
						speccy.KeyUp(uint(k.Keysym.Sym))
					}
							

				}
			}
		}

		speccy.RenderFrame()

		// FIXME: This auto-adjust-delay works well. BTW, it
		// could be better to take an average delay after N
		// rendered frames.

		delay = time.Nanoseconds() - lastFrame

		if delay <= frameDuration {
			time.Sleep(frameDuration - delay)
		}

		lastFrame = time.Nanoseconds()

		if showFPS {
			printFPS(applicationSurface, font, delay)
		}

		applicationSurface.Flip()

	}

	sdl.Quit()

}

func main() { run() }
