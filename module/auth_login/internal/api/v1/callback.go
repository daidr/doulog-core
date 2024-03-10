package v1

import (
	"github.com/daidr/doulog-core/lib/auth"
	"github.com/daidr/doulog-core/lib/format"
	"github.com/daidr/doulog-core/lib/utils"
	"github.com/daidr/doulog-core/module/auth_login/e"
	"github.com/daidr/doulog-core/module/auth_login/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Callback(c *gin.Context) {
	sp := utils.GetScope(c)
	state := c.Query("state")
	code := c.Query("code")

	token, callback, err := service.Callback(sp.DB, state, code)
	
	if callback == "" {
		c.Abort()
		// callback 不合法，不进行重定向，直接返回错误
		format.HTTP(c, e.ErrInvalidCallback, nil)
	}

	if err != nil {
		sp.Log.Debugw("failed to verify callback",
			"error", err,
			"state", state,
			"code", code)
		c.Abort()
		c.Redirect(http.StatusFound, auth.ReLoginWithMsg(callback, "Failed to verify callback"))
		return
	}

	c.Redirect(http.StatusFound, auth.ReLoginWithToken(callback, token))
}
