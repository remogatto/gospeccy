package spectrum

import (
	"testing"
	"spectrum/prettytest"
)

const scrFilename = "testdata/screen.scr"

func beforeAllSaveScreenTests(t *prettytest.T) {
	StartFullEmulation()
}

func afterAllSaveScreenTests(t *prettytest.T) {
	app.RequestExit()
	<-app.HasTerminated
}

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
		beforeAllSaveScreenTests,
		testMakeVideoMemoryDump,
		testMakeVideoMemoryDumpCmd,
		afterAllSaveScreenTests,
	)
}


