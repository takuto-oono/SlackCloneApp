package models

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func CreateTest(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		for i := 0; i < 10; i ++ {
			u := NewUser(rand.Uint32(), "createModelUserTest1" + strconv.Itoa(i), "pass")
			assert.Empty(t, u.Create())
		}
	})
}

func IsExistUserSameUsernameAndPasswordTest(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		for i := 0; i < 10; i ++ {
			u1 := NewUser(rand.Uint32(), "isExistUserSameUsernameAndPasswordTestUser1" + strconv.Itoa(i), "pass")
			assert.Empty(t, u1.Create())
			u2 := NewUser(rand.Uint32(), "isExistUserSameUsernameAndPasswordTestUser1" + strconv.Itoa(i), "pass")
			b, err := u2.IsExistUserSameUsernameAndPassword()
			assert.Empty(t, err)
			assert.Equal(t, true, b)
		}
	})
}