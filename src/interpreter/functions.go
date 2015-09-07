package interpreter

import (
	"bytes"
	"fmt"
	"github.com/remogatto/gospeccy/src/formats"
	"github.com/remogatto/gospeccy/src/spectrum"
	"github.com/sbinet/go-eval"
	"io/ioutil"
	"net/url"
	"strings"
	"time"
)

// Contains the names of all defined functions
var definedFunctions = make(map[string]byte)

func defineFunction(name string, t *eval.FuncType, f eval.FuncValue) {
	w.DefineVar(name, t, f)
	definedFunctions[name] = 0
}

type Function struct {
	Name       string         // Name of the variable used to access the function
	Type       *eval.FuncType // Function type
	Value      eval.FuncValue // The function itself
	Help_key   string         // Help
	Help_value string
}

var functionsToAdd []Function

func DefineFunction(f Function) {
	mutex.Lock()
	defer mutex.Unlock()

	if w == nil {
		// Postpone the function definition until interpreter initialization
		functionsToAdd = append(functionsToAdd, f)
	} else {
		defineFunction(f.Name, f.Type, f.Value)

		if (f.Help_key != "") && (f.Help_value != "") {
			help_keys = append(help_keys, f.Help_key)
			help_vals = append(help_vals, f.Help_value)
		}
	}
}

// ================
// Various commands
// ================

var help_keys []string
var help_vals []string

// Signature: func help()
func wrapper_help(t *eval.Thread, in []eval.Value, out []eval.Value) {
	fmt.Fprintf(stdout, "\nAvailable commands:\n")

	maxKeyLen := 1
	for i := 0; i < len(help_keys); i++ {
		if len(help_keys[i]) > maxKeyLen {
			maxKeyLen = len(help_keys[i])
		}
	}

	for i := 0; i < len(help_keys); i++ {
		fmt.Fprintf(stdout, "  %s", help_keys[i])
		for j := len(help_keys[i]); j < maxKeyLen; j++ {
			fmt.Fprintf(stdout, " ")
		}
		fmt.Fprintf(stdout, "  %s\n", help_vals[i])
	}
}

// Signature: func exit()
func wrapper_exit(t *eval.Thread, in []eval.Value, out []eval.Value) {
	// Implementation note:
	//   The following test has to be there only in cases in which something can go wrong.
	//   For example if the user tried to execute "exit(); audio(false)" then GoSpeccy would panic.
	//   An alternative way would be to actually terminate the whole program at the 1st statement - so that
	//   ["audio(false)" or whatever else] is not executed - alas this is somewhat problematic,
	//   since once the script "exit(); audio(false)" runs, it cannot be stopped halfway
	//   through its execution. Using "runtime.Goexit()" would solve this issue, but only partially,
	//   since it is potentially possible for the statement "audio(false)" to be hidden in a defer statement.
	//   So, the best option (until somebody implements a better one) is to convert the problematic commands
	//   into statements that are doing nothing while the application is in the process of being exited.
	if app.TerminationInProgress() || app.Terminated() {
		return
	}
	app.RequestExit()
}

// Signature: func vars() []string
func wrapper_vars(t *eval.Thread, in []eval.Value, out []eval.Value) {
	vars := make([]eval.Value, 0, len(intp.vars))

	for varName, _ := range intp.vars {
		s := string_value_t(varName)
		vars = append(vars, &s)
	}

	var slice eval.Slice
	base := array_value_t(vars)
	slice.Base = &base
	slice.Len = int64(len(vars))
	slice.Cap = int64(len(vars))

	out[0].(eval.SliceValue).Set(t, slice)
}

// Signature: func definedFunction(name string) bool
func wrapper_definedFunction(t *eval.Thread, in []eval.Value, out []eval.Value) {
	name := in[0].(eval.StringValue).Get(t)
	_, defined := definedFunctions[name]
	out[0].(eval.BoolValue).Set(t, defined)
}

