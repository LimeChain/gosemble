package types

import (
	"bytes"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	hash256Bytes = []sc.U8{
		0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37,
		0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37,
		0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37,
		0x37, 0x37,
	}

	expectedH256Hash = H256{sc.FixedSequence[sc.U8](hash256Bytes)}

	hash512Bytes = []sc.U8{
		0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37,
		0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37,
		0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37,
		0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37,
		0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37,
		0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37,
		0x37, 0x37,
	}

	expectedH512Hash = H512{sc.FixedSequence[sc.U8](hash512Bytes)}

	blake2bHashBytes = []sc.U8{
		0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37,
		0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37,
		0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37,
		0x37, 0x37,
	}

	expectedBlake2bHash = Blake2bHash{sc.FixedSequence[sc.U8](blake2bHashBytes)}
)

func Test_NewHash256(t *testing.T) {
	result := NewH256(hash256Bytes...)

	assert.Equal(t, expectedH256Hash, result)
}

func Test_NewHash256_InvalidLength(t *testing.T) {
	assert.PanicsWithValue(t, "H256 should be of size 32", func() {
		NewH256(hash256Bytes[1:32]...)
	})
}

func Test_Hash256_Encode(t *testing.T) {
	buf := &bytes.Buffer{}
	expectedH256Hash.Encode(buf)

	assert.Equal(t, sc.FixedSequenceU8ToBytes(hash256Bytes), buf.Bytes())
}

func Test_Hash256_Decode(t *testing.T) {
	buffer := bytes.NewBuffer(sc.FixedSequenceU8ToBytes(hash256Bytes))
	result := DecodeH256(buffer)

	assert.Equal(t, expectedH256Hash, result)
}

func Test_Hash256_Bytes(t *testing.T) {
	assert.Equal(t, sc.FixedSequenceU8ToBytes(hash256Bytes), expectedH256Hash.Bytes())
}

func Test_NewHash512(t *testing.T) {
	result := NewH512(hash512Bytes...)

	assert.Equal(t, expectedH512Hash, result)
}

func Test_NewHash512_InvalidLength(t *testing.T) {
	assert.PanicsWithValue(t, "H512 should be of size 64", func() {
		NewH512(hash512Bytes[1:63]...)
	})
}

func Test_Hash512_Encode(t *testing.T) {
	buf := &bytes.Buffer{}
	expectedH512Hash.Encode(buf)

	assert.Equal(t, sc.FixedSequenceU8ToBytes(hash512Bytes), buf.Bytes())
}

func Test_Hash512_Decode(t *testing.T) {
	buffer := bytes.NewBuffer(sc.FixedSequenceU8ToBytes(hash512Bytes))
	result := DecodeH512(buffer)

	assert.Equal(t, expectedH512Hash, result)
}

func Test_Hash512_Bytes(t *testing.T) {
	assert.Equal(t, sc.FixedSequenceU8ToBytes(hash512Bytes), expectedH512Hash.Bytes())
}

func Test_NewBlake2bHash(t *testing.T) {
	result := NewBlake2bHash(blake2bHashBytes...)

	assert.Equal(t, expectedBlake2bHash, result)
}

func Test_NewBlake2bHash_InvalidLength(t *testing.T) {
	assert.PanicsWithValue(t, "Blake2bHash should be of size 32", func() {
		NewBlake2bHash(blake2bHashBytes[1:32]...)
	})
}

func Test_Blake2bHash_Encode(t *testing.T) {
	buf := &bytes.Buffer{}
	expectedBlake2bHash.Encode(buf)

	assert.Equal(t, sc.FixedSequenceU8ToBytes(blake2bHashBytes), buf.Bytes())
}

func Test_Blake2bHash_Decode(t *testing.T) {
	buffer := bytes.NewBuffer(sc.FixedSequenceU8ToBytes(blake2bHashBytes))
	result := DecodeBlake2bHash(buffer)

	assert.Equal(t, expectedBlake2bHash, result)
}

func Test_Blake2bHash_Bytes(t *testing.T) {
	assert.Equal(t, sc.FixedSequenceU8ToBytes(blake2bHashBytes), expectedBlake2bHash.Bytes())
}
