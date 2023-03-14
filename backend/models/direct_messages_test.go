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

	dm := NewDirectMessage(randomstring.EnglishFrequencyString(99), rand.Uint32(), uint(rand.Uint64()))
	res := dm.Create()
	assert.NotEqual(t, uint(0), dm.ID)
	assert.Empty(t, res.Error)
}

func TestGetAllDMsByDLId(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("1 データが存在する場合", func(t *testing.T) {
		dmCount := 4
		dms := make([]*DirectMessage, dmCount)
		sendUserId1 := rand.Uint32()
		sendUserId2 := rand.Uint32()
		dlId := uint(rand.Uint64())

		for i := 0; i < dmCount; i++ {
			if i%2 == 0 {
				dms[i] = NewDirectMessage(randomstring.EnglishFrequencyString(40), sendUserId1, dlId)
			} else {
				dms[i] = NewDirectMessage(randomstring.EnglishFrequencyString(40), sendUserId2, dlId)
			}
		}

		for _, dm := range dms {
			assert.Empty(t, dm.Create().Error)
		}

		res, err := GetAllDMsByDLId(dlId)
		assert.Empty(t, err)
		assert.Equal(t, dmCount, len(res))

		for i := 0; i < dmCount-1; i++ {
			assert.False(t, res[i].CreatedAt.Before(res[i+1].CreatedAt))
		}
	})

	t.Run("2 データが存在しない場合", func(t *testing.T) {
		res, err := GetAllDMsByDLId(uint(rand.Uint64()))
		assert.Empty(t, err)
		assert.Empty(t, []DirectMessage{}, res)
	})
}

func TestGetDMById(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("1 データが存在する場合", func(t *testing.T) {
		dm := NewDirectMessage(randomstring.EnglishFrequencyString(100), rand.Uint32(), uint(rand.Uint64()))
		assert.Empty(t, dm.Create().Error)
		result, err := GetDMById(dm.ID)
		assert.Empty(t, err)
		assert.Equal(t, dm.ID, result.ID)
		assert.Equal(t, dm.Text, result.Text)
		assert.Equal(t, dm.SendUserId, result.SendUserId)
		assert.Equal(t, dm.DMLineId, result.DMLineId)
	})

	t.Run("2 データが存在しない場合", func(t *testing.T) {
		_, err := GetDMById(uint(rand.Uint64()))
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestUpdateDM(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	dm := NewDirectMessage(randomstring.EnglishFrequencyString(100), rand.Uint32(), uint(rand.Uint64()))
	assert.Empty(t, dm.Create().Error)
	newText := randomstring.EnglishFrequencyString(100)
	result, err := UpdateDM(dm.ID, newText)
	assert.Empty(t, err)
	assert.Equal(t, dm.ID, result.ID)
	assert.Equal(t, newText, result.Text)
	assert.Equal(t, dm.SendUserId, result.SendUserId)
	assert.Equal(t, dm.DMLineId, result.DMLineId)
	assert.NotEqual(t, result.CreatedAt, result.UpdatedAt)

}
