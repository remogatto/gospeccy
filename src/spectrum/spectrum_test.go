package spectrum

import (
	"testing"
	"os"
	"spectrum/prettytest"
)

const scrFilename = "testdata/screen.scr"
var speccy *Spectrum48k

func before(t *prettytest.T) {
	var err os.Error
	app := NewApplication()
	speccy, err = NewSpectrum48k(app, "testdata/48.rom")
	if err != nil {
		panic(err)
	}
}

func after(t *prettytest.T) {
	speccy = nil
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
		before,
		after,
		testMakeVideoMemoryDump,
		testMakeVideoMemoryDumpCmd,
	)
}


