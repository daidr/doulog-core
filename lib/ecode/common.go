// Package ecode 公共错误码
package ecode

var (
	Success          = conv(0)
	UnknownError     = conv(-1)
	InvalidParams    = conv(-2)
	RateLimit        = conv(-3)
	Unauthorized     = conv(-4)
	PermissionDenied = conv(-5)
)

func init() {
	Register(map[Code]string{
		Success:          "Ok",
		UnknownError:     "Unknown error",
		InvalidParams:    "Invalid parameters",
		RateLimit:        "Rate limit exceeded",
		Unauthorized:     "Unauthorized",
		PermissionDenied: "Permission denied",
	})
}
