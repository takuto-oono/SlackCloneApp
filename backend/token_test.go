package main

import (
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"backend/token"
)

func TestGenerateJWTToken(t *testing.T) {
	uuidWithHyphen := uuid.New()
	for i := 0; i < 100000; i++ {
		uuid := strings.Replace(uuidWithHyphen.String(), "-", "", -1)
		token, err := token.GenerateToken(uuid)
		assert.Empty(t, err)
		assert.NotEqual(t, "", token)
	}
}
