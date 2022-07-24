/*
Copyright © 2022 Yifan Gu <guyifan1121@gmail.com>

*/

package orgroam

type noopSqlPlanner struct{}

func (s *noopSqlPlanner) InsertNodeLinkTitleEntry(book *Book, outputPath string) error {
	return nil
}

func (s *noopSqlPlanner) InsertNodeLinkMarkEntry(book *Book, mark *Mark, outputPath string) error {
	return nil
}

func (s *noopSqlPlanner) InsertFileEntry(book *Book, fullpath string) error {
	return nil
}

func (s *noopSqlPlanner) CommitSql() error {
	return nil
}
