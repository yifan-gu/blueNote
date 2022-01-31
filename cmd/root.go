/*
Copyright © 2022 Yifan Gu <guyifan1121@gmail.com>

*/
package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/yifan-gu/klipping2org/pkg/config"
	"github.com/yifan-gu/klipping2org/pkg/db"
	"github.com/yifan-gu/klipping2org/pkg/parser"
)

var (
	cfg     config.Config
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:   "klipping2org",
	Short: "Convert Kindle exported clipping html(s) to org file(s)",
	Run:   run,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func run(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		cmd.Help()
		os.Exit(1)
	}

	if cfgFile != "" {
		if err := config.LoadConfig(cfgFile, &cfg); err != nil {
			log.Fatal(err)
		}
	}

	cfg.InputPath = args[0]
	if err := parser.ParseAndWrite(&cfg); err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "path to the config file, the command line parameters will override the config files if both are provided")
	rootCmd.PersistentFlags().StringVarP(&cfg.OutputDir, "output", "o", "./", "output directory")
	rootCmd.PersistentFlags().BoolVarP(&cfg.SplitBook, "split", "s", false, "split sub-sections into separate books")
	rootCmd.PersistentFlags().StringVarP(&cfg.RoamDir, "roam-dir", "r", config.DefaultRoamDir, "path to the org-roam directory")
	rootCmd.PersistentFlags().BoolVarP(&cfg.AuthorSubDir, "author-sub-dir", "a", true, "create sub-directory with the name of the author")

	rootCmd.PersistentFlags().BoolVar(&cfg.UpdateRoamDB, "update-roam-db", false, "automatically update the roam sqlite db for links")
	rootCmd.PersistentFlags().StringVarP(&cfg.RoamDBPath, "roam-db-path", "d", config.DefaultRoamDBPath, "path to the org-roam sqlite3 database")
	rootCmd.PersistentFlags().StringVar(&cfg.DBDriver, "db-driver", db.SqlDriverSqilite3, "the database driver to use")

	rootCmd.PersistentFlags().StringVar(&cfg.Parser, "parser", config.DefaultParser, "the parser to use")
	rootCmd.PersistentFlags().BoolVarP(&cfg.InsertRoamLink, "insert-roam-link", "l", true, "insert the roam links")

	rootCmd.PersistentFlags().BoolVarP(&cfg.PromptYesToAll, "yes-to-all", "y", false, "set yes to all prompt confirmation")
	rootCmd.PersistentFlags().BoolVarP(&cfg.PromptNoToAll, "no-to-all", "n", false, "set no to all prompt confirmation")
}
