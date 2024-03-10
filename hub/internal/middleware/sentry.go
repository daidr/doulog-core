package middleware

import (
	"github.com/daidr/doulog-core/hub/internal/logger"
	"github.com/getsentry/sentry-go"
	"github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

func Sentry() gin.HandlerFunc {
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              "https://03d229896f0ec1c97ce2ad3e9540f75b@o1355476.ingest.us.sentry.io/4506878866948096",
		EnableTracing:    true,
		AttachStacktrace: true,
		TracesSampleRate: 1.0,
	}); err != nil {
		logger.NameSpace("sentry").Fatalw("failed to init sentry", "error", err)
		return nil
	}
	return sentrygin.New(sentrygin.Options{
		Repanic: true,
	})
}
