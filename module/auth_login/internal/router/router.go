// Package router 模块内部路由定义
package router

import (
	"github.com/daidr/doulog-core/lib/middlewares"
	"github.com/daidr/doulog-core/lib/models"
	"github.com/daidr/doulog-core/module/auth_login/internal/api"
)

func Init(sp *models.Scope) {
	apiPath := sp.HTTP.Group("/")
	{
		apiPath.GET("/login/go", api.Go)             // OAuth登录跳转接口
		apiPath.GET("/login/callback", api.Callback) // OAuth登录回调接口
		apiPath.POST("/register", api.Register)      // 用户注册接口
		apiPath.POST("/login/password", api.PwLogin) // 用户密码登录接口
		apiPath.POST("/webauthn/reg/options", middlewares.Auth(true, false, false), api.WRegOptions)
		apiPath.POST("/webauthn/reg/finish", middlewares.Auth(true, false, false), api.WRegFinish)
		apiPath.POST("/webauthn/credentials", middlewares.Auth(true, false, false), api.ListWebAuthnCredentials)

	}
}
