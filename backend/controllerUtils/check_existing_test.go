package controllerUtils

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"backend/models"
)

func TestIsExistCAUByChannelIdAndUserId(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	cau := models.NewChannelsAndUses(27393769379, uint32(593593926), false)
	assert.Empty(t, cau.Create())

	b, err := IsExistCAUByChannelIdAndUserId(cau.ChannelId, cau.UserId)
	assert.Equal(t, true, b)
	assert.Empty(t, err)

	b, err = IsExistCAUByChannelIdAndUserId(-1, cau.UserId)
	assert.Equal(t, false, b)
	assert.NotEmpty(t, err)
}

func TestIsExistUserSameUsernameAndPassword(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	names := make([]string, 10)
	for i := 0; i < 10; i++ {
		names[i] = "testIsExistUserSameUsernameAndPasswordUsername" + strconv.Itoa(i)
	}

	for _, name := range names {
		u := models.NewUser(rand.Uint32(), name, "pass")
		assert.Empty(t, u.Create())
	}

	for _, name := range names {
		assert.Equal(t, true, IsExistUserSameUsernameAndPassword(name, "pass"))
		assert.Equal(t, false, IsExistUserSameUsernameAndPassword(name, "wrong pass"))
		assert.Equal(t, false, IsExistUserSameUsernameAndPassword(name+" wrong name", "pass"))
	}
}

func TestIsExistWorkspaceAndUser(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	w := models.NewWorkspace(0, "testIsExistWorkspaceAndUser", 4)
	w.Create()
	assert.Equal(t, true, IsExistWorkspaceById(w.ID))
	assert.Equal(t, false, IsExistWorkspaceById(-1))
}

func TestIsExistWAUByWorkspaceIdAndUserId(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	waus := make([]models.WorkspaceAndUsers, 10)
	for i := 0; i < 10; i++ {
		waus[i] = *models.NewWorkspaceAndUsers(int(rand.Uint32()), rand.Uint32(), rand.Int()%4+1)
	}
	for _, wau := range waus {
		assert.Empty(t, wau.Create())
	}

	for _, wau := range waus {
		assert.True(t, IsExistWAUByWorkspaceIdAndUserId(wau.WorkspaceId, wau.UserId))
		assert.False(t, IsExistWAUByWorkspaceIdAndUserId(rand.Int(), wau.UserId))
		assert.False(t, IsExistWAUByWorkspaceIdAndUserId(wau.WorkspaceId, rand.Uint32()))
	}
}
