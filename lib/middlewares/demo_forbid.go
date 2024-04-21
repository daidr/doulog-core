package middlewares

import (
	"github.com/daidr/doulog-core/lib/ecode"
	"github.com/daidr/doulog-core/lib/format"
	"github.com/daidr/doulog-core/lib/user_attr"
	"github.com/gin-gonic/gin"
)

// DemoForbid 测试用户禁止访问接口
func DemoForbid() gin.HandlerFunc {
	return func(c *gin.Context) {
		attr := c.GetInt("USER_ATTR")
		userAttr := user_attr.ParseUserAttr(attr)

		if userAttr.IsDemoUser {
			c.Abort()
			format.HTTP(c, ecode.PermissionDenied, nil)
			return
		}

	}
}
