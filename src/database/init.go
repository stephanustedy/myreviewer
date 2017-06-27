package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func Initialize() {
	db, err := sql.Open("sqlite3", "files/database/myreviewer.db")
	if err != nil {
		log.Fatal(err)
	}
}
