// +build linux freebsd

package sdl_output

import (
	intp "github.com/remogatto/gospeccy/src/interpreter"
	"github.com/sbinet/go-eval"
	"sync"
)

type userInterfaceSettings_t interface {
	Terminated() bool

	ResizeVideo(scale2x, fullscreen bool)
	ShowPaintedRegions(enable bool)
	EnableAudio(enable bool)
	SetAudioFreq(freq uint) // 0 means "default frequency"
	SetAudioQuality(hqAudio bool)
}

var uiSettings userInterfaceSettings_t

var mutex sync.Mutex

func setUI(ui userInterfaceSettings_t) {
	mutex.Lock()
	uiSettings = ui
	mutex.Unlock()
}

// Signature: func scale(n uint)
func wrapper_scale(t *eval.Thread, in []eval.Value, out []eval.Value) {
	if uiSettings.Terminated() {
		return
	}
	n := in[0].(eval.UintValue).Get(t)
	switch n {
	case 1:
		mutex.Lock()
		uiSettings.ResizeVideo(false, false)
		mutex.Unlock()
	case 2:
		mutex.Lock()
		uiSettings.ResizeVideo(true, false)
		mutex.Unlock()
	}
}

// Signature: func fullscreen(enable bool)
func wrapper_fullscreen(t *eval.Thread, in []eval.Value, out []eval.Value) {
	if uiSettings.Terminated() {
		return
	}
	enable := in[0].(eval.BoolValue).Get(t)
	if enable {
		mutex.Lock()
		uiSettings.ResizeVideo(true, true)
		mutex.Unlock()
	} else {
		mutex.Lock()
		uiSettings.ResizeVideo(true, false)
		mutex.Unlock()
	}
}

// Signature: func showPaint(enable bool)
func wrapper_showPaint(t *eval.Thread, in []eval.Value, out []eval.Value) {
	if uiSettings.Terminated() {
		return
	}

	enable := in[0].(eval.BoolValue).Get(t)

	mutex.Lock()
	uiSettings.ShowPaintedRegions(enable)
	mutex.Unlock()
}

// Signature: func audio(enable bool)
func wrapper_audio(t *eval.Thread, in []eval.Value, out []eval.Value) {
	if uiSettings.Terminated() {
		return
	}

	enable := in[0].(eval.BoolValue).Get(t)

	mutex.Lock()
	uiSettings.EnableAudio(enable)
	mutex.Unlock()
}

// Signature: func audioFreq(freq uint)
func wrapper_audioFreq(t *eval.Thread, in []eval.Value, out []eval.Value) {
	if uiSettings.Terminated() {
		return
	}

	freq := uint(in[0].(eval.UintValue).Get(t))

	mutex.Lock()
	uiSettings.SetAudioFreq(freq)
	mutex.Unlock()
}

// Signature: func audioHQ(enable bool)
func wrapper_audioHQ(t *eval.Thread, in []eval.Value, out []eval.Value) {
	if uiSettings.Terminated() {
		return
	}

	hqAudio := in[0].(eval.BoolValue).Get(t)

	mutex.Lock()
	uiSettings.SetAudioQuality(hqAudio)
	mutex.Unlock()
}

func defineFunctions() {
	{
		var functionSignature func(uint)
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_scale, functionSignature)
		intp.DefineFunction(intp.Function{
			Name:       "scale",
			Type:       funcType,
			Value:      funcValue,
			Help_key:   "scale(n uint)",
			Help_value: "Change the display scale (1 or 2)",
		})
	}
	{
		var functionSignature func(bool)
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_fullscreen, functionSignature)
		intp.DefineFunction(intp.Function{
			Name:       "fullscreen",
			Type:       funcType,
			Value:      funcValue,
			Help_key:   "fullscreen(enable bool)",
			Help_value: "Fullscreen on/off",
		})
	}
	{
		var functionSignature func(bool)
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_showPaint, functionSignature)
		intp.DefineFunction(intp.Function{
			Name:       "showPaint",
			Type:       funcType,
			Value:      funcValue,
			Help_key:   "showPaint(enable bool)",
			Help_value: "Show painted regions",
		})
	}
	{
		var functionSignature func(bool)
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_audio, functionSignature)
		intp.DefineFunction(intp.Function{
			Name:       "audio",
			Type:       funcType,
			Value:      funcValue,
			Help_key:   "audio(enable bool)",
			Help_value: "Enable or disable audio",
		})
	}
	{
		var functionSignature func(uint)
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_audioFreq, functionSignature)
		intp.DefineFunction(intp.Function{
			Name:       "audioFreq",
			Type:       funcType,
			Value:      funcValue,
			Help_key:   "audioFreq(freq uint)",
			Help_value: "Set audio playback frequency (0=default frequency)",
		})
	}
	{
		var functionSignature func(bool)
		funcType, funcValue := eval.FuncFromNativeTyped(wrapper_audioHQ, functionSignature)
		intp.DefineFunction(intp.Function{
			Name:       "audioHQ",
			Type:       funcType,
			Value:      funcValue,
			Help_key:   "audioHQ(enable bool)",
			Help_value: "Enable or disable high-quality audio",
		})
	}
}

func init() {
	defineFunctions()
}
