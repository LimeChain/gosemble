package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/utils"
)

type Authority struct {
	Id     PublicKey
	Weight sc.U64
}

func (a Authority) Encode(buffer *bytes.Buffer) error {
	return utils.EncodeEach(buffer,
		a.Id,
		a.Weight,
	)
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
