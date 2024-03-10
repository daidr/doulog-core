package auth

import (
	"fmt"
	"net/url"

	"github.com/daidr/doulog-core/lib/auth/internal/github"
	"github.com/daidr/doulog-core/lib/conf"
	"github.com/daidr/doulog-core/lib/utils"
)

func GetRedirectUrl(platform string) (string, string) {
	mark := utils.RandString(10)
	redirect := url.QueryEscape(conf.C.Server.Site + "/api/auth/login/v1/callback")
	state := url.QueryEscape(fmt.Sprintf("%s_%s", platform, mark))

	switch platform {
	case conf.AuthGitHub:
		return github.Redirect(redirect, state), mark
	}
	return "", ""
}
