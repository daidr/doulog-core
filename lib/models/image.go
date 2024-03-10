package models

type TImage struct {
	ID     uint64 `gorm:"primaryKey;autoIncrement;column:id;index"` // ID
	Url    string `gorm:"not null;column:url"`                      // 图片链接
	ETag   string `gorm:"not null;column:etag;unique"`              // 资源唯一标识，为文件名，七牛etag算法
	Src    int    `gorm:"type:smallint;not null;column:src"`        // 存储类型标识。1:本机，2:阿里云
	Width  int    `gorm:"type:smallint;not null;column:width"`      // 长
	Height int    `gorm:"type:smallint;not null;column:height"`     // 高
	MIME   string `gorm:"not null;column:mime"`                     // mime类型
	Author uint64 `gorm:"not null;column:author;index"`             // 上传者
	TimeHook
}

func (t *TImage) B() *BImage {
	return &BImage{
		Url:    t.Url,
		Width:  t.Width,
		Height: t.Height,
	}
}

type BImage struct {
	Url    string `json:"url"`    // 图片链接
	Width  int    `json:"width"`  // 长
	Height int    `json:"height"` // 高
}
