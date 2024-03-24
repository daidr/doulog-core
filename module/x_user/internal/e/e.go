// Package e 错误码定义
package e

import (
	"github.com/daidr/doulog-core/lib/ecode"
)

var (
	ErrGetUserInfo    = ecode.New(104001)
	ErrGetUserList    = ecode.New(104002)
	ErrUpdateName     = ecode.New(104003)
	ErrNameExisted    = ecode.New(104004)
	ErrUpdateEmail    = ecode.New(104005)
	ErrEmailExisted   = ecode.New(104006)
	ErrUpdateHomepage = ecode.New(104007)
)

var ECode = map[ecode.Code]string{
	ErrGetUserInfo:    "Failed to get user info",
	ErrGetUserList:    "Failed to get user list",
	ErrUpdateName:     "Failed to update name",
	ErrNameExisted:    "Name existed",
	ErrUpdateEmail:    "Failed to update email",
	ErrEmailExisted:   "Email existed",
	ErrUpdateHomepage: "Failed to update homepage",
}
