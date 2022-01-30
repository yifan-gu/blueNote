/*
Copyright © 2022 Yifan Gu <guyifan1121@gmail.com>

*/

package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	InputPath      string
	OutputDir      string
	SplitBook      bool
	RoamDir        string
	RoamDBPath     string
	InsertRoamLink bool
	AuthorSubDir   bool
}

const (
	DefaultRoamDir    = "~/org/roam"
	DefaultRoamDBPath = "~/.emacs.d/.local/etc/org-roam.db"
)

func LoadConfig(cfgFile string, cfg *Config) error {
	viper.SetDefault("OutputDir", "./")
	viper.SetDefault("SplitBook", false)
	viper.SetDefault("RoamDir", DefaultRoamDir)
	viper.SetDefault("RoamDBPath", DefaultRoamDBPath)
	viper.SetDefault("InsertRoamLink", true)
	viper.SetDefault("AuthorSubDir", false)

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
	cfg.RoamDBPath = viper.GetString("RoamDBPath")
	cfg.InsertRoamLink = viper.GetBool("InsertRoamLink")
	cfg.AuthorSubDir = viper.GetBool("AuthorSubDir")

	return nil
}
