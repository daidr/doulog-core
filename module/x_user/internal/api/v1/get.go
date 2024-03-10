package v1

import (
	"github.com/daidr/doulog-core/lib/ecode"
	"github.com/daidr/doulog-core/lib/format"
	"github.com/daidr/doulog-core/lib/utils"
	"github.com/daidr/doulog-core/module/x_user/internal/e"
	"github.com/daidr/doulog-core/module/x_user/internal/service/me"
	"github.com/gin-gonic/gin"
)

func GetSelfInfo(c *gin.Context) {
	sp := utils.GetScope(c)
	uid := c.GetUint64("UID")

	resp, err := me.GetSelfInfo(sp.DB, uid)
	if err != nil {
		sp.Log.Debugw("failed to get user info",
			"error", err,
			"uid", uid)
		format.HTTP(c, e.ErrGetSelfInfo, nil)
		return
	}
	format.HTTP(c, ecode.Success, resp)
}
