package parser

import (
	"github.com/har07/expat/etree/element"
)

type XML interface {
	ParseWhole(data string) (*element.E, error)
	Close() (*element.E, error)
	Feed(data string) error
	Free()
}

// ParseError contains error info on parsing tree
type ParseError struct {
	Desc   string
	Code   int
	Line   int
	Column int
}
