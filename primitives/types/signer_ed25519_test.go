package types

import (
	"bytes"
	"io"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	targetEd25519Signer = NewEd25519Signer(sc.BytesToSequenceU8(pubKeyEd25519Signer)...)
	invalidAddress      = []byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0}
)

func Test_Signer_Ed25519_Encode(t *testing.T) {
	buffer := &bytes.Buffer{}

	targetEd25519Signer.Encode(buffer)

	assert.Equal(t, pubKeyEd25519Signer, buffer.Bytes())
}

func Test_Signer_Ed25519_Bytes(t *testing.T) {
	assert.Equal(t, pubKeyEd25519Signer, targetEd25519Signer.Bytes())
}

func Test_DecodeEd25519_Signer(t *testing.T) {
	buffer := bytes.NewBuffer(pubKeyEd25519Signer)

	result, err := DecodeEd25519Signer(buffer)
	assert.NoError(t, err)

	assert.Equal(t, targetEd25519Signer, result)
}

func Test_DecodeEd25519_Signer_InvalidNumberOfBytes(t *testing.T) {
	buffer := bytes.NewBuffer(invalidAddress)

	_, err := DecodeEd25519Signer(buffer)
	assert.Error(t, err)
	assert.Equal(t, io.EOF, err)
}

func Test_Signer_Ed25519_New(t *testing.T) {
	newEd25519Signer := NewEd25519Signer(sc.BytesToSequenceU8(pubKeyEd25519Signer)...)
	assert.Equal(t, newEd25519Signer, targetEd25519Signer)
}

func Test_Signer_Ed25519_New_InvalidAddress(t *testing.T) {
	assert.PanicsWithValue(t,
		"Ed25519Signer should be of size 32",
		func() {
			NewEd25519Signer(sc.BytesToSequenceU8(invalidAddress)...)
		})
}
