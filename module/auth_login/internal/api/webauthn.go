package api

import (
	"github.com/daidr/doulog-core/lib/daos"
	"github.com/daidr/doulog-core/lib/ecode"
	"github.com/daidr/doulog-core/lib/format"
	"github.com/daidr/doulog-core/lib/utils"
	"github.com/daidr/doulog-core/module/auth_login/internal/dto"
	"github.com/daidr/doulog-core/module/auth_login/internal/e"
	"github.com/daidr/doulog-core/module/auth_login/internal/service"
	"github.com/gin-gonic/gin"
)

func WRegOptions(c *gin.Context) {
	sp := utils.GetScope(c)
	uid := c.GetUint64("UID")

	options, err := service.BeginWebAuthnRegOptions(sp, uid)

	if err != nil {
		sp.Log.Debugw("failed to begin webauthn registration",
			"error", err,
			"uid", uid)
		format.HTTP(c, e.ErrCreateWebAuthnChallenge, nil)
		return
	}

	format.HTTP(c, ecode.Success, options)
}

func WRegFinish(c *gin.Context) {
	sp := utils.GetScope(c)
	uid := c.GetUint64("UID")

	var req dto.WebAuthnRegFinishReq
	if err := c.ShouldBind(&req); err != nil {
		sp.Log.Debugw("failed to bind request",
			"error", err)
		format.HTTP(c, ecode.InvalidParams, nil)
		return
	}

	credential, err := service.FinishWebAuthnReg(sp, uid, req)

	if err != nil {
		sp.Log.Debugw("failed to finish webauthn registration",
			"error", err,
			"uid", uid)
		format.HTTP(c, e.ErrFinishWebAuthnVerify, nil)
		return
	}

	err = daos.NewUser(sp.DB).AddCredential(uid, credential)
	if err != nil {
		sp.Log.Debugw("failed to add credential to user",
			"error", err,
			"uid", uid)
		format.HTTP(c, e.ErrFinishWebAuthnVerify, nil)
		return
	}

	format.HTTP(c, ecode.Success, nil)
}

func ListWebAuthnCredentials(c *gin.Context) {
	sp := utils.GetScope(c)
	uid := c.GetUint64("UID")

	credentials, err := service.ListWebAuthnCredentials(sp, uid)

	if err != nil {
		sp.Log.Debugw("failed to list webauthn credentials",
			"error", err,
			"uid", uid)
		format.HTTP(c, ecode.UnknownError, nil)
		return
	}

	format.HTTP(c, ecode.Success, credentials)
}
