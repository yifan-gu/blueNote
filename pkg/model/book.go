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

// Book defines the details of a Book object, which also contains a list of marks.
type Book struct {
	Title  string
	Author string
	Marks  []Mark
}

// Mark defines the details of a mark object.
type Mark struct {
	Type      MarkType
	Section   string
	Location  Location
	Data      string
	UserNotes string
}

// Location defines the location of a mark in the book.
type Location struct {
	Chapter  string
	Page     int
	Location int
}
