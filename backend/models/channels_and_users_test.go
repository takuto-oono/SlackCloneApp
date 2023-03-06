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
	assert.Empty(t, cau.Create())

	cau.IsAdmin = false
	assert.NotEmpty(t, cau.Create())

	cau.ChannelId = rand.Int()
	assert.Empty(t, cau.Create())

	cau.ChannelId = channelId
	cau.UserId = rand.Uint32()
	assert.Empty(t, cau.Create())
}

func TestGetCAUByChannelIdAndUserId(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	channelId := rand.Int()
	userId := rand.Uint32()

	cau := NewChannelsAndUses(channelId, userId, false)
	assert.Empty(t, cau.Create())

	res, err := GetCAUByChannelIdAndUserId(channelId, userId)
	assert.Empty(t, err)
	assert.Equal(t, *cau, res)

	_, err = GetCAUByChannelIdAndUserId(-1, userId)
	assert.NotEmpty(t, err)
}

func TestIsExistCAUByChannelIdAndUserId(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	channelId := rand.Int()
	userId := rand.Uint32()

	cau := NewChannelsAndUses(channelId, userId, false)
	assert.Empty(t, cau.Create())

	assert.Equal(t, true, IsExistCAUByChannelIdAndUserId(channelId, userId))
	assert.Equal(t, false, IsExistCAUByChannelIdAndUserId(-1, userId))
}

func TestIsAdminUserInChannel(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	cau := NewChannelsAndUses(rand.Int(), rand.Uint32(), true)
	assert.Empty(t, cau.Create())

	assert.Equal(t, true, IsAdminUserInChannel(cau.ChannelId, cau.UserId))
	assert.Equal(t, false, IsAdminUserInChannel(cau.ChannelId, rand.Uint32()))

	cau = NewChannelsAndUses(cau.ChannelId, rand.Uint32(), false)
	assert.Empty(t, cau.Create())
	assert.Equal(t, false, IsAdminUserInChannel(cau.ChannelId, cau.UserId))
}

func TestDeleteUserFromChannel(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	cau := NewChannelsAndUses(rand.Int(), rand.Uint32(), true)
	assert.Empty(t, cau.Create())
	assert.Empty(t, cau.Delete())

	cau = NewChannelsAndUses(rand.Int(), rand.Uint32(), true)
	assert.Empty(t, cau.Create())
	cau.IsAdmin = false
	assert.Empty(t, cau.Delete())
	assert.Equal(t, false, IsExistCAUByChannelIdAndUserId(cau.ChannelId, cau.UserId))

	channelId := rand.Int()
	userId := rand.Uint32()
	cau = NewChannelsAndUses(channelId, userId, true)
	assert.Empty(t, cau.Create())
	assert.Equal(t, true, IsExistCAUByChannelIdAndUserId(channelId, userId))

	cau.ChannelId = -1
	assert.Empty(t, cau.Delete())
	assert.Equal(t, true, IsExistCAUByChannelIdAndUserId(channelId, userId))

	cau.ChannelId = channelId
	assert.Empty(t, cau.Delete())
	assert.Equal(t, false, IsExistCAUByChannelIdAndUserId(channelId, userId))
}

func TestDeleteCAUByChannelId(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	channelId := rand.Int()
	userIds := make([]uint32, 20)
	for i := 0; i < 20; i++ {
		userIds[i] = rand.Uint32()
	}

	for _, userId := range userIds {
		cau := NewChannelsAndUses(channelId, userId, false)
		assert.Empty(t, cau.Create())
	}

	assert.Empty(t, DeleteCAUByChannelId(channelId))

	for _, userId := range userIds {
		assert.Equal(t, false, IsExistCAUByChannelIdAndUserId(channelId, userId))
	}
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
			assert.Empty(t, cau.Create())
			caus[i] = *cau
		}
		res, err := GetCAUsByUserId(userId)
		assert.Empty(t, err)
		assert.Equal(t, cauCount, len(res))
		for _, cau := range caus {
			assert.Contains(t, res, cau)
		}
	})

	t.Run("2 データが存在しない場合", func(t *testing.T) {
		res, err := GetCAUsByUserId(rand.Uint32())
		assert.Empty(t, err)
		assert.Equal(t, 0, len(res))
	})
}
