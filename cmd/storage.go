/*
Copyright © 2022 Yifan Gu <guyifan1121@gmail.com>

*/

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yifan-gu/blueNote/pkg/config"
	"github.com/yifan-gu/blueNote/pkg/storage"
)

var storageConfig config.StorageConfig

var storageCmd = &cobra.Command{
	Use:   "storage",
	Short: "Operate on the storage directly",
	Run:   runStorage,
}

func runStorage(cmd *cobra.Command, args []string) {
	if len(args) > 1 {
		cmd.Help()
		os.Exit(1)
	}

	if storageConfig.ListStorages {
		printStoragesAndExit()
	}

	if len(args) < 1 {
		cmd.Help()
		fmt.Println(len(args))
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(storageCmd)

	storageCmd.PersistentFlags().BoolVar(&storageConfig.ListStorages, "list-storages", false, "list the supported storages")
	storageCmd.PersistentFlags().StringVar(&storageConfig.Storage, "storage", config.DefaultStorage, "the storage to use")
	storageCmd.PersistentFlags().StringVar(&storageConfig.Filter, "filter", "", "the filters for the storage CRUD operation, expecting a json format (e.g. \"{\"_id\":\"<id>\"}\")")

	registerStorages()

	storage.LoadConfigs(storageCmd)
	storage.LoadConfigs(serverCmd)
}
