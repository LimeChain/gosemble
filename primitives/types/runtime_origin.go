package types

import (
	sc "github.com/LimeChain/goscale"
)

const (
	RawOriginRoot sc.U8 = iota
	RawOriginSigned
	RawOriginNone
)

type RawOrigin struct {
	sc.VaryingData // [T AccountId]
}

func NewRawOriginRoot() RawOrigin {
	return RawOrigin{sc.NewVaryingData(RawOriginRoot)}
}

func NewRawOriginSigned(account Address32) RawOrigin {
	return RawOrigin{sc.NewVaryingData(RawOriginSigned, account)}
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

func (o RawOrigin) IsRootOrigin() sc.Bool {
	return o.VaryingData[0] == RawOriginRoot
}

func (o RawOrigin) IsSignedOrigin() sc.Bool {
	return o.VaryingData[0] == RawOriginSigned
}

func (o RawOrigin) IsNoneOrigin() sc.Bool {
	return o.VaryingData[0] == RawOriginNone
}

type RuntimeOrigin = RawOrigin
