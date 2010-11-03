package spectrum

import (
	"testing"
	pt "spectrum/prettytest"
	"os"
	"time"
)

func keyPress(logicalKeyCode uint) {
	speccy.Keyboard.KeyDown(logicalKeyCode)
	time.Sleep(int64(1e9/speccy.GetCurrentFPS()))
	speccy.Keyboard.KeyUp(logicalKeyCode)
}

// Send LOAD ""
func SendLoadBASICCommand() {
	keyPress(KEY_J)
	speccy.Keyboard.KeyDown(KEY_SymbolShift)
	keyPress(KEY_P)
	time.Sleep(int64(1e9))
	keyPress(KEY_P)
	speccy.Keyboard.KeyUp(KEY_SymbolShift)
	keyPress(KEY_Enter)
}

// Test Tape

var (
	tape *Tape
	err os.Error
)

func beforeTape(assert *pt.T) {
	tape, err = NewTapeFromFile("testdata/hello.tap")
}

func testNewTape(assert *pt.T) {
	assert.Nil(err)
	assert.NotNil(tape)
}

func testTapeAt(assert *pt.T) {
	assert.Equal(byte(0x00), tape.At(0))
	assert.Equal(byte(0xff), tape.At(0x13))
}

func TestTape(t *testing.T) {
	pt.Run(
		t,
		beforeTape,
		testNewTape,
		testTapeAt,
	)
}

// Test TapeDrive

func beforeTapeDrive(assert *pt.T) {
	tape, _ = NewTapeFromFile("testdata/hello.tap")
}

func testNewTapeDrive(assert *pt.T) {
	assert.NotNil(speccy.TapeDrive)
}

func testInsertTape(assert *pt.T) {
	speccy.TapeDrive.Insert(tape)
	assert.NotNil(speccy.TapeDrive.tape)
}

func testTapeDrivePlayStop(assert *pt.T) {
	speccy.TapeDrive.Play()
	speccy.TapeDrive.doPlay()
	speccy.TapeDrive.Stop()
	assert.Equal(byte(TAPE_DRIVE_PAUSE_STOP), speccy.TapeDrive.state)
	assert.Equal(0, speccy.TapeDrive.currBlockId)
}

func testTapeDriveLoad(assert *pt.T) {
	speccy.TapeDrive.Insert(tape)
	speccy.TapeDrive.Play()

	// SEND LOAD ""
	SendLoadBASICCommand()

	ok := <-speccy.TapeDrive.loadComplete
	assert.True(ok)

	speccy.TapeDrive.Stop()
	// Compare snapshots
}

func testTapeDriveLoadWithCustomLoader(assert *pt.T) {
	romLoaded = false
	speccy.reset()
	<-romLoadedCh

	tape, _ = NewTapeFromFile("testdata/Syntax09nF.tap")
	speccy.TapeDrive.Insert(tape)
	speccy.TapeDrive.Play()

	// SEND LOAD ""
	SendLoadBASICCommand()

	ok := <-speccy.TapeDrive.loadComplete
	assert.True(ok)

	// Compare snapshots
}

func TestTapeDrive(t *testing.T) {
	pt.Run(
		t,
		beforeAll,
		afterAll,
		beforeTapeDrive,

		testNewTapeDrive,
		testInsertTape,
		testTapeDrivePlayStop,
		testTapeDriveLoad,
		testTapeDriveLoadWithCustomLoader,
	)
}
