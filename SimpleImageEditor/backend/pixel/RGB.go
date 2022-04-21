package pixel

import "image/color"

type RGB struct {
	C color.RGBA
}

func (p RGB) Negative() color.Color {

	p.C.R = 255 - p.C.R
	p.C.G = 255 - p.C.G
	p.C.B = 255 - p.C.B

	return p
}

func (p RGB) RGBA() (r, g, b, a uint32) {

	return p.C.RGBA()
}
