package egdImage

import (
	"SimpleImageEditor/pixel"
	"errors"
	"golang.org/x/image/tiff"
	im "image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log"
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

	old := i.Image
	bounds := i.Image.Bounds()
	i.Image = NewRGB(bounds.Max.X, bounds.Max.Y)

	convert(i.Image.Bounds(), func(x int, y int) {
		i.Image.(*RGB).Set(x, y, pixel.RGB{C: old.At(x, y).(color.RGBA)})
	})

	i.Name = strings.Split(filepath.Base(path), ".")[0]

	i.PixelFormat = PixelFormatRGB

	return
}

func (i *Image) setter() func(int, int, color.Color) {

	switch i.PixelFormat {
	case PixelFormatYIQ:
		return i.Image.(*YIQ).Set
	case PixelFormatRGB:
		return i.Image.(*RGB).Set
	default:
		log.Fatalf("error on image setter: %v", errors.New("invalid image pixel type"))
	}

	return nil
}

func (i Image) YIQ() (Image, error) {

	if i.PixelFormat == PixelFormatYIQ {
		return Image{}, errors.New("image already in YIQ format")
	}

	var (
		bounds = i.Image.Bounds()

		res = Image{
			Image:       NewYIQ(bounds.Max.X, bounds.Max.Y),
			Name:        i.Name,
			PixelFormat: PixelFormatYIQ,
			ImageFormat: ImageFormatPNG,
		}

		cm = res.Image.ColorModel()
	)

	convert(bounds, func(x int, y int) { res.setter()(x, y, cm.Convert(i.Image.At(x, y))) })

	return res, nil
}

func (i Image) RGB() (Image, error) {

	if i.PixelFormat == PixelFormatRGB {
		return Image{}, errors.New("image already in RGB format")
	}

	var (
		bounds = i.Image.Bounds()

		res = Image{
			Image:       NewRGB(bounds.Max.X, bounds.Max.Y),
			Name:        i.Name,
			PixelFormat: PixelFormatRGB,
			ImageFormat: ImageFormatPNG,
		}

		cm = res.Image.ColorModel()
	)

	convert(bounds, func(x, y int) { res.setter()(x, y, cm.Convert(i.Image.At(x, y))) })

	return res, nil
}

func (i Image) Negative() (Image, error) {

	var (
		bounds = i.Image.Bounds()
		image  im.Image
	)

	switch i.PixelFormat {
	case PixelFormatYIQ:
		image = NewYIQ(bounds.Max.X, bounds.Max.Y)
	case PixelFormatRGB:
		image = NewRGB(bounds.Max.X, bounds.Max.Y)
	}

	var res = Image{
		Image:       image,
		Name:        i.Name,
		PixelFormat: i.PixelFormat,
		ImageFormat: ImageFormatPNG,
	}

	convert(bounds, func(x, y int) { res.setter()(x, y, i.Image.At(x, y).(pixel.Color).Negative()) })

	return res, nil

}

func convert(bounds im.Rectangle, closure func(int, int)) {

	wg := sync.WaitGroup{}
	wg.Add(bounds.Max.Y + 1)

	for j := bounds.Min.Y; j <= bounds.Max.Y; j++ {
		func() {

			var y = j

			go func() {
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
