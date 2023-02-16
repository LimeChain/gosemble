package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
)

const (
	TimestampErrorValidateTimestamp = sc.U8(iota)
	TimestampErrorTooFarInFuture
)

const (
	errInvalidTimestampType = "invalid TimestampError type"
)

type TimestampError sc.VaryingData

func NewTimestampError(values ...sc.Encodable) TimestampError {
	switch values[0] {
	case TimestampErrorValidateTimestamp:
		return TimestampError{sc.NewVaryingData(values...)}
	case TimestampErrorTooFarInFuture:
		return TimestampError{sc.NewVaryingData(values[0])}
	default:
		log.Critical(errInvalidTimestampType)
	}

	panic("unreachable")
}

func (te TimestampError) Encode(buffer *bytes.Buffer) {
	switch te[0] {
	case TimestampErrorTooFarInFuture:
		te[0].Encode(buffer)
	case TimestampErrorValidateTimestamp:
		te[0].Encode(buffer)
		te[1].Encode(buffer)
	default:
		log.Critical(errInvalidTimestampType)
	}
}

func (te TimestampError) Bytes() []byte {
	return sc.EncodedBytes(te)
}

func (te TimestampError) IsFatal() sc.Bool {
	switch te[0] {
	case TimestampErrorTooFarInFuture:
		return true
	default:
		return false
	}
}
