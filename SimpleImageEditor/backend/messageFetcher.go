package backend

import (
	image "SimpleImageEditor/backend/image"
	"SimpleImageEditor/common"
	"fmt"
	"os"
	"sync"
)

func Run(q <-chan rune, imPath <-chan common.Info, wg *sync.WaitGroup) {
	defer wg.Done()

	var (
		im     image.Image
		err    error
		imInfo = <-imPath
	)

	if err = im.Open(imInfo.Name); err != nil {
		panic(err)
	}

	for v := range q {
		switch v {
		case 'Y':
			if im, err = im.YIQ(); err != nil {
				Exit(err.Error())
			}
		case 'R':
			if im, err = im.RGB(); err != nil {
				Exit(err.Error())
			}
		case 'N':
			if im, err = im.Negative(); err != nil {
				Exit(err.Error())
			}
		case 'M':
			if im, err = im.Median(imInfo.Filter); err != nil {
				Exit(err.Error())
			}
		case 'F':
			if im, err = im.Filter(imInfo.Options); err != nil {
				Exit(err.Error())
			}
		}
	}

	if err = image.SaveImage(im); err != nil {
		Exit(err.Error())
	}
}

func Exit(err string) {
	fmt.Println(err)
	os.Exit(2)
}
