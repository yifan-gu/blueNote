/*
Copyright © 2022 Yifan Gu <guyifan1121@gmail.com>

*/

package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/google/uuid"
	"golang.org/x/net/html"

	"github.com/yifan-gu/kliping2org/pkg/config"
	"github.com/yifan-gu/kliping2org/pkg/db"
)

// TODO: chapter data, sectionHeading (done)
// struct of the entry (done)
// parse location (done)
// parse page number (done)
// write to file (done)
// parse multiple books (done)
//
// compute title (done)
// iterate entries, figure out pos val (done)
// insert database
// cmdline argument first
//
// optional roam db
// optional author dir
// check roam version

// default org roam dir (done)
// type for note (no-op, because all treated as level 1 headlines, done)
// ask for skip, replace all
//
// refactor book module, insert commit transaction
//

type MarkType int

const (
	MarkTypeHighlight MarkType = iota
	MarkTypeNote
)

var MarkTypeString = map[MarkType]string{
	MarkTypeHighlight: "Highlight",
	MarkTypeNote:      "Note",
}

type Location struct {
	Chapter  string
	Page     string
	Location string
}

func (l *Location) String() string {
	if l.Chapter != "" {
		return fmt.Sprintf("Chapter: %s Page: %s Loc: %s", l.Chapter, l.Page, l.Location)
	}
	return fmt.Sprintf("Page: %s Loc: %s", l.Page, l.Location)
}

type Mark struct {
	Type     string
	Section  string
	Location *Location
	Data     string
	Pos      int
	UUID     uuid.UUID
}

type Book struct {
	Title  string
	Author string
	Marks  []*Mark
	UUID   uuid.UUID
}

func generateOutputPath(b *Book, cfg *config.Config) string {
	filename := fmt.Sprintf("《%s》 by %s.org", b.Title, b.Author)
	return filepath.Join(cfg.OutputDir, filename)
}

func (b *Book) FormatOrg(sp *SqlPlanner, cfg *config.Config) ([]byte, error) {
	const orgTitleTpl = `:PROPERTIES:
:ID:       {{ .UUID }}
:END:
#+title: {{ .Title }}
#+filetags: :{{ .Author }}:

`
	const orgEntryTpl = `
* {{ .Data }}
:PROPERTIES:
:ID:       {{ .UUID }}
:END:
{{ .Type }} @
{{- if eq .Location.Chapter "" }}
Chapter: {{ .Section }}
{{ else }}
Section: {{ .Section }}
{{ end -}}
{{ .Location }}
`

	buf := new(bytes.Buffer)
	titleT := template.Must(template.New("template").Parse(orgTitleTpl))
	if err := titleT.Execute(buf, b); err != nil {
		return nil, fmt.Errorf("failed to execute org template for title: %v", err)
	}
	b.UUID = uuid.New()

	if err := sp.InsertNodeLinkTitleEntry(b, cfg.RoamDBPath, generateOutputPath(b, cfg)); err != nil {
		return nil, err
	}

	for _, mk := range b.Marks {
		mk.UUID = uuid.New()
		mk.Pos = len([]rune(buf.String())) + len("\n*")

		if err := sp.InsertNodeLinkMarkEntry(b, mk, cfg.RoamDBPath, generateOutputPath(b, cfg)); err != nil {
			return nil, err
		}

		entryT := template.Must(template.New("template").Parse(orgEntryTpl))
		if err := entryT.Execute(buf, mk); err != nil {
			return nil, fmt.Errorf("failed to execute org template for entries: %v", err)
		}
	}

	return buf.Bytes(), nil
}

