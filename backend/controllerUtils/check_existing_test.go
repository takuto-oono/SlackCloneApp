package controllerUtils

import (
	"backend/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsExistCAUByChannelIdAndUserId(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	cau := models.NewChannelsAndUses(27393769379, uint32(593593926), false)
	assert.Empty(t, cau.Create())

	b, err := IsExistCAUByChannelIdAndUserId(cau.ChannelId, cau.UserId)
	assert.Equal(t, true, b)
	assert.Empty(t, err)

	b, err = IsExistCAUByChannelIdAndUserId(-1, cau.UserId)
	assert.Equal(t, false, b)
	assert.NotEmpty(t, err)
}