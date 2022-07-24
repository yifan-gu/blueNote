/*
Copyright © 2022 Yifan Gu <guyifan1121@gmail.com>

*/

package exporter

import (
	"fmt"
	"log"
	"strings"

	"github.com/yifan-gu/BlueNote/pkg/config"
	"github.com/yifan-gu/BlueNote/pkg/model"
)

var registeredExporters map[string]Exporter

type Exporter interface {
	Name() string
	Export(cfg *config.Config, book *model.Book) error
}

func RegisterExporter(exporter Exporter) {
	name := strings.ToLower(exporter.Name())
	if registeredExporters == nil {
		registeredExporters = make(map[string]Exporter)
	}
	if _, ok := registeredExporters[name]; ok {
		panic(fmt.Sprintf("Name conflict for exporter %q", name))
	}
	registeredExporters[name] = exporter
}

func GetExporter(name string) Exporter {
	name = strings.ToLower(name)
	exporter, ok := registeredExporters[name]
	if !ok {
		log.Fatal(fmt.Errorf("unrecognized exporter type: %q", name))
	}
	return exporter
}
