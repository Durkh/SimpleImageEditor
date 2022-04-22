package pixel

import "image/color"

type Color interface {
	color.Color
	Negative() color.Color
}
