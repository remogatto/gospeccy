package console

import (
	"spectrum"
	"spectrum/formats"
	"spectrum/readline"
	"exp/eval"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"strings"
	"container/vector"
	"sync"
	"bytes"
	"time"
)

// ==============
// Some variables
// ==============

// These variables are set only once, before starting new goroutines,
// so there is no need for controlling concurrent access via a sync.Mutex
var app *spectrum.Application
var speccy *spectrum.Spectrum48k
var w *eval.World

const PROMPT = "gospeccy> "
const PROMPT_EMPTY = "          "

// Whether the terminal is currently showing a prompt string
var havePrompt = false
var havePrompt_mutex sync.Mutex

var ignoreStartupScript = false

const SCRIPT_DIRECTORY = "scripts"
const STARTUP_SCRIPT = "startup"

// ================
// Various commands
// ================

var help_keys vector.StringVector
var help_vals vector.StringVector

func printHelp() {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "\nAvailable commands:\n")

	maxKeyLen := 1
	for i := 0; i < help_keys.Len(); i++ {
		if len(help_keys[i]) > maxKeyLen {
			maxKeyLen = len(help_keys[i])
		}
	}

	for i := 0; i < help_keys.Len(); i++ {
		fmt.Fprintf(&buf, "    %s", help_keys[i])
		for j := len(help_keys[i]); j < maxKeyLen; j++ {
			fmt.Fprintf(&buf, " ")
		}
		fmt.Fprintf(&buf, "  %s\n", help_vals[i])
	}

	app.PrintfMsg("%s\n", buf.String())
}

// Signature: func help()
func wrapper_help(t *eval.Thread, in []eval.Value, out []eval.Value) {
	printHelp()
}

// Signature: func exit()
func wrapper_exit(t *eval.Thread, in []eval.Value, out []eval.Value) {
	// Implementation note:
	//   The following test has to be there only in cases in which something can go wrong.
	//   For example if the user would try to execute "exit(); sound(false)" then GoSpeccy would panic.
	//   An alternative way would be to actually terminate the whole program at the 1st statement - so that
	//   "sound(false)" or whatever is not executed - alas this is somewhat problematic,
	//   since once the script "exit(); sound(false)" runs, it cannot be stopped halfway
	//   through its execution. Using "runtime.Goexit()" would solve this issue, but only partially,
	//   since it is potentially possible for the statement "sound(false)" to be hidden in a defer statement.
	//   So, the best option (until somebody implements a better one) is to convert the problematic commands
	//   into statements that are doing nothing while the application is in the process of being exited.
	if app.TerminationInProgress() {
		return
	}

	app.RequestExit()
}

// Signature: func reset()
func wrapper_reset(t *eval.Thread, in []eval.Value, out []eval.Value) {
	if app.TerminationInProgress() {
		return
	}

	speccy.CommandChannel <- spectrum.Cmd_Reset{}
}

// Signature: func addSearchPath(path string)
func wrapper_addSearchPath(t *eval.Thread, in []eval.Value, out []eval.Value) {
	path := in[0].(eval.StringValue).Get(t)
	spectrum.AddCustomSearchPath(path)
}

// Signature: func load(path string)
func wrapper_load(t *eval.Thread, in []eval.Value, out []eval.Value) {
	if app.TerminationInProgress() {
		return
	}

	path := in[0].(eval.StringValue).Get(t)

	snapshot, err := formats.ReadSnapshot(spectrum.SnaPath(path))
	if err != nil {
		app.PrintfMsg("%s", err)
		return
	}

	errChan := make(chan os.Error)
	speccy.CommandChannel <- spectrum.Cmd_LoadSnapshot{path, snapshot, errChan}
	err = <-errChan
	if err != nil {
		app.PrintfMsg("%s", err)
	}
}

