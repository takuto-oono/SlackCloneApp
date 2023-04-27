package models

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTAM(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	tam := NewThreadAndMessage(uint(rand.Uint32()), uint(rand.Uint32()))
	assert.Empty(t, tam.Create(db))
}

func TestGetTAMsByThreadID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("1 データが存在する場合", func(t *testing.T) {
		tams := make([]ThreadAndMessage, 10)
		threadId := uint(rand.Uint32())
		for i := 0; i < 10; i ++ {
			tam := NewThreadAndMessage(threadId, uint(rand.Uint32()))
			tam.Create(db)
			tams[i] = *tam
		}
		res, err := GetTAMsByThreadId(db, threadId)
		assert.Empty(t, err)
		assert.Equal(t, 10, len(res))
		for _, tam := range tams {
			isExist := false
			for _, r := range res {
				if tam.MessageId == r.MessageId {
					isExist = true
					break
				}
			}
			assert.True(t, isExist)
		}
	})
	
	t.Run("2 データが存在しない場合", func(t *testing.T) {
		res, err := GetTAMsByThreadId(db, uint(rand.Uint32()))
		assert.Empty(t, err)
		assert.Equal(t, 0, len(res))
	})
}

func TestDeleteTAM(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	tam := NewThreadAndMessage(uint(rand.Uint32()), uint(rand.Uint32()))
	assert.Empty(t, tam.Create(db))
	assert.Empty(t, tam.Delete(db))
}
