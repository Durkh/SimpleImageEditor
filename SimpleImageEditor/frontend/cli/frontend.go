package cli

import (
	"SimpleImageEditor/backend/demo"
	"SimpleImageEditor/backend/parser"
	"SimpleImageEditor/common"
	"fmt"
	"os"
	"sync"
)

func Run(q chan<- rune, imPath chan<- common.Info, done chan<- bool, wg *sync.WaitGroup) {
	defer wg.Done()

	if len(os.Args) < 2 {
		Exit("digite os argumentos")
	}

	var (
		args       = os.Args[1:]
		err        error
		multipart  bool
		operations []rune
		imInfo     common.Info
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

			imInfo.Name = args[i+1]

			multipart = true
		case "-N":
			operations = append(operations, 'N')
		case "-F":
			operations = append(operations, 'F')

			if i+1 > len(args) {
				Exit("error: digite o caminho do filtro")
			}

			imInfo.Options = parser.ParseFileConfig(args[i+1])

			multipart = true
		case "-M":
			operations = append(operations, 'M')

			imInfo.Filter, err = parser.ParseMedianfilter(args[i+1])

			multipart = true
		case "-S":
			if imInfo.Options == nil {
				Exit("error: você não está passando um filtro")
			}

			imInfo.Options["sobel"] = true
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

	if imInfo.Name == "" {
		Exit("carregue uma imagem")
	}

	imPath <- imInfo

	for _, v := range operations {
		q <- v
	}

	done <- true

}

func Exit(err string) {
	fmt.Println(err)
	os.Exit(1)
}