// Signature: func save(path string)
func wrapper_save(t *eval.Thread, in []eval.Value, out []eval.Value) {
	if app.TerminationInProgress() {
		return
	}

	path := in[0].(eval.StringValue).Get(t)

	ch := make(chan *formats.FullSnapshot)
	speccy.CommandChannel <- spectrum.Cmd_MakeSnapshot{ch}

	fullSnapshot := <-ch

	data, err := fullSnapshot.EncodeSNA()
	if err != nil {
		app.PrintfMsg("%s", err)
		return
	}

	err = ioutil.WriteFile(path, data, 0600)
	if err != nil {
		app.PrintfMsg("%s", err)
	}

	if app.Verbose {
		app.PrintfMsg("wrote SNA snapshot \"%s\"", path)
	}
}

// Signature: func scale(n uint)
func wrapper_scale(t *eval.Thread, in []eval.Value, out []eval.Value) {
	if app.TerminationInProgress() {
		return
	}

	n := in[0].(eval.UintValue).Get(t)

	switch n {
	case 1:
		finished := make(chan byte)
		speccy.CommandChannel <- spectrum.Cmd_CloseAllDisplays{finished}
		<-finished

		speccy.CommandChannel <- spectrum.Cmd_AddDisplay{spectrum.NewSDLScreen(app)}

	case 2:
		finished := make(chan byte)
		speccy.CommandChannel <- spectrum.Cmd_CloseAllDisplays{finished}
		<-finished

		speccy.CommandChannel <- spectrum.Cmd_AddDisplay{spectrum.NewSDLScreen2x(app, /*fullscreen*/ false)}
	}
}

// Signature: func fps(n float)
func wrapper_fps(t *eval.Thread, in []eval.Value, out []eval.Value) {
	if app.TerminationInProgress() {
		return
	}

	fps := in[0].(eval.FloatValue).Get(t)
	speccy.CommandChannel <- spectrum.Cmd_SetFPS{float(fps)}
}

// Signature: func ula_accuracy(accurateEmulation bool)
func wrapper_ulaAccuracy(t *eval.Thread, in []eval.Value, out []eval.Value) {
	if app.TerminationInProgress() {
		return
	}

	accurateEmulation := in[0].(eval.BoolValue).Get(t)
	speccy.CommandChannel <- spectrum.Cmd_SetUlaEmulationAccuracy{accurateEmulation}
}

// Signature: func sound(enable bool)
func wrapper_sound(t *eval.Thread, in []eval.Value, out []eval.Value) {
	if app.TerminationInProgress() {
		return
	}

	enable := in[0].(eval.BoolValue).Get(t)

	if enable {
		audio, err := spectrum.NewSDLAudio(app)
		if err == nil {
			finished := make(chan byte)
			speccy.CommandChannel <- spectrum.Cmd_CloseAllAudioReceivers{finished}
			<-finished

			speccy.CommandChannel <- spectrum.Cmd_AddAudioReceiver{audio}
		} else {
			app.PrintfMsg("%s", err)
		}
	} else {
		finished := make(chan byte)
		speccy.CommandChannel <- spectrum.Cmd_CloseAllAudioReceivers{finished}
		<-finished
	}
}

// Signature: func wait(milliseconds uint)
func wrapper_wait(t *eval.Thread, in []eval.Value, out []eval.Value) {
	if app.TerminationInProgress() {
		return
	}

	milliseconds := uint(in[0].(eval.UintValue).Get(t))
	time.Sleep(1e6 * int64(milliseconds))
}

// Signature: func script(scriptName string)
func wrapper_script(t *eval.Thread, in []eval.Value, out []eval.Value) {
	scriptName := in[0].(eval.StringValue).Get(t)

	err := runScript(w, scriptName, /*optional*/ false)
	if err != nil {
		app.PrintfMsg("%s", err)
		return
	}
}

// Signature: func optionalScript(scriptName string)
func wrapper_optionalScript(t *eval.Thread, in []eval.Value, out []eval.Value) {
	scriptName := in[0].(eval.StringValue).Get(t)

	err := runScript(w, scriptName, /*optional*/ true)
	if err != nil {
		app.PrintfMsg("%s", err)
		return
	}
}

