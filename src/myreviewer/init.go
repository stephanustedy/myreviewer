package myreviewer

import (
	"log"
	"os"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func Initialize() {
	var err error
	dbPath := "files/database/myreviewer.db"

	// detect if file exists
	_, err = os.Stat(dbPath)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(dbPath)
		if err != nil {
			log.Println(err)
		}
		defer file.Close()
	}

	db, err = sql.Open("sqlite3", dbPath)
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
 			, username varchar(255) NOT null
 			, role int not null
 			, status int not null
 			, team_id int not null)`)
 
	if err != nil {
		log.Fatal(err)
	}

	//init table review
	_, err = db.Exec(`
 		create table IF NOT EXISTS 
 		review (
 			review_id integer PRIMARY KEY
 			, team_id int NOT null
 			, qa int NOT null
 			, developer int NOT null
 			, pull_request varchar(255) not null
 			, status int not null)`)
 
	if err != nil {
		log.Fatal(err)
	}

	//init table reviewer
	_, err = db.Exec(`
 		create table IF NOT EXISTS 
 		reviewer (
 			id integer PRIMARY key
 			, review_id integer not null 
 			, reviewer_id int NOT null
 			, status int not null
 			, last_notify timestamp
 			, notify_count int not null)`)
 
	if err != nil {
		log.Fatal(err)
	}

	InitializeCron()
}

