package cmd

import (
	"os"

	"github.com/yifan-gu/blueNote/pkg/exporter"
	jsonexporter "github.com/yifan-gu/blueNote/pkg/exporter/json"
	"github.com/yifan-gu/blueNote/pkg/exporter/mongodb"
	"github.com/yifan-gu/blueNote/pkg/exporter/orgroam"
	"github.com/yifan-gu/blueNote/pkg/parser"
	jsonparser "github.com/yifan-gu/blueNote/pkg/parser/json"
	"github.com/yifan-gu/blueNote/pkg/parser/kindlehtml"
	"github.com/yifan-gu/blueNote/pkg/parser/kindlemyclippings"
	"github.com/yifan-gu/blueNote/pkg/storage"
	mongodbStore "github.com/yifan-gu/blueNote/pkg/storage/mongodb"
	"github.com/yifan-gu/blueNote/pkg/util"
)

func registerParsers() {
	parser.RegisterParser(&kindlehtml.KindleHTMLParser{})
	parser.RegisterParser(&jsonparser.JSONParser{})
	parser.RegisterParser(&kindlemyclippings.KindleMyClippingsParser{})
}

func registerExporters() {
	exporter.RegisterExporter(&orgroam.OrgRoamExporter{})
	exporter.RegisterExporter(&jsonexporter.JSONExporter{})
	exporter.RegisterExporter(&mongodb.MongoDBExporter{})
}

func registerStorages() {
	storage.RegisterStorage(&mongodbStore.MongoDBStorage{})
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

func printStoragesAndExit() {
	for _, name := range storage.ListStorages() {
		util.Log(name)
	}
	os.Exit(0)
}
