package test

import (
	"container/vector"
	"fmt"
	"os"
	"prettytest"
	"regexp"
	"spectrum"
	"spectrum/console"
	"testing"
	"time"
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

type console_suite_t struct {
	prettytest.Suite
}

func (s *console_suite_t) beforeAll() {
	println("func (s *console_suite_t) beforeAll()")
	StartFullEmulation()

	console.IgnoreStartupScript = true
	console.Init(app, speccy)
}

func (s *console_suite_t) afterAll() {
	println("func (s *console_suite_t) after()")
	app.RequestExit()
	<-app.HasTerminated
}

func (s *console_suite_t) before() {
	println("func (s *console_suite_t) before()")
	app.SetMessageOutput(&testMessageOutput{new(vector.StringVector)})
	messageOutput = app.GetMessageOutput().(*testMessageOutput)
}

func (s *console_suite_t) should_allow_loading_tapes_using_ROM_routine() {
	err := console.RunString(fmt.Sprintf("load(\"%s\")", "testdata/hello.tap"))
	s.Nil(err)

	<-speccy.TapeDrive().LoadComplete()

	s.True(screenEqualTo("testdata/hello_tape_loaded.sna"))
}

func (s *console_suite_t) should_allow_accelerated_tape_load() {
	err := console.RunString("acceleratedLoad(true)")
	s.Nil(err)

	start := time.Nanoseconds()
	err = console.RunString(fmt.Sprintf("load(\"%s\")", "testdata/hello.tap"))
	s.Nil(err)

	<-speccy.TapeDrive().LoadComplete()

	err = console.RunString("acceleratedLoad(false)")

	s.True((time.Nanoseconds() - start) < 10e9)
	s.False(speccy.TapeDrive().AcceleratedLoad)
	s.True(screenEqualTo("testdata/hello_tape_loaded.sna"))
}

func (s *console_suite_t) should_honor_convention_over_configuration_when_loading_files() {
	spectrum.DefaultUserDir = "testdata/"

	err := console.RunString(fmt.Sprintf("load(\"%s\")", "hello.tap"))
	s.Nil(err)
	s.Equal(0, messageOutput.strings.Len())

	<-speccy.TapeDrive().LoadComplete()
	s.True(screenEqualTo("testdata/hello_tape_loaded.sna"))

	err = console.RunString(fmt.Sprintf("load(\"%s\")", "hello.zip"))
	s.Nil(err)
	s.Equal(0, messageOutput.strings.Len())
	s.True(screenEqualTo("testdata/hello_tape_loaded.sna"))

	err = console.RunString(fmt.Sprintf("load(\"%s\")", "hello.sna"))
	s.Nil(err)
	s.Equal(0, messageOutput.strings.Len())
	s.True(screenEqualTo("testdata/hello_tape_loaded.sna"))
}

func (s *console_suite_t) should_allow_reset() {
	err := console.RunString("reset()")
	s.Nil(err)

	<-speccy.ROMLoaded()

	s.True(screenEqualTo("testdata/system_rom_loaded.sna"))
}

func (s *console_suite_t) should_print_help() {
	err := console.RunString("help()")
	matched, _ := regexp.MatchString("Available commands:\n *.* *\n", messageOutput.String())

	s.Nil(err)
	s.True(matched)
}

func (s *console_suite_t) should_allow_printing_strings() {
	err := console.RunString(fmt.Sprintf("puts(\"%s\")", "Hello World!"))

	s.Nil(err)
	s.Equal("Hello World!\n", messageOutput.String())
}

func (s *console_suite_t) should_allow_taking_screenshots() {
	defer os.Remove("testdata/screenshot.scr")

	err := console.RunString(fmt.Sprintf("screenshot(\"%s\")", "testdata/screenshot.scr"))

	s.Nil(err)
	s.Equal(0, messageOutput.strings.Len())
	s.Path("testdata/screenshot.scr")
}

func (s *console_suite_t) should_allow_loading_scripts() {
	err := console.RunString(fmt.Sprintf("script(\"%s\")", "testdata/script"))

	s.True(err == nil)
	s.Equal(2, messageOutput.strings.Len())
	s.Equal("Hello World!\n", messageOutput.String())
}

func (s *console_suite_t) should_allow_keypress_events() {
	s.Pending()
}

func TestConsoleFeatures(t *testing.T) {
	prettytest.RunWithFormatter(
		t,
		&prettytest.BDDFormatter{"The console"},
		new(console_suite_t),
	)
}
