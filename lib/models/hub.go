package models

import (
	"github.com/gin-gonic/gin"
	goredis "github.com/go-redis/redis/v8"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Scope struct {
	HTTP  *gin.RouterGroup
	Log   *zap.SugaredLogger
	DB    *DB
	Cache *cache.Cache
}

type DB struct {
	PgSQL *gorm.DB
	Redis *goredis.Client
}

type Task struct {
	Spec string
	Cmd  func()
}

type OauthPayload struct {
	Id       string
	Login    string
	Name     string
	Email    string
	Homepage string
}
