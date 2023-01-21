package models

import (
	"fmt"

	"backend/config"
)

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	PassWord string `json:"password"`
}

func NewUser(id, name, password string) *User {
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

func GetAllUsers() ([]User, error) {
	cmd := fmt.Sprintf("SELECT id, name, password FROM %s", config.Config.UserTableName)
	rows, err := DbConnection.Query(cmd)
	if err != nil {
		return []User{}, err
	}
	defer rows.Close()

	var users = []User{}
	for rows.Next() {
		var user User
		rows.Scan(&user.ID, &user.Name, &user.PassWord)
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return []User{}, err
	}
	return users, nil
}

func GetUserById(uuid string) (User, error) {
	cmd := fmt.Sprintf("SELECT id, name, password FROM %s WHERE id = ?", config.Config.UserTableName)
	row := DbConnection.QueryRow(cmd, uuid)
	var user User
	err := row.Scan(&user.ID, &user.Name, &user.PassWord)
	if err != nil {
		return User{}, err
	}
	return user, err
}

func (user *User) UpdateUser() error {
	cmd := fmt.Sprintf("UPDATE %s SET name = ?, password = ? WHERE id = ?", config.Config.UserTableName)
	_, err := DbConnection.Exec(cmd, user.Name, user.PassWord, user.ID)
	return err
}

func (user *User) DeleteUser() error {
	cmd := fmt.Sprintf("DELETE FROM %s WHERE id = ?", config.Config.UserTableName)
	_, err := DbConnection.Exec(cmd, user.ID, user.Name, user.PassWord)
	return err
}
