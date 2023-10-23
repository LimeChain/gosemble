package types

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
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

func NewRawOriginSigned(address Address32) RawOrigin {
	return RawOrigin{sc.NewVaryingData(RawOriginSigned, address)}
}

func NewRawOriginNone() RawOrigin {
	return RawOrigin{sc.NewVaryingData(RawOriginNone)}
}

func RawOriginFrom(a sc.Option[Address32]) RawOrigin {
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

func (o RawOrigin) AsSigned() Address32 {
	if !o.IsSignedOrigin() {
		log.Critical("not a signed origin")
	}

	return o.VaryingData[1].(Address32)
}
