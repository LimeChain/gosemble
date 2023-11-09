package types

import (
	"bytes"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	pubKeyEd25519Signer = []byte{1, 1, 1, 1, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0}
	pubKeySr25519Signer = []byte{1, 1, 0, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0}
	pubKeyEcdsaSigner   = []byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0}

	ed25519Signer, _ = NewEd25519PublicKey(sc.BytesToSequenceU8(pubKeyEd25519Signer)...)
	sr25519Signer, _ = NewSr25519PublicKey(sc.BytesToSequenceU8(pubKeySr25519Signer)...)
	ecdsaSigner, _   = NewEcdsaPublicKey(sc.BytesToFixedSequenceU8(addr33Bytes)...)

	targetAccountIdEd25519 = New[PublicKey](ed25519Signer)
	targetAccountIdSr25519 = New[PublicKey](sr25519Signer)
	targetAccountIdEcdsa   = New[PublicKey](ecdsaSigner)
)

func Test_AccountId_Encode_Ed25519_Signer(t *testing.T) {
	buffer := &bytes.Buffer{}

	targetAccountIdEd25519.Encode(buffer)

	assert.Equal(t, pubKeyEd25519Signer, buffer.Bytes())
}

func Test_AccountId_Encode_Sr25519_Signer(t *testing.T) {
	buffer := &bytes.Buffer{}

	targetAccountIdSr25519.Encode(buffer)

	assert.Equal(t, pubKeySr25519Signer, buffer.Bytes())
}

func Test_AccountId_Encode_Ecdsa_Signer(t *testing.T) {
	buffer := &bytes.Buffer{}

	targetAccountIdEcdsa.Encode(buffer)

	assert.Equal(t, pubKeyEcdsaSigner, buffer.Bytes())
}

func Test_AccountId_Bytes(t *testing.T) {
	assert.Equal(t, pubKeyEd25519Signer, targetAccountIdEd25519.Bytes())
	assert.Equal(t, pubKeySr25519Signer, targetAccountIdSr25519.Bytes())
	assert.Equal(t, pubKeyEcdsaSigner, targetAccountIdEcdsa.Bytes())
}

func Test_DecodeAccountId_Ed25519_Signer(t *testing.T) {
	buffer := bytes.NewBuffer(pubKeyEd25519Signer)

	result, err := DecodeAccountId[testKeyType](buffer)
	assert.NoError(t, err)

	assert.Equal(t, targetAccountIdEd25519, result)
}

func Test_DecodeAccountId_Sr25519_Signer(t *testing.T) {
	buffer := bytes.NewBuffer(pubKeySr25519Signer)

	result, err := DecodeAccountId[Sr25519PublicKey](buffer)
	assert.NoError(t, err)

	assert.Equal(t, targetAccountIdSr25519, result)
}

func Test_DecodeAccountId_Ecdsa_Signer(t *testing.T) {
	buffer := bytes.NewBuffer(pubKeyEcdsaSigner)

	result, err := DecodeAccountId[EcdsaPublicKey](buffer)
	assert.NoError(t, err)

	assert.Equal(t, targetAccountIdEcdsa, result)
}
