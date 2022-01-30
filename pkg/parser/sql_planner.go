package parser

import (
	"bytes"
	"fmt"
	"path/filepath"
	"text/template"

	"github.com/yifan-gu/klipping2org/pkg/db"
)

type SqlPlanner struct {
	driver db.SqlInterface
	sql    []string
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

func (s *SqlPlanner) insertNodeLinkEntry(book *Book, roamDBPath, outputPath, uuid, data string, level, pos int) error {
	properties, err := generateProperties(book, uuid, outputPath, data)
	if err != nil {
		return err
	}

	sqlSentence := fmt.Sprintf(`INSERT INTO nodes (id, file, level, pos, todo, priority, scheduled, deadline, title, properties, olp) VALUES (%s, %s, %d, %d, "", "", "", "", %s, %s, "")`,
		uuid, outputPath, level, pos, book.Title, properties)

	s.sql = append(s.sql, sqlSentence)
	return nil
}

func (s *SqlPlanner) CommitSql() error {
	return s.driver.CommitTransaction(s.sql)
}

func generateProperties(book *Book, uuid, outputPath, data string) (string, error) {
	propertyTpl := `(("CATEGORY" . "{{ .Filename }}") ("ID" . "{{ .UUID }}") ("BLOCKED" . "") ("ALLTAGS" . #(":{{ .Author }}:" 1 {{ .AuthorEndPos }} (inherited t))) ("FILE" . "{{ .Fullpath }}") ("PRIORITY" . "B")`
	extraPropertyFormatTempalte := `("ITEM" . "{{ .Data }}")`
	propertyTplTailingString := `)`

	fullpath, err := filepath.Abs(outputPath)
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