// Signature: func reset()
func wrapper_reset(t *eval.Thread, in []eval.Value, out []eval.Value) {
	if app.TerminationInProgress() || app.Terminated() {
		return
	}
	romLoaded := make(chan (<-chan bool))
	speccy.CommandChannel <- spectrum.Cmd_Reset{romLoaded}
	<-(<-romLoaded)
}

// Signature: func addSearchPath(path string)
func wrapper_addSearchPath(t *eval.Thread, in []eval.Value, out []eval.Value) {
	path := in[0].(eval.StringValue).Get(t)
	spectrum.AddCustomSearchPath(path)
}

// Signature: func setDownloadPath(path string)
func wrapper_setDownloadPath(t *eval.Thread, in []eval.Value, out []eval.Value) {
	path := in[0].(eval.StringValue).Get(t)
	spectrum.SetDownloadPath(path)
}

func load(path string) {
	var program interface{}
	program, err := formats.ReadProgram(path)
	if err != nil {
		fmt.Fprintf(stdout, "%s\n", err)
		return
	}

	if _, isTAP := program.(*formats.TAP); isTAP {
		romLoaded := make(chan (<-chan bool))
		speccy.CommandChannel <- spectrum.Cmd_Reset{romLoaded}
		<-(<-romLoaded)
	}

	errChan := make(chan error)
	speccy.CommandChannel <- spectrum.Cmd_Load{path, program, errChan}

	err = <-errChan
	if err != nil {
		fmt.Fprintf(stdout, "%s\n", err)
		return
	}
}

// Signature: func load(path string)
func wrapper_load(t *eval.Thread, in []eval.Value, out []eval.Value) {
	if app.TerminationInProgress() || app.Terminated() {
		return
	}

	path := in[0].(eval.StringValue).Get(t)

	var err error
	path, err = spectrum.ProgramPath(path)
	if err != nil {
		fmt.Fprintf(stdout, "%s\n", err)
		return
	}

	load(path)
}

// Signature: func cmdLineArg() string
func wrapper_cmdLineArg(t *eval.Thread, in []eval.Value, out []eval.Value) {
	out[0].(eval.StringValue).Set(t, cmdLineArg)
}

// Signature: func save(path string)
func wrapper_save(t *eval.Thread, in []eval.Value, out []eval.Value) {
	if app.TerminationInProgress() || app.Terminated() {
		return
	}

	path := in[0].(eval.StringValue).Get(t)

	ch := make(chan *formats.FullSnapshot)
	speccy.CommandChannel <- spectrum.Cmd_MakeSnapshot{ch}

	fullSnapshot := <-ch

	data, err := fullSnapshot.EncodeSNA()
	if err != nil {
		fmt.Fprintf(stdout, "%s\n", err)
		return
	}

	err = ioutil.WriteFile(path, data, 0600)
	if err != nil {
		fmt.Fprintf(stdout, "%s\n", err)
	}

	if app.Verbose {
		fmt.Fprintf(stdout, "wrote SNA snapshot \"%s\"", path)
	}
}

// Signature: func fps(n float32)
func wrapper_fps(t *eval.Thread, in []eval.Value, out []eval.Value) {
	if app.TerminationInProgress() || app.Terminated() {
		return
	}

	fps := in[0].(eval.FloatValue).Get(t)
	speccy.CommandChannel <- spectrum.Cmd_SetFPS{float32(fps), nil}
}

// Signature: func ula_accuracy(accurateEmulation bool)
func wrapper_ulaAccuracy(t *eval.Thread, in []eval.Value, out []eval.Value) {
	if app.TerminationInProgress() || app.Terminated() {
		return
	}

	accurateEmulation := in[0].(eval.BoolValue).Get(t)
	speccy.CommandChannel <- spectrum.Cmd_SetUlaEmulationAccuracy{accurateEmulation}
}

