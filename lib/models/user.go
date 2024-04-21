package models

type TUser struct {
	Id        uint64 `gorm:"primaryKey;autoIncrement;column:id;index"` // Id
	Name      string `gorm:"not null;column:name;unique"`              // 昵称
	Email     string `gorm:"not null;column:email;unique"`             // 邮箱
	EmailHash string `gorm:"not null;column:email_hash;unique"`        // 邮箱 Hash
	Homepage  string `gorm:"not null;column:homepage"`                 // 个人主页
	Motto     string `gorm:"column:motto"`                             // 座右铭
	IsAdmin   bool   `gorm:"type:boolean;not null;column:is_admin"`    // 是否是管理员
	IsBanned  bool   `gorm:"type:boolean;column:is_banned"`            // 是否被封禁
	Attr      int    `gorm:"type:smallint;not null;column:attr"`       // 属性位
	TimeHook
}

func (t *TUser) B() *BUser {
	return &BUser{
		Id:        t.Id,
		Name:      t.Name,
		Email:     t.Email,
		EmailHash: t.EmailHash,
		Homepage:  t.Homepage,
		Motto:     t.Motto,
		IsAdmin:   t.IsAdmin,
		IsBanned:  t.IsBanned,
		CreatedAt: t.CreatedAt,
		Attr:      t.Attr,
	}
}

type BUser struct {
	Id        uint64 `json:"id"`        // Id
	Name      string `json:"name"`      // 昵称
	Email     string `json:"email"`     // 邮箱
	EmailHash string `json:"emailHash"` // 邮箱 Hash
	Homepage  string `json:"homepage"`  // 个人主页
	Motto     string `json:"motto"`     // 座右铭
	IsAdmin   bool   `json:"isAdmin"`   // 是否是管理员
	IsBanned  bool   `json:"isBanned"`  // 是否被封禁
	CreatedAt int64  `json:"createdAt"` // 创建时间
	Attr      int    `json:"attr"`      // 属性位
}
