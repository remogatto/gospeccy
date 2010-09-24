package formats

import (
	"os"
	"strings"
)

const (
	TStatesPerFrame = 69888
	InterruptLength = 32
)

type CpuState struct {
	A, F, B, C, D, E, H, L         byte
	A_, F_, B_, C_, D_, E_, H_, L_ byte
	IX, IY                         uint16
	I, R, IFF1, IFF2, IM           byte
	SP, PC                         uint16

	Tstate uint
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


type Snapshot interface {
	CpuState() CpuState
	UlaState() UlaState
	Memory() *[48 * 1024]byte
	RETN() bool
}

type SnapshotData []byte

func (data SnapshotData) Decode(filename string) (snapshot Snapshot, err os.Error) {
	switch {
	case strings.HasSuffix(strings.ToLower(filename), ".sna"):
		return data.DecodeSNA()

	case strings.HasSuffix(strings.ToLower(filename), ".z80"):
		return data.DecodeZ80()
	}

	return nil, os.NewError("unable to detect the snapshot format (missing or unknown filename extension)")
}


func splitWord(word uint16) (byte, byte) {
	return byte(word >> 8), byte(word)
}
