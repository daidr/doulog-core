package models

type TReply struct {
	ID  uint64 `gorm:"primaryKey;autoIncrement;column:id;index"` // ID
	AID uint64 `gorm:"not null;column:aid;index"`                // 评论区 article id
	// 根评论rpid
	// 若为一级评论则为0
	// 大于一级评论则为根评论id
	Root uint64 `gorm:"not null;column:root;index"`
	// 回复父评论rpid
	// 若为一级评论则为0
	// 若为二级评论则为根评论rpid
	// 大于二级评论为上一级评论rpid
	Parent  uint64 `gorm:"not null;column:parent;index"`
	Attr    int    `gorm:"type:smallint;not null;column:attr"` // 属性位
	Author  uint64 `gorm:"not null;column:author;index"`       // 发送者uid
	Content string `gorm:"not null;column:content"`            // 评论内容
	Hide    bool   `gorm:"type:boolean;not null;column:hide"`  // 违规评论被折叠。false:不折叠，true:折叠
	TimeHook
}

type BReply struct {
	ID uint64 `json:"rpid"` // ID
	// 根评论rpid
	// 若为一级评论则为0
	// 大于一级评论则为根评论id
	Root uint64 `json:"root"`
	// 回复父评论rpid
	// 若为一级评论则为0
	// 若为二级评论则为根评论rpid
	// 大于二级评论为上一级评论rpid
	Parent  uint64 `json:"parent"`
	Attr    int    `json:"attr"`    // 属性位
	Author  *BUser `json:"author"`  // 发送者uid
	Content string `json:"content"` // 评论内容
	Ctime   int64  `json:"ctime"`
	Hide    bool   `json:"hide"` // 违规评论被折叠。false:不折叠，true:折叠
}

type GeneralReply struct {
	*BReply
	Counter *ReplyCounter `json:"counter"`
}

type ReplyCounter struct {
	Reply int64 `json:"reply"`
}
