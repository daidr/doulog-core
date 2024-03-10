package auth

import (
	"github.com/daidr/doulog-core/lib/utils"
)

func ReLoginWithMsg(callback string, msg string) string {
	return utils.UrlAppend(callback, map[string][]string{
		"msg": {msg},
	})
}

func ReLoginWithToken(callback string, token string) string {
	return utils.UrlAppend(callback, map[string][]string{
		"token": {token},
	})
}
