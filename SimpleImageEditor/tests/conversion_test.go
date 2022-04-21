package tests

import (
	image "SimpleImageEditor/image"
	pixelColor "SimpleImageEditor/pixel"
	"errors"
	"fmt"
	"golang.org/x/image/tiff"
	"os"
	"sync"
	"testing"
)

func TestYIQ(T *testing.T) {

	f, err := os.Open("../../Imagens/testpat.1k.color.tif")
	if err != nil {
		T.Error(err)
	}

	orig, err := tiff.Decode(f)
	if err != nil {
		T.Error(err)
	}

	im := image.Image{
		Image:       orig,
		PixelFormat: image.PixelFormatRGB,
	}

	YIQ, err := im.YIQ()
	if err != nil {
		T.Error(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(YIQ.Image.Bounds().Max.X)
	for i := 0; i < YIQ.Image.Bounds().Max.X; i++ {

		func() {

			var x = i

			go func() {
				defer wg.Done()

				for y := 0; y < YIQ.Image.Bounds().Max.Y; y++ {
					old := pixelColor.YIQModel(orig.At(x, y))

					if n, o := YIQ.Image.At(x, y).(pixelColor.YIQColor).Y, old.(pixelColor.YIQColor).Y; n != o {
						T.Error(errors.New(fmt.Sprintf("R diferente em (%d, %d),\t\tesperado: %f, obtido: %f",
							x, y, o, n)))
					}

					if n, o := YIQ.Image.At(x, y).(pixelColor.YIQColor).I, old.(pixelColor.YIQColor).I; n != o {
						T.Error(errors.New(fmt.Sprintf("G diferente em (%d, %d),\t\tesperado: %f, obtido: %f",
							x, y, o, n)))
					}

					if n, o := YIQ.Image.At(x, y).(pixelColor.YIQColor).Q, old.(pixelColor.YIQColor).Q; n != o {
						T.Error(errors.New(fmt.Sprintf("B diferente em (%d, %d),\t\tesperado: %f, obtido: %f",
							x, y, o, n)))
					}
				}
			}()

		}()

	}

	wg.Wait()

	RGB, err := YIQ.RGB()
	if err != nil {
		T.Error(err)
	}

	wg.Add(RGB.Image.Bounds().Max.X)
	for i := 0; i < RGB.Image.Bounds().Max.X; i++ {

		func() {

			var x = i

			func() {
				defer wg.Done()

				for y := 0; y < RGB.Image.Bounds().Max.Y; y++ {
					oR, oG, oB, _ := orig.At(x, y).RGBA()
					shiftOrigR, shiftOrigG, shiftOrigB := oR>>8, oG>>8, oB>>8
					nR, nG, nB, _ := RGB.Image.At(x, y).RGBA()
					shiftNewR, shiftNewG, shiftNewB := nR>>8, nG>>8, nB>>8

					if shiftOrigR != shiftNewR {
						T.Error(errors.New(fmt.Sprintf("R diferente em (%d, %d),\t\tesperado: %d, obtido: %d", x, y, shiftOrigR, shiftNewR)))
					}

					if shiftOrigG != shiftNewG {
						T.Error(errors.New(fmt.Sprintf("G diferente em (%d, %d),\t\tesperado: %d, obtido: %d", x, y, shiftOrigG, shiftNewG)))
					}

					if shiftOrigB != shiftNewB {
						T.Error(errors.New(fmt.Sprintf("B diferente em (%d, %d),\t\tesperado: %d, obtido: %d", x, y, shiftOrigB, shiftNewB)))
					}
				}
			}()

		}()

	}

	wg.Wait()
}
