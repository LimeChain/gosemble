package types

import (
	"bytes"
	"encoding/hex"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

type testPublicKeyType = Ed25519PublicKey

var (
	expectedExtrinsicSignatureBytes, _ = hex.DecodeString(
		"000101010101010101010101010101010101010101010101010101010101010101000202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020200030000000005000000",
	)
)

var (
	signerAddressBytes = []byte{
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	}
	ed25519SignerOnesAddress, _ = NewEd25519PublicKey(sc.BytesToSequenceU8(signerAddressBytes)...)
	signer                      = NewMultiAddressId(
		NewAccountId[PublicKey](ed25519SignerOnesAddress),
	)

	signatureBytes = []byte{
		2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2,
		2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2,
		2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2,
		2, 2, 2, 2, 2, 2, 2,
	}
	signature = NewMultiSignatureEd25519(
		NewSignatureEd25519(
			sc.BytesToFixedSequenceU8(signatureBytes)...,
		),
	)

	extra = NewSignedExtra([]SignedExtension{
		newTestExtraCheck(false, sc.U32(3)),
		newTestExtraCheck(false, sc.U32(5)),
	})

	targetExtrinsicSignature = ExtrinsicSignature{
		Signer:    signer,
		Signature: signature,
		Extra:     extra,
	}
)

func Test_ExtrinsicSignature_Encode(t *testing.T) {
	buffer := &bytes.Buffer{}

	err := targetExtrinsicSignature.Encode(buffer)

	assert.NoError(t, err)
	assert.Equal(t, expectedExtrinsicSignatureBytes, buffer.Bytes())
}

func Test_ExtrinsicSignature_Bytes(t *testing.T) {
	assert.Equal(t, expectedExtrinsicSignatureBytes, targetExtrinsicSignature.Bytes())
}

func Test_DecodeExtrinsicSignature(t *testing.T) {
	buffer := bytes.NewBuffer(expectedExtrinsicSignatureBytes)

	signedExtraTemplate := NewSignedExtra(
		[]SignedExtension{
			newTestExtraCheck(false, sc.U32(0)),
			newTestExtraCheck(false, sc.U32(0)),
		},
	)

	result, err := DecodeExtrinsicSignature[testPublicKeyType](signedExtraTemplate, buffer)
	assert.NoError(t, err)

	assert.Equal(t, targetExtrinsicSignature, result)
}
