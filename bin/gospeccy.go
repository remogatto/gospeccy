package main

import (
	"spectrum"
	"sdl"
	"fmt"
	"flag"
	"os"
)

// Initialize system memory, ports and display objects. Memory and
// port objects depends by a display object that, in turn, implements
// the DisplayAccessor interface. System memory, ports and display are
// related because display memory is shared with system memory and
// ports control the display borders. Using Go's interfaces you could
// easily implement new display backend (e.g. exp/draw).
var (
	display = spectrum.NewSDLDisplay(sdl.SetVideoMode(320, 240, 32, 0))
	memory  = &spectrum.Memory{ Display: display }
	port    = &spectrum.Port{ Display: display }

	speccy  *spectrum.Spectrum48k
)

func run() {
	help := flag.Bool("h", false, "Show usage")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "gospeccy [options] [image.sna].\n\n")
		fmt.Fprintf(os.Stderr, "Options are:\n\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	speccy = spectrum.NewSpectrum48k(memory, port)
	
	if *help == true {
		flag.Usage()
		return
	}

	if flag.Arg(0) != "" {
		speccy.LoadSna(flag.Arg(0))
	}

	if sdl.Init(sdl.INIT_EVERYTHING) != 0 {
		panic(sdl.GetError())
	}

	if display == nil {
		panic(sdl.GetError())
	}

	sdl.WM_SetCaption("GoSpeccy", "")

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
		sdl.Delay(20)
	}

	sdl.Quit()

}

func main() {
	run()
}
