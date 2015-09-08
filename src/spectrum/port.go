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

package spectrum

import "github.com/remogatto/z80"

type FrameStatusOfPorts struct {
	shouldPlayTheTape bool
}

type BorderEvent struct {
	// The moment when the border color was changed.
	// It is the number of T-states since the beginning of the frame.
	TState int

	// The new border color
	Color byte
}

func (e *BorderEvent) GetTState() int {
	return e.TState
}

type BeeperEvent struct {
	// The moment when the beeper-event occurred.
	// It is the number of T-states since the beginning of the frame.
	TState int

	// The beeper level (0 .. MAX_AUDIO_LEVEL)
	Level byte
}

func (e *BeeperEvent) GetTState() int {
	return e.TState
}

type Ports struct {
	speccy *Spectrum48k

	borderEvents []BorderEvent

	beeperLevel  byte
	beeperEvents []BeeperEvent

	// Number of supposed reads from tapedrive port.
	// This counter is reset to 0 at the beginning of each frame.
	tapeReadCount uint
}

// If 'tapeReadCount' is equal to or above this threshold,
// the program running within the emulated machine probably wants to read data from the tape
const tapeReadCount_tapeAccessThreshold = 400

func NewPorts() *Ports {
	p := &Ports{}
	p.borderEvents = []BorderEvent{}
	p.beeperLevel = 0
	p.beeperEvents = []BeeperEvent{{TState: 0, Level: p.beeperLevel}}

	return p
}

func (p *Ports) init(speccy *Spectrum48k) {
	p.speccy = speccy
}

func (p *Ports) reset() {
	p.borderEvents = p.borderEvents[0:0]
	p.borderEvents = append(p.borderEvents, BorderEvent{TState: 0, Color: p.speccy.ula.getBorderColor()})

	p.beeperLevel = 0
	p.beeperEvents = p.beeperEvents[0:0]
	p.beeperEvents = append(p.beeperEvents, BeeperEvent{TState: 0, Level: p.beeperLevel})
}

func SameBorderEvents(l1, l2 []BorderEvent) bool {
	if len(l1) != len(l2) {
		return false
	}

	n := len(l1)
	for i := 0; i < n; i++ {
		e1 := l1[i]
		e2 := l2[i]
		if (e1.TState != e2.TState) || (e1.Color != e2.Color) {
			return false
		}
	}

	return true
}

func (p *Ports) frame_begin() {
	p.tapeReadCount = 0
}

func (p *Ports) frame_end() FrameStatusOfPorts {
	// Border events
	{
		// Determine the number of events overflowing the frame
		var numOverflow int
		{
			i := len(p.borderEvents)
			for (i > 0) && (p.borderEvents[i-1].TState >= TStatesPerFrame) {
				i--
			}
			numOverflow = len(p.borderEvents) - i
		}

		numEvents := len(p.borderEvents) - numOverflow

		overflow := make([]BorderEvent, numOverflow, numOverflow)
		copy(overflow[:], p.borderEvents[numEvents:])

		var colorAtTState0 byte
		if numOverflow == 0 {
			colorAtTState0 = p.speccy.ula.getBorderColor()
		} else if overflow[0].TState == TStatesPerFrame {
			colorAtTState0 = overflow[0].Color
		} else {
			// Use the Color of the last event that did NOT overflow.
			// Note: The fact that (numOverflow > 0) and (overflow[0].TState >= TStatesPerFrame) and
			// (there always exists an event with T-state value equal to 0)
			// implies that (numEvents > 0).
			colorAtTState0 = p.borderEvents[numEvents-1].Color
		}

		if (numOverflow > 0) && (overflow[0].TState == TStatesPerFrame) {
			p.borderEvents = p.borderEvents[0:0]
		} else {
			p.borderEvents = p.borderEvents[0:0]
			p.borderEvents = append(p.borderEvents, BorderEvent{TState: 0, Color: colorAtTState0})
		}

		// Replay the overflowing events
		for i := 0; i < numOverflow; i++ {
			p.borderEvents = append(p.borderEvents, BorderEvent{(overflow[i].TState - TStatesPerFrame), overflow[i].Color})
		}
	}

	// Beeper events
	{
		// Determine the number of events overflowing the frame
		var numOverflow int
		{
			i := len(p.beeperEvents)
			for (i > 0) && (p.beeperEvents[i-1].TState >= TStatesPerFrame) {
				i--
			}
			numOverflow = len(p.beeperEvents) - i
		}

		numEvents := len(p.beeperEvents) - numOverflow

		overflow := make([]BeeperEvent, numOverflow, numOverflow)
		copy(overflow[:], p.beeperEvents[numEvents:])

		var levelAtTState0 byte
		if numOverflow == 0 {
			levelAtTState0 = p.beeperLevel
		} else if overflow[0].TState == TStatesPerFrame {
			levelAtTState0 = overflow[0].Level
		} else {
			// Use the Level of the last event that did NOT overflow.
			// Note: The fact that (numOverflow > 0) and (overflow[0].TState >= TStatesPerFrame) and
			// (there always exists an event with T-state value equal to 0)
			// implies that (numEvents > 0).
			levelAtTState0 = p.beeperEvents[numEvents-1].Level
		}

		if (numOverflow > 0) && (overflow[0].TState == TStatesPerFrame) {
			p.beeperEvents = p.beeperEvents[0:0]
		} else {
			p.beeperEvents = p.beeperEvents[0:0]
			p.beeperEvents = append(p.beeperEvents, BeeperEvent{TState: 0, Level: levelAtTState0})
		}

		// Replay the overflowing events
		for i := 0; i < numOverflow; i++ {
			p.beeperEvents = append(p.beeperEvents, BeeperEvent{(overflow[i].TState - TStatesPerFrame), overflow[i].Level})
		}
	}

	return FrameStatusOfPorts{
		shouldPlayTheTape: (p.tapeReadCount >= tapeReadCount_tapeAccessThreshold),
	}
}

