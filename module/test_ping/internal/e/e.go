// Package e 模块内错误码定义，错误码应当全局唯一且大于0
// 程序启动时会自动检查，如果冲突则panic
package e

import "github.com/daidr/doulog-core/lib/ecode"

var (
	ErrPing        = ecode.New(100001)
	ErrRedisSet    = ecode.New(100002)
	ErrRedisGet    = ecode.New(100003)
	ErrHttpRequest = ecode.New(100004)
)

var ECode = map[ecode.Code]string{
	ErrPing:        "ping error",
	ErrRedisSet:    "failed to set redis value",
	ErrRedisGet:    "failed to get redis value",
	ErrHttpRequest: "failed to make http request",
}
