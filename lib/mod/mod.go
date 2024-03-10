package mod

import (
	"github.com/daidr/doulog-core/lib/models"
	"sync"
)

type Module interface {
	GetModuleInfo() Info

	// Module 的生命周期

	// Init 初始化
	// 待所有 Module 初始化完成后
	// 进行服务注册 Serve
	Init(sp *models.Scope)

	// PostInit 第二次初始化
	// 调用该函数时，所有 Module 都已完成第一段初始化过程
	// 方便进行跨Module调用
	PostInit(sp *models.Scope)

	// Serve 注册服务函数
	// 结束后调用 Start
	Serve(sp *models.Scope)

	// Start 启用Module
	// 此处调用为
	// ``` go
	// go Start()
	// ```
	// 结束后正式开启服务
	Start(sp *models.Scope)

	// Stop 应用结束时对所有 Module 进行通知
	// 在此进行资源回收
	Stop(sp *models.Scope, wg *sync.WaitGroup)
}
