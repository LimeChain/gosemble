//go:build !nonwasmenv

package log

import (
	"fmt"

	"github.com/LimeChain/gosemble/env"
	"github.com/LimeChain/gosemble/utils"
)

const (
	CriticalLevel = iota
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

const target = "runtime"

type TraceLogger interface {
	Trace(message string)
	Tracef(message string, a ...any)
}

type DebugLogger interface {
	TraceLogger
	Debug(message string)
	Debugf(message string, a ...any)
}

type WarnLogger interface {
	DebugLogger
	Warn(message string)
	Warnf(message string, a ...any)
}

type Logger struct {
	memUtils utils.WasmMemoryTranslator
}

func NewLogger() Logger {
	return Logger{
		memUtils: utils.NewMemoryTranslator(),
	}
}

func (l Logger) Critical(message string) {
	l.log(CriticalLevel, []byte(target), []byte(message))
	panic(message)
}

func (l Logger) Criticalf(message string, a ...any) {
	l.Critical(fmt.Sprintf(message, a...))
}

func (l Logger) Warn(message string) {
	l.log(WarnLevel, []byte(target), []byte(message))
}

func (l Logger) Warnf(message string, a ...any) {
	l.Warn(fmt.Sprintf(message, a...))
}

func (l Logger) Info(message string) {
	l.log(InfoLevel, []byte(target), []byte(message))
}

func (l Logger) Infof(message string, a ...any) {
	l.Info(fmt.Sprintf(message, a...))
}

func (l Logger) Debug(message string) {
	l.log(DebugLevel, []byte(target), []byte(message))
}

func (l Logger) Debugf(message string, a ...any) {
	l.Debug(fmt.Sprintf(message, a...))
}

func (l Logger) Trace(message string) {
	l.log(TraceLevel, []byte(target), []byte(message))
}

func (l Logger) Tracef(message string, a ...any) {
	l.Trace(fmt.Sprintf(message, a...))
}

func (l Logger) log(level int32, target []byte, message []byte) {
	targetOffsetSize := l.memUtils.BytesToOffsetAndSize(target)
	messageOffsetSize := l.memUtils.BytesToOffsetAndSize(message)
	env.ExtLoggingLogVersion1(level, targetOffsetSize, messageOffsetSize)
}
