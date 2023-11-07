package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type Authority struct {
	Id     PublicKey
	Weight sc.U64
}

func (a Authority) Encode(buffer *bytes.Buffer) error {
	err := a.Id.Encode(buffer)
	if err != nil {
		return err
	}
	return a.Weight.Encode(buffer)
}

func DecodeAuthority(buffer *bytes.Buffer) (Authority, error) {
	pk, err := DecodePublicKey(buffer)
	if err != nil {
		return Authority{}, err
	}
	weight, err := sc.DecodeU64(buffer)
	if err != nil {
		return Authority{}, err
	}
	return Authority{
		Id:     pk,
		Weight: weight,
	}, nil
}

func (a Authority) Bytes() []byte {
	return sc.EncodedBytes(a)
}
