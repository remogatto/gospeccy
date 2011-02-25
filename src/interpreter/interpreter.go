package interpreter

import (
	"bytes"
	"clingon"
	"container/vector"
	"exp/eval"
	"fmt"
	"go/token"
	"io/ioutil"
	"os"
	"spectrum"
	"spectrum/formats"
	"time"
)

// ==============
// Some variables
// ==============

// These variables are set only once, before starting new goroutines,
// so there is no need for controlling concurrent access via a sync.Mutex
var (
	app                 *spectrum.Application
	cmdLineArg          string // The 1st non-flag command-line argument, or empty string
	speccy              *spectrum.Spectrum48k
	renderer            spectrum.Renderer
	w                   *eval.World
	buffer              *bytes.Buffer
	IgnoreStartupScript = false
)

const (
	SCRIPT_DIRECTORY = "scripts"
	STARTUP_SCRIPT   = "startup"
)

type Interpreter struct{}

func (i *Interpreter) Run(console *clingon.Console, command string) os.Error {
	var err os.Error
	buffer = bytes.NewBufferString("")
	if command == "" {
		err = i.run(w, "", "help()")
		console.Print(buffer.String())
		return err
	}

	err = i.run(w, "", command)
	console.Print(buffer.String())
	return err
}

// Runs the specified Go source code in the context of 'w'
func (i *Interpreter) run(w *eval.World, path_orEmpty string, sourceCode string) os.Error {
	var err os.Error
	var code eval.Code

	fileSet := token.NewFileSet()
	if len(path_orEmpty) > 0 {
		fileSet.AddFile(path_orEmpty, fileSet.Base(), len(sourceCode))
	}

	code, err = w.Compile(fileSet, sourceCode)
	if err != nil {
		return err
	}

	_, err = code.Run()
	if err != nil {
		return err
	}

	return err
}

// ================
// Various commands
// ================

var help_keys vector.StringVector
var help_vals vector.StringVector

// Signature: func help()
func wrapper_help(t *eval.Thread, in []eval.Value, out []eval.Value) {
	fmt.Fprintf(buffer, "\nAvailable commands:\n")

	maxKeyLen := 1
	for i := 0; i < help_keys.Len(); i++ {
		if len(help_keys[i]) > maxKeyLen {
			maxKeyLen = len(help_keys[i])
		}
	}

	for i := 0; i < help_keys.Len(); i++ {
		fmt.Fprintf(buffer, "  %s", help_keys[i])
		for j := len(help_keys[i]); j < maxKeyLen; j++ {
			fmt.Fprintf(buffer, " ")
		}
		fmt.Fprintf(buffer, "  %s\n", help_vals[i])
	}
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
	speccy.CommandChannel <- spectrum.Cmd_Reset{nil}
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

	path = spectrum.ProgramPath(path)

	var program interface{}
	program, err := formats.ReadProgram(path)
	if err != nil {
		fmt.Fprintf(buffer, "%s\n", err)
		return
	}

	if _, isTAP := program.(*formats.TAP); isTAP {
		romLoaded := make(chan (<-chan bool))
		speccy.CommandChannel <- spectrum.Cmd_Reset{romLoaded}
		<-(<-romLoaded)
	}

	errChan := make(chan os.Error)
	speccy.CommandChannel <- spectrum.Cmd_Load{path, program, errChan}

	err = <-errChan
	if err != nil {
		fmt.Fprintf(buffer, "%s\n", err)
		return
	}
}

// Signature: func cmdLineArg() string
func wrapper_cmdLineArg(t *eval.Thread, in []eval.Value, out []eval.Value) {
	out[0].(eval.StringValue).Set(t, cmdLineArg)
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
		fmt.Fprintf(buffer, "%s\n", err)
		return
	}

	err = ioutil.WriteFile(path, data, 0600)
	if err != nil {
		fmt.Fprintf(buffer, "%s\n", err)
	}

	if app.Verbose {
		fmt.Fprintf(buffer, "wrote SNA snapshot \"%s\"", path)
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
		renderer.Resize(app, false, false)
	case 2:
		finished := make(chan byte)
		speccy.CommandChannel <- spectrum.Cmd_CloseAllDisplays{finished}
		<-finished
		renderer.Resize(app, true, false)
	}
}

