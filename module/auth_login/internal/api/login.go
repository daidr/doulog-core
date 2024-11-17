package api

import (
	"github.com/daidr/doulog-core/lib/ecode"
	"github.com/daidr/doulog-core/lib/format"
	"github.com/daidr/doulog-core/lib/utils"
	"github.com/daidr/doulog-core/module/auth_login/internal/dto"
	"github.com/daidr/doulog-core/module/auth_login/internal/e"
	"github.com/daidr/doulog-core/module/auth_login/internal/service"
	"github.com/gin-gonic/gin"
)

func PwLogin(c *gin.Context) {
	sp := utils.GetScope(c)
	var req dto.PwLoginReq
	if err := c.ShouldBind(&req); err != nil {
		utils.RespLogger(c).Debugw("failed to bind request", "error", err)
		format.HTTP(c, ecode.InvalidParams, nil)
		return
	}

	token, err := service.PwLogin(sp, &req)

	if err != nil {
		utils.RespLogger(c).Debugw("failed to login", "error", err)
		switch err.Error() {
		case "email not exists":
			format.HTTP(c, e.ErrEmailNotExists, nil)
		case "password error":
			format.HTTP(c, e.ErrPasswordWrong, nil)
		default:
			format.HTTP(c, ecode.UnknownError, nil)
		}
		return
	}

	format.HTTP(c, ecode.Success, token)
}
