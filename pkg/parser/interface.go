package parser

import (
	"fmt"
)

type Parser interface {
	Parse(inputPath string) (*Book, error)
}

const (
	ParserTypeHtmlClippingParser = "htmlclipping"
)

func NewParser(parser string) (Parser, error) {
	switch parser {
	case ParserTypeHtmlClippingParser:
		return &HtmlClippingParser{}, nil
	default:
		return nil, fmt.Errorf("unrecognized parser type: %q", parser)
	}
}
