package formats

import (
	"prettytest"
	"testing"
)

type testSuite struct {
	prettytest.Suite
}

func (t *testSuite) testReadProgramSnapshot() {
	program, err := ReadProgram("testdata/fire.sna")
	_, ok := program.(Snapshot)

	t.Nil(err)
	t.True(ok)
}

func (t *testSuite) testReadProgramTape() {
	program, err := ReadProgram("testdata/fire.tap")
	_, ok := program.(*TAP)

	t.Nil(err)
	t.True(ok)
}

func (t *testSuite) testReadProgramZIP() {
	program, err := ReadProgram("testdata/fire.sna.zip")
	_, ok := program.(Snapshot)

	t.Nil(err)
	t.True(ok)
}

func (t *testSuite) testReadSnapshot() {
	_, err := ReadProgram("testdata/fire.sna")
	t.Nil(err)
}

func TestFormats(t *testing.T) {
	prettytest.Run(t, new(testSuite))
}
