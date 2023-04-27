package models

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xyproto/randomstring"
	"gorm.io/gorm"
)

func TestCreateMessage(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("1 create channel message", func(t *testing.T) {
		text := randomstring.EnglishFrequencyString(30)
		channelId := rand.Int()
		userId := rand.Uint32()
		m := NewChannelMessage(text, channelId, userId)
		assert.Empty(t, m.Create(db))
		assert.NotEqual(t, 0, m.ID)
		assert.NotEqual(t, "", m.CreatedAt)
		m = NewChannelMessage(text, channelId, userId)
		assert.Empty(t, m.Create(db))
		assert.NotEqual(t, 0, m.ID)
		assert.NotEqual(t, "", m.CreatedAt)
	})

	t.Run("2 create direct message", func(t *testing.T) {
		text := randomstring.EnglishFrequencyString(30)
		dmLineId := uint(rand.Uint32())
		userId := rand.Uint32()
		m := NewDMMessage(text, dmLineId, userId)
		assert.Empty(t, m.Create(db))
		assert.NotEqual(t, 0, m.ID)
		assert.NotEqual(t, "", m.CreatedAt)
		m = NewDMMessage(text, dmLineId, userId)
		assert.Empty(t, m.Create(db))
		assert.NotEqual(t, 0, m.ID)
		assert.NotEqual(t, "", m.CreatedAt)
	})
}

func TestGetMessagesByChannelId(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("1 messageが存在する場合", func(t *testing.T) {
		testNum := 100
		texts := make([]string, testNum)
		for i := 0; i < testNum; i++ {
			texts[i] = randomstring.EnglishFrequencyString(30)
		}
		channelId := rand.Int()
		userId := rand.Uint32()

		for _, text := range texts {
			m := NewChannelMessage(text, channelId, userId)
			assert.Empty(t, m.Create(db))
		}

		messages, err := GetMessagesByChannelId(db, channelId)
		assert.Empty(t, err)
		assert.Equal(t, testNum, len(messages))
		for i := 0; i < testNum-1; i++ {
			assert.True(t, messages[i+1].CreatedAt.Before(messages[i].CreatedAt))
		}
		for _, m := range messages {
			assert.Equal(t, channelId, m.ChannelId)
			assert.Equal(t, userId, m.UserId)
			assert.Contains(t, texts, m.Text)
		}
	})

	t.Run("2 messageが存在しない場合", func(t *testing.T) {
		messages, err := GetMessagesByChannelId(db, -1)
		assert.Empty(t, err)
		assert.Equal(t, 0, len(messages))
	})
}

func TestGetMessagesByDlId(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("1 messageが存在する場合", func(t *testing.T) {
		testNum := 100
		texts := make([]string, testNum)
		for i := 0; i < testNum; i++ {
			texts[i] = randomstring.EnglishFrequencyString(30)
		}
		dlId := uint(rand.Uint32())
		userId := rand.Uint32()

		for _, text := range texts {
			m := NewDMMessage(text, dlId, userId)
			assert.Empty(t, m.Create(db))
		}

		messages, err := GetMessagesByDLId(db, dlId)
		assert.Empty(t, err)
		assert.Equal(t, testNum, len(messages))
		for i := 0; i < testNum-1; i++ {
			assert.True(t, messages[i+1].CreatedAt.Before(messages[i].CreatedAt))
		}
		for _, m := range messages {
			assert.Equal(t, dlId, m.DMLineId)
			assert.Equal(t, userId, m.UserId)
			assert.Contains(t, texts, m.Text)
		}
	})

	t.Run("2 messageが存在しない場合", func(t *testing.T) {
		messages, err := GetMessagesByDLId(db, uint(rand.Uint32()))
		assert.Empty(t, err)
		assert.Equal(t, 0, len(messages))
	})
}

func TestGetMessageById(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("1 データが存在する場合", func(t *testing.T) {
		m := NewChannelMessage(randomstring.EnglishFrequencyString(100), rand.Int(), rand.Uint32())
		assert.Empty(t, m.Create(db))
		res, err := GetMessageById(db, m.ID)
		assert.Empty(t, err)
		assert.Equal(t, m.ID, res.ID)
		assert.Equal(t, m.Text, res.Text)
		assert.Equal(t, m.ChannelId, res.ChannelId)
		assert.Equal(t, m.UserId, res.UserId)
		assert.Equal(t, m.CreatedAt.UTC(), res.CreatedAt)
		assert.Equal(t, m.UpdatedAt.UTC(), res.UpdatedAt)
	})

	t.Run("2 データが存在しない場合", func(t *testing.T) {
		_, err := GetMessageById(db, uint(rand.Uint32()))
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestUpdateMessageText(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("1 データが存在する場合", func(t *testing.T) {
		m := NewChannelMessage(randomstring.EnglishFrequencyString(100), rand.Int(), rand.Uint32())
		assert.Empty(t, m.Create(db))

		newText := randomstring.EnglishFrequencyString(100)
		res, err := UpdateMessageText(db, m.ID, newText)
		assert.Empty(t, err)
		assert.Equal(t, m.ID, res.ID)
		assert.Equal(t, newText, res.Text)
		assert.Equal(t, m.ChannelId, res.ChannelId)
		assert.Equal(t, m.UserId, res.UserId)
		assert.Equal(t, m.CreatedAt.UTC(), res.CreatedAt)
		assert.True(t, m.UpdatedAt.UTC().Before(res.UpdatedAt))

		res2, err := GetMessageById(db, m.ID)
		assert.Empty(t, err)
		assert.Equal(t, m.ID, res2.ID)
		assert.Equal(t, newText, res2.Text)
		assert.Equal(t, m.ChannelId, res2.ChannelId)
		assert.Equal(t, m.UserId, res2.UserId)
		assert.Equal(t, m.CreatedAt.UTC(), res2.CreatedAt)
		assert.Equal(t, res.UpdatedAt, res2.UpdatedAt)
	})

	t.Run("2 データが存在しない場合", func(t *testing.T) {
		_, err := UpdateMessageText(db, uint(rand.Uint32()), randomstring.EnglishFrequencyString(100))
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

}

func TestGetMessagesThreadId(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	m := NewChannelMessage(randomstring.EnglishFrequencyString(30), rand.Int(), rand.Uint32())
	assert.Empty(t, m.Create(db))
	message, _ := GetMessageById(db, m.ID)
	assert.Equal(t, uint(0), message.ThreadId)
}

func TestUpdateMessageThreadId(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	
	m := NewChannelMessage(randomstring.EnglishFrequencyString(30), rand.Int(), rand.Uint32())
	assert.Empty(t, m.Create(db))
	threadId := uint(rand.Uint32())
	assert.Empty(t, m.UpdateMessageThreadId(db, threadId))
	message, _ := GetMessageById(db, m.ID)
	assert.Equal(t, threadId, message.ThreadId)
}
