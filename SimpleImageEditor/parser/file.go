package parser

import (
	"bufio"
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

	}

	cfg.reader = bufio.NewReaderSize(f, 2048)
	info, _ := f.Stat()
	cfg.Size = uint64(info.Size())

	return &cfg
}
