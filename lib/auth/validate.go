package auth

import (
	"github.com/daidr/doulog-core/lib/auth/internal/github"
	"github.com/daidr/doulog-core/lib/conf"
	"github.com/daidr/doulog-core/lib/models"
	"github.com/pkg/errors"
)

// Validate 返回第三方平台用户unique id
func Validate(platform string, code string, state string) (*models.OauthPayload, error) {
	switch platform {
	case conf.AuthGitHub:
		return github.Validate(code, state)
	}
	return nil, errors.New("platform not matches")
}
