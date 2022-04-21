package egdImage

import (
	pixelColor "SimpleImageEditor/pixel"
	"errors"
	"golang.org/x/image/tiff"
	im "image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type pixelFormat byte
type imageFormat byte

const (
	PixelFormatRGB pixelFormat = iota
	PixelFormatYIQ

	ImageFormatPNG imageFormat = iota
	ImageFormatTIF
	ImageFormatJPG
)

type Image struct {
	Image       im.Image
	Name        string
	PixelFormat pixelFormat
	ImageFormat imageFormat
}

func (i *Image) Open(path string) (err error) {

	var (
		f *os.File
	)

	if f, err = os.Open(path); err != nil {
		return err
	}

	switch filepath.Ext(path) {
	case ".tif":
		fallthrough
	case ".tiff":
		i.Image, err = tiff.Decode(f)
		i.ImageFormat = ImageFormatTIF
	case ".png":
		i.Image, err = png.Decode(f)
		i.ImageFormat = ImageFormatPNG
	case ".jpg":
		fallthrough
	case ".jpeg":
		i.Image, err = jpeg.Decode(f)
		i.ImageFormat = ImageFormatJPG
	default:
		return errors.New("unknown format")
	}

	i.Name = strings.Split(filepath.Base(path), ".")[0]

	i.PixelFormat = PixelFormatRGB

	return
}

func (i Image) YIQ() (Image, error) {

	if i.PixelFormat == PixelFormatYIQ {
		return Image{}, errors.New("image already in YIQ format")
	}

	var (
		bounds = i.Image.Bounds()

		res = Image{
			Image:       pixelColor.NewYIQ(bounds.Max.X, bounds.Max.Y),
			Name:        i.Name,
			PixelFormat: PixelFormatYIQ,
			ImageFormat: ImageFormatPNG,
		}

		cm = res.Image.ColorModel()
	)

	convert(bounds, func(x int, y int) { res.Image.(*pixelColor.YIQ).Set(x, y, cm.Convert(i.Image.At(x, y))) })

	return res, nil
}

func (i Image) RGB() (Image, error) {

	if i.PixelFormat == PixelFormatRGB {
		return Image{}, errors.New("image already in RGB format")
	}

	var (
		bounds = i.Image.Bounds()

		res = Image{
			Image:       im.NewRGBA(im.Rect(0, 0, bounds.Max.X, bounds.Max.Y)),
			Name:        i.Name,
			PixelFormat: PixelFormatRGB,
			ImageFormat: ImageFormatPNG,
		}

		cm = res.Image.ColorModel()
	)

	convert(bounds, func(x, y int) { res.Image.(*im.RGBA).Set(x, y, cm.Convert(i.Image.At(x, y))) })

	return res, nil
}

func convert(bounds im.Rectangle, closure func(int, int)) {

	wg := sync.WaitGroup{}
	wg.Add(bounds.Max.Y + 1)

	for j := bounds.Min.Y; j <= bounds.Max.Y; j++ {
		func() {

			var y = j

			func() {
				defer wg.Done()

				for x := bounds.Min.X; x <= bounds.Max.X; x++ {
					closure(x, y)
				}

			}()
		}()
	}

	wg.Wait()

}

func SaveImage(im Image) error {

	f, err := os.Create(im.Name + "_YIQ_" + strconv.Itoa(int(time.Now().Unix())) + ".png")
	if err != nil {
		return err
	}

	err = png.Encode(f, im.Image)
	if err != nil {
		return err
	}

	return nil
}
