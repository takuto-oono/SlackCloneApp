package testdata

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/xyproto/randomstring"

	"backend/controllers"
	"backend/models"
)

func (td *TestData) createUsers() {
	signUpAndLogin := func(userName string) controllers.LoginResponse {
		controllers.SignUpTestFuncV2(userName, "abc123")
		_, lr := controllers.LoginTestFuncV2(userName, "abc123")
		return lr
	}

	// デフォルトで設定されていたデータが有れば、ユーザーを作成する
	func() {
		for i, lr := range td.users {
			if lr.Token == "" {
				td.users[i] = signUpAndLogin(lr.Username)
				if td.users[i].Token == "" {
					panic("failure signUp or login")
				}
			}
		}
	}()

	for i := 0; i < 1000; i++ {
		lr := signUpAndLogin(randomstring.EnglishFrequencyString(10))
		if lr.Token == "" {
			panic("failure signUp or login")
		}
		td.users = append(td.users, lr)
	}

	for _, lr := range td.users {
		td.jwtTokenMap[lr.UserId] = lr.Token
	}
}

func (td *TestData) createWorkspace() {
	errMes := "failure create workspace"
	if td.users[0].Token == "" {
		return
	}

	// デフォルトで設定したものを作成
	func() {
		for i, w := range td.workspaces {
			if w.ID == 0 {
				_, td.workspaces[i] = controllers.CreateWorkspaceTestFuncV2(w.Name, td.users[0].Token)
			}
			if td.workspaces[i].ID == 0 {
				panic(errMes)
			}
			td.workspaceAndUsers[td.workspaces[i].ID] = []models.WorkspaceAndUsers{
				models.WorkspaceAndUsers{
					UserId:      td.users[0].UserId,
					WorkspaceId: td.workspaces[i].ID,
					RoleId:      1,
				}}
		}
	}()

	for i := 0; i < 5; i++ {
		rr, w := controllers.CreateWorkspaceTestFuncV2(
			randomstring.EnglishFrequencyString(10),
			td.users[0].Token,
		)
		if rr.Code != http.StatusOK {
			panic("api err")
		}
		if w.ID == 0 {
			panic(errMes)
		}
		td.workspaceAndUsers[w.ID] = []models.WorkspaceAndUsers{
			models.WorkspaceAndUsers{
				UserId:      td.users[0].UserId,
				WorkspaceId: w.ID,
				RoleId:      1,
			}}
		td.workspaces = append(td.workspaces, w)
	}
}

func (td *TestData) addUserInWorkspace(workspaceID, userCnt int) {
	if !func(workspaceID int) bool {
		for _, w := range td.workspaces {
			if w.ID == workspaceID {
				return true
			}
		}
		return false
	}(workspaceID) {
		panic("no workspace")
	}

	for _, lr := range td.users {
		if len(td.workspaceAndUsers[workspaceID]) >= userCnt {
			break
		}
		if func(userID uint32) bool {
			for _, wau := range td.workspaceAndUsers[workspaceID] {
				if wau.UserId == userID {
					return false
				}
			}
			return true
		}(lr.UserId) {
			rr, wau := controllers.AddUserInWorkspaceV2(workspaceID, lr.UserId, 4, td.users[0].Token)
			if rr.Code != http.StatusOK {
				panic("api err")
			}
			td.workspaceAndUsers[workspaceID] = append(td.workspaceAndUsers[workspaceID], wau)
		}
	}
}

func (td *TestData) createChannels(workspaceID int) {
	if !func(workspaceID int) bool {
		if _, ok := td.workspaceAndUsers[workspaceID]; !ok {
			return false
		}
		for _, w := range td.workspaces {
			if w.ID == workspaceID {
				return true
			}
		}
		return false
	}(workspaceID) {
		panic("no workspace or no wau")
	}
	// channelは3個/人作ることにする
	// とりあえずはすべてpublicチャンネルにする
	isPrivate := false
	for _, wau := range td.workspaceAndUsers[workspaceID] {
		for i := 0; i < 7; i++ {
			rr, ch := controllers.CreateChannelTestFuncV2(randomstring.EnglishFrequencyString(10), "", &isPrivate, td.jwtTokenMap[wau.UserId], workspaceID)
			if rr.Code != http.StatusOK {
				fmt.Println("api err")
			}
			td.channels = append(td.channels, ch)
			cau := models.ChannelsAndUsers{
				ChannelId: ch.ID,
				UserId:    wau.UserId,
				IsAdmin:   true,
			}
			if _, ok := td.channelAndUsers[wau.UserId]; ok {
				td.channelAndUsers[wau.UserId] = append(td.channelAndUsers[wau.UserId], cau)
			} else {
				td.channelAndUsers[wau.UserId] = []models.ChannelsAndUsers{cau}
			}
		}
	}
}

func (td *TestData) addUserInChannel(workspaceID int) {
	for _, wau := range td.workspaceAndUsers[workspaceID] {
		for _, ch := range td.channelAndUsers[wau.UserId] {
			for i := 0; i < 20; i++ {
				lr := td.users[rand.Int()%len(td.workspaceAndUsers[workspaceID])]
				rr, cau := controllers.AddUserInChannelTestFuncV2(ch.ChannelId, lr.UserId, td.jwtTokenMap[wau.UserId])
				if rr.Code == http.StatusOK {
					td.channelAndUsers[cau.UserId] = append(td.channelAndUsers[cau.UserId], cau)
				} else if rr.Code == http.StatusForbidden || rr.Code == http.StatusConflict {
					break
				} else {
					panic("api err")
				}
			}
		}
	}
}

func (td *TestData) sendMessage(workspaceID int) {
	sendMessagePump := func(userID uint32) {
		for i := 0; i < 300; i++ {
			rr, m := controllers.SendMessageTestFuncV2(
				randomstring.EnglishFrequencyString(30),
				td.channelAndUsers[userID][rand.Int()%len(td.channelAndUsers[userID])].ChannelId,
				td.jwtTokenMap[userID],
				[]uint32{},
				time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			)
			if rr.Code == http.StatusOK {
				td.messages = append(td.messages, m)
			}
		}
	}

	var wg sync.WaitGroup
	for _, wau := range td.workspaceAndUsers[workspaceID] {
		wg.Add(1)
		go func(userID uint32) {
			defer wg.Done()
			sendMessagePump(userID)
		}(wau.UserId)
	}
	wg.Wait()
}

func (td *TestData) sendDM(workspaceID int) {
	sendDMPump := func(userID uint32) {
		for i := 0; i < 100; i++ {
			rr, m := controllers.SendDMTestFuncV2(
				randomstring.EnglishFrequencyString(30),
				td.jwtTokenMap[userID],
				td.workspaceAndUsers[workspaceID][rand.Int()%len(td.workspaceAndUsers[workspaceID])].UserId,
				workspaceID,
				[]uint32{},
				time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			)
			if rr.Code == http.StatusOK {
				td.dms = append(td.dms, m)
			}
		}
	}

	var wg sync.WaitGroup
	for _, wau := range td.workspaceAndUsers[workspaceID] {
		wg.Add(1)
		go func(userID uint32) {
			defer wg.Done()
			sendDMPump(userID)
		}(wau.UserId)
	}
	wg.Wait()
}
