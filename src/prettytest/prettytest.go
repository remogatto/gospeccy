/*

Copyright (c) 2010 Andrea Fazzi

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
"Software"), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

*/

/*

PrettyTest is a simple assertion testing library for golang. It aims
to simplify/prettify testing in golang.

It features:

* a simple assertion vocabulary for better readability
* colorful output

*/

package prettytest

import (
	"testing"
	"runtime"
	"fmt"
)

const FAIL = "\033[31;1mFAIL\033[0m"
const PASS = "\033[32;1mOK\033[0m"
const NOT_YET_IMPLEMENTED = "\033[33;1mNYI\033[0m"

const formatTag = "\t%s\t"

type Assertions struct {
	T   *testing.T
	Dry bool
}

func (assertion *Assertions) reportPASS(pass string) {
	if !assertion.Dry {
		fmt.Print(pass)
	}
}

func (assertion *Assertions) reportFAIL(fail, expected string) {
	if !assertion.Dry {
		fmt.Print(fail)
		assertion.T.Errorf(expected)
	}
}

// Assert that the expected value equals the actual value. Return true
// on success.
func (assertion *Assertions) Equal(exp, act interface{}) bool {
	if exp != act {
		_, fn, line, _ := runtime.Caller(1)
		assertion.reportFAIL(
			fmt.Sprintf(formatTag+"%s == %s\n", FAIL, act, exp),
			fmt.Sprintf("Expected %s but got %s -- %s:%d", exp, act, fn, line),
		)
		return false
	} else {
		assertion.reportPASS(fmt.Sprintf(formatTag+"%s == %s\n", PASS, act, exp))
	}
	return true
}

// Assert that the value is true.
func (assertion *Assertions) True(value bool) bool {
	if !value {
		_, fn, line, _ := runtime.Caller(1)
		assertion.reportFAIL(
			fmt.Sprintf(formatTag+"value == true\n", FAIL),
			fmt.Sprintf("Expected true but got false -- %s:%d", fn, line),
		)
		return false
	} else {
		assertion.reportPASS(fmt.Sprintf(formatTag+"value == true\n", PASS))
	}
	return true
}

// Assert that the value is false.
func (assertion *Assertions) False(value bool) bool {
	if value {
		_, fn, line, _ := runtime.Caller(1)
		assertion.reportFAIL(
			fmt.Sprintf(formatTag+"value == false\n", FAIL),
			fmt.Sprintf("Expected false but got true -- %s:%d", fn, line),
		)
		return false
	} else {
		assertion.reportPASS(fmt.Sprintf(formatTag+"value == false\n", PASS))
	}
	return true
}

func (assertion *Assertions) Pending(msg string) {
	fmt.Printf(formatTag+"%s\n", NOT_YET_IMPLEMENTED, msg)
}

// Run tests.
func Run(t *testing.T, description string, tests ...func(*Assertions)) {
	fmt.Printf("\n%s:\n", description)
	for _, test := range tests {
		test(&Assertions{t, false})
	}
}

// Run tests but don't emit output and don't fail on failing
// assertions.
func DryRun(t *testing.T, tests ...func(*Assertions)) {
	for _, test := range tests {
		test(&Assertions{t, true})
	}
}
