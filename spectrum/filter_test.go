/* 

Copyright (c) 2010 Andrea Fazzi

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
"Software"), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

*/

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

