package types

import (
	"bytes"
	"io"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	targetEd25519PublicKey = Ed25519PublicKey{
		FixedSequence: sc.BytesToFixedSequenceU8(pubKeyEd25519),
	}
	invalidAddress = []byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0}
)

func Test_PublicKey_Ed25519_Encode(t *testing.T) {
	buffer := &bytes.Buffer{}

	targetEd25519PublicKey.Encode(buffer)

	assert.Equal(t, pubKeyEd25519, buffer.Bytes())
}

func Test_PublicKey_Ed25519_Bytes(t *testing.T) {
	assert.Equal(t, pubKeyEd25519, targetEd25519PublicKey.Bytes())
}

func Test_DecodeEd25519_PublicKey(t *testing.T) {
	buffer := bytes.NewBuffer(pubKeyEd25519)

	result, err := DecodeEd25519PublicKey(buffer)
	assert.NoError(t, err)

	assert.Equal(t, targetEd25519PublicKey, result)
}

func Test_DecodeEd25519_PublicKey_InvalidNumberOfBytes(t *testing.T) {
	buffer := bytes.NewBuffer(invalidAddress)

	_, err := DecodeEd25519PublicKey(buffer)
	assert.Error(t, err)
	assert.Equal(t, io.EOF, err)
}

func Test_PublicKey_Ed25519_New(t *testing.T) {
	newEd25519PublicKey, err := NewEd25519PublicKey(sc.BytesToSequenceU8(pubKeyEd25519)...)
	assert.Nil(t, err)
	assert.Equal(t, newEd25519PublicKey, targetEd25519PublicKey)
}

func Test_PublicKey_Ed25519_New_InvalidAddress(t *testing.T) {
	expectedErr := newTypeError("Ed25519PublicKey")
	newEd25519PublicKey, err := NewEd25519PublicKey(sc.BytesToSequenceU8(invalidAddress)...)
	assert.Error(t, err)
	assert.Equal(t, Ed25519PublicKey{}, newEd25519PublicKey)
	assert.Equal(t, expectedErr, err)
}
