/*
Copyright © 2022 Yifan Gu <guyifan1121@gmail.com>

*/

package mongodb

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/yifan-gu/blueNote/pkg/config"
	"github.com/yifan-gu/blueNote/pkg/model"
	"github.com/yifan-gu/blueNote/pkg/util"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBExporter struct {
	username string
	password string
	host     string
	connOpt  string

	dbName         string
	collectionName string
}

func (e *MongoDBExporter) Name() string {
	return "mongodb"
}

func (e *MongoDBExporter) LoadConfigs(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&e.username, "mongodb.username", "", "username of the mongodb")
	cmd.PersistentFlags().StringVar(&e.password, "mongodb.password", "", "password of the mongodb")
	cmd.PersistentFlags().StringVar(&e.host, "mongodb.host", "localhost:27017", "host of the mongodb")
	cmd.PersistentFlags().StringVar(&e.connOpt, "mongodb.conn-opt", "", "connection option of the mongodb")
	cmd.PersistentFlags().StringVar(&e.dbName, "mongodb.database", "bluenote", "database to use in the mongodb")
	cmd.PersistentFlags().StringVar(&e.collectionName, "mongodb.collection", "marks", "the collection to use in the mongodb")
}

func (e *MongoDBExporter) constructConnectionURI() string {
	uri := "mongodb://"
	if e.username != "" && e.password != "" {
		uri += fmt.Sprintf("%s:%s@", e.username, e.password)
	}
	uri += fmt.Sprintf("%s/?%s", e.host, e.connOpt)
	return uri
}

func (e *MongoDBExporter) Export(cfg *config.ConvertConfig, books []*model.Book) error {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(e.constructConnectionURI()))
	if err != nil {
		return errors.Wrap(err, "")
	}
	defer client.Disconnect(ctx)

	coll := client.Database(e.dbName).Collection(e.collectionName)

	var totalInserted, alreadyExisted int
	for _, book := range books {
		for _, mark := range book.Marks {
			mk := mark.ToPersistenMark()
			err := coll.FindOne(ctx, bson.D{{"digest", mk.Digest}}).Err()
			if err == nil {
				alreadyExisted++
				continue
			}
			if err != mongo.ErrNoDocuments {
				return errors.Wrap(err, "")
			}
			_, err = coll.InsertOne(ctx, mk)
			if err != nil {
				return errors.Wrap(err, "")
			}
			totalInserted++
		}
	}
	util.Logf("Successfully loaded to mongodb database: %q\n", e.dbName)
	util.Logf("Total inserted: %v, already existed: %v\n", totalInserted, alreadyExisted)
	return nil
}
