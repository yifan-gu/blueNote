/*
Copyright © 2022 Yifan Gu <guyifan1121@gmail.com>

*/

package json

import (
	jsonenc "encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yifan-gu/blueNote/pkg/config"
	"github.com/yifan-gu/blueNote/pkg/model"
)

var books []Book

type Book struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Marks  []Mark `json:"marks"`
}

type Mark struct {
	Type      string   `json:"type"`
	Title     string   `json:"title"`
	Author    string   `json:"author"`
	Section   string   `json:"section"`
	Location  Location `json:"location"`
	Data      string   `json:"data"`
	UserNotes string   `json:"notes,omitempty"`
}

type Location struct {
	Chapter  string `json:"chapter"`
	Page     *int   `json:"page,omitempty"`
	Location *int   `json:"location,omitempty"`
}

type JSONExporter struct {
	prettyPrint bool
	indent      string
}

func (e *JSONExporter) Name() string {
	return "json"
}

func (e *JSONExporter) LoadConfigs(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVar(&e.prettyPrint, "json.pretty", false, "print the json with indent")
	cmd.PersistentFlags().StringVar(&e.indent, "json.indent", "  ", "sets the json indent")
}

func (e *JSONExporter) Export(cfg *config.Config, book *model.Book) error {
	bk := Book{
		Title:  book.Title,
		Author: book.Author,
	}
	for _, mk := range book.Marks {
		loc := Location{
			Chapter: mk.Location.Chapter,
		}
		if mk.Location.Page > 0 {
			loc.Page = &mk.Location.Page
		}
		if mk.Location.Location > 0 {
			loc.Location = &mk.Location.Location
		}

		bk.Marks = append(bk.Marks, Mark{
			Type:      model.MarkTypeString[mk.Type],
			Author:    book.Author,
			Title:     book.Title,
			Section:   mk.Section,
			Data:      mk.Data,
			UserNotes: mk.UserNotes,
			Location:  loc,
		})
	}

	if !cfg.SplitBook {
		return e.exportSingleBook(&bk)
	}

	books = append(books, bk)
	if cfg.CurrentBookIndex == cfg.TotalBookCnt-1 {
		return e.exportMultipleBooks(books)
	}
	return nil
}

func (e *JSONExporter) exportSingleBook(bk *Book) error {
	var b []byte
	var err error
	if e.prettyPrint {
		b, err = jsonenc.MarshalIndent(bk, "", e.indent)
	} else {
		b, err = jsonenc.Marshal(bk)
	}
	if err != nil {
		return fmt.Errorf("failed to marshal json: %v", err)
	}
	fmt.Println(string(b))
	return nil
}

func (e *JSONExporter) exportMultipleBooks(books []Book) error {
	var b []byte
	var err error
	if e.prettyPrint {
		b, err = jsonenc.MarshalIndent(books, "", e.indent)
	} else {
		b, err = jsonenc.Marshal(books)
	}
	if err != nil {
		return fmt.Errorf("failed to marshal json: %v", err)
	}
	fmt.Println(string(b))
	return nil
}
