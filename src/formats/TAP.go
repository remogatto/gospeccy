package formats

import "errors"

const (
	TAP_FILE_PROGRAM         = 0
	TAP_FILE_NUMBER_ARRAY    = 1
	TAP_FILE_CHARACTER_ARRAY = 2
	TAP_FILE_CODE            = 3

	TAP_BLOCK_HEADER = 0x00
	TAP_BLOCK_DATA   = 0xff
)

func joinBytes(h, l byte) uint16 {
	return uint16(l) | (uint16(h) << 8)
}

func checksum(data []byte) bool {
	sum := byte(0)
	for _, v := range data {
		sum ^= v
	}
	return sum == 0
}

type tapBlock interface {
	BlockType() byte // Usually returns TAP_BLOCK_HEADER or TAP_BLOCK_DATA
	Len() int        // Same as 'len(Data())'
	Data() []byte
	checksum() bool
}

type tapBlockHeader struct {
	data       []byte
	tapType    byte // Usually, the value is one of TAP_FILE_*
	filename   string
	length     uint16
	par1, par2 uint16
}

func (header *tapBlockHeader) BlockType() byte {
	return header.data[0]
}

func (header *tapBlockHeader) Len() int {
	return len(header.data)
}

func (header *tapBlockHeader) Data() []byte {
	return header.data
}

func (header *tapBlockHeader) checksum() bool {
	return checksum(header.data)
}

type tapBlockData []byte

func (data tapBlockData) BlockType() byte {
	return data[0]
}

func (data tapBlockData) Len() int {
	return len(data)
}

func (data tapBlockData) Data() []byte {
	return data
}

func (data tapBlockData) checksum() bool {
	return checksum(data)
}

type TAP struct {
	data   []byte
	blocks []tapBlock
}

func NewTAP(data []byte) (*TAP, error) {
	tap := &TAP{}

	err := tap.read(data)
	if err != nil {
		return nil, err
	}

	return tap, nil
}

func (tap *TAP) Len() uint {
	return uint(len(tap.data))
}

func (tap *TAP) At(pos uint) byte {
	return tap.data[pos]
}

func (tap *TAP) GetBlock(pos int) tapBlock {
	return tap.blocks[pos]
}

func readBlock_header(data []byte) *tapBlockHeader {
	header := new(tapBlockHeader)

	header.data = data
	header.tapType = data[1]
	header.filename = string(data[2:12])
	header.length = joinBytes(data[13], data[12])
	header.par1 = joinBytes(data[15], data[14])
	header.par2 = joinBytes(data[17], data[16])

	return header
}

func readBlock_data(data []byte) tapBlockData {
	return tapBlockData(data)
}

func (tap *TAP) readBlock(data []byte) (tapBlock, error) {
	var block tapBlock
	if data[0] == TAP_BLOCK_HEADER {
		block = readBlock_header(data)
	} else {
		block = readBlock_data(data)
	}

	if !block.checksum() {
		return nil, errors.New("checksum failed")
	}

	return block, nil
}

func (tap *TAP) read(data []byte) error {
	length := uint(len(data))
	if length == 0 {
		return errors.New("no TAP data to read")
	}

	pos := uint(0)
	for pos != length {
		if !(pos+1 <= length) {
			return errors.New("invalid TAP data")
		}

		blockLength := uint(joinBytes(data[pos+1], data[pos]))
		if blockLength == 0 {
			return errors.New("block size can't be 0")
		}

		pos += 2

		if !(pos+blockLength <= length) {
			return errors.New("invalid TAP data")
		}

		block, err := tap.readBlock(data[pos : pos+blockLength])
		if err != nil {
			return err
		}

		tap.blocks = append(tap.blocks, block)
		pos += blockLength
	}

	tap.data = make([]byte, len(data)-(len(tap.blocks)*2))
	c := 0
	for _, blk := range tap.blocks {
		for _, blk_data := range blk.Data() {
			tap.data[c] = blk_data
			c++
		}
	}
	if c != len(tap.data) {
		panic("assertion failed")
	}

	return nil
}
