package models

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func CreateTest(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	t.Run("1", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			u := NewUser(rand.Uint32(), "createModelUserTest1"+strconv.Itoa(i), "pass")
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
			name := "loginModelUserTest1" + strconv.Itoa(i)
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
