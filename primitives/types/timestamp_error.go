package types

import (
	sc "github.com/LimeChain/goscale"
)

const (
	TimestampErrorTooEarly sc.U8 = iota
	TimestampErrorTooFarInFuture
	TimestampErrorInvalid
)

type TimestampError struct {
	sc.VaryingData
}

func NewTimestampErrorTooEarly() TimestampError {
	return TimestampError{sc.NewVaryingData(TimestampErrorTooEarly)}
}

func NewTimestampErrorTooFarInFuture() TimestampError {
	return TimestampError{sc.NewVaryingData(TimestampErrorTooFarInFuture)}
}

func NewTimestampErrorInvalid() TimestampError {
	return TimestampError{sc.NewVaryingData(TimestampErrorInvalid)}
}

func (te TimestampError) IsFatal() sc.Bool {
	switch te.VaryingData[0] {
	case TimestampErrorTooEarly, TimestampErrorTooFarInFuture, TimestampErrorInvalid:
		return true
	default:
		return false
	}
}

func (te TimestampError) Error() string {
	switch te.VaryingData[0] {
	case TimestampErrorTooEarly:
		return "The time since the last timestamp is lower than the minimum period."
	case TimestampErrorTooFarInFuture:
		return "The timestamp of the block is too far in the future."
	case TimestampErrorInvalid:
		return "invalid inherent check for timestamp module"
	default:
		return NewTypeError("TimestampError").Error()
	}
}
