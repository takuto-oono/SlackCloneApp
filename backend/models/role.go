package models

import (
	"fmt"

	"backend/config"
)

type Role struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func NewRole(id int, name string) *Role {
	return &Role{ID: id, Name: name}
}

func (r *Role) Create() error {
	cmd := fmt.Sprintf("INSERT INTO %s (id, name) VALUES (?, ?)", config.Config.RoleTableName)
	_, err := DbConnection.Exec(cmd, r.ID, r.Name)
	if err != nil {
		fmt.Println(err.Error())
	}
	return err
}

func GetRoleById(id int) (Role, error) {
	cmd := fmt.Sprintf("SELECT id, name FROM %s WHERE id = ?", config.Config.RoleTableName)
	row := DbConnection.QueryRow(cmd, id)
	var r Role
	err := row.Scan(&r.ID, &r.Name)
	return r, err
}
