package spectrum

import (
	"fmt"
	"io/ioutil"
	"os"
)

const Spectrum48k_ROM_filepath = "roms/48.rom"
const TStatesPerFrame = 69888

type Spectrum48k struct {
	Cpu      *Z80
	Memory   MemoryAccessor
	Display  DisplayChannel // Initially nil
	Keyboard *Keyboard
	Ports    *Ports
}

// Create a new speccy object.
func NewSpectrum48k() (*Spectrum48k, os.Error) {
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

	speccy := &Spectrum48k{Cpu: z80, Memory: memory, Keyboard: keyboard, Display: nil, Ports: ports}

	return speccy, nil
}

func (speccy *Spectrum48k) SetDisplay(display DisplayChannel) {
	speccy.Display = display
}

// Execute the number of T-states corresponding to one screen frame
func (speccy *Spectrum48k) doOpcodes() {
	eventNextEvent = TStatesPerFrame
	speccy.Cpu.tstates = 0
	speccy.Cpu.doOpcodes()
}

func (speccy *Spectrum48k) interrupt() {
	speccy.Cpu.interrupt()
}

func (speccy *Spectrum48k) RenderFrame() {
	speccy.Ports.frame_begin(speccy.Memory.getBorder())
	speccy.doOpcodes()
	if speccy.Display != nil {
		speccy.Memory.sendScreenToDisplay(speccy.Display, speccy.Ports.borderEvents)
	}
	speccy.Ports.frame_releaseMemory()
	speccy.interrupt()
}

// Initialize state from the snapshot defined by the specified filename.
// Returns nil on success.
func (speccy *Spectrum48k) LoadSna(filename string) os.Error {
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
