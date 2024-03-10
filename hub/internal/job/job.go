package job

import (
	"github.com/daidr/doulog-core/hub/internal/logger"
	"github.com/daidr/doulog-core/lib/models"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"strings"
	"time"
)

func Init(tasks []*models.Task) (*cron.Cron, error) {
	l := &jobLogger{logger: logger.NameSpace("job")}

	c := cron.New(cron.WithLogger(l))

	for _, task := range tasks {
		if _, err := c.AddFunc(task.Spec, task.Cmd); err != nil {
			return nil, err
		}
	}

	return c, nil
}

type jobLogger struct {
	logger *zap.SugaredLogger
}

func (c *jobLogger) Info(msg string, kv ...interface{}) {
	keysAndValues := formatTimes(kv)
	c.logger.Debugf(
		formatString(len(keysAndValues)),
		append([]interface{}{msg}, keysAndValues...)...)
}

func (c *jobLogger) Error(err error, msg string, kv ...interface{}) {
	keysAndValues := formatTimes(kv)
	c.logger.Warnf(
		formatString(len(keysAndValues)+2),
		append([]interface{}{msg, "error", err}, keysAndValues...)...)
}

// formatString copied from default logger
func formatString(numKeysAndValues int) string {
	var sb strings.Builder
	sb.WriteString("%s")
	if numKeysAndValues > 0 {
		sb.WriteString(", ")
	}
	for i := 0; i < numKeysAndValues/2; i++ {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString("%v=%v")
	}
	return sb.String()
}

// formatTimes copied from default logger
func formatTimes(keysAndValues []interface{}) []interface{} {
	var formattedArgs []interface{}
	for _, arg := range keysAndValues {
		if t, ok := arg.(time.Time); ok {
			arg = t.Format(time.RFC3339)
		}
		formattedArgs = append(formattedArgs, arg)
	}
	return formattedArgs
}