// Signature: func screenshot(screenshotName string)
func wrapper_screenshot(t *eval.Thread, in []eval.Value, out []eval.Value) {
	if app.TerminationInProgress() {
		return
	}

	path := in[0].(eval.StringValue).Get(t)

	ch := make(chan []byte)
	speccy.CommandChannel <- spectrum.Cmd_MakeVideoMemoryDump{ ch }

	data := <-ch

	err := ioutil.WriteFile(path, data, 0600)

	if err != nil {
		app.PrintfMsg("%s", err)
	}

	if app.Verbose {
		app.PrintfMsg("wrote screenshot \"%s\"", path)
	}
}

// Signature: func puts(str string)
func wrapper_puts(t *eval.Thread, in []eval.Value, out []eval.Value) {
	str := in[0].(eval.StringValue).Get(t)
	app.PrintfMsg("%s", str)
}

// ==============
// Initialization
// ==============

func defineFunctions(w *eval.World) {
	{
		var functionSignature func()
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_help, functionSignature)
		w.DefineVar("help", funcType, funcValue)
		help_keys.Push("help()")
		help_vals.Push("This help")
	}

	{
		var functionSignature func()
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_exit, functionSignature)
		w.DefineVar("exit", funcType, funcValue)
		help_keys.Push("exit()")
		help_vals.Push("Terminate this program")
	}

	{
		var functionSignature func()
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_reset, functionSignature)
		w.DefineVar("reset", funcType, funcValue)
		help_keys.Push("reset()")
		help_vals.Push("Reset the emulated machine")
	}

	{
		var functionSignature func(string)
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_addSearchPath, functionSignature)
		w.DefineVar("addSearchPath", funcType, funcValue)
		help_keys.Push("addSearchPath(path string)")
		help_vals.Push("Append a path to the list of paths searched when loading snapshots")
	}

	{
		var functionSignature func(string)
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_load, functionSignature)
		w.DefineVar("load", funcType, funcValue)
		help_keys.Push("load(path string)")
		help_vals.Push("Load state from file (.SNA, .Z80, .Z80.ZIP, etc)")
	}

	{
		var functionSignature func(string)
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_save, functionSignature)
		w.DefineVar("save", funcType, funcValue)
		help_keys.Push("save(path string)")
		help_vals.Push("Save state to file (SNA format)")
	}

	{
		var functionSignature func(uint)
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_scale, functionSignature)
		w.DefineVar("scale", funcType, funcValue)
		help_keys.Push("scale(n uint)")
		help_vals.Push("Change the display scale")
	}

	{
		var functionSignature func(float)
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_fps, functionSignature)
		w.DefineVar("fps", funcType, funcValue)
		help_keys.Push("fps(n float)")
		help_vals.Push("Change the display refresh frequency")
	}

	{
		var functionSignature func(bool)
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_ulaAccuracy, functionSignature)
		w.DefineVar("ula_accuracy", funcType, funcValue)
		help_keys.Push("ula_accuracy(accurateEmulation bool)")
		help_vals.Push("Enable/disable accurate emulation of screen bitmap and screen attributes")
	}

	{
		var functionSignature func(bool)
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_sound, functionSignature)
		w.DefineVar("sound", funcType, funcValue)
		help_keys.Push("sound(enable bool)")
		help_vals.Push("Enable or disable sound")
	}

	{
		var functionSignature func(uint)
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_wait, functionSignature)
		w.DefineVar("wait", funcType, funcValue)
		help_keys.Push("wait(milliseconds uint)")
		help_vals.Push("Wait the specified amount of time before issuing the next command")
	}

	{
		var functionSignature func(string)
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_script, functionSignature)
		w.DefineVar("script", funcType, funcValue)
		help_keys.Push("script(scriptName string)")
		help_vals.Push("Load and evaluate the specified Go script")
	}

	{
		var functionSignature func(string)
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_optionalScript, functionSignature)
		w.DefineVar("optionalScript", funcType, funcValue)
		help_keys.Push("optionalScript(scriptName string)")
		help_vals.Push("Load (if found) and evaluate the specified Go script")
	}

	{
		var functionSignature func(string)
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_screenshot, functionSignature)
		w.DefineVar("screenshot", funcType, funcValue)
		help_keys.Push("screenshot(screenshotName string)")
		help_vals.Push("Take a screenshot of the current display")
	}
	{
		var functionSignature func(string)
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_puts, functionSignature)
		w.DefineVar("puts", funcType, funcValue)
		help_keys.Push("puts(str string)")
		help_vals.Push("Print the given string")
	}

}

