package parser

import (
	"bufio"
	"fmt"
	"os"
)

type config struct {
	reader *bufio.Reader
	Size   uint64
}

func GetConfigReader(filename string) *config {

	var (
		cfg config
		f   *os.File
		err error
	)

	if f, err = os.Open(filename); err != nil {
		fmt.Printf("%v\n", err.Error())
		os.Exit(1)
	}

	defer f.Close()

	cfg.reader = bufio.NewReaderSize(f, 2048)
	info, _ := f.Stat()
	cfg.Size = uint64(info.Size())

	return &cfg
}
