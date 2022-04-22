package egdImage

import (
	"SimpleImageEditor/pixel"
	"image"
	"image/color"
)

type YIQ struct {
	pixels       [][]pixel.YIQ
	xSize, ySize uint32
}

func (i YIQ) ColorModel() color.Model {

	return color.ModelFunc(YIQModel)
}

// YIQModel convert function
func YIQModel(c color.Color) color.Color {

	if _, ok := c.(pixel.YIQ); ok {
		return c
	}

	var (
		r, g, b, _ = c.RGBA()
		// back to 8bit values and cast to float
		realR, realG, realB = float64(r >> 8), float64(g >> 8), float64(b >> 8)
	)

	//https://en.wikipedia.org/wiki/YIQ#From_RGB_to_YIQ
	yiq := pixel.YIQ{
		Y: .299*realR + .587*realG + .114*realB,
		I: .596*realR - .274*realG - .322*realB,
		Q: .211*realR - .523*realG + .312*realB,
	}

	return yiq
}

func (i YIQ) Bounds() image.Rectangle {

	return image.Rect(0, 0, int(i.xSize), int(i.ySize))
}

func (i YIQ) At(x int, y int) color.Color {

	return i.pixels[y][x]
}

func (i *YIQ) Set(x, y int, c color.Color) {

	i.pixels[y][x] = YIQModel(c).(pixel.YIQ)
}

func NewYIQ(xSize, ySize int) *YIQ {

	pixels := make([][]pixel.YIQ, ySize+1)
	for i := range pixels {
		pixels[i] = make([]pixel.YIQ, xSize+1)
	}

	return &YIQ{
		pixels: pixels,
		xSize:  uint32(xSize),
		ySize:  uint32(ySize),
	}
}
