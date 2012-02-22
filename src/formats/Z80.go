package formats

import "errors"

type Z80 struct {
	cpu CpuState
	ula UlaState
	mem [48 * 1024]byte

	samRom                   bool
	issue2_emulation         bool
	doubleInterruptFrequency bool
	videoSynchronization     byte // 0..3
	joystick                 byte // 0..3
}

const (
	_Z80_V1_HEADER_SIZE  = 30
	_Z80_V2_HEADER_SIZE  = 30 + 2 + 23
	_Z80_V3_HEADER_SIZE  = 30 + 2 + 54
	_Z80_V3X_HEADER_SIZE = 30 + 2 + 55
)

// Decode [Z80 snapshot] from binary data
func (data SnapshotData) DecodeZ80() (*Z80, error) {
	if len(data) < _Z80_V1_HEADER_SIZE {
		return nil, errors.New("invalid Z80 snapshot")
	}

	PC := uint16(data[6]) | (uint16(data[7]) << 8)

	if PC != 0 {
		// Z80 version 1.xx
		return data.decodeZ80_v1()
	} else {
		if len(data) < _Z80_V2_HEADER_SIZE {
			return nil, errors.New("invalid Z80 snapshot")
		}

		extendedHeaderLength := uint16(data[30]) | (uint16(data[31]) << 8)

		switch _Z80_V1_HEADER_SIZE + 2 + extendedHeaderLength {
		case _Z80_V2_HEADER_SIZE:
			// Z80 version 2.01
			return data.decodeZ80_v2()

		case _Z80_V3_HEADER_SIZE, _Z80_V3X_HEADER_SIZE:
			// Z80 version 3.0x
			return data.decodeZ80_v3()
		}
	}

	return nil, errors.New("invalid Z80 snapshot, or unsupported Z80 snapshot version")
}

func (data SnapshotData) readHeader_v1(s *Z80) error {
	data12 := data[12]
	if data12 == 255 {
		data12 = 1
	}

	s.cpu.A = data[0]
	s.cpu.F = data[1]
	s.cpu.C = data[2]
	s.cpu.B = data[3]
	s.cpu.L = data[4]
	s.cpu.H = data[5]
	s.cpu.PC = uint16(data[6]) | (uint16(data[7]) << 8)
	s.cpu.SP = uint16(data[8]) | (uint16(data[9]) << 8)
	s.cpu.I = data[10]
	s.cpu.R = (data[11] & 0x7f) | ((data12 & 0x01) << 7)

	s.ula.Border = (data12 >> 1) & 0x07

	s.samRom = ((data12 & 0x10) != 0)

	s.cpu.E = data[13]
	s.cpu.D = data[14]
	s.cpu.C_ = data[15]
	s.cpu.B_ = data[16]
	s.cpu.E_ = data[17]
	s.cpu.D_ = data[18]
	s.cpu.L_ = data[19]
	s.cpu.H_ = data[20]
	s.cpu.A_ = data[21]
	s.cpu.F_ = data[22]
	s.cpu.IY = uint16(data[23]) | (uint16(data[24]) << 8)
	s.cpu.IX = uint16(data[25]) | (uint16(data[26]) << 8)

	if data[27] != 0 {
		s.cpu.IFF1 = 1
	} else {
		s.cpu.IFF1 = 0
	}

	if data[28] != 0 {
		s.cpu.IFF2 = 1
	} else {
		s.cpu.IFF2 = 0
	}

	switch IM := (data[29] & 0x03); IM {
	case 0, 1, 2:
		s.cpu.IM = IM
	default:
		return errors.New("invalid interrupt mode")
	}

	s.issue2_emulation = ((data[29] & 0x04) != 0)
	s.doubleInterruptFrequency = ((data[29] & 0x08) != 0)
	s.videoSynchronization = ((data[29] >> 4) & 0x03)
	s.joystick = ((data[29] >> 6) & 0x03)

	if s.samRom {
		return errors.New("unsupported feature: SamRom")
	}
	if s.issue2_emulation {
		return errors.New("unsupported feature: Issue 2 emulation")
	}

	return nil
}

