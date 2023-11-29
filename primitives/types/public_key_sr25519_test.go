package types

import (
	"bytes"
	"io"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	bytesPublicKeySr25519  = []byte{1, 1, 0, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0}
	targetSr25519PublicKey = Sr25519PublicKey{
		FixedSequence: sc.BytesToFixedSequenceU8(bytesPublicKeySr25519),
	}
)

func Test_NewSr25519PublicKey(t *testing.T) {
	newSr25519PublicKey, err := NewSr25519PublicKey(sc.BytesToSequenceU8(bytesPublicKeySr25519)...)
	assert.Nil(t, err)
	assert.Equal(t, newSr25519PublicKey, targetSr25519PublicKey)
}

func Test_NewSr25519PublicKey_InvalidAddress(t *testing.T) {
	expectedErr := newTypeError("Sr25519PublicKey")
	newSr25519PublicKey, err := NewSr25519PublicKey(sc.BytesToSequenceU8(invalidAddress)...)
	assert.Error(t, err)
	assert.Equal(t, Sr25519PublicKey{}, newSr25519PublicKey)
	assert.Equal(t, expectedErr, err)
}

func Test_Sr25519PublicKey_Encode(t *testing.T) {
	buffer := &bytes.Buffer{}

	err := targetSr25519PublicKey.Encode(buffer)
	assert.Nil(t, err)

	assert.Equal(t, bytesPublicKeySr25519, buffer.Bytes())
}

func Test_Sr25519PublicKey_Bytes(t *testing.T) {
	assert.Equal(t, bytesPublicKeySr25519, targetSr25519PublicKey.Bytes())
}

func Test_DecodeSr25519PublicKey(t *testing.T) {
	buffer := bytes.NewBuffer(bytesPublicKeySr25519)

	result, err := DecodeSr25519PublicKey(buffer)
	assert.NoError(t, err)

	assert.Equal(t, targetSr25519PublicKey, result)
}

func Test_DecodeSr25519PublicKey_InvalidNumberOfBytes(t *testing.T) {
	buffer := bytes.NewBuffer(invalidAddress)

	_, err := DecodeSr25519PublicKey(buffer)
	assert.Error(t, err)
	assert.Equal(t, io.EOF, err)
}

func Test_Sr25519PublicKey_SignatureType(t *testing.T) {
	assert.Equal(t, PublicKeySr25519, targetSr25519PublicKey.SignatureType())
}
