package models

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DbConnection *sql.DB

const (
	userTableName = "user_table"
)

func init() {
	driver := "sqlite3"
	dbName := "SlackCloneDB.sql"
	var err error
	DbConnection, err = sql.Open(driver, dbName)
	if err != nil {
		log.Fatalln(err)
	}
	cmd := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			num INT PRIMARY KEY NOT NULL,
			word STRING
		)
	`, "test_db1")
	DbConnection.Exec(cmd)
	
	cmd = fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id STRING PRIMARY KEY NOT NULL,
			name STRING,
			password STRING
		)
	`, userTableName)
	DbConnection.Exec(cmd)
}
