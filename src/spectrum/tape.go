package spectrum

import (
	"github.com/remogatto/gospeccy/src/formats"
	"io/ioutil"
	"sync"
)

const (
	TAPE_DRIVE_START = iota
	TAPE_DRIVE_STOP
	TAPE_DRIVE_PAUSE
	TAPE_DRIVE_LEADER
	TAPE_DRIVE_SYNC
	TAPE_DRIVE_NEWBIT
	TAPE_DRIVE_HALF2
	TAPE_DRIVE_PAUSE_STOP
	TAPE_DRIVE_NEWBYTE
)

const (
	TAPE_LEADER               = 2168
	TAPE_FIRST_SYNC           = 667
	TAPE_SECOND_SYNC          = 735
	TAPE_SET_BIT              = 1710
	TAPE_UNSET_BIT            = 855
	TAPE_HEADER_LEADER_PULSES = 8063
	TAPE_DATA_LEADER_PULSES   = 3223
	TAPE_PAUSE                = 3500000
)

const TAPE_ACCELERATION_IN_FPS = DefaultFPS * 20

type Tape struct {
	tap *formats.TAP
}

func NewTape(tap *formats.TAP) *Tape {
	return &Tape{tap}
}

func NewTapeFromFile(filename string) (*Tape, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	tap, err := formats.NewTAP(data)
	if err != nil {
		return nil, err
	}

	return &Tape{tap}, nil
}

func (tape *Tape) At(pos uint) byte {
	return tape.tap.At(pos)
}

type TapeDrive struct {
	AcceleratedLoad    bool
	NotifyLoadComplete bool

	speccy *Spectrum48k
	tape   *Tape

	pos                                   uint
	tstate, lastIn                        uint64
	earBit                                byte
	timeout                               int
	timeLastIn, currBlockLen, currBlockId int
	leaderPulses, bitTime                 uint16
	state, mask                           byte
	accelerating                          bool
	fpsBeforeAcceleration                 float32
	notifyCpuLoadCompleted                bool
	loadComplete                          chan bool

	mutex sync.RWMutex
}

func NewTapeDrive() *TapeDrive {
	return &TapeDrive{
		pos:          0,
		earBit:       0xbf,
		loadComplete: make(chan bool),
	}
}

func (tapeDrive *TapeDrive) init(speccy *Spectrum48k) {
	tapeDrive.speccy = speccy
}

func (tapeDrive *TapeDrive) Insert(tape *Tape) {
	tapeDrive.tape = tape
}

func (tapeDrive *TapeDrive) Play() {
	tapeDrive.speccy.readFromTape = true
	tapeDrive.pos = 0
	tapeDrive.state = TAPE_DRIVE_START
	tapeDrive.timeout = 0
	tapeDrive.timeLastIn = 0
}

func (tapeDrive *TapeDrive) Stop() {
	tapeDrive.speccy.readFromTape = false
	tapeDrive.pos = 0
	tapeDrive.state = TAPE_DRIVE_PAUSE_STOP
	tapeDrive.timeout = 0
	tapeDrive.timeLastIn = 0
	tapeDrive.currBlockId = 0
}

func (tapeDrive *TapeDrive) accelerate() {
	if !tapeDrive.accelerating {
		tapeDrive.accelerating = true
		go func() {
			oldFPS_chan := make(chan float32)
			tapeDrive.speccy.CommandChannel <- Cmd_SetFPS{TAPE_ACCELERATION_IN_FPS, oldFPS_chan}
			oldFPS := <-oldFPS_chan

			tapeDrive.mutex.Lock()
			tapeDrive.fpsBeforeAcceleration = oldFPS
			tapeDrive.mutex.Unlock()
		}()
	}
}

func (tapeDrive *TapeDrive) decelerate() {
	if tapeDrive.accelerating {
		tapeDrive.accelerating = false
		go func() {
			tapeDrive.mutex.RLock()
			fps := tapeDrive.fpsBeforeAcceleration
			tapeDrive.mutex.RUnlock()

			tapeDrive.speccy.CommandChannel <- Cmd_SetFPS{fps, nil}
		}()
	}
}

