package spectrum

import (
	"testing"
	"spectrum/prettytest"
	"os"
)

var speccy *Spectrum48k

func before(t *prettytest.Assertions) *prettytest.Assertions {
	var err os.Error
	app := NewApplication()
	speccy, err = NewSpectrum48k(app, "testdata/48.rom")
	if err != nil {
		panic(err)
	}
	return t
}

func testSaveScreen(t *prettytest.Assertions) *prettytest.Assertions {
	return t.Pending()
}

func TestSaveScreen(t *testing.T) {
	prettytest.Run(
		t,
		before,
		testSaveScreen,
	)
}


