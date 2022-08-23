/*
Copyright © 2022 Yifan Gu <guyifan1121@gmail.com>

*/

package mongodb

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/yifan-gu/blueNote/pkg/model"
	"github.com/yifan-gu/blueNote/pkg/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBStorage struct {
	cfg    *Config
	client *mongo.Client
	coll   *mongo.Collection
}

type Config struct {
	Username       string
	Password       string
	Host           string
	ConnOpt        string
	DBName         string
	CollectionName string
}

func (c *Config) constructConnectionURI() string {
	uri := "mongodb://"
	if c.Username != "" && c.Password != "" {
		uri += fmt.Sprintf("%s:%s@", c.Username, c.Password)
	}
	uri += fmt.Sprintf("%s/?%s", c.Host, c.ConnOpt)
	return uri
}

func NewMongoDBStorage(ctx context.Context, cfg *Config) storage.Storage {
	return &MongoDBStorage{cfg: cfg}
}

func (s *MongoDBStorage) Name() string {
	return "mongodb"
}

func (s *MongoDBStorage) LoadConfigs(cmd *cobra.Command) {
	if s.cfg == nil {
		s.cfg = &Config{}
	}
	cmd.PersistentFlags().StringVar(&s.cfg.Username, "mongodb.username", "", "username of the mongodb")
	cmd.PersistentFlags().StringVar(&s.cfg.Password, "mongodb.password", "", "password of the mongodb")
	cmd.PersistentFlags().StringVar(&s.cfg.Host, "mongodb.host", "localhost:27017", "host of the mongodb")
	cmd.PersistentFlags().StringVar(&s.cfg.ConnOpt, "mongodb.conn-opt", "", "connection option of the mongodb")
	cmd.PersistentFlags().StringVar(&s.cfg.DBName, "mongodb.database", "bluenote", "database to use in the mongodb")
	cmd.PersistentFlags().StringVar(&s.cfg.CollectionName, "mongodb.collection", "marks", "the collection to use in the mongodb")
}

func (s *MongoDBStorage) Connect(ctx context.Context) error {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(s.cfg.constructConnectionURI()))
	if err != nil {
		return errors.Wrap(err, "")
	}
	if err := client.Ping(ctx, nil); err != nil {
		return errors.Wrap(err, "")
	}
	s.client = client
	s.coll = client.Database(s.cfg.DBName).Collection(s.cfg.CollectionName)
	return nil
}

func (s *MongoDBStorage) CreateMark(ctx context.Context, mark *model.PersistentMark) error {
	if _, err := s.coll.InsertOne(ctx, mark); err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

func (s *MongoDBStorage) GetMarks(ctx context.Context, filter bson.M) ([]*model.PersistentMark, error) {
	var result []*model.PersistentMark

	cur, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	for cur.Next(ctx) {
		var mark model.PersistentMark
		if err := cur.Decode(&mark); err != nil {
			return nil, errors.Wrap(err, "")
		}
		result = append(result, &mark)
	}
	if err := cur.Err(); err != nil {
		return nil, errors.Wrap(err, "")
	}
	return result, nil
}

func (s *MongoDBStorage) UpdateMarks(ctx context.Context, filter, update bson.M) (int, error) {
	marks, err := s.GetMarks(ctx, filter)
	if err != nil {
		return 0, errors.Wrap(err, "")
	}
	for _, mk := range marks {
		if _, err := s.coll.UpdateByID(ctx, mk.ID, update); err != nil {
			return 0, errors.Wrap(err, "")
		}
	}
	return len(marks), nil
}

func (s *MongoDBStorage) UpdateOneMark(ctx context.Context, id string, update bson.M) error {
	if _, err := s.coll.UpdateByID(ctx, id, update); err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

func (s *MongoDBStorage) DeleteMarks(ctx context.Context, filter bson.M) (int, error) {
	result, err := s.coll.DeleteMany(ctx, filter)
	if err != nil {
		return 0, errors.Wrap(err, "")
	}
	return int(result.DeletedCount), nil
}

func (s *MongoDBStorage) DeleteOneMark(ctx context.Context, id string) error {
	result, err := s.coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return errors.Wrap(err, "")
	}
	if result.DeletedCount == 0 {
		return errors.New(fmt.Sprintf("no such mark found for id: %q", id))
	}
	return nil
}

func (s *MongoDBStorage) Close(ctx context.Context) error {
	if err := s.client.Disconnect(ctx); err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}
