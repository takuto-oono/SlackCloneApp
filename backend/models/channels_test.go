package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateChannel(t *testing.T) {
	name := "testCreateChannelModelName"
	description := "testCreateChannelModelDescription"
	is_private := true
	is_archive := false
	c := NewChannel(0, name, description, is_private, is_archive)
	assert.Empty(t, c.CreateChannel())
}