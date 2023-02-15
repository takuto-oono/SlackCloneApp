package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateChannelsAndWorkspaces(t *testing.T) {
	channelId := 36636564533
	workspaceId := 467068305803
	caw := NewChannelsAndWorkspaces(channelId, workspaceId)
	assert.Empty(t, caw.Create())

	caw.ChannelId = 435534623
	assert.Empty(t, caw.Create())

	caw.ChannelId = channelId
	caw.WorkspaceId = 54646435
	assert.NotEmpty(t, caw.Create())
}

func TestIsExistCAWByChannelIdAndWorkspaceId(t *testing.T) {
	channelId := 42415353
	workspaceId := 533246
	caw := NewChannelsAndWorkspaces(channelId, workspaceId)
	assert.Empty(t, caw.Create())
	assert.Equal(t, true, IsExistCAWByChannelIdAndWorkspaceId(channelId, workspaceId))
	assert.Equal(t, false, IsExistCAWByChannelIdAndWorkspaceId(-1, workspaceId))
	assert.Equal(t, false, IsExistCAWByChannelIdAndWorkspaceId(channelId, -1))
}

func TestDeleteCAW(t *testing.T) {
	channelId := 648947933
	workspaceId := 395798579239
	caw := NewChannelsAndWorkspaces(channelId, workspaceId)
	assert.Empty(t, caw.Create())
	assert.Empty(t, caw.Delete())
	assert.Equal(t, false, IsExistCAWByChannelIdAndWorkspaceId(channelId, workspaceId))
}
