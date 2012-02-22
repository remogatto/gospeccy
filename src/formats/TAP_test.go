package formats

import (
	"bytes"
	"io/ioutil"
	"path"
)

const testdataDir = "testdata"

var (
	tapCodeFn    = path.Join(testdataDir, "code.tap")
	tapProgramFn = path.Join(testdataDir, "hello.tap")
)

func (t *testSuite) TestReadTAP() {
	data, err := ioutil.ReadFile(tapCodeFn)
	t.Nil(err)
	tap, err := NewTAP(data)
	t.Nil(err)

	headerBlock := tap.blocks[0].(*tapBlockHeader)
	dataBlock := tap.blocks[1].(tapBlockData)

	t.Equal(23, int(tap.Len()))

	t.NotNil(headerBlock)
	t.Equal(byte(TAP_FILE_CODE), headerBlock.tapType)
	t.Equal("ROM       ", headerBlock.filename)
	t.Equal(uint16(2), headerBlock.length)
	t.Equal(byte(TAP_BLOCK_DATA), dataBlock[0])
}

func (t *testSuite) TestReadTAPError() {
	_, err := NewTAP(nil)
	t.NotNil(err)
}

// SAVE "ROM" CODE 0,2
func (t *testSuite) TestReadTAPCodeFile() {
	data, err := ioutil.ReadFile(tapCodeFn)
	t.Nil(err)
	tap, err := NewTAP(data)
	t.Nil(err)

	if !t.Failed() {
		headerBlock := tap.blocks[0].(*tapBlockHeader)
		dataBlock := tap.blocks[1].(tapBlockData)

		t.Equal(byte(TAP_FILE_CODE), headerBlock.tapType)
		t.Equal(uint16(0), headerBlock.par1)
		t.Equal(uint16(0x8000), headerBlock.par2)

		t.Equal(byte(TAP_BLOCK_DATA), dataBlock[0])
		t.Equal(byte(0xf3), dataBlock[1])
		t.Equal(byte(0xaf), dataBlock[2])
		t.Equal(byte(0xa3), dataBlock[3])
	}
}

// 10 PRINT "Hello World"
// SAVE "HELLO"
func (t *testSuite) TestReadTAPProgramFile() {
	data, err := ioutil.ReadFile(tapProgramFn)
	t.Nil(err)
	tap, err := NewTAP(data)
	t.Nil(err)

	if !t.Failed() {
		headerBlock := tap.blocks[0].(*tapBlockHeader)
		dataBlock := tap.blocks[1].(tapBlockData)

		t.Equal(byte(TAP_FILE_PROGRAM), headerBlock.tapType)
		t.Equal(uint16(0x8000), headerBlock.par1)
		t.Equal(uint16(0x14), headerBlock.par2)

		t.Equal(byte(TAP_BLOCK_DATA), dataBlock[0])
		t.True(bytes.Equal([]byte{
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
		},
			dataBlock[1:]))
	}
}

func (t *testSuite) TestReadTAPWithCustomLoader() {
	data, err := ioutil.ReadFile("testdata/fire.tap")
	t.Nil(err)
	tap, err := NewTAP(data)
	t.Nil(err)

	if !t.Failed() {
		headerBlock := tap.blocks[0].(*tapBlockHeader)
		dataBlock := tap.blocks[1].(tapBlockData)

		t.Equal(byte(TAP_FILE_PROGRAM), headerBlock.tapType)
		t.Equal(uint16(0x0a), headerBlock.par1)
		t.Equal(uint16(0x1e), headerBlock.par2)
		t.Equal(byte(TAP_BLOCK_DATA), dataBlock[0])
	}
}

func (t *testSuite) TestNewTAPFromFile() {
	data, err := ioutil.ReadFile(tapCodeFn)
	t.Nil(err)
	tap, err := NewTAP(data)
	t.Nil(err)

	if !t.Failed() {
		headerBlock := tap.blocks[0].(*tapBlockHeader)
		dataBlock := tap.blocks[1].(tapBlockData)

		t.Equal(byte(TAP_FILE_CODE), headerBlock.tapType)
		t.Equal(uint16(0), headerBlock.par1)
		t.Equal(uint16(0x8000), headerBlock.par2)

		t.Equal(byte(TAP_BLOCK_DATA), dataBlock[0])
		t.Equal(byte(0xf3), dataBlock[1])
		t.Equal(byte(0xaf), dataBlock[2])
		t.Equal(byte(0xa3), dataBlock[3])
	}
}

var tap *TAP

func (t *testSuite) Before() {
	data, _ := ioutil.ReadFile(tapProgramFn)
	tap, _ = NewTAP(data)
}

func (t *testSuite) TestTAPAt() {
	t.Equal(byte(0x00), tap.At(0))
	t.Equal(byte(0xff), tap.At(0x13))
}

func (t *testSuite) TestTAPGetBlock() {
	_, ok := tap.GetBlock(0).(*tapBlockHeader)
	t.True(ok)
	_, ok = tap.GetBlock(1).(tapBlockData)
	t.True(ok)
}

func (t *testSuite) TestTAPBlockLen() {
	t.Equal(19, tap.GetBlock(0).Len())
}
