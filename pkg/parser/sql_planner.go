package parser

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/template"

	"github.com/yifan-gu/klipping2org/pkg/db"
	"github.com/yifan-gu/klipping2org/pkg/util"
)

type SqlPlanner struct {
	driver db.SqlInterface
	sqls   []*db.SQL
}

func NewSqlPlanner(driver db.SqlInterface) *SqlPlanner {
	return &SqlPlanner{driver: driver}
}

func (s *SqlPlanner) InsertNodeLinkTitleEntry(book *Book, roamDBPath, outputPath string) error {
	return s.insertNodeLinkEntry(book, roamDBPath, outputPath, book.UUID.String(), "", 0, 1)
}

func (s *SqlPlanner) InsertNodeLinkMarkEntry(book *Book, mark *Mark, roamDBPath, outputPath string) error {
	return s.insertNodeLinkEntry(book, roamDBPath, outputPath, mark.UUID.String(), mark.Data, 1, mark.Pos)
}

func (s *SqlPlanner) InsertFileEntry(book *Book, fullpath string) error {
	hash, err := computeHash(fullpath)
	if err != nil {
		return err
	}

	atime, err := getAtime(fullpath)
	if err != nil {
		return err
	}
	mtime, err := getMtime(fullpath)
	if err != nil {
		return err
	}

	s.sqls = append(s.sqls, &db.SQL{
		Statement: "INSERT INTO files (file, title, hash, atime, mtime) VALUES (?, ?, ?, ?, ?)",
		Values:    []interface{}{quoteString(fullpath), quoteString(book.Title), quoteString(hash), atime, mtime},
	})
	return nil
}

func computeHash(fullpath string) (string, error) {
	f, err := os.Open(fullpath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func quoteString(str string) string {
	return fmt.Sprintf("%q", str)
}

func (s *SqlPlanner) insertNodeLinkEntry(book *Book, roamDBPath, outputPath, uuid, data string, level, pos int) error {
	properties, err := generateProperties(book, uuid, outputPath, data)
	if err != nil {
		return err
	}

	fullpath, err := util.ResolvePath(outputPath)
	if err != nil {
		return fmt.Errorf("failed to get full path of %s: %v", outputPath, err)
	}

	s.sqls = append(s.sqls, &db.SQL{
		Statement: "INSERT INTO nodes (id, file, level, pos, title, properties) VALUES(?, ?, ?, ?, ?, ?)",
		Values:    []interface{}{quoteString(uuid), quoteString(fullpath), level, pos, quoteString(data), properties},
	})
	s.sqls = append(s.sqls, &db.SQL{
		Statement: "INSERT INTO tags (node_id, tag) VALUES(?, ?)",
		Values:    []interface{}{quoteString(uuid), quoteString(book.Author)},
	})
	return nil
}

func (s *SqlPlanner) CommitSql() error {
	return s.driver.CommitTransaction(s.sqls)
}

func generateProperties(book *Book, uuid, outputPath, data string) (string, error) {
	propertyTpl := `(("CATEGORY" . "{{ .Filename }}") ("ID" . "{{ .UUID }}") ("BLOCKED" . "") ("ALLTAGS" . #(":{{ .Author }}:" 1 {{ .AuthorEndPos }} (inherited t))) ("FILE" . "{{ .Fullpath }}") ("PRIORITY" . "B")`
	extraPropertyFormatTempalte := `("ITEM" . "{{ .Data }}")`
	propertyTplTailingString := `)`

	fullpath, err := util.ResolvePath(outputPath)
	if err != nil {
		return "", fmt.Errorf("failed to get full path of %s: %v", outputPath, err)
	}

	v := struct {
		Filename     string
		UUID         string
		Author       string
		AuthorEndPos int
		Fullpath     string
		Data         string
	}{
		filepath.Base(outputPath),
		uuid,
		book.Author,
		len([]rune(book.Author)) + 1,
		fullpath,
		data,
	}

	if data != "" {
		propertyTpl = propertyTpl + extraPropertyFormatTempalte
	}
	propertyTpl = propertyTpl + propertyTplTailingString

	buf := new(bytes.Buffer)
	tpl := template.Must(template.New("template").Parse(propertyTpl))
	if err := tpl.Execute(buf, &v); err != nil {
		return "", fmt.Errorf("failed to execute peroperty template: %v", err)
	}

	return buf.String(), nil
}

func generateTitleProperties(book *Book, uuid, outputPath string) (string, error) {
	return generateProperties(book, uuid, outputPath, "")
}

func generateMarkProperties(book *Book, uuid, outputPath, data string) (string, error) {
	return generateProperties(book, uuid, outputPath, data)
}
