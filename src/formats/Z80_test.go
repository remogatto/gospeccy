package formats

import (
	"prettytest"
	"testing"
)

type z80_suite_t struct {
	prettytest.Suite
}

func (s *z80_suite_t) testDecodeZ80() {
	s.Pending()
}

func TestZ80Snapshot(t *testing.T) {
	prettytest.Run(
		t,
		new(z80_suite_t),
	)
}
