/*
Copyright © 2022 Yifan Gu <guyifan1121@gmail.com>

*/

package kindlehtml

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/yifan-gu/blueNote/pkg/model"
	"github.com/yifan-gu/blueNote/pkg/util"
	"golang.org/x/net/html"
)

var numberRegexp = regexp.MustCompile(`\d+`)

type KindleHTMLParser struct {
	authorOverride string
	titleOverride  string
	splitBook      bool
}

func (p *KindleHTMLParser) Name() string {
	return "kindle-html"
}

func (p *KindleHTMLParser) LoadConfigs(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&p.authorOverride, "kindle-html.author", "", "override the book author name")
	cmd.PersistentFlags().StringVar(&p.titleOverride, "kindle-html.title", "", "override the book title name")
	cmd.PersistentFlags().BoolVarP(&p.splitBook, "kindle-html.split", "s", false, "split sub-sections into separate books")
}

func (p *KindleHTMLParser) Parse(inputPath string) ([]*model.Book, error) {
	f, err := os.Open(inputPath)
	defer f.Close()
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	buf := bufio.NewReader(f)
	tokenizer := html.NewTokenizer(buf)

	var book model.Book
	var section string

	for {
		tt := tokenizer.Next()
		if tt == html.ErrorToken {
			if tokenizer.Err() == io.EOF {
				break
			}
			return nil, errors.Wrap(tokenizer.Err(), fmt.Sprintf("tokenize error for %q", inputPath))
		}

		token := tokenizer.Token()

		for _, attr := range token.Attr {
			if attr.Key != "class" {
				continue
			}
			switch attr.Val {
			case "bookTitle":
				tokenizer.Next()
				if p.titleOverride != "" {
					book.Title = p.titleOverride
				} else {
					book.Title = strings.TrimSpace(string(tokenizer.Raw()))
				}
			case "authors":
				tokenizer.Next()
				if p.authorOverride != "" {
					book.Author = p.authorOverride
				} else {
					book.Author = strings.Join(strings.Fields(strings.TrimSpace(string(tokenizer.Raw()))), ".")
				}
			case "sectionHeading":
				tokenizer.Next()
				section = strings.TrimSpace(string(tokenizer.Raw()))
			case "noteHeading":
				if err := handleNoteEntry(tokenizer, &book, section); err != nil {
					return nil, errors.Wrap(err, fmt.Sprintf("failed to handle notes for %q", inputPath))
				}
			}
		}
	}

	if p.splitBook {
		return splitBook(&book), nil
	}
	return []*model.Book{&book}, nil
}

// splitBook will turn a book into multiple  books.
// It's useful when the input is a book collection.
func splitBook(bk *model.Book) []*model.Book {
	var books []*model.Book
	var sectionTitles []string
	sectionMap := make(map[string][]*model.Mark)

	for _, mk := range bk.Marks {
		if mk.Section != "" {
			if _, ok := sectionMap[mk.Section]; !ok {
				sectionTitles = append(sectionTitles, mk.Section)
			}
			sectionMap[mk.Section] = append(sectionMap[mk.Section], &model.Mark{
				Type:      mk.Type,
				Title:     mk.Section,
				Author:    bk.Author,
				Section:   mk.Location.Chapter,
				Location:  mk.Location,
				Data:      mk.Data,
				UserNotes: mk.UserNotes,
			})
		}
	}

	for _, sectionTitle := range sectionTitles {
		books = append(books, &model.Book{
			Title:  sectionTitle,
			Author: bk.Author,
			Marks:  sectionMap[sectionTitle],
		})
	}
	if len(books) == 0 {
		books = []*model.Book{bk}
	}

	return books
}

func handleNextText(tokenizer *html.Tokenizer, f func(tokenizer *html.Tokenizer)) error {
	for {
		tt := tokenizer.Next()
		if tt == html.ErrorToken {
			return errors.Wrap(tokenizer.Err(), "")
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

func parseLocationWithoutChapter(data []byte) *model.Location {
	var page, location []byte
	var loc model.Location

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
	match := numberRegexp.FindSubmatch(page)
	if len(match) == 1 {
		pg, err := strconv.Atoi(string(match[0]))
		if err != nil {
			util.Fatal("Cannot parse page info", err)
		}
		loc.Page = &pg
	}

	match = numberRegexp.FindSubmatch(location)
	if len(match) == 1 {
		lc, err := strconv.Atoi(string(match[0]))
		if err != nil {
			util.Fatal("Cannot parse location info", err)
		}
		loc.Location = &lc
	}
	return &loc
}

func parseLocationWithChapter(chapterData, data []byte) *model.Location {
	chapter := bytes.TrimLeft(chapterData, ") -")
	chapter = bytes.TrimSpace(chapter)

	loc := parseLocationWithoutChapter(data)
	loc.Chapter = string(chapter)

	return loc
}

func parseLocation(data []byte) *model.Location {

	tuples := bytes.Split(data, []byte(">"))
	switch len(tuples) {
	case 1:
		return parseLocationWithoutChapter(tuples[0])
	case 2:
		return parseLocationWithChapter(tuples[0], tuples[1])
	default:
		util.Fatal(fmt.Sprintf("unexpected location format: %s", data))
		return nil
	}
}

func handleHighlight(tokenizer *html.Tokenizer, book *model.Book, section string) {
	mk := &model.Mark{
		Type:    model.MarkTypeHighlight,
		Title:   book.Title,
		Author:  book.Author,
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

func handleNote(tokenizer *html.Tokenizer, book *model.Book, section string) {
	mk := &model.Mark{
		Type:     model.MarkTypeNote,
		Title:    book.Title,
		Author:   book.Author,
		Section:  section,
		Location: parseLocation(tokenizer.Raw()),
	}

	handleNextText(tokenizer, func(tokenizer *html.Tokenizer) {
		mk.Data = book.Marks[len(book.Marks)-1].Data
		mk.UserNotes = string(tokenizer.Raw())
	})
	book.Marks = append(book.Marks, mk)
}

func handleNoteEntry(tokenizer *html.Tokenizer, book *model.Book, section string) error {
	return handleNextText(tokenizer, func(tokenizer *html.Tokenizer) {
		switch {
		case strings.HasPrefix(string(tokenizer.Raw()), "Highlight"):
			handleHighlight(tokenizer, book, section)
		case strings.HasPrefix(string(tokenizer.Raw()), "Note"):
			handleNote(tokenizer, book, section)
		}
	})
}
