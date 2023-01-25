package main

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"

	"backend/token"
)

func TestGetUserIdFromToken(t *testing.T) {
	for i := 0; i < 100000; i ++ {
		userId := rand.Uint32()
		jwtToken, _ := token.GenerateToken(userId)
		returnUserId, err := token.GetUserIdFromToken(jwtToken)
		assert.Empty(t, err)
		assert.Equal(t, returnUserId, userId)
	}
}

func TestGenerateJWTToken(t *testing.T) {
	for i := 0; i < 1000000; i++ {
		token, err := token.GenerateToken(rand.Uint32())
		assert.Empty(t, err)
		assert.NotEqual(t, "", token)
	}
}
