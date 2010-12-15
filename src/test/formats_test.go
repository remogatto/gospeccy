package test

import (
	"testing"
	pt "spectrum/prettytest"
	"spectrum/formats"
)

func should_support_SNA_format(t *pt.T) {
	program, err := formats.ReadProgram("testdata/hello.sna")
	t.Nil(err)

	err = speccy.Load(program)
	t.Nil(err)

	t.True(screenEqualTo("testdata/hello_tape_loaded.sna"))
}

func should_support_Z80_format(t *pt.T) {
	program, err := formats.ReadProgram("testdata/hello.z80")
	t.Nil(err)

	err = speccy.Load(program)
	t.Nil(err)

	t.True(screenEqualTo("testdata/hello_tape_loaded.sna"))
}

func should_support_TAP_format(t *pt.T) {
	program, err := formats.ReadProgram("testdata/hello.tap")
	t.Nil(err)

	err = speccy.Load(program)
	t.Nil(err)

	<-speccy.TapeDrive().LoadComplete()

	t.True(screenEqualTo("testdata/hello_tape_loaded.sna"))
}

func TestLoadFormats(t *testing.T) {
	pt.Describe(
		t,
		"The emulator",
		should_support_SNA_format,
		should_support_Z80_format,
		should_support_TAP_format,

		before,
		after,
	)
}
