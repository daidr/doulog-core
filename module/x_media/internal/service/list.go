package service

import (
	"github.com/daidr/doulog-core/lib/models"
	"github.com/daidr/doulog-core/module/x_media/internal/dto"
)

func ListAllMedia(db *models.DB, req *dto.MediaListReq) (*dto.MediaListResp, error) {
	var resp dto.MediaListResp
	var media []dto.MediaMetadataItemWithoutOwner
	var total int64
	var err error

	if req.Keyword != "" {
		if err = db.PgSQL.Model(&models.TMedia{}).Where("title like ? or alt like ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%").Count(&total).Limit(req.PageSize).Offset((req.Page - 1) * req.PageSize).Order("id desc").Find(&media).Error; err != nil {
			return nil, err
		}
	} else {
		if err = db.PgSQL.Model(&models.TMedia{}).Count(&total).Limit(req.PageSize).Offset((req.Page - 1) * req.PageSize).Order("id desc").Find(&media).Error; err != nil {
			return nil, err
		}
	}

	resp.Total = total
	resp.List = media

	return &resp, nil
}
