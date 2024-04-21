package service

import (
	"github.com/daidr/doulog-core/lib/daos"
	"github.com/daidr/doulog-core/lib/models"
	"github.com/daidr/doulog-core/module/x_user/internal/dto"
)

func GetUserInfo(db *models.DB, uid uint64) (*dto.GetUserInfoResp, error) {
	// 未登录返回空
	if uid == 0 {
		return &dto.GetUserInfoResp{
			Id:        0,
			Name:      "",
			Email:     "",
			EmailHash: "",
			Homepage:  "",
			Motto:     "",
			IsAdmin:   false,
			IsBanned:  false,
		}, nil
	}

	u, err := daos.NewUser(db).GetB(uid)
	if err != nil {
		return nil, err
	}

	return &dto.GetUserInfoResp{
		Id:        u.Id,
		Name:      u.Name,
		Email:     u.Email,
		EmailHash: u.EmailHash,
		Homepage:  u.Homepage,
		Motto:     u.Motto,
		IsAdmin:   u.IsAdmin,
		IsBanned:  u.IsBanned,
		CreatedAt: u.CreatedAt,
	}, nil
}
