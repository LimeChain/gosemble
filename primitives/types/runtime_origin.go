package types

import (
	sc "github.com/LimeChain/goscale"
)

const (
	RawOriginRoot sc.U8 = iota
	RawOriginSigned
	RawOriginNone
)

type RawOrigin = sc.VaryingData // [T AccountId]

func NewRawOriginRoot() RawOrigin {
	return sc.NewVaryingData(RawOriginRoot)
}

func NewRawOriginSigned(account Address32) RawOrigin {
	return sc.NewVaryingData(RawOriginSigned, account)
}

func NewRawOriginNone() RawOrigin {
	return sc.NewVaryingData(RawOriginNone)
}

func RawOriginFrom(a sc.Option[Address32]) RawOrigin {
	if a.HasValue {
		return NewRawOriginSigned(a.Value)
	} else {
		return NewRawOriginNone()
	}
}

type RuntimeOrigin = RawOrigin
