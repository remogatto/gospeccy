package formats

import (
	"testing"
	pt "spectrum/prettytest"
)

func testReadProgramSnapshot(t *pt.T) {
	program, err := ReadProgram("testdata/fire.sna")
	_, ok := program.(Snapshot)

	t.Nil(err)
	t.True(ok)
}

func testReadProgramTape(t *pt.T) {
	program, err := ReadProgram("testdata/fire.tap")
	_, ok := program.(*TAP)

	t.Nil(err)
	t.True(ok)
}

func testReadProgramZIP(t *pt.T) {
	program, err := ReadProgram("testdata/fire.sna.zip")
	_, ok := program.(Snapshot)

	t.Nil(err)
	t.True(ok)
}

func TestReadProgram(t *testing.T) {
	pt.Run(
		t,
		testReadProgramSnapshot,
		testReadProgramTape,
		testReadProgramZIP,
	)
}

func testReadSnapshot(t *pt.T) {
	_, err := ReadSnapshot("testdata/fire.sna")
	t.Nil(err)
}

func TestReadSnapshot(t *testing.T) {
	pt.Run(
		t,
		testReadSnapshot,
	)
}

