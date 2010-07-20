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

type Keyboard struct {
	keyStates [8]byte
	mutex     sync.RWMutex
}

func NewKeyboard() *Keyboard {
	k := &Keyboard{}

	// Initialize keyStates
	var row uint
	for row = 0; row < 8; row++ {
		k.SetKeyState(row, 0xff)
	}

	return k
}

func (keyboard *Keyboard) GetKeyState(row uint) byte {
	keyboard.mutex.RLock()
	keyState := keyboard.keyStates[row]
	keyboard.mutex.RUnlock()
	return keyState
}

func (keyboard *Keyboard) SetKeyState(row uint, state byte) {
	keyboard.mutex.Lock()
	keyboard.keyStates[row] = state
	keyboard.mutex.Unlock()
}

func (keyboard *Keyboard) KeyDown(keySym uint) {
	keyCode, ok := keyCodes[keySym]

	if ok {
		keyboard.mutex.Lock()
		keyboard.keyStates[keyCode.row] &= ^(keyCode.mask)
		keyboard.mutex.Unlock()
	}
}

func (keyboard *Keyboard) KeyUp(keySym uint) {
	keyCode, ok := keyCodes[keySym]

	if ok {
		keyboard.mutex.Lock()
		keyboard.keyStates[keyCode.row] |= (keyCode.mask)
		keyboard.mutex.Unlock()
	}
}

type keyCell struct {
	row, mask byte
}

var keyCodes = map[uint]keyCell{
	49: keyCell{row: 3, mask: 0x01}, /* 1 */
	50: keyCell{row: 3, mask: 0x02}, /* 2 */
	51: keyCell{row: 3, mask: 0x04}, /* 3 */
	52: keyCell{row: 3, mask: 0x08}, /* 4 */
	53: keyCell{row: 3, mask: 0x10}, /* 5 */
	54: keyCell{row: 4, mask: 0x10}, /* 6 */
	55: keyCell{row: 4, mask: 0x08}, /* 7 */
	56: keyCell{row: 4, mask: 0x04}, /* 8 */
	57: keyCell{row: 4, mask: 0x02}, /* 9 */
	48: keyCell{row: 4, mask: 0x01}, /* 0 */

	113: keyCell{row: 2, mask: 0x01}, /* Q */
	119: keyCell{row: 2, mask: 0x02}, /* W */
	101: keyCell{row: 2, mask: 0x04}, /* E */
	114: keyCell{row: 2, mask: 0x08}, /* R */
	116: keyCell{row: 2, mask: 0x10}, /* T */
	121: keyCell{row: 5, mask: 0x10}, /* Y */
	117: keyCell{row: 5, mask: 0x08}, /* U */
	105: keyCell{row: 5, mask: 0x04}, /* I */
	111: keyCell{row: 5, mask: 0x02}, /* O */
	112: keyCell{row: 5, mask: 0x01}, /* P */

	97:  keyCell{row: 1, mask: 0x01}, /* A */
	115: keyCell{row: 1, mask: 0x02}, /* S */
	100: keyCell{row: 1, mask: 0x04}, /* D */
	102: keyCell{row: 1, mask: 0x08}, /* F */
	103: keyCell{row: 1, mask: 0x10}, /* G */
	104: keyCell{row: 6, mask: 0x10}, /* H */
	106: keyCell{row: 6, mask: 0x08}, /* J */
	107: keyCell{row: 6, mask: 0x04}, /* K */
	108: keyCell{row: 6, mask: 0x02}, /* L */
	13:  keyCell{row: 6, mask: 0x01}, /* enter */

	304: keyCell{row: 0, mask: 0x01}, /* caps */
	122: keyCell{row: 0, mask: 0x02}, /* Z */
	120: keyCell{row: 0, mask: 0x04}, /* X */
	99:  keyCell{row: 0, mask: 0x08}, /* C */
	118: keyCell{row: 0, mask: 0x10}, /* V */
	98:  keyCell{row: 7, mask: 0x10}, /* B */
	110: keyCell{row: 7, mask: 0x08}, /* N */
	109: keyCell{row: 7, mask: 0x04}, /* M */
	306: keyCell{row: 7, mask: 0x02}, /* sym - gah, firefox screws up ctrl+key too */
	32:  keyCell{row: 7, mask: 0x01} /* space */}
