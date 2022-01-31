package parser

import (
	"fmt"

	"github.com/yifan-gu/klipping2org/pkg/db"
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

type SqlPlanner interface {
	InsertNodeLinkTitleEntry(book *Book, outputPath string) error
	InsertNodeLinkMarkEntry(book *Book, mark *Mark, outputPath string) error
	InsertFileEntry(book *Book, fullpath string) error
	CommitSql() error
}

func NewSqlPlanner(driver db.SqlInterface, updateRoamDB bool) SqlPlanner {
	if updateRoamDB {
		return &sqlPlanner{driver: driver}
	}
	return &dummySqlPlanner{}
}
