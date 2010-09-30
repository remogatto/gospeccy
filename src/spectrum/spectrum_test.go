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
	os.Remove("testdata/screen.scr")
}

func testMakeVideoMemoryDump(t *prettytest.T) {
	speccy.makeVideoMemoryDump(scrFilename)
	fileInfo, err := os.Stat(scrFilename)

	t.True(err == nil)
	t.Equal(int64(6912), fileInfo.Size)
}

func testMakeVideoMemoryDumpCmd(t *prettytest.T) {
	errChan := make(chan os.Error)
	speccy.CommandChannel <- Cmd_MakeVideoMemoryDump{ scrFilename, errChan }

	<-errChan

	fileInfo, err := os.Stat(scrFilename)
	
	t.True(err == nil)
	t.Equal(int64(6912), fileInfo.Size)
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


