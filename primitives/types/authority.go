package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type Authority struct {
	Id     PublicKey
	Weight sc.U64
}

func (a Authority) Encode(buffer *bytes.Buffer) {
	a.Id.Encode(buffer)
	a.Weight.Encode(buffer)
}

func DecodeAuthority(buffer *bytes.Buffer) (Authority, error) {
	weight, err := sc.DecodeU64(buffer)
	if err != nil {
		return Authority{}, err
	}
	return Authority{
		Id:     DecodePublicKey(buffer),
		Weight: weight,
	}, nil
}

func (a Authority) Bytes() []byte {
	return sc.EncodedBytes(a)
}
