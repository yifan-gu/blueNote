package parser

type mockSqlPlanner struct{}

func (s *mockSqlPlanner) InsertNodeLinkTitleEntry(book *Book, outputPath string) error {
	return nil
}

func (s *mockSqlPlanner) InsertNodeLinkMarkEntry(book *Book, mark *Mark, outputPath string) error {
	return nil
}

func (s *mockSqlPlanner) InsertFileEntry(book *Book, fullpath string) error {
	return nil
}

func (s *mockSqlPlanner) CommitSql() error {
	return nil
}
