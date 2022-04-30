package parser

import (
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"gonum.org/v1/gonum/mat"
	"log"
	"strconv"
	"strings"
)

func ParseFileConfig(path string) (args map[string]interface{}) {

	LexerInit(path)

	var tok *Token

	for {
		if tok = NextToken(); tok.tokenType == tokenEOF {
			break
		}

		p, err := parseToken(tok)
		if err != nil {
			log.Fatalf("error parsing config file: %v", err.Error())
		}

		for k, v := range p {
			args[k] = v
		}

	}

	var (
		filter, ok1 = args["filter"]
		pivot, ok2  = args["pivot"]
	)

	if !ok1 || !ok2 {
		panic("not enough arguments in config file")
	}

	if piv, fil := pivot.(struct {
		m int
		n int
	}), filter.(Filter).Size; uint64(piv.m) > fil.R || uint64(piv.n) > fil.C {
		panic("filtro inválido")
	}

	return
}

func parseToken(tok *Token) (map[string]interface{}, error) {

	var (
		res = make(map[string]interface{})
		err error
	)

	switch tok.tokenType {
	case tokenComment:
		return nil, nil
	case tokenFilter:
		if res["filter"], err = parseFilter(tok.value); err != nil {
			return nil, err
		}
	case tokenPivot:

		var x, y int

		if x, err = strconv.Atoi(string(tok.value[8])); err != nil {
			return nil, err
		}

		if y, err = strconv.Atoi(string(tok.value[10])); err != nil {
			return nil, err
		}

		res["pivot"] = struct {
			m int
			n int
		}{x, y}
	case tokenOffset:

		var offset int64

		_, err = fmt.Sscanf(tok.value, "offset=%d", &offset)
		if err != nil {
			return nil, err
		}

		res["offset"] = offset
	case tokenError:
		return nil, errors.New(tok.value)
	}

	return res, nil
}

type Filter struct {
	Size struct {
		R uint64
		C uint64
	}
	Filter *mat.VecDense
}

//TODO colocar caso de filtro não quadrado
func parseFilter(f string) (res Filter, err error) {

	rows := strings.Split(strings.TrimSuffix(strings.TrimPrefix(f, "filter=["), "]"), ";")

	res.Size.R = uint64(len(rows))

	var numbers []string

	for i := range rows {
		aux := strings.Split(rows[i], ",")

		if uint64(len(aux)) > res.Size.C {
			res.Size.C = uint64(len(aux))
		}

		numbers = append(numbers, aux...)
	}

	var (
		group         errgroup.Group
		filterNumbers = make([]float64, len(numbers))
	)

	for i := range numbers {
		func() {
			index := i

			group.Go(func() error {

				n, err := strconv.ParseFloat(strings.TrimSpace(numbers[index]), 64)
				if err != nil {
					return err
				}

				filterNumbers[index] = n

				return nil
			})

		}()
	}

	if err := group.Wait(); err != nil {
		return Filter{}, err
	}

	res.Filter = mat.NewVecDense(len(numbers), filterNumbers)

	return
}

func ParseMedianfilter(input string) (res Filter, err error) {
	// format 5x5

	if _, err = fmt.Sscanf(input, "%dx%d", &res.Size.R, &res.Size.C); err != nil {
		return Filter{}, err
	}

	if res.Size.R%2 == 0 || res.Size.C%2 == 0 {
		return Filter{}, errors.New("os filtro de média deve possuir componentes ímpares")
	}

	val := make([]float64, res.Size.R*res.Size.C)

	for i := range val {
		val[i] = 1 / (float64(res.Size.R * res.Size.C))
	}

	res.Filter = mat.NewVecDense(int(res.Size.R*res.Size.C), val)

	return
}
