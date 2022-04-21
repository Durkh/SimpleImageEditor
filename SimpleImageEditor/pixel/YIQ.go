package pixelColor

import (
	"image"
	"image/color"
	"math"
)

type YIQ struct {
	Pixels       [][]YIQColor
	xSize, ySize int
}

func (i YIQ) ColorModel() color.Model {

	return color.ModelFunc(YIQModel)
}

// YIQModel convert function
func YIQModel(c color.Color) color.Color {

	if _, ok := c.(YIQColor); ok {
		return c
	}

	var (
		r, g, b, _ = c.RGBA()
		// back to 8bit values and casted to float
		realR, realG, realB = float64(r >> 8), float64(g >> 8), float64(b >> 8)
	)

	//https://en.wikipedia.org/wiki/YIQ#From_RGB_to_YIQ
	yiq := YIQColor{
		Y: .299*realR + .587*realG + .114*realB,
		I: .596*realR - .274*realG - .322*realB,
		Q: .211*realR - .523*realG + .312*realB,
	}

	return yiq
}

func (i YIQ) Bounds() image.Rectangle {

	return image.Rect(0, 0, i.xSize, i.ySize)
}

func (i YIQ) At(x int, y int) color.Color {

	return i.Pixels[y][x]
}

func (i *YIQ) Set(x, y int, c color.Color) {

	i.Pixels[y][x] = YIQModel(c).(YIQColor)
}

func NewYIQ(xSize, ySize int) *YIQ {

	pixels := make([][]YIQColor, ySize+1)
	for i := range pixels {
		pixels[i] = make([]YIQColor, xSize+1)
	}

	return &YIQ{
		Pixels: pixels,
		xSize:  xSize,
		ySize:  ySize,
	}
}

type YIQColor struct {
	Y float64
	I float64
	Q float64
}

func (color YIQColor) RGBA() (r, g, b, a uint32) {

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
