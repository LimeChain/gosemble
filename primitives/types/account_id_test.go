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

	targetAccountIdEd25519 = AccountId{
		Ed25519Signer: NewEd25519Signer(sc.BytesToSequenceU8(pubKeyEd25519Signer)...),
	}
	targetAccountIdSr25519 = AccountId{
		Sr25519Signer: NewSr25519Signer(sc.BytesToSequenceU8(pubKeySr25519Signer)...),
	}
	targetAccountIdEcdsa = AccountId{EcdsaSigner: NewEcdsaSigner(sc.BytesToFixedSequenceU8(addr33Bytes)...)}
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

	result, err := DecodeAccountId[Ed25519Signer](buffer)
	assert.NoError(t, err)

	assert.Equal(t, targetAccountIdEd25519, result)
}

func Test_DecodeAccountId_Sr25519_Signer(t *testing.T) {
	buffer := bytes.NewBuffer(pubKeySr25519Signer)

	result, err := DecodeAccountId[Sr25519Signer](buffer)
	assert.NoError(t, err)

	assert.Equal(t, targetAccountIdSr25519, result)
}

func Test_DecodeAccountId_Ecdsa_Signer(t *testing.T) {
	buffer := bytes.NewBuffer(pubKeyEcdsaSigner)

	result, err := DecodeAccountId[EcdsaSigner](buffer)
	assert.NoError(t, err)

	assert.Equal(t, targetAccountIdEcdsa, result)
}

func Test_DecodeAccountId_PubKeyTypeNotSupported(t *testing.T) {
	buffer := bytes.NewBuffer(pubKeyEcdsaSigner)

	_, err := DecodeAccountId[sc.U8](buffer)
	assert.Error(t, err)

	assert.Equal(t, errorPubKeyNotSupported, err)
}
