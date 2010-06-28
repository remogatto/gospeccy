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
}

type Port struct {
	Display DisplayAccessor
}

func (p *Port) readPort(address uint16) byte {
	var result byte = 0xff
	p.contendPortPreio(address)

	if (address & 0x0001) == 0x0000 {
		// Read keyboard
		for row := 0; row < 8; row++ {
			if (address & (1 << (uint16(row) + 8))) == 0 { // bit held low, so scan this row
				result &= keyStates[row]
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

func (p *Port) writePort(address uint16, b byte) {
	p.contendPortPreio(address)

	if ((address & 0x0001) == 0) {
		p.Display.setBorderColor(palette[b & 0x07])
	}

	p.contendPortPostio(address)
}

func (p *Port) contendPortPreio(address uint16) {
	if (address & 0xc000) == 0x4000 {
	}
	tstates++
}

func (p *Port) contendPortPostio(address uint16) {
	if (address & 0x0001) != 0 {
		if (address & 0xc000) == 0x4000 {
		} else {
			tstates += 3
		}

	} else {
		tstates += 3

	}

}
