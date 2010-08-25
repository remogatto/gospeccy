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
)


type Console struct {
	app *Application
}

// ==============
// Some variables
// ==============

var console Console
var speccy *Spectrum48k


// ================
// Various commands
// ================

var help_keys vector.StringVector
var help_vals vector.StringVector

func printHelp() {
	fmt.Printf("\nAvailable commands:\n")

	maxKeyLen := 1
	for i := 0; i < help_keys.Len(); i++ {
		if len(help_keys[i]) > maxKeyLen {
			maxKeyLen = len(help_keys[i])
		}
	}

	for i := 0; i < help_keys.Len(); i++ {
		fmt.Printf("    %s", help_keys[i])
		for j := len(help_keys[i]); j < maxKeyLen; j++ {
			fmt.Print(" ")
		}
		fmt.Printf("  %s\n", help_vals[i])
	}

	fmt.Printf("\n")
}

// Signature: func help()
func wrapper_help(t *eval.Thread, in []eval.Value, out []eval.Value) {
	printHelp()
}

// Signature: func exit()
func wrapper_exit(t *eval.Thread, in []eval.Value, out []eval.Value) {
	console.app.RequestExit()
}

// Signature: func reset()
func wrapper_reset(t *eval.Thread, in []eval.Value, out []eval.Value) {
	speccy.CommandChannel <- Cmd_Reset{}
}

// Signature: func load(path string)
func wrapper_load(t *eval.Thread, in []eval.Value, out []eval.Value) {
	path := in[0].(eval.StringValue).Get(t)

	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	errChan := make(chan os.Error)
	speccy.CommandChannel <- Cmd_LoadSna{path, data, errChan}
	err = <-errChan
	if err != nil {
		fmt.Printf("%s\n", err)
	}
}

// Signature: func save(path string)
func wrapper_save(t *eval.Thread, in []eval.Value, out []eval.Value) {
	path := in[0].(eval.StringValue).Get(t)

	ch := make(chan Snapshot)
	speccy.CommandChannel <- Cmd_SaveSna{ch}

	var snapshot Snapshot = <-ch
	if snapshot.err != nil {
		fmt.Printf("%s\n", snapshot.err)
		return
	}

	err := ioutil.WriteFile(path, snapshot.data, 0600)
	if err != nil {
		fmt.Printf("%s\n", err)
	}

	if console.app.Verbose {
		fmt.Printf("wrote SNA snapshot \"%s\"\n", path)
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
	if fps < 0 {
		fps = DefaultFPS
	}
	speccy.FPS <- float(fps)
}

// Signature: func ULA_accuracy(accurateEmulation bool)
func wrapper_ulaAccuracy(t *eval.Thread, in []eval.Value, out []eval.Value) {
	accurateEmulation := in[0].(eval.BoolValue).Get(t)
	speccy.CommandChannel <- Cmd_SetUlaEmulationAccuracy{accurateEmulation}
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
}


// Runs the specified Go source code in the context of 'w'
func run(w *eval.World, sourceCode string) {
	// Avoids the need to put ";" at the end of the code
	sourceCode = strings.Join([]string{sourceCode, "\n"}, "")

	var err os.Error

	var code eval.Code
	code, err = w.Compile(sourceCode)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = code.Run()
	if err != nil {
		fmt.Println(err)
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
		}
	}
}

// Reads lines from os.Stdin and sends them through the channel 'code'.
//
// If no more input is available, an arbitrary value is sent through channel 'no_more_code'.
//
// This function is intended to be run in a separate goroutine.
func readCode(app *Application, code chan string, no_more_code chan byte) {
	handler := handler_t(0)
	InstallSignalHandler(handler)

	// BNF pattern: (string address)* nil
	readline_channel := make(chan *string)
	go func() {
		for {
			line := readline.ReadLine("gospeccy> ")
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
			fmt.Println()
			UninstallSignalHandler(handler)
			readline.FreeLineState()
			readline.CleanupAfterSignal()
			evtLoop.Pause <- 0

		case <-evtLoop.Terminate:
			if evtLoop.App().Verbose {
				println("readCode loop: exit")
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


// Reads Go code from os.Stdin and evaluates it.
//
// This function exits in two cases: if the application was terminated (from outside of this function),
// or if there is nothing more to read from os.Stdin. The latter can optionally cause the whole application to terminate.
func RunConsole(app *Application, _speccy *Spectrum48k, exitAppIfEndOfInput bool) {
	console = Console{app}
	speccy = _speccy

	w := eval.NewWorld()
	defineFunctions(w)

	// This should be printed before executing "go readCode(...)",
	// in order to ensure that this message *always* gets printed before printing the prompt
	fmt.Printf("Hint: Input an empty line to see available commands\n")

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
				println("console loop: exit")
			}
			evtLoop.Terminate <- 0
			return

		case code := <-code_chan:
			//fmt.Printf("code=\"%s\"\n", code)
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
