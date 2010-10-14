package formats

import (
	"testing"
	"spectrum/prettytest"
	"os"
	"path"
	"io/ioutil"
	"container/vector"
	"bytes"
)

const testdataDir = "testdata"

var (
	tapCodeFn = path.Join(testdataDir, "code.tap")
	tapProgramFn = path.Join(testdataDir, "hello.tap")
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

func checksum(exp byte, data []byte) bool {
	sum := data[0]
	for _, v := range data { sum ^= v }
	return exp == sum
}

type tapBlock interface {
	blockType() byte
	data() []byte
}

type tapBlockHeader struct {
	data []byte
	tapType byte
	filename string
	length uint16
	par1, par2 uint16
}

func (header *tapBlockHeader) blockType() byte {
	return header.data[0]
}

type tapBlockData []byte

func (data tapBlockData) blockType() byte {
	return data[0]
}

type TAP struct {
	blocks *vector.Vector
}

func NewTAP() *TAP {
	return &TAP{}
}

func (tap *TAP) readBlockHeader(data []byte) {
	tap.blocks.Push(new(tapBlockHeader))
	header := tap.blocks.Last().(*tapBlockHeader)

	header.data = data
	header.tapType = data[1]
	header.filename = string(data[2:12])
	header.length = joinBytes(data[13], data[12])
	header.par1 = joinBytes(data[15], data[14])
	header.par2 = joinBytes(data[17], data[16])
}

func (tap *TAP) readBlockData(data []byte) {
	tap.blocks.Push(tapBlockData(data))
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

	for pos = 0; pos < length; pos+=nextPos {
		blockLength = uint(joinBytes(data[pos+1], data[pos]))

		if blockLength == 0 {
			err = os.NewError("Block size can't be 0")
			n = int(pos)
			return
		}

		pos += 2
		blockData := data[pos:pos + blockLength]
		blockType := data[pos]

		switch blockType {
		case 0x00:
			tap.readBlockHeader(blockData)
			checksum(data[blockLength - 1], blockData)
			nextPos += blockLength
			n += int(blockLength) + 2
		case 0xff:
			tap.readBlockData(blockData)
			nextPos += blockLength
			n += int(blockLength) + 2
		}

	}

	return
}

func testReadTAP(assert *prettytest.T) {
	data, _ := ioutil.ReadFile(tapCodeFn)
	tap := NewTAP()
	n, err := tap.Read(data)

	headerBlock := tap.blocks.At(0).(*tapBlockHeader)
	dataBlock := tap.blocks.At(1).(tapBlockData)

	assert.True(err == nil)
	assert.Equal(27, n)

	assert.True(headerBlock != nil)
	assert.Equal(byte(TAP_FILE_CODE), headerBlock.tapType)
	assert.Equal("ROM       ", headerBlock.filename)
	assert.Equal(uint16(2), headerBlock.length)
	assert.Equal(byte(TAP_BLOCK_DATA), dataBlock[0])
}

func testReadTAPError(assert *prettytest.T) {
	tap := NewTAP()
	_, err := tap.Read(nil)
	assert.False(err == nil)
}

// SAVE "ROM" CODE 0,2
func testReadTAPCodeFile(assert *prettytest.T) {
	data, _ := ioutil.ReadFile(tapCodeFn)
	tap := NewTAP()
	_, err := tap.Read(data)

	assert.True(err == nil)

	if !assert.Failed() {
		headerBlock := tap.blocks.At(0).(*tapBlockHeader)
		dataBlock := tap.blocks.At(1).(tapBlockData)

		assert.Equal(byte(TAP_FILE_CODE), headerBlock.tapType)
		assert.Equal(uint16(0), headerBlock.par1)
		assert.Equal(uint16(0x8000), headerBlock.par2)

		assert.Equal(byte(TAP_BLOCK_DATA), dataBlock[0])
		assert.Equal(byte(0xf3), dataBlock[1])
		assert.Equal(byte(0xaf), dataBlock[2])
		assert.Equal(byte(0xa3), dataBlock[3])
	}
}

// 10 PRINT "Hello World"
// SAVE "HELLO"
func testReadTAPProgramFile(assert *prettytest.T) {
	data, _ := ioutil.ReadFile(tapProgramFn)
	tap := NewTAP()
	_, err := tap.Read(data)

	assert.True(err == nil)
	
	if !assert.Failed() {
		headerBlock := tap.blocks.At(0).(*tapBlockHeader)
		dataBlock := tap.blocks.At(1).(tapBlockData)

		assert.Equal(byte(TAP_FILE_PROGRAM), headerBlock.tapType)
		assert.Equal(uint16(0x8000), headerBlock.par1)
		assert.Equal(uint16(0x14), headerBlock.par2)

		assert.Equal(byte(TAP_BLOCK_DATA), dataBlock[0])
		assert.True(bytes.Equal([]byte{
			0x00, 0x0a, 
			0x10, 0x00,
			0x20, 0xf5,
			0x22, 0x48,
			0x65, 0x6c,
			0x6c, 0x6f,
			0x20, 0x57,
			0x6f, 0x72,
			0x6c, 0x64,
			0x22, 0x0d,
			0x1d,
		}, dataBlock[1:]))
	}
}

func TestLoadTAP(t *testing.T) {
	prettytest.Run(
		t,
		testReadTAP,
		testReadTAPError,
		testReadTAPCodeFile,
		testReadTAPProgramFile,
	)
}
