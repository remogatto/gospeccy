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

func partialScale(src SurfaceAccessor, dst SurfaceAccessor, startId, endId uint, ch chan bool) {
	var i, x, y, w, bpp, xDest, yDest uint
	var srcValue uint32

	bpp = src.Bpp()
	w = src.Width()
	
	y = startId / (w * bpp)

	for i = startId; i < endId; i += bpp {

		if x == w {
			x =  0
			y++
		}

		xDest = x * 2
		yDest = y * 2
		
		srcValue = src.getValueAt(i)

		dst.setPixelValue(xDest, yDest, srcValue)
		dst.setPixelValue(xDest + 1, yDest, srcValue)
		dst.setPixelValue(xDest, yDest + 1, srcValue)
		dst.setPixelValue(xDest + 1, yDest + 1, srcValue)

		x++

	}

	ch<-true

}

// Quick hack to double the size of the speccy screen. Try to run the
// executable with GOMAXPROCS>1 and see if it scales.
func Scale2x(src SurfaceAccessor, dst SurfaceAccessor) {
	var srcSize = src.SizeInBytes()
	var halfSrcSize = srcSize / 2
	
	ch := make(chan bool, 2)

	partialScale(src, dst, 0, halfSrcSize, ch)
	partialScale(src, dst, halfSrcSize, srcSize, ch)
	
	<-ch
	<-ch
}


