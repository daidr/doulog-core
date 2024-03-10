package hub

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/daidr/doulog-core/hub/internal/job"
	"github.com/daidr/doulog-core/lib/conf"
	"github.com/daidr/doulog-core/lib/mod"
	"github.com/daidr/doulog-core/lib/models"
	"github.com/gin-gonic/gin"
	goredis "github.com/go-redis/redis/v8"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	hub = &struct {
		HTTP  *gin.Engine
		Log   *zap.SugaredLogger
		PgSQL *gorm.DB
		Redis *goredis.Client
		Cron  *cron.Cron
		Cache *cache.Cache
	}{}
	modules   = make(map[string]*mod.Info)
	tasks     = make([]*models.Task, 0)
	scopes    = make(map[string]*models.Scope)
	server    *http.Server
	onceStart sync.Once
)

// Run 正式开启服务
func Run() {
	var err error
	// init job
	if hub.Cron, err = job.Init(tasks); err != nil {
		hub.Log.Fatalw("failed to init job", "err", err)
	}
	hub.Cron.Start()

	go func() {
		hub.Log.Info("http engine starting...")
		server = &http.Server{
			Addr:    fmt.Sprintf(":%d", conf.C.Server.Port),
			Handler: hub.HTTP,
		}

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zap.S().Infow("failed to start to listen and serve", "error", err, "port", conf.C.Server.Port)
		}
	}()
}

// StartService 启动服务
// 根据 Module 生命周期 此过程应在Login前调用
func StartService() {
	onceStart.Do(doStartService)
}

func doStartService() {
	for _, mi := range modules {
		mi.Instance.Init(scopes[mi.ID.String()])
	}
	for _, mi := range modules {
		mi.Instance.PostInit(scopes[mi.ID.String()])
	}
	for _, mi := range modules {
		mi.Instance.Serve(scopes[mi.ID.String()])
	}
	for _, mi := range modules {
		go mi.Instance.Start(scopes[mi.ID.String()])
	}
}

// Stop 停止所有服务
// 调用此函数并不会使Bot离线
func Stop() {
	hub.Log.Warn("stopping ...")
	hub.Cron.Stop()

	wg := sync.WaitGroup{}
	wg.Add(len(modules))
	for _, mi := range modules {
		mi.Instance.Stop(scopes[mi.ID.String()], &wg)
	}
	wg.Wait()
	ctx, stop := context.WithTimeout(context.Background(), 5*time.Second)
	defer stop()

	if err := server.Shutdown(ctx); err != nil {
		zap.S().Fatalw("server forced to shutdown", "error", err)
	}

	hub.Log.Info("stopped")
	modules = make(map[string]*mod.Info)
	scopes = make(map[string]*models.Scope)
}
