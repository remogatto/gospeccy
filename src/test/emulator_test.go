package test

import (
	"testing"
	"prettytest"
)

func TestEmulator(t *testing.T) {
	prettytest.RunWithFormatter(
		t,
		&prettytest.BDDFormatter{"The emulator"},
		new(testSuite),
	)
}
