package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
)

// драйверы
import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/sijms/go-ora"
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

type DatabaseConfig struct {
	Driver   string
	DSN      string // для облачных бд
	FilePath string // для файловых бд
}

type Database struct {
	pool *sql.DB
}

func IsSupported(driver string) bool {
	return driver == "sqlite" || driver == "postgres" || driver == "mysql" || driver == "oracle"
}

func NewDatabase(conf *DatabaseConfig) (*Database, error) {
	var Pool *sql.DB
	if !IsSupported(conf.Driver) {
		log.Fatalf("this type of database is not supported " + conf.Driver)
	}
	if conf.Driver == "sqlite" {
		if conf.FilePath == "" {
			panic("empty filepath for file-based db!")
		}
		if _, err := os.Stat(conf.FilePath); errors.Is(err, os.ErrNotExist) {
			_, err = os.Create(conf.FilePath)
			if err != nil {
				log.Fatalf(err.Error())
			}
			Pool, err = sql.Open("sqlite3", conf.FilePath)
			if err != nil {
				log.Fatalf(err.Error())
			}
			log.Println("Database created!")
		} else {
			Pool, err = sql.Open("sqlite3", conf.FilePath)
			if err != nil {
				log.Fatalf(err.Error())
			}
		}
	} else {
		pool, err := sql.Open(conf.Driver, conf.DSN)
		if err != nil {
			log.Fatalf(err.Error())
		}
		err = pool.Ping()
		if err != nil {
			log.Fatalf(err.Error())
		} else {
			log.Printf("Connected successfully to %s", conf.DSN)
		}
		Pool = pool
	}
	InitTable(Pool)
	return &Database{pool: Pool}, nil
}

func InitTable(db *sql.DB) {
	stmt := `CREATE TABLE Positions(Id INTEGER PRIMARY KEY, XCoord INTEGER, YCoord INTEGER)`
	_, err := db.Exec(stmt)
	if err != nil {
		log.Println(err.Error())
	}
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
