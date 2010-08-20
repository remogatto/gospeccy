package spectrum

import (
	"fmt"
	"io/ioutil"
	"os"
	"container/vector"
)

const TStatesPerFrame = 69888 // Number of T-states per frame
const InterruptLength = 32    // How long does an interrupt last in T-states
const DefaultFPS = 50.08
const DefaultRomPath = "roms/48.rom"

type DisplayInfo struct {
	displayReceiver DisplayReceiver

	// The index of the last frame sent to the 'displayReceiver', initially nil.
	lastFrame *uint
}

type Spectrum48k struct {
	Cpu      *Z80
	Memory   MemoryAccessor
	ula      *ULA
	Keyboard *Keyboard
	Ports    PortAccessor

	romPath string

	// Send a single value to this channel in order to change the
	// display refresh frequency.  By default, this channel
	// initially receives the value 'DefaultFPS'.
	FPS chan float

	CommandChannel chan interface{}

	// A vector of '*DisplayInfo', initially empty
	displays vector.Vector

	app *Application
}

// Create a new speccy object.
func NewSpectrum48k(app *Application, romPath string) (*Spectrum48k, os.Error) {
	memory := NewMemory()
	keyboard := NewKeyboard()
	ports := NewPorts()
	z80 := NewZ80(memory, ports)
	ula := NewULA()

	speccy := &Spectrum48k{
		Cpu:      z80,
		Memory:   memory,
		ula:      ula,
		Keyboard: keyboard,
		Ports:    ports,
		romPath:  romPath,
		displays: vector.Vector{},
		app:      app,
	}

	memory.init(speccy)
	z80.init(speccy)
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
type Cmd_RenderFrame struct{}
type Cmd_AddDisplay struct {
	Display DisplayReceiver
}
type Cmd_CloseAllDisplays struct{}
type Cmd_SetUlaEmulationAccuracy struct {
	accurateEmulation bool
}
type Cmd_LoadSna struct {
	InformalFilename string // This is only used for logging purposes
	Data             []byte // The SNA snapshot data
	ErrChan          chan os.Error
}
type Snapshot struct {
	data []byte // Constraint: (data == nil) != (err == nil)
	err  os.Error
}
type Cmd_SaveSna struct {
	Chan chan Snapshot
}

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

			case Cmd_RenderFrame:
				speccy.renderFrame()

			case Cmd_AddDisplay:
				speccy.addDisplay(cmd.Display)

			case Cmd_CloseAllDisplays:
				speccy.closeAllDisplays()

			case Cmd_SetUlaEmulationAccuracy:
				speccy.ula.setEmulationAccuracy(cmd.accurateEmulation)

			case Cmd_LoadSna:
				if speccy.app.Verbose {
					if len(cmd.InformalFilename) > 0 {
						fmt.Printf("loading SNA snapshot \"%s\"\n", cmd.InformalFilename)
					} else {
						fmt.Printf("loading a SNA snapshot\n")
					}
				}

				err := speccy.loadSna(cmd.Data)
				if cmd.ErrChan != nil {
					cmd.ErrChan <- err
				}

			case Cmd_SaveSna:
				data, err := speccy.Cpu.saveSna()
				if cmd.Chan != nil {
					cmd.Chan <- Snapshot{data, err}
				}
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
		rom48k, err := ioutil.ReadFile(speccy.romPath)
		if err != nil {
			return err
		}
		if len(rom48k) != 0x4000 {
			return os.NewError(fmt.Sprintf("ROM file \"%s\" has an invalid size", speccy.romPath))
		}

		for address, b := range rom48k {
			speccy.Memory.Write(uint16(address), b)
		}
	}

	return nil
}

func (speccy *Spectrum48k) Close() {
	speccy.Cpu.close()

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


// Initializes state from data in SNA format.
// Returns nil on success.
func (speccy *Spectrum48k) loadSna(data []byte) os.Error {
	err := speccy.reset()
	if err != nil {
		return err
	}

	err = speccy.Cpu.loadSna(data)
	if err != nil {
		return err
	}

	return nil
}
