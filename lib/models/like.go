package models

// TArticleLike 文章点赞数表
type TArticleLike struct {
	ID      uint64 `gorm:"primaryKey;autoIncrement;column:id;index"` // ID
	Article uint64 `gorm:"not null;column:article;index"`            // article id
	Likes   uint64 `gorm:"not null;column:likes;index"`              // likes
}
