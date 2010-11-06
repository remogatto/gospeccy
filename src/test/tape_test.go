package test

import (
	"testing"
	pt "spectrum/prettytest"
)

func should_load_tapes_using_ROM_routine(t *pt.T) {
	err := speccy.LoadTape("testdata/hello.tap")
	t.Nil(err)

	<-speccy.TapeDrive().LoadComplete()

	t.True(screenEqualTo("testdata/hello_tape_loaded.sna"))
}

func TestTapeFeatures(t *testing.T) {
	pt.Describe(
		t,
		"The emulator",
		should_load_tapes_using_ROM_routine,

		before,
		after,
	)
}
