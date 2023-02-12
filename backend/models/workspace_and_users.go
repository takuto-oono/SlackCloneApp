package models

import (
	"fmt"

	"backend/config"
)

type WorkspaceAndUsers struct {
	WorkspaceId int    `json:"workspace_id"`
	UserId      uint32 `json:"user_id"`
	RoleId      int    `json:"role_id"`
}

func NewWorkspaceAndUsers(workspaceId int, userId uint32, roleId int) *WorkspaceAndUsers {
	return &WorkspaceAndUsers{
		WorkspaceId: workspaceId,
		UserId:      userId,
		RoleId:      roleId,
	}
}

func (wau *WorkspaceAndUsers) Create() error {
	cmd := fmt.Sprintf("INSERT INTO %s (workspace_id, user_id, role_id) VALUES (?, ?, ?)", config.Config.WorkspaceAndUserTableName)
	_, err := DbConnection.Exec(cmd, wau.WorkspaceId, wau.UserId, wau.RoleId)
	return err
}

func GetWorkspaceAndUserByWorkspaceIdAndUserId(workspaceId int, userId uint32) (WorkspaceAndUsers, error) {
	cmd := fmt.Sprintf("SELECT workspace_id, user_id, role_id FROM %s WHERE workspace_id = ? AND user_id = ?", config.Config.WorkspaceAndUserTableName)
	row := DbConnection.QueryRow(cmd, workspaceId, userId)
	var wau WorkspaceAndUsers
	err := row.Scan(&wau.WorkspaceId, &wau.UserId, &wau.RoleId)
	if err != nil {
		return WorkspaceAndUsers{}, err
	}
	return wau, err
}

func (wau *WorkspaceAndUsers) IsExistWorkspaceAndUser() bool {
	var cmd string
	if wau.RoleId == 0 {
		cmd = fmt.Sprintf("SELECT workspace_id, user_id, role_id FROM %s WHERE workspace_id = ? AND user_id = ?", config.Config.WorkspaceAndUserTableName)
	} else {
		cmd = fmt.Sprintf("SELECT workspace_id, user_id, role_id FROM %s WHERE workspace_id = ? AND user_id = ? AND role_id = ?", config.Config.WorkspaceAndUserTableName)

	}
	rows, err := DbConnection.Query(cmd, wau.WorkspaceId, wau.UserId, wau.RoleId)
	if err != nil {
		return false
	}
	defer rows.Close()
	cnt := 0
	for rows.Next() {
		cnt += 1
	}
	return cnt == 1
}

func (wau *WorkspaceAndUsers) DeleteWorkspaceAndUser() error {
	cmd := fmt.Sprintf("DELETE FROM %s WHERE workspace_id = ? AND user_id = ? AND role_id = ?", config.Config.WorkspaceAndUserTableName)
	_, err := DbConnection.Exec(cmd, wau.WorkspaceId, wau.UserId, wau.RoleId)
	return err
}
