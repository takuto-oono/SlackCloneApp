package models

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"backend/config"
)

func TestNewWorkspace(t *testing.T) {
	for i := 0; i < 1000; i++ {
		id := i
		name := "testNewWorkspaceName" + strconv.Itoa(i)
		primary_owner_id := uint32(i) + 10
		w := NewWorkspace(id, name, primary_owner_id)
		assert.Equal(t, w.ID, id)
		assert.Equal(t, w.Name, name)
		assert.Equal(t, w.PrimaryOwnerId, primary_owner_id)
	}
}

func TestCreateWorkspace(t *testing.T) {
	// 正常な場合
	numbersOfTests := 1000
	names := make([]string, numbersOfTests)
	primaryOwnerIds := make([]uint32, numbersOfTests)
	for i := 0; i < numbersOfTests; i++ {
		names[i] = "testCreateWorkspace" + strconv.Itoa(i)
		primaryOwnerIds[i] = uint32(i) + 11
	}

	for i := 0; i < numbersOfTests; i++ {
		w := NewWorkspace(0, names[i], primaryOwnerIds[i])
		err := w.CreateWorkspace()
		assert.Empty(t, err)
	}

	cmd := fmt.Sprintf("SELECT id, name, workspace_primary_owner_id FROM %s WHERE name = ?", config.Config.WorkspaceTableName)
	for i := 0; i < numbersOfTests; i++ {
		row := DbConnection.QueryRow(cmd, names[i])
		var w Workspace
		err := row.Scan(&w.ID, &w.Name, &w.PrimaryOwnerId)
		assert.Empty(t, err)
		assert.NotEqual(t, w.ID, 0)
		assert.Equal(t, w.Name, names[i])
		assert.Equal(t, w.PrimaryOwnerId, primaryOwnerIds[i])
	}

	// workspace nameが既に存在する場合 error
	name := "testCreateWorkspaceDuplicate"
	w := NewWorkspace(0, name, uint32(2))
	err := w.CreateWorkspace()
	assert.Empty(t, err)
	w2 := NewWorkspace(0, name, uint32(1))
	err = w2.CreateWorkspace()
	assert.NotEmpty(t, err)
}
