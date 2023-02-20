package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateMessage(t *testing.T) {
	// if testing.Short() {
	// 	t.Skip("skipping test in short mode.")
	// }
	text := "test message"
	channelId := 7675071751
	userId := uint32(3571521121)
	m := NewMessage(text, channelId, userId)
	assert.Empty(t, m.Create())
	assert.NotEqual(t, 0, m.ID)
	assert.NotEqual(t, "", m.Date)
	m = NewMessage(text, channelId, userId)
	assert.Empty(t, m.Create())
	assert.NotEqual(t, 0, m.ID)
	assert.NotEqual(t, "", m.Date)
}
