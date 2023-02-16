//go:build !nonwasmenv

package log

import (
	"github.com/LimeChain/gosemble/env"
	"github.com/LimeChain/gosemble/utils"
)

const (
	Critical = iota
	Warn
	Info
	Debug
	Trace
)

func Log(level int32, target []byte, message []byte) {
	targetOffsetSize := utils.BytesToOffsetAndSize(target)
	messageOffsetSize := utils.BytesToOffsetAndSize(message)
	env.ExtLoggingLogVersion1(level, targetOffsetSize, messageOffsetSize)
}
