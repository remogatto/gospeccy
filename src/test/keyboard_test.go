package test

import (
	"spectrum"
)

func (t *testSuite) should_respond_to_keypress() {
	<-speccy.Keyboard.KeyPress(spectrum.KEY_R)
	<-speccy.Keyboard.KeyPress(spectrum.KEY_Enter)

	t.True(screenEqualTo("testdata/key_press_1_ok.sna"))
}

func (t *testSuite) should_respond_to_keypress_sequence() {
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
	t.True(screenEqualTo("testdata/key_press_sequence_1_ok.sna"))
}