func (data SnapshotData) decodeZ80_v1() (*Z80, error) {
	var s Z80
	var err error

	err = data.readHeader_v1(&s)
	if err != nil {
		return nil, err
	}

	var compressed bool
	{
		data12 := data[12]
		if data12 == 255 {
			data12 = 1
		}

		compressed = ((data12 & 0x20) != 0)
	}

	{
		var mem []byte

		if compressed {
			last0 := data[len(data)-4]
			last1 := data[len(data)-3]
			last2 := data[len(data)-2]
			last3 := data[len(data)-1]
			if !((last0 == 0x00) && (last1 == 0xED) && (last2 == 0xED) && (last3 == 0x00)) {
				return nil, errors.New("invalid Z80 snapshot: no end-marker")
			}

			mem = z80_decompress(data[30 : len(data)-4])
		} else {
			mem = data[30:]
		}

		if len(mem) != 48*1024 {
			return nil, errors.New("invalid Z80 snapshot")
		}

		for i := 0; i < (48 * 1024); i++ {
			s.mem[i] = mem[i]
		}
	}

	return &s, nil
}

func (data SnapshotData) decodeZ80_v2() (*Z80, error) {
	var s Z80
	var err error

	err = data.readHeader_v1(&s)
	if err != nil {
		return nil, err
	}

	if len(data) < _Z80_V2_HEADER_SIZE {
		return nil, errors.New("invalid Z80 snapshot")
	}

	extendedHeaderLength := uint16(data[30]) | (uint16(data[31]) << 8)
	if extendedHeaderLength != 23 {
		return nil, errors.New("invalid Z80 snapshot")
	}

	s.cpu.PC = uint16(data[32]) | (uint16(data[33]) << 8)

	hw_mode := data[34]
	switch hw_mode {
	case 0:
		// 48k
	case 1:
		// 48k + If.1
	default:
		return nil, errors.New("read Z80 snapshot version 2.01: unsupported hardware mode")
	}

	// data[35]: no meaning in 48k mode
	// data[36]: no meaning in 48k mode

	var modifyHardware bool = ((data[37] >> 7) != 0)
	if modifyHardware {
		return nil, errors.New("read Z80 snapshot version 2.01: unsupported hardware mode")
	}

	// rest of data[37]: ignored
	// data[38]: ignored
	// data[39..54]: ignored

	// Memory blocks
	{
		i := int(_Z80_V1_HEADER_SIZE + 2 + extendedHeaderLength)
		err = z80_loadMemBlocks(&s, data[i:])
		if err != nil {
			return nil, err
		}
	}

	return &s, nil
}

