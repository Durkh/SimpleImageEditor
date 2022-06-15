package egdImage

import (
	"SimpleImageEditor/backend/pixel"
	"image"
	"image/color"
)

type RGB struct {
	pixels       [][]pixel.RGB
	xSize, ySize uint32
}

func (i RGB) ColorModel() color.Model {

	return color.ModelFunc(RGBModel)
}

// RGBModel convert function
func RGBModel(c color.Color) color.Color {

	if _, ok := c.(pixel.RGB); ok {
		return c
	}

	return pixel.RGB{C: color.RGBAModel.Convert(c).(color.RGBA)}
}

func (i RGB) Bounds() image.Rectangle {

	return image.Rect(0, 0, int(i.xSize), int(i.ySize))
}

func (i RGB) At(x int, y int) color.Color {

	return i.pixels[y][x]
}

func (i *RGB) Set(x, y int, c color.Color) {

	i.pixels[y][x] = RGBModel(c).(pixel.RGB)
}

func NewRGB(xSize, ySize int) *RGB {

	pixels := make([][]pixel.RGB, ySize+1)
	for i := range pixels {
		pixels[i] = make([]pixel.RGB, xSize+1)
	}

	return &RGB{
		pixels: pixels,
		xSize:  uint32(xSize),
		ySize:  uint32(ySize),
	}
}
