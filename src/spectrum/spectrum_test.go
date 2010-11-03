package spectrum

import (
	"testing"
	"spectrum/formats"
	"spectrum/prettytest"
	"os"
)

const scrFilename = "testdata/screen.scr"

func testMakeVideoMemoryDump(t *prettytest.T) {
	t.Equal(6912, len(speccy.makeVideoMemoryDump()))
}

func testMakeVideoMemoryDumpCmd(t *prettytest.T) {
	ch := make(chan []byte)
	speccy.CommandChannel <- Cmd_MakeVideoMemoryDump{ ch }

	data := <-ch

	t.Equal(6912, len(data))
}

func TestSaveScreen(t *testing.T) {
	prettytest.Run(
		t,
		beforeAll,
		afterAll,
		testMakeVideoMemoryDump,
		testMakeVideoMemoryDumpCmd,
	)
}

func testLoadTape(t *prettytest.T) {
	tap, _ := formats.NewTAPFromFile("testdata/hello.tap")
	err := speccy.Load(tap)
	t.Nil(err)

	<-speccy.TapeDrive.loadComplete

	t.True(stateEqualTo("testdata/hello_tap_ok.sna"))
}

func testLoadTapeCmd(t *prettytest.T) {
	errCh := make(chan os.Error)
	program, _ := formats.NewTAPFromFile("testdata/hello.tap")
	speccy.CommandChannel <- Cmd_Load{ ErrChan: errCh, Program: program }
	t.Nil(<-errCh)

	<-speccy.TapeDrive.loadComplete

	t.True(stateEqualTo("testdata/hello_tap_ok.sna"))
}

func testLoadSnapshot(t *prettytest.T) {
	snapshot, _ := formats.ReadSnapshot("testdata/hello_tap_ok.sna")
	err := speccy.Load(snapshot)
	t.Nil(err)
	t.True(stateEqualTo("testdata/hello_tap_ok.sna"))
}

func testLoadSnapshotCmd(t *prettytest.T) {
	errCh := make(chan os.Error)
	program, _ := formats.ReadSnapshot("testdata/hello_tap_ok.sna")
	speccy.CommandChannel <- Cmd_Load{ ErrChan: errCh, Program: program }
	t.Nil(<-errCh)
	t.True(stateEqualTo("testdata/hello_tap_ok.sna"))
}

func TestLoad(t *testing.T) {
	prettytest.Run(
		t,
		before,
		after,
		testLoadTape,
		testLoadTapeCmd,
		testLoadSnapshot,
		testLoadSnapshotCmd,
	)
}
