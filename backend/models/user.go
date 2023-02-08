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
	cmd := fmt.Sprintf("SELECT id, name, password FROM %s WHERE name = ?", config.Config.UserTableName)
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

func GetUserById(id uint32) (User, error) {
	cmd := fmt.Sprintf("SELECT id, name, password FROM %s WHERE id = ?", config.Config.UserTableName)
	row := DbConnection.QueryRow(cmd, id)
	var user User
	err := row.Scan(&user.ID, &user.Name, &user.PassWord)
	if err != nil {
		return User{}, err
	}
	if user.Name == "" || user.PassWord == "" {
		err = fmt.Errorf("not found id")
		return User{}, err
	}
	return user, nil
}

func (u *User) IsExistUserSameUsernameAndPassword() (bool, error) {
	cmd := fmt.Sprintf("SELECT id, name, password FROM %s WHERE name = ? AND password = ?", config.Config.UserTableName)
	rows, err := DbConnection.Query(cmd, u.Name, u.PassWord)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	cnt := 0
	for rows.Next() {
		cnt += 1
	}
	return cnt != 0, nil
}