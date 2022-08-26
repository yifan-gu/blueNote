/*
Copyright © 2022 Yifan Gu <guyifan1121@gmail.com>

*/

package storage

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yifan-gu/blueNote/pkg/model"
	"github.com/yifan-gu/blueNote/pkg/util"
)

var registeredStorages map[string]Storage

type Storage interface {
	Name() string
	LoadConfigs(cmd *cobra.Command)
	Connect(ctx context.Context) error
	CreateMark(ctx context.Context, mark *model.Mark) (id string, err error)
	GetMarks(ctx context.Context, filter interface{}, limit int) ([]*model.Mark, error)
	UpdateMarks(ctx context.Context, filter interface{}, update *model.Mark) (ids []string, err error)
	UpdateOneMark(ctx context.Context, id string, update *model.Mark) error
	DeleteMarks(ctx context.Context, filter interface{}) (int, error)
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
