package models

import (
	"github.com/lib/pq"
	"gorm.io/datatypes"
)

type TMedia struct {
	Id             uint64         `gorm:"primaryKey;autoIncrement;index"`                // Id
	Ext            string         `gorm:"not null"`                                      // 扩展名
	Etag           string         `gorm:"not null;unique"`                               // 资源唯一标识，为文件名，七牛etag算法
	Src            int            `gorm:"type:smallint;not null"`                        // 存储类型标识。1:本机
	Width          int            `gorm:"type:smallint;not null"`                        // 长
	Height         int            `gorm:"type:smallint;not null"`                        // 高
	MIME           string         `gorm:"not null"`                                      // mime类型
	OwnerID        uint64         `gorm:"not null;index"`                                // 上传者外键
	Owner          TUser          `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"` // 上传者
	Title          string         `gorm:"not null"`                                      // 标题
	Alt            string         `gorm:"not null"`                                      // alt
	FileSizes      datatypes.JSON `gorm:"type:jsonb;not null"`                           // 文件大小信息
	ProminentColor pq.StringArray `gorm:"type:text[];not null"`                          // 主题色
	TimeHook
}

func (t *TMedia) B() *BMedia {
	return &BMedia{
		Ext:            t.Ext,
		Etag:           t.Etag,
		Title:          t.Title,
		Alt:            t.Alt,
		MIME:           t.MIME,
		Owner:          *t.Owner.B(),
		Width:          t.Width,
		Height:         t.Height,
		FileSizes:      t.FileSizes,
		CreatedAt:      t.CreatedAt,
		ProminentColor: t.ProminentColor,
	}
}

type BMedia struct {
	Id             uint64         `json:"id"`                                 // Id
	Ext            string         `json:"ext"`                                // 扩展名
	Etag           string         `json:"etag"`                               // 资源唯一标识，为文件名，七牛etag算法
	Title          string         `json:"title"`                              // 标题
	MIME           string         `json:"mime"`                               // mime类型
	Owner          BUser          `json:"owner"`                              // 上传者
	Alt            string         `json:"alt"`                                // alt
	Width          int            `json:"width"`                              // 长
	Height         int            `json:"height"`                             // 高
	FileSizes      datatypes.JSON `json:"extra" gorm:"type:jsonb"`            // 文件大小信息
	CreatedAt      int64          `json:"createdAt"`                          // 创建时间
	ProminentColor pq.StringArray `json:"prominent_color" gorm:"type:text[]"` // 主题色
}

type MediaFileSizes struct {
	Original  uint64 `json:"original"`
	Thumbnail uint64 `json:"thumbnail"`
}
