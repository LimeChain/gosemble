package types

import (
	"bytes"
	"io"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	targetSr25519Signer = Sr25519PublicKey{
		FixedSequence: sc.BytesToFixedSequenceU8(pubKeySr25519Signer),
	}
)

func Test_Signer_Sr25519_Encode(t *testing.T) {
	buffer := &bytes.Buffer{}

	targetSr25519Signer.Encode(buffer)

	assert.Equal(t, pubKeySr25519Signer, buffer.Bytes())
}

func Test_Signer_Sr25519_Bytes(t *testing.T) {
	assert.Equal(t, pubKeySr25519Signer, targetSr25519Signer.Bytes())
}

func Test_DecodeSr25519_Signer(t *testing.T) {
	buffer := bytes.NewBuffer(pubKeySr25519Signer)

	result, err := DecodeSr25519PublicKey(buffer)
	assert.NoError(t, err)

	assert.Equal(t, targetSr25519Signer, result)
}

func Test_DecodeSr25519_Signer_InvalidNumberOfBytes(t *testing.T) {
	buffer := bytes.NewBuffer(invalidAddress)

	_, err := DecodeSr25519PublicKey(buffer)
	assert.Error(t, err)
	assert.Equal(t, io.EOF, err)
}

func Test_Signer_Sr25519_New(t *testing.T) {
	newSr25519Signer, err := NewSr25519PublicKey(sc.BytesToSequenceU8(pubKeySr25519Signer)...)
	assert.Nil(t, err)
	assert.Equal(t, newSr25519Signer, targetSr25519Signer)
}

func Test_Signer_Sr25519__New_InvalidAddress(t *testing.T) {
	expectedErr := newTypeError("Sr25519PublicKey")
	newSr25519Signer, err := NewSr25519PublicKey(sc.BytesToSequenceU8(invalidAddress)...)
	assert.Error(t, err)
	assert.Equal(t, Sr25519PublicKey{}, newSr25519Signer)
	assert.Equal(t, expectedErr, err)
}
