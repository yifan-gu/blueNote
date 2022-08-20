/*
Copyright © 2022 Yifan Gu <guyifan1121@gmail.com>

*/
package cmd

import (
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/yifan-gu/blueNote/pkg/config"
	"github.com/yifan-gu/blueNote/pkg/exporter"
	"github.com/yifan-gu/blueNote/pkg/parser"
	"github.com/yifan-gu/blueNote/pkg/util"
)

var cfg config.ConvertConfig

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert reading notes and clippings",
	Run:   runConvert,
}

func runConvert(cmd *cobra.Command, args []string) {
	if len(args) == 0 || len(args) > 2 {
		cmd.Help()
		os.Exit(1)
	}

	if cfg.ListParsers {
		printParsersAndExit()
	}

	if cfg.ListExporters {
		printExportersAndExit()
	}

	cfg.InputPath = args[0]
	cfg.OutputDir = "./"
	if len(args) == 2 {
		cfg.OutputDir = args[1]
	}

	books, err := parser.GetParser(cfg.Parser).Parse(cfg.InputPath)
	if err != nil {
		stackTraceableErr, ok := err.(stackTracer)
		util.Log(errors.Cause(err))
		if ok {
			util.Logf("%+v\n", stackTraceableErr.StackTrace())
		}
		os.Exit(1)
	}

	exp := exporter.GetExporter(cfg.Exporter)
	if err := exp.Export(&cfg, books); err != nil {
		util.Log(errors.Cause(err))
		stackTraceableErr, ok := err.(stackTracer)
		if ok {
			util.Logf("%+v\n", stackTraceableErr.StackTrace())
		}
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(convertCmd)

	convertCmd.PersistentFlags().BoolVar(&cfg.ListParsers, "list-parsers", false, "list the supported parsers")
	convertCmd.PersistentFlags().BoolVar(&cfg.ListExporters, "list-exporters", false, "list the supported exporters")

	convertCmd.PersistentFlags().StringVarP(&cfg.Parser, "parser", "i", config.DefaultParser, "the parser to use")
	convertCmd.PersistentFlags().StringVarP(&cfg.Exporter, "exporter", "o", config.DefaultExporter, "the exporter to use")

	convertCmd.PersistentFlags().BoolVarP(&cfg.PromptYesToAll, "yes-to-all", "y", false, "set yes to all prompt confirmation")
	convertCmd.PersistentFlags().BoolVarP(&cfg.PromptNoToAll, "no-to-all", "n", false, "set no to all prompt confirmation")

	registerParsers()
	parser.LoadConfigs(convertCmd)

	registerExporters()
	exporter.LoadConfigs(convertCmd)
}