// Returns a copy of the list of border events.  The difference
// between [the T-state of the 1st event] and [the T-state of the last
// event] always equals to TStatesPerFrame (if the returned list is
// not empty).  
// 
// If the returned list is non-empty, its length is at least 2.
func (p *Ports) getBorderEvents() []BorderEvent {
	n := len(p.borderEvents)
	for (n > 0) && (p.borderEvents[n-1].TState > TStatesPerFrame) {
		n--
	}

	ret := make([]BorderEvent, n, n+1)
	copy(ret[0:n], p.borderEvents[0:n])

	if (n > 0) && (ret[n-1].TState < TStatesPerFrame) {
		ret = append(ret, BorderEvent{TStatesPerFrame, ret[n-1].Color})
	}

	return ret
}

// Returns a copy of the list of beeper events.  The difference
// between [the T-state of the 1st event] and [the T-state of the last
// event] always equals to TStatesPerFrame (if the returned list is
// not empty).
//
// If the returned list is non-empty, its length is at least 2.
func (p *Ports) getBeeperEvents() []BeeperEvent {
	n := len(p.beeperEvents)
	for (n > 0) && (p.beeperEvents[n-1].TState > TStatesPerFrame) {
		n--
	}

	ret := make([]BeeperEvent, n, n+1)
	copy(ret[0:n], p.beeperEvents[0:n])

	if (n > 0) && (ret[n-1].TState < TStatesPerFrame) {
		ret = append(ret, BeeperEvent{TStatesPerFrame, ret[n-1].Level})
	}

	return ret
}

func (p *Ports) ReadPort(address uint16) byte {
	return p.ReadPortInternal(address, true)
}

func (p *Ports) ReadPortInternal(address uint16, contend bool) byte {
	if contend {
		p.ContendPortPreio(address)
		p.ContendPortPostio(address)
	}

	var result byte = 0xff

	if (address & 0x0001) == 0x0000 {
		// Read keyboard
		var row uint
		for row = 0; row < 8; row++ {
			if (address & (1 << (uint16(row) + 8))) == 0 { // bit held low, so scan this row
				result &= p.speccy.Keyboard.GetKeyState(row)
			}
		}

		// Read tape
		if p.speccy.readFromTape && (address == 0x7ffe) {
			p.tapeReadCount++
			earBit := p.speccy.tapeDrive.getEarBit()
			result &= earBit
		}
	} else if (address & 0x00e0) == 0x0000 {
		result &= p.speccy.Joystick.GetState()
	} else {
		// Unassigned port
		result = 0xff
	}

	return result
}

func (p *Ports) WritePort(address uint16, b byte) {
	p.WritePortInternal(address, b, true)
}

func (p *Ports) WritePortInternal(address uint16, b byte, contend bool) {
	if contend {
		p.ContendPortPreio(address)
	}

	if (address & 0x0001) == 0 {
		color := (b & 0x07)

		// Modify the border only if it really changed
		if p.speccy.ula.getBorderColor() != color {
			p.speccy.ula.setBorderColor(color)

			last := len(p.borderEvents) - 1
			if p.borderEvents[last].TState == p.speccy.Cpu.Tstates {
				p.borderEvents[last].Color = color
			} else {
				p.borderEvents = append(p.borderEvents, BorderEvent{p.speccy.Cpu.Tstates, color})
			}
		}

		// EAR(bit 4) and MIC(bit 3) output
		newBeeperLevel := (b & 0x18) >> 3
		if p.speccy.readFromTape && !p.speccy.tapeDrive.AcceleratedLoad {
			if p.speccy.tapeDrive.earBit == 0xff {
				newBeeperLevel |= 2
			} else {
				newBeeperLevel &^= 2
			}
		}
		if p.beeperLevel != newBeeperLevel {
			p.beeperLevel = newBeeperLevel

			last := len(p.beeperEvents) - 1
			if p.beeperEvents[last].TState == p.speccy.Cpu.Tstates {
				p.beeperEvents[last].Level = newBeeperLevel
			} else {
				p.beeperEvents = append(p.beeperEvents, BeeperEvent{p.speccy.Cpu.Tstates, newBeeperLevel})
			}
		}
	}

	if contend {
		p.ContendPortPostio(address)
	}
}

func contendPort(z80 *z80.Z80, time int) {
	tstates_p := &z80.Tstates
	*tstates_p += int(delay_table[*tstates_p])
	*tstates_p += time
}

func (p *Ports) ContendPortPreio(address uint16) {
	if (address & 0xc000) == 0x4000 {
		contendPort(p.speccy.Cpu, 1)
	} else {
		p.speccy.Cpu.Tstates += 1
	}
}

func (p *Ports) ContendPortPostio(address uint16) {
	if (address & 0x0001) == 1 {
		if (address & 0xc000) == 0x4000 {
			contendPort(p.speccy.Cpu, 1)
			contendPort(p.speccy.Cpu, 1)
			contendPort(p.speccy.Cpu, 1)
		} else {
			p.speccy.Cpu.Tstates += 3
		}

	} else {
		contendPort(p.speccy.Cpu, 3)
	}
}
