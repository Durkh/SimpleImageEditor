package pixel

import "image/color"

type Color interface {
	color.Color
	Negative() color.Color
}

const (
	Max8BitPixelColor = 0xff
)
