/*
Copyright © 2022 Yifan Gu <guyifan1121@gmail.com>

*/

package model

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/yifan-gu/blueNote/pkg/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	// MarkTypeHighlight is a highlight marking.
	MarkTypeHighlight = "HIGHLIGHT"
	// MarkTypeNote is a note marking.
	MarkTypeNote = "NOTE"
)

// Book defines the details of a Book object, which also contains a list of marks.
type Book struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Marks  []Mark `json:"marks"`
}

// Mark defines the details of a mark object.
type Mark struct {
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	Section   string    `json:"section,omitempty"`
	Location  *Location `json:"location,omitempty"`
	Data      string    `json:"data,omitempty"`
	UserNotes string    `json:"notes,omitempty"`
}

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

// ToPersistenMark converts a mark to PersistentMark
func (m *Mark) ToPersistenMark() *PersistentMark {
	b, err := json.Marshal(m)
	if err != nil {
		util.Fatal("cannot marshal:", err)
	}
	return &PersistentMark{
		ID:      primitive.NewObjectID(),
		Digest:  fmt.Sprintf("%x", sha256.Sum256(b)),
		Type:    m.Type,
		Title:   m.Title,
		Author:  m.Author,
		Section: m.Section,
		Location: &Location{
			Chapter:  m.Location.Chapter,
			Page:     m.Location.Page,
			Location: m.Location.Location,
		},
		Data:      m.Data,
		UserNotes: m.UserNotes,
	}
}
