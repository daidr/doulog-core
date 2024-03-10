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
		apiV1.GET("/me", middlewares.Auth(false), v1.GetSelfInfo) // 获取自己的信息

		apiV1.POST("/me/change_name", middlewares.Auth(true), v1.ChangeNickName) // 修改昵称
	}
}
