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

type PortAccessor interface {
	readPort(address uint16) byte
	writePort(address uint16, b byte)
	readPortInternal(address uint16, contend bool) byte
	writePortInternal(address uint16, b byte, contend bool)
	contendPortPreio(address uint16)
	contendPortPostio(address uint16)

	reset()

	frame_begin()
	frame_end()

	// Returns a pointer to a linked-linked list of border events.
	// The difference between [the T-state of the 1st event] and [the T-state of the last event]
	// always equals to TStatesPerFrame.
	//
	// If the result is non-nil, the length of the returned list is at least 2.
	getBorderEvents_orNil() *BorderEvent

	// Returns a pointer to a linked-linked list of beeper events.
	// The difference between [the T-state of the 1st event] and [the T-state of the last event]
	// always equals to TStatesPerFrame.
	//
	// If the result is non-nil, the length of the returned list is at least 2.
	getBeeperEvents_orNil() *BeeperEvent
}


type BorderEvent struct {
	// The moment when the border color was changed.
	// It is the number of T-states since the beginning of the frame.
	tstate uint

	// The new border color
	color byte

	// Previous event, if any.
	// Constraint: (tstate >= previous_orNil.tstate)
	previous_orNil *BorderEvent
}

func (e *BorderEvent) TState() uint {
	return e.tstate
}

func (e *BorderEvent) Previous_orNil() Event {
	var ret Event
	if e.previous_orNil != nil {
		ret = e.previous_orNil
	}
	return ret
}

func (e1 *BorderEvent) Equals(e2 *BorderEvent) bool {
	var res bool

	if e1 == e2 {
		res = true
	} else if (e1 == nil) || (e2 == nil) {
		res = false
	} else {
		if (e1.tstate != e2.tstate) || (e1.color != e2.color) {
			res = false
		} else {
			res = e1.previous_orNil.Equals(e2.previous_orNil)
		}
	}

	return res
}


type BeeperEvent struct {
	// The moment when the beeper-event occurred.
	// It is the number of T-states since the beginning of the frame.
	tstate uint

	// The beeper level (0 .. MAX_AUDIO_LEVEL)
	level byte

	// Previous event, if any.
	// Constraint: (tstate >= previous_orNil.tstate)
	previous_orNil *BeeperEvent
}

func (e *BeeperEvent) TState() uint {
	return e.tstate
}

func (e *BeeperEvent) Previous_orNil() Event {
	var ret Event
	if e.previous_orNil != nil {
		ret = e.previous_orNil
	}
	return ret
}


type Ports struct {
	speccy *Spectrum48k

	borderEvents *BorderEvent // Might be nil

	beeperLevel  byte
	beeperEvents *BeeperEvent
}


func NewPorts() *Ports {
	p := &Ports{}
	p.beeperLevel = 0
	p.beeperEvents = &BeeperEvent{tstate: 0, level: p.beeperLevel, previous_orNil: nil}

	return p
}

func (p *Ports) init(speccy *Spectrum48k) {
	p.speccy = speccy
}

func (p *Ports) reset() {
	p.borderEvents = &BorderEvent{tstate: 0, color: p.speccy.ula.getBorderColor(), previous_orNil: nil}

	p.beeperLevel = 0
	p.beeperEvents = &BeeperEvent{tstate: 0, level: p.beeperLevel, previous_orNil: nil}
}


type Cmp_TStatesPerFrame struct{}

func (Cmp_TStatesPerFrame) isTrue(e Event) bool {
	return (e.TState() >= TStatesPerFrame)
}

type BeeperEvent_Array struct {
	events []*BeeperEvent
}

func (a *BeeperEvent_Array) Init(n int) {
	a.events = make([]*BeeperEvent, n)
}

func (a *BeeperEvent_Array) Set(i int, e Event) {
	a.events[i] = e.(*BeeperEvent)
}

type BorderEvent_Array struct {
	events []*BorderEvent
}

func (a *BorderEvent_Array) Init(n int) {
	a.events = make([]*BorderEvent, n)
}

func (a *BorderEvent_Array) Set(i int, e Event) {
	a.events[i] = e.(*BorderEvent)
}


func (p *Ports) frame_begin() {
	// Empty
}

