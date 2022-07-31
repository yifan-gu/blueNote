package cmd

import (
	"github.com/yifan-gu/blueNote/pkg/exporter"
	"github.com/yifan-gu/blueNote/pkg/exporter/json"
	"github.com/yifan-gu/blueNote/pkg/exporter/orgroam"
	"github.com/yifan-gu/blueNote/pkg/parser"
	"github.com/yifan-gu/blueNote/pkg/parser/kindlehtml"
)

func registerParsers() {
	parser.RegisterParser(&kindlehtml.KindleHTMLParser{})
}

func registerExporters() {
	exporter.RegisterExporter(&orgroam.OrgRoamExporter{})
	exporter.RegisterExporter(&json.JSONExporter{})
}
