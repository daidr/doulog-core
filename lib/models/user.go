package models

type TUser struct {
	Id       uint64 `gorm:"primaryKey;autoIncrement;column:id;index"` // ID
	Name     string `gorm:"not null;column:name;unique"`              // 昵称
	Email    string `gorm:"not null;column:email;unique"`             // 邮箱
	Homepage string `gorm:"not null;column:homepage"`                 // 个人主页
	IsAdmin  bool   `gorm:"type:boolean;not null;column:is_admin"`    // 是否是管理员
	IsBanned bool   `gorm:"type:boolean;column:is_banned"`            // 是否被封禁
	Attr     int    `gorm:"type:smallint;not null;column:attr"`       // 属性位
	TimeHook
}

type BUser struct {
	Id       uint64 `json:"id"`        // ID
	Name     string `json:"name"`      // 昵称
	Email    string `json:"email"`     // 邮箱
	Homepage string `json:"homepage"`  // 个人主页
	IsAdmin  bool   `json:"is_admin"`  // 是否是管理员
	IsBanned bool   `json:"is_banned"` // 是否被封禁
	Attr     int    `json:"attr"`      // 属性位
}
