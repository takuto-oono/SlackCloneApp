package models

import (
	"fmt"

	"backend/config"
)

type Channel struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsPrivate   bool   `json:"is_private"`
	IsArchive   bool   `json:"is_archive"`
	WorkspaceId int    `json:"workspace_id"`
}

func NewChannel(id int, name, description string, isPrivate, isArchive bool, workspaceId int) *Channel {
	return &Channel{
		ID:          id,
		Name:        name,
		Description: description,
		IsPrivate:   isPrivate,
		IsArchive:   isArchive,
		WorkspaceId: workspaceId,
	}
}

func (c *Channel) SetId() error {
	cmd := fmt.Sprintf("SELECT id FROM %s", config.Config.ChannelsTableName)
	rows, err := DbConnection.Query(cmd)
	if err != nil {
		return err
	}
	defer rows.Close()
	maxId := 0
	for rows.Next() {
		var id int
		rows.Scan(&id)
		if id > maxId {
			maxId = id
		}
	}
	c.ID = maxId + 1
	return nil
}

func (c *Channel) Create() error {
	if err := c.SetId(); err != nil {
		return err
	}
	cmd := fmt.Sprintf("INSERT INTO %s (id, name, description, is_private, is_archive, workspace_id) VALUES (?, ?, ?, ?, ?, ?)", config.Config.ChannelsTableName)
	_, err := DbConnection.Exec(cmd, c.ID, c.Name, c.Description, c.IsPrivate, c.IsArchive, c.WorkspaceId)
	return err
}

func (c *Channel) IsExistSameNameChannelInWorkspace(workspaceId int) (bool, error) {
	cmd := fmt.Sprintf("SELECT name FROM %s WHERE workspace_id = ?", config.Config.ChannelsTableName)
	rows, err := DbConnection.Query(cmd, c.WorkspaceId)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		var channelName string
		if err := rows.Scan(&channelName); err != nil {
			return false, err
		}
		if channelName == c.Name {
			return true, nil
		}
	}
	return false, nil
}

func GetChannelById(channelId int) (Channel, error) {
	cmd := fmt.Sprintf("SELECT id, name, description, is_private, is_archive, workspace_id FROM %s WHERE id = ?", config.Config.ChannelsTableName)
	row := DbConnection.QueryRow(cmd, channelId)
	var c Channel
	err := row.Scan(&c.ID, &c.Name, &c.Description, &c.IsPrivate, &c.IsArchive, &c.WorkspaceId)
	return c, err
}

func (c *Channel) Delete() error {
	cmd := fmt.Sprintf("DELETE FROM %s WHERE id = ? AND name = ? AND description = ? AND is_private = ? AND is_archive = ? AND workspace_id = ?", config.Config.ChannelsTableName)
	_, err := DbConnection.Exec(cmd, c.ID, c.Name, c.Description, c.IsPrivate, c.IsArchive, c.WorkspaceId)
	return err
}

func IsExistChannelByChannelIdAndWorkspaceId(channelId, workspaceId int) (bool, error) {
	cmd := fmt.Sprintf("SELECT * FROM %s WHERE id = ? AND workspace_id = ?", config.Config.ChannelsTableName)
	rows, err := DbConnection.Query(cmd, channelId, workspaceId)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	cnt := 0
	for rows.Next() {
		cnt++
	}
	return cnt == 1, nil
}

func (c *Channel) GetChannelByIdAndWorkspaceId() error {
	cmd := fmt.Sprintf("SELECT id, name, description, is_private, is_archive FROM %s WHERE id = ? AND workspace_id = ?", config.Config.ChannelsTableName)
	row := DbConnection.QueryRow(cmd, c.ID, c.WorkspaceId)
	err := row.Scan(&c.ID, &c.Name, &c.Description, &c.IsPrivate, &c.IsArchive)
	return err
}

func GetChannelsByWorkspaceId(workspaceId int) ([]Channel, error) {
	channels := make([]Channel, 0)
	cmd := fmt.Sprintf("SELECT id, name, description, is_private, is_archive, workspace_id FROM %s WHERE workspace_id = ?", config.Config.ChannelsTableName)
	rows, err := DbConnection.Query(cmd, workspaceId)
	if err != nil {
		return channels, err
	}
	defer rows.Close()
	for rows.Next() {
		var ch Channel
		if err := rows.Scan(&ch.ID, &ch.Name, &ch.Description, &ch.IsPrivate, &ch.IsArchive, &ch.WorkspaceId); err != nil {
			return channels, err
		}
		channels = append(channels, ch)
	}
	return channels, nil
}
