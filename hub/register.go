package hub

import (
	"fmt"
	"github.com/daidr/doulog-core/lib/mod"
	"github.com/daidr/doulog-core/lib/models"
)

// RegMod - 向全局添加 Module
func RegMod(mod mod.Module) {
	info := mod.GetModuleInfo()
	if info.Instance == nil {
		panic(fmt.Errorf("missing ModuleInfo.Instance"))
	}

	if _, ok := modules[info.ID.String()]; ok {
		panic(fmt.Errorf("module already registered: %s", info.ID))
	}
	modules[info.ID.String()] = &info
}

func RegJob(task *models.Task) {
	tasks = append(tasks, task)
}
