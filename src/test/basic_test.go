package test

func (t *testSuite) should_load_system_ROM() {
	t.True(screenEqualTo("testdata/system_rom_loaded.sna"))
}
