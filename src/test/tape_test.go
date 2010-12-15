package test

import (
	"testing"
	pt "spectrum/prettytest"
	"time"
)

func should_load_tapes_using_ROM_routine(t *pt.T) {
	err := speccy.LoadTape("testdata/hello.tap")
	t.Nil(err)

	<-speccy.TapeDrive().LoadComplete()

	t.True(screenEqualTo("testdata/hello_tape_loaded.sna"))
}

func should_support_accelerated_loading(t *pt.T) {
	start := time.Nanoseconds()
	speccy.TapeDrive().AcceleratedLoad = true
	err := speccy.LoadTape("testdata/hello.tap")
	t.Nil(err)

	<-speccy.TapeDrive().LoadComplete()

	t.True((time.Nanoseconds() - start) < 10e9)
	t.True(screenEqualTo("testdata/hello_tape_loaded.sna"))
}

func TestTapeFeatures(t *testing.T) {
	pt.Describe(
		t,
		"The emulator",
		//		should_load_tapes_using_ROM_routine,
		should_support_accelerated_loading,

		before,
		after,
	)
}
