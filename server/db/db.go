package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

const (
	dbfile = "productstore.db"
)

var DB *sql.DB

func ConnectToSqlite() *sql.DB {
	var err error
	currentWorkingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	filepath := filepath.Join(currentWorkingDirectory, dbfile)
	DB, err = sql.Open("sqlite3", filepath)
	if err != nil {
		log.Fatal(err)
	}
	return DB
}
