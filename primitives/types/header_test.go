package types

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/ChainSafe/gossamer/lib/common"
	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	expectBytesHeader, _ = hex.DecodeString("0f6d3477739f8a65886135f58c83ff7c2d4a8300a010dfc8b4c5d65ba37920bb04d9e8bf89bda43fb46914321c371add19b81ff92ad6923e8f189b52578074b073105165e71964828f2b8d1fd89904602cfb9b8930951d87eb249aa2d7c4b51ee7040661757261200000000000000000")

	parentHash     = common.MustHexToHash("0x0f6d3477739f8a65886135f58c83ff7c2d4a8300a010dfc8b4c5d65ba37920bb")
	stateRoot      = common.MustHexToHash("0xd9e8bf89bda43fb46914321c371add19b81ff92ad6923e8f189b52578074b073")
	extrinsicsRoot = common.MustHexToHash("0x105165e71964828f2b8d1fd89904602cfb9b8930951d87eb249aa2d7c4b51ee7")
	digest         = Digest{
		DigestTypePreRuntime: sc.FixedSequence[DigestItem]{
			DigestItem{
				Engine:  sc.BytesToFixedSequenceU8([]byte{'a', 'u', 'r', 'a'}),
				Payload: sc.BytesToSequenceU8(sc.U64(0).Bytes()),
			},
		},
	}

	targetHeader = Header{
		ParentHash:     NewBlake2bHash(sc.BytesToFixedSequenceU8(parentHash.ToBytes())...),
		Number:         1,
		StateRoot:      NewH256(sc.BytesToFixedSequenceU8(stateRoot.ToBytes())...),
		ExtrinsicsRoot: NewH256(sc.BytesToFixedSequenceU8(extrinsicsRoot.ToBytes())...),
		Digest:         digest,
	}
)

func Test_Header_Encode(t *testing.T) {
	buffer := &bytes.Buffer{}

	targetHeader.Encode(buffer)

	assert.Equal(t, expectBytesHeader, buffer.Bytes())
}

func Test_Header_Decode(t *testing.T) {
	buf := bytes.NewBuffer(expectBytesHeader)

	result := DecodeHeader(buf)

	assert.Equal(t, targetHeader.ParentHash, result.ParentHash)
	assert.Equal(t, targetHeader.Number, result.Number)
	assert.Equal(t, targetHeader.StateRoot, result.StateRoot)
	assert.Equal(t, targetHeader.ExtrinsicsRoot, result.ExtrinsicsRoot)
	assert.Equal(t, targetHeader.Digest, result.Digest)
}

func Test_Header_Bytes(t *testing.T) {
	assert.Equal(t, expectBytesHeader, targetHeader.Bytes())
}
