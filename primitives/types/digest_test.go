package types

import (
	"bytes"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	engineId = [4]byte{2, 3, 4, 5}
	n        = uint64(1)

	digestItem = DigestItem{
		Engine:  sc.BytesToFixedSequenceU8(engineId[:]),
		Payload: sc.BytesToSequenceU8(sc.U64(n).Bytes()),
	}
)

func Test_Decode_DigestTypeConsensusMessage(t *testing.T) {
	targetDigest := Digest{}
	targetDigest[DigestTypeConsensusMessage] = append(targetDigest[DigestTypeConsensusMessage], digestItem)

	buf := &bytes.Buffer{}
	err := targetDigest.Encode(buf)
	assert.NoError(t, err)

	digest, err := DecodeDigest(buf)
	assert.NoError(t, err)
	assert.Equal(t, targetDigest, digest)
}

func Test_Decode_DigestTypeSeal(t *testing.T) {
	targetDigest := Digest{}
	targetDigest[DigestTypeSeal] = append(targetDigest[DigestTypeSeal], digestItem)

	buf := &bytes.Buffer{}
	err := targetDigest.Encode(buf)
	assert.NoError(t, err)

	digest, err := DecodeDigest(buf)
	assert.NoError(t, err)
	assert.Equal(t, targetDigest, digest)
}

func Test_Decode_DigestTypePreRuntime(t *testing.T) {
	targetDigest := Digest{}
	targetDigest[DigestTypePreRuntime] = append(targetDigest[DigestTypePreRuntime], digestItem)

	buf := &bytes.Buffer{}
	err := targetDigest.Encode(buf)
	assert.NoError(t, err)

	digest, err := DecodeDigest(buf)
	assert.NoError(t, err)
	assert.Equal(t, targetDigest, digest)
}

func Test_Decode_DigestTypeRuntimeEnvironmentUpgraded(t *testing.T) {
	targetDigest := Digest{}

	buf := &bytes.Buffer{}
	err := targetDigest.Encode(buf)
	assert.NoError(t, err)

	digest, err := DecodeDigest(buf)
	assert.NoError(t, err)
	assert.Equal(t, targetDigest, digest)
}