func (p *Ports) frame_end() {
	// Border events
	{
		// Extract all events overflowing the frame
		overflow_array := &BorderEvent_Array{}
		EventListToArray_Ascending(p.borderEvents, overflow_array, Cmp_TStatesPerFrame{})

		overflow := overflow_array.events
		n := len(overflow)

		var colorAtTState0 byte
		if n == 0 {
			colorAtTState0 = p.speccy.ula.getBorderColor()
		} else if overflow[0].tstate == TStatesPerFrame {
			colorAtTState0 = overflow[0].color
		} else {
			// Note: The fact that (n>0) and (overflow[0].tstate >= TStatesPerFrame) and
			// (there always exists an event with T-state value equal to 0)
			// implies that (overflow[0].previous_orNil != nil).
			colorAtTState0 = overflow[0].previous_orNil.color
		}

		if (n > 0) && (overflow[0].tstate == TStatesPerFrame) {
			p.borderEvents = nil
		} else {
			p.borderEvents = &BorderEvent{tstate: 0, color: colorAtTState0, previous_orNil: nil}
		}

		// Replay the overflowing events
		for i := 0; i < n; i++ {
			p.borderEvents = &BorderEvent{(overflow[i].tstate - TStatesPerFrame), overflow[i].color, p.borderEvents}
		}
	}

	// Beeper events
	{
		// Extract all events overflowing the frame
		overflow_array := &BeeperEvent_Array{}
		EventListToArray_Ascending(p.beeperEvents, overflow_array, Cmp_TStatesPerFrame{})

		overflow := overflow_array.events
		n := len(overflow)

		var levelAtTState0 byte
		if n == 0 {
			levelAtTState0 = p.beeperLevel
		} else if overflow[0].tstate == TStatesPerFrame {
			levelAtTState0 = overflow[0].level
		} else {
			// Note: The fact that (n>0) and (overflow[0].tstate >= TStatesPerFrame) and
			// (there always exists an event with T-state value equal to 0)
			// implies that (overflow[0].previous_orNil != nil).
			levelAtTState0 = overflow[0].previous_orNil.level
		}

		if (n > 0) && (overflow[0].tstate == TStatesPerFrame) {
			p.beeperEvents = nil
		} else {
			p.beeperEvents = &BeeperEvent{tstate: 0, level: levelAtTState0, previous_orNil: nil}
		}

		// Replay the overflowing events
		for i := 0; i < n; i++ {
			p.beeperEvents = &BeeperEvent{(overflow[i].tstate - TStatesPerFrame), overflow[i].level, p.beeperEvents}
		}
	}
}

func (p *Ports) getBorderEvents_orNil() *BorderEvent {
	lastEvent := p.borderEvents

	for lastEvent.tstate > TStatesPerFrame {
		lastEvent = lastEvent.previous_orNil
	}

	if lastEvent.tstate < TStatesPerFrame {
		lastEvent = &BorderEvent{TStatesPerFrame, lastEvent.color, lastEvent}
	}

	return lastEvent
}

func (p *Ports) getBeeperEvents_orNil() *BeeperEvent {
	lastEvent := p.beeperEvents

	for lastEvent.tstate > TStatesPerFrame {
		lastEvent = lastEvent.previous_orNil
	}

	if lastEvent.tstate < TStatesPerFrame {
		lastEvent = &BeeperEvent{TStatesPerFrame, lastEvent.level, lastEvent}
	}

	return lastEvent
}

func (p *Ports) readPort(address uint16) byte {
	return p.readPortInternal(address, true)
}

func (p *Ports) readPortInternal(address uint16, contend bool) byte {
	if contend {
		p.contendPortPreio(address)
		p.contendPortPostio(address)
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
		if p.speccy.Cpu.readFromTape {
			earBit := p.speccy.TapeDrive.getEarBit()
			return result & earBit
		}
	} else if (address & 0x00e0) == 0x0000 {
		// Kempston joystick: treat this as attached but
		// unused (for the benefit of Manic Miner)
		result = 0x00
	} else {
		// Unassigned port
		result = 0xff
	}

	return result
}

func (p *Ports) writePort(address uint16, b byte) {
	p.writePortInternal(address, b, true)
}

func (p *Ports) writePortInternal(address uint16, b byte, contend bool) {
	if contend {
		p.contendPortPreio(address)
	}

	if (address & 0x0001) == 0 {
		color := (b & 0x07)

		// Modify the border only if it really changed
		if p.speccy.ula.getBorderColor() != color {
			p.speccy.ula.setBorderColor(color)
			if p.borderEvents.tstate == p.speccy.Cpu.tstates {
				p.borderEvents.color = color
			} else {
				p.borderEvents = &BorderEvent{p.speccy.Cpu.tstates, color, p.borderEvents}
			}
		}

		// EAR(bit 4) and MIC(bit 3) output
		newBeeperLevel := (b & 0x18) >> 3
		if p.speccy.Cpu.readFromTape {
			if p.speccy.TapeDrive.earBit == 0xff {
				newBeeperLevel = 3
			} else {
				newBeeperLevel = 1
			}
		}
		if p.beeperLevel != newBeeperLevel {
			p.beeperLevel = newBeeperLevel
			if p.beeperEvents.tstate == p.speccy.Cpu.tstates {
				p.beeperEvents.level = newBeeperLevel
			} else {
				p.beeperEvents = &BeeperEvent{p.speccy.Cpu.tstates, newBeeperLevel, p.beeperEvents}
			}
		}
	}

	if contend {
		p.contendPortPostio(address)
	}
}

func (p *Ports) contend(time uint) {
	tstates_p := &p.speccy.Cpu.tstates
	*tstates_p += uint(delay_table[*tstates_p])
	*tstates_p += time
}

func (p *Ports) contendPortPreio(address uint16) {
	if (address & 0xc000) == 0x4000 {
		p.contend(1)
	} else {
		p.speccy.Cpu.tstates += 1
	}
}

func (p *Ports) contendPortPostio(address uint16) {
	if (address & 0x0001) == 1 {
		if (address & 0xc000) == 0x4000 {
			p.contend(1)
			p.contend(1)
			p.contend(1)
		} else {
			p.speccy.Cpu.tstates += 3
		}

	} else {
		p.contend(3)
	}
}
