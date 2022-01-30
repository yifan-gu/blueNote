package db

type mockSqlite3Interface struct {
	sql string
	err error
}

func newSqlite3Mock(dbPath string) (*mockSqlite3Interface, error) {
	return &mockSqlite3Interface{}, nil
}

func (s *mockSqlite3Interface) CommitTransaction(sqls []string) error {
	if s.err != nil {
		return s.err
	}
	s.sql = createCommitSql(sqls)
	return nil
}

func (s *mockSqlite3Interface) Close() error {
	return s.err
}

func (s *mockSqlite3Interface) SetError(err error) {
	s.err = err
}

func (s *mockSqlite3Interface) GetExecutedSql() string {
	return s.sql
}
