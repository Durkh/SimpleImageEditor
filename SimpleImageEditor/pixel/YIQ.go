package pixel

import (
	"image/color"
	"math"
)

type YIQ struct {
	Y float64
	I float64
	Q float64
}

func (color YIQ) RGBA() (r, g, b, a uint32) {

	var y, i, q = color.Y, color.I, color.Q

	r = uint32(math.Round(y + .956*i + .621*q))
	r |= r << 8
	g = uint32(math.Round(y - .272*i - .647*q))
	g |= g << 8
	b = uint32(math.Round(y - 1.106*i + 1.703*q))
	b |= b << 8
	a = 0xff
	a |= a << 8

	return
}

func (color YIQ) Negative() color.Color {

	color.Y = 255 - color.Y

	return color
}
