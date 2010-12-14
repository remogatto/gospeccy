package formats

import (
	"bytes"
	"io/ioutil"
	"path"
	"prettytest"
	"testing"
)

const testdataDir = "testdata"

var (
	tapCodeFn    = path.Join(testdataDir, "code.tap")
	tapProgramFn = path.Join(testdataDir, "hello.tap")
)

type tap_suite_t struct {
	prettytest.Suite
}

func (s *tap_suite_t) testReadTAP() {
	data, _ := ioutil.ReadFile(tapCodeFn)
	tap := NewTAP()
	n, err := tap.Read(data)

	headerBlock := tap.blocks.At(0).(*tapBlockHeader)
	dataBlock := tap.blocks.At(1).(tapBlockData)

	s.Nil(err)
	s.Equal(27, n)

	s.NotNil(headerBlock)
	s.Equal(byte(TAP_FILE_CODE), headerBlock.tapType)
	s.Equal("ROM       ", headerBlock.filename)
	s.Equal(uint16(2), headerBlock.length)
	s.Equal(byte(TAP_BLOCK_DATA), dataBlock[0])
}

func (s *tap_suite_t) testReadTAPError() {
	tap := NewTAP()
	_, err := tap.Read(nil)
	s.NotNil(err)
}

// SAVE "ROM" CODE 0,2
func (s *tap_suite_t) testReadTAPCodeFile() {
	data, _ := ioutil.ReadFile(tapCodeFn)
	tap := NewTAP()
	_, err := tap.Read(data)

	s.Nil(err)

	if !s.Failed() {
		headerBlock := tap.blocks.At(0).(*tapBlockHeader)
		dataBlock := tap.blocks.At(1).(tapBlockData)

		s.Equal(byte(TAP_FILE_CODE), headerBlock.tapType)
		s.Equal(uint16(0), headerBlock.par1)
		s.Equal(uint16(0x8000), headerBlock.par2)

		s.Equal(byte(TAP_BLOCK_DATA), dataBlock[0])
		s.Equal(byte(0xf3), dataBlock[1])
		s.Equal(byte(0xaf), dataBlock[2])
		s.Equal(byte(0xa3), dataBlock[3])
	}
}

// 10 PRINT "Hello World"
// SAVE "HELLO"
func (s *tap_suite_t) testReadTAPProgramFile() {
	data, _ := ioutil.ReadFile(tapProgramFn)
	tap := NewTAP()
	_, err := tap.Read(data)

	s.Nil(err)

	if !s.Failed() {
		headerBlock := tap.blocks.At(0).(*tapBlockHeader)
		dataBlock := tap.blocks.At(1).(tapBlockData)

		s.Equal(byte(TAP_FILE_PROGRAM), headerBlock.tapType)
		s.Equal(uint16(0x8000), headerBlock.par1)
		s.Equal(uint16(0x14), headerBlock.par2)

		s.Equal(byte(TAP_BLOCK_DATA), dataBlock[0])
		s.True(bytes.Equal([]byte{
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

func (s *tap_suite_t) testReadTAPWithCustomLoader() {
	data, _ := ioutil.ReadFile("testdata/fire.tap")
	tap := NewTAP()
	_, err := tap.Read(data)

	s.Nil(err)

	if !s.Failed() {
		headerBlock := tap.blocks.At(0).(*tapBlockHeader)
		dataBlock := tap.blocks.At(1).(tapBlockData)

		s.Equal(byte(TAP_FILE_PROGRAM), headerBlock.tapType)
		s.Equal(uint16(0x8000), headerBlock.par1)
		s.Equal(uint16(0x14), headerBlock.par2)

		s.Equal(byte(TAP_BLOCK_DATA), dataBlock[0])
	}
}

func (s *tap_suite_t) testNewTAPFromFile() {
	tap, err := NewTAPFromFile(tapCodeFn)
	s.Nil(err)
	if !s.Failed() {
		headerBlock := tap.blocks.At(0).(*tapBlockHeader)
		dataBlock := tap.blocks.At(1).(tapBlockData)

		s.Equal(byte(TAP_FILE_CODE), headerBlock.tapType)
		s.Equal(uint16(0), headerBlock.par1)
		s.Equal(uint16(0x8000), headerBlock.par2)

		s.Equal(byte(TAP_BLOCK_DATA), dataBlock[0])
		s.Equal(byte(0xf3), dataBlock[1])
		s.Equal(byte(0xaf), dataBlock[2])
		s.Equal(byte(0xa3), dataBlock[3])
	}
}

var tap *TAP

func (s *tap_suite_t) before() {
	data, _ := ioutil.ReadFile(tapProgramFn)
	tap = NewTAP()
	tap.Read(data)
}

func (s *tap_suite_t) testTAPAt() {
	s.Equal(byte(0x00), tap.At(0))
	s.Equal(byte(0xff), tap.At(0x13))
}

func (s *tap_suite_t) testTAPGetBlock() {
	_, ok := tap.GetBlock(0).(*tapBlockHeader)
	s.True(ok)
	_, ok = tap.GetBlock(1).(tapBlockData)
	s.True(ok)
}

func (s *tap_suite_t) testTAPBlockLen() {
	s.Equal(19, tap.GetBlock(0).Len())
}

func TestTAP(t *testing.T) {
	prettytest.Run(
		t,
		new(tap_suite_t),
	)
}

