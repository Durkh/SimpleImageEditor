package egdImage

import (
	"SimpleImageEditor/parser"
	"SimpleImageEditor/pixel"
	"errors"
	"golang.org/x/image/tiff"
	"gonum.org/v1/gonum/mat"
	im "image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type pixelFormat byte
type imageFormat byte

const (
	pixelFormatRGB pixelFormat = iota
	pixelFormatYIQ

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

	defer f.Close()

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

	iterate(i.Image.Bounds(), func(x int, y int) {
		i.Image.(*RGB).Set(x, y, pixel.RGB{C: color.RGBAModel.Convert(old.At(x, y)).(color.RGBA)})
	})

	i.Name = strings.Split(filepath.Base(path), ".")[0]

	i.PixelFormat = pixelFormatRGB

	return
}

func (i *Image) setter() func(int, int, color.Color) {

	switch i.PixelFormat {
	case pixelFormatYIQ:
		return i.Image.(*YIQ).Set
	case pixelFormatRGB:
		return i.Image.(*RGB).Set
	default:
		log.Fatalf("error on image setter: %v", errors.New("invalid image pixel type"))
	}

	return nil
}

func (i Image) YIQ() (Image, error) {

	if i.PixelFormat == pixelFormatYIQ {
		return Image{}, errors.New("image already in YIQ format")
	}

	var (
		bounds = i.Image.Bounds()

		res = Image{
			Image:       NewYIQ(bounds.Max.X, bounds.Max.Y),
			Name:        i.Name,
			PixelFormat: pixelFormatYIQ,
			ImageFormat: ImageFormatPNG,
		}

		cm = res.Image.ColorModel()
	)

	iterate(bounds, func(x int, y int) {
		res.setter()(x, y, cm.Convert(i.Image.At(x, y)))
	})

	return res, nil
}

func (i Image) RGB() (Image, error) {

	if i.PixelFormat == pixelFormatRGB {
		return Image{}, errors.New("image already in RGB format")
	}

	var (
		bounds = i.Image.Bounds()

		res = Image{
			Image:       NewRGB(bounds.Max.X, bounds.Max.Y),
			Name:        i.Name,
			PixelFormat: pixelFormatRGB,
			ImageFormat: ImageFormatPNG,
		}

		cm = res.Image.ColorModel()
	)

	iterate(bounds, func(x, y int) {
		res.setter()(x, y, cm.Convert(i.Image.At(x, y)))
	})

	return res, nil
}

func (i Image) Negative() (Image, error) {

	var (
		bounds = i.Image.Bounds()
		image  im.Image
	)

	switch i.PixelFormat {
	case pixelFormatYIQ:
		image = NewYIQ(bounds.Max.X, bounds.Max.Y)
	case pixelFormatRGB:
		image = NewRGB(bounds.Max.X, bounds.Max.Y)
	}

	var res = Image{
		Image:       image,
		Name:        i.Name + "_neg",
		PixelFormat: i.PixelFormat,
		ImageFormat: ImageFormatPNG,
	}

	iterate(bounds, func(x, y int) {
		res.setter()(x, y, i.Image.At(x, y).(pixel.Color).Negative())
	})

	return res, nil

}

func (i Image) Filter(filterArgs map[string]interface{}) (Image, error) {

	var (
		bounds = i.Image.Bounds()

		res = Image{
			Image:       NewRGB(bounds.Max.X, bounds.Max.Y),
			Name:        i.Name + "_filter",
			PixelFormat: pixelFormatRGB,
			ImageFormat: ImageFormatPNG,
		}

		vecPiv = uint(math.Floor(
			float64(filterArgs["filter"].(parser.Filter).Filter.Len()) / float64(2),
		))

		//default function to assign the RGB values
		assignValues = func(res *Image, x, y int, r, g, b float64) {
			res.Image.(*RGB).Set(x, y, pixel.RGB{
				C: color.RGBA{
					R: uint8(math.Round(r)),
					G: uint8(math.Round(g)),
					B: uint8(math.Round(b)),
					A: 0xff,
				},
			})
		}

		histogramWG = sync.WaitGroup{}

		imagefilteredChannels [3][][]uint8
		min, max              [3]uint8

		mutex sync.Mutex
	)

	min[0], min[1], min[2] = 255, 255, 255

	if _, ok := filterArgs["sobel"]; ok {

		histogramWG.Add(3)

		for j := 0; j < 3; j++ {
			func() {
				index := j
				go func() {
					defer histogramWG.Done()
					imagefilteredChannels[index] = make([][]uint8, bounds.Max.X+1)
					for k := range imagefilteredChannels[index] {
						imagefilteredChannels[index][k] = make([]uint8, bounds.Max.Y+1)
					}
				}()
			}()
		}

		histogramWG.Wait()

		//get the values and check the biggest and the lowest
		assignValues = func(res *Image, x, y int, r, g, b float64) {

			// R channel assgnment and min/max check
			imagefilteredChannels[0][x][y] = uint8(math.Round(r))
			mutex.Lock()
			if imagefilteredChannels[0][x][y] < min[0] {
				min[0] = imagefilteredChannels[0][x][y]
			} else if imagefilteredChannels[0][x][y] > max[0] {
				max[0] = imagefilteredChannels[0][x][y]
			}
			mutex.Unlock()
			/////////////////////////////////////////

			// R channel assgnment and min/max check
			imagefilteredChannels[1][x][y] = uint8(math.Round(g))
			mutex.Lock()
			if imagefilteredChannels[1][x][y] < min[1] {
				min[1] = imagefilteredChannels[1][x][y]
			} else if imagefilteredChannels[1][x][y] > max[1] {
				max[1] = imagefilteredChannels[1][x][y]
			}
			mutex.Unlock()
			/////////////////////////////////////////

			// R channel assgnment and min/max check
			imagefilteredChannels[2][x][y] = uint8(math.Round(b))
			mutex.Lock()
			if imagefilteredChannels[2][x][y] < min[2] {
				min[2] = imagefilteredChannels[2][x][y]
			} else if imagefilteredChannels[2][x][y] > max[2] {
				max[2] = imagefilteredChannels[2][x][y]
			}
			mutex.Unlock()
			/////////////////////////////////////////
		}

	}

	iterate(bounds, func(x int, y int) {

		var (
			//get the image vectors to each channel
			vectorR, vectorG, vectorB = makeVector(x, y, bounds,
				uint(filterArgs["pivot"].(parser.Pivot).M), uint(filterArgs["pivot"].(parser.Pivot).N), vecPiv,
				filterArgs["filter"].(parser.Filter),
				func(channels [3][]float64, vecPiv, vecItr uint, xIt, yIt int, vecWG *sync.WaitGroup) {
					defer vecWG.Done()

					channels[0][vecPiv-vecItr] = float64(i.Image.At(xIt, yIt).(pixel.RGB).C.R)
					channels[1][vecPiv-vecItr] = float64(i.Image.At(xIt, yIt).(pixel.RGB).C.G)
					channels[2][vecPiv-vecItr] = float64(i.Image.At(xIt, yIt).(pixel.RGB).C.B)
				})

			r, g, b float64

			//function to apply filter and check for bounds
			apply = func(ch *float64, v *mat.VecDense, applyWG *sync.WaitGroup) {
				defer applyWG.Done()

				*ch = mat.Dot(v, filterArgs["filter"].(parser.Filter).Filter)

				if offset, ok := filterArgs["offset"]; ok {
					*ch += float64(offset.(uint64))
				}

				//bound check
				if *ch < 0 {
					*ch = 0
				} else if *ch > 0xff {
					*ch = 0xff
				}
			}

			applyWG = sync.WaitGroup{}
		)

		//apply filter to each channel
		applyWG.Add(3)

		go apply(&r, vectorR, &applyWG)
		go apply(&g, vectorG, &applyWG)
		go apply(&b, vectorB, &applyWG)

		applyWG.Wait()
		//////////////////////////////

		//assign the new values with bounds
		assignValues(&res, x, y, r, g, b)
		//////////////////////////////
	})

	//histogram expansion calculation
	if _, ok := filterArgs["sobel"]; ok {
		iterate(bounds, func(x int, y int) {
			res.Image.(*RGB).Set(x, y, pixel.RGB{
				C: color.RGBA{
					R: uint8(math.Round(
						float64(imagefilteredChannels[0][x][y]-min[0]) * float64(255) /
							float64(max[0]-min[0]))),
					G: uint8(math.Round(
						float64(imagefilteredChannels[1][x][y]-min[1]) * float64(255) /
							float64(max[1]-min[1]))),
					B: uint8(math.Round(
						float64(imagefilteredChannels[1][x][y]-min[1]) * float64(255) /
							float64(max[1]-min[1]))),
					A: 0xff,
				},
			})
		})
	}

	return res, nil
}

func (i Image) Mean(filter parser.Filter) (Image, error) {

	if i.PixelFormat != pixelFormatYIQ {
		return Image{}, errors.New("image has to be in YIQ")
	}

	var (
		//pivot is the middle element
		//the first element in the array is pivot-pivot and the last one is pivot+pivot
		pivotX = uint(math.Floor(float64(filter.Size.R) / float64(2)))
		pivotY = uint(math.Floor(float64(filter.Size.C) / float64(2)))
		vecPiv = uint(math.Floor(float64(filter.Filter.Len()) / float64(2)))

		bounds = i.Image.Bounds()

		res = Image{
			Image:       NewYIQ(bounds.Max.X, bounds.Max.Y),
			Name:        i.Name + "_mean",
			PixelFormat: pixelFormatYIQ,
			ImageFormat: ImageFormatPNG,
		}
	)

	iterate(bounds, func(x int, y int) {

		var (
			orig    = i.Image.At(x, y).(pixel.YIQ)
			numbers *[]float64
		)

		makeVector(x, y, bounds, pivotX, pivotY, vecPiv, filter,
			func(channels [3][]float64, vecPiv, vecItr uint, xIt, yIt int, wg *sync.WaitGroup) {
				defer wg.Done()
				//get the image numbers; I'm not using the function to assign to the image vectors
				channels[0][vecPiv-vecItr] = i.Image.At(xIt, yIt).(pixel.YIQ).Y
				//get the array address, only needs to be performed once
				numbers = &channels[0]
			})

		sort.Float64s(*numbers)

		res.Image.(*YIQ).Set(x, y, pixel.YIQ{
			Y: (*numbers)[len(*numbers)/2],
			I: orig.I,
			Q: orig.Q,
		})
	})

	return res, nil
}

func iterate(bounds im.Rectangle, closure func(int, int)) {

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

type assigner func([3][]float64, uint, uint, int, int, *sync.WaitGroup)

func makeVector(x, y int, bounds im.Rectangle, pivotX, pivotY, vecPiv uint, filter parser.Filter, assigner assigner) (*mat.VecDense, *mat.VecDense, *mat.VecDense) {

	var (
		//true image values
		imageValues [3][]float64
		//vector iterator, goes from [vecPiv -> -vecPiv]
		vecItr = vecPiv

		arrSizes = filter.Size.R * filter.Size.C

		wg = sync.WaitGroup{}
	)

	imageValues[0] = make([]float64, arrSizes)
	imageValues[1] = make([]float64, arrSizes)
	imageValues[2] = make([]float64, arrSizes)

	for xIt := x - int(pivotX); xIt <= x+int(pivotX); xIt++ {
		for yIt := y - int(pivotY); yIt <= y+int(pivotY); yIt++ {
			//if it goes out of bounds, set the value to zero
			if xIt < bounds.Min.X || xIt > bounds.Max.X || yIt < bounds.Min.Y || yIt > bounds.Max.Y {
				//as the arrays are already zero initialized, I just need to go to the other check

				//change the array iterator
				vecItr--

				continue
			}

			//assign the correct value to the vector
			wg.Add(1)
			go assigner(imageValues, vecPiv, vecItr, xIt, yIt, &wg)

			vecItr--
		}
	}

	wg.Wait()

	return mat.NewVecDense(len(imageValues[0]), imageValues[0]),
		mat.NewVecDense(len(imageValues[1]), imageValues[1]),
		mat.NewVecDense(len(imageValues[2]), imageValues[2])
}

func SaveImage(im Image) error {

	var (
		suffix string
		err    error
	)

	switch im.PixelFormat {
	case pixelFormatYIQ:
		suffix = "_YIQ"
		im, err = im.RGB()
		if err != nil {
			return err
		}
	case pixelFormatRGB:
		suffix = "_RGB"
	}

	f, err := os.Create(im.Name + suffix + "_" + strconv.Itoa(int(time.Now().Unix())) + ".png")
	if err != nil {
		return err
	}

	defer f.Close()

	err = png.Encode(f, im.Image)
	if err != nil {
		return err
	}

	return nil
}
