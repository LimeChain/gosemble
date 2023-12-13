package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type Pays sc.U8

const (
	// PaysYes Transactor will pay related fees.
	PaysYes Pays = iota

	// PaysNo Transactor will NOT pay related fees.
	PaysNo
)

func (p Pays) Encode(buffer *bytes.Buffer) error {
	return sc.U8(p).Encode(buffer)
}

func DecodePays(buffer *bytes.Buffer) (Pays, error) {
	b, err := sc.DecodeU8(buffer)
	if err != nil {
		return 0, err
	}

	switch Pays(b) {
	case PaysYes:
		return PaysYes, nil
	case PaysNo:
		return PaysNo, nil
	default:
		return 0, newTypeError("Pays")
	}
}

func (p Pays) Bytes() []byte {
	return sc.EncodedBytes(p)
}
