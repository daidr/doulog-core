// Package e 错误码定义
package e

import (
	"github.com/daidr/doulog-core/lib/ecode"
)

var (
	ErrGetSelfInfo     = ecode.New(104001)
	ErrChangeNickName  = ecode.New(104002)
	ErrNickNameExisted = ecode.New(104003)
)

var ECode = map[ecode.Code]string{
	ErrGetSelfInfo:     "Failed to get user info",
	ErrChangeNickName:  "Failed to change nick name",
	ErrNickNameExisted: "Nick name existed",
}
