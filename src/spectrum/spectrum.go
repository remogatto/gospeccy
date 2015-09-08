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
	"sync"
	"time"
	"github.com/remogatto/Go-PerfEvents"
	"github.com/remogatto/gospeccy/src/formats"
	"github.com/remogatto/z80"
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
	Cpu       *z80.Z80
	Memory    *Memory
	ula       *ULA
	Keyboard  *Keyboard
	Joystick  *Joystick
	tapeDrive *TapeDrive

	Ports *Ports

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

	readFromTape bool

	// The value is non-zero if a couple of the most recent frames
	// executed instructions which appeared to be reading from the tape
	shouldPlayTheTape int

	z80_instructionCounter     uint64 // Number of Z80 instructions executed
	z80_instructionsMeasured   uint64 // Number of Z80 instrs that can be related to 'hostCpu_instructionCounter'
	hostCpu_instructionCounter uint64
	perfCounter_hostCpuInstr   *perf.Counter // Can be nil (if creating the counter fails)
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
	z80 := z80.NewZ80(memory, ports)
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
	speccy.close()

	if speccy.app.Verbose {
		eff := speccy.GetEmulationEfficiency()
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
				cmd.Chan <- speccy.MakeSnapshot()

			case Cmd_MakeVideoMemoryDump:
				cmd.Chan <- speccy.makeVideoMemoryDump()

			case Cmd_SetAcceleratedLoad:
				speccy.tapeDrive.AcceleratedLoad = cmd.Enable

			}
		}
	}
}

