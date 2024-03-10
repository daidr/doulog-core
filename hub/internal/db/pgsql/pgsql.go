package pgsql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/driver/postgres"

	"github.com/daidr/doulog-core/hub/internal/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type dbLogger struct {
	Logger        *zap.SugaredLogger
	SlowThreshold time.Duration
	TraceStr      string
	TraceWarnStr  string
	TraceErrStr   string
}

func (d *dbLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	_ = level
	return d
}

func (d *dbLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	_ = ctx
	d.Logger.Info(msg, data)
}

func (d *dbLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	_ = ctx
	d.Logger.Warn(msg, data)
}

func (d *dbLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	_ = ctx
	d.Logger.Error(msg, data)
}

// Trace modified from gorm default logger
func (d *dbLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	_ = ctx
	elapsed := time.Since(begin)
	switch {
	case err != nil && (!errors.Is(err, gorm.ErrRecordNotFound)):
		sql, rows := fc()
		if rows == -1 {
			d.Logger.Warnf(d.TraceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			d.Logger.Warnf(d.TraceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > d.SlowThreshold && d.SlowThreshold != 0:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", d.SlowThreshold)
		if rows == -1 {
			d.Logger.Warnf(d.TraceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			d.Logger.Warnf(d.TraceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}
func newLogger(logger *zap.SugaredLogger) *dbLogger {
	return &dbLogger{
		Logger:        logger,
		SlowThreshold: 200 * time.Millisecond,
		TraceStr:      "%s\n[%.3fms] [rows:%v] %s",
		TraceWarnStr:  "%s %s\n[%.3fms] [rows:%v] %s",
		TraceErrStr:   "%s %s\n[%.3fms] [rows:%v] %s",
	}
}
func Init(host string, port int, user, password, db string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		host, port, user, password, db)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		Logger:      newLogger(logger.NameSpace("pgsql")),
		PrepareStmt: true,
		QueryFields: true,
	})
}
