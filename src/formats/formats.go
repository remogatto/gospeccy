// Decoding and encoding of ZX Spectrum emulator file formats
package formats

import (
	"errors"
	"io/ioutil"
	"path"
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

type Snapshot interface {
	CpuState() CpuState
	UlaState() UlaState
	Memory() *[48 * 1024]byte
}

type FullSnapshot struct {
	Cpu CpuState
	Ula UlaState
	Mem [48 * 1024]byte
}

func (s *FullSnapshot) CpuState() CpuState {
	return s.Cpu
}

func (s *FullSnapshot) UlaState() UlaState {
	return s.Ula
}

func (s *FullSnapshot) Memory() *[48 * 1024]byte {
	return &s.Mem
}

type SnapshotData []byte

type Archive interface {
	Filenames() []string
	Read(fileIndex int) ([]byte, error)
}

const (
	FORMAT_SNA = iota
	FORMAT_Z80
	FORMAT_TAP
)

const (
	ENCAPSULATION_NONE = iota
	ENCAPSULATION_ZIP
)

type FormatInfo struct {
	Format        int
	Encapsulation int
}

// Determines the format of the specified file based on its name,
// or based on the names of embedded files in case the file is an archive.
// Returns an error if the format could not be detected.
func DetectFormat(filePath string) (*FormatInfo, error) {
	return detectFormat(filePath, ENCAPSULATION_NONE, true)
}

func detectFormat(filePath string, encapsulation int, allowEncapsulation bool) (*FormatInfo, error) {
	ext := strings.ToLower(path.Ext(filePath))

	switch ext {
	case ".sna":
		return &FormatInfo{FORMAT_SNA, encapsulation}, nil

	case ".z80":
		return &FormatInfo{FORMAT_Z80, encapsulation}, nil

	case ".tap":
		return &FormatInfo{FORMAT_TAP, encapsulation}, nil

	case ".zip":
		if (encapsulation == ENCAPSULATION_NONE) && allowEncapsulation {
			archive, err := ReadZipFile(filePath)
			if err != nil {
				return nil, err
			}

			var embeddedFile_format *FormatInfo
			{
				n := 0
				for _, name := range archive.Filenames() {
					format, err := detectFormat(name, ENCAPSULATION_ZIP, false)
					if err == nil {
						embeddedFile_format = format
						n++
					}
				}

				if n == 0 {
					return nil, errors.New("the archive does not contain any supported files")
				}
				if n >= 2 {
					return nil, errors.New("the archive contains multiple supported files")
				}
			}

			return embeddedFile_format, nil
		} else {
			return nil, errors.New("unrecognized file format")
		}
	}

	return nil, errors.New("unrecognized file format")
}

// Decode a snapshot from binary data.
// The filename is a hint used to determine the snapshot format.
func (data SnapshotData) Decode(format int) (Snapshot, error) {
	switch format {
	case FORMAT_SNA:
		return data.DecodeSNA()

	case FORMAT_Z80:
		return data.DecodeZ80()
	}

	return nil, errors.New("unknown snapshot format")
}

func readZIP(filePath string) (interface{}, error) {
	archive, err := ReadZipFile(filePath)
	if err != nil {
		return nil, err
	}

	var embeddedFile_index int
	var embeddedFile_format *FormatInfo
	{
		n := 0
		for i, name := range archive.Filenames() {
			format, err := detectFormat(name, ENCAPSULATION_ZIP, false)
			if err == nil {
				embeddedFile_index = i
				embeddedFile_format = format
				n++
			}
		}

		if n == 0 {
			return nil, errors.New("the archive does not contain any supported program files")
		}
		if n >= 2 {
			return nil, errors.New("the archive contains multiple program files")
		}
	}

	var data []byte
	data, err = archive.Read(embeddedFile_index)
	if err != nil {
		return nil, err
	}

	if embeddedFile_format.Format == FORMAT_TAP {
		return NewTAP(data)
	}

	return SnapshotData(data).Decode(embeddedFile_format.Format)
}

// Read a program from the specified file.
// Return the program and errors if any.
// The file can be compressed.
func ReadProgram(filePath string) (interface{}, error) {
	ext := strings.ToLower(path.Ext(filePath))

	// ZIP archive
	if ext == ".zip" {
		return readZIP(filePath)
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var format *FormatInfo
	format, err = detectFormat(filePath, ENCAPSULATION_NONE, false)
	if err != nil {
		return nil, err
	}

	if format.Format == FORMAT_TAP {
		return NewTAP(data)
	}

	return SnapshotData(data).Decode(format.Format)
}

func splitWord(word uint16) (byte, byte) {
	return byte(word >> 8), byte(word)
}
