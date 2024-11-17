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

func Register(c *gin.Context) {
	sp := utils.GetScope(c)
	var req dto.RegisterReq
	if err := c.ShouldBind(&req); err != nil {
		utils.RespLogger(c).Debugw("failed to bind request", "error", err)
		format.HTTP(c, ecode.InvalidParams, nil)
		return
	}

	if err := service.Register(sp, &req); err != nil {
		utils.RespLogger(c).Errorw("failed to register", "error", err)
		// check error is email or username already exists
		if err.Error() == "email already exists" {
			format.HTTP(c, e.ErrUserExists, nil)
			return
		}

		format.HTTP(c, ecode.UnknownError, nil)
		return
	}

	format.HTTP(c, ecode.Success, nil)
}
