package tests

import (
	image "SimpleImageEditor/image"
	"SimpleImageEditor/parser"
	"SimpleImageEditor/pixel"
	"errors"
	"fmt"
	"math"
	"sync"
	"testing"
)

func TestNegative(T *testing.T) {

	im := image.Image{}
	if err := im.Open("../../Imagens/testpat.1k.color.tif"); err != nil {
		T.Error(err)
		panic(err)
	}

	neg, err := im.Negative()
	if err != nil {
		T.Error(err)
		panic(err)
	}
	wg := sync.WaitGroup{}

	wg.Add(neg.Image.Bounds().Max.X)
	for i := 0; i < neg.Image.Bounds().Max.X; i++ {

		func() {

			var x = i

			go func() {
				defer wg.Done()

				for y := 0; y < neg.Image.Bounds().Max.Y; y++ {
					oR, oG, oB, _ := im.Image.At(x, y).RGBA()
					shiftOrigR, shiftOrigG, shiftOrigB := oR>>8, oG>>8, oB>>8
					nR, nG, nB, _ := neg.Image.At(x, y).RGBA()
					shiftNewR, shiftNewG, shiftNewB := nR>>8, nG>>8, nB>>8

					if pixel.Max8BitPixelColor-shiftOrigR != shiftNewR {
						T.Error(errors.New(fmt.Sprintf("R diferente em (%d, %d),\t\tesperado: %d, obtido: %d", x, y, pixel.Max8BitPixelColor-shiftOrigR, shiftNewR)))
					}

					if pixel.Max8BitPixelColor-shiftOrigG != shiftNewG {
						T.Error(errors.New(fmt.Sprintf("G diferente em (%d, %d),\t\tesperado: %d, obtido: %d", x, y, pixel.Max8BitPixelColor-shiftOrigG, shiftNewG)))
					}

					if pixel.Max8BitPixelColor-shiftOrigB != shiftNewB {
						T.Error(errors.New(fmt.Sprintf("B diferente em (%d, %d),\t\tesperado: %d, obtido: %d", x, y, pixel.Max8BitPixelColor-shiftOrigB, shiftNewB)))
					}
				}
			}()

		}()

	}

	wg.Wait()

	YIQ, err := im.YIQ()
	if err != nil {
		T.Error(err)
		panic(err)
	}

	neg, err = YIQ.Negative()
	if err != nil {
		T.Error(err)
		panic(err)
	}

	wg.Add(YIQ.Image.Bounds().Max.X)
	for i := 0; i < YIQ.Image.Bounds().Max.X; i++ {

		func() {

			var x = i

			go func() {
				defer wg.Done()

				for y := 0; y < YIQ.Image.Bounds().Max.Y; y++ {
					old := image.YIQModel(YIQ.Image.At(x, y))
					newI := neg.Image.At(x, y)

					if n, o := uint(math.Round(newI.(pixel.YIQ).Y)), uint(math.Round(pixel.Max8BitPixelColor-old.(pixel.YIQ).Y)); n != o {
						T.Error(errors.New(fmt.Sprintf("Y diferente em (%d, %d),\t\tesperado: %d, obtido: %d",
							x, y, pixel.Max8BitPixelColor-o, n)))
					}

					if n, o := int(math.Round(neg.Image.At(x, y).(pixel.YIQ).I)), int(math.Round(old.(pixel.YIQ).I)); n != o {
						T.Error(errors.New(fmt.Sprintf("I diferente em (%d, %d),\t\tesperado: %d, obtido: %d",
							x, y, pixel.Max8BitPixelColor-o, n)))
					}

					if n, o := int(math.Round(neg.Image.At(x, y).(pixel.YIQ).Q)), int(math.Round(old.(pixel.YIQ).Q)); n != o {
						T.Error(errors.New(fmt.Sprintf("Q diferente em (%d, %d),\t\tesperado: %d, obtido: %d",
							x, y, pixel.Max8BitPixelColor-o, n)))
					}
				}
			}()

		}()

	}

	wg.Wait()

}

func TestMedian(T *testing.T) {

	yiq := image.NewYIQ(10, 10)
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			yiq.Set(i, j, pixel.YIQ{Y: float64(i + j)})
		}
	}

	im := image.Image{Image: yiq}
	im.PixelFormat = image.PixelFormatYIQ
	filter, err := parser.ParseMedianfilter("3x3")
	if err != nil {
		T.Error(err)
	}

	im, err = im.Median(filter)
	if err != nil {
		T.Error(err)
	}

	expected := []float64{0.000000, 1.000000, 2.000000, 3.000000, 4.000000, 5.000000, 6.000000, 7.000000, 8.000000,
		0.000000, 1.000000, 2.000000, 3.000000, 4.000000, 5.000000, 6.000000, 7.000000, 8.000000, 9.000000, 9.000000,
		2.000000, 3.000000, 4.000000, 5.000000, 6.000000, 7.000000, 8.000000, 9.000000, 10.000000, 10.000000, 3.000000,
		4.000000, 5.000000, 6.000000, 7.000000, 8.000000, 9.000000, 10.000000, 11.000000, 11.000000, 4.000000, 5.000000,
		6.000000, 7.000000, 8.000000, 9.000000, 10.000000, 11.000000, 12.000000, 12.000000, 5.000000, 6.000000, 7.000000,
		8.000000, 9.000000, 10.000000, 11.000000, 12.000000, 13.000000, 13.000000, 6.000000, 7.000000, 8.000000, 9.000000,
		10.000000, 11.000000, 12.000000, 13.000000, 14.000000, 14.000000, 7.000000, 8.000000, 9.000000, 10.000000, 11.000000,
		12.000000, 13.000000, 14.000000, 15.000000, 15.000000, 8.000000, 9.000000, 10.000000, 11.000000, 12.000000, 13.000000,
		14.000000, 15.000000, 16.000000, 16.000000, 0.000000, 9.000000, 10.000000, 11.000000, 12.000000, 13.000000, 14.000000,
		15.000000, 16.000000, 0.000000}

	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			if im.Image.At(i, j).(pixel.YIQ).Y != expected[i*10+j] {
				T.Errorf("expected: %f; got: %f", expected[i*10+j], im.Image.At(i, j).(pixel.YIQ).Y)
			}
		}
	}
}