// Signature: func wait(milliseconds uint)
func wrapper_wait(t *eval.Thread, in []eval.Value, out []eval.Value) {
	if app.TerminationInProgress() || app.Terminated() {
		return
	}

	milliseconds := in[0].(eval.UintValue).Get(t)
	time.Sleep(time.Millisecond * time.Duration(milliseconds))
}

// Signature: func script(scriptName string)
func wrapper_script(t *eval.Thread, in []eval.Value, out []eval.Value) {

	if app.TerminationInProgress() || app.Terminated() {
		return
	}

	path := in[0].(eval.StringValue).Get(t)

	var err error
	path, err = spectrum.ScriptPath(path)
	if err != nil {
		fmt.Fprintf(stdout, "%s\n", err)
		return
	}

	err = runScript(w, path, false /*optional*/)
	if err != nil {
		fmt.Fprintf(stdout, "%s\n", err)
		return
	}
}

// Signature: func optionalScript(scriptName string)
func wrapper_optionalScript(t *eval.Thread, in []eval.Value, out []eval.Value) {
	scriptName := in[0].(eval.StringValue).Get(t)

	err := runScript(w, scriptName, true /*optional*/)
	if err != nil {
		fmt.Fprintf(stdout, "%s\n", err)
		return
	}
}

// Signature: func screenshot(screenshotName string)
func wrapper_screenshot(t *eval.Thread, in []eval.Value, out []eval.Value) {
	if app.TerminationInProgress() || app.Terminated() {
		return
	}

	path := in[0].(eval.StringValue).Get(t)

	ch := make(chan []byte)
	speccy.CommandChannel <- spectrum.Cmd_MakeVideoMemoryDump{ch}

	data := <-ch

	err := ioutil.WriteFile(path, data, 0600)

	if err != nil {
		fmt.Fprintf(stdout, "%s\n", err)
	}

	if app.Verbose {
		fmt.Fprintf(stdout, "wrote screenshot \"%s\"", path)
	}
}

// Signature: func puts(str string)
func wrapper_puts(t *eval.Thread, in []eval.Value, out []eval.Value) {
	str := in[0].(eval.StringValue).Get(t)
	fmt.Fprintf(stdout, "%s", str)
}

// Signature: func acceleratedLoad(on bool)
func wrapper_acceleratedLoad(t *eval.Thread, in []eval.Value, out []eval.Value) {
	if app.TerminationInProgress() || app.Terminated() {
		return
	}

	enable := in[0].(eval.BoolValue).Get(t)
	speccy.CommandChannel <- spectrum.Cmd_SetAcceleratedLoad{enable}
}

type WOS struct {
	URL         string
	MachineType string
	Publication string
	Score       string
}

func url_printer(URL eval.Value) string {
	s := URL.(eval.StringValue).Get(nil)

	if len(s) > 60 {
		var buf bytes.Buffer

		i := 0
		for _, rune := range s {
			if i < 10 {
				buf.WriteRune(rune)
			} else if i == 10 {
				buf.WriteString("...")
			} else if (i > 10) && (i < len(s)-(60-3)) {
				// Nothing
			} else {
				buf.WriteRune(rune)
			}

			i++
		}

		s = buf.String()
	}

	return s
}

// Signature: func wosFind(pattern string) []WOS
func wrapper_wosFind(t *eval.Thread, in []eval.Value, out []eval.Value) {
	if app.TerminationInProgress() || app.Terminated() {
		return
	}

	pattern := in[0].(eval.StringValue).Get(t)
	pattern = strings.Replace(pattern, " ", "*", -1)

	var records []spectrum.WosRecord
	records, err := spectrum.WosQuery(app, "regexp="+url.QueryEscape(pattern))
	if err != nil {
		fmt.Fprintf(stdout, "%s", err)

		var emptySlice eval.Slice
		out[0].(eval.SliceValue).Set(t, emptySlice)
		return
	}

	var woss []eval.Value
	for _, r := range records {
		for _, url := range r.FtpFiles {
			mac_value := string_value_t(r.MachineType)
			pub_value := string_value_t(r.Publication)
			sco_value := string_value_t(r.Score)
			url_value := string_value_t(url)

			var wos struct_value_t
			wos.fields = []eval.Value{&url_value, &mac_value, &pub_value, &sco_value}
			wos.names = []string{"URL", "Machine type", "Publication", "Score"}
			wos.printers_orNil = []func(eval.Value) string{url_printer, nil, nil, nil}
			wos.hide_orNil = nil
			wos.printStyle = MULTI_LINE

			woss = append(woss, &wos)
		}
	}

	var slice eval.Slice
	base := array_value_t(woss)
	slice.Base = &base
	slice.Len = int64(len(woss))
	slice.Cap = int64(len(woss))

	out[0].(eval.SliceValue).Set(t, slice)
}

