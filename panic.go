package gox

import (
	"github.com/fidelfly/gox/logx"
)

//export
// CapturePanicAndRecover
// Must be called within **defer** code fragment
func CapturePanicAndRecover(messages ...string) {
	if err := recover(); err != nil {
		if panicErr, ok := err.(error); ok {
			logx.Error(panicErr)
		}
		if len(messages) > 0 {
			logx.Panic(messages)
		}
		logx.Panic(err)
		logx.Info("Panic recovered")
	}
}
