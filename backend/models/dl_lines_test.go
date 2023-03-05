package models

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateDMLine(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	dl := NewDMLine(int(rand.Uint64()), rand.Uint32(), rand.Uint32())
	res := dl.Create()
	assert.NotEqual(t, 0, dl.ID)
	assert.Empty(t, res.Error)
}

func TestGetByUserIds(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("1 データが存在する場合", func(t *testing.T) {
		userId1 := rand.Uint32()
		userId2 := rand.Uint32()
		workspaceId := int(rand.Uint64())
		dl := NewDMLine(workspaceId, userId1, userId2)
		res := dl.Create()
		assert.NotEqual(t, 0, dl.ID)
		assert.Empty(t, res.Error)
		dm_line, err := GetDLByUserIdsAndWorkspaceId(userId1, userId2, workspaceId)
		assert.Empty(t, err)
		assert.Equal(t, dl.ID, dm_line.ID)
		dm_line, err = GetDLByUserIdsAndWorkspaceId(userId2, userId1, workspaceId)
		assert.Empty(t, err)
		assert.Equal(t, dl.ID, dm_line.ID)
	})

	t.Run("2 データが存在しない場合", func(t *testing.T) {
		_, err := GetDLByUserIdsAndWorkspaceId(rand.Uint32(), rand.Uint32(), int(rand.Uint64()))
		assert.NotEmpty(t, err)
	})
}
