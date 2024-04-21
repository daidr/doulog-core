// Package router 模块内部路由定义
package router

import (
	"github.com/daidr/doulog-core/lib/middlewares"
	"github.com/daidr/doulog-core/lib/models"
	"github.com/daidr/doulog-core/module/x_user/internal/api"
)

func Init(sp *models.Scope) {
	apiPath := sp.HTTP.Group("/")
	{
		// field: name, email, homepage
		apiPath.PUT("/:field", middlewares.Auth(true, true, false), middlewares.DemoForbid(), api.UpdateUserInfo)      // 修改用户信息
		apiPath.PUT("/:field/:uid", middlewares.Auth(true, true, false), middlewares.DemoForbid(), api.UpdateUserInfo) // 修改用户信息

		apiPath.GET("", middlewares.Auth(false, false, false), api.GetUserInfo)     // 获取自己的信息
		apiPath.GET("/:uid", middlewares.Auth(false, true, false), api.GetUserInfo) // 获取用户信息

		apiPath.GET("/list", middlewares.Auth(true, true, true), api.ListAllUsers) // 获取所有用户列表 (admin)
	}
}
