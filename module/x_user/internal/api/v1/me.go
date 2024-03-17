package v1

import (
	"github.com/daidr/doulog-core/lib/ecode"
	"github.com/daidr/doulog-core/lib/format"
	"github.com/daidr/doulog-core/lib/utils"
	"github.com/daidr/doulog-core/module/x_user/internal/e"
	"github.com/daidr/doulog-core/module/x_user/internal/service/me"
	"github.com/gin-gonic/gin"
)

// ChangeNickName 修改昵称
func ChangeNickName(c *gin.Context) {
	req := struct {
		NewName string `json:"new_name" form:"new_name" binding:"required,min=1,max=25,sensitive,xss"` // 限制名字长度
	}{}

	sp := utils.GetScope(c)
	uid := c.GetUint64("UID")

	if err := c.ShouldBind(&req); err != nil {
		format.HTTPInvalidParams(c)
		return
	}

	// 昵称查重并修改
	ok, err := me.ChangeNickName(sp.DB, uid, req.NewName)

	if ok {
		format.HTTP(c, ecode.Success, nil)
		return
	}

	if err != nil {
		format.HTTP(c, e.ErrChangeNickName, nil)
		sp.Log.Debugw("failed to change user nick name",
			"error", err,
			"uid", uid,
			"new_name", req.NewName)
		return
	}

	format.HTTP(c, e.ErrNickNameExisted, nil)
}
