/*
Copyright © 2022 Yifan Gu <guyifan1121@gmail.com>

*/

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yifan-gu/blueNote/pkg/config"
	"github.com/yifan-gu/blueNote/pkg/util"
)

var rootCmd = &cobra.Command{
	Use:   "blueNote",
	Short: "note organizer",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		util.Fatal(err)
	}
}

func init() {
	convertCmd.PersistentFlags().BoolVarP(&config.GlobalCfg.PromptYesToAll, "yes-to-all", "y", false, "set yes to all prompt confirmation")
	convertCmd.PersistentFlags().BoolVarP(&config.GlobalCfg.PromptNoToAll, "no-to-all", "n", false, "set no to all prompt confirmation")
}
