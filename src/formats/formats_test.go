package formats

import (
	"github.com/remogatto/prettytest"
	"strings"
	"testing"
)

type testSuite struct {
	prettytest.Suite
}

func (t *testSuite) TestReadProgram_SNA() {
	program, err := ReadProgram("testdata/fire.sna")
	_, ok := program.(Snapshot)

	t.Nil(err)
	t.True(ok)
}

func (t *testSuite) TestReadProgram_Z80() {
	program, err := ReadProgram("testdata/fire.z80")
	_, ok := program.(Snapshot)

	t.Nil(err)
	t.True(ok)
}

func (t *testSuite) TestReadProgram_TAP() {
	program, err := ReadProgram("testdata/fire.tap")
	_, ok := program.(*TAP)

	t.Nil(err)
	t.True(ok)
}

func (t *testSuite) TestReadProgram_SNA_ZIP() {
	program, err := ReadProgram("testdata/fire.sna.zip")
	_, ok := program.(Snapshot)

	t.Nil(err)
	t.True(ok)
}

func (t *testSuite) TestReadProgram_ZIP_ambiguous() {
	_, err := ReadProgram("testdata/ambiguous.zip")
	t.NotNil(err)
	t.True(strings.Contains(err.Error(), "multiple"))
}

func TestFormats(t *testing.T) {
	prettytest.Run(t, new(testSuite))
}
