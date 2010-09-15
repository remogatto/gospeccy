package spectrum

import (
	"fmt"
)

// Prints a single-line message to 'os.Stdout' using 'fmt.Printf'.
// If the format string does not end with the new-line character,
// the new-line character is appended automatically.
//
// This is a simplistic version of function "PrintfMsg"
// defined in file "console.go". The purpose of this simplified
// version is to enable tests and benchmarking without having to
// install the C readline wrapper.
func PrintfMsg(format string, a ...interface{}) {
	appendNewLine := false
	if (len(format) == 0) || (format[len(format)-1] != '\n') {
		appendNewLine = true
	}

	fmt.Printf(format, a)

	if appendNewLine {
		fmt.Println()
	}
}
