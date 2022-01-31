package db

type mockSqlite3Interface struct {
	sqls []*SQL
	err  error
}

func newSqlite3Mock(dbPath string) (*mockSqlite3Interface, error) {
	return &mockSqlite3Interface{}, nil
}

func (s *mockSqlite3Interface) CommitTransaction(sqls []*SQL) error {
	if s.err != nil {
		return s.err
	}
	s.sqls = sqls
	return nil
}

func (s *mockSqlite3Interface) Close() error {
	return s.err
}

func (s *mockSqlite3Interface) SetError(err error) {
	s.err = err
}

func (s *mockSqlite3Interface) GetExecutedSql() []*SQL {
	return s.sqls
}
