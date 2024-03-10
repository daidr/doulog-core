package github

import (
	"fmt"
	"github.com/daidr/doulog-core/lib/conf"
)

func Redirect(redirect string, state string) string {
	return fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&state=%s&scope=%s",
		conf.C.Auth.GitHub.ClientID, redirect, state, "read:user%20user:email")
}
