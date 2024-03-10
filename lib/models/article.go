package models

// 文章(article)分为两种，一种是普通文章(post)，一种是页面(page)，两者共用id

type TArticle struct {
	ID          uint64 `gorm:"primaryKey;autoIncrement;column:id;index"` // id
	Slug        string `gorm:"not null;unique;column:slug"`              // slug
	Type        int    `gorm:"type:smallint;not null;column:type"`       // 类型，0为post，1为page
	Title       string `gorm:"not null;column:title"`                    // 标题
	Hero        uint64 `gorm:"column:hero"`                              // 封面图 image id
	Summary     string `gorm:"not null;column:summary"`                  // 简介
	Content     string `gorm:"not null;column:content"`                  // 内容
	Author      uint64 `gorm:"not null;column:author;index"`             // 发布者 uid
	Attr        int    `gorm:"type:smallint;not null;column:attr"`       // 属性位
	isDraft     bool   `gorm:"type:boolean;not null;column:is_draft"`    // 是否是草稿
	PublishedAt int64  `gorm:"not null;column:published_at"`             // 发布时间

	TimeHook
}

// TArticleTag tag绑定到article
type TArticleTag struct {
	ID      uint64 `gorm:"primaryKey;autoIncrement;column:id;index"` // id
	Article uint64 `gorm:"not null;column:article;index"`            // article id
	Tag     uint64 `gorm:"not null;column:tag;index"`                // tid
	TimeHook
}

type BArticle struct {
	ID      uint64  `json:"id"`
	Slug    string  `json:"slug"`
	Title   string  `json:"title"`
	Hero    *BImage `json:"hero,omitempty"`
	Content string  `json:"content"`
	Attr    int     `json:"attr"`
	Author  *BUser  `json:"author,omitempty"`
	Ctime   int64   `json:"ctime"`
	Mtime   int64   `json:"mtime"`
	IsDraft bool    `json:"is_draft"`
	PTime   int64   `json:"ptime"`
}

type GeneralArticle struct {
	*BArticle
	Stat    *ArticleStat    `json:"stat,omitempty"`
	Counter *ArticleCounter `json:"counter,omitempty"`
	Tags    []*BTag         `json:"tags,omitempty"`
}

type ArticleStat struct {
	Liked int `json:"liked"`
}

type ArticleCounter struct {
	View  uint64 `json:"view"`
	Like  uint64 `json:"like"`
	Reply int64  `json:"reply"`
}
