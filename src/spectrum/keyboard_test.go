package spectrum

import (
	"testing"
	"spectrum/prettytest"
)

func testCmd_KeyPress(t *prettytest.T) {
	done := make(chan bool)
	speccy.Keyboard.CommandChannel <- Cmd_KeyPress{ KEY_R, done }
	<-done
	speccy.Keyboard.CommandChannel <- Cmd_KeyPress{ KEY_Enter, done }
	<-done
	t.True(screenEqualTo("testdata/key_press_1_ok.sna"))
}

func testKeyPress(t *prettytest.T) {
	<-speccy.Keyboard.KeyPress(KEY_R)
	<-speccy.Keyboard.KeyPress(KEY_Enter)
	t.True(screenEqualTo("testdata/key_press_1_ok.sna"))
}

func testKeyPressSequence(t *prettytest.T) {
	k := speccy.Keyboard
	<-speccy.Keyboard.KeyPress(KEY_P)

	k.KeyDown(KEY_SymbolShift)
	<-speccy.Keyboard.KeyPress(KEY_P)
	k.KeyUp(KEY_SymbolShift)

	done := k.KeyPressSequence(KEY_H, KEY_E, KEY_L, KEY_L, KEY_O)
	for i := 0; i < 5; i++ { <-done }

	k.KeyDown(KEY_SymbolShift)
	<-speccy.Keyboard.KeyPress(KEY_P)
	k.KeyUp(KEY_SymbolShift)

	<-speccy.Keyboard.KeyPress(KEY_Enter)

	t.True(screenEqualTo("testdata/key_press_sequence_1_ok.sna"))
}

func TestKeyboard(t *testing.T) {
	prettytest.Run(
		t,
		before,
		after,
		testCmd_KeyPress,
		testKeyPress,
		testKeyPressSequence,
	)
}
