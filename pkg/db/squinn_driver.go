package db

import (
	"fmt"
	//"strings"

	"github.com/cvilsmeier/sqinn-go/sqinn"
)

const (
	sqlite3BeginTransactionSQL  = "BEGIN TRANSACTION"
	sqlite3CommitTransactionSQL = "COMMIT"
)

type SqinnWrapper struct {
	*sqinn.Sqinn
	dbPath string
}

func newSqinnWrapper(dbPath string) (*SqinnWrapper, error) {
	sq, err := sqinn.Launch(sqinn.Options{})
	if err != nil {
		return nil, fmt.Errorf("failed to launch sqinn sqlite3 driver: %v", err)
	}
	return &SqinnWrapper{sq, dbPath}, nil
}

func createCommitSql(sqls []*SQL) string {
	//sqls = append([]string{sqlite3BeginTransactionSQL}, sqls...)
	//sqls = append(sqls, sqlite3CommitTransactionSQL)
	//return strings.Join(sqls, ";\n")
	return ""
}

func (s *SqinnWrapper) CommitTransaction(sqls []*SQL) error {
	//if err := s.Open(s.dbPath); err != nil {
	//	return fmt.Errorf("failed to open sqlite3 database for %s: %v", s.dbPath, err)
	//}
	//defer s.Close()
	//
	//if _, err := s.ExecOne(createCommitSql(sqls)); err != nil {
	//	return fmt.Errorf("failed to execute sql insertion: %v", err)
	//}
	return nil
}

func (s *SqinnWrapper) Close() error {
	return s.Terminate()
}
