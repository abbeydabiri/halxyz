package main

import (
	"database/sql"
	"errors"

	//needed for sqlite3
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	Client *sqlx.DB
}

func (sqlite *SQLite) Open(dbfile string) (err error) {
	if dbfile == "" {
		err = errors.New("dbfile is empty")
		return
	}

	if sqlite.Client, err = sqlx.Open("sqlite3", "file:"+dbfile+"?cache=shared"); err != nil {
		return
	}

	err = sqlite.Client.Ping()
	return
}

func (sqlite *SQLite) Get(dest interface{}, query string, args ...interface{}) error {
	err := sqlite.Client.Get(dest, query, args...)
	return err
}

func (sqlite *SQLite) NamedExec(query string, arg interface{}) (sql.Result, error) {
	sqlResult, err := sqlite.Client.NamedExec(query, arg)
	return sqlResult, err
}

func (sqlite *SQLite) Exec(query string) (sql.Result, error) {
	sqlResult, err := sqlite.Client.Exec(query)
	return sqlResult, err
}
