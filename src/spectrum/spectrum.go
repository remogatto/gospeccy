package spectrum

import (
	"fmt"
	"io/ioutil"
	"os"
	"container/vector"
)

const Spectrum48k_ROM_filepath = "roms/48.rom"
const TStatesPerFrame = 69888 // Number of T-states per frame
const InterruptLength = 32    // How long does an interrupt last in T-states

type Spectrum48k struct {
	app      *Application
	Cpu      *Z80
	Memory   MemoryAccessor
	Displays vector.Vector // A vector of DisplayChannel, initially empty
	Keyboard *Keyboard
	Ports    *Ports
}

// Create a new speccy object.
func NewSpectrum48k(app *Application) (*Spectrum48k, os.Error) {
	memory := NewMemory()
	keyboard := NewKeyboard()
	ports := NewPorts(memory, keyboard)
	z80 := NewZ80(memory, ports)

	ports.z80 = z80
	memory.z80 = z80

	// Load the first 16k of memory with the ROM image
	{
		rom48k, err := ioutil.ReadFile(Spectrum48k_ROM_filepath)
		if err != nil {
			return nil, err
		}
		if len(rom48k) != 0x4000 {
			return nil, os.NewError(fmt.Sprintf("ROM file \"%s\" has an invalid size", Spectrum48k_ROM_filepath))
		}

		for address, b := range rom48k {
			memory.set(uint16(address), b)
		}
	}

	speccy := &Spectrum48k{app: app, Cpu: z80, Memory: memory, Keyboard: keyboard, Displays: vector.Vector{}, Ports: ports}

	return speccy, nil
}

func (speccy *Spectrum48k) Close() {
	speccy.Cpu.Close()

	if speccy.app.Verbose {
		eff := speccy.Cpu.GetEmulationEfficiency()
		if eff != 0 {
			fmt.Printf("emulation efficiency: %d host-CPU instructions per Z80 instruction\n", eff)
		} else {
			fmt.Printf("emulation efficiency: -\n")
		}
	}
}

func (speccy *Spectrum48k) AddDisplay(display DisplayChannel) {
	speccy.Displays.Push(display)
}

// Execute the number of T-states corresponding to one screen frame
func (speccy *Spectrum48k) doOpcodes() {
	eventNextEvent = TStatesPerFrame
	speccy.Cpu.tstates = (speccy.Cpu.tstates % TStatesPerFrame)
	speccy.Cpu.doOpcodes()
}

func (speccy *Spectrum48k) interrupt() {
	speccy.Cpu.interrupt()
}

func (speccy *Spectrum48k) RenderFrame() {
	speccy.Ports.frame_begin(speccy.Memory.getBorder())
	speccy.Memory.frame_begin()
	speccy.interrupt()
	speccy.doOpcodes()
	for _, display := range speccy.Displays {
		speccy.Memory.sendScreenToDisplay(display.(DisplayChannel), speccy.Ports.borderEvents)
	}
	speccy.Ports.frame_releaseMemory()
}

// Initialize state from the snapshot defined by the specified filename.
// Returns nil on success.
func (speccy *Spectrum48k) LoadSna(filename string) os.Error {
	if speccy.app.Verbose {
		fmt.Printf("loading snapshot \"%s\"\n", filename)
	}

	return speccy.Cpu.LoadSna(filename)
}

// func dumpRegisters() {
// 	fmt.Printf("%02x%02x %02x%02x %02x%02x %02x%02x %02x%02x %02x%02x %02x%02x %02x%02x %02x%02x %02x%02x %04x %04x\n",
// 		z80.a, z80.f, z80.b, z80.c, z80.d, z80.e, z80.h, z80.l, z80.a_, z80.f_, z80.b_, z80.c_, z80.d_, z80.e_, z80.h_, z80.l_, z80.ixh, z80.ixl, z80.iyh, z80.iyl, z80.sp, z80.pc)
// 	fmt.Printf("%02x %02x %d %d %d %d %d\n", z80.i, (z80.r7&0x80)|(z80.r&0x7f),
// 		z80.iff1, z80.iff2, z80.im, z80.halted, tstates)
// }

// func dumpMemory() {
// 	for i, val := range memory {
// 		fmt.Printf("%d %d\n", i, val)
// 	}
// }
