package pixel

import (
	"image/color"
	"math"
	"sync"
)

type YIQ struct {
	Y float64
	I float64
	Q float64
}

func (color YIQ) RGBA() (r, g, b, a uint32) {

	var (
		y, i, q = color.Y, color.I, color.Q
	)

	R := math.Round(y + .956*i + .621*q)
	G := math.Round(y - .272*i - .647*q)
	B := math.Round(y - 1.106*i + 1.703*q)
	a = 0xff

	bounds(&R, &G, &B)

	r |= uint32(R) << 8
	g |= uint32(G) << 8
	b |= uint32(B) << 8
	a |= a << 8

	return
}

func (color YIQ) Negative() color.Color {

	color.Y = 255 - math.Round(color.Y)

	return color
}

func bounds(r, g, b *float64) {
	var (
		bound = func(ch *float64, wg *sync.WaitGroup) {
			defer wg.Done()

			if *ch < 0 {
				*ch = 0
			} else if *ch > 0xff {
				*ch = 0xff
			}
		}

		wg = sync.WaitGroup{}
	)

	wg.Add(3)

	go bound(r, &wg)
	go bound(g, &wg)
	go bound(b, &wg)

	wg.Wait()
}
