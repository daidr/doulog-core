package mod

import (
	"github.com/daidr/doulog-core/lib/conf"
	"strings"
)

// ID 模块ID
// 请使用 小写 并用 _ 代替空格
// Example:
// - test.ping
type ID interface {
	Namespace() string // 一级路由命名空间，使用conf.RouterNSxxx
	Name() string      // 模块名
	String() string
}

type id struct {
	namespace string // 一级路由命名空间，使用conf.RouterNSxxx
	name      string // 模块名
}

func (id *id) Namespace() string {
	return id.namespace
}

func (id *id) Name() string {
	return id.name
}

// String 实现 Stringer 接口
// 凡是将 ID 转成 string 一律用这个接口, 无论是输出还是获取模块名
func (id *id) String() string {
	return id.namespace + "." + id.name
}

// NewModuleID 构造函数，统一生成 ID，避免非法模块名的存在
func NewModuleID(namespace conf.RouterNS, moduleName string) ID {
	ns := strings.TrimSpace(namespace.NS())
	name := strings.TrimSpace(moduleName)
	if ns == "" || name == "" {
		panic("模块名的namespace和name均不能为空白")
	}
	return &id{namespace: namespace.NS(), name: moduleName}
}
