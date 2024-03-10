// Package ecode 公共错误码
package ecode

var (
	Success       = conv(0)
	UnknownError  = conv(-1)
	InvalidParams = conv(-2)
	RateLimit     = conv(-3)
	Unauthorized  = conv(-4)
)

func init() {
	Register(map[Code]string{
		Success:       "ok",
		UnknownError:  "未知错误",
		InvalidParams: "参数错误或存在非法字符",
		RateLimit:     "访问频率过高",
		Unauthorized:  "身份验证失败",
	})
}
