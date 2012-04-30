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

// ZX Spectrum emulation core
package spectrum

import (
	"bytes"
	"errors"
	"github.com/remogatto/gospeccy/src/formats"
	"sync"
	"time"
)

const TStatesPerFrame = 69888 // Number of T-states per frame
const InterruptLength = 32    // How long does an interrupt last in T-states
const DefaultFPS = 50.08

type RomType int

const (
	ROM_UNKNOWN RomType = iota
	ROM_OPENSE          // OpenSE BASIC (http://www.worldofspectrum.org/infoseekid.cgi?id=0027510)
)

type DisplayInfo struct {
	displayReceiver DisplayReceiver

	// The index of the last frame sent to the 'displayReceiver', initially nil.
	lastFrame *uint

	// Number of frames sent to the DisplayReceiver
	numSentFrames uint

	// Total number of frames that the DisplayReceiver did not receive
	// because the send might block CPU emulation
	numMissedFrames uint

	missedChanges *DisplayData
}

type Spectrum48k struct {
	Cpu       *Z80
	Memory    MemoryAccessor
	ula       *ULA
	Keyboard  *Keyboard
	Joystick  *Joystick
	tapeDrive *TapeDrive

	Ports PortAccessor

	rom     [0x4000]byte
	romType RomType

	// The current display refresh frequency.
	// The initial value is 'DefaultFPS'.
	// It is always greater than 0.
	currentFPS       float32
	currentFPS_mutex sync.Mutex // To respect the Go memory model

	// A value received from this channel sets the display refresh frequency
	fpsCh chan float32

	// This buffered channel (if not nil) will receive at most one value.
	// The value 'true' sent through this channel indicates that the system ROM has been loaded.
	// The value 'false' sent through this channel indicates that the detection process did not finish.
	// The detection works with standard 48k ROM but probably doesn't work with custom ROMs.
	systemROMLoaded_orNil chan bool

	CommandChannel chan<- interface{}
	commandChannel <-chan interface{}

	// List of displays, initially empty
	displays []*DisplayInfo

	// List of audio receivers, initially empty
	audioReceivers []AudioReceiver

	// Register the state of FPS before accelerating tape loading
	fpsBeforeAccelerating float32

	app *Application
}

type Cmd_Reset struct {
	// This channel will receive [a channel X which will receive 'true'
	// when it is detected that the system ROM is loaded].
	// The channel X will receive 'false' if the emulated machine is reset before
	// the detection process ends.
	SystemROMLoaded_orNil chan<- <-chan bool
}
type Cmd_RenderFrame struct {
	// This channel (if not nil) will receive the real time when the rendering finished.
	//
	// If the list of host-machine displays is empty, the time only includes the emulation.
	// If the list of host-machine displays is non-empty, the time also includes host-machine
	// rendering. The sent time represents the moment of when the screen data reached
	// the host-machine display. On Linux this usually means: the moment right after
	// all pixels have been sent to the X server.
	// If there are multiple displays, the time includes the 1st display only.
	//
	// The time is obtained via a call to time.Nanoseconds().
	CompletionTime_orNil chan<- time.Time
}
type Cmd_GetNumDisplayReceivers struct {
	N chan<- uint
}
type Cmd_AddDisplay struct {
	Display DisplayReceiver
}
type Cmd_CloseAllDisplays struct {
	Finished chan<- byte
}
type Cmd_SetFPS struct {
	NewFPS       float32
	OldFPS_orNil chan<- float32
}
type Cmd_SetUlaEmulationAccuracy struct {
	AccurateEmulation bool
}
type Cmd_GetNumAudioReceivers struct {
	N chan<- uint
}
type Cmd_AddAudioReceiver struct {
	Receiver AudioReceiver
}
type Cmd_CloseAllAudioReceivers struct {
	Finished chan<- byte
}
type Cmd_LoadSnapshot struct {
	InformalFilename string // This is only used for logging purposes
	Snapshot         formats.Snapshot
	ErrChan          chan<- error
}
type Cmd_Load struct {
	InformalFilename string // This is only used for logging purposes
	Program          interface{}
	ErrChan          chan<- error
}
type Cmd_MakeSnapshot struct {
	Chan chan<- *formats.FullSnapshot
}
type Cmd_MakeVideoMemoryDump struct {
	Chan chan<- []byte
}
type Cmd_SetAcceleratedLoad struct {
	// Set accelerated tape load on/off
	Enable bool
}

