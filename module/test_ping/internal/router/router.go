package router

import (
	"github.com/daidr/doulog-core/lib/middlewares"
	"github.com/daidr/doulog-core/lib/models"
	v1 "github.com/daidr/doulog-core/module/test_ping/internal/api/v1"
	v2 "github.com/daidr/doulog-core/module/test_ping/internal/api/v2"
)

func Init(scope *models.Scope) {
	// 这里将完整演示模块内部分版本注册，可以根据实际情况不分模块内部版本
	scope.Log.Debug("scope has been passed to router.Init")

	// 你可以使用公共包中的Auth中间件来细粒度地控制哪些接口需要鉴权,但不应当重复使用，虽然这不会造成错误

	apiV1 := scope.HTTP.Group("/v1")
	{
		// 实际注册完整路由: <host>/api/test/ping/v1/
		apiV1.GET("/", middlewares.Auth(false), v1.Ping)
		// 实际注册完整路由: <host>/api/test/ping/v1/get_req_info
		apiV1.GET("/get_req_info", v1.GetReqInfo)
		// 实际注册完整路由: <host>/api/test/ping/v1/error
		apiV1.GET("/error", middlewares.Auth(false), v1.Error)
		apiV1.GET("/panic", v1.Panic)

		// 模块内使用Redis演示，你可以随意修改参数试试效果
		// <host>/api/test/ping/v1/redis_set?k=this-is-a-key&v=i-want-to-set-a-value
		// 20秒内访问redis_get可获得value，否则key过期
		apiV1.POST("/redis_set", v1.RedisSet)
		// <host>/api/test/ping/v1/redis_get?k=this-is-a-key
		apiV1.GET("/redis_get", v1.RedisGet)

		// 演示如何在模块内使用HTTP请求
		// <host>/api/test/ping/v1/http_request
		apiV1.GET("/http_request", v1.HttpRequest)
	}

	apiV2 := scope.HTTP.Group("/v2")
	{
		// 实际注册完整路由: <host>/api/test/ping/v2/
		apiV2.GET("/", v2.Ping)
		// 实际注册完整路由: <host>/api/test/ping/v2/get_req_info
		apiV2.GET("/get_req_info", middlewares.Auth(false), v2.GetReqInfo)
	}

}
