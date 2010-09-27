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

const LABEL_FAIL = "\033[31;1mFAIL\033[0m"
const LABEL_PASS = "\033[32;1mOK\033[0m"
const LABEL_PENDING= "\033[33;1mPENDING\033[0m"

const (
	STATUS_PASS = iota
	STATUS_FAIL
	STATUS_PENDING
)

const formatTag = "\t%s\t"

type callerInfo struct {
	name, fn string
	line int
}

type Assertions struct {
	T   *testing.T
	Status byte
	Dry bool

	callerInfo *callerInfo
}

func newCallerInfo(skip int) *callerInfo {
	pc, fn, line, ok := runtime.Caller(skip)
	if !ok {
		panic("An error occured while retrieving caller info!")
	}
	callerName := runtime.FuncForPC(pc).Name()
	return &callerInfo{callerName, fn, line}
}

func (assertion *Assertions) fail(exp, act interface{}, info *callerInfo) *Assertions {
	if !assertion.Dry {
		assertion.T.Errorf("Expected %s but got %s -- %s:%d\n", exp, act, info.fn, info.line)
	}
	assertion.Status = STATUS_FAIL
	return assertion
}

func (assertion *Assertions) setup() *Assertions {
	assertion.callerInfo = newCallerInfo(3)
	return assertion
}

// Assert that the expected value equals the actual value. Return true
// on success.
func (assertion *Assertions) Equal(exp, act interface{}) *Assertions {
	assertion.setup()
	if exp != act {
		return assertion.fail(exp, act, assertion.callerInfo)
	}
	return assertion
}

// Assert that the value is true.
func (assertion *Assertions) True(value bool) *Assertions {
	assertion.setup()
	if !value {
		return assertion.fail("true", "false", assertion.callerInfo)
	}
	return assertion
}

// Assert that the value is false.
func (assertion *Assertions) False(value bool) *Assertions {
	assertion.setup()
	if value {
		return assertion.fail("false", "true", assertion.callerInfo)
	}
	return assertion
}

// Mark the test function as pending.
func (assertion *Assertions) Pending() *Assertions {
	assertion.setup()
	assertion.Status = STATUS_PENDING
	return assertion
}

// Check if the test function has failed.
func (assertion *Assertions) IsFailed() bool {
	return assertion.Status == STATUS_FAIL
}

// Run tests.
func Run(t *testing.T, tests ...func(*Assertions) *Assertions) {
	pc, _, _, _ := runtime.Caller(1)
	callerName := runtime.FuncForPC(pc).Name()
	fmt.Printf("\n%s:\n", callerName)
	for _, test := range tests {
		assertions := &Assertions{t, STATUS_PASS, false, nil}
		test(assertions)
		switch assertions.Status {
		case STATUS_FAIL:
			fmt.Printf(formatTag+"%s\n", LABEL_FAIL, assertions.callerInfo.name)
		case STATUS_PASS:
			fmt.Printf(formatTag + "%s\n", LABEL_PASS, assertions.callerInfo.name)
		case STATUS_PENDING:
			fmt.Printf(formatTag + "%s\n", LABEL_PENDING, assertions.callerInfo.name)
			
		}
	}
}

// Run tests but don't emit output and don't fail on failing
// assertions.
func DryRun(t *testing.T, tests ...func(*Assertions) *Assertions) {
	for _, test := range tests {
		test(&Assertions{t, STATUS_PASS, true, nil})
	}
}