// Creates a new speccy object and starts its command-loop goroutine.
//
// The returned object's CommandChannel can be used to
// configure the emulated machine before starting the emulation-loop
// and also to configure the machine while the emulation-loop is running.
//
// To start the actual emulation-loop, create a separate goroutine for
// running the object's EmulatorLoop function.
func NewSpectrum48k(app *Application, rom [0x4000]byte) *Spectrum48k {
	memory := NewMemory()
	keyboard := NewKeyboard()
	joystick := NewJoystick()
	ports := NewPorts()
	z80 := NewZ80(memory, ports)
	ula := NewULA()

	tapeDrive := NewTapeDrive()

	speccy := &Spectrum48k{
		Cpu:            z80,
		Memory:         memory,
		ula:            ula,
		Keyboard:       keyboard,
		Joystick:       joystick,
		Ports:          ports,
		rom:            rom,
		romType:        ROM_UNKNOWN,
		displays:       make([]*DisplayInfo, 0),
		audioReceivers: make([]AudioReceiver, 0),
		app:            app,
		tapeDrive:      tapeDrive,
	}

	memory.init(speccy)
	keyboard.init(speccy)
	joystick.init(speccy)
	z80.init(ula, tapeDrive)
	ula.init(z80, memory, ports)
	ports.init(speccy)
	tapeDrive.init(speccy)

	speccy.reset(nil)

	speccy.currentFPS = DefaultFPS
	speccy.fpsCh = make(chan float32, 1)
	speccy.fpsCh <- DefaultFPS

	commandChannel := make(chan interface{})
	speccy.CommandChannel = commandChannel
	speccy.commandChannel = commandChannel
	go commandLoop(speccy)

	return speccy
}

// Turn off the machine
func (speccy *Spectrum48k) Close() {
	speccy.Cpu.close()

	if speccy.app.Verbose {
		eff := speccy.Cpu.GetEmulationEfficiency()
		if eff != 0 {
			speccy.app.PrintfMsg("emulation efficiency: %d host-CPU instructions per Z80 instruction", eff)
		} else {
			speccy.app.PrintfMsg("emulation efficiency: -")
		}

		for i, display := range speccy.displays {
			nSent := display.numSentFrames
			nMissed := display.numMissedFrames
			speccy.app.PrintfMsg("display #%d: %d shown frames, %d missed frames", i, nSent, nMissed)
		}
	}
}

// Get current FPS
func (speccy *Spectrum48k) GetCurrentFPS() float32 {
	speccy.currentFPS_mutex.Lock()
	fps := speccy.currentFPS
	speccy.currentFPS_mutex.Unlock()
	return fps
}

// Load a program (tape or snapshot)
func (speccy *Spectrum48k) load(program interface{}) error {
	var err error

	switch program := program.(type) {

	case formats.Snapshot:
		speccy.loadSnapshot(program.(formats.Snapshot))
	case *formats.TAP:
		speccy.loadTape(program)
	default:
		err = errors.New("Invalid program type.")
		return err
	}

	return err
}

// Return the TapeDrive instance
func (speccy *Spectrum48k) TapeDrive() *TapeDrive {
	return speccy.tapeDrive
}

// Sends 'Cmd_RenderFrame' commands to the 'speccy' object in regular intervals.
// The interval depends on the value of FPS (frames per second).
//
// This function should run in a separate goroutine.
func (speccy *Spectrum48k) EmulatorLoop() {
	evtLoop := speccy.app.NewEventLoop()
	app := evtLoop.App()

	fps := <-speccy.fpsCh
	ticker := time.NewTicker(time.Duration(1e9 / fps))

	// Render the 1st frame (the 2nd frame will be rendered after 1/FPS seconds)
	{
		completionTime := make(chan time.Time)
		speccy.CommandChannel <- Cmd_RenderFrame{completionTime}

		go func() {
			start := app.CreationTime
			end := <-completionTime
			if app.Verbose {
				app.PrintfMsg("first frame latency: %s", end.Sub(start))
			}
		}()
	}

	var newFPS_orMinusOne float32 = -1

	for {
		select {
		case <-evtLoop.Pause:
			ticker.Stop()
			Drain(ticker)
			evtLoop.Pause <- 0

		case <-evtLoop.Terminate:
			// Terminate this Go routine
			if app.Verbose {
				app.PrintfMsg("emulator loop: exit")
			}
			evtLoop.Terminate <- 0
			return

		case <-ticker.C:
			if newFPS_orMinusOne != -1 {
				newFPS := newFPS_orMinusOne
				newFPS_orMinusOne = -1

				if app.Verbose {
					app.PrintfMsg("setting FPS to %f", newFPS)
				}
				ticker.Stop()
				Drain(ticker)
				ticker = time.NewTicker(time.Duration(1e9 / newFPS))
				fps = newFPS
			}

			//app.PrintfMsg("%d", time.Now().UnixNano()/1e6)
			speccy.CommandChannel <- Cmd_RenderFrame{}

		case newFPS := <-speccy.fpsCh:
			if (newFPS != fps) && (newFPS > 0) {
				newFPS_orMinusOne = newFPS
			}
		}
	}
}

