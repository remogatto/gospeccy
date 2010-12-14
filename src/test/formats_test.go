package test

import (
	"prettytest"
	"spectrum/formats"
	"testing"
)

type formats_suite_t struct {
	test_suite_t
}

func (s *formats_suite_t) should_support_SNA_format() {
	program, err := formats.ReadProgram("testdata/hello.sna")
	s.Nil(err)

	err = speccy.Load(program)
	s.Nil(err)

	s.True(screenEqualTo("testdata/hello_tape_loaded.sna"))
}

func (s *formats_suite_t) should_support_Z80_format() {
	program, err := formats.ReadProgram("testdata/hello.z80")
	s.Nil(err)

	err = speccy.Load(program)
	s.Nil(err)

	s.True(screenEqualTo("testdata/hello_tape_loaded.sna"))
}

func (s *formats_suite_t) should_support_TAP_format() {
	program, err := formats.ReadProgram("testdata/hello.tap")
	s.Nil(err)

	err = speccy.Load(program)
	s.Nil(err)

	<-speccy.TapeDrive().LoadComplete()

	s.True(screenEqualTo("testdata/hello_tape_loaded.sna"))
}

func TestLoadFormats(t *testing.T) {
	prettytest.RunWithFormatter(
		t,
		&prettytest.BDDFormatter{"The formats"},
		new(formats_suite_t),
	)
}
