package spectrum

import (
	"bytes"
	"exp/eval"
	"fmt"
	"os"
	"strings"
	"container/vector"
)


type Console struct {
	app *Application
}

// ==============
// Some variables
// ==============

var console Console
var speccy *Spectrum48k

var exitted = false


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
	exitted = true
}

// Signature: func reset()
func wrapper_reset(t *eval.Thread, in []eval.Value, out []eval.Value) {
	speccy.CommandChannel <- Cmd_Reset{}
}

// Signature: func load(path string)
func wrapper_load(t *eval.Thread, in []eval.Value, out []eval.Value) {
	path := in[0].(eval.StringValue).Get(t)

	errChan := make(chan os.Error)
	speccy.CommandChannel <- Cmd_LoadSna{path, errChan}
	err := <-errChan
	if err != nil {
		fmt.Printf("%s\n", err)
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


// ==============
// Initialization
// ==============

func defineFunctions(w *eval.World) {
	{
		var help_functionSignature func()
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_help, help_functionSignature)
		w.DefineVar("help", funcType, funcValue)
		help_keys.Push("help()")
		help_vals.Push("This help")
	}

	{
		var exit_functionSignature func()
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_exit, exit_functionSignature)
		w.DefineVar("exit", funcType, funcValue)
		help_keys.Push("exit()")
		help_vals.Push("Terminate this program")
	}

	{
		var reset_functionSignature func()
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_reset, reset_functionSignature)
		w.DefineVar("reset", funcType, funcValue)
		help_keys.Push("reset()")
		help_vals.Push("Reset the emulated machine")
	}

	{
		var load_functionSignature func(string)
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_load, load_functionSignature)
		w.DefineVar("load", funcType, funcValue)
		help_keys.Push("load(path string)")
		help_vals.Push("Load .sna file")
	}

	{
		var scale_functionSignature func(uint)
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_scale, scale_functionSignature)
		w.DefineVar("scale", funcType, funcValue)
		help_keys.Push("scale(n uint)")
		help_vals.Push("Change the display scale")
	}

	{
		var fps_functionSignature func(float)
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_fps, fps_functionSignature)
		w.DefineVar("fps", funcType, funcValue)
		help_keys.Push("fps(n float)")
		help_vals.Push("Change the display refresh frequency")
	}
}


// Runs the specified Go source code in the context of 'w'
func run(w *eval.World, sourceCode string) {
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

// Reads lines from os.Stdin and sends them through the channel 'code'.
//
// If no more input is available, an arbitrary value is sent through channel 'no_more_code'
// and the control returns from this function.
//
// This function is intended to be run in a separate goroutine.
func readCode(code chan string, no_more_code chan byte) {
	var err os.Error
	for (err == nil) && !exitted {
		// Read a line of text (until a new-line character or an EOF)
		var buf bytes.Buffer
		for {
			b := make([]byte, 1)
			var n int
			n, err = os.Stdin.Read(b)

			// This goroutine got blocked on the 'os.Stdin.Read'.
			// In the meantime the application might have exitted.
			if exitted {
				no_more_code <- 0
				return
			}

			if (n == 0) && (err == os.EOF) {
				break
			}
			if err != nil {
				fmt.Printf("%s\n", err)
				break
			}
			if (len(b) > 0) && (b[0] == '\n') {
				break
			}

			buf.Write(b)
		}

		line := strings.TrimSpace(buf.String())

		code <- line
		<-code
	}

	no_more_code <- 0
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

	// Start a goroutine for reading code from os.Stdin.
	// The code pieces are being received from the channel 'code_chan'.
	code_chan := make(chan string)
	no_more_code := make(chan byte)
	go readCode(code_chan, no_more_code)

	fmt.Printf("Hint: Input an empty line to see available commands\n")

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
