package test

import (
	"spectrum/formats"
)

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
