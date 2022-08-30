/*
Copyright © 2022 Yifan Gu <guyifan1121@gmail.com>

*/

package mongodb

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

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
	Chapter  string `bson:"chapter,omitempty"`
	Page     *int   `bson:"page,omitempty"`
	Location *int   `bson:"location,omitempty"`
}

// PersistentMark defines the details of a mark object that will be stored in the databse.
type PersistentMark struct {
	ID             primitive.ObjectID `bson:"_id"`
	Type           string             `bson:"type"`
	Title          string             `bson:"title"`
	Author         string             `bson:"author"`
	Section        string             `bson:"section,omitempty"`
	Location       *Location          `bson:"location,omitempty"`
	Data           string             `bson:"data,omitempty"`
	UserNote       string             `bson:"note,omitempty"`
	Tags           []string           `bson:"tags,omitempty"`
	CreatedAt      *int64             `bson:"createdAt"`
	LastModifiedAt *int64             `bson:"lastModifiedAt"`
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

func (s *MongoDBStorage) CreateMark(ctx context.Context, mark *model.Mark) (string, error) {
	pm, err := MarkToPersistentMark(mark)
	if err != nil {
		return "", err
	}
	now := util.NowUnixMilli()
	pm.CreatedAt = &now
	pm.LastModifiedAt = &now
	result, err := s.coll.InsertOne(ctx, pm)
	if err != nil {
		return "", errors.Wrap(err, "")
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (s *MongoDBStorage) GetMarks(ctx context.Context, filter interface{}, limit int) ([]*model.Mark, error) {
	filterVal, err := parseFilter(filter)
	if err != nil {
		return nil, err
	}

	var result []*model.Mark
	if limit < 0 {
		limit = 0
	}
	cur, err := s.coll.Find(ctx, filterVal, options.Find().SetLimit(int64(limit)))
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

func (s *MongoDBStorage) UpdateMarks(ctx context.Context, filter interface{}, update *model.Mark) ([]string, error) {
	var ids []string
	marks, err := s.GetMarks(ctx, filter, 0)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	for _, mk := range marks {
		objectID, err := primitive.ObjectIDFromHex(mk.ID)
		if err != nil {
			return nil, errors.Wrap(err, "")
		}
		if _, err := s.coll.UpdateByID(ctx, objectID, constructUpdateFromMark(mk, update)); err != nil {
			return nil, errors.Wrap(err, "")
		}
		ids = append(ids, mk.ID)
	}
	return ids, nil
}

func (s *MongoDBStorage) UpdateOneMark(ctx context.Context, id string, update *model.Mark) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.Wrap(err, "")
	}
	marks, err := s.GetMarks(ctx, bson.M{"_id": objectID}, 0)
	if err != nil {
		return errors.Wrap(err, "")
	}
	if len(marks) != 0 {
		return errors.New(fmt.Sprintf("Expecting 1 mark for id %q, but saw %v", id, len(marks)))
	}
	if _, err := s.coll.UpdateByID(ctx, id, constructUpdateFromMark(marks[0], update)); err != nil {
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
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.Wrap(err, "")
	}
	result, err := s.coll.DeleteOne(ctx, bson.M{"_id": objectID})
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
		return parseFilterBSONM(val)
	default:
		return nil, errors.New(fmt.Sprintf("Invalid filter type %T, expecting \"string\" or \"bsons.M\"", val))
	}
}

func parseFilterBSONM(filter bson.M) (bson.M, error) {
	ret := filter
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
func MarkToPersistentMark(mark *model.Mark) (*PersistentMark, error) {
	ret := &PersistentMark{
		Type:    mark.Type,
		Title:   mark.Title,
		Author:  mark.Author,
		Section: mark.Section,
		Location: &Location{
			Chapter:  mark.Location.Chapter,
			Page:     mark.Location.Page,
			Location: mark.Location.Location,
		},
		Data:           mark.Data,
		UserNote:       mark.UserNote,
		Tags:           mark.Tags,
		CreatedAt:      mark.CreatedAt,
		LastModifiedAt: mark.LastModifiedAt,
	}
	if mark.ID != "" {
		id, err := primitive.ObjectIDFromHex(mark.ID)
		if err != nil {
			return nil, errors.Wrap(err, "")
		}
		ret.ID = id
	} else {
		ret.ID = primitive.NewObjectID()
	}
	return ret, nil
}

// PersistentMarkToMark converts a PersistentMark to a Mark.
func PersistentMarkToMark(pm *PersistentMark) *model.Mark {
	mark := &model.Mark{
		ID:             pm.ID.Hex(),
		Type:           pm.Type,
		Title:          pm.Title,
		Author:         pm.Author,
		Section:        pm.Section,
		Data:           pm.Data,
		UserNote:       pm.UserNote,
		Tags:           pm.Tags,
		CreatedAt:      pm.CreatedAt,
		LastModifiedAt: pm.LastModifiedAt,
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

func constructUpdateFromMark(original, update *model.Mark) bson.M {
	b := bson.M{}
	var modified bool

	if update.Type != "" && update.Type != original.Type {
		b["type"] = update.Type
		modified = true
	}
	if update.Title != "" && update.Title != original.Title {
		b["title"] = update.Title
		modified = true
	}
	if update.Author != "" && update.Author != original.Author {
		b["author"] = update.Author
		modified = true
	}
	if update.Section != "" && update.Section != original.Section {
		b["section"] = update.Section
		modified = true
	}
	if update.Location != nil {
		if update.Location.Chapter != "" && (original.Location == nil || update.Location.Chapter != original.Location.Chapter) {
			b["location.chapter"] = update.Location.Chapter
			modified = true
		}
		if update.Location.Page != nil && (original.Location == nil || original.Location.Page == nil || *update.Location.Page != *original.Location.Page) {
			b["location.page"] = update.Location.Page
			modified = true
		}
		if update.Location.Location != nil && (original.Location == nil || original.Location.Location == nil || *update.Location.Location != *original.Location.Location) {
			b["location.location"] = update.Location.Location
			modified = true
		}
	}
	if update.Data != "" && update.Data != original.Data {
		b["data"] = update.Data
		modified = true
	}
	if update.UserNote != "" && update.UserNote != original.UserNote {
		b["note"] = update.UserNote
		modified = true
	}
	if update.Tags != nil {
		sort.StringSlice(update.Tags).Sort()
		sort.StringSlice(original.Tags).Sort()
		if !util.StringSlicesEqual(update.Tags, original.Tags) {
			b["tags"] = update.Tags
			modified = true
		}
	}
	if modified {
		b["lastModifiedAt"] = util.NowUnixMilli()
	}
	return bson.M{"$set": b}
}
