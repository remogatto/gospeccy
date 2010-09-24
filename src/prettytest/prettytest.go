package prettytest

import (
	"testing"
	"runtime"
	"fmt"
)

const FAIL = "\033[31;1mFAIL\033[0m"
const PASS = "\033[32;1mOK\033[0m"
const formatTag = "\t%s\t"

type assertions struct {
	T *testing.T
	Dry bool
}

func (assertion *assertions) reportPASS(pass string) {
	if !assertion.Dry {
		fmt.Print(pass)
	}	
}

func (assertion *assertions) reportFAIL(fail, expected string) {
	if !assertion.Dry {
		fmt.Print(fail)
		assertion.T.Errorf(expected)
	}	
}

// Assert that the expected value equals the actual value. Return true
// on success.
func (assertion *assertions) Equal(exp, act interface{}) bool {
	if exp != act {
		_, fn, line, _ := runtime.Caller(1)
		assertion.reportFAIL(
			fmt.Sprintf(formatTag + "%s == %s\n", FAIL, act, exp), 
			fmt.Sprintf("Expected %s but got %s -- %s:%d", exp, act, fn, line),
		)
		return false
	} else {
		assertion.reportPASS(fmt.Sprintf(formatTag + "%s == %s\n", PASS, act, exp))
	}
	return true
}

// Assert that the value is true.
func (assertion *assertions) True(value bool) bool {
	if !value {
		_, fn, line, _ := runtime.Caller(1)
		assertion.reportFAIL(
			fmt.Sprintf(formatTag + "value == true\n", FAIL), 
			fmt.Sprintf("Expected true but got false -- %s:%d", fn, line),
		)
		return false
	} else {
		assertion.reportPASS(fmt.Sprintf(formatTag + "value == true\n", PASS))
	}
	return true
}

// Assert that the value is false.
func (assertion *assertions) False(value bool) bool {
	if value {
		_, fn, line, _ := runtime.Caller(1)
		assertion.reportFAIL(
			fmt.Sprintf(formatTag + "value == false\n", FAIL), 
			fmt.Sprintf("Expected false but got true -- %s:%d", fn, line),
		)
		return false
	} else {
		assertion.reportPASS(fmt.Sprintf(formatTag + "value == false\n", PASS))
	}
	return true
}

// Run tests.
func Run(t *testing.T, description string, tests... func (*assertions)) {
	fmt.Printf("\n%s:\n", description)
	for _, test := range tests {
		test(&assertions{t, false})
	}
}

// Run tests but don't emit output and don't fail on failing
// assertions.
func DryRun(t *testing.T, tests... func (*assertions)) {
	for _, test := range tests {
		test(&assertions{t, true})
	}
}

