package log

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	// ContextKey logrus of context key
	ContextKey = loggerKey("_loggerKey_context")
)

var requestKeys []string

type contextHook struct {
	levels []logrus.Level
}

// SetRequestKeys set request query key or header keys
func SetRequestKeys(keys ...string) {
	requestKeys = keys
}

func newContextHook(levels ...logrus.Level) logrus.Hook {
	hook := contextHook{
		levels: levels,
	}
	if len(hook.levels) == 0 {
		hook.levels = logrus.AllLevels
	}
	return &hook
}

// Levels implement levels
func (hook contextHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire implement fire
func (contextHook) Fire(entry *logrus.Entry) error {
	if entry.Context != nil {
		for _, k := range requestKeys {
			if v := entry.Context.Value(k); v != nil {
				entry.Data[k] = v
			}
		}
	}
	return nil
}

// ContextWithLogrus set logrus into context
func ContextWithLogrus() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// set headers or query to context value
		for _, k := range requestKeys {
			if v := ctx.GetHeader(k); v != "" {
				ctx.Set(k, v)
			}
			if v, ok := ctx.GetQuery(k); ok {
				ctx.Set(k, v)
			}
		}

		// new entry
		entry := FromContext(ctx)
		// set context into entry
		entry = entry.WithContext(ctx)

		// set entry into gin context
		ctx.Set(string(ContextKey), entry)

		ctx.Next()
	}
}
