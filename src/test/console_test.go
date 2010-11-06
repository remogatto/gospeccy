package test

import (
	"testing"
	"fmt"
	"os"
	"regexp"
	"container/vector"
	pt "spectrum/prettytest"
	"spectrum"
	"spectrum/console"
)

type testMessageOutput struct {
	strings *vector.StringVector
}

var messageOutput *testMessageOutput

func (out *testMessageOutput) Clear() {
	out.strings = new(vector.StringVector)
}

func (out *testMessageOutput) String() string {
	outputString := ""
	for _, str := range *out.strings {
		outputString += str
	}
	return outputString
}

func (out *testMessageOutput) PrintfMsg(format string, a ...interface{}) {
	out.strings.Push(fmt.Sprintf(format, a...))

	appendNewLine := false
	if (len(format) == 0) || (format[len(format)-1] != '\n') {
		appendNewLine = true
	}

	if appendNewLine {
		out.strings.Push("\n")
	}

}

func beforeAllConsole(t *pt.T) {
	beforeAll(t)
	console.Init(app, speccy)
}

func beforeConsole(t *pt.T) {
	app.SetMessageOutput(&testMessageOutput{ new(vector.StringVector) })
	messageOutput = app.GetMessageOutput().(*testMessageOutput)
}

func should_allow_loading_tapes_using_ROM_routine(t *pt.T) {
	err := console.RunString(fmt.Sprintf("load(\"%s\")", "testdata/hello.tap"))
	t.Nil(err)

	<-speccy.TapeDrive().LoadComplete()

	t.True(screenEqualTo("testdata/hello_tape_loaded.sna"))
}

func should_honor_convention_over_configuration_when_loading_files(t *pt.T) {
	spectrum.DefaultUserDir = "testdata/"

	err := console.RunString(fmt.Sprintf("load(\"%s\")", "hello.tap"))
	t.Nil(err)
	t.Equal(0, messageOutput.strings.Len())

	<-speccy.TapeDrive().LoadComplete()
	t.True(screenEqualTo("testdata/hello_tape_loaded.sna"))

	err = console.RunString(fmt.Sprintf("load(\"%s\")", "hello.zip"))
	t.Nil(err)
	t.Equal(0, messageOutput.strings.Len())
	t.True(screenEqualTo("testdata/hello_tape_loaded.sna"))

	err = console.RunString(fmt.Sprintf("load(\"%s\")", "hello.sna"))
	t.Nil(err)
	t.Equal(0, messageOutput.strings.Len())
	t.True(screenEqualTo("testdata/hello_tape_loaded.sna"))
}

func should_allow_reset(t *pt.T) {
	err := console.RunString("reset()")
	t.Nil(err)
	
	<-speccy.ROMLoaded()

	t.True(screenEqualTo("testdata/system_rom_loaded.sna"))
}

func should_print_help(t *pt.T) {
	err := console.RunString("help()")
	matched, _ := regexp.MatchString("Available commands:\n *.* *\n", messageOutput.String())

	t.Nil(err)
	t.True(matched)
}

func should_allow_printing_strings(t *pt.T) {
	err := console.RunString(fmt.Sprintf("puts(\"%s\")", "Hello World!"))
	
	t.Nil(err)
	t.Equal("Hello World!\n", messageOutput.String())
}

func should_allow_taking_screenshots(t *pt.T) {
	defer os.Remove("testdata/screenshot.scr")

	err := console.RunString(fmt.Sprintf("screenshot(\"%s\")", "testdata/screenshot.scr"))

	t.Nil(err)
	t.Equal(0, messageOutput.strings.Len())
	t.Path("testdata/screenshot.scr")
}

func should_allow_loading_scripts(t *pt.T) {
	err := console.RunString(fmt.Sprintf("script(\"%s\")", "testdata/script"))
	
	t.True(err == nil)
	t.Equal(2, messageOutput.strings.Len())
	t.Equal("Hello World!\n", messageOutput.String())
}

func should_allow_keypress_events(t *pt.T) {
	t.Pending()
}

func TestConsoleFeatures(t *testing.T) {
	pt.Describe(
		t,
		"The console",
		should_allow_loading_tapes_using_ROM_routine,
		should_honor_convention_over_configuration_when_loading_files,
		should_allow_reset,
		should_print_help,
		should_allow_printing_strings,
		should_allow_taking_screenshots,
		should_allow_loading_scripts,
		should_allow_keypress_events,

		beforeAllConsole,
		beforeConsole,
		afterAll,
	)
}
