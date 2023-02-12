package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateChannelAndUsers(t *testing.T) {
	channelId := 5553155345262
	userId := uint32(452526)
	isAdmin := true

	cau := NewChannelsAndUses(channelId, userId, isAdmin)
	assert.Empty(t, cau.CreateChannelAndUsers())
	
	cau.IsAdmin = false
	assert.NotEmpty(t, cau.CreateChannelAndUsers())

	cau.ChannelId = 5347595792
	assert.Empty(t, cau.CreateChannelAndUsers())

	cau.ChannelId = channelId
	cau.UserId = 53534626
	assert.Empty(t, cau.CreateChannelAndUsers())
}