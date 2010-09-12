package spectrum

import (
	"exp/eval"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"container/vector"
	"⚛readline"
	"sync"
	"bytes"
)


// ==============
// Some variables
// ==============

// These variables are set only once, before starting new goroutines,
// so there is no need for controlling concurrent access via a sync.Mutex
var app *Application
var speccy *Spectrum48k

const PROMPT = "gospeccy> "
const PROMPT_EMPTY = "          "

// Whether the terminal is currently showing a prompt string
var havePrompt = false
var havePrompt_mutex sync.Mutex


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

	PrintfMsg("%s\n", buf.String())
}

// Signature: func help()
func wrapper_help(t *eval.Thread, in []eval.Value, out []eval.Value) {
	printHelp()
}

// Signature: func exit()
func wrapper_exit(t *eval.Thread, in []eval.Value, out []eval.Value) {
	app.RequestExit()
}

// Signature: func reset()
func wrapper_reset(t *eval.Thread, in []eval.Value, out []eval.Value) {
	speccy.CommandChannel <- Cmd_Reset{}
}

// Signature: func load(path string)
func wrapper_load(t *eval.Thread, in []eval.Value, out []eval.Value) {
	path := in[0].(eval.StringValue).Get(t)

	data, err := ioutil.ReadFile(SnaPath(path))
	if err != nil {
		PrintfMsg("%s", err)
		return
	}

	errChan := make(chan os.Error)
	speccy.CommandChannel <- Cmd_LoadSna{path, data, errChan}
	err = <-errChan
	if err != nil {
		PrintfMsg("%s", err)
	}
}

// Signature: func save(path string)
func wrapper_save(t *eval.Thread, in []eval.Value, out []eval.Value) {
	path := in[0].(eval.StringValue).Get(t)

	ch := make(chan Snapshot)
	speccy.CommandChannel <- Cmd_SaveSna{ch}

	var snapshot Snapshot = <-ch
	if snapshot.Err != nil {
		PrintfMsg("%s", snapshot.Err)
		return
	}

	err := ioutil.WriteFile(path, snapshot.Data, 0600)
	if err != nil {
		PrintfMsg("%s", err)
	}

	if app.Verbose {
		PrintfMsg("wrote SNA snapshot \"%s\"", path)
	}
}

// Signature: func scale(n uint)
func wrapper_scale(t *eval.Thread, in []eval.Value, out []eval.Value) {
	n := in[0].(eval.UintValue).Get(t)

	switch n {
	case 1:
		speccy.CommandChannel <- Cmd_CloseAllDisplays{}
		speccy.CommandChannel <- Cmd_AddDisplay{NewSDLScreen(speccy.app)}

	case 2:
		speccy.CommandChannel <- Cmd_CloseAllDisplays{}
		speccy.CommandChannel <- Cmd_AddDisplay{NewSDLScreen2x(speccy.app, /*fullscreen*/ false)}
	}
}

// Signature: func fps(n float)
func wrapper_fps(t *eval.Thread, in []eval.Value, out []eval.Value) {
	fps := in[0].(eval.FloatValue).Get(t)
	speccy.CommandChannel <- Cmd_SetFPS{float(fps)}
}

// Signature: func ULA_accuracy(accurateEmulation bool)
func wrapper_ulaAccuracy(t *eval.Thread, in []eval.Value, out []eval.Value) {
	accurateEmulation := in[0].(eval.BoolValue).Get(t)
	speccy.CommandChannel <- Cmd_SetUlaEmulationAccuracy{accurateEmulation}
}

