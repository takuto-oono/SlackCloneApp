package models

import (
	"fmt"

	"backend/config"
)

type Workspace struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	PrimaryOwnerId uint32 `json:"primary_owner_id"`
}

func NewWorkspace(id int, name string, primaryOwnerId uint32) *Workspace {
	return &Workspace{
		ID:             id,
		Name:           name,
		PrimaryOwnerId: primaryOwnerId,
	}
}

func (w *Workspace) CreateWorkspace() error {
	cmd := fmt.Sprintf("SELECT * FROM %s", config.Config.WorkspaceTableName)
	rows, err := DbConnection.Query(cmd)
	if err != nil {
		return err
	}
	defer rows.Close()
	cnt := 0
	for rows.Next() {
		cnt ++
	}
	w.ID = cnt + 1

	cmd = fmt.Sprintf("INSERT INTO %s (id, name, workspace_primary_owner_id) VALUES (?, ?, ?)", config.Config.WorkspaceTableName)
	_, err = DbConnection.Exec(cmd, w.ID, w.Name, w.PrimaryOwnerId)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}

func GetWorkspaceByName(name string) (Workspace, error) {
	if name == "" {
		return Workspace{}, fmt.Errorf("not found workspace name")
	}
	cmd := fmt.Sprintf("SELECT id, name, workspace_primary_owner_id FROM %s WHERE name = ?", config.Config.WorkspaceTableName)
	row := DbConnection.QueryRow(cmd, name)
	var w Workspace
	err := row.Scan(&w.ID, &w.Name, &w.PrimaryOwnerId)
	if err != nil {
		return Workspace{}, err
	}

	if w.ID == 0 || w.Name == "" || w.PrimaryOwnerId == 0 {
		return Workspace{}, fmt.Errorf("find empty fields")
	}

	if w.Name != name {
		return Workspace{}, fmt.Errorf("find wrong workspace")
	}

	return w, nil
}
