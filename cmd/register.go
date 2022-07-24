package cmd

import (
	"github.com/yifan-gu/BlueNote/pkg/exporter"
	"github.com/yifan-gu/BlueNote/pkg/exporter/orgroam"
	"github.com/yifan-gu/BlueNote/pkg/parser"
	"github.com/yifan-gu/BlueNote/pkg/parser/kindlehtml"
)

func registerParsers() {
	parser.RegisterParser(&kindlehtml.KindleHTMLParser{})
}

func registerExporters() {
	exporter.RegisterExporter(&orgroam.OrgRoamExporter{})
}