func (b *Book) Split() []*Book {
	var books []*Book
	var sectionTitles []string
	sectionMap := make(map[string][]*Mark)

	for _, mk := range b.Marks {
		if mk.Section != "" {
			loc := &Location{
				Page:     mk.Location.Page,
				Location: mk.Location.Location,
			}
			if _, ok := sectionMap[mk.Section]; !ok {
				sectionTitles = append(sectionTitles, mk.Section)
			}
			sectionMap[mk.Section] = append(sectionMap[mk.Section], &Mark{
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

func handleNextText(tokenizer *html.Tokenizer, f func(tokenizer *html.Tokenizer)) error {
	for {
		tt := tokenizer.Next()
		if tt == html.ErrorToken {
			return tokenizer.Err()
		}
		if tokenizer.Token().Type == html.TextToken {
			break
		}
	}
	if f != nil {
		f(tokenizer)
	}
	return nil
}

func trimLocationString(location []byte) []byte {
	for i := range location {
		if location[i] < '0' || location[i] > '9' {
			return location[0:i]
		}
	}
	return location
}

func parseLocationWithoutChapter(data []byte) *Location {
	var page, location []byte

	pageMarker, locMarker := []byte("Page"), []byte("Location")
	tuples := bytes.Fields(data)
	for i, tp := range tuples {
		switch {
		case bytes.Equal(tp, pageMarker):
			page = tuples[i+1]
		case bytes.Equal(tp, locMarker):
			location = tuples[i+1]
		}
	}
	location = trimLocationString(location)

	return &Location{Page: string(page), Location: string(location)}
}

func parseLocationWithChapter(chapterData, data []byte) *Location {
	chapter := bytes.TrimLeft(chapterData, ") -")
	chapter = bytes.TrimSpace(chapter)

	loc := parseLocationWithoutChapter(data)
	loc.Chapter = string(chapter)

	return loc
}

func parseLocation(data []byte) *Location {

	tuples := bytes.Split(data, []byte(">"))
	switch len(tuples) {
	case 1:
		return parseLocationWithoutChapter(tuples[0])
	case 2:
		return parseLocationWithChapter(tuples[0], tuples[1])
	default:
		panic(fmt.Sprintf("unexpected location format: %s", data))
	}

}

func handleHighlight(tokenizer *html.Tokenizer, book *Book, section string) {
	mk := &Mark{
		Type:    MarkTypeString[MarkTypeHighlight],
		Section: section,
	}

	handleNextText(tokenizer, nil)
	handleNextText(tokenizer, func(tokenizer *html.Tokenizer) {
		mk.Location = parseLocation(tokenizer.Raw())
	})
	handleNextText(tokenizer, func(tokenizer *html.Tokenizer) {
		mk.Data = string(tokenizer.Raw())
	})
	book.Marks = append(book.Marks, mk)
}

func handleNote(tokenizer *html.Tokenizer, book *Book, section string) {
	mk := &Mark{
		Type:     MarkTypeString[MarkTypeNote],
		Section:  section,
		Location: parseLocation(tokenizer.Raw()),
	}

	handleNextText(tokenizer, func(tokenizer *html.Tokenizer) {
		mk.Data = string(tokenizer.Raw())
	})
	book.Marks = append(book.Marks, mk)
}

func handleNoteEntry(tokenizer *html.Tokenizer, book *Book, section string) error {
	return handleNextText(tokenizer, func(tokenizer *html.Tokenizer) {
		switch {
		case strings.HasPrefix(string(tokenizer.Raw()), "Highlight"):
			handleHighlight(tokenizer, book, section)
		case strings.HasPrefix(string(tokenizer.Raw()), "Note"):
			handleNote(tokenizer, book, section)
		}
	})
}

func parseBook(inputPath string) (*Book, error) {
	f, err := os.Open(inputPath)
	defer f.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to open %q: %v", inputPath, err)
	}

	buf := bufio.NewReader(f)
	tokenizer := html.NewTokenizer(buf)

	var book Book
	var section string

	for {
		tt := tokenizer.Next()
		if tt == html.ErrorToken {
			if tokenizer.Err() == io.EOF {
				break
			}
			return nil, fmt.Errorf("tokenize error for %q: %v", inputPath, tokenizer.Err())
		}

		token := tokenizer.Token()

		for _, attr := range token.Attr {
			if attr.Key != "class" {
				continue
			}
			switch attr.Val {
			case "bookTitle":
				tokenizer.Next()
				book.Title = strings.TrimSpace(string(tokenizer.Raw()))
			case "authors":
				tokenizer.Next()
				book.Author = strings.Join(strings.Fields(strings.TrimSpace(string(tokenizer.Raw()))), ".")
			case "sectionHeading":
				tokenizer.Next()
				section = strings.TrimSpace(string(tokenizer.Raw()))
			case "noteHeading":
				if err := handleNoteEntry(tokenizer, &book, section); err != nil {
					return nil, fmt.Errorf("failed to handle notes for %q: %v", inputPath, err)
				}
			}
		}
	}
	return &book, nil
}

func writeRunesToFile(fullpath string, runes []rune) error {
	f, err := os.OpenFile(fullpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open or create file %s: %v", fullpath, err)
	}
	defer f.Close()

	buf := bufio.NewWriter(f)
	defer buf.Flush()

	for i := range runes {
		_, err = fmt.Fprintf(buf, "%c", runes[i])
		if err != nil {
			return fmt.Errorf("failed to write to file %s: %v", fullpath, err)
		}
	}
	return nil
}

func parseAndWrite(inputPath string, cfg *config.Config) error {
	var books []*Book

	sq, err := db.NewSqlInterface(cfg.RoamDBPath, cfg.DBDriver)
	if err != nil {
		return err
	}
	defer sq.Close()

	book, err := parseBook(inputPath)
	if err != nil {
		return err
	}

	if cfg.SplitBook {
		books = book.Split()
	} else {
		books = []*Book{book}
	}

	for _, bk := range books {
		sp := NewSqlPlanner(sq)
		b, err := bk.FormatOrg(sp, cfg)
		if err != nil {
			return err
		}

		// To fix non-English encoding issue.
		r := []rune(string(b))
		fullpath := generateOutputPath(bk, cfg)
		if err := writeRunesToFile(fullpath, r); err != nil {
			return err
		}
		if err := sp.CommitSql(); err != nil {
			return err
		}
	}

	return nil
}

func ParseAndWrite(cfg *config.Config) error {
	f, err := os.Open(cfg.InputPath)
	defer f.Close()

	if err != nil {
		return fmt.Errorf("failed to open %q: %v", cfg.InputPath, err)
	}

	fi, err := f.Stat()
	if err != nil {
		return fmt.Errorf("failed to read stat %q: %v", cfg.InputPath, err)
	}

	if fi.IsDir() {
		if err := filepath.Walk(cfg.InputPath, func(path string, file os.FileInfo, err error) error {
			return parseAndWrite(path, cfg)
		}); err != nil {
			return err
		}

		return nil
	}
	return parseAndWrite(cfg.InputPath, cfg)
}
