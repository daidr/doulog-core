package service

import (
	"github.com/daidr/doulog-core/lib/daos"
	"github.com/daidr/doulog-core/lib/models"
)

func GetMediaMetadata(db *models.DB, pid uint64) (*models.BMedia, error) {
	media, err := daos.NewMedia(db).GetB(pid)
	if err != nil {
		return nil, err
	}

	return media, nil
}
