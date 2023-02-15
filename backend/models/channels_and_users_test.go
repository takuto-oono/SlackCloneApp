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
	assert.Empty(t, cau.Create())

	cau.IsAdmin = false
	assert.NotEmpty(t, cau.Create())

	cau.ChannelId = 5347595792
	assert.Empty(t, cau.Create())

	cau.ChannelId = channelId
	cau.UserId = 53534626
	assert.Empty(t, cau.Create())
}

func TestIsExistCAUByChannelIdAndUserId(t *testing.T) {
	channelId := 532446463423234
	userId := uint32(3535422)

	cau := NewChannelsAndUses(channelId, userId, false)
	assert.Empty(t, cau.Create())

	assert.Equal(t, true, IsExistCAUByChannelIdAndUserId(channelId, userId))
	assert.Equal(t, false, IsExistCAUByChannelIdAndUserId(-1, userId))
}

func TestIsAdminUserInChannel(t *testing.T) {
	cau := NewChannelsAndUses(532446463423234434, uint32(53663732), true)
	assert.Empty(t, cau.Create())

	assert.Equal(t, true, IsAdminUserInChannel(cau.ChannelId, cau.UserId))
	assert.Equal(t, false, IsAdminUserInChannel(cau.ChannelId, 46433))

	cau = NewChannelsAndUses(cau.ChannelId, 5236632, false)
	assert.Empty(t, cau.Create())
	assert.Equal(t, false, IsAdminUserInChannel(cau.ChannelId, cau.UserId))
}

func TestDeleteUserFromChannel(t *testing.T) {
	cau := NewChannelsAndUses(9479923, 646433, true)
	assert.Empty(t, cau.Create())
	assert.Empty(t, cau.Delete())

	cau = NewChannelsAndUses(6464333, 79797, true)
	assert.Empty(t, cau.Create())
	cau.IsAdmin = false
	assert.Empty(t, cau.Delete())
	assert.Equal(t, false, IsExistCAUByChannelIdAndUserId(cau.ChannelId, cau.UserId))

	channelId := 37597692793
	userId := uint32(35362622)
	cau = NewChannelsAndUses(channelId, userId, true)
	assert.Empty(t, cau.Create())
	assert.Equal(t, true, IsExistCAUByChannelIdAndUserId(channelId, userId))

	cau.ChannelId = -1
	assert.Empty(t, cau.Delete())
	assert.Equal(t, true, IsExistCAUByChannelIdAndUserId(channelId, userId))

	cau.ChannelId = channelId
	assert.Empty(t, cau.Delete())
	assert.Equal(t, false, IsExistCAUByChannelIdAndUserId(channelId, userId))
}

func TestDeleteCAUByChannelId(t *testing.T) {
	channelId := 525325
	userIds := make([]uint32, 20)
	for i := 0; i < 20; i++ {
		userIds[i] = uint32(758932 + i*10)
	}

	for _, userId := range userIds {
		cau := NewChannelsAndUses(channelId, userId, false)
		assert.Empty(t, cau.Create())
	}

	assert.Empty(t, DeleteCAUByChannelId(channelId))

	for _, userId := range userIds {
		assert.Equal(t, false, IsExistCAUByChannelIdAndUserId(channelId, userId))
	}
}
