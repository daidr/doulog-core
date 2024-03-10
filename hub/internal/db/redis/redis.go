package redis

import (
	"context"
	"fmt"
	"github.com/daidr/doulog-core/hub/internal/logger"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type redisLogger struct {
	Logger *zap.SugaredLogger
}

func Init(host string, port int, auth string, db int) (*redis.Client, error) {
	c := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		DB:       db,
		Password: auth,
	})
	_, err := c.Ping(context.Background()).Result()
	redis.SetLogger(&redisLogger{Logger: logger.NameSpace("redis")})
	return c, err
}

func (l *redisLogger) Printf(ctx context.Context, format string, v ...interface{}) {
	_ = ctx
	l.Logger.Infof(format, v)
}
