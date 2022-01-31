/*
Copyright © 2022 Yifan Gu <guyifan1121@gmail.com>

*/

package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	InputPath    string
	OutputDir    string
	SplitBook    bool
	RoamDir      string
	AuthorSubDir bool

	UpdateRoamDB bool
	RoamDBPath   string
	DBDriver     string

	Parser         string
	InsertRoamLink bool

	PromptYesToAll bool
	PromptNoToAll  bool
}

const (
	DefaultRoamDir    = "~/org/roam"
	DefaultRoamDBPath = "~/.emacs.d/.local/etc/org-roam.db"
	DefaultParser     = "htmlclipping"
)

func LoadConfig(cfgFile string, cfg *Config, cmd *cobra.Command) error {
	viper.BindPFlag("OUTPUT_DIR", cmd.PersistentFlags().Lookup("output-dir"))
	viper.BindPFlag("SPLIT_BOOK", cmd.PersistentFlags().Lookup("split"))
	viper.BindPFlag("ROAM_DIR", cmd.PersistentFlags().Lookup("roam-dir"))
	viper.BindPFlag("AUTHOR_SUBDIR", cmd.PersistentFlags().Lookup("author-sub-dir"))

	viper.BindPFlag("UPDATE_ROAM_DB", cmd.PersistentFlags().Lookup("update-roam-db"))
	viper.BindPFlag("ROAM_DB_PATH", cmd.PersistentFlags().Lookup("roam-db-path"))
	viper.BindPFlag("DB_DRIVER", cmd.PersistentFlags().Lookup("db-driver"))

	viper.BindPFlag("PARSER", cmd.PersistentFlags().Lookup("parser"))
	viper.BindPFlag("INSERT_ROAM_LINK", cmd.PersistentFlags().Lookup("insert-roam-link"))

	viper.BindPFlag("PROMPT_YES_TO_ALL", cmd.PersistentFlags().Lookup("yes-to-all"))
	viper.BindPFlag("PROMPT_NO_TO_ALL", cmd.PersistentFlags().Lookup("no-to-all"))

	f, err := os.Open(cfgFile)
	if err != nil {
		return fmt.Errorf("failed to open config file %s: %v", cfgFile, err)
	}
	defer f.Close()

	viper.SetConfigType(filepath.Ext(cfgFile)[1:])
	if err := viper.ReadConfig(f); err != nil {
		return fmt.Errorf("failed to read config file %s: %v", cfgFile, err)
	}

	cfg.OutputDir = viper.GetString("OUTPUT_DIR")
	cfg.SplitBook = viper.GetBool("SPLIT_BOOK")
	cfg.RoamDir = viper.GetString("ROAM_DIR")
	cfg.AuthorSubDir = viper.GetBool("AUTHOR_SUBDIR")

	cfg.UpdateRoamDB = viper.GetBool("UPDATE_ROAM_DB")
	cfg.RoamDBPath = viper.GetString("ROAM_DB_PATH")
	cfg.DBDriver = viper.GetString("DB_DRIVER")

	cfg.Parser = viper.GetString("PARSER")
	cfg.InsertRoamLink = viper.GetBool("INSERT_ROAM_LINK")

	cfg.PromptYesToAll = viper.GetBool("PROMPT_YES_TO_ALL")
	cfg.PromptNoToAll = viper.GetBool("PROMPT_NO_TO_ALL")

	return nil
}
