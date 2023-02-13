package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateChannelsAndWorkspaces(t *testing.T) {
	channelId := 36636564533
	workspaceId := 467068305803
	caw := NewChannelsAndWorkspaces(channelId, workspaceId)
	assert.Empty(t, caw.CreateChannelsAndWorkspaces())

	caw.ChannelId = 435534623
	assert.Empty(t, caw.CreateChannelsAndWorkspaces())

	caw.ChannelId = channelId
	caw.WorkspaceId = 54646435
	assert.NotEmpty(t, caw.CreateChannelsAndWorkspaces())
}

func TestIsExistCAWByChannelIdAndWorkspaceId(t *testing.T) {
	channelId := 42415353
	workspaceId := 533246
	caw := NewChannelsAndWorkspaces(channelId, workspaceId)
	assert.Empty(t, caw.CreateChannelsAndWorkspaces())
	assert.Equal(t, true, IsExistCAWByChannelIdAndWorkspaceId(channelId, workspaceId))
	assert.Equal(t, false, IsExistCAWByChannelIdAndWorkspaceId(-1, workspaceId))
	assert.Equal(t, false, IsExistCAWByChannelIdAndWorkspaceId(channelId, -1))
}
