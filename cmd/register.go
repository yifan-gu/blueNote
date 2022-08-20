package cmd

import (
	"os"

	"github.com/yifan-gu/blueNote/pkg/exporter"
	jsonexporter "github.com/yifan-gu/blueNote/pkg/exporter/json"
	"github.com/yifan-gu/blueNote/pkg/exporter/orgroam"
	"github.com/yifan-gu/blueNote/pkg/parser"
	jsonparser "github.com/yifan-gu/blueNote/pkg/parser/json"
	"github.com/yifan-gu/blueNote/pkg/parser/kindlehtml"
	"github.com/yifan-gu/blueNote/pkg/util"
)

func registerParsers() {
	parser.RegisterParser(&kindlehtml.KindleHTMLParser{})
	parser.RegisterParser(&jsonparser.JSONParser{})
}

func registerExporters() {
	exporter.RegisterExporter(&orgroam.OrgRoamExporter{})
	exporter.RegisterExporter(&jsonexporter.JSONExporter{})
}

func printParsersAndExit() {
	for _, name := range parser.ListParsers() {
		util.Log(name)
	}
	os.Exit(0)
}

func printExportersAndExit() {
	for _, name := range exporter.ListExporters() {
		util.Log(name)
	}
	os.Exit(0)
}