// Runs the specified Go source code in the context of 'w'
func run(w *eval.World, sourceCode string) os.Error {
	var err os.Error
	var code eval.Code

	code, err = w.Compile(sourceCode)
	if err != nil {
		app.PrintfMsg("%s", err)
		return err
	}

	_, err = code.Run()
	if err != nil {
		app.PrintfMsg("%s", err)
		return err
	}

	return err
}

// Loads and evaluates the specified Go script
func runScript(w *eval.World, scriptName string, optional bool) os.Error {
	fileName := scriptName + ".go"
	data, err := ioutil.ReadFile(spectrum.ScriptPath(fileName))
	if err != nil {
		if !optional {
			return err
		} else {
			return nil
		}
	}

	var buf bytes.Buffer
	buf.Write(data)
	run(w, buf.String())

	return nil
}


type handler_t byte

func (h handler_t) HandleSignal(s signal.Signal) {
	switch ss := s.(type) {
	case signal.UnixSignal:
		switch ss {
		case signal.SIGQUIT, signal.SIGTERM, signal.SIGALRM, signal.SIGTSTP, signal.SIGTTIN, signal.SIGTTOU:
			readline.CleanupAfterSignal()
			if ss == signal.SIGTSTP {
				syscall.Kill(os.Getpid(), int(signal.SIGSTOP))
			}

		case signal.SIGINT:
			readline.FreeLineState()
			readline.CleanupAfterSignal()

		case signal.SIGWINCH:
			readline.ResizeTerminal()
		}
	}
}

// Reads lines from os.Stdin and sends them through the channel 'code'.
//
// If no more input is available, an arbitrary value is sent through channel 'no_more_code'.
//
// This function is intended to be run in a separate goroutine.
func readCode(app *spectrum.Application, code chan string, no_more_code chan<- byte) {
	handler := handler_t(0)
	spectrum.InstallSignalHandler(handler)

	// BNF pattern: (string address)* nil
	readline_channel := make(chan *string)
	go func() {
		prevMsgOut := app.SetMessageOutput(&consoleMessageOutput{})

		for {
			havePrompt_mutex.Lock()
			havePrompt = true
			havePrompt_mutex.Unlock()

			if app.TerminationInProgress() {
				break
			}

			line := readline.ReadLine(PROMPT)

			havePrompt_mutex.Lock()
			havePrompt = false
			havePrompt_mutex.Unlock()

			readline_channel <- line
			if line == nil {
				break
			} else {
				<-readline_channel
			}
		}

		app.SetMessageOutput(prevMsgOut)
	}()

	evtLoop := app.NewEventLoop()
	for {
		select {
		case <-evtLoop.Pause:
			spectrum.UninstallSignalHandler(handler)

			havePrompt_mutex.Lock()
			if havePrompt && len(PROMPT) > 0 {
				fmt.Printf("\r%s\r", PROMPT_EMPTY)
				havePrompt = false
			}
			havePrompt_mutex.Unlock()

			readline.FreeLineState()
			readline.CleanupAfterSignal()

			evtLoop.Pause <- 0

		case <-evtLoop.Terminate:
			if evtLoop.App().Verbose {
				app.PrintfMsg("readCode loop: exit")
			}
			evtLoop.Terminate <- 0
			return

		case lineP := <-readline_channel:
			// EOF
			if lineP == nil {
				no_more_code <- 0
				evtLoop.Delete()
				continue
			}

			line := strings.TrimSpace(*lineP)

			if len(line) > 0 {
				readline.AddHistory(line)
			}

			code <- line
			<-code

			readline_channel <- nil
		}
	}
}

