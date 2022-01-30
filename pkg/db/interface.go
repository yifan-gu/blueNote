package db

import (
	"fmt"
)

type SqlInterface interface {
	CommitTransaction([]string) error
	Close() error
}

const (
	SqlDriverSqinn = "sqinn"
	SqlDriverMock  = "mock"
)

func NewSqlInterface(dbPath, driver string) (SqlInterface, error) {
	switch driver {
	case SqlDriverSqinn:
		return newSqinnWrapper(dbPath)
	case SqlDriverMock:
		return newSqlite3Mock(dbPath)
	default:
		return nil, fmt.Errorf("unrecognized sql driver type: %q", driver)
	}
}
