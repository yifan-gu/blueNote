package parser

type dummySqlPlanner struct{}

func (s *dummySqlPlanner) InsertNodeLinkTitleEntry(book *Book, outputPath string) error {
	return nil
}

func (s *dummySqlPlanner) InsertNodeLinkMarkEntry(book *Book, mark *Mark, outputPath string) error {
	return nil
}

func (s *dummySqlPlanner) InsertFileEntry(book *Book, fullpath string) error {
	return nil
}

func (s *dummySqlPlanner) CommitSql() error {
	return nil
}
