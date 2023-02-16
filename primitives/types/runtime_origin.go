package types

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
)

// The system itself ordained this dispatch to happen: this is the highest privilege level.
type RootOrigin struct {
	sc.Empty
}

// It is signed by some public key and we provide the `AccountId`.
type SignedOrigin struct {
	sc.Empty
	AccountId Address32
}

// It is signed by nobody, can be either:
// * included and agreed upon by the validators anyway,
// * or unsigned transaction validated by a pallet.
type NoneOrigin struct {
	sc.Empty
}

type RawOrigin sc.VaryingData // [T AccountId]

func NewRawOrigin(value sc.Encodable) RawOrigin {
	switch value.(type) {
	case RootOrigin, SignedOrigin, NoneOrigin:
		return RawOrigin(sc.NewVaryingData(value))
	default:
		log.Critical("invalid RawOrigin type")
	}

	panic("unreachable")
}

func RawOriginFrom(a sc.Option[Address32]) RawOrigin {
	if a.HasValue {
		return NewRawOrigin(SignedOrigin{AccountId: a.Value})
	} else {
		return NewRawOrigin(NoneOrigin{})
	}
}

type RuntimeOrigin = RawOrigin
