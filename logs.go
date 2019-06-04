package fxgo

import (
	"fmt"
	"io"
	"os"

	"github.com/fidelfly/fxgo/logx"
)

func init() {
	logx.SetOutput(os.Stdout)
	logx.SetFormatter(logx.PlainFormatter("2006-01-02 15:04:05"))
}

type LogConfig struct {
	LogLevel  string
	LogPath   string
	LogFile   string
	MaxSize   int
	MaxAge    int
	MaxBackup int
	Compress  bool
	Stdout    bool
}

type ConsoleOutput struct {
}

func (so ConsoleOutput) Info(args ...interface{}) {
	fmt.Println(args...)
}

func (so ConsoleOutput) Infof(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (so ConsoleOutput) Error(args ...interface{}) {
	fmt.Println(args...)
}

func (so ConsoleOutput) Errorf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (so ConsoleOutput) Warn(args ...interface{}) {
	fmt.Println(args...)
}

func (so ConsoleOutput) Warnf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (so ConsoleOutput) Debug(args ...interface{}) {
	fmt.Println(args...)
}

func (so ConsoleOutput) Debugf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (so ConsoleOutput) Panic(args ...interface{}) {
	fmt.Println(args...)
}

func (so ConsoleOutput) Panicf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

//export
func SetupLogs(config *LogConfig) {
	configLogger(logx.StandardLogger(), config)
}

func configLogger(logger *logx.Logger, config *LogConfig) {
	level, err0 := logx.ParseLevel(config.LogLevel)
	if err0 != nil {
		level = logx.WarnLevel
	}
	logger.SetLevel(level)

	if len(config.LogFile) == 0 {
		logger.SetOutput(os.Stdout)
	} else {
		logPath := config.LogPath
		if len(logPath) == 0 {
			logPath = "."
		}
		rotate := logx.RotateLog(fmt.Sprintf("%s/%s", logPath, config.LogFile), config.MaxSize, config.MaxBackup, config.MaxAge, config.Compress)

		if config.Stdout {
			logger.SetOutput(io.MultiWriter(os.Stdout, rotate))
		} else {
			logger.SetOutput(rotate)
		}
	}
}

//export
func NewLog(config *LogConfig) *logx.Logger {
	logger := logx.New()
	configLogger(logger, config)
	return logger
}
