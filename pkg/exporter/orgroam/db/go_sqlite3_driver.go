package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"github.com/yifan-gu/blueNote/pkg/util"
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
		return errors.Wrap(err, fmt.Sprintf("failed to open sqlite3 database for %s", s.dbPath))
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}

	for _, sq := range sqls {
		stmt, err := db.Prepare(sq.Statement)
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return errors.Wrap(fmt.Errorf("failed to prepare statement for %q: %v, unable to rollback: %v", sq.Statement, err, rollbackErr), "")
			}
			return errors.Wrap(err, fmt.Sprintf("failed to prepare statement for %q", sq.Statement))
		}

		_, err = tx.Stmt(stmt).Exec(sq.Values...)
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return errors.Wrap(fmt.Errorf("failed to exec statement for %q: %v, unable to rollback: %v", sq, err, rollbackErr), "")
			}
			return errors.Wrap(err, fmt.Sprintf("failed to exec statement for %q", sq))
		}
	}
	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}

func (s *GoSqlite3Driver) Close() error {
	return nil
}
