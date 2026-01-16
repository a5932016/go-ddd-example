package log

import (
	"context"

	"github.com/sirupsen/logrus"
)

// Fields type, used to pass to `WithFields`.
type Fields = logrus.Fields

type loggerKey string

type logFormatter struct {
	logrus.JSONFormatter
}

func (mf logFormatter) Format(e *logrus.Entry) ([]byte, error) {
	mf.JSONFormatter.TimestampFormat = "2006/01/02 15:04:05"
	e.Time = e.Time.UTC()
	return mf.JSONFormatter.Format(e)
}

func init() {
	logrus.SetFormatter(&logFormatter{})
	logrus.AddHook(newSourceHook())
	logrus.AddHook(newContextHook())
}

// SetLevel set log level
func SetLevel(level string) error {
	lv, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}
	logrus.SetLevel(lv)
	return nil
}

// Options from context options
type Options struct {
	Tags map[string]string
}

// FromContext get entry frome context
func FromContext(ctx context.Context, options ...Options) *logrus.Entry {
	var logger *logrus.Entry
	if l := ctx.Value(string(ContextKey)); l != nil {
		logger = l.(*logrus.Entry)
	} else {
		logger = logrus.NewEntry(logrus.StandardLogger())
	}

	if len(options) > 0 {
		option := options[0]

		if len(option.Tags) > 0 {
			var fields logrus.Fields
			for k, v := range option.Tags {
				fields[k] = v
			}
			logger = logrus.WithFields(fields)
		}
	}

	return logger
}

// WithContext creates an entry from the standard logger and adds a context to it.
func WithContext(ctx context.Context) *logrus.Entry {
	return logrus.WithContext(ctx)
}

// WithError creates an entry from the standard logger and adds an error to it, using the value defined in ErrorKey as key.
func WithError(err error) *logrus.Entry {
	return logrus.WithError(err)
}

// WithField creates an entry from the standard logger and adds a field to
// it. If you want multiple fields, use `WithFields`.
func WithField(key string, value interface{}) *logrus.Entry {
	return logrus.WithField(key, value)
}

// WithFields creates an entry from the standard logger and adds multiple
// fields to it. This is simply a helper for `WithField`, invoking it
// once for each field.
func WithFields(fields Fields) *logrus.Entry {
	return logrus.WithFields(fields)
}

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	logrus.Debug(args...)
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	logrus.Info(args...)
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

// Warning logs a message at level Warn on the standard logger.
func Warning(args ...interface{}) {
	logrus.Warning(args...)
}

// Warningf logs a message at level Warn on the standard logger.
func Warningf(format string, args ...interface{}) {
	logrus.Warningf(format, args...)
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	logrus.Error(args...)
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}

// Panic logs a message at level Panic on the standard logger.
func Panic(args ...interface{}) {
	logrus.Panic(args...)
}

// Panicf logs a message at level Panic on the standard logger.
func Panicf(format string, args ...interface{}) {
	logrus.Panicf(format, args...)
}

// Fatal logs a message at level Fatal on the standard logger.
func Fatal(args ...interface{}) {
	logrus.Fatal(args...)
}

// Fatalf logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatalf(format string, args ...interface{}) {
	logrus.Fatalf(format, args...)
}

// GetLevel returns the standard logger level.
func GetLevel() string {
	return logrus.GetLevel().String()
}
