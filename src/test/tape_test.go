package test

import (
	"prettytest"
	"testing"
	"time"
)

type tape_test_t struct {
	test_suite_t
}

func (s *tape_test_t) should_load_tapes_using_ROM_routine() {
	err := speccy.LoadTape("testdata/hello.tap")
	s.Nil(err)

	<-speccy.TapeDrive().LoadComplete()

	s.True(screenEqualTo("testdata/hello_tape_loaded.sna"))
}

func (s *tape_test_t) should_support_accelerated_loading() {
	start := time.Nanoseconds()
	speccy.TapeDrive().AcceleratedLoad = true
	err := speccy.LoadTape("testdata/hello.tap")
	s.Nil(err)

	<-speccy.TapeDrive().LoadComplete()

	s.True((time.Nanoseconds() - start) < 10e9)
	s.True(screenEqualTo("testdata/hello_tape_loaded.sna"))
}

func TestTapeFeatures(t *testing.T) {
	prettytest.RunWithFormatter(
		t,
		&prettytest.BDDFormatter{"The tape"},
		new(tape_test_t),
	)
}
