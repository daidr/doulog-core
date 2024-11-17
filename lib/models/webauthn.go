package models

import (
	"github.com/go-webauthn/webauthn/webauthn"
	"gorm.io/datatypes"
)

type TWebAuthnCredential struct {
	ID         uint64                                  `gorm:"primaryKey;autoIncrement;column:id;index"`                        // ID
	UserID     uint64                                  `gorm:"column:user_id;index"`                                            // 所属用户外键
	User       TUser                                   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"` // 所属用户
	Label      string                                  `gorm:"not null;column:label"`                                           // 标签
	LastUsedAt int64                                   `gorm:"column:last_used_at"`                                             // 最后使用时间
	Credential datatypes.JSONType[webauthn.Credential] `gorm:"column:credential"`                                               // 凭证
	TimeHook
}
