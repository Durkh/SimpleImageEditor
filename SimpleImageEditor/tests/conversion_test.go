package tests

import (
	image "SimpleImageEditor/image"
	"SimpleImageEditor/pixel"
	"errors"
	"fmt"
	"math"
	"sync"
	"testing"
)

func TestYIQ(T *testing.T) {

	im := image.Image{}
	if err := im.Open("../../Imagens/testpat.1k.color.tif"); err != nil {
		T.Error(err)
		panic(err)
	}

	YIQ, err := im.YIQ()
	if err != nil {
		T.Error(err)
		panic(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(YIQ.Image.Bounds().Max.X)
	for i := 0; i < YIQ.Image.Bounds().Max.X; i++ {

		func() {

			var x = i

			go func() {
				defer wg.Done()

				for y := 0; y < YIQ.Image.Bounds().Max.Y; y++ {
					old := image.YIQModel(im.Image.At(x, y))

					if n, o := YIQ.Image.At(x, y).(pixel.YIQ).Y, old.(pixel.YIQ).Y; n != o {
						T.Error(errors.New(fmt.Sprintf("Y diferente em (%d, %d),\t\tesperado: %f, obtido: %f",
							x, y, o, n)))
					}

					if n, o := YIQ.Image.At(x, y).(pixel.YIQ).I, old.(pixel.YIQ).I; n != o {
						T.Error(errors.New(fmt.Sprintf("I diferente em (%d, %d),\t\tesperado: %f, obtido: %f",
							x, y, o, n)))
					}

					if n, o := YIQ.Image.At(x, y).(pixel.YIQ).Q, old.(pixel.YIQ).Q; n != o {
						T.Error(errors.New(fmt.Sprintf("Q diferente em (%d, %d),\t\tesperado: %f, obtido: %f",
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
		panic(err)
	}

	wg.Add(RGB.Image.Bounds().Max.X)
	for i := 0; i < RGB.Image.Bounds().Max.X; i++ {

		func() {

			var x = i

			go func() {
				defer wg.Done()

				for y := 0; y < RGB.Image.Bounds().Max.Y; y++ {
					oR, oG, oB, _ := im.Image.At(x, y).RGBA()
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
