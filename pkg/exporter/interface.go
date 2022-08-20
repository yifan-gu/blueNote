/*
Copyright © 2022 Yifan Gu <guyifan1121@gmail.com>

*/

package exporter

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yifan-gu/blueNote/pkg/config"
	"github.com/yifan-gu/blueNote/pkg/model"
	"github.com/yifan-gu/blueNote/pkg/util"
)

var registeredExporters map[string]Exporter

type Exporter interface {
	Name() string
	LoadConfigs(cmd *cobra.Command)
	Export(cfg *config.ConvertConfig, book *model.Book) error
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
		util.Fatal(fmt.Errorf("unrecognized exporter type: %q", name))
	}
	return exporter
}

func ListExporters() []string {
	var names []string
	for _, exporter := range registeredExporters {
		names = append(names, exporter.Name())
	}
	return names
}

func LoadConfigs(cmd *cobra.Command) {
	for _, exporter := range registeredExporters {
		exporter.LoadConfigs(cmd)
	}
}
