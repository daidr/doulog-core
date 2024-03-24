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

// UpdateUserInfo 更新用户信息
func UpdateUserInfo(c *gin.Context) {
	var req model.UpdateUserInfoReq
	var uri model.UpdateUserInfoUri

	sp := utils.GetScope(c)
	uid := c.GetUint64("UID")
	isAdmin := c.GetBool("ADMIN")

	if err := c.ShouldBind(&req); err != nil {
		format.HTTPInvalidParams(c)
		return
	}

	if err := c.ShouldBindUri(&uri); err != nil {
		format.HTTPInvalidParams(c)
		return
	}

	// 未登录
	if uid == 0 {
		format.HTTP(c, ecode.Unauthorized, nil)
		return
	}

	// 非管理员且修改的不是自己的信息
	if uid != uri.TargetUid && !isAdmin && uri.TargetUid != 0 {
		format.HTTP(c, ecode.PermissionDenied, nil)
		return
	}

	var finalTargetUid uint64
	if uri.TargetUid == 0 {
		finalTargetUid = uid
	} else {
		finalTargetUid = uri.TargetUid
	}

	// 修改用户信息
	if uri.TargetField == "name" {
		if req.NewName == "" {
			format.HTTP(c, ecode.InvalidParams, nil)
			return
		}

		ok, err := service.UpdateName(sp.DB, finalTargetUid, req.NewName)

		if ok {
			format.HTTP(c, ecode.Success, nil)
			return
		}

		if err != nil {
			format.HTTP(c, e.ErrUpdateName, nil)
			sp.Log.Debugw("failed to update user info",
				"error", err,
				"uid", uid,
				"new_name", req.NewName)
			return
		}

		format.HTTP(c, e.ErrNameExisted, nil)
	}

	if uri.TargetField == "email" {
		if req.NewEmail == "" {
			format.HTTP(c, ecode.InvalidParams, nil)
			return
		}

		ok, err := service.UpdateEmail(sp.DB, finalTargetUid, req.NewEmail)

		if ok {
			format.HTTP(c, ecode.Success, nil)
			return
		}

		if err != nil {
			format.HTTP(c, e.ErrUpdateEmail, nil)
			sp.Log.Debugw("failed to update user info",
				"error", err,
				"uid", uid,
				"new_email", req.NewEmail)
			return
		}

		format.HTTP(c, e.ErrEmailExisted, nil)
	}

	if uri.TargetField == "homepage" {
		if req.NewHomepage == "" {
			format.HTTP(c, ecode.InvalidParams, nil)
			return
		}

		ok, err := service.UpdateHomepage(sp.DB, finalTargetUid, req.NewHomepage)

		if ok {
			format.HTTP(c, ecode.Success, nil)
			return
		}

		if err != nil {
			format.HTTP(c, e.ErrUpdateHomepage, nil)
			sp.Log.Debugw("failed to update user info",
				"error", err,
				"uid", uid,
				"new_homepage", req.NewHomepage)
			return
		}
	}

}
