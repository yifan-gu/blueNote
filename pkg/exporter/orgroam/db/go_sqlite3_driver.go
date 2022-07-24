package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/yifan-gu/BlueNote/pkg/util"
)

type GoSqlite3Driver struct {
	dbPath string
}

func newGoSqlite3Driver(dbPath string) (*GoSqlite3Driver, error) {
	fullpath, err := util.ResolvePath(dbPath)
	if err != nil {
		return nil, err
	}
	return &GoSqlite3Driver{dbPath: fullpath}, nil
}

func (s *GoSqlite3Driver) CommitTransaction(sqls []*SQL) error {
	db, err := sql.Open("sqlite3", s.dbPath)
	if err != nil {
		return fmt.Errorf("failed to open sqlite3 database for %s: %v", s.dbPath, err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	for _, sq := range sqls {
		stmt, err := db.Prepare(sq.Statement)
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return fmt.Errorf("failed to prepare statement for %q: %v, unable to rollback: %v", sq.Statement, err, rollbackErr)
			}
			return fmt.Errorf("failed to prepare statement for %q: %v", sq.Statement, err)
		}

		_, err = tx.Stmt(stmt).Exec(sq.Values...)
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return fmt.Errorf("failed to exec statement for %q: %v, unable to rollback: %v", sq, err, rollbackErr)
			}
			return fmt.Errorf("failed to exec statement for %q: %v", sq, err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func (s *GoSqlite3Driver) Close() error {
	return nil
}
