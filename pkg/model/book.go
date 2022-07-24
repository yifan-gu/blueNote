/*
Copyright © 2022 Yifan Gu <guyifan1121@gmail.com>

*/

package model

// MarkType is the type of the marking
type MarkType string

const (
	// MarkTypeHighlight is a highlight marking.
	MarkTypeHighlight MarkType = "HIGHLIGHT"
	// MarkTypeNote is a note marking.
	MarkTypeNote MarkType = "NOTE"
)

// MarkTypeString is a map from MarkType to their readable string representitives.
var MarkTypeString = map[MarkType]string{
	MarkTypeHighlight: "HIGHLIGHT",
	MarkTypeNote:      "NOTE",
}

// Location defines the location of a mark in the book.
type Location struct {
	Chapter  string
	Page     string
	Location string
}

// Mark defines the details of a mark object.
type Mark struct {
	Type     MarkType
	Section  string
	Location Location
	Data     string
}

// Book defines the details of a Book object, which also contains a list of marks.
type Book struct {
	Title  string
	Author string
	Marks  []Mark
}

// Split will turn a book into multiple  books.
// It's useful when an e-book is a collection.
func (b *Book) Split() []*Book {
	var books []*Book
	var sectionTitles []string
	sectionMap := make(map[string][]Mark)

	for _, mk := range b.Marks {
		if mk.Section != "" {
			loc := Location{
				Page:     mk.Location.Page,
				Location: mk.Location.Location,
			}
			if _, ok := sectionMap[mk.Section]; !ok {
				sectionTitles = append(sectionTitles, mk.Section)
			}
			sectionMap[mk.Section] = append(sectionMap[mk.Section], Mark{
				Type:     mk.Type,
				Section:  mk.Location.Chapter,
				Location: loc,
				Data:     mk.Data,
			})
		}
	}

	for _, sectionTitle := range sectionTitles {
		books = append(books, &Book{
			Title:  sectionTitle,
			Author: b.Author,
			Marks:  sectionMap[sectionTitle],
		})
	}
	if len(books) == 0 {
		books = []*Book{b}
	}

	return books
}
