package test

import (
	"testing"
	"spectrum/prettytest"
	"spectrum"
)

func should_respond_to_keypress(t *prettytest.T) {
	<-speccy.Keyboard.KeyPress(spectrum.KEY_R)
	<-speccy.Keyboard.KeyPress(spectrum.KEY_Enter)

	t.True(screenEqualTo("testdata/key_press_1_ok.sna"))
}

func should_respond_to_keypress_sequence(t *prettytest.T) {
	k := speccy.Keyboard
	<-speccy.Keyboard.KeyPress(spectrum.KEY_P)

	k.KeyDown(spectrum.KEY_SymbolShift)
	<-speccy.Keyboard.KeyPress(spectrum.KEY_P)
	k.KeyUp(spectrum.KEY_SymbolShift)

	done := k.KeyPressSequence(
		spectrum.KEY_H, 
		spectrum.KEY_E,
		spectrum.KEY_L,
		spectrum.KEY_L,
		spectrum.KEY_O,
	)
	for i := 0; i < 5; i++ { <-done }

	k.KeyDown(spectrum.KEY_SymbolShift)
	<-speccy.Keyboard.KeyPress(spectrum.KEY_P)
	k.KeyUp(spectrum.KEY_SymbolShift)

	<-speccy.Keyboard.KeyPress(spectrum.KEY_Enter)
	t.True(screenEqualTo("testdata/key_press_sequence_1_ok.sna"))
}

func TestKeyboard(t *testing.T) {
	prettytest.Describe(
		t,
		"The keyboard",
		should_respond_to_keypress,
		should_respond_to_keypress_sequence,

		before,
		after,

	)
}