func Init(_app *spectrum.Application, _speccy *spectrum.Spectrum48k) {
	if app != nil {
		panic("running multiple consoles is unsupported")
	}

	app = _app
	speccy = _speccy
	w = eval.NewWorld()

	defineFunctions(w)

	// Run the startup script
	var err os.Error

	err = runScript(w, STARTUP_SCRIPT, /*optional*/ ignoreStartupScript)
	if err != nil {
		app.PrintfMsg("%s", err)
		app.RequestExit()
		return
	}
}


// Reads lines of Go code from standard input and evaluates the code.
//
// This function exits in two cases: if the application was terminated (from outside of this function),
// or if there is nothing more to read from os.Stdin. The latter can optionally cause the whole application
// to terminate (controlled by the 'exitAppIfEndOfInput' parameter).
func Run(exitAppIfEndOfInput bool) {
	// This should be printed before executing "go readCode(...)",
	// in order to ensure that this message *always* gets printed before printing the prompt
	app.PrintfMsg("Hint: Input an empty line to see available commands")

	// Start a goroutine for reading code from os.Stdin.
	// The code pieces are being received from the channel 'code_chan'.
	code_chan := make(chan string)
	no_more_code := make(chan byte)
	go readCode(app, code_chan, no_more_code)

	// Loop pattern: (read code, run code)+ (terminate app)?
	evtLoop := app.NewEventLoop()
	for {
		select {
		case <-evtLoop.Pause:
			evtLoop.Pause <- 0

		case <-evtLoop.Terminate:
			// Exit this function
			if app.Verbose {
				app.PrintfMsg("console loop: exit")
			}
			evtLoop.Terminate <- 0
			return

		case code := <-code_chan:
			//app.PrintfMsg("code=\"%s\"", code)
			if len(code) > 0 {
				run(w, code)
			} else {
				printHelp()
			}
			code_chan <- "<next>"

		case <-no_more_code:
			if exitAppIfEndOfInput {
				app.RequestExit()
			} else {
				evtLoop.Delete()
			}
		}
	}
}


type consoleMessageOutput struct {
	// This mutex is used to serialize the multiple calls to fmt.Printf
	// used in function PrintfMsg. Otherwise, a concurrent entry to PrintfMsg
	// would cause undesired interleaving of fmt.Printf calls.
	mutex sync.Mutex
}

// Prints a single-line message to 'os.Stdout' using 'fmt.Printf'.
// If the format string does not end with the new-line character,
// the new-line character is appended automatically.
//
// Using this function instead of 'fmt.Printf', 'println', etc,
// ensures proper redisplay of the current command line.
func (out *consoleMessageOutput) PrintfMsg(format string, a ...interface{}) {
	out.mutex.Lock()
	{
		havePrompt_mutex.Lock()
		if havePrompt && len(PROMPT) > 0 {
			fmt.Printf("\r%s\r", PROMPT_EMPTY)
		}
		havePrompt_mutex.Unlock()

		appendNewLine := false
		if (len(format) == 0) || (format[len(format)-1] != '\n') {
			appendNewLine = true
		}

		fmt.Printf(format, a...)
		if appendNewLine {
			fmt.Println()
		}

		havePrompt_mutex.Lock()
		if havePrompt {
			if (app == nil) || !app.TerminationInProgress() {
				readline.OnNewLine()
				readline.Redisplay()
				// 'havePrompt' remains to have the value 'true'
			} else {
				havePrompt = false
			}
		}
		havePrompt_mutex.Unlock()
	}
	out.mutex.Unlock()
}


