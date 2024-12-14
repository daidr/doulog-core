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

func RemoveWebAuthnCredential(c *gin.Context) {
	sp := utils.GetScope(c)
	uid := c.GetUint64("UID")

	var req dto.WebAuthnCredentialDeleteReq

	if err := c.ShouldBind(&req); err != nil {
		sp.Log.Debugw("failed to bind request",
			"error", err)
		format.HTTP(c, ecode.InvalidParams, nil)
		return
	}

	cid := req.ID

	err := service.DeleteWebAuthnCredential(sp, uid, cid)

	if err != nil {
		sp.Log.Debugw("failed to delete webauthn credential",
			"error", err,
			"uid", uid,
			"cid", cid)
		switch err.Error() {
		case "credential not found":
			format.HTTP(c, e.ErrCredentialNotExists, nil)
		default:
			format.HTTP(c, ecode.UnknownError, nil)
		}
		return
	}

	format.HTTP(c, ecode.Success, nil)
}

func RenameWebAuthnCredential(c *gin.Context) {
	sp := utils.GetScope(c)
	uid := c.GetUint64("UID")

	var req dto.WebAuthnCredentialRenameReq

	if err := c.ShouldBind(&req); err != nil {
		sp.Log.Debugw("failed to bind request",
			"error", err)
		format.HTTP(c, ecode.InvalidParams, nil)
		return
	}

	cid := req.ID
	newLabel := req.Label

	err := service.RenameWebAuthnCredential(sp, uid, cid, newLabel)

	if err != nil {
		sp.Log.Debugw("failed to rename webauthn credential",
			"error", err,
			"uid", uid,
			"cid", cid)
		switch err.Error() {
		case "credential not found":
			format.HTTP(c, e.ErrCredentialNotExists, nil)
		default:
			format.HTTP(c, ecode.UnknownError, nil)
		}
		return
	}

	format.HTTP(c, ecode.Success, nil)
}

func WDiscoverLoginOptions(c *gin.Context) {
	sp := utils.GetScope(c)
	options, err := service.BeginWebAuthnDiscoverLoginOptions(sp)

	if err != nil {
		sp.Log.Debugw("failed to begin webauthn discover login",
			"error", err)
		format.HTTP(c, e.ErrCreateWebAuthnChallenge, nil)
		return
	}

	format.HTTP(c, ecode.Success, options)
}

func WDiscoverLoginFinish(c *gin.Context) {
	sp := utils.GetScope(c)

	var req dto.WebAuthnLoginFinishReq
	if err := c.ShouldBind(&req); err != nil {
		sp.Log.Debugw("failed to bind request",
			"error", err)
		format.HTTP(c, ecode.InvalidParams, nil)
		return
	}

	token, err := service.FinishWebAuthnDiscoverLogin(sp, req)

	if err != nil {
		sp.Log.Debugw("failed to finish webauthn discover login",
			"error", err)
		format.HTTP(c, e.ErrFinishWebAuthnVerify, nil)
		return
	}

	format.HTTP(c, ecode.Success, token)
}
