package demo

import (
	"errors"
	cv "gocv.io/x/gocv"
	"golang.org/x/sync/errgroup"
	"image"
	"image/color"
	"strconv"
	"time"
)

func Demo(picPath, corrPath string) (err error) {

	var (
		pic, grey, corr cv.Mat
		res1            cv.Mat
		region          image.Rectangle
		group           errgroup.Group
		rectColor       = color.RGBA{
			R: 5,
			G: 166,
			B: 207,
			A: 255,
		}
	)

	//load images
	group.Go(func() error {
		if pic = cv.IMRead(picPath, cv.IMReadColor); pic.Empty() {
			return errors.New("error: não foi possível ler a imagem")
		}
		return nil
	})

	group.Go(func() error {
		if corr = cv.IMRead(corrPath, cv.IMReadGrayScale); corr.Empty() {
			return errors.New("error: não foi possível ler a imagem")
		}
		return nil
	})

	if err = group.Wait(); err != nil {
		return nil
	}

	////////////////////////////////////

	// convert image to greyscale
	grey = cv.NewMat()
	cv.CvtColor(pic, &grey, cv.ColorBGRAToGray)

	//do the correlation and find the max correlation point
	if region, err = correlate(grey.Clone(), corr); err != nil {
		return err
	}

	//add rectangle to the image
	res1 = pic.Clone()
	cv.Rectangle(&res1, region, rectColor, 2)

	// save the first image
	if ok := cv.IMWrite("demo_correlation_1_"+strconv.Itoa(int(time.Now().Unix()))+".png", res1); !ok {
		return errors.New("error: não foi possível salvar a primeira imagem")
	}

	//remove the correlation
	cv.Rectangle(&grey, region, color.RGBA{}, int(cv.Filled))

	// do a second correlation
	if region, err = correlate(grey, corr); err != nil {
		return err
	}

	//add rectangle to the image
	cv.Rectangle(&pic, region, rectColor, 2)

	// save the second image
	if ok := cv.IMWrite("demo_correlation_2_"+strconv.Itoa(int(time.Now().Unix()))+".png", pic); !ok {
		return errors.New("error: não foi possível salvar a segunda imagem")
	}

	return nil
}

func correlate(im, templ cv.Mat) (image.Rectangle, error) {

	var (
		correlation = cv.NewMat()
		region      image.Rectangle
	)

	cv.MatchTemplate(im, templ, &correlation, cv.TmCcoeffNormed, cv.NewMat())

	_, _, _, maxLoc := cv.MinMaxLoc(correlation)

	region = image.Rect(maxLoc.X, maxLoc.Y, maxLoc.X+templ.Cols(), maxLoc.Y+templ.Rows())

	return region, nil
}
