package formats

import (
	"io/ioutil"
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
}


type SnapshotData []byte

type Archive interface {
	Filenames() []string
	Read(fileIndex int) ([]byte, os.Error)
}


const (
	FORMAT_SNA = iota
	FORMAT_Z80
)

// Derives the snapshot format from the filename,
// or returns -1 if the type could not be detected
func typeFromSuffix(filename string) (int, os.Error) {
	fName := strings.ToLower(filename)

	switch {
	case strings.HasSuffix(fName, ".sna"):
		return FORMAT_SNA, nil

	case strings.HasSuffix(fName, ".z80"):
		return FORMAT_Z80, nil
	}

	return -1, os.NewError("unable to detect the snapshot format (missing or unknown filename extension)")
}

// Decode a snapshot from binary data.
// The filename is a hint used to determine the snapshot format.
func (data SnapshotData) Decode(format int) (Snapshot, os.Error) {
	switch format {
	case FORMAT_SNA:
		return data.DecodeSNA()

	case FORMAT_Z80:
		return data.DecodeZ80()
	}

	return nil, os.NewError("unknown snapshot format")
}

// Read a snapshot from the specified file.
// The file can be compressed.
func ReadSnapshot(filePath string) (Snapshot, os.Error) {
	fName := strings.ToLower(filePath)

	// ZIP archive
	if strings.HasSuffix(fName, ".zip") {
		archive, err := ReadZipFile(filePath)
		if err != nil {
			return nil, err
		}

		var archive_fileIndex int
		var archive_snapshotFormat int
		{
			n := 0
			for i, name := range archive.Filenames() {
				format, err := typeFromSuffix(name)
				if err == nil {
					archive_fileIndex = i
					archive_snapshotFormat = format
					n++
				}
			}

			if n == 0 {
				return nil, os.NewError("the archive does not contain any supported snapshot files")
			}
			if n >= 2 {
				return nil, os.NewError("the archive contains multiple snapshot files")
			}
		}

		var data []byte
		data, err = archive.Read(archive_fileIndex)
		if err != nil {
			return nil, err
		}

		return SnapshotData(data).Decode(archive_snapshotFormat)
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var format int
	format, err = typeFromSuffix(filePath)
	if err != nil {
		return nil, err
	}

	return SnapshotData(data).Decode(format)
}


func splitWord(word uint16) (byte, byte) {
	return byte(word >> 8), byte(word)
}
