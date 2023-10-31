package types

import (
	"bytes"
	"testing"

	"github.com/ChainSafe/gossamer/lib/common"
	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	hash256 = common.MustHexToHash("0x37373737373737373737373737373737373737373737373737373737373737")

	hash256Sequence = sc.BytesToFixedSequenceU8(hash256.ToBytes())

	expectedH256Hash = H256{hash256Sequence}
)

func Test_NewHash256(t *testing.T) {
	result := NewH256(hash256Sequence...)

	assert.Equal(t, expectedH256Hash, result)
}

func Test_NewHash256_InvalidLength(t *testing.T) {
	assert.PanicsWithValue(t, "H256 should be of size 32", func() {
		NewH256(hash256Sequence[1:32]...)
	})
}

func Test_Hash256_Encode(t *testing.T) {
	buf := &bytes.Buffer{}
	expectedH256Hash.Encode(buf)

	assert.Equal(t, sc.FixedSequenceU8ToBytes(hash256Sequence), buf.Bytes())
}

func Test_Hash256_Decode(t *testing.T) {
	buffer := bytes.NewBuffer(sc.FixedSequenceU8ToBytes(hash256Sequence))
	result, err := DecodeH256(buffer)
	assert.Nil(t, err)

	assert.Equal(t, expectedH256Hash, result)
}

func Test_Hash256_Bytes(t *testing.T) {
	assert.Equal(t, sc.FixedSequenceU8ToBytes(hash256Sequence), expectedH256Hash.Bytes())
}
