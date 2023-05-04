package models

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateMention(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	mention := NewMention(rand.Uint32(), uint(rand.Uint32()))
	assert.Empty(t, mention.Create(db))
}

func TestGetMentionsByMessageID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	messageID := uint(rand.Uint32())
	mentions := make([]Mention, 10)
	for i := 0; i < 10; i++ {
		mention := NewMention(rand.Uint32(), messageID)
		mention.Create(db)
		mentions[i] = *mention
	}
	res, err := GetMentionsByMessageID(db, messageID)
	assert.Empty(t, err)
	assert.Equal(t, 10, len(res))
	for _, men := range mentions {
		isExist := false
		for _, r := range res {
			if men.UserID == r.UserID {
				isExist = true
				break
			}
		}
		assert.True(t, isExist)
	}
}

func TestGetMentionsByUserID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	userID := rand.Uint32()
	mentions := make([]Mention, 10)
	for i := 0; i < 10; i++ {
		mention := NewMention(userID, uint(rand.Uint32()))
		mention.Create(db)
		mentions[i] = *mention
	}
	res, err := GetMentionsByUserID(db, userID)
	assert.Empty(t, err)
	assert.Equal(t, 10, len(res))
	for _, men := range mentions {
		isExist := false
		for _, r := range res {
			if men.MessageID == r.MessageID {
				isExist = true
				break
			}
		}
		assert.True(t, isExist)
	}
}
