package spectrum

import (
	"os"
	"io/ioutil"
	"spectrum/formats"
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
	TAPE_LEADER = 2168
	TAPE_FIRST_SYNC = 667
	TAPE_SECOND_SYNC = 735
	TAPE_SET_BIT = 1710
	TAPE_UNSET_BIT = 855
	TAPE_HEADER_LEADER_PULSES = 8063
	TAPE_DATA_LEADER_PULSES = 3223
	TAPE_PAUSE = 3500000
)

type Tape struct {
	tap *formats.TAP
}

func NewTape(tap *formats.TAP) *Tape {
	return &Tape{ tap }
}

func NewTapeFromFile(filename string) (*Tape, os.Error) {
	data, err := ioutil.ReadFile(filename)
	tap := formats.NewTAP()
	_, err = tap.Read(data)
	tape := &Tape{ tap }
	return tape, err
}

func (tape *Tape) At(pos uint) byte {
	return tape.tap.At(pos)
}

type TapeDrive struct {
	Status byte
	tape *Tape
	pos int
	tstate, lastIn uint64
	earBit byte
	timeout int
	timeLastIn, currBlockLen, currBlockId int
	leaderPulses, bitTime uint16
	state, mask byte
	speccy *Spectrum48k
	loadComplete chan bool
}

func NewTapeDrive() *TapeDrive {
	return &TapeDrive{

	Status: TAPE_DRIVE_STOP, 
	pos: 0,
	earBit: 0xbf,
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
	tapeDrive.speccy.Cpu.readFromTape = true
	tapeDrive.pos = 0
	tapeDrive.state = TAPE_DRIVE_START
	tapeDrive.timeout = 0
	tapeDrive.timeLastIn = 0
}

func (tapeDrive *TapeDrive) Stop() {
	tapeDrive.speccy.Cpu.readFromTape = false
	tapeDrive.pos = 0
	tapeDrive.state = TAPE_DRIVE_PAUSE_STOP
	tapeDrive.timeout = 0
	tapeDrive.timeLastIn = 0
	tapeDrive.currBlockId = 0
}

func (tapeDrive *TapeDrive) doPlay() {
        now := tapeDrive.speccy.ula.frame * TStatesPerFrame + tapeDrive.speccy.Cpu.tstates

        tapeDrive.timeout -= (int(now) - tapeDrive.timeLastIn)
        tapeDrive.timeLastIn = int(now)

        if tapeDrive.timeout > 0 {
		return
        }

        tapeDrive.timeout = 0

        switch tapeDrive.state {
        case TAPE_DRIVE_START:
                tapeDrive.currBlockLen = tapeDrive.tape.tap.GetBlock(tapeDrive.currBlockId).Len()

		if tapeDrive.currBlockId == 0 {
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
			break
                }
                tapeDrive.timeout = TAPE_FIRST_SYNC
                tapeDrive.state = TAPE_DRIVE_SYNC

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
                if ((tapeDrive.tape.At(uint(tapeDrive.pos)) & tapeDrive.mask) == 0) {
			tapeDrive.bitTime = TAPE_UNSET_BIT
                } else {
			tapeDrive.bitTime = TAPE_SET_BIT
                }
                tapeDrive.timeout = int(tapeDrive.bitTime)
                tapeDrive.state = TAPE_DRIVE_HALF2
                break
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
                break;
        case TAPE_DRIVE_PAUSE:
		if tapeDrive.earBit == 0xbf {
			tapeDrive.earBit = 0xff
		} else {
			tapeDrive.earBit = 0xbf
		}
                tapeDrive.timeout = TAPE_PAUSE
                tapeDrive.state = TAPE_DRIVE_PAUSE_STOP
                break
        case TAPE_DRIVE_PAUSE_STOP:
                if tapeDrive.pos == int(tapeDrive.tape.tap.Len()) {
			tapeDrive.state = TAPE_DRIVE_STOP
			tapeDrive.timeout = 1
			tapeDrive.speccy.Cpu.readFromTape = false
			// tapeDrive.loadComplete <- true
                } else {
			tapeDrive.currBlockId++
			tapeDrive.state = TAPE_DRIVE_START
                }
        }

}

func (tapeDrive *TapeDrive) getEarBit() uint8 {
	return tapeDrive.earBit
}
