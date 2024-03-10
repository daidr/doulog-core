package v1

import (
	"context"
	"fmt"
	"github.com/daidr/doulog-core/lib/ecode"
	"github.com/daidr/doulog-core/lib/format"
	"github.com/daidr/doulog-core/lib/request"
	"github.com/daidr/doulog-core/lib/utils"
	"github.com/daidr/doulog-core/module/test_ping/internal/e"
	"github.com/gin-gonic/gin"
	"time"
)

// Ping <host>/api/test/ping/v1/
func Ping(c *gin.Context) {
	// 获取scope
	scope := utils.GetScope(c)
	scope.Log.Debug("scope has benn passed to Ping v1 handler")
	// 每个scope.Log都会携带框架劫持的trace参数
	scope.Log.Debugw("test debugw", "field1", "1", "field2", "2", "field3", "3", "trace", "tracetest")

	format.HTTP(c, ecode.Success, gin.H{
		"ver":  "v1",
		"msg":  "pong",
		"time": time.Now().Format("2006-01-02 15:04:05"),
	})
}

// GetReqInfo <host>/api/test/ping/v1/get_req_info
func GetReqInfo(c *gin.Context) {
	format.HTTP(c, ecode.Success, gin.H{
		"ver":  "v1",
		"ip":   c.ClientIP(),
		"ua":   c.Request.UserAgent(),
		"uri":  c.Request.RequestURI,
		"time": time.Now().Format("2006-01-02 15:04:05"),
	})
}

// Error <host>/api/test/ping/v1/error
func Error(c *gin.Context) {
	format.HTTP(c, e.ErrPing, nil)
}

func Panic(c *gin.Context) {
	_ = c
	panic("test.ping panic")
}

// RedisSet 这里仅为演示，业务中需要对参数进行校验
func RedisSet(c *gin.Context) {
	sp := utils.GetScope(c)
	// 获取想要设置的value
	v := c.Query("v")

	// Redis的Key必须使用format.Gen()进行构建
	result := sp.DB.Redis.Set(context.Background(),
		format.Key.Gen("test", "redis", "set", c.Query("k")),
		c.Query("v"),
		20*time.Second)

	if result.Err() != nil {
		format.HTTP(c, e.ErrRedisSet, nil)
		sp.Log.Debugw("failed to set kv", "error", result.Err(), "value", v)
		return
	}

	sp.Log.Debugw("redis set succ")
	format.HTTP(c, ecode.Success, gin.H{
		"result": result.String(),
	})
}

// RedisGet 这里仅为演示，业务中需要对参数进行校验
func RedisGet(c *gin.Context) {
	sp := utils.GetScope(c)

	result := sp.DB.Redis.Get(context.Background(),
		format.Key.Gen("test", "redis", "set", c.Query("k")))
	if result.Err() != nil {
		format.HTTP(c, e.ErrRedisGet, nil)
		sp.Log.Debugw("failed to get kv", "error", result.Err())
		return
	}

	// 在业务中应当进行错误处理，这只是测试应当保证已经Set
	v, _ := result.Result()
	format.HTTP(c, ecode.Success, gin.H{
		"result": result.String(),
		"value":  v,
	})
}

func HttpRequest(c *gin.Context) {
	// 记录一下请求时间
	start := time.Now()

	sp := utils.GetScope(c)

	sp.Log.Debugw("start http request...")

	// 所有HTTP请求必须使用 request 包中的 HTTP() 获得一个请求实例，可以在一个handler中非并发地复用
	// HTTP请求使用 https://github.com/guonaihong/gout ，用法很简单，看一下文档就会了

	// 以一言请求为例
	text := ""
	err := request.HTTP().GET("https://v1.hitokoto.cn/?encode=text").BindBody(&text).Do()
	if err != nil {
		format.HTTP(c, e.ErrHttpRequest, nil)
		sp.Log.Debugw("failed to request v1.hitokoto.cn", "error", err)
		return
	}

	format.HTTP(c, ecode.Success, gin.H{
		"text": text,
		"time": fmt.Sprintf("%dms", time.Since(start).Milliseconds()),
	})
}
