package test

import (
	"github.com/remogatto/gospeccy/src/formats"
	"github.com/remogatto/gospeccy/src/spectrum"
	"github.com/remogatto/prettytest"
	"io/ioutil"
	"testing"
	"time"
)

// System ROM

func (t *testSuite) Should_load_system_ROM() {
	romLoaded := make(chan (<-chan bool))
	speccy.CommandChannel <- spectrum.Cmd_Reset{romLoaded}
	<-(<-romLoaded)

	t.True(screenEqualTo("testdata/system_rom_loaded.sna"))
}

// Keyboard

func (t *testSuite) Should_respond_to_keypress() {
	<-speccy.Keyboard.KeyPress(spectrum.KEY_R)
	<-speccy.Keyboard.KeyPress(spectrum.KEY_Enter)
	t.True(screenEqualTo("testdata/key_press_1_ok.sna"))
}

func (t *testSuite) Should_respond_to_keypress_sequence() {
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

func (t *testSuite) Should_load_tapes_using_ROM_routine() {
	filename := "testdata/hello.tap"
	data, err := ioutil.ReadFile(filename)
	t.Nil(err)
	tap, err := formats.NewTAP(data)
	t.Nil(err)

	// Reset
	romLoaded := make(chan (<-chan bool))
	speccy.CommandChannel <- spectrum.Cmd_Reset{romLoaded}
	<-(<-romLoaded)

	errChan := make(chan error)
	speccy.CommandChannel <- spectrum.Cmd_Load{ /*informalFileName*/ filename, tap, errChan}
	t.Nil(<-errChan)

	<-speccy.TapeDrive().LoadComplete()

	t.True(screenEqualTo("testdata/hello_tape_loaded.sna"))
}

func (t *testSuite) Should_support_accelerated_loading() {
	filename := "testdata/hello.tap"
	data, err := ioutil.ReadFile(filename)
	t.Nil(err)
	tap, err := formats.NewTAP(data)
	t.Nil(err)

	// Reset
	romLoaded := make(chan (<-chan bool))
	speccy.CommandChannel <- spectrum.Cmd_Reset{romLoaded}
	<-(<-romLoaded)

	start := time.Now()
	speccy.TapeDrive().AcceleratedLoad = true

	errChan := make(chan error)
	speccy.CommandChannel <- spectrum.Cmd_Load{ /*informalFileName*/ filename, tap, errChan}
	t.Nil(<-errChan)

	<-speccy.TapeDrive().LoadComplete()

	t.True(time.Now().Sub(start).Nanoseconds() < 10e9)
	t.True(screenEqualTo("testdata/hello_tape_loaded.sna"))
}

// // Formats

func (t *testSuite) Should_support_SNA_format() {
	filename := "testdata/hello.sna"
	snapshot, err := formats.ReadProgram(filename)
	t.Nil(err)

	errChan := make(chan error)
	speccy.CommandChannel <- spectrum.Cmd_Load{ /*informalFileName*/ filename, snapshot, errChan}
	t.Nil(<-errChan)

	t.True(screenEqualTo("testdata/hello_tape_loaded.sna"))
}

func (t *testSuite) Should_support_Z80_format() {
	filename := "testdata/hello.z80"
	snapshot, err := formats.ReadProgram(filename)
	t.Nil(err)

	errChan := make(chan error)
	speccy.CommandChannel <- spectrum.Cmd_Load{ /*informalFileName*/ filename, snapshot, errChan}
	t.Nil(<-errChan)

	t.True(screenEqualTo("testdata/hello_tape_loaded.sna"))
}

func (t *testSuite) Should_support_TAP_format() {
	filename := "testdata/hello.tap"
	data, err := ioutil.ReadFile(filename)
	t.Nil(err)
	tap, err := formats.NewTAP(data)
	t.Nil(err)

	// Reset
	romLoaded := make(chan (<-chan bool))
	speccy.CommandChannel <- spectrum.Cmd_Reset{romLoaded}
	<-(<-romLoaded)

	errChan := make(chan error)
	speccy.CommandChannel <- spectrum.Cmd_Load{ /*informalFileName*/ filename, tap, errChan}
	t.Nil(<-errChan)

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
