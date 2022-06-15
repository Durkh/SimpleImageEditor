package tests

import (
	"SimpleImageEditor/backend/image"
	"SimpleImageEditor/backend/pixel"
	"image/color"
	"math/rand"
	"testing"
)

func TestPixelConversion(T *testing.T) {

	var (
		targets  = make([]pixel.YIQ, 5)
		expected = make([]color.RGBA, 5)
	)

	for i := range expected {
		expected[i] = color.RGBA{
			R: uint8(rand.Int() % 0xff),
			G: uint8(rand.Int() % 0xff),
			B: uint8(rand.Int() % 0xff),
			A: 0xff,
		}
	}

	for i := range targets {
		targets[i] = egdImage.YIQModel(expected[i]).(pixel.YIQ)
	}

	for i := range targets {
		r, g, b, _ := targets[i].RGBA()
		if val := uint8(r >> 8); val != expected[i].R {
			T.Errorf("R diferente em %d, \tesperado: %d, obtido: %d", i, expected[i].R, val)
		}
		if val := uint8(g >> 8); val != expected[i].G {
			T.Errorf("G diferente em %d, \tesperado: %d, obtido: %d", i, expected[i].G, val)
		}
		if val := uint8(b >> 8); val != expected[i].B {
			T.Errorf("B diferente em %d, \tesperado: %d, obtido: %d", i, expected[i].B, val)
		}
	}

}

func TestBounds(T *testing.T) {

	var (
		targets = []pixel.YIQ{{-50, .0000000000000005, -771}, {423789, 980713245, 8972345}}
	)

	for i := range targets {
		r, g, b, _ := targets[i].RGBA()
		if r > 0xffff || (r>>8) > 0xff {
			T.Error("erro de arredondamento no R")
		}
		if g > 0xffff || (r>>8) > 0xff {
			T.Error("erro de arredondamento no G")
		}
		if b > 0xffff || (r>>8) > 0xff {
			T.Error("erro de arredondamento no B")
		}
	}

}
