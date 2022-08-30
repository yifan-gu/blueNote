/*
Copyright © 2022 Yifan Gu <guyifan1121@gmail.com>

*/

package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yifan-gu/blueNote/pkg/storage"
	"github.com/yifan-gu/blueNote/pkg/util"
)

var storageGetCmd = &cobra.Command{
	Use:   "get",
	Short: "get marks from the storage",
	Run:   runStorageGet,
}

var storageGetLimit int

func runStorageGet(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	if len(args) != 0 {
		cmd.Help()
		os.Exit(1)
	}

	store := storage.GetStorages(storageConfig.Storage)
	if err := store.Connect(ctx); err != nil {
		util.StackTraceErrorAndExit(err)
	}
	defer store.Close(ctx)

	if storageConfig.Filter == "" {
		util.Fatal("Missing parameters for --filter")
	}

	marks, err := store.GetMarks(ctx, storageConfig.Filter, storageGetLimit)
	if err != nil {
		util.StackTraceErrorAndExit(err)
	}
	b, err := json.MarshalIndent(marks, "", "  ")
	if err != nil {
		util.Fatal(err)
	}
	fmt.Println(string(b))
}

func init() {
	storageCmd.AddCommand(storageGetCmd)
	storageGetCmd.PersistentFlags().IntVar(&storageGetLimit, "limit", 0, "set the maximum number of marks to return, set 0 or negative to return all")
}
