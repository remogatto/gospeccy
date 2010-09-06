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

	// Total number of frames that the DisplayReceiver did not receive
	// because the send might block CPU emulation
	numMissedFrames uint
}

type Spectrum48k struct {
	Cpu      *Z80
	Memory   MemoryAccessor
	ula      *ULA
	Keyboard *Keyboard
	Ports    PortAccessor

	romPath string

	// The current display refresh frequency.
	// The initial value if 'DefaultFPS'.
	// It is always greater than 0.
	currentFPS float

	// A value received from this channel indicates the new display refresh frequency.
	// By default, this channel initially receives the value 'DefaultFPS'.
	FPS <-chan float
	fps chan float

	CommandChannel chan<- interface{}
	commandChannel <-chan interface{}

	// A vector of '*DisplayInfo', initially empty
	displays vector.Vector

	// A vector of 'AudioReceiver', initially empty
	audioReceivers vector.Vector

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
		Cpu:            z80,
		Memory:         memory,
		ula:            ula,
		Keyboard:       keyboard,
		Ports:          ports,
		romPath:        romPath,
		displays:       vector.Vector{},
		audioReceivers: vector.Vector{},
		app:            app,
	}

	memory.init(speccy)
	z80.init(speccy)
	ula.init(speccy)
	ports.init(speccy)

	err := speccy.reset()
	if err != nil {
		return nil, err
	}

	speccy.currentFPS = DefaultFPS
	speccy.fps = make(chan float, 1)
	speccy.FPS = speccy.fps
	speccy.fps <- DefaultFPS

	commandChannel := make(chan interface{})
	speccy.CommandChannel = commandChannel
	speccy.commandChannel = commandChannel
	go commandLoop(speccy)

	return speccy, nil
}

type Cmd_Reset struct{}
type Cmd_RenderFrame struct {
	// This channel (if not nil) will receive the time when the WHOLE rendering finished.
	// The time is obtained via a call to time.Nanoseconds().
	CompletionTime_orNil chan<- int64
}
type Cmd_AddDisplay struct {
	Display DisplayReceiver
}
type Cmd_CloseAllDisplays struct{}
type Cmd_SetFPS struct {
	NewFPS float
}
type Cmd_SetUlaEmulationAccuracy struct {
	AccurateEmulation bool
}
type Cmd_AddAudioReceiver struct {
	Receiver AudioReceiver
}
type Cmd_CloseAllAudioReceivers struct{}
type Cmd_LoadSna struct {
	InformalFilename string // This is only used for logging purposes
	Data             []byte // The SNA snapshot data
	ErrChan          chan<- os.Error
}
type Snapshot struct {
	Data []byte // Constraint: (data == nil) != (err == nil)
	Err  os.Error
}
type Cmd_SaveSna struct {
	Chan chan<- Snapshot
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
				PrintfMsg("command loop: exit")
			}
			evtLoop.Terminate <- 0
			return

		case untyped_cmd := <-speccy.commandChannel:
			switch cmd := untyped_cmd.(type) {
			case Cmd_Reset:
				speccy.reset()

			case Cmd_RenderFrame:
				speccy.renderFrame(cmd.CompletionTime_orNil)

			case Cmd_AddDisplay:
				speccy.addDisplay(cmd.Display)

			case Cmd_CloseAllDisplays:
				speccy.closeAllDisplays()

			case Cmd_SetFPS:
				speccy.setFPS(cmd.NewFPS)

			case Cmd_SetUlaEmulationAccuracy:
				speccy.ula.setEmulationAccuracy(cmd.AccurateEmulation)

			case Cmd_AddAudioReceiver:
				speccy.addAudioReceiver(cmd.Receiver)

			case Cmd_CloseAllAudioReceivers:
				speccy.closeAllAudioReceivers()

			case Cmd_LoadSna:
				if speccy.app.Verbose {
					if len(cmd.InformalFilename) > 0 {
						PrintfMsg("loading SNA snapshot \"%s\"", cmd.InformalFilename)
					} else {
						PrintfMsg("loading a SNA snapshot")
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
			PrintfMsg("emulation efficiency: %d host-CPU instructions per Z80 instruction", eff)
		} else {
			PrintfMsg("emulation efficiency: -")
		}

		for i, display := range speccy.displays {
			PrintfMsg("display #%d: %d missed frames", i, display.(*DisplayInfo).numMissedFrames)
		}
	}
}

func (speccy *Spectrum48k) addDisplay(display DisplayReceiver) {
	d := &DisplayInfo{displayReceiver: display, lastFrame: nil, numMissedFrames: 0}
	speccy.displays.Push(d)
}

func (speccy *Spectrum48k) closeAllDisplays() {
	var displays vector.Vector
	{
		displays = speccy.displays
		speccy.displays = vector.Vector{}
	}

	for _, d := range displays {
		d.(*DisplayInfo).displayReceiver.close()
	}
}

func (speccy *Spectrum48k) setFPS(newFPS float) {
	if newFPS <= 1.0 {
		newFPS = DefaultFPS
	}

	if newFPS != speccy.currentFPS {
		speccy.currentFPS = newFPS
		speccy.fps <- newFPS
	}
}

func (speccy *Spectrum48k) addAudioReceiver(receiver AudioReceiver) {
	speccy.audioReceivers.Push(receiver)
}

func (speccy *Spectrum48k) closeAllAudioReceivers() {
	var audioReceivers vector.Vector
	{
		audioReceivers = speccy.audioReceivers
		speccy.audioReceivers = vector.Vector{}
	}

	for _, r := range audioReceivers {
		r.(AudioReceiver).close()
	}
}

// Execute the number of T-states corresponding to one screen frame
func (speccy *Spectrum48k) doOpcodes() {
	speccy.Cpu.eventNextEvent = TStatesPerFrame
	speccy.Cpu.tstates = (speccy.Cpu.tstates % TStatesPerFrame)
	speccy.Cpu.doOpcodes()
}

func (speccy *Spectrum48k) renderFrame(completionTime_orNil chan<- int64) {
	speccy.Ports.frame_begin()
	speccy.ula.frame_begin()

	speccy.Cpu.interrupt()
	speccy.doOpcodes()

	// Send display data to display backend(s)
	firstDisplay := true
	for _, display := range speccy.displays {
		var tm chan<- int64
		if firstDisplay {
			tm = completionTime_orNil
		} else {
			tm = nil
		}
		speccy.ula.sendScreenToDisplay(display.(*DisplayInfo), tm)
		firstDisplay = false
	}

	// Send audio data to audio backend(s)
	{
		audioData := AudioData{
			fps:          speccy.currentFPS,
			beeperEvents: speccy.Ports.getBeeperEvents(),
		}

		for _, audioReceiver := range speccy.audioReceivers {
			audioReceiver.(AudioReceiver).getAudioDataChannel() <- &audioData
		}
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
