package models

import (
	"fmt"

	"backend/config"
)

type ChannelsAndWorkspaces struct {
	ChannelId   int
	WorkspaceId int
}

func NewChannelsAndWorkspaces(channelId, workspaceId int) *ChannelsAndWorkspaces {
	return &ChannelsAndWorkspaces{
		ChannelId:   channelId,
		WorkspaceId: workspaceId,
	}
}

func (caw *ChannelsAndWorkspaces) Create() error {
	cmd := fmt.Sprintf("INSERT INTO %s (channel_id, workspace_id) VALUES (?, ?)", config.Config.ChannelsAndWorkspaceTableName)
	_, err := DbConnection.Exec(cmd, caw.ChannelId, caw.WorkspaceId)
	return err
}

func FindChannelIdsByWorkspaceId(workspaceId int) ([]int, error) {
	cmd := fmt.Sprintf("SELECT channel_id FROM %s WHERE workspace_id = ?", config.Config.ChannelsAndWorkspaceTableName)
	rows, err := DbConnection.Query(cmd, workspaceId)
	if err != nil {
		return []int{}, err
	}
	defer rows.Close()
	res := make([]int, 0)
	for rows.Next() {
		var channelId int
		err := rows.Scan(&channelId)
		if err != nil {
			return []int{}, err
		}
		res = append(res, channelId)
	}
	return res, nil
}

func IsExistCAWByChannelIdAndWorkspaceId(channelId, workspaceId int) bool {
	cmd := fmt.Sprintf("SELECT * FROM %s WHERE channel_id = ? AND workspace_id = ?", config.Config.ChannelsAndWorkspaceTableName)
	rows, err := DbConnection.Query(cmd, channelId, workspaceId)
	if err != nil {
		return false
	}
	defer rows.Close()
	cnt := 0
	for rows.Next() {
		cnt++
	}
	return cnt == 1
}

func (caw *ChannelsAndWorkspaces) Delete() error {
	cmd := fmt.Sprintf("DELETE FROM %s WHERE channel_id = ? AND workspace_id = ?", config.Config.ChannelsAndWorkspaceTableName)
	_, err := DbConnection.Exec(cmd, caw.ChannelId, caw.WorkspaceId)
	return err
}
