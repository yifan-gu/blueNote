package parser

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestInsertNodeLinkTitleEntry(t *testing.T) {
	testUUID := uuid.MustParse("4b3d008c-e1a5-4227-bbaf-3fd23a858114")
	book := &Book{
		Title:  "书名",
		Author: "作者",
		UUID:   testUUID,
	}
	roamDBPath := "/path/to/db"
	outputPath := "/path/to/output.org"

	sp := &SqlPlanner{}
	err := sp.InsertNodeLinkTitleEntry(book, roamDBPath, outputPath)

	exp := []string{
		`INSERT INTO nodes (id, file, level, pos, title, properties) VALUES(?, ?, ?, ?, ?, ?)["4b3d008c-e1a5-4227-bbaf-3fd23a858114" "/path/to/output.org" 0 1 "" (("CATEGORY" . "output.org") ("ID" . "4b3d008c-e1a5-4227-bbaf-3fd23a858114") ("BLOCKED" . "") ("ALLTAGS" . #(":作者:" 1 3 (inherited t))) ("FILE" . "/path/to/output.org") ("PRIORITY" . "B"))]`,
		`INSERT INTO tags (node_id, tag) VALUES(?, ?)["4b3d008c-e1a5-4227-bbaf-3fd23a858114" "作者"]`,
	}

	assert.NoError(t, err)
	assert.Equal(t, len(exp), len(sp.sqls))

	for i := range exp {
		assert.Equal(t, exp[i], sp.sqls[i].String(), "case #i")
	}
}

func TestInsertNodeLinkMarkEntry(t *testing.T) {
	testUUID := uuid.MustParse("4b3d008c-e1a5-4227-bbaf-3fd23a858114")
	testUUID2 := uuid.MustParse("b2b4cfff-abd1-4432-b7b0-cda98c50e1a1")
	book := &Book{
		Title:  "书名",
		Author: "作者",
		UUID:   testUUID,
	}
	mark := &Mark{
		UUID: testUUID2,
		Data: "这是一段标记",
	}
	roamDBPath := "/path/to/db"
	outputPath := "/path/to/output.org"

	sp := &SqlPlanner{}
	err := sp.InsertNodeLinkMarkEntry(book, mark, roamDBPath, outputPath)

	exp := []string{
		`INSERT INTO nodes (id, file, level, pos, title, properties) VALUES(?, ?, ?, ?, ?, ?)["b2b4cfff-abd1-4432-b7b0-cda98c50e1a1" "/path/to/output.org" 1 0 "这是一段标记" (("CATEGORY" . "output.org") ("ID" . "b2b4cfff-abd1-4432-b7b0-cda98c50e1a1") ("BLOCKED" . "") ("ALLTAGS" . #(":作者:" 1 3 (inherited t))) ("FILE" . "/path/to/output.org") ("PRIORITY" . "B")("ITEM" . "这是一段标记"))]`,
		`INSERT INTO tags (node_id, tag) VALUES(?, ?)["b2b4cfff-abd1-4432-b7b0-cda98c50e1a1" "作者"]`,
	}

	assert.NoError(t, err)
	assert.Equal(t, len(exp), len(sp.sqls))

	for i := range exp {
		assert.Equal(t, exp[i], sp.sqls[i].String(), "case #i")
	}
}
