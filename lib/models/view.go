package models

// TArticleViews 文章浏览计数表，Redis会定时同步
type TArticleViews struct {
	ID      uint64 `gorm:"primaryKey;autoIncrement;column:id;index"`
	Article uint64 `gorm:"not null;column:article;index"`
	Views   uint64 `gorm:"not null;column:views"`
}
