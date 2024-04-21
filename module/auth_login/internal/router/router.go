// Package router 模块内部路由定义
package router

import (
	"github.com/daidr/doulog-core/lib/models"
	"github.com/daidr/doulog-core/module/auth_login/internal/api"
)

func Init(sp *models.Scope) {
	apiPath := sp.HTTP.Group("/")
	{
		apiPath.GET("/go", api.Go)             // 登录跳转接口
		apiPath.GET("/callback", api.Callback) // 登录回调接口
		apiPath.GET("/demo", api.DemoLogin)    // 登录回调接口
	}
}
