package test

import (
	"time"
)

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
