package models

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"backend/config"
)

var DbConnection *sql.DB

func init() {
	driver := config.Config.Driver
	dbName := config.Config.DbName
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
			id INT PRIMARY KEY NOT NULL,
			name STRING NOT NULL UNIQUE,
			password STRING NOT NULL
		)
	`, config.Config.UserTableName)
	DbConnection.Exec(cmd)
}
