package models

type TOauth struct {
	ID       uint64 `gorm:"primaryKey;autoIncrement;column:id;index"`
	User     uint64 `gorm:"not null;column:user;index"`    // 对应User表的ID
	Platform string `gorm:"not null;column:platform"`      // 平台类型。定义在常量中
	OpenID   string `gorm:"not null;column:open_id;index"` // 第三方平台唯一id
	TimeHook
}
