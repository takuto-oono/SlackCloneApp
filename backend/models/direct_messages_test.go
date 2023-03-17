package models

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xyproto/randomstring"
	"gorm.io/gorm"
)

func TestCreateDirectMessage(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	dm := NewDirectMessage(randomstring.EnglishFrequencyString(99), rand.Uint32(), uint(rand.Uint32()))
	err := dm.Create(db)
	assert.Empty(t, err)
}

func TestGetAllDMsByDLId(t *testing.T) {
	// if testing.Short() {
	// 	t.Skip("skipping test in short mode.")
	// }

	t.Run("1 データが存在する場合", func(t *testing.T) {
		dmCount := 4
		dms := make([]*DirectMessage, dmCount)
		sendUserId1 := rand.Uint32()
		sendUserId2 := rand.Uint32()
		dlId := uint(rand.Uint32())

		for i := 0; i < dmCount; i++ {
			if i%2 == 0 {
				dms[i] = NewDirectMessage(randomstring.EnglishFrequencyString(40), sendUserId1, dlId)
			} else {
				dms[i] = NewDirectMessage(randomstring.EnglishFrequencyString(40), sendUserId2, dlId)
			}
		}

		for _, dm := range dms {
			assert.Empty(t, dm.Create(db))
		}

		res, err := GetAllDMsByDLId(db, dlId)
		assert.Empty(t, err)
		assert.Equal(t, dmCount, len(res))

		for i := 0; i < dmCount-1; i++ {
			assert.False(t, res[i].CreatedAt.Before(res[i+1].CreatedAt))
		}
	})

	t.Run("2 データが存在しない場合", func(t *testing.T) {
		res, err := GetAllDMsByDLId(db, uint(rand.Uint64()))
		assert.Empty(t, err)
		assert.Empty(t, []DirectMessage{}, res)
	})
}

func TestGetDMById(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("1 データが存在する場合", func(t *testing.T) {
		dm := NewDirectMessage(randomstring.EnglishFrequencyString(100), rand.Uint32(), uint(rand.Uint32()))
		assert.Empty(t, dm.Create(db))
		result, err := GetDMById(db, dm.ID)
		assert.Empty(t, err)
		assert.Equal(t, dm.ID, result.ID)
		assert.Equal(t, dm.Text, result.Text)
		assert.Equal(t, dm.SendUserId, result.SendUserId)
		assert.Equal(t, dm.DMLineId, result.DMLineId)
	})

	t.Run("2 データが存在しない場合", func(t *testing.T) {
		_, err := GetDMById(db, uint(rand.Uint64()))
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestUpdateDM(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	dm := NewDirectMessage(randomstring.EnglishFrequencyString(100), rand.Uint32(), uint(rand.Uint32()))
	assert.Empty(t, dm.Create(db))
	newText := randomstring.EnglishFrequencyString(100)
	result, err := UpdateDM(db, dm.ID, newText)
	assert.Empty(t, err)
	assert.Equal(t, dm.ID, result.ID)
	assert.Equal(t, newText, result.Text)
	assert.Equal(t, dm.SendUserId, result.SendUserId)
	assert.Equal(t, dm.DMLineId, result.DMLineId)
	assert.NotEqual(t, result.CreatedAt, result.UpdatedAt)
}

func TestDeleteDM(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("1 データが存在する場合", func(t *testing.T) {
		dm := NewDirectMessage(randomstring.EnglishFrequencyString(30), rand.Uint32(), uint(rand.Uint32()))
		assert.Empty(t, dm.Create(db))
		err := dm.DeleteDM(db)
		assert.Empty(t, err)
	})
}
