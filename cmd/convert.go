/*
Copyright © 2022 Yifan Gu <guyifan1121@gmail.com>

*/

package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/yifan-gu/blueNote/pkg/config"
	"github.com/yifan-gu/blueNote/pkg/exporter"
	"github.com/yifan-gu/blueNote/pkg/parser"
	"github.com/yifan-gu/blueNote/pkg/util"
)

var convertConfig config.ConvertConfig

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert reading notes and clippings",
	Run:   runConvert,
}

func runConvert(cmd *cobra.Command, args []string) {
	if len(args) > 2 {
		cmd.Help()
		os.Exit(1)
	}

	if convertConfig.ListParsers {
		printParsersAndExit()
	}

	if convertConfig.ListExporters {
		printExportersAndExit()
	}

	if len(args) != 0 {
		convertConfig.InputPath = args[0]
	} else {
		if convertConfig.Parser != "json" {
			cmd.Help()
			os.Exit(1)
		}
	}

	convertConfig.OutputDir = "./"
	if len(args) == 2 {
		convertConfig.OutputDir = args[1]
	}

	books, err := parser.GetParser(convertConfig.Parser).Parse(convertConfig.InputPath)
	if err != nil {
		util.StackTraceErrorAndExit(err)
	}

	exp := exporter.GetExporter(convertConfig.Exporter)
	if err := exp.Export(&convertConfig, books); err != nil {
		util.StackTraceErrorAndExit(err)
	}
}

func init() {
	rootCmd.AddCommand(convertCmd)

	convertCmd.PersistentFlags().BoolVar(&convertConfig.ListParsers, "list-parsers", false, "list the supported parsers")
	convertCmd.PersistentFlags().BoolVar(&convertConfig.ListExporters, "list-exporters", false, "list the supported exporters")
	convertCmd.PersistentFlags().StringVarP(&convertConfig.Parser, "parser", "i", config.DefaultParser, "the parser to use")
	convertCmd.PersistentFlags().StringVarP(&convertConfig.Exporter, "exporter", "o", config.DefaultExporter, "the exporter to use")

	registerParsers()
	parser.LoadConfigs(convertCmd)

	registerExporters()
	exporter.LoadConfigs(convertCmd)
}
