package formats

type CpuState struct {
	A, F, B, C, D, E, H, L         byte
	A_, F_, B_, C_, D_, E_, H_, L_ byte
	IX, IY                         uint16
	I, R, IFF1, IFF2, IM           byte
	SP, PC                         uint16
}

type UlaState struct {
	// 0..7
	Border byte
}

type FullSnapshot struct {
	Cpu    CpuState
	Ula    UlaState
	Memory [48 * 1024]byte
}

type SnapshotData []byte

func splitWord(word uint16) (byte, byte) {
	return byte(word >> 8), byte(word)
}
