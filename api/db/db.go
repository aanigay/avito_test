package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
)

type Database struct {
	db *sql.DB
}

func NewDatabase() (*Database, error) {
	db, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}

func (d *Database) Close() {
	err := d.db.Close()
	if err != nil {
		log.Fatal("closing db connection error")
		return
	}
}

func (d *Database) GetDB() *sql.DB {
	return d.db
}
