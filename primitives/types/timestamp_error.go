package types

import (
	"bytes"
	sc "github.com/LimeChain/goscale"
)

const (
	TimestampError_ValidAtTimestamp = sc.U8(iota)
	TimestampError_TooFarInFuture
)

type TimestampError sc.VaryingData

func NewTimestampError(values ...sc.Encodable) TimestampError {
	switch values[0] {
	case TimestampError_ValidAtTimestamp:
		return TimestampError{sc.NewVaryingData(values...)}
	case TimestampError_TooFarInFuture:
		return TimestampError{sc.NewVaryingData(values[0])}
	default:
		panic("invalid InherentError")
	}
}

func (te TimestampError) Encode(buffer *bytes.Buffer) {
	switch te[0] {
	case TimestampError_TooFarInFuture:
		te[0].Encode(buffer)
	case TimestampError_ValidAtTimestamp:
		te[0].Encode(buffer)
		te[1].Encode(buffer)
	default:
		panic("invalid TimestampError type")
	}
}

func (te TimestampError) Bytes() []byte {
	return sc.EncodedBytes(te)
}

func (te TimestampError) IsFatal() sc.Bool {
	switch te[0] {
	case TimestampError_TooFarInFuture:
		return true
	default:
		return false
	}
}
