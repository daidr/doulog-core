package me

import (
	"github.com/daidr/doulog-core/lib/daos"
	"github.com/daidr/doulog-core/lib/models"
	"github.com/daidr/doulog-core/module/x_user/internal/model"
)

func GetSelfInfo(db *models.DB, uid uint64) (*model.GetSelfInfoResp, error) {
	// 未登录返回空
	if uid == 0 {
		return &model.GetSelfInfoResp{
			Id:        0,
			Name:      "",
			Email:     "",
			EmailHash: "",
			IsAdmin:   false,
			IsBanned:  false,
		}, nil
	}

	u, err := daos.NewUser(db).GetB(uid)
	if err != nil {
		return nil, err
	}

	return &model.GetSelfInfoResp{
		Id:        u.Id,
		Name:      u.Name,
		Email:     u.Email,
		EmailHash: u.EmailHash,
		IsAdmin:   u.IsAdmin,
		IsBanned:  u.IsBanned,
	}, nil
}
