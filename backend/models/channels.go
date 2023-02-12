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
}

func NewChannel(id int, name, description string, isPrivate, isArchive bool) *Channel {
	return &Channel{
		ID:          id,
		Name:        name,
		Description: description,
		IsPrivate:   isPrivate,
		IsArchive:   isArchive,
	}
}

func (c *Channel) CreateChannel() error {
	cmd := fmt.Sprintf("SELECT * FROM %s", config.Config.ChannelsTableName)
	rows, err := DbConnection.Query(cmd)
	if err != nil {
		return err
	}
	defer rows.Close()
	cnt := 0
	for rows.Next() {
		cnt++
	}
	c.ID = cnt + 1

	cmd = fmt.Sprintf("INSERT INTO %s (id, name, description, is_private, is_archive) VALUES (?, ?, ?, ?, ?)", config.Config.ChannelsTableName)
	_, err = DbConnection.Exec(cmd, c.ID, c.Name, c.Description, c.IsPrivate, c.IsArchive)
	return err
}

func (c *Channel) IsExistSameNameChannelInWorkspace(workspaceId int) (bool, error) {
	channelIds, err := FindChannelIdsByWorkspaceId(workspaceId)
	if err != nil {
		return false, err
	}
	
	cmd := fmt.Sprintf("SELECT name FROM %s WHERE id = ?", config.Config.ChannelsTableName)
	for _, channelId := range channelIds {
		row := DbConnection.QueryRow(cmd, channelId)
		var channelName string
		err := row.Scan(&channelName)
		if err != nil {
			return false, err
		}
		if channelName == c.Name {
			return true, nil
		}
	}
	return false, nil
}
