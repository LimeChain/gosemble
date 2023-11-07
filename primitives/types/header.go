package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/utils"
)

type Header struct {
	ParentHash     Blake2bHash
	Number         sc.U64
	StateRoot      H256
	ExtrinsicsRoot H256
	Digest         Digest
}

func (h Header) Encode(buffer *bytes.Buffer) error {
	return utils.EncodeEach(buffer,
		h.ParentHash,
		sc.ToCompact(h.Number),
		h.StateRoot,
		h.ExtrinsicsRoot,
		h.Digest,
	)
}

func (h Header) Bytes() []byte {
	return sc.EncodedBytes(h)
}

func DecodeHeader(buffer *bytes.Buffer) (Header, error) {
	parentHash, err := DecodeBlake2bHash(buffer)
	if err != nil {
		return Header{}, err
	}
	blockNumber, err := sc.DecodeCompact(buffer)
	if err != nil {
		return Header{}, err
	}
	stateRoot, err := DecodeH256(buffer)
	if err != nil {
		return Header{}, err
	}
	extrinsicRoot, err := DecodeH256(buffer)
	if err != nil {
		return Header{}, err
	}
	digest, err := DecodeDigest(buffer)
	if err != nil {
		return Header{}, err
	}

	return Header{
		ParentHash:     parentHash,
		Number:         sc.U64(blockNumber.ToBigInt().Uint64()),
		StateRoot:      stateRoot,
		ExtrinsicsRoot: extrinsicRoot,
		Digest:         digest,
	}, nil
}
