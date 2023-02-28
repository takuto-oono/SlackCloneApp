package models

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"

	"backend/config"
)

var DbConnection *sql.DB

func init() {
	// driver := config.Config.Driver
	driver := "postgres"
	// dbName := config.Config.DbName
	var err error
	DbConnection, err = sql.Open(driver, "host=db  user=postgres password=postgres dbname=postgres sslmode=disable")
	if err != nil {
		fmt.Println(err)
	}

	// create users table
	cmd := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id bigint PRIMARY KEY NOT NULL,
			name varchar(100) NOT NULL,
			password varchar(100) NOT NULL
		)
	`, config.Config.UserTableName)
	_, err = DbConnection.Exec(cmd)
	fmt.Println(err)

	// create workspace table
	cmd = fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id bigint PRIMARY KEY NOT NULL,
			name varchar(100) NOT NULL UNIQUE,
			workspace_primary_owner_id bigint not NULL
		)
	`, config.Config.WorkspaceTableName)
	DbConnection.Exec(cmd)

	// create workspace and user table
	cmd = fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			workspace_id bigint NOT NULL,
			user_id bigint NOT NULL,
			role_id bigint NOT NULL,
			PRIMARY KEY (workspace_id, user_id)
		)
	`, config.Config.WorkspaceAndUserTableName)
	DbConnection.Exec(cmd)

	// create role table
	cmd = fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id bigint PRIMARY KEY NOT NULL,
			name varchar(100) NOT NULL
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
			id bigint PRIMARY KEY NOT NULL,
			name varchar(100) NOT NULL,
			description varchar(100),
			is_private boolean NOT NULL,
			is_archive boolean NOT NULL,
			workspace_id bigint NOT NULL
		)
	`, config.Config.ChannelsTableName)
	DbConnection.Exec(cmd)

	// create channels_and_users table
	cmd = fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s
		(
			channel_id bigint NOT NULL,
			user_id bigint NOT NULL,
			is_admin boolean NOT NULL,
			PRIMARY KEY (channel_id, user_id)
		)
	`, config.Config.ChannelsAndUserTableName)
	DbConnection.Exec(cmd)

	// create messages table
	cmd = fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s 
		(
			id bigint PRIMARY KEY NOT NULL,
			text varchar(100) NOT NULL,
			date varchar(100) NOT NULL,
			channel_id bigint NOT NULL,
			user_id bigint NOT NULL
		)
	`, config.Config.MessagesTableName)
	DbConnection.Exec(cmd)
}
