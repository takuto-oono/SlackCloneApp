package models

import (
	"fmt"

	"backend/config"
)

type User struct {
	ID       uint32 `json:"id"`
	Name     string `json:"name"`
	PassWord string `json:"password"`
}

func NewUser(id uint32, name, password string) *User {
	return &User{ID: id, Name: name, PassWord: password}
}

func (user *User) Create() error {
	cmd := fmt.Sprintf("INSERT INTO %s (id, name, password) VALUES (?, ?, ?)", config.Config.UserTableName)
	_, err := DbConnection.Exec(cmd, user.ID, user.Name, user.PassWord)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}

func GetUserByName(name string) (User, error) {
	cmd := fmt.Sprintf("SELECT id, name, password FROM %s WHERE name = ?", config.Config.UserTableName )
	row := DbConnection.QueryRow(cmd, name)
	var user User
	err := row.Scan(&user.ID, &user.Name, &user.PassWord)
	if err != nil {
		return User{}, err
	}
	if user.Name == "" || user.PassWord == "" {
		err = fmt.Errorf("not found name = %s", name)
		return User{}, err
	}
	return user, nil
}
