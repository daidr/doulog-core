package models

import "gorm.io/plugin/soft_delete"

type TimeHook struct {
	CreatedAt int64 `gorm:"not null;column:created_at" json:"-"`
	UpdatedAt int64 `gorm:"not null;column:updated_at" json:"-"`
	DeletedAt soft_delete.DeletedAt
}

type PageDto struct {
	Page     int `form:"page" json:"page"`
	PageSize int `form:"page_size" json:"page_size"`
}
