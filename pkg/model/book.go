/*
Copyright © 2022 Yifan Gu <guyifan1121@gmail.com>

*/

package model

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

const (
	// MarkTypeHighlight is a highlight marking.
	MarkTypeHighlight = "HIGHLIGHT"
	// MarkTypeNote is a note marking.
	MarkTypeNote = "NOTE"
	// MarkTypeBookmark is a bookmark marking.
	MarkTypeBookmark = "BOOKMARK"
)

var (
	typeMaps = map[string]struct{}{
		MarkTypeHighlight: struct{}{},
		MarkTypeNote:      struct{}{},
		MarkTypeBookmark:  struct{}{},
	}
)

// Book defines the details of a Book object, which also contains a list of marks.
type Book struct {
	BookID    string  `json:"bookId,omitempty"` // Unique ID for the book, this is a better way to identify a book than using its title.
	Title     string  `json:"title"`
	Author    string  `json:"author"` // TODO(yifan): Change to Authors
	Publisher string  `json:"publisher,omitempty"`
	Date      string  `json:"date,omitempty"`
	ISBN      string  `json:"isbn,omitempty"`
	Marks     []*Mark `json:"marks,omitempty"`
}

// Mark defines the details of a mark object.
type Mark struct {
	ID             string    `json:"id,omitempty"`
	BookID         string    `json:"bookId,omitempty"` // Unique ID for the book, also the foreign key to the book info in the database.
	Type           string    `json:"type"`
	Title          string    `json:"title"`
	Author         string    `json:"author"`
	Section        string    `json:"section,omitempty"`
	Location       *Location `json:"location,omitempty"`
	Data           string    `json:"data,omitempty"`
	UserNote       string    `json:"note,omitempty"`
	Tags           []string  `json:"tags,omitempty"`
	CreatedAt      *int64    `json:"createdAt,omitempty"`
	LastModifiedAt *int64    `json:"lastModifiedAt,omitempty"`
}

// Location defines the location of a mark in the book.
type Location struct {
	Chapter  string `json:"chapter,omitempty"`
	Page     *int   `json:"page,omitempty"`
	Location *int   `json:"location,omitempty"`
}

func isSupportedType(typ string) bool {
	_, ok := typeMaps[typ]
	return ok
}

func ValidateMark(m *Mark) error {
	if !isSupportedType(m.Type) {
		return errors.New(fmt.Sprintf("Type %v is not supported", m.Type))
	}
	if m.Data == "" && m.UserNote == "" {
		return errors.New("Expect 'data' or 'note' to be set")
	}
	return nil
}

// ByTitle is a type that implements sort.Interface for []*Book based on the Title field.
type ByTitle []*Book

// Len is the number of elements in the collection.
func (b ByTitle) Len() int {
	return len(b)
}

// Less reports whether the element at index i should sort before the element at index j.
func (b ByTitle) Less(i, j int) bool {
	// Compare titles case-insensitively
	return strings.ToLower(b[i].Title) < strings.ToLower(b[j].Title)
}

// Swap swaps the elements at indexes i and j.
func (b ByTitle) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

// SortBooksByTitle sorts a slice of books by title.
func SortBooksByTitle(books []*Book) {
	sort.Sort(ByTitle(books))
}