func commandLoop(speccy *Spectrum48k) {
	evtLoop := speccy.app.NewEventLoop()
	for {
		select {
		case <-evtLoop.Pause:
			// Unblock the goroutine that is waiting for the end of ROM initialization
			if speccy.systemROMLoaded_orNil != nil {
				// Note: This is a buffered channel, so the send won't block
				speccy.systemROMLoaded_orNil <- false
				speccy.systemROMLoaded_orNil = nil
			}

			speccy.Close()
			evtLoop.Pause <- 0

		case <-evtLoop.Terminate:
			// Terminate this Go routine
			if evtLoop.App().Verbose {
				evtLoop.App().PrintfMsg("command loop: exit")
			}
			evtLoop.Terminate <- 0
			return

		case untyped_cmd := <-speccy.commandChannel:
			switch cmd := untyped_cmd.(type) {
			case Cmd_Reset:
				speccy.reset(cmd.SystemROMLoaded_orNil)

			case Cmd_RenderFrame:
				// Ugly hack to check whenever the system ROM has been loaded after a reset.
				// I bet this won't work with custom ROMs.
				if (speccy.Cpu.PC() == 0x10ac) && (speccy.systemROMLoaded_orNil != nil) {
					// Note: This is a buffered channel, so the send won't block
					speccy.systemROMLoaded_orNil <- true
					speccy.systemROMLoaded_orNil = nil
				}

				speccy.renderFrame(cmd.CompletionTime_orNil)

			case Cmd_GetNumDisplayReceivers:
				cmd.N <- uint(len(speccy.displays))

			case Cmd_AddDisplay:
				speccy.addDisplay(cmd.Display)

			case Cmd_CloseAllDisplays:
				go func() {
					speccy.closeAllDisplays()
					cmd.Finished <- 0
				}()

			case Cmd_SetFPS:
				speccy.currentFPS_mutex.Lock()
				{
					if cmd.OldFPS_orNil != nil {
						cmd.OldFPS_orNil <- speccy.currentFPS
					}

					newFPS := cmd.NewFPS
					if newFPS <= 1.0 {
						newFPS = DefaultFPS
					}

					if newFPS != speccy.currentFPS {
						speccy.currentFPS = newFPS

						go func() {
							speccy.fpsCh <- newFPS
						}()
					}
				}
				speccy.currentFPS_mutex.Unlock()

			case Cmd_SetUlaEmulationAccuracy:
				speccy.ula.setEmulationAccuracy(cmd.AccurateEmulation)

			case Cmd_GetNumAudioReceivers:
				cmd.N <- uint(len(speccy.audioReceivers))

			case Cmd_AddAudioReceiver:
				speccy.addAudioReceiver(cmd.Receiver)

			case Cmd_CloseAllAudioReceivers:
				go func() {
					speccy.closeAllAudioReceivers()
					cmd.Finished <- 0
				}()

			case Cmd_LoadSnapshot:
				if speccy.app.Verbose {
					if len(cmd.InformalFilename) > 0 {
						speccy.app.PrintfMsg("loading snapshot \"%s\"", cmd.InformalFilename)
					} else {
						speccy.app.PrintfMsg("loading a snapshot")
					}
				}

				err := speccy.loadSnapshot(cmd.Snapshot)

				if cmd.ErrChan != nil {
					cmd.ErrChan <- err
				}

			case Cmd_Load:
				if speccy.app.Verbose {
					if len(cmd.InformalFilename) > 0 {
						speccy.app.PrintfMsg("loading program \"%s\"", cmd.InformalFilename)
					} else {
						speccy.app.PrintfMsg("loading a program")
					}
				}

				err := speccy.load(cmd.Program)

				if cmd.ErrChan != nil {
					cmd.ErrChan <- err
				}

			case Cmd_MakeSnapshot:
				cmd.Chan <- speccy.Cpu.MakeSnapshot()

			case Cmd_MakeVideoMemoryDump:
				cmd.Chan <- speccy.makeVideoMemoryDump()

			case Cmd_SetAcceleratedLoad:
				speccy.tapeDrive.AcceleratedLoad = cmd.Enable

			}
		}
	}
}

