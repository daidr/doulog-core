package mod

// Info 模块信息
type Info struct {
	// ID 模块的名称
	// 应全局唯一
	ID ID

	// Instance 返回 Module
	Instance Module
}

func (mi Info) String() string {
	return mi.ID.String()
}
