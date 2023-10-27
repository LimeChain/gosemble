package types

import (
	"bytes"
	"testing"

	"github.com/ChainSafe/gossamer/lib/common"
	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	blake2bHash = common.MustHexToHash("0x37373737373737373737373737373737373737373737373737373737373737")

	blake2bHashSequence = sc.BytesToFixedSequenceU8(blake2bHash.ToBytes())

	expectedBlake2bHash = Blake2bHash{blake2bHashSequence}
)

func Test_NewBlake2bHash(t *testing.T) {
	result := NewBlake2bHash(blake2bHashSequence...)

	assert.Equal(t, expectedBlake2bHash, result)
}

func Test_NewBlake2bHash_InvalidLength(t *testing.T) {
	assert.PanicsWithValue(t, "Blake2bHash should be of size 32", func() {
		NewBlake2bHash(blake2bHashSequence[1:32]...)
	})
}

func Test_Blake2bHash_Encode(t *testing.T) {
	buf := &bytes.Buffer{}
	expectedBlake2bHash.Encode(buf)

	assert.Equal(t, sc.FixedSequenceU8ToBytes(blake2bHashSequence), buf.Bytes())
}

func Test_Blake2bHash_Decode(t *testing.T) {
	buffer := bytes.NewBuffer(sc.FixedSequenceU8ToBytes(blake2bHashSequence))
	result := DecodeBlake2bHash(buffer)

	assert.Equal(t, expectedBlake2bHash, result)
}

func Test_Blake2bHash_Bytes(t *testing.T) {
	assert.Equal(t, sc.FixedSequenceU8ToBytes(blake2bHashSequence), expectedBlake2bHash.Bytes())
}
