package models

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xyproto/randomstring"
)

func TestCreateChannel(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	name := randomstring.EnglishFrequencyString(30)
	description := randomstring.EnglishFrequencyString(30)
	is_private := true
	is_archive := false
	workspaceId := rand.Int()
	c := NewChannel(0, name, description, is_private, is_archive, workspaceId)
	assert.Empty(t, c.Create())
}

func TestGetChannelById(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	name := randomstring.EnglishFrequencyString(30)
	description := randomstring.EnglishFrequencyString(30)
	is_private := true
	is_archive := false
	workspaceId := rand.Int()
	c := NewChannel(0, name, description, is_private, is_archive, workspaceId)
	assert.Empty(t, c.Create())
	assert.NotEqual(t, 0, c.ID)
	c2, err := GetChannelById(c.ID)
	assert.Empty(t, err)
	assert.Equal(t, *c, c2)

	_, err = GetChannelById(-1)
	assert.NotEmpty(t, err)
}

func TestDeleteChannel(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	c := NewChannel(0, randomstring.EnglishFrequencyString(30), "", true, false, rand.Int())
	assert.Empty(t, c.Create())
	channelId := c.ID
	assert.Empty(t, c.Delete())

	_, err := GetChannelById(channelId)
	assert.NotEmpty(t, err)
}

func TestGetChannelsByWorkspaceId(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	// 1. データが存在する場合
	// 2. データが存在しない場合

	t.Run("1 データが存在する場合", func(t *testing.T) {
		channelCount := 10
		workspaceId := int(rand.Uint64())
		channels := make([]Channel, channelCount)
		for i := 0; i < channelCount; i++ {
			channelName := randomstring.EnglishFrequencyString(30)
			ch := NewChannel(0, channelName, "des", false, false, workspaceId)
			assert.Empty(t, ch.Create())
			channels[i] = *ch
		}
		chs, err := GetChannelsByWorkspaceId(workspaceId)
		assert.Empty(t, err)
		assert.Equal(t, len(channels), len(chs))
		for _, ch := range channels {
			assert.Contains(t, chs, ch)
		}
	})

	t.Run("2 データが存在しない場合", func(t *testing.T) {
		chs, err := GetChannelsByWorkspaceId(int(rand.Uint64()))
		assert.Empty(t, err)
		assert.Equal(t, 0, len(chs))
	})
}
