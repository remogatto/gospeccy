// Wrapper around the GNU readline(3) library
package readline

// #include <stdio.h>
// #include <stdlib.h>
// #include <readline/readline.h>
// #include <readline/history.h>
import "C"
import "unsafe"

func init() {
	C.rl_catch_signals = 0
}

// Calls the C readline function.
// If prompt is an empty string, the C function will receive NULL.
func ReadLine(prompt string) *string {
	var p *C.char

	//readline allows an empty prompt(NULL)
	if prompt != "" {
		p = C.CString(prompt)
	}

	ret := C.readline(p)

	if p != nil {
		C.free(unsafe.Pointer(p))
	}

	if ret == nil {
		// EOF
		return nil
	}

	s := C.GoString(ret)
	C.free(unsafe.Pointer(ret))
	return &s
}

func CleanupAfterSignal() {
	C.rl_cleanup_after_signal()
}

func FreeLineState() {
	C.rl_free_line_state()
}

func ResetAfterSignal() {
	C.rl_reset_after_signal()
}

func AddHistory(s string) {
	p := C.CString(s)
	C.add_history(p)
	C.free(unsafe.Pointer(p))
}

func ClearHistory() {
	C.clear_history()
}
