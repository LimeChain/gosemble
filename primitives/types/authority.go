package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type Authority struct {
	Id     AccountId[PublicKey]
	Weight sc.U64
}

func (a Authority) Encode(buffer *bytes.Buffer) error {
	return sc.EncodeEach(buffer,
		a.Id,
		a.Weight,
	)
}

func DecodeAuthority[T PublicKey](buffer *bytes.Buffer) (Authority, error) {
	pk, err := DecodeAccountId[T](buffer)
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
