// Package router 模块内部路由定义
package router

import (
	"github.com/daidr/doulog-core/lib/localstorage"
	"github.com/daidr/doulog-core/lib/middlewares"
	"github.com/daidr/doulog-core/lib/models"
	"github.com/daidr/doulog-core/module/x_media/internal/api"
)

func Init(sp *models.Scope) {
	apiPath := sp.HTTP.Group("/")
	{
		apiPath.POST("/upload", middlewares.Auth(true, true, true), middlewares.DemoForbid(), api.MediaUpload) // 上传图片
		apiPath.GET("/metadata/:id", api.GetMediaMetadata)                                                     // 查看图片信息
		apiPath.GET("/:id", api.GetMediaById)                                                                  // 下载图片
		apiPath.GET("/list", middlewares.Auth(true, true, true), api.ListAllMedia)                             // 获取所有图片列表 (admin)
		sp.HubHTTP.Static("/store/"+localstorage.MediaStore.Biz, localstorage.MediaStore.Root)
	}
}
