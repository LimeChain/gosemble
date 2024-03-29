package types

import (
	"bytes"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	hash512Sequence = sc.BytesToFixedSequenceU8(append(hash256.ToBytes(), hash256.ToBytes()...))

	expectedH512Hash = H512{hash512Sequence}
)

func Test_NewHash512(t *testing.T) {
	result, err := NewH512(hash512Sequence...)

	assert.NoError(t, err)
	assert.Equal(t, expectedH512Hash, result)
}

func Test_NewHash512_InvalidLength(t *testing.T) {
	result, err := NewH512(hash512Sequence[1:63]...)

	assert.Error(t, err)
	assert.Equal(t, "H512 should be of size 64", err.Error())
	assert.Equal(t, H512{}, result)
}

func Test_Hash512_Encode(t *testing.T) {
	buf := &bytes.Buffer{}

	err := expectedH512Hash.Encode(buf)

	assert.NoError(t, err)
	assert.Equal(t, sc.FixedSequenceU8ToBytes(hash512Sequence), buf.Bytes())
}

func Test_Hash512_Decode(t *testing.T) {
	buffer := bytes.NewBuffer(sc.FixedSequenceU8ToBytes(hash512Sequence))
	result, err := DecodeH512(buffer)
	assert.NoError(t, err)

	assert.Equal(t, expectedH512Hash, result)
}

func Test_Hash512_Bytes(t *testing.T) {
	assert.Equal(t, sc.FixedSequenceU8ToBytes(hash512Sequence), expectedH512Hash.Bytes())
}
