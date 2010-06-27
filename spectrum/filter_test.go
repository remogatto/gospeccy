package spectrum

import (
	"testing"
)

var (
	src = []uint32{ 1,  2,  3,  4,
                        5,  6,  7,  8,
	                9,  10, 11, 12,
	                13, 14, 15, 16 }

	expected = []uint32{ 1,   1,   2,   2,  3,  3,  4,  4,
	                     1,   1,   2,   2,  3,  3,  4,  4,
	                     5,   5,   6,   6,  7,  7,  8,  8,
	                     5,   5,   6,   6,  7,  7,  8,  8,
	                     9,   9,   10,  10, 11, 11, 12, 12,
	                     9,   9,   10,  10, 11, 11, 12, 12,
	                     13,  13,  14,  14, 15, 15, 16, 16,
	                     13,  13,  14,  14, 15, 15, 16, 16 }

	dst [16*4]uint32
)

type testSurface struct {
	data []uint32
	w, h uint
}

func (s *testSurface) Width() uint {
	return s.w
}

func (s *testSurface) Height() uint {
	return s.h
}

func (s *testSurface) Bpp() uint {
	return 1
}

func (s *testSurface) SizeInBytes() uint {
	return s.Width() * s.Height() * s.Bpp()
}

func (s *testSurface) Size() uint {
	return s.Width() * s.Height()
}

func (s *testSurface) getValueAt(id uint) uint32 {
	return s.data[id]
}

func (s *testSurface) setPixelValue(x, y uint, value uint32) {
	s.data[s.Width() * y + x] = value
}

func (s *testSurface) setValueAt(id uint, value uint32) { }
func (s *testSurface) setPixel(x, y uint, color [3]byte) { }
func (s *testSurface) setPixelAt(id uint, color [3]byte) { }

func TestScale2x(t *testing.T) {
	Scale2x(&testSurface{ data: src, w: 4, h: 4}, &testSurface{ data: &dst, w: 8, h: 8 })

	for i, val := range expected {
		if val != dst[i] {
			t.Errorf("Expected %d at %d but got %d", val, i, dst[i])
		}
	}
}

