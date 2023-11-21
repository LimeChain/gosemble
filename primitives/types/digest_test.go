package types

import (
	"bytes"
	"encoding/hex"
	"io"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	engineId          = sc.BytesToFixedSequenceU8([]byte{2, 3, 4, 5})
	digestItemMessage = sc.BytesToSequenceU8([]byte{'a', 'b', 'c'})

	digestItemPreRuntime = NewDigestItemPreRuntime(
		engineId,
		digestItemMessage,
	)
	digestItemSeal = NewDigestItemSeal(engineId, digestItemMessage)
)

var (
	expectDigest = NewDigest(sc.Sequence[DigestItem]{
		digestItemPreRuntime,
		digestItemSeal,
	})
	expectBytesDigest, _ = hex.DecodeString("0806020304050c61626305020304050c616263")
)

func Test_DecodeDigest(t *testing.T) {
	buffer := bytes.NewBuffer(expectBytesDigest)

	result, err := DecodeDigest(buffer)

	assert.NoError(t, err)
	assert.Equal(t, expectDigest, result)
}

func Test_DecodeDigest_Fails_Compact(t *testing.T) {
	buffer := &bytes.Buffer{}

	result, err := DecodeDigest(buffer)

	assert.Equal(t, io.EOF, err)
	assert.Equal(t, Digest{}, result)
}

func Test_DecodeDigest_Fails_DigestItem(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(4)

	result, err := DecodeDigest(buffer)

	assert.Equal(t, io.EOF, err)
	assert.Equal(t, Digest{}, result)
}

func Test_DecodeDigest_PreRuntimes(t *testing.T) {
	expect := sc.Sequence[DigestPreRuntime]{
		NewDigestPreRuntime(engineId, digestItemMessage),
	}

	target := NewDigest(sc.Sequence[DigestItem]{
		digestItemSeal,
		NewDigestItemRuntimeEnvironmentUpgrade(),
		NewDigestItemOther(digestItemMessage),
		digestItemPreRuntime,
		NewDigestItemOther(digestItemMessage),
	})

	result, err := target.PreRuntimes()

	assert.NoError(t, err)
	assert.Equal(t, expect, result)
}

func Test_DecodeDigest_OnlyPreRuntimes(t *testing.T) {
	expect := NewDigest(sc.Sequence[DigestItem]{digestItemPreRuntime})

	target := NewDigest(sc.Sequence[DigestItem]{
		digestItemSeal,
		NewDigestItemRuntimeEnvironmentUpgrade(),
		NewDigestItemOther(digestItemMessage),
		digestItemPreRuntime,
		NewDigestItemOther(digestItemMessage),
	})

	result := target.OnlyPreRuntimes()

	assert.Equal(t, expect, result)
}
