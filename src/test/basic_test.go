package test

import (
	"testing"
	pt "spectrum/prettytest"
)

func should_load_system_ROM(t *pt.T) {
	t.True(screenEqualTo("testdata/system_rom_loaded.sna"))
}

func TestBasicFeature(t *testing.T) {
	pt.Describe(
		t,
		"The emulator",
		should_load_system_ROM,

		before,
		after,
	)
}
