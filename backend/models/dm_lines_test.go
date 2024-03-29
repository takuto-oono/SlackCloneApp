package models

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCreateDMLine(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	dl := NewDMLine(int(rand.Uint64()), rand.Uint32(), rand.Uint32())
	res := dl.Create(db)
	assert.NotEqual(t, 0, dl.ID)
	assert.Empty(t, res)
}

func TestGetDLsByUserIdAndWorkspaceId(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("1 データが存在する場合", func(t *testing.T) {
		userId := rand.Uint32()
		userId2s := make([]uint32, 10)
		workspaceId := int(rand.Uint32())
		for i := 0; i < 10; i++ {
			userId2s[i] = rand.Uint32()
		}
		for i := 0; i < 10; i++ {
			dl := NewDMLine(workspaceId, userId, userId2s[i])
			dl.Create(db)
		}
		res, err := GetDLsByUserIdAndWorkspaceId(db, userId, workspaceId)
		assert.Empty(t, err)
		assert.Equal(t, 10, len(res))

		for _, r := range res {
			assert.Equal(t, workspaceId, r.WorkspaceId)
			assert.True(t, r.UserId1 == userId || r.UserId2 == userId)
		}
	})

	t.Run("2 データが存在しない場合", func(t *testing.T) {
		res, err := GetDLsByUserIdAndWorkspaceId(db, rand.Uint32(), int(rand.Uint32()))
		assert.Empty(t, err)
		assert.Equal(t, 0, len(res))
	})
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
		res := dl.Create(db)
		assert.NotEqual(t, uint(0), dl.ID)
		assert.Empty(t, res)
		dm_line, err := GetDLByUserIdsAndWorkspaceId(db, userId1, userId2, workspaceId)
		assert.Empty(t, err)
		assert.Equal(t, dl.ID, dm_line.ID)
		dm_line, err = GetDLByUserIdsAndWorkspaceId(db, userId2, userId1, workspaceId)
		assert.Empty(t, err)
		assert.Equal(t, dl.ID, dm_line.ID)
	})

	t.Run("2 データが存在しない場合", func(t *testing.T) {
		_, err := GetDLByUserIdsAndWorkspaceId(db, rand.Uint32(), rand.Uint32(), int(rand.Uint64()))
		assert.NotEmpty(t, err)
	})
}

func TestGetDLById(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("1 データが存在する場合", func(t *testing.T) {
		dl := NewDMLine(rand.Int(), rand.Uint32(), rand.Uint32())
		assert.Empty(t, dl.Create(db))
		assert.NotEqual(t, uint(0), dl.ID)
		res, err := GetDLById(db, dl.ID)
		assert.Empty(t, err)
		assert.Equal(t, dl.ID, res.ID)
		assert.Equal(t, dl.WorkspaceId, res.WorkspaceId)
		assert.Equal(t, dl.UserId1, res.UserId1)
		assert.Equal(t, dl.UserId2, res.UserId2)
	})

	t.Run("2 データが存在しない場合", func(t *testing.T) {
		_, err := GetDLById(db, uint(rand.Uint64()))
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}
