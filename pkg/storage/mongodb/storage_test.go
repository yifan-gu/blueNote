/*
Copyright © 2022 Yifan Gu <guyifan1121@gmail.com>

*/

package mongodb

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yifan-gu/blueNote/pkg/model"
	"github.com/yifan-gu/blueNote/pkg/util"
	"go.mongodb.org/mongo-driver/bson"
)

func TestConstructUpdateFromMark(t *testing.T) {
	page10, loc100 := 10, 100
	page42, loc420 := 42, 420

	originalMark1 := model.Mark{
		Type:    "HIGHLIGHT",
		Title:   "Title T",
		Author:  "Author A",
		Section: "Section S",
		Location: &model.Location{
			Chapter:  "Chapter C",
			Page:     &page10,
			Location: &loc100,
		},
		Data:     "Data D",
		UserNote: "Note N",
		Tags:     []string{"Tag D", "Tag C", "Tag B", "Tag A"},
	}
	tests := []struct {
		original *model.Mark
		update   *model.Mark
		result   bson.M
	}{
		{
			original: &originalMark1,
			update:   &model.Mark{},
			result:   bson.M{"$set": bson.M{}},
		},
		{
			original: &originalMark1,
			update:   &originalMark1,
			result:   bson.M{"$set": bson.M{}},
		},
		{
			original: &originalMark1,
			update:   &model.Mark{Type: "NOTE"},
			result:   bson.M{"$set": bson.M{"type": "NOTE", "lastModifiedAt": int64(1)}},
		},
		{
			original: &originalMark1,
			update:   &model.Mark{Title: "Title I"},
			result:   bson.M{"$set": bson.M{"title": "Title I", "lastModifiedAt": int64(2)}},
		},
		{
			original: &originalMark1,
			update:   &model.Mark{Author: "Author U"},
			result:   bson.M{"$set": bson.M{"author": "Author U", "lastModifiedAt": int64(3)}},
		},
		{
			original: &originalMark1,
			update:   &model.Mark{Section: "Section E"},
			result:   bson.M{"$set": bson.M{"section": "Section E", "lastModifiedAt": int64(4)}},
		},
		{
			original: &originalMark1,
			update:   &model.Mark{Location: &model.Location{Chapter: "Chapter H"}},
			result:   bson.M{"$set": bson.M{"location.chapter": "Chapter H", "lastModifiedAt": int64(5)}},
		},
		{
			original: &originalMark1,
			update:   &model.Mark{Location: &model.Location{Page: &page42}},
			result:   bson.M{"$set": bson.M{"location.page": &page42, "lastModifiedAt": int64(6)}},
		},
		{
			original: &originalMark1,
			update:   &model.Mark{Location: &model.Location{Location: &loc420}},
			result:   bson.M{"$set": bson.M{"location.location": &loc420, "lastModifiedAt": int64(7)}},
		},
		{
			original: &originalMark1,
			update:   &model.Mark{Data: "Data A"},
			result:   bson.M{"$set": bson.M{"data": "Data A", "lastModifiedAt": int64(8)}},
		},
		{
			original: &originalMark1,
			update:   &model.Mark{UserNote: "Note O"},
			result:   bson.M{"$set": bson.M{"note": "Note O", "lastModifiedAt": int64(9)}},
		},
		{
			original: &originalMark1,
			update:   &model.Mark{Tags: []string{"tag d", "tag c", "tag b", "tag a"}},
			result:   bson.M{"$set": bson.M{"tags": []string{"tag a", "tag b", "tag c", "tag d"}, "lastModifiedAt": int64(10)}},
		},
		{
			original: &originalMark1,
			update: &model.Mark{
				Type:    "NOTE",
				Title:   "Title I",
				Author:  "Author U",
				Section: "Section E",
				Location: &model.Location{
					Chapter:  "Chapter H",
					Page:     &page42,
					Location: &loc420,
				},
				Data:     "Data A",
				UserNote: "Note O",
				Tags:     []string{"tag d", "tag c", "tag b", "tag a"}},
			result: bson.M{"$set": bson.M{
				"type":              "NOTE",
				"title":             "Title I",
				"author":            "Author U",
				"section":           "Section E",
				"location.chapter":  "Chapter H",
				"location.page":     &page42,
				"location.location": &loc420,
				"data":              "Data A",
				"note":              "Note O",
				"tags":              []string{"tag a", "tag b", "tag c", "tag d"},
				"lastModifiedAt":    int64(11),
			}},
		},
		{
			original: &model.Mark{},
			update: &model.Mark{
				Type:    "NOTE",
				Title:   "Title I",
				Author:  "Author U",
				Section: "Section E",
				Location: &model.Location{
					Chapter:  "Chapter H",
					Page:     &page42,
					Location: &loc420,
				},
				Data:     "Data A",
				UserNote: "Note O",
				Tags:     []string{"tag d", "tag c", "tag b", "tag a"}},
			result: bson.M{"$set": bson.M{
				"type":              "NOTE",
				"title":             "Title I",
				"author":            "Author U",
				"section":           "Section E",
				"location.chapter":  "Chapter H",
				"location.page":     &page42,
				"location.location": &loc420,
				"data":              "Data A",
				"note":              "Note O",
				"tags":              []string{"tag a", "tag b", "tag c", "tag d"},
				"lastModifiedAt":    int64(12),
			}},
		},
	}

	util.UseFakeClock()
	util.ResetFakeClock()

	for i, tt := range tests {
		result := constructUpdateFromMark(tt.original, tt.update)
		assert.Equal(t, tt.result, result, fmt.Sprintf("Invalid result for test case #%d", i))
	}
}
