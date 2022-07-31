/*
Copyright © 2022 Yifan Gu <guyifan1121@gmail.com>

*/
package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/yifan-gu/blueNote/pkg/config"
	"github.com/yifan-gu/blueNote/pkg/exporter"
	"github.com/yifan-gu/blueNote/pkg/model"
	"github.com/yifan-gu/blueNote/pkg/parser"
)

var cfg config.GlobalConfig

var rootCmd = &cobra.Command{
	Use:   "blueNote",
	Short: "Convert reading notes and clippings",
	Run:   run,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func run(cmd *cobra.Command, args []string) {
	if len(args) < 1 || len(args) > 2 {
		cmd.Help()
		os.Exit(1)
	}

	cfg.InputPath = args[0]
	cfg.OutputDir = "./"
	if len(args) == 2 {
		cfg.OutputDir = args[1]
	}

	book, err := parser.GetParser(cfg.Parser).Parse(cfg.InputPath)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.Author != "" {
		book.Author = cfg.Author
	}
	if cfg.Title != "" {
		book.Title = cfg.Title
	}

	books := []*model.Book{book}
	if cfg.SplitBook {
		books = book.Split()
	}
	cfg.TotalBookCnt = len(books)

	exp := exporter.GetExporter(cfg.Exporter)
	for i, bk := range books {
		cfg.CurrentBookIndex = i
		if err := exp.Export(&cfg, bk); err != nil {
			log.Fatal(err)
		}
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&cfg.SplitBook, "split", "s", false, "split sub-sections into separate books")
	rootCmd.PersistentFlags().BoolVarP(&cfg.AuthorSubDir, "author-sub-dir", "a", true, "create sub-directory with the name of the author")

	rootCmd.PersistentFlags().StringVar(&cfg.Author, "author", "", "override the book author name")
	rootCmd.PersistentFlags().StringVar(&cfg.Title, "title", "", "override the book title name")

	rootCmd.PersistentFlags().StringVarP(&cfg.Parser, "parser", "i", config.DefaultParser, "the parser to use")
	rootCmd.PersistentFlags().StringVarP(&cfg.Exporter, "exporter", "o", config.DefaultExporter, "the exporter to use")

	rootCmd.PersistentFlags().BoolVarP(&cfg.PromptYesToAll, "yes-to-all", "y", false, "set yes to all prompt confirmation")
	rootCmd.PersistentFlags().BoolVarP(&cfg.PromptNoToAll, "no-to-all", "n", false, "set no to all prompt confirmation")

	registerParsers()
	parser.LoadConfigs(rootCmd)

	registerExporters()
	exporter.LoadConfigs(rootCmd)
}
