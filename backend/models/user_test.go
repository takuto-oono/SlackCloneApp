package models

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xyproto/randomstring"
)

func CreateTest(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	t.Run("1", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			u := NewUser(rand.Uint32(), randomstring.EnglishFrequencyString(30), "pass")
			assert.Empty(t, u.Create())
		}
	})
}

func GetUserByNameAndPasswordTest(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	t.Run("1", func(t *testing.T) {
		password := "pass"
		for i := 0; i < 10; i++ {
			id := rand.Uint32()
			name := randomstring.EnglishFrequencyString(30)
			u := NewUser(id, name, password)
			assert.Empty(t, u.Create())
			u1, err := GetUserByNameAndPassword(name, password)
			assert.Empty(t, err)
			assert.Equal(t, *u, u1)
			_, err = GetUserByNameAndPassword("wrong name", password)
			assert.NotEmpty(t, err)
			_, err = GetUserByNameAndPassword(name, "wrong password")
			assert.NotEmpty(t, err)
		}
	})
}

func TestGetUsers(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	_, err := GetUsers()
	assert.Empty(t, err)
}
