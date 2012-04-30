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

import (
	"sync"
	"time"
)

type rowState struct {
	row, state byte
}

type Cmd_KeyPress struct {
	logicalKeyCode uint
	done           chan bool
}

type Cmd_SendLoad struct {
	romType RomType
}

type Keyboard struct {
	speccy    *Spectrum48k
	keyStates [8]byte
	mutex     sync.RWMutex

	CommandChannel chan interface{}
}

func NewKeyboard() *Keyboard {
	keyboard := &Keyboard{}
	keyboard.reset()

	keyboard.CommandChannel = make(chan interface{})

	return keyboard
}

func (keyboard *Keyboard) init(speccy *Spectrum48k) {
	keyboard.speccy = speccy
	go keyboard.commandLoop()
}

func (keyboard *Keyboard) delayAfterKeyDown() {
	// Sleep for 1 frame
	time.Sleep(1e9 / time.Duration(keyboard.speccy.GetCurrentFPS()))
}

func (keyboard *Keyboard) delayAfterKeyUp() {
	// Sleep for 10 frames
	time.Sleep(10 * 1e9 / time.Duration(keyboard.speccy.GetCurrentFPS()))
}

func (keyboard *Keyboard) commandLoop() {
	evtLoop := keyboard.speccy.app.NewEventLoop()
	for {
		select {

		case <-evtLoop.Pause:
			evtLoop.Pause <- 0

		case <-evtLoop.Terminate:
			// Terminate this Go routine
			if evtLoop.App().Verbose {
				evtLoop.App().PrintfMsg("keyboard command loop: exit")
			}
			evtLoop.Terminate <- 0
			return

		case untyped_cmd := <-keyboard.CommandChannel:
			switch cmd := untyped_cmd.(type) {
			case Cmd_KeyPress:
				keyboard.KeyDown(cmd.logicalKeyCode)
				keyboard.delayAfterKeyDown()
				keyboard.KeyUp(cmd.logicalKeyCode)
				keyboard.delayAfterKeyUp()
				cmd.done <- true

			case Cmd_SendLoad:
				if cmd.romType == ROM_OPENSE {
					// Sleep for 30 frames
					time.Sleep(30e9 / time.Duration(keyboard.speccy.GetCurrentFPS()))

					// l o a d
					for _, keycode := range []uint{KEY_L, KEY_O, KEY_A, KEY_D} {
						keyboard.KeyDown(keycode)
						keyboard.delayAfterKeyDown()
						keyboard.KeyUp(keycode)
						keyboard.delayAfterKeyUp()
					}

					// " "
					keyboard.KeyDown(KEY_SymbolShift)
					{
						keyboard.KeyDown(KEY_P)
						keyboard.delayAfterKeyDown()
						keyboard.KeyUp(KEY_P)
						keyboard.delayAfterKeyUp()

						keyboard.KeyDown(KEY_P)
						keyboard.delayAfterKeyDown()
						keyboard.KeyUp(KEY_P)
						keyboard.delayAfterKeyUp()
					}
					keyboard.KeyUp(KEY_SymbolShift)

					keyboard.KeyDown(KEY_Enter)
					keyboard.delayAfterKeyDown()
					keyboard.KeyUp(KEY_Enter)
				} else {
					// LOAD
					keyboard.KeyDown(KEY_J)
					keyboard.delayAfterKeyDown()
					keyboard.KeyUp(KEY_J)
					keyboard.delayAfterKeyUp()

					// " "
					keyboard.KeyDown(KEY_SymbolShift)
					{
						keyboard.KeyDown(KEY_P)
						keyboard.delayAfterKeyDown()
						keyboard.KeyUp(KEY_P)
						keyboard.delayAfterKeyUp()

						keyboard.KeyDown(KEY_P)
						keyboard.delayAfterKeyDown()
						keyboard.KeyUp(KEY_P)
						keyboard.delayAfterKeyUp()
					}
					keyboard.KeyUp(KEY_SymbolShift)

					keyboard.KeyDown(KEY_Enter)
					keyboard.delayAfterKeyDown()
					keyboard.KeyUp(KEY_Enter)
				}
			}
		}
	}

}