func (data SnapshotData) decodeZ80_v3() (*Z80, error) {
	var s Z80
	var err error

	err = data.readHeader_v1(&s)
	if err != nil {
		return nil, err
	}

	if len(data) < _Z80_V3X_HEADER_SIZE {
		return nil, errors.New("invalid Z80 snapshot")
	}

	extendedHeaderLength := uint16(data[30]) | (uint16(data[31]) << 8)
	if !((extendedHeaderLength == 54) || (extendedHeaderLength == 55)) {
		return nil, errors.New("invalid Z80 snapshot")
	}

	s.cpu.PC = uint16(data[32]) | (uint16(data[33]) << 8)

	hw_mode := data[34]
	switch hw_mode {
	case 0:
		// 48k
	case 1:
		// 48k + If.1
	default:
		return nil, errors.New("read Z80 snapshot version 3.0x: unsupported hardware mode")
	}

	// data[35]: no meaning in 48k mode
	// data[36]: no meaning in 48k mode

	var modifyHardware bool = ((data[37] >> 7) != 0)
	if modifyHardware {
		return nil, errors.New("read Z80 snapshot version 3.0x: unsupported hardware mode")
	}

	// rest of data[37]: ignored
	// data[38]: ignored
	// data[39..54]: ignored

	tstate_low := uint(data[55]) | (uint(data[56]) << 8)
	tstate_hi := uint(data[57] & 0x03)
	const T4 = TStatesPerFrame / 4
	s.cpu.Tstate = ((tstate_hi-3)%4)*T4 + (T4 - (tstate_low % T4) - 1)

	// data[58]: always ignored

	MGT_rom_paged := (data[59] == 0xff)
	multiface_rom_paged := (data[60] == 0xff)
	rom0_writable := (data[61] == 0xff)
	rom1_writable := (data[62] == 0xff)

	// data[63..72]: ignored
	// data[73..82]: ignored
	// data[83]: ignored
	// data[84]: ignored
	// data[85]: ignored

	if extendedHeaderLength == 55 {
		// data[86]: ignored
	}

	if MGT_rom_paged {
		return nil, errors.New("read Z80 snapshot version 3.0x: unsupported feature: MGT ROM paging")
	}
	if multiface_rom_paged {
		return nil, errors.New("read Z80 snapshot version 3.0x: unsupported feature: Multiface ROM paging")
	}
	if rom0_writable {
		return nil, errors.New("read Z80 snapshot version 3.0x: unsupported feature: RAM 0..8191")
	}
	if rom1_writable {
		return nil, errors.New("read Z80 snapshot version 3.0x: unsupported feature: RAM 8192..16383")
	}

	// Memory blocks
	{
		i := int(_Z80_V1_HEADER_SIZE + 2 + extendedHeaderLength)
		err = z80_loadMemBlocks(&s, data[i:])
		if err != nil {
			return nil, err
		}
	}

	return &s, nil
}

func z80_loadMemBlocks(s *Z80, data []byte) error {
	pages := make(map[byte]([]byte))

	i := 0
	for i+3 <= len(data) {
		length := int(data[i+0]) | (int(data[i+1]) << 8)
		page := data[i+2]

		i += 3

		compressed := true
		if length == 0xFFFF {
			compressed = false
			length = 0x4000
		}

		if !(i+length <= len(data)) {
			return errors.New("invalid Z80 snapshot")
		}

		if !compressed {
			pages[page] = data[i:(i + length)]
		} else {
			decompressedBlock := z80_decompress(data[i:(i + length)])
			pages[page] = decompressedBlock
		}

		i += length
	}

	if i != len(data) {
		return errors.New("invalid Z80 snapshot")
	}

	if len(pages) != 3 {
		return errors.New("invalid Z80 snapshot")
	}

	for page, pageData := range pages {
		var addr, length int

		switch page {
		case 8:
			addr = 0x4000
			length = 0x4000
		case 4:
			addr = 0x8000
			length = 0x4000
		case 5:
			addr = 0xc000
			length = 0x4000
		default:
			return errors.New("invalid Z80 snapshot")
		}

		if len(pageData) != length {
			return errors.New("invalid Z80 snapshot")
		}

		for i := 0; i < length; i++ {
			s.mem[addr+i-0x4000] = pageData[i]
		}
	}

	return nil
}

func z80_decompress(in []byte) []byte {
	// The input is decompressed in 2 phases:
	//  1. Determine output size
	//  2. Decompress

	len_in := len(in)
	i := 0
	j := 0
	for i < len_in {
		if i+4 <= len_in {
			if (in[i+0] == 0xED) && (in[i+1] == 0xED) {
				count := in[i+2]
				j += int(count)
				i += 4
				continue
			}
		}

		i++
		j++
	}

	len_out := j
	out := make([]byte, len_out)

	i = 0
	j = 0
	for i < len_in {
		if i+4 <= len_in {
			if (in[i+0] == 0xED) && (in[i+1] == 0xED) {
				count := in[i+2]
				value := in[i+3]

				for jj := byte(0); jj < count; jj++ {
					out[j] = value
					j++
				}

				i += 4
				continue
			}
		}

		out[j] = in[i]
		i++
		j++
	}

	return out
}

func (s *Z80) CpuState() CpuState {
	return s.cpu
}

func (s *Z80) UlaState() UlaState {
	return s.ula
}

func (s *Z80) Memory() *[48 * 1024]byte {
	return &s.mem
}
