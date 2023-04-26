package models

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCreateThread(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	thread := NewThread(uint(rand.Uint32()))
	assert.Empty(t, thread.Create(db))
}

func TestEditUpdatedAt(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	thread := NewThread(uint(rand.Uint32()))
	assert.Empty(t, thread.Create(db))
	assert.Empty(t, thread.EditUpdatedAt(db))
}

func TestGetThreadById(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("1 データがある場合", func(t *testing.T) {
		thread := NewThread(uint(rand.Uint32()))
		assert.Empty(t, thread.Create(db))
		res, err := GetThreadById(db, thread.ID)
		assert.Empty(t, err)
		assert.Equal(t, thread.ID, res.ID)
		assert.Equal(t, thread.ParentMessageId, res.ParentMessageId)
	})

	t.Run("2 データがない場合", func(t *testing.T) {
		_, err := GetThreadById(db, uint(rand.Uint32()))
		assert.Error(t, gorm.ErrRecordNotFound, err)
	})
}

func TestGetThreadByParentMessageId(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("1 データがある場合", func(t *testing.T) {
		thread := NewThread(uint(rand.Uint32()))
		assert.Empty(t, thread.Create(db))
		res, err := GetThreadByParentMessageId(db, thread.ParentMessageId)
		assert.Empty(t, err)
		assert.Equal(t, thread.ID, res.ID)
		assert.Equal(t, thread.ParentMessageId, res.ParentMessageId)
	})

	t.Run("2 データがない場合", func(t *testing.T) {
		_, err := GetThreadByParentMessageId(db, uint(rand.Uint32()))
		assert.Error(t, gorm.ErrRecordNotFound, err)
	})
}
