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

func IsExistUserSameUsernameAndPasswordTest(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	t.Run("1", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			u1 := NewUser(rand.Uint32(), "isExistUserSameUsernameAndPasswordTestUser1"+strconv.Itoa(i), "pass")
			assert.Empty(t, u1.Create())
			u2 := NewUser(rand.Uint32(), "isExistUserSameUsernameAndPasswordTestUser1"+strconv.Itoa(i), "pass")
			b, err := u2.IsExistUserSameUsernameAndPassword()
			assert.Empty(t, err)
			assert.Equal(t, true, b)
		}
	})
}

func GetUserByNameAndPasswordTest(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	t.Run("1", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			id := rand.Uint32()
			name := "loginModelUserTest1" + strconv.Itoa(i)
			u := NewUser(id, name, "pass")
			assert.Empty(t, u.Create())
			u1 := NewUser(0, name, "pass")
			assert.Empty(t, u1.GetUserByNameAndPassword())
			assert.Equal(t, u.ID, u1.ID)
		}
	})
}
