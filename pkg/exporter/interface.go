/*
Copyright © 2022 Yifan Gu <guyifan1121@gmail.com>

*/

package exporter

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yifan-gu/blueNote/pkg/config"
	"github.com/yifan-gu/blueNote/pkg/model"
)

var registeredExporters map[string]Exporter

type Exporter interface {
	Name() string
	LoadConfigs(cmd *cobra.Command)
	Export(cfg *config.GlobalConfig, book *model.Book) error
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

func LoadConfigs(cmd *cobra.Command) {
	for _, exporter := range registeredExporters {
		exporter.LoadConfigs(cmd)
	}
}
