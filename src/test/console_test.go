package test

import (
	"fmt"
	"github.com/remogatto/gospeccy/src/spectrum"
	pt "github.com/remogatto/prettytest"
	"os"
	"regexp"
	"testing"
	"time"
)

func (t *cliTestSuite) Should_allow_loading_tapes_using_ROM_routine() {
	console.PutCommand(fmt.Sprintf("load(\"%s\")", "testdata/hello.tap"))
	<-speccy.TapeDrive().LoadComplete()
	t.True(screenEqualTo("testdata/hello_tape_loaded.sna"))
}

func (t *cliTestSuite) Should_allow_accelerated_tape_load() {
	console.PutCommand("acceleratedLoad(true)")

	start := time.Now()
	console.PutCommand(fmt.Sprintf("load(\"%s\")", "testdata/hello.tap"))

	<-speccy.TapeDrive().LoadComplete()

	console.PutCommand("acceleratedLoad(false)")

	t.True(time.Since(start).Nanoseconds() < 10e9)
	t.Not(t.True(speccy.TapeDrive().AcceleratedLoad))
	t.True(screenEqualTo("testdata/hello_tape_loaded.sna"))
}

func (t *cliTestSuite) Should_honor_convention_over_configuration_when_loading_tap() {
	spectrum.DefaultUserDir = "testdata/"
	console.PutCommand(fmt.Sprintf("load(\"%s\")", "hello.tap"))
	<-speccy.TapeDrive().LoadComplete()
	t.True(screenEqualTo("testdata/hello_tape_loaded.sna"))
}

func (t *cliTestSuite) Should_honor_convention_over_configuration_when_loading_zip() {
	spectrum.DefaultUserDir = "testdata/"
	console.PutCommand(fmt.Sprintf("load(\"%s\")", "hello.zip"))
	t.True(screenEqualTo("testdata/hello_tape_loaded.sna"))
}

func (t *cliTestSuite) Should_honor_convention_over_configuration_when_loading_sna() {
	spectrum.DefaultUserDir = "testdata/"
	console.PutCommand(fmt.Sprintf("load(\"%s\")", "hello.sna"))
	t.True(screenEqualTo("testdata/hello_tape_loaded.sna"))
}

func (t *cliTestSuite) Should_reset_speccy() {
	console.PutCommand("reset()")
	t.True(screenEqualTo("testdata/system_rom_loaded.sna"))
}

func (t *cliTestSuite) Should_print_help() {
	console.PutCommand("help()")
	matched, _ := regexp.MatchString("Available commands:\n *.* *\n", console.String())
	t.True(matched)
}

func (t *cliTestSuite) Should_allow_printing_strings() {
	console.PutCommand(fmt.Sprintf("puts(\"%s\")", "Hello World!"))
	r, _ := regexp.Compile("Hello World!")
	t.True(len(r.FindAllString(console.String(), -1)) > 1)
}

func (t *cliTestSuite) Should_allow_taking_screenshots() {
	defer os.Remove("testdata/screenshot.scr")
	console.PutCommand(fmt.Sprintf("screenshot(\"%s\")", "testdata/screenshot.scr"))
	t.Path("testdata/screenshot.scr")
}

func (t *cliTestSuite) Should_allow_loading_scripts() {
	spectrum.DefaultUserDir = "testdata/"
	console.PutCommand(fmt.Sprintf("script(\"%s\")", "script"))
	r, _ := regexp.Compile("Hello World!")
	t.Equal(1, len(r.FindAllString(console.String(), -1)))
}

func (t *cliTestSuite) Should_allow_keypress_events() {
	t.Pending()
}

func TestConsoleFeatures(t *testing.T) {
	pt.RunWithFormatter(
		t,
		&pt.BDDFormatter{Description: "The CLI"},
		new(cliTestSuite),
	)
}