func (tapeDrive *TapeDrive) doPlay() (endOfBlock bool) {
	now := int(tapeDrive.speccy.ula.frame)*TStatesPerFrame + tapeDrive.speccy.Cpu.Tstates

	tapeDrive.timeout -= now - tapeDrive.timeLastIn
	tapeDrive.timeLastIn = now

	if tapeDrive.timeout > 0 {
		return
	}

	tapeDrive.timeout = 0

	if tapeDrive.AcceleratedLoad {
		tapeDrive.accelerate()
	} else {
		tapeDrive.decelerate()
	}

	switch tapeDrive.state {
	case TAPE_DRIVE_START:
		currBlock := tapeDrive.tape.tap.GetBlock(tapeDrive.currBlockId)
		tapeDrive.currBlockLen = currBlock.Len()

		if currBlock.BlockType() == formats.TAP_BLOCK_HEADER {
			tapeDrive.leaderPulses = TAPE_HEADER_LEADER_PULSES
		} else {
			tapeDrive.leaderPulses = TAPE_DATA_LEADER_PULSES
		}

		tapeDrive.earBit = 0xbf
		tapeDrive.timeout = TAPE_LEADER
		tapeDrive.state = TAPE_DRIVE_LEADER
	case TAPE_DRIVE_LEADER:
		if tapeDrive.earBit == 0xbf {
			tapeDrive.earBit = 0xff
		} else {
			tapeDrive.earBit = 0xbf
		}
		tapeDrive.leaderPulses--
		if tapeDrive.leaderPulses > 0 {
			tapeDrive.timeout = TAPE_LEADER
		} else {
			tapeDrive.timeout = TAPE_FIRST_SYNC
			tapeDrive.state = TAPE_DRIVE_SYNC
		}

	case TAPE_DRIVE_SYNC:
		if tapeDrive.earBit == 0xbf {
			tapeDrive.earBit = 0xff
		} else {
			tapeDrive.earBit = 0xbf
		}
		tapeDrive.timeout = TAPE_SECOND_SYNC
		tapeDrive.state = TAPE_DRIVE_NEWBYTE

	case TAPE_DRIVE_NEWBYTE:
		tapeDrive.mask = 0x80
		fallthrough

	case TAPE_DRIVE_NEWBIT:
		if tapeDrive.earBit == 0xbf {
			tapeDrive.earBit = 0xff
		} else {
			tapeDrive.earBit = 0xbf
		}
		if (tapeDrive.tape.At(tapeDrive.pos) & tapeDrive.mask) == 0 {
			tapeDrive.bitTime = TAPE_UNSET_BIT
		} else {
			tapeDrive.bitTime = TAPE_SET_BIT
		}
		tapeDrive.timeout = int(tapeDrive.bitTime)
		tapeDrive.state = TAPE_DRIVE_HALF2

	case TAPE_DRIVE_HALF2:
		if tapeDrive.earBit == 0xbf {
			tapeDrive.earBit = 0xff
		} else {
			tapeDrive.earBit = 0xbf
		}
		tapeDrive.timeout = int(tapeDrive.bitTime)
		tapeDrive.mask >>= 1
		if tapeDrive.mask == 0 {
			tapeDrive.pos++
			tapeDrive.currBlockLen--
			if tapeDrive.currBlockLen > 0 {
				tapeDrive.state = TAPE_DRIVE_NEWBYTE
			} else {
				tapeDrive.state = TAPE_DRIVE_PAUSE
			}
		} else {
			tapeDrive.state = TAPE_DRIVE_NEWBIT
		}

	case TAPE_DRIVE_PAUSE:
		if tapeDrive.earBit == 0xbf {
			tapeDrive.earBit = 0xff
		} else {
			tapeDrive.earBit = 0xbf
		}

		endOfBlock = true
		tapeDrive.decelerate()

		if tapeDrive.pos < tapeDrive.tape.tap.Len() {
			tapeDrive.timeout = TAPE_PAUSE
			tapeDrive.state = TAPE_DRIVE_PAUSE_STOP
		} else {
			tapeDrive.timeout = 1
			tapeDrive.state = TAPE_DRIVE_STOP

			tapeDrive.speccy.readFromTape = false
			tapeDrive.notifyCpuLoadCompleted = true
		}

	case TAPE_DRIVE_PAUSE_STOP:
		tapeDrive.currBlockId++
		tapeDrive.state = TAPE_DRIVE_START
	}

	return endOfBlock
}

func (tapeDrive *TapeDrive) getEarBit() uint8 {
	return tapeDrive.earBit
}

func (tapeDrive *TapeDrive) LoadComplete() <-chan bool {
	return tapeDrive.loadComplete
}
