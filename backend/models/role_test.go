package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsExist4Roles(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	roleNames := []string{
		"Workspace Primary Owner",
		"Workspace Owners",
		"Workspace Admins",
		"Full members",
	}
	for i, n := range roleNames {
		r, _ := GetRoleById(db, i + 1)
		assert.Equal(t, n, r.Name)
	}
}
