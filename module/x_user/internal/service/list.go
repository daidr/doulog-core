package service

import (
	"github.com/daidr/doulog-core/lib/models"
	"github.com/daidr/doulog-core/module/x_user/internal/dto"
)

func ListAllUsers(db *models.DB, req *dto.UserListReq) (*dto.UserListResp, error) {
	var resp dto.UserListResp
	var users []dto.GetUserInfoResp
	var total int64
	var err error

	if req.Keyword != "" {
		if err = db.PgSQL.Model(&models.TUser{}).Where("name like ? or email like ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%").Count(&total).Limit(req.PageSize).Offset((req.Page - 1) * req.PageSize).Order("id asc").Find(&users).Error; err != nil {
			return nil, err
		}
	} else {
		if err = db.PgSQL.Model(&models.TUser{}).Count(&total).Limit(req.PageSize).Offset((req.Page - 1) * req.PageSize).Order("id asc").Find(&users).Error; err != nil {
			return nil, err
		}
	}

	resp.Total = total
	resp.List = users

	return &resp, nil
}
