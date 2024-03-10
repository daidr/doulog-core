package middleware

import (
	"github.com/daidr/doulog-core/hub/internal/logger"
	"github.com/daidr/doulog-core/lib/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"time"
)

func RequestLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.NameSpace("middleware")

		trace := uuid.New().String()
		// 用于下面Trace的获取
		c.Set("trace", trace)
		start := time.Now()
		log.Infof("new request...trace id: %s", trace)

		c.Next()

		log.Infof("%s | %3d | %13v | %15s | %s | %s |",
			trace,
			c.Writer.Status(),
			time.Now().Sub(start),
			c.ClientIP(),
			c.Request.Method,
			c.Request.RequestURI)
	}
}

// Hijack 劫持scope的中间件，可以在进入模块handler前做预处理
func Hijack(scope *models.Scope) gin.HandlerFunc {
	return func(c *gin.Context) {
		l := scope.Log.With("_trace", c.MustGet("trace").(string))

		c.Set("scope", &models.Scope{
			HTTP: scope.HTTP,
			// 将上层中间件的trace id传递到handler的scope.Log中，使其始终携带trace输出
			Log: l,
			DB: &models.DB{
				PgSQL: scope.DB.PgSQL,
				Redis: scope.DB.Redis,
			},
			Cache: scope.Cache,
		})
	}
}
