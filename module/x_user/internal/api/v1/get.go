package v1

import (
	"github.com/daidr/doulog-core/lib/ecode"
	"github.com/daidr/doulog-core/lib/format"
	"github.com/daidr/doulog-core/lib/utils"
	"github.com/daidr/doulog-core/module/x_user/internal/e"
	"github.com/daidr/doulog-core/module/x_user/internal/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetSelfInfo(c *gin.Context) {
	sp := utils.GetScope(c)
	uid := c.GetUint64("UID")

	resp, err := service.GetUserInfo(sp.DB, uid)
	if err != nil {
		sp.Log.Debugw("failed to get user info",
			"error", err,
			"uid", uid)
		format.HTTP(c, e.ErrGetUserInfo, nil)
		return
	}
	format.HTTP(c, ecode.Success, resp)
}

func GetUserInfo(c *gin.Context) {
	sp := utils.GetScope(c)
	uid := c.Param("uid")

	uidInt, err := strconv.ParseUint(uid, 10, 64)
	if err != nil {
		format.HTTPInvalidParams(c)
		return
	}

	resp, err := service.GetUserInfo(sp.DB, uidInt)
	if err != nil {
		sp.Log.Debugw("failed to get user info",
			"error", err,
			"uid", uid)
		format.HTTP(c, e.ErrGetUserInfo, nil)
		return
	}
	format.HTTP(c, ecode.Success, resp)
}
