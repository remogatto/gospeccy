package formats

import "errors"

type SNA struct {
	cpu CpuState
	ula UlaState
	mem [48 * 1024]byte
}

// Decode SNA from binary data
func (data SnapshotData) DecodeSNA() (*SNA, error) {
	if len(data) != 49179 {
		return nil, errors.New("snapshot has invalid size")
	}

	var s SNA

	// Populate registers
	s.cpu.I = data[0]
	s.cpu.L_ = data[1]
	s.cpu.H_ = data[2]
	s.cpu.E_ = data[3]
	s.cpu.D_ = data[4]
	s.cpu.C_ = data[5]
	s.cpu.B_ = data[6]
	s.cpu.F_ = data[7]
	s.cpu.A_ = data[8]
	s.cpu.L = data[9]
	s.cpu.H = data[10]
	s.cpu.E = data[11]
	s.cpu.D = data[12]
	s.cpu.C = data[13]
	s.cpu.B = data[14]
	s.cpu.IY = uint16(data[15]) | (uint16(data[16]) << 8)
	s.cpu.IX = uint16(data[17]) | (uint16(data[18]) << 8)

	if (data[19] & 0x04) != 0 {
		s.cpu.IFF1 = 1
	} else {
		s.cpu.IFF1 = 0
	}
	s.cpu.IFF2 = s.cpu.IFF1

	s.cpu.R = data[20]

	s.cpu.F = data[21]
	s.cpu.A = data[22]
	s.cpu.SP = uint16(data[23]) | (uint16(data[24]) << 8)

	switch IM := data[25]; IM {
	case 0, 1, 2:
		s.cpu.IM = IM
	default:
		return nil, errors.New("invalid interrupt mode")
	}

	s.ula.Border = data[26] & 0x07

	for i := 0; i < 0xc000; i++ {
		s.mem[i] = data[i+27]
	}

	s.cpu.Tstate = 0

	// Start by executing RETN at address 0x72 in ROM
	s.cpu.PC = 0x72

	return &s, nil
}

// Turn snapshot into binary data (SNA format)
func (s *FullSnapshot) EncodeSNA() ([]byte, error) {
	var data [49179]byte

	// Save registers
	data[0] = s.Cpu.I
	data[1] = s.Cpu.L_
	data[2] = s.Cpu.H_
	data[3] = s.Cpu.E_
	data[4] = s.Cpu.D_
	data[5] = s.Cpu.C_
	data[6] = s.Cpu.B_
	data[7] = s.Cpu.F_
	data[8] = s.Cpu.A_
	data[9] = s.Cpu.L
	data[10] = s.Cpu.H
	data[11] = s.Cpu.E
	data[12] = s.Cpu.D
	data[13] = s.Cpu.C
	data[14] = s.Cpu.B
	data[15] = byte(s.Cpu.IY & 0xff)
	data[16] = byte(s.Cpu.IY >> 8)
	data[17] = byte(s.Cpu.IX & 0xff)
	data[18] = byte(s.Cpu.IX >> 8)

	if s.Cpu.IFF1 != 0 {
		data[19] = (1 << 2)
	} else {
		data[19] = (0 << 2)
	}

	data[20] = s.Cpu.R

	data[21] = s.Cpu.F
	data[22] = s.Cpu.A

	sp_afterSimulatedPushPC := s.Cpu.SP - 2
	if (sp_afterSimulatedPushPC < 0x4000) || (sp_afterSimulatedPushPC > 0xfffe) {
		// We would be saving the PC to ROM or outside of memory
		return nil, errors.New("failed to simulate a RETN")
	}

	data[23] = byte(sp_afterSimulatedPushPC & 0xff)
	data[24] = byte(sp_afterSimulatedPushPC >> 8)
	data[25] = s.Cpu.IM

	// Border color
	data[26] = s.Ula.Border & 0x07

	// Memory
	for i := 0; i < 0xc000; i++ {
		data[i+27] = s.Mem[i]
	}

	// Push PC
	pch, pcl := splitWord(s.Cpu.PC)
	data[(sp_afterSimulatedPushPC-0x4000+0)+27] = pcl
	data[(sp_afterSimulatedPushPC-0x4000+1)+27] = pch

	return data[:], nil
}

func (s *SNA) CpuState() CpuState {
	return s.cpu
}

func (s *SNA) UlaState() UlaState {
	return s.ula
}

func (s *SNA) Memory() *[48 * 1024]byte {
	return &s.mem
}
