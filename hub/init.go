package hub

import (
	"fmt"
	"github.com/daidr/doulog-core/hub/internal/db/pgsql"
	"github.com/daidr/doulog-core/hub/internal/db/redis"
	"github.com/daidr/doulog-core/hub/internal/logger"
	"github.com/daidr/doulog-core/hub/internal/middleware"
	"github.com/daidr/doulog-core/lib/conf"
	"github.com/daidr/doulog-core/lib/models"
	"github.com/daidr/doulog-core/lib/search"
	"github.com/patrickmn/go-cache"
	"log"
	"time"

	"github.com/daidr/doulog-core/lib/utils"
	"github.com/daidr/doulog-core/lib/validatorx"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap/zapcore"
)

// Init 快速初始化
func Init() {
	var err error
	if err = conf.Init(); err != nil {
		log.Fatalf("failed to read config: %s", err)
		return
	}
	logger.Init(utils.IF(conf.C.Debug, zapcore.DebugLevel, zapcore.InfoLevel).(zapcore.LevelEnabler))

	hub.Log = logger.NameSpace("hub")
	hub.Cache = cache.New(cache.NoExpiration, 1*time.Minute)

	hub.Log.Info("init config successfully...")

	initDB()
	initMods()

	if err = search.Init(); err != nil {
		hub.Log.Fatalw("failed to init search",
			"error", err)
	}
}

func initDB() {
	mqc := conf.C.PgSQL
	m, err := pgsql.Init(mqc.Host, mqc.Port, mqc.Username, mqc.Password, mqc.Database)
	if err != nil {
		hub.Log.Fatalw("failed to init postgres",
			"error", err)
	}
	if err = m.AutoMigrate(
		&models.TArticle{},
		&models.TArticleTag{},
		&models.TImage{},
		&models.TArticleLike{},
		&models.TOauth{},
		&models.TReply{},
		&models.TTag{},
		&models.TUser{},
		&models.TArticleViews{},
	); err != nil {
		hub.Log.Fatalw("failed to auto migrate table",
			"error", err)
	}

	rc := conf.C.Redis
	r, err := redis.Init(rc.Host, rc.Port, rc.Auth, rc.Database)
	if err != nil {
		hub.Log.Fatalw("failed to init redis",
			"error", err)
	}
	hub.PgSQL = m
	hub.Redis = r

	hub.Log.Info("init db successfully...")
}

func initMods() {
	// gin.SetMode(util.IF(conf.C.Debug, gin.DebugMode, gin.ReleaseMode).(string))
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(middleware.Recovery(), middleware.RequestLog())
	// enable sentry in prod
	if !conf.C.Debug {
		engine.Use(middleware.Sentry())
	}

	// register custom validator
	validatorx.MustRegister("xss", validatorx.NewXSSFunc())
	validatorx.MustRegister("sensitive", validatorx.NewSensitiveFunc())

	hub.HTTP = engine
	api := engine.Group("/api")

	// 注册路由命名空间，根据不同命名空间添加需要的中间件
	ns := make(map[string]*gin.RouterGroup)
	ns[conf.RouterNSMain.NS()] = api.Group(conf.RouterNSMain.NS())
	ns[conf.RouterNSTest.NS()] = api.Group(conf.RouterNSTest.NS())
	ns[conf.RouterNSAuth.NS()] = api.Group(conf.RouterNSAuth.NS())

	// 注册模块路由命名空间 /:namespace/:module
	for _, info := range modules {
		g := ns[info.ID.Namespace()].Group(fmt.Sprintf("/%s", info.ID.Name()))
		hub.Log.Debugw("module register router",
			"mod", info.ID.String(),
			"route", g.BasePath())

		scope := &models.Scope{
			HTTP: g,
			Log:  logger.NameSpace(info.ID.String()),
			DB: &models.DB{
				PgSQL: hub.PgSQL,
				Redis: hub.Redis,
			},
			Cache: hub.Cache,
		}
		g.Use(middleware.Hijack(scope))

		scopes[info.ID.String()] = scope
	}

	hub.Log.Info("init router and mods successfully...")
}
