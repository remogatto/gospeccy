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
const DefaultFPS = 50.08

type DisplayInfo struct {
	displayReceiver DisplayReceiver

	// The index of the last frame sent to the 'displayReceiver', initially nil.
	lastFrame *uint
}

type Spectrum48k struct {
	app    *Application
	Cpu    *Z80
	Memory MemoryAccessor
	ula    *ULA

	displays vector.Vector // A vector of '*DisplayInfo', initially empty

	Keyboard *Keyboard
	Ports    PortAccessor

	// Send a single value to this channel in order to change the display refresh frequency.
	// By default, this channel initially receives the value 'DefaultFPS'.
	FPS chan float

	CommandChannel chan interface{}
}

// Create a new speccy object.
func NewSpectrum48k(app *Application) (*Spectrum48k, os.Error) {
	memory := NewMemory()
	keyboard := NewKeyboard()
	ports := NewPorts()
	z80 := NewZ80(memory, ports)
	ula := NewULA()

	speccy := &Spectrum48k{app: app, Cpu: z80, Memory: memory, ula: ula, Keyboard: keyboard, displays: vector.Vector{}, Ports: ports}

	memory.init(speccy)
	ula.init(speccy)
	ports.init(speccy)

	err := speccy.reset()
	if err != nil {
		return nil, err
	}

	speccy.FPS = make(chan float, 1)
	speccy.FPS <- DefaultFPS

	speccy.CommandChannel = make(chan interface{})
	go commandLoop(speccy)

	return speccy, nil
}

type Cmd_Reset struct{}
type Cmd_LoadSna struct {
	Filename string
	ErrChan  chan os.Error
}
type Cmd_RenderFrame struct{}
type Cmd_AddDisplay struct {
	Display DisplayReceiver
}
type Cmd_CloseAllDisplays struct{}

func commandLoop(speccy *Spectrum48k) {
	evtLoop := speccy.app.NewEventLoop()
	for {
		select {
		case <-evtLoop.Pause:
			speccy.Close()
			evtLoop.Pause <- 0

		case <-evtLoop.Terminate:
			// Terminate this Go routine
			if evtLoop.App().Verbose {
				println("command loop: exit")
			}
			evtLoop.Terminate <- 0
			return

		case untyped_cmd := <-speccy.CommandChannel:
			switch cmd := untyped_cmd.(type) {
			case Cmd_Reset:
				speccy.reset()

			case Cmd_LoadSna:
				err := speccy.loadSna(cmd.Filename)
				if cmd.ErrChan != nil {
					cmd.ErrChan <- err
				}

			case Cmd_RenderFrame:
				speccy.renderFrame()

			case Cmd_AddDisplay:
				speccy.addDisplay(cmd.Display)

			case Cmd_CloseAllDisplays:
				speccy.closeAllDisplays()
			}
		}
	}
}

func (speccy *Spectrum48k) reset() os.Error {
	speccy.Cpu.reset()
	speccy.Memory.reset()
	speccy.ula.reset()
	speccy.Keyboard.reset()
	speccy.Ports.reset()

	// Load the first 16k of memory with the ROM image
	{
		rom48k, err := ioutil.ReadFile(Spectrum48k_ROM_filepath)
		if err != nil {
			return err
		}
		if len(rom48k) != 0x4000 {
			return os.NewError(fmt.Sprintf("ROM file \"%s\" has an invalid size", Spectrum48k_ROM_filepath))
		}

		for address, b := range rom48k {
			speccy.Memory.Write(uint16(address), b)
		}
	}

	return nil
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

func (speccy *Spectrum48k) addDisplay(display DisplayReceiver) {
	speccy.displays.Push(&DisplayInfo{display, nil})
}

func (speccy *Spectrum48k) closeAllDisplays() {
	var displays vector.Vector
	{
		displays = speccy.displays
		speccy.displays = vector.Vector{}
	}

	for _, display := range displays {
		display.(*DisplayInfo).displayReceiver.close()
	}
}

// Execute the number of T-states corresponding to one screen frame
func (speccy *Spectrum48k) doOpcodes() {
	speccy.Cpu.eventNextEvent = TStatesPerFrame
	speccy.Cpu.tstates = (speccy.Cpu.tstates % TStatesPerFrame)
	speccy.Cpu.doOpcodes()
}

func (speccy *Spectrum48k) renderFrame() {
	speccy.Ports.frame_begin()
	speccy.ula.frame_begin()

	speccy.Cpu.interrupt()
	speccy.doOpcodes()

	for _, display := range speccy.displays {
		speccy.ula.sendScreenToDisplay(display.(*DisplayInfo))
	}

	speccy.Ports.frame_releaseMemory()
}

// Initialize state from the snapshot defined by the specified filename.
// Returns nil on success.
func (speccy *Spectrum48k) loadSna(filename string) os.Error {
	if speccy.app.Verbose {
		fmt.Printf("loading snapshot \"%s\"\n", filename)
	}

	err := speccy.reset()
	if err != nil {
		return err
	}

	err = speccy.Cpu.LoadSna(filename)
	if err != nil {
		return err
	}

	return nil
}