// Signature: func fullscreen(enable bool)
func wrapper_fullscreen(t *eval.Thread, in []eval.Value, out []eval.Value) {
	if app.TerminationInProgress() {
		return
	}
	enable := in[0].(eval.BoolValue).Get(t)
	if enable {
		finished := make(chan byte)
		speccy.CommandChannel <- spectrum.Cmd_CloseAllDisplays{finished}
		<-finished
		renderer.Resize(app, true, true)
	} else {
		finished := make(chan byte)
		speccy.CommandChannel <- spectrum.Cmd_CloseAllDisplays{finished}
		<-finished
		renderer.Resize(app, true, false)
	}
}

// Signature: func fps(n float32)
func wrapper_fps(t *eval.Thread, in []eval.Value, out []eval.Value) {
	if app.TerminationInProgress() {
		return
	}

	fps := in[0].(eval.FloatValue).Get(t)
	speccy.CommandChannel <- spectrum.Cmd_SetFPS{float32(fps), nil}
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
			fmt.Fprintf(buffer, "%s\n", err)
			return
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
		fmt.Fprintf(buffer, "%s\n", err)
		return
	}
}

// Signature: func optionalScript(scriptName string)
func wrapper_optionalScript(t *eval.Thread, in []eval.Value, out []eval.Value) {
	scriptName := in[0].(eval.StringValue).Get(t)

	err := runScript(w, scriptName, /*optional*/ true)
	if err != nil {
		fmt.Fprintf(buffer, "%s\n", err)
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
	speccy.CommandChannel <- spectrum.Cmd_MakeVideoMemoryDump{ch}

	data := <-ch

	err := ioutil.WriteFile(path, data, 0600)

	if err != nil {
		fmt.Fprintf(buffer, "%s\n", err)
	}

	if app.Verbose {
		app.PrintfMsg("wrote screenshot \"%s\"", path)
	}
}

// Signature: func puts(str string)
func wrapper_puts(t *eval.Thread, in []eval.Value, out []eval.Value) {
	str := in[0].(eval.StringValue).Get(t)
	fmt.Fprintf(buffer, "%s", str)
}

// Signature: func acceleratedLoad(on bool)
func wrapper_acceleratedLoad(t *eval.Thread, in []eval.Value, out []eval.Value) {
	value := in[0].(eval.BoolValue).Get(t)
	speccy.EnableAcceleratedLoad(value)
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
		help_vals.Push("Append to the paths searched when loading snapshots")
	}
	{
		var functionSignature func() string
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_cmdLineArg, functionSignature)
		w.DefineVar("cmdLineArg", funcType, funcValue)
		help_keys.Push("cmdLineArg() string)")
		help_vals.Push("The 1st non-flag command-line argument, or an empty string")
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
		help_vals.Push("Change the display scale (1 or 2)")
	}
	{
		var functionSignature func(float32)
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_fps, functionSignature)
		w.DefineVar("fps", funcType, funcValue)
		help_keys.Push("fps(n float32)")
		help_vals.Push("Change the display refresh frequency (0=default FPS)")
	}
	{
		var functionSignature func(bool)
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_ulaAccuracy, functionSignature)
		w.DefineVar("ula", funcType, funcValue)
		help_keys.Push("ula(accurateEmulation bool)")
		help_vals.Push("Enable/disable accurate ULA emulation")
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
		help_vals.Push("Wait before executing the next command")
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
	{
		var functionSignature func(bool)
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_fullscreen, functionSignature)
		w.DefineVar("fullscreen", funcType, funcValue)
		help_keys.Push("fullscreen(enable bool)")
		help_vals.Push("Fullscreen on/off")
	}
	{
		var functionSignature func(bool)
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_acceleratedLoad, functionSignature)
		w.DefineVar("acceleratedLoad", funcType, funcValue)
		help_keys.Push("acceleratedLoad(on bool)")
		help_vals.Push("Set accelerated tape load on/off")
	}

}

// Loads and evaluates the specified Go script
func runScript(w *eval.World, scriptName string, optional bool) os.Error {
	var i Interpreter
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
	i.run(w, fileName, buf.String())

	return nil
}

func Init(_app *spectrum.Application, _cmdLineArg string, _speccy *spectrum.Spectrum48k, _renderer spectrum.Renderer) {
	app = _app
	cmdLineArg = _cmdLineArg
	speccy = _speccy
	renderer = _renderer

	if w == nil {
		w = eval.NewWorld()
		defineFunctions(w)

		// Run the startup script
		var err os.Error

		err = runScript(w, STARTUP_SCRIPT, /*optional*/ IgnoreStartupScript)
		if err != nil {
			app.PrintfMsg("%s", err)
			app.RequestExit()
			return
		}
	}
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
