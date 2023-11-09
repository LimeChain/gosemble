package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type DigestItem struct {
	Engine  sc.FixedSequence[sc.U8]
	Payload sc.Sequence[sc.U8]
}

func (di DigestItem) Encode(buffer *bytes.Buffer) error {
	return sc.EncodeEach(buffer,
		di.Engine,
		di.Payload,
	)
}

func (di DigestItem) Bytes() []byte {
	return sc.EncodedBytes(di)
}

func DecodeDigestItem(buffer *bytes.Buffer) (DigestItem, error) {
	engine, err := sc.DecodeFixedSequence[sc.U8](4, buffer)
	if err != nil {
		return DigestItem{}, err
	}
	payload, err := sc.DecodeSequence[sc.U8](buffer)
	if err != nil {
		return DigestItem{}, err
	}
	return DigestItem{
		Engine:  engine,
		Payload: payload,
	}, nil
}
