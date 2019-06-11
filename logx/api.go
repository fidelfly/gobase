package logx

import (
	"io"

	"github.com/sirupsen/logrus"
)

var std = &Logger{logrus.StandardLogger()}

func StandardLogger() *Logger {
	return std
}

func SetStandard(logger *Logger) {
	std = logger
}

// SetOutput sets the standard logger output.
func SetOutput(out io.Writer) {
	logrus.SetOutput(out)
}

// SetFormatter sets the standard logger Formatter.
func SetFormatter(formatter Formatter) {
	logrus.SetFormatter(&logrusFormatter{formatter})
}

//export
func SetLogrusFormatter(formatter logrus.Formatter) {
	logrus.SetFormatter(formatter)
}

//export
// SetLevel sets the standard logger level.
func SetLevel(level Level) {
	logrus.SetLevel(logrus.Level(level))
}

//export
// GetLevel returns the standard logger level.
func GetLevel() Level {
	return Level(logrus.GetLevel())
}

//export
// AddHook adds a hook to the standard logger hooks.
func AddHook(hook logrus.Hook) {
	logrus.AddHook(hook)
}

//export
// WithError creates an entry from the standard logger and adds an error to it, using the value defined in ErrorKey as key.
func WithError(err error) *Entry {
	return std.WithField(ErrorKey, err)
}

//export
// WithField creates an entry from the standard logger and adds a field to
// it. If you want multiple fields, use `WithFields`.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithField(key string, value interface{}) *Entry {
	return std.WithField(key, value)
}

//export
// WithFields creates an entry from the standard logger and adds multiple
// fields to it. This is simply a helper for `WithField`, invoking it
// once for each field.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithFields(fields Fields) *Entry {
	return std.WithFields(fields)
}

//export
func Info(args ...interface{}) {
	logrus.Info(args...)
}

//export
func Infof(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

//export
func Error(args ...interface{}) {
	logrus.Error(args...)
}

//export
func Errorf(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}

//export
func Warn(args ...interface{}) {
	logrus.Warn(args...)
}

//export
func Warnf(format string, args ...interface{}) {
	logrus.Warnf(format, args...)
}

//export
func Debug(args ...interface{}) {
	logrus.Debug(args...)
}

//export
func Debugf(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}

//export
func Panic(args ...interface{}) {
	logrus.Panic(args...)
}

//export
func Panicf(format string, args ...interface{}) {
	logrus.Panicf(format, args...)
}

//export
func CaptureErrorWith(logger StdLog, args ...interface{}) {
	if len(args) > 0 {
		for _, arg := range args {
			if arg != nil {
				if err, ok := arg.(error); ok {
					logger.Error(err)
				}
			}
		}
	}
}

//export
func CaptureError(args ...interface{}) {
	CaptureErrorWith(StandardLogger(), args...)
}
