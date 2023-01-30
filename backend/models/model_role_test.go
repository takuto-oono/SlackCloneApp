package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsExist4Roles(t *testing.T) {
	roleNames := []string{
		"Workspace Primary Owner",
		"Workspace Owners",
		"Workspace Admins",
		"Full members",
	}
	for i, n := range roleNames {
		r, _ := GetRoleById(i + 1)
		assert.Equal(t, r.Name, n)
	}
}