func (k *Keyboard) reset() {
	// Initialize 'k.keyStates'
	for row := uint(0); row < 8; row++ {
		k.SetKeyState(row, 0xff)
	}
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

func (keyboard *Keyboard) KeyDown(logicalKeyCode uint) {
	keyCode, ok := keyCodes[logicalKeyCode]

	if ok {
		keyboard.mutex.Lock()
		keyboard.keyStates[keyCode.row] &= ^(keyCode.mask)
		keyboard.mutex.Unlock()
	}
}

func (keyboard *Keyboard) KeyUp(logicalKeyCode uint) {
	keyCode, ok := keyCodes[logicalKeyCode]

	if ok {
		keyboard.mutex.Lock()
		keyboard.keyStates[keyCode.row] |= (keyCode.mask)
		keyboard.mutex.Unlock()
	}
}

func (keyboard *Keyboard) KeyPress(logicalKeyCode uint) chan bool {
	done := make(chan bool)
	keyboard.CommandChannel <- Cmd_KeyPress{logicalKeyCode, done}
	return done
}

func (keyboard *Keyboard) KeyPressSequence(logicalKeyCodes ...uint) chan bool {
	done := make(chan bool, len(logicalKeyCodes))
	for _, keyCode := range logicalKeyCodes {
		keyboard.CommandChannel <- Cmd_KeyPress{keyCode, done}
	}
	return done
}

// Logical key codes
const (
	KEY_1 = iota
	KEY_2
	KEY_3
	KEY_4
	KEY_5
	KEY_6
	KEY_7
	KEY_8
	KEY_9
	KEY_0

	KEY_Q
	KEY_W
	KEY_E
	KEY_R
	KEY_T
	KEY_Y
	KEY_U
	KEY_I
	KEY_O
	KEY_P

	KEY_A
	KEY_S
	KEY_D
	KEY_F
	KEY_G
	KEY_H
	KEY_J
	KEY_K
	KEY_L
	KEY_Enter

	KEY_CapsShift
	KEY_Z
	KEY_X
	KEY_C
	KEY_V
	KEY_B
	KEY_N
	KEY_M
	KEY_SymbolShift
	KEY_Space
)

type keyCell struct {
	row, mask byte
}

var keyCodes = map[uint]keyCell{
	KEY_1: keyCell{row: 3, mask: 0x01},
	KEY_2: keyCell{row: 3, mask: 0x02},
	KEY_3: keyCell{row: 3, mask: 0x04},
	KEY_4: keyCell{row: 3, mask: 0x08},
	KEY_5: keyCell{row: 3, mask: 0x10},
	KEY_6: keyCell{row: 4, mask: 0x10},
	KEY_7: keyCell{row: 4, mask: 0x08},
	KEY_8: keyCell{row: 4, mask: 0x04},
	KEY_9: keyCell{row: 4, mask: 0x02},
	KEY_0: keyCell{row: 4, mask: 0x01},

	KEY_Q: keyCell{row: 2, mask: 0x01},
	KEY_W: keyCell{row: 2, mask: 0x02},
	KEY_E: keyCell{row: 2, mask: 0x04},
	KEY_R: keyCell{row: 2, mask: 0x08},
	KEY_T: keyCell{row: 2, mask: 0x10},
	KEY_Y: keyCell{row: 5, mask: 0x10},
	KEY_U: keyCell{row: 5, mask: 0x08},
	KEY_I: keyCell{row: 5, mask: 0x04},
	KEY_O: keyCell{row: 5, mask: 0x02},
	KEY_P: keyCell{row: 5, mask: 0x01},

	KEY_A:     keyCell{row: 1, mask: 0x01},
	KEY_S:     keyCell{row: 1, mask: 0x02},
	KEY_D:     keyCell{row: 1, mask: 0x04},
	KEY_F:     keyCell{row: 1, mask: 0x08},
	KEY_G:     keyCell{row: 1, mask: 0x10},
	KEY_H:     keyCell{row: 6, mask: 0x10},
	KEY_J:     keyCell{row: 6, mask: 0x08},
	KEY_K:     keyCell{row: 6, mask: 0x04},
	KEY_L:     keyCell{row: 6, mask: 0x02},
	KEY_Enter: keyCell{row: 6, mask: 0x01},

	KEY_CapsShift:   keyCell{row: 0, mask: 0x01},
	KEY_Z:           keyCell{row: 0, mask: 0x02},
	KEY_X:           keyCell{row: 0, mask: 0x04},
	KEY_C:           keyCell{row: 0, mask: 0x08},
	KEY_V:           keyCell{row: 0, mask: 0x10},
	KEY_B:           keyCell{row: 7, mask: 0x10},
	KEY_N:           keyCell{row: 7, mask: 0x08},
	KEY_M:           keyCell{row: 7, mask: 0x04},
	KEY_SymbolShift: keyCell{row: 7, mask: 0x02},
	KEY_Space:       keyCell{row: 7, mask: 0x01},
}

var SDL_KeyMap = map[string][]uint{
	"0": []uint{KEY_0},
	"1": []uint{KEY_1},
	"2": []uint{KEY_2},
	"3": []uint{KEY_3},
	"4": []uint{KEY_4},
	"5": []uint{KEY_5},
	"6": []uint{KEY_6},
	"7": []uint{KEY_7},
	"8": []uint{KEY_8},
	"9": []uint{KEY_9},

	"a": []uint{KEY_A},
	"b": []uint{KEY_B},
	"c": []uint{KEY_C},
	"d": []uint{KEY_D},
	"e": []uint{KEY_E},
	"f": []uint{KEY_F},
	"g": []uint{KEY_G},
	"h": []uint{KEY_H},
	"i": []uint{KEY_I},
	"j": []uint{KEY_J},
	"k": []uint{KEY_K},
	"l": []uint{KEY_L},
	"m": []uint{KEY_M},
	"n": []uint{KEY_N},
	"o": []uint{KEY_O},
	"p": []uint{KEY_P},
	"q": []uint{KEY_Q},
	"r": []uint{KEY_R},
	"s": []uint{KEY_S},
	"t": []uint{KEY_T},
	"u": []uint{KEY_U},
	"v": []uint{KEY_V},
	"w": []uint{KEY_W},
	"x": []uint{KEY_X},
	"y": []uint{KEY_Y},
	"z": []uint{KEY_Z},

	"return":      []uint{KEY_Enter},
	"space":       []uint{KEY_Space},
	"left shift":  []uint{KEY_CapsShift},
	"right shift": []uint{KEY_CapsShift},
	"left ctrl":   []uint{KEY_SymbolShift},
	"right ctrl":  []uint{KEY_SymbolShift},

	//"escape":    []uint{KEY_CapsShift, KEY_1},
	//"caps lock": []uint{KEY_CapsShift, KEY_2}, // FIXME: SDL never sends the sdl.KEYUP event
	"left":      []uint{KEY_CapsShift, KEY_5},
	"down":      []uint{KEY_CapsShift, KEY_6},
	"up":        []uint{KEY_CapsShift, KEY_7},
	"right":     []uint{KEY_CapsShift, KEY_8},
	"backspace": []uint{KEY_CapsShift, KEY_0},

	"-": []uint{KEY_SymbolShift, KEY_J},
	//"_": []uint{KEY_SymbolShift, KEY_0},
	"=": []uint{KEY_SymbolShift, KEY_L},
	//"+": []uint{KEY_SymbolShift, KEY_K},
	"[": []uint{KEY_SymbolShift, KEY_8}, // Maps to "("
	"]": []uint{KEY_SymbolShift, KEY_9}, // Maps to ")"
	";": []uint{KEY_SymbolShift, KEY_O},
	//":": []uint{KEY_SymbolShift, KEY_Z},
	"'": []uint{KEY_SymbolShift, KEY_7},
	//"\"": []uint{KEY_SymbolShift, KEY_P},
	",": []uint{KEY_SymbolShift, KEY_N},
	".": []uint{KEY_SymbolShift, KEY_M},
	"/": []uint{KEY_SymbolShift, KEY_V},
	//"<": []uint{KEY_SymbolShift, KEY_R},
	//">": []uint{KEY_SymbolShift, KEY_T},
	//"?": []uint{KEY_SymbolShift, KEY_C},

	// Keypad
	"[0]": []uint{KEY_0},
	"[1]": []uint{KEY_1},
	"[2]": []uint{KEY_2},
	"[3]": []uint{KEY_3},
	"[4]": []uint{KEY_4},
	"[5]": []uint{KEY_5},
	"[6]": []uint{KEY_6},
	"[7]": []uint{KEY_7},
	"[8]": []uint{KEY_8},
	"[9]": []uint{KEY_9},
	"[*]": []uint{KEY_SymbolShift, KEY_B},
	"[-]": []uint{KEY_SymbolShift, KEY_J},
	"[+]": []uint{KEY_SymbolShift, KEY_K},
	"[/]": []uint{KEY_SymbolShift, KEY_V},
}

func init() {
	if len(keyCodes) != 40 {
		panic("invalid keyboard specification")
	}

	// Make sure we are able to press every button on the Spectrum keyboard
	used := make(map[uint]bool)
	for logicalKeyCode := range keyCodes {
		used[logicalKeyCode] = false
	}
	for _, seq := range SDL_KeyMap {
		if len(seq) == 1 {
			used[seq[0]] = true
		}
	}
	for _, isUsed := range used {
		if !isUsed {
			panic("some key is missing in the SDL keymap")
		}
	}
}
