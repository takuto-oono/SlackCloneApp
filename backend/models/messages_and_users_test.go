package models

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateMAU(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	mau := NewMessageAndUser(uint(rand.Uint32()), rand.Uint32(), false)
	assert.Empty(t, mau.Create(db))
}

func TestGetMAUByUserIDAndIsRead(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	userID := rand.Uint32()
	maus := make([]MessageAndUser, 10)
	for i := 0; i < 10; i++ {
		mau := NewMessageAndUser(uint(rand.Uint32()), userID, (i%2 == 0))
		assert.Empty(t, mau.Create(db))
		maus[i] = *mau
	}

	res, err := GetMAUByUserIDAndIsRead(db, userID, false)
	assert.Empty(t, err)
	assert.Equal(t, 5, len(res))
	for _, r := range res {
		assert.False(t, r.IsRead)
		isExist := false
		for _, mau := range maus {
			if mau.MessageID == r.MessageID {
				isExist = true
				break
			}
		}
		assert.True(t, isExist)
	}
}

func TestUpdateIsRead(t *testing.T) {
	mau := NewMessageAndUser(uint(rand.Uint32()), rand.Uint32(), false)
	assert.Empty(t, mau.Create(db))

	assert.Empty(t, mau.UpdateIsRead(db, true))
	assert.True(t, mau.IsRead)
	
	res, err := GetMAUByUserIDAndIsRead(db, mau.UserID, true)
	assert.Empty(t, err)
	assert.Equal(t, 1, len(res))
}
