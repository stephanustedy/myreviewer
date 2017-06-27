package database

import (
	"log"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func Initialize() {
	db, err := sql.Open("sqlite3", "files/database/myreviewer.db")
	if err != nil {
		log.Fatal(err)
	}

	result, err := db.Exec(
 "create table IF NOT EXISTS employee (employeeID integer PRIMARY KEY,name varchar(255) NOT null,age int, person_id int, FOREIGN KEY (person_id) REFERENCES persons(id), CONSTRAINT uc_empID UNIQUE (employeeID, person_id, name))",)
 
if err != nil {
 log.Fatal(err)
}
log.Println(result)
}
