package testdata

import (
	"backend/controllers"
	"backend/models"
)

type TestData struct {
	users             []controllers.LoginResponse
	workspaces        []models.Workspace
	channels          []models.Channel
	messages          []models.Message
	workspaceAndUsers map[int]([]models.WorkspaceAndUsers)
	jwtTokenMap       map[uint32]string
	channelAndUsers   map[uint32]([]models.ChannelsAndUsers)
	dms               []models.Message
}

func NewTestData() *TestData {
	td := &TestData{
		users:             make([]controllers.LoginResponse, 0),
		workspaces:        make([]models.Workspace, 0),
		channels:          make([]models.Channel, 0),
		messages:          make([]models.Message, 0),
		workspaceAndUsers: map[int][]models.WorkspaceAndUsers{},
		jwtTokenMap:       map[uint32]string{},
		channelAndUsers:   map[uint32][]models.ChannelsAndUsers{},
		dms:               make([]models.Message, 0),
	}

	setDefaultUsers := func(userNames []string) {
		for _, name := range userNames {
			td.users = append(td.users, controllers.LoginResponse{
				Username: name,
				Token:    "",
				UserId:   uint32(0),
			})
		}
	}

	setDefaultWorkspaces := func(workspaceNames []string) {
		for _, name := range workspaceNames {
			td.workspaces = append(td.workspaces, models.Workspace{
				Name: name,
			})
		}
	}
	setDefaultUsers([]string{"user1", "user2", "user3"})
	setDefaultWorkspaces([]string{"testWorkspace"})
	return td
}