// Signature: func wosDownload(wos WOS) string
func wrapper_wosDownload(t *eval.Thread, in []eval.Value, out []eval.Value) {
	if app.TerminationInProgress() || app.Terminated() {
		return
	}

	var url string = in[0].(eval.StructValue).Field(t, 0).(eval.StringValue).Get(t)
	filePath, err := spectrum.WosGet(app, stdout, url)
	if err != nil {
		fmt.Fprintf(stdout, "%s", err)
		out[0].(eval.StringValue).Set(t, "")
		return
	}

	out[0].(eval.StringValue).Set(t, filePath)
}

// Signature: func wosLoad(wos WOS)
func wrapper_wosLoad(t *eval.Thread, in []eval.Value, out []eval.Value) {
	if app.TerminationInProgress() || app.Terminated() {
		return
	}

	var url string = in[0].(eval.StructValue).Field(t, 0).(eval.StringValue).Get(t)
	filePath, err := spectrum.WosGet(app, stdout, url)
	if err != nil {
		fmt.Fprintf(stdout, "%s", err)
		return
	}

	load(filePath)
}

// ==============
// Initialization
// ==============

func defineFunctions(w *eval.World) {
	{
		var functionSignature func()
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_help, functionSignature)
		defineFunction("help", funcType, funcValue)
		help_keys = append(help_keys, "help()")
		help_vals = append(help_vals, "This help")
	}
	{
		var functionSignature func()
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_exit, functionSignature)
		defineFunction("exit", funcType, funcValue)
		help_keys = append(help_keys, "exit()")
		help_vals = append(help_vals, "Terminate this program")
	}
	{
		var functionSignature func() []string
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_vars, functionSignature)
		defineFunction("vars", funcType, funcValue)
		help_keys = append(help_keys, "vars()")
		help_vals = append(help_vals, "Get the names of all variables")
	}
	{
		var functionSignature func()
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_reset, functionSignature)
		defineFunction("reset", funcType, funcValue)
		help_keys = append(help_keys, "reset()")
		help_vals = append(help_vals, "Reset the emulated machine")
	}
	{
		var functionSignature func(string) bool
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_definedFunction, functionSignature)
		defineFunction("definedFunction", funcType, funcValue)
		help_keys = append(help_keys, "definedFunction(name string) bool")
		help_vals = append(help_vals, "Returns whether a Go function exists")
	}
	{
		var functionSignature func(string)
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_addSearchPath, functionSignature)
		defineFunction("addSearchPath", funcType, funcValue)
		help_keys = append(help_keys, "addSearchPath(path string)")
		help_vals = append(help_vals, "Append to the paths searched when loading snapshots, scripts, etc")
	}
	{
		var functionSignature func(string)
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_setDownloadPath, functionSignature)
		defineFunction("setDownloadPath", funcType, funcValue)
		help_keys = append(help_keys, "setDownloadPath(path string)")
		help_vals = append(help_vals, `Set path where to download files (""=default path)`)
	}
	{
		var functionSignature func() string
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_cmdLineArg, functionSignature)
		defineFunction("cmdLineArg", funcType, funcValue)
		help_keys = append(help_keys, "cmdLineArg() string)")
		help_vals = append(help_vals, "The 1st non-flag command-line argument, or an empty string")
	}
	{
		var functionSignature func(string)
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_load, functionSignature)
		defineFunction("load", funcType, funcValue)
		help_keys = append(help_keys, "load(path string)")
		help_vals = append(help_vals, "Load state from file (.SNA, .Z80, .Z80.ZIP, etc)")
	}
	{
		var functionSignature func(string)
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_save, functionSignature)
		defineFunction("save", funcType, funcValue)
		help_keys = append(help_keys, "save(path string)")
		help_vals = append(help_vals, "Save state to file (SNA format)")
	}
	{
		var functionSignature func(float32)
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_fps, functionSignature)
		defineFunction("fps", funcType, funcValue)
		help_keys = append(help_keys, "fps(n float32)")
		help_vals = append(help_vals, "Change the display refresh frequency (0=default FPS)")
	}
	{
		var functionSignature func(bool)
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_ulaAccuracy, functionSignature)
		defineFunction("ula", funcType, funcValue)
		help_keys = append(help_keys, "ula(accurateEmulation bool)")
		help_vals = append(help_vals, "Enable/disable accurate ULA emulation")
	}
	{
		var functionSignature func(uint)
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_wait, functionSignature)
		defineFunction("wait", funcType, funcValue)
		help_keys = append(help_keys, "wait(milliseconds uint)")
		help_vals = append(help_vals, "Wait before executing the next command")
	}
	{
		var functionSignature func(string)
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_script, functionSignature)
		defineFunction("script", funcType, funcValue)
		help_keys = append(help_keys, "script(scriptName string)")
		help_vals = append(help_vals, "Load and evaluate the specified Go script")
	}
	{
		var functionSignature func(string)
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_optionalScript, functionSignature)
		defineFunction("optionalScript", funcType, funcValue)
		help_keys = append(help_keys, "optionalScript(scriptName string)")
		help_vals = append(help_vals, "Load (if found) and evaluate the specified Go script")
	}
	{
		var functionSignature func(string)
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_screenshot, functionSignature)
		defineFunction("screenshot", funcType, funcValue)
		help_keys = append(help_keys, "screenshot(screenshotName string)")
		help_vals = append(help_vals, "Take a screenshot of the current display")
	}
	{
		var functionSignature func(string)
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_puts, functionSignature)
		defineFunction("puts", funcType, funcValue)
		help_keys = append(help_keys, "puts(str string)")
		help_vals = append(help_vals, "Print the given string")
	}
	{
		var functionSignature func(bool)
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_acceleratedLoad, functionSignature)
		defineFunction("acceleratedLoad", funcType, funcValue)
		help_keys = append(help_keys, "acceleratedLoad(on bool)")
		help_vals = append(help_vals, "Set accelerated tape load on/off")
	}
	{
		var functionSignature func(string) []WOS
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_wosFind, functionSignature)
		defineFunction("wosFind", funcType, funcValue)
		help_keys = append(help_keys, "wosFind(pattern string) []WOS")
		help_vals = append(help_vals, "Find tapes and snapshots on worldofspectrum.org")
	}
	{
		var functionSignature func(WOS) string
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_wosDownload, functionSignature)
		defineFunction("wosDownload", funcType, funcValue)
		help_keys = append(help_keys, "wosDownload(wos WOS) string")
		help_vals = append(help_vals, "Download from worldofspectrum.org")
	}
	{
		var functionSignature func(WOS)
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_wosLoad, functionSignature)
		defineFunction("wosLoad", funcType, funcValue)
		help_keys = append(help_keys, "wosLoad(wos WOS)")
		help_vals = append(help_vals, "Same as load(wosDownload(wos))")
	}

	for _, f := range functionsToAdd {
		defineFunction(f.Name, f.Type, f.Value)

		if (f.Help_key != "") && (f.Help_value != "") {
			help_keys = append(help_keys, f.Help_key)
			help_vals = append(help_vals, f.Help_value)
		}
	}
	functionsToAdd = nil
}
