/*
Copyright © 2022 Yifan Gu <guyifan1121@gmail.com>

*/

package mongodb

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/yifan-gu/blueNote/pkg/model"
	"github.com/yifan-gu/blueNote/pkg/storage"
	"github.com/yifan-gu/blueNote/pkg/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Location defines the location of a mark in the book.
type Location struct {
	Chapter  string `json:"chapter,omitempty" bson:"chapter,omitempty"`
	Page     *int   `json:"page,omitempty" bson:"page,omitempty"`
	Location *int   `json:"location,omitempty" bson:"location,omitempty"`
}

// PersistentMark defines the details of a mark object that will be stored in the databse.
type PersistentMark struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Digest    string             `json:"digest" bson:"digest"`
	Type      string             `json:"type" bson:"type"`
	Title     string             `json:"title" bson:"title"`
	Author    string             `json:"author" bson:"author"`
	Section   string             `json:"section,omitempty" bson:"section,omitempty"`
	Location  *Location          `json:"location,omitempty" bson:"location,omitempty"`
	Data      string             `json:"data,omitempty" bson:"data,omitempty"`
	UserNotes string             `json:"notes,omitempty" bson:"notes,omitempty"`
	Tags      []string           `json:"tags,omitempty" bson:"tags,omitempty"`
}

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

func (s *MongoDBStorage) CreateMark(ctx context.Context, mark *model.Mark) error {
	if _, err := s.coll.InsertOne(ctx, MarkToPersistentMark(mark)); err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

func (s *MongoDBStorage) GetMarks(ctx context.Context, filter interface{}) ([]*model.Mark, error) {
	filterVal, err := parseFilter(filter)
	if err != nil {
		return nil, err
	}

	var result []*model.Mark
	cur, err := s.coll.Find(ctx, filterVal)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	for cur.Next(ctx) {
		var mark PersistentMark
		if err := cur.Decode(&mark); err != nil {
			return nil, errors.Wrap(err, "")
		}
		result = append(result, PersistentMarkToMark(&mark))
	}
	if err := cur.Err(); err != nil {
		return nil, errors.Wrap(err, "")
	}
	return result, nil
}

func (s *MongoDBStorage) UpdateMarks(ctx context.Context, filter interface{}, update *model.Mark) (int, error) {
	filterVal, err := parseFilter(filter)
	if err != nil {
		return 0, err
	}
	marks, err := s.GetMarks(ctx, filterVal)
	if err != nil {
		return 0, errors.Wrap(err, "")
	}
	for _, mk := range marks {
		if _, err := s.coll.UpdateByID(ctx, mk.ID, constructUpdateFromMark(update)); err != nil {
			return 0, errors.Wrap(err, "")
		}
	}
	return len(marks), nil
}

func (s *MongoDBStorage) UpdateOneMark(ctx context.Context, id string, update *model.Mark) error {
	if _, err := s.coll.UpdateByID(ctx, id, constructUpdateFromMark(update)); err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

func (s *MongoDBStorage) DeleteMarks(ctx context.Context, filter interface{}) (int, error) {
	filterVal, err := parseFilter(filter)
	if err != nil {
		return 0, err
	}
	result, err := s.coll.DeleteMany(ctx, filterVal)
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

func parseFilter(filter interface{}) (bson.M, error) {
	switch val := filter.(type) {
	case string:
		return parseFilterString(val)
	case bson.M:
		return val, nil
	default:
		return nil, errors.New(fmt.Sprintf("Invalid filter type %T, expecting \"string\" or \"bsons.M\"", val))
	}
}

func parseFilterString(filter string) (bson.M, error) {
	ret := bson.M{}
	if err := json.Unmarshal([]byte(filter), &ret); err != nil {
		return nil, errors.Wrap(err, "")
	}

	// Convert "_id" from string to primitive.ObjectID
	id := ret["_id"]
	if id != nil {
		hexID, ok := id.(string)
		if ok {
			objID, err := primitive.ObjectIDFromHex(hexID)
			if err != nil {
				return nil, errors.Wrap(err, "")
			}
			ret["_id"] = objID
		}
	}
	return ret, nil
}

// MarkToPersistentMark converts a Mark to a PersistentMark
func MarkToPersistentMark(mark *model.Mark) *PersistentMark {
	b, err := json.Marshal(mark)
	if err != nil {
		util.Fatal("cannot marshal:", err)
	}
	return &PersistentMark{
		ID:      primitive.NewObjectID(),
		Digest:  fmt.Sprintf("%x", sha256.Sum256(b)),
		Type:    mark.Type,
		Title:   mark.Title,
		Author:  mark.Author,
		Section: mark.Section,
		Location: &Location{
			Chapter:  mark.Location.Chapter,
			Page:     mark.Location.Page,
			Location: mark.Location.Location,
		},
		Data:      mark.Data,
		UserNotes: mark.UserNotes,
		Tags:      mark.Tags,
	}
}

// PersistentMarkToMark converts a PersistentMark to a Mark.
func PersistentMarkToMark(pm *PersistentMark) *model.Mark {
	mark := &model.Mark{
		ID:        pm.ID.Hex(),
		Type:      pm.Type,
		Title:     pm.Title,
		Author:    pm.Author,
		Section:   pm.Section,
		Data:      pm.Data,
		UserNotes: pm.UserNotes,
		Tags:      pm.Tags,
	}
	if pm.Location != nil {
		mark.Location = &model.Location{
			Chapter:  pm.Location.Chapter,
			Page:     pm.Location.Page,
			Location: pm.Location.Location,
		}
	}
	return mark
}

func constructUpdateFromMark(mark *model.Mark) bson.M {
	var update bson.M
	if mark.Type != "" {
		update["type"] = mark.Type
	}
	if mark.Title != "" {
		update["title"] = mark.Title
	}
	if mark.Author != "" {
		update["author"] = mark.Author
	}
	if mark.Section != "" {
		update["section"] = mark.Section
	}
	if mark.Location != nil {
		if mark.Location.Chapter != "" {
			update["location.chapter"] = mark.Location.Chapter
		}
		if mark.Location.Page != nil {
			update["location.page"] = mark.Location.Page
		}
		if mark.Location.Location != nil {
			update["location.location"] = mark.Location.Location
		}
	}
	if mark.Data != "" {
		update["data"] = mark.Data
	}
	if mark.UserNotes != "" {
		update["notes"] = mark.UserNotes
	}
	if mark.Tags != nil {
		update["tags"] = mark.Tags
	}
	return bson.M{"$set": update}
}
