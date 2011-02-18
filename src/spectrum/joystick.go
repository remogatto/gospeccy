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

import "sync"

const (
	KEMPSTON_FIRE = iota
	KEMPSTON_UP
	KEMPSTON_DOWN
	KEMPSTON_LEFT
	KEMPSTON_RIGHT
)

var kempstonMask = map[uint]byte{
	KEMPSTON_FIRE:  0x0010,
	KEMPSTON_UP:    0x0008,
	KEMPSTON_DOWN:  0x0004,
	KEMPSTON_LEFT:  0x0002,
	KEMPSTON_RIGHT: 0x0001,
}

type Joystick struct {
	speccy *Spectrum48k
	state  byte
	mutex  sync.RWMutex
}

func NewJoystick() *Joystick {
	joystick := &Joystick{}
	joystick.reset()
	return joystick
}

func (joystick *Joystick) init(speccy *Spectrum48k) {
	joystick.speccy = speccy
}

func (joystick *Joystick) reset() {
	joystick.SetState(0x0)
}

func (joystick *Joystick) GetState() byte {
	joystick.mutex.RLock()
	state := joystick.state
	joystick.mutex.RUnlock()
	return state
}

func (joystick *Joystick) SetState(state byte) {
	joystick.mutex.Lock()
	joystick.state = state
	joystick.mutex.Unlock()
}

func (joystick *Joystick) KempstonDown(logicalCode uint) {
	joystick.mutex.Lock()
	joystick.state |= kempstonMask[logicalCode]
	joystick.mutex.Unlock()
}

func (joystick *Joystick) KempstonUp(logicalCode uint) {
	joystick.mutex.Lock()
	joystick.state &= ^kempstonMask[logicalCode]
	joystick.mutex.Unlock()
}
