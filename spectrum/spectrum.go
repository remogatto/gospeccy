package spectrum

type Spectrum48k struct {
	Cpu    *Z80
	Memory MemoryAccessor
	Port   PortAccessor
}

func NewSpectrum48k(memory MemoryAccessor, port PortAccessor) *Spectrum48k {
	// Load the built-in ROM image into memory
	for address, b := range rom48k {
		memory.set(uint(address), b)
	}

	// Initialize keyStates
	for row := 0; row < 8; row++ {
		keyStates[byte(row)] = 0xff
	}

	return &Spectrum48k{Cpu: NewZ80(memory, port), Memory: memory, Port: port}
}

// Execute 69888 T-states
func (speccy *Spectrum48k) doOpcodes() {
	eventNextEvent = 69888
	tstates = 0
	speccy.Cpu.doOpcodes()
}

func (speccy *Spectrum48k) interrupt() {
	speccy.Cpu.interrupt()
}

func (speccy *Spectrum48k) RenderFrame() {
	speccy.doOpcodes()
	speccy.Memory.renderScreen()
	speccy.interrupt()
}

func (speccy *Spectrum48k) LoadSna(filename string) {
	speccy.Cpu.LoadSna(filename)
}

func (speccy *Spectrum48k) KeyDown(keySym uint) {
	keyCode, ok := keyCodes[keySym]
	if ok {
		keyStates[keyCode.row] &= ^(keyCode.mask)
	}
}

func (speccy *Spectrum48k) KeyUp(keySym uint) {
	keyCode, ok := keyCodes[keySym]
	if ok {
		keyStates[keyCode.row] |= (keyCode.mask)
	}
}

// func dumpRegisters() {
// 	fmt.Printf("%02x%02x %02x%02x %02x%02x %02x%02x %02x%02x %02x%02x %02x%02x %02x%02x %02x%02x %02x%02x %04x %04x\n",
// 		z80.a, z80.f, z80.b, z80.c, z80.d, z80.e, z80.h, z80.l, z80.a_, z80.f_, z80.b_, z80.c_, z80.d_, z80.e_, z80.h_, z80.l_, z80.ixh, z80.ixl, z80.iyh, z80.iyl, z80.sp, z80.pc)
// 	fmt.Printf("%02x %02x %d %d %d %d %d\n", z80.i, (z80.r7&0x80)|(z80.r&0x7f),
// 		z80.iff1, z80.iff2, z80.im, z80.halted, tstates)
// }

// func dumpMemory() {
// 	for i, val := range memory {
// 		fmt.Printf("%d %d\n", i, val)
// 	}
// }