func (speccy *Spectrum48k) reset(systemROMLoaded_orNil chan<- <-chan bool) error {
	speccy.Cpu.reset()
	speccy.Memory.reset()
	speccy.ula.reset()
	speccy.Keyboard.reset()
	speccy.Ports.reset()

	if speccy.systemROMLoaded_orNil != nil {
		speccy.systemROMLoaded_orNil <- false
		speccy.systemROMLoaded_orNil = nil
	}
	if systemROMLoaded_orNil != nil {
		// Create a buffered channel and send it to the goroutine which requested the reset
		speccy.systemROMLoaded_orNil = make(chan bool, 1)
		systemROMLoaded_orNil <- speccy.systemROMLoaded_orNil
	}

	// Copy the ROM image into the first 16k of memory
	copy(speccy.Memory.Data()[0:0x4000], speccy.rom[:])

	// ROM type detection
	if bytes.Contains(speccy.rom[:], []byte("1981 Nine Tiles Networks")) {
		speccy.romType = ROM_OPENSE
	}

	// OpenSE BASIC initializes almost immediately
	if (speccy.systemROMLoaded_orNil != nil) && (speccy.romType == ROM_OPENSE) {
		speccy.systemROMLoaded_orNil <- true
		speccy.systemROMLoaded_orNil = nil
	}

	return nil
}

func (speccy *Spectrum48k) addDisplay(display DisplayReceiver) {
	d := &DisplayInfo{
		displayReceiver: display,
		lastFrame:       nil,
		missedChanges:   nil,
	}

	speccy.displays = append(speccy.displays, d)
}

func (speccy *Spectrum48k) closeAllDisplays() {
	displays := speccy.displays
	speccy.displays = make([]*DisplayInfo, 0)

	for i, d := range displays {
		d.displayReceiver.Close()
		if speccy.app.Verbose {
			nSent := d.numSentFrames
			nMissed := d.numMissedFrames
			speccy.app.PrintfMsg("display #%d: %d shown frames, %d missed frames", i, nSent, nMissed)
		}
	}
}

func (speccy *Spectrum48k) addAudioReceiver(receiver AudioReceiver) {
	speccy.audioReceivers = append(speccy.audioReceivers, receiver)
}

func (speccy *Spectrum48k) closeAllAudioReceivers() {
	audioReceivers := speccy.audioReceivers
	speccy.audioReceivers = make([]AudioReceiver, 0)

	for _, r := range audioReceivers {
		r.Close()
	}
}

func (speccy *Spectrum48k) renderFrame(completionTime_orNil chan<- time.Time) {
	speccy.Ports.frame_begin()
	speccy.ula.frame_begin()

	// Execute instructions corresponding to one screen frame
	speccy.Cpu.tstates = (speccy.Cpu.tstates % TStatesPerFrame)
	speccy.Cpu.interrupt()
	speccy.Cpu.eventNextEvent = TStatesPerFrame
	speccy.Cpu.doOpcodes()

	// Send display data to display backend(s)
	if len(speccy.displays) > 0 {
		firstDisplay := true
		for _, display := range speccy.displays {
			var tm chan<- time.Time
			if firstDisplay {
				tm = completionTime_orNil
			} else {
				tm = nil
			}
			speccy.ula.sendScreenToDisplay(display, tm)
			firstDisplay = false
		}
	} else {
		if completionTime_orNil != nil {
			completionTime_orNil <- time.Now()
		}
	}

	// Send audio data to audio backend(s)
	if len(speccy.audioReceivers) > 0 {
		audioData := AudioData{
			FPS:          speccy.currentFPS,
			BeeperEvents: speccy.Ports.getBeeperEvents(),
		}

		for _, audioReceiver := range speccy.audioReceivers {
			audioReceiver.GetAudioDataChannel() <- &audioData
		}
	}

	portFrameStatus := speccy.Ports.frame_end()

	if portFrameStatus.shouldPlayTheTape {
		speccy.Cpu.shouldPlayTheTape = 75
	} else {
		if speccy.Cpu.shouldPlayTheTape > 0 {
			speccy.Cpu.shouldPlayTheTape--
		}
	}
}

// Initializes state from the specified snapshot.
// Returns nil on success.
func (speccy *Spectrum48k) loadSnapshot(s formats.Snapshot) error {
	speccy.reset(nil)

	err := speccy.Cpu.loadSnapshot(s)
	if err != nil {
		return err
	}

	return nil
}

// Load the given tape
func (speccy *Spectrum48k) loadTape(tap *formats.TAP) {
	speccy.tapeDrive.Insert(NewTape(tap))
	speccy.tapeDrive.Stop()
	speccy.sendLOADCommand()
	speccy.tapeDrive.Play()
}

// Send LOAD ""
func (speccy *Spectrum48k) sendLOADCommand() {
	speccy.Keyboard.CommandChannel <- Cmd_SendLoad{speccy.romType}
}

func (speccy *Spectrum48k) makeVideoMemoryDump() []byte {
	return speccy.Memory.Data()[0x4000 : 0x4000+6912]
}
