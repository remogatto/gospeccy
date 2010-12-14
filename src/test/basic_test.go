package test

import (
	"prettytest"
	"testing"
)

type basic_suite_t struct {
	test_suite_t
}

func (s *basic_suite_t) should_load_system_ROM() {
	s.True(screenEqualTo("testdata/system_rom_loaded.sna"))
}

func TestBasicFeature(t *testing.T) {
	prettytest.RunWithFormatter(
		t,
		&prettytest.BDDFormatter{"Basic emulator tests"},
		new(basic_suite_t),
	)
}
