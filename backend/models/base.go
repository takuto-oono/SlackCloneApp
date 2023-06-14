package models

import (
	"database/sql"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DbConnection *sql.DB
var db *gorm.DB

func ExportDB() *gorm.DB {
	return db
}

func initRoleTable(db *gorm.DB) error {
	roleNames := []string{
		"Workspace Primary Owner",
		"Workspace Owners",
		"Workspace Admins",
		"Full members",
	}

	tx := db.Begin()
	if err := tx.Error; err != nil {
		return err
	}

	for i, n := range roleNames {
		r := NewRole(i+1, n)
		err := r.Create(tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func init() {
	var err error
	db, err = gorm.Open(mysql.New(mysql.Config{
		DSN: "MYSQL:MYSQL@tcp(MYSQL)/MYSQL?charset=utf8mb4&parseTime=True&loc=Local",
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	fmt.Println("------------------")
	fmt.Println(db)
	fmt.Println("------------------")

	if err := db.AutoMigrate(&User{}); err != nil {
		fmt.Println(err)
	}

	if err := db.AutoMigrate(&Role{}); err != nil {
		fmt.Println(err)
	}

	if err := initRoleTable(db); err != nil {
		fmt.Println(err)
	}

	if err := db.AutoMigrate(&Workspace{}); err != nil {
		fmt.Println(err)
	}

	if err := db.AutoMigrate(&WorkspaceAndUsers{}); err != nil {
		fmt.Println(err)
	}

	if err := db.AutoMigrate(&Channel{}); err != nil {
		fmt.Println(err)
	}

	if err := db.AutoMigrate(&ChannelsAndUsers{}); err != nil {
		fmt.Println(err)
	}

	if err := db.AutoMigrate(&DMLine{}); err != nil {
		fmt.Println(err)
	}

	if err := db.AutoMigrate(&Message{}); err != nil {
		fmt.Println(err)
	}

	if err := db.AutoMigrate(&Thread{}); err != nil {
		fmt.Println(err)
	}

	if err := db.AutoMigrate(&ThreadAndMessage{}); err != nil {
		fmt.Println(err)
	}

	if err := db.AutoMigrate(&ThreadAndUser{}); err != nil {
		fmt.Println(err)
	}

	if err := db.AutoMigrate(&Mention{}); err != nil {
		fmt.Println(err)
	}

	if err := db.AutoMigrate(&MessageAndUser{}); err != nil {
		fmt.Println(err)
	}
}
