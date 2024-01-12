package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

const (
	RawOriginRoot sc.U8 = iota
	RawOriginSigned
	RawOriginNone
)

type RawOrigin struct {
	sc.VaryingData
}

func NewRawOriginRoot() RawOrigin {
	return RawOrigin{sc.NewVaryingData(RawOriginRoot)}
}

func NewRawOriginSigned(address AccountId) RawOrigin {
	return RawOrigin{sc.NewVaryingData(RawOriginSigned, address)}
}

func NewRawOriginNone() RawOrigin {
	return RawOrigin{sc.NewVaryingData(RawOriginNone)}
}

func RawOriginFrom(a sc.Option[AccountId]) RawOrigin {
	if a.HasValue {
		return NewRawOriginSigned(a.Value)
	} else {
		return NewRawOriginNone()
	}
}

func (o RawOrigin) IsRootOrigin() bool {
	return o.VaryingData[0] == RawOriginRoot
}

func (o RawOrigin) IsSignedOrigin() bool {
	return o.VaryingData[0] == RawOriginSigned
}

func (o RawOrigin) IsNoneOrigin() bool {
	return o.VaryingData[0] == RawOriginNone
}

func (o RawOrigin) AsSigned() (AccountId, error) {
	if !o.IsSignedOrigin() {
		return AccountId{}, newTypeError("RawOrigin")
	}

	return o.VaryingData[1].(AccountId), nil
}

func DecodeRawOrigin(buffer *bytes.Buffer) (RawOrigin, error) {
	b, err := sc.DecodeU8(buffer)
	if err != nil {
		return RawOrigin{}, err
	}

	switch b {
	case RawOriginRoot:
		return NewRawOriginRoot(), nil
	case RawOriginSigned:
		address, err := DecodeAccountId(buffer)
		if err != nil {
			return RawOrigin{}, err
		}
		return NewRawOriginSigned(address), nil
	case RawOriginNone:
		return NewRawOriginNone(), nil
	default:
		return RawOrigin{}, newTypeError("RawOrigin")
	}
}
