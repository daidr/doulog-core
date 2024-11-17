package models

import (
	"encoding/binary"
	"github.com/go-webauthn/webauthn/webauthn"
)

type TUser struct {
	ID          uint64                `gorm:"primaryKey;autoIncrement;column:id;index"`                        // ID
	Name        string                `gorm:"not null;column:name"`                                            // 昵称
	Email       string                `gorm:"not null;column:email;unique"`                                    // 邮箱
	EmailHash   string                `gorm:"not null;column:email_hash;unique"`                               // 邮箱 Hash
	Password    string                `gorm:"not null;column:password"`                                        // 密码
	Homepage    string                `gorm:"column:homepage"`                                                 // 个人主页
	Motto       string                `gorm:"column:motto"`                                                    // 座右铭
	IsAdmin     bool                  `gorm:"type:boolean;not null;column:is_admin"`                           // 是否是管理员
	IsBanned    bool                  `gorm:"type:boolean;column:is_banned"`                                   // 是否被封禁
	Attr        int                   `gorm:"type:smallint;not null;column:attr"`                              // 属性位
	Credentials []TWebAuthnCredential `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"` // WebAuthn 凭证
	TimeHook
}

func uint64ToWebAuthnID(id uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, id)
	return b
}

func (t *TUser) WebAuthnID() []byte {
	return uint64ToWebAuthnID(t.ID)
}

func (t *TUser) WebAuthnName() string {
	return t.Name
}

func (t *TUser) WebAuthnDisplayName() string {
	return t.Name
}

func (t *TUser) WebAuthnCredentials() []webauthn.Credential {
	credentials := make([]webauthn.Credential, len(t.Credentials))
	for i, c := range t.Credentials {
		credentials[i] = c.Credential.Data()
	}
	return credentials
}

func (t *TUser) B() *BUser {
	return &BUser{
		Id:        t.ID,
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
