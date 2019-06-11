package logx

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

type Entry struct {
	*logrus.Entry
}

type Fields logrus.Fields

type Level logrus.Level

const ErrorKey = "error"

func (level Level) String() string {
	switch level {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warning"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	case PanicLevel:
		return "panic"
	}

	return "unknown"
}

// ParseLevel takes a string level and returns the Logrus log level constant.
func ParseLevel(lvl string) (Level, error) {
	switch lvl {
	case "panic":
		return PanicLevel, nil
	case "fatal":
		return FatalLevel, nil
	case "error":
		return ErrorLevel, nil
	case "warn", "warning":
		return WarnLevel, nil
	case "info":
		return InfoLevel, nil
	case "debug":
		return DebugLevel, nil
	}

	var l Level
	return l, fmt.Errorf("not a valid logrus Level: %q", lvl)
}

const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel Level = iota
	// FatalLevel level. Logs and then calls `os.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking httprxr.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
)

func New() *Logger {
	return &Logger{
		logrus.New(),
	}
}

func (logger *Logger) SetLevel(level Level) {
	logger.Logger.SetLevel(logrus.Level(level))
}

func NewEntry(logger *Logger) *Entry {
	return &Entry{logrus.NewEntry(logger.Logger)}
}

func (entry *Entry) WithError(err error) *Entry {
	return &Entry{entry.Entry.WithError(err)}
}

func (entry *Entry) WithField(key string, value interface{}) *Entry {
	return entry.WithFields(Fields{key: value})
}

func (entry *Entry) Log(v ...interface{}) {
	entry.Info(v...)
}
func (entry *Entry) Logf(format string, v ...interface{}) {
	entry.Infof(format, v...)
}

// Add a map of fields to the Entry.
func (entry *Entry) WithFields(fields Fields) *Entry {
	return &Entry{entry.Entry.WithFields(logrus.Fields(fields))}
}

func (logger *Logger) WithField(key string, value interface{}) *Entry {
	return NewEntry(logger).WithField(key, value)
}

// Adds a struct of fields to the log entry. All it does is call `WithField` for
// each `Field`.
func (logger *Logger) WithFields(fields Fields) *Entry {
	return NewEntry(logger).WithFields(fields)
}

func (logger *Logger) Log(v ...interface{}) {
	logger.Info(v...)
}
func (logger *Logger) Logf(format string, v ...interface{}) {
	logger.Infof(format, v...)
}

type StdLog interface {
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
}
