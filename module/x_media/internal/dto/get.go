package dto

import (
	"github.com/daidr/doulog-core/lib/models"
	"github.com/lib/pq"
	"gorm.io/datatypes"
)

type MediaListReq struct {
	Keyword string `form:"keyword" json:"keyword"`

	models.PageDto
}

type MediaMetadataItemWithoutOwner struct {
	Id             uint64         `json:"id"`
	Ext            string         `json:"ext"`
	Etag           string         `json:"etag"`
	Title          string         `json:"title"`
	Alt            string         `json:"alt"`
	MIME           string         `json:"mime"`
	Width          int            `json:"width"`
	Height         int            `json:"height"`
	FileSizes      datatypes.JSON `json:"extra" gorm:"type:jsonb"`
	ProminentColor pq.StringArray `json:"prominentColor" gorm:"type:text[];"`
	CreatedAt      int64          `json:"createdAt"`
}

type MediaListResp struct {
	Total int64                           `json:"total"`
	List  []MediaMetadataItemWithoutOwner `json:"list"`
}
