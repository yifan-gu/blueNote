/*
Copyright © 2022 Yifan Gu <guyifan1121@gmail.com>

*/
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yifan-gu/blueNote/pkg/storage"
	"github.com/yifan-gu/blueNote/pkg/util"
)

var storageDelCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete marks from the storage",
	Run:   runStorageDelete,
}

func runStorageDelete(cmd *cobra.Command, args []string) {
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

	cnt, err := store.DeleteMarks(ctx, storageConfig.Filter)
	if err != nil {
		util.StackTraceErrorAndExit(err)
	}
	fmt.Println("Total deleted:", cnt)
}

func init() {
	storageCmd.AddCommand(storageDelCmd)
}
