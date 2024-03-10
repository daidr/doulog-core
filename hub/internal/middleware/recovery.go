package middleware

import (
	"github.com/daidr/doulog-core/lib/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				utils.GetScope(c).Log.Error(err)
				c.Status(http.StatusInternalServerError)
				c.Abort()
			}
		}()
		c.Next()
	}
}
