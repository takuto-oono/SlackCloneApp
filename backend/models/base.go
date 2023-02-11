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

	// create users table
	cmd := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id INT PRIMARY KEY NOT NULL,
			name STRING NOT NULL,
			password STRING NOT NULL
		)
	`, config.Config.UserTableName)
	DbConnection.Exec(cmd)

	// create workspace table
	cmd = fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id INT PRIMARY KEY NOT NULL,
			name STRING NOT NULL UNIQUE,
			workspace_primary_owner_id STRING not NULL
		)
	`, config.Config.WorkspaceTableName)
	DbConnection.Exec(cmd)

	// create workspace and user table
	cmd = fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			workspace_id INT NOT NULL,
			user_id INT NOT NULL,
			role_id INT NOT NULL,
			PRIMARY KEY (workspace_id, user_id)
		)
	`, config.Config.WorkspaceAndUserTableName)
	DbConnection.Exec(cmd)

	// create role table
	cmd = fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id INT PRIMARY KEY NOT NULL,
			name STRING NOT NULL
		)
	`, config.Config.RoleTableName)
	DbConnection.Exec(cmd)

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
			is_archive BOOLEAN NOT NULL
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
	DbConnection.Exec(cmd)

	// create channels_and_workspaces table
	cmd = fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s
		(
			channel_id INT PRIMARY KEY NOT NULL,
			workspace_id INT NOT NULL
		)
	`, config.Config.ChannelsAndWorkspaceTableName)
	DbConnection.Exec(cmd)

}
