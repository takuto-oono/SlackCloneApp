package models

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCreateTAU(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	tau := NewThreadAndUser(rand.Uint32(), uint(rand.Uint32()))
	assert.Empty(t, tau.Create(db))
}

func TestGetTAUByThreadIdAndUserId(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	tau := NewThreadAndUser(rand.Uint32(), uint(rand.Uint32()))
	assert.Empty(t, tau.Create(db))
	res, err := GetTAUByThreadIdAndUserId(db, tau.UserId, tau.ThreadId)
	assert.Empty(t, err)
	assert.Equal(t, tau.UserId, res.UserId)
	assert.Equal(t, tau.ThreadId, res.ThreadId)

	_, err = GetTAUByThreadIdAndUserId(db, rand.Uint32(), uint(rand.Uint32()))
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestGetTAUsByUserId(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	userId := rand.Uint32()
	taus := make([]ThreadAndUser, 10)
	for i := 0; i < 10; i++ {
		tau := NewThreadAndUser(userId, uint(rand.Uint32()))
		tau.Create(db)
		taus[i] = *tau
	}
	res, err := GetTAUsByUserId(db, userId)
	assert.Empty(t, err)
	assert.Equal(t, 10, len(res))

	for _, r := range res {
		assert.Equal(t, userId, r.UserId)
		isExist := false
		for _, tau := range taus {
			if tau.ThreadId == r.ThreadId {
				isExist = true
				break
			}
		}
		assert.True(t, isExist)
	}
}

func TestGetTAUsByThreadId(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	threadId := uint(rand.Uint32())
	taus := make([]ThreadAndUser, 10)
	for i := 0; i < 10; i++ {
		tau := NewThreadAndUser(rand.Uint32(), threadId)
		tau.Create(db)
		taus[i] = *tau
	}
	res, err := GetTAUsByThreadId(db, threadId)
	assert.Empty(t, err)
	assert.Equal(t, 10, len(res))

	for _, r := range res {
		assert.Equal(t, threadId, r.ThreadId)
		isExist := false
		for _, tau := range taus {
			if tau.UserId == r.UserId {
				isExist = true
				break
			}
		}
		assert.True(t, isExist)
	}
}
