package test

import (
	"time"
	"spectrum"
	"spectrum/formats"
	"testing"
	"prettytest"
)

// System ROM

func (t *testSuite) should_load_system_ROM() {
	t.True(screenEqualTo("testdata/system_rom_loaded.sna"))
}

// Keyboard

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

// // Tapedrive

func (t *testSuite) should_load_tapes_using_ROM_routine() {
	err := speccy.LoadTape("testdata/hello.tap")
	t.Nil(err)

	<-speccy.TapeDrive().LoadComplete()

	t.True(screenEqualTo("testdata/hello_tape_loaded.sna"))
}

func (t *testSuite) should_support_accelerated_loading() {
	start := time.Nanoseconds()
	speccy.TapeDrive().AcceleratedLoad = true
	err := speccy.LoadTape("testdata/hello.tap")
	t.Nil(err)

	<-speccy.TapeDrive().LoadComplete()

	t.True((time.Nanoseconds() - start) < 10e9)
	t.True(screenEqualTo("testdata/hello_tape_loaded.sna"))
}

// // Formats

func (t *testSuite) should_support_SNA_format() {
	program, err := formats.ReadProgram("testdata/hello.sna")
	t.Nil(err)
	
	err = speccy.Load(program)
	t.Nil(err)

	t.True(screenEqualTo("testdata/hello_tape_loaded.sna"))
}

func (t *testSuite) should_support_Z80_format() {
	program, err := formats.ReadProgram("testdata/hello.z80")
	t.Nil(err)

	err = speccy.Load(program)
	t.Nil(err)

	t.True(screenEqualTo("testdata/hello_tape_loaded.sna"))
}

func (t *testSuite) should_support_TAP_format() {
	program, err := formats.ReadProgram("testdata/hello.tap")
	t.Nil(err)

	err = speccy.Load(program)
	t.Nil(err)

	<-speccy.TapeDrive().LoadComplete()

	t.True(screenEqualTo("testdata/hello_tape_loaded.sna"))
}

func TestEmulator(t *testing.T) {
	prettytest.RunWithFormatter(
		t,
		&prettytest.BDDFormatter{"The emulator"},
		new(testSuite),
	)
}
