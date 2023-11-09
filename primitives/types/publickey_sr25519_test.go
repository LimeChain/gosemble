package types

import (
	"bytes"
	"io"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	targetSr25519PublicKey = Sr25519PublicKey{
		FixedSequence: sc.BytesToFixedSequenceU8(pubKeySr25519),
	}
)

func Test_PublicKey_Sr25519_Encode(t *testing.T) {
	buffer := &bytes.Buffer{}

	targetSr25519PublicKey.Encode(buffer)

	assert.Equal(t, pubKeySr25519, buffer.Bytes())
}

func Test_PublicKey_Sr25519_Bytes(t *testing.T) {
	assert.Equal(t, pubKeySr25519, targetSr25519PublicKey.Bytes())
}

func Test_DecodeSr25519_PublicKey(t *testing.T) {
	buffer := bytes.NewBuffer(pubKeySr25519)

	result, err := DecodeSr25519PublicKey(buffer)
	assert.NoError(t, err)

	assert.Equal(t, targetSr25519PublicKey, result)
}

func Test_DecodeSr25519_PublicKey_InvalidNumberOfBytes(t *testing.T) {
	buffer := bytes.NewBuffer(invalidAddress)

	_, err := DecodeSr25519PublicKey(buffer)
	assert.Error(t, err)
	assert.Equal(t, io.EOF, err)
}

func Test_PublicKey_Sr25519_New(t *testing.T) {
	newSr25519PublicKey, err := NewSr25519PublicKey(sc.BytesToSequenceU8(pubKeySr25519)...)
	assert.Nil(t, err)
	assert.Equal(t, newSr25519PublicKey, targetSr25519PublicKey)
}

func Test_PublicKey_Sr25519__New_InvalidAddress(t *testing.T) {
	expectedErr := newTypeError("Sr25519PublicKey")
	newSr25519PublicKey, err := NewSr25519PublicKey(sc.BytesToSequenceU8(invalidAddress)...)
	assert.Error(t, err)
	assert.Equal(t, Sr25519PublicKey{}, newSr25519PublicKey)
	assert.Equal(t, expectedErr, err)
}
