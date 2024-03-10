package test_ping

import (
	"github.com/daidr/doulog-core/hub"
	"github.com/daidr/doulog-core/lib/conf"
	"github.com/daidr/doulog-core/lib/ecode"
	"github.com/daidr/doulog-core/lib/mod"
	"github.com/daidr/doulog-core/lib/models"
	"github.com/daidr/doulog-core/module/test_ping/internal/e"
	"github.com/daidr/doulog-core/module/test_ping/internal/router"
	"github.com/daidr/doulog-core/module/test_ping/internal/util"
	"sync"
)

func init() {
	ins = &ping{}
	hub.RegMod(ins)
}

var ins *ping

type ping struct{}

func (m *ping) GetModuleInfo() mod.Info {
	return mod.Info{
		// ID不应当随意设置，需根据前端对接的接口进行设置，且不能轻易改动
		// NameSpace 为一级路由命名空间，api将大致分为 test,x,auth...
		// ModuleName 为二级路由命名空间，api举例分为 /test/ping,/x/article/get....
		// Note: RouterNSMain对应的路由为x
		ID:       mod.NewModuleID(conf.RouterNSTest, "ping"),
		Instance: ins,
	}
}

func (m *ping) Init(scope *models.Scope) {
	// 每个模块初始化应当初始化模块内部的错误码
	ecode.Register(e.ECode)
	scope.Log.Debug("init ping ...")
	// 初始化过程
	// 在此处可以进行 Module 的初始化配置
	// 如配置读取
}

func (m *ping) PostInit(scope *models.Scope) {
	// 调用模块内util演示
	scope.Log.Debugf("module post init, info: %s", util.Ping())
	// 第二次初始化
	// 再次过程中可以进行跨Module的动作
	// 如通用数据库等等
}

func (m *ping) Serve(scope *models.Scope) {
	scope.Log.Debug("start serving /test/ping/*...")
	// 这是一个完整的项目结构组成，如果项目结构较为简单，可以直接在这里进行路由注册
	router.Init(scope)
}

func (m *ping) Start(scope *models.Scope) {
	scope.Log.Debug("module start...")
	// 此函数会新开携程进行调用
	// ```go
	// 		go exampleModule.Start()
	// ```

	// 可以利用此部分进行后台操作
	// 如http服务器等等
}

func (m *ping) Stop(scope *models.Scope, wg *sync.WaitGroup) {
	scope.Log.Debug("module stop")
	// 别忘了解锁
	defer wg.Done()
	// 结束部分
	// 一般调用此函数时，程序接收到 os.Interrupt 信号
	// 即将退出
	// 在此处应该释放相应的资源或者对状态进行保存
}
