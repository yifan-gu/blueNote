/*
Copyright © 2022 Yifan Gu <guyifan1121@gmail.com>

*/

package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/google/uuid"

	"github.com/yifan-gu/klipping2org/pkg/config"
	"github.com/yifan-gu/klipping2org/pkg/db"
	"github.com/yifan-gu/klipping2org/pkg/util"
)

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
	if cfg.AuthorSubDir {
		return filepath.Join(cfg.OutputDir, b.Author, filename)
	}
	return filepath.Join(cfg.OutputDir, filename)
}

func (b *Book) FormatOrg(sp SqlPlanner, cfg *config.Config) ([]byte, error) {
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

	b.UUID = uuid.New()
	buf := new(bytes.Buffer)
	titleT := template.Must(template.New("template").Parse(orgTitleTpl))
	if err := titleT.Execute(buf, b); err != nil {
		return nil, fmt.Errorf("failed to execute org template for title: %v", err)
	}

	if err := sp.InsertNodeLinkTitleEntry(b, generateOutputPath(b, cfg)); err != nil {
		return nil, err
	}

	for _, mk := range b.Marks {
		mk.UUID = uuid.New()
		mk.Pos = len([]rune(buf.String())) + len("\n*")

		if err := sp.InsertNodeLinkMarkEntry(b, mk, generateOutputPath(b, cfg)); err != nil {
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

	parser, err := NewParser(cfg.Parser)
	if err != nil {
		return err
	}

	book, err := parser.Parse(inputPath)
	if err != nil {
		return err
	}

	if cfg.SplitBook {
		books = book.Split()
	} else {
		books = []*Book{book}
	}

	for _, bk := range books {
		sp := NewSqlPlanner(sq, cfg.UpdateRoamDB)
		b, err := bk.FormatOrg(sp, cfg)
		if err != nil {
			return err
		}

		// To fix non-English encoding issue.
		r := []rune(string(b))

		fullpath, err := util.ResolvePath(generateOutputPath(bk, cfg))
		if err != nil {
			return err
		}
		dir := filepath.Dir(fullpath)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			confirm, err := PromptConfirmation(cfg, fmt.Sprintf("directory %s doesn't exit, create?", dir))
			if err != nil {
				return err
			}
			if confirm {
				if err := os.MkdirAll(dir, 0755); err != nil {
					return fmt.Errorf("failed to create dir %q: %v", dir, err)
				}
			}
		}

		if _, err := os.Stat(fullpath); err == nil || !os.IsNotExist(err) {
			confirm, err := PromptConfirmation(cfg, fmt.Sprintf("file %s already exits, replace?", fullpath))
			if err != nil {
				return err
			}
			if !confirm {
				continue
			}
		}

		if err := writeRunesToFile(fullpath, r); err != nil {
			return err
		}

		if err := sp.InsertFileEntry(bk, fullpath); err != nil {
			return err
		}

		if err := sp.CommitSql(); err != nil {
			return err
		}

		fmt.Println("Successfully created:", fullpath)
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
