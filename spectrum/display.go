package spectrum

var palette [16][3]byte = [16][3]byte{
	[3]byte{0, 0, 0},
	[3]byte{0, 0, 192},
	[3]byte{192, 0, 0},
	[3]byte{192, 0, 192},
	[3]byte{0, 192, 0},
	[3]byte{0, 192, 192},
	[3]byte{192, 192, 0},
	[3]byte{192, 192, 192},
	[3]byte{0, 0, 0},
	[3]byte{0, 0, 255},
	[3]byte{255, 0, 0},
	[3]byte{255, 0, 255},
	[3]byte{0, 255, 0},
	[3]byte{0, 255, 255},
	[3]byte{255, 255, 0},
	[3]byte{255, 255, 255}}

type SurfaceAccessor interface {
	Width() uint
	Height() uint
	SizeInBytes() uint
	Bpp() uint

	getValueAt(id uint) uint32
	setValueAt(id uint, value uint32)

	setPixelAt(address uint, color [3]byte)
	setPixelValue(x, y uint, value uint32)
	setPixel(x, y uint, color [3]byte)
}


type DisplayAccessor interface {
	setPixel(x, y uint, color [3]byte)
	setPixelAt(address uint, color [3]byte)
	setBorderColor(color [3]byte)
	flush()
}

