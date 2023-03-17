package models

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateChannelAndUsers(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	channelId := rand.Int()
	userId := rand.Uint32()
	isAdmin := true

	cau := NewChannelsAndUses(channelId, userId, isAdmin)
	assert.Empty(t, cau.Create(db))

	cau.IsAdmin = false
	assert.NotEmpty(t, cau.Create(db))

	cau.ChannelId = rand.Int()
	assert.Empty(t, cau.Create(db))

	cau.ChannelId = channelId
	cau.UserId = rand.Uint32()
	assert.Empty(t, cau.Create(db))
}

func TestGetCAUByChannelIdAndUserId(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	channelId := rand.Int()
	userId := rand.Uint32()

	cau := NewChannelsAndUses(channelId, userId, false)
	assert.Empty(t, cau.Create(db))

	res, err := GetCAUByChannelIdAndUserId(db, channelId, userId)
	assert.Empty(t, err)
	assert.Equal(t, *cau, res)

	_, err = GetCAUByChannelIdAndUserId(db, -1, userId)
	assert.NotEmpty(t, err)
}

func TestDeleteCAU(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	cau := NewChannelsAndUses(rand.Int(), rand.Uint32(), true)
	assert.Empty(t, cau.Create(db))
	assert.Empty(t, cau.Delete(db))
}

func TestDeleteCAUsByChannelId(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	cau := NewChannelsAndUses(rand.Int(), rand.Uint32(), true)
	assert.Empty(t, cau.Create(db))
	assert.Empty(t, cau.Delete(db))

	cau = NewChannelsAndUses(rand.Int(), rand.Uint32(), true)
	assert.Empty(t, cau.Create(db))
	cau.IsAdmin = false
	assert.Empty(t, cau.Delete(db))

	channelId := rand.Int()
	userId := rand.Uint32()
	cau = NewChannelsAndUses(channelId, userId, true)
	assert.Empty(t, cau.Create(db))

	cau.ChannelId = -1
	assert.Empty(t, cau.Delete(db))

	cau.ChannelId = channelId
	assert.Empty(t, cau.Delete(db))
}

func TestGetCAUsByUserId(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	// 1. データが存在する場合
	// 2. データが存在しない場合

	t.Run("1 データが存在する場合", func(t *testing.T) {
		cauCount := 10
		userId := rand.Uint32()
		caus := make([]ChannelsAndUsers, cauCount)
		for i := 0; i < cauCount; i++ {
			cau := NewChannelsAndUses(int(rand.Uint64()), userId, false)
			assert.Empty(t, cau.Create(db))
			caus[i] = *cau
		}
		res, err := GetCAUsByUserId(db, userId)
		assert.Empty(t, err)
		assert.Equal(t, cauCount, len(res))
		for _, cau := range caus {
			assert.Contains(t, res, cau)
		}
	})

	t.Run("2 データが存在しない場合", func(t *testing.T) {
		res, err := GetCAUsByUserId(db, rand.Uint32())
		assert.Empty(t, err)
		assert.Equal(t, 0, len(res))
	})
}

func TestGetCAUsByChannelId(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	// 1. データが存在する場合
	// 2. データが存在しない場合

	t.Run("1 データが存在する場合", func(t *testing.T) {
		cauCount := 10
		channelId := int(rand.Uint32())
		caus := make([]ChannelsAndUsers, cauCount)
		for i := 0; i < cauCount; i++ {
			cau := NewChannelsAndUses(channelId, rand.Uint32(), false)
			assert.Empty(t, cau.Create(db))
			caus[i] = *cau
		}
		res, err := GetCAUsByChannelId(db, channelId)
		assert.Empty(t, err)
		assert.Equal(t, cauCount, len(res))
		for _, cau := range caus {
			assert.Contains(t, res, cau)
		}
	})

	t.Run("2 データが存在しない場合", func(t *testing.T) {
		res, err := GetCAUsByChannelId(db, int(rand.Uint32()))
		assert.Empty(t, err)
		assert.Equal(t, 0, len(res))
	})
}
