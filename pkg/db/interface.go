package db

import (
	"fmt"
)

type SqlInterface interface {
	CommitTransaction([]*SQL) error
	Close() error
}

const (
	SqlDriverSqilite3 = "sqlite3"
	SqlDriverMock     = "mock"
)

type SQL struct {
	Statement string
	Values    []interface{}
}

func (s *SQL) String() string {
	return fmt.Sprint(s.Statement, s.Values)
}

func NewSqlInterface(dbPath, driver string) (SqlInterface, error) {
	switch driver {
	case SqlDriverSqilite3:
		return newGoSqlite3Driver(dbPath)
	case SqlDriverMock:
		return newSqlite3Mock(dbPath)
	default:
		return nil, fmt.Errorf("unrecognized sql driver type: %q", driver)
	}
}
