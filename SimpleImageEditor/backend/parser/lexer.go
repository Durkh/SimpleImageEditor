package parser

import (
	"bytes"
	"errors"
	"io"
	"log"
	"regexp"
	"strings"
)

type tokenType uint8

const (
	tokenComment tokenType = iota
	tokenFilter
	tokenPivot
	tokenOffset
	tokenEOF
	tokenError
)

type Token struct {
	tokenType tokenType
	value     string
}

var (
	buffer *config
)

func LexerInit(path string) {

	buffer = GetConfigReader(path)
}

func fetchFragment() string {

	var (
		strBuf     bytes.Buffer
		incomplete = true
		b          []byte
		err        error
	)

	for incomplete {
		if b, incomplete, err = buffer.reader.ReadLine(); err == nil {
			strBuf.Write(b)
		} else if errors.Is(err, io.EOF) {
			break
		} else {
			log.Fatalf("error reading file line: %v", err.Error())
		}
	}

	return strBuf.String()
}

func NextToken() (tok *Token) {

	tok = new(Token)

	var (
		ok  bool
		err error
	)

	fragment := fetchFragment()

	if ok, err = regexp.MatchString(`^#(.*)$`, fragment); err == nil && ok {
		tok.tokenType = tokenComment
	} else if ok, err = regexp.MatchString(`^filter=\[`, fragment); err == nil && ok {

		if !strings.ContainsRune(fragment, ']') {

			var (
				strBuf bytes.Buffer
				end    bool
			)
			strBuf.Write([]byte(fragment))

			for !end {
				nFrag := fetchFragment()
				strBuf.Write([]byte(nFrag))
				if ok, err := regexp.MatchString(`(\s*((\d+(\.\d+)?)|\.\d+)[,|;])+\s*((\d+(\.\d+)?)|\.\d+)([;|\]])$`, strings.TrimSpace(nFrag)); err == nil && !ok {

					tok.tokenType = tokenError
					tok.value = "error on line: \t" + nFrag + "\n invalid filter syntax"
					return
				}
				end = strings.ContainsRune(nFrag, ']')
			}
			fragment = strBuf.String()
		}

		tok.tokenType = tokenFilter

	} else if ok, err = regexp.MatchString(`^pivot=\[\d+,\d+]$`, fragment); err == nil && ok {
		tok.tokenType = tokenPivot
	} else if ok, err = regexp.MatchString(`^offset=\d+$`, fragment); err == nil && ok {
		tok.tokenType = tokenOffset
	} else if fragment == "" {
		tok.tokenType = tokenEOF
	} else {
		tok.tokenType = tokenError
		fragment = "error on line: \t" + fragment
	}

	tok.value = strings.TrimSpace(fragment)

	return
}
