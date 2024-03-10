// Package router 模块内部路由定义
package router

import (
	"github.com/daidr/doulog-core/lib/models"
	v1 "github.com/daidr/doulog-core/module/auth_login/internal/api/v1"
)

func Init(sp *models.Scope) {
	apiV1 := sp.HTTP.Group("/v1")
	{
		apiV1.GET("/go", v1.Go)             // 登录跳转接口
		apiV1.GET("/callback", v1.Callback) // 登录回调接口
	}
}
