package controllerUtils

import (
	"backend/models"
)

type UserInfoInWorkspace struct {
	ID     uint32 `json:"id"`
	Name   string `json:"name"`
	RoleId int    `json:"role_id"`
}

func GetWorkspacesByUserId(userId uint32) ([]models.Workspace, error) {
	// 引数で指定したuserIdのuserが所属しているworkspaceのstructを配列にして返す

	// 結果を保存する配列を作成
	workspaces := make([]models.Workspace, 0)

	// workspaces_and_users tableからuserIdが等しいものをすべて取得
	waus, err := models.GetWAUsByUserId(db, userId)
	if err != nil {
		return workspaces, err
	}

	// workspaces tableからidが等しいworkspaceの情報を取得する。
	for _, wau := range waus {
		workspace, err := models.GetWorkspaceById(db, wau.WorkspaceId)
		if err != nil {
			return workspaces, err
		}
		workspaces = append(workspaces, workspace)
	}
	return workspaces, err
}

func GetChannelsByUserIdAndWorkspaceId(userId uint32, workspaceId int) ([]models.Channel, error) {
	// 指定されたworkspaceの中からuserが所属しているchannelのstructを配列にして返す

	res := make([]models.Channel, 0)

	// 指定されたworkspaceに存在するすべてのchannelを取得(userが所属していないchannelも含まれている)
	chs, err := models.GetChannelsByWorkspaceId(db, workspaceId)
	if err != nil {
		return res, err
	}

	// 指定されたuserが所属しているChannelAndUsersをすべて取得(他のworkspaceの物も含まれている)
	caus, err := models.GetCAUsByUserId(db, userId)
	if err != nil {
		return res, err
	}

	// workspaceに存在するchannelとuserが所属しているchannelで同じものがあればスライスに追加する
	for _, ch := range chs {
		for _, cau := range caus {
			if ch.ID == cau.ChannelId {
				res = append(res, ch)
				break
			}
		}
	}
	return res, nil
}

func GetUserInWorkspace(workspaceId int) ([]UserInfoInWorkspace, error) {
	// workspace内のuserの情報を配列にして返す
	// userぞれぞれの情報は返り値のstructを参照
	// アクセスしたuserの情報も含めて返す

	res := make([]UserInfoInWorkspace, 0)

	// workspaces_and_users tableからworkspace_idが等しいものをすべて取得する
	waus, err := models.GetWAUsByWorkspaceId(db, workspaceId)
	if err != nil {
		return res, err
	}

	// users tableの全情報を取得する
	users, err := models.GetUsers(db)
	if err != nil {
		return res, err
	}

	// 2つのデータからuser_idが等しいものの組み合わせを見つけて、res配列に追加する
	for _, wau := range waus {
		for _, user := range users {
			if user.ID == wau.UserId {
				res = append(res, UserInfoInWorkspace{
					ID:     user.ID,
					Name:   user.Name,
					RoleId: wau.RoleId,
				})
				break
			}
		}
	}
	return res, nil
}

func GetThreadsByUserSortedByEditedTime(userId uint32) ([]models.Thread, error) {
	var ths []models.Thread
	// userが所属しているthreadの情報を取得する
	taus, err := models.GetTAUsByUserId(db, userId)
	if err != nil {
		return ths, err
	}

	// thread tableから情報を取得する
	for _, tau := range taus {
		th, err := models.GetThreadById(db, tau.ThreadId)
		if err != nil {
			return ths, err
		}
		ths = append(ths, th)
	}

	// 更新時間でソート
	for i := 0; i < len(ths); i++ {
		for j := i + 1; j < len(ths); j++ {
			if ths[i].UpdatedAt.Before(ths[j].UpdatedAt) {
				ths[i], ths[j] = ths[j], ths[i]
			}
		}
	}

	return ths, nil
}
