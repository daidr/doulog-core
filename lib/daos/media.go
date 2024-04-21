package daos

import (
	"errors"
	"github.com/daidr/doulog-core/lib/models"
)

type Media struct {
	db *models.DB
}

func NewMedia(db *models.DB) *Media {
	return &Media{db: db}
}

func (d *Media) Get(pid uint64) (*models.TMedia, error) {
	pic := models.TMedia{}
	if err := d.db.PgSQL.Preload("Owner").First(&pic, pid).Error; err != nil {
		return nil, err
	}
	return &pic, nil
}

func (d *Media) GetB(pid uint64) (*models.BMedia, error) {
	pic, err := d.Get(pid)
	if err != nil {
		return nil, err
	}

	return pic.B(), nil
}

func (d *Media) Add(pic *models.TMedia) error {
	return d.db.PgSQL.Create(pic).Error
}

func (d *Media) Exist(pids ...uint64) error {
	var num int64
	if err := d.db.PgSQL.Model(&models.TMedia{}).Where("id in (?)", pids).Count(&num).Error; err != nil {
		return err
	}
	if num != int64(len(pids)) {
		return errors.New("tag num is not matches")
	}
	return nil
}
