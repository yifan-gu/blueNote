/*
Copyright © 2022 Yifan Gu <guyifan1121@gmail.com>

*/

package parser

import (
	"fmt"
	"log"
	"strings"

	"github.com/yifan-gu/BlueNote/pkg/model"
)

var registeredParsers map[string]Parser

type Parser interface {
	Name() string
	Parse(inputPath string) (*model.Book, error)
}

func RegisterParser(parser Parser) {
	name := strings.ToLower(parser.Name())
	if registeredParsers == nil {
		registeredParsers = make(map[string]Parser)
	}
	if _, ok := registeredParsers[name]; ok {
		panic(fmt.Sprintf("Name conflict for parser %q", name))
	}
	registeredParsers[name] = parser
}

func GetParser(name string) Parser {
	name = strings.ToLower(name)
	parser, ok := registeredParsers[name]
	if !ok {
		log.Fatal(fmt.Errorf("unrecognized parser type: %q", name))
	}
	return parser
}
