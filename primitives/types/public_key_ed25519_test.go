package types

import (
	"bytes"
	"io"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	bytesPublicKeyEd25519  = []byte{1, 1, 1, 1, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0}
	targetEd25519PublicKey = Ed25519PublicKey{
		FixedSequence: sc.BytesToFixedSequenceU8(bytesPublicKeyEd25519),
	}
	invalidAddress = []byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0}
)

func Test_NewEd25519PublicKey(t *testing.T) {
	newEd25519PublicKey, err := NewEd25519PublicKey(sc.BytesToSequenceU8(bytesPublicKeyEd25519)...)
	assert.Nil(t, err)
	assert.Equal(t, newEd25519PublicKey, targetEd25519PublicKey)
}

func Test_NewEd25519PublicKey_InvalidAddress(t *testing.T) {
	expectedErr := newTypeError("Ed25519PublicKey")
	newEd25519PublicKey, err := NewEd25519PublicKey(sc.BytesToSequenceU8(invalidAddress)...)
	assert.Error(t, err)
	assert.Equal(t, Ed25519PublicKey{}, newEd25519PublicKey)
	assert.Equal(t, expectedErr, err)
}

func Test_Ed25519PublicKey_Encode(t *testing.T) {
	buffer := &bytes.Buffer{}

	err := targetEd25519PublicKey.Encode(buffer)
	assert.Nil(t, err)

	assert.Equal(t, bytesPublicKeyEd25519, buffer.Bytes())
}

func Test_Ed25519PublicKey_Bytes(t *testing.T) {
	assert.Equal(t, bytesPublicKeyEd25519, targetEd25519PublicKey.Bytes())
}

func Test_DecodeEd25519PublicKey(t *testing.T) {
	buffer := bytes.NewBuffer(bytesPublicKeyEd25519)

	result, err := DecodeEd25519PublicKey(buffer)
	assert.NoError(t, err)

	assert.Equal(t, targetEd25519PublicKey, result)
}

func Test_DecodeEd25519PublicKey_InvalidNumberOfBytes(t *testing.T) {
	buffer := bytes.NewBuffer(invalidAddress)

	_, err := DecodeEd25519PublicKey(buffer)
	assert.Error(t, err)
	assert.Equal(t, io.EOF, err)
}

func Test_Ed25519PublicKey_SignatureType(t *testing.T) {
	assert.Equal(t, PublicKeyEd25519, targetEd25519PublicKey.SignatureType())
}
