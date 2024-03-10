package main

import (
	"github.com/daidr/doulog-core/hub"
	"os"
	"os/signal"

	/** Test */
	_ "github.com/daidr/doulog-core/module/test_ping" // 测试

	/** Auth */
	_ "github.com/daidr/doulog-core/module/auth_login" // 登录

	/** X */
	_ "github.com/daidr/doulog-core/module/x_user" // 用户相关
)

func main() {
	hub.Init()
	hub.StartService()
	hub.Run()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch
	hub.Stop()
}
