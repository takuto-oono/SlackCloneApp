package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateChannel(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	name := "testCreateChannelModelName"
	description := "testCreateChannelModelDescription"
	is_private := true
	is_archive := false
	workspaceId := 7397593
	c := NewChannel(0, name, description, is_private, is_archive, workspaceId)
	assert.Empty(t, c.Create())
}

func TestGetChannelById(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	name := "testGetChannelByIdName"
	description := "testGetChannelByIdDescription"
	is_private := true
	is_archive := false
	workspaceId := 3999526
	c := NewChannel(0, name, description, is_private, is_archive, workspaceId)
	assert.Empty(t, c.Create())
	assert.NotEqual(t, 0, c.ID)
	c2, err := GetChannelById(c.ID)
	assert.Empty(t, err)
	assert.Equal(t, *c, c2)

	_, err = GetChannelById(-1)
	assert.NotEmpty(t, err)
}

func TestDeleteChannel(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	c := NewChannel(0, "testDeleteChannelName", "", true, false, 58430850380)
	assert.Empty(t, c.Create())
	channelId := c.ID
	assert.Empty(t, c.Delete())

	_, err := GetChannelById(channelId)
	assert.NotEmpty(t, err)
}
