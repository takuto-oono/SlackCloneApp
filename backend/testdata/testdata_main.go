package testdata

import (
	"backend/models"
)

func TestDataMain() {
	models.DeleteAllTableData()
	td := NewTestData()
	td.createUsers()
	td.createWorkspace()
	td.addUserInWorkspace(td.workspaces[0].ID, 100)
	td.createChannels(td.workspaces[0].ID)
	td.addUserInChannel(td.workspaces[0].ID)
	td.sendMessage(td.workspaces[0].ID)
	td.sendDM(td.workspaces[0].ID)
}
