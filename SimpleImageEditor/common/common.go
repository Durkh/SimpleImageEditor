package common

import "SimpleImageEditor/backend/parser"

type Options map[string]interface{}
type Name = string

type Info = struct {
	Name
	parser.Filter
	Options
}
