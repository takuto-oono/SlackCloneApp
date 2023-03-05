package models

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"backend/config"
)

var DbConnection *sql.DB
var db *gorm.DB

func init() {
	driver := config.Config.Driver
	dbName := config.Config.DbName
	var err error
	DbConnection, err = sql.Open(driver, dbName)
	if err != nil {
		fmt.Println(err)
	}

	db, err = gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("------------------")
	fmt.Println(db)
	fmt.Println("------------------")

	// create users table
	cmd := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id INT PRIMARY KEY NOT NULL,
			name STRING NOT NULL,
			password STRING NOT NULL
		)
	`, config.Config.UserTableName)
	_, err = DbConnection.Exec(cmd)
	fmt.Println(err)

	// create workspace table
	cmd = fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id INT PRIMARY KEY NOT NULL,
			name STRING NOT NULL UNIQUE,
			workspace_primary_owner_id STRING not NULL
		)
	`, config.Config.WorkspaceTableName)
	_, err = DbConnection.Exec(cmd)
	fmt.Println(err)

	// create workspace and user table
	cmd = fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			workspace_id INT NOT NULL,
			user_id INT NOT NULL,
			role_id INT NOT NULL,
			PRIMARY KEY (workspace_id, user_id)
		)
	`, config.Config.WorkspaceAndUserTableName)
	_, err = DbConnection.Exec(cmd)
	fmt.Println(err)

	// create role table
	cmd = fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id INT PRIMARY KEY NOT NULL,
			name STRING NOT NULL
		)
	`, config.Config.RoleTableName)
	_, err = DbConnection.Exec(cmd)
	fmt.Println(err)

	// insert 4 roles in roles table
	roleNames := []string{
		"Workspace Primary Owner",
		"Workspace Owners",
		"Workspace Admins",
		"Full members",
	}
	for i, n := range roleNames {
		r := NewRole(i+1, n)
		err := r.Create()
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	// create channels table
	cmd = fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s 
		(
			id INT PRIMARY KEY NOT NULL,
			name STRING NOT NULL,
			description STRING,
			is_private BOOLEAN NOT NULL,
			is_archive BOOLEAN NOT NULL,
			workspace_id INT NOT NULL
		)
	`, config.Config.ChannelsTableName)
	_, err = DbConnection.Exec(cmd)
	fmt.Println(err)

	// create channels_and_users table
	cmd = fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s
		(
			channel_id INT NOT NULL,
			user_id INT NOT NULL,
			is_admin BOOLEAN NOT NULL,
			PRIMARY KEY (channel_id, user_id)
		)
	`, config.Config.ChannelsAndUserTableName)
	_, err = DbConnection.Exec(cmd)
	fmt.Println(err)

	// create messages table
	cmd = fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s 
		(
			id INT PRIMARY KEY NOT NULL,
			text STRING NOT NULL,
			date STRING NOT NULL,
			channel_id INT NOT NULL,
			user_id INT NOT NULL
		)
	`, config.Config.MessagesTableName)
	_, err = DbConnection.Exec(cmd)
	fmt.Println(err)
	
	// create direct_messages table
	// cmd = fmt.Sprintf(`
	// 	CREATE TABLE IF NOT EXISTS %s
	// 	(
	// 		id INT PRIMARY KEY NOT NULL,
	// 		text STRING NOT NULL,
	// 		date STRING NOT NULL,
	// 		send_user_id INT NOT NULL,
	// 		dm_line INT NOT NULL
	// 	)
	// `, config.Config.DirectMessagesTableName)
	// db.Exec(cmd)
	db.AutoMigrate(&DirectMessage{})
	
	// create dm_lines table
	// cmd = fmt.Sprintf(`
	// 	CREATE TABLE IF NOT EXISTS %s
	// 	(
	// 		id INT PRIMARY KEY NOT NULL,
	// 		workspace_id INT NOT NULL,
	// 		user_id_1 INT NOT NULL,
	// 		user_id_2 INT NOT NULL
	// 	)
	// `, config.Config.DMLinesTableName)
	// db.Exec(cmd)
	db.AutoMigrate(&DMLine{})
}
