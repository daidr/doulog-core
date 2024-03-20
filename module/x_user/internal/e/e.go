// Package e 错误码定义
package e

import (
	"github.com/daidr/doulog-core/lib/ecode"
)

var (
	ErrGetUserInfo     = ecode.New(104001)
	ErrChangeNickName  = ecode.New(104002)
	ErrNickNameExisted = ecode.New(104003)
	ErrGetUserList     = ecode.New(104004)
)

var ECode = map[ecode.Code]string{
	ErrGetUserInfo:     "Failed to get user info",
	ErrChangeNickName:  "Failed to change nick name",
	ErrNickNameExisted: "Nick name existed",
	ErrGetUserList:     "Failed to get user list",
}
