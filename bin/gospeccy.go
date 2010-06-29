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
	"fmt"
	"flag"
	"os"
)

var (
	display *spectrum.SDLDisplay
	doubledDisplay *spectrum.SDLDoubledDisplay

	applicationScreen *spectrum.SDLSurface
	port *spectrum.Port
	memory *spectrum.Memory
	speccy *spectrum.Spectrum48k

	sdlMode uint32
	deltat int64
)

// Big game loop block. Need a bit of refactoring I guess :)
func run() {
	help := flag.Bool("h", false, "Show usage")
	scale := flag.Bool("d", false, "Double size")
	fullscreen:= flag.Bool("f", false, "Fullscreen (enable double size by default)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "GOSpeccy - A simple ZX Spectrum 48k Emulator written in GO\n\n")
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
		doubledDisplay = spectrum.NewSDLDoubledDisplay(sdl.SetVideoMode(640, 480, 32, 0))
		memory  = &spectrum.Memory{ Display: doubledDisplay }
		port    = &spectrum.Port{ Display: doubledDisplay }
		speccy = spectrum.NewSpectrum48k(memory, port)
		applicationScreen = doubledDisplay.ScreenSurface 
	} else {
		display = spectrum.NewSDLDisplay(sdl.SetVideoMode(320, 240, 32, 0))
		memory  = &spectrum.Memory{ Display: display }
		port    = &spectrum.Port{ Display: display }
		speccy = spectrum.NewSpectrum48k(memory, port)
		applicationScreen = display.ScreenSurface

	}

	if applicationScreen.Surface == nil {
		panic(sdl.GetError())
	}

	if flag.Arg(0) != "" {
		fmt.Println(flag.Arg(0))
		speccy.LoadSna(flag.Arg(0))
	}

	sdl.WM_SetCaption("GoSpeccy - ZX Spectrum Emulator", "")

	running := true

	for running {

		e := &sdl.Event{}

		for e.Poll() {
			switch e.Type {
			case sdl.QUIT:
				running = false
				break
			case sdl.KEYDOWN, sdl.KEYUP:
				println("")
				println(e.Keyboard().Keysym.Sym, ": ", sdl.GetKeyName(sdl.Key(e.Keyboard().Keysym.Sym)))

				fmt.Printf("%04x ", e.Type)

				for i := 0; i < len(e.Pad0); i++ {
					fmt.Printf("%02x ", e.Pad0[i])
				}
				println()

				k := e.Keyboard()

				fmt.Printf("Type: %02x Which: %02x State: %02x Pad: %02x\n", k.Type, k.Which, k.State, k.Pad0[0])
				fmt.Printf("Scancode: %02x Sym: %08x Mod: %04x Unicode: %04x\n", k.Keysym.Scancode, k.Keysym.Sym, k.Keysym.Mod, k.Keysym.Unicode)

				switch k.Keysym.Sym {

				case 27:
					running = false
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

		applicationScreen.Surface.Flip()

		sdl.Delay(20)
	}

	sdl.Quit()

}

func main() {
	run()
}
