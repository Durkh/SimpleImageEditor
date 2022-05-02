package cli

import (
	"SimpleImageEditor/demo"
	image "SimpleImageEditor/image"
	"SimpleImageEditor/parser"
	"fmt"
	"os"
)

func Run() {

	if len(os.Args) < 2 {
		Exit("digite os argumentos")
	}

	var (
		args         = os.Args[1:]
		err          error
		multipart    bool
		operations   []rune
		im           image.Image
		config       map[string]interface{}
		medianFilter parser.Filter
	)

	for i := range args {

		if multipart {
			multipart = false
			continue
		}

		switch args[i] {
		case "-I":
			if i+1 > len(args) {
				Exit("error: digite o caminho da imagem")
			}

			if err := im.Open(args[i+1]); err != nil {
				panic(err)
			}

			multipart = true
		case "-N":
			operations = append(operations, 'N')
		case "-F":
			operations = append(operations, 'F')

			if i+1 > len(args) {
				Exit("error: digite o caminho do filtro")
			}

			config = parser.ParseFileConfig(args[i+1])

			multipart = true
		case "-M":
			operations = append(operations, 'M')

			medianFilter, err = parser.ParseMedianfilter(args[i+1])

			multipart = true
		case "-S":
			if config == nil {
				Exit("error: você não está passando um filtro")
			}

			config["sobel"] = true
		case "YIQ":
			operations = append(operations, 'Y')
		case "RGB":
			operations = append(operations, 'R')

		case "DEMO":
			if i+2 > len(args) {
				Exit("error: digite o caminho da imagem e do template")
			}

			if err = demo.Demo(args[i+1], args[i+2]); err != nil {
				Exit(err.Error())
			}

			os.Exit(0)
		}
	}

	if im.Name == "" {
		Exit("carregue uma imagem")
	}

	for _, v := range operations {
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
			if im, err = im.Mean(medianFilter); err != nil {
				Exit(err.Error())
			}
		case 'F':
			if im, err = im.Filter(config); err != nil {
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
	os.Exit(1)
}
