package x_user

import (
	"github.com/daidr/doulog-core/hub"
	"github.com/daidr/doulog-core/lib/conf"
	"github.com/daidr/doulog-core/lib/ecode"
	"github.com/daidr/doulog-core/lib/mod"
	"github.com/daidr/doulog-core/lib/models"
	"github.com/daidr/doulog-core/module/x_user/internal/e"
	"github.com/daidr/doulog-core/module/x_user/internal/router"
	"sync"
)

func init() {
	ins = &user{}
	hub.RegMod(ins)
}

var ins *user

type user struct{}

func (m *user) GetModuleInfo() mod.Info {
	return mod.Info{
		ID:       mod.NewModuleID(conf.RouterNSMain, "user"),
		Instance: ins,
	}
}

func (m *user) Init(scope *models.Scope) {
	ecode.Register(e.ECode)
	_ = scope
	// 初始化过程
	// 在此处可以进行 Module 的初始化配置
	// 如配置读取
}

func (m *user) PostInit(scope *models.Scope) {
	_ = scope
	// 第二次初始化
	// 再次过程中可以进行跨Module的动作
	// 如通用数据库等等
}

func (m *user) Serve(scope *models.Scope) {
	// 注册服务函数部分
	router.Init(scope)
}

func (m *user) Start(scope *models.Scope) {
	_ = scope
	// 此函数会新开携程进行调用
	// ```go
	// 		go exampleModule.Start()
	// ```

	// 可以利用此部分进行后台操作
	// 如http服务器等等
}

func (m *user) Stop(scope *models.Scope, wg *sync.WaitGroup) {
	_ = scope
	// 别忘了解锁
	defer wg.Done()
	// 结束部分
	// 一般调用此函数时，程序接收到 os.Interrupt 信号
	// 即将退出
	// 在此处应该释放相应的资源或者对状态进行保存
}
