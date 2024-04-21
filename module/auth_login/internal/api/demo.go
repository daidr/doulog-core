package api

import (
	"github.com/daidr/doulog-core/lib/auth"
	"github.com/daidr/doulog-core/lib/format"
	"github.com/daidr/doulog-core/lib/utils"
	"github.com/daidr/doulog-core/module/auth_login/e"
	"github.com/daidr/doulog-core/module/auth_login/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func DemoLogin(c *gin.Context) {
	sp := utils.GetScope(c)
	callback := c.Query("callback")

	token := service.DemoLogin(sp.DB)

	if callback == "" {
		c.Abort()
		// callback 不合法，不进行重定向，直接返回错误
		format.HTTP(c, e.ErrInvalidCallback, nil)
	}

	c.Redirect(http.StatusFound, auth.ReLoginWithToken(callback, token))
}
