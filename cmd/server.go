/*
Copyright © 2022 Yifan Gu <guyifan1121@gmail.com>

*/

package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	"github.com/yifan-gu/blueNote/pkg/config"
	"github.com/yifan-gu/blueNote/pkg/server"
	"github.com/yifan-gu/blueNote/pkg/storage"
	"github.com/yifan-gu/blueNote/pkg/util"
)

var serverConfig config.ServerConfig

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the backend server to serve marks",
	Run:   runServer,
}

func runServer(cmd *cobra.Command, args []string) {
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

	server.NewServer(&serverConfig, store).Run()
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.PersistentFlags().StringVar(&serverConfig.ListenAddr, "server.addr", "localhost:11212", "The port to listen for the server.")
}
