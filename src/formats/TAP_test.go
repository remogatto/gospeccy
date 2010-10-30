package formats

import (
	"testing"
	"spectrum/prettytest"
	"path"
	"io/ioutil"
	"bytes"
)

const testdataDir = "testdata"

var (
	tapCodeFn = path.Join(testdataDir, "code.tap")
	tapProgramFn = path.Join(testdataDir, "hello.tap")
)

func testReadTAP(assert *prettytest.T) {
	data, _ := ioutil.ReadFile(tapCodeFn)
	tap := NewTAP()
	n, err := tap.Read(data)

	headerBlock := tap.blocks.At(0).(*tapBlockHeader)
	dataBlock := tap.blocks.At(1).(tapBlockData)

	assert.Nil(err)
	assert.Equal(27, n)

	assert.NotNil(headerBlock)
	assert.Equal(byte(TAP_FILE_CODE), headerBlock.tapType)
	assert.Equal("ROM       ", headerBlock.filename)
	assert.Equal(uint16(2), headerBlock.length)
	assert.Equal(byte(TAP_BLOCK_DATA), dataBlock[0])
}

func testReadTAPError(assert *prettytest.T) {
	tap := NewTAP()
	_, err := tap.Read(nil)
	assert.NotNil(err)
}

// SAVE "ROM" CODE 0,2
func testReadTAPCodeFile(assert *prettytest.T) {
	data, _ := ioutil.ReadFile(tapCodeFn)
	tap := NewTAP()
	_, err := tap.Read(data)

	assert.Nil(err)

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

	assert.Nil(err)
	
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

func testReadTAPWithCustomLoader(assert *prettytest.T) {
	data, _ := ioutil.ReadFile("testdata/fire.tap")
	tap := NewTAP()
	_, err := tap.Read(data)

	assert.Nil(err)
	
	if !assert.Failed() {
		headerBlock := tap.blocks.At(0).(*tapBlockHeader)
		dataBlock := tap.blocks.At(1).(tapBlockData)

		assert.Equal(byte(TAP_FILE_PROGRAM), headerBlock.tapType)
		assert.Equal(uint16(0x8000), headerBlock.par1)
		assert.Equal(uint16(0x14), headerBlock.par2)

		assert.Equal(byte(TAP_BLOCK_DATA), dataBlock[0])
	}
}

func TestLoadTAP(t *testing.T) {
	prettytest.Run(
		t,
		testReadTAP,
		testReadTAPError,
		testReadTAPCodeFile,
		testReadTAPProgramFile,
		testReadTAPWithCustomLoader,
	)
}

var tap *TAP

func before(assert *prettytest.T) {
	data, _ := ioutil.ReadFile(tapProgramFn)
	tap = NewTAP()
	tap.Read(data)
}

func testTAPAt(assert *prettytest.T) {
	assert.Equal(byte(0x00), tap.At(0))
	assert.Equal(byte(0xff), tap.At(0x13))
}

func testTAPGetBlock(assert *prettytest.T) {
	_, ok := tap.GetBlock(0).(*tapBlockHeader)
	assert.True(ok)
	_, ok = tap.GetBlock(1).(tapBlockData)
	assert.True(ok)
}

func testTAPBlockLen(assert *prettytest.T) {
	assert.Equal(19, tap.GetBlock(0).Len())
}

func TestTAPAccessors(t *testing.T) {
	prettytest.Run(
		t,
		before,
		testTAPAt,
		testTAPBlockLen,
		testTAPGetBlock,
	)
		
}

