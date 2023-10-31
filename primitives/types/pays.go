package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
)

const (
	// PaysYes Transactor will pay related fees.
	PaysYes sc.U8 = iota

	// PaysNo Transactor will NOT pay related fees.
	PaysNo
)

type Pays = sc.VaryingData

func NewPaysYes() Pays {
	return sc.NewVaryingData(PaysYes)
}

func NewPaysNo() Pays {
	return sc.NewVaryingData(PaysNo)
}

func DecodePays(buffer *bytes.Buffer) (Pays, error) {
	b, err := sc.DecodeU8(buffer)
	if err != nil {
		return Pays{}, err
	}

	switch b {
	case PaysYes:
		return NewPaysYes(), nil
	case PaysNo:
		return NewPaysNo(), nil
	default:
		log.Critical("invalid Pays type")
	}

	panic("unreachable")
}