// Lines below will be uncommented when/if the keypress console
// command will be implemented.

// type uintV uint

// func newUint(v uint) *uintV { 
// 	vp := uintV(v)
// 	return &vp
// }


// func (v *uintV) String() string { return fmt.Sprint(*v) }
// func (v *uintV) Assign(t *eval.Thread, o eval.Value) { *v = uintV(o.(eval.UintValue).Get(t)) }
// func (v *uintV) Get(*eval.Thread) uint64 { return uint64(*v) }
// func (v *uintV) Set(t *eval.Thread, x uint64) { *v = uintV(x) }

// func defineConstants(w *eval.World) {
// 	w.DefineConst("KEY_1", eval.UintType, newUint(0))
// 	w.DefineConst("KEY_2", eval.UintType, newUint(1))
// 	w.DefineConst("KEY_3", eval.UintType, newUint(2))
// 	w.DefineConst("KEY_4", eval.UintType, newUint(3))
// 	w.DefineConst("KEY_5", eval.UintType, newUint(4))
// 	w.DefineConst("KEY_6", eval.UintType, newUint(5))
// 	w.DefineConst("KEY_7", eval.UintType, newUint(6))
// 	w.DefineConst("KEY_8", eval.UintType, newUint(7))
// 	w.DefineConst("KEY_9", eval.UintType, newUint(8))
// 	w.DefineConst("KEY_0", eval.UintType, newUint(9))
// 	w.DefineConst("KEY_Q", eval.UintType, newUint(10))
// 	w.DefineConst("KEY_W", eval.UintType, newUint(11))
// 	w.DefineConst("KEY_E", eval.UintType, newUint(12))
// 	w.DefineConst("KEY_R", eval.UintType, newUint(13))
// 	w.DefineConst("KEY_T", eval.UintType, newUint(14))
// 	w.DefineConst("KEY_Y", eval.UintType, newUint(15))
// 	w.DefineConst("KEY_U", eval.UintType, newUint(16))
// 	w.DefineConst("KEY_I", eval.UintType, newUint(17))
// 	w.DefineConst("KEY_O", eval.UintType, newUint(18))
// 	w.DefineConst("KEY_P", eval.UintType, newUint(19))

// 	w.DefineConst("KEY_A", eval.UintType, newUint(20))
// 	w.DefineConst("KEY_S", eval.UintType, newUint(21))
// 	w.DefineConst("KEY_D", eval.UintType, newUint(22))
// 	w.DefineConst("KEY_F", eval.UintType, newUint(23))
// 	w.DefineConst("KEY_G", eval.UintType, newUint(24))
// 	w.DefineConst("KEY_H", eval.UintType, newUint(25))
// 	w.DefineConst("KEY_J", eval.UintType, newUint(26))
// 	w.DefineConst("KEY_K", eval.UintType, newUint(27))
// 	w.DefineConst("KEY_L", eval.UintType, newUint(28))
// 	w.DefineConst("KEY_Enter", eval.UintType, newUint(29))

// 	w.DefineConst("KEY_CapsShift", eval.UintType, newUint(30))
// 	w.DefineConst("KEY_Z", eval.UintType, newUint(31))
// 	w.DefineConst("KEY_X", eval.UintType, newUint(32))
// 	w.DefineConst("KEY_C", eval.UintType, newUint(33))
// 	w.DefineConst("KEY_V", eval.UintType, newUint(34))
// 	w.DefineConst("KEY_B", eval.UintType, newUint(35))
// 	w.DefineConst("KEY_N", eval.UintType, newUint(36))
// 	w.DefineConst("KEY_M", eval.UintType, newUint(37))
// 	w.DefineConst("KEY_SymbolShift", eval.UintType, newUint(38))
// 	w.DefineConst("KEY_Space", eval.UintType, newUint(39))
// }
