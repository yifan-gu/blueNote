/*
Copyright © 2022 Yifan Gu <guyifan1121@gmail.com>

*/
package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/yifan-gu/blueNote/pkg/util"
)

var rootCmd = &cobra.Command{
	Use:   "blueNote",
	Short: "note organizer",
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		util.Fatal(err)
	}
}
