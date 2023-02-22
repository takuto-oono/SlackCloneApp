package token

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUserIdFromToken(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	for i := 0; i < 1000; i++ {
		userId := rand.Uint32()
		jwtToken, _ := GenerateToken(userId)
		returnUserId, err := GetUserIdFromToken(jwtToken)
		assert.Empty(t, err)
		fmt.Println(returnUserId, userId)
		assert.Equal(t, returnUserId, userId)
	}
}

func TestGenerateJWTToken(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	for i := 0; i < 1000; i++ {
		token, err := GenerateToken(rand.Uint32())
		assert.Empty(t, err)
		assert.NotEqual(t, "", token)
	}
}
