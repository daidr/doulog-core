package service

import (
	"github.com/daidr/doulog-core/lib/daos"
	"github.com/daidr/doulog-core/lib/models"
	"github.com/daidr/doulog-core/module/x_user/internal/model"
)

func GetUserInfo(db *models.DB, uid uint64) (*model.GetUserInfoResp, error) {
	// 未登录返回空
	if uid == 0 {
		return &model.GetUserInfoResp{
			Id:        0,
			Name:      "",
			Email:     "",
			EmailHash: "",
			Homepage:  "",
			IsAdmin:   false,
			IsBanned:  false,
		}, nil
	}

	u, err := daos.NewUser(db).GetB(uid)
	if err != nil {
		return nil, err
	}

	return &model.GetUserInfoResp{
		Id:        u.Id,
		Name:      u.Name,
		Email:     u.Email,
		EmailHash: u.EmailHash,
		Homepage:  u.Homepage,
		IsAdmin:   u.IsAdmin,
		IsBanned:  u.IsBanned,
		CreatedAt: u.CreatedAt,
	}, nil
}
