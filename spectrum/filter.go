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

// Quick hack to double the size of the speccy screen.
func Scale2x(src SurfaceAccessor, dst SurfaceAccessor) {
	var srcSize = src.SizeInBytes()
	var halfSrcSize = srcSize / 2
	
	ch := make(chan bool, 2)

	partialScale(src, dst, 0, halfSrcSize, ch)
	partialScale(src, dst, halfSrcSize, srcSize, ch)
	
	<-ch
	<-ch
}


