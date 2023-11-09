package types

import (
	"bytes"
	"io"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	targetEcdsaSigner = EcdsaPublicKey{
		FixedSequence: sc.BytesToFixedSequenceU8(pubKeyEcdsaSigner),
	}
)

func Test_Signer_Ecdsa_Encode(t *testing.T) {
	buffer := &bytes.Buffer{}

	targetEcdsaSigner.Encode(buffer)

	assert.Equal(t, pubKeyEcdsaSigner, buffer.Bytes())
}

func Test_Signer_Ecdsa_Bytes(t *testing.T) {
	assert.Equal(t, pubKeyEcdsaSigner, targetEcdsaSigner.Bytes())
}

func Test_DecodeEcdsa_Signer(t *testing.T) {
	buffer := bytes.NewBuffer(pubKeyEcdsaSigner)

	result, err := DecodeEcdsaPublicKey(buffer)
	assert.NoError(t, err)

	assert.Equal(t, targetEcdsaSigner, result)
}

func Test_DecodeEcdsa_Signer_InvalidNumberOfBytes(t *testing.T) {
	buffer := bytes.NewBuffer(invalidAddress)

	_, err := DecodeEcdsaPublicKey(buffer)
	assert.Error(t, err)
	assert.Equal(t, io.EOF, err)
}

func Test_Signer_Ecdsa_New(t *testing.T) {
	newEcdsaSigner, err := NewEcdsaPublicKey(sc.BytesToSequenceU8(pubKeyEcdsaSigner)...)
	assert.Nil(t, err)
	assert.Equal(t, newEcdsaSigner, targetEcdsaSigner)
}

func Test_Signer_Ecdsa__New_InvalidAddress(t *testing.T) {
	expectedErr := newTypeError("EcdsaPublicKey")
	newEcdsaSigner, err := NewEcdsaPublicKey(sc.BytesToSequenceU8(invalidAddress)...)
	assert.Error(t, err)
	assert.Equal(t, EcdsaPublicKey{}, newEcdsaSigner)
	assert.Equal(t, expectedErr, err)
}
