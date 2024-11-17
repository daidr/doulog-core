package service

import (
	"github.com/daidr/doulog-core/lib/daos"
	"github.com/daidr/doulog-core/lib/models"
	"github.com/daidr/doulog-core/module/auth_login/internal/dto"
	"github.com/pkg/errors"
)

func Register(sp *models.Scope, req *dto.RegisterReq) error {
	var exists bool
	var err error
	if exists, err = daos.NewUser(sp.DB).CheckEmailExists(req.Email); err != nil {
		return err
	}
	if exists {
		return errors.New("email already exists")
	}

	if err := daos.NewUser(sp.DB).CreateUser(req.Name, req.Email, req.Password, false); err != nil {
		return err
	}

	return nil
}
