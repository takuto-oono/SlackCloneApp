package models

import (
	"fmt"

	"backend/config"
	"backend/utils"
)

type Message struct {
	ID        int    `json:"id"`
	Text      string `json:"text"`
	Date      string `json:"date"`
	ChannelId int    `json:"channel_id"`
	UserId    uint32 `json:"user_id"`
}

func NewMessage(text string, channelId int, userId uint32) *Message {
	return &Message{
		ID:        0,
		Text:      text,
		Date:      "",
		ChannelId: channelId,
		UserId:    userId,
	}
}

func (m *Message) SetID() error {
	cmd := fmt.Sprintf("SELECT id FROM %s", config.Config.MessagesTableName)
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
	m.ID = maxId + 1
	return nil
}

func (m *Message) SetDate() {
	m.Date = utils.GetCurrentTime()
}

func (m *Message) Create() error {
	if err := m.SetID(); err != nil {
		return err
	}
	m.SetDate()
	cmd := fmt.Sprintf("INSERT INTO %s (id, text, date, channel_id, user_id) VALUES ($1, $2, $3, $4, $5)", config.Config.MessagesTableName)
	_, err := DbConnection.Exec(cmd, m.ID, m.Text, m.Date, m.ChannelId, m.UserId)
	return err
}

func GetMessagesByChannelId(channelId int) ([]Message, error) {
	res := make([]Message, 0)
	cmd := fmt.Sprintf("SELECT id, text, date, channel_id, user_id FROM %s WHERE channel_id = $1 ORDER BY date DESC", config.Config.MessagesTableName)
	rows, err := DbConnection.Query(cmd, channelId)
	if err != nil {
		return res, err
	}
	defer rows.Close()
	for rows.Next() {
		var m Message
		err := rows.Scan(
			&m.ID,
			&m.Text,
			&m.Date,
			&m.ChannelId,
			&m.UserId,
		)
		if err != nil {
			return res, err
		}
		res = append(res, m)
	}
	return res, nil
}
