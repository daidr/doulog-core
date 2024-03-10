package v2

import (
	"github.com/daidr/doulog-core/lib/ecode"
	"github.com/daidr/doulog-core/lib/format"
	"github.com/daidr/doulog-core/lib/utils"
	"github.com/gin-gonic/gin"
	"time"
)

func Ping(c *gin.Context) {
	format.HTTP(c, ecode.Success, gin.H{
		"ver":  "v2",
		"msg":  "pong",
		"time": time.Now().Format("2006-01-02 15:04:05"),
	})
}

// GetReqInfo <host>/api/test/ping/v2/get_req_info
func GetReqInfo(c *gin.Context) {
	scope := utils.GetScope(c)
	scope.Log.Debug("scope has been passed to get req info v2")
	format.HTTP(c, ecode.Success, gin.H{
		"ver":  "v2",
		"ip":   c.ClientIP(),
		"ua":   c.Request.UserAgent(),
		"uri":  c.Request.RequestURI,
		"time": time.Now().Format("2006-01-02 15:04:05"),
	})
}
