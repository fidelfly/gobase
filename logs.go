package gox

import (
	"fmt"
)

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
