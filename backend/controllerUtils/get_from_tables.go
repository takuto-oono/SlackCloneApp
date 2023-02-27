package controllerUtils

import "backend/models"

func GetWorkspacesByUserId(userId uint32) ([]models.Workspace, error) {
	// 引数で指定したuserIdのuserが所属しているworkspaceのstructを配列にして返す

	// 結果を保存する配列を作成
	workspaces := make([]models.Workspace, 0)

	// workspaces_and_users tableからuserIdが等しいものをすべて取得
	waus, err := models.GetWAUsByUserId(userId)
	if err != nil {
		return workspaces, err
	}

	// workspaces tableからidが等しいworkspaceの情報を取得する。
	for _, wau := range waus {
		workspace, err := models.GetWorkspaceById(wau.WorkspaceId)
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
	chs, err := models.GetChannelsByWorkspaceId(workspaceId)
	if err != nil {
		return res, err
	}

	// 指定されたuserが所属しているChannelAndUsersをすべて取得(他のworkspaceの物も含まれている)
	caus, err := models.GetCAUsByUserId(userId)
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
