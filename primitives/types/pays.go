package types

import (
	"bytes"
	"reflect"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
)

const (
	// Transactor will pay related fees.
	PaysYes = PaysValue(iota)

	// Transactor will NOT pay related fees.
	PaysNo
)

type PaysValue sc.U8

func (pv PaysValue) Encode(buffer *bytes.Buffer) {
	sc.U8(pv).Encode(buffer)
}

func DecodePaysValue(buffer *bytes.Buffer) PaysValue {
	return PaysValue(sc.DecodeU8(buffer))
}

func (pv PaysValue) Bytes() []byte {
	return sc.EncodedBytes(pv)
}

type Pays sc.VaryingData

func NewPays(value PaysValue) Pays {
	return Pays(sc.NewVaryingData(value))
}

func (p Pays) Encode(buffer *bytes.Buffer) {
	switch reflect.TypeOf(p) {
	case reflect.TypeOf(NewPays(PaysYes)):
		sc.U8(0).Encode(buffer)
	case reflect.TypeOf(NewPays(PaysNo)):
		sc.U8(1).Encode(buffer)
	default:
		log.Critical("invalid Pays type")
	}
}

func DecodePays(buffer *bytes.Buffer) Pays {
	b := sc.DecodeU8(buffer)

	switch b {
	case 0:
		return NewPays(PaysYes)
	case 1:
		return NewPays(PaysNo)
	default:
		log.Critical("invalid Pays type")
	}

	panic("unreachable")
}

func (p Pays) Bytes() []byte {
	return sc.EncodedBytes(p)
}
