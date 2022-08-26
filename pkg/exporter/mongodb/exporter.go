/*
Copyright © 2022 Yifan Gu <guyifan1121@gmail.com>

*/

package mongodb

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/yifan-gu/blueNote/pkg/config"
	"github.com/yifan-gu/blueNote/pkg/model"
	"github.com/yifan-gu/blueNote/pkg/storage/mongodb"
	"github.com/yifan-gu/blueNote/pkg/util"
)

type MongoDBExporter struct {
	mongodbConfig mongodb.Config
}

func (e *MongoDBExporter) Name() string {
	return "mongodb"
}

func (e *MongoDBExporter) LoadConfigs(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&e.mongodbConfig.Username, "mongodb.username", "", "username of the mongodb")
	cmd.PersistentFlags().StringVar(&e.mongodbConfig.Password, "mongodb.password", "", "password of the mongodb")
	cmd.PersistentFlags().StringVar(&e.mongodbConfig.Host, "mongodb.host", "localhost:27017", "host of the mongodb")
	cmd.PersistentFlags().StringVar(&e.mongodbConfig.ConnOpt, "mongodb.conn-opt", "", "connection option of the mongodb")
	cmd.PersistentFlags().StringVar(&e.mongodbConfig.DBName, "mongodb.database", "bluenote", "database to use in the mongodb")
	cmd.PersistentFlags().StringVar(&e.mongodbConfig.CollectionName, "mongodb.collection", "marks", "the collection to use in the mongodb")
}

func (e *MongoDBExporter) Export(cfg *config.ConvertConfig, books []*model.Book) error {
	ctx := context.Background()

	conn := mongodb.NewMongoDBStorage(ctx, &e.mongodbConfig)
	if err := conn.Connect(ctx); err != nil {
		return err
	}
	defer conn.Close(ctx)

	var totalInserted int
	for _, book := range books {
		for _, mark := range book.Marks {
			id, err := conn.CreateMark(ctx, mark)
			if err != nil {
				return err
			}
			util.Logf("Mark created with id: %s\n", id)
			totalInserted++
		}
	}
	util.Logf("Successfully loaded to mongodb, (database: %s, collection: %s)\n", e.mongodbConfig.DBName, e.mongodbConfig.CollectionName)
	util.Logf("Total inserted: %v\n", totalInserted)
	return nil
}
