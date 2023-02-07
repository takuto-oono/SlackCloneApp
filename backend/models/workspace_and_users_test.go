package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWorkspaceAndUsers(t *testing.T) {
	for i := 0; i < 10; i++ {
		workspaceId := i + 1
		for j := 0; j < 100; j++ {
			userId := uint32(j + 1)
			roleId := j%4 + 1
			wau := NewWorkspaceAndUsers(workspaceId, uint32(userId), roleId)
			assert.Equal(t, workspaceId, wau.WorkspaceId)
			assert.Equal(t, userId, wau.UserId)
			assert.Equal(t, roleId, wau.RoleId)
		}
	}
}

func TestCreate(t *testing.T) {
	for i := 0; i < 10; i++ {
		workspaceId := i + 1
		for j := 0; j < 100; j++ {
			userId := uint32(j + 1)
			roleId := j%4 + 1
			wau := NewWorkspaceAndUsers(workspaceId, uint32(userId), roleId)
			assert.Empty(t, wau.Create())
			assert.NotEmpty(t, wau.Create())
			wau = NewWorkspaceAndUsers(workspaceId, uint32(userId), (roleId+1)%4)
			assert.NotEmpty(t, wau.Create())
		}
	}
}

func TestGetWorkspaceAndUserByWorkspaceIdAndUserId(t *testing.T) {
	for i := 0; i < 10; i++ {
		workspaceId := i + 1 + 100
		for j := 0; j < 100; j++ {
			userId := uint32(j + 1)
			roleId := j%4 + 1
			wau := NewWorkspaceAndUsers(workspaceId, userId, roleId)
			wau.Create()
			getWau, err := GetWorkspaceAndUserByWorkspaceIdAndUserId(workspaceId, userId)
			assert.Empty(t, err)
			assert.Equal(t, workspaceId, getWau.WorkspaceId)
			assert.Equal(t, userId, getWau.UserId)
			assert.Equal(t, roleId, getWau.RoleId)

			_, err = GetWorkspaceAndUserByWorkspaceIdAndUserId(workspaceId+1000, userId)
			assert.NotEmpty(t, err)
		}
	}
}

func TestDeleteWorkspaceAndUser(t *testing.T) {
	for i := 0; i < 10; i++ {
		workspaceId := i + 200
		for j := 0; j < 100; j++ {
			userId := uint32(j + 2)
			roleId := j%4 + 1
			wau := NewWorkspaceAndUsers(workspaceId, userId, roleId)
			assert.Empty(t, wau.Create())
			_, err := GetWorkspaceAndUserByWorkspaceIdAndUserId(wau.WorkspaceId, wau.UserId)
			assert.Empty(t, err)
			err = wau.DeleteWorkspaceAndUser()
			assert.Empty(t, err)
			_, err = GetWorkspaceAndUserByWorkspaceIdAndUserId(wau.WorkspaceId, wau.UserId)
			assert.NotEmpty(t, err)
			err = wau.DeleteWorkspaceAndUser()
			assert.Empty(t, err)
		}
	}
}
