//go:build !nonwasmenv

package log

import (
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

func Critical(target string, message string) {
	log(CriticalLevel, []byte(target), []byte(message))
}

func Warn(target string, message string) {
	log(WarnLevel, []byte(target), []byte(message))
}

func Info(target string, message string) {
	log(InfoLevel, []byte(target), []byte(message))
}

func Debug(target string, message string) {
	log(DebugLevel, []byte(target), []byte(message))
}

func Trace(target string, message string) {
	log(TraceLevel, []byte(target), []byte(message))
}

func log(level int32, target []byte, message []byte) {
	targetOffsetSize := utils.BytesToOffsetAndSize(target)
	messageOffsetSize := utils.BytesToOffsetAndSize(message)
	env.ExtLoggingLogVersion1(level, targetOffsetSize, messageOffsetSize)
}
