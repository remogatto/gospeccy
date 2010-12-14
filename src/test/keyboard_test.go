package test

import (
	"prettytest"
	"spectrum"
	"testing"
)

type keyboard_suite_t struct {
	test_suite_t
}

func (s *keyboard_suite_t) should_respond_to_keypress() {
	<-speccy.Keyboard.KeyPress(spectrum.KEY_R)
	<-speccy.Keyboard.KeyPress(spectrum.KEY_Enter)

	s.True(screenEqualTo("testdata/key_press_1_ok.sna"))
}

func (s *keyboard_suite_t) should_respond_to_keypress_sequence() {
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
	for i := 0; i < 5; i++ {
		<-done
	}

	k.KeyDown(spectrum.KEY_SymbolShift)
	<-speccy.Keyboard.KeyPress(spectrum.KEY_P)
	k.KeyUp(spectrum.KEY_SymbolShift)

	<-speccy.Keyboard.KeyPress(spectrum.KEY_Enter)
	s.True(screenEqualTo("testdata/key_press_sequence_1_ok.sna"))
}

func TestKeyboard(t *testing.T) {
	prettytest.RunWithFormatter(
		t,
		&prettytest.BDDFormatter{"The keyboard"},
		new(keyboard_suite_t),
	)
}
