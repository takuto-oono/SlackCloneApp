package models

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xyproto/randomstring"
)

func TestNewWorkspace(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	for i := 0; i < 10; i++ {
		name := randomstring.EnglishFrequencyString(30)
		primary_owner_id := rand.Uint32()
		w := NewWorkspace(name, primary_owner_id)
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
		w := NewWorkspace(names[i], primaryOwnerIds[i])
		err := w.Create(db)
		assert.Empty(t, err)
	}

	// workspace nameが既に存在する場合 error
	name := randomstring.EnglishFrequencyString(30)
	w := NewWorkspace(name, rand.Uint32())
	err := w.Create(db)
	assert.Empty(t, err)
	w2 := NewWorkspace(name, rand.Uint32())
	err = w2.Create(db)
	assert.NotEmpty(t, err)
}

func TestRenameWorkspaceName(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	w := NewWorkspace(randomstring.EnglishFrequencyString(30), rand.Uint32())
	w.Create(db)
	newName := randomstring.EnglishFrequencyString(30)
	_, err := UpdateWorkspaceName(db, w.ID, newName)
	assert.Empty(t, err)
	w2, err := GetWorkspaceById(db, w.ID)
	assert.Empty(t, err)
	assert.Equal(t, w2.ID, w.ID)
	assert.Equal(t, w2.Name, newName)
	assert.Equal(t, w2.PrimaryOwnerId, w.PrimaryOwnerId)
}
