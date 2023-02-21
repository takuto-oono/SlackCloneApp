package models

import (
	"backend/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateMessage(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	text := "test message"
	channelId := 7675071751
	userId := uint32(3571521121)
	m := NewMessage(text, channelId, userId)
	assert.Empty(t, m.Create())
	assert.NotEqual(t, 0, m.ID)
	assert.NotEqual(t, "", m.Date)
	m = NewMessage(text, channelId, userId)
	assert.Empty(t, m.Create())
	assert.NotEqual(t, 0, m.ID)
	assert.NotEqual(t, "", m.Date)
}

func TestGetMessagesByChannelId(t *testing.T) {
	// if testing.Short() {
	// 	t.Skip("skipping test in short mode.")
	// }

	t.Run("1 messageが存在する場合", func(t *testing.T) {
		testNum := 100
		texts := make([]string, testNum)
		for i := 0; i < testNum; i++ {
			texts[i] = "testGetMessagesByChannelId"
		}
		channelId := 37590793729
		userId := uint32(536357900)

		for _, text := range texts {
			m := NewMessage(text, channelId, userId)
			assert.Empty(t, m.Create())
		}

		messages, err := GetMessagesByChannelId(channelId)
		assert.Empty(t, err)
		assert.Equal(t, testNum, len(messages))
		for i := 0; i < testNum-1; i++ {
			d1, err1 := utils.TimeFromString(messages[i].Date)
			d2, err2 := utils.TimeFromString(messages[i+1].Date)
			assert.Empty(t, err1)
			assert.Empty(t, err2)
			assert.True(t, d2.Before(d1))
		}
		for _, m := range messages {
			assert.Equal(t, channelId, m.ChannelId)
			assert.Equal(t, userId, m.UserId)
			assert.Equal(t, texts[0], m.Text)
		}	
	})

	t.Run("2 messageが存在しない場合", func(t *testing.T) {
		messages, err := GetMessagesByChannelId(-1)
		assert.Empty(t, err)
		assert.Equal(t, 0, len(messages))
	})
}
