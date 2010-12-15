package formats

import (
	"testing"
	"spectrum/prettytest"
)

func testDecodeZ80(assert *prettytest.T) {
	assert.Pending()
}

func TestZ80Snapshot(t *testing.T) {
	prettytest.Run(
		t,
		testDecodeZ80,
	)
}
