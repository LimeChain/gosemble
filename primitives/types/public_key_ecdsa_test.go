package types

import (
	"bytes"
	"io"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	pubKeyEcdsa          = []byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0}
	targetEcdsaPublicKey = EcdsaPublicKey{
		FixedSequence: sc.BytesToFixedSequenceU8(pubKeyEcdsa),
	}
)

func Test_NewEcdsaPublicKey(t *testing.T) {
	newEcdsaPublicKey, err := NewEcdsaPublicKey(sc.BytesToSequenceU8(pubKeyEcdsa)...)
	assert.Nil(t, err)
	assert.Equal(t, newEcdsaPublicKey, targetEcdsaPublicKey)
}

func Test_NewEcdsaPublicKey_InvalidAddress(t *testing.T) {
	expectedErr := newTypeError("EcdsaPublicKey")
	newEcdsaPublicKey, err := NewEcdsaPublicKey(sc.BytesToSequenceU8(invalidAddress)...)
	assert.Error(t, err)
	assert.Equal(t, EcdsaPublicKey{}, newEcdsaPublicKey)
	assert.Equal(t, expectedErr, err)
}

func Test_EcdsaPublicKey_Encode(t *testing.T) {
	buffer := &bytes.Buffer{}

	err := targetEcdsaPublicKey.Encode(buffer)
	assert.Nil(t, err)

	assert.Equal(t, pubKeyEcdsa, buffer.Bytes())
}

func Test_EcdsaPublicKey_Bytes(t *testing.T) {
	assert.Equal(t, pubKeyEcdsa, targetEcdsaPublicKey.Bytes())
}

func Test_DecodeEcdsaPublicKey(t *testing.T) {
	buffer := bytes.NewBuffer(pubKeyEcdsa)

	result, err := DecodeEcdsaPublicKey(buffer)
	assert.NoError(t, err)

	assert.Equal(t, targetEcdsaPublicKey, result)
}

func Test_DecodeEcdsaPublicKey_InvalidNumberOfBytes(t *testing.T) {
	buffer := bytes.NewBuffer(invalidAddress)

	_, err := DecodeEcdsaPublicKey(buffer)
	assert.Error(t, err)
	assert.Equal(t, io.EOF, err)
}

func Test_EcdsaPublicKey_SignatureType(t *testing.T) {
	assert.Equal(t, PublicKeyEcdsa, targetEcdsaPublicKey.SignatureType())
}
