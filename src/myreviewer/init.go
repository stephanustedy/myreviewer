package myreviewer

import (
	"log"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func Initialize() {
	var err error
	db, err = sql.Open("sqlite3", "files/database/myreviewer.db")
	if err != nil {
		log.Fatal(err)
	}

	//init table team
	_, err = db.Exec(`
 		create table IF NOT EXISTS 
 		team (
 			team_id integer PRIMARY KEY
 			, name varchar(255) NOT null
 			, webhook varchar(500) not null
 			, status int not null
 			, channel varchar(255) not null)`)
 
	if err != nil {
		log.Fatal(err)
	}

	//init table member
	_, err = db.Exec(`
 		create table IF NOT EXISTS 
 		member (
 			member_id integer PRIMARY KEY
 			, name varchar(255) NOT null
 			, status int not null)`)
 
	if err != nil {
		log.Fatal(err)
	}
}

