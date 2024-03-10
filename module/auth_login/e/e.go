// Package e 错误码定义
package e

import (
	"github.com/daidr/doulog-core/lib/ecode"
)

var (
	ErrInvalidCallback = ecode.New(101001)
)

var ECode = map[ecode.Code]string{
	ErrInvalidCallback: "Invalid callback",
}
