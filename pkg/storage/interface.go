/*
Copyright © 2022 Yifan Gu <guyifan1121@gmail.com>

*/

package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/yifan-gu/blueNote/pkg/model"
	"github.com/yifan-gu/blueNote/pkg/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var registeredStorages map[string]Storage

type Filter struct {
	model.PersistentMark
}

type Storage interface {
	Name() string
	LoadConfigs(cmd *cobra.Command)
	Connect(ctx context.Context) error
	CreateMark(ctx context.Context, mark *model.PersistentMark) error
	GetMarks(ctx context.Context, filter bson.M) ([]*model.PersistentMark, error)
	UpdateMarks(ctx context.Context, filter, update bson.M) (int, error)
	UpdateOneMark(ctx context.Context, id string, update bson.M) error
	DeleteMarks(ctx context.Context, filter bson.M) (int, error)
	DeleteOneMark(ctx context.Context, id string) error
	Close(ctx context.Context) error
}

func RegisterStorage(storage Storage) {
	name := strings.ToLower(storage.Name())
	if registeredStorages == nil {
		registeredStorages = make(map[string]Storage)
	}
	if _, ok := registeredStorages[name]; ok {
		util.Fatal(fmt.Errorf("Name conflict for storage %q", name))
	}
	registeredStorages[name] = storage
}

func GetStorages(name string) Storage {
	name = strings.ToLower(name)
	storage, ok := registeredStorages[name]
	if !ok {
		util.Fatal(fmt.Errorf("unrecognized storage type: %q", name))
	}
	return storage
}

func ListStorages() []string {
	var names []string
	for _, storage := range registeredStorages {
		names = append(names, storage.Name())
	}
	return names
}

func LoadConfigs(cmd *cobra.Command) {
	for _, storage := range registeredStorages {
		storage.LoadConfigs(cmd)
	}
}

func ParseFilterString(filter string) (bson.M, error) {
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
