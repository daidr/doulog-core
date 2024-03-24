package v1

import (
	"github.com/daidr/doulog-core/lib/ecode"
	"github.com/daidr/doulog-core/lib/format"
	"github.com/daidr/doulog-core/lib/utils"
	"github.com/daidr/doulog-core/module/x_user/internal/e"
	"github.com/daidr/doulog-core/module/x_user/internal/service"
	"github.com/gin-gonic/gin"
)

func GetUserInfo(c *gin.Context) {
	sp := utils.GetScope(c)

	req := struct {
		TargetUid uint64 `uri:"uid"`
	}{}

	if err := c.ShouldBindUri(&req); err != nil {
		format.HTTPInvalidParams(c)
		return
	}

	uid := c.GetUint64("UID")
	isAdmin := c.GetBool("ADMIN")

	if req.TargetUid != uid && !isAdmin && req.TargetUid != 0 {
		format.HTTP(c, ecode.PermissionDenied, nil)
		return
	}

	var finalTargetUid uint64
	if req.TargetUid == 0 {
		finalTargetUid = uid
	} else {
		finalTargetUid = req.TargetUid
	}

	resp, err := service.GetUserInfo(sp.DB, finalTargetUid)
	if err != nil {
		sp.Log.Debugw("failed to get user info",
			"error", err,
			"targetUid", req.TargetUid)
		format.HTTP(c, e.ErrGetUserInfo, nil)
		return
	}
	format.HTTP(c, ecode.Success, resp)
}
