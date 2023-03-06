package models

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xyproto/randomstring"

	"backend/config"
)

func TestNewWorkspace(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	for i := 0; i < 10; i++ {
		id := rand.Int()
		name := randomstring.EnglishFrequencyString(30)
		primary_owner_id := rand.Uint32()
		w := NewWorkspace(id, name, primary_owner_id)
		assert.Equal(t, id, w.ID)
		assert.Equal(t, name, w.Name)
		assert.Equal(t, primary_owner_id, w.PrimaryOwnerId)
	}
}

func TestCreateWorkspace(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	// 正常な場合
	numbersOfTests := 1
	names := make([]string, numbersOfTests)
	primaryOwnerIds := make([]uint32, numbersOfTests)
	for i := 0; i < numbersOfTests; i++ {
		names[i] = randomstring.EnglishFrequencyString(30)
		primaryOwnerIds[i] = rand.Uint32()
	}

	for i := 0; i < numbersOfTests; i++ {
		w := NewWorkspace(0, names[i], primaryOwnerIds[i])
		err := w.Create()
		assert.Empty(t, err)
	}

	cmd := fmt.Sprintf("SELECT id, name, workspace_primary_owner_id FROM %s WHERE name = $1", config.Config.WorkspaceTableName)
	for i := 0; i < numbersOfTests; i++ {
		row := DbConnection.QueryRow(cmd, names[i])
		var w Workspace
		err := row.Scan(&w.ID, &w.Name, &w.PrimaryOwnerId)
		assert.Empty(t, err)
		assert.NotEqual(t, 0, w.ID)
		assert.Equal(t, names[i], w.Name)
		assert.Equal(t, primaryOwnerIds[i], w.PrimaryOwnerId)
	}

	// workspace nameが既に存在する場合 error
	name := randomstring.EnglishFrequencyString(30)
	w := NewWorkspace(0, name, rand.Uint32())
	err := w.Create()
	assert.Empty(t, err)
	w2 := NewWorkspace(0, name, rand.Uint32())
	err = w2.Create()
	assert.NotEmpty(t, err)
}

func TestRenameWorkspaceName(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	w := NewWorkspace(rand.Int(), randomstring.EnglishFrequencyString(30), 3)
	w.Create()
	w.Name = randomstring.EnglishFrequencyString(30)
	err := w.RenameWorkspaceName()
	assert.Empty(t, err)
	w2, err := GetWorkspaceById(w.ID)
	assert.Empty(t, err)
	assert.Equal(t, w2.ID, w.ID)
	assert.Equal(t, w2.Name, w.Name)
	assert.Equal(t, w2.PrimaryOwnerId, w.PrimaryOwnerId)
}

