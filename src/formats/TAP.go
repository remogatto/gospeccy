package formats

import (
	"os"
	"io/ioutil"
	"container/vector"
)

const (
	TAP_FILE_PROGRAM = iota
	TAP_FILE_NUMBER
	TAP_FILE_CHARACTER_ARRAY
	TAP_FILE_CODE
	TAP_BLOCK_HEADER = 0x00
	TAP_BLOCK_DATA = 0xff
)

func joinBytes(h, l byte) uint16 {
	return uint16(l) | (uint16(h) << 8)
}

func checksum(data []byte) bool {
	exp := data[len(data) - 1]
	sum := byte(0)
	for _, v := range data[0:len(data) - 1] { sum ^= v }
	return exp == sum
}

type tapBlock interface {
	blockType() byte
	checksum() bool
	Len() int
	Data() []byte
}

type tapBlockHeader struct {
	data []byte
	tapType byte
	filename string
	length uint16
	par1, par2 uint16
}

func (header *tapBlockHeader) Len() int {
	return len(header.data)
}

func (header *tapBlockHeader) Data() []byte {
	return header.data
}

func (header *tapBlockHeader) blockType() byte {
	return header.data[0]
}

func (header *tapBlockHeader) checksum() bool {
	return checksum(header.data)
}

type tapBlockData []byte

func (data tapBlockData) Len() int {
	return len(data)
}

func (data tapBlockData) Data() []byte {
	return data
}

func (data tapBlockData) blockType() byte {
	return data[0]
}

func (data tapBlockData) checksum() bool {
	return checksum(data)
}

type TAP struct {
	data []byte
	blocks *vector.Vector
}

func NewTAP() *TAP {
	return &TAP{}
}

func NewTAPFromFile(filename string) (*TAP, os.Error) {
	data, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	tap := NewTAP()
	_, err = tap.Read(data)

	if err != nil {
		return nil, err
	}

	return tap, err
}

func (tap *TAP) Len() uint {
	return uint(len(tap.data))
}

func (tap *TAP) At(pos uint) byte {
	return tap.data[pos]
}

func (tap *TAP) GetBlock(pos int) tapBlock {
	return tap.blocks.At(pos).(tapBlock)
}

func (tap *TAP) readBlockHeader(data []byte) tapBlock {
	tap.blocks.Push(new(tapBlockHeader))
	header := tap.blocks.Last().(*tapBlockHeader)

	header.data = data
	header.tapType = data[1]
	header.filename = string(data[2:12])
	header.length = joinBytes(data[13], data[12])
	header.par1 = joinBytes(data[15], data[14])
	header.par2 = joinBytes(data[17], data[16])

	return header
}

func (tap *TAP) readBlockData(data []byte) tapBlock {
	tap.blocks.Push(tapBlockData(data))
	return tap.blocks.Last().(tapBlockData)
}

func (tap *TAP) readBlock(data []byte) (block tapBlock, err os.Error) {
	if data[0] == 0x00 {
		block = tap.readBlockHeader(data)
	} else {
		block = tap.readBlockData(data)
	}
	if !block.checksum() {
		err = os.NewError("Checksum failed")
	}
	return
}

func (tap *TAP) Read(data []byte) (n int, err os.Error) {
	var (
		length, blockLength uint
		pos, nextPos uint
	)

	if len(data) == 0 {
		err = os.NewError("No TAP data to read!")
		return
	}

	tap.blocks = new(vector.Vector)
	length = uint(len(data))

	for pos = 0; pos < length; pos += nextPos {
		blockLength = uint(joinBytes(data[pos+1], data[pos]))

		if blockLength == 0 {
			err = os.NewError("Block size can't be 0")
			n = int(pos)
			return
		}

		pos += 2
		_, err = tap.readBlock(data[pos:pos + blockLength])
		nextPos = blockLength
		n += int(blockLength) + 2
	}

	tap.data = make([]byte, len(data) - (tap.blocks.Len() * 2))
	var c = 0
	
	for _, blk := range *tap.blocks {
		for _, v := range blk.(tapBlock).Data() {
			tap.data[c] = v
			c++
		}
	}

	return
}



