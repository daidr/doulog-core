package v1

import (
	"github.com/daidr/doulog-core/lib/format"
	"github.com/daidr/doulog-core/module/auth_login/e"
	"net/http"

	"github.com/daidr/doulog-core/lib/auth"
	"github.com/daidr/doulog-core/lib/utils"
	"github.com/daidr/doulog-core/module/auth_login/internal/service"
	"github.com/gin-gonic/gin"
)

func Go(c *gin.Context) {
	platform := c.Query("platform")
	callback := c.Query("callback")

	// check if callback does not start with allowed callback
	if !utils.IsAllowedFrontendCallback(callback) {
		utils.RespLogger(c).Debugw("invalid callback",
			"callback", callback)
		c.Abort()
		// callback 不合法，不进行重定向，直接返回错误
		format.HTTP(c, e.ErrInvalidCallback, nil)
		return
	}

	redirect, err := service.Go(utils.GetScope(c).DB, platform, callback)
	if err != nil {
		utils.RespLogger(c).Debugw("failed to get redirect url",
			"error", err)
		c.Abort()
		c.Redirect(http.StatusFound, auth.ReLoginWithMsg(callback, "Failed to get redirect url"))
		return
	}

	c.Redirect(http.StatusFound, redirect)
}
