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

// func (td *TestData) createUserData() {
// 	for i := 0; i < 1000; i ++ {
// 		td.users = append(td.users, controllers.LoginResponse{
// 			Username: randomstring.EnglishFrequencyString(30),
// 			Token: "",
// 			UserId: uint32(0),
// 		})
// 	}
// }

// func (td *TestData) createChannelData() {
// 	for i := 0; i < 10000; i ++ {
// 		if i % 5 == 0 {
// 			td.channels = append(td.channels, *&models.Channel{
// 				Name: randomstring.EnglishFrequencyString(20),
// 				Description: randomstring.EnglishFrequencyString(30),
// 				IsPrivate: true,
// 				IsArchive: true,
// 			})
// 		} else {
// 			td.channels = append(td.channels, *&models.Channel{
// 				Name: randomstring.EnglishFrequencyString(20),
// 				Description: randomstring.EnglishFrequencyString(30),
// 				IsPrivate: false,
// 				IsArchive: true,
// 			})
// 		}
// 	}
// }
