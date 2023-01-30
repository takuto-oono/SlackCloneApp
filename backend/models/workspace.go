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
	cmd := fmt.Sprintf("SELECT COUNT(*) FROM %s", config.Config.WorkspaceTableName)
	cntColumns, err := DbConnection.Exec(cmd)
	fmt.Println(cntColumns)
	w.ID = cntColumns + 1

	cmd = fmt.Sprintf("INSERT INTO %s (id, name, workspace_primary_owner_id) VALUES (?, ?, ?)", config.Config.WorkspaceTableName)
	_, err := DbConnection.Exec(cmd, w.ID, w.Name, w.PrimaryOwnerId)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}
