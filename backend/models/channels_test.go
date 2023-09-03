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
	c := NewChannel(name, description, is_private, is_archive, workspaceId)
	assert.Empty(t, c.Create(db))
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
	c := NewChannel(name, description, is_private, is_archive, workspaceId)
	assert.Empty(t, c.Create(db))
	assert.NotEqual(t, 0, c.ID)
	c2, err := GetChannelById(db, c.ID)
	assert.Empty(t, err)
	assert.Equal(t, *c, c2)

	_, err = GetChannelById(db, -1)
	assert.NotEmpty(t, err)
}

func TestDeleteChannel(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	c := NewChannel(randomstring.EnglishFrequencyString(30), "", true, false, rand.Int())
	assert.Empty(t, c.Create(db))
	channelId := c.ID
	assert.Empty(t, c.Delete(db))

	_, err := GetChannelById(db, channelId)
	assert.NotEmpty(t, err)
}

func TestGetAllChannels(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("1 データが存在する場合", func(t *testing.T) {
		chCnt := 100
		channels := make([]Channel, chCnt)
		for i := 0; i < chCnt; i++ {
			ch := NewChannel(
				randomstring.EnglishFrequencyString(30),
				"",
				false,
				false,
				int(rand.Uint64()),
			)
			assert.Empty(t, ch.Create(db))
			channels[i] = *ch
		}
		res, err := GetAllChannels(db)
		assert.Empty(t, err)
		assert.Equal(t, chCnt, len(res))
		for _, ch := range channels {
			assert.Contains(t, res, ch)
		}
	})
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
			ch := NewChannel(channelName, "des", false, false, workspaceId)
			assert.Empty(t, ch.Create(db))
			channels[i] = *ch
		}
		chs, err := GetChannelsByWorkspaceId(db, workspaceId)
		assert.Empty(t, err)
		assert.Equal(t, len(channels), len(chs))
		for _, ch := range channels {
			assert.Contains(t, chs, ch)
		}
	})

	t.Run("2 データが存在しない場合", func(t *testing.T) {
		chs, err := GetChannelsByWorkspaceId(db, int(rand.Uint64()))
		assert.Empty(t, err)
		assert.Equal(t, 0, len(chs))
	})
}
