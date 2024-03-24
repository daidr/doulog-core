// Package router 模块内部路由定义
package router

import (
	"github.com/daidr/doulog-core/lib/middlewares"
	"github.com/daidr/doulog-core/lib/models"
	v1 "github.com/daidr/doulog-core/module/x_user/internal/api/v1"
)

func Init(sp *models.Scope) {
	apiV1 := sp.HTTP.Group("/v1")
	{
		// field: name, email, homepage
		apiV1.PUT("/:field", middlewares.Auth(true, true, false), v1.UpdateUserInfo)      // 修改用户信息
		apiV1.PUT("/:field/:uid", middlewares.Auth(true, true, false), v1.UpdateUserInfo) // 修改用户信息

		apiV1.GET("", middlewares.Auth(false, false, false), v1.GetUserInfo)     // 获取自己的信息
		apiV1.GET("/:uid", middlewares.Auth(false, true, false), v1.GetUserInfo) // 获取用户信息

		apiV1.GET("/list", middlewares.Auth(true, true, true), v1.ListAllUsers) // 获取所有用户列表 (admin)
	}
}
