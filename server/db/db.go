package db

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

type Storable interface {
	StoreSelf(db *sql.DB) error
	SelfScan(row *sql.Rows) error
}

type StorableList interface {
	At(i int) Storable
	Len() int
	AppendEmpty()
}

type Database struct {
	pool *sql.DB
}

func NewDatabase(path string) (*Database, error) {
	var Pool *sql.DB
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		_, err = os.Create(path)
		if err != nil {
			log.Fatalf(err.Error())
		}
		Pool, err = sql.Open("sqlite3", path)
		if err != nil {
			log.Fatalf(err.Error())
		}
		log.Println("Database created!")
	} else {
		Pool, err = sql.Open("sqlite3", path)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}
	InitTable(Pool)
	return &Database{pool: Pool}, nil
}

func InitTable(db *sql.DB) {
	stmt := `CREATE TABLE Positions(Id INTEGER PRIMARY KEY, XCoord INTEGER, YCoord INTEGER)`
	_, _ = db.Exec(stmt)
}

func (db *Database) GetAll(dest StorableList) error {
	stmt := `SELECT Id, XCoord, YCoord FROM Positions`
	rows, err := db.pool.Query(stmt)
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}(rows)
	if err != nil {
		return errors.New(fmt.Sprintf("cant retrieve data %v", err.Error()))
	}
	i := 0
	for rows.Next() {
		if i >= dest.Len() {
			dest.AppendEmpty()
		}
		err = dest.At(i).SelfScan(rows)
		i++
		if err != nil {
			return errors.New(fmt.Sprintf("cant retrieve data %v", err.Error()))
		}
	}
	return nil
}

func (db *Database) Store(data Storable) error {
	err := data.StoreSelf(db.pool)
	if err != nil {
		return errors.New(fmt.Sprintf("cant store value %v", err.Error()))
	}
	return nil
}

func (db *Database) GetStorable(id int64, dest Storable) error {
	stmt := `SELECT Id, XCoord, YCoord FROM Positions WHERE Id=?`
	res, err := db.pool.Query(stmt, id)
	defer func(res *sql.Rows) {
		err = res.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}(res)
	if err != nil {
		return errors.New(fmt.Sprintf("cant retrieve data %v", err.Error()))
	}
	for res.Next() {
		err = dest.SelfScan(res)
		if err != nil {
			return errors.New(fmt.Sprintf("cant retrieve data %v", err.Error()))
		}
	}
	return nil
}
