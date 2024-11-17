package service

import (
	"github.com/daidr/doulog-core/lib/daos"
	"github.com/daidr/doulog-core/lib/models"
	"github.com/daidr/doulog-core/lib/utils"
	"github.com/daidr/doulog-core/module/auth_login/internal/dto"
	authUtils "github.com/daidr/doulog-core/module/auth_login/internal/utils"
	"github.com/pkg/errors"
)

func PwLogin(sp *models.Scope, req *dto.PwLoginReq) (string, error) {
	// check if email exists
	var user *models.TUser
	var err error
	if user, err = daos.NewUser(sp.DB).GetByEmail(req.Email); err != nil {
		return "", errors.New("email not exists")
	}

	// check password
	if !utils.CompareBcrypt(req.Password, user.Password) {
		return "", errors.New("password error")
	}

	token := authUtils.SetToken(sp.DB, user.ID)
	return token, nil
}
