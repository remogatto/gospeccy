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
	contendPortPreio(address uint16)
	contendPortPostio(address uint16)

	reset()

	frame_begin()
	frame_releaseMemory()

	// This method may return nil
	getBorderEvents() *BorderEvent

	// This method may return nil
	getBeeperEvents() *BeeperEvent
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

type BeeperEvent struct {
	// The moment when the beeper-event occurred.
	// It is the number of T-states since the beginning of the frame.
	tstate uint

	// The beeper level (0 or 1)
	level byte

	// Previous event, if any.
	// Constraint: (tstate >= previous_orNil.tstate)
	previous_orNil *BeeperEvent
}

type Ports struct {
	speccy *Spectrum48k

	borderEvents *BorderEvent // Might be nil

	beeperLevel  byte
	beeperEvents *BeeperEvent // Might be nil
}


func NewPorts() *Ports {
	return &Ports{}
}

func (p *Ports) init(speccy *Spectrum48k) {
	p.speccy = speccy
}

func (p *Ports) reset() {
	p.borderEvents = nil
	p.beeperLevel = 0
	p.beeperEvents = nil
}


func (p *Ports) frame_begin() {
	borderColor := p.speccy.ula.getBorderColor()
	p.borderEvents = &BorderEvent{tstate: 0, color: borderColor, previous_orNil: nil}
	p.beeperEvents = &BeeperEvent{tstate: 0, level: p.beeperLevel, previous_orNil: nil}
}

func (p *Ports) frame_releaseMemory() {
	// Release memory
	p.borderEvents = nil
	p.beeperEvents = nil
}

func (p *Ports) getBorderEvents() *BorderEvent {
	return p.borderEvents
}

func (p *Ports) getBeeperEvents() *BeeperEvent {
	return p.beeperEvents
}


func (p *Ports) readPort(address uint16) byte {
	var result byte = 0xff
	p.contendPortPreio(address)

	if (address & 0x0001) == 0x0000 {
		// Read keyboard
		var row uint
		for row = 0; row < 8; row++ {
			if (address & (1 << (uint16(row) + 8))) == 0 { // bit held low, so scan this row
				result &= p.speccy.Keyboard.GetKeyState(row)
			}
		}
		return result
	} else if (address & 0x00e0) == 0x0000 {
		// Kempston joystick: treat this as attached but
		// unused (for the benefit of Manic Miner)
		return 0x00
	} else {
		return 0xff // Unassigned port
	}

	p.contendPortPostio(address)

	return result
}

func (p *Ports) writePort(address uint16, b byte) {
	p.contendPortPreio(address)

	if (address & 0x0001) == 0 {
		color := (b & 0x07)

		// Modify the border only if it really changed
		if p.speccy.ula.getBorderColor() != color {
			p.speccy.ula.setBorderColor(color)
			p.borderEvents = &BorderEvent{p.speccy.Cpu.tstates, color, p.borderEvents}
		}

		// EAR output
		newBeeperLevel := (b & 0x10) >> 4
		if p.beeperLevel != newBeeperLevel {
			p.beeperLevel = newBeeperLevel
			p.beeperEvents = &BeeperEvent{p.speccy.Cpu.tstates, newBeeperLevel, p.beeperEvents}
		}
	}

	p.contendPortPostio(address)
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
