package formats

import (
	"testing"
	"spectrum/prettytest"
)

func testDecodeZ80(assert *prettytest.Assertions) {
	assert.Pending("testDecodeZ80")
}

func TestZ80Snapshot(t *testing.T) {
	prettytest.Run(
		t,
		"TestZ80Snapshot",
		testDecodeZ80,
	)
}
