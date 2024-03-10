package auth_login

import (
	"github.com/daidr/doulog-core/hub"
	"github.com/daidr/doulog-core/lib/conf"
	"github.com/daidr/doulog-core/lib/ecode"
	"github.com/daidr/doulog-core/lib/mod"
	"github.com/daidr/doulog-core/lib/models"
	"github.com/daidr/doulog-core/module/auth_login/e"
	"github.com/daidr/doulog-core/module/auth_login/internal/router"
	"sync"
)

func init() {
	ins = &login{}
	hub.RegMod(ins)
}

var ins *login

type login struct{}

func (m *login) GetModuleInfo() mod.Info {
	return mod.Info{
		ID:       mod.NewModuleID(conf.RouterNSAuth, "login"),
		Instance: ins,
	}
}

func (m *login) Init(scope *models.Scope) {
	ecode.Register(e.ECode)
	_ = scope

	// Register Cron Job
	// hub.RegJob(&models.Task{
	// 	Spec: "@every 1h",
	// 	Cmd:  task.Job(),
	// })

	// 初始化过程
	// 在此处可以进行 Module 的初始化配置
	// 如配置读取
}

func (m *login) PostInit(scope *models.Scope) {
	_ = scope
	// 第二次初始化
	// 再次过程中可以进行跨Module的动作
	// 如通用数据库等等
}

func (m *login) Serve(scope *models.Scope) {
	// 注册服务函数部分
	router.Init(scope)
}

func (m *login) Start(scope *models.Scope) {
	_ = scope
	// 此函数会新开携程进行调用
	// ```go
	// 		go exampleModule.Start()
	// ```

	// 可以利用此部分进行后台操作
	// 如http服务器等等
}

func (m *login) Stop(scope *models.Scope, wg *sync.WaitGroup) {
	_ = scope
	// 别忘了解锁
	defer wg.Done()
	// 结束部分
	// 一般调用此函数时，程序接收到 os.Interrupt 信号
	// 即将退出
	// 在此处应该释放相应的资源或者对状态进行保存
}
