package main

import (
	"database/sql"
	"errors"
	"testing"
)

type SQLiteMock struct {
	Client interface{}
}

func (sqlite *SQLiteMock) Open(dbfile string) (err error) {
	if dbfile == "" {
		return errors.New("dbfile is empty")
	}
	return nil
}

func (sqlite *SQLiteMock) Get(dest interface{}, query string, args ...interface{}) (err error) {

	switch p := dest.(type) {
	case *Trigger:
		*p = Trigger{"BAT", "0x0000000000000000000000000000000000000000",
			"0x0d8775f648430679a709e98d2b0cb6250d2887ef", "balanceOf(address)(uint256)", "", 0,
		}
	default:
		panic("Unexpected type")
	}
	return nil
}

func (sqlite *SQLiteMock) NamedExec(query string, arg interface{}) (sql.Result, error) {
	return nil, nil
}

func (sqlite *SQLiteMock) Exec(query string) (sql.Result, error) {
	return nil, nil
}

func TestMockOpen(t *testing.T) {
	db := SQLiteMock{}
	err := db.Open("hal.xyz")
	ok(t, err)
}
