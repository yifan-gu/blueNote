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
	"golang.org/x/net/html"
)

var numberRegexp = regexp.MustCompile(`\d+`)

type KindleHTMLParser struct {
	author    string
	title     string
	splitBook bool
}

func (p *KindleHTMLParser) Name() string {
	return "kindle-html"
}

func (p *KindleHTMLParser) LoadConfigs(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&p.author, "kindlehtml.author", "", "override the book author name")
	cmd.PersistentFlags().StringVar(&p.title, "kindlehtml.title", "", "override the book title name")
	cmd.PersistentFlags().BoolVarP(&p.splitBook, "kindlehtml.split", "s", false, "split sub-sections into separate books")
}

func (p *KindleHTMLParser) Parse(inputPath string) (*model.Book, error) {
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
				book.Title = strings.TrimSpace(string(tokenizer.Raw()))
			case "authors":
				tokenizer.Next()
				book.Author = strings.Join(strings.Fields(strings.TrimSpace(string(tokenizer.Raw()))), ".")
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
	return &book, nil
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

func parseLocationWithoutChapter(data []byte) model.Location {
	var err error
	var pg, lc = -1, -1
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
	match := numberRegexp.FindSubmatch(page)
	if len(match) == 1 {
		pg, err = strconv.Atoi(string(match[0]))
		if err != nil {
			panic(err)
		}
	}

	match = numberRegexp.FindSubmatch(location)
	if len(match) == 1 {
		lc, err = strconv.Atoi(string(match[0]))
		if err != nil {
			panic(err)
		}
	}
	return model.Location{Page: pg, Location: lc}
}

func parseLocationWithChapter(chapterData, data []byte) model.Location {
	chapter := bytes.TrimLeft(chapterData, ") -")
	chapter = bytes.TrimSpace(chapter)

	loc := parseLocationWithoutChapter(data)
	loc.Chapter = string(chapter)

	return loc
}

func parseLocation(data []byte) model.Location {

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

func handleHighlight(tokenizer *html.Tokenizer, book *model.Book, section string) {
	mk := model.Mark{
		Type:    model.MarkTypeHighlight,
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
	mk := model.Mark{
		Type:     model.MarkTypeNote,
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
