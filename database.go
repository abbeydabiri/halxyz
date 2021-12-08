package main

import (
	"database/sql"
	"errors"

	//needed for sqlite3
	_ "github.com/mattn/go-sqlite3"
)

// interface to mock sqlx.DB
type SQLClient interface {
	Open(dbfile string) (err error)
	Get(dest interface{}, query string, args ...interface{}) error
	NamedExec(query string, arg interface{}) (sql.Result, error)
	Exec(query string) (sql.Result, error)
}

type Database struct {
	Client SQLClient
}

func (database *Database) Connect(dbfile string) (err error) {
	err = database.Client.Open(dbfile)
	if err == nil {
		err = database.createTableTriggers()
	}
	return
}

func (database *Database) createTableTriggers() error {
	sqlCreateTable := `CREATE TABLE IF NOT EXISTS triggers (
		triggername text NOT NULL,
		useraddress text NOT NULL,
		contractaddress text NOT NULL,
		method text NOT NULL,
		status text NOT NULL,
		number int NOT NULL
	);`
	_, err := database.Client.Exec(sqlCreateTable)
	return err
}

func (database *Database) Find(TriggerName, UserAddress string) (trigger *Trigger, err error) {
	if TriggerName == "" {
		err = errors.New("TriggerName is empty")
		return
	}

	if UserAddress == "" {
		err = errors.New("UserAddress is empty")
		return
	}
	trigger = &Trigger{}
	err = database.Client.Get(trigger, "SELECT * FROM triggers WHERE triggername = ? and useraddress = ? limit 1", TriggerName, UserAddress)
	return
}

func (database *Database) Save(trigger *Trigger) (err error) {

	if trigger.TriggerName == "" {
		err = errors.New("TriggerName is empty")
		return
	}

	if trigger.UserAddress == "" {
		err = errors.New("UserAddress is empty")
		return
	}

	SearchTrigger, err := database.Find(trigger.TriggerName, trigger.UserAddress)
	if err != nil {
		return
	}

	if SearchTrigger.TriggerName != "" && SearchTrigger.TriggerName == trigger.TriggerName {
		sqlUpdate := "UPDATE triggers SET triggername = :TriggerName, useraddress = :UserAddress, contractaddress = :ContractAddress, "
		sqlUpdate += " method = :Method, status = :Status, number = :Number WHERE triggername = :TriggerName and useraddress = :UserAddress"
		_, err = database.Client.NamedExec(sqlUpdate, trigger)
	} else {
		sqlInsert := "INSERT INTO triggers (triggername, useraddress, contractaddress, method, status, number) "
		sqlInsert += "VALUES (:TriggerName, :UserAddress, :ContractAddress, :Method, :Status, :Number) "
		_, err = database.Client.NamedExec(sqlInsert, trigger)
	}

	return
}
