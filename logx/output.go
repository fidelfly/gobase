package logx

import (
	"io"

	"github.com/natefinch/lumberjack"
)

func RotateLog(filename string, maxSize int, maxBackup, maxAge int, compress bool) io.Writer {
	if maxSize == 0 {
		maxSize = 1
	}
	if maxBackup == 0 {
		maxBackup = 7
	}
	if maxAge == 0 {
		maxAge = 7
	}
	return &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
		Compress:   compress,
		LocalTime:  true,
	}
}
