package v1

import (
	"github.com/daidr/doulog-core/lib/ecode"
	"github.com/daidr/doulog-core/lib/format"
	"github.com/daidr/doulog-core/lib/utils"
	"github.com/daidr/doulog-core/module/x_user/internal/e"
	"github.com/daidr/doulog-core/module/x_user/internal/model"
	"github.com/daidr/doulog-core/module/x_user/internal/service"
	"github.com/gin-gonic/gin"
)

func ListAllUsers(c *gin.Context) {
	sp := utils.GetScope(c)
	var req model.UserListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.RespLogger(c).Debugw("failed to bind request", "error", err)
		format.HTTP(c, ecode.InvalidParams, nil)
		return
	}

	if req.Page < 1 {
		req.Page = 1
	}

	if req.PageSize < 10 {
		req.PageSize = 10
	}

	if req.PageSize > 100 {
		req.PageSize = 10
	}

	resp, err := service.ListAllUsers(sp.DB, &req)

	if err != nil {
		sp.Log.Debugw("failed to list all users",
			"error", err,
			"req", req)
		format.HTTP(c, e.ErrGetUserInfo, nil)
		return
	}

	format.HTTP(c, ecode.Success, resp)
}
