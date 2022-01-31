/*
Copyright © 2022 Yifan Gu <guyifan1121@gmail.com>

*/

package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
	"github.com/yifan-gu/klipping2org/pkg/db"
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

func LoadConfig(cfgFile string, cfg *Config) error {
	viper.SetDefault("OutputDir", "./")
	viper.SetDefault("SplitBook", false)
	viper.SetDefault("RoamDir", DefaultRoamDir)
	viper.SetDefault("AuthorSubDir", true)

	viper.SetDefault("UpdateRoamDB", false)
	viper.SetDefault("RoamDBPath", DefaultRoamDBPath)
	viper.SetDefault("DBDriver", db.SqlDriverSqilite3)

	viper.SetDefault("Parser", DefaultParser)
	viper.SetDefault("InsertRoamLink", true)

	viper.SetDefault("PromptYesToAll", false)
	viper.SetDefault("PromptNoToAll", false)

	f, err := os.Open(cfgFile)
	if err != nil {
		return fmt.Errorf("failed to open config file %s: %v", cfgFile, err)
	}
	defer f.Close()

	if err := viper.ReadConfig(f); err != nil {
		return fmt.Errorf("failed to read config file %s: %v", cfgFile, err)
	}
	cfg.OutputDir = viper.GetString("OutputDir")
	cfg.SplitBook = viper.GetBool("SplitBook")
	cfg.RoamDir = viper.GetString("RoamDir")
	cfg.AuthorSubDir = viper.GetBool("AuthorSubDir")

	cfg.UpdateRoamDB = viper.GetBool("UpdateRoamDB")
	cfg.RoamDBPath = viper.GetString("RoamDBPath")
	cfg.DBDriver = viper.GetString("DBDriver")

	cfg.Parser = viper.GetString("Parser")
	cfg.InsertRoamLink = viper.GetBool("InsertRoamLink")

	cfg.PromptYesToAll = viper.GetBool("PromptYesToAll")
	cfg.PromptNoToAll = viper.GetBool("PromptNoToAll")

	return nil
}
