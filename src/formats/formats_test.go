package formats

import (
	"prettytest"
	"testing"
)

type formats_suite_t struct {
	prettytest.Suite
}

func (s *formats_suite_t) testReadProgramSnapshot() {
	program, err := ReadProgram("testdata/fire.sna")
	_, ok := program.(Snapshot)

	s.Nil(err)
	s.True(ok)
}

func (s *formats_suite_t) testReadProgramTape() {
	program, err := ReadProgram("testdata/fire.tap")
	_, ok := program.(*TAP)

	s.Nil(err)
	s.True(ok)
}

func (s *formats_suite_t) testReadProgramZIP() {
	program, err := ReadProgram("testdata/fire.sna.zip")
	_, ok := program.(Snapshot)

	s.Nil(err)
	s.True(ok)
}

func (s *formats_suite_t) testReadSnapshot() {
	_, err := ReadProgram("testdata/fire.sna")
	s.Nil(err)
}

func TestRead(t *testing.T) {
	prettytest.Run(
		t,
		new(formats_suite_t),
	)
}