// Signature: func sound(enable bool)
func wrapper_sound(t *eval.Thread, in []eval.Value, out []eval.Value) {
	enable := in[0].(eval.BoolValue).Get(t)

	if enable {
		audio, err := NewSDLAudio(speccy.app)
		if err == nil {
			speccy.CommandChannel <- Cmd_CloseAllAudioReceivers{}
			speccy.CommandChannel <- Cmd_AddAudioReceiver{audio}
		} else {
			PrintfMsg("%s", err)
			app.RequestExit()
		}
	} else {
		speccy.CommandChannel <- Cmd_CloseAllAudioReceivers{}
	}
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
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_load, functionSignature)
		w.DefineVar("load", funcType, funcValue)
		help_keys.Push("load(path string)")
		help_vals.Push("Load state from file (SNA format)")
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
		w.DefineVar("ULA_accuracy", funcType, funcValue)
		help_keys.Push("ULA_accuracy(accurateEmulation bool)")
		help_vals.Push("Enable/disable accurate emulation of screen bitmap and screen attributes")
	}

	{
		var functionSignature func(bool)
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_sound, functionSignature)
		w.DefineVar("sound", funcType, funcValue)
		help_keys.Push("sound(enable bool)")
		help_vals.Push("Enable or disable sound")
	}
}


// Runs the specified Go source code in the context of 'w'
func run(w *eval.World, sourceCode string) {
	// Avoids the need to put ";" at the end of the code
	sourceCode = strings.Join([]string{sourceCode, "\n"}, "")

	var err os.Error

	var code eval.Code
	code, err = w.Compile(sourceCode)
	if err != nil {
		PrintfMsg("%s", err)
		return
	}

	_, err = code.Run()
	if err != nil {
		PrintfMsg("%s", err)
		return
	}
}


type handler_t byte

func (h handler_t) HandleSignal(s signal.Signal) {
	switch ss := s.(type) {
	case signal.UnixSignal:
		switch ss {
		case signal.SIGQUIT, signal.SIGTERM, signal.SIGALRM, signal.SIGTSTP, signal.SIGTTIN, signal.SIGTTOU:
			readline.CleanupAfterSignal()

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
func readCode(app *Application, code chan string, no_more_code chan<- byte) {
	handler := handler_t(0)
	InstallSignalHandler(handler)

	// BNF pattern: (string address)* nil
	readline_channel := make(chan *string)
	go func() {
		for {
			havePrompt_mutex.Lock()
			havePrompt = true
			havePrompt_mutex.Unlock()

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
	}()

	evtLoop := app.NewEventLoop()
	for {
		select {
		case <-evtLoop.Pause:
			UninstallSignalHandler(handler)

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
				PrintfMsg("readCode loop: exit")
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


// Reads lines of Go code from standard input and evaluates the code.
//
// This function exits in two cases: if the application was terminated (from outside of this function),
// or if there is nothing more to read from os.Stdin. The latter can optionally cause the whole application
// to terminate (controlled by the 'exitAppIfEndOfInput' parameter).
func RunConsole(_app *Application, _speccy *Spectrum48k, exitAppIfEndOfInput bool) {
	app = _app
	speccy = _speccy

	w := eval.NewWorld()
	defineFunctions(w)

	// This should be printed before executing "go readCode(...)",
	// in order to ensure that this message *always* gets printed before printing the prompt
	PrintfMsg("Hint: Input an empty line to see available commands")

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
			if evtLoop.App().Verbose {
				PrintfMsg("console loop: exit")
			}
			evtLoop.Terminate <- 0
			return

		case code := <-code_chan:
			//PrintfMsg("code=\"%s\"", code)
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


// Prints a single-line message to 'os.Stdout' using 'fmt.Printf'.
// If the format string does not end with the new-line character,
// the new-line character is appended automatically.
//
// Using this function instead of 'fmt.Printf', 'println', etc,
// ensures proper redisplay of the current command line.
func PrintfMsg(format string, a ...interface{}) {
	msg_mutex.Lock()
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

		fmt.Printf(format, a)
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
	msg_mutex.Unlock()
}

// This mutex is used to serialize the multiple calls to fmt.Printf
// used in function PrintfMsg. Otherwise, a concurrent entry to PrintfMsg
// would cause undesired interleaving of fmt.Printf calls.
var msg_mutex sync.Mutex