func (speccy *Spectrum48k) reset(systemROMLoaded_orNil chan<- <-chan bool) error {
	speccy.Cpu.Reset()
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

// Returns the average number of host-CPU instructions required to execute one Z80 instruction.
// Returns zero if this information is not available.
func (speccy *Spectrum48k) GetEmulationEfficiency() uint {
	var eff uint
	if speccy.z80_instructionsMeasured > 0 {
		eff = uint(speccy.hostCpu_instructionCounter / speccy.z80_instructionsMeasured)
	} else {
		eff = 0
	}
	return eff
}

func (speccy *Spectrum48k) close() {
	if speccy.perfCounter_hostCpuInstr != nil {
		speccy.perfCounter_hostCpuInstr.Close()
		speccy.perfCounter_hostCpuInstr = nil
	}
}

// Initializes state from the specified snapshot.
// Returns nil on success.
func (speccy *Spectrum48k) loadSnapshot(s formats.Snapshot) error {
	speccy.reset(nil)

	cpu := s.CpuState()
	ula := s.UlaState()
	mem := s.Memory()

	// Populate registers
	speccy.Cpu.A = cpu.A
	speccy.Cpu.F = cpu.F
	speccy.Cpu.B = cpu.B
	speccy.Cpu.C = cpu.C
	speccy.Cpu.D = cpu.D
	speccy.Cpu.E = cpu.E
	speccy.Cpu.H = cpu.H
	speccy.Cpu.L = cpu.L
	speccy.Cpu.A_ = cpu.A_
	speccy.Cpu.F_ = cpu.F_
	speccy.Cpu.B_ = cpu.B_
	speccy.Cpu.C_ = cpu.C_
	speccy.Cpu.D_ = cpu.D_
	speccy.Cpu.E_ = cpu.E_
	speccy.Cpu.H_ = cpu.H_
	speccy.Cpu.L_ = cpu.L_
	speccy.Cpu.IXL = byte(cpu.IX & 0xff)
	speccy.Cpu.IXH = byte(cpu.IX >> 8)
	speccy.Cpu.IYL = byte(cpu.IY & 0xff)
	speccy.Cpu.IYH = byte(cpu.IY >> 8)

	speccy.Cpu.I = cpu.I
	speccy.Cpu.IFF1 = cpu.IFF1
	speccy.Cpu.IFF2 = cpu.IFF2
	speccy.Cpu.IM = cpu.IM

	speccy.Cpu.R = uint16(cpu.R & 0x7f)
	speccy.Cpu.R7 = cpu.R & 0x80

	speccy.Cpu.SetPC(cpu.PC)
	speccy.Cpu.SetSP(cpu.SP)

	// Border color
	speccy.Ports.WritePortInternal(0xfe, ula.Border&0x07, false /*contend*/)

	// Populate memory
	copy(speccy.Memory.Data()[0x4000:], mem[:])

	speccy.Cpu.Tstates = int(cpu.Tstate)

	return nil
}

func (speccy *Spectrum48k) MakeSnapshot() *formats.FullSnapshot {
	var s formats.FullSnapshot

	// Save registers
	s.Cpu.A = speccy.Cpu.A
	s.Cpu.F = speccy.Cpu.F
	s.Cpu.B = speccy.Cpu.B
	s.Cpu.C = speccy.Cpu.C
	s.Cpu.D = speccy.Cpu.D
	s.Cpu.E = speccy.Cpu.E
	s.Cpu.H = speccy.Cpu.H
	s.Cpu.L = speccy.Cpu.L
	s.Cpu.A_ = speccy.Cpu.A_
	s.Cpu.F_ = speccy.Cpu.F_
	s.Cpu.B_ = speccy.Cpu.B_
	s.Cpu.C_ = speccy.Cpu.C_
	s.Cpu.D_ = speccy.Cpu.D_
	s.Cpu.E_ = speccy.Cpu.E_
	s.Cpu.H_ = speccy.Cpu.H_
	s.Cpu.L_ = speccy.Cpu.L_
	s.Cpu.IX = uint16(speccy.Cpu.IXL) | (uint16(speccy.Cpu.IXH) << 8)
	s.Cpu.IY = uint16(speccy.Cpu.IYL) | (uint16(speccy.Cpu.IYH) << 8)

	s.Cpu.I = speccy.Cpu.I
	s.Cpu.IFF1 = speccy.Cpu.IFF1
	s.Cpu.IFF2 = speccy.Cpu.IFF2
	s.Cpu.IM = speccy.Cpu.IM

	s.Cpu.R = byte(speccy.Cpu.R & 0x7f) | (speccy.Cpu.R7 & 0x80)

	s.Cpu.SP = speccy.Cpu.SP()
	s.Cpu.PC = speccy.Cpu.PC()

	// Border color
	s.Ula.Border = speccy.ula.getBorderColor() & 0x07

	// Memory
	copy(s.Mem[:], speccy.Memory.Data()[0x4000:])

	return &s
}

func (speccy *Spectrum48k) doOpcodes() {
	var ttid_start int
	if speccy.perfCounter_hostCpuInstr != nil {
		ttid_start = speccy.perfCounter_hostCpuInstr.Gettid()
	} else {
		ttid_start = -1
	}

	var hostCpu_instrCount_start uint64 = 0
	var hostCpu_instrCount_startErr error = nil
	if speccy.perfCounter_hostCpuInstr != nil {
		hostCpu_instrCount_start, hostCpu_instrCount_startErr = speccy.perfCounter_hostCpuInstr.Read()
	}

	var z80_localInstructionCounter uint = 0

	// Main instruction emulation loop
	{
		var readFromTape bool = (speccy.readFromTape && (speccy.shouldPlayTheTape > 0) && (speccy.tapeDrive != nil))

		if speccy.tapeDrive != nil && speccy.tapeDrive.NotifyLoadComplete && speccy.tapeDrive.notifyCpuLoadCompleted {
			speccy.tapeDrive.notifyCpuLoadCompleted = false
			speccy.tapeDrive.loadComplete <- true
		}

		if !readFromTape {
			if speccy.tapeDrive != nil {
				speccy.tapeDrive.decelerate()
			}
		}

		for (speccy.Cpu.Tstates < speccy.Cpu.EventNextEvent) && !speccy.Cpu.Halted {
			speccy.Memory.ContendRead(speccy.Cpu.PC(), 4)
			opcode := speccy.Memory.ReadByteInternal(speccy.Cpu.PC())

			speccy.Cpu.R = (speccy.Cpu.R + 1) & 0x7f
			speccy.Cpu.IncPC(1)

			z80_localInstructionCounter++

			z80.OpcodesMap[opcode](speccy.Cpu)

			if readFromTape {
				endOfBlock := speccy.tapeDrive.doPlay()
				if endOfBlock {
					readFromTape = false
					speccy.shouldPlayTheTape = 0
					speccy.tapeDrive.decelerate()
				}
			}
		}

		if speccy.Cpu.Halted {
			speccy.shouldPlayTheTape = 0
			if speccy.tapeDrive != nil {
				speccy.tapeDrive.decelerate()
			}

			// Repeat emulating the HALT instruction until 'speccy.Cpu.eventNextEvent'
			for speccy.Cpu.Tstates < speccy.Cpu.EventNextEvent {
				speccy.Memory.ContendRead(speccy.Cpu.PC(), 4)

				speccy.Cpu.R = (speccy.Cpu.R + 1) & 0x7f
				z80_localInstructionCounter++
			}
		}
	}

	// Update emulation efficiency counters
	if speccy.perfCounter_hostCpuInstr != nil {
		ttid_end := speccy.perfCounter_hostCpuInstr.Gettid()

		var hostCpu_instrCount_end uint64
		var hostCpu_instrCount_endErr error
		hostCpu_instrCount_end, hostCpu_instrCount_endErr = speccy.perfCounter_hostCpuInstr.Read()

		speccy.z80_instructionCounter += uint64(z80_localInstructionCounter)

		/*if z80_localInstructionCounter > 0 {
		 println( z80_localInstructionCounter, hostCpu_instrCount_start, hostCpu_instrCount_end,
		 hostCpu_instrCount_end-hostCpu_instrCount_start,
		 (hostCpu_instrCount_end - hostCpu_instrCount_start) / uint64(z80_localInstructionCounter) )
		 }*/

		if (hostCpu_instrCount_startErr == nil) &&
			(hostCpu_instrCount_endErr == nil) &&
			(ttid_start == ttid_end) &&
			(z80_localInstructionCounter > 0) &&
			(hostCpu_instrCount_end > hostCpu_instrCount_start) {

			avg := uint((hostCpu_instrCount_end - hostCpu_instrCount_start) / uint64(z80_localInstructionCounter))

			// It may happen that the measured values are invalid.
			// The primary cause of this is that the Go runtime
			// can move a goroutine to a different OS thread,
			// without notifying us when it does so.
			// The majority of these cases is detected by (ttid_start == ttid_end) constraint.
			eff := speccy.GetEmulationEfficiency()
			bogusMeasurement := (avg < eff/4) || ((eff > 0) && (avg > eff*4))

			if !bogusMeasurement {
				speccy.z80_instructionsMeasured += uint64(z80_localInstructionCounter)
				speccy.hostCpu_instructionCounter += (hostCpu_instrCount_end - hostCpu_instrCount_start)
			}
		}
	}

}

func (speccy *Spectrum48k) renderFrame(completionTime_orNil chan<- time.Time) {
	speccy.Ports.frame_begin()
	speccy.ula.frame_begin()

	// Execute instructions corresponding to one screen frame
	speccy.Cpu.Tstates = (speccy.Cpu.Tstates % TStatesPerFrame)
	speccy.Cpu.Interrupt()
	speccy.Cpu.EventNextEvent = TStatesPerFrame
	speccy.doOpcodes()

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
		speccy.shouldPlayTheTape = 75
	} else {
		if speccy.shouldPlayTheTape > 0 {
			speccy.shouldPlayTheTape--
		}
	}
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
