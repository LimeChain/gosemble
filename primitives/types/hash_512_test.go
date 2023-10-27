package types

import (
	"bytes"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	hash512Sequence  = sc.BytesToFixedSequenceU8(append(hash256.ToBytes(), hash256.ToBytes()...))
	expectedH512Hash = H512{hash512Sequence}
)

func Test_NewHash512(t *testing.T) {
	result := NewH512(hash512Sequence...)

	assert.Equal(t, expectedH512Hash, result)
}

func Test_NewHash512_InvalidLength(t *testing.T) {
	assert.PanicsWithValue(t, "H512 should be of size 64", func() {
		NewH512(hash512Sequence[1:63]...)
	})
}

func Test_Hash512_Encode(t *testing.T) {
	buf := &bytes.Buffer{}
	expectedH512Hash.Encode(buf)

	assert.Equal(t, sc.FixedSequenceU8ToBytes(hash512Sequence), buf.Bytes())
}

func Test_Hash512_Decode(t *testing.T) {
	buffer := bytes.NewBuffer(sc.FixedSequenceU8ToBytes(hash512Sequence))
	result := DecodeH512(buffer)

	assert.Equal(t, expectedH512Hash, result)
}

func Test_Hash512_Bytes(t *testing.T) {
	assert.Equal(t, sc.FixedSequenceU8ToBytes(hash512Sequence), expectedH512Hash.Bytes())
}
